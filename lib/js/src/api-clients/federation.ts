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

export default class Federation {
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

  setAccessTokenFn (fn: () => string | undefined): Federation {
    this.accessTokenFn = fn
    return this
  }

  setHeaders (headers?: Headers): Federation {
    if (typeof headers === 'object') {
      this.headers = headers
    }

    return this
  }

  setHeader (name: string, value: string | undefined): Federation {
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

  // Initialize the handshake step with node B
  async nodeHandshakeInitialize (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      nodeID,
      pairToken,
      sharedNodeID,
      authToken,
    } = (a as KV) || {}
    if (!nodeID) {
      throw Error('field nodeID is empty')
    }
    if (!pairToken) {
      throw Error('field pairToken is empty')
    }
    if (!sharedNodeID) {
      throw Error('field sharedNodeID is empty')
    }
    if (!authToken) {
      throw Error('field authToken is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.nodeHandshakeInitializeEndpoint({
        nodeID,
      }),
    }
    cfg.data = {
      pairToken,
      sharedNodeID,
      authToken,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  nodeHandshakeInitializeEndpoint (a: KV): string {
    const {
      nodeID,
    } = a || {}
    return `/nodes/${nodeID}/handshake`
  }

  // Search federated nodes
  async nodeSearch (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      query,
      status,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.nodeSearchEndpoint(),
    }
    cfg.params = {
      query,
      status,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  nodeSearchEndpoint (): string {
    return '/nodes/'
  }

  // Create a new federation node
  async nodeCreate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      baseURL,
      name,
      contact,
      pairingURI,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.nodeCreateEndpoint(),
    }
    cfg.data = {
      baseURL,
      name,
      contact,
      pairingURI,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  nodeCreateEndpoint (): string {
    return '/nodes/'
  }

  // Read a federation node
  async nodeRead (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      nodeID,
    } = (a as KV) || {}
    if (!nodeID) {
      throw Error('field nodeID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.nodeReadEndpoint({
        nodeID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  nodeReadEndpoint (a: KV): string {
    const {
      nodeID,
    } = a || {}
    return `/nodes/${nodeID}`
  }

  // Creates new sharable federation URI
  async nodeGenerateUri (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      nodeID,
    } = (a as KV) || {}
    if (!nodeID) {
      throw Error('field nodeID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.nodeGenerateUriEndpoint({
        nodeID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  nodeGenerateUriEndpoint (a: KV): string {
    const {
      nodeID,
    } = a || {}
    return `/nodes/${nodeID}/uri`
  }

  // Updates existing node
  async nodeUpdate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      nodeID,
      name,
      contact,
      baseURL,
    } = (a as KV) || {}
    if (!nodeID) {
      throw Error('field nodeID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.nodeUpdateEndpoint({
        nodeID,
      }),
    }
    cfg.data = {
      name,
      contact,
      baseURL,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  nodeUpdateEndpoint (a: KV): string {
    const {
      nodeID,
    } = a || {}
    return `/nodes/${nodeID}`
  }

  // Deletes node
  async nodeDelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      nodeID,
    } = (a as KV) || {}
    if (!nodeID) {
      throw Error('field nodeID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.nodeDeleteEndpoint({
        nodeID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  nodeDeleteEndpoint (a: KV): string {
    const {
      nodeID,
    } = a || {}
    return `/nodes/${nodeID}`
  }

  // Undeletes a node
  async nodeUndelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      nodeID,
    } = (a as KV) || {}
    if (!nodeID) {
      throw Error('field nodeID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.nodeUndeleteEndpoint({
        nodeID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  nodeUndeleteEndpoint (a: KV): string {
    const {
      nodeID,
    } = a || {}
    return `/nodes/${nodeID}/undelete`
  }

  // Initialize the pairing process between the two nodes
  async nodePair (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      nodeID,
    } = (a as KV) || {}
    if (!nodeID) {
      throw Error('field nodeID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.nodePairEndpoint({
        nodeID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  nodePairEndpoint (a: KV): string {
    const {
      nodeID,
    } = a || {}
    return `/nodes/${nodeID}/pair`
  }

  // Confirm the requested handshake
  async nodeHandshakeConfirm (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      nodeID,
    } = (a as KV) || {}
    if (!nodeID) {
      throw Error('field nodeID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.nodeHandshakeConfirmEndpoint({
        nodeID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  nodeHandshakeConfirmEndpoint (a: KV): string {
    const {
      nodeID,
    } = a || {}
    return `/nodes/${nodeID}/handshake-confirm`
  }

  // Complete the handshake
  async nodeHandshakeComplete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      nodeID,
      authToken,
    } = (a as KV) || {}
    if (!nodeID) {
      throw Error('field nodeID is empty')
    }
    if (!authToken) {
      throw Error('field authToken is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.nodeHandshakeCompleteEndpoint({
        nodeID,
      }),
    }
    cfg.data = {
      authToken,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  nodeHandshakeCompleteEndpoint (a: KV): string {
    const {
      nodeID,
    } = a || {}
    return `/nodes/${nodeID}/handshake-complete`
  }

  // Exposed settings for module
  async manageStructureReadExposed (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      nodeID,
      moduleID,
    } = (a as KV) || {}
    if (!nodeID) {
      throw Error('field nodeID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.manageStructureReadExposedEndpoint({
        nodeID, moduleID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  manageStructureReadExposedEndpoint (a: KV): string {
    const {
      nodeID,
      moduleID,
    } = a || {}
    return `/nodes/${nodeID}/modules/${moduleID}/exposed`
  }

  // Add module to federation
  async manageStructureCreateExposed (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      nodeID,
      composeModuleID,
      composeNamespaceID,
      name,
      handle,
      fields,
    } = (a as KV) || {}
    if (!nodeID) {
      throw Error('field nodeID is empty')
    }
    if (!composeModuleID) {
      throw Error('field composeModuleID is empty')
    }
    if (!composeNamespaceID) {
      throw Error('field composeNamespaceID is empty')
    }
    if (!name) {
      throw Error('field name is empty')
    }
    if (!handle) {
      throw Error('field handle is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.manageStructureCreateExposedEndpoint({
        nodeID,
      }),
    }
    cfg.data = {
      composeModuleID,
      composeNamespaceID,
      name,
      handle,
      fields,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  manageStructureCreateExposedEndpoint (a: KV): string {
    const {
      nodeID,
    } = a || {}
    return `/nodes/${nodeID}/modules/`
  }

  // Update already exposed module
  async manageStructureUpdateExposed (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      nodeID,
      moduleID,
      composeModuleID,
      composeNamespaceID,
      name,
      handle,
      fields,
    } = (a as KV) || {}
    if (!nodeID) {
      throw Error('field nodeID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    if (!composeModuleID) {
      throw Error('field composeModuleID is empty')
    }
    if (!composeNamespaceID) {
      throw Error('field composeNamespaceID is empty')
    }
    if (!name) {
      throw Error('field name is empty')
    }
    if (!handle) {
      throw Error('field handle is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.manageStructureUpdateExposedEndpoint({
        nodeID, moduleID,
      }),
    }
    cfg.data = {
      composeModuleID,
      composeNamespaceID,
      name,
      handle,
      fields,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  manageStructureUpdateExposedEndpoint (a: KV): string {
    const {
      nodeID,
      moduleID,
    } = a || {}
    return `/nodes/${nodeID}/modules/${moduleID}/exposed`
  }

  // Remove from federation
  async manageStructureRemoveExposed (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      nodeID,
      moduleID,
    } = (a as KV) || {}
    if (!nodeID) {
      throw Error('field nodeID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.manageStructureRemoveExposedEndpoint({
        nodeID, moduleID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  manageStructureRemoveExposedEndpoint (a: KV): string {
    const {
      nodeID,
      moduleID,
    } = a || {}
    return `/nodes/${nodeID}/modules/${moduleID}/exposed`
  }

  // Shared settings for module
  async manageStructureReadShared (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      nodeID,
      moduleID,
    } = (a as KV) || {}
    if (!nodeID) {
      throw Error('field nodeID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.manageStructureReadSharedEndpoint({
        nodeID, moduleID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  manageStructureReadSharedEndpoint (a: KV): string {
    const {
      nodeID,
      moduleID,
    } = a || {}
    return `/nodes/${nodeID}/modules/${moduleID}/shared`
  }

  // Add fields mappings to federated module
  async manageStructureCreateMappings (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      nodeID,
      moduleID,
      composeModuleID,
      composeNamespaceID,
      fields,
    } = (a as KV) || {}
    if (!nodeID) {
      throw Error('field nodeID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    if (!composeModuleID) {
      throw Error('field composeModuleID is empty')
    }
    if (!composeNamespaceID) {
      throw Error('field composeNamespaceID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'put',
      url: this.manageStructureCreateMappingsEndpoint({
        nodeID, moduleID,
      }),
    }
    cfg.data = {
      composeModuleID,
      composeNamespaceID,
      fields,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  manageStructureCreateMappingsEndpoint (a: KV): string {
    const {
      nodeID,
      moduleID,
    } = a || {}
    return `/nodes/${nodeID}/modules/${moduleID}/mapped`
  }

  // Fields mappings for module
  async manageStructureReadMappings (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      nodeID,
      moduleID,
      composeModuleID,
    } = (a as KV) || {}
    if (!nodeID) {
      throw Error('field nodeID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.manageStructureReadMappingsEndpoint({
        nodeID, moduleID,
      }),
    }
    cfg.params = {
      composeModuleID,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  manageStructureReadMappingsEndpoint (a: KV): string {
    const {
      nodeID,
      moduleID,
    } = a || {}
    return `/nodes/${nodeID}/modules/${moduleID}/mapped`
  }

  // List of shared/exposed/mapped modules
  async manageStructureListAll (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      nodeID,
      shared,
      exposed,
      mapped,
    } = (a as KV) || {}
    if (!nodeID) {
      throw Error('field nodeID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.manageStructureListAllEndpoint({
        nodeID,
      }),
    }
    cfg.params = {
      shared,
      exposed,
      mapped,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  manageStructureListAllEndpoint (a: KV): string {
    const {
      nodeID,
    } = a || {}
    return `/nodes/${nodeID}/modules/`
  }

  // List all exposed modules changes
  async syncStructureReadExposedInternal (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      nodeID,
      lastSync,
      query,
      limit,
      pageCursor,
      sort,
    } = (a as KV) || {}
    if (!nodeID) {
      throw Error('field nodeID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.syncStructureReadExposedInternalEndpoint({
        nodeID,
      }),
    }
    cfg.params = {
      lastSync,
      query,
      limit,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  syncStructureReadExposedInternalEndpoint (a: KV): string {
    const {
      nodeID,
    } = a || {}
    return `/nodes/${nodeID}/modules/exposed/`
  }

  // List all exposed modules changes in activity streams format
  async syncStructureReadExposedSocial (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      nodeID,
      lastSync,
      query,
      limit,
      pageCursor,
      sort,
    } = (a as KV) || {}
    if (!nodeID) {
      throw Error('field nodeID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.syncStructureReadExposedSocialEndpoint({
        nodeID,
      }),
    }
    cfg.params = {
      lastSync,
      query,
      limit,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  syncStructureReadExposedSocialEndpoint (a: KV): string {
    const {
      nodeID,
    } = a || {}
    return `/nodes/${nodeID}/modules/exposed/activity-stream`
  }

  // List all record changes
  async syncDataReadExposedAll (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      nodeID,
      lastSync,
      query,
      limit,
      pageCursor,
      sort,
    } = (a as KV) || {}
    if (!nodeID) {
      throw Error('field nodeID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.syncDataReadExposedAllEndpoint({
        nodeID,
      }),
    }
    cfg.params = {
      lastSync,
      query,
      limit,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  syncDataReadExposedAllEndpoint (a: KV): string {
    const {
      nodeID,
    } = a || {}
    return `/nodes/${nodeID}/modules/exposed/records/`
  }

  // List all records per module
  async syncDataReadExposedInternal (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      nodeID,
      moduleID,
      lastSync,
      query,
      limit,
      pageCursor,
      sort,
    } = (a as KV) || {}
    if (!nodeID) {
      throw Error('field nodeID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.syncDataReadExposedInternalEndpoint({
        nodeID, moduleID,
      }),
    }
    cfg.params = {
      lastSync,
      query,
      limit,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  syncDataReadExposedInternalEndpoint (a: KV): string {
    const {
      nodeID,
      moduleID,
    } = a || {}
    return `/nodes/${nodeID}/modules/${moduleID}/records/`
  }

  // List all records per module in activitystreams format
  async syncDataReadExposedSocial (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      nodeID,
      moduleID,
      lastSync,
      query,
      limit,
      pageCursor,
      sort,
    } = (a as KV) || {}
    if (!nodeID) {
      throw Error('field nodeID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.syncDataReadExposedSocialEndpoint({
        nodeID, moduleID,
      }),
    }
    cfg.params = {
      lastSync,
      query,
      limit,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  syncDataReadExposedSocialEndpoint (a: KV): string {
    const {
      nodeID,
      moduleID,
    } = a || {}
    return `/nodes/${nodeID}/modules/${moduleID}/records/activity-stream/`
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
