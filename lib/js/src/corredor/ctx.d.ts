import * as apiClients from '../api-clients'
import { SystemHelper, ComposeHelper } from './helpers'
import { Logger } from 'pino'
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
    config?: Config;
    systemAPI?: apiClients.System;
    composeAPI?: apiClients.Compose;
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
export declare class Ctx {
  protected args: BaseArgs
  protected config?: Config
  protected logger: Logger
  protected systemAPI?: apiClients.System
  protected composeAPI?: apiClients.Compose
  constructor(args: BaseArgs, logger: Logger, a?: CtxInitArgs);
  /**
     * Alias for log, to make developer's life easier <3
     */
  get console(): Logger;
  /**
     * Alias for log, to make developer's life easier <3
     */
  get log(): Logger;
  /**
     * Returns promise with the current user (if authToken argument was given)
     *
     * This is a temporary solution that decodes the userID from the access token (JWT)
     * and fetches the user info
     *
     * @returns {Promise<User>}
     */
  get $authUser(): Promise<User>;
  /**
     * Configures and returns system API client
     */
  get SystemAPI(): apiClients.System;
  /**
     * Configures and returns compose API client
     */
  get ComposeAPI(): apiClients.Compose;
  /**
     * Configures and returns system helper
     */
  get System(): SystemHelper;
  /**
     * Configures and returns compose helper
     */
  get Compose(): ComposeHelper;
  /**
     *
     */
  get frontendBaseURL(): string | undefined;
}
export {}
