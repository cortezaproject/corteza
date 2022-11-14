// @todo option to allow multiple entries
// @todo option to allow duplicates
import { ModuleField, Registry, Options, defaultOptions } from './base'
import { Apply, CortezaID, NoID } from '../../../cast'

const kind = 'Record'

interface RecordOptions extends Options {
  moduleID: string;
  labelField: string;
  recordLabelField: string;
  queryFields: Array<string>;
  selectType: string;
  multiDelimiter: string;
  prefilter?: string;
}

const defaults = (): Readonly<RecordOptions> => Object.freeze({
  ...defaultOptions(),
  moduleID: NoID,
  labelField: '',
  recordLabelField: '',
  queryFields: [],
  selectType: '',
  multiDelimiter: '\n',
  prefilter: undefined,
})

export class ModuleFieldRecord extends ModuleField {
  readonly kind = kind

  options: RecordOptions = { ...defaults() }

  constructor (i?: Partial<ModuleFieldRecord>) {
    super(i)
    this.applyOptions(i?.options)
  }

  applyOptions (o?: Partial<RecordOptions>): void {
    if (!o) return
    super.applyOptions(o)

    Apply(this.options, o, CortezaID, 'moduleID')
    Apply(this.options, o, String, 'labelField', 'recordLabelField', 'selectType', 'multiDelimiter', 'prefilter')
    Apply(this.options, o, (o) => {
      if (!o) {
        return []
      }
      if (!Array.isArray(o)) {
        return [o]
      }
      return o
    }, 'queryFields')
  }
}

Registry.set(kind, ModuleFieldRecord)
