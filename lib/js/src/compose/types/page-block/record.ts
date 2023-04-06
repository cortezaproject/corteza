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
  magnifyOption: string;
  recordSelectorDisplayOption: string;
  referenceField?: string;
  referenceModuleID?: string;
}

const defaults: Readonly<Options> = Object.freeze({
  fields: [],
  fieldConditions: [],
  magnifyOption: '',
  recordSelectorDisplayOption: 'sameTab',
  referenceField: '',
  referenceModuleID: undefined
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

    Apply(this.options, o, String, 'magnifyOption', 'recordSelectorDisplayOption', 'referenceField', 'referenceModuleID')

    if (o.fields) {
      this.options.fields = o.fields
    }

    if (o.fieldConditions) {
      this.options.fieldConditions = o.fieldConditions
    }
  }
}

Registry.set(kind, PageBlockRecord)
