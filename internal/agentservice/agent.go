package agentservice

import (
	"errors"
	"time"

	"github.com/leg100/otf/internal"
	"golang.org/x/exp/slog"
)

const (
	Busy    Status = "busy"
	Idle    Status = "idle"
	Exited  Status = "exited"
	Errored Status = "errored"
	Unknown Status = "unknown"
)

type (
	Agent struct {
		ID           string
		Status       Status
		IPAddress    string
		Version      string
		Name         *string
		External     bool
		LastSeen     time.Time
		Organization *string
	}

	Status string
)

func newAgent(opts RegisterOptions) (*Agent, error) {
	if opts.External && opts.Organization == nil {
		return nil, errors.New("external agent must specify an organization")
	}
	return &Agent{
		ID:           internal.NewID("agent"),
		Status:       Idle,
		IPAddress:    opts.IPAddress,
		Version:      opts.Version,
		Name:         opts.Name,
		External:     opts.External,
		LastSeen:     internal.CurrentTimestamp(),
		Organization: opts.Organization,
	}, nil
}

func (a *Agent) IsActive() bool {
	switch a.Status {
	case Busy, Idle, Unknown:
		return true
	default:
		return false
	}
}

func (a *Agent) LogValue() slog.Value {
	attrs := []slog.Attr{
		slog.String("id", a.ID),
		slog.String("status", string(a.Status)),
		slog.String("ip_address", a.IPAddress),
	}
	return slog.GroupValue(attrs...)
}
