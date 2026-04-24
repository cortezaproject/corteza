export type ScenarioType = 'static_message' | 'conversation' | 'form'

export interface Scenario {
  id: string
  name: string
  type: ScenarioType
  agentID?: string
  config: any
}

export interface Styling {
  logoURL: string
  logoAttachmentID?: string
  fontFamily: string
  fontSizes: { base: string; small: string; heading: string }
  colors: {
    primary: string
    primaryText: string
    header: string
    headerText: string
    background: string
    text: string
    userBubble: string
    agentBubble: string
  }
  launcher: { iconURL: string; iconAttachmentID?: string; iconVisible: boolean; label: string; buttonLabel: string; size: string; shape: string; position: string; startOpen: boolean }
}

export interface ChatbotConfig {
  styling: Styling
  scenarios: Scenario[]
  handoff: { enabled: boolean; notImplemented: boolean }
}

export interface Session {
  sessionID: string
  token: string
  conversationID: string
}
