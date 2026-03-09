<template>
  <AutoComplete
    :model-value="selectedProvider"
    @update:model-value="onSelect"
    :suggestions="suggestions"
    :option-label="getOptionLabel"
    :placeholder="placeholder"
    :disabled="disabled"
    :loading="loading"
    @complete="search"
    class="w-full"
    dropdown
    showClear
  >
    <template #option="{ option }">
      {{ getOptionLabel(option) }}
    </template>
  </AutoComplete>
</template>

<script setup>
import AutoComplete from 'primevue/autocomplete'
import { debounce } from 'lodash-es'
import { inject, onBeforeUnmount, onMounted, ref, watch } from 'vue'

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

const $SystemAPI = inject('$SystemAPI')

const suggestions = ref([])
const selectedProvider = ref(null)
const loading = ref(false)

// Store cancel function for current request
let cancelCurrentRequest = null

function getOptionLabel(provider) {
  if (!provider) return ''
  return provider.handle || provider.meta?.short || provider.provider || provider.llmProviderID
}

const search = debounce(async event => {
  // Cancel previous request if pending
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
    cancelCurrentRequest = null
  }

  loading.value = true
  try {
    const { response, cancel } = $SystemAPI.llmProviderListCancellable({
      provider: event.query || undefined,
      status: 'active',
    })
    cancelCurrentRequest = cancel

    const result = await response()
    suggestions.value = Array.isArray(result) ? result : result.set || []
  } catch (e) {
    // Ignore cancelled requests
    if (e?.message !== 'canceled') {
      suggestions.value = []
    }
  } finally {
    loading.value = false
    cancelCurrentRequest = null
  }
}, 300)

function onSelect(value) {
  selectedProvider.value = value
  emit('update:modelValue', value?.llmProviderID || null)
}

async function loadProviderById(llmProviderID) {
  if (!llmProviderID || llmProviderID === '0' || !$SystemAPI) return
  loading.value = true
  try {
    const provider = await $SystemAPI.llmProviderRead({ llmProviderID })
    selectedProvider.value = provider
  } catch {
    // Provider not found or API error
  } finally {
    loading.value = false
  }
}

watch(
  () => props.modelValue,
  newVal => {
    if (
      newVal &&
      newVal !== '0' &&
      (!selectedProvider.value || selectedProvider.value.llmProviderID !== newVal)
    ) {
      loadProviderById(newVal)
    } else if (!newVal || newVal === '0') {
      selectedProvider.value = null
    }
  },
  { immediate: true },
)

onMounted(() => {
  if (props.modelValue && props.modelValue !== '0') {
    loadProviderById(props.modelValue)
  }
})

onBeforeUnmount(() => {
  // Cancel any pending request
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
  search.cancel()
})
</script>
