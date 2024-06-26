import { PageBlock, PageBlockInput, Registry } from './base'
import { Apply } from '../../../cast'

const kind = 'Record'

interface FieldCondition {
  field: string;
  condition: string;
}

interface Options {
  fields: unknown[];
  fieldConditions: FieldCondition[];
  recordSelectorShowAddRecordButton: boolean;
  magnifyOption: string;
  recordSelectorDisplayOption: string;
  recordSelectorAddRecordDisplayOption: string;
  referenceField?: string;
  referenceModuleID?: string;
  inlineRecordEditEnabled: boolean;
  horizontalFieldLayoutEnabled: boolean;
  recordFieldLayoutOption: string;
}

const defaults: Readonly<Options> = Object.freeze({
  fields: [],
  fieldConditions: [],
  recordSelectorShowAddRecordButton: false,
  magnifyOption: '',
  recordSelectorDisplayOption: 'sameTab',
  recordSelectorAddRecordDisplayOption: 'sameTab',
  referenceField: '',
  referenceModuleID: undefined,
  inlineRecordEditEnabled: false,
  horizontalFieldLayoutEnabled: false,
  recordFieldLayoutOption: 'default',
})

export class PageBlockRecord extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor (i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return

    Apply(this.options, o, String, 'magnifyOption', 'recordSelectorDisplayOption', 'recordSelectorAddRecordDisplayOption', 'referenceField', 'referenceModuleID', 'recordFieldLayoutOption')
    Apply(this.options, o, Boolean, 'recordSelectorShowAddRecordButton', 'inlineRecordEditEnabled', 'horizontalFieldLayoutEnabled')

    if (o.fields) {
      this.options.fields = o.fields
    }

    if (o.fieldConditions) {
      this.options.fieldConditions = o.fieldConditions
    }
  }
}

Registry.set(kind, PageBlockRecord)
