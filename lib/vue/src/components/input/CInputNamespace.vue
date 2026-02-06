<template>
  <CInputSelect
    :model-value="selectedNamespace"
    @update:model-value="onSelect"
    :options="options"
    :option-label="getOptionLabel"
    :placeholder="placeholder"
    :disabled="disabled"
    :loading="loading"
    @search="onSearch"
  >
    <template #option="{ option }">
      {{ option.name }}
    </template>
  </CInputSelect>
</template>

<script setup>
import { debounce } from 'lodash-es'
import { inject, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import CInputSelect from './CInputSelect.vue'

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

const $ComposeAPI = inject('$ComposeAPI')

const options = ref([])
const selectedNamespace = ref(null)
const loading = ref(false)

// Store cancel function for current request
let cancelCurrentRequest = null

function getOptionLabel(namespace) {
  if (!namespace) return ''
  return namespace.name || namespace.slug || namespace.namespaceID
}

async function fetchNamespaces(query = '') {
  if (!$ComposeAPI) return

  // Cancel previous request if pending
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
    cancelCurrentRequest = null
  }

  loading.value = true
  try {
    const { response, cancel } = $ComposeAPI.namespaceListCancellable({
      query,
      limit: 100,
    })
    cancelCurrentRequest = cancel

    const result = await response()
    options.value = result.set || []
  } catch (e) {
    // Ignore cancelled requests
    if (e?.message !== 'canceled') {
      options.value = []
    }
  } finally {
    loading.value = false
    cancelCurrentRequest = null
  }
}

const debouncedFetch = debounce((query) => {
  fetchNamespaces(query)
}, 200)

function onSearch(query) {
  // If empty query (dropdown click) and we already have options, don't refetch
  // This uses the preloaded data instead
  if (!query && options.value.length > 0) {
    return
  }
  debouncedFetch(query)
}

function onSelect(value) {
  selectedNamespace.value = value
  emit('update:modelValue', value?.namespaceID || null)
}

async function loadNamespaceById(namespaceID) {
  if (!namespaceID || !$ComposeAPI) return

  // First check if already in options
  const existing = options.value.find(ns => ns.namespaceID === namespaceID)
  if (existing) {
    selectedNamespace.value = existing
    return
  }

  // Otherwise fetch it
  loading.value = true
  try {
    const namespace = await $ComposeAPI.namespaceRead({ namespaceID })
    selectedNamespace.value = namespace
    // Add to options if not present
    if (!options.value.find(ns => ns.namespaceID === namespaceID)) {
      options.value = [...options.value, namespace]
    }
  } catch (_e) {
    // Namespace not found or API error
  } finally {
    loading.value = false
  }
}

watch(() => props.modelValue, (newVal) => {
  if (newVal && (!selectedNamespace.value || selectedNamespace.value.namespaceID !== newVal)) {
    loadNamespaceById(newVal)
  } else if (!newVal) {
    selectedNamespace.value = null
  }
}, { immediate: true })

onMounted(() => {
  fetchNamespaces()
  if (props.modelValue) {
    loadNamespaceById(props.modelValue)
  }
})

onBeforeUnmount(() => {
  // Cancel any pending request
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
  debouncedFetch.cancel()
})
</script>
