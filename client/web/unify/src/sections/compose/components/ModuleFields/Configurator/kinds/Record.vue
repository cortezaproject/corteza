<template>
  <div class="flex flex-col gap-6">
    <!-- Module selector -->
    <CFormGroup :label="$t('field.kind.record.moduleLabel')" required>
      <Select
        v-model="field.options.moduleID"
        :options="moduleOptions"
        option-label="name"
        option-value="moduleID"
        :placeholder="$t('field.kind.record.modulePlaceholder')"
        filter
        class="w-full"
        @update:model-value="onModuleChange"
      />
    </CFormGroup>

    <template v-if="selectedModule">
      <CFormGroup :label="$t('field.kind.record.moduleField')">
        <Select
          v-model="field.options.labelField"
          :options="fieldOptions"
          option-label="text"
          option-value="value"
          :placeholder="$t('field.kind.record.pickField')"
          show-clear
          class="w-full"
          @update:model-value="onLabelFieldChange"
        />
      </CFormGroup>

      <!-- Record label field (when label field is a Record field itself) -->
      <CFormGroup
        :label="$t('field.kind.record.fieldFromModuleField')"
        v-if="labelField && labelField.kind === 'Record'"
      >
        <Select
          v-model="field.options.recordLabelField"
          :options="labelFieldOptions"
          option-label="text"
          option-value="value"
          :disabled="!labelFieldModule"
          :placeholder="$t('field.kind.record.pickField')"
          show-clear
          class="w-full"
        />
      </CFormGroup>

      <!-- Query fields -->
      <CFormGroup :label="$t('field.kind.record.queryFieldsLabel')">
        <MultiSelect
          v-model="field.options.queryFields"
          :options="queryFieldOptions"
          option-label="text"
          option-value="value"
          :placeholder="$t('field.kind.record.queryFieldsPlaceholder')"
          filter
          class="w-full"
        />
      </CFormGroup>

      <!-- Prefilter -->
      <CFormGroup :label="$t('field.kind.record.prefilterLabel')">
        <CInputExpression
          ref="prefilterInput"
          v-model="field.options.prefilter"
          dialect="ql"
          :scope="scope"
          :query-fields="selectedModule?.fields || []"
          :placeholder="$t('field.kind.record.prefilterPlaceholder')"
        />
        <CExpressionHint :scope="scope" @insert="prefilterInput?.insert($event)" />
      </CFormGroup>
    </template>

    <!-- Multi-value select type -->
    <template v-if="field.isMulti">
      <CFormGroup :label="$t('field.kind.select.optionType.label')">
        <div class="flex flex-col gap-2">
          <div v-for="opt in selectTypeOptions" :key="opt.value" class="flex items-center gap-2">
            <RadioButton
              :input-id="`recordSelectType-${opt.value}`"
              v-model="field.options.selectType"
              :value="opt.value"
              @update:model-value="onSelectTypeChange"
            />
            <label :for="`recordSelectType-${opt.value}`" class="cursor-pointer">
              {{ opt.label }}
            </label>
          </div>
        </div>
      </CFormGroup>

      <div v-if="showAllowDuplicates" class="flex items-center gap-2">
        <Checkbox
          :model-value="!field.options.isUniqueMultiValue"
          inputId="recordAllowDuplicates"
          :binary="true"
          @update:model-value="field.options.isUniqueMultiValue = !$event"
        />
        <label for="recordAllowDuplicates" class="cursor-pointer">
          {{ $t('field.kind.select.allow-duplicates') }}
        </label>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { buildScope, useModuleStore } from '@planetcrust/human-vue'

const { t } = useI18n()

const field = inject('fieldDraft')

const moduleStore = useModuleStore()

// `${record...}` is the record being edited, so it resolves against the module
// this field belongs to; the prefilter itself queries the module the field
// points at.
const fieldModule = inject('fieldModule', null)
const prefilterInput = ref(null)
const scope = computed(() =>
  buildScope({ recordModule: fieldModule?.value || null, hasRecord: true }),
)

// Non-queryable field kinds (same as in lib/js/src/compose/types/module-field/base.ts)
const nonQueryableFieldKinds = ['Number', 'Record', 'User', 'Bool', 'DateTime', 'File', 'Geometry']

const moduleOptions = computed(() => moduleStore.set || [])

const selectedModule = computed(() => {
  if (!field.value.options.moduleID || field.value.options.moduleID === '0') return null
  return moduleStore.getByID(field.value.options.moduleID) || null
})

const fieldOptions = computed(() => {
  if (!selectedModule.value) return []
  return selectedModule.value.fields
    .map(f => ({ value: f.name, text: f.label || f.name, kind: f.kind }))
    .sort((a, b) => a.text.localeCompare(b.text))
})

const queryFieldOptions = computed(() =>
  fieldOptions.value.filter(f => !nonQueryableFieldKinds.includes(f.kind)),
)

// Resolve the actual label field object from the selected module
const labelField = computed(() => {
  if (!field.value.options.labelField || !selectedModule.value) return null
  return selectedModule.value.fields.find(f => f.name === field.value.options.labelField) || null
})

// Resolve the module that the label field points to (when label field is a Record type)
const labelFieldModule = computed(() => {
  if (!labelField.value || labelField.value.kind !== 'Record') return null
  return moduleStore.getByID(labelField.value.options?.moduleID) || null
})

// Field options from the related module (for recordLabelField)
const labelFieldOptions = computed(() => {
  if (!labelFieldModule.value) return []
  return labelFieldModule.value.fields
    .map(f => ({ value: f.name, text: f.label || f.name }))
    .sort((a, b) => a.text.localeCompare(b.text))
})

const duplicatesAllowedTypes = ['default', 'each']
const showAllowDuplicates = computed(() =>
  duplicatesAllowedTypes.includes(field.value.options.selectType),
)

const selectTypeOptions = computed(() => [
  { value: 'default', label: t('field.kind.select.optionType.default') },
  { value: 'multiple', label: t('field.kind.select.optionType.multiple') },
  { value: 'each', label: t('field.kind.select.optionType.each') },
])

function onModuleChange() {
  field.value.options.labelField = ''
  field.value.options.queryFields = []
  field.value.options.prefilter = ''
}

function onLabelFieldChange() {
  field.value.options.queryFields = []
  field.value.options.prefilter = ''
  field.value.options.recordLabelField = ''
}

function onSelectTypeChange(val) {
  if (!duplicatesAllowedTypes.includes(val)) {
    field.value.options.isUniqueMultiValue = true
  }
}
</script>
