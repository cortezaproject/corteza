<template>
  <div class="flex flex-col gap-3">
    <Fieldset :legend="$t('block.recordList.record.generalLabel')">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <CFormGroup :label="$t('block.general.module')">
          <InputText :model-value="moduleName" disabled class="w-full" />
        </CFormGroup>

        <CFormGroup
          :label="$t('block.record.referenceRecordField')"
          :description="$t('block.record.referenceRecordFieldDescription')"
        >
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
        </CFormGroup>

        <CFormGroup :label="$t('block.record.fieldsLayoutMode.label')">
          <Select
            v-model="layoutMode"
            :options="fieldLayoutOptions"
            option-label="text"
            option-value="value"
            class="w-full"
          />
        </CFormGroup>

        <CInputToggleCard
          v-if="layoutMode !== 'noWrap'"
          v-model="horizontalLayout"
          :label="$t('block.record.horizontalFormLayout')"
          :description="$t('block.record.horizontalFormLayoutDescription')"
        />
      </div>
    </Fieldset>

    <template v-if="fieldPickerModule">
      <Divider />

      <Fieldset :legend="$t('block.general.fields')">
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
      </Fieldset>

      <Divider />

      <CInputToggleCard
        v-model="inlineEditEnabled"
        :label="$t('block.record.inlineEdit.enabled')"
        :description="$t('block.record.inlineEdit.description')"
      />

      <CInputToggleCard
        v-model="inlineCopyEnabled"
        :label="$t('block.record.inlineCopy.enabled')"
        :description="$t('block.record.inlineCopy.description')"
      />

      <Divider />

      <Fieldset :legend="$t('block.record.recordDisplay.label')">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <CFormGroup :label="$t('block.record.recordSelectorDisplayOptions')">
            <Select
              v-model="recordSelectorDisplayOption"
              :options="displayOptions"
              option-label="text"
              option-value="value"
              class="w-full"
            />
          </CFormGroup>

          <CInputToggleCard
            v-model="recordSelectorShowAddRecordButton"
            :label="$t('block.record.recordSelectorCanAddRecord')"
            :description="$t('block.record.recordSelectorCanAddRecordDescription')"
          />

          <CFormGroup :label="$t('block.record.recordSelectorAddRecordDisplayOption')">
            <Select
              v-model="recordSelectorAddRecordDisplayOption"
              :options="displayOptions"
              option-label="text"
              option-value="value"
              class="w-full"
              :disabled="!recordSelectorShowAddRecordButton"
            />
          </CFormGroup>
        </div>
      </Fieldset>

      <Divider />

      <Fieldset>
        <template #legend>
          <div class="flex items-center gap-2">
            <span>{{ $t('block.record.fieldConditions.label') }}</span>
            <i
              class="pi pi-exclamation-triangle text-orange-500 text-sm"
              v-tooltip="$t('block.record.fieldConditions.tooltip.performance')"
            />
          </div>
        </template>

        <div class="flex flex-col gap-3">
          <CInputToggleCard
            v-model="clearConditionalFieldsOnHide"
            :label="$t('block.record.fieldConditions.clearAllOnHide')"
            :description="$t('block.record.fieldConditions.clearAllOnHideDescription')"
          />

          <div class="flex flex-col gap-2">
            <div v-for="(condition, i) in fieldConditions" :key="i" class="flex items-center gap-2">
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
                <label :for="`clearOnHide-${i}`" class="text-xs text-muted-color">
                  {{ $t('block.record.fieldConditions.clearOnHide') }}
                </label>
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
      </Fieldset>
    </template>
  </div>
</template>

<script setup>
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useModuleStore } from '@planetcrust/human-vue'

const { t } = useI18n()

const props = defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const block = inject('blockDraft')
const $ComposeAPI = inject('$ComposeAPI')

function patchOptions(patch) {
  if (!block.value.options) block.value.options = {}
  Object.assign(block.value.options, patch)
}

const moduleStore = useModuleStore()

// The page's module (base module for this record block)
const selectedModule = computed(() => {
  const moduleID = props.page?.moduleID
  if (!moduleID) return null
  return moduleStore.getByID(moduleID) || null
})

// Reference module (resolved from the reference field's target module)
const referenceModule = ref(null)

// Local ref so fieldPickerModule reacts immediately on selection, without waiting
// for the parent to process the update:block emit and push props back down.
const localReferenceFieldID = ref(block.value.options?.referenceField || null)

// The module used for the field picker: reference module (if set) or the page's module
const fieldPickerModule = computed(() => {
  if (localReferenceFieldID.value && referenceModule.value) {
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
  patchOptions({ [key]: value })
}

const inlineEditEnabled = computed({
  get: () => !!block.value.options?.inlineRecordEditEnabled,
  set: v => updateOptions('inlineRecordEditEnabled', v),
})

const inlineCopyEnabled = computed({
  get: () => !!block.value.options?.inlineRecordCopyEnabled,
  set: v => updateOptions('inlineRecordCopyEnabled', v),
})

const horizontalLayout = computed({
  get: () => !!block.value.options?.horizontalFieldLayoutEnabled,
  set: v => updateOptions('horizontalFieldLayoutEnabled', v),
})

const layoutMode = computed({
  get: () => block.value.options?.recordFieldLayoutOption || 'default',
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
  get: () => localReferenceFieldID.value,
  set: v => updateReferenceModule(v),
})

function resolveReferenceModule(moduleID) {
  if (!moduleID || moduleID === '0') return
  // Use synchronous cache lookup first to avoid an async tick before fieldPickerModule updates
  const cached = moduleStore.getByID(moduleID)
  if (cached) {
    referenceModule.value = cached
    return
  }
  moduleStore
    .findByID({ namespaceID: props.namespace.namespaceID, moduleID })
    .then(mod => {
      referenceModule.value = mod
    })
    .catch(e => {
      console.warn('Failed to resolve reference module', moduleID, e)
    })
}

// When selecting a reference field, resolve the target module and persist to block options
function updateReferenceModule(fieldID) {
  if (!fieldID) {
    localReferenceFieldID.value = null
    referenceModule.value = null
    patchOptions({ referenceField: null, referenceModuleID: null, fields: [] })
    return
  }

  localReferenceFieldID.value = fieldID

  const field = recordSelectorFields.value.find(f => f.fieldID === fieldID)
  const moduleID = field?.options?.moduleID

  if (moduleID && moduleID !== '0') {
    const cached = moduleStore.getByID(moduleID)
    if (cached) {
      referenceModule.value = cached
      patchOptions({ referenceField: fieldID, referenceModuleID: moduleID, fields: [] })
    } else {
      patchOptions({ referenceField: fieldID, fields: [] })
      moduleStore
        .findByID({ namespaceID: props.namespace.namespaceID, moduleID })
        .then(mod => {
          referenceModule.value = mod
          patchOptions({ referenceField: fieldID, referenceModuleID: mod.moduleID, fields: [] })
        })
        .catch(e => {
          console.warn('Failed to resolve reference module', moduleID, e)
        })
    }
  } else {
    patchOptions({ referenceField: fieldID, fields: [] })
  }
}

// On mount (and when selectedModule becomes available), restore reference module from saved options
watch(
  [() => block.value.options?.referenceField, selectedModule],
  ([fieldID]) => {
    localReferenceFieldID.value = fieldID || null
    if (!fieldID || !selectedModule.value) {
      if (!fieldID) referenceModule.value = null
      return
    }
    const field = recordSelectorFields.value.find(f => f.fieldID === fieldID)
    resolveReferenceModule(field?.options?.moduleID)
  },
  { immediate: true },
)

// --- Field picker ---

const selectedFieldNames = ref([])

watch(
  () => block.value.options?.fields,
  fields => {
    if (fields?.length) {
      selectedFieldNames.value = fields.map(f => f.name ?? f)
    } else {
      selectedFieldNames.value = []
    }
  },
  { immediate: true },
)

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

// --- Field conditions ---

const fieldConditions = computed(() => {
  return block.value.options?.fieldConditions || []
})

const conditionFieldOptions = computed(() => {
  const fields = block.value.options?.fields || []
  if (!fieldPickerModule.value) return []

  // Use configured fields if available, otherwise all module fields
  const moduleFields = fieldPickerModule.value.fields || []
  const sourceFields =
    fields.length > 0
      ? moduleFields.filter(f => fields.map(ff => ff.name ?? ff).includes(f.name))
      : moduleFields

  return sourceFields.map(f => ({
    value: f.fieldID || f.name,
    text: f.label || f.name,
  }))
})

function addCondition() {
  const conditions = [
    ...fieldConditions.value,
    { field: undefined, condition: '', clearOnHide: false },
  ]
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

const clearConditionalFieldsOnHide = computed({
  get: () => !!block.value.options?.clearConditionalFieldsOnHide,
  set: v => updateOptions('clearConditionalFieldsOnHide', v),
})

const recordSelectorDisplayOption = computed({
  get: () => block.value.options?.recordSelectorDisplayOption || 'sameTab',
  set: v => updateOptions('recordSelectorDisplayOption', v),
})

const recordSelectorAddRecordDisplayOption = computed({
  get: () => block.value.options?.recordSelectorAddRecordDisplayOption || 'sameTab',
  set: v => updateOptions('recordSelectorAddRecordDisplayOption', v),
})

const recordSelectorShowAddRecordButton = computed({
  get: () => !!block.value.options?.recordSelectorShowAddRecordButton,
  set: v => updateOptions('recordSelectorShowAddRecordButton', v),
})
</script>
