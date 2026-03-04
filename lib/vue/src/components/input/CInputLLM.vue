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
  >
    <template #option="{ option }">
      {{ getOptionLabel(option) }}
    </template>
  </AutoComplete>
</template>

<script setup>
import AutoComplete from 'primevue/autocomplete'
import { debounce } from 'lodash-es'
import { inject, onMounted, ref, watch } from 'vue'

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

function getOptionLabel(provider) {
  if (!provider) return ''
  return provider.handle || provider.meta?.short || provider.provider || provider.llmProviderID
}

const search = debounce(async event => {
  loading.value = true
  try {
    const response = await $SystemAPI.llmProviderList({
      provider: event.query || undefined,
      status: 'active',
    })
    suggestions.value = Array.isArray(response) ? response : response.set || []
  } catch {
    suggestions.value = []
  } finally {
    loading.value = false
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
</script>
