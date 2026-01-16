import { Apply, CortezaID, ISO8601Date, NoID } from '../../cast'
import { IsOf } from '../../guards'

interface Meta {
  name?: string;
  description?: string;
  visual?: Record<string, unknown>;
  subWorkflow?: boolean;
}

interface PartialWorkflow extends Partial<Omit<Workflow, 'createdAt' | 'updatedAt' | 'deletedAt' | 'suspendedAt'>> {
  meta?: Partial<Meta>;
  createdAt?: string|number|Date;
  updatedAt?: string|number|Date;
  deletedAt?: string|number|Date;
}

export class Workflow {
  public workflowID = NoID
  public handle = ''
  public enabled = false
  public trace = false
  public keepSessions = 0
  public labels: Record<string, string> = {}
  public meta: Meta = {
    name: '',
    description: '',
    visual: {},
    subWorkflow: false,
  };

  public scope?: Record<string, unknown> = undefined
  public steps?: unknown[] = undefined
  public paths?: unknown[] = undefined
  public issues?: unknown[] = undefined

  public runAs = NoID
  public ownedBy = NoID;
  public createdBy = NoID;
  public updatedBy = NoID;
  public deletedBy = NoID;
  public createdAt?: Date = undefined
  public updatedAt?: Date = undefined
  public deletedAt?: Date = undefined

  public canGrant = false
  public canUpdateWorkflow = false
  public canDeleteWorkflow = false
  public canExecuteWorkflow = false

  constructor (w?: PartialWorkflow) {
    this.apply(w)
  }

  apply (w?: PartialWorkflow): void {
    Apply(this, w, CortezaID, 'workflowID')
    Apply(this, w, String, 'handle')

    Apply(this, w, Boolean, 'enabled', 'trace')
    Apply(this, w, Number, 'keepSessions')

    Apply(this, w, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')
    Apply(this, w, CortezaID, 'runAs', 'ownedBy', 'createdBy', 'updatedBy', 'deletedBy')

    Apply(this, w, Boolean, 'canGrant', 'canUpdateWorkflow', 'canDeleteWorkflow', 'canExecuteWorkflow')

    if (IsOf(w, 'meta')) {
      this.meta = { ...w.meta }
    }

    if (IsOf(w, 'labels')) {
      this.labels = { ...w.labels }
    }

    if (IsOf(w, 'scope')) {
      this.scope = w.scope
    }

    if (IsOf(w, 'steps')) {
      this.steps = w.steps
    }

    if (IsOf(w, 'paths')) {
      this.paths = w.paths
    }

    if (IsOf(w, 'issues')) {
      this.issues = w.issues
    }
  }

  /**
   * Returns resource ID
   */
  get resourceID (): string {
    return `${this.resourceType}:${this.workflowID}`
  }

  /**
   * Resource type
   */
  get resourceType (): string {
    return 'automation:workflow'
  }
}
