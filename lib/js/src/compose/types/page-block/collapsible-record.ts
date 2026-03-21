import { PageBlock, PageBlockInput, Registry } from './base'
import { Apply } from '../../../cast'

const kind = 'CollapsibleRecord'

interface FieldCondition {
  field: string
  condition: string
  clearOnHide?: boolean
}

interface Options {
  // Header section - title
  titleExpression: string
  titleAlignment: 'left' | 'center' | 'right'

  // Header section - subtitle
  subtitleExpression: string
  subtitleAlignment: 'left' | 'center' | 'right'
  subtitleShowWhenCollapsed: boolean

  // Body section - full-width field
  bodyField: string

  // Inline editing for body field
  inlineEditEnabled: boolean

  // Other fields section
  fields: unknown[]
  fieldConditions: FieldCondition[]
  clearConditionalFieldsOnHide: boolean
  otherFieldsPosition: 'above' | 'below' | 'default'
  horizontalFieldLayoutEnabled: boolean
  recordFieldLayoutOption: string

  // Collapse state
  defaultCollapsed: boolean
}

const defaults: Readonly<Options> = Object.freeze({
  titleExpression: '',
  titleAlignment: 'left',
  subtitleExpression: '',
  subtitleAlignment: 'left',
  subtitleShowWhenCollapsed: false,
  bodyField: '',
  inlineEditEnabled: false,
  fields: [],
  fieldConditions: [],
  clearConditionalFieldsOnHide: false,
  otherFieldsPosition: 'default',
  horizontalFieldLayoutEnabled: false,
  recordFieldLayoutOption: 'default',
  defaultCollapsed: false,
})

export class PageBlockCollapsibleRecord extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor (i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return

    Apply(this.options, o, String, 'titleExpression', 'titleAlignment', 'subtitleExpression', 'subtitleAlignment', 'bodyField', 'otherFieldsPosition', 'recordFieldLayoutOption')
    Apply(this.options, o, Boolean, 'subtitleShowWhenCollapsed', 'inlineEditEnabled', 'horizontalFieldLayoutEnabled', 'clearConditionalFieldsOnHide', 'defaultCollapsed')

    if (o.fields) {
      this.options.fields = o.fields
    }

    if (o.fieldConditions) {
      this.options.fieldConditions = o.fieldConditions
    }
  }
}

Registry.set(kind, PageBlockCollapsibleRecord)
