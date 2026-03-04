<template>
  <AutoComplete
    :model-value="selectedModel"
    @update:model-value="onSelect"
    :suggestions="filteredModels"
    :option-label="getOptionLabel"
    :placeholder="placeholder"
    :disabled="disabled || !llmProviderID || llmProviderID === '0'"
    :loading="loading"
    @complete="onComplete"
    class="w-full"
    dropdown
  >
    <template #option="{ option }">
      {{ getOptionLabel(option) }}
    </template>
  </AutoComplete>
</template>

<script setup>
import AutoComplete from 'primevue/autocomplete'
import { inject, onMounted, ref, watch } from 'vue'

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

const allModels = ref([])
const filteredModels = ref([])
const selectedModel = ref(null)
const loading = ref(false)

function getOptionLabel(model) {
  if (!model) return ''
  if (typeof model === 'string') return model
  return model.name || model.model || model.id || ''
}

async function fetchModels() {
  if (!props.llmProviderID || props.llmProviderID === '0' || !$SystemAPI) {
    allModels.value = []
    filteredModels.value = []
    return
  }

  loading.value = true
  try {
    const response = await $SystemAPI.llmProviderModels({
      llmProviderID: props.llmProviderID,
    })
    const models = Array.isArray(response) ? response : response.models || response.set || []
    allModels.value = models
    filteredModels.value = [...models]
  } catch {
    allModels.value = []
    filteredModels.value = []
  } finally {
    loading.value = false
  }
}

function onComplete(event) {
  const query = (event.query || '').toLowerCase()
  if (!query) {
    filteredModels.value = [...allModels.value]
  } else {
    filteredModels.value = allModels.value.filter(m => {
      const label = getOptionLabel(m).toLowerCase()
      return label.includes(query)
    })
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

  // Find matching model object in allModels
  const match = allModels.value.find(m => {
    if (typeof m === 'string') return m === props.modelValue
    return (m.model || m.name || m.id) === props.modelValue
  })

  selectedModel.value = match || props.modelValue
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

watch(allModels, () => {
  restoreSelection()
})

onMounted(() => {
  if (props.llmProviderID && props.llmProviderID !== '0') {
    fetchModels()
  }
})
</script>
