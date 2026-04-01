<template>
  <Select
    :model-value="modelValue"
    @update:model-value="onSelect"
    :options="options"
    option-label="label"
    option-value="workflowID"
    :placeholder="placeholder"
    :disabled="disabled"
    :loading="loading"
    class="w-full"
    filter
    :filter-fields="['label', 'handle', 'workflowID']"
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

function getOptionLabel(workflow) {
  if (!workflow) return ''
  return workflow.meta?.name || workflow.handle || workflow.workflowID
}

async function fetchWorkflows() {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
    cancelCurrentRequest = null
  }

  loading.value = true
  try {
    const { response, cancel } = $AutomationAPI.workflowListCancellable({
      limit: 100,
    })
    cancelCurrentRequest = cancel

    const result = await response()
    const workflows = Array.isArray(result) ? result : result.set || []
    
    options.value = workflows.map(w => ({ ...w, label: getOptionLabel(w) }))

    if (props.modelValue && props.modelValue !== '0') {
      if (!options.value.some(w => w.workflowID === props.modelValue)) {
        loadWorkflowById(props.modelValue)
      }
    }
  } finally {
    loading.value = false
    cancelCurrentRequest = null
  }
}

function onShow() {
  if (options.value.length === 0) {
    fetchWorkflows()
  }
}

function onSelect(workflowID) {
  emit('update:modelValue', workflowID || null)
}

async function loadWorkflowById(workflowID) {
  if (!workflowID || workflowID === '0' || !$AutomationAPI) return
  loading.value = true
  try {
    const workflow = await $AutomationAPI.workflowRead({ workflowID })
    if (workflow) {
      workflow.label = getOptionLabel(workflow)
      if (!options.value.find(w => w.workflowID === workflowID)) {
        options.value = [...options.value, workflow]
      }
    }
  } catch {
    // Workflow not found or API error
  } finally {
    loading.value = false
  }
}

watch(
  () => props.modelValue,
  newVal => {
    if (
      newVal &&
      newVal !== '0'
    ) {
      loadWorkflowById(newVal)
    }
  },
  { immediate: true },
)

onMounted(() => {
  fetchWorkflows()
  if (props.modelValue && props.modelValue !== '0') {
    loadWorkflowById(props.modelValue)
  }
})

onBeforeUnmount(() => {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
})
</script>
