<template>
  <Select
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
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useComposeResourceStore } from '../../stores/useComposeResourceStore'

defineOptions({ inheritAttrs: false })

const props = defineProps({
  modelValue: {
    type: [String, Number],
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
})

const emit = defineEmits(['update:modelValue'])

const store = useComposeResourceStore()

const options = ref([])
const selectedModule = ref(null)
const loading = ref(false)

// Store cancel function for current request
let cancelCurrentRequest = null

function getOptionLabel(module) {
  if (!module) return ''
  return module.name || module.handle || module.moduleID
}

async function fetchModules() {
  if (!props.namespaceID) {
    options.value = []
    return
  }

  // Cancel previous request if pending
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
    cancelCurrentRequest = null
  }

  loading.value = true
  try {
    const { response, cancel } = store.searchModules(props.namespaceID, {
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
  if (options.value.length === 0 && props.namespaceID) {
    fetchModules()
  }
}

function onSelect(value) {
  selectedModule.value = value
  emit('update:modelValue', value?.moduleID || null)
}

async function loadModuleById(moduleID) {
  if (!moduleID || !props.namespaceID) return

  // First check if already in options
  const existing = options.value.find(m => m.moduleID === moduleID)
  if (existing) {
    selectedModule.value = existing
    return
  }

  // Resolve through store (cache-first)
  loading.value = true
  try {
    const module = await store.resolveModule(props.namespaceID, moduleID)
    if (module) {
      selectedModule.value = module
      if (!options.value.find(m => m.moduleID === moduleID)) {
        options.value = [...options.value, module]
      }
    }
  } catch (_e) {
    // Module not found or API error
  } finally {
    loading.value = false
  }
}

// Watch for namespace changes - clear selection and reload
watch(
  () => props.namespaceID,
  (newVal, oldVal) => {
    if (oldVal && newVal !== oldVal) {
      selectedModule.value = null
      options.value = []
      emit('update:modelValue', null)
    }
    if (newVal) {
      fetchModules()
    }
  },
)

watch(
  () => props.modelValue,
  newVal => {
    if (newVal && (!selectedModule.value || selectedModule.value.moduleID !== newVal)) {
      loadModuleById(newVal)
    } else if (!newVal) {
      selectedModule.value = null
    }
  },
  { immediate: true },
)

onMounted(() => {
  if (props.namespaceID) {
    fetchModules()
  }
  if (props.modelValue && props.namespaceID) {
    loadModuleById(props.modelValue)
  }
})

onBeforeUnmount(() => {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
})
</script>
