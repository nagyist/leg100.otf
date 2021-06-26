package ots

import (
	"fmt"
	"time"

	"github.com/google/jsonapi"
	tfe "github.com/leg100/go-tfe"
)

const (
	MaxPlanLogsLimit = 65536
)

type PlanService interface {
	GetPlan(id string) (*Plan, error)
	UpdatePlanStatus(id string, status tfe.PlanStatus) (*Plan, error)
	FinishPlan(id string, opts *PlanFinishOptions) (*Plan, error)

	GetPlanLogs(id string, opts PlanLogOptions) ([]byte, error)
	UploadPlanLogs(id string, logs []byte) error
}

// Plan represents a Terraform Enterprise plan.
type Plan struct {
	ID                   string                `jsonapi:"primary,plans"`
	HasChanges           bool                  `jsonapi:"attr,has-changes"`
	LogReadURL           string                `jsonapi:"attr,log-read-url"`
	ResourceAdditions    int                   `jsonapi:"attr,resource-additions"`
	ResourceChanges      int                   `jsonapi:"attr,resource-changes"`
	ResourceDestructions int                   `jsonapi:"attr,resource-destructions"`
	Status               tfe.PlanStatus        `jsonapi:"attr,status"`
	StatusTimestamps     *PlanStatusTimestamps `jsonapi:"attr,status-timestamps"`

	// Relations
	//Exports []*PlanExport `jsonapi:"relation,exports"`
}

// PlanLogOptions are send by the client configuring which and how many logs to
// be send to it.
type PlanLogOptions struct {
	// The maximum number of bytes of logs to return to the client
	Limit int `schema:"limit"`

	// The start position in the logs from which to send to the client
	Offset int `schema:"offset"`
}

// PlanFinishOptions represents the options for finishing a plan.
type PlanFinishOptions struct {
	// Type is a public field utilized by JSON:API to set the resource type via
	// the field tag.  It is not a user-defined value and does not need to be
	// set.  https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,plans"`

	ResourceAdditions    int `jsonapi:"attr,resource-additions"`
	ResourceChanges      int `jsonapi:"attr,resource-changes"`
	ResourceDestructions int `jsonapi:"attr,resource-destructions"`
}

func (p *Plan) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("/api/v2/plans/%s", p.ID),
	}
}

// PlanStatusTimestamps holds the timestamps for individual plan statuses.
type PlanStatusTimestamps struct {
	CanceledAt      time.Time `json:"canceled-at,rfc3339"`
	ErroredAt       time.Time `json:"errored-at,rfc3339"`
	FinishedAt      time.Time `json:"finished-at,rfc3339"`
	ForceCanceledAt time.Time `json:"force-canceled-at,rfc3339"`
	QueuedAt        time.Time `json:"queued-at,rfc3339"`
	StartedAt       time.Time `json:"started-at,rfc3339"`
}

func NewPlanID() string {
	return fmt.Sprintf("plan-%s", GenerateRandomString(16))
}
