<template>
  <Select
    :model-value="modelValue"
    @update:model-value="onSelect"
    :options="options"
    option-label="label"
    option-value="id"
    :placeholder="placeholder"
    :disabled="disabled || !llmProviderID || llmProviderID === '0'"
    :form-control="{ novalidate: true }"
    :loading="loading"
    :size="size"
    class="w-full"
    filter
    :filter-fields="['label', 'id']"
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
  size: {
    type: String,
    default: null,
  },
})

const emit = defineEmits(['update:modelValue'])

const $SystemAPI = inject('$SystemAPI')

const options = ref([])
const loading = ref(false)
const selectedModel = ref(null)

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
    const rawModels = Array.isArray(result) ? result : result.models || result.set || []

    options.value = rawModels.map(m => {
      const isStr = typeof m === 'string'
      const id = isStr ? m : m.model || m.name || m.id
      const label = getOptionLabel(m)
      return isStr ? { id, label } : { ...m, id, label }
    })
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

function onSelect(modelID) {
  emit('update:modelValue', modelID || null)
}

// Ensure the specific string/model value gets emitted if available
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
