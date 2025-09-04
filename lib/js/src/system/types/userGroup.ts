import { Apply, CortezaID, ISO8601Date, NoID } from '../../cast'
import { IsOf, AreStrings } from '../../guards'

interface PartialUserGroup extends Partial<Omit<UserGroup, 'createdAt' | 'updatedAt' | 'deletedAt' | 'suspendedAt'>> {
  createdAt?: string|number|Date;
  updatedAt?: string|number|Date;
  deletedAt?: string|number|Date;
  suspendedAt?: string|number|Date;
}

interface Meta {
  description: string;
  short: string;
}

const defaultMeta = {
  description: '',
  short: '',
}

export class UserGroup {
  public userGroupID = NoID
  public handle = ''
  public labels: object = {}

  public meta: Meta = { ...defaultMeta }

  public selfID = NoID

  public canGrant = false
  public canUpdateUserGroup = false
  public canDeleteUserGroup = false
  public canManageMembersOnUserGroup = false;

  public createdAt?: Date = undefined
  public updatedAt?: Date = undefined
  public deletedAt?: Date = undefined
  public suspendedAt?: Date = undefined
  public roles?: Array<string>

  constructor (u?: PartialUserGroup) {
    this.apply(u)
  }

  apply (u?: PartialUserGroup): void {
    Apply(this, u, CortezaID, 'userGroupID', 'selfID')
    Apply(this, u, String, 'handle')
    Apply(this, u, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt', 'suspendedAt')
    Apply(this, u, Boolean, 'canGrant', 'canUpdateUserGroup', 'canDeleteUserGroup', 'canManageMembersOnUserGroup')

    if (u?.roles) {
      this.roles = []
      if (AreStrings(u.roles)) {
        this.roles = u.roles
      }
    }

    if (IsOf(u, 'meta')) {
      this.meta = { ...u.meta }
    }

    if (!this.meta) {
      this.meta = { ...defaultMeta }
    }

    if (IsOf(u, 'labels')) {
      this.labels = { ...u.labels }
    }
  }

  /**
   * Returns resource ID
   */
  get resourceID (): string {
    return `${this.resourceType}:${this.userGroupID}`
  }

  /**
   * Resource type
   */
  get resourceType (): string {
    return 'system:user-group'
  }

  get fts (): string {
    return [
      this.meta.short,
      this.handle,
      this.userGroupID,
    ].join(' ').toLocaleLowerCase()
  }

  clone (): UserGroup {
    return new UserGroup(JSON.parse(JSON.stringify(this)))
  }

  properties (): string[] {
    return [
      'userGroupID',
      'handle',
      'labels',
      'canGrant',
      'canUpdateUserGroup',
      'canDeleteUserGroup',
      'canManageMembersOnUserGroup',
      'createdAt',
      'updatedAt',
      'deletedAt',
      'suspendedAt',
      'roles',
    ]
  }
}
