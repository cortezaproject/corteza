<template>
  <div v-if="!hideToasts" class="pointer-events-none fixed right-4 top-[calc(var(--topbar-height)+1rem)] z-[1200] flex w-[min(28rem,calc(100vw-2rem))] flex-col gap-3">
    <div
      v-for="entry in toasts"
      :key="entry.prompt.stateID"
      class="pointer-events-auto rounded-xl border bg-surface-0 p-4 shadow-lg"
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
import { computed, getCurrentInstance, inject, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useWorkflowPromptsStore } from '../../stores/useWorkflowPromptsStore'
import definitions from './kinds'
import { pVal } from './utils'

const props = defineProps({
  hideToasts: {
    type: Boolean,
    default: false,
  },
})

const store = useWorkflowPromptsStore()
const $AutomationAPI = inject('$AutomationAPI')
const instance = getCurrentInstance()

const passivePrompts = ref([])
const hasFocus = ref(document.hasFocus())
let observer = 0

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

const activePrompts = computed(() => withComponents.value.filter(({ passive }) => !passive))
const toasts = computed(() => props.hideToasts ? [] : [...passivePrompts.value, ...activePrompts.value])

watch(withHandlers, async handlers => {
  if (!handlers.length || !$AutomationAPI) {
    return
  }

  const { handler, prompt } = handlers[0]
  await store.resume($AutomationAPI, prompt, {})
  await handler?.call(instance?.proxy, prompt.payload)
})

watch(
  withComponents,
  next => {
    next.forEach(entry => {
      if (entry.passive && !passivePrompts.value.some(({ prompt }) => prompt.stateID === entry.prompt.stateID)) {
        passivePrompts.value.push(entry)
      }
    })
  },
  { immediate: true },
)

function getTitle(prompt) {
  return pVal(prompt.payload, 'title', 'Workflow prompt')
}

async function resumePrompt(prompt, input) {
  if (!$AutomationAPI) {
    return
  }

  const keep = !!input?.keep
  const payload = keep ? {} : input
  await store.resume($AutomationAPI, prompt, payload)
}

async function handleHide(entry) {
  if (entry.passive) {
    passivePrompts.value = passivePrompts.value.filter(({ prompt }) => prompt.stateID !== entry.prompt.stateID)
    if ($AutomationAPI) {
      await store.clear(entry.prompt)
    }
    return
  }

  if ($AutomationAPI) {
    await store.cancel($AutomationAPI, entry.prompt)
  }
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
})
</script>
