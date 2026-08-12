import type { ChatbotConfig, Scenario, StepStartPayload } from './types'

export type Message =
  | { role: 'user'; content: string }
  | { role: 'agent'; content: string }
  | { role: 'operator'; content: string; operator?: string }
  | { role: 'system'; content: string }

// HandoffState tracks the visible state of an operator handover. The engine
// does not orchestrate the handoff itself — it only reflects what the server
// reports via SSE events.
export type HandoffState =
  | { phase: 'idle' }
  | { phase: 'requested'; handoffID: string }
  | { phase: 'active'; handoffID: string; operator?: string }

export type EngineEvent =
  | { type: 'step_start'; payload: StepStartPayload }
  | { type: 'step_complete'; scenarioID: string; scenarioIndex: number }
  | { type: 'message'; message: Message }
  | { type: 'typing'; on: boolean }
  | { type: 'agent_start' }
  | { type: 'agent_delta'; text: string }
  | { type: 'agent_end' }
  | { type: 'form_error'; scenarioID: string; errors: Record<string, string> }
  | { type: 'handoff_change'; state: HandoffState }
  | { type: 'session_closed' }
  | { type: 'error'; error: string }

type Handler = (ev: EngineEvent) => void

// Engine is a thin pub/sub. It holds the *current step* state derived from the
// last `step_start` event delivered by the server. The DOM layer subscribes;
// the transport layer (api + SSE) calls into `userInput` / `submitForm` etc.
// via the WidgetUI callbacks.
export class Engine {
  private cfg: ChatbotConfig
  private handlers: Handler[] = []
  private agentOpen = false
  private currentStep: StepStartPayload | null = null
  private handoff: HandoffState = { phase: 'idle' }
  messages: Message[] = []

  constructor(cfg: ChatbotConfig) {
    this.cfg = cfg
  }

  on(h: Handler) {
    this.handlers.push(h)
  }

  emit(ev: EngineEvent) {
    this.handlers.forEach(h => h(ev))
  }

  get current(): StepStartPayload | null {
    return this.currentStep
  }

  get handoffState(): HandoffState {
    return this.handoff
  }

  // Lookup a scenario by ID from the public config (used for static fallbacks
  // when no step_start payload has been received yet, e.g. legacy clients).
  scenarioByID(id: string): Scenario | undefined {
    return this.cfg.scenarios.find(s => s.id === id)
  }

  handleStepStart(p: StepStartPayload) {
    this.currentStep = p
    if (this.agentOpen) this.endAgent()
    this.emit({ type: 'step_start', payload: p })
  }

  handleStepComplete(scenarioID: string, scenarioIndex: number) {
    this.emit({ type: 'step_complete', scenarioID, scenarioIndex })
  }

  handleFormError(scenarioID: string, errors: Record<string, string>) {
    this.emit({ type: 'form_error', scenarioID, errors })
  }

  handleHandoffRequested(handoffID: string) {
    this.handoff = { phase: 'requested', handoffID }
    this.emit({ type: 'handoff_change', state: this.handoff })
  }

  handleHandoffActive(handoffID: string, operator?: string) {
    this.handoff = { phase: 'active', handoffID, operator }
    this.emit({ type: 'handoff_change', state: this.handoff })
  }

  handleHandoffComplete() {
    // If an operator was actually connected, surface a system message so the
    // visitor sees the conversation ended (otherwise the widget silently
    // hands control back to the agent and the transition is opaque).
    const wasActive = this.handoff.phase === 'active'
    this.handoff = { phase: 'idle' }
    this.emit({ type: 'handoff_change', state: this.handoff })
    if (wasActive) {
      this.pushMessage({
        role: 'system',
        content: 'Conversation with operator has ended.',
      })
    }
  }

  handleSessionClosed() {
    if (this.agentOpen) this.endAgent()
    this.emit({ type: 'session_closed' })
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
