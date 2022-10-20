package e2e

import (
	"os/exec"
	"testing"

	"github.com/leg100/otf"
	"github.com/stretchr/testify/require"
)

// TestE2E demonstrates a user with write permissions on a workspace interacting
// with the workspace via the terraform CLI.
func TestWritePermission(t *testing.T) {
	addBuildsToPath(t)

	// First we need to setup an organization with a user who is both in the
	// owners team and the devops team.
	org := otf.NewTestOrganization(t)
	owners := otf.NewTeam("owners", org)
	devops := otf.NewTeam("devops", org)
	boss := otf.NewUser("boss", otf.WithOrganizationMemberships(org), otf.WithTeamMemberships(owners, devops))
	hostname := startDaemon(t, boss)
	url := "https://" + hostname

	// create workspace via web - note this also syncs the org and owner
	allocater := newBrowserAllocater(t)
	workspace := createWebWorkspace(t, allocater, url, org)

	// assign write permissions to devops team
	addWorkspacePermission(t, allocater, url, org.Name(), workspace, devops.Name(), "write")

	// setup non-owner user - note we start another daemon because this is the
	// only way at present that an additional user can be seeded for testing.
	engineer := otf.NewUser("engineer", otf.WithOrganizationMemberships(org), otf.WithTeamMemberships(devops))
	hostname = startDaemon(t, engineer)

	engineerToken := createAPIToken(t, hostname)
	login(t, hostname, engineerToken)

	root := newRootModule(t, hostname, org.Name(), workspace)

	// terraform init
	cmd := exec.Command("terraform", "init", "-no-color")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	t.Log(string(out))
	require.NoError(t, err)

	// terraform plan
	cmd = exec.Command("terraform", "plan", "-no-color")
	cmd.Dir = root
	out, err = cmd.CombinedOutput()
	t.Log(string(out))
	require.NoError(t, err)
	require.Contains(t, string(out), "Plan: 1 to add, 0 to change, 0 to destroy.")

	// terraform apply
	cmd = exec.Command("terraform", "apply", "-no-color", "-auto-approve")
	cmd.Dir = root
	out, err = cmd.CombinedOutput()
	t.Log(string(out))
	require.NoError(t, err)
	require.Contains(t, string(out), "Apply complete! Resources: 1 added, 0 changed, 0 destroyed.")

	// terraform destroy
	cmd = exec.Command("terraform", "destroy", "-no-color", "-auto-approve")
	cmd.Dir = root
	out, err = cmd.CombinedOutput()
	require.NoError(t, err)
	t.Log(string(out))
	require.Contains(t, string(out), "Apply complete! Resources: 0 added, 0 changed, 1 destroyed.")

	// lock workspace
	cmd = exec.Command("otf", "workspaces", "lock", workspace, "--organization", org.Name())
	cmd.Dir = root
	out, err = cmd.CombinedOutput()
	t.Log(string(out))
	require.NoError(t, err)

	// unlock workspace
	cmd = exec.Command("otf", "workspaces", "unlock", workspace, "--organization", org.Name())
	cmd.Dir = root
	out, err = cmd.CombinedOutput()
	t.Log(string(out))
	require.NoError(t, err)

	// list workspaces
	cmd = exec.Command("otf", "workspaces", "list", "--organization", org.Name())
	cmd.Dir = root
	out, err = cmd.CombinedOutput()
	t.Log(string(out))
	require.NoError(t, err)
	require.Contains(t, string(out), workspace)
}
