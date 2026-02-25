interface PartialApplication extends Partial<
  Omit<Application, 'createdAt' | 'updatedAt' | 'deletedAt' | 'lastUsedAt'>
> {
  createdAt?: string | number | Date
  updatedAt?: string | number | Date
  deletedAt?: string | number | Date
}
interface Unify {
  name: string
  listed: boolean
  url: string
  config: string
  iconID: string
  logoID: string
}
export declare class Application {
  applicationID: undefined
  name: string
  ownerID?: number
  enabled: boolean
  weight?: number
  unify?: Unify
  canGrant: boolean
  canUpdateApplication: boolean
  canDeleteApplication: boolean
  createdAt?: Date
  updatedAt?: Date
  deletedAt?: Date
  constructor(r?: PartialApplication)
  apply(r?: PartialApplication): void
  /**
   * Returns resource ID
   */
  get resourceID(): string
  /**
   * Resource type
   */
  get resourceType(): string
  clone(): Application
}
export {}
