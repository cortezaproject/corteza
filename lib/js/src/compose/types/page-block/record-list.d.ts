import { PageBlock, PageBlockInput } from './base'
import { Compose as ComposeAPI } from '../../../api-clients'
import { Module } from '../module'
import { Button } from './types'
interface FilterPreset {
  name: string
  filter: unknown[]
  roles: string[]
}
declare enum SummaryMetric {
  Min = 'min',
  Max = 'max',
  Avg = 'avg',
  Sum = 'sum',
  EmptyCount = 'emptyCount',
  NotEmptyCount = 'notEmptyCount',
  UniqueCount = 'uniqueCount',
  Earliest = 'earliest',
  Latest = 'latest',
}
interface Summary {
  label: string
  field: string[]
  metric: SummaryMetric
  roles: string[]
}
export interface Options {
  moduleID: string
  prefilter: string
  presort: string
  fields: unknown[]
  hideHeader: boolean
  hideAddButton: boolean
  hideImportButton: boolean
  hideConfigureFieldsButton: boolean
  hideSearch: boolean
  hidePaging: boolean
  hideSorting: boolean
  hideFiltering: boolean
  hideRecordReminderButton: boolean
  hideRecordCloneButton: boolean
  hideRecordEditButton: boolean
  hideRecordViewButton: boolean
  hideRecordPermissionsButton: boolean
  enableRecordPageNavigation: boolean
  allowExport: boolean
  perPage: number
  recordDisplayOption: string
  recordSelectorDisplayOption: string
  addRecordDisplayOption: string
  magnifyOption: string
  fullPageNavigation: boolean
  showTotalCount: boolean
  showDeletedRecordsOption: boolean
  customFilterPresets: boolean
  refreshRate: number
  showRefresh: boolean
  editable: boolean
  draggable?: boolean
  positionField?: string
  refField?: string
  editFields?: unknown[]
  linkToParent: boolean
  openInNewTab: boolean
  selectable: boolean
  selectMode: 'multi' | 'single' | 'range'
  selectionButtons: Array<Button>
  bulkRecordEditEnabled: boolean
  inlineRecordEditEnabled: boolean
  inlineRecordEditAllowAddField: boolean
  inlineValueFiltering: boolean
  filterPresets: FilterPreset[]
  showRecordPerPageOption: boolean
  openRecordInEditMode: boolean
  customSummaries: boolean
  summaries: Summary[]
  textStyles: {
    wrappedFields: Array<string>
  }
}
export declare class PageBlockRecordList extends PageBlock {
  readonly kind = 'RecordList'
  options: Options
  constructor(i?: PageBlockInput)
  applyOptions(o?: Partial<Options>): void
  fetch(
    api: ComposeAPI,
    recordListModule: Module,
    filter: {
      [_: string]: unknown
    },
  ): Promise<object>
}
export {}
