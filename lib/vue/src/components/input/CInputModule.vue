<template>
  <!-- Multi-select mode -->
  <MultiSelect
    v-if="multiple"
    :model-value="selectedModules"
    @update:model-value="onMultiSelect"
    :options="options"
    :option-label="getOptionLabel"
    :placeholder="placeholder"
    :disabled="disabled || !namespaceID"
    :loading="loading"
    class="w-full"
    display="chip"
    filter
    :filter-fields="['name', 'handle', 'moduleID']"
    fluid
    @show="onShow"
  >
    <template #option="{ option }">
      <span>{{ option.name }}</span>
    </template>
  </MultiSelect>

  <!-- Single-select mode -->
  <Select
    v-else
    :model-value="selectedModule"
    @update:model-value="onSelect"
    :options="options"
    :option-label="getOptionLabel"
    data-key="moduleID"
    :placeholder="placeholder"
    :disabled="disabled || !namespaceID"
    :loading="loading"
    class="w-full"
    filter
    :filter-fields="['name', 'handle', 'moduleID']"
    fluid
    showClear
    @show="onShow"
  >
    <template #option="{ option }">
      <span>{{ option.name }}</span>
    </template>
  </Select>
</template>

<script setup>
import { computed, inject, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useComposeResourceStore } from '../../stores/useComposeResourceStore'

defineOptions({ inheritAttrs: false })

const props = defineProps({
  modelValue: {
    type: [String, Number, Array],
    default: null,
  },
  namespaceID: {
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
  multiple: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue'])

// Store-first: prefer the app-wide shared module store (multi-namespace cache;
// modulesFor/loadFor are scoped to a namespace and don't disturb the active
// one). Fall back to the lib resource store for standalone/legacy usage.
const sharedStore = inject('$moduleStore', null)
const fallbackStore = useComposeResourceStore()

const fallbackOptions = ref([])
const selectedModule = ref(null)
const selectedModules = ref([])
const loading = ref(false)

let cancelCurrentRequest = null

// Options for THIS picker's namespace — from the shared store's per-namespace
// cache when available, otherwise the fallback search results.
const options = computed(() =>
  sharedStore ? sharedStore.modulesFor(props.namespaceID) : fallbackOptions.value,
)

function getOptionLabel(module) {
  if (!module) return ''
  return module.name || module.handle || module.moduleID
}

async function fetchModules() {
  if (!props.namespaceID) {
    if (!sharedStore) fallbackOptions.value = []
    return
  }

  loading.value = true
  try {
    if (sharedStore) {
      // Cache-guarded; returns immediately if this namespace is already loaded.
      await sharedStore.loadFor(props.namespaceID)
      return
    }
    if (cancelCurrentRequest) {
      cancelCurrentRequest()
      cancelCurrentRequest = null
    }
    const { response, cancel } = fallbackStore.searchModules(props.namespaceID, {
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
  if (props.namespaceID) {
    fetchModules()
  }
}

// --- Single-select ---

function onSelect(value) {
  selectedModule.value = value
  emit('update:modelValue', value?.moduleID || null)
}

async function resolveModule(moduleID) {
  if (sharedStore) {
    return sharedStore.findByID({ namespaceID: props.namespaceID, moduleID })
  }
  const mod = await fallbackStore.resolveModule(props.namespaceID, moduleID)
  if (mod && !fallbackOptions.value.find(m => m.moduleID === mod.moduleID)) {
    fallbackOptions.value = [...fallbackOptions.value, mod]
  }
  return mod
}

async function loadModuleById(moduleID) {
  if (!moduleID || !props.namespaceID) return

  const existing = options.value.find(m => m.moduleID === moduleID)
  if (existing) {
    selectedModule.value = existing
    return
  }

  loading.value = true
  try {
    const module = await resolveModule(moduleID)
    if (module) selectedModule.value = module
  } catch (_e) {
    // Module not found or API error
  } finally {
    loading.value = false
  }
}

// --- Multi-select ---

function onMultiSelect(modules) {
  selectedModules.value = modules
  emit('update:modelValue', modules.map(m => m.moduleID))
}

async function loadModulesById(moduleIDs) {
  if (!moduleIDs?.length || !props.namespaceID) return

  const resolved = []
  for (const id of moduleIDs) {
    let mod = options.value.find(m => m.moduleID === id || m.moduleID === String(id))
    if (!mod) {
      try {
        mod = await resolveModule(id)
      } catch (_e) {
        // skip
      }
    }
    if (mod) resolved.push(mod)
  }
  selectedModules.value = resolved
}

// --- Watchers ---

watch(
  () => props.namespaceID,
  (newVal, oldVal) => {
    if (oldVal && newVal !== oldVal) {
      selectedModule.value = null
      selectedModules.value = []
      if (!sharedStore) fallbackOptions.value = []
      emit('update:modelValue', props.multiple ? [] : null)
    }
    if (newVal) {
      fetchModules()
    }
  },
)

watch(
  () => props.modelValue,
  newVal => {
    if (props.multiple) {
      const ids = newVal || []
      if (ids.length && options.value.length) {
        syncMultiSelection(ids)
      } else if (!ids.length) {
        selectedModules.value = []
      }
    } else {
      if (newVal && (!selectedModule.value || selectedModule.value.moduleID !== newVal)) {
        loadModuleById(newVal)
      } else if (!newVal) {
        selectedModule.value = null
      }
    }
  },
  { immediate: true },
)

// When options load, sync multi-selection
watch(
  () => options.value.length,
  () => {
    if (props.multiple && props.modelValue?.length && options.value.length) {
      syncMultiSelection(props.modelValue)
    }
  },
)

function syncMultiSelection(ids) {
  selectedModules.value = options.value.filter(m =>
    ids.includes(m.moduleID) || ids.includes(String(m.moduleID)),
  )
  const missing = ids.filter(id =>
    !selectedModules.value.find(m => m.moduleID === id || m.moduleID === String(id)),
  )
  if (missing.length) {
    loadModulesById(missing)
  }
}

onMounted(() => {
  if (props.namespaceID) {
    fetchModules()
  }
  if (props.multiple && props.modelValue?.length && props.namespaceID) {
    loadModulesById(props.modelValue)
  } else if (!props.multiple && props.modelValue && props.namespaceID) {
    loadModuleById(props.modelValue)
  }
})

onBeforeUnmount(() => {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
})
</script>
