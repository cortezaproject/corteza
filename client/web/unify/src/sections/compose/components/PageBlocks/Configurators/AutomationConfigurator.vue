<template>
  <div class="flex flex-col gap-3">
    <AutomationButtonsEditor
      :buttons="buttons"
      :scope="scope"
      :expr-scope="exprScope"
      :is-record-page="isRecordPage"
      :namespace="namespace"
      :page="page"
      :module="pageModule"
      @update:buttons="onButtonsUpdate"
    />
    <CExpressionHint :scope="scope" :insertable="false" />
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import { useModuleStore } from '@planetcrust/human-vue'
import AutomationButtonsEditor from '../Shared/AutomationButtonsEditor.vue'
import { useExpressionScope } from '@/sections/compose/composables/useExpressionScope'

const props = defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const { scope, exprScope, isRecordPage } = useExpressionScope({
  page: computed(() => props.page),
})

const moduleStore = useModuleStore()

const block = inject('blockDraft')

const buttons = computed(() => block.value.options?.buttons || [])

// The module of the record page the block sits on — what a script's module
// constraints are matched against.
const pageModule = computed(() => {
  const moduleID = props.page?.moduleID
  if (!moduleID || moduleID === '0') return null
  return moduleStore.getByID(moduleID) || null
})

function onButtonsUpdate(next) {
  if (!block.value.options) block.value.options = {}
  block.value.options.buttons = next
}
</script>
