<template>
  <Select
    :model-value="modelValue"
    @update:model-value="onSelect"
    :options="options"
    option-label="label"
    option-value="recordID"
    :placeholder="effectivePlaceholder"
    :disabled="disabled || !resolvedNamespaceID || !moduleID"
    :loading="loading"
    class="w-full"
    filter
    :filter-fields="['label', 'recordID']"
    fluid
    showClear
    @filter="onFilter"
    @show="onShow"
  >
    <template v-if="hasNextPage || hasPrevPage" #footer>
      <div class="flex justify-between items-center px-3 py-2 border-t border-surface">
        <Button
          icon="pi pi-angle-left"
          text
          size="small"
          :disabled="!hasPrevPage"
          @click="goToPage(false)"
        />
        <Button
          icon="pi pi-angle-right"
          text
          size="small"
          :disabled="!hasNextPage"
          @click="goToPage(true)"
        />
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
  // CQL prefilter to apply to record listing
  prefilter: {
    type: String,
    default: '',
  },
  // Fields to search in when user types a query
  queryFields: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['update:modelValue'])

const $ComposeAPI = inject('$ComposeAPI')
const $namespace = inject('$namespace', null)

// Resolve namespaceID from prop, falling back to injected $namespace
const resolvedNamespaceID = computed(
  () => props.namespaceID || $namespace?.value?.namespaceID || '',
)

const options = ref([])
const loading = ref(false)

// Pagination state
const nextPageCursor = ref('')
const prevPageCursor = ref('')
const hasNextPage = ref(false)
const hasPrevPage = ref(false)
const currentSearchQuery = ref('')

// Store cancel function for current request
let cancelCurrentRequest = null
let searchTimeout = null

const effectivePlaceholder = computed(() => {
  if (!resolvedNamespaceID.value) {
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

/**
 * Build a CQL filter string from search query and queryFields
 */
function buildQueryFilter(searchQuery) {
  const parts = []

  // Add prefilter if provided
  if (props.prefilter) {
    parts.push(`(${props.prefilter})`)
  }

  // Add search across queryFields
  if (searchQuery && props.queryFields.length > 0) {
    const searchParts = props.queryFields.map(
      field => `${field} LIKE '%${searchQuery.replace(/'/g, "\\'")}%'`,
    )
    parts.push(`(${searchParts.join(' OR ')})`)
  }

  return parts.join(' AND ')
}

async function fetchRecords(searchQuery = '', pageCursor = '') {
  if (!resolvedNamespaceID.value || !props.moduleID || !$ComposeAPI) {
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
    const params = {
      namespaceID: resolvedNamespaceID.value,
      moduleID: props.moduleID,
      limit: 15,
    }

    if (props.labelField) {
      params.sort = `${props.labelField} ASC`
    }

    const filter = buildQueryFilter(searchQuery)
    if (filter) {
      params.filter = filter
    }

    if (pageCursor) {
      params.pageCursor = pageCursor
    }

    const result = await $ComposeAPI.recordList(params)
    const records = result.set || []
    
    options.value = records.map(r => ({ ...r, label: getOptionLabel(r) }))
    
    nextPageCursor.value = result.filter?.nextPage || ''
    prevPageCursor.value = result.filter?.prevPage || ''
    hasNextPage.value = !!nextPageCursor.value
    hasPrevPage.value = !!prevPageCursor.value

    if (props.modelValue && resolvedNamespaceID.value && props.moduleID) {
      if (!options.value.some(r => r.recordID === props.modelValue)) {
        loadRecordById(props.modelValue)
      }
    }
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

function onFilter(event) {
  currentSearchQuery.value = event.value || ''
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    fetchRecords(currentSearchQuery.value)
  }, 300)
}

function onShow() {
  if (resolvedNamespaceID.value && props.moduleID) {
    fetchRecords(currentSearchQuery.value)
  }
}

function goToPage(next) {
  const cursor = next ? nextPageCursor.value : prevPageCursor.value
  if (cursor) {
    fetchRecords(currentSearchQuery.value, cursor)
  }
}

function onSelect(recordID) {
  emit('update:modelValue', recordID || null)
}

async function loadRecordById(recordID) {
  if (!recordID || !resolvedNamespaceID.value || !props.moduleID || !$ComposeAPI) return

  // First check if already in options
  const existing = options.value.find(r => r.recordID === recordID)
  if (existing) {
    return
  }

  // Otherwise fetch it
  loading.value = true
  try {
    const record = await $ComposeAPI.recordRead({
      namespaceID: resolvedNamespaceID.value,
      moduleID: props.moduleID,
      recordID,
    })
    
    if (record) {
      record.label = getOptionLabel(record)
      if (!options.value.find(r => r.recordID === recordID)) {
        options.value = [...options.value, record]
      }
    }
  } catch {
    // Record not found or API error
  } finally {
    loading.value = false
  }
}

// Watch for namespace/module changes - clear selection and reload records
watch(
  () => [resolvedNamespaceID.value, props.moduleID],
  ([newNs, newMod], [oldNs, oldMod]) => {
    if ((oldNs && newNs !== oldNs) || (oldMod && newMod !== oldMod)) {
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
    if (newVal) {
      loadRecordById(newVal)
    }
  },
  { immediate: true },
)

onMounted(() => {
  if (resolvedNamespaceID.value && props.moduleID) {
    fetchRecords()
  }
  if (props.modelValue && resolvedNamespaceID.value && props.moduleID) {
    loadRecordById(props.modelValue)
  }
})

onBeforeUnmount(() => {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }
})
</script>

