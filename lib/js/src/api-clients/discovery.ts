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

function stdResolve (response: AxiosResponse<CortezaResponse>): KV|Promise<never> {
  if (response.data.error) {
    return Promise.reject(response.data.error)
  } else {
    return response.data.response as KV
  }
}

export default class Discovery {
  protected baseURL?: string;
  protected accessTokenFn?: () => (string | undefined);
  protected headers: Headers = {};

  constructor ({ baseURL, headers, accessTokenFn }: Ctor) {
    this.baseURL = baseURL
    this.accessTokenFn = accessTokenFn
    this.headers = {
      'Content-Type': 'application/json',
    }

    this.setHeaders(headers)
  }

  setAccessTokenFn (fn: () => string | undefined): Discovery {
    this.accessTokenFn = fn
    return this
  }

  setHeaders (headers?: Headers): Discovery {
    if (typeof headers === 'object') {
      this.headers = headers
    }

    return this
  }

  setHeader (name: string, value: string | undefined): Discovery {
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

  async query (a: KV, extra: AxiosRequestConfig = {}): Promise<KV> {
    const {
      query = '',
      from,
      size,
      resourceTypes,
    } = a || {}

    const params = new URLSearchParams()
    if (resourceTypes && Array.isArray(resourceTypes)) {
      resourceTypes.forEach(t => params.append('resourceTypes', t))
    }

    if (from) params.append('from', from.toString())
    if (size) params.append('size', size.toString())

    const cfg: AxiosRequestConfig = {
      ...extra,
      method: 'get',
      url: `/?q=${query}`,
      params,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  queryCancellable (a: KV, extra: AxiosRequestConfig = {}): { response: () => Promise<KV>; cancel: () => void } {
    const cancelTokenSource = axios.CancelToken.source()
    const options = { ...extra, cancelToken: cancelTokenSource.token }

    return {
      response: () => this.query(a, options),
      cancel: () => {
        cancelTokenSource.cancel()
      },
    }
  }
}
