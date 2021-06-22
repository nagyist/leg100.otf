package ots

import (
	"context"
	"fmt"
	"net/url"
	"time"

	tfe "github.com/hashicorp/go-tfe"
)

type RunService interface {
	CreateRun(opts *tfe.RunCreateOptions) (*Run, error)
	ApplyRun(id string, opts *tfe.RunApplyOptions) error
	GetRun(id string) (*Run, error)
	ListRuns(workspaceID string, opts ListOptions) (*RunList, error)
	DiscardRun(id string, opts *tfe.RunDiscardOptions) error
	CancelRun(id string, opts *tfe.RunCancelOptions) error
	ForceCancelRun(id string, opts *tfe.RunForceCancelOptions) error
}

// RunList represents a list of runs.
type RunList struct {
	*Pagination
	Items []*Run
}

// Run represents a Terraform Enterprise run.
type Run struct {
	ID                     string               `jsonapi:"primary,runs"`
	Actions                *RunActions          `jsonapi:"attr,actions"`
	CreatedAt              time.Time            `jsonapi:"attr,created-at,iso8601"`
	ForceCancelAvailableAt time.Time            `jsonapi:"attr,force-cancel-available-at,iso8601"`
	HasChanges             bool                 `jsonapi:"attr,has-changes"`
	IsDestroy              bool                 `jsonapi:"attr,is-destroy"`
	Message                string               `jsonapi:"attr,message"`
	Permissions            *RunPermissions      `jsonapi:"attr,permissions"`
	PositionInQueue        int                  `jsonapi:"attr,position-in-queue"`
	Refresh                bool                 `jsonapi:"attr,refresh"`
	RefreshOnly            bool                 `jsonapi:"attr,refresh-only"`
	ReplaceAddrs           []string             `jsonapi:"attr,replace-addrs,omitempty"`
	Source                 tfe.RunSource        `jsonapi:"attr,source"`
	Status                 tfe.RunStatus        `jsonapi:"attr,status"`
	StatusTimestamps       *RunStatusTimestamps `jsonapi:"attr,status-timestamps"`
	TargetAddrs            []string             `jsonapi:"attr,target-addrs,omitempty"`

	// Relations
	Apply                *Apply                `jsonapi:"relation,apply"`
	ConfigurationVersion *ConfigurationVersion `jsonapi:"relation,configuration-version"`
	CostEstimate         *CostEstimate         `jsonapi:"relation,cost-estimate"`
	CreatedBy            *User                 `jsonapi:"relation,created-by"`
	Plan                 *Plan                 `jsonapi:"relation,plan"`
	PolicyChecks         []*PolicyCheck        `jsonapi:"relation,policy-checks"`
	Workspace            *Workspace            `jsonapi:"relation,workspace"`
}

// RunActions represents the run actions.
type RunActions struct {
	IsCancelable      bool `jsonapi:"attr,is-cancelable"`
	IsConfirmable     bool `jsonapi:"attr,is-confirmable"`
	IsDiscardable     bool `jsonapi:"attr,is-discardable"`
	IsForceCancelable bool `jsonapi:"attr,is-force-cancelable"`
}

// RunPermissions represents the run permissions.
type RunPermissions struct {
	CanApply        bool `jsonapi:"attr,can-apply"`
	CanCancel       bool `jsonapi:"attr,can-cancel"`
	CanDiscard      bool `jsonapi:"attr,can-discard"`
	CanForceCancel  bool `jsonapi:"attr,can-force-cancel"`
	CanForceExecute bool `jsonapi:"attr,can-force-execute"`
}

// RunStatusTimestamps holds the timestamps for individual run statuses.
type RunStatusTimestamps struct {
	ErroredAt            time.Time `jsonapi:"attr,errored-at,rfc3339"`
	FinishedAt           time.Time `jsonapi:"attr,finished-at,rfc3339"`
	QueuedAt             time.Time `jsonapi:"attr,queued-at,rfc3339"`
	StartedAt            time.Time `jsonapi:"attr,started-at,rfc3339"`
	ApplyingAt           time.Time `jsonapi:"attr,applying-at,rfc3339"`
	AppliedAt            time.Time `jsonapi:"attr,applied-at,rfc3339"`
	PlanningAt           time.Time `jsonapi:"attr,planning-at,rfc3339"`
	PlannedAt            time.Time `jsonapi:"attr,planned-at,rfc3339"`
	PlannedAndFinishedAt time.Time `jsonapi:"attr,planned-and-finished-at,rfc3339"`
	PlanQueuabledAt      time.Time `jsonapi:"attr,plan-queueable-at,rfc3339"`
}

// RunListOptions represents the options for listing runs.
type RunListOptions struct {
	ListOptions

	// A list of relations to include. See available resources:
	// https://www.terraform.io/docs/cloud/api/run.html#available-related-resources
	Include *string `url:"include"`
}

// List all the runs of the given workspace.
func (s *runs) List(ctx context.Context, workspaceID string, options RunListOptions) (*RunList, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s/runs", url.QueryEscape(workspaceID))
	req, err := s.client.newRequest("GET", u, &options)
	if err != nil {
		return nil, err
	}

	rl := &RunList{}
	err = s.client.do(ctx, req, rl)
	if err != nil {
		return nil, err
	}

	return rl, nil
}

// RunCreateOptions represents the options for creating a new run.
type RunCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,runs"`

	// Specifies if this plan is a destroy plan, which will destroy all
	// provisioned resources.
	IsDestroy *bool `jsonapi:"attr,is-destroy,omitempty"`

	// Refresh determines if the run should
	// update the state prior to checking for differences
	Refresh *bool `jsonapi:"attr,refresh,omitempty"`

	// RefreshOnly determines whether the run should ignore config changes
	// and refresh the state only
	RefreshOnly *bool `jsonapi:"attr,refresh-only,omitempty"`

	// Specifies the message to be associated with this run.
	Message *string `jsonapi:"attr,message,omitempty"`

	// Specifies the configuration version to use for this run. If the
	// configuration version object is omitted, the run will be created using the
	// workspace's latest configuration version.
	ConfigurationVersion *ConfigurationVersion `jsonapi:"relation,configuration-version"`

	// Specifies the workspace where the run will be executed.
	Workspace *Workspace `jsonapi:"relation,workspace"`

	// If non-empty, requests that Terraform should create a plan including
	// actions only for the given objects (specified using resource address
	// syntax) and the objects they depend on.
	//
	// This capability is provided for exceptional circumstances only, such as
	// recovering from mistakes or working around existing Terraform
	// limitations. Terraform will generally mention the -target command line
	// option in its error messages describing situations where setting this
	// argument may be appropriate. This argument should not be used as part
	// of routine workflow and Terraform will emit warnings reminding about
	// this whenever this property is set.
	TargetAddrs []string `jsonapi:"attr,target-addrs,omitempty"`

	// If non-empty, requests that Terraform create a plan that replaces
	// (destroys and then re-creates) the objects specified by the given
	// resource addresses.
	ReplaceAddrs []string `jsonapi:"attr,replace-addrs,omitempty"`
}
