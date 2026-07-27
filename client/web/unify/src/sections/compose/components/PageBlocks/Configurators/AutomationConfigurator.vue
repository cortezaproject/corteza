<template>
  <div class="flex flex-col gap-3">
    <AutomationButtonsEditor :buttons="buttons" @update:buttons="onButtonsUpdate" />
    <InterpolationFootnote :is-record-page="isRecordPage" />
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import AutomationButtonsEditor from '../Shared/AutomationButtonsEditor.vue'
import InterpolationFootnote from '@/sections/compose/components/Common/InterpolationFootnote.vue'

const props = defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const isRecordPage = computed(() => !!props.page?.moduleID && props.page.moduleID !== '0')

const block = inject('blockDraft')

const buttons = computed(() => block.value.options?.buttons || [])

function onButtonsUpdate(next) {
  if (!block.value.options) block.value.options = {}
  block.value.options.buttons = next
}
</script>
