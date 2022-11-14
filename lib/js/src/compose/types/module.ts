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

/**
 * System fields that are present in every object.
 */
export const systemFields = Object.freeze([
  { isSystem: true, name: 'recordID', label: 'Record ID', kind: 'String' },
  { isSystem: true, name: 'ownedBy', label: 'Owned by', kind: 'User' },
  { isSystem: true, name: 'createdBy', label: 'Created by', kind: 'User' },
  { isSystem: true, name: 'createdAt', label: 'Created at', kind: 'DateTime' },
  { isSystem: true, name: 'updatedBy', label: 'Updated by', kind: 'User' },
  { isSystem: true, name: 'updatedAt', label: 'Updated at', kind: 'DateTime' },
  { isSystem: true, name: 'deletedBy', label: 'Deleted by', kind: 'User' },
  { isSystem: true, name: 'deletedAt', label: 'Deleted at', kind: 'DateTime' },
].map(f => ModuleFieldMaker(f)))

interface PartialModule extends Partial<Omit<Module, 'fields' | 'meta' | 'labels' | 'createdAt' | 'updatedAt' | 'deletedAt'>> {
  fields?: Array<Partial<ModuleField>> | Array<ModuleField>;
  meta?: Partial<Meta>;

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
  public meta: object = {};

  public labels: object = {};

  public createdAt?: Date = undefined;
  public updatedAt?: Date = undefined;
  public deletedAt?: Date = undefined;

  public canUpdateModule = false;
  public canDeleteModule = false;
  public canCreateRecord = false;
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
      this.meta = { ...m.meta }
    }

    if (IsOf(m, 'labels')) {
      this.labels = { ...m.labels }
    }

    Apply(this, m, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')
    Apply(this, m, Boolean,
      'canUpdateModule',
      'canDeleteModule',
      'canCreateRecord',
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
      requested = (requested as ModuleField[]).map((f: ModuleField) => f.name)
    }

    const out: ModuleField[] = []

    for (const r of requested) {
      const sf = this.systemFields().find(f => r === f.name)
      if (sf) {
        out.push(sf)
        continue
      }

      const mf = this.fields.find(f => r === f.name)
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
