import axios from 'axios'

export default function (service, opt = {}) {
  if (!opt.baseURL) {
    // @ts-ignore
    if (!window.CortezaDiscoveryAPI) {
      throw new Error('config.js missing or window.CortezaDiscoveryAPI not set')
    }

    // @ts-ignore
    opt.baseURL = `${window.CortezaDiscoveryAPI}/`
  }

  return function (Vue) {
    if (!opt.accessTokenFn) {
      /**
       * Checking if auth plugin was initialized before and
       * hooking on to it's accessTokenFn
       */
      opt.accessTokenFn = Vue.prototype.$auth.accessTokenFn
    }

    // @ts-ignore
    // makes Vue.$<service>API (Vue.$SystemAPI, Vue.$ComposeAPI, Vue.$FederationAPI, Vue.$AutomationAPI) available
    Vue.prototype.$DiscoveryAPI = new Searcher(opt)
  }
}

function stdResolve (response) {
  if (response.data.error) {
    return Promise.reject(response.data.error)
  } else {
    return response.data.response
  }
}

class Searcher {
  baseURL;
  accessTokenFn;
  headers = {};

  constructor ({ baseURL, headers, accessTokenFn }) {
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

  setAccessTokenFn (fn) {
    this.accessTokenFn = fn
    return this
  }

  setHeaders (headers) {
    if (typeof headers === 'object') {
      this.headers = headers
    }

    return this
  }

  setHeader (name, value) {
    if (value === undefined) {
      delete this.headers[name]
    } else {
      this.headers[name] = value
    }

    return this
  }

  api () {
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
  async query (a, extra = {}) {
    const {
      modules,
      namespaces,
      from,
      size,
    } = a || {}

    const params = new URLSearchParams()

    // Filter
    if (modules?.length > 0) modules.forEach(m => params.append('moduleAggs', m))
    if (namespaces?.length > 0) namespaces.forEach(n => params.append('namespaceAggs', n))

    // Pagination
    if (from) params.append('from', from)
    if (size) params.append('size', size)

    const cfg = {
      ...extra,
      method: 'get',
      url: this.queryEndpoint(a),
      params,
    }

    return this.api().request(cfg).then(result => stdResolve(result))
  }

  queryEndpoint (a) {
    const {
      query = '',
    } = a || {}

    return `/?q=${query}`
  }
}
