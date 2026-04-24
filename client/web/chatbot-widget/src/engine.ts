import type { ChatbotConfig, Scenario } from './types'

export type Message = { role: 'user' | 'agent' | 'system'; content: string }

export type EngineEvent =
  | { type: 'scenario'; scenario: Scenario }
  | { type: 'message'; message: Message }
  | { type: 'typing'; on: boolean }
  | { type: 'agent_start' }
  | { type: 'agent_delta'; text: string }
  | { type: 'agent_end' }
  | { type: 'error'; error: string }

type Handler = (ev: EngineEvent) => void

// Engine is a tiny event-bus + scenario state machine. The DOM layer
// subscribes; the transport layer (api + SSE) calls into `userInput` /
// `advance`.
export class Engine {
  private cfg: ChatbotConfig
  private currentID: string
  private handlers: Handler[] = []
  private agentOpen = false
  private autoAdvanceTimer: ReturnType<typeof setTimeout> | null = null
  messages: Message[] = []

  constructor(cfg: ChatbotConfig) {
    this.cfg = cfg
    this.currentID = cfg.scenarios[0]?.id || ''
  }

  on(h: Handler) {
    this.handlers.push(h)
  }

  emit(ev: EngineEvent) {
    this.handlers.forEach(h => h(ev))
  }

  get current(): Scenario | null {
    return this.cfg.scenarios.find(s => s.id === this.currentID) || null
  }

  start() {
    const s = this.current
    if (s) this.enterScenario(s)
  }

  advance(nextID?: string) {
    if (nextID) {
      const s = this.cfg.scenarios.find(x => x.id === nextID)
      if (!s) return
      this.currentID = s.id
      this.enterScenario(s)
      return
    }
    const idx = this.cfg.scenarios.findIndex(x => x.id === this.currentID)
    const nxt = idx >= 0 ? this.cfg.scenarios[idx + 1] : null
    if (!nxt) return
    this.currentID = nxt.id
    this.enterScenario(nxt)
  }

  private enterScenario(s: Scenario) {
    if (this.autoAdvanceTimer) {
      clearTimeout(this.autoAdvanceTimer)
      this.autoAdvanceTimer = null
    }
    this.emit({ type: 'scenario', scenario: s })
    const idx = this.cfg.scenarios.findIndex(x => x.id === s.id)
    const hasNext = idx >= 0 && idx + 1 < this.cfg.scenarios.length
    if (s.type === 'static_message' && hasNext) {
      const ms = Number((s.config as { autoAdvanceMs?: number } | undefined)?.autoAdvanceMs)
      const delay = Number.isFinite(ms) && ms >= 0 ? ms : 500
      this.autoAdvanceTimer = setTimeout(() => {
        this.autoAdvanceTimer = null
        this.advance()
      }, delay)
    }
  }

  pushMessage(m: Message) {
    if (m.role === 'user' && this.agentOpen) this.endAgent()
    this.messages.push(m)
    this.emit({ type: 'message', message: m })
  }

  startAgent() {
    if (this.agentOpen) return
    this.agentOpen = true
    this.emit({ type: 'agent_start' })
  }

  appendAgentDelta(text: string) {
    if (!this.agentOpen) this.startAgent()
    this.emit({ type: 'agent_delta', text })
  }

  endAgent() {
    if (!this.agentOpen) return
    this.agentOpen = false
    this.emit({ type: 'agent_end' })
  }
}
