<template>
  <div ref="hostRef" class="hb-preview-host" />
</template>

<script setup>
import { inject, onBeforeUnmount, ref, watch } from 'vue'
import { Engine } from 'human-webapp-chatbot-widget/engine'
import { WidgetUI } from 'human-webapp-chatbot-widget/ui'
import { PreviewClient } from '@/sections/chatbot/composables/usePreviewClient'

const props = defineProps({
  chatbot: { type: Object, required: true },
})

const $SystemAPI = inject('$SystemAPI')

const hostRef = ref(null)
let ui = null
let engine = null
let client = null
let session = null
let stream = null
let userClosed = false
let autoCloseTimer = null

// Drop empty-string numeric ID fields; the server type uses uint64 with
// `,string,omitempty` which rejects "" during decode.
function cleanNumericID(v) {
  return v === '' || v == null ? undefined : v
}

function buildConfig() {
  const cb = props.chatbot
  const scenarios = (cb.scenarios || []).map(s => ({
    id: s.id,
    name: s.name,
    type: s.type || 'conversation',
    agentID: cleanNumericID(s.agentID),
    config: s.config || {},
    automation: s.automation,
  }))
  const stylingSrc = cb.styling || {}
  const launcherSrc = stylingSrc.launcher || {}
  const styling = {
    ...stylingSrc,
    logoAttachmentID: cleanNumericID(stylingSrc.logoAttachmentID),
    launcher: {
      ...launcherSrc,
      iconAttachmentID: cleanNumericID(launcherSrc.iconAttachmentID),
      position: 'bottom-right',
    },
  }
  // Identity is propagated so the unified inbox can resolve preview sessions
  // back to their source chatbot (otherwise BE decode lands at `Chatbot.ID=0`
  // and the inbox can't look up the name).
  return {
    chatbotID: cleanNumericID(cb.chatbotID),
    handle: cb.handle || '',
    name: cb.name || '',
    styling,
    scenarios,
    handoff: { enabled: !!cb.handoff?.enabled, notImplemented: false },
  }
}

function clearAutoClose() {
  if (autoCloseTimer) {
    clearTimeout(autoCloseTimer)
    autoCloseTimer = null
  }
}

function parseEvent(ev) {
  try {
    return JSON.parse(ev.data || '{}')
  } catch {
    return null
  }
}

function attachStream(s) {
  s.addEventListener('open', () => console.debug('[preview] SSE open'))
  s.addEventListener('error', e => console.warn('[preview] SSE error', e))
  s.addEventListener('token', ev => {
    const data = parseEvent(ev)
    const text = data?.text || data?.token || ''
    if (!text) return
    engine.emit({ type: 'typing', on: false })
    engine.appendAgentDelta(text)
  })
  s.addEventListener('done', () => {
    engine.emit({ type: 'typing', on: false })
    engine.endAgent()
  })
  s.addEventListener('agent_error', ev => {
    const p = parseEvent(ev)
    engine.emit({ type: 'typing', on: false })
    engine.endAgent()
    engine.emit({ type: 'error', error: p?.error || 'agent error' })
  })
  s.addEventListener('step_start', ev => {
    const p = parseEvent(ev)
    console.debug('[preview] step_start', p)
    if (!p) return
    engine.handleStepStart(p)
    if (p.type === 'static_message') {
      const scenarios = props.chatbot.scenarios || []
      const isLast = p.scenarioIndex >= scenarios.length - 1
      const c = p.config || {}
      clearAutoClose()
      if (!isLast && typeof c.autoAdvanceMs === 'number' && c.autoAdvanceMs >= 0) {
        autoCloseTimer = setTimeout(() => void handleEndConversation(), c.autoAdvanceMs)
      } else if (isLast && typeof c.autoCloseAfterMs === 'number' && c.autoCloseAfterMs > 0) {
        autoCloseTimer = setTimeout(() => void closeSession(), c.autoCloseAfterMs)
      }
    }
  })
  s.addEventListener('step_complete', ev => {
    const p = parseEvent(ev)
    if (p) engine.handleStepComplete(p.scenarioID, p.scenarioIndex)
  })
  s.addEventListener('form_error', ev => {
    const p = parseEvent(ev)
    if (p) engine.handleFormError(p.scenarioID, p.errors || {})
  })
  s.addEventListener('handoff_requested', ev => {
    const p = parseEvent(ev)
    if (p) engine.handleHandoffRequested(p.handoffID)
  })
  s.addEventListener('handoff_active', ev => {
    const p = parseEvent(ev)
    if (p) engine.handleHandoffActive(p.handoffID, p.operator)
  })
  s.addEventListener('handoff_complete', () => {
    engine.handleHandoffComplete()
  })
  s.addEventListener('operator_message', ev => {
    const p = parseEvent(ev)
    if (p) engine.pushMessage({ role: 'operator', content: p.content, operator: p.operator })
  })
  s.addEventListener('session_closed', () => {
    clearAutoClose()
    engine.handleSessionClosed()
  })
  s.addEventListener('step_hook_skipped', ev => {
    const p = parseEvent(ev)
    if (!p) return
    engine.emit({
      type: 'message',
      message: {
        role: 'system',
        content: `[preview] ${p.phase} hook ${p.resource} (async=${p.async}) — not executed`,
      },
    })
  })
  s.onerror = () => {
    engine.emit({ type: 'typing', on: false })
    engine.endAgent()
  }
}

async function openSession() {
  if (session) return session
  const cfg = buildConfig()
  // Forward the full snapshot (identity + config) so the BE-side preview
  // session carries the source chatbot's ID / name / handle. Without these
  // the unified inbox can't resolve preview rows back to a chatbot.
  const s = await client.openSession(cfg)
  session = s
  stream = client.openStream(s.sessionID)
  attachStream(stream)

  // Wait until the EventSource has actually connected (onopen fired) before
  // asking the server to emit the first step_start — the bus has no
  // per-subscriber buffering, so an early emit would be silently dropped.
  await waitForStreamOpen(stream)
  try {
    await client.startSession(s.sessionID)
  } catch (err) {
    engine?.emit({ type: 'error', error: err?.message || 'start failed' })
  }
  return session
}

function waitForStreamOpen(es) {
  return new Promise(resolve => {
    if (es.readyState === EventSource.OPEN) {
      resolve()
      return
    }
    const done = () => {
      es.removeEventListener('open', done)
      resolve()
    }
    es.addEventListener('open', done)
    // Hard cap; never block startup longer than 1s even on slow proxies.
    setTimeout(done, 1000)
  })
}

async function closeSession() {
  if (!session) return
  try {
    await client.closeSession(session.sessionID)
  } catch {
    /* ignore */
  }
}

async function handleUserInput(text) {
  if (!engine || !session) return
  engine.pushMessage({ role: 'user', content: text })
  const handoff = engine.handoffState.phase
  if (handoff === 'idle') engine.emit({ type: 'typing', on: true })
  try {
    await client.sendMessage(session.sessionID, text)
  } catch (err) {
    engine.emit({ type: 'typing', on: false })
    engine.emit({ type: 'error', error: err?.message || 'send failed' })
  }
}

async function handleFormSubmit(values) {
  if (!session) return
  try {
    const errors = await client.submitForm(session.sessionID, values)
    if (errors) engine.handleFormError(engine.current?.scenarioID || '', errors)
  } catch (err) {
    engine.emit({ type: 'error', error: err?.message || 'submit failed' })
  }
}

async function handleConsentDecision(accepted) {
  if (!session) return
  try {
    await client.submitConsent(session.sessionID, accepted)
  } catch (err) {
    engine.emit({ type: 'error', error: err?.message || 'consent failed' })
  }
}

async function handleRequestHandoff() {
  if (!session) return
  try {
    await client.requestHandoff(session.sessionID)
  } catch (err) {
    engine.emit({ type: 'error', error: err?.message || 'handoff failed' })
  }
}

async function handleCancelHandoff() {
  if (!session) return
  try {
    await client.closeHandoff(session.sessionID)
  } catch (err) {
    engine.emit({ type: 'error', error: err?.message || 'cancel failed' })
  }
}

async function handleEndConversation() {
  if (!session) return
  try {
    await client.advanceStep(session.sessionID)
  } catch (err) {
    engine.emit({ type: 'error', error: err?.message || 'end failed' })
  }
}

function mount() {
  if (!hostRef.value) return
  destroy()
  const cfg = buildConfig()
  engine = new Engine(cfg)
  ui = new WidgetUI(cfg, engine, {
    container: hostRef.value,
    contained: true,
    startOpen: !userClosed,
    onToggle: open => {
      userClosed = !open
    },
  })
  ui.onUserInput = handleUserInput
  ui.onFormSubmit = handleFormSubmit
  ui.onConsentDecision = handleConsentDecision
  ui.onRequestHandoff = handleRequestHandoff
  ui.onCancelHandoff = handleCancelHandoff
  ui.onEndConversation = handleEndConversation
  ui.onCloseSession = closeSession

  client = new PreviewClient($SystemAPI)
  void openSession()
}

function destroy() {
  clearAutoClose()
  if (stream) {
    stream.close()
    stream = null
  }
  if (session && client) {
    void client.closeSession(session.sessionID).catch(() => {})
  }
  session = null
  if (ui) {
    ui.destroy()
    ui = null
  }
  engine = null
  client = null
}

// Re-mount on chatbot edits, but coalesce rapid keystrokes/sliders so each
// change doesn't tear down + reopen a live SSE session.
const remountDelayMs = 400
let remountTimer = null

function scheduleRemount() {
  if (remountTimer) clearTimeout(remountTimer)
  remountTimer = setTimeout(() => {
    remountTimer = null
    mount()
  }, remountDelayMs)
}

watch(
  () => JSON.stringify(props.chatbot),
  () => scheduleRemount(),
  { immediate: false },
)

watch(hostRef, host => {
  if (host) mount()
})

onBeforeUnmount(() => {
  if (remountTimer) clearTimeout(remountTimer)
  destroy()
})
</script>

<style scoped>
.hb-preview-host {
  position: relative;
  width: 100%;
  height: 100%;
  /* No overflow clipping here so the widget panel's box-shadow can extend
     past the host into the gap between this column and the config card.
     The panel itself still has its own overflow:hidden so message content
     stays clipped to its rounded box. */
}
</style>
