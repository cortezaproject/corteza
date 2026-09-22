// Ctx comes from its own module rather than the package root: it is extended at
// module scope, which the package root cannot be without every consumer of this
// library evaluating it.
import { Ctx } from '@planetcrust/human-js/src/corredor/ctx'
import type { BaseArgs } from '@planetcrust/human-js/src/corredor/shared'
import type { apiClients, compose } from '@planetcrust/human-js'
import { ComposeUIHelper, type ComposeUIContext, type ComposeUIToast } from './compose-ui'
import { consoleLogger } from './logger'

export interface ComposeCtxServices {
  systemAPI?: apiClients.System
  composeAPI?: apiClients.Compose
  pages?: () => compose.Page[]
  toast?: ComposeUIToast
  routePusher?: (_to: object) => void
}

/**
 * Corredor exec context for scripts running inside Compose
 */
export class ComposeCtx extends Ctx {
  protected services: ComposeCtxServices
  protected composeUI: ComposeUIHelper

  constructor(args: BaseArgs, services: ComposeCtxServices = {}) {
    super(args, consoleLogger(), {
      systemAPI: services.systemAPI,
      composeAPI: services.composeAPI,
      config: { frontend: { baseURL: window.location.origin } },
    })

    this.services = services

    const resources = args as unknown as ComposeUIContext

    this.composeUI = new ComposeUIHelper({
      $namespace: resources.$namespace,
      $module: resources.$module,
      $record: resources.$record,
      pages: services.pages,
      toast: services.toast,
      routePusher: services.routePusher,
    })
  }

  /**
   * Clones the context and uses new arguments
   */
  withArgs(args: BaseArgs): ComposeCtx {
    Object.assign(args, this.args)
    return new ComposeCtx(args, this.services)
  }

  get ComposeUI(): ComposeUIHelper {
    return this.composeUI
  }
}
