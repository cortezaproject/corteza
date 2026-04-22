import { Apply, HumanID, ISO8601Date, NoID } from '../../cast'
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
  temperature: number | null
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

interface AgentChatbotFontSizes {
  base: string
  small: string
  heading: string
}

interface AgentChatbotColors {
  primary: string
  primaryText: string
  background: string
  text: string
  userBubble: string
  agentBubble: string
}

interface AgentChatbotLauncher {
  iconURL: string
  iconVisible: boolean
  label: string
  buttonLabel: string
  size: string
  shape: string
}

interface AgentChatbotStyling {
  logoURL: string
  fontFamily: string
  fontSizes: AgentChatbotFontSizes
  colors: AgentChatbotColors
  launcher: AgentChatbotLauncher
}

interface AgentChatbotScenario {
  id: string
  name: string
  type: string
  config: unknown
}

interface AgentChatbotHandoff {
  enabled: boolean
  targetRoles: string[]
}

interface AgentChatbot {
  enabled: boolean
  widgetKey: string
  allowedOrigins: string[]
  sessionTTL: string
  handoff: AgentChatbotHandoff
  styling: AgentChatbotStyling
  scenarios: AgentChatbotScenario[]
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

  public chatbot: AgentChatbot = {
    enabled: false,
    widgetKey: '',
    allowedOrigins: [],
    sessionTTL: '2h',
    handoff: {
      enabled: false,
      targetRoles: [],
    },
    styling: {
      logoURL: '',
      fontFamily: '',
      fontSizes: { base: '14px', small: '12px', heading: '16px' },
      colors: {
        primary: '#09344E',
        primaryText: '#ffffff',
        background: '#ffffff',
        text: '#111827',
        userBubble: '#09344E',
        agentBubble: '#f4f4f5',
      },
      launcher: {
        iconURL: '',
        iconVisible: true,
        label: '',
        buttonLabel: '',
        size: '56px',
        shape: 'circle',
      },
    },
    scenarios: [],
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
    Apply(this, o, HumanID, 'agentID', 'createdBy', 'updatedBy', 'deletedBy')
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
      const m = o.execution?.model || {}
      this.execution = {
        model: {
          ...this.execution.model,
          ...m,
          temperature: 'temperature' in m ? m.temperature : (this.agentID ? null : this.execution.model.temperature),
        },
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

    if (IsOf(o, 'chatbot')) {
      const cb = o.chatbot || {}
      this.chatbot = {
        ...this.chatbot,
        ...cb,
        allowedOrigins: Array.isArray(cb.allowedOrigins) ? cb.allowedOrigins : this.chatbot.allowedOrigins,
        handoff: {
          ...this.chatbot.handoff,
          ...(cb.handoff || {}),
          targetRoles: Array.isArray(cb.handoff?.targetRoles)
            ? cb.handoff.targetRoles
            : this.chatbot.handoff.targetRoles,
        },
        styling: {
          ...this.chatbot.styling,
          ...(cb.styling || {}),
          fontSizes: { ...this.chatbot.styling.fontSizes, ...((cb.styling || {}).fontSizes || {}) },
          colors: { ...this.chatbot.styling.colors, ...((cb.styling || {}).colors || {}) },
          launcher: { ...this.chatbot.styling.launcher, ...((cb.styling || {}).launcher || {}) },
        },
        scenarios: Array.isArray(cb.scenarios) ? cb.scenarios : this.chatbot.scenarios,
      }
    }
  }

  clone(): Agent {
    return new Agent(JSON.parse(JSON.stringify(this)))
  }
}
