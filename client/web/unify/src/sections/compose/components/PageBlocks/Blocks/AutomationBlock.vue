<template>
  <PageBlock :block="block">
    <div v-if="buttons.length" class="flex flex-wrap gap-2 p-3">
      <Button
        v-for="(btn, i) in buttons"
        :key="i"
        :label="buttonLabel(btn.label)"
        :severity="mapVariant(btn.variant)"
        :loading="processingIDs.includes(i)"
        :disabled="processingIDs.includes(i)"
        class="flex-auto min-w-[150px] whitespace-normal"
        @click="handleButton(btn, i)"
      />
    </div>
    <div v-else class="flex items-center justify-center h-full p-3 text-muted-color italic">
      {{ $t('block.automation.noScripts') }}
    </div>
  </PageBlock>
</template>

<script setup>
import { computed, ref, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import PageBlock from './PageBlock.vue'
import { evaluatePrefilter } from '../../../lib/record-filter'

const { t } = useI18n()

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const $toast = inject('$toast')
const $auth = inject('$auth', {})
const $AutomationAPI = inject('$AutomationAPI', null)

const processingIDs = ref([])
const buttons = computed(() => props.block.options?.buttons || [])

function mapVariant(variant) {
  const map = {
    primary: undefined,
    secondary: 'secondary',
    light: 'secondary',
    dark: 'contrast',
    success: 'success',
    danger: 'danger',
    warning: 'warn',
    info: 'info',
  }
  return map[variant] || undefined
}

function buttonLabel(label = '') {
  try {
    const record = props.record
    const user = $auth?.user || {}
    return evaluatePrefilter(label, {
      record,
      user,
      recordID: record?.recordID || '0',
      ownerID: record?.ownedBy || '0',
      userID: user?.userID || '0',
    })
  } catch {
    return label
  }
}

async function handleButton(btn, index) {
  processingIDs.value.push(index)

  try {
    if (btn.automationID && $AutomationAPI) {
      // Execute TAQ (NG Automation)
      const input = {}
      if (props.namespace?.namespaceID) {
        input.namespace = { '@type': 'ComposeNamespace', '@value': props.namespace }
      }
      if (props.record?.recordID) {
        input.record = { '@type': 'ComposeRecord', '@value': props.record }
      }

      await $AutomationAPI.ngAutomationExec({
        automationID: btn.automationID,
        input,
      })
    } else if (btn.workflowID && $AutomationAPI) {
      // Execute workflow
      const input = {}

      if (props.namespace?.namespaceID) {
        input.namespace = { '@type': 'ComposeNamespace', '@value': props.namespace }
      }
      if (props.page?.pageID) {
        input.page = { '@type': 'ComposePage', '@value': props.page }
      }
      if (props.record?.recordID) {
        input.record = { '@type': 'ComposeRecord', '@value': props.record }
      }

      await $AutomationAPI.workflowExec({
        workflowID: btn.workflowID,
        stepID: btn.stepID || '0',
        input,
      })

    } else if (!btn.workflowID && !btn.automationID) {
      $toast?.toastInfo(t('block.automation.noScript'))
    }
  } catch (e) {
    console.error('Automation execution failed:', e)
    $toast?.toastErrorHandler(t('block.automation.executionFailed'))(e)
  } finally {
    processingIDs.value = processingIDs.value.filter(id => id !== index)
  }
}
</script>

