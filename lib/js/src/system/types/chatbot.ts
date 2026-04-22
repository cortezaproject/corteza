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
  iconVisible: boolean
  label: string
  buttonLabel: string
  size: string
  shape: string
  position: string
}

interface ChatbotStyling {
  logoURL: string
  fontFamily: string
  fontSizes: ChatbotFontSizes
  colors: ChatbotColors
  launcher: ChatbotLauncher
}

interface ChatbotScenario {
  id: string
  name: string
  type: string
  agentID?: string
  config: unknown
}

interface ChatbotHandoff {
  enabled: boolean
  targetRoles: string[]
}

interface PartialChatbot extends Partial<
  Omit<Chatbot, 'createdAt' | 'updatedAt' | 'deletedAt'>
> {
  createdAt?: string | number | Date
  updatedAt?: string | number | Date
  deletedAt?: string | number | Date
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
    targetRoles: [],
  }

  public styling: ChatbotStyling = {
    logoURL: '',
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
      iconVisible: true,
      label: '',
      buttonLabel: '',
      size: '56px',
      shape: 'circle',
      position: 'bottom-right',
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
        targetRoles: Array.isArray(o.handoff?.targetRoles)
          ? o.handoff.targetRoles
          : this.handoff.targetRoles,
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
          }))
        : this.scenarios
    }
  }

  clone(): Chatbot {
    return new Chatbot(JSON.parse(JSON.stringify(this)))
  }
}
