<template>
  <div class="flex gap-2">
    <Button
      v-for="(btn, i) in buttons"
      :key="i"
      :label="evaluatedLabel(btn) || '-'"
      :severity="mapVariant(btn.variant)"
      :loading="processingIDs.includes(i)"
      :disabled="processingIDs.includes(i)"
      size="small"
      @click.prevent="handleButton(btn, i)"
    />
  </div>
</template>

<script setup>
import { ref, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import { evaluatePrefilter } from '../../../lib/record-filter'

const { t } = useI18n()

const props = defineProps({
  buttons: { type: Array, default: () => [] },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  module: { type: Object, default: () => null },
  record: { type: Object, default: undefined },
  records: { type: Array, default: () => [] },
  filter: { type: String, default: '' },
})

const emit = defineEmits(['refresh'])

const $toast = inject('$toast', null)
const $auth = inject('$auth', {})
const $AutomationAPI = inject('$AutomationAPI', null)

const processingIDs = ref([])

const variantSeverityMap = {
  primary: undefined,
  secondary: 'secondary',
  success: 'success',
  danger: 'danger',
  warning: 'warn',
  info: 'info',
}
const mapVariant = key => variantSeverityMap[key]

function evaluatedLabel(btn) {
  try {
    const record = props.record
    const user = $auth?.user || {}
    return evaluatePrefilter(btn.label || '', {
      record,
      user,
      recordID: record?.recordID || '0',
      ownerID: record?.ownedBy || '0',
      userID: user?.userID || '0',
    })
  } catch {
    return btn.label
  }
}

function buildInput() {
  const input = {}
  if (props.namespace?.namespaceID) {
    input.namespace = { '@type': 'ComposeNamespace', '@value': props.namespace }
  }
  if (props.page?.pageID) {
    input.page = { '@type': 'ComposePage', '@value': props.page }
  }
  if (props.module?.moduleID) {
    input.module = { '@type': 'ComposeModule', '@value': props.module }
  }
  if (props.record?.recordID) {
    input.record = { '@type': 'ComposeRecord', '@value': props.record }
  }
  if (props.records?.length) {
    input.selected = props.records
  }
  if (props.filter) {
    input.filter = props.filter
  }
  return input
}

async function handleButton(btn, index) {
  if (!$AutomationAPI) return

  processingIDs.value.push(index)

  try {
    const input = buildInput()

    if (btn.automationID) {
      await $AutomationAPI.ngAutomationExec({ automationID: btn.automationID, input })
    } else if (btn.workflowID) {
      await $AutomationAPI.workflowExec({
        workflowID: btn.workflowID,
        stepID: btn.stepID || '0',
        input,
      })
    } else {
      $toast?.toastInfo?.(t('block.automation.noScript'))
      return
    }

    emit('refresh')
  } catch (e) {
    console.error('Automation execution failed:', e)
    $toast?.toastErrorHandler?.(t('block.automation.executionFailed'))(e)
  } finally {
    processingIDs.value = processingIDs.value.filter(id => id !== index)
  }
}
</script>
