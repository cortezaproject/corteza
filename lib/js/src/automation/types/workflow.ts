import { Apply, CortezaID, ISO8601Date, NoID } from '../../cast'
import { IsOf } from '../../guards'

interface Meta {
  name: '';
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
  public labels: object = {}
  public meta: object = {};

  public runAs = NoID
  public ownedBy = NoID;
  public createdBy = NoID;
  public createdAt?: Date = undefined
  public updatedAt?: Date = undefined
  public deletedAt?: Date = undefined

  constructor (w?: PartialWorkflow) {
    this.apply(w)
  }

  apply (w?: PartialWorkflow): void {
    Apply(this, w, CortezaID, 'workflowID')
    Apply(this, w, String, 'handle')

    Apply(this, w, Boolean, 'enabled')

    Apply(this, w, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')
    Apply(this, w, CortezaID, 'runAs', 'ownedBy', 'createdBy')

    if (IsOf(w, 'meta')) {
      this.meta = { ...w.meta }
    }

    if (IsOf(w, 'labels')) {
      this.labels = { ...w.labels }
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
