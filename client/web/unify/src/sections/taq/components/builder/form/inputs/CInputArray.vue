<template>
  <div class="flex flex-col gap-3">
    <div
      v-for="(row, index) in rows"
      :key="row.target"
      class="rounded-lg border bg-[--p-content-background] p-3"
    >
      <div class="flex items-center gap-1 mb-2">
        <label class="text-sm font-medium text-color mb-0 flex-1">
          {{ t('builder.arrayInput.item') }} {{ index + 1 }}
        </label>
        <CInputDelete
          text
          size="small"
          :message="
            t('builder.arrayInput.deleteConfirm', 'Are you sure you want to remove this item?')
          "
          :header="t('builder.arrayInput.itemHead', 'Item')"
          @confirm="removeRow(index)"
        />
      </div>

      <!-- Reference chip mode -->
      <CReferenceChip
        v-if="row.ref"
        :label="formatRefLabel(row.ref)"
        @click="emit('toggleRowReference', row.target)"
        @clear="clearRowRef(index)"
      />

      <!-- Literal value mode -->
      <div v-else class="flex gap-1 items-center">
        <CFieldEditor
          :field="{ kind: 'String', name: 'value' }"
          :model-value="row.value"
          :disabled="disabled"
          class="w-full"
          @update:model-value="onFieldValueUpdate(row, $event)"
        />
        <Button
          icon="pi pi-link"
          text
          rounded
          size="small"
          :severity="isRowReferenceActive(row.target) ? 'primary' : 'secondary'"
          :title="t('builder.form.referenceToggle')"
          @click.stop="emit('toggleRowReference', row.target)"
        />
      </div>
    </div>

    <Button
      :label="t('builder.arrayInput.addRow', 'Add Item')"
      icon="pi pi-plus"
      severity="secondary"
      outlined
      size="small"
      :disabled="disabled"
      @click="addRow"
    />
  </div>
</template>

<script setup>
import { inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import CReferenceChip from '../CReferenceChip.vue'
import { CFieldEditor } from '@planetcrust/human-vue/src/components/field'
import { components } from '@planetcrust/human-vue'

const { CInputDelete } = components

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({}),
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue', 'toggleRowReference'])
const { t } = useI18n()
const rows = ref([])

// Unique ID for each target key to ensure stable references during array pushes
let uidCounter = 0
function generateNextTarget() {
  return `item_${++uidCounter}`
}

// Injected from Builder.vue — tracks which argument has the reference panel open
const activeReferenceArgument = inject('activeReferenceArgument', ref(null))

function isRowReferenceActive(target) {
  if (!activeReferenceArgument.value || !target) return false
  return activeReferenceArgument.value.target === target
}

// Convert modelValue { [target]: { value, scope?, source? } } to rows array
function modelValueToRows(val) {
  if (!val || typeof val !== 'object') return []
  return Object.entries(val).map(([target, data]) => {
    if (data && typeof data === 'object' && data.scope) {
      return { target, value: '', ref: { scope: data.scope, source: data.source } }
    }
    const value = data && typeof data === 'object' ? data.value : data
    return { target, value: value ?? '' }
  })
}

// Convert rows to modelValue format (ordered object)
function rowsToModelValue() {
  const result = {}
  for (const row of rows.value) {
    if (row.ref) {
      result[row.target] = { value: '', scope: row.ref.scope, source: row.ref.source }
    } else {
      result[row.target] = { value: row.value ?? '' }
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
 * Set a reference on a specific row by target id.
 * Called from outside (via expose) when user selects a reference from the panel.
 */
function setRowReference(target, { scope, source }) {
  const row = rows.value.find(r => r.target === target)
  if (row) {
    row.ref = { scope, source }
    row.value = ''
    emitValue()
  }
}

function addRow() {
  rows.value.push({ target: generateNextTarget(), value: '' })
  emitValue()
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

defineExpose({ setRowReference })

// Initialize rows from modelValue
watch(
  () => props.modelValue,
  val => {
    const incoming = modelValueToRows(val)
    const current = rowsToModelValue()
    // Update UID counter to be beyond the highest incoming index to avoid collisions
    incoming.forEach(row => {
      if (row.target.startsWith('item_')) {
        const id = parseInt(row.target.split('_')[1])
        if (!isNaN(id) && id > uidCounter) uidCounter = id
      }
    })

    // Only update if externally changed to avoid cursor jumps
    if (JSON.stringify(val) !== JSON.stringify(current)) {
      rows.value = incoming.length ? incoming : []
    }
  },
  { immediate: true },
)
</script>
