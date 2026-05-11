import { Apply, HumanID, ISO8601Date, NoID } from '../../cast'
import { IsOf } from '../../guards'
import { Connection } from './connection'

interface ConfiguredConnectionParam {
  scope: string[]
  name: string
  value: string
}

interface ConfiguredConnectionConfig {
  namespaceID: string
  dalConnectionID: string
  credentialID: string
  params: ConfiguredConnectionParam[]
}

interface PartialConfiguredConnection extends Partial<
  Omit<ConfiguredConnection, 'createdAt' | 'updatedAt' | 'deletedAt'>
> {
  createdAt?: string | number | Date
  updatedAt?: string | number | Date
  deletedAt?: string | number | Date
}

export class ConfiguredConnection {
  public configurationID = NoID
  public connectionID = NoID

  public name = ''
  public status = ''

  public connection: Connection = new Connection()
  public config: ConfiguredConnectionConfig = {
    namespaceID: NoID,
    dalConnectionID: NoID,
    credentialID: NoID,
    params: [],
  }

  public labels: Record<string, string> = {}

  public canUpdateConfiguredConnection = false
  public canDeleteConfiguredConnection = false

  public createdAt?: Date = undefined
  public createdBy?: string = undefined
  public updatedAt?: Date = undefined
  public updatedBy?: string = undefined
  public deletedAt?: Date = undefined
  public deletedBy?: string = undefined

  constructor(cc?: PartialConfiguredConnection) {
    this.apply(cc)
  }

  apply(cc?: PartialConfiguredConnection): void {
    Apply(this, cc, HumanID, 'configurationID', 'connectionID')
    Apply(this, cc, String, 'name', 'status', 'createdBy', 'updatedBy', 'deletedBy')

    Apply(this, cc, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')
    Apply(this, cc, Boolean, 'canUpdateConfiguredConnection', 'canDeleteConfiguredConnection')

    if (cc && IsOf(cc, 'connection')) {
      this.connection = new Connection(cc.connection)
    }

    if (cc && IsOf(cc, 'config')) {
      this.config = { ...this.config, ...cc.config }
      if (cc.config?.params) {
        this.config.params = [...cc.config.params]
      }
    }

    if (cc && IsOf(cc, 'labels')) {
      this.labels = { ...cc.labels }
    }
  }

  /**
   * Returns resource ID
   */
  get resourceID(): string {
    return `${this.resourceType}:${this.configurationID}`
  }

  /**
   * Resource type
   */
  get resourceType(): string {
    return 'system:configured-connection'
  }

  clone(): ConfiguredConnection {
    return new ConfiguredConnection(JSON.parse(JSON.stringify(this)))
  }
}
