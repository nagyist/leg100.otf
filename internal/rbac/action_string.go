// Code generated by "stringer -type Action ./internal/rbac"; DO NOT EDIT.

package rbac

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[WatchAction-0]
	_ = x[CreateOrganizationAction-1]
	_ = x[UpdateOrganizationAction-2]
	_ = x[GetOrganizationAction-3]
	_ = x[ListOrganizationsAction-4]
	_ = x[GetEntitlementsAction-5]
	_ = x[DeleteOrganizationAction-6]
	_ = x[CreateVCSProviderAction-7]
	_ = x[GetVCSProviderAction-8]
	_ = x[ListVCSProvidersAction-9]
	_ = x[DeleteVCSProviderAction-10]
	_ = x[CreateAgentTokenAction-11]
	_ = x[ListAgentTokensAction-12]
	_ = x[DeleteAgentTokenAction-13]
	_ = x[CreateOrganizationTokenAction-14]
	_ = x[DeleteOrganizationTokenAction-15]
	_ = x[CreateRunTokenAction-16]
	_ = x[CreateModuleAction-17]
	_ = x[CreateModuleVersionAction-18]
	_ = x[UpdateModuleAction-19]
	_ = x[ListModulesAction-20]
	_ = x[GetModuleAction-21]
	_ = x[DeleteModuleAction-22]
	_ = x[DeleteModuleVersionAction-23]
	_ = x[CreateWorkspaceVariableAction-24]
	_ = x[UpdateWorkspaceVariableAction-25]
	_ = x[ListWorkspaceVariablesAction-26]
	_ = x[GetWorkspaceVariableAction-27]
	_ = x[DeleteWorkspaceVariableAction-28]
	_ = x[CreateVariableSetAction-29]
	_ = x[UpdateVariableSetAction-30]
	_ = x[ListVariableSetsAction-31]
	_ = x[GetVariableSetAction-32]
	_ = x[DeleteVariableSetAction-33]
	_ = x[AddVariableToSetAction-34]
	_ = x[RemoveVariableFromSetAction-35]
	_ = x[ApplyVariableSetToWorkspacesAction-36]
	_ = x[DeleteVariableSetFromWorkspacesAction-37]
	_ = x[GetRunAction-38]
	_ = x[ListRunsAction-39]
	_ = x[ApplyRunAction-40]
	_ = x[CreateRunAction-41]
	_ = x[DiscardRunAction-42]
	_ = x[DeleteRunAction-43]
	_ = x[CancelRunAction-44]
	_ = x[EnqueuePlanAction-45]
	_ = x[StartPhaseAction-46]
	_ = x[FinishPhaseAction-47]
	_ = x[PutChunkAction-48]
	_ = x[TailLogsAction-49]
	_ = x[GetPlanFileAction-50]
	_ = x[UploadPlanFileAction-51]
	_ = x[GetLockFileAction-52]
	_ = x[UploadLockFileAction-53]
	_ = x[ListWorkspacesAction-54]
	_ = x[GetWorkspaceAction-55]
	_ = x[CreateWorkspaceAction-56]
	_ = x[DeleteWorkspaceAction-57]
	_ = x[SetWorkspacePermissionAction-58]
	_ = x[UnsetWorkspacePermissionAction-59]
	_ = x[UpdateWorkspaceAction-60]
	_ = x[ListTagsAction-61]
	_ = x[DeleteTagsAction-62]
	_ = x[TagWorkspacesAction-63]
	_ = x[AddTagsAction-64]
	_ = x[RemoveTagsAction-65]
	_ = x[ListWorkspaceTags-66]
	_ = x[LockWorkspaceAction-67]
	_ = x[UnlockWorkspaceAction-68]
	_ = x[ForceUnlockWorkspaceAction-69]
	_ = x[CreateStateVersionAction-70]
	_ = x[ListStateVersionsAction-71]
	_ = x[GetStateVersionAction-72]
	_ = x[DeleteStateVersionAction-73]
	_ = x[RollbackStateVersionAction-74]
	_ = x[DownloadStateAction-75]
	_ = x[GetStateVersionOutputAction-76]
	_ = x[CreateConfigurationVersionAction-77]
	_ = x[ListConfigurationVersionsAction-78]
	_ = x[GetConfigurationVersionAction-79]
	_ = x[DownloadConfigurationVersionAction-80]
	_ = x[DeleteConfigurationVersionAction-81]
	_ = x[CreateUserAction-82]
	_ = x[ListUsersAction-83]
	_ = x[GetUserAction-84]
	_ = x[DeleteUserAction-85]
	_ = x[CreateTeamAction-86]
	_ = x[UpdateTeamAction-87]
	_ = x[GetTeamAction-88]
	_ = x[ListTeamsAction-89]
	_ = x[DeleteTeamAction-90]
	_ = x[AddTeamMembershipAction-91]
	_ = x[RemoveTeamMembershipAction-92]
	_ = x[CreateNotificationConfigurationAction-93]
	_ = x[UpdateNotificationConfigurationAction-94]
	_ = x[ListNotificationConfigurationsAction-95]
	_ = x[GetNotificationConfigurationAction-96]
	_ = x[DeleteNotificationConfigurationAction-97]
}

const _Action_name = "WatchActionCreateOrganizationActionUpdateOrganizationActionGetOrganizationActionListOrganizationsActionGetEntitlementsActionDeleteOrganizationActionCreateVCSProviderActionGetVCSProviderActionListVCSProvidersActionDeleteVCSProviderActionCreateAgentTokenActionListAgentTokensActionDeleteAgentTokenActionCreateOrganizationTokenActionDeleteOrganizationTokenActionCreateRunTokenActionCreateModuleActionCreateModuleVersionActionUpdateModuleActionListModulesActionGetModuleActionDeleteModuleActionDeleteModuleVersionActionCreateWorkspaceVariableActionUpdateWorkspaceVariableActionListWorkspaceVariablesActionGetWorkspaceVariableActionDeleteWorkspaceVariableActionCreateVariableSetActionUpdateVariableSetActionListVariableSetsActionGetVariableSetActionDeleteVariableSetActionAddVariableToSetActionRemoveVariableFromSetActionApplyVariableSetToWorkspacesActionDeleteVariableSetFromWorkspacesActionGetRunActionListRunsActionApplyRunActionCreateRunActionDiscardRunActionDeleteRunActionCancelRunActionEnqueuePlanActionStartPhaseActionFinishPhaseActionPutChunkActionTailLogsActionGetPlanFileActionUploadPlanFileActionGetLockFileActionUploadLockFileActionListWorkspacesActionGetWorkspaceActionCreateWorkspaceActionDeleteWorkspaceActionSetWorkspacePermissionActionUnsetWorkspacePermissionActionUpdateWorkspaceActionListTagsActionDeleteTagsActionTagWorkspacesActionAddTagsActionRemoveTagsActionListWorkspaceTagsLockWorkspaceActionUnlockWorkspaceActionForceUnlockWorkspaceActionCreateStateVersionActionListStateVersionsActionGetStateVersionActionDeleteStateVersionActionRollbackStateVersionActionDownloadStateActionGetStateVersionOutputActionCreateConfigurationVersionActionListConfigurationVersionsActionGetConfigurationVersionActionDownloadConfigurationVersionActionDeleteConfigurationVersionActionCreateUserActionListUsersActionGetUserActionDeleteUserActionCreateTeamActionUpdateTeamActionGetTeamActionListTeamsActionDeleteTeamActionAddTeamMembershipActionRemoveTeamMembershipActionCreateNotificationConfigurationActionUpdateNotificationConfigurationActionListNotificationConfigurationsActionGetNotificationConfigurationActionDeleteNotificationConfigurationAction"

var _Action_index = [...]uint16{0, 11, 35, 59, 80, 103, 124, 148, 171, 191, 213, 236, 258, 279, 301, 330, 359, 379, 397, 422, 440, 457, 472, 490, 515, 544, 573, 601, 627, 656, 679, 702, 724, 744, 767, 789, 816, 850, 887, 899, 913, 927, 942, 958, 973, 988, 1005, 1021, 1038, 1052, 1066, 1083, 1103, 1120, 1140, 1160, 1178, 1199, 1220, 1248, 1278, 1299, 1313, 1329, 1348, 1361, 1377, 1394, 1413, 1434, 1460, 1484, 1507, 1528, 1552, 1578, 1597, 1624, 1656, 1687, 1716, 1750, 1782, 1798, 1813, 1826, 1842, 1858, 1874, 1887, 1902, 1918, 1941, 1967, 2004, 2041, 2077, 2111, 2148}

func (i Action) String() string {
	if i < 0 || i >= Action(len(_Action_index)-1) {
		return "Action(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Action_name[_Action_index[i]:_Action_index[i+1]]
}
