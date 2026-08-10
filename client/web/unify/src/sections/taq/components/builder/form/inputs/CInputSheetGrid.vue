<template>
  <div class="flex flex-col gap-3">
    <div v-if="loading" class="text-sm text-muted-color">
      {{ t('builder.sheetGrid.loading', 'Loading columns…') }}
    </div>

    <template v-else-if="columns.length">
      <!-- Pick which columns to fill (shared across all rows) -->
      <div class="flex flex-wrap gap-1 items-center">
        <Chip
          v-for="col in pickedInOrder"
          :key="col"
          :label="columnLabel(col)"
          removable
          @remove="removeColumn(col)"
        />
        <Select
          v-if="available.length"
          :model-value="null"
          :options="available"
          option-label="label"
          option-value="value"
          :placeholder="t('builder.sheetGrid.addColumn', '+ Add column')"
          :disabled="disabled"
          class="min-w-40"
          @update:model-value="addColumn"
        />
      </div>

      <!-- One card per row -->
      <div
        v-for="(row, ri) in rows"
        :key="row.rowId"
        class="rounded-lg border bg-[--p-content-background] p-3 flex flex-col gap-2"
      >
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium text-color">
            {{ t('builder.sheetGrid.row', 'Row') }} {{ ri + 1 }}
          </span>
          <Button
            icon="pi pi-trash"
            text
            rounded
            size="small"
            severity="secondary"
            :disabled="disabled"
            @click="removeRow(ri)"
          />
        </div>

        <div v-if="!pickedInOrder.length" class="text-xs text-muted-color">
          {{ t('builder.sheetGrid.pickColumns', 'Add a column above to start filling rows.') }}
        </div>

        <div v-for="col in pickedInOrder" :key="col" class="flex flex-col gap-1">
          <label class="text-xs font-medium text-muted-color">{{ columnLabel(col) }}</label>

          <CReferenceChip
            v-if="isRef(row.rowId, col)"
            :label="refLabel(row.rowId, col)"
            @click="toggleRef(row.rowId, col)"
            @clear="clearRef(row.rowId, col)"
          />
          <div v-else class="flex gap-1 items-center">
            <InputText
              :model-value="valueFor(row.rowId, col)"
              :disabled="disabled"
              class="w-full"
              @update:model-value="onCellUpdate(row.rowId, col, $event)"
            />
            <Button
              icon="pi pi-link"
              text
              rounded
              size="small"
              :severity="isActive(row.rowId, col) ? 'primary' : 'secondary'"
              :title="t('builder.form.referenceToggle')"
              @click.stop="toggleRef(row.rowId, col)"
            />
          </div>
        </div>
      </div>

      <Button
        :label="t('builder.sheetGrid.addRow', 'Add row')"
        icon="pi pi-plus"
        severity="secondary"
        outlined
        size="small"
        :disabled="disabled"
        @click="addRow"
      />
    </template>

    <!-- No columns: show a hint -->
    <div v-else class="text-sm text-muted-color">{{ fallbackHint }}</div>
  </div>
</template>

<script setup>
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import CReferenceChip from '../CReferenceChip.vue'
import { useSheetColumns } from '@/sections/taq/composables/useSheetColumns'

const props = defineProps({
  modelValue: { type: Object, default: () => ({}) },
  disabled: { type: Boolean, default: false },
})

const emit = defineEmits(['update:modelValue', 'toggleRowReference'])
const { t } = useI18n()

const { columns, loading, spreadsheetId } = useSheetColumns()
const activeReferenceArgument = inject('activeReferenceArgument', ref(null))

const fallbackHint = computed(() =>
  !spreadsheetId.value
    ? t('builder.sheetGrid.pickSheet', 'Pick a spreadsheet and tab to load its columns.')
    : t('builder.sheetGrid.noColumns', 'No header row found — add one to the sheet.'),
)

const CELL = /^r(\d+)c(\d+)$/

function keyFor(rowId, col) {
  return `r${rowId}c${col}`
}
function columnLabel(col) {
  return columns.value[col] || `${t('builder.sheetGrid.column', 'Column')} ${col + 1}`
}

// rows: [{ rowId }]. picked: which columns are shown.
const rows = ref([])
const picked = ref(new Set())
let uid = -1

const pickedInOrder = computed(() => columns.value.map((_, i) => i).filter(i => picked.value.has(i)))
const available = computed(() =>
  columns.value
    .map((_, i) => ({ label: columnLabel(i), value: i }))
    .filter(o => !picked.value.has(o.value)),
)

function entryFor(rowId, col) {
  return props.modelValue?.[keyFor(rowId, col)] || {}
}
function valueFor(rowId, col) {
  return entryFor(rowId, col).value ?? ''
}
function isRef(rowId, col) {
  return !!entryFor(rowId, col).scope
}
function refLabel(rowId, col) {
  return entryFor(rowId, col).source || t('builder.form.reference', 'Reference')
}
function isActive(rowId, col) {
  return activeReferenceArgument.value?.target === keyFor(rowId, col)
}
function toggleRef(rowId, col) {
  emit('toggleRowReference', keyFor(rowId, col))
}

// Keep existing cells of live rows and picked columns; apply one override.
function emitCells(overrideKey, entry) {
  const live = new Set(rows.value.map(r => r.rowId))
  const next = {}
  for (const [key, val] of Object.entries(props.modelValue || {})) {
    const m = CELL.exec(key)
    if (m && live.has(+m[1]) && picked.value.has(+m[2])) next[key] = val
  }
  if (overrideKey) next[overrideKey] = entry
  emit('update:modelValue', next)
}

function onCellUpdate(rowId, col, val) {
  emitCells(keyFor(rowId, col), { value: val })
}
function clearRef(rowId, col) {
  emitCells(keyFor(rowId, col), { value: '' })
}
function addRow() {
  rows.value.push({ rowId: ++uid })
}
function removeRow(ri) {
  rows.value.splice(ri, 1)
  emitCells()
}
function addColumn(col) {
  if (col == null) return
  picked.value = new Set(picked.value).add(col)
}
function removeColumn(col) {
  const next = new Set(picked.value)
  next.delete(col)
  picked.value = next
  emitCells()
}

// Seed rows and picked columns from a saved value.
let seeded = false
watch(
  [columns, () => props.modelValue],
  () => {
    const val = props.modelValue || {}
    const rowIds = new Set()
    const cols = new Set()
    for (const key of Object.keys(val)) {
      const m = CELL.exec(key)
      if (!m) continue
      rowIds.add(+m[1])
      const e = val[key]
      if (e && (e.value || e.scope)) cols.add(+m[2])
    }

    const sortedRows = [...rowIds].sort((a, b) => a - b)
    sortedRows.forEach(id => {
      if (id > uid) uid = id
    })
    const shown = new Set(rows.value.map(r => r.rowId))
    const missing = sortedRows.filter(id => !shown.has(id))
    if (missing.length) {
      rows.value = [...rows.value, ...missing.map(rowId => ({ rowId }))].sort(
        (a, b) => a.rowId - b.rowId,
      )
    }
    if (rows.value.length === 0) rows.value = [{ rowId: ++uid }]

    if (!seeded && columns.value.length) {
      picked.value = cols
      seeded = true
    }
  },
  { immediate: true },
)
</script>
