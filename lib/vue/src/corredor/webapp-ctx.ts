import { corredor } from '@cortezaproject/corteza-js'
import pino from 'pino'

/**
 * Bare-minimum webapp helper
 */
export default class WebappCtx extends corredor.Ctx {
  constructor (args: corredor.BaseArgs) {
    super(args, pino(), {})
  }

  /**
   * Clones context and uses new arguments
   */
  withArgs (args: corredor.BaseArgs): WebappCtx {
    return new WebappCtx(args)
  }
}
