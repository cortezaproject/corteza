<template>
  <Select
    :model-value="selectedProvider"
    @update:model-value="onSelect"
    :options="options"
    :option-label="getOptionLabel"
    :placeholder="placeholder"
    :disabled="disabled"
    :loading="loading"
    :size="size"
    class="w-full"
    filter
    fluid
    showClear
    @show="onShow"
  >
    <template #option="{ option }">
      {{ getOptionLabel(option) }}
    </template>
    <!-- Let consumers add actions (e.g. "Add provider") inside the dropdown. -->
    <template v-if="$slots.footer" #footer>
      <slot name="footer" />
    </template>
  </Select>
</template>

<script setup>
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
  size: {
    type: String,
    default: null,
  },
})

const emit = defineEmits(['update:modelValue'])

const $SystemAPI = inject('$SystemAPI')

const options = ref([])
const selectedProvider = ref(null)
const loading = ref(false)

let cancelCurrentRequest = null

function getOptionLabel(provider) {
  if (!provider) return ''
  return provider.meta?.short || provider.handle || provider.provider || provider.llmProviderID
}

async function fetchProviders() {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
    cancelCurrentRequest = null
  }

  loading.value = true
  try {
    const { response, cancel } = $SystemAPI.llmProviderListCancellable({
      status: 'active',
      sort: 'handle ASC',
    })
    cancelCurrentRequest = cancel

    const result = await response()
    const all = Array.isArray(result) ? result : result.set || []
    options.value = all
      .filter(p => !p.deletedAt && p.status === 'active')
      .sort((a, b) => (getOptionLabel(a) || '').localeCompare(getOptionLabel(b) || ''))
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
  if (options.value.length === 0) {
    fetchProviders()
  }
}

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
    if (!options.value.find(p => p.llmProviderID === llmProviderID)) {
      options.value = [...options.value, provider]
    }
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
  fetchProviders()
  if (props.modelValue && props.modelValue !== '0') {
    loadProviderById(props.modelValue)
  }
})

onBeforeUnmount(() => {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
})
</script>
