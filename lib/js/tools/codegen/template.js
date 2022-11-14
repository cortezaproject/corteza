export const template = `
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

export default class {{className}} {
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

  setAccessTokenFn (fn: () => string | undefined): {{className}} {
    this.accessTokenFn = fn
    return this
  }

  setHeaders (headers?: Headers): {{className}} {
    if (typeof headers === 'object') {
      this.headers = headers
    }

    return this
  }

  setHeader (name: string, value: string | undefined): {{className}} {
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

{{#endpoints}}
  // {{title}}{{#description}}
  // {{description}}{{/description}}
  async {{fname}} ({{#if fargs}}a: KV, {{/if}}extra: AxiosRequestConfig = {}): Promise<KV> {
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
    return this.api().request(cfg).then(result => stdResolve(result))
  }
  
  {{fname}}Endpoint ({{#if pathParams}}a: KV{{/if}}): string {
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
