interface PartialUser extends Partial<
  Omit<User, 'createdAt' | 'updatedAt' | 'deletedAt' | 'suspendedAt'>
> {
  createdAt?: string | number | Date
  updatedAt?: string | number | Date
  deletedAt?: string | number | Date
  suspendedAt?: string | number | Date
}
interface UserMeta {
  preferredLanguage?: string
  securityPolicy?: SecurityPolicy
  avatarID?: string
  avatarKind?: string
  avatarColor?: string
  avatarBgColor?: string
  theme?: string
}
interface SecurityPolicy {
  mfa: MFA
}
interface MFA {
  enforcedEmailOTP: boolean
  enforcedTOTP: boolean
}
export declare class User {
  userID: string
  handle: string
  username: string
  email: string
  name: string
  emailConfirmed: boolean
  labels: object
  meta: UserMeta
  canGrant: boolean
  canUpdateUser: boolean
  canDeleteUser: boolean
  createdAt?: Date
  updatedAt?: Date
  deletedAt?: Date
  suspendedAt?: Date
  roles?: Array<string>
  constructor(u?: PartialUser)
  apply(u?: PartialUser): void
  /**
   * Returns resource ID
   */
  get resourceID(): string
  /**
   * Resource type
   */
  get resourceType(): string
  get fts(): string
  clone(): User
  properties(): string[]
}
export {}
