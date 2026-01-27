import { Caster } from './shared'
/**
 * Handles arguments, passed to the script
 *
 * By convention variables holding "current" resources are prefixed with dollar ($) sign.
 * For example, before/after triggers for record will call registered scripts with $record, $module
 * and $namespace, holding current record, it's module and namespace.
 *
 * All these variables are casted (if passed as an argument) to proper types ($record => Record, $module => Module, ...)
 */
export declare class Args {
  constructor(args: {
        [_: string]: unknown;
    }, caster?: Caster);
}
/**
 * Handles arguments, passed to the script but preserves references to the original objects
 *
 * By convention variables holding "current" resources are prefixed with dollar ($) sign.
 * For example, before/after triggers for record will call registered scripts with $record, $module
 * and $namespace, holding current record, it's module and namespace.
 *
 * These variables are not additionally casted, since in order to preserve references they should
 * already be in the correct type.
 */
export declare class ArgsProxy {
  constructor(args: {
        [_: string]: unknown;
    }, caster?: Caster);
}
