<template>
  <Select
    :model-value="selectedRecord"
    @update:model-value="onSelect"
    :options="options"
    :option-label="getOptionLabel"
    :placeholder="effectivePlaceholder"
    :disabled="disabled || !namespaceID || !moduleID"
    :loading="loading"
    class="w-full"
    filter
    fluid
    showClear
    @show="onShow"
  >
    <template #option="{ option }">
      <div class="flex items-center gap-2">
        <span>{{ getOptionLabel(option) }}</span>
      </div>
    </template>
  </Select>
</template>

<script setup>
import { computed, inject, onBeforeUnmount, onMounted, ref, watch } from 'vue'

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
  moduleID: {
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
  if (!props.namespaceID) {
    return 'Select a namespace first'
  }
  if (!props.moduleID) {
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

async function fetchRecords() {
  if (!props.namespaceID || !props.moduleID || !$ComposeAPI) {
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
      namespaceID: props.namespaceID,
      moduleID: props.moduleID,
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

function onShow() {
  if (options.value.length === 0 && props.namespaceID && props.moduleID) {
    fetchRecords()
  }
}

function onSelect(value) {
  selectedRecord.value = value
  emit('update:modelValue', value?.recordID || null)
}

async function loadRecordById(recordID) {
  if (!recordID || !props.namespaceID || !props.moduleID || !$ComposeAPI) return

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
      namespaceID: props.namespaceID,
      moduleID: props.moduleID,
      recordID,
    })
    selectedRecord.value = record
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
  () => [props.namespaceID, props.moduleID],
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
  if (props.namespaceID && props.moduleID) {
    fetchRecords()
  }
  if (props.modelValue && props.namespaceID && props.moduleID) {
    loadRecordById(props.modelValue)
  }
})

onBeforeUnmount(() => {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
})
</script>
