<template>
  <CInputSelect
    :model-value="selectedRecord"
    @update:model-value="onSelect"
    :options="options"
    :option-label="getOptionLabel"
    :placeholder="effectivePlaceholder"
    :disabled="disabled || !namespaceId || !moduleId"
    :loading="loading"
    :complete-on-focus="true"
    @search="onSearch"
  >
    <template #option="{ option }">
      <div class="flex items-center gap-2">
        <span>{{ getOptionLabel(option) }}</span>
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
  namespaceId: {
    type: [String, Number],
    default: null,
  },
  moduleId: {
    type: [String, Number],
    default: null,
  },
  // Field to use as display label (e.g., 'name', 'title')
  labelField: {
    type: String,
    default: '',
  },
  placeholder: {
    type: String,
    default: 'Select a record',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue'])

const $ComposeAPI = inject('$ComposeAPI')

const options = ref([])
const selectedRecord = ref(null)
const loading = ref(false)

// Store cancel function for current request
let cancelCurrentRequest = null

const effectivePlaceholder = computed(() => {
  if (!props.namespaceId) {
    return 'Select a namespace first'
  }
  if (!props.moduleId) {
    return 'Select a module first'
  }
  return props.placeholder
})

function getOptionLabel(record) {
  if (!record) return ''

  // Try labelField if provided
  if (props.labelField && record.values) {
    const value = record.values.find(v => v.name === props.labelField)
    if (value?.value) return value.value
  }

  // Fallback to first value with content, or recordID
  if (record.values?.length) {
    const firstValue = record.values.find(v => v.value)
    if (firstValue?.value) return firstValue.value
  }

  return `Record ${record.recordID}`
}

async function fetchRecords(query = '') {
  if (!props.namespaceId || !props.moduleId || !$ComposeAPI) {
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
    const { response, cancel } = $ComposeAPI.recordListCancellable({
      namespaceID: props.namespaceId,
      moduleID: props.moduleId,
      query,
      limit: 50,
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

const debouncedFetch = debounce(query => {
  fetchRecords(query)
}, 200)

function onSearch(query) {
  debouncedFetch(query)
}

function onSelect(value) {
  selectedRecord.value = value
  emit('update:modelValue', value?.recordID || null)
}

async function loadRecordById(recordID) {
  if (!recordID || !props.namespaceId || !props.moduleId || !$ComposeAPI) return

  // First check if already in options
  const existing = options.value.find(r => r.recordID === recordID)
  if (existing) {
    selectedRecord.value = existing
    return
  }

  // Otherwise fetch it
  loading.value = true
  try {
    const record = await $ComposeAPI.recordRead({
      namespaceID: props.namespaceId,
      moduleID: props.moduleId,
      recordID,
    })
    selectedRecord.value = record
    // Add to options if not present
    if (!options.value.find(r => r.recordID === recordID)) {
      options.value = [...options.value, record]
    }
  } catch {
    // Record not found or API error
  } finally {
    loading.value = false
  }
}

// Watch for namespace/module changes - clear selection and reload records
watch(
  () => [props.namespaceId, props.moduleId],
  ([newNs, newMod], [oldNs, oldMod]) => {
    if ((oldNs && newNs !== oldNs) || (oldMod && newMod !== oldMod)) {
      selectedRecord.value = null
      emit('update:modelValue', null)
    }
    if (newNs && newMod) {
      fetchRecords()
    } else {
      options.value = []
    }
  },
)

watch(
  () => props.modelValue,
  newVal => {
    if (newVal && (!selectedRecord.value || selectedRecord.value.recordID !== newVal)) {
      loadRecordById(newVal)
    } else if (!newVal) {
      selectedRecord.value = null
    }
  },
  { immediate: true },
)

onMounted(() => {
  if (props.namespaceId && props.moduleId) {
    fetchRecords()
  }
  if (props.modelValue && props.namespaceId && props.moduleId) {
    loadRecordById(props.modelValue)
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
