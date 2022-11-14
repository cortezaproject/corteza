/* eslint-disable padded-blocks */

// This is a generated file.
// See README.md file for update instructions

import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'

interface KV {
  [header: string]: unknown;
}

interface Headers {
  [header: string]: string;
}

interface Ctor {
  baseURL?: string;
  accessTokenFn?: () => string | undefined
  headers?: Headers;
}

interface CortezaResponse {
  error?: string;
  response?: unknown;
}

interface ExtraConfig {
  headers?: Headers;
}

function stdResolve (response: AxiosResponse<CortezaResponse>): KV|Promise<never> {
  if (response.data.error) {
    return Promise.reject(response.data.error)
  } else {
    return response.data.response as KV
  }
}

export default class System {
  protected baseURL?: string;
  protected accessTokenFn?: () => (string | undefined);
  protected headers: Headers = {};

  constructor ({ baseURL, headers, accessTokenFn }: Ctor) {
    this.baseURL = baseURL
    this.accessTokenFn = accessTokenFn
    this.headers = {
      /**
       * All we send is JSON
       */
      'Content-Type': 'application/json',
    }

    this.setHeaders(headers)
  }

  setAccessTokenFn (fn: () => string | undefined): System {
    this.accessTokenFn = fn
    return this
  }

  setHeaders (headers?: Headers): System {
    if (typeof headers === 'object') {
      this.headers = headers
    }

    return this
  }

  setHeader (name: string, value: string | undefined): System {
    if (value === undefined) {
      delete this.headers[name]
    } else {
      this.headers[name] = value
    }

    return this
  }

  api (): AxiosInstance {
    const headers = { ...this.headers }
    const accessToken = this.accessTokenFn ? this.accessTokenFn() : undefined
    if (accessToken) {
      headers.Authorization = 'Bearer ' + accessToken
    }

    return axios.create({
      withCredentials: true,
      baseURL: this.baseURL,
      headers,
    })
  }

  // Impersonate a user
  async authImpersonate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      userID,
    } = (a as KV) || {}
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.authImpersonateEndpoint(),
    }
    cfg.data = {
      userID,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  authImpersonateEndpoint (): string {
    return '/auth/impersonate'
  }

  // List clients
  async authClientList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      handle,
      deleted,
      labels,
      limit,
      pageCursor,
      sort,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.authClientListEndpoint(),
    }
    cfg.params = {
      handle,
      deleted,
      labels,
      limit,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  authClientListEndpoint (): string {
    return '/auth/clients/'
  }

  // Create client
  async authClientCreate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      handle,
      meta,
      validGrant,
      redirectURI,
      scope,
      trusted,
      enabled,
      validFrom,
      expiresAt,
      security,
      labels,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.authClientCreateEndpoint(),
    }
    cfg.data = {
      handle,
      meta,
      validGrant,
      redirectURI,
      scope,
      trusted,
      enabled,
      validFrom,
      expiresAt,
      security,
      labels,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  authClientCreateEndpoint (): string {
    return '/auth/clients/'
  }

  // Update user details
  async authClientUpdate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      clientID,
      handle,
      meta,
      validGrant,
      redirectURI,
      scope,
      trusted,
      enabled,
      validFrom,
      expiresAt,
      security,
      labels,
    } = (a as KV) || {}
    if (!clientID) {
      throw Error('field clientID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.authClientUpdateEndpoint({
        clientID,
      }),
    }
    cfg.data = {
      handle,
      meta,
      validGrant,
      redirectURI,
      scope,
      trusted,
      enabled,
      validFrom,
      expiresAt,
      security,
      labels,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  authClientUpdateEndpoint (a: KV): string {
    const {
      clientID,
    } = a || {}
    return `/auth/clients/${clientID}`
  }

  // Read client details
  async authClientRead (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      clientID,
    } = (a as KV) || {}
    if (!clientID) {
      throw Error('field clientID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.authClientReadEndpoint({
        clientID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  authClientReadEndpoint (a: KV): string {
    const {
      clientID,
    } = a || {}
    return `/auth/clients/${clientID}`
  }

  // Remove client
  async authClientDelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      clientID,
    } = (a as KV) || {}
    if (!clientID) {
      throw Error('field clientID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.authClientDeleteEndpoint({
        clientID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  authClientDeleteEndpoint (a: KV): string {
    const {
      clientID,
    } = a || {}
    return `/auth/clients/${clientID}`
  }

  // Undelete client
  async authClientUndelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      clientID,
    } = (a as KV) || {}
    if (!clientID) {
      throw Error('field clientID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.authClientUndeleteEndpoint({
        clientID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  authClientUndeleteEndpoint (a: KV): string {
    const {
      clientID,
    } = a || {}
    return `/auth/clients/${clientID}/undelete`
  }

  // Regenerate client&#x27;s secret
  async authClientRegenerateSecret (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      clientID,
    } = (a as KV) || {}
    if (!clientID) {
      throw Error('field clientID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.authClientRegenerateSecretEndpoint({
        clientID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  authClientRegenerateSecretEndpoint (a: KV): string {
    const {
      clientID,
    } = a || {}
    return `/auth/clients/${clientID}/secret`
  }

  // Exposes client&#x27;s secret
  async authClientExposeSecret (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      clientID,
    } = (a as KV) || {}
    if (!clientID) {
      throw Error('field clientID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.authClientExposeSecretEndpoint({
        clientID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  authClientExposeSecretEndpoint (a: KV): string {
    const {
      clientID,
    } = a || {}
    return `/auth/clients/${clientID}/secret`
  }

  // List settings
  async settingsList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      prefix,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.settingsListEndpoint(),
    }
    cfg.params = {
      prefix,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  settingsListEndpoint (): string {
    return '/settings/'
  }

  // Update settings
  async settingsUpdate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      values,
    } = (a as KV) || {}
    if (!values) {
      throw Error('field values is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'patch',
      url: this.settingsUpdateEndpoint(),
    }
    cfg.data = {
      values,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  settingsUpdateEndpoint (): string {
    return '/settings/'
  }

  // Get a value for a key
  async settingsGet (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      key,
      ownerID,
    } = (a as KV) || {}
    if (!key) {
      throw Error('field key is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.settingsGetEndpoint({
        key,
      }),
    }
    cfg.params = {
      ownerID,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  settingsGetEndpoint (a: KV): string {
    const {
      key,
    } = a || {}
    return `/settings/${key}`
  }

  // Set value for specific setting
  async settingsSet (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      key,
      upload,
      ownerID,
    } = (a as KV) || {}
    if (!key) {
      throw Error('field key is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.settingsSetEndpoint({
        key,
      }),
    }
    cfg.data = {
      upload,
      ownerID,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  settingsSetEndpoint (a: KV): string {
    const {
      key,
    } = a || {}
    return `/settings/${key}`
  }

  // Current compose settings
  async settingsCurrent (extra: AxiosRequestConfig = {}): Promise<KV> {

    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.settingsCurrentEndpoint(),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  settingsCurrentEndpoint (): string {
    return '/settings/current'
  }

  // List roles
  async roleList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      query,
      memberID,
      deleted,
      archived,
      labels,
      limit,
      pageCursor,
      sort,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.roleListEndpoint(),
    }
    cfg.params = {
      query,
      memberID,
      deleted,
      archived,
      labels,
      limit,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  roleListEndpoint (): string {
    return '/roles/'
  }

  // Update role details
  async roleCreate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      name,
      handle,
      members,
      meta,
      labels,
    } = (a as KV) || {}
    if (!name) {
      throw Error('field name is empty')
    }
    if (!handle) {
      throw Error('field handle is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.roleCreateEndpoint(),
    }
    cfg.data = {
      name,
      handle,
      members,
      meta,
      labels,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  roleCreateEndpoint (): string {
    return '/roles/'
  }

  // Update role details
  async roleUpdate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      roleID,
      name,
      handle,
      members,
      meta,
      labels,
    } = (a as KV) || {}
    if (!roleID) {
      throw Error('field roleID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.roleUpdateEndpoint({
        roleID,
      }),
    }
    cfg.data = {
      name,
      handle,
      members,
      meta,
      labels,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  roleUpdateEndpoint (a: KV): string {
    const {
      roleID,
    } = a || {}
    return `/roles/${roleID}`
  }

  // Read role details and memberships
  async roleRead (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      roleID,
    } = (a as KV) || {}
    if (!roleID) {
      throw Error('field roleID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.roleReadEndpoint({
        roleID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  roleReadEndpoint (a: KV): string {
    const {
      roleID,
    } = a || {}
    return `/roles/${roleID}`
  }

  // Remove role
  async roleDelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      roleID,
    } = (a as KV) || {}
    if (!roleID) {
      throw Error('field roleID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.roleDeleteEndpoint({
        roleID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  roleDeleteEndpoint (a: KV): string {
    const {
      roleID,
    } = a || {}
    return `/roles/${roleID}`
  }

  // Archive role
  async roleArchive (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      roleID,
    } = (a as KV) || {}
    if (!roleID) {
      throw Error('field roleID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.roleArchiveEndpoint({
        roleID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  roleArchiveEndpoint (a: KV): string {
    const {
      roleID,
    } = a || {}
    return `/roles/${roleID}/archive`
  }

  // Unarchive role
  async roleUnarchive (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      roleID,
    } = (a as KV) || {}
    if (!roleID) {
      throw Error('field roleID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.roleUnarchiveEndpoint({
        roleID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  roleUnarchiveEndpoint (a: KV): string {
    const {
      roleID,
    } = a || {}
    return `/roles/${roleID}/unarchive`
  }

  // Undelete role
  async roleUndelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      roleID,
    } = (a as KV) || {}
    if (!roleID) {
      throw Error('field roleID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.roleUndeleteEndpoint({
        roleID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  roleUndeleteEndpoint (a: KV): string {
    const {
      roleID,
    } = a || {}
    return `/roles/${roleID}/undelete`
  }

  // Move role to different organisation
  async roleMove (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      roleID,
      organisationID,
    } = (a as KV) || {}
    if (!roleID) {
      throw Error('field roleID is empty')
    }
    if (!organisationID) {
      throw Error('field organisationID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.roleMoveEndpoint({
        roleID,
      }),
    }
    cfg.data = {
      organisationID,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  roleMoveEndpoint (a: KV): string {
    const {
      roleID,
    } = a || {}
    return `/roles/${roleID}/move`
  }

  // Merge one role into another
  async roleMerge (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      roleID,
      destination,
    } = (a as KV) || {}
    if (!roleID) {
      throw Error('field roleID is empty')
    }
    if (!destination) {
      throw Error('field destination is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.roleMergeEndpoint({
        roleID,
      }),
    }
    cfg.data = {
      destination,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  roleMergeEndpoint (a: KV): string {
    const {
      roleID,
    } = a || {}
    return `/roles/${roleID}/merge`
  }

  // Returns all role members
  async roleMemberList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      roleID,
    } = (a as KV) || {}
    if (!roleID) {
      throw Error('field roleID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.roleMemberListEndpoint({
        roleID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  roleMemberListEndpoint (a: KV): string {
    const {
      roleID,
    } = a || {}
    return `/roles/${roleID}/members`
  }

  // Add member to a role
  async roleMemberAdd (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      roleID,
      userID,
    } = (a as KV) || {}
    if (!roleID) {
      throw Error('field roleID is empty')
    }
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.roleMemberAddEndpoint({
        roleID, userID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  roleMemberAddEndpoint (a: KV): string {
    const {
      roleID,
      userID,
    } = a || {}
    return `/roles/${roleID}/member/${userID}`
  }

  // Remove member from a role
  async roleMemberRemove (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      roleID,
      userID,
    } = (a as KV) || {}
    if (!roleID) {
      throw Error('field roleID is empty')
    }
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.roleMemberRemoveEndpoint({
        roleID, userID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  roleMemberRemoveEndpoint (a: KV): string {
    const {
      roleID,
      userID,
    } = a || {}
    return `/roles/${roleID}/member/${userID}`
  }

  // Fire system:role trigger
  async roleTriggerScript (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      roleID,
      script,
      args,
    } = (a as KV) || {}
    if (!roleID) {
      throw Error('field roleID is empty')
    }
    if (!script) {
      throw Error('field script is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.roleTriggerScriptEndpoint({
        roleID,
      }),
    }
    cfg.data = {
      script,
      args,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  roleTriggerScriptEndpoint (a: KV): string {
    const {
      roleID,
    } = a || {}
    return `/roles/${roleID}/trigger`
  }

  // Search users (Directory)
  async userList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      userID,
      roleID,
      query,
      username,
      email,
      handle,
      kind,
      incDeleted,
      incSuspended,
      deleted,
      suspended,
      labels,
      limit,
      pageCursor,
      sort,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.userListEndpoint(),
    }
    cfg.params = {
      userID,
      roleID,
      query,
      username,
      email,
      handle,
      kind,
      incDeleted,
      incSuspended,
      deleted,
      suspended,
      labels,
      limit,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  userListEndpoint (): string {
    return '/users/'
  }

  // Create user
  async userCreate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      email,
      name,
      handle,
      kind,
      labels,
    } = (a as KV) || {}
    if (!email) {
      throw Error('field email is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.userCreateEndpoint(),
    }
    cfg.data = {
      email,
      name,
      handle,
      kind,
      labels,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  userCreateEndpoint (): string {
    return '/users/'
  }

  // Update user details
  async userUpdate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      userID,
      email,
      name,
      handle,
      kind,
      labels,
    } = (a as KV) || {}
    if (!userID) {
      throw Error('field userID is empty')
    }
    if (!email) {
      throw Error('field email is empty')
    }
    if (!name) {
      throw Error('field name is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.userUpdateEndpoint({
        userID,
      }),
    }
    cfg.data = {
      email,
      name,
      handle,
      kind,
      labels,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  userUpdateEndpoint (a: KV): string {
    const {
      userID,
    } = a || {}
    return `/users/${userID}`
  }

  // Patch user (experimental)
  async userPartialUpdate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      userID,
    } = (a as KV) || {}
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'patch',
      url: this.userPartialUpdateEndpoint({
        userID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  userPartialUpdateEndpoint (a: KV): string {
    const {
      userID,
    } = a || {}
    return `/users/${userID}`
  }

  // Read user details
  async userRead (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      userID,
    } = (a as KV) || {}
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.userReadEndpoint({
        userID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  userReadEndpoint (a: KV): string {
    const {
      userID,
    } = a || {}
    return `/users/${userID}`
  }

  // Remove user
  async userDelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      userID,
    } = (a as KV) || {}
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.userDeleteEndpoint({
        userID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  userDeleteEndpoint (a: KV): string {
    const {
      userID,
    } = a || {}
    return `/users/${userID}`
  }

  // Suspend user
  async userSuspend (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      userID,
    } = (a as KV) || {}
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.userSuspendEndpoint({
        userID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  userSuspendEndpoint (a: KV): string {
    const {
      userID,
    } = a || {}
    return `/users/${userID}/suspend`
  }

  // Unsuspend user
  async userUnsuspend (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      userID,
    } = (a as KV) || {}
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.userUnsuspendEndpoint({
        userID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  userUnsuspendEndpoint (a: KV): string {
    const {
      userID,
    } = a || {}
    return `/users/${userID}/unsuspend`
  }

  // Undelete user
  async userUndelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      userID,
    } = (a as KV) || {}
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.userUndeleteEndpoint({
        userID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  userUndeleteEndpoint (a: KV): string {
    const {
      userID,
    } = a || {}
    return `/users/${userID}/undelete`
  }

  // Set&#x27;s or changes user&#x27;s password
  async userSetPassword (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      userID,
      password,
    } = (a as KV) || {}
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.userSetPasswordEndpoint({
        userID,
      }),
    }
    cfg.data = {
      password,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  userSetPasswordEndpoint (a: KV): string {
    const {
      userID,
    } = a || {}
    return `/users/${userID}/password`
  }

  // Add member to a role
  async userMembershipList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      userID,
    } = (a as KV) || {}
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.userMembershipListEndpoint({
        userID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  userMembershipListEndpoint (a: KV): string {
    const {
      userID,
    } = a || {}
    return `/users/${userID}/membership`
  }

  // Add role to a user
  async userMembershipAdd (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      roleID,
      userID,
    } = (a as KV) || {}
    if (!roleID) {
      throw Error('field roleID is empty')
    }
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.userMembershipAddEndpoint({
        roleID, userID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  userMembershipAddEndpoint (a: KV): string {
    const {
      roleID,
      userID,
    } = a || {}
    return `/users/${userID}/membership/${roleID}`
  }

  // Remove role from a user
  async userMembershipRemove (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      roleID,
      userID,
    } = (a as KV) || {}
    if (!roleID) {
      throw Error('field roleID is empty')
    }
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.userMembershipRemoveEndpoint({
        roleID, userID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  userMembershipRemoveEndpoint (a: KV): string {
    const {
      roleID,
      userID,
    } = a || {}
    return `/users/${userID}/membership/${roleID}`
  }

  // Fire system:user trigger
  async userTriggerScript (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      userID,
      script,
      args,
    } = (a as KV) || {}
    if (!userID) {
      throw Error('field userID is empty')
    }
    if (!script) {
      throw Error('field script is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.userTriggerScriptEndpoint({
        userID,
      }),
    }
    cfg.data = {
      script,
      args,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  userTriggerScriptEndpoint (a: KV): string {
    const {
      userID,
    } = a || {}
    return `/users/${userID}/trigger`
  }

  // Remove all auth sessions of user
  async userSessionsRemove (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      userID,
    } = (a as KV) || {}
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.userSessionsRemoveEndpoint({
        userID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  userSessionsRemoveEndpoint (a: KV): string {
    const {
      userID,
    } = a || {}
    return `/users/${userID}/sessions`
  }

  // Export users
  async userExport (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      filename,
      inclRoleMembership,
      inclRoles,
    } = (a as KV) || {}
    if (!filename) {
      throw Error('field filename is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.userExportEndpoint({
        filename,
      }),
    }
    cfg.params = {
      inclRoleMembership,
      inclRoles,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  userExportEndpoint (a: KV): string {
    const {
      filename,
    } = a || {}
    return `/users/export/${filename}.zip`
  }

  // Import users
  async userImport (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      upload,
    } = (a as KV) || {}
    if (!upload) {
      throw Error('field upload is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.userImportEndpoint(),
    }
    cfg.data = {
      upload,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  userImportEndpoint (): string {
    return '/users/import'
  }

  // List applications
  async applicationList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      name,
      query,
      deleted,
      labels,
      flags,
      incFlags,
      limit,
      pageCursor,
      sort,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.applicationListEndpoint(),
    }
    cfg.params = {
      name,
      query,
      deleted,
      labels,
      flags,
      incFlags,
      limit,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  applicationListEndpoint (): string {
    return '/application/'
  }

  // Create application
  async applicationCreate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      name,
      enabled,
      weight,
      unify,
      config,
      labels,
    } = (a as KV) || {}
    if (!name) {
      throw Error('field name is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.applicationCreateEndpoint(),
    }
    cfg.data = {
      name,
      enabled,
      weight,
      unify,
      config,
      labels,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  applicationCreateEndpoint (): string {
    return '/application/'
  }

  // Update user details
  async applicationUpdate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      applicationID,
      name,
      enabled,
      weight,
      unify,
      config,
      labels,
    } = (a as KV) || {}
    if (!applicationID) {
      throw Error('field applicationID is empty')
    }
    if (!name) {
      throw Error('field name is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.applicationUpdateEndpoint({
        applicationID,
      }),
    }
    cfg.data = {
      name,
      enabled,
      weight,
      unify,
      config,
      labels,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  applicationUpdateEndpoint (a: KV): string {
    const {
      applicationID,
    } = a || {}
    return `/application/${applicationID}`
  }

  // Upload application assets
  async applicationUpload (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      upload,
    } = (a as KV) || {}
    if (!upload) {
      throw Error('field upload is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.applicationUploadEndpoint(),
    }
    cfg.data = {
      upload,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  applicationUploadEndpoint (): string {
    return '/application/upload'
  }

  // Flag application
  async applicationFlagCreate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      applicationID,
      flag,
      ownedBy,
    } = (a as KV) || {}
    if (!applicationID) {
      throw Error('field applicationID is empty')
    }
    if (!flag) {
      throw Error('field flag is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.applicationFlagCreateEndpoint({
        applicationID, flag, ownedBy,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  applicationFlagCreateEndpoint (a: KV): string {
    const {
      applicationID,
      flag,
      ownedBy,
    } = a || {}
    return `/application/${applicationID}/flag/${ownedBy}/${flag}`
  }

  // Unflag application
  async applicationFlagDelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      applicationID,
      flag,
      ownedBy,
    } = (a as KV) || {}
    if (!applicationID) {
      throw Error('field applicationID is empty')
    }
    if (!flag) {
      throw Error('field flag is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.applicationFlagDeleteEndpoint({
        applicationID, flag, ownedBy,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  applicationFlagDeleteEndpoint (a: KV): string {
    const {
      applicationID,
      flag,
      ownedBy,
    } = a || {}
    return `/application/${applicationID}/flag/${ownedBy}/${flag}`
  }

  // Read application details
  async applicationRead (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      applicationID,
      incFlags,
    } = (a as KV) || {}
    if (!applicationID) {
      throw Error('field applicationID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.applicationReadEndpoint({
        applicationID,
      }),
    }
    cfg.params = {
      incFlags,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  applicationReadEndpoint (a: KV): string {
    const {
      applicationID,
    } = a || {}
    return `/application/${applicationID}`
  }

  // Remove application
  async applicationDelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      applicationID,
    } = (a as KV) || {}
    if (!applicationID) {
      throw Error('field applicationID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.applicationDeleteEndpoint({
        applicationID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  applicationDeleteEndpoint (a: KV): string {
    const {
      applicationID,
    } = a || {}
    return `/application/${applicationID}`
  }

  // Undelete application
  async applicationUndelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      applicationID,
    } = (a as KV) || {}
    if (!applicationID) {
      throw Error('field applicationID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.applicationUndeleteEndpoint({
        applicationID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  applicationUndeleteEndpoint (a: KV): string {
    const {
      applicationID,
    } = a || {}
    return `/application/${applicationID}/undelete`
  }

  // Fire system:application trigger
  async applicationTriggerScript (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      applicationID,
      script,
      args,
    } = (a as KV) || {}
    if (!applicationID) {
      throw Error('field applicationID is empty')
    }
    if (!script) {
      throw Error('field script is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.applicationTriggerScriptEndpoint({
        applicationID,
      }),
    }
    cfg.data = {
      script,
      args,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  applicationTriggerScriptEndpoint (a: KV): string {
    const {
      applicationID,
    } = a || {}
    return `/application/${applicationID}/trigger`
  }

  // Reorder applications
  async applicationReorder (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      applicationIDs,
    } = (a as KV) || {}
    if (!applicationIDs) {
      throw Error('field applicationIDs is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.applicationReorderEndpoint(),
    }
    cfg.data = {
      applicationIDs,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  applicationReorderEndpoint (): string {
    return '/application/reorder'
  }

  // Retrieve defined permissions
  async permissionsList (extra: AxiosRequestConfig = {}): Promise<KV> {

    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.permissionsListEndpoint(),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  permissionsListEndpoint (): string {
    return '/permissions/'
  }

  // Effective rules for current user
  async permissionsEffective (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      resource,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.permissionsEffectiveEndpoint(),
    }
    cfg.params = {
      resource,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  permissionsEffectiveEndpoint (): string {
    return '/permissions/effective'
  }

  // Retrieve role permissions
  async permissionsRead (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      roleID,
    } = (a as KV) || {}
    if (!roleID) {
      throw Error('field roleID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.permissionsReadEndpoint({
        roleID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  permissionsReadEndpoint (a: KV): string {
    const {
      roleID,
    } = a || {}
    return `/permissions/${roleID}/rules`
  }

  // Remove all defined role permissions
  async permissionsDelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      roleID,
    } = (a as KV) || {}
    if (!roleID) {
      throw Error('field roleID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.permissionsDeleteEndpoint({
        roleID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  permissionsDeleteEndpoint (a: KV): string {
    const {
      roleID,
    } = a || {}
    return `/permissions/${roleID}/rules`
  }

  // Update permission settings
  async permissionsUpdate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      roleID,
      rules,
    } = (a as KV) || {}
    if (!roleID) {
      throw Error('field roleID is empty')
    }
    if (!rules) {
      throw Error('field rules is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'patch',
      url: this.permissionsUpdateEndpoint({
        roleID,
      }),
    }
    cfg.data = {
      rules,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  permissionsUpdateEndpoint (a: KV): string {
    const {
      roleID,
    } = a || {}
    return `/permissions/${roleID}/rules`
  }

  // Clone permission settings to a role
  async permissionsClone (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      roleID,
      cloneToRoleID,
    } = (a as KV) || {}
    if (!roleID) {
      throw Error('field roleID is empty')
    }
    if (!cloneToRoleID) {
      throw Error('field cloneToRoleID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.permissionsCloneEndpoint({
        roleID,
      }),
    }
    cfg.params = {
      cloneToRoleID,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  permissionsCloneEndpoint (a: KV): string {
    const {
      roleID,
    } = a || {}
    return `/permissions/${roleID}/rules/clone`
  }

  // List/read reminders
  async reminderList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      reminderID,
      resource,
      assignedTo,
      scheduledFrom,
      scheduledUntil,
      scheduledOnly,
      excludeDismissed,
      includeDeleted,
      limit,
      pageCursor,
      sort,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.reminderListEndpoint(),
    }
    cfg.params = {
      reminderID,
      resource,
      assignedTo,
      scheduledFrom,
      scheduledUntil,
      scheduledOnly,
      excludeDismissed,
      includeDeleted,
      limit,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  reminderListEndpoint (): string {
    return '/reminder/'
  }

  // Add new reminder
  async reminderCreate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      resource,
      assignedTo,
      payload,
      remindAt,
    } = (a as KV) || {}
    if (!resource) {
      throw Error('field resource is empty')
    }
    if (!assignedTo) {
      throw Error('field assignedTo is empty')
    }
    if (!payload) {
      throw Error('field payload is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.reminderCreateEndpoint(),
    }
    cfg.data = {
      resource,
      assignedTo,
      payload,
      remindAt,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  reminderCreateEndpoint (): string {
    return '/reminder/'
  }

  // Update reminder
  async reminderUpdate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      reminderID,
      resource,
      assignedTo,
      payload,
      remindAt,
    } = (a as KV) || {}
    if (!reminderID) {
      throw Error('field reminderID is empty')
    }
    if (!resource) {
      throw Error('field resource is empty')
    }
    if (!assignedTo) {
      throw Error('field assignedTo is empty')
    }
    if (!payload) {
      throw Error('field payload is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.reminderUpdateEndpoint({
        reminderID,
      }),
    }
    cfg.data = {
      resource,
      assignedTo,
      payload,
      remindAt,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  reminderUpdateEndpoint (a: KV): string {
    const {
      reminderID,
    } = a || {}
    return `/reminder/${reminderID}`
  }

  // Read reminder by ID
  async reminderRead (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      reminderID,
    } = (a as KV) || {}
    if (!reminderID) {
      throw Error('field reminderID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.reminderReadEndpoint({
        reminderID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  reminderReadEndpoint (a: KV): string {
    const {
      reminderID,
    } = a || {}
    return `/reminder/${reminderID}`
  }

  // Delete reminder
  async reminderDelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      reminderID,
    } = (a as KV) || {}
    if (!reminderID) {
      throw Error('field reminderID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.reminderDeleteEndpoint({
        reminderID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  reminderDeleteEndpoint (a: KV): string {
    const {
      reminderID,
    } = a || {}
    return `/reminder/${reminderID}`
  }

  // Dismiss reminder
  async reminderDismiss (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      reminderID,
    } = (a as KV) || {}
    if (!reminderID) {
      throw Error('field reminderID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'patch',
      url: this.reminderDismissEndpoint({
        reminderID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  reminderDismissEndpoint (a: KV): string {
    const {
      reminderID,
    } = a || {}
    return `/reminder/${reminderID}/dismiss`
  }

  // Snooze reminder
  async reminderSnooze (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      reminderID,
      remindAt,
    } = (a as KV) || {}
    if (!reminderID) {
      throw Error('field reminderID is empty')
    }
    if (!remindAt) {
      throw Error('field remindAt is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'patch',
      url: this.reminderSnoozeEndpoint({
        reminderID,
      }),
    }
    cfg.data = {
      remindAt,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  reminderSnoozeEndpoint (a: KV): string {
    const {
      reminderID,
    } = a || {}
    return `/reminder/${reminderID}/snooze`
  }

  // Attachment details
  async attachmentRead (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      kind,
      attachmentID,
      sign,
      userID,
    } = (a as KV) || {}
    if (!kind) {
      throw Error('field kind is empty')
    }
    if (!attachmentID) {
      throw Error('field attachmentID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.attachmentReadEndpoint({
        kind, attachmentID,
      }),
    }
    cfg.params = {
      sign,
      userID,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  attachmentReadEndpoint (a: KV): string {
    const {
      kind,
      attachmentID,
    } = a || {}
    return `/attachment/${kind}/${attachmentID}`
  }

  // Delete attachment
  async attachmentDelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      kind,
      attachmentID,
      sign,
      userID,
    } = (a as KV) || {}
    if (!kind) {
      throw Error('field kind is empty')
    }
    if (!attachmentID) {
      throw Error('field attachmentID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.attachmentDeleteEndpoint({
        kind, attachmentID,
      }),
    }
    cfg.params = {
      sign,
      userID,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  attachmentDeleteEndpoint (a: KV): string {
    const {
      kind,
      attachmentID,
    } = a || {}
    return `/attachment/${kind}/${attachmentID}`
  }

  // Serves attached file
  async attachmentOriginal (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      kind,
      attachmentID,
      name,
      sign,
      userID,
      download,
    } = (a as KV) || {}
    if (!kind) {
      throw Error('field kind is empty')
    }
    if (!attachmentID) {
      throw Error('field attachmentID is empty')
    }
    if (!name) {
      throw Error('field name is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.attachmentOriginalEndpoint({
        kind, attachmentID, name,
      }),
    }
    cfg.params = {
      sign,
      userID,
      download,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  attachmentOriginalEndpoint (a: KV): string {
    const {
      kind,
      attachmentID,
      name,
    } = a || {}
    return `/attachment/${kind}/${attachmentID}/original/${name}`
  }

  // Serves preview of an attached file
  async attachmentPreview (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      kind,
      attachmentID,
      ext,
      sign,
      userID,
    } = (a as KV) || {}
    if (!kind) {
      throw Error('field kind is empty')
    }
    if (!attachmentID) {
      throw Error('field attachmentID is empty')
    }
    if (!ext) {
      throw Error('field ext is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.attachmentPreviewEndpoint({
        kind, attachmentID, ext,
      }),
    }
    cfg.params = {
      sign,
      userID,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  attachmentPreviewEndpoint (a: KV): string {
    const {
      kind,
      attachmentID,
      ext,
    } = a || {}
    return `/attachment/${kind}/${attachmentID}/preview.${ext}`
  }

  // List templates
  async templateList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      handle,
      type,
      ownerID,
      partial,
      deleted,
      labels,
      limit,
      pageCursor,
      sort,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.templateListEndpoint(),
    }
    cfg.params = {
      handle,
      type,
      ownerID,
      partial,
      deleted,
      labels,
      limit,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  templateListEndpoint (): string {
    return '/template/'
  }

  // Create template
  async templateCreate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      handle,
      language,
      type,
      partial,
      meta,
      template,
      ownerID,
      labels,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.templateCreateEndpoint(),
    }
    cfg.data = {
      handle,
      language,
      type,
      partial,
      meta,
      template,
      ownerID,
      labels,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  templateCreateEndpoint (): string {
    return '/template/'
  }

  // Read template
  async templateRead (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      templateID,
    } = (a as KV) || {}
    if (!templateID) {
      throw Error('field templateID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.templateReadEndpoint({
        templateID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  templateReadEndpoint (a: KV): string {
    const {
      templateID,
    } = a || {}
    return `/template/${templateID}`
  }

  // Update template
  async templateUpdate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      templateID,
      handle,
      language,
      type,
      partial,
      meta,
      template,
      ownerID,
      labels,
    } = (a as KV) || {}
    if (!templateID) {
      throw Error('field templateID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.templateUpdateEndpoint({
        templateID,
      }),
    }
    cfg.data = {
      handle,
      language,
      type,
      partial,
      meta,
      template,
      ownerID,
      labels,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  templateUpdateEndpoint (a: KV): string {
    const {
      templateID,
    } = a || {}
    return `/template/${templateID}`
  }

  // Delete template
  async templateDelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      templateID,
    } = (a as KV) || {}
    if (!templateID) {
      throw Error('field templateID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.templateDeleteEndpoint({
        templateID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  templateDeleteEndpoint (a: KV): string {
    const {
      templateID,
    } = a || {}
    return `/template/${templateID}`
  }

  // Undelete template
  async templateUndelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      templateID,
    } = (a as KV) || {}
    if (!templateID) {
      throw Error('field templateID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.templateUndeleteEndpoint({
        templateID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  templateUndeleteEndpoint (a: KV): string {
    const {
      templateID,
    } = a || {}
    return `/template/${templateID}/undelete`
  }

  // Render drivers
  async templateRenderDrivers (extra: AxiosRequestConfig = {}): Promise<KV> {

    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.templateRenderDriversEndpoint(),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  templateRenderDriversEndpoint (): string {
    return '/template/render/drivers'
  }

  // Render template
  async templateRender (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      templateID,
      filename,
      ext,
      variables,
      options,
    } = (a as KV) || {}
    if (!templateID) {
      throw Error('field templateID is empty')
    }
    if (!filename) {
      throw Error('field filename is empty')
    }
    if (!ext) {
      throw Error('field ext is empty')
    }
    if (!variables) {
      throw Error('field variables is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.templateRenderEndpoint({
        templateID, filename, ext,
      }),
    }
    cfg.data = {
      variables,
      options,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  templateRenderEndpoint (a: KV): string {
    const {
      templateID,
      filename,
      ext,
    } = a || {}
    return `/template/${templateID}/render/${filename}.${ext}`
  }

  // List reports
  async reportList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      handle,
      deleted,
      labels,
      limit,
      pageCursor,
      sort,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.reportListEndpoint(),
    }
    cfg.params = {
      handle,
      deleted,
      labels,
      limit,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  reportListEndpoint (): string {
    return '/reports/'
  }

  // Create report
  async reportCreate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      handle,
      meta,
      scenarios,
      sources,
      blocks,
      labels,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.reportCreateEndpoint(),
    }
    cfg.data = {
      handle,
      meta,
      scenarios,
      sources,
      blocks,
      labels,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  reportCreateEndpoint (): string {
    return '/reports/'
  }

  // Update report
  async reportUpdate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      reportID,
      handle,
      meta,
      scenarios,
      sources,
      blocks,
      labels,
    } = (a as KV) || {}
    if (!reportID) {
      throw Error('field reportID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.reportUpdateEndpoint({
        reportID,
      }),
    }
    cfg.data = {
      handle,
      meta,
      scenarios,
      sources,
      blocks,
      labels,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  reportUpdateEndpoint (a: KV): string {
    const {
      reportID,
    } = a || {}
    return `/reports/${reportID}`
  }

  // Read report details
  async reportRead (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      reportID,
    } = (a as KV) || {}
    if (!reportID) {
      throw Error('field reportID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.reportReadEndpoint({
        reportID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  reportReadEndpoint (a: KV): string {
    const {
      reportID,
    } = a || {}
    return `/reports/${reportID}`
  }

  // Remove report
  async reportDelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      reportID,
    } = (a as KV) || {}
    if (!reportID) {
      throw Error('field reportID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.reportDeleteEndpoint({
        reportID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  reportDeleteEndpoint (a: KV): string {
    const {
      reportID,
    } = a || {}
    return `/reports/${reportID}`
  }

  // Undelete report
  async reportUndelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      reportID,
    } = (a as KV) || {}
    if (!reportID) {
      throw Error('field reportID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.reportUndeleteEndpoint({
        reportID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  reportUndeleteEndpoint (a: KV): string {
    const {
      reportID,
    } = a || {}
    return `/reports/${reportID}/undelete`
  }

  // Describe report
  async reportDescribe (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      sources,
      steps,
      describe,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.reportDescribeEndpoint(),
    }
    cfg.data = {
      sources,
      steps,
      describe,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  reportDescribeEndpoint (): string {
    return '/reports/describe'
  }

  // Run report
  async reportRun (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      reportID,
      frames,
    } = (a as KV) || {}
    if (!reportID) {
      throw Error('field reportID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.reportRunEndpoint({
        reportID,
      }),
    }
    cfg.data = {
      frames,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  reportRunEndpoint (a: KV): string {
    const {
      reportID,
    } = a || {}
    return `/reports/${reportID}/run`
  }

  // List system statistics
  async statsList (extra: AxiosRequestConfig = {}): Promise<KV> {

    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.statsListEndpoint(),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  statsListEndpoint (): string {
    return '/stats/'
  }

  // List all available automation scripts for system resources
  async automationList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      resourceTypePrefixes,
      resourceTypes,
      eventTypes,
      excludeInvalid,
      excludeClientScripts,
      excludeServerScripts,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.automationListEndpoint(),
    }
    cfg.params = {
      resourceTypePrefixes,
      resourceTypes,
      eventTypes,
      excludeInvalid,
      excludeClientScripts,
      excludeServerScripts,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  automationListEndpoint (): string {
    return '/automation/'
  }

  // Serves client scripts bundle
  async automationBundle (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      bundle,
      type,
      ext,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.automationBundleEndpoint({
        bundle, type, ext,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  automationBundleEndpoint (a: KV): string {
    const {
      bundle,
      type,
      ext,
    } = a || {}
    return `/automation/${bundle}-${type}.${ext}`
  }

  // Triggers execution of a specific script on a system service level
  async automationTriggerScript (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      script,
      args,
    } = (a as KV) || {}
    if (!script) {
      throw Error('field script is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.automationTriggerScriptEndpoint(),
    }
    cfg.data = {
      script,
      args,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  automationTriggerScriptEndpoint (): string {
    return '/automation/trigger'
  }

  // Action log events
  async actionlogList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      from,
      to,
      beforeActionID,
      resource,
      action,
      actorID,
      limit,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.actionlogListEndpoint(),
    }
    cfg.params = {
      from,
      to,
      beforeActionID,
      resource,
      action,
      actorID,
      limit,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  actionlogListEndpoint (): string {
    return '/actionlog/'
  }

  // Messaging queues
  async queuesList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      query,
      limit,
      pageCursor,
      sort,
      deleted,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.queuesListEndpoint(),
    }
    cfg.params = {
      query,
      limit,
      pageCursor,
      sort,
      deleted,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  queuesListEndpoint (): string {
    return '/queues/'
  }

  // Create messaging queue
  async queuesCreate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      queue,
      consumer,
      meta,
    } = (a as KV) || {}
    if (!queue) {
      throw Error('field queue is empty')
    }
    if (!consumer) {
      throw Error('field consumer is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.queuesCreateEndpoint(),
    }
    cfg.data = {
      queue,
      consumer,
      meta,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  queuesCreateEndpoint (): string {
    return '/queues'
  }

  // Messaging queue details
  async queuesRead (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      queueID,
    } = (a as KV) || {}
    if (!queueID) {
      throw Error('field queueID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.queuesReadEndpoint({
        queueID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  queuesReadEndpoint (a: KV): string {
    const {
      queueID,
    } = a || {}
    return `/queues/${queueID}`
  }

  // Update queue details
  async queuesUpdate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      queueID,
      queue,
      consumer,
      meta,
    } = (a as KV) || {}
    if (!queueID) {
      throw Error('field queueID is empty')
    }
    if (!queue) {
      throw Error('field queue is empty')
    }
    if (!consumer) {
      throw Error('field consumer is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.queuesUpdateEndpoint({
        queueID,
      }),
    }
    cfg.data = {
      queue,
      consumer,
      meta,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  queuesUpdateEndpoint (a: KV): string {
    const {
      queueID,
    } = a || {}
    return `/queues/${queueID}`
  }

  // Messaging queue delete
  async queuesDelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      queueID,
    } = (a as KV) || {}
    if (!queueID) {
      throw Error('field queueID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.queuesDeleteEndpoint({
        queueID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  queuesDeleteEndpoint (a: KV): string {
    const {
      queueID,
    } = a || {}
    return `/queues/${queueID}`
  }

  // Messaging queue undelete
  async queuesUndelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      queueID,
    } = (a as KV) || {}
    if (!queueID) {
      throw Error('field queueID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.queuesUndeleteEndpoint({
        queueID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  queuesUndeleteEndpoint (a: KV): string {
    const {
      queueID,
    } = a || {}
    return `/queues/${queueID}/undelete`
  }

  // List routes
  async apigwRouteList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      routeID,
      query,
      deleted,
      disabled,
      labels,
      limit,
      pageCursor,
      sort,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.apigwRouteListEndpoint(),
    }
    cfg.params = {
      routeID,
      query,
      deleted,
      disabled,
      labels,
      limit,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  apigwRouteListEndpoint (): string {
    return '/apigw/route/'
  }

  // Create route
  async apigwRouteCreate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      endpoint,
      method,
      enabled,
      group,
      meta,
    } = (a as KV) || {}
    if (!endpoint) {
      throw Error('field endpoint is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.apigwRouteCreateEndpoint(),
    }
    cfg.data = {
      endpoint,
      method,
      enabled,
      group,
      meta,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  apigwRouteCreateEndpoint (): string {
    return '/apigw/route'
  }

  // Update route details
  async apigwRouteUpdate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      routeID,
      endpoint,
      method,
      enabled,
      group,
      meta,
    } = (a as KV) || {}
    if (!routeID) {
      throw Error('field routeID is empty')
    }
    if (!endpoint) {
      throw Error('field endpoint is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.apigwRouteUpdateEndpoint({
        routeID,
      }),
    }
    cfg.data = {
      endpoint,
      method,
      enabled,
      group,
      meta,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  apigwRouteUpdateEndpoint (a: KV): string {
    const {
      routeID,
    } = a || {}
    return `/apigw/route/${routeID}`
  }

  // Read route details
  async apigwRouteRead (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      routeID,
    } = (a as KV) || {}
    if (!routeID) {
      throw Error('field routeID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.apigwRouteReadEndpoint({
        routeID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  apigwRouteReadEndpoint (a: KV): string {
    const {
      routeID,
    } = a || {}
    return `/apigw/route/${routeID}`
  }

  // Remove route
  async apigwRouteDelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      routeID,
    } = (a as KV) || {}
    if (!routeID) {
      throw Error('field routeID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.apigwRouteDeleteEndpoint({
        routeID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  apigwRouteDeleteEndpoint (a: KV): string {
    const {
      routeID,
    } = a || {}
    return `/apigw/route/${routeID}`
  }

  // Undelete route
  async apigwRouteUndelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      routeID,
    } = (a as KV) || {}
    if (!routeID) {
      throw Error('field routeID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.apigwRouteUndeleteEndpoint({
        routeID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  apigwRouteUndeleteEndpoint (a: KV): string {
    const {
      routeID,
    } = a || {}
    return `/apigw/route/${routeID}/undelete`
  }

  // List filters
  async apigwFilterList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      routeID,
      deleted,
      disabled,
      limit,
      pageCursor,
      sort,
    } = (a as KV) || {}
    if (!routeID) {
      throw Error('field routeID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.apigwFilterListEndpoint(),
    }
    cfg.params = {
      routeID,
      deleted,
      disabled,
      limit,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  apigwFilterListEndpoint (): string {
    return '/apigw/filter/'
  }

  // Create filter
  async apigwFilterCreate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      routeID,
      weight,
      kind,
      ref,
      enabled,
      params,
    } = (a as KV) || {}
    if (!routeID) {
      throw Error('field routeID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.apigwFilterCreateEndpoint(),
    }
    cfg.data = {
      routeID,
      weight,
      kind,
      ref,
      enabled,
      params,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  apigwFilterCreateEndpoint (): string {
    return '/apigw/filter'
  }

  // Update filter details
  async apigwFilterUpdate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      filterID,
      routeID,
      weight,
      kind,
      ref,
      enabled,
      params,
    } = (a as KV) || {}
    if (!filterID) {
      throw Error('field filterID is empty')
    }
    if (!routeID) {
      throw Error('field routeID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.apigwFilterUpdateEndpoint({
        filterID,
      }),
    }
    cfg.data = {
      routeID,
      weight,
      kind,
      ref,
      enabled,
      params,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  apigwFilterUpdateEndpoint (a: KV): string {
    const {
      filterID,
    } = a || {}
    return `/apigw/filter/${filterID}`
  }

  // Read filter details
  async apigwFilterRead (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      filterID,
    } = (a as KV) || {}
    if (!filterID) {
      throw Error('field filterID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.apigwFilterReadEndpoint({
        filterID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  apigwFilterReadEndpoint (a: KV): string {
    const {
      filterID,
    } = a || {}
    return `/apigw/filter/${filterID}`
  }

  // Remove filter
  async apigwFilterDelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      filterID,
    } = (a as KV) || {}
    if (!filterID) {
      throw Error('field filterID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.apigwFilterDeleteEndpoint({
        filterID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  apigwFilterDeleteEndpoint (a: KV): string {
    const {
      filterID,
    } = a || {}
    return `/apigw/filter/${filterID}`
  }

  // Undelete filter
  async apigwFilterUndelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      filterID,
    } = (a as KV) || {}
    if (!filterID) {
      throw Error('field filterID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.apigwFilterUndeleteEndpoint({
        filterID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  apigwFilterUndeleteEndpoint (a: KV): string {
    const {
      filterID,
    } = a || {}
    return `/apigw/filter/${filterID}/undelete`
  }

  // Filter definitions
  async apigwFilterDefFilter (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      kind,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.apigwFilterDefFilterEndpoint(),
    }
    cfg.params = {
      kind,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  apigwFilterDefFilterEndpoint (): string {
    return '/apigw/filter/def'
  }

  // Proxy auth definitions
  async apigwFilterDefProxyAuth (extra: AxiosRequestConfig = {}): Promise<KV> {

    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.apigwFilterDefProxyAuthEndpoint(),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  apigwFilterDefProxyAuthEndpoint (): string {
    return '/apigw/filter/proxy_auth/def'
  }

  // List aggregated list of routes
  async apigwProfilerAggregation (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      path,
      before,
      sort,
      limit,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.apigwProfilerAggregationEndpoint(),
    }
    cfg.params = {
      path,
      before,
      sort,
      limit,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  apigwProfilerAggregationEndpoint (): string {
    return '/apigw/profiler/'
  }

  // List hits per route
  async apigwProfilerRoute (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      routeID,
      path,
      before,
      sort,
      limit,
    } = (a as KV) || {}
    if (!routeID) {
      throw Error('field routeID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.apigwProfilerRouteEndpoint({
        routeID,
      }),
    }
    cfg.params = {
      path,
      before,
      sort,
      limit,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  apigwProfilerRouteEndpoint (a: KV): string {
    const {
      routeID,
    } = a || {}
    return `/apigw/profiler/route/${routeID}`
  }

  // Hit details
  async apigwProfilerHit (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      hitID,
    } = (a as KV) || {}
    if (!hitID) {
      throw Error('field hitID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.apigwProfilerHitEndpoint({
        hitID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  apigwProfilerHitEndpoint (a: KV): string {
    const {
      hitID,
    } = a || {}
    return `/apigw/profiler/hit/${hitID}`
  }

  // List resources translations
  async localeListResource (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      lang,
      resource,
      resourceType,
      ownerID,
      deleted,
      limit,
      pageCursor,
      sort,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.localeListResourceEndpoint(),
    }
    cfg.params = {
      lang,
      resource,
      resourceType,
      ownerID,
      deleted,
      limit,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  localeListResourceEndpoint (): string {
    return '/locale/resource'
  }

  // Create resource translation
  async localeCreateResource (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      lang,
      resource,
      key,
      place,
      message,
      ownerID,
    } = (a as KV) || {}
    if (!lang) {
      throw Error('field lang is empty')
    }
    if (!resource) {
      throw Error('field resource is empty')
    }
    if (!key) {
      throw Error('field key is empty')
    }
    if (!message) {
      throw Error('field message is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.localeCreateResourceEndpoint(),
    }
    cfg.data = {
      lang,
      resource,
      key,
      place,
      message,
      ownerID,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  localeCreateResourceEndpoint (): string {
    return '/locale/resource'
  }

  // Update resource translation
  async localeUpdateResource (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      translationID,
      lang,
      resource,
      key,
      place,
      message,
      ownerID,
    } = (a as KV) || {}
    if (!translationID) {
      throw Error('field translationID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.localeUpdateResourceEndpoint({
        translationID,
      }),
    }
    cfg.data = {
      lang,
      resource,
      key,
      place,
      message,
      ownerID,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  localeUpdateResourceEndpoint (a: KV): string {
    const {
      translationID,
    } = a || {}
    return `/locale/resource/${translationID}`
  }

  // Read resource translation details
  async localeReadResource (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      translationID,
    } = (a as KV) || {}
    if (!translationID) {
      throw Error('field translationID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.localeReadResourceEndpoint({
        translationID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  localeReadResourceEndpoint (a: KV): string {
    const {
      translationID,
    } = a || {}
    return `/locale/resource/${translationID}`
  }

  // Remove resource translation
  async localeDeleteResource (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      translationID,
    } = (a as KV) || {}
    if (!translationID) {
      throw Error('field translationID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.localeDeleteResourceEndpoint({
        translationID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  localeDeleteResourceEndpoint (a: KV): string {
    const {
      translationID,
    } = a || {}
    return `/locale/resource/${translationID}`
  }

  // Undelete resource translation
  async localeUndeleteResource (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      translationID,
    } = (a as KV) || {}
    if (!translationID) {
      throw Error('field translationID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.localeUndeleteResourceEndpoint({
        translationID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  localeUndeleteResourceEndpoint (a: KV): string {
    const {
      translationID,
    } = a || {}
    return `/locale/resource/${translationID}/undelete`
  }

  // List all available languages
  async localeList (extra: AxiosRequestConfig = {}): Promise<KV> {

    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.localeListEndpoint(),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  localeListEndpoint (): string {
    return '/locale/'
  }

  // List all available translation in a language for a specific webapp
  async localeGet (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      lang,
      application,
    } = (a as KV) || {}
    if (!lang) {
      throw Error('field lang is empty')
    }
    if (!application) {
      throw Error('field application is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.localeGetEndpoint({
        lang, application,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  localeGetEndpoint (a: KV): string {
    const {
      lang,
      application,
    } = a || {}
    return `/locale/${lang}/${application}`
  }

}
