import { ModuleField, Registry } from './base'
export { ModuleFieldBool } from './bool'
export { ModuleFieldDateTime } from './datetime'
export { ModuleFieldEmail } from './email'
export { ModuleFieldFile } from './file'
export { ModuleFieldSelect } from './select'
export { ModuleFieldNumber } from './number'
export { ModuleFieldRecord } from './record'
export { ModuleFieldString } from './string'
export { ModuleFieldUrl } from './url'
export { ModuleFieldUser } from './user'
export { ModuleFieldGeometry } from './geometry'

export function ModuleFieldMaker(i: { kind?: string }): ModuleField {
  if (!i.kind) {
    return new ModuleField(i)
  }

  if (!Registry.has(i.kind)) {
    // Unknown field kind — return a plain ModuleField so the module can still load
    return new ModuleField(i)
  }

  return new (Registry.get(i.kind) as typeof ModuleField)(i)
}

export { Registry as ModuleFieldRegistry, ModuleField }
