<template>
  <Dialog
    v-model:visible="opened"
    modal
    :header="currentTitle"
    :style="{ width: 'min(42rem, calc(100vw - 2rem))' }"
  >
    <component
      :is="current.component"
      v-if="current"
      :payload="current.prompt.payload"
      :loading="store.isLoading"
      @submit="handleSubmit"
    />

    <div v-else class="flex flex-col gap-2">
      <Button
        v-for="entry in list"
        :key="entry.prompt.stateID"
        severity="secondary"
        variant="text"
        class="!justify-start"
        @click="store.activate(entry.prompt)"
      >
        <span class="flex items-baseline gap-2 min-w-0">
          <span class="truncate">{{ entry.title }}</span>
          <time
            v-if="entry.age"
            class="text-muted-color text-sm shrink-0"
            :datetime="entry.prompt.createdAt"
          >
            {{ entry.age }}
          </time>
        </span>
      </Button>
    </div>

    <template v-if="current" #footer>
      <Button
        severity="secondary"
        variant="text"
        :label="tF('general.label.back', 'Back to list')"
        @click="store.activate(true)"
      />
    </template>
  </Dialog>
</template>

<script setup>
import { computed, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import { useWorkflowPromptsStore } from '../../stores/useWorkflowPromptsStore'
import definitions from './kinds'
import { pVal } from './utils'

const store = useWorkflowPromptsStore()
const $AutomationAPI = inject('$AutomationAPI')
const { t, te } = useI18n()
const tF = (key, fallback) => te(key) ? t(key) : fallback

const opened = computed({
  get: () => store.isActive,
  set: value => {
    if (!value) {
      store.deactivate()
    } else {
      store.activate()
    }
  },
})

const list = computed(() => {
  return store.prompts
    .filter(({ ref }) => !!definitions[ref]?.component && !definitions[ref]?.passive)
    .map(prompt => ({
      ...definitions[prompt.ref],
      prompt,
      title: pVal(prompt.payload, 'title', tF('prompt.title.single', 'Workflow prompt')),
      age: relativeTime(prompt.createdAt),
    }))
})

const current = computed(() => {
  if (!store.current) {
    return undefined
  }

  return list.value.find(({ prompt }) => prompt.stateID === store.current?.stateID)
})

const currentTitle = computed(() => current.value?.title || tF('prompt.title.list', 'Workflow prompts'))

async function handleSubmit(input) {
  if (!$AutomationAPI || !current.value) {
    return
  }

  await store.resume($AutomationAPI, current.value.prompt, input)
}

function relativeTime(value) {
  if (!value) return ''
  const then = new Date(value).getTime()
  if (Number.isNaN(then)) return ''
  const diffSec = Math.round((then - Date.now()) / 1000)
  const abs = Math.abs(diffSec)
  const rtf = new Intl.RelativeTimeFormat(undefined, { numeric: 'auto' })
  if (abs < 60) return rtf.format(diffSec, 'second')
  if (abs < 3600) return rtf.format(Math.round(diffSec / 60), 'minute')
  if (abs < 86400) return rtf.format(Math.round(diffSec / 3600), 'hour')
  return rtf.format(Math.round(diffSec / 86400), 'day')
}
</script>
