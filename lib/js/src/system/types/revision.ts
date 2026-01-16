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
  public resource = ''
  public revision = 0
  public operation: RevisionOperation = 'created'
  public status: RevisionStatus = ''
  public userID = NoID
  public changes: RevisionChange[] = []
  public comment = ''
  public deletedAt?: Date = undefined
  public deletedBy = NoID
  public record?: unknown = undefined

  constructor (r?: PartialRevision) {
    this.apply(r)
  }

  apply (r?: PartialRevision): void {
    if (!r) return

    Apply(this, r, String, 'changeID', 'resource')
    Apply(this, r, CortezaID, 'userID', 'deletedBy')
    Apply(this, r, ISO8601Date, 'timestamp', 'deletedAt')
    Apply(this, r, Number, 'revision')
    Apply(this, r, String, 'operation', 'status', 'comment')
    Apply(this, r, (v: unknown) => v, 'record')

    if (IsOf(r, 'changes') && Array.isArray(r.changes)) {
      this.changes = r.changes.map(c => ({
        key: c.key || '',
        old: Array.isArray(c.old) ? c.old : [],
        new: Array.isArray(c.new) ? c.new : [],
      }))
    }
  }

  get resourceIdentifier (): string {
    return this.resource || `system:revision:${this.changeID}`
  }

  get isDraft (): boolean {
    return this.status === 'draft'
  }

  clone (): Revision {
    return new Revision(JSON.parse(JSON.stringify(this)))
  }
}
