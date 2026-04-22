<template>
  <div ref="hostRef" class="hb-preview-host" />
</template>

<script setup>
import { inject, onBeforeUnmount, ref, watch } from 'vue'
import { Engine } from 'human-webapp-chatbot/engine'
import { WidgetUI } from 'human-webapp-chatbot/ui'

const props = defineProps({
  chatbot: { type: Object, required: true },
  agentID: { type: String, default: '' },
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

  if (!props.agentID) {
    engine.pushMessage({
      role: 'system',
      content: 'Save the agent to chat with it in the preview.',
    })
    return
  }
  if (sending) return
  sending = true
  engine.emit({ type: 'typing', on: true })

  try {
    const res = await $SystemAPI.agentExec({
      agentID: props.agentID,
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

watch(
  () => props.agentID,
  () => {
    conversationID = null
  },
)

onBeforeUnmount(destroy)
</script>

<style scoped>
.hb-preview-host {
  position: relative;
  flex: 1;
  min-height: 0;
  overflow: hidden;
}
</style>
