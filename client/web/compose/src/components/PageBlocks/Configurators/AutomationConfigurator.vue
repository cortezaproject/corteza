<template>
  <AutomationButtonsEditor
    :buttons="buttons"
    @update:buttons="onButtonsUpdate"
  />
</template>

<script setup>
import { computed } from 'vue'
import AutomationButtonsEditor from '../Shared/AutomationButtonsEditor.vue'

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const emit = defineEmits(['update:block'])

const buttons = computed(() => props.block.options?.buttons || [])

function onButtonsUpdate(next) {
  emit('update:block', {
    ...props.block,
    options: { ...props.block.options, buttons: next },
  })
}
</script>
