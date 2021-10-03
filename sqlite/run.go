package sqlite

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/jmoiron/sqlx"
	"github.com/leg100/otf"
	"gorm.io/gorm"
)

var _ otf.RunStore = (*RunDB)(nil)

type RunDB struct {
	*sqlx.DB
}

func NewRunDB(db *sqlx.DB) *RunDB {
	return &RunDB{
		DB: db,
	}
}

// Create persists a Run to the DB.
func (db RunDB) Create(run *otf.Run) (*otf.Run, error) {
	run.CreatedAt = time.Now()
	run.UpdatedAt = time.Now()

	clauses := map[string]interface{}{
		"configuration_version_id":  run.ConfigurationVersion.ID,
		"created_at":                run.CreatedAt,
		"external_id":               run.ID,
		"force_cancel_available_at": run.ForceCancelAvailableAt,
		"is_destroy":                run.IsDestroy,
		"message":                   run.Message,
		"refresh":                   run.Refresh,
		"refresh_only":              run.RefreshOnly,
		"replace_addrs":             run.ReplaceAddrs,
		"status":                    run.Status,
		"target_addrs":              run.TargetAddrs,
		"updated_at":                run.UpdatedAt,
		"workspace_id":              run.Workspace.ID,
	}

	result, err := sq.Insert("runs").SetMap(clauses).RunWith(db).Exec()
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	timestamps := sq.Insert("run_timestamps").Columns("id", "status", "timestamp")
	for status, timestamp := range run.StatusTimestamps {
		timestamps = timestamps.Values(id, status, timestamp)
	}
	_, err = timestamps.RunWith(db).Exec()
	if err != nil {
		return nil, err
	}

	// plan#create(run_id)

	// apply#create(run_id)

	return run, nil
}

// Update persists an updated Run to the DB. The existing run is fetched from
// the DB, the supplied func is invoked on the run, and the updated run is
// persisted back to the DB. The returned Run includes any changes, including a
// new UpdatedAt value.
func (db RunDB) Update(id string, fn func(*otf.Run) error) (*otf.Run, error) {
	tx := db.MustBegin()
	defer tx.Rollback()

	run, err := getRun(tx, otf.RunGetOptions{ID: &id})
	if err != nil {
		return nil, err
	}

	before := otf.Run{}
	copier.Copy(&before, run)

	// Update obj using client-supplied fn
	if err := fn(run); err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	setIfChanged(before.Status, run.Status, updates, "status")

	if len(updates) == 0 {
		return run, nil
	}

	updates["updated_at"] = time.Now()

	_, err = sq.Update("runs").SetMap(updates).Where("id = ?", run.ID).RunWith(db).Exec()
	if err != nil {
		return nil, err
	}

	// Insert new status timestamps.
	for status, timestamp := range run.StatusTimestamps {
		if _, ok := before.StatusTimestamps[status]; ok {
			// Skip: status timestamp is not new
			continue
		}
		_, err := sq.Insert("run_timestamps").Columns("id", "status", "timestamp").
			Values(run.ID, status, timestamp).RunWith(db).Exec()
		if err != nil {
			return nil, err
		}
	}

	return run, tx.Commit()
}

func (db RunDB) List(opts otf.RunListOptions) (*otf.RunList, error) {
	query := sq.Select("*").From("runs")

	// Optionally filter by workspace
	if opts.WorkspaceID != nil {
		ws, err := getWorkspace(db.DB, otf.WorkspaceSpecifier{ID: opts.WorkspaceID})
		if err != nil {
			return nil, err
		}

		query = query.Where("workspace_id = ?", ws.Model.ID)
	}

	// Optionally filter by statuses
	if len(opts.Statuses) > 0 {
		query = query.Where(sq.Eq{"status": opts.Statuses})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rl := otf.RunList{}
	var count int

	rows, err := db.Queryx(sql, args)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		run, err := scanRun(db.DB, rows)
		if err != nil {
			return nil, err
		}

		rl.Items = append(rl.Items, run)
		count++
	}

	rl.Pagination = otf.NewPagination(opts.ListOptions, count)

	return &rl, nil
}

// Get retrieves a Run domain obj
func (db RunDB) Get(opts otf.RunGetOptions) (*otf.Run, error) {
	return getRun(db.DB, opts)
}

func getRun(db sqlx.Queryer, opts otf.RunGetOptions) (*otf.Run, error) {
	query := sq.Select("runs.*").From("runs")

	switch {
	case opts.ID != nil:
		query = query.Where("external_id = ?", opts.ID)
	case opts.PlanID != nil:
		query = query.Join("plans ON plans.run_id = runs.id").Where("plans.external_id = ?", opts.PlanID)
	case opts.ApplyID != nil:
		query = query.Join("applies ON applies.run_id = runs.id").Where("applies.external_id = ?", opts.ApplyID)
	default:
		return nil, otf.ErrInvalidRunGetOptions
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	run, err := scanRun(db, db.QueryRowx(sql, args))
	if err != nil {
		return nil, err
	}

	// Attach workspace to run
	run.Workspace, err = getWorkspace(db, otf.WorkspaceSpecifier{InternalID: &run.Workspace.Model.ID})
	if err != nil {
		return nil, err
	}

	return run, nil
}

func scanRun(db sqlx.Queryer, scannable StructScannable) (*otf.Run, error) {
	type plan struct {
		metadata

		otf.Resources
		Status         otf.PlanStatus
		LogsBlobID     string `db:"logs_blob_id"`
		PlanFileBlobID string `db:"plan_file_blob_id"`
		PlanJSONBlobID string `db:"plan_json_blob_id"`
		RunID          uint   `db:"run_id"`
	}

	type result struct {
		metadata

		ForceCancelAvailableAt time.Time `db:"force_cancel_available_at"`
		IsDestroy              bool      `db:"is_destroy"`
		Message                string
		PositionInQueue        int `db:"position_in_queue"`
		Refresh                bool
		RefreshOnly            bool `db:"refresh_only"`
		Status                 otf.RunStatus
		ReplaceAddrs           CSV  `db:"replace_addrs"`
		TargetAddrs            CSV  `db:"target_addrs"`
		WorkspaceID            uint `db:"workspace_id"`
		ConfigurationVersionID uint `db:"configuration_version_id"`

		plan
	}

	res := result{}
	if err := scannable.StructScan(res); err != nil {
		return nil, err
	}

	run := otf.Run{
		ID: res.ExternalID,
		Model: gorm.Model{
			ID:        res.ID,
			CreatedAt: res.CreatedAt,
			UpdatedAt: res.UpdatedAt,
		},
		ForceCancelAvailableAt: res.ForceCancelAvailableAt,
		IsDestroy:              res.IsDestroy,
		Message:                res.Message,
		PositionInQueue:        res.PositionInQueue,
		Refresh:                res.Refresh,
		RefreshOnly:            res.RefreshOnly,
		ReplaceAddrs:           res.ReplaceAddrs,
		TargetAddrs:            res.TargetAddrs,
		Status:                 res.Status,
		StatusTimestamps:       make(map[otf.RunStatus]time.Time),
		Workspace: &otf.Workspace{
			Model: gorm.Model{
				ID: res.WorkspaceID,
			},
		},
		ConfigurationVersion: &otf.ConfigurationVersion{
			Model: gorm.Model{
				ID: res.ConfigurationVersionID,
			},
		},
		Plan: &otf.Plan{
			ID: res.ExternalID,
			Model: gorm.Model{
				ID:        res.plan.ID,
				CreatedAt: res.plan.CreatedAt,
				UpdatedAt: res.plan.UpdatedAt,
			},
			Status: res.plan.Status,
		},
	}

	// Fetch and attach timestamps to run
	sql, args, err := sq.Select("status, timestamp").From("run_timestamps").Where("id", res.ID).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := db.Queryx(sql, args)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var timestamp time.Time
		var status otf.RunStatus
		if err := rows.Scan(&status, &timestamp); err != nil {
			return nil, err
		}
		run.StatusTimestamps[status] = timestamp
	}

	// Relations
	Plan * Plan
	Apply * Apply
	Workspace * Workspace
	ConfigurationVersion * ConfigurationVersion
}
