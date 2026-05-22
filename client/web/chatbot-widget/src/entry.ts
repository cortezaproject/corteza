import { WidgetAPI } from './api'
import { Engine } from './engine'
import { WidgetUI } from './ui'
import type {
  ChatbotConfig,
  FormErrorPayload,
  HandoffActivePayload,
  HandoffRequestedPayload,
  OperatorMessagePayload,
  Session,
  StepCompletePayload,
  StepStartPayload,
} from './types'

function findSelf(): HTMLScriptElement | null {
  const current = document.currentScript as HTMLScriptElement | null
  if (current && current.getAttribute('data-widget-key')) return current
  const withKey = document.querySelector<HTMLScriptElement>(
    'script[data-widget-key][src*="widget.js"]',
  )
  if (withKey) return withKey
  return document.querySelector<HTMLScriptElement>('script[src*="widget.js"]')
}

async function boot() {
  console.log('[human-chatbot] boot start')
  const script = findSelf()
  if (!script) {
    console.warn('[human-chatbot] no <script> tag found')
    return
  }
  let widgetKey = script.getAttribute('data-widget-key') || ''
  if (!widgetKey && script.src) {
    try {
      widgetKey = new URL(script.src).searchParams.get('k') || ''
    } catch {
      /* ignore */
    }
  }
  if (!widgetKey) {
    console.warn('[human-chatbot] missing data-widget-key (and no ?k= in src)')
    return
  }

  const api = new WidgetAPI(script.src, widgetKey)
  let cfg: ChatbotConfig
  try {
    cfg = await api.config()
    console.log('[human-chatbot] config loaded', cfg)
  } catch (err) {
    console.warn('[human-chatbot] cannot load config', err)
    return
  }

  const absolutize = (u: string) => {
    if (!u) return u
    if (/^https?:\/\//i.test(u)) return u
    if (u.startsWith('/')) return api.origin + u
    return u
  }
  if (cfg?.styling) {
    cfg.styling.logoURL = absolutize(cfg.styling.logoURL)
    if (cfg.styling.launcher) {
      cfg.styling.launcher.iconURL = absolutize(cfg.styling.launcher.iconURL)
    }
  }

  let engine: Engine
  let ui: WidgetUI
  try {
    engine = new Engine(cfg)
    ui = new WidgetUI(cfg, engine, { startOpen: !!cfg.styling?.launcher?.startOpen })
    console.log('[human-chatbot] UI mounted')
  } catch (err) {
    console.error('[human-chatbot] UI construction failed', err)
    return
  }

  let session: Session | null = null
  let stream: EventSource | null = null
  let autoCloseTimer: ReturnType<typeof setTimeout> | null = null

  function clearAutoClose() {
    if (autoCloseTimer) {
      clearTimeout(autoCloseTimer)
      autoCloseTimer = null
    }
  }

  function attachStream(s: EventSource) {
    s.addEventListener('token', ev => {
      try {
        const data = JSON.parse((ev as MessageEvent).data || '{}')
        const text = data.text || data.token || ''
        if (!text) return
        engine.emit({ type: 'typing', on: false })
        engine.appendAgentDelta(text)
      } catch {
        /* ignore */
      }
    })
    s.addEventListener('done', () => {
      engine.emit({ type: 'typing', on: false })
      engine.endAgent()
    })
    s.addEventListener('agent_error', ev => {
      const p = parseEvent<{ error: string }>(ev)
      engine.emit({ type: 'typing', on: false })
      engine.endAgent()
      engine.emit({ type: 'error', error: p?.error || 'agent error' })
    })
    s.addEventListener('step_start', ev => {
      const p = parseEvent<StepStartPayload>(ev)
      if (!p) return
      engine.handleStepStart(p)
      // Static message auto-progression. Non-last steps advance to the next
      // scenario after `autoAdvanceMs`; the last step auto-closes after
      // `autoCloseAfterMs`. Either field is optional.
      if (p.type === 'static_message') {
        const c = (p.config as { autoAdvanceMs?: number; autoCloseAfterMs?: number } | undefined) || {}
        const lastIdx = cfg.scenarios.length - 1
        const isLast = p.scenarioIndex >= lastIdx
        clearAutoClose()
        if (!isLast && typeof c.autoAdvanceMs === 'number' && c.autoAdvanceMs >= 0) {
          autoCloseTimer = setTimeout(() => { void advanceStep() }, c.autoAdvanceMs)
        } else if (isLast && typeof c.autoCloseAfterMs === 'number' && c.autoCloseAfterMs > 0) {
          autoCloseTimer = setTimeout(() => { void closeSession() }, c.autoCloseAfterMs)
        }
      }
    })
    s.addEventListener('step_complete', ev => {
      const p = parseEvent<StepCompletePayload>(ev)
      if (!p) return
      engine.handleStepComplete(p.scenarioID, p.scenarioIndex)
    })
    s.addEventListener('form_error', ev => {
      const p = parseEvent<FormErrorPayload>(ev)
      if (!p) return
      engine.handleFormError(p.scenarioID, p.errors)
    })
    s.addEventListener('handoff_requested', ev => {
      const p = parseEvent<HandoffRequestedPayload>(ev)
      if (!p) return
      engine.handleHandoffRequested(p.handoffID)
    })
    s.addEventListener('handoff_active', ev => {
      const p = parseEvent<HandoffActivePayload>(ev)
      if (!p) return
      engine.handleHandoffActive(p.handoffID, p.operator)
    })
    s.addEventListener('handoff_complete', () => {
      engine.handleHandoffComplete()
    })
    s.addEventListener('operator_message', ev => {
      const p = parseEvent<OperatorMessagePayload>(ev)
      if (!p) return
      engine.pushMessage({ role: 'operator', content: p.content, operator: p.operator })
    })
    s.addEventListener('user_message', _ev => {
      // The user already sees their own message locally (echoed by ui.ts on
      // submit). This event is for the operator console; ignore on the widget.
    })
    s.addEventListener('session_closed', _ev => {
      clearAutoClose()
      engine.handleSessionClosed()
    })
    s.addEventListener('step_error', ev => {
      const p = parseEvent<{ error: string }>(ev)
      engine.emit({ type: 'error', error: p?.error || 'step error' })
    })
    s.onerror = () => {
      engine.emit({ type: 'typing', on: false })
      engine.endAgent()
    }
  }

  async function ensureSession(): Promise<Session> {
    if (session) return session
    const s = await api.openSession()
    api.setToken(s.token)
    session = s
    stream = api.openStream(s.sessionID)
    attachStream(stream)
    return session
  }

  async function closeSession() {
    if (!session) return
    try {
      await api.closeSession(session.sessionID)
    } catch (err) {
      console.warn('[human-chatbot] close session failed', err)
    }
  }

  async function advanceStep() {
    if (!session) return
    try {
      await api.advanceStep(session.sessionID)
    } catch (err) {
      console.warn('[human-chatbot] advance step failed', err)
    }
  }

  ui.onUserInput = async text => {
    const current = engine.current
    if (!current) return
    engine.pushMessage({ role: 'user', content: text })

    // During handoff, the message is relayed to the operator (no agent typing).
    const handoff = engine.handoffState.phase
    if (handoff === 'idle') engine.emit({ type: 'typing', on: true })

    try {
      const s = await ensureSession()
      await api.sendMessage(s.sessionID, text)
    } catch (err: any) {
      engine.emit({ type: 'typing', on: false })
      engine.emit({ type: 'error', error: err?.message || 'send failed' })
    }
  }

  ui.onFormSubmit = async values => {
    try {
      const s = await ensureSession()
      const errors = await api.submitForm(s.sessionID, values)
      if (errors) {
        engine.handleFormError(engine.current?.scenarioID || '', errors)
      }
    } catch (err: any) {
      engine.emit({ type: 'error', error: err?.message || 'submit failed' })
    }
  }

  ui.onConsentDecision = async accepted => {
    try {
      const s = await ensureSession()
      await api.submitConsent(s.sessionID, accepted)
    } catch (err: any) {
      engine.emit({ type: 'error', error: err?.message || 'consent failed' })
    }
  }

  ui.onRequestHandoff = async () => {
    try {
      const s = await ensureSession()
      await api.requestHandoff(s.sessionID)
    } catch (err: any) {
      engine.emit({ type: 'error', error: err?.message || 'handoff failed' })
    }
  }

  ui.onCancelHandoff = async () => {
    const st = engine.handoffState
    if (st.phase === 'idle') return
    try {
      const s = await ensureSession()
      await api.closeHandoff(s.sessionID, st.handoffID)
    } catch (err: any) {
      engine.emit({ type: 'error', error: err?.message || 'cancel failed' })
    }
  }

  ui.onEndConversation = async () => {
    try {
      const s = await ensureSession()
      await api.advanceStep(s.sessionID)
    } catch (err: any) {
      engine.emit({ type: 'error', error: err?.message || 'end failed' })
    }
  }

  ui.onCloseSession = closeSession

  // Open the session immediately so the server can fire step_start for the
  // first scenario over SSE. No client-side scenario tracking.
  void ensureSession()
}

function parseEvent<T>(ev: Event): T | null {
  try {
    return JSON.parse((ev as MessageEvent).data || '{}') as T
  } catch {
    return null
  }
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', () => {
    void boot()
  })
} else {
  void boot()
}
