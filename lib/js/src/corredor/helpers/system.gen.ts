// This is a generated file.
// Run `pnpm run codegen` in lib/js to rebuild it.
//
// A method here mirrors one endpoint of the REST spec, so a change the server
// makes to that endpoint reaches scripts as soon as this file is regenerated.
// Over the raw API client a method adds two things and nothing else: every ID
// in the endpoint's path also accepts the object it names or its handle, and a
// response is cast into its resource class where one exists. The curated helper
// that extends this class is where a breaking endpoint change is absorbed.

import { AxiosRequestConfig } from 'axios'
import { System as SystemAPI } from '../../api-clients'
import { Attachment } from '../../shared'
import {
  Agent,
  Application,
  AuthClient,
  Chatbot,
  ChatbotSession,
  Connection,
  LlmProvider,
  Notification,
  Project,
  ProjectAiSystem,
  Reminder,
  Role,
  Template,
  User,
  UserGroup,
} from '../../system'
import { Args, IdArgs, KV, ListResponse, RefSpec, argsOf, castSet, resolveRef } from './shared'

// Where a handle, a slug or a name arriving in place of an ID is looked up.
const refs: RefSpec = {
  authClient: { idField: 'clientID', list: 'authClientList', by: ['handle'], parents: [] },
  role: { idField: 'roleID', list: 'roleList', by: ['handle', 'name'], parents: [] },
  userGroup: { idField: 'userGroupID', list: 'userGroupList', by: [], parents: [] },
  user: { idField: 'userID', list: 'userList', by: ['handle', 'email'], parents: [] },
  application: { idField: 'applicationID', list: 'applicationList', by: ['name'], parents: [] },
  reminder: { idField: 'reminderID', list: 'reminderList', by: [], parents: [] },
  notification: { idField: 'notificationID', list: 'notificationList', by: [], parents: [] },
  attachment: { idField: 'attachmentID', list: undefined, by: [], parents: [] },
  template: { idField: 'templateID', list: 'templateList', by: ['handle'], parents: [] },
  report: { idField: 'reportID', list: 'reportList', by: ['handle'], parents: [] },
  connection: { idField: 'connectionID', list: 'connectionList', by: ['handle'], parents: [] },
  configuredConnection: {
    idField: 'connectionID',
    list: 'configuredConnectionList',
    by: [],
    parents: [],
  },
  projectIncident: { idField: 'incidentID', list: 'projectIncidentList', by: [], parents: [] },
  projectFeature: { idField: 'featureID', list: 'projectFeatureList', by: [], parents: [] },
  projectPrivacy: { idField: 'privacyID', list: 'projectPrivacyList', by: [], parents: [] },
  projectTask: { idField: 'taskID', list: 'projectTaskList', by: [], parents: [] },
  projectReview: { idField: 'reviewID', list: 'projectReviewList', by: [], parents: [] },
  projectBacklogItem: {
    idField: 'backlogItemID',
    list: 'projectBacklogItemList',
    by: [],
    parents: [],
  },
  agent: { idField: 'agentID', list: 'agentList', by: ['handle'], parents: [] },
  llmProvider: { idField: 'llmProviderID', list: 'llmProviderList', by: [], parents: [] },
  knowledgeBase: {
    idField: 'knowledgeBaseID',
    list: 'knowledgeBaseList',
    by: ['handle'],
    parents: [],
  },
  aiConversation: { idField: 'aiConversationID', list: 'aiConversationList', by: [], parents: [] },
  chatbot: { idField: 'chatbotID', list: 'chatbotList', by: ['handle'], parents: [] },
  chatbotSession: { idField: 'sessionID', list: 'chatbotSessionList', by: [], parents: [] },
  project: { idField: 'projectID', list: 'projectList', by: ['handle'], parents: [] },
  projectAiSystem: {
    idField: 'projectAiSystemID',
    list: 'projectAiSystemList',
    by: ['handle'],
    parents: ['projectID'],
  },
  projectFriaScenario: {
    idField: 'projectFriaScenarioID',
    list: 'projectFriaScenarioList',
    by: [],
    parents: ['projectID'],
  },
  tenant: { idField: 'tenantID', list: 'tenantList', by: ['handle'], parents: [] },
}

export default class GeneratedSystemHelper {
  readonly SystemAPI: SystemAPI

  constructor(ctx: { SystemAPI: SystemAPI }) {
    this.SystemAPI = ctx.SystemAPI
  }

  // Turns an object, an ID or a handle into the ID the endpoint wants.
  protected resolveRef(entity: string, value: unknown, context: KV = {}): Promise<unknown> {
    return resolveRef(this.SystemAPI, refs, entity, value, context)
  }

  // Impersonate a user
  authImpersonate(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.authImpersonate(args, extra)
  }

  // List clients
  // Mirrors SystemAPI.authClientList.
  findAuthClients(
    a: Args = {},
    extra: AxiosRequestConfig = {},
  ): Promise<ListResponse<AuthClient[], KV>> {
    const args = argsOf(a)
    return this.SystemAPI.authClientList(args, extra).then(r => castSet(r, v => new AuthClient(v)))
  }

  // Create client
  // Mirrors SystemAPI.authClientCreate.
  createAuthClient(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<AuthClient> {
    const args = argsOf(a)
    return this.SystemAPI.authClientCreate(args, extra).then(r => new AuthClient(r))
  }

  // Update user details
  // Mirrors SystemAPI.authClientUpdate.
  async updateAuthClient(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<AuthClient> {
    const args = argsOf(a, 'clientID')
    args.clientID = await this.resolveRef('authClient', args.clientID, args)
    return this.SystemAPI.authClientUpdate(args, extra).then(r => new AuthClient(r))
  }

  // Read client details
  // Mirrors SystemAPI.authClientRead.
  async findAuthClientByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<AuthClient> {
    const args = argsOf(a, 'clientID')
    args.clientID = await this.resolveRef('authClient', args.clientID, args)
    return this.SystemAPI.authClientRead(args, extra).then(r => new AuthClient(r))
  }

  // Remove client
  // Mirrors SystemAPI.authClientDelete.
  async deleteAuthClient(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'clientID')
    args.clientID = await this.resolveRef('authClient', args.clientID, args)
    return this.SystemAPI.authClientDelete(args, extra)
  }

  // Undelete client
  // Mirrors SystemAPI.authClientUndelete.
  async undeleteAuthClient(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'clientID')
    args.clientID = await this.resolveRef('authClient', args.clientID, args)
    return this.SystemAPI.authClientUndelete(args, extra)
  }

  // Regenerate client's secret
  async authClientRegenerateSecret(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'clientID')
    args.clientID = await this.resolveRef('authClient', args.clientID, args)
    return this.SystemAPI.authClientRegenerateSecret(args, extra)
  }

  // Exposes client's secret
  async authClientExposeSecret(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'clientID')
    args.clientID = await this.resolveRef('authClient', args.clientID, args)
    return this.SystemAPI.authClientExposeSecret(args, extra)
  }

  // Evaluate expressions
  expressionEvaluate(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.expressionEvaluate(args, extra)
  }

  // List settings
  settingsList(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.settingsList(args, extra)
  }

  // Update settings
  settingsUpdate(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.settingsUpdate(args, extra)
  }

  // Get a value for a key
  settingsGet(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.settingsGet(args, extra)
  }

  // Set value for specific setting
  settingsSet(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.settingsSet(args, extra)
  }

  // Current compose settings
  settingsCurrent(extra: AxiosRequestConfig = {}): Promise<KV> {
    return this.SystemAPI.settingsCurrent(extra)
  }

  // Update role details
  // Mirrors SystemAPI.roleCreate.
  createRole(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<Role> {
    const args = argsOf(a)
    return this.SystemAPI.roleCreate(args, extra).then(r => new Role(r))
  }

  // Update role details
  // Mirrors SystemAPI.roleUpdate.
  async updateRole(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<Role> {
    const args = argsOf(a, 'roleID')
    args.roleID = await this.resolveRef('role', args.roleID, args)
    return this.SystemAPI.roleUpdate(args, extra).then(r => new Role(r))
  }

  // Archive role
  async roleArchive(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'roleID')
    args.roleID = await this.resolveRef('role', args.roleID, args)
    return this.SystemAPI.roleArchive(args, extra)
  }

  // Unarchive role
  async roleUnarchive(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'roleID')
    args.roleID = await this.resolveRef('role', args.roleID, args)
    return this.SystemAPI.roleUnarchive(args, extra)
  }

  // Undelete role
  // Mirrors SystemAPI.roleUndelete.
  async undeleteRole(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'roleID')
    args.roleID = await this.resolveRef('role', args.roleID, args)
    return this.SystemAPI.roleUndelete(args, extra)
  }

  // Move role to different organisation
  async roleMove(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'roleID')
    args.roleID = await this.resolveRef('role', args.roleID, args)
    return this.SystemAPI.roleMove(args, extra)
  }

  // Merge one role into another
  async roleMerge(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'roleID')
    args.roleID = await this.resolveRef('role', args.roleID, args)
    return this.SystemAPI.roleMerge(args, extra)
  }

  // Returns all role members
  async roleMemberList(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'roleID')
    args.roleID = await this.resolveRef('role', args.roleID, args)
    return this.SystemAPI.roleMemberList(args, extra)
  }

  // Add user group to a user group
  async roleMemberAddGroup(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.roleID = await this.resolveRef('role', args.roleID, args)
    args.userGroupID = await this.resolveRef('userGroup', args.userGroupID, args)
    return this.SystemAPI.roleMemberAddGroup(args, extra)
  }

  // Add member to a role
  async roleMemberAdd(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.roleID = await this.resolveRef('role', args.roleID, args)
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.roleMemberAdd(args, extra)
  }

  // Remove member from a role
  async roleMemberRemove(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.roleID = await this.resolveRef('role', args.roleID, args)
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.roleMemberRemove(args, extra)
  }

  // Remove user group from a role
  async roleMemberRemoveGroup(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.roleID = await this.resolveRef('role', args.roleID, args)
    args.userGroupID = await this.resolveRef('userGroup', args.userGroupID, args)
    return this.SystemAPI.roleMemberRemoveGroup(args, extra)
  }

  // Fire system:role trigger
  async roleTriggerScript(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'roleID')
    args.roleID = await this.resolveRef('role', args.roleID, args)
    return this.SystemAPI.roleTriggerScript(args, extra)
  }

  // Clone permission settings to a role
  async roleCloneRules(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'roleID')
    args.roleID = await this.resolveRef('role', args.roleID, args)
    return this.SystemAPI.roleCloneRules(args, extra)
  }

  // List user groups
  // Mirrors SystemAPI.userGroupList.
  findUserGroups(
    a: Args = {},
    extra: AxiosRequestConfig = {},
  ): Promise<ListResponse<UserGroup[], KV>> {
    const args = argsOf(a)
    return this.SystemAPI.userGroupList(args, extra).then(r => castSet(r, v => new UserGroup(v)))
  }

  // Update user groups details
  // Mirrors SystemAPI.userGroupCreate.
  createUserGroup(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<UserGroup> {
    const args = argsOf(a)
    return this.SystemAPI.userGroupCreate(args, extra).then(r => new UserGroup(r))
  }

  // Update user group details
  // Mirrors SystemAPI.userGroupUpdate.
  async updateUserGroup(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<UserGroup> {
    const args = argsOf(a, 'userGroupID')
    args.userGroupID = await this.resolveRef('userGroup', args.userGroupID, args)
    return this.SystemAPI.userGroupUpdate(args, extra).then(r => new UserGroup(r))
  }

  // Read user group details and memberships
  // Mirrors SystemAPI.userGroupRead.
  async findUserGroupByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<UserGroup> {
    const args = argsOf(a, 'userGroupID')
    args.userGroupID = await this.resolveRef('userGroup', args.userGroupID, args)
    return this.SystemAPI.userGroupRead(args, extra).then(r => new UserGroup(r))
  }

  // Remove user group
  // Mirrors SystemAPI.userGroupDelete.
  async deleteUserGroup(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'userGroupID')
    args.userGroupID = await this.resolveRef('userGroup', args.userGroupID, args)
    return this.SystemAPI.userGroupDelete(args, extra)
  }

  // Undelete user group
  // Mirrors SystemAPI.userGroupUndelete.
  async undeleteUserGroup(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'userGroupID')
    args.userGroupID = await this.resolveRef('userGroup', args.userGroupID, args)
    return this.SystemAPI.userGroupUndelete(args, extra)
  }

  // Returns all user group members
  async userGroupMemberList(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'userGroupID')
    args.userGroupID = await this.resolveRef('userGroup', args.userGroupID, args)
    return this.SystemAPI.userGroupMemberList(args, extra)
  }

  // Add member to a user group
  async userGroupMemberAdd(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.userGroupID = await this.resolveRef('userGroup', args.userGroupID, args)
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.userGroupMemberAdd(args, extra)
  }

  // Create user
  // Mirrors SystemAPI.userCreate.
  createUser(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<User> {
    const args = argsOf(a)
    return this.SystemAPI.userCreate(args, extra).then(r => new User(r))
  }

  // Update user details
  // Mirrors SystemAPI.userUpdate.
  async updateUser(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<User> {
    const args = argsOf(a, 'userID')
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.userUpdate(args, extra).then(r => new User(r))
  }

  // Patch user (experimental)
  async userPartialUpdate(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'userID')
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.userPartialUpdate(args, extra)
  }

  // Suspend user
  async userSuspend(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'userID')
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.userSuspend(args, extra)
  }

  // Unsuspend user
  async userUnsuspend(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'userID')
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.userUnsuspend(args, extra)
  }

  // Undelete user
  // Mirrors SystemAPI.userUndelete.
  async undeleteUser(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'userID')
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.userUndelete(args, extra)
  }

  // Set's or changes user's password
  async userSetPassword(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'userID')
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.userSetPassword(args, extra)
  }

  // Add member to a role
  async userMembershipList(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'userID')
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.userMembershipList(args, extra)
  }

  // Add role to a user
  async userMembershipAdd(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.roleID = await this.resolveRef('role', args.roleID, args)
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.userMembershipAdd(args, extra)
  }

  // Remove role from a user
  async userMembershipRemove(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.roleID = await this.resolveRef('role', args.roleID, args)
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.userMembershipRemove(args, extra)
  }

  // Fire system:user trigger
  async userTriggerScript(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'userID')
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.userTriggerScript(args, extra)
  }

  // Remove all auth sessions of user
  async userSessionsRemove(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'userID')
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.userSessionsRemove(args, extra)
  }

  // List user's credentials
  async userListCredentials(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'userID')
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.userListCredentials(args, extra)
  }

  // List user's credentials
  async userDeleteCredentials(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.userDeleteCredentials(args, extra)
  }

  // User's profile avatar
  async userProfileAvatar(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'userID')
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.userProfileAvatar(args, extra)
  }

  // User profile avatar initial
  async userProfileAvatarInitial(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'userID')
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.userProfileAvatarInitial(args, extra)
  }

  // delete user's profile avatar
  async userDeleteAvatar(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'userID')
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.userDeleteAvatar(args, extra)
  }

  // Export users
  userExport(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.userExport(args, extra)
  }

  // Import users
  userImport(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.userImport(args, extra)
  }

  // Search drivers
  dalDriverList(extra: AxiosRequestConfig = {}): Promise<KV> {
    return this.SystemAPI.dalDriverList(extra)
  }

  // Search sensitivity levels
  dalSensitivityLevelList(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.dalSensitivityLevelList(args, extra)
  }

  // Create sensitivity level
  dalSensitivityLevelCreate(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.dalSensitivityLevelCreate(args, extra)
  }

  // Update sensitivity details
  dalSensitivityLevelUpdate(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'sensitivityLevelID')
    return this.SystemAPI.dalSensitivityLevelUpdate(args, extra)
  }

  // Read connection details
  dalSensitivityLevelRead(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'sensitivityLevelID')
    return this.SystemAPI.dalSensitivityLevelRead(args, extra)
  }

  // Remove sensitivity level
  dalSensitivityLevelDelete(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'sensitivityLevelID')
    return this.SystemAPI.dalSensitivityLevelDelete(args, extra)
  }

  // Undelete sensitivity level
  dalSensitivityLevelUndelete(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'sensitivityLevelID')
    return this.SystemAPI.dalSensitivityLevelUndelete(args, extra)
  }

  // Search schema alterations
  dalSchemaAlterationList(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.dalSchemaAlterationList(args, extra)
  }

  // Read alteration details
  dalSchemaAlterationRead(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'alterationID')
    return this.SystemAPI.dalSchemaAlterationRead(args, extra)
  }

  // Apply alterations
  dalSchemaAlterationApply(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.dalSchemaAlterationApply(args, extra)
  }

  // Dismiss alterations
  dalSchemaAlterationDismiss(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.dalSchemaAlterationDismiss(args, extra)
  }

  // Search connections (Directory)
  dalConnectionList(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.dalConnectionList(args, extra)
  }

  // Create connection
  dalConnectionCreate(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.dalConnectionCreate(args, extra)
  }

  // Update connection details
  dalConnectionUpdate(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'connectionID')
    return this.SystemAPI.dalConnectionUpdate(args, extra)
  }

  // Read connection details
  dalConnectionRead(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'connectionID')
    return this.SystemAPI.dalConnectionRead(args, extra)
  }

  // Remove connection
  dalConnectionDelete(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'connectionID')
    return this.SystemAPI.dalConnectionDelete(args, extra)
  }

  // Undelete connection
  dalConnectionUndelete(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'connectionID')
    return this.SystemAPI.dalConnectionUndelete(args, extra)
  }

  // List applications
  // Mirrors SystemAPI.applicationList.
  findApplications(
    a: Args = {},
    extra: AxiosRequestConfig = {},
  ): Promise<ListResponse<Application[], KV>> {
    const args = argsOf(a)
    return this.SystemAPI.applicationList(args, extra).then(r =>
      castSet(r, v => new Application(v)),
    )
  }

  // Create application
  // Mirrors SystemAPI.applicationCreate.
  createApplication(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<Application> {
    const args = argsOf(a)
    return this.SystemAPI.applicationCreate(args, extra).then(r => new Application(r))
  }

  // Update user details
  // Mirrors SystemAPI.applicationUpdate.
  async updateApplication(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<Application> {
    const args = argsOf(a, 'applicationID')
    args.applicationID = await this.resolveRef('application', args.applicationID, args)
    return this.SystemAPI.applicationUpdate(args, extra).then(r => new Application(r))
  }

  // Read the HTML source of a custom application
  async applicationSourceRead(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'applicationID')
    args.applicationID = await this.resolveRef('application', args.applicationID, args)
    return this.SystemAPI.applicationSourceRead(args, extra)
  }

  // Replace the HTML source of a custom application
  async applicationSourceSet(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'applicationID')
    args.applicationID = await this.resolveRef('application', args.applicationID, args)
    return this.SystemAPI.applicationSourceSet(args, extra)
  }

  // Upload application assets
  applicationUpload(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.applicationUpload(args, extra)
  }

  // Flag application
  async applicationFlagCreate(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.applicationID = await this.resolveRef('application', args.applicationID, args)
    return this.SystemAPI.applicationFlagCreate(args, extra)
  }

  // Unflag application
  async applicationFlagDelete(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.applicationID = await this.resolveRef('application', args.applicationID, args)
    return this.SystemAPI.applicationFlagDelete(args, extra)
  }

  // Read application details
  // Mirrors SystemAPI.applicationRead.
  async findApplicationByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<Application> {
    const args = argsOf(a, 'applicationID')
    args.applicationID = await this.resolveRef('application', args.applicationID, args)
    return this.SystemAPI.applicationRead(args, extra).then(r => new Application(r))
  }

  // Remove application
  // Mirrors SystemAPI.applicationDelete.
  async deleteApplication(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'applicationID')
    args.applicationID = await this.resolveRef('application', args.applicationID, args)
    return this.SystemAPI.applicationDelete(args, extra)
  }

  // Undelete application
  // Mirrors SystemAPI.applicationUndelete.
  async undeleteApplication(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'applicationID')
    args.applicationID = await this.resolveRef('application', args.applicationID, args)
    return this.SystemAPI.applicationUndelete(args, extra)
  }

  // Fire system:application trigger
  async applicationTriggerScript(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'applicationID')
    args.applicationID = await this.resolveRef('application', args.applicationID, args)
    return this.SystemAPI.applicationTriggerScript(args, extra)
  }

  // Reorder applications
  applicationReorder(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.applicationReorder(args, extra)
  }

  // List labels
  labelList(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.labelList(args, extra)
  }

  // Delete label
  labelDelete(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.labelDelete(args, extra)
  }

  // Retrieve defined permissions
  permissionsList(extra: AxiosRequestConfig = {}): Promise<KV> {
    return this.SystemAPI.permissionsList(extra)
  }

  // Effective rules for current user
  permissionsEffective(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.permissionsEffective(args, extra)
  }

  // Evaluate rules for given user/role combo
  permissionsTrace(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.permissionsTrace(args, extra)
  }

  // Retrieve role permissions
  async permissionsRead(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'roleID')
    args.roleID = await this.resolveRef('role', args.roleID, args)
    return this.SystemAPI.permissionsRead(args, extra)
  }

  // Remove all defined role permissions
  async permissionsDelete(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'roleID')
    args.roleID = await this.resolveRef('role', args.roleID, args)
    return this.SystemAPI.permissionsDelete(args, extra)
  }

  // Update permission settings
  async permissionsUpdate(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'roleID')
    args.roleID = await this.resolveRef('role', args.roleID, args)
    return this.SystemAPI.permissionsUpdate(args, extra)
  }

  // List/read reminders
  // Mirrors SystemAPI.reminderList.
  findReminders(
    a: Args = {},
    extra: AxiosRequestConfig = {},
  ): Promise<ListResponse<Reminder[], KV>> {
    const args = argsOf(a)
    return this.SystemAPI.reminderList(args, extra).then(r => castSet(r, v => new Reminder(v)))
  }

  // Add new reminder
  // Mirrors SystemAPI.reminderCreate.
  createReminder(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<Reminder> {
    const args = argsOf(a)
    return this.SystemAPI.reminderCreate(args, extra).then(r => new Reminder(r))
  }

  // Update reminder
  // Mirrors SystemAPI.reminderUpdate.
  async updateReminder(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<Reminder> {
    const args = argsOf(a, 'reminderID')
    args.reminderID = await this.resolveRef('reminder', args.reminderID, args)
    return this.SystemAPI.reminderUpdate(args, extra).then(r => new Reminder(r))
  }

  // Read reminder by ID
  // Mirrors SystemAPI.reminderRead.
  async findReminderByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<Reminder> {
    const args = argsOf(a, 'reminderID')
    args.reminderID = await this.resolveRef('reminder', args.reminderID, args)
    return this.SystemAPI.reminderRead(args, extra).then(r => new Reminder(r))
  }

  // Delete reminder
  // Mirrors SystemAPI.reminderDelete.
  async deleteReminder(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'reminderID')
    args.reminderID = await this.resolveRef('reminder', args.reminderID, args)
    return this.SystemAPI.reminderDelete(args, extra)
  }

  // Dismiss reminder
  async reminderDismiss(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'reminderID')
    args.reminderID = await this.resolveRef('reminder', args.reminderID, args)
    return this.SystemAPI.reminderDismiss(args, extra)
  }

  // Undismiss reminder
  async reminderUndismiss(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'reminderID')
    args.reminderID = await this.resolveRef('reminder', args.reminderID, args)
    return this.SystemAPI.reminderUndismiss(args, extra)
  }

  // Snooze reminder
  async reminderSnooze(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'reminderID')
    args.reminderID = await this.resolveRef('reminder', args.reminderID, args)
    return this.SystemAPI.reminderSnooze(args, extra)
  }

  // List/read notifications
  // Mirrors SystemAPI.notificationList.
  findNotifications(
    a: Args = {},
    extra: AxiosRequestConfig = {},
  ): Promise<ListResponse<Notification[], KV>> {
    const args = argsOf(a)
    return this.SystemAPI.notificationList(args, extra).then(r =>
      castSet(r, v => new Notification(v)),
    )
  }

  // Add new notification
  // Mirrors SystemAPI.notificationCreate.
  createNotification(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<Notification> {
    const args = argsOf(a)
    return this.SystemAPI.notificationCreate(args, extra).then(r => new Notification(r))
  }

  // Update notification
  // Mirrors SystemAPI.notificationUpdate.
  async updateNotification(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<Notification> {
    const args = argsOf(a, 'notificationID')
    args.notificationID = await this.resolveRef('notification', args.notificationID, args)
    return this.SystemAPI.notificationUpdate(args, extra).then(r => new Notification(r))
  }

  // Read notification by ID
  // Mirrors SystemAPI.notificationRead.
  async findNotificationByID(
    a: IdArgs = {},
    extra: AxiosRequestConfig = {},
  ): Promise<Notification> {
    const args = argsOf(a, 'notificationID')
    args.notificationID = await this.resolveRef('notification', args.notificationID, args)
    return this.SystemAPI.notificationRead(args, extra).then(r => new Notification(r))
  }

  // Delete notification
  // Mirrors SystemAPI.notificationDelete.
  async deleteNotification(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'notificationID')
    args.notificationID = await this.resolveRef('notification', args.notificationID, args)
    return this.SystemAPI.notificationDelete(args, extra)
  }

  // Mark notification as read
  async notificationMarkAsRead(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'notificationID')
    args.notificationID = await this.resolveRef('notification', args.notificationID, args)
    return this.SystemAPI.notificationMarkAsRead(args, extra)
  }

  // Mark notification as unread
  async notificationMarkAsUnread(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'notificationID')
    args.notificationID = await this.resolveRef('notification', args.notificationID, args)
    return this.SystemAPI.notificationMarkAsUnread(args, extra)
  }

  // Mark all notifications as read for current user
  notificationMarkAllAsRead(extra: AxiosRequestConfig = {}): Promise<KV> {
    return this.SystemAPI.notificationMarkAllAsRead(extra)
  }

  // Mark all notifications as unread for current user
  notificationMarkAllAsUnread(extra: AxiosRequestConfig = {}): Promise<KV> {
    return this.SystemAPI.notificationMarkAllAsUnread(extra)
  }

  // Attachment details
  // Mirrors SystemAPI.attachmentRead.
  async findAttachmentByID(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<Attachment> {
    const args = argsOf(a)
    args.attachmentID = await this.resolveRef('attachment', args.attachmentID, args)
    return this.SystemAPI.attachmentRead(args, extra).then(r => new Attachment(r))
  }

  // Delete attachment
  // Mirrors SystemAPI.attachmentDelete.
  async deleteAttachment(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.attachmentID = await this.resolveRef('attachment', args.attachmentID, args)
    return this.SystemAPI.attachmentDelete(args, extra)
  }

  // Serves attached file
  async attachmentOriginal(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.attachmentID = await this.resolveRef('attachment', args.attachmentID, args)
    return this.SystemAPI.attachmentOriginal(args, extra)
  }

  // Serves preview of an attached file
  async attachmentPreview(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.attachmentID = await this.resolveRef('attachment', args.attachmentID, args)
    return this.SystemAPI.attachmentPreview(args, extra)
  }

  // List templates
  // Mirrors SystemAPI.templateList.
  findTemplates(
    a: Args = {},
    extra: AxiosRequestConfig = {},
  ): Promise<ListResponse<Template[], KV>> {
    const args = argsOf(a)
    return this.SystemAPI.templateList(args, extra).then(r => castSet(r, v => new Template(v)))
  }

  // Create template
  // Mirrors SystemAPI.templateCreate.
  createTemplate(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<Template> {
    const args = argsOf(a)
    return this.SystemAPI.templateCreate(args, extra).then(r => new Template(r))
  }

  // Read template
  // Mirrors SystemAPI.templateRead.
  async findTemplateByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<Template> {
    const args = argsOf(a, 'templateID')
    args.templateID = await this.resolveRef('template', args.templateID, args)
    return this.SystemAPI.templateRead(args, extra).then(r => new Template(r))
  }

  // Update template
  // Mirrors SystemAPI.templateUpdate.
  async updateTemplate(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<Template> {
    const args = argsOf(a, 'templateID')
    args.templateID = await this.resolveRef('template', args.templateID, args)
    return this.SystemAPI.templateUpdate(args, extra).then(r => new Template(r))
  }

  // Delete template
  // Mirrors SystemAPI.templateDelete.
  async deleteTemplate(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'templateID')
    args.templateID = await this.resolveRef('template', args.templateID, args)
    return this.SystemAPI.templateDelete(args, extra)
  }

  // Undelete template
  // Mirrors SystemAPI.templateUndelete.
  async undeleteTemplate(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'templateID')
    args.templateID = await this.resolveRef('template', args.templateID, args)
    return this.SystemAPI.templateUndelete(args, extra)
  }

  // Render drivers
  templateRenderDrivers(extra: AxiosRequestConfig = {}): Promise<KV> {
    return this.SystemAPI.templateRenderDrivers(extra)
  }

  // Render template
  async templateRender(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.templateID = await this.resolveRef('template', args.templateID, args)
    return this.SystemAPI.templateRender(args, extra)
  }

  // List reports
  // Mirrors SystemAPI.reportList.
  findReports(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.reportList(args, extra)
  }

  // Create report
  // Mirrors SystemAPI.reportCreate.
  createReport(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.reportCreate(args, extra)
  }

  // Update report
  // Mirrors SystemAPI.reportUpdate.
  async updateReport(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'reportID')
    args.reportID = await this.resolveRef('report', args.reportID, args)
    return this.SystemAPI.reportUpdate(args, extra)
  }

  // Read report details
  // Mirrors SystemAPI.reportRead.
  async findReportByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'reportID')
    args.reportID = await this.resolveRef('report', args.reportID, args)
    return this.SystemAPI.reportRead(args, extra)
  }

  // Remove report
  // Mirrors SystemAPI.reportDelete.
  async deleteReport(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'reportID')
    args.reportID = await this.resolveRef('report', args.reportID, args)
    return this.SystemAPI.reportDelete(args, extra)
  }

  // Undelete report
  // Mirrors SystemAPI.reportUndelete.
  async undeleteReport(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'reportID')
    args.reportID = await this.resolveRef('report', args.reportID, args)
    return this.SystemAPI.reportUndelete(args, extra)
  }

  // Describe report
  reportDescribe(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.reportDescribe(args, extra)
  }

  // Run report
  async reportRun(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'reportID')
    args.reportID = await this.resolveRef('report', args.reportID, args)
    return this.SystemAPI.reportRun(args, extra)
  }

  // List system statistics
  statsList(extra: AxiosRequestConfig = {}): Promise<KV> {
    return this.SystemAPI.statsList(extra)
  }

  // List all available automation scripts for system resources
  automationList(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.automationList(args, extra)
  }

  // Serves client scripts bundle
  automationBundle(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.automationBundle(args, extra)
  }

  // Triggers execution of a specific script on a system service level
  automationTriggerScript(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.automationTriggerScript(args, extra)
  }

  // Action log events
  actionlogList(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.actionlogList(args, extra)
  }

  // Aggregated action log report
  actionlogReport(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.actionlogReport(args, extra)
  }

  // Messaging queues
  queuesList(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.queuesList(args, extra)
  }

  // Create messaging queue
  queuesCreate(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.queuesCreate(args, extra)
  }

  // Messaging queue details
  queuesRead(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'queueID')
    return this.SystemAPI.queuesRead(args, extra)
  }

  // Update queue details
  queuesUpdate(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'queueID')
    return this.SystemAPI.queuesUpdate(args, extra)
  }

  // Messaging queue delete
  queuesDelete(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'queueID')
    return this.SystemAPI.queuesDelete(args, extra)
  }

  // Messaging queue undelete
  queuesUndelete(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'queueID')
    return this.SystemAPI.queuesUndelete(args, extra)
  }

  // List routes
  apigwRouteList(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.apigwRouteList(args, extra)
  }

  // Create route
  apigwRouteCreate(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.apigwRouteCreate(args, extra)
  }

  // Update route details
  apigwRouteUpdate(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'routeID')
    return this.SystemAPI.apigwRouteUpdate(args, extra)
  }

  // Read route details
  apigwRouteRead(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'routeID')
    return this.SystemAPI.apigwRouteRead(args, extra)
  }

  // Remove route
  apigwRouteDelete(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'routeID')
    return this.SystemAPI.apigwRouteDelete(args, extra)
  }

  // Undelete route
  apigwRouteUndelete(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'routeID')
    return this.SystemAPI.apigwRouteUndelete(args, extra)
  }

  // List filters
  apigwFilterList(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.apigwFilterList(args, extra)
  }

  // Create filter
  apigwFilterCreate(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.apigwFilterCreate(args, extra)
  }

  // Update filter details
  apigwFilterUpdate(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'filterID')
    return this.SystemAPI.apigwFilterUpdate(args, extra)
  }

  // Read filter details
  apigwFilterRead(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'filterID')
    return this.SystemAPI.apigwFilterRead(args, extra)
  }

  // Remove filter
  apigwFilterDelete(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'filterID')
    return this.SystemAPI.apigwFilterDelete(args, extra)
  }

  // Undelete filter
  apigwFilterUndelete(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'filterID')
    return this.SystemAPI.apigwFilterUndelete(args, extra)
  }

  // Filter definitions
  apigwFilterDefFilter(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.apigwFilterDefFilter(args, extra)
  }

  // Proxy auth definitions
  apigwFilterDefProxyAuth(extra: AxiosRequestConfig = {}): Promise<KV> {
    return this.SystemAPI.apigwFilterDefProxyAuth(extra)
  }

  // List aggregated list of routes
  apigwProfilerAggregation(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.apigwProfilerAggregation(args, extra)
  }

  // List hits per route
  apigwProfilerRoute(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'routeID')
    return this.SystemAPI.apigwProfilerRoute(args, extra)
  }

  // Hit details
  apigwProfilerHit(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'hitID')
    return this.SystemAPI.apigwProfilerHit(args, extra)
  }

  // Purge all profiler hits
  apigwProfilerPurgeAll(extra: AxiosRequestConfig = {}): Promise<KV> {
    return this.SystemAPI.apigwProfilerPurgeAll(extra)
  }

  // Purge route profiler hits
  apigwProfilerPurge(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'routeID')
    return this.SystemAPI.apigwProfilerPurge(args, extra)
  }

  // List resources translations
  localeListResource(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.localeListResource(args, extra)
  }

  // Create resource translation
  localeCreateResource(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.localeCreateResource(args, extra)
  }

  // Update resource translation
  localeUpdateResource(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'translationID')
    return this.SystemAPI.localeUpdateResource(args, extra)
  }

  // Read resource translation details
  localeReadResource(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'translationID')
    return this.SystemAPI.localeReadResource(args, extra)
  }

  // Remove resource translation
  localeDeleteResource(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'translationID')
    return this.SystemAPI.localeDeleteResource(args, extra)
  }

  // Undelete resource translation
  localeUndeleteResource(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'translationID')
    return this.SystemAPI.localeUndeleteResource(args, extra)
  }

  // List all available languages
  localeList(extra: AxiosRequestConfig = {}): Promise<KV> {
    return this.SystemAPI.localeList(extra)
  }

  // List all available translation in a language for a specific webapp
  localeGet(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.localeGet(args, extra)
  }

  // List connections for data privacy
  dataPrivacyConnectionList(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.dataPrivacyConnectionList(args, extra)
  }

  // List data privacy requests
  dataPrivacyRequestList(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.dataPrivacyRequestList(args, extra)
  }

  // Create data privacy request
  dataPrivacyRequestCreate(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.dataPrivacyRequestCreate(args, extra)
  }

  // Get details about specific request
  dataPrivacyRequestRead(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'requestID')
    return this.SystemAPI.dataPrivacyRequestRead(args, extra)
  }

  // Update data privacy request status
  dataPrivacyRequestUpdateStatus(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.dataPrivacyRequestUpdateStatus(args, extra)
  }

  // List data privacy request comments
  dataPrivacyRequestCommentList(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'requestID')
    return this.SystemAPI.dataPrivacyRequestCommentList(args, extra)
  }

  // Create data privacy request comment
  dataPrivacyRequestCommentCreate(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'requestID')
    return this.SystemAPI.dataPrivacyRequestCommentCreate(args, extra)
  }

  // List connections
  // Mirrors SystemAPI.connectionList.
  findConnections(
    a: Args = {},
    extra: AxiosRequestConfig = {},
  ): Promise<ListResponse<Connection[], KV>> {
    const args = argsOf(a)
    return this.SystemAPI.connectionList(args, extra).then(r => castSet(r, v => new Connection(v)))
  }

  // Create connection
  // Mirrors SystemAPI.connectionCreate.
  createConnection(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<Connection> {
    const args = argsOf(a)
    return this.SystemAPI.connectionCreate(args, extra).then(r => new Connection(r))
  }

  // Update connection
  // Mirrors SystemAPI.connectionUpdate.
  async updateConnection(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<Connection> {
    const args = argsOf(a, 'connectionID')
    args.connectionID = await this.resolveRef('connection', args.connectionID, args)
    return this.SystemAPI.connectionUpdate(args, extra).then(r => new Connection(r))
  }

  // Read connection
  // Mirrors SystemAPI.connectionRead.
  async findConnectionByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<Connection> {
    const args = argsOf(a, 'connectionID')
    args.connectionID = await this.resolveRef('connection', args.connectionID, args)
    return this.SystemAPI.connectionRead(args, extra).then(r => new Connection(r))
  }

  // Delete connection
  // Mirrors SystemAPI.connectionDelete.
  async deleteConnection(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'connectionID')
    args.connectionID = await this.resolveRef('connection', args.connectionID, args)
    return this.SystemAPI.connectionDelete(args, extra)
  }

  // Undelete connection
  // Mirrors SystemAPI.connectionUndelete.
  async undeleteConnection(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'connectionID')
    args.connectionID = await this.resolveRef('connection', args.connectionID, args)
    return this.SystemAPI.connectionUndelete(args, extra)
  }

  // Enable connection (lock draft local connection to active)
  async connectionEnable(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'connectionID')
    args.connectionID = await this.resolveRef('connection', args.connectionID, args)
    return this.SystemAPI.connectionEnable(args, extra)
  }

  // Generate connection via builder agent
  connectionGenerate(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.connectionGenerate(args, extra)
  }

  // Import a connection from the catalog into the local store
  connectionImport(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.connectionImport(args, extra)
  }

  // Configure a connection (create configured connection)
  async connectionConfigure(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'connectionID')
    args.connectionID = await this.resolveRef('connection', args.connectionID, args)
    return this.SystemAPI.connectionConfigure(args, extra)
  }

  // Update connection configuration
  async connectionUpdateConfiguration(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.connectionID = await this.resolveRef('connection', args.connectionID, args)
    return this.SystemAPI.connectionUpdateConfiguration(args, extra)
  }

  // List configured connections
  // Mirrors SystemAPI.configuredConnectionList.
  findConfiguredConnections(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.configuredConnectionList(args, extra)
  }

  // Read configured connection
  // Mirrors SystemAPI.configuredConnectionRead.
  async findConfiguredConnectionByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'connectionID')
    args.connectionID = await this.resolveRef('configuredConnection', args.connectionID, args)
    return this.SystemAPI.configuredConnectionRead(args, extra)
  }

  // Delete configured connection
  // Mirrors SystemAPI.configuredConnectionDelete.
  async deleteConfiguredConnection(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'connectionID')
    args.connectionID = await this.resolveRef('configuredConnection', args.connectionID, args)
    return this.SystemAPI.configuredConnectionDelete(args, extra)
  }

  // Enable configured connection (validate, provision, lock to active)
  async configuredConnectionEnable(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'connectionID')
    args.connectionID = await this.resolveRef('configuredConnection', args.connectionID, args)
    return this.SystemAPI.configuredConnectionEnable(args, extra)
  }

  // Check configured connection health (connectivity, auth, probe)
  async configuredConnectionCheck(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'connectionID')
    args.connectionID = await this.resolveRef('configuredConnection', args.connectionID, args)
    return this.SystemAPI.configuredConnectionCheck(args, extra)
  }

  // List available MCP tools
  mcpListTools(extra: AxiosRequestConfig = {}): Promise<KV> {
    return this.SystemAPI.mcpListTools(extra)
  }

  // Check SMTP server configuration settings
  smtpConfigurationCheckerCheck(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.smtpConfigurationCheckerCheck(args, extra)
  }

  // List project incidents
  // Mirrors SystemAPI.projectIncidentList.
  findProjectIncidents(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.projectIncidentList(args, extra)
  }

  // Create project incident
  // Mirrors SystemAPI.projectIncidentCreate.
  createProjectIncident(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.projectIncidentCreate(args, extra)
  }

  // Read project incident details
  // Mirrors SystemAPI.projectIncidentRead.
  async findProjectIncidentByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'incidentID')
    args.incidentID = await this.resolveRef('projectIncident', args.incidentID, args)
    return this.SystemAPI.projectIncidentRead(args, extra)
  }

  // Update project incident details
  // Mirrors SystemAPI.projectIncidentUpdate.
  async updateProjectIncident(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'incidentID')
    args.incidentID = await this.resolveRef('projectIncident', args.incidentID, args)
    return this.SystemAPI.projectIncidentUpdate(args, extra)
  }

  // Delete project incident
  // Mirrors SystemAPI.projectIncidentDelete.
  async deleteProjectIncident(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'incidentID')
    args.incidentID = await this.resolveRef('projectIncident', args.incidentID, args)
    return this.SystemAPI.projectIncidentDelete(args, extra)
  }

  // List project features
  // Mirrors SystemAPI.projectFeatureList.
  findProjectFeatures(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.projectFeatureList(args, extra)
  }

  // Create project feature
  // Mirrors SystemAPI.projectFeatureCreate.
  createProjectFeature(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.projectFeatureCreate(args, extra)
  }

  // Read project feature details
  // Mirrors SystemAPI.projectFeatureRead.
  async findProjectFeatureByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'featureID')
    args.featureID = await this.resolveRef('projectFeature', args.featureID, args)
    return this.SystemAPI.projectFeatureRead(args, extra)
  }

  // Update project feature details
  // Mirrors SystemAPI.projectFeatureUpdate.
  async updateProjectFeature(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'featureID')
    args.featureID = await this.resolveRef('projectFeature', args.featureID, args)
    return this.SystemAPI.projectFeatureUpdate(args, extra)
  }

  // Delete project feature
  // Mirrors SystemAPI.projectFeatureDelete.
  async deleteProjectFeature(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'featureID')
    args.featureID = await this.resolveRef('projectFeature', args.featureID, args)
    return this.SystemAPI.projectFeatureDelete(args, extra)
  }

  // List project privacy items
  // Mirrors SystemAPI.projectPrivacyList.
  findProjectPrivacies(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.projectPrivacyList(args, extra)
  }

  // Create project privacy item
  // Mirrors SystemAPI.projectPrivacyCreate.
  createProjectPrivacy(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.projectPrivacyCreate(args, extra)
  }

  // Read project privacy item details
  // Mirrors SystemAPI.projectPrivacyRead.
  async findProjectPrivacyByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'privacyID')
    args.privacyID = await this.resolveRef('projectPrivacy', args.privacyID, args)
    return this.SystemAPI.projectPrivacyRead(args, extra)
  }

  // Update project privacy item details
  // Mirrors SystemAPI.projectPrivacyUpdate.
  async updateProjectPrivacy(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'privacyID')
    args.privacyID = await this.resolveRef('projectPrivacy', args.privacyID, args)
    return this.SystemAPI.projectPrivacyUpdate(args, extra)
  }

  // Delete project privacy item
  // Mirrors SystemAPI.projectPrivacyDelete.
  async deleteProjectPrivacy(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'privacyID')
    args.privacyID = await this.resolveRef('projectPrivacy', args.privacyID, args)
    return this.SystemAPI.projectPrivacyDelete(args, extra)
  }

  // List project tasks
  // Mirrors SystemAPI.projectTaskList.
  findProjectTasks(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.projectTaskList(args, extra)
  }

  // Create project task
  // Mirrors SystemAPI.projectTaskCreate.
  createProjectTask(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.projectTaskCreate(args, extra)
  }

  // Read project task details
  // Mirrors SystemAPI.projectTaskRead.
  async findProjectTaskByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'taskID')
    args.taskID = await this.resolveRef('projectTask', args.taskID, args)
    return this.SystemAPI.projectTaskRead(args, extra)
  }

  // Update project task details
  // Mirrors SystemAPI.projectTaskUpdate.
  async updateProjectTask(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'taskID')
    args.taskID = await this.resolveRef('projectTask', args.taskID, args)
    return this.SystemAPI.projectTaskUpdate(args, extra)
  }

  // Delete project task
  // Mirrors SystemAPI.projectTaskDelete.
  async deleteProjectTask(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'taskID')
    args.taskID = await this.resolveRef('projectTask', args.taskID, args)
    return this.SystemAPI.projectTaskDelete(args, extra)
  }

  // List project reviews
  // Mirrors SystemAPI.projectReviewList.
  findProjectReviews(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.projectReviewList(args, extra)
  }

  // Create project review
  // Mirrors SystemAPI.projectReviewCreate.
  createProjectReview(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.projectReviewCreate(args, extra)
  }

  // Read project review details
  // Mirrors SystemAPI.projectReviewRead.
  async findProjectReviewByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'reviewID')
    args.reviewID = await this.resolveRef('projectReview', args.reviewID, args)
    return this.SystemAPI.projectReviewRead(args, extra)
  }

  // Update project review details
  // Mirrors SystemAPI.projectReviewUpdate.
  async updateProjectReview(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'reviewID')
    args.reviewID = await this.resolveRef('projectReview', args.reviewID, args)
    return this.SystemAPI.projectReviewUpdate(args, extra)
  }

  // Delete project review
  // Mirrors SystemAPI.projectReviewDelete.
  async deleteProjectReview(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'reviewID')
    args.reviewID = await this.resolveRef('projectReview', args.reviewID, args)
    return this.SystemAPI.projectReviewDelete(args, extra)
  }

  // List project backlog items
  // Mirrors SystemAPI.projectBacklogItemList.
  findProjectBacklogItems(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.projectBacklogItemList(args, extra)
  }

  // Create project backlog item
  // Mirrors SystemAPI.projectBacklogItemCreate.
  createProjectBacklogItem(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.projectBacklogItemCreate(args, extra)
  }

  // Read project backlog item details
  // Mirrors SystemAPI.projectBacklogItemRead.
  async findProjectBacklogItemByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'backlogItemID')
    args.backlogItemID = await this.resolveRef('projectBacklogItem', args.backlogItemID, args)
    return this.SystemAPI.projectBacklogItemRead(args, extra)
  }

  // Update project backlog item details
  // Mirrors SystemAPI.projectBacklogItemUpdate.
  async updateProjectBacklogItem(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'backlogItemID')
    args.backlogItemID = await this.resolveRef('projectBacklogItem', args.backlogItemID, args)
    return this.SystemAPI.projectBacklogItemUpdate(args, extra)
  }

  // Delete project backlog item
  // Mirrors SystemAPI.projectBacklogItemDelete.
  async deleteProjectBacklogItem(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'backlogItemID')
    args.backlogItemID = await this.resolveRef('projectBacklogItem', args.backlogItemID, args)
    return this.SystemAPI.projectBacklogItemDelete(args, extra)
  }

  // Aggregated project category report
  projectReportReport(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.projectReportReport(args, extra)
  }

  // Kanban board columns (true totals + first page of cards) across all six work-item types
  projectBoardBoard(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.projectBoardBoard(args, extra)
  }

  // List agents
  // Mirrors SystemAPI.agentList.
  findAgents(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<ListResponse<Agent[], KV>> {
    const args = argsOf(a)
    return this.SystemAPI.agentList(args, extra).then(r => castSet(r, v => new Agent(v)))
  }

  // Create agent
  // Mirrors SystemAPI.agentCreate.
  createAgent(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<Agent> {
    const args = argsOf(a)
    return this.SystemAPI.agentCreate(args, extra).then(r => new Agent(r))
  }

  // Read agent details
  // Mirrors SystemAPI.agentRead.
  async findAgentByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<Agent> {
    const args = argsOf(a, 'agentID')
    args.agentID = await this.resolveRef('agent', args.agentID, args)
    return this.SystemAPI.agentRead(args, extra).then(r => new Agent(r))
  }

  // Update agent details
  // Mirrors SystemAPI.agentUpdate.
  async updateAgent(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<Agent> {
    const args = argsOf(a, 'agentID')
    args.agentID = await this.resolveRef('agent', args.agentID, args)
    return this.SystemAPI.agentUpdate(args, extra).then(r => new Agent(r))
  }

  // Delete agent
  // Mirrors SystemAPI.agentDelete.
  async deleteAgent(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'agentID')
    args.agentID = await this.resolveRef('agent', args.agentID, args)
    return this.SystemAPI.agentDelete(args, extra)
  }

  // Undelete agent
  // Mirrors SystemAPI.agentUndelete.
  async undeleteAgent(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'agentID')
    args.agentID = await this.resolveRef('agent', args.agentID, args)
    return this.SystemAPI.agentUndelete(args, extra)
  }

  // Execute agent
  async agentExec(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'agentID')
    args.agentID = await this.resolveRef('agent', args.agentID, args)
    return this.SystemAPI.agentExec(args, extra)
  }

  // List LLM providers
  // Mirrors SystemAPI.llmProviderList.
  findLlmProviders(
    a: Args = {},
    extra: AxiosRequestConfig = {},
  ): Promise<ListResponse<LlmProvider[], KV>> {
    const args = argsOf(a)
    return this.SystemAPI.llmProviderList(args, extra).then(r =>
      castSet(r, v => new LlmProvider(v)),
    )
  }

  // Create LLM provider
  // Mirrors SystemAPI.llmProviderCreate.
  createLlmProvider(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<LlmProvider> {
    const args = argsOf(a)
    return this.SystemAPI.llmProviderCreate(args, extra).then(r => new LlmProvider(r))
  }

  // Read LLM provider
  // Mirrors SystemAPI.llmProviderRead.
  async findLlmProviderByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<LlmProvider> {
    const args = argsOf(a, 'llmProviderID')
    args.llmProviderID = await this.resolveRef('llmProvider', args.llmProviderID, args)
    return this.SystemAPI.llmProviderRead(args, extra).then(r => new LlmProvider(r))
  }

  // Update LLM provider
  // Mirrors SystemAPI.llmProviderUpdate.
  async updateLlmProvider(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<LlmProvider> {
    const args = argsOf(a, 'llmProviderID')
    args.llmProviderID = await this.resolveRef('llmProvider', args.llmProviderID, args)
    return this.SystemAPI.llmProviderUpdate(args, extra).then(r => new LlmProvider(r))
  }

  // Delete LLM provider
  // Mirrors SystemAPI.llmProviderDelete.
  async deleteLlmProvider(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'llmProviderID')
    args.llmProviderID = await this.resolveRef('llmProvider', args.llmProviderID, args)
    return this.SystemAPI.llmProviderDelete(args, extra)
  }

  // List available models for an LLM provider
  async llmProviderModels(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'llmProviderID')
    args.llmProviderID = await this.resolveRef('llmProvider', args.llmProviderID, args)
    return this.SystemAPI.llmProviderModels(args, extra)
  }

  // Validate LLM provider connection
  async llmProviderValidate(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'llmProviderID')
    args.llmProviderID = await this.resolveRef('llmProvider', args.llmProviderID, args)
    return this.SystemAPI.llmProviderValidate(args, extra)
  }

  // List knowledge bases
  // Mirrors SystemAPI.knowledgeBaseList.
  findKnowledgeBases(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.knowledgeBaseList(args, extra)
  }

  // Create knowledge base
  // Mirrors SystemAPI.knowledgeBaseCreate.
  createKnowledgeBase(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.knowledgeBaseCreate(args, extra)
  }

  // Read knowledge base
  // Mirrors SystemAPI.knowledgeBaseRead.
  async findKnowledgeBaseByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'knowledgeBaseID')
    args.knowledgeBaseID = await this.resolveRef('knowledgeBase', args.knowledgeBaseID, args)
    return this.SystemAPI.knowledgeBaseRead(args, extra)
  }

  // Update knowledge base
  // Mirrors SystemAPI.knowledgeBaseUpdate.
  async updateKnowledgeBase(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'knowledgeBaseID')
    args.knowledgeBaseID = await this.resolveRef('knowledgeBase', args.knowledgeBaseID, args)
    return this.SystemAPI.knowledgeBaseUpdate(args, extra)
  }

  // Delete knowledge base
  // Mirrors SystemAPI.knowledgeBaseDelete.
  async deleteKnowledgeBase(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'knowledgeBaseID')
    args.knowledgeBaseID = await this.resolveRef('knowledgeBase', args.knowledgeBaseID, args)
    return this.SystemAPI.knowledgeBaseDelete(args, extra)
  }

  // Undelete knowledge base
  // Mirrors SystemAPI.knowledgeBaseUndelete.
  async undeleteKnowledgeBase(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'knowledgeBaseID')
    args.knowledgeBaseID = await this.resolveRef('knowledgeBase', args.knowledgeBaseID, args)
    return this.SystemAPI.knowledgeBaseUndelete(args, extra)
  }

  // List AI conversations
  // Mirrors SystemAPI.aiConversationList.
  findAiConversations(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.aiConversationList(args, extra)
  }

  // Read AI conversation details
  // Mirrors SystemAPI.aiConversationRead.
  async findAiConversationByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'aiConversationID')
    args.aiConversationID = await this.resolveRef('aiConversation', args.aiConversationID, args)
    return this.SystemAPI.aiConversationRead(args, extra)
  }

  // Delete AI conversation
  // Mirrors SystemAPI.aiConversationDelete.
  async deleteAiConversation(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'aiConversationID')
    args.aiConversationID = await this.resolveRef('aiConversation', args.aiConversationID, args)
    return this.SystemAPI.aiConversationDelete(args, extra)
  }

  // Undelete AI conversation
  // Mirrors SystemAPI.aiConversationUndelete.
  async undeleteAiConversation(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'aiConversationID')
    args.aiConversationID = await this.resolveRef('aiConversation', args.aiConversationID, args)
    return this.SystemAPI.aiConversationUndelete(args, extra)
  }

  // Continue AI conversation
  async aiConversationContinue(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'aiConversationID')
    args.aiConversationID = await this.resolveRef('aiConversation', args.aiConversationID, args)
    return this.SystemAPI.aiConversationContinue(args, extra)
  }

  // List chatbots
  // Mirrors SystemAPI.chatbotList.
  findChatbots(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<ListResponse<Chatbot[], KV>> {
    const args = argsOf(a)
    return this.SystemAPI.chatbotList(args, extra).then(r => castSet(r, v => new Chatbot(v)))
  }

  // Create chatbot
  // Mirrors SystemAPI.chatbotCreate.
  createChatbot(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<Chatbot> {
    const args = argsOf(a)
    return this.SystemAPI.chatbotCreate(args, extra).then(r => new Chatbot(r))
  }

  // Read chatbot details
  // Mirrors SystemAPI.chatbotRead.
  async findChatbotByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<Chatbot> {
    const args = argsOf(a, 'chatbotID')
    args.chatbotID = await this.resolveRef('chatbot', args.chatbotID, args)
    return this.SystemAPI.chatbotRead(args, extra).then(r => new Chatbot(r))
  }

  // Update chatbot details
  // Mirrors SystemAPI.chatbotUpdate.
  async updateChatbot(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<Chatbot> {
    const args = argsOf(a, 'chatbotID')
    args.chatbotID = await this.resolveRef('chatbot', args.chatbotID, args)
    return this.SystemAPI.chatbotUpdate(args, extra).then(r => new Chatbot(r))
  }

  // Delete chatbot
  // Mirrors SystemAPI.chatbotDelete.
  async deleteChatbot(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'chatbotID')
    args.chatbotID = await this.resolveRef('chatbot', args.chatbotID, args)
    return this.SystemAPI.chatbotDelete(args, extra)
  }

  // Undelete chatbot
  // Mirrors SystemAPI.chatbotUndelete.
  async undeleteChatbot(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'chatbotID')
    args.chatbotID = await this.resolveRef('chatbot', args.chatbotID, args)
    return this.SystemAPI.chatbotUndelete(args, extra)
  }

  // Regenerate chatbot widget key
  async chatbotRegenerateWidgetKey(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'chatbotID')
    args.chatbotID = await this.resolveRef('chatbot', args.chatbotID, args)
    return this.SystemAPI.chatbotRegenerateWidgetKey(args, extra)
  }

  // Upload chatbot asset
  async chatbotUploadAsset(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'chatbotID')
    args.chatbotID = await this.resolveRef('chatbot', args.chatbotID, args)
    return this.SystemAPI.chatbotUploadAsset(args, extra)
  }

  // List sessions (widget + preview, aggregated)
  // Mirrors SystemAPI.chatbotSessionList.
  findChatbotSessions(
    a: Args = {},
    extra: AxiosRequestConfig = {},
  ): Promise<ListResponse<ChatbotSession[], KV>> {
    const args = argsOf(a)
    return this.SystemAPI.chatbotSessionList(args, extra).then(r =>
      castSet(r, v => new ChatbotSession(v)),
    )
  }

  // Read session with steps
  // Mirrors SystemAPI.chatbotSessionRead.
  async findChatbotSessionByID(
    a: IdArgs = {},
    extra: AxiosRequestConfig = {},
  ): Promise<ChatbotSession> {
    const args = argsOf(a, 'sessionID')
    args.sessionID = await this.resolveRef('chatbotSession', args.sessionID, args)
    return this.SystemAPI.chatbotSessionRead(args, extra).then(r => new ChatbotSession(r))
  }

  // Force-advance session step (admin)
  async chatbotSessionAdvanceStep(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'sessionID')
    args.sessionID = await this.resolveRef('chatbotSession', args.sessionID, args)
    return this.SystemAPI.chatbotSessionAdvanceStep(args, extra)
  }

  // Close session (admin)
  async chatbotSessionClose(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'sessionID')
    args.sessionID = await this.resolveRef('chatbotSession', args.sessionID, args)
    return this.SystemAPI.chatbotSessionClose(args, extra)
  }

  // Accept pending handoff
  async chatbotSessionHandoffAccept(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'sessionID')
    args.sessionID = await this.resolveRef('chatbotSession', args.sessionID, args)
    return this.SystemAPI.chatbotSessionHandoffAccept(args, extra)
  }

  // Send operator message
  async chatbotSessionOperatorMessage(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'sessionID')
    args.sessionID = await this.resolveRef('chatbotSession', args.sessionID, args)
    return this.SystemAPI.chatbotSessionOperatorMessage(args, extra)
  }

  // Complete handoff
  async chatbotSessionHandoffComplete(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'sessionID')
    args.sessionID = await this.resolveRef('chatbotSession', args.sessionID, args)
    return this.SystemAPI.chatbotSessionHandoffComplete(args, extra)
  }

  // SSE stream of session events
  async chatbotSessionStream(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'sessionID')
    args.sessionID = await this.resolveRef('chatbotSession', args.sessionID, args)
    return this.SystemAPI.chatbotSessionStream(args, extra)
  }

  // List projects
  // Mirrors SystemAPI.projectList.
  findProjects(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<ListResponse<Project[], KV>> {
    const args = argsOf(a)
    return this.SystemAPI.projectList(args, extra).then(r => castSet(r, v => new Project(v)))
  }

  // Create project
  // Mirrors SystemAPI.projectCreate.
  createProject(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<Project> {
    const args = argsOf(a)
    return this.SystemAPI.projectCreate(args, extra).then(r => new Project(r))
  }

  // Read project details
  // Mirrors SystemAPI.projectRead.
  async findProjectByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<Project> {
    const args = argsOf(a, 'projectID')
    args.projectID = await this.resolveRef('project', args.projectID, args)
    return this.SystemAPI.projectRead(args, extra).then(r => new Project(r))
  }

  // Update project details
  // Mirrors SystemAPI.projectUpdate.
  async updateProject(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<Project> {
    const args = argsOf(a, 'projectID')
    args.projectID = await this.resolveRef('project', args.projectID, args)
    return this.SystemAPI.projectUpdate(args, extra).then(r => new Project(r))
  }

  // Remove project
  // Mirrors SystemAPI.projectDelete.
  async deleteProject(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'projectID')
    args.projectID = await this.resolveRef('project', args.projectID, args)
    return this.SystemAPI.projectDelete(args, extra)
  }

  // Undelete project
  // Mirrors SystemAPI.projectUndelete.
  async undeleteProject(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'projectID')
    args.projectID = await this.resolveRef('project', args.projectID, args)
    return this.SystemAPI.projectUndelete(args, extra)
  }

  // Archive project
  async projectArchive(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'projectID')
    args.projectID = await this.resolveRef('project', args.projectID, args)
    return this.SystemAPI.projectArchive(args, extra)
  }

  // Restore an archived project
  async projectUnarchive(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'projectID')
    args.projectID = await this.resolveRef('project', args.projectID, args)
    return this.SystemAPI.projectUnarchive(args, extra)
  }

  // List project members
  async projectListMembers(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'projectID')
    args.projectID = await this.resolveRef('project', args.projectID, args)
    return this.SystemAPI.projectListMembers(args, extra)
  }

  // Add project member
  async projectAddMember(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'projectID')
    args.projectID = await this.resolveRef('project', args.projectID, args)
    return this.SystemAPI.projectAddMember(args, extra)
  }

  // Update project member role
  async projectUpdateMember(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.projectID = await this.resolveRef('project', args.projectID, args)
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.projectUpdateMember(args, extra)
  }

  // Remove project member
  async projectRemoveMember(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.projectID = await this.resolveRef('project', args.projectID, args)
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.projectRemoveMember(args, extra)
  }

  // Get project resource graph
  async projectGraph(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'projectID')
    args.projectID = await this.resolveRef('project', args.projectID, args)
    return this.SystemAPI.projectGraph(args, extra)
  }

  // Create a draft revision of the project
  async projectCreateRevision(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'projectID')
    args.projectID = await this.resolveRef('project', args.projectID, args)
    return this.SystemAPI.projectCreateRevision(args, extra)
  }

  // List project revisions
  async projectListRevisions(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'projectID')
    args.projectID = await this.resolveRef('project', args.projectID, args)
    return this.SystemAPI.projectListRevisions(args, extra)
  }

  // Submit a draft revision for publish approval
  async projectRequestApproval(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'projectID')
    args.projectID = await this.resolveRef('project', args.projectID, args)
    return this.SystemAPI.projectRequestApproval(args, extra)
  }

  // Approve a draft revision for publishing
  async projectGrantApproval(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'projectID')
    args.projectID = await this.resolveRef('project', args.projectID, args)
    return this.SystemAPI.projectGrantApproval(args, extra)
  }

  // Send a draft revision back instead of approving it
  async projectRejectApproval(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'projectID')
    args.projectID = await this.resolveRef('project', args.projectID, args)
    return this.SystemAPI.projectRejectApproval(args, extra)
  }

  // Compute deployment plan (diff, risk, suggested mapping) for a draft revision
  async projectGetDeploymentPlan(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'projectID')
    args.projectID = await this.resolveRef('project', args.projectID, args)
    return this.SystemAPI.projectGetDeploymentPlan(args, extra)
  }

  // Publish a draft project
  async projectPublish(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'projectID')
    args.projectID = await this.resolveRef('project', args.projectID, args)
    return this.SystemAPI.projectPublish(args, extra)
  }

  // List AI systems
  // Mirrors SystemAPI.projectAiSystemList.
  async findProjectAiSystems(
    a: IdArgs = {},
    extra: AxiosRequestConfig = {},
  ): Promise<ListResponse<ProjectAiSystem[], KV>> {
    const args = argsOf(a, 'projectID')
    args.projectID = await this.resolveRef('project', args.projectID, args)
    return this.SystemAPI.projectAiSystemList(args, extra).then(r =>
      castSet(r, v => new ProjectAiSystem(v)),
    )
  }

  // Create AI system
  // Mirrors SystemAPI.projectAiSystemCreate.
  async createProjectAiSystem(
    a: IdArgs = {},
    extra: AxiosRequestConfig = {},
  ): Promise<ProjectAiSystem> {
    const args = argsOf(a, 'projectID')
    args.projectID = await this.resolveRef('project', args.projectID, args)
    return this.SystemAPI.projectAiSystemCreate(args, extra).then(r => new ProjectAiSystem(r))
  }

  // Read AI system
  // Mirrors SystemAPI.projectAiSystemRead.
  async findProjectAiSystemByID(
    a: Args = {},
    extra: AxiosRequestConfig = {},
  ): Promise<ProjectAiSystem> {
    const args = argsOf(a)
    args.projectID = await this.resolveRef('project', args.projectID, args)
    args.projectAiSystemID = await this.resolveRef('projectAiSystem', args.projectAiSystemID, args)
    return this.SystemAPI.projectAiSystemRead(args, extra).then(r => new ProjectAiSystem(r))
  }

  // Update AI system
  // Mirrors SystemAPI.projectAiSystemUpdate.
  async updateProjectAiSystem(
    a: Args = {},
    extra: AxiosRequestConfig = {},
  ): Promise<ProjectAiSystem> {
    const args = argsOf(a)
    args.projectID = await this.resolveRef('project', args.projectID, args)
    args.projectAiSystemID = await this.resolveRef('projectAiSystem', args.projectAiSystemID, args)
    return this.SystemAPI.projectAiSystemUpdate(args, extra).then(r => new ProjectAiSystem(r))
  }

  // Delete AI system
  // Mirrors SystemAPI.projectAiSystemDelete.
  async deleteProjectAiSystem(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.projectID = await this.resolveRef('project', args.projectID, args)
    args.projectAiSystemID = await this.resolveRef('projectAiSystem', args.projectAiSystemID, args)
    return this.SystemAPI.projectAiSystemDelete(args, extra)
  }

  // Add resource to AI system
  async projectAiSystemEntryAdd(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.projectID = await this.resolveRef('project', args.projectID, args)
    args.projectAiSystemID = await this.resolveRef('projectAiSystem', args.projectAiSystemID, args)
    return this.SystemAPI.projectAiSystemEntryAdd(args, extra)
  }

  // Remove resource from AI system
  async projectAiSystemEntryRemove(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.projectID = await this.resolveRef('project', args.projectID, args)
    args.projectAiSystemID = await this.resolveRef('projectAiSystem', args.projectAiSystemID, args)
    return this.SystemAPI.projectAiSystemEntryRemove(args, extra)
  }

  // List FRIA risk scenarios
  // Mirrors SystemAPI.projectFriaScenarioList.
  async findProjectFriaScenarios(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'projectID')
    args.projectID = await this.resolveRef('project', args.projectID, args)
    return this.SystemAPI.projectFriaScenarioList(args, extra)
  }

  // Create FRIA risk scenario
  // Mirrors SystemAPI.projectFriaScenarioCreate.
  async createProjectFriaScenario(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'projectID')
    args.projectID = await this.resolveRef('project', args.projectID, args)
    return this.SystemAPI.projectFriaScenarioCreate(args, extra)
  }

  // Read FRIA risk scenario
  // Mirrors SystemAPI.projectFriaScenarioRead.
  async findProjectFriaScenarioByID(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.projectID = await this.resolveRef('project', args.projectID, args)
    args.projectFriaScenarioID = await this.resolveRef(
      'projectFriaScenario',
      args.projectFriaScenarioID,
      args,
    )
    return this.SystemAPI.projectFriaScenarioRead(args, extra)
  }

  // Update FRIA risk scenario
  // Mirrors SystemAPI.projectFriaScenarioUpdate.
  async updateProjectFriaScenario(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.projectID = await this.resolveRef('project', args.projectID, args)
    args.projectFriaScenarioID = await this.resolveRef(
      'projectFriaScenario',
      args.projectFriaScenarioID,
      args,
    )
    return this.SystemAPI.projectFriaScenarioUpdate(args, extra)
  }

  // Delete FRIA risk scenario
  // Mirrors SystemAPI.projectFriaScenarioDelete.
  async deleteProjectFriaScenario(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.projectID = await this.resolveRef('project', args.projectID, args)
    args.projectFriaScenarioID = await this.resolveRef(
      'projectFriaScenario',
      args.projectFriaScenarioID,
      args,
    )
    return this.SystemAPI.projectFriaScenarioDelete(args, extra)
  }

  // List tenants
  // Mirrors SystemAPI.tenantList.
  findTenants(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.tenantList(args, extra)
  }

  // Create tenant
  // Mirrors SystemAPI.tenantCreate.
  createTenant(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.tenantCreate(args, extra)
  }

  // Read tenant details
  // Mirrors SystemAPI.tenantRead.
  async findTenantByID(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'tenantID')
    args.tenantID = await this.resolveRef('tenant', args.tenantID, args)
    return this.SystemAPI.tenantRead(args, extra)
  }

  // Update tenant details
  // Mirrors SystemAPI.tenantUpdate.
  async updateTenant(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'tenantID')
    args.tenantID = await this.resolveRef('tenant', args.tenantID, args)
    return this.SystemAPI.tenantUpdate(args, extra)
  }

  // Remove tenant
  // Mirrors SystemAPI.tenantDelete.
  async deleteTenant(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'tenantID')
    args.tenantID = await this.resolveRef('tenant', args.tenantID, args)
    return this.SystemAPI.tenantDelete(args, extra)
  }

  // Undelete tenant
  // Mirrors SystemAPI.tenantUndelete.
  async undeleteTenant(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'tenantID')
    args.tenantID = await this.resolveRef('tenant', args.tenantID, args)
    return this.SystemAPI.tenantUndelete(args, extra)
  }

  // Suspend tenant
  async tenantSuspend(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'tenantID')
    args.tenantID = await this.resolveRef('tenant', args.tenantID, args)
    return this.SystemAPI.tenantSuspend(args, extra)
  }

  // Activate tenant
  async tenantActivate(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'tenantID')
    args.tenantID = await this.resolveRef('tenant', args.tenantID, args)
    return this.SystemAPI.tenantActivate(args, extra)
  }

  // Archive tenant
  async tenantArchive(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'tenantID')
    args.tenantID = await this.resolveRef('tenant', args.tenantID, args)
    return this.SystemAPI.tenantArchive(args, extra)
  }

  // List tenant members
  async tenantListMembers(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'tenantID')
    args.tenantID = await this.resolveRef('tenant', args.tenantID, args)
    return this.SystemAPI.tenantListMembers(args, extra)
  }

  // Invite tenant member
  async tenantAddMember(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'tenantID')
    args.tenantID = await this.resolveRef('tenant', args.tenantID, args)
    return this.SystemAPI.tenantAddMember(args, extra)
  }

  // Update tenant member
  async tenantUpdateMember(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.tenantID = await this.resolveRef('tenant', args.tenantID, args)
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.tenantUpdateMember(args, extra)
  }

  // Remove tenant member
  async tenantRemoveMember(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.tenantID = await this.resolveRef('tenant', args.tenantID, args)
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.tenantRemoveMember(args, extra)
  }

  // Suspend tenant member
  async tenantSuspendMember(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.tenantID = await this.resolveRef('tenant', args.tenantID, args)
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.tenantSuspendMember(args, extra)
  }

  // Activate tenant member
  async tenantActivateMember(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    args.tenantID = await this.resolveRef('tenant', args.tenantID, args)
    args.userID = await this.resolveRef('user', args.userID, args)
    return this.SystemAPI.tenantActivateMember(args, extra)
  }

  // List DML connections
  dmlConnectionList(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.dmlConnectionList(args, extra)
  }

  // Create DML connection
  dmlConnectionCreate(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.dmlConnectionCreate(args, extra)
  }

  // Read DML connection
  dmlConnectionRead(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'connectionID')
    return this.SystemAPI.dmlConnectionRead(args, extra)
  }

  // List external models for a DML connection
  dmlConnectionModels(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'connectionID')
    return this.SystemAPI.dmlConnectionModels(args, extra)
  }

  // Create DML mapping
  dmlMappingCreate(a: IdArgs = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a, 'connectionID')
    return this.SystemAPI.dmlMappingCreate(args, extra)
  }

  // Read DML mapping
  dmlMappingRead(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.dmlMappingRead(args, extra)
  }

  // Update DML mapping
  dmlMappingUpdate(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.dmlMappingUpdate(args, extra)
  }

  // Delete DML mapping
  dmlMappingDelete(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.dmlMappingDelete(args, extra)
  }

  // Start DML import run (applies schema then imports data)
  dmlImportRun(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.dmlImportRun(args, extra)
  }

  // Read DML import run status
  dmlImportRunRead(a: Args = {}, extra: AxiosRequestConfig = {}): Promise<KV> {
    const args = argsOf(a)
    return this.SystemAPI.dmlImportRunRead(args, extra)
  }
}
