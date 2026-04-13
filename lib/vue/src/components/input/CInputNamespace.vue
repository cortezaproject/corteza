<template>
  <Select
    :model-value="selectedNamespace"
    @update:model-value="onSelect"
    :options="options"
    :option-label="getOptionLabel"
    data-key="namespaceID"
    :placeholder="placeholder"
    :disabled="disabled"
    :loading="loading"
    class="w-full"
    filter
    :filter-fields="['name', 'slug', 'namespaceID']"
    fluid
    showClear
    @show="onShow"
  >
    <template #option="{ option }">
      {{ option.name }}
    </template>
  </Select>
</template>

<script setup>
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useComposeResourceStore } from '../../stores/useComposeResourceStore'

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

const store = useComposeResourceStore()

const options = ref([])
const selectedNamespace = ref(null)
const loading = ref(false)

// Store cancel function for current request
let cancelCurrentRequest = null

function getOptionLabel(namespace) {
  if (!namespace) return ''
  return namespace.name || namespace.slug || namespace.namespaceID
}

async function fetchNamespaces() {
  // Cancel previous request if pending
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
    cancelCurrentRequest = null
  }

  loading.value = true
  try {
    const { response, cancel } = store.searchNamespaces({
      query: '',
      limit: 100,
      sort: 'name ASC',
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

function onShow() {
  fetchNamespaces()
}

function onSelect(value) {
  selectedNamespace.value = value
  emit('update:modelValue', value?.namespaceID || null)
}

async function loadNamespaceById(namespaceID) {
  if (!namespaceID) return

  // First check if already in options
  const existing = options.value.find(ns => ns.namespaceID === namespaceID)
  if (existing) {
    selectedNamespace.value = existing
    return
  }

  // Resolve through store (cache-first)
  loading.value = true
  try {
    const namespace = await store.resolveNamespace(namespaceID)
    if (namespace) {
      selectedNamespace.value = namespace
      // Add to options if not present
      if (!options.value.find(ns => ns.namespaceID === namespaceID)) {
        options.value = [...options.value, namespace]
      }
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
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
})
</script>
