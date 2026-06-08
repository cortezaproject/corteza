<template>
  <div v-if="!hideToasts" class="pointer-events-none fixed left-1/2 top-[calc(var(--topbar-height)+1rem)] z-[1200] flex w-[min(28rem,calc(100vw-2rem))] -translate-x-1/2 flex-col gap-3">
    <div
      v-for="entry in toasts"
      :key="entry.prompt.stateID"
      class="pointer-events-auto rounded-xl border border-surface bg-emphasis p-4 shadow-lg"
    >
      <div class="mb-2 flex items-start justify-between gap-3">
        <strong class="min-w-0 break-words">
          {{ getTitle(entry.prompt) }}
        </strong>
        <Button
          icon="pi pi-times"
          severity="secondary"
          variant="text"
          rounded
          @click="handleHide(entry)"
        />
      </div>

      <component
        :is="entry.component"
        v-if="entry.component"
        :payload="entry.prompt.payload"
        :loading="store.isLoading"
        @submit="resumePrompt(entry.prompt, $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, getCurrentInstance, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useWorkflowPromptsStore } from '../../stores/useWorkflowPromptsStore'
import definitions from './kinds'
import { pVal } from './utils'

const props = defineProps({
  hideToasts: {
    type: Boolean,
    default: false,
  },
})

const DEFAULT_TIMEOUT_SEC = 7

const store = useWorkflowPromptsStore()
const instance = getCurrentInstance()
const { t, te } = useI18n()
const tF = (key: string, fallback: string): string => te(key) ? t(key) : fallback

const passivePrompts = ref<any[]>([])
const hasFocus = ref(document.hasFocus())
let observer = 0
const passiveTimers = new Map<string, number>()

const withHandlers = computed(() => {
  return (hasFocus.value ? store.prompts : [])
    .filter(({ ref }) => !!definitions[ref]?.handler)
    .map(prompt => ({ ...definitions[prompt.ref], prompt }))
})

const withComponents = computed(() => {
  return (hasFocus.value ? store.prompts : [])
    .filter(({ ref }) => !!definitions[ref]?.component)
    .map(prompt => ({ ...definitions[prompt.ref], prompt }))
})

const activePrompts = computed(() => withComponents.value.filter(({ passive }: any) => !passive))

const toasts = computed(() => {
  if (props.hideToasts) return []
  const live = store.isActive ? passivePrompts.value : [...passivePrompts.value, ...activePrompts.value]
  return live
})

watch(withHandlers, async handlers => {
  if (!handlers.length) {
    return
  }

  const { handler, prompt } = handlers[0]
  await store.resume(prompt, {})
  await handler?.call(instance?.proxy, prompt.payload)
})

watch(
  withComponents,
  next => {
    next.forEach(entry => {
      if (entry.passive && !passivePrompts.value.some(({ prompt }) => prompt.stateID === entry.prompt.stateID)) {
        passivePrompts.value.push(entry)
        schedulePassiveAutoHide(entry)
      }
    })
  },
  { immediate: true },
)

function getTitle(prompt) {
  return pVal(prompt.payload, 'title', tF('prompt.title.single', 'Workflow prompt'))
}

function schedulePassiveAutoHide(entry) {
  const stateID = entry.prompt.stateID
  if (passiveTimers.has(stateID)) return
  const timeoutSec = pVal(entry.prompt.payload, 'timeout', DEFAULT_TIMEOUT_SEC) as number
  if (!timeoutSec || timeoutSec <= 0) return
  const handle = window.setTimeout(() => {
    passiveTimers.delete(stateID)
    removePassive(entry.prompt)
  }, timeoutSec * 1000)
  passiveTimers.set(stateID, handle)
}

function clearPassiveTimer(stateID: string) {
  const handle = passiveTimers.get(stateID)
  if (handle) {
    window.clearTimeout(handle)
    passiveTimers.delete(stateID)
  }
}

async function removePassive(prompt) {
  passivePrompts.value = passivePrompts.value.filter(({ prompt: p }) => p.stateID !== prompt.stateID)
  await store.clear(prompt)
}

async function resumePrompt(prompt, input) {
  const keep = !!input?.keep
  const payload = keep ? {} : input
  await store.resume(prompt, payload)
}

async function handleHide(entry) {
  clearPassiveTimer(entry.prompt.stateID)
  if (entry.passive) {
    await removePassive(entry.prompt)
    return
  }

  await store.cancel(entry.prompt)
}

onMounted(() => {
  observer = window.setInterval(() => {
    hasFocus.value = document.hasFocus()
  }, 1000)
})

onBeforeUnmount(() => {
  if (observer) {
    window.clearInterval(observer)
  }
  passiveTimers.forEach(handle => window.clearTimeout(handle))
  passiveTimers.clear()
})
</script>
