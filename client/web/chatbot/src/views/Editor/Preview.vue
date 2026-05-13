<template>
  <div ref="hostRef" class="hb-preview-host" />
</template>

<script setup>
import { inject, onBeforeUnmount, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { Engine } from 'human-webapp-chatbot-widget/engine'
import { WidgetUI } from 'human-webapp-chatbot-widget/ui'

const props = defineProps({
  chatbot: { type: Object, required: true },
})

const { t } = useI18n()
const $SystemAPI = inject('$SystemAPI')

const hostRef = ref(null)
let ui = null
let engine = null
let conversationID = null
let sending = false
let userClosed = false

function apiOrigin() {
  try {
    return new URL($SystemAPI.baseURL).origin
  } catch {
    return window.location.origin
  }
}

function absolutize(u) {
  if (!u) return u
  if (/^(https?:|blob:|data:)/i.test(u)) return u
  if (u.startsWith('/')) return apiOrigin() + u
  return u
}

function buildConfig() {
  const cb = props.chatbot
  const scenarios = (cb.scenarios || []).map(s => ({
    id: s.id,
    name: s.name,
    type: s.type || 'conversation',
    agentID: s.agentID,
    config: s.config || {},
  }))
  const styling = {
    ...(cb.styling || {}),
    logoURL: absolutize(cb.styling?.logoURL),
    launcher: {
      ...((cb.styling && cb.styling.launcher) || {}),
      iconURL: absolutize(cb.styling?.launcher?.iconURL),
      position: 'bottom-right',
    },
  }
  return {
    styling,
    scenarios,
    handoff: { enabled: !!cb.handoff?.enabled, notImplemented: true },
  }
}

async function handleUserInput(text) {
  if (!engine) return
  engine.pushMessage({ role: 'user', content: text })

  const current = engine.current
  const agentID = current?.agentID
  if (!agentID) {
    engine.pushMessage({
      role: 'system',
      content: t('chatbot.editor.preview.selectAgentHint'),
    })
    return
  }
  if (sending) return
  sending = true
  engine.emit({ type: 'typing', on: true })

  try {
    const res = await $SystemAPI.agentExec({
      agentID,
      input: text,
      ...(conversationID ? { conversationID } : {}),
    })
    if (res?.conversationID) conversationID = res.conversationID
    const out =
      res?.output ||
      (typeof res === 'string' ? res : res?.response?.text || '')
    engine.emit({ type: 'typing', on: false })
    if (out) engine.pushMessage({ role: 'agent', content: out })
  } catch (err) {
    engine.emit({ type: 'typing', on: false })
    engine.pushMessage({
      role: 'system',
      content: t('chatbot.editor.preview.requestFailed', {
        reason: err?.message || t('chatbot.editor.preview.requestFailedReason'),
      }),
    })
  } finally {
    sending = false
  }
}

function mount() {
  if (!hostRef.value) return
  destroy()
  const cfg = buildConfig()
  conversationID = null
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
  ui.onFormSubmit = () => {
    const idx = cfg.scenarios.findIndex(s => s.id === engine.current?.id)
    const next = idx >= 0 ? cfg.scenarios[idx + 1] : null
    if (next) {
      engine.advance(next.id)
    } else {
      engine.emit({ type: 'message', message: { role: 'system', content: '✓ Done' } })
    }
  }
  engine.start()
}

function destroy() {
  if (ui) {
    ui.destroy()
    ui = null
  }
  engine = null
  conversationID = null
}

watch(
  () => JSON.stringify(props.chatbot),
  () => mount(),
  { immediate: false },
)

watch(hostRef, host => {
  if (host) mount()
})

onBeforeUnmount(destroy)
</script>

<style scoped>
.hb-preview-host {
  position: relative;
  width: 100%;
  height: 100%;
  overflow: hidden;
}
</style>
