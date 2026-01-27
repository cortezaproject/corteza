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
export declare function ModuleFieldMaker(i: {
    kind?: string;
}): ModuleField;
export { Registry as ModuleFieldRegistry, ModuleField }
