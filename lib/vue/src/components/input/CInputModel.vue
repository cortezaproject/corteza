<template>
  <Select
    :model-value="selectedModel"
    @update:model-value="onSelect"
    :options="options"
    :option-label="getOptionLabel"
    :placeholder="placeholder"
    :disabled="disabled || !llmProviderID || llmProviderID === '0'"
    :loading="loading"
    class="w-full"
    filter
    fluid
    showClear
    @show="onShow"
  >
    <template #option="{ option }">
      {{ getOptionLabel(option) }}
    </template>
  </Select>
</template>

<script setup>
import { inject, onBeforeUnmount, onMounted, ref, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: null,
  },
  llmProviderID: {
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

const $SystemAPI = inject('$SystemAPI')

const options = ref([])
const selectedModel = ref(null)
const loading = ref(false)

let cancelCurrentRequest = null

function getOptionLabel(model) {
  if (!model) return ''
  if (typeof model === 'string') return model
  return model.name || model.model || model.id || ''
}

async function fetchModels() {
  if (!props.llmProviderID || props.llmProviderID === '0' || !$SystemAPI) {
    options.value = []
    return
  }

  if (cancelCurrentRequest) {
    cancelCurrentRequest()
    cancelCurrentRequest = null
  }

  loading.value = true
  try {
    const { response, cancel } = $SystemAPI.llmProviderModelsCancellable({
      llmProviderID: props.llmProviderID,
    })
    cancelCurrentRequest = cancel

    const result = await response()
    options.value = Array.isArray(result) ? result : result.models || result.set || []
  } catch (e) {
    if (e?.message !== 'canceled') {
      options.value = []
    }
  } finally {
    loading.value = false
    cancelCurrentRequest = null
  }
}

function onShow() {
  if (options.value.length === 0 && props.llmProviderID && props.llmProviderID !== '0') {
    fetchModels()
  }
}

function onSelect(value) {
  selectedModel.value = value

  if (!value) {
    emit('update:modelValue', null)
  } else if (typeof value === 'string') {
    emit('update:modelValue', value)
  } else {
    emit('update:modelValue', value.model || value.name || value.id || null)
  }
}

// Restore selection from modelValue string
function restoreSelection() {
  if (!props.modelValue) {
    selectedModel.value = null
    return
  }

  const match = options.value.find(m => {
    if (typeof m === 'string') return m === props.modelValue
    return (m.model || m.name || m.id) === props.modelValue
  })

  selectedModel.value = match || null
}

// When provider changes, re-fetch models and clear selection
watch(
  () => props.llmProviderID,
  (newVal, oldVal) => {
    if (newVal !== oldVal) {
      selectedModel.value = null
      emit('update:modelValue', null)
      fetchModels()
    }
  },
)

watch(
  () => props.modelValue,
  () => {
    restoreSelection()
  },
)

watch(options, () => {
  restoreSelection()
})

onMounted(() => {
  if (props.llmProviderID && props.llmProviderID !== '0') {
    fetchModels()
  }
})

onBeforeUnmount(() => {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
})
</script>
