<template>
  <div class="flex flex-col gap-3">
    <div
      v-for="(row, index) in rows"
      :key="index"
      class="rounded-lg border bg-[--p-content-background] p-3"
    >
      <div class="flex items-center gap-1 mb-2">
        <Select
          v-model="row.field"
          :options="getAvailableFields(index)"
          option-label="label"
          option-value="name"
          :placeholder="t('builder.fieldValueMap.selectField')"
          :disabled="disabled || !moduleID"
          class="flex-1"
          @update:model-value="onRowChange"
        />
        <CInputDelete
          text
          size="small"
          :message="t('builder.fieldValueMap.deleteConfirm')"
          :header="row.field || t('builder.fieldValueMap.selectField')"
          @confirm="removeRow(index)"
        />
      </div>

      <!-- Reference chip mode -->
      <CReferenceChip
        v-if="row.ref"
        :label="formatRefLabel(row.ref)"
        @click="emit('toggleRowReference', row.field)"
        @clear="clearRowRef(index)"
      />

      <!-- Literal value mode -->
      <div v-else class="flex gap-1 items-center">
        <CFieldEditor
          :field="getFieldDef(row.field)"
          :model-value="row.value"
          :disabled="disabled || !row.field"
          :add-label="t('builder.fieldValueMap.addValue')"
          class="w-full"
          @update:model-value="onFieldValueUpdate(row, $event)"
        />
        <Button
          v-if="row.field"
          icon="pi pi-link"
          text
          rounded
          size="small"
          :severity="isRowReferenceActive(row.field) ? 'primary' : 'secondary'"
          :title="t('builder.form.referenceToggle')"
          @click.stop="emit('toggleRowReference', row.field)"
        />
      </div>
    </div>

    <Button
      :label="t('builder.fieldValueMap.addRow')"
      icon="pi pi-plus"
      severity="secondary"
      outlined
      size="small"
      :disabled="disabled || !moduleID"
      @click="addRow"
    />
  </div>
</template>

<script setup>
import { inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import CReferenceChip from '../CReferenceChip.vue'
import { CFieldEditor } from '@planetcrust/human-vue/src/components/field'
import { components, useModuleStore } from '@planetcrust/human-vue'
import { systemFields as moduleSystemFields } from '@planetcrust/human-js/src/compose/types/module'

const { CInputDelete } = components

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({}),
  },
  namespaceID: {
    type: [String, Number],
    default: null,
  },
  moduleID: {
    type: [String, Number],
    default: null,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue', 'toggleRowReference'])
const { t } = useI18n()
const moduleStore = useModuleStore()
const fields = ref([])
const rows = ref([])
const loading = ref(false)

// Injected from Builder.vue — tracks which argument has the reference panel open
const activeReferenceArgument = inject('activeReferenceArgument', ref(null))
function isRowReferenceActive(fieldName) {
  if (!activeReferenceArgument.value || !fieldName) return false
  return activeReferenceArgument.value.target === fieldName
}

// Convert modelValue { [field]: { value, scope?, source? } } to rows
function modelValueToRows(val) {
  if (!val || typeof val !== 'object') return []
  return Object.entries(val).map(([field, data]) => {
    if (data && typeof data === 'object' && data.scope) {
      return { field, value: '', ref: { scope: data.scope, source: data.source } }
    }
    const value = data && typeof data === 'object' ? data.value : data
    return { field, value: value ?? '' }
  })
}

// Convert rows to modelValue format
function rowsToModelValue() {
  const result = {}
  for (const row of rows.value) {
    if (row.field) {
      if (row.ref) {
        result[row.field] = { value: '', scope: row.ref.scope, source: row.ref.source }
      } else {
        result[row.field] = { value: row.value ?? '' }
      }
    }
  }
  return result
}

function formatRefLabel(ref) {
  if (!ref?.source) return 'Reference'
  return ref.source
}

function clearRowRef(index) {
  rows.value[index].ref = null
  rows.value[index].value = ''
  emitValue()
}

/**
 * Set a reference on a specific row by field name.
 * Called from outside (via expose) when user selects a reference from the panel.
 */
function setRowReference(fieldName, { scope, source }) {
  const row = rows.value.find(r => r.field === fieldName)
  if (row) {
    row.ref = { scope, source }
    row.value = ''
    emitValue()
  }
}

function getAvailableFields(currentIndex) {
  const selectedFields = new Set(
    rows.value
      .filter((_, i) => i !== currentIndex)
      .map(r => r.field)
      .filter(Boolean),
  )
  return fields.value.filter(f => !selectedFields.has(f.name))
}

function getFieldDef(fieldName) {
  if (!fieldName) {
    return { kind: 'String', name: '', label: '', options: {}, isMulti: false, isRequired: false }
  }
  const f = fields.value.find(fd => fd.name === fieldName)
  if (!f) {
    return {
      kind: 'String',
      name: fieldName,
      label: fieldName,
      options: {},
      isMulti: false,
      isRequired: false,
    }
  }
  return f
}

function addRow() {
  rows.value.push({ field: null, value: '' })
}

function removeRow(index) {
  rows.value.splice(index, 1)
  emitValue()
}

function onFieldValueUpdate(row, val) {
  row.value = val
  onRowChange()
}

function onRowChange() {
  emitValue()
}

function emitValue() {
  emit('update:modelValue', rowsToModelValue())
}

async function fetchModuleFields() {
  if (!props.namespaceID || !props.moduleID) {
    fields.value = []
    return
  }

  loading.value = true
  try {
    const mod = await moduleStore.findByID({ namespaceID: props.namespaceID, moduleID: props.moduleID })
    if (!mod) {
      fields.value = []
      return
    }
    const WRITABLE_SYSTEM_FIELDS = new Set(['ownedBy'])

    const systemEntries = moduleSystemFields
      .filter(f => WRITABLE_SYSTEM_FIELDS.has(f.name))
      .map(f => ({
        name: f.name,
        label: f.label || f.name,
        kind: f.kind,
        options: f.options || {},
        isMulti: !!f.isMulti,
        isRequired: !!f.isRequired,
      }))

    fields.value = [
      ...(mod.fields || [])
        .filter(f => !f.isSystem)
        .map(f => ({
          name: f.name,
          label: f.label || f.name,
          kind: f.kind,
          options: f.options || {},
          isMulti: !!f.isMulti,
          isRequired: !!f.isRequired,
        })),
      ...systemEntries,
    ]
  } catch {
    fields.value = []
  } finally {
    loading.value = false
  }
}

defineExpose({ setRowReference })

// Initialize rows from modelValue
watch(
  () => props.modelValue,
  val => {
    const incoming = modelValueToRows(val)
    const current = rowsToModelValue()
    // Only update if externally changed to avoid cursor jumps
    if (JSON.stringify(val) !== JSON.stringify(current)) {
      rows.value = incoming.length ? incoming : []
    }
  },
  { immediate: true },
)

// Re-fetch fields when module changes; clear rows
watch(
  () => props.moduleID,
  (newVal, oldVal) => {
    if (oldVal && newVal !== oldVal) {
      rows.value = []
      emitValue()
    }
    if (newVal && props.namespaceID) {
      fetchModuleFields()
    } else {
      fields.value = []
    }
  },
  { immediate: true },
)

// Re-fetch when namespace changes
watch(
  () => props.namespaceID,
  (newVal, oldVal) => {
    if (oldVal && newVal !== oldVal) {
      fields.value = []
      rows.value = []
      emitValue()
    }
  },
)
</script>
