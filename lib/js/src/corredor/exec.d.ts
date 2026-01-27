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
export declare function Exec(script: ExecutableScript, args: BaseArgs, ctx: Ctx): Promise<Results>;
export {}
