export const template = `
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

interface HumanErrorDetail {
  kind?: string
  message?: string
  meta?: { [key: string]: unknown }
}

interface HumanError {
  message?: string
  details?: HumanErrorDetail[]
}

export interface HumanApiError extends Error {
  details?: HumanErrorDetail[]
}

interface HumanResponse {
  error?: string | HumanError
  response?: unknown
}

// The API reports failure in the body of a 200, as either a bare string or an
// { message } object. Both reach callers as an Error so that a catch block can
// read .message and stringify to something a user can act on.
//
// \`details\` rides along where the API sends it: record validation reports one
// entry per issue, each naming the field it belongs to, and the summary message
// on its own ("2 issue(s) found") says nothing a caller can place or act on.
function stdResolve(response: AxiosResponse<HumanResponse>): KV | Promise<never> {
  if (response.data.error) {
    const err = response.data.error
    if (typeof err === 'string') {
      return Promise.reject(new Error(err))
    }

    const out: HumanApiError = new Error(err.message || 'unknown error')
    if (Array.isArray(err.details) && err.details.length > 0) {
      out.details = err.details
    }

    return Promise.reject(out)
  } else if (response.data.response) {
    return response.data.response as KV
  } else {
    return response.data as KV
  }
}

export default class {{className}} {
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

  setAccessTokenFn(fn: () => string | undefined): {{className}} {
    this.accessTokenFn = fn
    return this
  }

  setHeaders(headers?: Headers): {{className}} {
    if (typeof headers === 'object') {
      this.headers = headers
    }

    return this
  }

  setHeader(name: string, value: string | undefined): {{className}} {
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

{{#endpoints}}
  // {{title}}{{#description}}
  // {{description}}{{/description}}
  async {{fname}}({{#if fargs}}a: KV, {{/if}}extra: AxiosRequestConfig = {}): Promise<KV> {
    {{#if fargs}}const {
      {{#fargs}}
      {{.}},
      {{/fargs}}
    } = (a as KV) || {}{{/if}}
    {{#required}}
    if (!{{.}}) {
      throw Error('field {{.}} is empty')
    }
    {{/required}}
    const cfg: AxiosRequestConfig = {
      ...extra,
      method: '{{method}}',
      url: this.{{fname}}Endpoint({{#if pathParams}}{
        {{#pathParams}}{{.}}, {{/pathParams}}
      }{{/if}}),
    }
    {{#hasParams}}cfg.params = {
      {{#params}}
      {{.}},
      {{/params}}
    }
    {{/hasParams}}{{#hasData}}cfg.data = {
      {{#data}}
      {{.}},
      {{/data}}
    }{{/hasData}}
    return this.api()
      .request(cfg)
      .then(result => stdResolve(result))
  }

  {{fname}}Cancellable(
    {{#if fargs}}a: KV, {{/if}}extra: AxiosRequestConfig = {},
  ): { response: (a: KV, extra?: AxiosRequestConfig) => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.{{fname}}({{#if fargs}}a, {{/if}}options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }

  {{fname}}Endpoint({{#if pathParams}}a: KV{{/if}}): string {
  {{#if pathParams}}
    const {
      {{#pathParams}}
      {{.}},
      {{/pathParams}}
    } = a || {}
    return \`{{path}}\`
  {{else}}
    return '{{path}}'
  {{/if}}
  }

{{/endpoints}}
}
`
