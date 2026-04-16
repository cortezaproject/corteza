import { Apply, CortezaID, ISO8601Date, NoID } from '../../cast'
import { IsOf } from '../../guards'

interface AgentMeta {
  short: string
  description: string
  sidebarRoles: string[]
}

interface AgentBehaviorKnowledgeBase {
  id: string
  [key: string]: unknown
}

interface AgentBehavior {
  systemPrompt: string
  guardrails: string[]
  injectSystemContext: boolean
  knowledgeBases: AgentBehaviorKnowledgeBase[]
  treatyCLEnabled: boolean
  tclTemperature: number
  tclArticles: string[]
}

interface AgentExecutionModel {
  llmProviderID: string
  model: string
  temperature: number
}

interface AgentExecutionLimits {
  maxIterations: number
  contextWindow: number
  outputTokens: number
  timeout: string
  softLimitRatio: number
}

interface AgentExecution {
  model: AgentExecutionModel
  limits: AgentExecutionLimits
}

interface AgentAccessAllow {
  namespaceID: string
  moduleIDs: string[]
}

interface AgentAccessTool {
  name: string
  description: string
  allow: AgentAccessAllow[] | null
  context?: {
    defaults?: Record<string, unknown>
    overrides?: Record<string, unknown>
  }
}

interface AgentAccessTAQ {
  id: string
  description: string
  params?: Record<string, string>
}

interface AgentAccessWorkflow {
  id: string
  description: string
}

interface AgentAccessContext {
  namespace: string
  module: string
  defaults?: Record<string, unknown>
}

interface AgentAccess {
  context: AgentAccessContext
  tools: AgentAccessTool[]
  taqs: AgentAccessTAQ[]
  workflows: AgentAccessWorkflow[]
}

interface AgentInvocationUser {
  enabled: boolean
}

interface AgentInvocationSystem {
  enabled: boolean
  serviceAccount: string
  inputSchema: unknown
  outputFormat: string
}

interface AgentInvocation {
  user: AgentInvocationUser
  system: AgentInvocationSystem
}

interface PartialAgent extends Partial<
  Omit<Agent, 'createdAt' | 'updatedAt' | 'deletedAt'>
> {
  createdAt?: string | number | Date
  updatedAt?: string | number | Date
  deletedAt?: string | number | Date
}

export class Agent {
  public agentID = NoID
  public handle = ''
  public status = 'active'
  public revision = 0
  public labels: Record<string, string> = {}

  public meta: AgentMeta = {
    short: '',
    description: '',
    sidebarRoles: [],
  }

  public behavior: AgentBehavior = {
    systemPrompt: '',
    guardrails: [],
    injectSystemContext: true,
    knowledgeBases: [],
    treatyCLEnabled: false,
    tclTemperature: 5,
    tclArticles: [],
  }

  public execution: AgentExecution = {
    model: {
      llmProviderID: NoID,
      model: '',
      temperature: 0.7,
    },
    limits: {
      maxIterations: 10,
      contextWindow: 10000,
      outputTokens: 500,
      timeout: '30s',
      softLimitRatio: 0.8,
    },
  }

  public access: AgentAccess = {
    context: {
      namespace: '',
      module: '',
    },
    tools: [],
    taqs: [],
    workflows: [],
  }

  public invocation: AgentInvocation = {
    user: { enabled: true },
    system: {
      enabled: false,
      serviceAccount: NoID,
      inputSchema: null,
      outputFormat: '',
    },
  }

  public createdAt?: Date = undefined
  public updatedAt?: Date = undefined
  public deletedAt?: Date = undefined

  public createdBy = NoID
  public updatedBy = NoID
  public deletedBy = NoID

  public canGrant = false
  public canUpdateAgent = false
  public canDeleteAgent = false

  constructor(o?: PartialAgent) {
    this.apply(o)
  }

  apply(o?: PartialAgent): void {
    Apply(this, o, CortezaID, 'agentID', 'createdBy', 'updatedBy', 'deletedBy')
    Apply(this, o, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')
    Apply(this, o, String, 'handle', 'status')
    Apply(this, o, Number, 'revision')
    Apply(this, o, Boolean, 'canGrant', 'canUpdateAgent', 'canDeleteAgent')

    if (IsOf(o, 'labels')) {
      this.labels = { ...o.labels }
    }

    if (IsOf(o, 'meta')) {
      this.meta = {
        ...this.meta,
        ...o.meta,
        sidebarRoles: Array.isArray(o.meta?.sidebarRoles) ? o.meta.sidebarRoles : this.meta.sidebarRoles,
      }
    }

    if (IsOf(o, 'behavior')) {
      this.behavior = {
        ...this.behavior,
        ...o.behavior,
        guardrails: Array.isArray(o.behavior?.guardrails) ? o.behavior.guardrails : this.behavior.guardrails,
        knowledgeBases: Array.isArray(o.behavior?.knowledgeBases) ? o.behavior.knowledgeBases : this.behavior.knowledgeBases,
        tclArticles: Array.isArray(o.behavior?.tclArticles) ? o.behavior.tclArticles : this.behavior.tclArticles,
      }
    }

    if (IsOf(o, 'execution')) {
      this.execution = {
        model: { ...this.execution.model, ...(o.execution?.model || {}) },
        limits: { ...this.execution.limits, ...(o.execution?.limits || {}) },
      }
    }

    if (IsOf(o, 'access')) {
      this.access = {
        context: { ...this.access.context, ...(o.access?.context || {}) },
        tools: Array.isArray(o.access?.tools) ? o.access.tools : this.access.tools,
        taqs: Array.isArray(o.access?.taqs) ? o.access.taqs : this.access.taqs,
        workflows: Array.isArray(o.access?.workflows) ? o.access.workflows : this.access.workflows,
      }
    }

    if (IsOf(o, 'invocation')) {
      this.invocation = {
        user: { ...this.invocation.user, ...(o.invocation?.user || {}) },
        system: { ...this.invocation.system, ...(o.invocation?.system || {}) },
      }
    }
  }

  clone(): Agent {
    return new Agent(JSON.parse(JSON.stringify(this)))
  }
}
