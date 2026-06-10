// This is a generated file.
// See README.md file for update instructions

import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'

interface KV {
  [header: string]: unknown
}

interface Headers {
  [header: string]: string
}

interface Ctor {
  baseURL?: string
  accessTokenFn?: () => string | undefined
  headers?: Headers
}

interface HumanResponse {
  error?: string
  response?: unknown
}

function stdResolve(response: AxiosResponse<HumanResponse>): KV | Promise<never> {
  if (response.data.error) {
    return Promise.reject(response.data.error)
  } else if (response.data.response) {
    return response.data.response as KV
  } else {
    return response.data as KV
  }
}

export default class System {
  protected baseURL?: string
  protected accessTokenFn?: () => string | undefined
  protected headers: Headers = {}

  constructor({ baseURL, headers, accessTokenFn }: Ctor) {
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

  setAccessTokenFn(fn: () => string | undefined): System {
    this.accessTokenFn = fn
    return this
  }

  setHeaders(headers?: Headers): System {
    if (typeof headers === 'object') {
      this.headers = headers
    }

    return this
  }

  setHeader(name: string, value: string | undefined): System {
    if (value === undefined) {
      delete this.headers[name]
    } else {
      this.headers[name] = value
    }

    return this
  }

  api(): AxiosInstance {
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
  async authImpersonate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { userID } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  authImpersonateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.authImpersonate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  authImpersonateEndpoint(): string {
    return '/auth/impersonate'
  }

  // List clients
  async authClientList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { handle, deleted, labels, limit, incTotal, pageCursor, sort } = (a as KV) || {}
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
      incTotal,
      pageCursor,
      sort,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  authClientListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.authClientList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  authClientListEndpoint(): string {
    return '/auth/clients/'
  }

  // Create client
  async authClientCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  authClientCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.authClientCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  authClientCreateEndpoint(): string {
    return '/auth/clients/'
  }

  // Update user details
  async authClientUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
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
      updatedAt,
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
      updatedAt,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  authClientUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.authClientUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  authClientUpdateEndpoint(a: KV): string {
    const { clientID } = a || {}
    return `/auth/clients/${clientID}`
  }

  // Read client details
  async authClientRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { clientID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  authClientReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.authClientRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  authClientReadEndpoint(a: KV): string {
    const { clientID } = a || {}
    return `/auth/clients/${clientID}`
  }

  // Remove client
  async authClientDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { clientID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  authClientDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.authClientDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  authClientDeleteEndpoint(a: KV): string {
    const { clientID } = a || {}
    return `/auth/clients/${clientID}`
  }

  // Undelete client
  async authClientUndelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { clientID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  authClientUndeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.authClientUndelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  authClientUndeleteEndpoint(a: KV): string {
    const { clientID } = a || {}
    return `/auth/clients/${clientID}/undelete`
  }

  // Regenerate client&#x27;s secret
  async authClientRegenerateSecret(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { clientID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  authClientRegenerateSecretCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.authClientRegenerateSecret(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  authClientRegenerateSecretEndpoint(a: KV): string {
    const { clientID } = a || {}
    return `/auth/clients/${clientID}/secret`
  }

  // Exposes client&#x27;s secret
  async authClientExposeSecret(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { clientID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  authClientExposeSecretCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.authClientExposeSecret(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  authClientExposeSecretEndpoint(a: KV): string {
    const { clientID } = a || {}
    return `/auth/clients/${clientID}/secret`
  }

  // Evaluate expressions
  async expressionEvaluate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { variables, expressions } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.expressionEvaluateEndpoint(),
    }
    cfg.data = {
      variables,
      expressions,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  expressionEvaluateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.expressionEvaluate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  expressionEvaluateEndpoint(): string {
    return '/expressions/evaluate'
  }

  // List settings
  async settingsList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { prefix } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.settingsListEndpoint(),
    }
    cfg.params = {
      prefix,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  settingsListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.settingsList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  settingsListEndpoint(): string {
    return '/settings/'
  }

  // Update settings
  async settingsUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { values } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  settingsUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.settingsUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  settingsUpdateEndpoint(): string {
    return '/settings/'
  }

  // Get a value for a key
  async settingsGet(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { key, ownerID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  settingsGetCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.settingsGet(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  settingsGetEndpoint(a: KV): string {
    const { key } = a || {}
    return `/settings/${key}`
  }

  // Set value for specific setting
  async settingsSet(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { key, upload, ownerID } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  settingsSetCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.settingsSet(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  settingsSetEndpoint(a: KV): string {
    const { key } = a || {}
    return `/settings/${key}`
  }

  // Current compose settings
  async settingsCurrent(extra: AxiosRequestConfig = {}): Promise<KV> {
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.settingsCurrentEndpoint(),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  settingsCurrentCancellable(extra: AxiosRequestConfig = {}): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.settingsCurrent(options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  settingsCurrentEndpoint(): string {
    return '/settings/current'
  }

  // List roles
  async roleList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      query,
      memberID,
      userGroupID,
      roleID,
      deleted,
      archived,
      labels,
      limit,
      incTotal,
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
      userGroupID,
      roleID,
      deleted,
      archived,
      labels,
      limit,
      incTotal,
      pageCursor,
      sort,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  roleListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.roleList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  roleListEndpoint(): string {
    return '/roles/'
  }

  // Update role details
  async roleCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { name, handle, members, meta, labels } = (a as KV) || {}
    if (!name) {
      throw Error('field name is empty')
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  roleCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.roleCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  roleCreateEndpoint(): string {
    return '/roles/'
  }

  // Update role details
  async roleUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { roleID, name, handle, members, meta, labels, updatedAt } = (a as KV) || {}
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
      updatedAt,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  roleUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.roleUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  roleUpdateEndpoint(a: KV): string {
    const { roleID } = a || {}
    return `/roles/${roleID}`
  }

  // Read role details and memberships
  async roleRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { roleID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  roleReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.roleRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  roleReadEndpoint(a: KV): string {
    const { roleID } = a || {}
    return `/roles/${roleID}`
  }

  // Remove role
  async roleDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { roleID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  roleDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.roleDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  roleDeleteEndpoint(a: KV): string {
    const { roleID } = a || {}
    return `/roles/${roleID}`
  }

  // Archive role
  async roleArchive(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { roleID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  roleArchiveCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.roleArchive(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  roleArchiveEndpoint(a: KV): string {
    const { roleID } = a || {}
    return `/roles/${roleID}/archive`
  }

  // Unarchive role
  async roleUnarchive(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { roleID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  roleUnarchiveCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.roleUnarchive(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  roleUnarchiveEndpoint(a: KV): string {
    const { roleID } = a || {}
    return `/roles/${roleID}/unarchive`
  }

  // Undelete role
  async roleUndelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { roleID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  roleUndeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.roleUndelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  roleUndeleteEndpoint(a: KV): string {
    const { roleID } = a || {}
    return `/roles/${roleID}/undelete`
  }

  // Move role to different organisation
  async roleMove(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { roleID, organisationID } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  roleMoveCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.roleMove(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  roleMoveEndpoint(a: KV): string {
    const { roleID } = a || {}
    return `/roles/${roleID}/move`
  }

  // Merge one role into another
  async roleMerge(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { roleID, destination } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  roleMergeCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.roleMerge(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  roleMergeEndpoint(a: KV): string {
    const { roleID } = a || {}
    return `/roles/${roleID}/merge`
  }

  // Returns all role members
  async roleMemberList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { roleID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  roleMemberListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.roleMemberList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  roleMemberListEndpoint(a: KV): string {
    const { roleID } = a || {}
    return `/roles/${roleID}/members`
  }

  // Add user group to a user group
  async roleMemberAddGroup(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { roleID, userGroupID } = (a as KV) || {}
    if (!roleID) {
      throw Error('field roleID is empty')
    }
    if (!userGroupID) {
      throw Error('field userGroupID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.roleMemberAddGroupEndpoint({
        roleID,
        userGroupID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  roleMemberAddGroupCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.roleMemberAddGroup(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  roleMemberAddGroupEndpoint(a: KV): string {
    const { roleID, userGroupID } = a || {}
    return `/roles/${roleID}/member-g/${userGroupID}`
  }

  // Add member to a role
  async roleMemberAdd(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { roleID, userID } = (a as KV) || {}
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
        roleID,
        userID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  roleMemberAddCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.roleMemberAdd(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  roleMemberAddEndpoint(a: KV): string {
    const { roleID, userID } = a || {}
    return `/roles/${roleID}/member/${userID}`
  }

  // Remove member from a role
  async roleMemberRemove(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { roleID, userID } = (a as KV) || {}
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
        roleID,
        userID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  roleMemberRemoveCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.roleMemberRemove(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  roleMemberRemoveEndpoint(a: KV): string {
    const { roleID, userID } = a || {}
    return `/roles/${roleID}/member/${userID}`
  }

  // Remove user group from a role
  async roleMemberRemoveGroup(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { roleID, userGroupID } = (a as KV) || {}
    if (!roleID) {
      throw Error('field roleID is empty')
    }
    if (!userGroupID) {
      throw Error('field userGroupID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.roleMemberRemoveGroupEndpoint({
        roleID,
        userGroupID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  roleMemberRemoveGroupCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.roleMemberRemoveGroup(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  roleMemberRemoveGroupEndpoint(a: KV): string {
    const { roleID, userGroupID } = a || {}
    return `/roles/${roleID}/member-g/${userGroupID}`
  }

  // Fire system:role trigger
  async roleTriggerScript(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { roleID, script, args } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  roleTriggerScriptCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.roleTriggerScript(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  roleTriggerScriptEndpoint(a: KV): string {
    const { roleID } = a || {}
    return `/roles/${roleID}/trigger`
  }

  // Clone permission settings to a role
  async roleCloneRules(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { roleID, cloneToRoleID } = (a as KV) || {}
    if (!roleID) {
      throw Error('field roleID is empty')
    }
    if (!cloneToRoleID) {
      throw Error('field cloneToRoleID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.roleCloneRulesEndpoint({
        roleID,
      }),
    }
    cfg.params = {
      cloneToRoleID,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  roleCloneRulesCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.roleCloneRules(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  roleCloneRulesEndpoint(a: KV): string {
    const { roleID } = a || {}
    return `/roles/${roleID}/rules/clone`
  }

  // List user groups
  async userGroupList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      query,
      memberID,
      userGroupID,
      deleted,
      archived,
      labels,
      limit,
      incTotal,
      pageCursor,
      sort,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.userGroupListEndpoint(),
    }
    cfg.params = {
      query,
      memberID,
      userGroupID,
      deleted,
      archived,
      labels,
      limit,
      incTotal,
      pageCursor,
      sort,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userGroupListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userGroupList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userGroupListEndpoint(): string {
    return '/user-groups/'
  }

  // Update user groups details
  async userGroupCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { handle, members, config, meta, labels } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.userGroupCreateEndpoint(),
    }
    cfg.data = {
      handle,
      members,
      config,
      meta,
      labels,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userGroupCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userGroupCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userGroupCreateEndpoint(): string {
    return '/user-groups/'
  }

  // Update user group details
  async userGroupUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { userGroupID, handle, members, config, meta, labels, updatedAt } = (a as KV) || {}
    if (!userGroupID) {
      throw Error('field userGroupID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.userGroupUpdateEndpoint({
        userGroupID,
      }),
    }
    cfg.data = {
      handle,
      members,
      config,
      meta,
      labels,
      updatedAt,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userGroupUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userGroupUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userGroupUpdateEndpoint(a: KV): string {
    const { userGroupID } = a || {}
    return `/user-groups/${userGroupID}`
  }

  // Read user group details and memberships
  async userGroupRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { userGroupID } = (a as KV) || {}
    if (!userGroupID) {
      throw Error('field userGroupID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.userGroupReadEndpoint({
        userGroupID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userGroupReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userGroupRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userGroupReadEndpoint(a: KV): string {
    const { userGroupID } = a || {}
    return `/user-groups/${userGroupID}`
  }

  // Remove user group
  async userGroupDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { userGroupID } = (a as KV) || {}
    if (!userGroupID) {
      throw Error('field userGroupID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.userGroupDeleteEndpoint({
        userGroupID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userGroupDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userGroupDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userGroupDeleteEndpoint(a: KV): string {
    const { userGroupID } = a || {}
    return `/user-groups/${userGroupID}`
  }

  // Undelete user group
  async userGroupUndelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { userGroupID } = (a as KV) || {}
    if (!userGroupID) {
      throw Error('field userGroupID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.userGroupUndeleteEndpoint({
        userGroupID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userGroupUndeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userGroupUndelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userGroupUndeleteEndpoint(a: KV): string {
    const { userGroupID } = a || {}
    return `/user-groups/${userGroupID}/undelete`
  }

  // Returns all user group members
  async userGroupMemberList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { userGroupID } = (a as KV) || {}
    if (!userGroupID) {
      throw Error('field userGroupID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.userGroupMemberListEndpoint({
        userGroupID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userGroupMemberListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userGroupMemberList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userGroupMemberListEndpoint(a: KV): string {
    const { userGroupID } = a || {}
    return `/user-groups/${userGroupID}/members`
  }

  // Add member to a user group
  async userGroupMemberAdd(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { userGroupID, userID } = (a as KV) || {}
    if (!userGroupID) {
      throw Error('field userGroupID is empty')
    }
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.userGroupMemberAddEndpoint({
        userGroupID,
        userID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userGroupMemberAddCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userGroupMemberAdd(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userGroupMemberAddEndpoint(a: KV): string {
    const { userGroupID, userID } = a || {}
    return `/user-groups/${userGroupID}/member/${userID}`
  }

  // Search users (Directory)
  async userList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
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
      incTotal,
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
      incTotal,
      pageCursor,
      sort,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userListEndpoint(): string {
    return '/users/'
  }

  // Create user
  async userCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { email, name, handle, userGroupID, kind, labels, meta } = (a as KV) || {}
    if (!email) {
      throw Error('field email is empty')
    }
    if (!userGroupID) {
      throw Error('field userGroupID is empty')
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
      userGroupID,
      kind,
      labels,
      meta,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userCreateEndpoint(): string {
    return '/users/'
  }

  // Update user details
  async userUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { userID, email, name, handle, userGroupID, kind, labels, meta, updatedAt } =
      (a as KV) || {}
    if (!userID) {
      throw Error('field userID is empty')
    }
    if (!email) {
      throw Error('field email is empty')
    }
    if (!userGroupID) {
      throw Error('field userGroupID is empty')
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
      userGroupID,
      kind,
      labels,
      meta,
      updatedAt,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userUpdateEndpoint(a: KV): string {
    const { userID } = a || {}
    return `/users/${userID}`
  }

  // Patch user (experimental)
  async userPartialUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { userID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userPartialUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userPartialUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userPartialUpdateEndpoint(a: KV): string {
    const { userID } = a || {}
    return `/users/${userID}`
  }

  // Read user details
  async userRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { userID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userReadEndpoint(a: KV): string {
    const { userID } = a || {}
    return `/users/${userID}`
  }

  // Remove user
  async userDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { userID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userDeleteEndpoint(a: KV): string {
    const { userID } = a || {}
    return `/users/${userID}`
  }

  // Suspend user
  async userSuspend(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { userID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userSuspendCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userSuspend(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userSuspendEndpoint(a: KV): string {
    const { userID } = a || {}
    return `/users/${userID}/suspend`
  }

  // Unsuspend user
  async userUnsuspend(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { userID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userUnsuspendCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userUnsuspend(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userUnsuspendEndpoint(a: KV): string {
    const { userID } = a || {}
    return `/users/${userID}/unsuspend`
  }

  // Undelete user
  async userUndelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { userID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userUndeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userUndelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userUndeleteEndpoint(a: KV): string {
    const { userID } = a || {}
    return `/users/${userID}/undelete`
  }

  // Set&#x27;s or changes user&#x27;s password
  async userSetPassword(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { userID, password } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userSetPasswordCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userSetPassword(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userSetPasswordEndpoint(a: KV): string {
    const { userID } = a || {}
    return `/users/${userID}/password`
  }

  // Add member to a role
  async userMembershipList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { userID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userMembershipListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userMembershipList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userMembershipListEndpoint(a: KV): string {
    const { userID } = a || {}
    return `/users/${userID}/membership`
  }

  // Add role to a user
  async userMembershipAdd(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { roleID, userID } = (a as KV) || {}
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
        roleID,
        userID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userMembershipAddCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userMembershipAdd(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userMembershipAddEndpoint(a: KV): string {
    const { roleID, userID } = a || {}
    return `/users/${userID}/membership/${roleID}`
  }

  // Remove role from a user
  async userMembershipRemove(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { roleID, userID } = (a as KV) || {}
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
        roleID,
        userID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userMembershipRemoveCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userMembershipRemove(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userMembershipRemoveEndpoint(a: KV): string {
    const { roleID, userID } = a || {}
    return `/users/${userID}/membership/${roleID}`
  }

  // Fire system:user trigger
  async userTriggerScript(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { userID, script, args } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userTriggerScriptCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userTriggerScript(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userTriggerScriptEndpoint(a: KV): string {
    const { userID } = a || {}
    return `/users/${userID}/trigger`
  }

  // Remove all auth sessions of user
  async userSessionsRemove(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { userID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userSessionsRemoveCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userSessionsRemove(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userSessionsRemoveEndpoint(a: KV): string {
    const { userID } = a || {}
    return `/users/${userID}/sessions`
  }

  // List user&#x27;s credentials
  async userListCredentials(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { userID } = (a as KV) || {}
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.userListCredentialsEndpoint({
        userID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userListCredentialsCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userListCredentials(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userListCredentialsEndpoint(a: KV): string {
    const { userID } = a || {}
    return `/users/${userID}/credentials`
  }

  // List user&#x27;s credentials
  async userDeleteCredentials(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { userID, credentialsID } = (a as KV) || {}
    if (!userID) {
      throw Error('field userID is empty')
    }
    if (!credentialsID) {
      throw Error('field credentialsID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.userDeleteCredentialsEndpoint({
        userID,
        credentialsID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userDeleteCredentialsCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userDeleteCredentials(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userDeleteCredentialsEndpoint(a: KV): string {
    const { userID, credentialsID } = a || {}
    return `/users/${userID}/credentials/${credentialsID}`
  }

  // User&#x27;s profile avatar
  async userProfileAvatar(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { userID, upload, width, height } = (a as KV) || {}
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.userProfileAvatarEndpoint({
        userID,
      }),
    }
    cfg.data = {
      upload,
      width,
      height,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userProfileAvatarCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userProfileAvatar(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userProfileAvatarEndpoint(a: KV): string {
    const { userID } = a || {}
    return `/users/${userID}/avatar`
  }

  // User profile avatar initial
  async userProfileAvatarInitial(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { userID, avatarColor, avatarBgColor } = (a as KV) || {}
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.userProfileAvatarInitialEndpoint({
        userID,
      }),
    }
    cfg.data = {
      avatarColor,
      avatarBgColor,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userProfileAvatarInitialCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userProfileAvatarInitial(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userProfileAvatarInitialEndpoint(a: KV): string {
    const { userID } = a || {}
    return `/users/${userID}/avatar-initial`
  }

  // delete user&#x27;s profile avatar
  async userDeleteAvatar(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { userID } = (a as KV) || {}
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.userDeleteAvatarEndpoint({
        userID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userDeleteAvatarCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userDeleteAvatar(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userDeleteAvatarEndpoint(a: KV): string {
    const { userID } = a || {}
    return `/users/${userID}/avatar`
  }

  // Export users
  async userExport(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { filename, inclRoleMembership, inclRoles } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userExportCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userExport(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userExportEndpoint(a: KV): string {
    const { filename } = a || {}
    return `/users/export/${filename}.zip`
  }

  // Import users
  async userImport(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { upload } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  userImportCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.userImport(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  userImportEndpoint(): string {
    return '/users/import'
  }

  // Search drivers
  async dalDriverList(extra: AxiosRequestConfig = {}): Promise<KV> {
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.dalDriverListEndpoint(),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dalDriverListCancellable(extra: AxiosRequestConfig = {}): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dalDriverList(options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dalDriverListEndpoint(): string {
    return '/dal/drivers/'
  }

  // Search sensitivity levels
  async dalSensitivityLevelList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { sensitivityLevelID, deleted, incTotal } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.dalSensitivityLevelListEndpoint(),
    }
    cfg.params = {
      sensitivityLevelID,
      deleted,
      incTotal,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dalSensitivityLevelListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dalSensitivityLevelList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dalSensitivityLevelListEndpoint(): string {
    return '/dal/sensitivity-levels/'
  }

  // Create sensitivity level
  async dalSensitivityLevelCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { handle, level, meta } = (a as KV) || {}
    if (!level) {
      throw Error('field level is empty')
    }
    if (!meta) {
      throw Error('field meta is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.dalSensitivityLevelCreateEndpoint(),
    }
    cfg.data = {
      handle,
      level,
      meta,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dalSensitivityLevelCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dalSensitivityLevelCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dalSensitivityLevelCreateEndpoint(): string {
    return '/dal/sensitivity-levels/'
  }

  // Update sensitivity details
  async dalSensitivityLevelUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { sensitivityLevelID, handle, level, meta, updatedAt } = (a as KV) || {}
    if (!sensitivityLevelID) {
      throw Error('field sensitivityLevelID is empty')
    }
    if (!level) {
      throw Error('field level is empty')
    }
    if (!meta) {
      throw Error('field meta is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.dalSensitivityLevelUpdateEndpoint({
        sensitivityLevelID,
      }),
    }
    cfg.data = {
      handle,
      level,
      meta,
      updatedAt,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dalSensitivityLevelUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dalSensitivityLevelUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dalSensitivityLevelUpdateEndpoint(a: KV): string {
    const { sensitivityLevelID } = a || {}
    return `/dal/sensitivity-levels/${sensitivityLevelID}`
  }

  // Read connection details
  async dalSensitivityLevelRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { sensitivityLevelID } = (a as KV) || {}
    if (!sensitivityLevelID) {
      throw Error('field sensitivityLevelID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.dalSensitivityLevelReadEndpoint({
        sensitivityLevelID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dalSensitivityLevelReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dalSensitivityLevelRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dalSensitivityLevelReadEndpoint(a: KV): string {
    const { sensitivityLevelID } = a || {}
    return `/dal/sensitivity-levels/${sensitivityLevelID}`
  }

  // Remove sensitivity level
  async dalSensitivityLevelDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { sensitivityLevelID } = (a as KV) || {}
    if (!sensitivityLevelID) {
      throw Error('field sensitivityLevelID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.dalSensitivityLevelDeleteEndpoint({
        sensitivityLevelID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dalSensitivityLevelDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dalSensitivityLevelDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dalSensitivityLevelDeleteEndpoint(a: KV): string {
    const { sensitivityLevelID } = a || {}
    return `/dal/sensitivity-levels/${sensitivityLevelID}`
  }

  // Undelete sensitivity level
  async dalSensitivityLevelUndelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { sensitivityLevelID } = (a as KV) || {}
    if (!sensitivityLevelID) {
      throw Error('field sensitivityLevelID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.dalSensitivityLevelUndeleteEndpoint({
        sensitivityLevelID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dalSensitivityLevelUndeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dalSensitivityLevelUndelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dalSensitivityLevelUndeleteEndpoint(a: KV): string {
    const { sensitivityLevelID } = a || {}
    return `/dal/sensitivity-levels/${sensitivityLevelID}/undelete`
  }

  // Search schema alterations
  async dalSchemaAlterationList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      alterationID,
      batchID,
      resource,
      resourceType,
      kind,
      deleted,
      completed,
      dismissed,
      incTotal,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.dalSchemaAlterationListEndpoint(),
    }
    cfg.params = {
      alterationID,
      batchID,
      resource,
      resourceType,
      kind,
      deleted,
      completed,
      dismissed,
      incTotal,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dalSchemaAlterationListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dalSchemaAlterationList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dalSchemaAlterationListEndpoint(): string {
    return '/dal/schema/alterations/'
  }

  // Read alteration details
  async dalSchemaAlterationRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { alterationID } = (a as KV) || {}
    if (!alterationID) {
      throw Error('field alterationID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.dalSchemaAlterationReadEndpoint({
        alterationID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dalSchemaAlterationReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dalSchemaAlterationRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dalSchemaAlterationReadEndpoint(a: KV): string {
    const { alterationID } = a || {}
    return `/dal/schema/alterations/${alterationID}`
  }

  // Apply alterations
  async dalSchemaAlterationApply(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { alterationID } = (a as KV) || {}
    if (!alterationID) {
      throw Error('field alterationID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.dalSchemaAlterationApplyEndpoint(),
    }
    cfg.params = {
      alterationID,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dalSchemaAlterationApplyCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dalSchemaAlterationApply(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dalSchemaAlterationApplyEndpoint(): string {
    return '/dal/schema/alterations/apply'
  }

  // Dismiss alterations
  async dalSchemaAlterationDismiss(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { alterationID } = (a as KV) || {}
    if (!alterationID) {
      throw Error('field alterationID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.dalSchemaAlterationDismissEndpoint(),
    }
    cfg.params = {
      alterationID,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dalSchemaAlterationDismissCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dalSchemaAlterationDismiss(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dalSchemaAlterationDismissEndpoint(): string {
    return '/dal/schema/alterations/dismiss'
  }

  // Search connections (Directory)
  async dalConnectionList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { connectionID, handle, type, deleted, incTotal } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.dalConnectionListEndpoint(),
    }
    cfg.params = {
      connectionID,
      handle,
      type,
      deleted,
      incTotal,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dalConnectionListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dalConnectionList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dalConnectionListEndpoint(): string {
    return '/dal/connections/'
  }

  // Create connection
  async dalConnectionCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { handle, type, meta, config } = (a as KV) || {}
    if (!type) {
      throw Error('field type is empty')
    }
    if (!meta) {
      throw Error('field meta is empty')
    }
    if (!config) {
      throw Error('field config is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.dalConnectionCreateEndpoint(),
    }
    cfg.data = {
      handle,
      type,
      meta,
      config,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dalConnectionCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dalConnectionCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dalConnectionCreateEndpoint(): string {
    return '/dal/connections/'
  }

  // Update connection details
  async dalConnectionUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { connectionID, handle, type, meta, config, updatedAt } = (a as KV) || {}
    if (!connectionID) {
      throw Error('field connectionID is empty')
    }
    if (!type) {
      throw Error('field type is empty')
    }
    if (!meta) {
      throw Error('field meta is empty')
    }
    if (!config) {
      throw Error('field config is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.dalConnectionUpdateEndpoint({
        connectionID,
      }),
    }
    cfg.data = {
      handle,
      type,
      meta,
      config,
      updatedAt,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dalConnectionUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dalConnectionUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dalConnectionUpdateEndpoint(a: KV): string {
    const { connectionID } = a || {}
    return `/dal/connections/${connectionID}`
  }

  // Read connection details
  async dalConnectionRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { connectionID } = (a as KV) || {}
    if (!connectionID) {
      throw Error('field connectionID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.dalConnectionReadEndpoint({
        connectionID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dalConnectionReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dalConnectionRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dalConnectionReadEndpoint(a: KV): string {
    const { connectionID } = a || {}
    return `/dal/connections/${connectionID}`
  }

  // Remove connection
  async dalConnectionDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { connectionID } = (a as KV) || {}
    if (!connectionID) {
      throw Error('field connectionID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.dalConnectionDeleteEndpoint({
        connectionID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dalConnectionDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dalConnectionDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dalConnectionDeleteEndpoint(a: KV): string {
    const { connectionID } = a || {}
    return `/dal/connections/${connectionID}`
  }

  // Undelete connection
  async dalConnectionUndelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { connectionID } = (a as KV) || {}
    if (!connectionID) {
      throw Error('field connectionID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.dalConnectionUndeleteEndpoint({
        connectionID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dalConnectionUndeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dalConnectionUndelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dalConnectionUndeleteEndpoint(a: KV): string {
    const { connectionID } = a || {}
    return `/dal/connections/${connectionID}/undelete`
  }

  // List applications
  async applicationList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { name, query, deleted, labels, flags, incFlags, limit, incTotal, pageCursor, sort } =
      (a as KV) || {}
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
      incTotal,
      pageCursor,
      sort,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  applicationListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.applicationList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  applicationListEndpoint(): string {
    return '/application/'
  }

  // Create application
  async applicationCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { name, enabled, weight, unify, config, labels } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  applicationCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.applicationCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  applicationCreateEndpoint(): string {
    return '/application/'
  }

  // Update user details
  async applicationUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { applicationID, name, enabled, weight, unify, config, labels, updatedAt } =
      (a as KV) || {}
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
      updatedAt,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  applicationUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.applicationUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  applicationUpdateEndpoint(a: KV): string {
    const { applicationID } = a || {}
    return `/application/${applicationID}`
  }

  // Upload application assets
  async applicationUpload(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { upload } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  applicationUploadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.applicationUpload(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  applicationUploadEndpoint(): string {
    return '/application/upload'
  }

  // Flag application
  async applicationFlagCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { applicationID, flag, ownedBy } = (a as KV) || {}
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
        applicationID,
        flag,
        ownedBy,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  applicationFlagCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.applicationFlagCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  applicationFlagCreateEndpoint(a: KV): string {
    const { applicationID, flag, ownedBy } = a || {}
    return `/application/${applicationID}/flag/${ownedBy}/${flag}`
  }

  // Unflag application
  async applicationFlagDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { applicationID, flag, ownedBy } = (a as KV) || {}
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
        applicationID,
        flag,
        ownedBy,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  applicationFlagDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.applicationFlagDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  applicationFlagDeleteEndpoint(a: KV): string {
    const { applicationID, flag, ownedBy } = a || {}
    return `/application/${applicationID}/flag/${ownedBy}/${flag}`
  }

  // Read application details
  async applicationRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { applicationID, incFlags } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  applicationReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.applicationRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  applicationReadEndpoint(a: KV): string {
    const { applicationID } = a || {}
    return `/application/${applicationID}`
  }

  // Remove application
  async applicationDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { applicationID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  applicationDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.applicationDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  applicationDeleteEndpoint(a: KV): string {
    const { applicationID } = a || {}
    return `/application/${applicationID}`
  }

  // Undelete application
  async applicationUndelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { applicationID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  applicationUndeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.applicationUndelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  applicationUndeleteEndpoint(a: KV): string {
    const { applicationID } = a || {}
    return `/application/${applicationID}/undelete`
  }

  // Fire system:application trigger
  async applicationTriggerScript(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { applicationID, script, args } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  applicationTriggerScriptCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.applicationTriggerScript(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  applicationTriggerScriptEndpoint(a: KV): string {
    const { applicationID } = a || {}
    return `/application/${applicationID}/trigger`
  }

  // Reorder applications
  async applicationReorder(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { applicationIDs } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  applicationReorderCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.applicationReorder(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  applicationReorderEndpoint(): string {
    return '/application/reorder'
  }

  // List labels
  async labelList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { kind, name, value, limit } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.labelListEndpoint(),
    }
    cfg.params = {
      kind,
      name,
      value,
      limit,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  labelListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.labelList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  labelListEndpoint(): string {
    return '/label/'
  }

  // Delete label
  async labelDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { name, kind } = (a as KV) || {}
    if (!name) {
      throw Error('field name is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.labelDeleteEndpoint({
        name,
      }),
    }
    cfg.params = {
      kind,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  labelDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.labelDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  labelDeleteEndpoint(a: KV): string {
    const { name } = a || {}
    return `/label/${name}`
  }

  // Retrieve defined permissions
  async permissionsList(extra: AxiosRequestConfig = {}): Promise<KV> {
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.permissionsListEndpoint(),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  permissionsListCancellable(extra: AxiosRequestConfig = {}): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.permissionsList(options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  permissionsListEndpoint(): string {
    return '/permissions/'
  }

  // Effective rules for current user
  async permissionsEffective(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { resource } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.permissionsEffectiveEndpoint(),
    }
    cfg.params = {
      resource,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  permissionsEffectiveCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.permissionsEffective(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  permissionsEffectiveEndpoint(): string {
    return '/permissions/effective'
  }

  // Evaluate rules for given user/role combo
  async permissionsTrace(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { resource, userID, roleID } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.permissionsTraceEndpoint(),
    }
    cfg.params = {
      resource,
      userID,
      roleID,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  permissionsTraceCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.permissionsTrace(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  permissionsTraceEndpoint(): string {
    return '/permissions/trace'
  }

  // Retrieve role permissions
  async permissionsRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { roleID, resource } = (a as KV) || {}
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
    cfg.params = {
      resource,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  permissionsReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.permissionsRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  permissionsReadEndpoint(a: KV): string {
    const { roleID } = a || {}
    return `/permissions/${roleID}/rules`
  }

  // Remove all defined role permissions
  async permissionsDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { roleID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  permissionsDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.permissionsDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  permissionsDeleteEndpoint(a: KV): string {
    const { roleID } = a || {}
    return `/permissions/${roleID}/rules`
  }

  // Update permission settings
  async permissionsUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { roleID, rules } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  permissionsUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.permissionsUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  permissionsUpdateEndpoint(a: KV): string {
    const { roleID } = a || {}
    return `/permissions/${roleID}/rules`
  }

  // List/read reminders
  async reminderList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  reminderListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.reminderList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  reminderListEndpoint(): string {
    return '/reminder/'
  }

  // Add new reminder
  async reminderCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { resource, assignedTo, payload, remindAt } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  reminderCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.reminderCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  reminderCreateEndpoint(): string {
    return '/reminder/'
  }

  // Update reminder
  async reminderUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { reminderID, resource, assignedTo, payload, remindAt } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  reminderUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.reminderUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  reminderUpdateEndpoint(a: KV): string {
    const { reminderID } = a || {}
    return `/reminder/${reminderID}`
  }

  // Read reminder by ID
  async reminderRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { reminderID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  reminderReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.reminderRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  reminderReadEndpoint(a: KV): string {
    const { reminderID } = a || {}
    return `/reminder/${reminderID}`
  }

  // Delete reminder
  async reminderDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { reminderID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  reminderDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.reminderDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  reminderDeleteEndpoint(a: KV): string {
    const { reminderID } = a || {}
    return `/reminder/${reminderID}`
  }

  // Dismiss reminder
  async reminderDismiss(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { reminderID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  reminderDismissCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.reminderDismiss(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  reminderDismissEndpoint(a: KV): string {
    const { reminderID } = a || {}
    return `/reminder/${reminderID}/dismiss`
  }

  // Undismiss reminder
  async reminderUndismiss(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { reminderID } = (a as KV) || {}
    if (!reminderID) {
      throw Error('field reminderID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'patch',
      url: this.reminderUndismissEndpoint({
        reminderID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  reminderUndismissCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.reminderUndismiss(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  reminderUndismissEndpoint(a: KV): string {
    const { reminderID } = a || {}
    return `/reminder/${reminderID}/undismiss`
  }

  // Snooze reminder
  async reminderSnooze(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { reminderID, remindAt } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  reminderSnoozeCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.reminderSnooze(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  reminderSnoozeEndpoint(a: KV): string {
    const { reminderID } = a || {}
    return `/reminder/${reminderID}/snooze`
  }

  // List/read notifications
  async notificationList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { notificationID, kind, read, deleted, limit, pageCursor, sort } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.notificationListEndpoint(),
    }
    cfg.params = {
      notificationID,
      kind,
      read,
      deleted,
      limit,
      pageCursor,
      sort,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  notificationListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.notificationList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  notificationListEndpoint(): string {
    return '/notification/'
  }

  // Add new notification
  async notificationCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { kind, config, recipient } = (a as KV) || {}
    if (!kind) {
      throw Error('field kind is empty')
    }
    if (!config) {
      throw Error('field config is empty')
    }
    if (!recipient) {
      throw Error('field recipient is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.notificationCreateEndpoint(),
    }
    cfg.data = {
      kind,
      config,
      recipient,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  notificationCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.notificationCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  notificationCreateEndpoint(): string {
    return '/notification/'
  }

  // Update notification
  async notificationUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { notificationID, kind, config, recipient } = (a as KV) || {}
    if (!notificationID) {
      throw Error('field notificationID is empty')
    }
    if (!kind) {
      throw Error('field kind is empty')
    }
    if (!config) {
      throw Error('field config is empty')
    }
    if (!recipient) {
      throw Error('field recipient is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.notificationUpdateEndpoint({
        notificationID,
      }),
    }
    cfg.data = {
      kind,
      config,
      recipient,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  notificationUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.notificationUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  notificationUpdateEndpoint(a: KV): string {
    const { notificationID } = a || {}
    return `/notification/${notificationID}`
  }

  // Read notification by ID
  async notificationRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { notificationID } = (a as KV) || {}
    if (!notificationID) {
      throw Error('field notificationID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.notificationReadEndpoint({
        notificationID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  notificationReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.notificationRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  notificationReadEndpoint(a: KV): string {
    const { notificationID } = a || {}
    return `/notification/${notificationID}`
  }

  // Delete notification
  async notificationDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { notificationID } = (a as KV) || {}
    if (!notificationID) {
      throw Error('field notificationID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.notificationDeleteEndpoint({
        notificationID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  notificationDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.notificationDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  notificationDeleteEndpoint(a: KV): string {
    const { notificationID } = a || {}
    return `/notification/${notificationID}`
  }

  // Mark notification as read
  async notificationMarkAsRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { notificationID } = (a as KV) || {}
    if (!notificationID) {
      throw Error('field notificationID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'patch',
      url: this.notificationMarkAsReadEndpoint({
        notificationID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  notificationMarkAsReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.notificationMarkAsRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  notificationMarkAsReadEndpoint(a: KV): string {
    const { notificationID } = a || {}
    return `/notification/${notificationID}/read`
  }

  // Mark notification as unread
  async notificationMarkAsUnread(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { notificationID } = (a as KV) || {}
    if (!notificationID) {
      throw Error('field notificationID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'patch',
      url: this.notificationMarkAsUnreadEndpoint({
        notificationID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  notificationMarkAsUnreadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.notificationMarkAsUnread(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  notificationMarkAsUnreadEndpoint(a: KV): string {
    const { notificationID } = a || {}
    return `/notification/${notificationID}/unread`
  }

  // Mark all notifications as read for current user
  async notificationMarkAllAsRead(extra: AxiosRequestConfig = {}): Promise<KV> {
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'patch',
      url: this.notificationMarkAllAsReadEndpoint(),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  notificationMarkAllAsReadCancellable(extra: AxiosRequestConfig = {}): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.notificationMarkAllAsRead(options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  notificationMarkAllAsReadEndpoint(): string {
    return '/notification/all/read'
  }

  // Mark all notifications as unread for current user
  async notificationMarkAllAsUnread(extra: AxiosRequestConfig = {}): Promise<KV> {
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'patch',
      url: this.notificationMarkAllAsUnreadEndpoint(),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  notificationMarkAllAsUnreadCancellable(extra: AxiosRequestConfig = {}): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.notificationMarkAllAsUnread(options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  notificationMarkAllAsUnreadEndpoint(): string {
    return '/notification/all/unread'
  }

  // Attachment details
  async attachmentRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { kind, attachmentID, sign, userID } = (a as KV) || {}
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
        kind,
        attachmentID,
      }),
    }
    cfg.params = {
      sign,
      userID,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  attachmentReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.attachmentRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  attachmentReadEndpoint(a: KV): string {
    const { kind, attachmentID } = a || {}
    return `/attachment/${kind}/${attachmentID}`
  }

  // Delete attachment
  async attachmentDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { kind, attachmentID, sign, userID } = (a as KV) || {}
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
        kind,
        attachmentID,
      }),
    }
    cfg.params = {
      sign,
      userID,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  attachmentDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.attachmentDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  attachmentDeleteEndpoint(a: KV): string {
    const { kind, attachmentID } = a || {}
    return `/attachment/${kind}/${attachmentID}`
  }

  // Serves attached file
  async attachmentOriginal(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { kind, attachmentID, name, sign, userID, download } = (a as KV) || {}
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
        kind,
        attachmentID,
        name,
      }),
    }
    cfg.params = {
      sign,
      userID,
      download,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  attachmentOriginalCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.attachmentOriginal(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  attachmentOriginalEndpoint(a: KV): string {
    const { kind, attachmentID, name } = a || {}
    return `/attachment/${kind}/${attachmentID}/original/${name}`
  }

  // Serves preview of an attached file
  async attachmentPreview(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { kind, attachmentID, ext, sign, userID } = (a as KV) || {}
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
        kind,
        attachmentID,
        ext,
      }),
    }
    cfg.params = {
      sign,
      userID,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  attachmentPreviewCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.attachmentPreview(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  attachmentPreviewEndpoint(a: KV): string {
    const { kind, attachmentID, ext } = a || {}
    return `/attachment/${kind}/${attachmentID}/preview.${ext}`
  }

  // List templates
  async templateList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      query,
      handle,
      type,
      ownerID,
      partial,
      deleted,
      labels,
      limit,
      incTotal,
      pageCursor,
      sort,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.templateListEndpoint(),
    }
    cfg.params = {
      query,
      handle,
      type,
      ownerID,
      partial,
      deleted,
      labels,
      limit,
      incTotal,
      pageCursor,
      sort,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  templateListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.templateList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  templateListEndpoint(): string {
    return '/template/'
  }

  // Create template
  async templateCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { handle, language, type, partial, meta, template, ownerID, labels } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  templateCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.templateCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  templateCreateEndpoint(): string {
    return '/template/'
  }

  // Read template
  async templateRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { templateID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  templateReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.templateRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  templateReadEndpoint(a: KV): string {
    const { templateID } = a || {}
    return `/template/${templateID}`
  }

  // Update template
  async templateUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
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
      updatedAt,
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
      updatedAt,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  templateUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.templateUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  templateUpdateEndpoint(a: KV): string {
    const { templateID } = a || {}
    return `/template/${templateID}`
  }

  // Delete template
  async templateDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { templateID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  templateDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.templateDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  templateDeleteEndpoint(a: KV): string {
    const { templateID } = a || {}
    return `/template/${templateID}`
  }

  // Undelete template
  async templateUndelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { templateID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  templateUndeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.templateUndelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  templateUndeleteEndpoint(a: KV): string {
    const { templateID } = a || {}
    return `/template/${templateID}/undelete`
  }

  // Render drivers
  async templateRenderDrivers(extra: AxiosRequestConfig = {}): Promise<KV> {
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.templateRenderDriversEndpoint(),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  templateRenderDriversCancellable(extra: AxiosRequestConfig = {}): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.templateRenderDrivers(options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  templateRenderDriversEndpoint(): string {
    return '/template/render/drivers'
  }

  // Render template
  async templateRender(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { templateID, filename, ext, variables, options } = (a as KV) || {}
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
        templateID,
        filename,
        ext,
      }),
    }
    cfg.data = {
      variables,
      options,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  templateRenderCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.templateRender(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  templateRenderEndpoint(a: KV): string {
    const { templateID, filename, ext } = a || {}
    return `/template/${templateID}/render/${filename}.${ext}`
  }

  // List reports
  async reportList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { handle, query, deleted, labels, limit, incTotal, pageCursor, sort } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.reportListEndpoint(),
    }
    cfg.params = {
      handle,
      query,
      deleted,
      labels,
      limit,
      incTotal,
      pageCursor,
      sort,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  reportListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.reportList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  reportListEndpoint(): string {
    return '/reports/'
  }

  // Create report
  async reportCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { handle, meta, scenarios, sources, blocks, labels } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  reportCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.reportCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  reportCreateEndpoint(): string {
    return '/reports/'
  }

  // Update report
  async reportUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { reportID, handle, meta, scenarios, sources, blocks, labels, updatedAt } =
      (a as KV) || {}
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
      updatedAt,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  reportUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.reportUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  reportUpdateEndpoint(a: KV): string {
    const { reportID } = a || {}
    return `/reports/${reportID}`
  }

  // Read report details
  async reportRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { reportID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  reportReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.reportRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  reportReadEndpoint(a: KV): string {
    const { reportID } = a || {}
    return `/reports/${reportID}`
  }

  // Remove report
  async reportDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { reportID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  reportDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.reportDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  reportDeleteEndpoint(a: KV): string {
    const { reportID } = a || {}
    return `/reports/${reportID}`
  }

  // Undelete report
  async reportUndelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { reportID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  reportUndeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.reportUndelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  reportUndeleteEndpoint(a: KV): string {
    const { reportID } = a || {}
    return `/reports/${reportID}/undelete`
  }

  // Describe report
  async reportDescribe(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { sources, steps, describe } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  reportDescribeCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.reportDescribe(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  reportDescribeEndpoint(): string {
    return '/reports/describe'
  }

  // Run report
  async reportRun(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { reportID, frames } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  reportRunCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.reportRun(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  reportRunEndpoint(a: KV): string {
    const { reportID } = a || {}
    return `/reports/${reportID}/run`
  }

  // List system statistics
  async statsList(extra: AxiosRequestConfig = {}): Promise<KV> {
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.statsListEndpoint(),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  statsListCancellable(extra: AxiosRequestConfig = {}): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.statsList(options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  statsListEndpoint(): string {
    return '/stats/'
  }

  // List all available automation scripts for system resources
  async automationList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  automationListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.automationList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  automationListEndpoint(): string {
    return '/automation/'
  }

  // Serves client scripts bundle
  async automationBundle(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { bundle, type, ext } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.automationBundleEndpoint({
        bundle,
        type,
        ext,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  automationBundleCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.automationBundle(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  automationBundleEndpoint(a: KV): string {
    const { bundle, type, ext } = a || {}
    return `/automation/${bundle}-${type}.${ext}`
  }

  // Triggers execution of a specific script on a system service level
  async automationTriggerScript(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { script, args } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  automationTriggerScriptCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.automationTriggerScript(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  automationTriggerScriptEndpoint(): string {
    return '/automation/trigger'
  }

  // Action log events
  async actionlogList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { from, to, beforeActionID, resource, action, actorID, origin, limit } = (a as KV) || {}
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
      origin,
      limit,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  actionlogListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.actionlogList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  actionlogListEndpoint(): string {
    return '/actionlog/'
  }

  // Messaging queues
  async queuesList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { query, limit, incTotal, pageCursor, sort, deleted } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.queuesListEndpoint(),
    }
    cfg.params = {
      query,
      limit,
      incTotal,
      pageCursor,
      sort,
      deleted,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  queuesListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.queuesList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  queuesListEndpoint(): string {
    return '/queues/'
  }

  // Create messaging queue
  async queuesCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { queue, consumer, meta } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  queuesCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.queuesCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  queuesCreateEndpoint(): string {
    return '/queues'
  }

  // Messaging queue details
  async queuesRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { queueID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  queuesReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.queuesRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  queuesReadEndpoint(a: KV): string {
    const { queueID } = a || {}
    return `/queues/${queueID}`
  }

  // Update queue details
  async queuesUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { queueID, queue, consumer, meta, updatedAt } = (a as KV) || {}
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
      updatedAt,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  queuesUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.queuesUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  queuesUpdateEndpoint(a: KV): string {
    const { queueID } = a || {}
    return `/queues/${queueID}`
  }

  // Messaging queue delete
  async queuesDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { queueID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  queuesDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.queuesDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  queuesDeleteEndpoint(a: KV): string {
    const { queueID } = a || {}
    return `/queues/${queueID}`
  }

  // Messaging queue undelete
  async queuesUndelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { queueID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  queuesUndeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.queuesUndelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  queuesUndeleteEndpoint(a: KV): string {
    const { queueID } = a || {}
    return `/queues/${queueID}/undelete`
  }

  // List routes
  async apigwRouteList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { routeID, query, deleted, disabled, labels, limit, incTotal, pageCursor, sort } =
      (a as KV) || {}
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
      incTotal,
      pageCursor,
      sort,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  apigwRouteListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.apigwRouteList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  apigwRouteListEndpoint(): string {
    return '/apigw/route/'
  }

  // Create route
  async apigwRouteCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { endpoint, method, enabled, group, meta } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  apigwRouteCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.apigwRouteCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  apigwRouteCreateEndpoint(): string {
    return '/apigw/route'
  }

  // Update route details
  async apigwRouteUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { routeID, endpoint, method, enabled, group, meta, updatedAt } = (a as KV) || {}
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
      updatedAt,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  apigwRouteUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.apigwRouteUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  apigwRouteUpdateEndpoint(a: KV): string {
    const { routeID } = a || {}
    return `/apigw/route/${routeID}`
  }

  // Read route details
  async apigwRouteRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { routeID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  apigwRouteReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.apigwRouteRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  apigwRouteReadEndpoint(a: KV): string {
    const { routeID } = a || {}
    return `/apigw/route/${routeID}`
  }

  // Remove route
  async apigwRouteDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { routeID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  apigwRouteDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.apigwRouteDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  apigwRouteDeleteEndpoint(a: KV): string {
    const { routeID } = a || {}
    return `/apigw/route/${routeID}`
  }

  // Undelete route
  async apigwRouteUndelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { routeID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  apigwRouteUndeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.apigwRouteUndelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  apigwRouteUndeleteEndpoint(a: KV): string {
    const { routeID } = a || {}
    return `/apigw/route/${routeID}/undelete`
  }

  // List filters
  async apigwFilterList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { routeID, deleted, disabled, limit, pageCursor, sort } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  apigwFilterListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.apigwFilterList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  apigwFilterListEndpoint(): string {
    return '/apigw/filter/'
  }

  // Create filter
  async apigwFilterCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { routeID, weight, kind, ref, enabled, params } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  apigwFilterCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.apigwFilterCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  apigwFilterCreateEndpoint(): string {
    return '/apigw/filter'
  }

  // Update filter details
  async apigwFilterUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { filterID, routeID, weight, kind, ref, enabled, params, updatedAt } = (a as KV) || {}
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
      updatedAt,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  apigwFilterUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.apigwFilterUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  apigwFilterUpdateEndpoint(a: KV): string {
    const { filterID } = a || {}
    return `/apigw/filter/${filterID}`
  }

  // Read filter details
  async apigwFilterRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { filterID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  apigwFilterReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.apigwFilterRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  apigwFilterReadEndpoint(a: KV): string {
    const { filterID } = a || {}
    return `/apigw/filter/${filterID}`
  }

  // Remove filter
  async apigwFilterDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { filterID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  apigwFilterDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.apigwFilterDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  apigwFilterDeleteEndpoint(a: KV): string {
    const { filterID } = a || {}
    return `/apigw/filter/${filterID}`
  }

  // Undelete filter
  async apigwFilterUndelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { filterID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  apigwFilterUndeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.apigwFilterUndelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  apigwFilterUndeleteEndpoint(a: KV): string {
    const { filterID } = a || {}
    return `/apigw/filter/${filterID}/undelete`
  }

  // Filter definitions
  async apigwFilterDefFilter(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { kind } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.apigwFilterDefFilterEndpoint(),
    }
    cfg.params = {
      kind,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  apigwFilterDefFilterCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.apigwFilterDefFilter(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  apigwFilterDefFilterEndpoint(): string {
    return '/apigw/filter/def'
  }

  // Proxy auth definitions
  async apigwFilterDefProxyAuth(extra: AxiosRequestConfig = {}): Promise<KV> {
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.apigwFilterDefProxyAuthEndpoint(),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  apigwFilterDefProxyAuthCancellable(extra: AxiosRequestConfig = {}): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.apigwFilterDefProxyAuth(options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  apigwFilterDefProxyAuthEndpoint(): string {
    return '/apigw/filter/proxy_auth/def'
  }

  // List aggregated list of routes
  async apigwProfilerAggregation(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { path, before, sort, limit } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  apigwProfilerAggregationCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.apigwProfilerAggregation(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  apigwProfilerAggregationEndpoint(): string {
    return '/apigw/profiler/'
  }

  // List hits per route
  async apigwProfilerRoute(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { routeID, path, before, sort, limit } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  apigwProfilerRouteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.apigwProfilerRoute(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  apigwProfilerRouteEndpoint(a: KV): string {
    const { routeID } = a || {}
    return `/apigw/profiler/route/${routeID}`
  }

  // Hit details
  async apigwProfilerHit(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { hitID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  apigwProfilerHitCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.apigwProfilerHit(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  apigwProfilerHitEndpoint(a: KV): string {
    const { hitID } = a || {}
    return `/apigw/profiler/hit/${hitID}`
  }

  // Purge all profiler hits
  async apigwProfilerPurgeAll(extra: AxiosRequestConfig = {}): Promise<KV> {
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.apigwProfilerPurgeAllEndpoint(),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  apigwProfilerPurgeAllCancellable(extra: AxiosRequestConfig = {}): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.apigwProfilerPurgeAll(options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  apigwProfilerPurgeAllEndpoint(): string {
    return '/apigw/profiler/purge'
  }

  // Purge route profiler hits
  async apigwProfilerPurge(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { routeID } = (a as KV) || {}
    if (!routeID) {
      throw Error('field routeID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.apigwProfilerPurgeEndpoint({
        routeID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  apigwProfilerPurgeCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.apigwProfilerPurge(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  apigwProfilerPurgeEndpoint(a: KV): string {
    const { routeID } = a || {}
    return `/apigw/profiler/purge/${routeID}`
  }

  // List resources translations
  async localeListResource(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { lang, resource, resourceType, ownerID, deleted, limit, pageCursor, sort } =
      (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  localeListResourceCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.localeListResource(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  localeListResourceEndpoint(): string {
    return '/locale/resource'
  }

  // Create resource translation
  async localeCreateResource(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { lang, resource, key, place, message, ownerID } = (a as KV) || {}
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
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  localeCreateResourceCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.localeCreateResource(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  localeCreateResourceEndpoint(): string {
    return '/locale/resource'
  }

  // Update resource translation
  async localeUpdateResource(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { translationID, lang, resource, key, place, message, ownerID, updatedAt } =
      (a as KV) || {}
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
      updatedAt,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  localeUpdateResourceCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.localeUpdateResource(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  localeUpdateResourceEndpoint(a: KV): string {
    const { translationID } = a || {}
    return `/locale/resource/${translationID}`
  }

  // Read resource translation details
  async localeReadResource(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { translationID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  localeReadResourceCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.localeReadResource(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  localeReadResourceEndpoint(a: KV): string {
    const { translationID } = a || {}
    return `/locale/resource/${translationID}`
  }

  // Remove resource translation
  async localeDeleteResource(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { translationID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  localeDeleteResourceCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.localeDeleteResource(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  localeDeleteResourceEndpoint(a: KV): string {
    const { translationID } = a || {}
    return `/locale/resource/${translationID}`
  }

  // Undelete resource translation
  async localeUndeleteResource(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { translationID } = (a as KV) || {}
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

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  localeUndeleteResourceCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.localeUndeleteResource(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  localeUndeleteResourceEndpoint(a: KV): string {
    const { translationID } = a || {}
    return `/locale/resource/${translationID}/undelete`
  }

  // List all available languages
  async localeList(extra: AxiosRequestConfig = {}): Promise<KV> {
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.localeListEndpoint(),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  localeListCancellable(extra: AxiosRequestConfig = {}): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.localeList(options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  localeListEndpoint(): string {
    return '/locale/'
  }

  // List all available translation in a language for a specific webapp
  async localeGet(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { lang, application } = (a as KV) || {}
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
        lang,
        application,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  localeGetCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.localeGet(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  localeGetEndpoint(a: KV): string {
    const { lang, application } = a || {}
    return `/locale/${lang}/${application}`
  }

  // List connections for data privacy
  async dataPrivacyConnectionList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { connectionID, handle, type, deleted } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.dataPrivacyConnectionListEndpoint(),
    }
    cfg.params = {
      connectionID,
      handle,
      type,
      deleted,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dataPrivacyConnectionListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dataPrivacyConnectionList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dataPrivacyConnectionListEndpoint(): string {
    return '/data-privacy/connection/'
  }

  // List data privacy requests
  async dataPrivacyRequestList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { requestedBy, query, kind, status, limit, pageCursor, sort } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.dataPrivacyRequestListEndpoint(),
    }
    cfg.params = {
      requestedBy,
      query,
      kind,
      status,
      limit,
      pageCursor,
      sort,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dataPrivacyRequestListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dataPrivacyRequestList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dataPrivacyRequestListEndpoint(): string {
    return '/data-privacy/requests/'
  }

  // Create data privacy request
  async dataPrivacyRequestCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { kind, payload } = (a as KV) || {}
    if (!kind) {
      throw Error('field kind is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.dataPrivacyRequestCreateEndpoint(),
    }
    cfg.data = {
      kind,
      payload,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dataPrivacyRequestCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dataPrivacyRequestCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dataPrivacyRequestCreateEndpoint(): string {
    return '/data-privacy/requests/'
  }

  // Get details about specific request
  async dataPrivacyRequestRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { requestID } = (a as KV) || {}
    if (!requestID) {
      throw Error('field requestID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.dataPrivacyRequestReadEndpoint({
        requestID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dataPrivacyRequestReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dataPrivacyRequestRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dataPrivacyRequestReadEndpoint(a: KV): string {
    const { requestID } = a || {}
    return `/data-privacy/requests/${requestID}`
  }

  // Update data privacy request status
  async dataPrivacyRequestUpdateStatus(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { requestID, status } = (a as KV) || {}
    if (!requestID) {
      throw Error('field requestID is empty')
    }
    if (!status) {
      throw Error('field status is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'patch',
      url: this.dataPrivacyRequestUpdateStatusEndpoint({
        requestID,
        status,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dataPrivacyRequestUpdateStatusCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dataPrivacyRequestUpdateStatus(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dataPrivacyRequestUpdateStatusEndpoint(a: KV): string {
    const { requestID, status } = a || {}
    return `/data-privacy/requests/${requestID}/status/${status}`
  }

  // List data privacy request comments
  async dataPrivacyRequestCommentList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { requestID, limit, pageCursor, sort } = (a as KV) || {}
    if (!requestID) {
      throw Error('field requestID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.dataPrivacyRequestCommentListEndpoint({
        requestID,
      }),
    }
    cfg.params = {
      limit,
      pageCursor,
      sort,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dataPrivacyRequestCommentListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dataPrivacyRequestCommentList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dataPrivacyRequestCommentListEndpoint(a: KV): string {
    const { requestID } = a || {}
    return `/data-privacy/requests/${requestID}/comments/`
  }

  // Create data privacy request comment
  async dataPrivacyRequestCommentCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { requestID, comment } = (a as KV) || {}
    if (!requestID) {
      throw Error('field requestID is empty')
    }
    if (!comment) {
      throw Error('field comment is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.dataPrivacyRequestCommentCreateEndpoint({
        requestID,
      }),
    }
    cfg.data = {
      comment,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  dataPrivacyRequestCommentCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.dataPrivacyRequestCommentCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  dataPrivacyRequestCommentCreateEndpoint(a: KV): string {
    const { requestID } = a || {}
    return `/data-privacy/requests/${requestID}/comments/`
  }

  // List connections
  async connectionList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      handle,
      status,
      query,
      tags,
      source,
      deleted,
      labels,
      limit,
      incTotal,
      pageCursor,
      sort,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.connectionListEndpoint(),
    }
    cfg.params = {
      handle,
      status,
      query,
      tags,
      source,
      deleted,
      labels,
      limit,
      incTotal,
      pageCursor,
      sort,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  connectionListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.connectionList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  connectionListEndpoint(): string {
    return '/connections/'
  }

  // Create connection
  async connectionCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { handle, meta, service, resources, standardOperations, operations, labels } =
      (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.connectionCreateEndpoint(),
    }
    cfg.data = {
      handle,
      meta,
      service,
      resources,
      standardOperations,
      operations,
      labels,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  connectionCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.connectionCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  connectionCreateEndpoint(): string {
    return '/connections/'
  }

  // Update connection
  async connectionUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      connectionID,
      handle,
      meta,
      service,
      resources,
      standardOperations,
      operations,
      labels,
      updatedAt,
    } = (a as KV) || {}
    if (!connectionID) {
      throw Error('field connectionID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.connectionUpdateEndpoint({
        connectionID,
      }),
    }
    cfg.data = {
      handle,
      meta,
      service,
      resources,
      standardOperations,
      operations,
      labels,
      updatedAt,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  connectionUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.connectionUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  connectionUpdateEndpoint(a: KV): string {
    const { connectionID } = a || {}
    return `/connections/${connectionID}`
  }

  // Read connection
  async connectionRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { connectionID } = (a as KV) || {}
    if (!connectionID) {
      throw Error('field connectionID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.connectionReadEndpoint({
        connectionID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  connectionReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.connectionRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  connectionReadEndpoint(a: KV): string {
    const { connectionID } = a || {}
    return `/connections/${connectionID}`
  }

  // Delete connection
  async connectionDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { connectionID } = (a as KV) || {}
    if (!connectionID) {
      throw Error('field connectionID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.connectionDeleteEndpoint({
        connectionID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  connectionDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.connectionDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  connectionDeleteEndpoint(a: KV): string {
    const { connectionID } = a || {}
    return `/connections/${connectionID}`
  }

  // Undelete connection
  async connectionUndelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { connectionID } = (a as KV) || {}
    if (!connectionID) {
      throw Error('field connectionID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.connectionUndeleteEndpoint({
        connectionID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  connectionUndeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.connectionUndelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  connectionUndeleteEndpoint(a: KV): string {
    const { connectionID } = a || {}
    return `/connections/${connectionID}/undelete`
  }

  // Enable connection (lock draft local connection to active)
  async connectionEnable(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { connectionID } = (a as KV) || {}
    if (!connectionID) {
      throw Error('field connectionID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.connectionEnableEndpoint({
        connectionID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  connectionEnableCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.connectionEnable(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  connectionEnableEndpoint(a: KV): string {
    const { connectionID } = a || {}
    return `/connections/${connectionID}/enable`
  }

  // Generate connection via builder agent
  async connectionGenerate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { prompt, context } = (a as KV) || {}
    if (!prompt) {
      throw Error('field prompt is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.connectionGenerateEndpoint(),
    }
    cfg.data = {
      prompt,
      context,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  connectionGenerateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.connectionGenerate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  connectionGenerateEndpoint(): string {
    return '/connections/generate'
  }

  // Import a catalog connection
  async connectionImport(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { catalogID } = (a as KV) || {}
    if (!catalogID) {
      throw Error('field catalogID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.connectionImportEndpoint(),
    }
    cfg.data = {
      catalogID,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  connectionImportCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.connectionImport(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  connectionImportEndpoint(): string {
    return '/connections/import'
  }

  // Configure a connection (create configured connection)
  async connectionConfigure(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { connectionID, catalogID, name, config, labels } = (a as KV) || {}
    if (!connectionID) {
      throw Error('field connectionID is empty')
    }
    if (!name) {
      throw Error('field name is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.connectionConfigureEndpoint({
        connectionID,
      }),
    }
    cfg.data = {
      catalogID,
      name,
      config,
      labels,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  connectionConfigureCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.connectionConfigure(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  connectionConfigureEndpoint(a: KV): string {
    const { connectionID } = a || {}
    return `/connections/${connectionID}/configure`
  }

  // Update connection configuration
  async connectionUpdateConfiguration(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { connectionID, configuredConnectionID, name, config, labels } = (a as KV) || {}
    if (!connectionID) {
      throw Error('field connectionID is empty')
    }
    if (!configuredConnectionID) {
      throw Error('field configuredConnectionID is empty')
    }
    if (!name) {
      throw Error('field name is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'patch',
      url: this.connectionUpdateConfigurationEndpoint({
        connectionID,
        configuredConnectionID,
      }),
    }
    cfg.data = {
      name,
      config,
      labels,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  connectionUpdateConfigurationCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.connectionUpdateConfiguration(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  connectionUpdateConfigurationEndpoint(a: KV): string {
    const { connectionID, configuredConnectionID } = a || {}
    return `/connections/${connectionID}/configure/${configuredConnectionID}`
  }

  // List configured connections
  async configuredConnectionList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { connectionID, status, query, deleted, labels, limit, incTotal, pageCursor, sort } =
      (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.configuredConnectionListEndpoint(),
    }
    cfg.params = {
      connectionID,
      status,
      query,
      deleted,
      labels,
      limit,
      incTotal,
      pageCursor,
      sort,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  configuredConnectionListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.configuredConnectionList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  configuredConnectionListEndpoint(): string {
    return '/configured-connections/'
  }

  // Read configured connection
  async configuredConnectionRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { connectionID } = (a as KV) || {}
    if (!connectionID) {
      throw Error('field connectionID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.configuredConnectionReadEndpoint({
        connectionID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  configuredConnectionReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.configuredConnectionRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  configuredConnectionReadEndpoint(a: KV): string {
    const { connectionID } = a || {}
    return `/configured-connections/${connectionID}`
  }

  // Delete configured connection
  async configuredConnectionDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { connectionID } = (a as KV) || {}
    if (!connectionID) {
      throw Error('field connectionID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.configuredConnectionDeleteEndpoint({
        connectionID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  configuredConnectionDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.configuredConnectionDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  configuredConnectionDeleteEndpoint(a: KV): string {
    const { connectionID } = a || {}
    return `/configured-connections/${connectionID}`
  }

  // Enable configured connection (validate, provision, lock to active)
  async configuredConnectionEnable(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { connectionID } = (a as KV) || {}
    if (!connectionID) {
      throw Error('field connectionID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.configuredConnectionEnableEndpoint({
        connectionID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  configuredConnectionEnableCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.configuredConnectionEnable(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  configuredConnectionEnableEndpoint(a: KV): string {
    const { connectionID } = a || {}
    return `/configured-connections/${connectionID}/enable`
  }

  // Check configured connection health (connectivity, auth, probe)
  async configuredConnectionCheck(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { connectionID } = (a as KV) || {}
    if (!connectionID) {
      throw Error('field connectionID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.configuredConnectionCheckEndpoint({
        connectionID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  configuredConnectionCheckCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.configuredConnectionCheck(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  configuredConnectionCheckEndpoint(a: KV): string {
    const { connectionID } = a || {}
    return `/configured-connections/${connectionID}/check`
  }

  // List available MCP tools
  async mcpListTools(extra: AxiosRequestConfig = {}): Promise<KV> {
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.mcpListToolsEndpoint(),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  mcpListToolsCancellable(extra: AxiosRequestConfig = {}): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.mcpListTools(options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  mcpListToolsEndpoint(): string {
    return '/mcp/tools'
  }

  // Check SMTP server configuration settings
  async smtpConfigurationCheckerCheck(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { host, port, recipients, username, password, tlsInsecure, tlsServerName } =
      (a as KV) || {}
    if (!host) {
      throw Error('field host is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.smtpConfigurationCheckerCheckEndpoint(),
    }
    cfg.data = {
      host,
      port,
      recipients,
      username,
      password,
      tlsInsecure,
      tlsServerName,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  smtpConfigurationCheckerCheckCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.smtpConfigurationCheckerCheck(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  smtpConfigurationCheckerCheckEndpoint(): string {
    return '/smtp/configuration-checker/'
  }

  // Get TCL master list
  async agentTclMasterList(extra: AxiosRequestConfig = {}): Promise<KV> {
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.agentTclMasterListEndpoint(),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  agentTclMasterListCancellable(extra: AxiosRequestConfig = {}): {
    response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>
    cancel: () => void
  } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.agentTclMasterList(options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  agentTclMasterListEndpoint(): string {
    return '/agents/tcl'
  }

  // List agents
  async agentList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { query, handle, status, deleted, labels, limit, incTotal, pageCursor, sort } =
      (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.agentListEndpoint(),
    }
    cfg.params = {
      query,
      handle,
      status,
      deleted,
      labels,
      limit,
      incTotal,
      pageCursor,
      sort,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  agentListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.agentList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  agentListEndpoint(): string {
    return '/agents/'
  }

  // Create agent
  async agentCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { handle, status, meta, behavior, execution, access, invocation, labels } =
      (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.agentCreateEndpoint(),
    }
    cfg.data = {
      handle,
      status,
      meta,
      behavior,
      execution,
      access,
      invocation,
      labels,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  agentCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.agentCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  agentCreateEndpoint(): string {
    return '/agents'
  }

  // Read agent details
  async agentRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { agentID } = (a as KV) || {}
    if (!agentID) {
      throw Error('field agentID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.agentReadEndpoint({
        agentID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  agentReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.agentRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  agentReadEndpoint(a: KV): string {
    const { agentID } = a || {}
    return `/agents/${agentID}`
  }

  // Update agent details
  async agentUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      agentID,
      handle,
      status,
      meta,
      behavior,
      execution,
      access,
      invocation,
      labels,
      updatedAt,
    } = (a as KV) || {}
    if (!agentID) {
      throw Error('field agentID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.agentUpdateEndpoint({
        agentID,
      }),
    }
    cfg.data = {
      handle,
      status,
      meta,
      behavior,
      execution,
      access,
      invocation,
      labels,
      updatedAt,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  agentUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.agentUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  agentUpdateEndpoint(a: KV): string {
    const { agentID } = a || {}
    return `/agents/${agentID}`
  }

  // Delete agent
  async agentDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { agentID } = (a as KV) || {}
    if (!agentID) {
      throw Error('field agentID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.agentDeleteEndpoint({
        agentID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  agentDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.agentDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  agentDeleteEndpoint(a: KV): string {
    const { agentID } = a || {}
    return `/agents/${agentID}`
  }

  // Undelete agent
  async agentUndelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { agentID } = (a as KV) || {}
    if (!agentID) {
      throw Error('field agentID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.agentUndeleteEndpoint({
        agentID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  agentUndeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.agentUndelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  agentUndeleteEndpoint(a: KV): string {
    const { agentID } = a || {}
    return `/agents/${agentID}/undelete`
  }

  // Execute agent
  async agentExec(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { agentID, input, conversationID, context } = (a as KV) || {}
    if (!agentID) {
      throw Error('field agentID is empty')
    }
    if (!input) {
      throw Error('field input is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.agentExecEndpoint({
        agentID,
      }),
    }
    cfg.data = {
      input,
      conversationID,
      context,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  agentExecCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.agentExec(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  agentExecEndpoint(a: KV): string {
    const { agentID } = a || {}
    return `/agents/${agentID}/exec`
  }

  // List LLM providers
  async llmProviderList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { provider, status } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.llmProviderListEndpoint(),
    }
    cfg.params = {
      provider,
      status,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  llmProviderListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.llmProviderList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  llmProviderListEndpoint(): string {
    return '/llm-providers/'
  }

  // Create LLM provider
  async llmProviderCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { handle, provider, status, apiKey, meta, config } = (a as KV) || {}
    if (!provider) {
      throw Error('field provider is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.llmProviderCreateEndpoint(),
    }
    cfg.data = {
      handle,
      provider,
      status,
      apiKey,
      meta,
      config,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  llmProviderCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.llmProviderCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  llmProviderCreateEndpoint(): string {
    return '/llm-providers/'
  }

  // Read LLM provider
  async llmProviderRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { llmProviderID } = (a as KV) || {}
    if (!llmProviderID) {
      throw Error('field llmProviderID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.llmProviderReadEndpoint({
        llmProviderID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  llmProviderReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.llmProviderRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  llmProviderReadEndpoint(a: KV): string {
    const { llmProviderID } = a || {}
    return `/llm-providers/${llmProviderID}`
  }

  // Update LLM provider
  async llmProviderUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { llmProviderID, handle, provider, status, apiKey, meta, config } = (a as KV) || {}
    if (!llmProviderID) {
      throw Error('field llmProviderID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.llmProviderUpdateEndpoint({
        llmProviderID,
      }),
    }
    cfg.data = {
      handle,
      provider,
      status,
      apiKey,
      meta,
      config,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  llmProviderUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.llmProviderUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  llmProviderUpdateEndpoint(a: KV): string {
    const { llmProviderID } = a || {}
    return `/llm-providers/${llmProviderID}`
  }

  // Delete LLM provider
  async llmProviderDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { llmProviderID } = (a as KV) || {}
    if (!llmProviderID) {
      throw Error('field llmProviderID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.llmProviderDeleteEndpoint({
        llmProviderID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  llmProviderDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.llmProviderDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  llmProviderDeleteEndpoint(a: KV): string {
    const { llmProviderID } = a || {}
    return `/llm-providers/${llmProviderID}`
  }

  // List available models for an LLM provider
  async llmProviderModels(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { llmProviderID } = (a as KV) || {}
    if (!llmProviderID) {
      throw Error('field llmProviderID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.llmProviderModelsEndpoint({
        llmProviderID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  llmProviderModelsCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.llmProviderModels(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  llmProviderModelsEndpoint(a: KV): string {
    const { llmProviderID } = a || {}
    return `/llm-providers/${llmProviderID}/models`
  }

  // Validate LLM provider connection
  async llmProviderValidate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { llmProviderID } = (a as KV) || {}
    if (!llmProviderID) {
      throw Error('field llmProviderID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.llmProviderValidateEndpoint({
        llmProviderID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  llmProviderValidateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.llmProviderValidate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  llmProviderValidateEndpoint(a: KV): string {
    const { llmProviderID } = a || {}
    return `/llm-providers/${llmProviderID}/validate`
  }

  // List knowledge bases
  async knowledgeBaseList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { query, handle, deleted, limit, incTotal, pageCursor, sort } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.knowledgeBaseListEndpoint(),
    }
    cfg.params = {
      query,
      handle,
      deleted,
      limit,
      incTotal,
      pageCursor,
      sort,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  knowledgeBaseListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.knowledgeBaseList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  knowledgeBaseListEndpoint(): string {
    return '/knowledge-bases/'
  }

  // Create knowledge base
  async knowledgeBaseCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { handle, title, description, context } = (a as KV) || {}
    if (!title) {
      throw Error('field title is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.knowledgeBaseCreateEndpoint(),
    }
    cfg.data = {
      handle,
      title,
      description,
      context,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  knowledgeBaseCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.knowledgeBaseCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  knowledgeBaseCreateEndpoint(): string {
    return '/knowledge-bases'
  }

  // Read knowledge base
  async knowledgeBaseRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { knowledgeBaseID } = (a as KV) || {}
    if (!knowledgeBaseID) {
      throw Error('field knowledgeBaseID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.knowledgeBaseReadEndpoint({
        knowledgeBaseID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  knowledgeBaseReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.knowledgeBaseRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  knowledgeBaseReadEndpoint(a: KV): string {
    const { knowledgeBaseID } = a || {}
    return `/knowledge-bases/${knowledgeBaseID}`
  }

  // Update knowledge base
  async knowledgeBaseUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { knowledgeBaseID, handle, title, description, context, updatedAt } = (a as KV) || {}
    if (!knowledgeBaseID) {
      throw Error('field knowledgeBaseID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.knowledgeBaseUpdateEndpoint({
        knowledgeBaseID,
      }),
    }
    cfg.data = {
      handle,
      title,
      description,
      context,
      updatedAt,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  knowledgeBaseUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.knowledgeBaseUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  knowledgeBaseUpdateEndpoint(a: KV): string {
    const { knowledgeBaseID } = a || {}
    return `/knowledge-bases/${knowledgeBaseID}`
  }

  // Delete knowledge base
  async knowledgeBaseDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { knowledgeBaseID } = (a as KV) || {}
    if (!knowledgeBaseID) {
      throw Error('field knowledgeBaseID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.knowledgeBaseDeleteEndpoint({
        knowledgeBaseID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  knowledgeBaseDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.knowledgeBaseDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  knowledgeBaseDeleteEndpoint(a: KV): string {
    const { knowledgeBaseID } = a || {}
    return `/knowledge-bases/${knowledgeBaseID}`
  }

  // Undelete knowledge base
  async knowledgeBaseUndelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { knowledgeBaseID } = (a as KV) || {}
    if (!knowledgeBaseID) {
      throw Error('field knowledgeBaseID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.knowledgeBaseUndeleteEndpoint({
        knowledgeBaseID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  knowledgeBaseUndeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.knowledgeBaseUndelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  knowledgeBaseUndeleteEndpoint(a: KV): string {
    const { knowledgeBaseID } = a || {}
    return `/knowledge-bases/${knowledgeBaseID}/undelete`
  }

  // List AI conversations
  async aiConversationList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { agentID, deleted, limit, incTotal, pageCursor, sort } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.aiConversationListEndpoint(),
    }
    cfg.params = {
      agentID,
      deleted,
      limit,
      incTotal,
      pageCursor,
      sort,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  aiConversationListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.aiConversationList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  aiConversationListEndpoint(): string {
    return '/ai-conversations/'
  }

  // Read AI conversation details
  async aiConversationRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { aiConversationID } = (a as KV) || {}
    if (!aiConversationID) {
      throw Error('field aiConversationID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.aiConversationReadEndpoint({
        aiConversationID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  aiConversationReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.aiConversationRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  aiConversationReadEndpoint(a: KV): string {
    const { aiConversationID } = a || {}
    return `/ai-conversations/${aiConversationID}`
  }

  // Delete AI conversation
  async aiConversationDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { aiConversationID } = (a as KV) || {}
    if (!aiConversationID) {
      throw Error('field aiConversationID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.aiConversationDeleteEndpoint({
        aiConversationID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  aiConversationDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.aiConversationDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  aiConversationDeleteEndpoint(a: KV): string {
    const { aiConversationID } = a || {}
    return `/ai-conversations/${aiConversationID}`
  }

  // Undelete AI conversation
  async aiConversationUndelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { aiConversationID } = (a as KV) || {}
    if (!aiConversationID) {
      throw Error('field aiConversationID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.aiConversationUndeleteEndpoint({
        aiConversationID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  aiConversationUndeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.aiConversationUndelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  aiConversationUndeleteEndpoint(a: KV): string {
    const { aiConversationID } = a || {}
    return `/ai-conversations/${aiConversationID}/undelete`
  }

  // Continue AI conversation
  async aiConversationContinue(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { aiConversationID, input } = (a as KV) || {}
    if (!aiConversationID) {
      throw Error('field aiConversationID is empty')
    }
    if (!input) {
      throw Error('field input is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.aiConversationContinueEndpoint({
        aiConversationID,
      }),
    }
    cfg.data = {
      input,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  aiConversationContinueCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.aiConversationContinue(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  aiConversationContinueEndpoint(a: KV): string {
    const { aiConversationID } = a || {}
    return `/ai-conversations/${aiConversationID}/continue`
  }

  // List chatbots
  async chatbotList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { query, handle, deleted, labels, limit, incTotal, pageCursor, sort } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.chatbotListEndpoint(),
    }
    cfg.params = {
      query,
      handle,
      deleted,
      labels,
      limit,
      incTotal,
      pageCursor,
      sort,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  chatbotListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.chatbotList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  chatbotListEndpoint(): string {
    return '/chatbots/'
  }

  // Create chatbot
  async chatbotCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      handle,
      name,
      enabled,
      sessionTTL,
      allowedOrigins,
      handoff,
      styling,
      scenarios,
      labels,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.chatbotCreateEndpoint(),
    }
    cfg.data = {
      handle,
      name,
      enabled,
      sessionTTL,
      allowedOrigins,
      handoff,
      styling,
      scenarios,
      labels,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  chatbotCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.chatbotCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  chatbotCreateEndpoint(): string {
    return '/chatbots'
  }

  // Read chatbot details
  async chatbotRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { chatbotID } = (a as KV) || {}
    if (!chatbotID) {
      throw Error('field chatbotID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.chatbotReadEndpoint({
        chatbotID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  chatbotReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.chatbotRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  chatbotReadEndpoint(a: KV): string {
    const { chatbotID } = a || {}
    return `/chatbots/${chatbotID}`
  }

  // Update chatbot details
  async chatbotUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      chatbotID,
      handle,
      name,
      enabled,
      sessionTTL,
      allowedOrigins,
      handoff,
      styling,
      scenarios,
      labels,
      updatedAt,
    } = (a as KV) || {}
    if (!chatbotID) {
      throw Error('field chatbotID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.chatbotUpdateEndpoint({
        chatbotID,
      }),
    }
    cfg.data = {
      handle,
      name,
      enabled,
      sessionTTL,
      allowedOrigins,
      handoff,
      styling,
      scenarios,
      labels,
      updatedAt,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  chatbotUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.chatbotUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  chatbotUpdateEndpoint(a: KV): string {
    const { chatbotID } = a || {}
    return `/chatbots/${chatbotID}`
  }

  // Delete chatbot
  async chatbotDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { chatbotID } = (a as KV) || {}
    if (!chatbotID) {
      throw Error('field chatbotID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.chatbotDeleteEndpoint({
        chatbotID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  chatbotDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.chatbotDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  chatbotDeleteEndpoint(a: KV): string {
    const { chatbotID } = a || {}
    return `/chatbots/${chatbotID}`
  }

  // Undelete chatbot
  async chatbotUndelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { chatbotID } = (a as KV) || {}
    if (!chatbotID) {
      throw Error('field chatbotID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.chatbotUndeleteEndpoint({
        chatbotID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  chatbotUndeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.chatbotUndelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  chatbotUndeleteEndpoint(a: KV): string {
    const { chatbotID } = a || {}
    return `/chatbots/${chatbotID}/undelete`
  }

  // Regenerate chatbot widget key
  async chatbotRegenerateWidgetKey(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { chatbotID } = (a as KV) || {}
    if (!chatbotID) {
      throw Error('field chatbotID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.chatbotRegenerateWidgetKeyEndpoint({
        chatbotID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  chatbotRegenerateWidgetKeyCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.chatbotRegenerateWidgetKey(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  chatbotRegenerateWidgetKeyEndpoint(a: KV): string {
    const { chatbotID } = a || {}
    return `/chatbots/${chatbotID}/regenerate-key`
  }

  // Upload chatbot asset
  async chatbotUploadAsset(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { chatbotID, upload } = (a as KV) || {}
    if (!chatbotID) {
      throw Error('field chatbotID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.chatbotUploadAssetEndpoint({
        chatbotID,
      }),
    }
    cfg.data = {
      upload,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  chatbotUploadAssetCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.chatbotUploadAsset(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  chatbotUploadAssetEndpoint(a: KV): string {
    const { chatbotID } = a || {}
    return `/chatbots/${chatbotID}/upload-asset`
  }

  // List sessions (widget + preview, aggregated)
  async chatbotSessionList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { chatbotID, status, source, limit, incTotal, pageCursor, sort } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.chatbotSessionListEndpoint(),
    }
    cfg.params = {
      chatbotID,
      status,
      source,
      limit,
      incTotal,
      pageCursor,
      sort,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  chatbotSessionListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.chatbotSessionList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  chatbotSessionListEndpoint(): string {
    return '/chatbots/sessions/'
  }

  // Read session with steps
  async chatbotSessionRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { sessionID } = (a as KV) || {}
    if (!sessionID) {
      throw Error('field sessionID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.chatbotSessionReadEndpoint({
        sessionID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  chatbotSessionReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.chatbotSessionRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  chatbotSessionReadEndpoint(a: KV): string {
    const { sessionID } = a || {}
    return `/chatbots/sessions/${sessionID}`
  }

  // Force-advance session step (admin)
  async chatbotSessionAdvanceStep(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { sessionID } = (a as KV) || {}
    if (!sessionID) {
      throw Error('field sessionID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.chatbotSessionAdvanceStepEndpoint({
        sessionID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  chatbotSessionAdvanceStepCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.chatbotSessionAdvanceStep(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  chatbotSessionAdvanceStepEndpoint(a: KV): string {
    const { sessionID } = a || {}
    return `/chatbots/sessions/${sessionID}/advance-step`
  }

  // Close session (admin)
  async chatbotSessionClose(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { sessionID } = (a as KV) || {}
    if (!sessionID) {
      throw Error('field sessionID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.chatbotSessionCloseEndpoint({
        sessionID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  chatbotSessionCloseCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.chatbotSessionClose(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  chatbotSessionCloseEndpoint(a: KV): string {
    const { sessionID } = a || {}
    return `/chatbots/sessions/${sessionID}/close`
  }

  // Accept pending handoff
  async chatbotSessionHandoffAccept(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { sessionID, operator } = (a as KV) || {}
    if (!sessionID) {
      throw Error('field sessionID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.chatbotSessionHandoffAcceptEndpoint({
        sessionID,
      }),
    }
    cfg.data = {
      operator,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  chatbotSessionHandoffAcceptCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.chatbotSessionHandoffAccept(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  chatbotSessionHandoffAcceptEndpoint(a: KV): string {
    const { sessionID } = a || {}
    return `/chatbots/sessions/${sessionID}/handoff-accept`
  }

  // Send operator message
  async chatbotSessionOperatorMessage(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { sessionID, message, operator } = (a as KV) || {}
    if (!sessionID) {
      throw Error('field sessionID is empty')
    }
    if (!message) {
      throw Error('field message is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.chatbotSessionOperatorMessageEndpoint({
        sessionID,
      }),
    }
    cfg.data = {
      message,
      operator,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  chatbotSessionOperatorMessageCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.chatbotSessionOperatorMessage(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  chatbotSessionOperatorMessageEndpoint(a: KV): string {
    const { sessionID } = a || {}
    return `/chatbots/sessions/${sessionID}/operator-message`
  }

  // Complete handoff
  async chatbotSessionHandoffComplete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { sessionID } = (a as KV) || {}
    if (!sessionID) {
      throw Error('field sessionID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.chatbotSessionHandoffCompleteEndpoint({
        sessionID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  chatbotSessionHandoffCompleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.chatbotSessionHandoffComplete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  chatbotSessionHandoffCompleteEndpoint(a: KV): string {
    const { sessionID } = a || {}
    return `/chatbots/sessions/${sessionID}/handoff-complete`
  }

  // SSE stream of session events
  async chatbotSessionStream(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { sessionID } = (a as KV) || {}
    if (!sessionID) {
      throw Error('field sessionID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.chatbotSessionStreamEndpoint({
        sessionID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  chatbotSessionStreamCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.chatbotSessionStream(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  chatbotSessionStreamEndpoint(a: KV): string {
    const { sessionID } = a || {}
    return `/chatbots/sessions/${sessionID}/stream`
  }

  // List projects
  async projectList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { query, handle, status, deleted, labels, limit, incTotal, pageCursor, sort } =
      (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.projectListEndpoint(),
    }
    cfg.params = {
      query,
      handle,
      status,
      deleted,
      labels,
      limit,
      incTotal,
      pageCursor,
      sort,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  projectListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.projectList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  projectListEndpoint(): string {
    return '/projects/'
  }

  // Create project
  async projectCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { handle, status, config, meta, labels } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.projectCreateEndpoint(),
    }
    cfg.data = {
      handle,
      status,
      config,
      meta,
      labels,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  projectCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.projectCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  projectCreateEndpoint(): string {
    return '/projects'
  }

  // Read project details
  async projectRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { projectID } = (a as KV) || {}
    if (!projectID) {
      throw Error('field projectID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.projectReadEndpoint({
        projectID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  projectReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.projectRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  projectReadEndpoint(a: KV): string {
    const { projectID } = a || {}
    return `/projects/${projectID}`
  }

  // Update project details
  async projectUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { projectID, handle, status, config, meta, updatedAt, labels } = (a as KV) || {}
    if (!projectID) {
      throw Error('field projectID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.projectUpdateEndpoint({
        projectID,
      }),
    }
    cfg.data = {
      handle,
      status,
      config,
      meta,
      updatedAt,
      labels,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  projectUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.projectUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  projectUpdateEndpoint(a: KV): string {
    const { projectID } = a || {}
    return `/projects/${projectID}`
  }

  // Remove project
  async projectDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { projectID } = (a as KV) || {}
    if (!projectID) {
      throw Error('field projectID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.projectDeleteEndpoint({
        projectID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  projectDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.projectDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  projectDeleteEndpoint(a: KV): string {
    const { projectID } = a || {}
    return `/projects/${projectID}`
  }

  // Undelete project
  async projectUndelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { projectID } = (a as KV) || {}
    if (!projectID) {
      throw Error('field projectID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.projectUndeleteEndpoint({
        projectID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  projectUndeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.projectUndelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  projectUndeleteEndpoint(a: KV): string {
    const { projectID } = a || {}
    return `/projects/${projectID}/undelete`
  }

  // List project members
  async projectListMembers(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { projectID } = (a as KV) || {}
    if (!projectID) {
      throw Error('field projectID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.projectListMembersEndpoint({
        projectID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  projectListMembersCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.projectListMembers(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  projectListMembersEndpoint(a: KV): string {
    const { projectID } = a || {}
    return `/projects/${projectID}/members`
  }

  // Add project member
  async projectAddMember(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { projectID, userID, rolePreset } = (a as KV) || {}
    if (!projectID) {
      throw Error('field projectID is empty')
    }
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.projectAddMemberEndpoint({
        projectID,
      }),
    }
    cfg.data = {
      userID,
      rolePreset,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  projectAddMemberCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.projectAddMember(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  projectAddMemberEndpoint(a: KV): string {
    const { projectID } = a || {}
    return `/projects/${projectID}/members`
  }

  // Update project member role
  async projectUpdateMember(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { projectID, userID, rolePreset } = (a as KV) || {}
    if (!projectID) {
      throw Error('field projectID is empty')
    }
    if (!userID) {
      throw Error('field userID is empty')
    }
    if (!rolePreset) {
      throw Error('field rolePreset is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.projectUpdateMemberEndpoint({
        projectID,
        userID,
      }),
    }
    cfg.data = {
      rolePreset,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  projectUpdateMemberCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.projectUpdateMember(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  projectUpdateMemberEndpoint(a: KV): string {
    const { projectID, userID } = a || {}
    return `/projects/${projectID}/members/${userID}`
  }

  // Remove project member
  async projectRemoveMember(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { projectID, userID } = (a as KV) || {}
    if (!projectID) {
      throw Error('field projectID is empty')
    }
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.projectRemoveMemberEndpoint({
        projectID,
        userID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  projectRemoveMemberCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.projectRemoveMember(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  projectRemoveMemberEndpoint(a: KV): string {
    const { projectID, userID } = a || {}
    return `/projects/${projectID}/members/${userID}`
  }

  // Get project resource graph
  async projectGraph(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { projectID } = (a as KV) || {}
    if (!projectID) {
      throw Error('field projectID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.projectGraphEndpoint({
        projectID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  projectGraphCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.projectGraph(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  projectGraphEndpoint(a: KV): string {
    const { projectID } = a || {}
    return `/projects/${projectID}/graph`
  }

  // Save governance step form values
  async projectGovernanceSave(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { projectID, stepKey, values } = (a as KV) || {}
    if (!projectID) {
      throw Error('field projectID is empty')
    }
    if (!stepKey) {
      throw Error('field stepKey is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.projectGovernanceSaveEndpoint({
        projectID,
        stepKey,
      }),
    }
    cfg.data = {
      values,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  projectGovernanceSaveCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.projectGovernanceSave(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  projectGovernanceSaveEndpoint(a: KV): string {
    const { projectID, stepKey } = a || {}
    return `/projects/${projectID}/governance/${stepKey}`
  }

  // Transition governance step status
  async projectGovernanceTransition(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { projectID, stepKey, action, note } = (a as KV) || {}
    if (!projectID) {
      throw Error('field projectID is empty')
    }
    if (!stepKey) {
      throw Error('field stepKey is empty')
    }
    if (!action) {
      throw Error('field action is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.projectGovernanceTransitionEndpoint({
        projectID,
        stepKey,
      }),
    }
    cfg.data = {
      action,
      note,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  projectGovernanceTransitionCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.projectGovernanceTransition(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  projectGovernanceTransitionEndpoint(a: KV): string {
    const { projectID, stepKey } = a || {}
    return `/projects/${projectID}/governance/${stepKey}/transition`
  }

  // List project groups
  async projectGroupList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { projectID, query, handle, projectGroupID, deleted, limit, incTotal, pageCursor, sort } =
      (a as KV) || {}
    if (!projectID) {
      throw Error('field projectID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.projectGroupListEndpoint({
        projectID,
      }),
    }
    cfg.params = {
      query,
      handle,
      projectGroupID,
      deleted,
      limit,
      incTotal,
      pageCursor,
      sort,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  projectGroupListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.projectGroupList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  projectGroupListEndpoint(a: KV): string {
    const { projectID } = a || {}
    return `/projects/${projectID}/groups/`
  }

  // Create project group
  async projectGroupCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { projectID, handle, name, description } = (a as KV) || {}
    if (!projectID) {
      throw Error('field projectID is empty')
    }
    if (!handle) {
      throw Error('field handle is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.projectGroupCreateEndpoint({
        projectID,
      }),
    }
    cfg.data = {
      handle,
      name,
      description,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  projectGroupCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.projectGroupCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  projectGroupCreateEndpoint(a: KV): string {
    const { projectID } = a || {}
    return `/projects/${projectID}/groups/`
  }

  // Read project group
  async projectGroupRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { projectID, projectGroupID } = (a as KV) || {}
    if (!projectID) {
      throw Error('field projectID is empty')
    }
    if (!projectGroupID) {
      throw Error('field projectGroupID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.projectGroupReadEndpoint({
        projectID,
        projectGroupID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  projectGroupReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.projectGroupRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  projectGroupReadEndpoint(a: KV): string {
    const { projectID, projectGroupID } = a || {}
    return `/projects/${projectID}/groups/${projectGroupID}`
  }

  // Update project group
  async projectGroupUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { projectID, projectGroupID, handle, name, description, updatedAt } = (a as KV) || {}
    if (!projectID) {
      throw Error('field projectID is empty')
    }
    if (!projectGroupID) {
      throw Error('field projectGroupID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.projectGroupUpdateEndpoint({
        projectID,
        projectGroupID,
      }),
    }
    cfg.data = {
      handle,
      name,
      description,
      updatedAt,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  projectGroupUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.projectGroupUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  projectGroupUpdateEndpoint(a: KV): string {
    const { projectID, projectGroupID } = a || {}
    return `/projects/${projectID}/groups/${projectGroupID}`
  }

  // Delete project group
  async projectGroupDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { projectID, projectGroupID } = (a as KV) || {}
    if (!projectID) {
      throw Error('field projectID is empty')
    }
    if (!projectGroupID) {
      throw Error('field projectGroupID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.projectGroupDeleteEndpoint({
        projectID,
        projectGroupID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  projectGroupDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.projectGroupDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  projectGroupDeleteEndpoint(a: KV): string {
    const { projectID, projectGroupID } = a || {}
    return `/projects/${projectID}/groups/${projectGroupID}`
  }

  // Add resource to group
  async projectGroupEntryAdd(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { projectID, projectGroupID, resourceRef } = (a as KV) || {}
    if (!projectID) {
      throw Error('field projectID is empty')
    }
    if (!projectGroupID) {
      throw Error('field projectGroupID is empty')
    }
    if (!resourceRef) {
      throw Error('field resourceRef is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.projectGroupEntryAddEndpoint({
        projectID,
        projectGroupID,
        resourceRef,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  projectGroupEntryAddCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.projectGroupEntryAdd(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  projectGroupEntryAddEndpoint(a: KV): string {
    const { projectID, projectGroupID, resourceRef } = a || {}
    return `/projects/${projectID}/groups/${projectGroupID}/entries/${resourceRef}`
  }

  // Remove resource from group
  async projectGroupEntryRemove(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { projectID, projectGroupID, resourceRef } = (a as KV) || {}
    if (!projectID) {
      throw Error('field projectID is empty')
    }
    if (!projectGroupID) {
      throw Error('field projectGroupID is empty')
    }
    if (!resourceRef) {
      throw Error('field resourceRef is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.projectGroupEntryRemoveEndpoint({
        projectID,
        projectGroupID,
        resourceRef,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  projectGroupEntryRemoveCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.projectGroupEntryRemove(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  projectGroupEntryRemoveEndpoint(a: KV): string {
    const { projectID, projectGroupID, resourceRef } = a || {}
    return `/projects/${projectID}/groups/${projectGroupID}/entries/${resourceRef}`
  }

  // List tenants
  async tenantList(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { query, handle, status, deleted, labels, limit, incTotal, pageCursor, sort } =
      (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.tenantListEndpoint(),
    }
    cfg.params = {
      query,
      handle,
      status,
      deleted,
      labels,
      limit,
      incTotal,
      pageCursor,
      sort,
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  tenantListCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.tenantList(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  tenantListEndpoint(): string {
    return '/tenants/'
  }

  // Create tenant
  async tenantCreate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { handle, status, config, meta, labels } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.tenantCreateEndpoint(),
    }
    cfg.data = {
      handle,
      status,
      config,
      meta,
      labels,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  tenantCreateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.tenantCreate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  tenantCreateEndpoint(): string {
    return '/tenants'
  }

  // Read tenant details
  async tenantRead(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { tenantID } = (a as KV) || {}
    if (!tenantID) {
      throw Error('field tenantID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.tenantReadEndpoint({
        tenantID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  tenantReadCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.tenantRead(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  tenantReadEndpoint(a: KV): string {
    const { tenantID } = a || {}
    return `/tenants/${tenantID}`
  }

  // Update tenant details
  async tenantUpdate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { tenantID, handle, status, config, meta, updatedAt, labels } = (a as KV) || {}
    if (!tenantID) {
      throw Error('field tenantID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.tenantUpdateEndpoint({
        tenantID,
      }),
    }
    cfg.data = {
      handle,
      status,
      config,
      meta,
      updatedAt,
      labels,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  tenantUpdateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.tenantUpdate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  tenantUpdateEndpoint(a: KV): string {
    const { tenantID } = a || {}
    return `/tenants/${tenantID}`
  }

  // Remove tenant
  async tenantDelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { tenantID } = (a as KV) || {}
    if (!tenantID) {
      throw Error('field tenantID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.tenantDeleteEndpoint({
        tenantID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  tenantDeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.tenantDelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  tenantDeleteEndpoint(a: KV): string {
    const { tenantID } = a || {}
    return `/tenants/${tenantID}`
  }

  // Undelete tenant
  async tenantUndelete(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { tenantID } = (a as KV) || {}
    if (!tenantID) {
      throw Error('field tenantID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.tenantUndeleteEndpoint({
        tenantID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  tenantUndeleteCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.tenantUndelete(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  tenantUndeleteEndpoint(a: KV): string {
    const { tenantID } = a || {}
    return `/tenants/${tenantID}/undelete`
  }

  // Suspend tenant
  async tenantSuspend(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { tenantID } = (a as KV) || {}
    if (!tenantID) {
      throw Error('field tenantID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.tenantSuspendEndpoint({
        tenantID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  tenantSuspendCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.tenantSuspend(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  tenantSuspendEndpoint(a: KV): string {
    const { tenantID } = a || {}
    return `/tenants/${tenantID}/suspend`
  }

  // Activate tenant
  async tenantActivate(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { tenantID } = (a as KV) || {}
    if (!tenantID) {
      throw Error('field tenantID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.tenantActivateEndpoint({
        tenantID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  tenantActivateCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.tenantActivate(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  tenantActivateEndpoint(a: KV): string {
    const { tenantID } = a || {}
    return `/tenants/${tenantID}/activate`
  }

  // Archive tenant
  async tenantArchive(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { tenantID } = (a as KV) || {}
    if (!tenantID) {
      throw Error('field tenantID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.tenantArchiveEndpoint({
        tenantID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  tenantArchiveCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.tenantArchive(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  tenantArchiveEndpoint(a: KV): string {
    const { tenantID } = a || {}
    return `/tenants/${tenantID}/archive`
  }

  // List tenant members
  async tenantListMembers(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { tenantID } = (a as KV) || {}
    if (!tenantID) {
      throw Error('field tenantID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.tenantListMembersEndpoint({
        tenantID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  tenantListMembersCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.tenantListMembers(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  tenantListMembersEndpoint(a: KV): string {
    const { tenantID } = a || {}
    return `/tenants/${tenantID}/members`
  }

  // Invite tenant member
  async tenantAddMember(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { tenantID, userID, role } = (a as KV) || {}
    if (!tenantID) {
      throw Error('field tenantID is empty')
    }
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.tenantAddMemberEndpoint({
        tenantID,
      }),
    }
    cfg.data = {
      userID,
      role,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  tenantAddMemberCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.tenantAddMember(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  tenantAddMemberEndpoint(a: KV): string {
    const { tenantID } = a || {}
    return `/tenants/${tenantID}/members`
  }

  // Update tenant member
  async tenantUpdateMember(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { tenantID, userID, role, status } = (a as KV) || {}
    if (!tenantID) {
      throw Error('field tenantID is empty')
    }
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.tenantUpdateMemberEndpoint({
        tenantID,
        userID,
      }),
    }
    cfg.data = {
      role,
      status,
    }
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  tenantUpdateMemberCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.tenantUpdateMember(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  tenantUpdateMemberEndpoint(a: KV): string {
    const { tenantID, userID } = a || {}
    return `/tenants/${tenantID}/members/${userID}`
  }

  // Remove tenant member
  async tenantRemoveMember(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { tenantID, userID } = (a as KV) || {}
    if (!tenantID) {
      throw Error('field tenantID is empty')
    }
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.tenantRemoveMemberEndpoint({
        tenantID,
        userID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  tenantRemoveMemberCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.tenantRemoveMember(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  tenantRemoveMemberEndpoint(a: KV): string {
    const { tenantID, userID } = a || {}
    return `/tenants/${tenantID}/members/${userID}`
  }

  // Suspend tenant member
  async tenantSuspendMember(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { tenantID, userID } = (a as KV) || {}
    if (!tenantID) {
      throw Error('field tenantID is empty')
    }
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.tenantSuspendMemberEndpoint({
        tenantID,
        userID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  tenantSuspendMemberCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.tenantSuspendMember(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  tenantSuspendMemberEndpoint(a: KV): string {
    const { tenantID, userID } = a || {}
    return `/tenants/${tenantID}/members/${userID}/suspend`
  }

  // Activate tenant member
  async tenantActivateMember(a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const { tenantID, userID } = (a as KV) || {}
    if (!tenantID) {
      throw Error('field tenantID is empty')
    }
    if (!userID) {
      throw Error('field userID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.tenantActivateMemberEndpoint({
        tenantID,
        userID,
      }),
    }

    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  tenantActivateMemberCancellable(
    a: KV,
    extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.tenantActivateMember(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  tenantActivateMemberEndpoint(a: KV): string {
    const { tenantID, userID } = a || {}
    return `/tenants/${tenantID}/members/${userID}/activate`
  }
}
