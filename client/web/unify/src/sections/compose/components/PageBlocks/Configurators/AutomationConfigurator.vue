<template>
  <div class="flex flex-col gap-3">
    <AutomationButtonsEditor :buttons="buttons" @update:buttons="onButtonsUpdate" />
    <CExpressionHint :scope="scope" :insertable="false" />
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import AutomationButtonsEditor from '../Shared/AutomationButtonsEditor.vue'
import { useExpressionScope } from '@/sections/compose/composables/useExpressionScope'

const props = defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const { scope } = useExpressionScope({ page: computed(() => props.page) })

const block = inject('blockDraft')

const buttons = computed(() => block.value.options?.buttons || [])

function onButtonsUpdate(next) {
  if (!block.value.options) block.value.options = {}
  block.value.options.buttons = next
}
</script>
