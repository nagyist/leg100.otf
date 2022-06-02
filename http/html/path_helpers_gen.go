package html

import (
	"fmt"
	"html/template"
)

func loginPath() string {
	return "/login"
}

func logoutPath() string {
	return "/logout"
}

func getProfilePath() string {
	return "/profile"
}

func listSessionPath() string {
	return "/profile/sessions"
}

func revokeSessionPath() string {
	return "/profile/sessions/revoke"
}

func listTokenPath() string {
	return "/profile/tokens"
}

func deleteTokenPath() string {
	return "/profile/tokens/delete"
}

func newTokenPath() string {
	return "/profile/tokens/new"
}

func createTokenPath() string {
	return "/profile/tokens/create"
}

func listOrganizationPath() string {
	return "/organizations"
}

func newOrganizationPath() string {
	return "/organizations/new"
}

func createOrganizationPath() string {
	return "/organizations/create"
}

func getOrganizationPath(organizationName string) string {
	return fmt.Sprintf("/organizations/%s", organizationName)
}

func getOrganizationOverviewPath(organizationName string) string {
	return fmt.Sprintf("/organizations/%s/overview", organizationName)
}

func editOrganizationPath(organizationName string) string {
	return fmt.Sprintf("/organizations/%s/edit", organizationName)
}

func updateOrganizationPath(organizationName string) string {
	return fmt.Sprintf("/organizations/%s/update", organizationName)
}

func deleteOrganizationPath(organizationName string) string {
	return fmt.Sprintf("/organizations/%s/delete", organizationName)
}

func listWorkspacePath(organizationName string) string {
	return fmt.Sprintf("/organizations/%s/workspaces", organizationName)
}

func newWorkspacePath(organizationName string) string {
	return fmt.Sprintf("/organizations/%s/workspaces/new", organizationName)
}

func createWorkspacePath(organizationName string) string {
	return fmt.Sprintf("/organizations/%s/workspaces/create", organizationName)
}

func getWorkspacePath(organizationName, workspaceName string) string {
	return fmt.Sprintf("/organizations/%s/workspaces/%s", organizationName, workspaceName)
}

func editWorkspacePath(organizationName, workspaceName string) string {
	return fmt.Sprintf("/organizations/%s/workspaces/%s/edit", organizationName, workspaceName)
}

func updateWorkspacePath(organizationName, workspaceName string) string {
	return fmt.Sprintf("/organizations/%s/workspaces/%s/update", organizationName, workspaceName)
}

func deleteWorkspacePath(organizationName, workspaceName string) string {
	return fmt.Sprintf("/organizations/%s/workspaces/%s/delete", organizationName, workspaceName)
}

func lockWorkspacePath(organizationName, workspaceName string) string {
	return fmt.Sprintf("/organizations/%s/workspaces/%s/lock", organizationName, workspaceName)
}

func unlockWorkspacePath(organizationName, workspaceName string) string {
	return fmt.Sprintf("/organizations/%s/workspaces/%s/unlock", organizationName, workspaceName)
}

func listRunPath(organizationName, workspaceName string) string {
	return fmt.Sprintf("/organizations/%s/workspaces/%s/runs", organizationName, workspaceName)
}

func newRunPath(organizationName, workspaceName string) string {
	return fmt.Sprintf("/organizations/%s/workspaces/%s/runs/new", organizationName, workspaceName)
}

func createRunPath(organizationName, workspaceName string) string {
	return fmt.Sprintf("/organizations/%s/workspaces/%s/runs/create", organizationName, workspaceName)
}

func getRunPath(organizationName, workspaceName, runId string) string {
	return fmt.Sprintf("/organizations/%s/workspaces/%s/runs/%s", organizationName, workspaceName, runId)
}

func getPlanPath(organizationName, workspaceName, runId string) string {
	return fmt.Sprintf("/organizations/%s/workspaces/%s/runs/%s/plan", organizationName, workspaceName, runId)
}

func getApplyPath(organizationName, workspaceName, runId string) string {
	return fmt.Sprintf("/organizations/%s/workspaces/%s/runs/%s/apply", organizationName, workspaceName, runId)
}

func deleteRunPath(organizationName, workspaceName, runId string) string {
	return fmt.Sprintf("/organizations/%s/workspaces/%s/runs/%s/delete", organizationName, workspaceName, runId)
}

func addHelpersToFuncMap(m template.FuncMap) {
	m["loginPath"] = loginPath
	m["logoutPath"] = logoutPath
	m["getProfilePath"] = getProfilePath
	m["listSessionPath"] = listSessionPath
	m["revokeSessionPath"] = revokeSessionPath
	m["listTokenPath"] = listTokenPath
	m["deleteTokenPath"] = deleteTokenPath
	m["newTokenPath"] = newTokenPath
	m["createTokenPath"] = createTokenPath
	m["listOrganizationPath"] = listOrganizationPath
	m["newOrganizationPath"] = newOrganizationPath
	m["createOrganizationPath"] = createOrganizationPath
	m["getOrganizationPath"] = getOrganizationPath
	m["getOrganizationOverviewPath"] = getOrganizationOverviewPath
	m["editOrganizationPath"] = editOrganizationPath
	m["updateOrganizationPath"] = updateOrganizationPath
	m["deleteOrganizationPath"] = deleteOrganizationPath
	m["listWorkspacePath"] = listWorkspacePath
	m["newWorkspacePath"] = newWorkspacePath
	m["createWorkspacePath"] = createWorkspacePath
	m["getWorkspacePath"] = getWorkspacePath
	m["editWorkspacePath"] = editWorkspacePath
	m["updateWorkspacePath"] = updateWorkspacePath
	m["deleteWorkspacePath"] = deleteWorkspacePath
	m["lockWorkspacePath"] = lockWorkspacePath
	m["unlockWorkspacePath"] = unlockWorkspacePath
	m["listRunPath"] = listRunPath
	m["newRunPath"] = newRunPath
	m["createRunPath"] = createRunPath
	m["getRunPath"] = getRunPath
	m["getPlanPath"] = getPlanPath
	m["getApplyPath"] = getApplyPath
	m["deleteRunPath"] = deleteRunPath
}
