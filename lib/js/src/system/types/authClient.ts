import { Apply, CortezaID, ISO8601Date, NoID } from "../../cast";
import { IsOf } from "../../guards";

interface PartialAuthClient
  extends Partial<
    Omit<AuthClient, "createdAt" | "updatedAt" | "deletedAt" | "lastUsedAt">
  > {
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

export class AuthClient {
  public authClientID = NoID;
  public handle = "";
  public scope = "profile api";
  public redirectURI = "";
  public validGrant = 'authorization_code';
  
  public meta: AuthClientMeta = {
    name: "",
    description: "",
  };

  public security: DefSecurity = {
    impersonateUser: "0",
    permittedRoles: [],
    prohibitedRoles: [],
    forcedRoles: [],
  };

  public enabled = true;
  public trusted = false;

  public validFrom?: Date = undefined;
  public expiresAt?: Date = undefined;
  public createdAt?: Date = undefined;
  public updatedAt?: Date = undefined;
  public deletedAt?: Date = undefined;

  public createdBy = NoID;
  public updatedBy = NoID;
  public deletedBy = NoID;

  public canDeleteAuthClient = false;
  public canGrant = false;
  public canUpdateAuthClient = false;

  constructor (o?: PartialAuthClient) {
    this.apply(o)
  }

  apply (o?: PartialAuthClient): void {
    Apply(this, o, CortezaID, 'authClientID')
    Apply(this, o, ISO8601Date, 'validFrom', 'expiresAt', 'createdAt', 'updatedAt', 'deletedAt');
    Apply(this, o, String, 'handle', 'scope', 'redirectURI', 'validGrant');
    Apply(this, o, Boolean, 'enabled', 'trusted', 'canDeleteAuthClient', 'canGrant', 'canUpdateAuthClient');
    

    if (IsOf(o, 'meta')) {
        this.meta = { ...o.meta }
    }

    if (IsOf(o, 'security')) {
        this.security = { 
            ...this.security,
            ...o.security 
        }
    }

    Apply(this, o, CortezaID, 'createdBy', 'updatedBy', 'deletedBy');
  }

  clone(): AuthClient {
    return new AuthClient(JSON.parse(JSON.stringify(this)));
  }
}
