import { Apply, HumanID, ISO8601Date, NoID } from '../../cast'
import { IsOf } from '../../guards'

interface ChatbotFontSizes {
  base: string
  small: string
  heading: string
}

interface ChatbotColors {
  primary: string
  primaryText: string
  header: string
  headerText: string
  background: string
  text: string
  userBubble: string
  agentBubble: string
}

interface ChatbotLauncher {
  iconURL: string
  iconAttachmentID: string
  iconVisible: boolean
  label: string
  buttonLabel: string
  size: string
  shape: string
  position: string
  startOpen: boolean
}

interface ChatbotStyling {
  logoURL: string
  logoAttachmentID: string
  fontFamily: string
  fontSizes: ChatbotFontSizes
  colors: ChatbotColors
  launcher: ChatbotLauncher
}

interface ChatbotAutomationHook {
  automation: string
  async: boolean
}

interface ChatbotScenarioAutomation {
  before?: ChatbotAutomationHook
  after?: ChatbotAutomationHook
}

interface ChatbotHandoffAutomation {
  onRequested?: ChatbotAutomationHook
  onAccepted?: ChatbotAutomationHook
}

interface ChatbotScenario {
  id: string
  name: string
  type: string
  agentID?: string
  config: unknown
  automation?: ChatbotScenarioAutomation
}

interface ChatbotHandoff {
  enabled: boolean
  automation?: ChatbotHandoffAutomation
}

interface PartialChatbot extends Partial<Omit<Chatbot, 'createdAt' | 'updatedAt' | 'deletedAt'>> {
  createdAt?: string | number | Date
  updatedAt?: string | number | Date
  deletedAt?: string | number | Date
}

// normalizeHook accepts either the legacy string shape (resource id) or the
// current {automation, async} object shape and yields the canonical object
// shape. Empty / missing inputs become a zero hook.
function normalizeHook(input: unknown): ChatbotAutomationHook {
  if (typeof input === 'string') {
    return { automation: input, async: false }
  }
  if (input && typeof input === 'object') {
    const h = input as Partial<ChatbotAutomationHook>
    return {
      automation: typeof h.automation === 'string' ? h.automation : '',
      async: !!h.async,
    }
  }
  return { automation: '', async: false }
}

function normalizeScenarioAutomation(input: unknown): ChatbotScenarioAutomation {
  const a = (input || {}) as Partial<ChatbotScenarioAutomation>
  return {
    before: normalizeHook(a.before),
    after: normalizeHook(a.after),
  }
}

export class Chatbot {
  public chatbotID = NoID
  public handle = ''
  public name = ''
  public enabled = false
  public widgetKey = ''
  public allowedOrigins: string[] = []
  public sessionTTL = '2h'
  public labels: Record<string, string> = {}

  public handoff: ChatbotHandoff = {
    enabled: false,
  }

  public styling: ChatbotStyling = {
    logoURL: '',
    logoAttachmentID: '',
    fontFamily: '',
    fontSizes: { base: '14px', small: '12px', heading: '16px' },
    colors: {
      primary: '#09344E',
      primaryText: '#ffffff',
      header: '#09344E',
      headerText: '#ffffff',
      background: '#ffffff',
      text: '#111827',
      userBubble: '#09344E',
      agentBubble: '#f4f4f5',
    },
    launcher: {
      iconURL: '',
      iconAttachmentID: '',
      iconVisible: true,
      label: '',
      buttonLabel: '',
      size: '56px',
      shape: 'circle',
      position: 'bottom-right',
      startOpen: false,
    },
  }

  public scenarios: ChatbotScenario[] = []

  public createdAt?: Date = undefined
  public updatedAt?: Date = undefined
  public deletedAt?: Date = undefined

  public createdBy = NoID
  public updatedBy = NoID
  public deletedBy = NoID

  public canGrant = false
  public canUpdateChatbot = false
  public canDeleteChatbot = false

  constructor(o?: PartialChatbot) {
    this.apply(o)
  }

  apply(o?: PartialChatbot): void {
    Apply(this, o, HumanID, 'chatbotID', 'createdBy', 'updatedBy', 'deletedBy')
    Apply(this, o, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')
    Apply(this, o, String, 'handle', 'name', 'widgetKey', 'sessionTTL')
    Apply(this, o, Boolean, 'enabled', 'canGrant', 'canUpdateChatbot', 'canDeleteChatbot')

    if (IsOf(o, 'labels')) {
      this.labels = { ...o.labels }
    }

    if (IsOf(o, 'allowedOrigins')) {
      this.allowedOrigins = Array.isArray(o.allowedOrigins) ? o.allowedOrigins : this.allowedOrigins
    }

    if (IsOf(o, 'handoff')) {
      this.handoff = {
        ...this.handoff,
        ...(o.handoff || {}),
      }
    }

    if (IsOf(o, 'styling')) {
      const s = o.styling || {}
      this.styling = {
        ...this.styling,
        ...s,
        fontSizes: { ...this.styling.fontSizes, ...(s.fontSizes || {}) },
        colors: { ...this.styling.colors, ...(s.colors || {}) },
        launcher: { ...this.styling.launcher, ...(s.launcher || {}) },
      }
    }

    if (IsOf(o, 'scenarios')) {
      this.scenarios = Array.isArray(o.scenarios)
        ? o.scenarios.map(s => ({
            ...s,
            agentID: s.agentID ? HumanID(s.agentID) : undefined,
            automation: normalizeScenarioAutomation(s.automation),
          }))
        : this.scenarios
    }
  }

  clone(): Chatbot {
    return new Chatbot(JSON.parse(JSON.stringify(this)))
  }
}

export interface ChatbotSessionHandoff {
  id: string
  status: string
  startedAt?: Date
  closedAt?: Date
}

interface PartialChatbotSession extends Partial<Omit<ChatbotSession, 'createdAt' | 'updatedAt'>> {
  createdAt?: string | number | Date
  updatedAt?: string | number | Date
}

// ChatbotSession mirrors the unified wire shape returned by
// /chatbots/sessions. `id` is a string because preview-source sessions use an
// opaque (non-numeric) identifier; widget-source sessions use the DB id as a
// numeric string. RBAC flags are populated per-row by the server.
export class ChatbotSession {
  public id = ''
  public source = ''
  public chatbotID = NoID
  public status = ''
  public currentStep = 0
  public conversationID = NoID

  public handoff?: ChatbotSessionHandoff = undefined

  public canView = false
  public canManage = false
  public canManageHandoff = false

  public createdAt?: Date = undefined
  public updatedAt?: Date = undefined

  constructor(o?: PartialChatbotSession) {
    this.apply(o)
  }

  apply(o?: PartialChatbotSession): void {
    if (!o) return
    Apply(this, o, String, 'id', 'source', 'status')
    Apply(this, o, HumanID, 'chatbotID', 'conversationID')
    Apply(this, o, Number, 'currentStep')
    Apply(this, o, Boolean, 'canView', 'canManage', 'canManageHandoff')
    Apply(this, o, ISO8601Date, 'createdAt', 'updatedAt')

    if (IsOf(o, 'handoff') && o.handoff) {
      const h = o.handoff as Partial<ChatbotSessionHandoff> & {
        startedAt?: string | Date
        closedAt?: string | Date
      }
      this.handoff = {
        id: String(h.id || ''),
        status: String(h.status || ''),
        startedAt: h.startedAt ? new Date(h.startedAt) : undefined,
        closedAt: h.closedAt ? new Date(h.closedAt) : undefined,
      }
    }
  }
}

interface PartialChatbotSessionStep extends Partial<Omit<ChatbotSessionStep, 'createdAt'>> {
  createdAt?: string | number | Date
}

export class ChatbotSessionStep {
  public id = ''
  public scenarioIndex = 0
  public scenarioID = ''
  public conversationID = NoID
  public status = ''
  public hookLog: string[] = []

  public createdAt?: Date = undefined

  constructor(o?: PartialChatbotSessionStep) {
    this.apply(o)
  }

  apply(o?: PartialChatbotSessionStep): void {
    if (!o) return
    Apply(this, o, String, 'id', 'scenarioID', 'status')
    Apply(this, o, HumanID, 'conversationID')
    Apply(this, o, Number, 'scenarioIndex')
    Apply(this, o, ISO8601Date, 'createdAt')

    if (IsOf(o, 'hookLog')) {
      this.hookLog = Array.isArray(o.hookLog) ? o.hookLog.map(String) : []
    }
  }
}
