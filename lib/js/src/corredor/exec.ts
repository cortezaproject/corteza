import { BaseArgs } from './shared'
import { Ctx } from './ctx'

export interface ScriptExecFn {
    (args: BaseArgs, ctx?: Ctx): unknown;
}

export interface ExecutableScript {
  exec: ScriptExecFn;
}

interface Results {
  result?: unknown;
}

/**
 * Script executor
 *
 * @param script - Script to be executed
 * @param args - Arguments for the script
 * @param ctx - Exec context (exec function's 2nd param)
 */
export async function Exec (script: ExecutableScript, args: BaseArgs, ctx: Ctx): Promise<Results> {
  try {
    // Wrap exec() with Promise.resolve - we do not know if function is async or not.
    return Promise.resolve(script.exec(args, ctx)).then((rval: unknown): object => {
      let result = {}

      if (rval === false) {
        // Abort when returning false!
        throw new Error('Aborted')
      }

      if (typeof rval === 'object' && rval && rval.constructor.name === 'Object') {
        // Expand returned values into result if function returned a plain javascript object
        result = { ...rval }
      } else if (rval !== undefined) {
        // If anything usable was returned, stack it under 'result' property
        result = { result: rval }
      }

      // Wrap returning value
      return result
    })
  } catch (e) {
    return Promise.reject(e)
  }
}
