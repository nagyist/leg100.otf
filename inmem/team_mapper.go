package inmem

import (
	"context"
	"fmt"
	"sync"

	"github.com/leg100/otf"
)

type teamMapper struct {
	mu sync.Mutex
	// map team id to organization name
	idOrgMap map[string]string
	// map qualified team name to team id
	nameIDMap map[otf.TeamQualifiedName]string
}

func newTeamMapper() *teamMapper {
	return &teamMapper{
		idOrgMap:  make(map[string]string),
		nameIDMap: make(map[otf.TeamQualifiedName]string),
	}
}

func (m *teamMapper) populate(ctx context.Context, svc otf.TeamService) error {
	opts := otf.TeamListOptions{}
	var allocated bool
	for {
		listing, err := svc.ListTeams(ctx, opts)
		if err != nil {
			return fmt.Errorf("populating team mapper: %w", err)
		}
		if !allocated {
			m.idOrgMap = make(map[string]string, listing.TotalCount())
			m.nameIDMap = make(map[otf.TeamQualifiedName]string, listing.TotalCount())
			allocated = true
		}
		for _, ws := range listing.Items {
			m.addWithoutLock(ws)
		}
		if listing.NextPage() == nil {
			break
		}
		opts.PageNumber = *listing.NextPage()
	}
	return nil
}

func (m *teamMapper) add(ws *otf.Team) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.addWithoutLock(ws)
}

func (m *teamMapper) addWithoutLock(ws *otf.Team) {
	m.idOrgMap[ws.ID()] = ws.OrganizationName()
	m.nameIDMap[ws.QualifiedName()] = ws.ID()
}

// update the mapping for a team that has been renamed
func (m *teamMapper) update(ws *otf.Team) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// we don't have the old name to hand, so we have to enumerate every entry
	// and look for a team with a matching name.
	for qualified, id := range m.nameIDMap {
		if ws.ID() == id {
			// remove old entry
			delete(m.nameIDMap, qualified)
			// add new entry
			m.nameIDMap[ws.QualifiedName()] = ws.ID()
			return
		}
	}
}

// LookupTeamID looks up the ID of the team given its name and
// organization name.
func (m *teamMapper) lookupID(org, name string) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.nameIDMap[otf.TeamQualifiedName{
		Name:         name,
		Organization: org,
	}]
}

func (m *teamMapper) remove(ws *otf.Team) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.idOrgMap, ws.ID())
	delete(m.nameIDMap, ws.QualifiedName())
}

func (m *teamMapper) lookupOrganizationByID(teamID string) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.idOrgMap[teamID]
}

func (m *teamMapper) lookupOrganizationBySpec(spec otf.TeamSpec) (string, bool) {
	if spec.OrganizationName != nil {
		return *spec.OrganizationName, true
	} else if spec.ID != nil {
		m.mu.Lock()
		defer m.mu.Unlock()

		return m.idOrgMap[*spec.ID], true
	} else {
		return "", false
	}
}
