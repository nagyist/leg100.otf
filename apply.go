package ots

import (
	"fmt"
	"time"

	"github.com/google/jsonapi"
	tfe "github.com/leg100/go-tfe"
)

type ApplyService interface {
	GetApply(id string) (*Apply, error)
}

// Apply represents a Terraform Enterprise apply.
type Apply struct {
	ID                   string                 `jsonapi:"primary,applies"`
	LogReadURL           string                 `jsonapi:"attr,log-read-url"`
	ResourceAdditions    int                    `jsonapi:"attr,resource-additions"`
	ResourceChanges      int                    `jsonapi:"attr,resource-changes"`
	ResourceDestructions int                    `jsonapi:"attr,resource-destructions"`
	Status               tfe.ApplyStatus        `jsonapi:"attr,status"`
	StatusTimestamps     *ApplyStatusTimestamps `jsonapi:"attr,status-timestamps"`
}

func (a *Apply) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("/api/v2/applies/%s", a.ID),
	}
}

// ApplyStatusTimestamps holds the timestamps for individual apply statuses.
type ApplyStatusTimestamps struct {
	CanceledAt      time.Time `json:"canceled-at,rfc3339"`
	ErroredAt       time.Time `json:"errored-at,rfc3339"`
	FinishedAt      time.Time `json:"finished-at,rfc3339"`
	ForceCanceledAt time.Time `json:"force-canceled-at,rfc3339"`
	QueuedAt        time.Time `json:"queued-at,rfc3339"`
	StartedAt       time.Time `json:"started-at,rfc3339"`
}

func NewApplyID() string {
	return fmt.Sprintf("apply-%s", GenerateRandomString(16))
}
