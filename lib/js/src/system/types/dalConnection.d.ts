interface PartialDalConnection extends Partial<Omit<DalConnection, 'createdAt' | 'updatedAt' | 'deletedAt'>> {
    createdAt?: string | number | Date;
    updatedAt?: string | number | Date;
    deletedAt?: string | number | Date;
}
interface ConnectionMetaProperty {
    enabled: boolean;
    notes: string;
}
interface ConnectionMetaProperties {
    dataAtRestEncryption?: ConnectionMetaProperty;
    dataAtRestProtection?: ConnectionMetaProperty;
    dataAtTransitEncryption?: ConnectionMetaProperty;
    dataRestoration?: ConnectionMetaProperty;
}
interface ConnectionMeta {
    name: string;
    ownership: string;
    location?: object;
    properties?: ConnectionMetaProperties;
}
interface ConnectionConfigDAL {
    type?: string;
    params?: object;
    modelIdent?: string;
    modelIdentCheck?: Array<string>;
}
interface ConnectionConfigPrivacy {
    sensitivityLevelID: string;
}
interface ConnectionConfig {
    privacy: ConnectionConfigPrivacy;
    dal?: ConnectionConfigDAL;
}
export declare class DalConnection {
  connectionID: string
  handle: string
  type: string
  meta: ConnectionMeta
  config: ConnectionConfig
  issues: never[]
  labels: never[]
  createdAt?: Date
  updatedAt?: Date
  deletedAt?: Date
  createdBy: string
  updatedBy: string
  deletedBy: string
  canDeleteConnection: boolean
  canManageDalConfig: boolean
  constructor(dc?: PartialDalConnection);
  apply(dc?: PartialDalConnection): void;
  clone(): DalConnection;
}
export {}
