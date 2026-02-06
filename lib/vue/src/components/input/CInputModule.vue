<template>
  <CInputSelect
    :model-value="selectedModule"
    @update:model-value="onSelect"
    :options="options"
    :option-label="getOptionLabel"
    :placeholder="effectivePlaceholder"
    :disabled="disabled || !namespaceID"
    :loading="loading"
    @search="onSearch"
  >
    <template #option="{ option }">
      <div class="flex items-center gap-2">
        <span>{{ option.name }}</span>
        <span v-if="option.handle" class="text-muted-color text-xs">{{ option.handle }}</span>
      </div>
    </template>
  </CInputSelect>
</template>

<script setup>
import { debounce } from 'lodash-es'
import { computed, inject, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import CInputSelect from './CInputSelect.vue'

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

const $ComposeAPI = inject('$ComposeAPI')

const options = ref([])
const selectedModule = ref(null)
const loading = ref(false)

// Store cancel function for current request
let cancelCurrentRequest = null

const effectivePlaceholder = computed(() => {
  if (!props.namespaceID) {
    return 'Select a namespace first'
  }
  return props.placeholder
})

function getOptionLabel(module) {
  if (!module) return ''
  return module.name || module.handle || module.moduleID
}

async function fetchModules(query = '') {
  if (!props.namespaceID || !$ComposeAPI) {
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
    const { response, cancel } = $ComposeAPI.moduleListCancellable({
      namespaceID: props.namespaceID,
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
  fetchModules(query)
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
  selectedModule.value = value
  emit('update:modelValue', value?.moduleID || null)
}

async function loadModuleById(moduleID) {
  if (!moduleID || !props.namespaceID || !$ComposeAPI) return

  // First check if already in options
  const existing = options.value.find(m => m.moduleID === moduleID)
  if (existing) {
    selectedModule.value = existing
    return
  }

  // Otherwise fetch it
  loading.value = true
  try {
    const module = await $ComposeAPI.moduleRead({
      namespaceID: props.namespaceID,
      moduleID,
    })
    selectedModule.value = module
    // Add to options if not present
    if (!options.value.find(m => m.moduleID === moduleID)) {
      options.value = [...options.value, module]
    }
  } catch (_e) {
    // Module not found or API error
  } finally {
    loading.value = false
  }
}

// Watch for namespace changes - clear selection and reload modules
watch(() => props.namespaceID, (newVal, oldVal) => {
  if (oldVal && newVal !== oldVal) {
    selectedModule.value = null
    emit('update:modelValue', null)
  }
  if (newVal) {
    fetchModules()
  } else {
    options.value = []
  }
})

watch(() => props.modelValue, (newVal) => {
  if (newVal && (!selectedModule.value || selectedModule.value.moduleID !== newVal)) {
    loadModuleById(newVal)
  } else if (!newVal) {
    selectedModule.value = null
  }
}, { immediate: true })

onMounted(() => {
  if (props.namespaceID) {
    fetchModules()
  }
  if (props.modelValue && props.namespaceID) {
    loadModuleById(props.modelValue)
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
