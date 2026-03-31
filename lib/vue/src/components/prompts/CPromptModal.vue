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

    <div v-else class="flex flex-col gap-3">
      <Button
        v-for="entry in list"
        :key="entry.prompt.stateID"
        severity="secondary"
        variant="text"
        class="!justify-start"
        @click="store.activate(entry.prompt)"
      >
        {{ entry.title }}
      </Button>
    </div>

    <template v-if="current" #footer>
      <Button
        severity="secondary"
        variant="text"
        :label="t('general.label.back')"
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
const { t } = useI18n()

const opened = computed({
  get: () => store.isActive,
  set: value => {
    if (!value) {
      store.deactivate()
    }
  },
})

const list = computed(() => {
  return store.prompts
    .filter(({ ref }) => !!definitions[ref]?.component && !definitions[ref]?.passive)
    .map(prompt => ({
      ...definitions[prompt.ref],
      prompt,
      title: pVal(prompt.payload, 'title', 'Workflow prompt'),
    }))
})

const current = computed(() => {
  if (!store.current) {
    return undefined
  }

  return list.value.find(({ prompt }) => prompt.stateID === store.current?.stateID)
})

const currentTitle = computed(() => current.value?.title || 'Workflow prompts')

async function handleSubmit(input) {
  if (!$AutomationAPI || !current.value) {
    return
  }

  await store.resume($AutomationAPI, current.value.prompt, input)
}
</script>
