import { Apply, CortezaID, ISO8601Date, NoID } from '../../cast'
import { IsOf } from '../../guards'

interface PartialApplication extends Partial<Omit<Application, 'createdAt' | 'updatedAt' | 'deletedAt' | 'lastUsedAt'>> {
  createdAt?: string|number|Date;
  updatedAt?: string|number|Date;
  deletedAt?: string|number|Date;
}

interface Unify {
  name: string;
  listed: boolean,
  url: string,
  config: string,
  iconID: string,
  logoID: string
}

export class Application {
  public applicationID = undefined
  public name = ''
  public ownerID?: number = 0;
  public enabled = false
  public weight?: number = 0;

  public unify?: Unify = {
    name: '',
    listed: false,
    url: '',
    config: '',
    iconID: NoID,
    logoID: NoID,
  };

  public canGrant: boolean = true;
  public canUpdateApplication: boolean = true;
  public canDeleteApplication: boolean = true;
  public createdAt?: Date = undefined
  public updatedAt?: Date = undefined
  public deletedAt?: Date = undefined

  constructor (r?: PartialApplication) {
    this.apply(r)
  }

  apply (r?: PartialApplication): void {
    Apply(this, r, CortezaID, 'applicationID')
    Apply(this, r, String, 'name')
    Apply(this, r, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')
    Apply(this, r, Number, 'weight', 'ownerID')
    Apply(this, r, Boolean, 'enabled', 'canGrant', 'canUpdateApplication', 'canDeleteApplication')


    if (r && IsOf(r, 'unify')) {
      this.unify = r.unify
    }
  }

  /**
   * Returns resource ID
   */
  get resourceID (): string {
    return `${this.resourceType}:${this.applicationID}`
  }

  /**
   * Resource type
   */
  get resourceType (): string {
    return 'system:application'
  }

  clone (): Application {
    return new Application(JSON.parse(JSON.stringify(this)))
  }
}
