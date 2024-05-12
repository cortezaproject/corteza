import { merge } from 'lodash'
import { ModuleField, ModuleFieldMaker } from './module-field'
import { CortezaID, NoID, ISO8601Date, Apply } from '../../cast'
import { AreObjects, AreStrings, IsOf } from '../../guards'
import { Namespace } from './namespace'

const propNamespace = Symbol('namespace')

interface MetaAdmin {
  fields: string[];
}

interface MetaUi {
  admin: MetaAdmin;
}

interface Meta {
  ui: MetaUi;
}

type systemFieldEncoding = null | { omit: true } | { ident: string }

interface Config {
  dal: {
    connectionID: string;
    // operations
    // constraints
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
  }[]
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
export const systemFields = Object.freeze([
  { isSystem: true, name: 'recordID', label: 'Record ID', kind: 'String' },
  { isSystem: true, name: 'ownedBy', label: 'Owned by', kind: 'User' },
  { isSystem: true, name: 'createdBy', label: 'Created by', kind: 'User' },
  { isSystem: true, name: 'createdAt', label: 'Created at', kind: 'DateTime' },
  { isSystem: true, name: 'updatedBy', label: 'Updated by', kind: 'User' },
  { isSystem: true, name: 'updatedAt', label: 'Updated at', kind: 'DateTime' },
  { isSystem: true, name: 'revision', label: 'Revision', kind: 'Number' },
  { isSystem: true, name: 'deletedBy', label: 'Deleted by', kind: 'User' },
  { isSystem: true, name: 'deletedAt', label: 'Deleted at', kind: 'DateTime' },
].map(f => ModuleFieldMaker(f)))

interface PartialModule extends Partial<Omit<Module, 'fields' | 'meta' | 'labels' | 'createdAt' | 'updatedAt' | 'deletedAt'>> {
  fields?: Array<Partial<ModuleField>> | Array<ModuleField>;
  meta?: Partial<Meta>;
  config?: Partial<Config>;
  issues?: Array<string>;
  labels?: Partial<object>;

  createdAt?: string|number|Date;
  updatedAt?: string|number|Date;
  deletedAt?: string|number|Date;
}

export class Module {
  public moduleID = NoID;
  public namespaceID = NoID;
  public name = '';
  public handle = '';
  public fields: Array<ModuleField> = [];
  public issues: Array<string> = [];

  public config: Partial<Config> = {
    dal: {
      connectionID: NoID,
      ident: '',
      systemFieldEncoding: {
        id: null,
        revision: null,
        moduleID: null,
        namespaceID: null,
        ownedBy: null,
        createdBy: null,
        createdAt: null,
        updatedBy: null,
        updatedAt: null,
        deletedBy: null,
        deletedAt: null,
      },
    },

    privacy: {
      sensitivityLevelID: NoID,
      usageDisclosure: '',
    },

    discovery: {
      public: {
        result: [
          {
            lang: '',
            fields: [],
          },
        ],
      },
      private: {
        result: [
          {
            lang: '',
            fields: [],
          },
        ],
      },
      protected: {
        result: [
          {
            lang: '',
            fields: [],
          },
        ],
      },
    },

    recordRevisions: {
      enabled: false,
      ident: '',
    },

    recordDeDup: {
      rules: [],
    },
  }

  public meta: object = {};
  public labels: object = {};

  public createdAt?: Date = undefined;
  public updatedAt?: Date = undefined;
  public deletedAt?: Date = undefined;

  public canUpdateModule = false;
  public canDeleteModule = false;
  public canCreateRecord = false;
  public canCreateOwnedRecord = false;
  public canGrant = false;

  private [propNamespace]?: Namespace

  constructor (i?: PartialModule, ns?: Namespace) {
    if (ns) {
      this.namespace = ns
    }

    this.apply(i)
  }

  clone (): Module {
    return new Module(JSON.parse(JSON.stringify(this)), this.namespace)
  }

  apply (m?: PartialModule): void {
    if (!m) return

    if (this.namespace && m.namespaceID && m.namespaceID !== this.namespace.namespaceID) {
      throw new Error('module can not change namespace')
    }

    Apply(this, m, CortezaID, 'moduleID', 'namespaceID')
    Apply(this, m, String, 'name', 'handle')

    if (IsOf(m, 'fields')) {
      this.fields = []
      if (AreObjects(m.fields)) {
        // We're very permissive here -- array of (empty) objects is all we need
        // to create fields.
        this.fields = m.fields.map((b: { kind?: string }) => ModuleFieldMaker(b))
      }
    }

    if (IsOf(m, 'meta')) {
      if (m.meta.ui && m.meta.ui.admin && m.meta.ui.admin.fields) {
        if (!AreStrings(m.meta.ui.admin.fields)) {
          m.meta.ui.admin.fields = m.meta.ui.admin.fields.map((f: any) => f.fieldID)
        }
      }

      this.meta = { ...m.meta }
    }

    if (IsOf(m, 'config')) {
      this.config = merge({}, this.config, m.config)
    }

    if (IsOf(m, 'labels')) {
      this.labels = { ...m.labels }
    }

    if (IsOf(m, 'issues')) {
      this.issues = m.issues
    }

    Apply(this, m, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')
    Apply(this, m, Boolean,
      'canUpdateModule',
      'canDeleteModule',
      'canCreateRecord',
      'canCreateOwnedRecord',
      'canGrant',
    )
  }

  /**
   * Returns resource ID
   */
  get resourceID (): string {
    return `${this.resourceType}:${this.moduleID}`
  }

  /**
   * Resource type
   */
  get resourceType (): string {
    return 'compose:module'
  }

  public get namespace (): Namespace {
    return this[propNamespace] as Namespace
  }

  public set namespace (ns: Namespace) {
    if (this[propNamespace]) {
      if ((this[propNamespace] as Namespace).namespaceID !== ns.namespaceID) {
        throw new Error('namespace for this module already set')
      }
    }

    this.namespaceID = ns.namespaceID

    if (Object.isFrozen(ns)) {
      this[propNamespace] = ns
    } else {
      // Making a copy and freezing it
      this[propNamespace] = Object.freeze(new Namespace(ns))
    }

    this[propNamespace] = ns
  }

  /**
   * Returns fields from module, filtered and order as requested
   */
  filterFields (requested?: string[] | Array<ModuleField>): Array<ModuleField> {
    if (!requested || requested.length === 0) {
      return []
    }

    if (!AreStrings(requested)) {
      requested = (requested as ModuleField[]).map((f: ModuleField) => f.name || f.fieldID)
    }

    const out: ModuleField[] = []

    for (const r of requested) {
      const sf = this.systemFields().find(f => r === f.name || r === f.fieldID)
      if (sf) {
        out.push(sf)
        continue
      }

      const mf = this.fields.find(f => r === f.name || r === f.fieldID)
      if (mf) {
        out.push(mf)
      }
    }

    return out
  }

  public findField (name: string): ModuleField|undefined {
    const r = this.filterFields([name])
    return r && r.length > 0 ? r[0] : undefined
  }

  fieldNames (): readonly string[] {
    return this.fields.map(f => f.name)
  }

  systemFields (): readonly ModuleField[] {
    return systemFields
  }

  export (): Module {
    return this
  }

  import (): Module {
    return this
  }
}
