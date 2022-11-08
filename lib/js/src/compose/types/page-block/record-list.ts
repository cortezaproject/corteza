import { PageBlock, PageBlockInput, Registry } from './base'
import { Apply, CortezaID, NoID } from '../../../cast'
import { Compose as ComposeAPI } from '../../../api-clients'
import { Module } from '../module'
import { Button } from './types'

const kind = 'RecordList'
interface Options {
  moduleID: string;
  prefilter: string;
  presort: string;
  fields: unknown[];
  hideHeader: boolean;
  hideAddButton: boolean;
  hideImportButton: boolean;
  hideSearch: boolean;
  hidePaging: boolean;
  hideSorting: boolean;
  hideRecordReminderButton: boolean;
  hideRecordCloneButton: boolean;
  hideRecordEditButton: boolean;
  hideRecordViewButton: boolean;
  hideRecordPermissionsButton: boolean;
  allowExport: boolean;
  perPage: number;
  recordDisplayOption: string;

  fullPageNavigation: boolean;
  showTotalCount: boolean;

  // Record-lines
  editable: boolean;
  draggable?: boolean;
  positionField?: string;
  refField?: string;
  editFields?: unknown[];

  // When adding a new record, link it to parent when available
  linkToParent: boolean;

  // Should records be opened in a new tab
  // legacy field that has been removed but we keep it for backwards compatibility
  openInNewTab: boolean;

  // Are table rows selectable
  selectable: boolean;
  selectMode: 'multi' | 'single' | 'range';

  // Ordered list of buttons to display in the block
  selectionButtons: Array<Button>;
}

const defaults: Readonly<Options> = Object.freeze({
  moduleID: NoID,
  prefilter: '',
  presort: '',
  fields: [],
  hideHeader: false,
  hideAddButton: false,
  hideImportButton: false,
  hideSearch: false,
  hidePaging: false,
  hideSorting: false,
  hideRecordReminderButton: true,
  hideRecordCloneButton: true,
  hideRecordEditButton: false,
  hideRecordViewButton: true,
  hideRecordPermissionsButton: true,
  allowExport: false,
  perPage: 20,
  recordDisplayOption: 'sameTab',

  fullPageNavigation: true,
  showTotalCount: true,

  editable: false,
  draggable: false,
  positionField: undefined,
  refField: undefined,
  editFields: [],

  linkToParent: true,

  openInNewTab: false,

  selectable: true,
  selectMode: 'multi',

  selectionButtons: [],
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
    Apply(this.options, o, String, 'prefilter', 'presort', 'selectMode', 'positionField', 'refField', 'recordDisplayOption')
    Apply(this.options, o, Number, 'perPage')

    if (o.fields) {
      this.options.fields = o.fields
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
      'hideSearch',
      'hidePaging',
      'fullPageNavigation',
      'showTotalCount',
      'hideSorting',
      'allowExport',
      'selectable',
      'hideRecordReminderButton',
      'hideRecordCloneButton',
      'hideRecordEditButton',
      'hideRecordViewButton',
      'hideRecordPermissionsButton',
      'editable',
      'draggable',
      'linkToParent',
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
