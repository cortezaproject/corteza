<template>
  <div class="flex flex-col gap-3">
    <div v-if="loading" class="text-sm text-muted-color">
      {{ t('builder.sheetRow.loading', 'Loading columns…') }}
    </div>

    <template v-else-if="columns.length">
      <!-- Only the columns the user picked -->
      <div
        v-for="i in pickedInOrder"
        :key="i"
        class="rounded-lg border bg-[--p-content-background] p-3"
      >
        <div class="flex items-center justify-between mb-2">
          <label class="text-sm font-medium text-color">{{ columnLabel(i) }}</label>
          <Button
            icon="pi pi-times"
            text
            rounded
            size="small"
            severity="secondary"
            :disabled="disabled"
            @click="removeColumn(i)"
          />
        </div>

        <CReferenceChip
          v-if="isRef(i)"
          :label="refLabel(i)"
          @click="toggleRef(i)"
          @clear="clearRef(i)"
        />
        <div v-else class="flex gap-1 items-center">
          <InputText
            :model-value="valueFor(i)"
            :disabled="disabled"
            class="w-full"
            @update:model-value="onValueUpdate(i, $event)"
          />
          <Button
            icon="pi pi-link"
            text
            rounded
            size="small"
            :severity="isActive(i) ? 'primary' : 'secondary'"
            :title="t('builder.form.referenceToggle')"
            @click.stop="toggleRef(i)"
          />
        </div>
      </div>

      <!-- Add another column -->
      <Select
        v-if="available.length"
        :model-value="null"
        :options="available"
        option-label="label"
        option-value="value"
        :placeholder="t('builder.sheetRow.addColumn', '+ Add column')"
        :disabled="disabled"
        class="w-full"
        @update:model-value="addColumn"
      />
    </template>

    <!-- No columns: fall back to a plain list -->
    <div v-else class="flex flex-col gap-2">
      <div v-if="error" class="text-sm text-red-500">{{ error }}</div>
      <div v-else class="text-sm text-muted-color">{{ fallbackHint }}</div>
      <CInputArray
        :model-value="modelValue"
        :disabled="disabled"
        @update:model-value="emit('update:modelValue', $event)"
      />
    </div>
  </div>
</template>

<script setup>
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import CInputArray from './CInputArray.vue'
import CReferenceChip from '../CReferenceChip.vue'
import { useSheetColumns } from '@/sections/taq/composables/useSheetColumns'

const props = defineProps({
  modelValue: { type: Object, default: () => ({}) },
  disabled: { type: Boolean, default: false },
})

const emit = defineEmits(['update:modelValue', 'toggleRowReference'])
const { t } = useI18n()

const { columns, loading, error, spreadsheetId } = useSheetColumns()
const activeReferenceArgument = inject('activeReferenceArgument', ref(null))

const fallbackHint = computed(() =>
  !spreadsheetId.value
    ? t('builder.sheetRow.pickSheet', 'Pick a spreadsheet and tab to load its columns.')
    : t('builder.sheetRow.noColumns', 'No header row found — add one, or enter values below.'),
)

function keyFor(index) {
  return `col_${index}`
}
function columnLabel(index) {
  return columns.value[index] || `${t('builder.sheetRow.column', 'Column')} ${index + 1}`
}

// Which columns are shown. Seeded from columns that already have a value.
const picked = ref(new Set())
let seeded = false
watch(
  [columns, () => props.modelValue],
  () => {
    if (seeded || !columns.value.length) return
    const s = new Set()
    columns.value.forEach((_, i) => {
      const e = props.modelValue?.[keyFor(i)]
      if (e && (e.value || e.scope)) s.add(i)
    })
    picked.value = s
    seeded = true
  },
  { immediate: true },
)

const pickedInOrder = computed(() => columns.value.map((_, i) => i).filter(i => picked.value.has(i)))
const available = computed(() =>
  columns.value
    .map((_, i) => ({ label: columnLabel(i), value: i }))
    .filter(o => !picked.value.has(o.value)),
)

function addColumn(index) {
  if (index == null) return
  picked.value = new Set(picked.value).add(index)
}
function removeColumn(index) {
  const next = new Set(picked.value)
  next.delete(index)
  picked.value = next
  emitColumns(keyFor(index), { value: '' })
}

function entryFor(index) {
  return props.modelValue?.[keyFor(index)] || {}
}
function valueFor(index) {
  return entryFor(index).value ?? ''
}
function isRef(index) {
  return !!entryFor(index).scope
}
function refLabel(index) {
  return entryFor(index).source || t('builder.form.reference', 'Reference')
}
function isActive(index) {
  return activeReferenceArgument.value?.target === keyFor(index)
}
function toggleRef(index) {
  emit('toggleRowReference', keyFor(index))
}

// Emit every column in order (unpicked stay blank), applying one override.
function emitColumns(overrideKey, entry) {
  const next = {}
  columns.value.forEach((_, i) => {
    const key = keyFor(i)
    next[key] = key === overrideKey ? entry : (props.modelValue?.[key] || { value: '' })
  })
  emit('update:modelValue', next)
}

function onValueUpdate(index, val) {
  emitColumns(keyFor(index), { value: val })
}
function clearRef(index) {
  emitColumns(keyFor(index), { value: '' })
}
</script>
