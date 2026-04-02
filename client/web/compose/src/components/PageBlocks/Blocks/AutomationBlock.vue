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
    if (btn.workflowID && $AutomationAPI) {
      // Execute workflow
      const input = []

      // Pass context as input parameters
      if (props.namespace?.namespaceID) {
        input.push({ name: 'namespace', value: JSON.stringify({ namespaceID: props.namespace.namespaceID }) })
      }
      if (props.page?.pageID) {
        input.push({ name: 'page', value: JSON.stringify({ pageID: props.page.pageID }) })
      }
      if (props.record?.recordID) {
        input.push({ name: 'record', value: JSON.stringify(props.record) })
      }

      await $AutomationAPI.workflowExec({
        workflowID: btn.workflowID,
        stepID: btn.stepID || '0',
        input,
      })

      $toast?.toastSuccess(t('block.automation.executionSuccess'))
    } else if (!btn.workflowID) {
      $toast?.toastInfo(t('block.automation.noScript'))
    } else {
      $toast?.toastWarning(t('block.automation.noScript'))
    }
  } catch (e) {
    console.error('Automation execution failed:', e)
    $toast?.toastDanger(t('block.automation.executionFailed'))
  } finally {
    processingIDs.value = processingIDs.value.filter(id => id !== index)
  }
}
</script>

