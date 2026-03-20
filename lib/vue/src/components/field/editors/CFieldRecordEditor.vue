<template>
  <!-- Multi-value: MultiSelect mode (selectType === 'multiple') -->
  <MultiSelect
    v-if="isMultipleType"
    :model-value="multiValue"
    :options="recordOptions"
    :option-label="getOptionLabel"
    option-value="recordID"
    :placeholder="$t('field.kind.record.suggestionPlaceholder')"
    :disabled="disabled"
    :loading="loading"
    filter
    class="w-full"
    @filter="onFilter"
    @update:model-value="onMultiSelectChange"
  >
    <template #option="{ option }">
      {{ getOptionLabel(option) }}
    </template>
  </MultiSelect>

  <!-- Single-value or default/each multi: CInputRecord per slot -->
  <CInputRecord
    v-else
    :model-value="modelValue"
    :namespace-i-d="namespaceID"
    :module-i-d="field.options.moduleID"
    :label-field="field.options.labelField"
    :prefilter="field.options.prefilter || ''"
    :query-fields="field.options.queryFields || []"
    placeholder=""
    :disabled="disabled"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script setup>
import { computed, inject, onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import CInputRecord from '../../input/CInputRecord.vue'

const { t: $t } = useI18n()

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
  modelValue: {
    type: [String, Array],
    default: '',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  namespace: {
    type: Object,
    default: () => ({}),
  },
})

const emit = defineEmits(['update:modelValue'])

const $ComposeAPI = inject('$ComposeAPI')
const $namespace = inject('$namespace', null)

const namespaceID = computed(
  () => props.namespace?.namespaceID || $namespace?.value?.namespaceID || '',
)

const isMultipleType = computed(
  () => props.field.isMulti && props.field.options?.selectType === 'multiple',
)

// --- MultiSelect mode state ---
const recordOptions = ref([])
const loading = ref(false)
let cancelCurrentRequest = null
let searchTimeout = null

const multiValue = computed(() => {
  if (Array.isArray(props.modelValue)) return props.modelValue.filter(Boolean)
  return props.modelValue ? [props.modelValue] : []
})

function getOptionLabel(record) {
  if (!record) return ''
  const lf = props.field.options?.labelField
  if (lf && record.values) {
    const v = record.values.find(v => v.name === lf)
    if (v?.value) return v.value
  }
  if (record.values?.length) {
    const first = record.values.find(v => v.value)
    if (first?.value) return first.value
  }
  return `Record ${record.recordID}`
}

function buildFilter(searchQuery) {
  const parts = []
  if (props.field.options?.prefilter) {
    parts.push(`(${props.field.options.prefilter})`)
  }
  const qf = props.field.options?.queryFields || []
  if (searchQuery && qf.length > 0) {
    const sp = qf.map(f => `${f} LIKE '%${searchQuery.replace(/'/g, "\\'")}%'`)
    parts.push(`(${sp.join(' OR ')})`)
  }
  return parts.join(' AND ')
}

async function fetchRecords(searchQuery = '') {
  const moduleID = props.field.options?.moduleID
  if (!namespaceID.value || !moduleID || !$ComposeAPI) return

  if (cancelCurrentRequest) {
    cancelCurrentRequest()
    cancelCurrentRequest = null
  }

  loading.value = true
  try {
    const params = {
      namespaceID: namespaceID.value,
      moduleID,
      limit: 50,
    }
    const filter = buildFilter(searchQuery)
    if (filter) params.filter = filter

    const { response, cancel } = $ComposeAPI.recordListCancellable(params)
    cancelCurrentRequest = cancel
    const result = await response()
    recordOptions.value = result.set || []
  } catch (e) {
    if (e?.message !== 'canceled') recordOptions.value = []
  } finally {
    loading.value = false
    cancelCurrentRequest = null
  }
}

function onFilter(event) {
  const q = event.value || ''
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => fetchRecords(q), 300)
}

function onMultiSelectChange(values) {
  if (props.field.options?.isUniqueMultiValue) {
    values = [...new Set(values)]
  }
  emit('update:modelValue', values)
}

// Resolve existing record IDs so MultiSelect displays labels
async function resolveExisting() {
  const ids = multiValue.value
  if (!ids.length || !$ComposeAPI) return
  const moduleID = props.field.options?.moduleID
  if (!namespaceID.value || !moduleID) return

  for (const id of ids) {
    if (id && !recordOptions.value.find(r => r.recordID === id)) {
      try {
        const record = await $ComposeAPI.recordRead({
          namespaceID: namespaceID.value,
          moduleID,
          recordID: id,
        })
        if (record) {
          recordOptions.value = [...recordOptions.value, record]
        }
      } catch {
        // ignore
      }
    }
  }
}

onMounted(async () => {
  if (isMultipleType.value) {
    await fetchRecords()
    await resolveExisting()
  }
})

onBeforeUnmount(() => {
  if (cancelCurrentRequest) cancelCurrentRequest()
  if (searchTimeout) clearTimeout(searchTimeout)
})
</script>
