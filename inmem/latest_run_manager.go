package inmem

import (
	"context"
	"fmt"

	"github.com/leg100/otf"
)

var _ otf.LatestRunManager = (*LatestRunManager)(nil)

// LatestRunManager maintains in memory the latest run for each workspace.
type LatestRunManager struct {
	// mapping of workspace ID to ID of latest run - nil means the workspace
	// does not have a latest run
	mapping map[string]*string

	// for subscribing to run events to relay to Watch() consumers
	events otf.EventService
}

func NewLatestRunManager(svc otf.WorkspaceService, events otf.EventService) (*LatestRunManager, error) {
	m := &LatestRunManager{
		events: events,
	}

	// Retrieve latest run for each workspace
	opts := otf.WorkspaceListOptions{}
	for {
		listing, err := svc.List(otf.ContextWithAppUser(), opts)
		if err != nil {
			return nil, fmt.Errorf("retrieving latest runs: %w", err)
		}
		if m.mapping == nil {
			m.mapping = make(map[string]*string, listing.TotalCount())
		}
		for _, ws := range listing.Items {
			m.mapping[ws.ID()] = ws.LatestRunID()
		}
		if listing.NextPage() == nil {
			break
		}
		opts.PageNumber = *listing.NextPage()
	}

	return m, nil
}

func (m *LatestRunManager) Set(ctx context.Context, workspaceID string, run *otf.Run) {
	m.mapping[workspaceID] = otf.String(run.ID())
}

func (m *LatestRunManager) Watch(ctx context.Context, workspaceID string) (<-chan *otf.Run, error) {
	sub, err := m.events.Subscribe("latest")
	if err != nil {
		return nil, err
	}
	c := make(chan *otf.Run, 0)
	go func() {
		for {
			select {
			case <-ctx.Done():
				// context cancelled
				sub.Close()
				return
			case event, ok := <-sub.C():
				if !ok {
					// sender closed channel
					return
				}
				run, ok := event.Payload.(*otf.Run)
				if !ok {
					// skip non-run events
					continue
				}
				if run.WorkspaceID() != workspaceID {
					// skip runs for a different workspace
					continue
				}
				if otf.String(run.ID()) == m.mapping[workspaceID] {
					c <- run
				}
			}
		}
	}()

	return c, nil
}