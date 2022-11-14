import * as apiClients from '../api-clients'
import { SystemHelper, ComposeHelper } from './helpers'
import { BaseLogger } from 'pino'
import { BaseArgs } from './shared'
import { User } from '../system'

export interface ConfigCServers {
  system?: ConfigServer;
  compose?: ConfigServer;
}

export interface ConfigServer {
  apiBaseURL?: string;
}

export interface ConfigFrontend {
  baseURL?: string;
}

export interface Config {
  cServers?: ConfigCServers;
  frontend?: ConfigFrontend;
}

interface CtxInitArgs {
  config?:
    Config;

  systemAPI?:
    apiClients.System;

  composeAPI?:
    apiClients.Compose;
}

/**
 * Handles script execution context
 *
 * Context accepts pre-assembled *API props or it construct them fly from passed config
 *
 * Naming convention for properties:
 *  - Corteza classes, high-level helpers, API clients are upper cased
 *  - low-level helpers are lower cased
 *  - simple scalar are lower cased
 *  - $authUser is the only one prefixed with the dollar sign for historical reasons
 */
export class Ctx {
  protected args: BaseArgs;
  protected config?: Config;

  protected logger: BaseLogger;

  protected systemAPI?:
    apiClients.System;

  protected composeAPI?:
    apiClients.Compose;

  constructor (args: BaseArgs, logger: BaseLogger, a?: CtxInitArgs) {
    this.args = args
    this.logger = logger

    if (a) {
      Object.assign(this, a)
    }
  }

  /**
   * Alias for log, to make developer's life easier <3
   */
  get console (): BaseLogger {
    return this.logger
  }

  /**
   * Alias for log, to make developer's life easier <3
   */
  get log (): BaseLogger {
    return this.logger
  }

  /**
   * Returns promise with the current user (if authToken argument was given)
   *
   * This is a temporary solution that decodes the userID from the access token (JWT)
   * and fetches the user info
   *
   * @returns {Promise<User>}
   */
  get $authUser (): Promise<User> {
    const [, payload] = this.args.authToken.split('.')
    // eslint-disable-next-line node/no-deprecated-api
    const buf = new Buffer(payload, 'base64')
    const { sub: userID } = JSON.parse(buf.toString('ascii'))

    return this.SystemAPI.userRead({ userID }).then(r => new User(r))
  }

  /**
   * Configures and returns system API client
   */
  get SystemAPI (): apiClients.System {
    if (!this.systemAPI) {
      if (!this.config?.cServers?.system) {
        throw new Error('configuration for corteza system server missing')
      }

      this.systemAPI = new apiClients.System({
        baseURL: this.config.cServers.system.apiBaseURL,
        accessTokenFn: (): string => this.args.authToken,
      })
    }

    return this.systemAPI
  }

  /**
   * Configures and returns compose API client
   */
  get ComposeAPI (): apiClients.Compose {
    if (!this.composeAPI) {
      if (!this.config?.cServers?.compose) {
        throw new Error('configuration for corteza compose server missing')
      }

      this.composeAPI = new apiClients.Compose({
        baseURL: this.config.cServers.compose.apiBaseURL,
        accessTokenFn: (): string => this.args.authToken,
      })
    }

    return this.composeAPI
  }

  /**
   * Configures and returns system helper
   */
  get System (): SystemHelper {
    return new SystemHelper({ SystemAPI: this.SystemAPI, ...this.args })
  }

  /**
   * Configures and returns compose helper
   */
  get Compose (): ComposeHelper {
    return new ComposeHelper({ ComposeAPI: this.ComposeAPI, ...this.args })
  }

  /**
   *
   */
  get frontendBaseURL (): string|undefined {
    return this.config?.frontend?.baseURL
  }
}
