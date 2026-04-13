<template>
  <Select
    :model-value="modelValue"
    @update:model-value="onSelect"
    :options="options"
    option-label="label"
    option-value="automationID"
    :placeholder="placeholder"
    :disabled="disabled"
    :loading="loading"
    class="w-full"
    filter
    :filter-fields="['label', 'handle', 'automationID']"
    fluid
    showClear
    @show="onShow"
  />
</template>

<script setup>
import { inject, onBeforeUnmount, onMounted, ref, watch } from 'vue'

defineOptions({ inheritAttrs: false })

const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: null,
  },
  placeholder: {
    type: String,
    default: '',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue'])

const $AutomationAPI = inject('$AutomationAPI')

const options = ref([])
const loading = ref(false)

let cancelCurrentRequest = null

function getOptionLabel(automation) {
  if (!automation) return ''
  return automation.meta?.short || automation.handle || automation.automationID
}

async function fetchAutomations() {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
    cancelCurrentRequest = null
  }

  loading.value = true
  try {
    const { response, cancel } = $AutomationAPI.ngAutomationListCancellable({
      limit: 100,
      sort: 'handle ASC',
    })
    cancelCurrentRequest = cancel

    const result = await response()
    const automations = Array.isArray(result) ? result : result.set || []

    options.value = automations
      .map(a => ({ ...a, label: getOptionLabel(a) }))
      .sort((a, b) => (a.label || '').localeCompare(b.label || ''))

    if (props.modelValue && props.modelValue !== '0') {
      if (!options.value.some(a => a.automationID === props.modelValue)) {
        loadAutomationById(props.modelValue)
      }
    }
  } finally {
    loading.value = false
    cancelCurrentRequest = null
  }
}

function onShow() {
  if (options.value.length === 0) {
    fetchAutomations()
  }
}

function onSelect(automationID) {
  emit('update:modelValue', automationID || null)
}

async function loadAutomationById(automationID) {
  if (!automationID || automationID === '0' || !$AutomationAPI) return
  loading.value = true
  try {
    const automation = await $AutomationAPI.ngAutomationRead({ automationID })
    if (automation) {
      automation.label = getOptionLabel(automation)
      if (!options.value.find(a => a.automationID === automationID)) {
        options.value = [...options.value, automation]
      }
    }
  } catch {
    // TAQ not found or API error
  } finally {
    loading.value = false
  }
}

watch(
  () => props.modelValue,
  newVal => {
    if (newVal && newVal !== '0') {
      loadAutomationById(newVal)
    }
  },
  { immediate: true },
)

onMounted(() => {
  fetchAutomations()
  if (props.modelValue && props.modelValue !== '0') {
    loadAutomationById(props.modelValue)
  }
})

onBeforeUnmount(() => {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
})
</script>
