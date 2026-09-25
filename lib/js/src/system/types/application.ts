import { Apply, HumanID, ISO8601Date, NoID } from '../../cast'
import { IsOf } from '../../guards'

interface PartialApplication extends Partial<
  Omit<Application, 'createdAt' | 'updatedAt' | 'deletedAt' | 'lastUsedAt'>
> {
  createdAt?: string | number | Date
  updatedAt?: string | number | Date
  deletedAt?: string | number | Date
}

interface Meta {
  description: string
}

interface Unify {
  name: string
  listed: boolean
  url: string
  config: string
  iconID: string
  logoID: string
  // '' serves a shell section or a link; 'custom' serves the application's own
  // HTML source in the app view.
  kind?: string
  // The instance-wide home application; at most one holds it.
  home?: boolean
}

// What is known about a custom application's source without loading it.
interface SourceMeta {
  hash?: string
  size: number
  namespace?: string
  modules?: string[]
  writes?: string[]
  updatedAt?: string
  updatedBy?: string
}

export class Application {
  public applicationID = undefined
  public name = ''
  public ownerID = NoID
  public enabled = false
  public weight?: number = 0

  public meta?: Meta = {
    description: '',
  }

  public unify?: Unify = {
    name: '',
    listed: false,
    url: '',
    config: '',
    iconID: NoID,
    logoID: NoID,
    kind: '',
    home: false,
  }

  public sourceMeta?: SourceMeta = { size: 0 }

  public canGrant: boolean = true
  public canAccessApplication: boolean = true
  public canUpdateApplication: boolean = true
  public canDeleteApplication: boolean = true
  public canManageSourceOnApplication: boolean = false
  public createdAt?: Date = undefined
  public updatedAt?: Date = undefined
  public deletedAt?: Date = undefined

  constructor(r?: PartialApplication) {
    this.apply(r)
  }

  apply(r?: PartialApplication): void {
    Apply(this, r, HumanID, 'applicationID', 'ownerID')
    Apply(this, r, String, 'name')
    Apply(this, r, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')
    Apply(this, r, Number, 'weight')
    Apply(
      this,
      r,
      Boolean,
      'enabled',
      'canGrant',
      'canAccessApplication',
      'canUpdateApplication',
      'canDeleteApplication',
      'canManageSourceOnApplication',
    )

    if (r && IsOf(r, 'meta')) {
      this.meta = { ...this.meta, ...r.meta }
    }

    if (r && IsOf(r, 'unify')) {
      this.unify = { ...this.unify, ...r.unify }
    }

    if (r && IsOf(r, 'sourceMeta')) {
      this.sourceMeta = { size: 0, ...r.sourceMeta }
    }
  }

  /**
   * Returns resource ID
   */
  get resourceID(): string {
    return `${this.resourceType}:${this.applicationID}`
  }

  /**
   * Resource type
   */
  get resourceType(): string {
    return 'system:application'
  }

  clone(): Application {
    return new Application(JSON.parse(JSON.stringify(this)))
  }
}
