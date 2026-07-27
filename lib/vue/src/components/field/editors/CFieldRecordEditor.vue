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
    :record-label-field="field.options.recordLabelField || ''"
    :prefilter="field.options.prefilter || ''"
    :query-fields="field.options.queryFields || []"
    placeholder=""
    :disabled="disabled"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script setup>
import { computed, inject, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { compose } from '@planetcrust/human-js'
import { useModuleStore } from '../../../stores/useModuleStore'
import { useRecordStore } from '../../../stores/useRecordStore'
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
const $recordContext = inject('$recordContext', null)
const $Auth = inject('$Auth', null)
const moduleStore = useModuleStore()
const recordStore = useRecordStore()

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

// When labelField is itself a Record field, labels come from the nested
// module's recordLabelField instead of the raw nested recordID
const nestedModuleID = computed(() => {
  const { moduleID, labelField, recordLabelField } = props.field.options || {}
  if (!moduleID || !labelField || !recordLabelField) return ''
  const lf = moduleStore.getByID(moduleID)?.fields?.find(f => f.name === labelField)
  return lf?.kind === 'Record' ? lf.options?.moduleID || '' : ''
})

function nestedRecordLabel(nestedID) {
  if (!nestedModuleID.value || !nestedID) return ''
  const rec = recordStore.getByID(nestedID)
  if (!rec) return ''
  const rlf = props.field.options.recordLabelField
  const v = Array.isArray(rec.values)
    ? rec.values.find(v => v.name === rlf)?.value
    : rec.values?.[rlf]
  return (Array.isArray(v) ? v.filter(Boolean).join(', ') : v) || ''
}

function resolveNestedLabels(records) {
  if (!nestedModuleID.value || !records.length) return
  const lf = props.field.options.labelField
  const recordIDs = records.map(r => r.values?.find?.(v => v.name === lf)?.value).filter(Boolean)
  if (!recordIDs.length) return
  recordStore.resolveRecordLabels({
    namespaceID: namespaceID.value,
    moduleID: nestedModuleID.value,
    recordIDs,
  })
}

function getOptionLabel(record) {
  if (!record) return ''
  const lf = props.field.options?.labelField
  if (lf && record.values) {
    const v = record.values.find(v => v.name === lf)
    if (v?.value) return nestedRecordLabel(v.value) || v.value
  }
  if (record.values?.length) {
    const first = record.values.find(v => v.value)
    if (first?.value) return first.value
  }
  return `Record ${record.recordID}`
}

// Interpolates ${record...}/${recordID}/${ownerID}/${userID} expressions in the
// prefilter. Falls back to the raw prefilter on failure so existing (non-templated)
// configurations can't start throwing.
function resolvePrefilter() {
  const prefilter = props.field.options?.prefilter
  if (!prefilter) return ''

  try {
    const record = $recordContext?.value || null
    const user = $Auth?.user || {}
    return compose.interpolateTemplate(prefilter, {
      record,
      user,
      recordID: record?.recordID || '0',
      ownerID: record?.ownedBy || '0',
      userID: user?.userID || '0',
    })
  } catch {
    return prefilter
  }
}

function buildFilter(searchQuery) {
  const parts = []
  const prefilter = resolvePrefilter()
  if (prefilter) {
    parts.push(`(${prefilter})`)
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
    resolveNestedLabels(recordOptions.value)
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
        // re-check after await: concurrent resolves may have added it meanwhile
        if (record && !recordOptions.value.find(r => r.recordID === id)) {
          recordOptions.value = [...recordOptions.value, record]
        }
      } catch {
        // ignore
      }
    }
  }

  resolveNestedLabels(recordOptions.value)
}

onMounted(async () => {
  if (isMultipleType.value) {
    await fetchRecords()
    await resolveExisting()
  }
})

// Resolve selected IDs whenever the value changes (workflow-prefilled
// values, external mutations) — not just at mount
watch(multiValue, () => {
  if (isMultipleType.value) resolveExisting()
})

onBeforeUnmount(() => {
  if (cancelCurrentRequest) cancelCurrentRequest()
  if (searchTimeout) clearTimeout(searchTimeout)
})
</script>
