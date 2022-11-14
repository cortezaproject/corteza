import { Apply, CortezaID, ISO8601Date, NoID } from '../../cast'
import { IsOf } from '../../guards'

interface PartialTemplate extends Partial<Omit<Template, 'createdAt' | 'updatedAt' | 'deletedAt' | 'lastUsedAt'>> {
  createdAt?: string|number|Date;
  updatedAt?: string|number|Date;
  deletedAt?: string|number|Date;
  lastUsedAt?: string|number|Date;
}

interface Meta {
  short?: string;
  description?: string;
}

export class Template {
  public templateID = NoID
  public handle = ''
  public language = ''
  public type = 'text/html'
  public partial = false
  public meta: Meta = {}
  public template = ''
  public labels: object = {}
  public ownerID = NoID
  public createdAt?: Date = undefined
  public updatedAt?: Date = undefined
  public deletedAt?: Date = undefined
  public lastUsedAt?: Date = undefined

  constructor (r?: PartialTemplate) {
    this.apply(r)
  }

  apply (r?: PartialTemplate): void {
    Apply(this, r, CortezaID, 'templateID', 'ownerID')

    Apply(this, r, String, 'handle', 'language', 'type', 'template')
    Apply(this, r, Boolean, 'partial')

    if (r && IsOf(r, 'meta')) {
      this.meta = r.meta
    }

    Apply(this, r, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt', 'lastUsedAt')

    if (IsOf(r, 'labels')) {
      this.labels = { ...r.labels }
    }
  }

  /**
   * Returns resource ID
   */
  get resourceID (): string {
    return `${this.resourceType}:${this.templateID}`
  }

  /**
   * Resource type
   */
  get resourceType (): string {
    return 'system:template'
  }
}
