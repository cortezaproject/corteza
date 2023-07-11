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
  accessTokenFn?: () => string | undefined;
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
      subWorkflow,
      labels,
      limit,
      incTotal,
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
      subWorkflow,
      labels,
      limit,
      incTotal,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  workflowListCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.workflowList(a, options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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

  workflowCreateCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.workflowCreate(a, options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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
      updatedAt,
    } = (a as KV) || {}
    if (!workflowID) {
      throw Error('field workflowID is empty')
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
      updatedAt,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  workflowUpdateCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.workflowUpdate(a, options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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

  workflowReadCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.workflowRead(a, options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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

  workflowDeleteCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.workflowDelete(a, options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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

  workflowUndeleteCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.workflowUndelete(a, options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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

  workflowTestCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.workflowTest(a, options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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

  workflowExecCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.workflowExec(a, options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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

  triggerListCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.triggerList(a, options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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

  triggerCreateCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.triggerCreate(a, options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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
      updatedAt,
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
      updatedAt,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  triggerUpdateCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.triggerUpdate(a, options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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

  triggerReadCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.triggerRead(a, options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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

  triggerDeleteCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.triggerDelete(a, options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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

  triggerUndeleteCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.triggerUndelete(a, options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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
      incTotal,
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
      incTotal,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  sessionListCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.sessionList(a, options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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

  sessionReadCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.sessionRead(a, options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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

  sessionCancelCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.sessionCancel(a, options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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

  sessionListPromptsCancellable (extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.sessionListPrompts(options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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

  sessionResumeStateCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.sessionResumeState(a, options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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

  functionListCancellable (extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.functionList(options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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

  typeListCancellable (extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.typeList(options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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

  eventTypesListCancellable (extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.eventTypesList(options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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

  permissionsListCancellable (extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.permissionsList(options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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

  permissionsEffectiveCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.permissionsEffective(a, options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
  }

  permissionsEffectiveEndpoint (): string {
    return '/permissions/effective'
  }

  // Evaluate rules for given user/role combo
  async permissionsTrace (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      resource,
      userID,
      roleID,
    } = (a as KV) || {}
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

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  permissionsTraceCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.permissionsTrace(a, options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
  }

  permissionsTraceEndpoint (): string {
    return '/permissions/trace'
  }

  // Retrieve role permissions
  async permissionsRead (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      roleID,
      resource,
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
    cfg.params = {
      resource,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  permissionsReadCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.permissionsRead(a, options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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

  permissionsDeleteCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.permissionsDelete(a, options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
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

  permissionsUpdateCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void; } {
    const cancelTokenSource = axios.CancelToken.source();
    let options = {...extra, cancelToken: cancelTokenSource.token }

    return {
        response: () => this.permissionsUpdate(a, options),
        cancel: () => {
          cancelTokenSource.cancel();
        }
    }
  }

  permissionsUpdateEndpoint (a: KV): string {
    const {
      roleID,
    } = a || {}
    return `/permissions/${roleID}/rules`
  }

}
