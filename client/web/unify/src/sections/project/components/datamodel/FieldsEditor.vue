<template>
  <CFormList
    v-model="fields"
    :columns="columns"
    empty-message="No fields yet."
    :hide-remove="disabled"
  >
    <template #row="{ item: f }">
      <InputText
        :ref="el => nameRefs.set(f.id, el)"
        v-model="f.name"
        placeholder="Field name"
        size="small"
        fluid
        :disabled="disabled"
        :invalid="!!errors[f.id]?.name"
      />
      <Select
        v-model="f.type"
        :options="FIELD_TYPES"
        option-label="label"
        option-value="id"
        size="small"
        fluid
        :disabled="disabled"
        @update:model-value="onType(f)"
      />
      <!-- Unclassified is value null; the placeholder renders it as the
           selected-looking default since Select shows nothing for null. -->
      <Select
        v-if="showSensitivity"
        v-model="f.sensitivity"
        :options="SENSITIVITY_OPTIONS"
        option-label="label"
        option-value="id"
        size="small"
        fluid
        placeholder="Unclassified"
        :disabled="disabled"
      />
      <div class="flex justify-center">
        <ToggleSwitch v-model="f.required" :disabled="disabled" />
      </div>
      <div class="flex justify-center">
        <ToggleSwitch v-model="f.multi" :disabled="disabled" />
      </div>
    </template>

    <!-- Type-specific settings (record reference target, select options) -->
    <template #extra="{ item: f }">
      <Fieldset
        v-if="isRecordRef(f.type)"
        legend="Record"
        class="mt-1"
        :pt="{ content: { class: '!p-3' } }"
      >
        <div class="flex flex-wrap gap-5">
          <CFormGroup label="Target module" required>
            <Select
              :model-value="f.targetModuleId"
              :options="moduleOptions"
              option-label="name"
              option-value="id"
              size="small"
              placeholder="Select module"
              class="w-72"
              :disabled="disabled"
              :invalid="!!errors[f.id]?.target"
              @update:model-value="v => onTarget(f, v)"
            />
          </CFormGroup>
          <CFormGroup label="Label field" description="Represents the referenced record">
            <Select
              v-model="f.labelField"
              :options="labelFieldOptions(f)"
              option-label="label"
              option-value="value"
              size="small"
              placeholder="First field"
              show-clear
              class="w-72"
              :disabled="disabled || !f.targetModuleId"
            />
          </CFormGroup>
        </div>
      </Fieldset>

      <Fieldset
        v-else-if="f.type === 'Select'"
        legend="Options"
        class="mt-1"
        :pt="{ content: { class: '!p-3' } }"
      >
        <div class="flex flex-wrap items-center gap-2">
          <Chip
            v-for="(opt, i) in f.selectOptions"
            :key="`${opt}-${i}`"
            :label="opt"
            :removable="!disabled"
            class="!text-sm"
            @remove="f.selectOptions.splice(i, 1)"
          />
          <InputText
            v-if="!disabled"
            v-model="optionDrafts[f.id]"
            size="small"
            placeholder="Add option, press Enter"
            class="w-56"
            :invalid="!!errors[f.id]?.options"
            @keydown.enter.prevent="addOption(f)"
          />
        </div>
      </Fieldset>
    </template>
  </CFormList>
</template>

<script setup>
import { FIELD_TYPES, isRecordRef } from '@/sections/project/config/fieldTypes'
import { SENSITIVITY_OPTIONS } from '@/sections/project/config/sensitivity'
import { fieldName } from '@/sections/project/utils/fields'
import { computed, reactive } from 'vue'

// Staged field list (two-way; the dialog owns add, CFormList owns remove).
const fields = defineModel({ type: Array, required: true })

const props = defineProps({
  // Modules a Record field can point at.
  moduleOptions: { type: Array, default: () => [] },
  disabled: { type: Boolean, default: false },
  // Sensitivity classification is part of governance — free builds hide it.
  showSensitivity: { type: Boolean, default: true },
  // Per-field validation errors keyed by field id:
  // { [id]: { name?, target?, options? } }. Empty until the dialog's Save runs.
  errors: { type: Object, default: () => ({}) },
})

const columns = computed(() => [
  { label: 'Name', required: true, width: 'minmax(12rem, 1fr)' },
  { label: 'Type', width: '13rem' },
  ...(props.showSensitivity ? [{ label: 'Sensitivity', width: '13rem' }] : []),
  { label: 'Required', width: '6rem', headerClass: 'text-center' },
  { label: 'Multiple', width: '6rem', headerClass: 'text-center' },
])

// Clear type-specific settings when the kind changes away from them.
function onType(f) {
  if (!isRecordRef(f.type)) {
    f.targetModuleId = null
    f.labelField = null
  }
  if (f.type !== 'Select') f.selectOptions = []
}

// Changing the target invalidates the chosen label field.
function onTarget(f, moduleId) {
  f.targetModuleId = moduleId
  f.labelField = null
}

// Fields of the referenced module, by their (predicted) machine name — the
// same derivation the store uses when persisting, so unsaved fields work too.
function labelFieldOptions(f) {
  const target = props.moduleOptions.find(m => m.id === f.targetModuleId)
  return (target?.fields || []).map(tf => ({ label: tf.name, value: fieldName(tf.name) }))
}

// Per-field draft text for the Select option input, keyed by field id.
const optionDrafts = reactive({})

// Name-input refs by field id, so the dialog can focus a freshly added row.
const nameRefs = new Map()
defineExpose({
  focusName(id) {
    const c = nameRefs.get(id)
    const el = c?.$el || c
    el?.focus?.()
  },
})

function addOption(f) {
  const v = (optionDrafts[f.id] || '').trim()
  if (!v) return
  if (!f.selectOptions) f.selectOptions = []
  if (!f.selectOptions.includes(v)) f.selectOptions.push(v)
  optionDrafts[f.id] = ''
}
</script>
