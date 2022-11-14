/* eslint-disable @typescript-eslint/ban-ts-ignore */
import { compose, apiClients, corredor } from '@cortezaproject/corteza-js'
import ComposeUIHelper from './compose-ui'
import pino from 'pino'

interface Vue {
  $SystemAPI: apiClients.System;
  $ComposeAPI: apiClients.Compose;

  $store: { getters: { [_: string]: Array<compose.Page> } };
  $emit: unknown;
  $router: { push: unknown };
}

/**
 * Extends corredor exec context with compose UI helper
 */
export default class ComposeCtx extends corredor.Ctx {
  protected emitter: unknown
  protected routePusher: unknown
  protected pages: Array<compose.Page> = []

  protected composeUI: ComposeUIHelper

  protected vue: Vue

  // @todo remove ts-ignore flags
  constructor (args: corredor.BaseArgs, vue: Vue) {
    // @ts-ignore
    super(args, pino(), {})

    this.vue = vue

    // @ts-ignore
    this.systemAPI = vue.$SystemAPI
    // @ts-ignore
    this.composeAPI = vue.$ComposeAPI

    this.composeUI = new ComposeUIHelper({
      ...this.args,

      pages: vue.$store.getters['page/set'],
      // @ts-ignore
      emitter: (name, params): void => vue.$emit(name, params),
      // @ts-ignore
      routePusher: (params): void => vue.$router.push(params),
    })

    if (!this.config) {
      this.config = {}
    }

    this.config.frontend = {
      baseURL: `${window.location.protocol}://${window.location.host}`,
    }
  }

  /**
   * Clones context and uses new arguments
   */
  withArgs (args: corredor.BaseArgs): ComposeCtx {
    Object.assign(args, this.args)
    return new ComposeCtx(args, this.vue)
  }

  get ComposeUI (): ComposeUIHelper {
    return this.composeUI
  }
}
