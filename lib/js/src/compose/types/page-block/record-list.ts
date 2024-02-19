import { PageBlock, PageBlockInput, Registry } from './base'
import { Apply, CortezaID, NoID } from '../../../cast'
import { Compose as ComposeAPI } from '../../../api-clients'
import { Module } from '../module'
import { Button } from './types'

const kind = 'RecordList'

interface FilterPreset {
  name: string;
  filter: unknown[];
  roles: string[];
}

export interface Options {
  moduleID: string;
  prefilter: string;
  presort: string;
  fields: unknown[];
  hideHeader: boolean;
  hideAddButton: boolean;
  hideImportButton: boolean;
  hideConfigureFieldsButton: boolean;
  hideSearch: boolean;
  hidePaging: boolean;
  hideSorting: boolean;
  hideFiltering: boolean;
  hideRecordReminderButton: boolean;
  hideRecordCloneButton: boolean;
  hideRecordEditButton: boolean;
  hideRecordViewButton: boolean;
  hideRecordPermissionsButton: boolean;
  enableRecordPageNavigation: boolean;
  allowExport: boolean;
  perPage: number;
  recordDisplayOption: string;
  recordSelectorDisplayOption: string;
  magnifyOption: string;

  fullPageNavigation: boolean;
  showTotalCount: boolean;
  showDeletedRecordsOption: boolean;
  customFilterPresets: boolean;
  refreshRate: number;
  showRefresh: boolean;

  // Record-lines
  editable: boolean;
  draggable?: boolean;
  positionField?: string;
  refField?: string; // When adding a new record, prefill refField value with parent record ID
  editFields?: unknown[];
  linkToParent: boolean; // Legacy

  // Should records be opened in a new tab
  // legacy field that has been removed but we keep it for backwards compatibility
  openInNewTab: boolean;

  // Are table rows selectable
  selectable: boolean;
  selectMode: 'multi' | 'single' | 'range';

  // Ordered list of buttons to display in the block
  selectionButtons: Array<Button>;

  bulkRecordEditEnabled: boolean;
  inlineRecordEditEnabled: boolean;
  filterPresets: FilterPreset[];
  showRecordPerPageOption: boolean;
  openRecordInEditMode: boolean;
}

const defaults: Readonly<Options> = Object.freeze({
  moduleID: NoID,
  prefilter: '',
  presort: 'createdAt DESC',
  fields: [],
  hideHeader: false,
  hideAddButton: false,
  hideImportButton: false,
  hideConfigureFieldsButton: true,
  hideSearch: false,
  hidePaging: false,
  hideSorting: false,
  hideFiltering: false,
  hideRecordReminderButton: false,
  hideRecordCloneButton: false,
  hideRecordEditButton: false,
  hideRecordViewButton: false,
  hideRecordPermissionsButton: false,
  enableRecordPageNavigation: false,
  allowExport: false,
  perPage: 20,
  recordDisplayOption: 'sameTab',
  recordSelectorDisplayOption: 'sameTab',
  magnifyOption: '',

  fullPageNavigation: false,
  showTotalCount: false,
  showDeletedRecordsOption: false,
  customFilterPresets: false,

  editable: false,
  draggable: false,
  positionField: undefined,
  refField: undefined,
  editFields: [],

  linkToParent: false,

  openInNewTab: false,

  selectable: true,
  selectMode: 'multi',

  selectionButtons: [],
  refreshRate: 0,
  showRefresh: false,

  bulkRecordEditEnabled: true,
  inlineRecordEditEnabled: false,
  filterPresets: [],
  showRecordPerPageOption: false,
  openRecordInEditMode: false,
})

export class PageBlockRecordList extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor (i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return

    Apply(this.options, o, CortezaID, 'moduleID')
    Apply(this.options, o, String, 'prefilter', 'presort', 'selectMode', 'positionField', 'refField', 'recordDisplayOption', 'magnifyOption', 'recordSelectorDisplayOption')
    Apply(this.options, o, Number, 'perPage', 'refreshRate')

    if (o.fields) {
      this.options.fields = o.fields
    }

    if (o.filterPresets) {
      this.options.filterPresets = o.filterPresets
    }

    if (o.editFields) {
      this.options.editFields = o.editFields
    }

    if (o.openInNewTab) {
      this.options.recordDisplayOption = 'newTab'
    }

    Apply(this.options, o, Boolean,
      'hideHeader',
      'hideAddButton',
      'hideImportButton',
      'hideConfigureFieldsButton',
      'hideSearch',
      'hidePaging',
      'hideFiltering',
      'fullPageNavigation',
      'showTotalCount',
      'showDeletedRecordsOption',
      'customFilterPresets',
      'hideSorting',
      'allowExport',
      'selectable',
      'hideRecordReminderButton',
      'hideRecordCloneButton',
      'hideRecordEditButton',
      'hideRecordViewButton',
      'hideRecordPermissionsButton',
      'enableRecordPageNavigation',
      'editable',
      'draggable',
      'linkToParent',
      'showRefresh',
      'bulkRecordEditEnabled',
      'inlineRecordEditEnabled',
      'showRecordPerPageOption',
      'openRecordInEditMode',
    )

    if (o.selectionButtons) {
      this.options.selectionButtons = o.selectionButtons.map(b => new Button(b))
    }
  }

  async fetch (api: ComposeAPI, recordListModule: Module, filter: {[_: string]: unknown}): Promise<object> {
    if (recordListModule.moduleID !== this.options.moduleID) {
      throw Error('Module incompatible, module mismatch')
    }

    filter.moduleID = this.options.moduleID
    filter.namespaceID = recordListModule.namespaceID

    return api
      .recordList(filter)
      .then(r => {
        const { set: records, filter } = r as { filter: object; set: object[] }
        return { records, filter }
      })
  }
}

Registry.set(kind, PageBlockRecordList)
