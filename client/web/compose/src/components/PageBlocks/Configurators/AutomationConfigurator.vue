<template>
  <AutomationButtonsEditor
    :buttons="buttons"
    @update:buttons="onButtonsUpdate"
  />
</template>

<script setup>
import { computed, inject } from 'vue'
import AutomationButtonsEditor from '../Shared/AutomationButtonsEditor.vue'

defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const block = inject('blockDraft')

const buttons = computed(() => block.value.options?.buttons || [])

function onButtonsUpdate(next) {
  if (!block.value.options) block.value.options = {}
  block.value.options.buttons = next
}
</script>
