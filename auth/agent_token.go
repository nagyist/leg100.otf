package auth

import (
	"fmt"
	"time"

	"github.com/leg100/otf"
	"github.com/leg100/otf/rbac"
)

// agentToken is an long-lived authentication token for an external agent.
type agentToken struct {
	id           string
	createdAt    time.Time
	token        string
	description  string
	organization string
}

func (t *agentToken) ID() string           { return t.id }
func (t *agentToken) String() string       { return t.id }
func (t *agentToken) Token() string        { return t.token }
func (t *agentToken) CreatedAt() time.Time { return t.createdAt }
func (t *agentToken) Description() string  { return t.description }
func (t *agentToken) Organization() string { return t.organization }

func (*agentToken) CanAccessSite(action rbac.Action) bool {
	// agent cannot carry out site-level actions
	return false
}

func (t *agentToken) CanAccessOrganization(action rbac.Action, name string) bool {
	return t.organization == name
}

func (t *agentToken) CanAccessWorkspace(action rbac.Action, policy *otf.WorkspacePolicy) bool {
	// agent can access anything within its organization
	return t.organization == policy.Organization
}

func newAgentToken(opts otf.CreateAgentTokenOptions) (*agentToken, error) {
	if opts.Organization == "" {
		return nil, fmt.Errorf("organization name cannot be an empty string")
	}
	if opts.Description == "" {
		return nil, fmt.Errorf("description cannot be an empty string")
	}
	t, err := otf.GenerateAuthToken("agent")
	if err != nil {
		return nil, fmt.Errorf("generating token: %w", err)
	}
	token := agentToken{
		id:           otf.NewID("at"),
		createdAt:    otf.CurrentTimestamp(),
		token:        t,
		description:  opts.Description,
		organization: opts.Organization,
	}
	return &token, nil
}