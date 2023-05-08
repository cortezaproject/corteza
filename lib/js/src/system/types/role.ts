import { Apply, CortezaID, ISO8601Date, NoID } from '../../cast'
import { IsOf, AreStrings } from '../../guards'

interface PartialRole extends Partial<Omit<Role, 'createdAt' | 'updatedAt' | 'deletedAt' | 'archivedAt'>> {
  createdAt?: string|number|Date;
  updatedAt?: string|number|Date;
  deletedAt?: string|number|Date;
  archivedAt?: string|number|Date;
}

interface Meta {
  description: string;
  context: MetaContext;
}

interface MetaContext {
  resourceTypes: Array<string>;
  expr: string;
}

const defaultMeta = {
  description: '',
  context: {
    resourceTypes: [],
    expr: '',
  },
}

export class Role {
  public roleID = NoID
  public name = ''
  public handle = ''
  public members: string[] = []
  public labels: object = {}
  public meta: Meta = { ...defaultMeta }

  public createdAt?: Date = undefined
  public updatedAt?: Date = undefined
  public deletedAt?: Date = undefined
  public archivedAt?: Date = undefined

  public isSystem = false;
  public isClosed = false;
  public isBypass = false;
  public canGrant = false;
  public canUpdateRole = false;
  public canDeleteRole = false;
  public canManageMembersOnRole = false;

  constructor (r?: PartialRole) {
    this.apply(r)
  }

  apply (r?: PartialRole): void {
    Apply(this, r, CortezaID, 'roleID')

    Apply(this, r, String, 'name', 'handle')

    if (r?.members) {
      this.members = []
      if (AreStrings(r.members)) {
        this.members = r.members
      }
    }

    if (IsOf(r, 'meta')) {
      this.meta = { ...r.meta }
    }

    if (!this.meta) {
      this.meta = { ...defaultMeta }
    }

    if (!this.meta.context) {
      this.meta.context = { ...defaultMeta.context }
    }

    if (IsOf(r, 'labels')) {
      this.labels = { ...r.labels }
    }

    Apply(this, r, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt', 'archivedAt')
    Apply(this, r, Boolean, 'isSystem', 'isClosed', 'isBypass', 'canGrant', 'canUpdateRole', 'canDeleteRole', 'canManageMembersOnRole')
  }

  /**
   * Returns resource ID
   */
  get resourceID (): string {
    return `${this.resourceType}:${this.roleID}`
  }

  /**
   * Resource type
   */
  get resourceType (): string {
    return 'system:role'
  }

  get isContext (): boolean {
    return this.meta?.context?.expr?.length > 0
  }

  clone (): Role {
    return new Role(JSON.parse(JSON.stringify(this)))
  }
}
