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

export default class Automation {
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

  setAccessTokenFn (fn: () => string | undefined): Automation {
    this.accessTokenFn = fn
    return this
  }

  setHeaders (headers?: Headers): Automation {
    if (typeof headers === 'object') {
      this.headers = headers
    }

    return this
  }

  setHeader (name: string, value: string | undefined): Automation {
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

  // List workflows
  async workflowList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      workflowID,
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
      url: this.workflowListEndpoint(),
    }
    cfg.params = {
      workflowID,
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

  workflowListEndpoint (): string {
    return '/workflows/'
  }

  // Create workflow
  async workflowCreate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      handle,
      labels,
      meta,
      enabled,
      trace,
      keepSessions,
      scope,
      steps,
      paths,
      runAs,
      ownedBy,
    } = (a as KV) || {}
    if (!handle) {
      throw Error('field handle is empty')
    }
    if (!runAs) {
      throw Error('field runAs is empty')
    }
    if (!ownedBy) {
      throw Error('field ownedBy is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.workflowCreateEndpoint(),
    }
    cfg.data = {
      handle,
      labels,
      meta,
      enabled,
      trace,
      keepSessions,
      scope,
      steps,
      paths,
      runAs,
      ownedBy,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  workflowCreateEndpoint (): string {
    return '/workflows/'
  }

  // Update triger details
  async workflowUpdate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      workflowID,
      handle,
      labels,
      meta,
      enabled,
      trace,
      keepSessions,
      scope,
      steps,
      paths,
      runAs,
      ownedBy,
    } = (a as KV) || {}
    if (!workflowID) {
      throw Error('field workflowID is empty')
    }
    if (!handle) {
      throw Error('field handle is empty')
    }
    if (!runAs) {
      throw Error('field runAs is empty')
    }
    if (!ownedBy) {
      throw Error('field ownedBy is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.workflowUpdateEndpoint({
        workflowID,
      }),
    }
    cfg.data = {
      handle,
      labels,
      meta,
      enabled,
      trace,
      keepSessions,
      scope,
      steps,
      paths,
      runAs,
      ownedBy,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  workflowUpdateEndpoint (a: KV): string {
    const {
      workflowID,
    } = a || {}
    return `/workflows/${workflowID}`
  }

  // Read workflow details
  async workflowRead (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      workflowID,
    } = (a as KV) || {}
    if (!workflowID) {
      throw Error('field workflowID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.workflowReadEndpoint({
        workflowID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  workflowReadEndpoint (a: KV): string {
    const {
      workflowID,
    } = a || {}
    return `/workflows/${workflowID}`
  }

  // Remove workflow
  async workflowDelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      workflowID,
    } = (a as KV) || {}
    if (!workflowID) {
      throw Error('field workflowID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.workflowDeleteEndpoint({
        workflowID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  workflowDeleteEndpoint (a: KV): string {
    const {
      workflowID,
    } = a || {}
    return `/workflows/${workflowID}`
  }

  // Undelete workflow
  async workflowUndelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      workflowID,
    } = (a as KV) || {}
    if (!workflowID) {
      throw Error('field workflowID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.workflowUndeleteEndpoint({
        workflowID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  workflowUndeleteEndpoint (a: KV): string {
    const {
      workflowID,
    } = a || {}
    return `/workflows/${workflowID}/undelete`
  }

  // Test workflow details
  async workflowTest (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      workflowID,
      scope,
      runAs,
    } = (a as KV) || {}
    if (!workflowID) {
      throw Error('field workflowID is empty')
    }
    if (!runAs) {
      throw Error('field runAs is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.workflowTestEndpoint({
        workflowID,
      }),
    }
    cfg.data = {
      scope,
      runAs,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  workflowTestEndpoint (a: KV): string {
    const {
      workflowID,
    } = a || {}
    return `/workflows/${workflowID}/test`
  }

  // Executes workflow on a specific step (must be orphan step and connected to &#x27;onManual&#x27; trigger)
  async workflowExec (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      workflowID,
      stepID,
      input,
      trace,
      wait,
      async,
    } = (a as KV) || {}
    if (!workflowID) {
      throw Error('field workflowID is empty')
    }
    if (!stepID) {
      throw Error('field stepID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.workflowExecEndpoint({
        workflowID,
      }),
    }
    cfg.data = {
      stepID,
      input,
      trace,
      wait,
      async,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  workflowExecEndpoint (a: KV): string {
    const {
      workflowID,
    } = a || {}
    return `/workflows/${workflowID}/exec`
  }

  // List triggers
  async triggerList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      triggerID,
      workflowID,
      deleted,
      disabled,
      eventType,
      resourceType,
      query,
      labels,
      limit,
      pageCursor,
      sort,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.triggerListEndpoint(),
    }
    cfg.params = {
      triggerID,
      workflowID,
      deleted,
      disabled,
      eventType,
      resourceType,
      query,
      labels,
      limit,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  triggerListEndpoint (): string {
    return '/triggers/'
  }

  // Create trigger
  async triggerCreate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      eventType,
      resourceType,
      enabled,
      workflowID,
      workflowStepID,
      input,
      labels,
      meta,
      constraints,
      ownedBy,
    } = (a as KV) || {}
    if (!eventType) {
      throw Error('field eventType is empty')
    }
    if (!resourceType) {
      throw Error('field resourceType is empty')
    }
    if (!workflowID) {
      throw Error('field workflowID is empty')
    }
    if (!workflowStepID) {
      throw Error('field workflowStepID is empty')
    }
    if (!ownedBy) {
      throw Error('field ownedBy is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.triggerCreateEndpoint(),
    }
    cfg.data = {
      eventType,
      resourceType,
      enabled,
      workflowID,
      workflowStepID,
      input,
      labels,
      meta,
      constraints,
      ownedBy,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  triggerCreateEndpoint (): string {
    return '/triggers/'
  }

  // Update trigger details
  async triggerUpdate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      triggerID,
      eventType,
      resourceType,
      enabled,
      workflowID,
      workflowStepID,
      input,
      labels,
      meta,
      constraints,
      ownedBy,
    } = (a as KV) || {}
    if (!triggerID) {
      throw Error('field triggerID is empty')
    }
    if (!eventType) {
      throw Error('field eventType is empty')
    }
    if (!resourceType) {
      throw Error('field resourceType is empty')
    }
    if (!workflowID) {
      throw Error('field workflowID is empty')
    }
    if (!workflowStepID) {
      throw Error('field workflowStepID is empty')
    }
    if (!ownedBy) {
      throw Error('field ownedBy is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.triggerUpdateEndpoint({
        triggerID,
      }),
    }
    cfg.data = {
      eventType,
      resourceType,
      enabled,
      workflowID,
      workflowStepID,
      input,
      labels,
      meta,
      constraints,
      ownedBy,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  triggerUpdateEndpoint (a: KV): string {
    const {
      triggerID,
    } = a || {}
    return `/triggers/${triggerID}`
  }

  // Read trigger details
  async triggerRead (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      triggerID,
    } = (a as KV) || {}
    if (!triggerID) {
      throw Error('field triggerID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.triggerReadEndpoint({
        triggerID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  triggerReadEndpoint (a: KV): string {
    const {
      triggerID,
    } = a || {}
    return `/triggers/${triggerID}`
  }

  // Remove trigger
  async triggerDelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      triggerID,
    } = (a as KV) || {}
    if (!triggerID) {
      throw Error('field triggerID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.triggerDeleteEndpoint({
        triggerID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  triggerDeleteEndpoint (a: KV): string {
    const {
      triggerID,
    } = a || {}
    return `/triggers/${triggerID}`
  }

  // Undelete trigger
  async triggerUndelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      triggerID,
    } = (a as KV) || {}
    if (!triggerID) {
      throw Error('field triggerID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.triggerUndeleteEndpoint({
        triggerID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  triggerUndeleteEndpoint (a: KV): string {
    const {
      triggerID,
    } = a || {}
    return `/triggers/${triggerID}/undelete`
  }

  // List sessions
  async sessionList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      sessionID,
      workflowID,
      createdBy,
      completed,
      status,
      eventType,
      resourceType,
      limit,
      pageCursor,
      sort,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.sessionListEndpoint(),
    }
    cfg.params = {
      sessionID,
      workflowID,
      createdBy,
      completed,
      status,
      eventType,
      resourceType,
      limit,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  sessionListEndpoint (): string {
    return '/sessions/'
  }

  // Read session details
  async sessionRead (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      sessionID,
    } = (a as KV) || {}
    if (!sessionID) {
      throw Error('field sessionID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.sessionReadEndpoint({
        sessionID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  sessionReadEndpoint (a: KV): string {
    const {
      sessionID,
    } = a || {}
    return `/sessions/${sessionID}`
  }

  // Cancel session
  async sessionCancel (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      sessionID,
    } = (a as KV) || {}
    if (!sessionID) {
      throw Error('field sessionID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.sessionCancelEndpoint({
        sessionID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  sessionCancelEndpoint (a: KV): string {
    const {
      sessionID,
    } = a || {}
    return `/sessions/${sessionID}/cancel`
  }

  // Returns pending prompts from all sessions
  async sessionListPrompts (extra: AxiosRequestConfig = {}): Promise<KV> {

    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.sessionListPromptsEndpoint(),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  sessionListPromptsEndpoint (): string {
    return '/sessions/prompts'
  }

  // Resume session
  async sessionResumeState (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      sessionID,
      stateID,
      input,
    } = (a as KV) || {}
    if (!sessionID) {
      throw Error('field sessionID is empty')
    }
    if (!stateID) {
      throw Error('field stateID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.sessionResumeStateEndpoint({
        sessionID, stateID,
      }),
    }
    cfg.data = {
      input,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  sessionResumeStateEndpoint (a: KV): string {
    const {
      sessionID,
      stateID,
    } = a || {}
    return `/sessions/${sessionID}/state/${stateID}`
  }

  // Available workflow functions
  async functionList (extra: AxiosRequestConfig = {}): Promise<KV> {

    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.functionListEndpoint(),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  functionListEndpoint (): string {
    return '/functions/'
  }

  // Available workflow types
  async typeList (extra: AxiosRequestConfig = {}): Promise<KV> {

    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.typeListEndpoint(),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  typeListEndpoint (): string {
    return '/types/'
  }

  // Available workflow types
  async eventTypesList (extra: AxiosRequestConfig = {}): Promise<KV> {

    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.eventTypesListEndpoint(),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  eventTypesListEndpoint (): string {
    return '/event-types/'
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

}
