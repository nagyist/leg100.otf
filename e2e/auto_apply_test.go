package e2e

import (
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
	"github.com/leg100/otf/cloud"
	"github.com/stretchr/testify/require"
)

// TestAutoApply tests auto-apply functionality.
func TestAutoApply(t *testing.T) {
	addBuildsToPath(t)

	wd, err := os.Getwd()
	require.NoError(t, err)
	t.Setenv("SSL_CERT_DIR", path.Join(wd, "./fixtures"))
	t.Logf("SSL_CERT_DIR=%s", os.Getenv("SSL_CERT_DIR"))

	workspace := t.Name() // workspace name reflects test name
	org := uuid.NewString()

	// Build and start a daemon
	user := cloud.User{
		Name:          uuid.NewString(),
		Organizations: []string{org},
		Teams: []cloud.Team{
			{Name: "owners", Organization: org},
		},
	}
	daemon := &daemon{}
	daemon.withGithubUser(&user)
	hostname := daemon.start(t)

	// create browser
	ctx, cancel := chromedp.NewContext(allocator)
	defer cancel()

	// login and create workspace and enable auto-apply
	err = chromedp.Run(ctx, chromedp.Tasks{
		githubLoginTasks(t, hostname, user.Name),
		createWorkspaceTasks(t, hostname, org, workspace),
		chromedp.Tasks{
			// go to workspace
			chromedp.Navigate(workspacePath(hostname, org, workspace)),
			screenshot(t),
			// go to workspace settings
			chromedp.Click(`//a[text()='settings']`, chromedp.NodeVisible),
			screenshot(t),
			// enable auto-apply
			chromedp.Click("input#auto_apply", chromedp.NodeVisible, chromedp.ByQuery),
			screenshot(t),
			// submit form
			chromedp.Click(`//button[text()='Save changes']`, chromedp.NodeVisible),
			screenshot(t),
			// confirm workspace updated
			matchText(t, ".flash-success", "updated workspace"),
		},
		// terraform login in order to fetch token before running CLI below
		terraformLoginTasks(t, hostname),
	})
	require.NoError(t, err)

	// create terraform config
	configPath := newRootModule(t, hostname, org, workspace)

	// terraform init
	cmd := exec.Command("terraform", "init", "-no-color")
	cmd.Dir = configPath
	out, err := cmd.CombinedOutput()
	t.Log(string(out))
	require.NoError(t, err)

	// terraform apply - note we are not passing the -auto-approve flag yet we
	// expect it to auto-apply because the workspace is set to auto-apply.
	cmd = exec.Command("terraform", "apply", "-no-color")
	cmd.Dir = configPath
	out, err = cmd.CombinedOutput()
	t.Log(string(out))
	require.NoError(t, err)
	require.Contains(t, string(out), "Apply complete! Resources: 1 added, 0 changed, 0 destroyed.")
}
