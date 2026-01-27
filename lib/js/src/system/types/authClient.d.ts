interface PartialAuthClient extends Partial<Omit<AuthClient, 'createdAt' | 'updatedAt' | 'deletedAt' | 'lastUsedAt'>> {
    createdAt?: string | number | Date;
    updatedAt?: string | number | Date;
    deletedAt?: string | number | Date;
    lastUsedAt?: string | number | Date;
}
interface AuthClientMeta {
    name: string;
    description: string;
}
interface DefSecurity {
    impersonateUser: string;
    permittedRoles: Array<string>;
    prohibitedRoles: Array<string>;
    forcedRoles: Array<string>;
}
export declare class AuthClient {
  authClientID: string
  handle: string
  scope: string
  redirectURI: string
  validGrant: string
  meta: AuthClientMeta
  security: DefSecurity
  enabled: boolean
  trusted: boolean
  validFrom?: Date
  expiresAt?: Date
  createdAt?: Date
  updatedAt?: Date
  deletedAt?: Date
  createdBy: string
  updatedBy: string
  deletedBy: string
  canDeleteAuthClient: boolean
  canGrant: boolean
  canUpdateAuthClient: boolean
  constructor(o?: PartialAuthClient);
  apply(o?: PartialAuthClient): void;
  clone(): AuthClient;
}
export {}
