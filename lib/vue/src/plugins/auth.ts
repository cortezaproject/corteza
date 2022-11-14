import axios, { AxiosInstance } from 'axios'
import { Make } from '../libs/url'
import { system } from '@cortezaproject/corteza-js'
import { PluginFunction } from 'vue'

const accessToken = Symbol('accessToken')
const user = Symbol('user')

/**
 * This is an endpoint of an oauth2 authorization-code flow + refresh token exchange
 * for default client.
 *
 * If you are concerned about security, this is not much different to using a dedicated backend
 * for this SPA that only redirects to a set of allowed URIs.
 */
const oauth2FlowURL = '/oauth2/default-client'
const oauth2InfoURL = '/oauth2/info'
const oauth2Scope = 'profile api'

const storeKeyFlowStarted = 'auth.flow-started'
const storeKeyFinalState = 'auth.state.final'
const storeKeyRefreshToken = 'auth.refresh-token'

const maxStartAttempts = 5

// signature copied from dom definition
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type eventListenerSignature = <K extends keyof WindowEventMap>(type: K, listener: (this: Window, ev: WindowEventMap[K]) => any, options?: boolean | AddEventListenerOptions) => void

interface AuthInfo {
  accessTokenFn: () => string | undefined;
  user: system.User;
}

interface OAuth2TokenResponse {
  aud: string;
  sub: string;
  scope: string;
  access_token: string;
  refresh_token: string;
  expires_in: number;

  name?: string;
  handle?: string;
  email?: string;
  preferred_language?: string;
}

interface PluginOpts {
  cortezaAuthURL: string;
  callbackURL: string;
}

interface AuthCtor {
  app: string;

  /**
   * when true, use console as a logger, no-op otherwise.
   */
  verbose: boolean;

  /**
   * where the authe backend is
   */
  cortezaAuthURL: string;

  /**
   * URL we'll be listening to for callbacksa
   */
  callbackURL: string;

  /**
   * used for redirection
   */
  location: Location;

  /**
   * used for storing
   */
  sessionStorage: Storage;

  /**
   * used for event listeners
   */
  registerEventListener: eventListenerSignature;

  /**
   * Static string with entry-point URL stored at app init
   * so that there is no risk of changes when Vue router gets it's hands on it
   */
  entrypointURL: string;

  /**
   * multiply factor for token expiration
   * this will tell internal refresh system how much
   * before the token expiration we'll refresh the access toke
   *
   * keep in mind that access token is exchanged on every app load
   */
  refreshFactor: number;
}

interface Logger {
  debug(...data: unknown[]): void;
  info(...data: unknown[]): void;
  error(...data: unknown[]): void;
}

export class Auth {
  /**
   * Access token is only stored here (in-memory)!
   * we do not want to keep it in the local store
   */
  private [accessToken]?: string

  /**
   * Access token is only stored here (in-memory)!
   * we do not want to keep it in the local store
   */
  private [user]?: system.User

  /**
   * Name of the app that is using the auth plugin
   */
  readonly app: string

  readonly refreshFactor: number
  readonly verbose: boolean
  readonly cortezaAuthURL: string
  readonly callbackURL: string
  readonly location: Location
  readonly sessionStorage: Storage
  readonly registerEventListener: eventListenerSignature

  /**
   * Application entrypoint URL
   */
  readonly entrypointURL: string

  /**
   * Keeps track of timeout callback in case we re-run it before it timesout
   * @private
   */
  private refreshTimeout?: number

  private $emit?: (event: string, ...args: unknown[]) => unknown

  constructor ({ app, verbose, cortezaAuthURL, callbackURL, entrypointURL, location, sessionStorage, refreshFactor, registerEventListener }: AuthCtor) {
    if (refreshFactor >= 1 || refreshFactor <= 0) {
      throw new Error('refreshFactor should be between 0 and 1')
    }

    this.app = app
    this.verbose = verbose
    this.cortezaAuthURL = cortezaAuthURL
    this.callbackURL = callbackURL
    this.location = location
    this.sessionStorage = sessionStorage
    this.registerEventListener = registerEventListener
    this.refreshFactor = refreshFactor
    this.entrypointURL = entrypointURL

    this.log.debug('initialized auth plugin', {
      app,
      cortezaAuthURL,
      callbackURL,
      entrypointURL,
    })
  }

  vue (vue: Vue): Auth {
    this.$emit = (event, ...args): void => { vue.$emit(event, ...args) }
    return this
  }

  get axios (): AxiosInstance {
    return axios.create({ baseURL: this.cortezaAuthURL })
  }

  /**
   * wrapper for console (when in debug mode) or a simple no-op obj
   */
  get log (): Logger {
    if (this.verbose) {
      return console
    }

    // eslint-disable-next-line @typescript-eslint/no-empty-function
    const noop = (): void => {}

    return {
      debug: noop,
      info: noop,
      error: noop,
    }
  }

  /**
   * Returns function that returns current access token
   */
  get accessTokenFn (): () => string | undefined {
    return (): string | undefined => { return this[accessToken] }
  }

  /**
   * Handles initial authentication check
   *
   * handle function should be called immediately when application is created
   * it checks whether app was requested on an URL with /auth/callback at the end
   * if there is an error or code passed and handles that request appropriately:
   *
   *  .../auth/callback?code=... exchanged authorization code for access token
   *  .../auth/callback?error=... renders an error that we got from the oauth2 provider
   *
   * If handle was called without /auth/callback or without params mentioned above:
   *   if user is not authorized, redirect to the configured path to start oauth2 flow
   *   if user is authorized, continue with execution
   */
  async handle (req: URL = new URL(this.entrypointURL)): Promise<AuthInfo | null> {
    this.log.info('handling authentication')

    // State management
    const dup = this.handleStateManagement()
    if (dup) {
      this.log.debug('duplicate tab: unauthorized')
      throw new Error('Unauthenticated')
    }

    // Handle auth callback requests
    const params = new URLSearchParams(req.search)
    if (this.isCallback(req.pathname) && (params.has('error') || params.has('code'))) {
      if (params.has('error')) {
        throw new Error(params.get('error') || 'authentication flow failed with error')
      }

      this.log.info('handling authentication callback')
      return this.handleCallbackRoute(params.get('state'), (params.has('code') ? params.get('code') as string : ''))
    }

    // Handle auth from the current system state
    this.log.info('handling authentication from state')
    return this.handleState()
  }

  /**
   * Flagging session storage
   *
   * Challenge with sessions and browser tabs:
   * Each browser tab & window interacts with an isolated session. When user clicks on a link to and wants to open
   * it in a new window or a tab, or when tab is duplicated, session contents are copied!
   *
   * Consequences of that are that two (or more) tabs end up with the same session
   * and the same refresh token, and we need to detect if we're dealing with refresh token in an old or new  session.
   *
   * With this function we start the final state and flag the session. This way we'll
   * know, after the redirection to the final location, if this is a final stage or not and if the refresh token
   * belongs to this session or not.
   */
  handleStateManagement (): boolean {
    // See if this is a duplicate
    const dup = this.sessionStorage.getItem(storeKeyFinalState) !== null
    window.sessionStorage.setItem(storeKeyFinalState, Date.now().toString())
    return dup
  }

  bindListeners (): void {
    // binding multiple listeners for cases where some browser refuser
    // to emit one of them.
    this.registerEventListener('pagehide', () => {
      this.cleanFlags()
    })

    this.registerEventListener('unload', () => {
      this.cleanFlags()
    })

    this.registerEventListener('beforeunload', () => {
      this.cleanFlags()
    })
  }

  cleanFlags (): void {
    this.sessionStorage.removeItem(storeKeyFinalState)
  }

  /**
   * Called when refresh token is re-fetched.
   *
   * Cleanup aux items in the session store
   */
  completeFinalState (): void {
    this.sessionStorage.removeItem(storeKeyFlowStarted)

    const stateKey = /^auth\.state\.\w+\.location$/
    for (let i = 0; i < this.sessionStorage.length; i++) {
      const key = this.sessionStorage.key(i)
      if (key !== null && stateKey.test(key)) {
        this.sessionStorage.removeItem(key)
      }
    }
  }

  /**
   * Exchanges the auth parameters for access & refresh token.
   *
   * If the parameters are correct and exchange is successful, the refresh token
   * gets stored in localStorage for further use, the access token and current user get stored
   * in-memory.
   *
   * Function will throw null when user is unauthenticated
   */
  async handleCallbackRoute (state: string|null, code: string): Promise<AuthInfo | null> {
    let finalLocation = this.entrypointURL

    if (state) {
      const storeKeyStateLocation = `auth.state.${state}.location`
      const tmp = this.sessionStorage.getItem(storeKeyStateLocation)
      if (tmp === null) {
        console.warn('state does not match, restarting authentication flow')
        this.startAuthenticationFlow()
        return null
      }

      if (!this.isCallback(tmp)) {
        // if by some coincidence we got callback URL to finalLocation
        // we'll silently ignore it and redirect user back to entrypoint
        finalLocation = tmp
      }

      this.sessionStorage.removeItem(storeKeyStateLocation)
    }

    this.log.info('authorization code received', code)
    const rsp = await this.exchangeCode(code)

    this.log.info('redirecting back to final destination', finalLocation)
    this.cleanFlags()
    this.location.assign(finalLocation)
    return rsp
  }

  /**
   * Checks current auth state; is access token loaded OR do we have a refresh token we can use
   *
   * check uses system API client verify given/current JWT
   *
   * If JWT is valid, it is stored into local storage alongside
   * loaded user.
   *
   * We're explicitly passing systemAPI to minimize plugin initialization complexity
   *
   * Function will throw null when user is unauthenticated
   */
  async handleState (): Promise<AuthInfo | null> {
    this.log.info('checking authentication')

    if (this[accessToken]) {
      this.log.info('access token found')

      const headers = { Authorization: `Bearer ${this[accessToken]}` }

      this.log.info('fetching authentication info from ' + oauth2InfoURL)

      return this.axios.get(oauth2InfoURL, { headers }).then(({ data }) => {
        this.log.info('data fetch form info endpoint', { oauth2InfoURL, headers, data })

        const authUser = new system.User({
          userID: data.sub,
          name: data.name,
          email: data.email,
          handle: data.username,
        })

        if (data.preferred_language) {
          authUser.meta.preferredLanguage = data.preferred_language || 'en'
        }

        this[user] = authUser

        this.bindListeners()
        return data
      }).catch((error) => {
        this.log.error('data fetch form info endpoint failed', { oauth2InfoURL, headers, error })
        // assume invalid JWT and remove it
        this[accessToken] = undefined
        throw new Error('Unauthenticated')
      })
    }

    const refreshToken = this.sessionStorage.getItem(storeKeyRefreshToken)
    if (refreshToken) {
      this.log.debug('refresh token found', { refreshToken })

      /**
       * Only exchange refresh token if this is the final state (see startFinalState function for more details)
       *
       * If this is a duplicated-session, an error will be thrown and authentication will be (probably)
       * restarted by the caller.
       */
      this.log.info('refreshing token', refreshToken)

      /**
       * Refresh token found in the storage,
       * let's use it to get new access token
       */
      return this.exchangeRefresh(refreshToken)
        .then(r => {
          this.bindListeners()
          return r
        })
    }

    throw new Error('Unauthenticated')
  }

  logout (): void {
    this.pruneStore()

    this.location.assign(Make({
      url: `${this.cortezaAuthURL}/logout`,
      query: { back: this.location.toString() },
    }))
  }

  /**
   * Starts new authentication flow
   *
   * It generates simple rand state to harden security and to
   * keep track of before-flow-start location of the user
   */
  startAuthenticationFlow (): void {
    this.log.debug('starting new authentication flow')

    this.cleanFlags()
    this.incFlowCounter()

    const state = Math.random().toString(36).substring(2)
    this.sessionStorage.setItem(`auth.state.${state}.location`, this.getRedirect(this.location.toString()))

    this.location.assign(Make({
      url: `${this.cortezaAuthURL}` + oauth2FlowURL,
      query: {
        // eslint-disable-next-line @typescript-eslint/camelcase
        redirect_uri: this.callbackURL,
        scope: oauth2Scope,
        state,
      },
    }))
  }

  getRedirect (url: string): string {
    const u = new URL(url)

    // In case someone started the flow on a callback route, default to the root
    // of the webapp.
    if (this.isCallback(u.pathname)) {
      u.pathname = ''
      u.search = ''
      u.hash = ''
    }

    return u.toString()
  }

  isCallback (url: string): boolean {
    return /\/auth\/callback$/.test(url)
  }

  /**
   * protects against too many tries when we try to auto-fix the "state does not match" error
   * by restarting the aut flow.
   */
  private incFlowCounter (): void {
    const aux = this.sessionStorage.getItem(storeKeyFlowStarted)
    if (aux === null) {
      this.sessionStorage.setItem(storeKeyFlowStarted, '1')
      return
    }

    const count = parseInt(this.sessionStorage.getItem(storeKeyFlowStarted) as string)
    if (count >= maxStartAttempts) {
      // Too many start attempts
      this.sessionStorage.removeItem(storeKeyFlowStarted)
      throw new Error('could not start authentication flow, too many attempts')
    }

    this.sessionStorage.setItem(
      storeKeyFlowStarted,
      (count + 1).toString(),
    )
  }

  /**
   * Exchanges authorization code for access and refresh tokens
   */
  private async exchangeCode (code = ''): Promise<AuthInfo> {
    return this.oauth2token({
      code: code,
      scope: oauth2Scope,
      // eslint-disable-next-line @typescript-eslint/camelcase
      redirect_uri: this.callbackURL,
    }).then((oa2tr) => this.procTokenResponse(oa2tr))
  }

  /**
   * Exchanges refresh token for new access and new refresh token
   *
   * After successful token exchange, we call response processing function
   * to update internals & stored values
   *
   * @param refreshToken
   */
  private async exchangeRefresh (refreshToken: string): Promise<AuthInfo | null> {
    /**
     * Finalize
     */
    this.completeFinalState()

    return this.oauth2token({
      // eslint-disable-next-line @typescript-eslint/camelcase
      refresh_token: refreshToken || '',
    }).then((oa2tr) => this.procTokenResponse(oa2tr))
      .catch((err) => {
        const { response: { data: { error = undefined } = {} } = {} } = err
        if (error === 'invalid_grant') {
          this.pruneStore()
          throw new Error('Unauthenticated')
        }
        throw err
      })
  }

  /**
   * Processes fetched token and stores it
   *
   * Access token is stored only to instance of this object
   * Refresh token is stored only to local store
   *
   * @param oa2tkn OAuth2 token response
   * @private
   */
  private procTokenResponse (oa2tkn: OAuth2TokenResponse): AuthInfo {
    this.log.debug('new token', oa2tkn)

    if (this.refreshTimeout) {
      window.clearTimeout(this.refreshTimeout)
    }

    const timeout = oa2tkn.expires_in * this.refreshFactor

    this.log.debug('setting up refresh timeout callback', {
      // eslint-disable-next-line @typescript-eslint/camelcase
      expires_in: oa2tkn.expires_in,
      timeout,
    })

    // Schedule next refresh
    this.refreshTimeout = window.setTimeout(async () => {
      this.log.debug('refreshing token')
      this.exchangeRefresh(oa2tkn.refresh_token)
        .catch((err) => {
          this.log.error('refresh token exchange failed', err)
          this.startAuthenticationFlow()
        })
    }, 1000 * timeout)

    this.sessionStorage.setItem(storeKeyRefreshToken, oa2tkn.refresh_token)

    const u = new system.User({
      userID: oa2tkn.sub,
      name: oa2tkn.name,
      handle: oa2tkn.handle,
      email: oa2tkn.email,
    })

    if (oa2tkn.preferred_language) {
      u.meta.preferredLanguage = oa2tkn.preferred_language
    }

    this[accessToken] = oa2tkn.access_token
    this[user] = u

    if (this.$emit) {
      this.$emit('auth-token-processed', {
        user: u,
        accessToken: this[accessToken],
      })
    }

    return {
      accessTokenFn: (): string | undefined => { return this[accessToken] },
      user: u,
    }
  }

  /**
   * oauth2token exchanges authorization code or refresh token for (new) access token
   *
   * @param payload
   * @private
   */
  private async oauth2token (payload: Record<string, string>): Promise<OAuth2TokenResponse> {
    const data = new URLSearchParams()

    this.log.debug('exchanging for token', payload)

    Object.entries(payload).forEach(([key, value]) => {
      data.set(key, value)
    })

    const config = {
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    }

    return this.axios.post(oauth2FlowURL, data, config).then(({ data }) => data)
  }

  private pruneStore (): void {
    this[accessToken] = undefined
    this[user] = undefined
    this.sessionStorage.clear()
  }

  get accessToken (): string | undefined {
    return this[accessToken]
  }

  get user (): system.User | undefined {
    return this[user]
  }
}

export default function (): PluginFunction<PluginOpts> {
  return function (Vue, opts): void {
    let {
      app = '',
      rootApp = false,
      cortezaAuthURL = '',
      callbackURL = '',
      verbose = undefined,
      refreshFactor = 0.75,
      entrypointURL = window.location.toString(),
      location = window.location,
      sessionStorage = window.sessionStorage,
      registerEventListener = window.addEventListener.bind(window),
    } = (opts || {}) as Partial<AuthCtor & { rootApp: boolean }>

    if (!cortezaAuthURL) {
      /**
       * cortezaAuthURL not explicitly set, try to auto-configure from properties set on window variable
       * (most likely through config.js)
       */

      // eslint-disable-next-line @typescript-eslint/ban-ts-ignore
      // @ts-ignore
      const { CortezaAPI = undefined, CortezaAuth = undefined } = window

      switch (true) {
        case !!CortezaAuth:
          /**
           * Corteza authentication endpoints location is set explicitly:
           */
          cortezaAuthURL = CortezaAuth
          break
        case !!CortezaAPI && /\/api$/.test(CortezaAPI):
          /**
           * Corteza API base-url is explicitly set and string ends with /api,
           * do a leap of faith and replace it with /auth, so that
           * corteza.example.tld/api becomes corteza.example.tld/auth
           */
          cortezaAuthURL = CortezaAPI.replace('/api', '/auth')
          break
        case !!CortezaAPI:
          /**
           * Corteza API base-url is explicitly set. Since it does not end with /api
           * we will assume api is served directly on root of that domain and we'll just append the /auth suffix
           * so that corteza.example.tld becomes corteza.example.tld/auth
           */
          cortezaAuthURL = CortezaAPI + '/auth'
          break
        default:
          throw new Error('failed to configure auth cortezaAuthURL')
      }
    }

    if (!callbackURL) {
      if (!app) {
        throw new Error('can not construct callbackURL; specify \'callbackURL\' or \'app\' property')
      }

      // eslint-disable-next-line @typescript-eslint/ban-ts-ignore
      // @ts-ignore
      const { CortezaWebapp = undefined } = window
      const callbackPath = 'auth/callback'

      if (CortezaWebapp) {
        // construct redirect URL fallback from configured corteza webapp
        callbackURL = Make({ url: `${CortezaWebapp}` })
      } else {
        // Try to get callbackURL from <base> tag's href value
        const baseTags = document.getElementsByTagName('base')
        if (baseTags.length === 1) {
          callbackURL = baseTags[0].href
        }

        if (!callbackURL) {
          // construct redirect URL fallback from current location
          // note: host contains port, hostname does not!
          const { protocol, host } = location
          callbackURL = `${protocol}//${host}`
        }
      }

      if (!rootApp) {
        callbackURL = callbackURL.replace(/\/$/, '') + `/${app}`
      }

      callbackURL = callbackURL.replace(/\/$/, '') + `/${callbackPath}`
    }

    if (verbose === undefined) {
      // enable debug (when not expl. disabled on localhost)
      verbose = location.hostname === 'localhost' ||
        window.location.search.includes('verboseAuth') ||
        !!window.localStorage.getItem('auth.verbose') ||
        !!window.sessionStorage.getItem('auth.verbose')
    }

    if (verbose) {
      console.debug({
        app,
        verbose,
        cortezaAuthURL,
        callbackURL,
        location,
        sessionStorage,
        entrypointURL,
        refreshFactor,
      })
    }

    Vue.prototype.$auth = new Auth({
      app,
      verbose,
      cortezaAuthURL,
      callbackURL,
      location,
      sessionStorage,
      entrypointURL,
      refreshFactor,
      registerEventListener,
    })
  }
}
