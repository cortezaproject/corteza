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

          <CInputSwitch v-model="showAddButton" :label="$t('block.recordList.record.hideAddButton')" />

          <CInputSwitch v-model="selectable" :label="$t('block.recordList.selectable')" />

          <CInputSwitch v-model="showImport" :label="$t('block.recordList.record.hideImportButton')" />

          <CInputSwitch v-model="showExport" :label="$t('block.recordList.record.hideExportButton')" />
        </div>
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
  { value: 'sameTab', text: t('block.record.openInSameTab') },
  { value: 'newTab', text: t('block.record.openInNewTab') },
  { value: 'modal', text: t('block.record.openInModal') },
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

const showExport = computed({
  get: () => !props.block.options?.hideExportButton,
  set: v => updateOptions('hideExportButton', !v),
})

// --- Field picker ---

const selectedFieldNames = ref([])

// Initialize selectedFieldNames from block options
watch(() => props.block.options?.fields, (fields) => {
  if (fields?.length) {
    selectedFieldNames.value = fields.map(f => f.name ?? f)
  }
}, { immediate: true })

// Available fields that are not yet selected
const availableFields = computed(() => {
  if (!recordListModule.value) return []
  const selected = new Set(selectedFieldNames.value)
  return (recordListModule.value.fields || []).filter(f => !selected.has(f.name))
})

// Selected fields in order
const selectedFields = computed(() => {
  if (!recordListModule.value) return []
  return selectedFieldNames.value
    .map(name => (recordListModule.value.fields || []).find(f => f.name === name))
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
</script>
