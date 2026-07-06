<template>
  <CFormList
    v-model="fields"
    :columns="columns"
    :empty-message="$t('project.dataModel.emptyFields')"
    :hide-remove="disabled"
  >
    <template #row="{ item: f }">
      <InputText
        :ref="el => nameRefs.set(f.id, el)"
        v-model="f.name"
        :placeholder="$t('project.field.fieldNamePlaceholder')"
        size="small"
        fluid
        :disabled="disabled"
        :invalid="!!errors[f.id]?.name"
      />
      <Select
        v-model="f.type"
        :options="fieldTypeOptions"
        option-label="label"
        option-value="id"
        size="small"
        fluid
        :disabled="disabled"
        @update:model-value="onType(f)"
      >
        <template #value="{ value, placeholder }">
          <span v-if="findFieldTypeOption(value)" class="inline-flex items-center gap-2">
            <i
              v-if="findFieldTypeOption(value).icon"
              :class="[findFieldTypeOption(value).icon, 'text-xs text-muted-color']"
            />
            {{ findFieldTypeOption(value).label }}
          </span>
          <span v-else>{{ placeholder }}</span>
        </template>
        <template #option="{ option }">
          <span class="inline-flex items-center gap-2">
            <i v-if="option.icon" :class="[option.icon, 'text-xs text-muted-color']" />
            {{ option.label }}
          </span>
        </template>
      </Select>
      <!-- Unclassified is value null; the placeholder renders it as the
           selected-looking default since Select shows nothing for null. -->
      <Select
        v-if="showSensitivity"
        v-model="f.sensitivity"
        :options="sensitivityOptions"
        option-label="label"
        option-value="id"
        size="small"
        fluid
        :placeholder="$t('project.sensitivity.unclassified')"
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
        :legend="$t('project.field.record')"
        class="mt-1"
        :pt="{ content: { class: '!p-3' } }"
      >
        <div class="flex flex-wrap gap-5">
          <CFormGroup :label="$t('project.field.targetModule')" required>
            <Select
              :model-value="f.targetModuleId"
              :options="moduleOptions"
              option-label="name"
              option-value="id"
              size="small"
              :placeholder="$t('project.field.targetModulePlaceholder')"
              class="w-72"
              :disabled="disabled"
              :invalid="!!errors[f.id]?.target"
              @update:model-value="v => onTarget(f, v)"
            />
          </CFormGroup>
          <CFormGroup :label="$t('project.field.labelField')" :description="$t('project.field.labelFieldDescription')">
            <Select
              v-model="f.labelField"
              :options="labelFieldOptions(f)"
              option-label="label"
              option-value="value"
              size="small"
              :placeholder="$t('project.field.labelFieldPlaceholder')"
              show-clear
              class="w-72"
              :disabled="disabled || !f.targetModuleId"
            />
          </CFormGroup>
        </div>
      </Fieldset>

      <Fieldset
        v-else-if="f.type === 'Select'"
        :legend="$t('project.field.options')"
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
            :placeholder="$t('project.field.addOption')"
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
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

// Staged field list (two-way; the dialog owns add, CFormList owns remove).
const fields = defineModel({ type: Array, required: true })

// Localized option lists for the Selects (id/value is the persisted value).
const fieldTypeOptions = computed(() =>
  FIELD_TYPES.map(ft => ({ id: ft.id, label: t(ft.labelKey), icon: ft.icon })),
)
// Resolve a selected field-type id back to its option (icon + label) for the Select's value slot.
const findFieldTypeOption = id => fieldTypeOptions.value.find(o => o.id === id)
const sensitivityOptions = computed(() =>
  SENSITIVITY_OPTIONS.map(o => ({ id: o.id, label: t(o.labelKey) })),
)

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
  { label: t('project.field.columns.name'), required: true, width: 'minmax(12rem, 1fr)' },
  { label: t('project.field.columns.type'), width: '13rem' },
  ...(props.showSensitivity ? [{ label: t('project.field.columns.sensitivity'), width: '13rem' }] : []),
  { label: t('project.field.columns.required'), width: '6rem', headerClass: 'text-center' },
  { label: t('project.field.columns.multiple'), width: '6rem', headerClass: 'text-center' },
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
