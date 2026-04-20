import { Apply, HumanID, ISO8601Date, NoID } from '../../cast'
import { IsOf } from '../../guards'

interface PartialLlmProvider extends Partial<
  Omit<LlmProvider, 'createdAt' | 'updatedAt' | 'deletedAt'>
> {
  createdAt?: string | number | Date
  updatedAt?: string | number | Date
  deletedAt?: string | number | Date
}

interface LlmProviderMeta {
  short: string
  description: string
}

interface LlmProviderConfig {
  promptURL: string
  model: string
  temperature: number
  maxTokens: number
  timeout: string
}

export class LlmProvider {
  public llmProviderID = NoID
  public handle = ''
  public status = ''
  public provider = ''
  public credentialID = NoID

  public meta: LlmProviderMeta = {
    short: '',
    description: '',
  }

  public config: LlmProviderConfig = {
    promptURL: '',
    model: '',
    temperature: 0,
    maxTokens: 0,
    timeout: '',
  }

  public createdAt?: Date = undefined
  public updatedAt?: Date = undefined
  public deletedAt?: Date = undefined

  public createdBy = NoID
  public updatedBy = NoID
  public deletedBy = NoID

  public canUpdateLlmProvider = false
  public canDeleteLlmProvider = false
  public canGrant = false

  constructor(o?: PartialLlmProvider) {
    this.apply(o)
  }

  apply(o?: PartialLlmProvider): void {
    Apply(this, o, HumanID, 'llmProviderID', 'credentialID')
    Apply(this, o, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')
    Apply(this, o, String, 'handle', 'status', 'provider')
    Apply(this, o, Boolean, 'canUpdateLlmProvider', 'canDeleteLlmProvider', 'canGrant')

    if (IsOf(o, 'meta')) {
      this.meta = { ...this.meta, ...o.meta }
    }

    if (IsOf(o, 'config')) {
      this.config = { ...this.config, ...o.config }
    }

    Apply(this, o, HumanID, 'createdBy', 'updatedBy', 'deletedBy')
  }

  clone(): LlmProvider {
    return new LlmProvider(JSON.parse(JSON.stringify(this)))
  }
}
