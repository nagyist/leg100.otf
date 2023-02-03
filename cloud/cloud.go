// Package cloud provides types for use with cloud providers.
package cloud

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/leg100/otf/events"
)

// Cloud is an external provider of various cloud services e.g. identity provider, VCS
// repositories etc.
type Cloud interface {
	NewClient(context.Context, ClientOptions) (Client, error)
	HandlerFactory
}

type Service interface {
	GetCloudConfig(name string) (Config, error)
	ListCloudConfigs() []Config
}

// HandlerFactory makes VCS event handlers
type HandlerFactory interface {
	NewHandler(HandlerOptions) http.Handler
}

type HandlerOptions struct {
	Secret    string
	WebhookID uuid.UUID

	events.PubSubService
}

// Repo is a VCS repository belonging to a cloud
//
type Repo struct {
	Identifier string `schema:"identifier,required"` // <repo_owner>/<repo_name>
	Branch     string `schema:"branch,required"`     // default branch
}

func (r Repo) ID() string { return r.Identifier }
