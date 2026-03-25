<template>
  <div class="flex flex-col gap-5">
    <!-- General -->
    <div class="flex flex-col gap-3">
      <h5 class="text-lg font-semibold text-primary m-0">
        {{ $t('block.recordList.record.generalLabel') }}
      </h5>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <!-- Module (read-only since it comes from the page) -->
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.general.module') }}</label>
          <InputText :model-value="moduleName" disabled class="w-full" />
        </div>

        <!-- Inline record edit -->
        <CInputSwitch v-model="inlineEditEnabled" :label="$t('block.record.inlineEdit.enabled')" />

        <!-- Allow adding fields to inline edit -->
        <CInputSwitch v-model="inlineRecordEditAllowAddField" :label="$t('block.record.inlineEdit.allowAddField')" />

        <!-- Horizontal form layout -->
        <CInputSwitch v-model="horizontalLayout" :label="$t('block.record.horizontalFormLayout')" :disabled="layoutMode === 'noWrap'" />

        <!-- Field layout mode -->
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.record.fieldsLayoutMode.label') }}</label>
          <Select
            v-model="layoutMode"
            :options="fieldLayoutOptions"
            option-label="text"
            option-value="value"
            class="w-full"
          />
        </div>

        <!-- Reference record field -->
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.record.referenceRecordField') }}</label>
          <Select
            v-model="referenceField"
            :options="recordSelectorFields"
            :option-label="f => f.label || f.name"
            option-value="fieldID"
            :placeholder="$t('block.record.referenceRecordFieldPlaceholder')"
            class="w-full"
            show-clear
            :disabled="!selectedModule"
          />
          <small class="text-muted-color">{{ $t('block.record.referenceRecordFieldDescription') }}</small>
        </div>
      </div>
    </div>

    <template v-if="fieldPickerModule">
      <Divider />

      <!-- Fields -->
      <div class="flex flex-col gap-3">
        <h5 class="text-lg font-semibold text-primary m-0">
          {{ $t('block.general.fields') }}
        </h5>

        <CFieldPicker
          :all-fields="allModuleFields"
          :model-value="selectedFieldNames"
          :available-label="$t('field.selector.available')"
          :selected-label="$t('field.selector.selected')"
          :select-all-label="$t('field.selector.selectAll')"
          :unselect-all-label="$t('field.selector.unselectAll')"
          :search-placeholder="$t('field.selector.search')"
          :no-items-label="$t('field.no-items-found')"
          @update:model-value="onFieldPickerUpdate"
        />
      </div>

      <template v-if="isRecordFieldUsedConfigured">
        <Divider />

        <!-- Record Display Options (only shown when Record-kind fields are used) -->
        <div class="flex flex-col gap-3">
          <h5 class="text-lg font-semibold text-primary m-0">
            {{ $t('block.record.recordDisplay.label') }}
          </h5>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
            <div class="flex flex-col gap-1">
              <label class="text-primary font-medium text-sm">{{ $t('block.record.recordSelectorDisplayOptions') }}</label>
              <Select
                v-model="recordSelectorDisplayOption"
                :options="displayOptions"
                option-label="text"
                option-value="value"
                class="w-full"
              />
            </div>

            <CInputSwitch v-model="recordSelectorShowAddRecordButton" :label="$t('block.record.recordSelectorCanAddRecord')" />

            <div class="flex flex-col gap-1">
              <label class="text-primary font-medium text-sm">{{ $t('block.record.recordSelectorAddRecordDisplayOption') }}</label>
              <Select
                v-model="recordSelectorAddRecordDisplayOption"
                :options="displayOptions"
                option-label="text"
                option-value="value"
                class="w-full"
                :disabled="!recordSelectorShowAddRecordButton"
              />
            </div>
          </div>
        </div>
      </template>

      <Divider />

      <!-- Field conditions -->
      <div class="flex flex-col gap-3">
        <div class="flex items-center gap-2">
          <h5 class="text-lg font-semibold text-primary m-0">
            {{ $t('block.record.fieldConditions.label') }}
          </h5>
          <i class="pi pi-exclamation-triangle text-orange-500 text-sm" v-tooltip="$t('block.record.fieldConditions.tooltip.performance')" />
        </div>

        <!-- Clear all on hide -->
        <CInputSwitch v-model="clearConditionalFieldsOnHide" :label="$t('block.record.fieldConditions.clearAllOnHide')" />
        <small class="text-muted-color -mt-2">{{ $t('block.record.fieldConditions.clearAllOnHideDescription') }}</small>

        <div class="flex flex-col gap-2">
          <div
            v-for="(condition, i) in fieldConditions"
            :key="i"
            class="flex items-center gap-2"
          >
            <Select
              v-model="condition.field"
              :options="conditionFieldOptions"
              option-label="text"
              option-value="value"
              :placeholder="$t('block.record.fieldConditions.selectPlaceholder')"
              class="flex-1"
              filter
            />
            <InputText
              v-model="condition.condition"
              :placeholder="$t('block.record.fieldConditions.placeholder')"
              class="flex-1"
            />
            <div class="flex items-center gap-1">
              <Checkbox
                :model-value="condition.clearOnHide || false"
                binary
                :input-id="`clearOnHide-${i}`"
                @update:model-value="updateConditionClearOnHide(i, $event)"
              />
              <label :for="`clearOnHide-${i}`" class="text-xs text-muted-color">{{ $t('block.record.fieldConditions.clearOnHide') }}</label>
            </div>
            <Button
              icon="pi pi-trash"
              severity="danger"
              text
              size="small"
              @click="removeCondition(i)"
            />
          </div>
        </div>

        <Button
          :label="$t('general.label.add')"
          icon="pi pi-plus"
          severity="secondary"
          size="small"
          class="self-start"
          @click="addCondition"
        />

        <small class="text-muted-color">
          {{ $t('block.record.fieldConditions.description') }}
        </small>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useModuleStore } from '@/stores/module'

const { t } = useI18n()

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const emit = defineEmits(['update:block'])

const moduleStore = useModuleStore()

// The page's module (base module for this record block)
const selectedModule = computed(() => {
  const moduleID = props.page?.moduleID
  if (!moduleID) return null
  return moduleStore.getByID(moduleID) || null
})

// Reference module (resolved from the reference field's target module)
const referenceModule = ref(null)

// The module used for the field picker: reference module (if set) or the page's module
const fieldPickerModule = computed(() => {
  if (referenceField.value && referenceModule.value) {
    return referenceModule.value
  }
  return selectedModule.value
})

const moduleName = computed(() => selectedModule.value?.name || t('block.record.noModule'))

const fieldLayoutOptions = [
  { value: 'default', text: t('block.record.fieldsLayoutMode.default') },
  { value: 'noWrap', text: t('block.record.fieldsLayoutMode.noWrap') },
  { value: 'wrap', text: t('block.record.fieldsLayoutMode.wrap') },
]

function updateOptions(key, value) {
  emit('update:block', {
    ...props.block,
    options: { ...props.block.options, [key]: value },
  })
}

const inlineEditEnabled = computed({
  get: () => !!props.block.options?.inlineRecordEditEnabled,
  set: v => updateOptions('inlineRecordEditEnabled', v),
})

const horizontalLayout = computed({
  get: () => !!props.block.options?.horizontalFieldLayoutEnabled,
  set: v => updateOptions('horizontalFieldLayoutEnabled', v),
})

const layoutMode = computed({
  get: () => props.block.options?.recordFieldLayoutOption || 'default',
  set: v => {
    updateOptions('recordFieldLayoutOption', v)
    if (v === 'noWrap') {
      updateOptions('horizontalFieldLayoutEnabled', false)
    }
  },
})

// --- Reference field ---

// All single-value Record-kind fields on the page's module (for reference field selector)
const recordSelectorFields = computed(() => {
  if (!selectedModule.value) return []
  return (selectedModule.value.fields || []).filter(f => f.kind === 'Record' && !f.isMulti)
})

const referenceField = computed({
  get: () => props.block.options?.referenceField || null,
  set: v => {
    updateReferenceModule(v)
  },
})

// When selecting a reference field, resolve the target module and set referenceModuleID
function updateReferenceModule(fieldID) {
  if (!fieldID) {
    // Cleared — reset reference field, module, and fields
    updateOptions('referenceField', null)
    updateOptions('referenceModuleID', null)
    referenceModule.value = null
    return
  }

  const field = recordSelectorFields.value.find(f => f.fieldID === fieldID)
  const moduleID = field?.options?.moduleID

  if (moduleID) {
    moduleStore.findByID({ namespaceID: props.namespace.namespaceID, moduleID }).then(mod => {
      referenceModule.value = mod
      updateOptions('referenceField', fieldID)
      updateOptions('referenceModuleID', mod.moduleID)
      // Reset fields when reference module changes
      updateOptions('fields', [])
    })
  } else {
    updateOptions('referenceField', fieldID)
  }
}

// On mount, resolve existing reference field
watch(() => props.block.options?.referenceField, (fieldID) => {
  if (fieldID && selectedModule.value) {
    const field = recordSelectorFields.value.find(f => f.fieldID === fieldID)
    const moduleID = field?.options?.moduleID
    if (moduleID) {
      moduleStore.findByID({ namespaceID: props.namespace.namespaceID, moduleID }).then(mod => {
        referenceModule.value = mod
      })
    }
  }
}, { immediate: true })

// --- Field picker ---

const selectedFieldNames = ref([])

watch(() => props.block.options?.fields, (fields) => {
  if (fields?.length) {
    selectedFieldNames.value = fields.map(f => f.name ?? f)
  } else {
    selectedFieldNames.value = []
  }
}, { immediate: true })

// All fields: regular + system with translated labels
const allModuleFields = computed(() => {
  if (!fieldPickerModule.value) return []
  const regular = fieldPickerModule.value.fields || []
  const system = (fieldPickerModule.value.systemFields?.() || []).map(f => ({
    ...f,
    label: t(`field.system.${f.name}`, f.label || f.name),
    isSystem: true,
  }))
  return [...regular, ...system]
})

function onFieldPickerUpdate(names) {
  selectedFieldNames.value = names
  updateOptions('fields', names)
}

// --- Record display options visibility ---
// Only show record display section when configured fields include Record-kind fields
const isRecordFieldUsedConfigured = computed(() => {
  if (!fieldPickerModule.value) return false
  const configuredFields = props.block.options?.fields || []
  if (configuredFields.length === 0) {
    // No fields configured — check all module fields
    return (fieldPickerModule.value.fields || []).some(f => f.kind === 'Record')
  }
  // Check if any configured field is of kind Record
  const names = configuredFields.map(f => f.name ?? f)
  return (fieldPickerModule.value.fields || [])
    .filter(f => names.includes(f.name))
    .some(f => f.kind === 'Record')
})

// --- Field conditions ---

const fieldConditions = computed(() => {
  return props.block.options?.fieldConditions || []
})

const conditionFieldOptions = computed(() => {
  const fields = props.block.options?.fields || []
  if (!fieldPickerModule.value) return []

  // Use configured fields if available, otherwise all module fields
  const moduleFields = fieldPickerModule.value.fields || []
  const sourceFields = fields.length > 0
    ? moduleFields.filter(f => fields.map(ff => ff.name ?? ff).includes(f.name))
    : moduleFields

  return sourceFields.map(f => ({
    value: f.fieldID || f.name,
    text: f.label || f.name,
  }))
})

function addCondition() {
  const conditions = [...fieldConditions.value, { field: undefined, condition: '', clearOnHide: false }]
  updateOptions('fieldConditions', conditions)
}

function removeCondition(index) {
  const conditions = [...fieldConditions.value]
  conditions.splice(index, 1)
  updateOptions('fieldConditions', conditions)
}

function updateConditionClearOnHide(index, value) {
  const conditions = [...fieldConditions.value]
  conditions[index] = { ...conditions[index], clearOnHide: value }
  updateOptions('fieldConditions', conditions)
}

// --- Display options ---
const displayOptions = [
  { value: 'sameTab', text: t('block.record.openInSameTab') },
  { value: 'newTab', text: t('block.record.openInNewTab') },
  { value: 'modal', text: t('block.record.openInModal') },
]

const inlineRecordEditAllowAddField = computed({
  get: () => !!props.block.options?.inlineRecordEditAllowAddField,
  set: v => updateOptions('inlineRecordEditAllowAddField', v),
})

const clearConditionalFieldsOnHide = computed({
  get: () => !!props.block.options?.clearConditionalFieldsOnHide,
  set: v => updateOptions('clearConditionalFieldsOnHide', v),
})

const recordSelectorDisplayOption = computed({
  get: () => props.block.options?.recordSelectorDisplayOption || 'sameTab',
  set: v => updateOptions('recordSelectorDisplayOption', v),
})

const recordSelectorAddRecordDisplayOption = computed({
  get: () => props.block.options?.recordSelectorAddRecordDisplayOption || 'sameTab',
  set: v => updateOptions('recordSelectorAddRecordDisplayOption', v),
})

const recordSelectorShowAddRecordButton = computed({
  get: () => !!props.block.options?.recordSelectorShowAddRecordButton,
  set: v => updateOptions('recordSelectorShowAddRecordButton', v),
})
</script>
