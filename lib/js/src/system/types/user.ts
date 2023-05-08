import { Apply, CortezaID, ISO8601Date, NoID } from '../../cast'
import { IsOf, AreStrings } from '../../guards'

interface PartialUser extends Partial<Omit<User, 'createdAt' | 'updatedAt' | 'deletedAt' | 'suspendedAt'>> {
  createdAt?: string|number|Date;
  updatedAt?: string|number|Date;
  deletedAt?: string|number|Date;
  suspendedAt?: string|number|Date;
}

interface UserMeta {
  preferredLanguage?: string;
  securityPolicy?: SecurityPolicy;
  avatarID?: string;
  avatarKind?: string;
  avatarColor?: string;
  avatarBgColor?: string;
}

interface SecurityPolicy {
  mfa: MFA;
}

interface MFA {
  enforcedEmailOTP: boolean;
  enforcedTOTP: boolean;
}

export class User {
  public userID = NoID
  public handle = ''
  public username = ''
  public email = ''
  public name = ''
  public emailConfirmed = false
  public labels: object = {}
  public meta: UserMeta = {
    preferredLanguage: 'en',
    securityPolicy: {
      mfa: {
        enforcedEmailOTP: false,
        enforcedTOTP: false,
      },
    },
    avatarID: NoID,
    avatarKind: '',
    avatarColor: '',
    avatarBgColor: '',
  }

  public canGrant = false
  public canUpdateUser = false
  public canDeleteUser = false
  public createdAt?: Date = undefined
  public updatedAt?: Date = undefined
  public deletedAt?: Date = undefined
  public suspendedAt?: Date = undefined
  public roles?: Array<string>

  constructor (u?: PartialUser) {
    this.apply(u)
  }

  apply (u?: PartialUser): void {
    Apply(this, u, CortezaID, 'userID')
    Apply(this, u, String, 'handle', 'username', 'email', 'name')
    Apply(this, u, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt', 'suspendedAt')
    Apply(this, u, Boolean, 'emailConfirmed', 'canGrant', 'canUpdateUser', 'canDeleteUser')

    if (u?.roles) {
      this.roles = []
      if (AreStrings(u.roles)) {
        this.roles = u.roles
      }
    }

    if (IsOf(u, 'meta')) {
      this.meta = { ...u.meta }
    }

    if (IsOf(u, 'labels')) {
      this.labels = { ...u.labels }
    }
  }

  /**
   * Returns resource ID
   */
  get resourceID (): string {
    return `${this.resourceType}:${this.userID}`
  }

  /**
   * Resource type
   */
  get resourceType (): string {
    return 'system:user'
  }

  get fts (): string {
    return [
      this.name,
      this.username,
      this.handle,
      this.email,
      this.userID,
    ].join(' ').toLocaleLowerCase()
  }

  clone (): User {
    return new User(JSON.parse(JSON.stringify(this)))
  }
}
