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

export default class Compose {
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

  setAccessTokenFn (fn: () => string | undefined): Compose {
    this.accessTokenFn = fn
    return this
  }

  setHeaders (headers?: Headers): Compose {
    if (typeof headers === 'object') {
      this.headers = headers
    }

    return this
  }

  setHeader (name: string, value: string | undefined): Compose {
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

  // List namespaces
  async namespaceList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      query,
      slug,
      limit,
      incTotal,
      labels,
      pageCursor,
      sort,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.namespaceListEndpoint(),
    }
    cfg.params = {
      query,
      slug,
      limit,
      incTotal,
      labels,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  namespaceListEndpoint (): string {
    return '/namespace/'
  }

  // Create namespace
  async namespaceCreate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      name,
      labels,
      slug,
      enabled,
      meta,
    } = (a as KV) || {}
    if (!name) {
      throw Error('field name is empty')
    }
    if (!meta) {
      throw Error('field meta is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.namespaceCreateEndpoint(),
    }
    cfg.data = {
      name,
      labels,
      slug,
      enabled,
      meta,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  namespaceCreateEndpoint (): string {
    return '/namespace/'
  }

  // Read namespace
  async namespaceRead (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.namespaceReadEndpoint({
        namespaceID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  namespaceReadEndpoint (a: KV): string {
    const {
      namespaceID,
    } = a || {}
    return `/namespace/${namespaceID}`
  }

  // Update namespace
  async namespaceUpdate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      name,
      slug,
      enabled,
      meta,
      labels,
      updatedAt,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!name) {
      throw Error('field name is empty')
    }
    if (!meta) {
      throw Error('field meta is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.namespaceUpdateEndpoint({
        namespaceID,
      }),
    }
    cfg.data = {
      name,
      slug,
      enabled,
      meta,
      labels,
      updatedAt,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  namespaceUpdateEndpoint (a: KV): string {
    const {
      namespaceID,
    } = a || {}
    return `/namespace/${namespaceID}`
  }

  // Delete namespace
  async namespaceDelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.namespaceDeleteEndpoint({
        namespaceID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  namespaceDeleteEndpoint (a: KV): string {
    const {
      namespaceID,
    } = a || {}
    return `/namespace/${namespaceID}`
  }

  // Upload namespace assets
  async namespaceUpload (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      upload,
    } = (a as KV) || {}
    if (!upload) {
      throw Error('field upload is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.namespaceUploadEndpoint(),
    }
    cfg.data = {
      upload,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  namespaceUploadEndpoint (): string {
    return '/namespace/upload'
  }

  // Clone compose namespace
  async namespaceClone (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      name,
      slug,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!name) {
      throw Error('field name is empty')
    }
    if (!slug) {
      throw Error('field slug is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.namespaceCloneEndpoint({
        namespaceID,
      }),
    }
    cfg.data = {
      name,
      slug,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  namespaceCloneEndpoint (a: KV): string {
    const {
      namespaceID,
    } = a || {}
    return `/namespace/${namespaceID}/clone`
  }

  // Export compose namespace
  async namespaceExport (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      filename,
      ext,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!filename) {
      throw Error('field filename is empty')
    }
    if (!ext) {
      throw Error('field ext is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.namespaceExportEndpoint({
        namespaceID, filename, ext,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  namespaceExportEndpoint (a: KV): string {
    const {
      namespaceID,
      filename,
      ext,
    } = a || {}
    return `/namespace/${namespaceID}/export/${filename}.zip`
  }

  // Initiate namespace import session
  async namespaceImportInit (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      upload,
    } = (a as KV) || {}
    if (!upload) {
      throw Error('field upload is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.namespaceImportInitEndpoint(),
    }
    cfg.data = {
      upload,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  namespaceImportInitEndpoint (): string {
    return '/namespace/import'
  }

  // Run namespace import
  async namespaceImportRun (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      sessionID,
      name,
      slug,
    } = (a as KV) || {}
    if (!sessionID) {
      throw Error('field sessionID is empty')
    }
    if (!name) {
      throw Error('field name is empty')
    }
    if (!slug) {
      throw Error('field slug is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.namespaceImportRunEndpoint({
        sessionID,
      }),
    }
    cfg.data = {
      name,
      slug,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  namespaceImportRunEndpoint (a: KV): string {
    const {
      sessionID,
    } = a || {}
    return `/namespace/import/${sessionID}`
  }

  // Fire compose:namespace trigger
  async namespaceTriggerScript (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      script,
      args,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!script) {
      throw Error('field script is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.namespaceTriggerScriptEndpoint({
        namespaceID,
      }),
    }
    cfg.data = {
      script,
      args,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  namespaceTriggerScriptEndpoint (a: KV): string {
    const {
      namespaceID,
    } = a || {}
    return `/namespace/${namespaceID}/trigger`
  }

  // List translation
  async namespaceListTranslations (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.namespaceListTranslationsEndpoint({
        namespaceID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  namespaceListTranslationsEndpoint (a: KV): string {
    const {
      namespaceID,
    } = a || {}
    return `/namespace/${namespaceID}/translation`
  }

  // Update translation
  async namespaceUpdateTranslations (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      translations,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!translations) {
      throw Error('field translations is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'patch',
      url: this.namespaceUpdateTranslationsEndpoint({
        namespaceID,
      }),
    }
    cfg.data = {
      translations,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  namespaceUpdateTranslationsEndpoint (a: KV): string {
    const {
      namespaceID,
    } = a || {}
    return `/namespace/${namespaceID}/translation`
  }

  // List available pages
  async pageList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      selfID,
      moduleID,
      query,
      handle,
      labels,
      limit,
      pageCursor,
      sort,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.pageListEndpoint({
        namespaceID,
      }),
    }
    cfg.params = {
      selfID,
      moduleID,
      query,
      handle,
      labels,
      limit,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  pageListEndpoint (a: KV): string {
    const {
      namespaceID,
    } = a || {}
    return `/namespace/${namespaceID}/page/`
  }

  // Create page
  async pageCreate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      selfID,
      moduleID,
      title,
      handle,
      description,
      weight,
      labels,
      visible,
      blocks,
      config,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!title) {
      throw Error('field title is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.pageCreateEndpoint({
        namespaceID,
      }),
    }
    cfg.data = {
      selfID,
      moduleID,
      title,
      handle,
      description,
      weight,
      labels,
      visible,
      blocks,
      config,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  pageCreateEndpoint (a: KV): string {
    const {
      namespaceID,
    } = a || {}
    return `/namespace/${namespaceID}/page/`
  }

  // Get page details
  async pageRead (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      pageID,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!pageID) {
      throw Error('field pageID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.pageReadEndpoint({
        namespaceID, pageID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  pageReadEndpoint (a: KV): string {
    const {
      namespaceID,
      pageID,
    } = a || {}
    return `/namespace/${namespaceID}/page/${pageID}`
  }

  // Get page all (non-record) pages, hierarchically
  async pageTree (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.pageTreeEndpoint({
        namespaceID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  pageTreeEndpoint (a: KV): string {
    const {
      namespaceID,
    } = a || {}
    return `/namespace/${namespaceID}/page/tree`
  }

  // Update page
  async pageUpdate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      pageID,
      selfID,
      moduleID,
      title,
      handle,
      description,
      weight,
      labels,
      visible,
      blocks,
      config,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!pageID) {
      throw Error('field pageID is empty')
    }
    if (!title) {
      throw Error('field title is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.pageUpdateEndpoint({
        namespaceID, pageID,
      }),
    }
    cfg.data = {
      selfID,
      moduleID,
      title,
      handle,
      description,
      weight,
      labels,
      visible,
      blocks,
      config,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  pageUpdateEndpoint (a: KV): string {
    const {
      namespaceID,
      pageID,
    } = a || {}
    return `/namespace/${namespaceID}/page/${pageID}`
  }

  // Reorder pages
  async pageReorder (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      selfID,
      pageIDs,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!selfID) {
      throw Error('field selfID is empty')
    }
    if (!pageIDs) {
      throw Error('field pageIDs is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.pageReorderEndpoint({
        namespaceID, selfID,
      }),
    }
    cfg.data = {
      pageIDs,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  pageReorderEndpoint (a: KV): string {
    const {
      namespaceID,
      selfID,
    } = a || {}
    return `/namespace/${namespaceID}/page/${selfID}/reorder`
  }

  // Delete page
  async pageDelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      pageID,
      strategy,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!pageID) {
      throw Error('field pageID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.pageDeleteEndpoint({
        namespaceID, pageID,
      }),
    }
    cfg.params = {
      strategy,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  pageDeleteEndpoint (a: KV): string {
    const {
      namespaceID,
      pageID,
    } = a || {}
    return `/namespace/${namespaceID}/page/${pageID}`
  }

  // Uploads attachment to page
  async pageUpload (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      pageID,
      upload,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!pageID) {
      throw Error('field pageID is empty')
    }
    if (!upload) {
      throw Error('field upload is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.pageUploadEndpoint({
        namespaceID, pageID,
      }),
    }
    cfg.data = {
      upload,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  pageUploadEndpoint (a: KV): string {
    const {
      namespaceID,
      pageID,
    } = a || {}
    return `/namespace/${namespaceID}/page/${pageID}/attachment`
  }

  // Fire compose:page trigger
  async pageTriggerScript (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      pageID,
      script,
      args,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!pageID) {
      throw Error('field pageID is empty')
    }
    if (!script) {
      throw Error('field script is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.pageTriggerScriptEndpoint({
        namespaceID, pageID,
      }),
    }
    cfg.data = {
      script,
      args,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  pageTriggerScriptEndpoint (a: KV): string {
    const {
      namespaceID,
      pageID,
    } = a || {}
    return `/namespace/${namespaceID}/page/${pageID}/trigger`
  }

  // List page translation
  async pageListTranslations (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      pageID,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!pageID) {
      throw Error('field pageID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.pageListTranslationsEndpoint({
        namespaceID, pageID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  pageListTranslationsEndpoint (a: KV): string {
    const {
      namespaceID,
      pageID,
    } = a || {}
    return `/namespace/${namespaceID}/page/${pageID}/translation`
  }

  // Update page translation
  async pageUpdateTranslations (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      pageID,
      translations,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!pageID) {
      throw Error('field pageID is empty')
    }
    if (!translations) {
      throw Error('field translations is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'patch',
      url: this.pageUpdateTranslationsEndpoint({
        namespaceID, pageID,
      }),
    }
    cfg.data = {
      translations,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  pageUpdateTranslationsEndpoint (a: KV): string {
    const {
      namespaceID,
      pageID,
    } = a || {}
    return `/namespace/${namespaceID}/page/${pageID}/translation`
  }

  // List modules
  async moduleList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      query,
      name,
      handle,
      limit,
      incTotal,
      pageCursor,
      labels,
      sort,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.moduleListEndpoint({
        namespaceID,
      }),
    }
    cfg.params = {
      query,
      name,
      handle,
      limit,
      incTotal,
      pageCursor,
      labels,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  moduleListEndpoint (a: KV): string {
    const {
      namespaceID,
    } = a || {}
    return `/namespace/${namespaceID}/module/`
  }

  // Create module
  async moduleCreate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      name,
      handle,
      config,
      meta,
      fields,
      labels,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!name) {
      throw Error('field name is empty')
    }
    if (!meta) {
      throw Error('field meta is empty')
    }
    if (!fields) {
      throw Error('field fields is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.moduleCreateEndpoint({
        namespaceID,
      }),
    }
    cfg.data = {
      name,
      handle,
      config,
      meta,
      fields,
      labels,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  moduleCreateEndpoint (a: KV): string {
    const {
      namespaceID,
    } = a || {}
    return `/namespace/${namespaceID}/module/`
  }

  // Read module
  async moduleRead (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      moduleID,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.moduleReadEndpoint({
        namespaceID, moduleID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  moduleReadEndpoint (a: KV): string {
    const {
      namespaceID,
      moduleID,
    } = a || {}
    return `/namespace/${namespaceID}/module/${moduleID}`
  }

  // Update module
  async moduleUpdate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      moduleID,
      name,
      handle,
      config,
      meta,
      fields,
      labels,
      updatedAt,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    if (!name) {
      throw Error('field name is empty')
    }
    if (!meta) {
      throw Error('field meta is empty')
    }
    if (!fields) {
      throw Error('field fields is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.moduleUpdateEndpoint({
        namespaceID, moduleID,
      }),
    }
    cfg.data = {
      name,
      handle,
      config,
      meta,
      fields,
      labels,
      updatedAt,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  moduleUpdateEndpoint (a: KV): string {
    const {
      namespaceID,
      moduleID,
    } = a || {}
    return `/namespace/${namespaceID}/module/${moduleID}`
  }

  // Delete module
  async moduleDelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      moduleID,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.moduleDeleteEndpoint({
        namespaceID, moduleID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  moduleDeleteEndpoint (a: KV): string {
    const {
      namespaceID,
      moduleID,
    } = a || {}
    return `/namespace/${namespaceID}/module/${moduleID}`
  }

  // Fire compose:module trigger
  async moduleTriggerScript (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      moduleID,
      script,
      args,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    if (!script) {
      throw Error('field script is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.moduleTriggerScriptEndpoint({
        namespaceID, moduleID,
      }),
    }
    cfg.data = {
      script,
      args,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  moduleTriggerScriptEndpoint (a: KV): string {
    const {
      namespaceID,
      moduleID,
    } = a || {}
    return `/namespace/${namespaceID}/module/${moduleID}/trigger`
  }

  // List moudle translation
  async moduleListTranslations (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      moduleID,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.moduleListTranslationsEndpoint({
        namespaceID, moduleID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  moduleListTranslationsEndpoint (a: KV): string {
    const {
      namespaceID,
      moduleID,
    } = a || {}
    return `/namespace/${namespaceID}/module/${moduleID}/translation`
  }

  // Update module translation
  async moduleUpdateTranslations (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      moduleID,
      translations,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    if (!translations) {
      throw Error('field translations is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'patch',
      url: this.moduleUpdateTranslationsEndpoint({
        namespaceID, moduleID,
      }),
    }
    cfg.data = {
      translations,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  moduleUpdateTranslationsEndpoint (a: KV): string {
    const {
      namespaceID,
      moduleID,
    } = a || {}
    return `/namespace/${namespaceID}/module/${moduleID}/translation`
  }

  // Generates report from module records
  async recordReport (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      moduleID,
      metrics,
      dimensions,
      filter,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    if (!dimensions) {
      throw Error('field dimensions is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.recordReportEndpoint({
        namespaceID, moduleID,
      }),
    }
    cfg.params = {
      metrics,
      dimensions,
      filter,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  recordReportEndpoint (a: KV): string {
    const {
      namespaceID,
      moduleID,
    } = a || {}
    return `/namespace/${namespaceID}/module/${moduleID}/record/report`
  }

  // List/read records from module section
  async recordList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      moduleID,
      query,
      meta,
      deleted,
      limit,
      incTotal,
      incPageNavigation,
      pageCursor,
      sort,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.recordListEndpoint({
        namespaceID, moduleID,
      }),
    }
    cfg.params = {
      query,
      meta,
      deleted,
      limit,
      incTotal,
      incPageNavigation,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  recordListEndpoint (a: KV): string {
    const {
      namespaceID,
      moduleID,
    } = a || {}
    return `/namespace/${namespaceID}/module/${moduleID}/record/`
  }

  // Initiate record import session
  async recordImportInit (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      moduleID,
      upload,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    if (!upload) {
      throw Error('field upload is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.recordImportInitEndpoint({
        namespaceID, moduleID,
      }),
    }
    cfg.data = {
      upload,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  recordImportInitEndpoint (a: KV): string {
    const {
      namespaceID,
      moduleID,
    } = a || {}
    return `/namespace/${namespaceID}/module/${moduleID}/record/import`
  }

  // Run record import
  async recordImportRun (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      moduleID,
      sessionID,
      fields,
      onError,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    if (!sessionID) {
      throw Error('field sessionID is empty')
    }
    if (!fields) {
      throw Error('field fields is empty')
    }
    if (!onError) {
      throw Error('field onError is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'patch',
      url: this.recordImportRunEndpoint({
        namespaceID, moduleID, sessionID,
      }),
    }
    cfg.data = {
      fields,
      onError,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  recordImportRunEndpoint (a: KV): string {
    const {
      namespaceID,
      moduleID,
      sessionID,
    } = a || {}
    return `/namespace/${namespaceID}/module/${moduleID}/record/import/${sessionID}`
  }

  // Get import progress
  async recordImportProgress (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      moduleID,
      sessionID,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    if (!sessionID) {
      throw Error('field sessionID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.recordImportProgressEndpoint({
        namespaceID, moduleID, sessionID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  recordImportProgressEndpoint (a: KV): string {
    const {
      namespaceID,
      moduleID,
      sessionID,
    } = a || {}
    return `/namespace/${namespaceID}/module/${moduleID}/record/import/${sessionID}`
  }

  // Exports records that match
  async recordExport (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      moduleID,
      filename,
      ext,
      filter,
      fields,
      timezone,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    if (!ext) {
      throw Error('field ext is empty')
    }
    if (!fields) {
      throw Error('field fields is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.recordExportEndpoint({
        namespaceID, moduleID, filename, ext,
      }),
    }
    cfg.params = {
      filter,
      fields,
      timezone,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  recordExportEndpoint (a: KV): string {
    const {
      namespaceID,
      moduleID,
      filename,
      ext,
    } = a || {}
    return `/namespace/${namespaceID}/module/${moduleID}/record/export${filename}.${ext}`
  }

  // Executes server-side procedure over one or more module records
  async recordExec (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      moduleID,
      procedure,
      args,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    if (!procedure) {
      throw Error('field procedure is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.recordExecEndpoint({
        namespaceID, moduleID, procedure,
      }),
    }
    cfg.data = {
      args,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  recordExecEndpoint (a: KV): string {
    const {
      namespaceID,
      moduleID,
      procedure,
    } = a || {}
    return `/namespace/${namespaceID}/module/${moduleID}/record/exec/${procedure}`
  }

  // Create record in module section
  async recordCreate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      moduleID,
      values,
      ownedBy,
      records,
      meta,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.recordCreateEndpoint({
        namespaceID, moduleID,
      }),
    }
    cfg.data = {
      values,
      ownedBy,
      records,
      meta,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  recordCreateEndpoint (a: KV): string {
    const {
      namespaceID,
      moduleID,
    } = a || {}
    return `/namespace/${namespaceID}/module/${moduleID}/record/`
  }

  // Read records by ID from module section
  async recordRead (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      moduleID,
      recordID,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    if (!recordID) {
      throw Error('field recordID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.recordReadEndpoint({
        namespaceID, moduleID, recordID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  recordReadEndpoint (a: KV): string {
    const {
      namespaceID,
      moduleID,
      recordID,
    } = a || {}
    return `/namespace/${namespaceID}/module/${moduleID}/record/${recordID}`
  }

  // Update records in module section
  async recordUpdate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      moduleID,
      recordID,
      values,
      ownedBy,
      meta,
      records,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    if (!recordID) {
      throw Error('field recordID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.recordUpdateEndpoint({
        namespaceID, moduleID, recordID,
      }),
    }
    cfg.data = {
      values,
      ownedBy,
      meta,
      records,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  recordUpdateEndpoint (a: KV): string {
    const {
      namespaceID,
      moduleID,
      recordID,
    } = a || {}
    return `/namespace/${namespaceID}/module/${moduleID}/record/${recordID}`
  }

  // Delete record row from module section
  async recordBulkDelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      moduleID,
      recordIDs,
      truncate,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.recordBulkDeleteEndpoint({
        namespaceID, moduleID,
      }),
    }
    cfg.data = {
      recordIDs,
      truncate,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  recordBulkDeleteEndpoint (a: KV): string {
    const {
      namespaceID,
      moduleID,
    } = a || {}
    return `/namespace/${namespaceID}/module/${moduleID}/record/`
  }

  // Delete record row from module section
  async recordDelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      moduleID,
      recordID,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    if (!recordID) {
      throw Error('field recordID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.recordDeleteEndpoint({
        namespaceID, moduleID, recordID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  recordDeleteEndpoint (a: KV): string {
    const {
      namespaceID,
      moduleID,
      recordID,
    } = a || {}
    return `/namespace/${namespaceID}/module/${moduleID}/record/${recordID}`
  }

  // Uploads attachment and validates it against record field requirements
  async recordUpload (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      moduleID,
      recordID,
      fieldName,
      upload,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    if (!fieldName) {
      throw Error('field fieldName is empty')
    }
    if (!upload) {
      throw Error('field upload is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.recordUploadEndpoint({
        namespaceID, moduleID,
      }),
    }
    cfg.data = {
      recordID,
      fieldName,
      upload,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  recordUploadEndpoint (a: KV): string {
    const {
      namespaceID,
      moduleID,
    } = a || {}
    return `/namespace/${namespaceID}/module/${moduleID}/record/attachment`
  }

  // Fire compose:record trigger
  async recordTriggerScript (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      moduleID,
      recordID,
      script,
      values,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    if (!recordID) {
      throw Error('field recordID is empty')
    }
    if (!script) {
      throw Error('field script is empty')
    }
    if (!values) {
      throw Error('field values is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.recordTriggerScriptEndpoint({
        namespaceID, moduleID, recordID,
      }),
    }
    cfg.data = {
      script,
      values,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  recordTriggerScriptEndpoint (a: KV): string {
    const {
      namespaceID,
      moduleID,
      recordID,
    } = a || {}
    return `/namespace/${namespaceID}/module/${moduleID}/record/${recordID}/trigger`
  }

  // Fire compose:record trigger
  async recordTriggerScriptOnList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      moduleID,
      script,
      args,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    if (!script) {
      throw Error('field script is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.recordTriggerScriptOnListEndpoint({
        namespaceID, moduleID,
      }),
    }
    cfg.data = {
      script,
      args,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  recordTriggerScriptOnListEndpoint (a: KV): string {
    const {
      namespaceID,
      moduleID,
    } = a || {}
    return `/namespace/${namespaceID}/module/${moduleID}/record/trigger`
  }

  // List record revisions
  async recordRevisions (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      moduleID,
      recordID,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!moduleID) {
      throw Error('field moduleID is empty')
    }
    if (!recordID) {
      throw Error('field recordID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.recordRevisionsEndpoint({
        namespaceID, moduleID, recordID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  recordRevisionsEndpoint (a: KV): string {
    const {
      namespaceID,
      moduleID,
      recordID,
    } = a || {}
    return `/namespace/${namespaceID}/module/${moduleID}/record/${recordID}/revisions`
  }

  // List records for data privacy
  async dataPrivacyRecordList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      sensitivityLevelID,
      connectionID,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.dataPrivacyRecordListEndpoint(),
    }
    cfg.params = {
      sensitivityLevelID,
      connectionID,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  dataPrivacyRecordListEndpoint (): string {
    return '/data-privacy/record'
  }

  // List modules
  async dataPrivacyModuleList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      connectionID,
      limit,
      pageCursor,
      sort,
    } = (a as KV) || {}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.dataPrivacyModuleListEndpoint(),
    }
    cfg.params = {
      connectionID,
      limit,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  dataPrivacyModuleListEndpoint (): string {
    return '/data-privacy/module'
  }

  // List/read charts
  async chartList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      query,
      handle,
      labels,
      limit,
      incTotal,
      pageCursor,
      sort,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.chartListEndpoint({
        namespaceID,
      }),
    }
    cfg.params = {
      query,
      handle,
      labels,
      limit,
      incTotal,
      pageCursor,
      sort,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  chartListEndpoint (a: KV): string {
    const {
      namespaceID,
    } = a || {}
    return `/namespace/${namespaceID}/chart/`
  }

  // List/read charts
  async chartCreate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      config,
      name,
      handle,
      labels,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!config) {
      throw Error('field config is empty')
    }
    if (!name) {
      throw Error('field name is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.chartCreateEndpoint({
        namespaceID,
      }),
    }
    cfg.data = {
      config,
      name,
      handle,
      labels,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  chartCreateEndpoint (a: KV): string {
    const {
      namespaceID,
    } = a || {}
    return `/namespace/${namespaceID}/chart/`
  }

  // Read charts by ID
  async chartRead (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      chartID,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!chartID) {
      throw Error('field chartID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.chartReadEndpoint({
        namespaceID, chartID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  chartReadEndpoint (a: KV): string {
    const {
      namespaceID,
      chartID,
    } = a || {}
    return `/namespace/${namespaceID}/chart/${chartID}`
  }

  // Add/update charts
  async chartUpdate (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      chartID,
      config,
      name,
      handle,
      labels,
      updatedAt,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!chartID) {
      throw Error('field chartID is empty')
    }
    if (!config) {
      throw Error('field config is empty')
    }
    if (!name) {
      throw Error('field name is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.chartUpdateEndpoint({
        namespaceID, chartID,
      }),
    }
    cfg.data = {
      config,
      name,
      handle,
      labels,
      updatedAt,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  chartUpdateEndpoint (a: KV): string {
    const {
      namespaceID,
      chartID,
    } = a || {}
    return `/namespace/${namespaceID}/chart/${chartID}`
  }

  // Delete chart
  async chartDelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      chartID,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!chartID) {
      throw Error('field chartID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.chartDeleteEndpoint({
        namespaceID, chartID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  chartDeleteEndpoint (a: KV): string {
    const {
      namespaceID,
      chartID,
    } = a || {}
    return `/namespace/${namespaceID}/chart/${chartID}`
  }

  // List chart translation
  async chartListTranslations (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      chartID,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!chartID) {
      throw Error('field chartID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.chartListTranslationsEndpoint({
        namespaceID, chartID,
      }),
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  chartListTranslationsEndpoint (a: KV): string {
    const {
      namespaceID,
      chartID,
    } = a || {}
    return `/namespace/${namespaceID}/chart/${chartID}/translation`
  }

  // Update chart translation
  async chartUpdateTranslations (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      namespaceID,
      chartID,
      translations,
    } = (a as KV) || {}
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!chartID) {
      throw Error('field chartID is empty')
    }
    if (!translations) {
      throw Error('field translations is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'patch',
      url: this.chartUpdateTranslationsEndpoint({
        namespaceID, chartID,
      }),
    }
    cfg.data = {
      translations,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  chartUpdateTranslationsEndpoint (a: KV): string {
    const {
      namespaceID,
      chartID,
    } = a || {}
    return `/namespace/${namespaceID}/chart/${chartID}/translation`
  }

  // Send email from the Compose
  async notificationEmailSend (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      to,
      cc,
      replyTo,
      subject,
      content,
      remoteAttachments,
    } = (a as KV) || {}
    if (!to) {
      throw Error('field to is empty')
    }
    if (!content) {
      throw Error('field content is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'post',
      url: this.notificationEmailSendEndpoint(),
    }
    cfg.data = {
      to,
      cc,
      replyTo,
      subject,
      content,
      remoteAttachments,
    }
    return this.api().request(cfg).then(result => stdResolve(result))
  }

  notificationEmailSendEndpoint (): string {
    return '/notification/email'
  }

  // List, filter all page attachments
  async attachmentList (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      kind,
      namespaceID,
      sign,
      userID,
      pageID,
      moduleID,
      recordID,
      fieldName,
      limit,
      pageCursor,
    } = (a as KV) || {}
    if (!kind) {
      throw Error('field kind is empty')
    }
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.attachmentListEndpoint({
        kind, namespaceID,
      }),
    }
    cfg.params = {
      sign,
      userID,
      pageID,
      moduleID,
      recordID,
      fieldName,
      limit,
      pageCursor,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  attachmentListEndpoint (a: KV): string {
    const {
      kind,
      namespaceID,
    } = a || {}
    return `/namespace/${namespaceID}/attachment/${kind}/`
  }

  // Attachment details
  async attachmentRead (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      kind,
      namespaceID,
      attachmentID,
      sign,
      userID,
    } = (a as KV) || {}
    if (!kind) {
      throw Error('field kind is empty')
    }
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!attachmentID) {
      throw Error('field attachmentID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: this.attachmentReadEndpoint({
        kind, namespaceID, attachmentID,
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
      namespaceID,
      attachmentID,
    } = a || {}
    return `/namespace/${namespaceID}/attachment/${kind}/${attachmentID}`
  }

  // Delete attachment
  async attachmentDelete (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      kind,
      namespaceID,
      attachmentID,
      sign,
      userID,
    } = (a as KV) || {}
    if (!kind) {
      throw Error('field kind is empty')
    }
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
    }
    if (!attachmentID) {
      throw Error('field attachmentID is empty')
    }
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'delete',
      url: this.attachmentDeleteEndpoint({
        kind, namespaceID, attachmentID,
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
      namespaceID,
      attachmentID,
    } = a || {}
    return `/namespace/${namespaceID}/attachment/${kind}/${attachmentID}`
  }

  // Serves attached file
  async attachmentOriginal (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      kind,
      namespaceID,
      attachmentID,
      name,
      sign,
      userID,
      download,
    } = (a as KV) || {}
    if (!kind) {
      throw Error('field kind is empty')
    }
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
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
        kind, namespaceID, attachmentID, name,
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
      namespaceID,
      attachmentID,
      name,
    } = a || {}
    return `/namespace/${namespaceID}/attachment/${kind}/${attachmentID}/original/${name}`
  }

  // Serves preview of an attached file
  async attachmentPreview (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      kind,
      namespaceID,
      attachmentID,
      ext,
      sign,
      userID,
    } = (a as KV) || {}
    if (!kind) {
      throw Error('field kind is empty')
    }
    if (!namespaceID) {
      throw Error('field namespaceID is empty')
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
        kind, namespaceID, attachmentID, ext,
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
      namespaceID,
      attachmentID,
      ext,
    } = a || {}
    return `/namespace/${namespaceID}/attachment/${kind}/${attachmentID}/preview.${ext}`
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

  // List all available automation scripts for compose resources
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

}
