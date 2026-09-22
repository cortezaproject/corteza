<template>
  <div :class="containerClass">
    <Button
      v-for="(btn, i) in buttons"
      :key="i"
      :label="evaluatedLabel(btn) || '-'"
      :severity="mapVariant(btn.variant)"
      :loading="processingIDs.includes(i)"
      :disabled="processingIDs.includes(i)"
      :size="size"
      :class="buttonClass"
      @click.prevent="handleButton(btn, i)"
    />
  </div>
</template>

<script setup>
import { ref, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import { compose } from '@planetcrust/human-js'
import { evaluatePrefilter } from '../../../lib/record-filter'
import { scriptConstraintMatcher } from '../../../lib/script-events'

const { t } = useI18n()

const props = defineProps({
  buttons: { type: Array, default: () => [] },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  module: { type: Object, default: () => null },
  record: { type: Object, default: undefined },
  records: { type: Array, default: () => [] },
  filter: { type: String, default: '' },
  containerClass: { type: String, default: 'flex gap-2' },
  buttonClass: { type: String, default: '' },
  size: { type: String, default: null },
})

const emit = defineEmits(['refresh'])

const $toast = inject('$toast', null)
const $Auth = inject('$Auth', {})
const $AutomationAPI = inject('$AutomationAPI', null)
const $ScriptBus = inject('$ScriptBus', null)

const processingIDs = ref([])

const variantSeverityMap = {
  primary: undefined,
  secondary: 'secondary',
  light: 'secondary',
  dark: 'contrast',
  success: 'success',
  danger: 'danger',
  warning: 'warn',
  info: 'info',
}
const mapVariant = key => variantSeverityMap[key]

function evaluatedLabel(btn) {
  try {
    const record = props.record
    const user = $Auth?.user || {}
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

// Every entry is a typed envelope: the exec endpoint decodes `input` into
// expr.Vars, which reads only {"@type":…,"@value":…} and rejects the whole
// request over a bare value.
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
    input.selected = { '@type': 'Array', '@value': props.records }
  }
  if (props.filter) {
    input.filter = { '@type': 'String', '@value': props.filter }
  }
  return input
}

// The event a Corredor script is triggered with — shaped by the resource its
// trigger is bound to, and carrying the whole page context as arguments.
function buildScriptEvent(btn) {
  const args = {
    namespace: props.namespace,
    selected: props.records,
    filter: props.filter,
  }

  if (props.module) args.module = props.module
  if (props.page?.pageID) args.page = props.page

  const match = scriptConstraintMatcher({ namespace: props.namespace, module: props.module })

  switch (btn.resourceType) {
    case 'compose:record':
      if (!props.record || !props.module) {
        $toast?.toastWarning?.(t('block.automation.noRecord'))
        return null
      }
      return compose.RecordEvent(props.record, { match, args })
    case 'compose:module':
      return compose.ModuleEvent(props.module, { match, args })
    case 'compose:namespace':
      return compose.NamespaceEvent(props.namespace, { match, args })
    case 'compose:page':
      return compose.PageEvent(props.page, { match, args })
    default:
      return compose.ComposeEvent({ match, args })
  }
}

// The bus knows which scripts are client and which are server ones, and runs
// each where it belongs.
async function dispatchScript(btn) {
  if (!$ScriptBus) return

  const ev = buildScriptEvent(btn)
  if (!ev) return

  await $ScriptBus.Dispatch(ev, btn.script)
}

async function handleButton(btn, index) {
  if (!$AutomationAPI && !btn.script) return

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
    } else if (btn.script) {
      await dispatchScript(btn)
    } else {
      $toast?.toastInfo?.(t('block.automation.noScript'))
      return
    }

    emit('refresh')
  } catch (e) {
    console.error('Automation execution failed:', e)
    $toast?.toastErrorHandler?.(
      t(btn.script ? 'notification.automation.scriptFailed' : 'block.automation.executionFailed'),
    )(e)
  } finally {
    processingIDs.value = processingIDs.value.filter(id => id !== index)
  }
}
</script>
