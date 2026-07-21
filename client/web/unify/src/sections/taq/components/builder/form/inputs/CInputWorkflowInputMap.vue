<template>
  <div class="flex flex-col gap-3">
    <div
      v-if="!workflowID"
      class="text-sm text-muted-color"
    >
      {{ t('builder.workflowInputMap.selectWorkflowFirst') }}
    </div>

    <div
      v-else-if="loading"
      class="flex items-center gap-2 text-sm text-muted-color"
    >
      <i class="pi pi-spin pi-spinner" />
      {{ t('builder.workflowInputMap.loading') }}
    </div>

    <div
      v-else-if="rows.length === 0"
      class="text-sm text-muted-color"
    >
      {{ t('builder.workflowInputMap.noInputs') }}
    </div>

    <div
      v-for="(row, index) in rows"
      :key="row.field"
      class="rounded-lg border bg-[--p-content-background] p-3"
    >
      <div class="flex items-center gap-1 mb-2">
        <span class="text-sm font-medium text-color flex-1 capitalize">
          {{ row.label }}
          <span v-if="row.required" class="text-red-500">*</span>
        </span>
        <span v-if="row.types?.length" class="text-xs text-muted-color">
          {{ row.types.join(', ') }}
        </span>
      </div>

      <!-- Reference chip mode -->
      <CReferenceChip
        v-if="row.ref"
        :label="formatRefLabel(row.ref)"
        @click="toggleRowReference(row)"
        @clear="clearRowRef(index)"
      />

      <!-- Literal value mode -->
      <div v-else class="flex gap-1 items-center">
        <InputText
          :model-value="row.value"
          :disabled="disabled"
          class="w-full"
          @update:model-value="onValueUpdate(row, $event)"
        />
        <Button
          icon="pi pi-link"
          text
          rounded
          size="small"
          :severity="isRowReferenceActive(row.field) ? 'primary' : 'secondary'"
          :title="t('builder.form.referenceToggle')"
          @click.stop="toggleRowReference(row)"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import CReferenceChip from '../CReferenceChip.vue'

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({}),
  },
  workflowID: {
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
const $AutomationAPI = inject('$AutomationAPI')

const fields = ref([])
const rows = ref([])
const loading = ref(false)

// Injected from Builder.vue — tracks which argument has the reference panel open
const activeReferenceArgument = inject('activeReferenceArgument', ref(null))
function isRowReferenceActive(fieldName) {
  if (!activeReferenceArgument.value || !fieldName) return false
  return activeReferenceArgument.value.target === fieldName
}

// Build rows from the declared inputs, hydrating each with any existing modelValue.
function buildRows() {
  rows.value = fields.value.map(f => {
    const data = props.modelValue?.[f.name]
    if (data && typeof data === 'object' && data.scope) {
      return { ...f, value: '', ref: { scope: data.scope, source: data.source } }
    }
    const value = data && typeof data === 'object' ? data.value : data
    return { ...f, value: value ?? '', ref: null }
  })
}

// Convert rows back to modelValue { [name]: { value } | { value:'', scope, source } }
function rowsToModelValue() {
  const result = {}
  for (const row of rows.value) {
    if (row.ref) {
      result[row.field] = { value: '', scope: row.ref.scope, source: row.ref.source }
    } else if (row.value !== '' && row.value != null) {
      result[row.field] = { value: row.value }
    }
  }
  return result
}

function emitValue() {
  emit('update:modelValue', rowsToModelValue())
}

function formatRefLabel(ref) {
  return ref?.source || 'Reference'
}

function onValueUpdate(row, val) {
  row.value = val
  emitValue()
}

function clearRowRef(index) {
  rows.value[index].ref = null
  rows.value[index].value = ''
  emitValue()
}

function toggleRowReference(row) {
  emit('toggleRowReference', { target: row.field, types: row.types || [] })
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

async function fetchWorkflowInputs() {
  if (!props.workflowID || !$AutomationAPI) {
    fields.value = []
    buildRows()
    return
  }

  loading.value = true
  try {
    const wf = await $AutomationAPI.workflowRead({ workflowID: props.workflowID })
    fields.value = (wf?.meta?.input || [])
      .filter(f => f.name)
      .map(f => ({
        field: f.name,
        label: f.label || f.name,
        types: f.types || [],
        required: !!f.required,
      }))
  } catch {
    fields.value = []
  } finally {
    loading.value = false
  }
  buildRows()
}

defineExpose({ setRowReference })

// Re-fetch inputs when the selected workflow changes
watch(
  () => props.workflowID,
  (newVal, oldVal) => {
    if (oldVal && newVal !== oldVal) {
      rows.value = []
      emitValue()
    }
    fetchWorkflowInputs()
  },
  { immediate: true },
)

// Re-hydrate rows when modelValue changes externally (e.g. a reference is applied)
watch(
  () => props.modelValue,
  val => {
    if (JSON.stringify(val) !== JSON.stringify(rowsToModelValue())) {
      buildRows()
    }
  },
)
</script>
