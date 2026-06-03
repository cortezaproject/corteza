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
import { computed, inject, onBeforeUnmount, onMounted, ref, watch } from 'vue'
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

// Store-first: prefer the app-wide shared namespace store (preloaded at boot,
// no per-open refetch). Fall back to the lib resource store when the shared
// store isn't provided (e.g. standalone/legacy usage).
const sharedStore = inject('$namespaceStore', null)
const fallbackStore = useComposeResourceStore()

const fallbackOptions = ref([])
const selectedNamespace = ref(null)
const loading = ref(false)

// Options come straight from the shared store's reactive (preloaded) set when
// available; otherwise from the fallback store's search results.
const options = computed(() => (sharedStore ? sharedStore.set : fallbackOptions.value))

// Store cancel function for the fallback request
let cancelCurrentRequest = null

function getOptionLabel(namespace) {
  if (!namespace) return ''
  return namespace.name || namespace.slug || namespace.namespaceID
}

async function ensureLoaded() {
  loading.value = true
  try {
    if (sharedStore) {
      // Cache-guarded — returns immediately if already loaded.
      await sharedStore.load()
      return
    }
    if (cancelCurrentRequest) {
      cancelCurrentRequest()
      cancelCurrentRequest = null
    }
    const { response, cancel } = fallbackStore.searchNamespaces({
      query: '',
      limit: 100,
      sort: 'name ASC',
    })
    cancelCurrentRequest = cancel
    const result = await response()
    fallbackOptions.value = result.set || []
  } catch (e) {
    if (e?.message !== 'canceled' && !sharedStore) {
      fallbackOptions.value = []
    }
  } finally {
    loading.value = false
    cancelCurrentRequest = null
  }
}

function onShow() {
  ensureLoaded()
}

function onSelect(value) {
  selectedNamespace.value = value
  emit('update:modelValue', value?.namespaceID || null)
}

async function loadNamespaceById(namespaceID) {
  if (!namespaceID) return

  const existing = options.value.find(ns => ns.namespaceID === namespaceID)
  if (existing) {
    selectedNamespace.value = existing
    return
  }

  // Resolve through the store (cache-first).
  loading.value = true
  try {
    const namespace = sharedStore
      ? await sharedStore.findByID({ namespaceID })
      : await fallbackStore.resolveNamespace(namespaceID)
    if (namespace) {
      selectedNamespace.value = namespace
      // For the fallback store, keep its local option list in sync.
      if (!sharedStore && !fallbackOptions.value.find(ns => ns.namespaceID === namespaceID)) {
        fallbackOptions.value = [...fallbackOptions.value, namespace]
      }
    }
  } catch (_e) {
    // Namespace not found or API error
  } finally {
    loading.value = false
  }
}

watch(
  () => props.modelValue,
  newVal => {
    if (newVal && (!selectedNamespace.value || selectedNamespace.value.namespaceID !== newVal)) {
      loadNamespaceById(newVal)
    } else if (!newVal) {
      selectedNamespace.value = null
    }
  },
  { immediate: true },
)

onMounted(() => {
  ensureLoaded()

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
