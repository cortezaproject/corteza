<template>
  <div class="flex flex-col gap-5">
    <!-- General -->
    <div class="flex flex-col gap-3">
      <h5 class="text-lg font-semibold text-primary m-0">
        {{ $t('block.recordList.record.generalLabel') }}
      </h5>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <!-- Module selection -->
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.general.module') }}</label>
          <Select
            v-model="moduleID"
            :options="modules"
            option-label="name"
            option-value="moduleID"
            :placeholder="$t('block.recordList.modulePlaceholder')"
            class="w-full"
            filter
          />
        </div>

        <!-- Inline editing (when module selected) -->
        <CInputSwitch v-if="recordListModule" v-model="editable" :label="$t('block.recordList.record.inlineEditorAllow')" />
      </div>
    </div>

    <template v-if="recordListModule">
      <Divider />

      <!-- Fields -->
      <div class="flex flex-col gap-3">
        <h5 class="text-lg font-semibold text-primary m-0">
          {{ $t('block.general.fields') }}
        </h5>
        <small class="text-muted-color">{{ $t('block.recordList.moduleFieldsFootnote') }}</small>

        <PickList
          v-model="fieldPickerModel"
          data-key="name"
          breakpoint="768px"
          :pt="{
            list: { style: 'height: 300px' },
          }"
        >
          <template #option="{ option }">
            {{ option.label || option.name }}
          </template>
        </PickList>
      </div>

      <Divider />

      <!-- Prefilter & Search -->
      <div class="flex flex-col gap-3">
        <h5 class="text-lg font-semibold text-primary m-0">
          {{ $t('block.recordList.record.prefilterLabel') }}
        </h5>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <CInputSwitch v-model="showSearch" :label="$t('block.recordList.record.prefilterHideSearch')" />

          <CInputSwitch v-model="showFiltering" :label="$t('block.recordList.record.filterHide')" />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.recordList.record.prefilterLabel') }}</label>
          <Textarea v-model="prefilter" rows="3" class="w-full" :placeholder="$t('block.recordList.record.prefilterPlaceholder')" />
          <small class="text-muted-color">
            {{ $t('block.recordList.record.prefilterFootnote') }}
          </small>
        </div>
      </div>

      <Divider />

      <!-- Sorting -->
      <div class="flex flex-col gap-3">
        <h5 class="text-lg font-semibold text-primary m-0">
          {{ $t('block.recordList.record.presortLabel') }}
        </h5>

        <CInputSwitch v-model="showSorting" :label="$t('block.recordList.record.presortHideSort')" />

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.recordList.record.presortInputLabel') }}</label>
          <InputText v-model="presort" class="w-full" :placeholder="$t('block.recordList.record.presortPlaceholder')" />
          <small class="text-muted-color">
            {{ $t('block.recordList.record.presortFootnote') }}
          </small>
        </div>
      </div>

      <Divider />

      <!-- Paging -->
      <div class="flex flex-col gap-3">
        <h5 class="text-lg font-semibold text-primary m-0">
          {{ $t('block.recordList.record.pagingLabel') }}
        </h5>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <CInputSwitch v-model="showPaging" :label="$t('block.recordList.record.hidePaging')" />

          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">{{ $t('block.recordList.record.perPage') }}</label>
            <InputNumber v-model="perPage" :min="1" :max="1000" class="w-full" />
          </div>

          <CInputSwitch v-model="showTotalCount" :label="$t('block.recordList.record.showTotalCount')" />
        </div>
      </div>

      <Divider />

      <!-- Records -->
      <div class="flex flex-col gap-3">
        <h5 class="text-lg font-semibold text-primary m-0">
          {{ $t('block.recordList.record.recordsLabel') }}
        </h5>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">{{ $t('block.recordList.record.recordDisplayOptions') }}</label>
            <Select
              v-model="recordDisplayOption"
              :options="recordDisplayOptions"
              option-label="text"
              option-value="value"
              class="w-full"
            />
          </div>

          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">{{ $t('block.recordList.record.recordSelectorDisplayOptions') }}</label>
            <Select
              v-model="recordSelectorDisplayOption"
              :options="recordDisplayOptions"
              option-label="text"
              option-value="value"
              class="w-full"
            />
          </div>

          <CInputSwitch v-model="showAddButton" :label="$t('block.recordList.record.hideAddButton')" />

          <CInputSwitch v-model="selectable" :label="$t('block.recordList.selectable')" />

          <CInputSwitch v-model="showImport" :label="$t('block.recordList.record.hideImportButton')" />

          <CInputSwitch v-model="allowExport" :label="$t('block.recordList.export.allow')" />

          <CInputSwitch v-model="showConfigureFieldsButton" :label="$t('block.recordList.hideConfigureFieldsButton')" />
        </div>
      </div>

      <Divider />

      <!-- Row Action Buttons -->
      <div class="flex flex-col gap-3">
        <h5 class="text-lg font-semibold text-primary m-0">
          {{ $t('block.recordList.record.rowActionButtons') }}
        </h5>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <CInputSwitch v-model="hideRecordViewButton" :label="$t('block.recordList.hideRecordViewButton')" />
          <CInputSwitch v-model="hideRecordEditButton" :label="$t('block.recordList.hideRecordEditButton')" />
          <CInputSwitch v-model="hideRecordCloneButton" :label="$t('block.recordList.hideRecordCloneButton')" />
          <CInputSwitch v-model="hideRecordReminderButton" :label="$t('block.recordList.hideRecordReminderButton')" />
          <CInputSwitch v-model="hideRecordPermissionsButton" :label="$t('block.recordList.hideRecordPermissionsButton')" />
          <CInputSwitch v-model="hideRecordDeleteButton" :label="$t('block.recordList.hideRecordDeleteButton')" />
        </div>
      </div>

      <Divider />

      <!-- Inline Editing -->
      <div class="flex flex-col gap-3">
        <h5 class="text-lg font-semibold text-primary m-0">
          {{ $t('block.recordList.record.inlineEditLabel') }}
        </h5>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <CInputSwitch v-model="inlineRecordEditEnabled" :label="$t('block.recordList.record.inlineRecordEditEnabled')" />
          <CInputSwitch v-model="inlineRecordEditAllowAddField" :label="$t('block.recordList.record.inlineRecordEditAllowAddField')" :disabled="!inlineRecordEditEnabled" />
          <CInputSwitch v-model="bulkRecordEditEnabled" :label="$t('block.recordList.record.bulkRecordEditEnabled')" />
          <CInputSwitch v-model="inlineValueFiltering" :label="$t('block.recordList.record.inlineValueFiltering')" />
          <CInputSwitch v-model="openRecordInEditMode" :label="$t('block.recordList.record.openRecordInEditMode')" />
        </div>
      </div>

      <Divider />

      <!-- Advanced Display -->
      <div class="flex flex-col gap-3">
        <h5 class="text-lg font-semibold text-primary m-0">
          {{ $t('block.recordList.record.advancedDisplay') }}
        </h5>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <CInputSwitch v-model="showDeletedRecordsOption" :label="$t('block.recordList.record.showDeletedRecordsOption')" />
          <CInputSwitch v-model="showRecordPerPageOption" :label="$t('block.recordList.record.showRecordPerPageOption')" />
          <CInputSwitch v-model="fullPageNavigation" :label="$t('block.recordList.record.fullPageNavigation')" />
          <CInputSwitch v-model="enableRecordPageNavigation" :label="$t('block.recordList.record.enableRecordPageNavigation')" />
          <CInputSwitch v-model="showRefresh" :label="$t('block.general.refresh.show')" />
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">{{ $t('block.recordList.record.addRecordDisplayOption') }}</label>
            <Select
              v-model="addRecordDisplayOption"
              :options="recordCreateOptions"
              option-label="text"
              option-value="value"
              class="w-full"
            />
          </div>

          <div v-if="showRefresh" class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">{{ $t('block.general.refreshRate') }}</label>
            <InputNumber v-model="refreshRate" :min="0" class="w-full" />
            <small class="text-muted-color">{{ $t('block.general.refreshRateFootnote') }}</small>
          </div>
        </div>
      </div>

      <Divider />

      <!-- Filter Presets -->
      <div class="flex flex-col gap-3">
        <h5 class="text-lg font-semibold text-primary m-0">
          {{ $t('block.recordList.record.filterPresets') }}
        </h5>

        <CInputSwitch v-model="customFilterPresets" :label="$t('block.recordList.record.enableFilterPresets')" />

        <template v-if="customFilterPresets">
          <div
            v-for="(preset, i) in filterPresets"
            :key="i"
            class="flex flex-col gap-2 p-3 border border-surface rounded-border"
          >
            <div class="flex items-center justify-between">
              <span class="text-sm font-medium">{{ $t('block.recordList.record.preset') }} {{ i + 1 }}</span>
              <Button icon="pi pi-trash" severity="danger" text size="small" @click="removeFilterPreset(i)" />
            </div>
            <InputText
              :model-value="preset.name || ''"
              :placeholder="$t('block.recordList.record.presetName')"
              class="w-full"
              @update:model-value="updateFilterPreset(i, 'name', $event)"
            />
            <Textarea
              :model-value="preset.filter || ''"
              :placeholder="$t('block.recordList.record.presetFilter')"
              rows="2"
              class="w-full"
              @update:model-value="updateFilterPreset(i, 'filter', $event)"
            />
            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">{{ $t('block.recordList.record.presetRoles') }}</label>
              <CInputRole
                :model-value="preset.roles || []"
                multiple
                class="w-full"
                @update:model-value="updateFilterPreset(i, 'roles', $event)"
              />
            </div>
          </div>

          <Button
            :label="$t('general.label.add')"
            icon="pi pi-plus"
            severity="secondary"
            size="small"
            class="self-start"
            @click="addFilterPreset"
          />
        </template>
      </div>

      <Divider />

      <!-- Summaries -->
      <div class="flex flex-col gap-3">
        <h5 class="text-lg font-semibold text-primary m-0">
          {{ $t('block.recordList.record.summaries') }}
        </h5>

        <CInputSwitch v-model="customSummaries" :label="$t('block.recordList.record.enableSummaries')" />

        <template v-if="customSummaries">
          <div
            v-for="(summary, i) in summaries"
            :key="i"
            class="flex flex-col gap-2 p-3 border border-surface rounded-border"
          >
            <div class="flex items-center justify-between">
              <span class="text-sm font-medium">{{ $t('block.recordList.record.summary') }} {{ i + 1 }}</span>
              <Button icon="pi pi-trash" severity="danger" text size="small" @click="removeSummary(i)" />
            </div>
            <InputText
              :model-value="summary.label || ''"
              :placeholder="$t('block.recordList.record.summaryLabel')"
              class="w-full"
              @update:model-value="updateSummary(i, 'label', $event)"
            />
            <div class="grid grid-cols-2 gap-2">
              <Select
                :model-value="summary.field || ''"
                :options="moduleFieldOptions"
                option-label="text"
                option-value="value"
                :placeholder="$t('block.recordList.record.summaryField')"
                class="w-full"
                @update:model-value="updateSummary(i, 'field', $event)"
              />
              <Select
                :model-value="summary.metric || ''"
                :options="summaryMetrics"
                option-label="label"
                option-value="value"
                :placeholder="$t('block.recordList.record.summaryMetric')"
                class="w-full"
                @update:model-value="updateSummary(i, 'metric', $event)"
              />
            </div>
            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">{{ $t('block.recordList.record.summaryRoles') }}</label>
              <CInputRole
                :model-value="summary.roles || []"
                multiple
                class="w-full"
                @update:model-value="updateSummary(i, 'roles', $event)"
              />
            </div>
          </div>

          <Button
            :label="$t('general.label.add')"
            icon="pi pi-plus"
            severity="secondary"
            size="small"
            class="self-start"
            @click="addSummary"
          />
        </template>
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
const modules = computed(() => moduleStore.set || [])

const recordListModule = computed(() => {
  if (!moduleID.value) return null
  return moduleStore.getByID(moduleID.value) || null
})

const recordDisplayOptions = [
  { value: 'sameTab', text: t('block.recordList.record.openInSameTab') },
  { value: 'newTab', text: t('block.recordList.record.openInNewTab') },
  { value: 'modal', text: t('block.recordList.record.openInModal') },
  { value: 'doNothing', text: t('block.recordList.record.doNothing') },
]

const recordCreateOptions = [
  { value: 'sameTab', text: t('block.recordList.record.createInSameTab') },
  { value: 'newTab', text: t('block.recordList.record.createInNewTab') },
  { value: 'modal', text: t('block.recordList.record.createInModal') },
]

// Helper to update block options
function updateOptions(key, value) {
  emit('update:block', {
    ...props.block,
    options: { ...props.block.options, [key]: value },
  })
}

// --- Computed options bindings ---

const moduleID = computed({
  get: () => props.block.options?.moduleID,
  set: v => updateOptions('moduleID', v),
})

const editable = computed({
  get: () => !!props.block.options?.editable,
  set: v => updateOptions('editable', v),
})

const prefilter = computed({
  get: () => props.block.options?.prefilter || '',
  set: v => updateOptions('prefilter', v),
})

const presort = computed({
  get: () => props.block.options?.presort || '',
  set: v => updateOptions('presort', v),
})

const perPage = computed({
  get: () => props.block.options?.perPage ?? 20,
  set: v => updateOptions('perPage', v),
})

const recordDisplayOption = computed({
  get: () => props.block.options?.recordDisplayOption || 'sameTab',
  set: v => updateOptions('recordDisplayOption', v),
})

// Inverted toggles (show = !hide)
const showSearch = computed({
  get: () => !props.block.options?.hideSearch,
  set: v => updateOptions('hideSearch', !v),
})

const showFiltering = computed({
  get: () => !props.block.options?.hideFiltering,
  set: v => updateOptions('hideFiltering', !v),
})

const showSorting = computed({
  get: () => !props.block.options?.hideSorting,
  set: v => updateOptions('hideSorting', !v),
})

const showPaging = computed({
  get: () => !props.block.options?.hidePaging,
  set: v => updateOptions('hidePaging', !v),
})

const showAddButton = computed({
  get: () => !props.block.options?.hideAddButton,
  set: v => updateOptions('hideAddButton', !v),
})

const showTotalCount = computed({
  get: () => !!props.block.options?.showTotalCount,
  set: v => updateOptions('showTotalCount', v),
})

const selectable = computed({
  get: () => props.block.options?.selectable !== false,
  set: v => updateOptions('selectable', v),
})

const showImport = computed({
  get: () => !props.block.options?.hideImportButton,
  set: v => updateOptions('hideImportButton', !v),
})

const allowExport = computed({
  get: () => props.block.options?.allowExport !== false,
  set: v => updateOptions('allowExport', v),
})

const showConfigureFieldsButton = computed({
  get: () => !props.block.options?.hideConfigureFieldsButton,
  set: v => updateOptions('hideConfigureFieldsButton', !v),
})

const recordSelectorDisplayOption = computed({
  get: () => props.block.options?.recordSelectorDisplayOption || 'sameTab',
  set: v => updateOptions('recordSelectorDisplayOption', v),
})

const showRefresh = computed({
  get: () => !!props.block.options?.showRefresh,
  set: v => updateOptions('showRefresh', v),
})

const refreshRate = computed({
  get: () => props.block.options?.refreshRate ?? 0,
  set: v => updateOptions('refreshRate', v),
})

// --- Field picker ---

const selectedFieldNames = ref([])

// Initialize selectedFieldNames from block options
watch(() => props.block.options?.fields, (fields) => {
  if (fields?.length) {
    selectedFieldNames.value = fields.map(f => f.name ?? f)
  }
}, { immediate: true })

// All fields: regular + system with translated labels
const allModuleFields = computed(() => {
  if (!recordListModule.value) return []
  const regular = recordListModule.value.fields || []
  const system = (recordListModule.value.systemFields?.() || []).map(f => ({
    ...f,
    label: t(`field.system.${f.name}`, f.label || f.name),
    isSystem: true,
  }))
  return [...regular, ...system]
})

// Available fields that are not yet selected
const availableFields = computed(() => {
  const selected = new Set(selectedFieldNames.value)
  return allModuleFields.value.filter(f => !selected.has(f.name))
})

// Selected fields in order
const selectedFields = computed(() => {
  return selectedFieldNames.value
    .map(name => allModuleFields.value.find(f => f.name === name))
    .filter(Boolean)
})

// PickList model: [available, selected]
const fieldPickerModel = computed({
  get: () => [availableFields.value, selectedFields.value],
  set: (val) => {
    const [, selected] = val
    selectedFieldNames.value = selected.map(f => f.name)
    updateOptions('fields', selected.map(f => f.name))
  },
})

const summaryMetrics = [
  { value: 'sum', label: t('block.recordList.summaries.metrics.sum.label') },
  { value: 'min', label: t('block.recordList.summaries.metrics.min.label') },
  { value: 'max', label: t('block.recordList.summaries.metrics.max.label') },
  { value: 'avg', label: t('block.recordList.summaries.metrics.avg.label') },
  { value: 'emptyCount', label: t('block.recordList.summaries.metrics.emptyCount.label') },
  { value: 'notEmptyCount', label: t('block.recordList.summaries.metrics.notEmptyCount.label') },
  { value: 'uniqueCount', label: t('block.recordList.summaries.metrics.uniqueCount.label') },
]

const moduleFieldOptions = computed(() => {
  return allModuleFields.value.map(f => ({
    value: f.name,
    text: f.label || f.name,
  }))
})

// Row action buttons — use full keys matching the type definition
const hideRecordViewButton = computed({
  get: () => !!props.block.options?.hideRecordViewButton,
  set: v => updateOptions('hideRecordViewButton', v),
})
const hideRecordEditButton = computed({
  get: () => !!props.block.options?.hideRecordEditButton,
  set: v => updateOptions('hideRecordEditButton', v),
})
const hideRecordCloneButton = computed({
  get: () => !!props.block.options?.hideRecordCloneButton,
  set: v => updateOptions('hideRecordCloneButton', v),
})
const hideRecordReminderButton = computed({
  get: () => !!props.block.options?.hideRecordReminderButton,
  set: v => updateOptions('hideRecordReminderButton', v),
})
const hideRecordPermissionsButton = computed({
  get: () => !!props.block.options?.hideRecordPermissionsButton,
  set: v => updateOptions('hideRecordPermissionsButton', v),
})
const hideRecordDeleteButton = computed({
  get: () => !!props.block.options?.hideRecordDeleteButton,
  set: v => updateOptions('hideRecordDeleteButton', v),
})

// Inline editing
const inlineRecordEditEnabled = computed({
  get: () => !!props.block.options?.inlineRecordEditEnabled,
  set: v => updateOptions('inlineRecordEditEnabled', v),
})
const inlineRecordEditAllowAddField = computed({
  get: () => !!props.block.options?.inlineRecordEditAllowAddField,
  set: v => updateOptions('inlineRecordEditAllowAddField', v),
})
const bulkRecordEditEnabled = computed({
  get: () => !!props.block.options?.bulkRecordEditEnabled,
  set: v => updateOptions('bulkRecordEditEnabled', v),
})
const inlineValueFiltering = computed({
  get: () => !!props.block.options?.inlineValueFiltering,
  set: v => updateOptions('inlineValueFiltering', v),
})
const openRecordInEditMode = computed({
  get: () => !!props.block.options?.openRecordInEditMode,
  set: v => updateOptions('openRecordInEditMode', v),
})

// Advanced display
const showDeletedRecordsOption = computed({
  get: () => !!props.block.options?.showDeletedRecordsOption,
  set: v => updateOptions('showDeletedRecordsOption', v),
})
const showRecordPerPageOption = computed({
  get: () => !!props.block.options?.showRecordPerPageOption,
  set: v => updateOptions('showRecordPerPageOption', v),
})
const fullPageNavigation = computed({
  get: () => !!props.block.options?.fullPageNavigation,
  set: v => updateOptions('fullPageNavigation', v),
})
const enableRecordPageNavigation = computed({
  get: () => !!props.block.options?.enableRecordPageNavigation,
  set: v => updateOptions('enableRecordPageNavigation', v),
})
const addRecordDisplayOption = computed({
  get: () => props.block.options?.addRecordDisplayOption || 'sameTab',
  set: v => updateOptions('addRecordDisplayOption', v),
})

// Filter presets
const customFilterPresets = computed({
  get: () => !!props.block.options?.customFilterPresets,
  set: v => updateOptions('customFilterPresets', v),
})
const filterPresets = computed(() => props.block.options?.filterPresets || [])

function addFilterPreset() {
  updateOptions('filterPresets', [...filterPresets.value, { name: '', filter: '', roles: [] }])
}
function removeFilterPreset(i) {
  const updated = [...filterPresets.value]
  updated.splice(i, 1)
  updateOptions('filterPresets', updated)
}
function updateFilterPreset(i, key, value) {
  const updated = [...filterPresets.value]
  updated[i] = { ...updated[i], [key]: value }
  updateOptions('filterPresets', updated)
}

// Summaries
const customSummaries = computed({
  get: () => !!props.block.options?.customSummaries,
  set: v => updateOptions('customSummaries', v),
})
const summaries = computed(() => props.block.options?.summaries || [])

function addSummary() {
  updateOptions('summaries', [...summaries.value, { label: '', field: '', metric: '', roles: [] }])
}
function removeSummary(i) {
  const updated = [...summaries.value]
  updated.splice(i, 1)
  updateOptions('summaries', updated)
}
function updateSummary(i, key, value) {
  const updated = [...summaries.value]
  updated[i] = { ...updated[i], [key]: value }
  updateOptions('summaries', updated)
}
</script>
