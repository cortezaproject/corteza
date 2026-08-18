import { PageBlock, PageBlockInput, Registry, BlockIssue, isUnsetID } from './base'
import { Apply, HumanID, NoID } from '../../../cast'

const kind = 'RecordOrganizer'
interface Options {
  moduleID: string
  labelField: string
  descriptionField: string
  filter: string
  positionField: string
  groupField: string
  group: string
  refreshRate: number
  showRefresh: boolean
  magnifyOption: string
  displayOption: string
  addRecordDisplayOption: string
}

const defaults: Readonly<Options> = Object.freeze({
  moduleID: NoID,
  labelField: '',
  descriptionField: '',
  filter: '',
  positionField: '',
  groupField: '',
  group: '',
  refreshRate: 0,
  showRefresh: false,
  magnifyOption: '',
  displayOption: 'sameTab',
  addRecordDisplayOption: '',
})

export class PageBlockRecordOrganizer extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor(i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions(o?: Partial<Options>): void {
    if (!o) return

    Apply(this.options, o, HumanID, 'moduleID')
    Apply(
      this.options,
      o,
      String,
      'labelField',
      'descriptionField',
      'filter',
      'positionField',
      'groupField',
      'group',
      'magnifyOption',
      'displayOption',
      'addRecordDisplayOption',
    )
    Apply(this.options, o, Number, 'refreshRate')
    Apply(this.options, o, Boolean, 'showRefresh')
  }

  validate(): Array<BlockIssue> {
    const ee = super.validate()
    if (isUnsetID(this.options.moduleID)) {
      ee.push({ option: 'moduleID', labelKey: 'block.issue.module' })
    }
    return ee
  }
}

Registry.set(kind, PageBlockRecordOrganizer)
