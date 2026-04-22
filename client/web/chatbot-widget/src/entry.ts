import { WidgetAPI } from './api'
import { Engine } from './engine'
import { WidgetUI } from './ui'

// Locate our own <script> tag so we can read attributes + derive API base.
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
  let cfg
  try {
    cfg = await api.config()
    console.log('[human-chatbot] config loaded', cfg)
  } catch (err) {
    console.warn('[human-chatbot] cannot load config', err)
    return
  }

  let engine: Engine
  let ui: WidgetUI
  try {
    engine = new Engine(cfg)
    ui = new WidgetUI(cfg, engine)
    console.log('[human-chatbot] UI mounted')
  } catch (err) {
    console.error('[human-chatbot] UI construction failed', err)
    return
  }

  let session: { sessionID: string; conversationID: string } | null = null
  let stream: EventSource | null = null

  async function ensureSession() {
    if (session) return session
    const s = await api.openSession()
    api.setToken(s.token)
    session = { sessionID: s.sessionID, conversationID: s.conversationID }

    stream = api.openStream(s.sessionID)
    stream.addEventListener('token', ev => {
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
    stream.addEventListener('done', () => {
      engine.emit({ type: 'typing', on: false })
      engine.endAgent()
    })
    stream.onerror = () => {
      engine.emit({ type: 'typing', on: false })
      engine.endAgent()
    }
    return session
  }

  ui.onUserInput = async text => {
    const current = engine.current
    if (!current || current.type !== 'conversation') return
    engine.pushMessage({ role: 'user', content: text })
    engine.emit({ type: 'typing', on: true })
    try {
      const s = await ensureSession()
      await api.sendMessage(s.sessionID, current.id, text)
    } catch (err: any) {
      engine.emit({ type: 'typing', on: false })
      engine.emit({ type: 'error', error: err?.message || 'send failed' })
    }
  }

  ui.onFormSubmit = _values => {
    engine.advance()
  }

  engine.on(ev => {
    if (ev.type !== 'scenario') return
    const s = ev.scenario
    if (s.type !== 'conversation') return
    const prompt = (s.config as { initialPrompt?: string } | undefined)?.initialPrompt
    if (!prompt) return
    ;(async () => {
      engine.emit({ type: 'typing', on: true })
      try {
        const sess = await ensureSession()
        await api.sendMessage(sess.sessionID, s.id, prompt)
      } catch (err: any) {
        engine.emit({ type: 'typing', on: false })
        engine.emit({ type: 'error', error: err?.message || 'send failed' })
      }
    })()
  })

  engine.start()
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', () => {
    void boot()
  })
} else {
  void boot()
}
