import { Apply, CortezaID, ISO8601Date, NoID } from '../../cast'
import { IsOf } from '../../guards'

interface PartialConnection extends Partial<
  Omit<Connection, 'createdAt' | 'updatedAt' | 'deletedAt'>
> {
  createdAt?: string | number | Date
  updatedAt?: string | number | Date
  deletedAt?: string | number | Date
}

interface ConnectionMeta {
  short: string
  description: string
  icon: string
  tags: string[]
}

interface ConnectionPlaceholder {
  name: string
  type: string
  description: string
  required: boolean
  default: string
  options?: string[]
}

interface ConnectionTemplate {
  value: string
  placeholders?: ConnectionPlaceholder[]
}

interface ConnectionAuth {
  method: string
  params?: Record<string, ConnectionTemplate>
}

interface ConnectionService {
  baseURL: ConnectionTemplate
  protocol: string
  contentType: string
  headers?: Record<string, ConnectionTemplate>
  auth: ConnectionAuth
}

export class Connection {
  public connectionID = NoID
  public handle = ''
  public revision = 0
  public status = ''

  public meta: ConnectionMeta = {
    short: '',
    description: '',
    icon: '',
    tags: [],
  }

  public service: ConnectionService = {
    baseURL: { value: '' },
    protocol: '',
    contentType: '',
    auth: { method: 'none' },
  }

  public resources: any[] = []
  public standardOperations: any = {}
  public operations: any[] = []
  public derivedParams: any[] = []

  public labels: Record<string, string> = {}

  public canUpdateConnection = false
  public canDeleteConnection = false

  public createdAt?: Date = undefined
  public updatedBy?: string = undefined
  public updatedAt?: Date = undefined
  public deletedBy?: string = undefined
  public deletedAt?: Date = undefined

  constructor(c?: PartialConnection) {
    this.apply(c)
  }

  apply(c?: PartialConnection): void {
    Apply(this, c, CortezaID, 'connectionID')
    Apply(this, c, String, 'handle', 'status', 'updatedBy', 'deletedBy')
    Apply(this, c, Number, 'revision')

    Apply(this, c, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')
    Apply(this, c, Boolean, 'canUpdateConnection', 'canDeleteConnection')

    if (c && IsOf(c, 'meta')) {
      this.meta = { ...this.meta, ...c.meta }
    }

    if (c && IsOf(c, 'service')) {
      this.service = { ...this.service, ...c.service }
    }

    if (c && IsOf(c, 'resources')) {
      this.resources = c.resources || []
    }
    if (c && IsOf(c, 'operations')) {
      this.operations = c.operations || []
    }
    if (c && IsOf(c, 'standardOperations')) {
      this.standardOperations = c.standardOperations || {}
    }
    if (c && IsOf(c, 'derivedParams')) {
      this.derivedParams = c.derivedParams || []
    }

    if (c && IsOf(c, 'labels')) {
      this.labels = { ...c.labels }
    }
  }

  /**
   * Returns resource ID
   */
  get resourceID(): string {
    return `${this.resourceType}:${this.connectionID}`
  }

  /**
   * Resource type
   */
  get resourceType(): string {
    return 'system:connection'
  }

  clone(): Connection {
    return new Connection(JSON.parse(JSON.stringify(this)))
  }
}
