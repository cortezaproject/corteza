import { Apply, CortezaID, ISO8601Date, NoID } from '../../cast'
import { IsOf } from '../../guards'

export type RevisionStatus = '' | 'draft'
export type RevisionOperation = 'created' | 'updated' | 'soft-deleted' | 'undeleted' | 'hard-deleted'

export interface RevisionChange {
  key: string
  old: unknown[]
  new: unknown[]
}

interface PartialRevision extends Partial<Omit<Revision, 'timestamp' | 'deletedAt'>> {
  timestamp?: string | number | Date
  deletedAt?: string | number | Date
}

export class Revision {
  public changeID = ''
  public timestamp?: Date = undefined
  public resourceID = NoID
  public resourceType = ''
  public revision = 0
  public operation: RevisionOperation = 'created'
  public status: RevisionStatus = ''
  public userID = NoID
  public changes: RevisionChange[] = []
  public comment = ''
  public deletedAt?: Date = undefined
  public deletedBy = NoID

  constructor (r?: PartialRevision) {
    this.apply(r)
  }

  apply (r?: PartialRevision): void {
    if (!r) return

    Apply(this, r, String, 'changeID')
    Apply(this, r, CortezaID, 'resourceID', 'userID', 'deletedBy')
    Apply(this, r, ISO8601Date, 'timestamp', 'deletedAt')
    Apply(this, r, Number, 'revision')
    Apply(this, r, String, 'resourceType', 'operation', 'status', 'comment')

    if (IsOf(r, 'changes') && Array.isArray(r.changes)) {
      this.changes = r.changes.map(c => ({
        key: c.key || '',
        old: Array.isArray(c.old) ? c.old : [],
        new: Array.isArray(c.new) ? c.new : [],
      }))
    }
  }

  get resourceIdentifier (): string {
    return `${this.resourceTypeIdentifier}:${this.changeID}`
  }

  get resourceTypeIdentifier (): string {
    return 'system:revision'
  }

  get isDraft (): boolean {
    return this.status === 'draft'
  }

  clone (): Revision {
    return new Revision(JSON.parse(JSON.stringify(this)))
  }
}
