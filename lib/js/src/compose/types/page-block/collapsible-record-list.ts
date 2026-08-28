import { PageBlock, PageBlockInput, Registry } from './base'
import { Apply, CortezaID, NoID } from '../../../cast'
import { Compose as ComposeAPI } from '../../../api-clients'
import { Module } from '../module'

const kind = 'CollapsibleRecordList'

interface FieldCondition {
  field: string
  condition: string
  clearOnHide?: boolean
}

interface Options {
  // Module and basic list configuration
  moduleID: string
  prefilter: string
  presort: string
  perPage: number
  refField?: string

  // Header section - title
  titleExpression: string
  titleAlignment: 'left' | 'center' | 'right'

  // Header section - subtitle
  subtitleExpression: string
  subtitleAlignment: 'left' | 'center' | 'right'
  subtitleShowWhenCollapsed: boolean

  // Body section - full-width field
  bodyField: string

  // Other fields section
  fields: unknown[]
  fieldConditions: FieldCondition[]
  clearConditionalFieldsOnHide: boolean
  otherFieldsPosition: 'above' | 'below' | 'default'
  horizontalFieldLayoutEnabled: boolean
  recordFieldLayoutOption: string
  searchableFields: string[]

  // Collapse state per record (stored locally)
  defaultCollapsed: boolean

  // Display options
  hideHeader: boolean
  hidePaging: boolean
  hideSorting: boolean
  hideSearch: boolean
  hideFiltering: boolean

  // Navigation
  enableRecordPageNavigation: boolean
  fullPageNavigation: boolean
  showRecordPerPageOption: boolean

  // Actions
  hideRecordViewButton: boolean
  hideRecordEditButton: boolean
  hideRecordDeleteButton: boolean

  // Record display
  recordDisplayOption: string

  // Show total count
  showTotalCount: boolean
}

const defaults: Readonly<Options> = Object.freeze({
  // Module and basic list configuration
  moduleID: NoID,
  prefilter: '',
  presort: 'createdAt DESC',
  perPage: 10,
  refField: undefined,

  // Header section - title
  titleExpression: '',
  titleAlignment: 'left',

  // Header section - subtitle
  subtitleExpression: '',
  subtitleAlignment: 'left',
  subtitleShowWhenCollapsed: false,

  // Body section
  bodyField: '',

  // Other fields section
  fields: [],
  fieldConditions: [],
  clearConditionalFieldsOnHide: false,
  otherFieldsPosition: 'default',
  horizontalFieldLayoutEnabled: false,
  recordFieldLayoutOption: 'default',
  searchableFields: [],

  // Collapse state
  defaultCollapsed: false,

  // Display options
  hideHeader: false,
  hidePaging: false,
  hideSorting: false,
  hideSearch: false,
  hideFiltering: false,

  // Navigation
  enableRecordPageNavigation: true,
  fullPageNavigation: false,
  showRecordPerPageOption: false,

  // Actions
  hideRecordViewButton: false,
  hideRecordEditButton: false,
  hideRecordDeleteButton: false,

  // Record display
  recordDisplayOption: 'modal',

  // Show total count
  showTotalCount: true,
})

export class PageBlockCollapsibleRecordList extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor (i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return

    Apply(this.options, o, CortezaID, 'moduleID')

    Apply(this.options, o, String,
      'titleExpression',
      'titleAlignment',
      'subtitleExpression',
      'subtitleAlignment',
      'bodyField',
      'otherFieldsPosition',
      'recordFieldLayoutOption',
      'prefilter',
      'presort',
      'refField',
      'recordDisplayOption',
    )

    Apply(this.options, o, Number, 'perPage')

    Apply(this.options, o, Boolean,
      'subtitleShowWhenCollapsed',
      'horizontalFieldLayoutEnabled',
      'clearConditionalFieldsOnHide',
      'defaultCollapsed',
      'hideHeader',
      'hidePaging',
      'hideSorting',
      'hideSearch',
      'hideFiltering',
      'enableRecordPageNavigation',
      'fullPageNavigation',
      'showRecordPerPageOption',
      'hideRecordViewButton',
      'hideRecordEditButton',
      'hideRecordDeleteButton',
      'showTotalCount',
    )

    if (o.fields) {
      this.options.fields = o.fields
    }

    if (o.fieldConditions) {
      this.options.fieldConditions = o.fieldConditions
    }

    if (o.searchableFields) {
      this.options.searchableFields = o.searchableFields
    }
  }

  async fetch (api: ComposeAPI, recordListModule: Module, filter: {[_: string]: unknown}): Promise<object> {
    if (recordListModule.moduleID !== this.options.moduleID) {
      throw Error('Module incompatible, module mismatch')
    }

    filter.moduleID = this.options.moduleID
    filter.namespaceID = recordListModule.namespaceID
    filter.sort = this.options.presort

    if (this.options.prefilter) {
      filter.filter = this.options.prefilter
    }

    return api
      .recordList(filter)
      .then(r => {
        const { set: records, filter } = r as { filter: object; set: object[] }
        return { records, filter }
      })
  }
}

Registry.set(kind, PageBlockCollapsibleRecordList)
