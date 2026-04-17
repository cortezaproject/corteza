import { User, Role, Application, SinkResponse, SinkRequest } from '../system'
import { Module, Page, Namespace, Record } from '../compose'
import { Caster, GenericCaster, GenericCasterFreezer } from './shared'

interface RecordCasterCaller {
  $module: Module
}

/**
 * Record type caster
 *
 * Record arg is a bit special, it takes 2 params (record itself + record's module)
 */
function recordCaster(this: RecordCasterCaller, val: unknown): Record | undefined {
  if (val) {
    try {
      return new Record(this.$module, val as object)
    } catch (e) {
      console.error(e)
    }
  }

  return undefined
}

function recordCasterFreezer(this: RecordCasterCaller, val: unknown): Readonly<Record> | undefined {
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
 * HumanTypes map helps ExecArgs class with translation of (special) arguments
 * to their respected types
 *
 * There's noe need to set/define casters for old* arguments,
 * It's auto-magically done by Args class
 */
export const HumanTypes: Caster = new Map()

HumanTypes.set('authUser', GenericCasterFreezer(User))
HumanTypes.set('invoker', GenericCasterFreezer(User))
HumanTypes.set('module', GenericCaster(Module))
HumanTypes.set('oldModule', GenericCasterFreezer(Module))
HumanTypes.set('page', GenericCaster(Page))
HumanTypes.set('oldPage', GenericCasterFreezer(Page))
HumanTypes.set('namespace', GenericCaster(Namespace))
HumanTypes.set('oldNamespace', GenericCasterFreezer(Namespace))
HumanTypes.set('application', GenericCaster(Application))
HumanTypes.set('oldApplication', GenericCasterFreezer(Application))
HumanTypes.set('user', GenericCaster(User))
HumanTypes.set('oldUser', GenericCasterFreezer(User))
HumanTypes.set('role', GenericCaster(Role))
HumanTypes.set('oldRole', GenericCasterFreezer(Role))
HumanTypes.set('record', recordCaster)
HumanTypes.set('oldRecord', recordCasterFreezer)
HumanTypes.set('request', GenericCasterFreezer(SinkRequest))
HumanTypes.set('response', GenericCaster(SinkResponse))
