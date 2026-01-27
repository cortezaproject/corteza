import { ModuleField } from './module-field'
import { Namespace } from './namespace'
declare const propNamespace: unique symbol
interface MetaAdmin {
    fields: string[];
}
interface MetaUi {
    admin: MetaAdmin;
}
interface Meta {
    ui: MetaUi;
}
type systemFieldEncoding = null | {
    omit: true;
} | {
    ident: string;
};
interface Config {
    dal: {
        connectionID: string;
        ident: string;
        systemFieldEncoding: {
            id: systemFieldEncoding;
            revision: systemFieldEncoding;
            moduleID: systemFieldEncoding;
            namespaceID: systemFieldEncoding;
            ownedBy: systemFieldEncoding;
            createdBy: systemFieldEncoding;
            createdAt: systemFieldEncoding;
            updatedBy: systemFieldEncoding;
            updatedAt: systemFieldEncoding;
            deletedBy: systemFieldEncoding;
            deletedAt: systemFieldEncoding;
        };
    };
    privacy: {
        sensitivityLevelID: string;
        usageDisclosure: string;
    };
    discovery: {
        public: ConfigDiscoveryAccess;
        private: ConfigDiscoveryAccess;
        protected: ConfigDiscoveryAccess;
    };
    recordRevisions: {
        enabled: boolean;
        ident: string;
    };
    recordDeDup: {
        rules: RecordDeDupRule[];
    };
}
interface ConfigDiscoveryAccess {
    result: {
        lang: string;
        fields: string[];
    }[];
}
interface Constraint {
    attribute: string;
    modifier: string;
    multiValue: string;
    type: string;
}
interface RecordDeDupRule {
    name?: string;
    strict: boolean;
    constraints: Constraint[];
}
/**
 * System fields that are present in every record.
 */
export declare const systemFields: readonly ModuleField[]
interface PartialModule extends Partial<Omit<Module, 'fields' | 'meta' | 'labels' | 'createdAt' | 'updatedAt' | 'deletedAt'>> {
    fields?: Array<Partial<ModuleField>> | Array<ModuleField>;
    meta?: Partial<Meta>;
    config?: Partial<Config>;
    issues?: Array<string>;
    labels?: Partial<object>;
    createdAt?: string | number | Date;
    updatedAt?: string | number | Date;
    deletedAt?: string | number | Date;
}
export declare class Module {
  moduleID: string
  namespaceID: string
  name: string
  handle: string
  fields: Array<ModuleField>
  issues: Array<string>
  config: Partial<Config>
  meta: Meta
  labels: object
  createdAt?: Date
  updatedAt?: Date
  deletedAt?: Date
  canUpdateModule: boolean
  canDeleteModule: boolean
  canCreateRecord: boolean
  canCreateOwnedRecord: boolean
  canGrant: boolean
  private [propNamespace]?
  constructor(i?: PartialModule, ns?: Namespace);
  clone(): Module;
  apply(m?: PartialModule): void;
  /**
     * Returns resource ID
     */
  get resourceID(): string;
  /**
     * Resource type
     */
  get resourceType(): string;
  get namespace(): Namespace;
  set namespace(ns: Namespace);
  /**
     * Returns fields from module, filtered and order as requested
     */
  filterFields(requested?: string[] | Array<ModuleField>): Array<ModuleField>;
  findField(name: string): ModuleField | undefined;
  fieldNames(): readonly string[];
  systemFields(): readonly ModuleField[];
  export(): Module;
  import(): Module;
}
export {}
