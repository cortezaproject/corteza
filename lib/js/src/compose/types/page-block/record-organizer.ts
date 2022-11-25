import { PageBlock, PageBlockInput, Registry } from './base'
import { Apply, CortezaID, NoID } from '../../../cast'

const kind = 'RecordOrganizer'
interface Options {
  moduleID: string;
  labelField: string;
  descriptionField: string;
  filter: string;
  positionField: string;
  groupField: string;
  group: string;
  refreshRate: number;
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
})

export class PageBlockRecordOrganizer extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor (i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return

    Apply(this.options, o, CortezaID, 'moduleID')
    Apply(this.options, o, String, 'labelField', 'descriptionField', 'filter', 'positionField', 'groupField', 'group')
    Apply(this.options, o, Number, 'refreshRate')
  }
}

Registry.set(kind, PageBlockRecordOrganizer)
