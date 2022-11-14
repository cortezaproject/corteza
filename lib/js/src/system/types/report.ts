import { Step, Block, FilterDefinition } from '../../reporter'
import { Apply, CortezaID, ISO8601Date, NoID } from '../../cast'
import { IsOf } from '../../guards'

interface PartialReport extends Partial<Omit<Report, 'steps' | 'blocks' | 'scenarios' | 'createdAt' | 'createdBy' | 'updatedAt' | 'updatedBy' | 'deletedAt' | 'deletedBy'>> {
  sources?: Array<ReportDataSource>;
  blocks?: Array<unknown|Block>;
  scenarios?: Array<ReportScenario>;
  createdAt?: string|number|Date;
  createdBy?: string;
  updatedAt?: string|number|Date;
  updatedBy?: string;
  deletedAt?: string|number|Date;
  deletedBy?: string;
}

interface Meta {
  name?: string;
  description?: string;
  tags?: Array<string>;
}

interface ReportDataSource {
  meta?: Object;
  step: Step;
}

interface ReportScenario {
  scenarioID: string;
  label: string;
  datasource: string;
  filter: FilterDefinition;
}

export class Report {
  public reportID = NoID
  public handle = ''
  public meta: Meta = {}
  public sources: Array<ReportDataSource> = []
  public blocks: Array<Block> = []
  public scenarios: Array<ReportScenario> = []

  public labels: object = {}
  public createdAt?: Date = undefined
  public createdBy?: string = undefined
  public updatedAt?: Date = undefined
  public updatedBy?: string = undefined
  public deletedAt?: Date = undefined
  public deletedBy?: string = undefined

  public canReadReport = false;
  public canUpdateReport = false;
  public canDeleteReport = false;
  public canGrant = false;
  public canRunReport = false;

  constructor (r?: PartialReport) {
    this.apply(r)
  }

  apply (r?: PartialReport): void {
    Apply(this, r, CortezaID, 'reportID')

    Apply(this, r, String, 'handle')

    if (r && IsOf(r, 'meta')) {
      this.meta = r.meta
    }


    this.sources = []

    for (const s of r?.sources || []) {
      s.step = s.step as Step
      this.sources.push(s as ReportDataSource)
    }

    if (r?.blocks) {
      this.blocks = []
      for (const p of r.blocks) {
        this.blocks.push(new Block(p as Block))
      }
    }

    if (r?.scenarios) {
      this.scenarios = []
      for (const s of r.scenarios) {
        this.scenarios.push(s as ReportScenario)
      }
    }

    if (IsOf(r, 'labels')) {
      this.labels = { ...r.labels }
    }

    Apply(this, r, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')
    Apply(this, r, CortezaID, 'createdBy', 'updatedBy', 'deletedBy')
    Apply(this, r, Boolean,
      'canReadReport',
      'canUpdateReport',
      'canDeleteReport', 
      'canGrant',
      'canRunReport',
    )
  }
}
