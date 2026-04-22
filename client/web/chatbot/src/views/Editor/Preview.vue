<template>
  <div ref="hostRef" class="hb-preview-host" />
</template>

<script setup>
import { inject, onBeforeUnmount, ref, watch } from 'vue'
import { Engine } from 'human-webapp-chatbot-widget/engine'
import { WidgetUI } from 'human-webapp-chatbot-widget/ui'

const props = defineProps({
  chatbot: { type: Object, required: true },
})

const $SystemAPI = inject('$SystemAPI')

const hostRef = ref(null)
let ui = null
let engine = null
let conversationID = null
let sending = false

function buildConfig() {
  const cb = props.chatbot
  const scenarios = (cb.scenarios || []).map(s => ({
    id: s.id,
    name: s.name,
    type: s.type || 'conversation',
    agentID: s.agentID,
    config: s.config || {},
  }))
  return {
    styling: cb.styling,
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
      content: 'Select an agent on the active conversation scenario to preview it.',
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
      content: 'Error: ' + (err?.message || 'request failed'),
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
    startOpen: true,
  })
  ui.onUserInput = handleUserInput
  ui.onFormSubmit = () => engine.advance()
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
