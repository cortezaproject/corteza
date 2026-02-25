import { ModuleField, Options } from './base'
interface RecordOptions extends Options {
  moduleID: string
  labelField: string
  recordLabelField: string
  queryFields: Array<string>
  selectType: string
  multiDelimiter: string
  isUniqueMultiValue: boolean
  prefilter?: string
}
export declare class ModuleFieldRecord extends ModuleField {
  readonly kind = 'Record'
  options: RecordOptions
  constructor(i?: Partial<ModuleFieldRecord>)
  applyOptions(o?: Partial<RecordOptions>): void
}
export {}
