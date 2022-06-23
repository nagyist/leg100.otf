package app

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/leg100/otf"
	"github.com/leg100/otf/inmem"
	"github.com/leg100/otf/sql"
)

var _ otf.ApplyService = (*ApplyService)(nil)

type ApplyService struct {
	proxy otf.ChunkStore
	db    *sql.DB
	otf.EventService
	logr.Logger
}

func NewApplyService(db *sql.DB, logger logr.Logger, es otf.EventService, cache otf.Cache) (*ApplyService, error) {
	proxy, err := inmem.NewChunkProxy(cache, db.ApplyLogStore())
	if err != nil {
		return nil, fmt.Errorf("constructing chunk proxy: %w", err)
	}
	return &ApplyService{
		proxy:        proxy,
		db:           db,
		EventService: es,
		Logger:       logger,
	}, nil
}

func (s ApplyService) Get(ctx context.Context, id string) (*otf.Apply, error) {
	run, err := s.db.GetRun(ctx, otf.RunGetOptions{ApplyID: &id})
	if err != nil {
		return nil, err
	}
	return run.Apply(), nil
}
