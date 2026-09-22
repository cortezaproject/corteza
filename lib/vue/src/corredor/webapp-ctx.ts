// Ctx comes from its own module rather than the package root: it is extended at
// module scope, which the package root cannot be without every consumer of this
// library evaluating it.
import { Ctx } from '@planetcrust/human-js/src/corredor/ctx'
import type { BaseArgs } from '@planetcrust/human-js/src/corredor/shared'
import type { apiClients } from '@planetcrust/human-js'
import { consoleLogger } from './logger'

export interface WebappCtxServices {
  systemAPI?: apiClients.System
  composeAPI?: apiClients.Compose
}

/**
 * Bare-minimum Corredor exec context for scripts running outside Compose
 */
export class WebappCtx extends Ctx {
  protected services: WebappCtxServices

  constructor(args: BaseArgs, services: WebappCtxServices = {}) {
    super(args, consoleLogger(), {
      systemAPI: services.systemAPI,
      composeAPI: services.composeAPI,
      config: { frontend: { baseURL: window.location.origin } },
    })

    this.services = services
  }

  /**
   * Clones the context and uses new arguments
   */
  withArgs(args: BaseArgs): WebappCtx {
    Object.assign(args, this.args)
    return new WebappCtx(args, this.services)
  }
}
