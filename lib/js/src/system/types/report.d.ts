import { Step, Block, FilterDefinition } from '../../reporter'
interface PartialReport extends Partial<
  Omit<
    Report,
    | 'steps'
    | 'blocks'
    | 'scenarios'
    | 'createdAt'
    | 'createdBy'
    | 'updatedAt'
    | 'updatedBy'
    | 'deletedAt'
    | 'deletedBy'
  >
> {
  sources?: Array<ReportDataSource>
  blocks?: Array<unknown | Block>
  scenarios?: Array<ReportScenario>
  createdAt?: string | number | Date
  createdBy?: string
  updatedAt?: string | number | Date
  updatedBy?: string
  deletedAt?: string | number | Date
  deletedBy?: string
}
interface Meta {
  name?: string
  description?: string
  tags?: Array<string>
}
interface ReportDataSource {
  meta?: object
  step: Step
}
interface ReportScenario {
  scenarioID: string
  label: string
  datasource: string
  filter: FilterDefinition
}
export declare class Report {
  reportID: string
  handle: string
  meta: Meta
  sources: Array<ReportDataSource>
  blocks: Array<Block>
  scenarios: Array<ReportScenario>
  labels: object
  createdAt?: Date
  createdBy?: string
  updatedAt?: Date
  updatedBy?: string
  deletedAt?: Date
  deletedBy?: string
  canReadReport: boolean
  canUpdateReport: boolean
  canDeleteReport: boolean
  canGrant: boolean
  canRunReport: boolean
  constructor(r?: PartialReport)
  apply(r?: PartialReport): void
  clone(): Report
}
export {}
