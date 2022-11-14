/* eslint-disable @typescript-eslint/ban-ts-ignore */

import { User, Role, Application, SinkResponse, SinkRequest } from '../system'
import { Module, Page, Namespace, Record } from '../compose'
import { Caster, GenericCaster, GenericCasterFreezer } from './shared'

interface RecordCasterCaller {
  $module: Module;
}

/**
 * Record type caster
 *
 * Record arg is a bit special, it takes 2 params (record itself + record's module)
 */
function recordCaster (this: RecordCasterCaller, val: unknown): Record|undefined {
  if (val) {
    try {
      return new Record(this.$module, val as object)
    } catch (e) {
      console.error(e)
    }
  }

  return undefined
}

function recordCasterFreezer (this: RecordCasterCaller, val: unknown): Readonly<Record>|undefined {
  if (val) {
    try {
      return Object.freeze(new Record(this.$module, val as object))
    } catch (e) {
      console.error(e)
    }
  }

  return undefined
}

/**
 * CortezaTypes map helps ExecArgs class with translation of (special) arguments
 * to their respected types
 *
 * There's noe need to set/define casters for old* arguments,
 * It's auto-magically done by Args class
 */
export const CortezaTypes: Caster = new Map()

CortezaTypes.set('authUser', GenericCasterFreezer(User))
CortezaTypes.set('invoker', GenericCasterFreezer(User))
CortezaTypes.set('module', GenericCaster(Module))
CortezaTypes.set('oldModule', GenericCasterFreezer(Module))
CortezaTypes.set('page', GenericCaster(Page))
CortezaTypes.set('oldPage', GenericCasterFreezer(Page))
CortezaTypes.set('namespace', GenericCaster(Namespace))
CortezaTypes.set('oldNamespace', GenericCasterFreezer(Namespace))
CortezaTypes.set('application', GenericCaster(Application))
CortezaTypes.set('oldApplication', GenericCasterFreezer(Application))
CortezaTypes.set('user', GenericCaster(User))
CortezaTypes.set('oldUser', GenericCasterFreezer(User))
CortezaTypes.set('role', GenericCaster(Role))
CortezaTypes.set('oldRole', GenericCasterFreezer(Role))
CortezaTypes.set('record', recordCaster)
CortezaTypes.set('oldRecord', recordCasterFreezer)
CortezaTypes.set('request', GenericCasterFreezer(SinkRequest))
CortezaTypes.set('response', GenericCaster(SinkResponse))
