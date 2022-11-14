import { merge } from 'lodash'
import { Apply, CortezaID, ISO8601Date, NoID } from '../../cast'
import { IsOf } from '../../guards'

interface PartialDalConnection extends Partial<Omit<DalConnection, 'createdAt' | 'updatedAt' | 'deletedAt'>> {
  createdAt?: string|number|Date;
  updatedAt?: string|number|Date;
  deletedAt?: string|number|Date;
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

export class DalConnection {
  public connectionID = NoID
  public handle = ''
  public type = 'corteza::system:dal-connection'
  public meta: ConnectionMeta = {
    name: '',
    ownership: '',
    location: {
      properties: { name: '' },
      geometry: {
        coordinates: [],
        type: '',
      },
    },
    properties: {
      dataAtRestEncryption: {
        enabled: false,
        notes: '',
      },
      dataAtRestProtection: {
        enabled: false,
        notes: '',
      },
      dataAtTransitEncryption: {
        enabled: false,
        notes: '',
      },
      dataRestoration: {
        enabled: false,
        notes: '',
      },
    },
  }

  public config: ConnectionConfig = {
    privacy: { sensitivityLevelID: NoID },
    dal: {},
  }

  public issues = []
  public labels = []

  public createdAt?: Date = undefined
  public updatedAt?: Date = undefined
  public deletedAt?: Date = undefined

  public createdBy = ''
  public updatedBy = ''
  public deletedBy = ''

  public canDeleteConnection = false
  public canManageDalConfig = false

  constructor (dc?: PartialDalConnection) {
    this.apply(dc)
  }

  apply (dc?: PartialDalConnection): void {
    Apply(this, dc, CortezaID, 'connectionID')
    Apply(this, dc, String, 'handle', 'type')
    Apply(this, dc, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')
    Apply(this, dc, CortezaID, 'createdBy', 'updatedBy', 'deletedBy')
    Apply(this, dc, Boolean, 'canDeleteConnection', 'canManageDalConfig')

    if (IsOf(dc, 'meta')) {
      this.meta = merge(this.meta, dc.meta)
    }

    if (IsOf(dc, 'config')) {
      this.config = { ...dc.config.privacy }

      if (this.connectionID !== NoID && this.canManageDalConfig) {
        this.config = {
          dal: {
            type: 'corteza::dal:connection:dsn',
            params: { dsn: '' },
            modelIdent: '',
            modelIdentCheck: [],
          },
          ...dc.config,
        }
      }
    }

    if (dc?.issues) {
      this.issues = []
      for (const i of dc.issues) {
        this.issues.push(i)
      }
    }

    if (dc?.labels) {
      this.labels = []
      for (const l of dc.labels) {
        this.labels.push(l)
      }
    }
  }
}
