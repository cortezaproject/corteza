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
import { computed, onMounted, ref, watch } from 'vue'
import { NoID } from '@planetcrust/human-js'
import { useModuleStore } from '../../stores/useModuleStore'

// An unset module reads as NoID once it has been through a block's options, and
// fetching '0' is a request the server rejects — a console error for a picker
// that simply has nothing selected yet.
const isUnset = id => !id || String(id) === NoID

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

const store = useModuleStore()
const selectedModule = ref(null)
const selectedModules = ref([])
const loading = ref(false)

const options = computed(() => store.modulesFor(props.namespaceID))

function getOptionLabel(module) {
  if (!module) return ''
  return module.name || module.handle || module.moduleID
}

async function fetchModules() {
  if (!props.namespaceID) return

  loading.value = true
  try {
    await store.loadFor(props.namespaceID)
  } catch {
    // ignore
  } finally {
    loading.value = false
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

async function loadModuleById(moduleID) {
  if (isUnset(moduleID) || !props.namespaceID) return

  const existing = options.value.find(m => m.moduleID === moduleID)
  if (existing) {
    selectedModule.value = existing
    return
  }

  loading.value = true
  try {
    const module = await store.findByID({ namespaceID: props.namespaceID, moduleID })
    if (module) {
      selectedModule.value = store.set.find(m => m.moduleID === moduleID) || module
    }
  } catch {
    // module not found
  } finally {
    loading.value = false
  }
}

// --- Multi-select ---

function onMultiSelect(modules) {
  selectedModules.value = modules
  emit(
    'update:modelValue',
    modules.map(m => m.moduleID),
  )
}

async function loadModulesById(moduleIDs) {
  if (!moduleIDs?.length || !props.namespaceID) return

  const resolved = []
  for (const id of moduleIDs) {
    if (isUnset(id)) continue
    let mod = options.value.find(m => m.moduleID === id || m.moduleID === String(id))
    if (!mod) {
      try {
        mod = await store.findByID({ namespaceID: props.namespaceID, moduleID: id })
      } catch {
        // skip
      }
    }
    if (mod) resolved.push(mod)
  }
  selectedModules.value = resolved
}

function syncMultiSelection(ids) {
  selectedModules.value = options.value.filter(
    m => ids.includes(m.moduleID) || ids.includes(String(m.moduleID)),
  )
  const missing = ids.filter(
    id => !selectedModules.value.find(m => m.moduleID === id || m.moduleID === String(id)),
  )
  if (missing.length) {
    loadModulesById(missing)
  }
}

// --- Watchers ---

watch(
  () => props.namespaceID,
  (newVal, oldVal) => {
    if (oldVal && newVal !== oldVal) {
      selectedModule.value = null
      selectedModules.value = []
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

watch(
  () => options.value.length,
  () => {
    if (props.multiple && props.modelValue?.length && options.value.length) {
      syncMultiSelection(props.modelValue)
    }
  },
)

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
</script>
