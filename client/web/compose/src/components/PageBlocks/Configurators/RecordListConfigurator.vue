<template>
  <Tabs v-model:value="activeTab">
    <TabList>
      <Tab value="general">{{ $t('block.recordList.record.generalLabel') }}</Tab>
      <Tab value="automation">{{ $t('block.automation.tabs.workflows') }}</Tab>
    </TabList>
    <TabPanels class="px-0 pb-0">
      <TabPanel value="general">
        <div class="flex flex-col gap-3">
          <!-- General -->
          <Panel :header="$t('block.recordList.record.generalLabel')" toggleable>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
              <div class="flex flex-col gap-1">
                <label class="text-primary font-medium text-sm">
                  {{ $t('block.general.module') }}
                </label>
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

              <CInputToggleCard
                v-if="recordListModule"
                v-model="editable"
                v-tooltip.top="$t('block.recordList.record.inlineEditorAllow')"
                :label="$t('block.recordList.record.inlineEditorAllow')"
              />

              <!-- refField (parent link) — only on record pages -->
              <div
                v-if="recordListModule && onRecordPage"
                class="flex flex-col gap-1 md:col-span-2"
              >
                <label class="text-primary font-medium text-sm">
                  {{ $t('block.recordList.refField.label') }}
                </label>
                <Select
                  v-model="refField"
                  :options="parentFields"
                  option-label="label"
                  option-value="name"
                  :placeholder="$t('general.label.none')"
                  class="w-full"
                  show-clear
                />
                <small class="text-muted-color">
                  {{ $t('block.recordList.refField.footnote') }}
                </small>
              </div>
            </div>
          </Panel>

          <template v-if="recordListModule">
            <!-- Fields -->
            <Panel :header="$t('block.general.fields')" toggleable>
              <div class="flex flex-col gap-3">
                <small class="text-muted-color">
                  {{ $t('block.recordList.moduleFieldsFootnote') }}
                </small>
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
            </Panel>

            <!-- Inline Editor -->
            <Panel v-if="editable" :header="$t('block.recordList.record.inlineEditor')" toggleable>
              <div class="flex flex-col gap-1">
                <label class="text-primary font-medium text-sm">
                  {{ $t('block.recordList.editFields') }}
                </label>
                <CFieldPicker
                  :all-fields="editableFieldSubset"
                  :model-value="selectedEditFieldNames"
                  :available-label="$t('field.selector.available')"
                  :selected-label="$t('field.selector.selected')"
                  :select-all-label="$t('field.selector.selectAll')"
                  :unselect-all-label="$t('field.selector.unselectAll')"
                  :search-placeholder="$t('field.selector.search')"
                  :no-items-label="$t('field.no-items-found')"
                  @update:model-value="onEditFieldPickerUpdate"
                />
              </div>
            </Panel>

            <!-- Prefilter & Search -->
            <Panel :header="$t('block.recordList.record.prefilterLabel')" toggleable>
              <div class="flex flex-col gap-3">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                  <CInputToggleCard
                    v-model="showSearch"
                    :label="$t('block.recordList.record.prefilterHideSearch')"
                  />
                  <CInputToggleCard
                    v-model="showFiltering"
                    :label="$t('block.recordList.record.filterHide')"
                  />
                </div>

                <div v-if="showSearch" class="flex flex-col gap-1">
                  <label class="text-primary font-medium text-sm">
                    {{ $t('block.recordList.record.searchableFields') }}
                  </label>
                  <CFieldPicker
                    :all-fields="queryableFields"
                    :model-value="selectedSearchableFieldNames"
                    :available-label="$t('field.selector.available')"
                    :selected-label="$t('field.selector.selected')"
                    :select-all-label="$t('field.selector.selectAll')"
                    :unselect-all-label="$t('field.selector.unselectAll')"
                    :search-placeholder="$t('field.selector.search')"
                    :no-items-label="$t('field.no-items-found')"
                    @update:model-value="onSearchableFieldPickerUpdate"
                  />
                  <small class="text-muted-color">
                    {{ $t('block.recordList.record.searchableFieldsFootnote') }}
                  </small>
                </div>

                <div class="flex flex-col gap-1">
                  <label class="text-primary font-medium text-sm">
                    {{ $t('block.recordList.record.prefilterLabel') }}
                  </label>
                  <Textarea
                    v-model="prefilter"
                    rows="3"
                    class="w-full"
                    :placeholder="$t('block.recordList.record.prefilterPlaceholder')"
                  />
                  <small class="text-muted-color">
                    {{ $t('block.recordList.record.prefilterFootnote') }}
                  </small>
                </div>

                <!-- Filter Presets -->
                <div class="flex flex-col gap-3">
                  <CInputToggleCard
                    v-model="customFilterPresets"
                    :label="$t('block.recordList.record.enableFilterPresets')"
                  />
                  <template v-if="customFilterPresets">
                    <div
                      v-for="(preset, i) in filterPresets"
                      :key="i"
                      class="flex flex-col gap-2 p-3 border border-surface rounded-border"
                    >
                      <div class="flex items-center justify-between gap-2">
                        <InputText
                          :model-value="preset.name || ''"
                          :placeholder="$t('block.recordList.record.presetName')"
                          class="flex-1"
                          @update:model-value="updateFilterPreset(i, 'name', $event)"
                        />
                        <RecordListFilter
                          :model-value="normalizePresetFilter(preset.filter)"
                          :module="recordListModule"
                          :show-reset="true"
                          @update:model-value="updateFilterPreset(i, 'filter', $event)"
                        />
                        <Button
                          icon="pi pi-trash"
                          severity="danger"
                          text
                          size="small"
                          @click="removeFilterPreset(i)"
                        />
                      </div>
                      <div class="flex flex-col gap-1">
                        <label class="text-sm text-muted-color">
                          {{ $t('block.recordList.record.presetRoles') }}
                        </label>
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
              </div>
            </Panel>

            <!-- Sorting -->
            <Panel :header="$t('block.recordList.record.presortLabel')" toggleable>
              <div class="flex flex-col gap-3">
                <CInputToggleCard
                  v-model="showSorting"
                  :label="$t('block.recordList.record.presortHideSort')"
                />
                <div class="flex flex-col gap-2">
                  <label class="text-primary font-medium text-sm">
                    {{ $t('block.recordList.record.presortInputLabel') }}
                  </label>

                  <div
                    v-for="(item, i) in presortItems"
                    :key="i"
                    class="flex items-center gap-2 flex-nowrap"
                  >
                    <Select
                      :model-value="item.field"
                      :options="sortableFields"
                      option-label="label"
                      option-value="name"
                      :placeholder="$t('block.recordList.record.presortPlaceholder')"
                      class="flex-1 min-w-0"
                      filter
                      show-clear
                      @update:model-value="updatePresortItem(i, 'field', $event)"
                    />
                    <Select
                      :model-value="item.descending"
                      :options="sortDirections"
                      option-label="label"
                      option-value="value"
                      class="shrink-0 w-auto"
                      :disabled="!item.field"
                      @update:model-value="updatePresortItem(i, 'descending', $event)"
                    />
                    <Button
                      icon="pi pi-trash"
                      severity="danger"
                      text
                      size="small"
                      class="shrink-0"
                      :disabled="presortItems.length === 1"
                      @click="removePresortItem(i)"
                    />
                  </div>

                  <Button
                    :label="$t('general.label.add')"
                    icon="pi pi-plus"
                    severity="secondary"
                    size="small"
                    class="self-start"
                    @click="addPresortItem"
                  />
                </div>
              </div>
            </Panel>

            <!-- Paging -->
            <Panel :header="$t('block.recordList.record.pagingLabel')" toggleable>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                <CInputToggleCard
                  v-model="showPaging"
                  :label="$t('block.recordList.record.hidePaging')"
                />
                <CInputToggleCard
                  v-model="fullPageNavigation"
                  v-tooltip.top="$t('block.recordList.record.fullPageNavigationHint')"
                  :label="$t('block.recordList.record.fullPageNavigation')"
                />
                <div class="flex flex-col gap-1">
                  <label
                    v-tooltip.top="$t('block.recordList.record.perPageHint')"
                    class="text-primary font-medium text-sm"
                  >
                    {{ $t('block.recordList.record.perPage') }}
                    <i class="pi pi-info-circle text-muted-color text-xs ml-1" />
                  </label>
                  <InputNumber v-model="perPage" :min="1" :max="1000" class="w-full" />
                </div>
                <CInputToggleCard
                  v-model="showRecordPerPageOption"
                  :label="$t('block.recordList.record.showRecordPerPageOption')"
                />
                <CInputToggleCard
                  v-model="showTotalCount"
                  :label="$t('block.recordList.record.showTotalCount')"
                />
              </div>
            </Panel>

            <!-- Summaries -->
            <Panel :header="$t('block.recordList.record.summaries')" toggleable>
              <div class="flex flex-col gap-3">
                <CInputToggleCard
                  v-model="customSummaries"
                  :label="$t('block.recordList.record.enableSummaries')"
                />
                <template v-if="customSummaries">
                  <div
                    v-for="(summary, i) in summaries"
                    :key="i"
                    class="flex flex-col gap-2 p-3 border border-surface rounded-border"
                  >
                    <div class="flex items-center justify-between">
                      <span class="text-sm font-medium">
                        {{ $t('block.recordList.record.summary') }} {{ i + 1 }}
                      </span>
                      <Button
                        icon="pi pi-trash"
                        severity="danger"
                        text
                        size="small"
                        @click="removeSummary(i)"
                      />
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
                      <label class="text-sm text-muted-color">
                        {{ $t('block.recordList.record.summaryRoles') }}
                      </label>
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
            </Panel>

            <!-- Records -->
            <Panel :header="$t('block.recordList.record.recordsLabel')" toggleable>
              <div class="flex flex-col gap-3">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                  <div class="flex flex-col gap-1">
                    <label class="text-primary font-medium text-sm">
                      {{ $t('block.recordList.record.recordDisplayOptions') }}
                    </label>
                    <Select
                      v-model="recordDisplayOption"
                      :options="recordDisplayOptions"
                      option-label="text"
                      option-value="value"
                      class="w-full"
                    />
                  </div>
                  <div class="flex flex-col gap-1">
                    <label class="text-primary font-medium text-sm">
                      {{ $t('block.recordList.record.recordSelectorDisplayOptions') }}
                    </label>
                    <Select
                      v-model="recordSelectorDisplayOption"
                      :options="recordDisplayOptions"
                      option-label="text"
                      option-value="value"
                      class="w-full"
                    />
                  </div>
                  <CInputToggleCard
                    v-model="showAddButton"
                    :label="$t('block.recordList.record.hideAddButton')"
                  />
                  <div class="flex flex-col gap-1">
                    <label class="text-primary font-medium text-sm">
                      {{ $t('block.recordList.record.addRecordDisplayOption') }}
                    </label>
                    <Select
                      v-model="addRecordDisplayOption"
                      :options="recordCreateOptions"
                      option-label="text"
                      option-value="value"
                      class="w-full"
                      :disabled="!showAddButton"
                    />
                  </div>
                  <CInputToggleCard
                    v-model="openRecordInEditMode"
                    :label="$t('block.recordList.record.openRecordInEditMode')"
                  />
                  <CInputToggleCard v-model="selectable" :label="$t('block.recordList.selectable')" />
                  <CInputToggleCard
                    v-model="showImport"
                    :label="$t('block.recordList.record.hideImportButton')"
                  />
                  <CInputToggleCard
                    v-model="allowExport"
                    :label="$t('block.recordList.export.allow')"
                  />
                  <CInputToggleCard
                    v-model="showConfigureFieldsButton"
                    :label="$t('block.recordList.hideConfigureFieldsButton')"
                  />
                  <CInputToggleCard
                    v-model="inlineRecordEditEnabled"
                    :label="$t('block.recordList.record.inlineRecordEditEnabled')"
                  />
                  <CInputToggleCard
                    v-if="inlineRecordEditEnabled"
                    v-model="inlineRecordEditAllowAddField"
                    :label="$t('block.recordList.record.inlineRecordEditAllowAddField')"
                  />
                  <CInputToggleCard
                    v-model="bulkRecordEditEnabled"
                    :label="$t('block.recordList.record.bulkRecordEditEnabled')"
                  />
                  <CInputToggleCard
                    v-model="inlineValueFiltering"
                    :label="$t('block.recordList.record.inlineValueFiltering')"
                  />
                  <CInputToggleCard
                    v-model="enableRecordPageNavigation"
                    v-tooltip.top="$t('block.recordList.record.enableRecordPageNavigationHint')"
                    :label="$t('block.recordList.record.enableRecordPageNavigation')"
                  />
                  <CInputToggleCard
                    v-model="showDeletedRecordsOption"
                    :label="$t('block.recordList.record.showDeletedRecordsOption')"
                  />
                </div>

                <!-- Inline Edit Fields Picker -->
                <div v-if="inlineRecordEditEnabled" class="flex flex-col gap-1">
                  <label class="text-primary font-medium text-sm">
                    {{ $t('block.recordList.record.inlineEditFields') }}
                  </label>
                  <CFieldPicker
                    :all-fields="editableFieldSubset"
                    :model-value="selectedInlineEditFieldNames"
                    :available-label="$t('field.selector.available')"
                    :selected-label="$t('field.selector.selected')"
                    :select-all-label="$t('field.selector.selectAll')"
                    :unselect-all-label="$t('field.selector.unselectAll')"
                    :search-placeholder="$t('field.selector.search')"
                    :no-items-label="$t('field.no-items-found')"
                    @update:model-value="onInlineEditFieldPickerUpdate"
                  />
                </div>

                <!-- Row Action Buttons -->
                <div class="flex flex-col gap-2">
                  <label class="text-primary font-medium text-sm">
                    {{ $t('block.recordList.record.rowActionButtons') }}
                  </label>
                  <div class="flex flex-col gap-2">
                    <div class="flex items-center gap-2">
                      <Checkbox
                        v-model="hideRecordViewButton"
                        binary
                        input-id="hideRecordViewButton"
                      />
                      <label for="hideRecordViewButton" class="text-sm">
                        {{ $t('block.recordList.hideRecordViewButton') }}
                      </label>
                    </div>
                    <div class="flex items-center gap-2">
                      <Checkbox
                        v-model="hideRecordEditButton"
                        binary
                        input-id="hideRecordEditButton"
                      />
                      <label for="hideRecordEditButton" class="text-sm">
                        {{ $t('block.recordList.hideRecordEditButton') }}
                      </label>
                    </div>
                    <div class="flex items-center gap-2">
                      <Checkbox
                        v-model="hideRecordCloneButton"
                        binary
                        input-id="hideRecordCloneButton"
                      />
                      <label for="hideRecordCloneButton" class="text-sm">
                        {{ $t('block.recordList.hideRecordCloneButton') }}
                      </label>
                    </div>
                    <div class="flex items-center gap-2">
                      <Checkbox
                        v-model="hideRecordReminderButton"
                        binary
                        input-id="hideRecordReminderButton"
                      />
                      <label for="hideRecordReminderButton" class="text-sm">
                        {{ $t('block.recordList.hideRecordReminderButton') }}
                      </label>
                    </div>
                    <div class="flex items-center gap-2">
                      <Checkbox
                        v-model="hideRecordPermissionsButton"
                        binary
                        input-id="hideRecordPermissionsButton"
                      />
                      <label for="hideRecordPermissionsButton" class="text-sm">
                        {{ $t('block.recordList.hideRecordPermissionsButton') }}
                      </label>
                    </div>
                    <div class="flex items-center gap-2">
                      <Checkbox
                        v-model="hideRecordDeleteButton"
                        binary
                        input-id="hideRecordDeleteButton"
                      />
                      <label for="hideRecordDeleteButton" class="text-sm">
                        {{ $t('block.recordList.hideRecordDeleteButton') }}
                      </label>
                    </div>
                  </div>
                </div>
              </div>
            </Panel>
          </template>
        </div>
      </TabPanel>

      <TabPanel value="automation">
        <AutomationButtonsEditor
          :buttons="selectionButtons"
          @update:buttons="onSelectionButtonsUpdate"
        />
      </TabPanel>
    </TabPanels>
  </Tabs>
</template>

<script setup>
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useModuleStore } from '@/stores/module'
import RecordListFilter from '@/components/Common/RecordListFilter.vue'
import AutomationButtonsEditor from '../Shared/AutomationButtonsEditor.vue'

const { t } = useI18n()

const props = defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const block = inject('blockDraft')

const moduleStore = useModuleStore()
const modules = computed(() => moduleStore.set || [])

const activeTab = ref('general')

const recordListModule = computed(() => {
  if (!moduleID.value) return null
  return moduleStore.getByID(moduleID.value) || null
})

// The parent module — the module of the current record page (if any)
const parentModule = computed(() => {
  const pageModuleID = props.page?.moduleID
  if (!pageModuleID || pageModuleID === '0') return null
  return moduleStore.getByID(pageModuleID) || null
})

const onRecordPage = computed(() => !!parentModule.value)

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
  if (!block.value.options) block.value.options = {}
  block.value.options[key] = value
}

// --- Computed options bindings ---

const moduleID = computed({
  get: () => block.value.options?.moduleID,
  set: v => updateOptions('moduleID', v),
})

const editable = computed({
  get: () => !!block.value.options?.editable,
  set: v => updateOptions('editable', v),
})

// Field subset available for inline editing: selected display fields or all module fields
const editableFieldSubset = computed(() => {
  if (!recordListModule.value) return []
  const selected = block.value.options?.fields || []
  if (selected.length) {
    return allModuleFields.value.filter(f => selected.some(s => (s.name ?? s) === f.name))
  }
  return allModuleFields.value
})

const prefilter = computed({
  get: () => block.value.options?.prefilter || '',
  set: v => updateOptions('prefilter', v),
})

const perPage = computed({
  get: () => block.value.options?.perPage ?? 20,
  set: v => updateOptions('perPage', v),
})

const recordDisplayOption = computed({
  get: () => block.value.options?.recordDisplayOption || 'sameTab',
  set: v => updateOptions('recordDisplayOption', v),
})

// Inverted toggles (show = !hide)
const showSearch = computed({
  get: () => !block.value.options?.hideSearch,
  set: v => updateOptions('hideSearch', !v),
})

const showFiltering = computed({
  get: () => !block.value.options?.hideFiltering,
  set: v => updateOptions('hideFiltering', !v),
})

const showSorting = computed({
  get: () => !block.value.options?.hideSorting,
  set: v => updateOptions('hideSorting', !v),
})

const showPaging = computed({
  get: () => !block.value.options?.hidePaging,
  set: v => updateOptions('hidePaging', !v),
})

const showAddButton = computed({
  get: () => !block.value.options?.hideAddButton,
  set: v => updateOptions('hideAddButton', !v),
})

const showTotalCount = computed({
  get: () => block.value.options?.showTotalCount !== false,
  set: v => updateOptions('showTotalCount', v),
})

const selectable = computed({
  get: () => block.value.options?.selectable !== false,
  set: v => updateOptions('selectable', v),
})

const showImport = computed({
  get: () => !block.value.options?.hideImportButton,
  set: v => updateOptions('hideImportButton', !v),
})

const allowExport = computed({
  get: () => block.value.options?.allowExport !== false,
  set: v => updateOptions('allowExport', v),
})

const showConfigureFieldsButton = computed({
  get: () => !block.value.options?.hideConfigureFieldsButton,
  set: v => updateOptions('hideConfigureFieldsButton', !v),
})

const recordSelectorDisplayOption = computed({
  get: () => block.value.options?.recordSelectorDisplayOption || 'sameTab',
  set: v => updateOptions('recordSelectorDisplayOption', v),
})

// --- refField ---
const refField = computed({
  get: () => block.value.options?.refField || null,
  set: v => updateOptions('refField', v || undefined),
})

const parentFields = computed(() => {
  if (!recordListModule.value || !parentModule.value) return []
  const parentID = parentModule.value.moduleID
  return (recordListModule.value.fields || [])
    .filter(f => f.kind === 'Record' && f.options?.moduleID === parentID)
    .map(f => ({ name: f.name, label: f.label || f.name }))
})

// --- Field picker ---

const selectedFieldNames = ref([])
const selectedEditFieldNames = ref([])
const selectedSearchableFieldNames = ref([])
const selectedInlineEditFieldNames = ref([])

// Initialize field name refs from block options
watch(
  () => block.value.options?.fields,
  fields => {
    if (fields?.length) {
      selectedFieldNames.value = fields.map(f => f.name ?? f)
    }
  },
  { immediate: true },
)

watch(
  () => block.value.options?.editFields,
  fields => {
    if (fields?.length) {
      selectedEditFieldNames.value = fields.map(f => f.name ?? f)
    }
  },
  { immediate: true },
)

watch(
  () => block.value.options?.searchableFields,
  fields => {
    if (fields?.length) {
      selectedSearchableFieldNames.value = fields.map(f => f.name ?? f)
    } else {
      selectedSearchableFieldNames.value = []
    }
  },
  { immediate: true },
)

watch(
  () => block.value.options?.inlineEditFields,
  fields => {
    if (fields?.length) {
      selectedInlineEditFieldNames.value = fields.map(f => f.name ?? f)
    } else {
      selectedInlineEditFieldNames.value = []
    }
  },
  { immediate: true },
)

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

const queryableFields = computed(() => {
  if (!allModuleFields.value) return []
  return allModuleFields.value.filter(f => f.isQueryable)
})

const sortableFields = computed(() => {
  if (!allModuleFields.value) return []
  return allModuleFields.value
    .filter(f => !f.isMulti)
    .map(f => ({ name: f.name, label: `${f.label || f.name} (${f.name})` }))
})

function onFieldPickerUpdate(names) {
  selectedFieldNames.value = names
  updateOptions('fields', names)
  // Keep editFields restricted to still-selected display fields
  if (block.value.options?.editFields?.length) {
    const filtered = block.value.options.editFields.filter(ef =>
      names.some(n => (ef.name ?? ef) === n),
    )
    updateOptions('editFields', filtered)
  }
  // Keep inlineEditFields restricted to still-selected display fields
  if (block.value.options?.inlineEditFields?.length) {
    const filtered = block.value.options.inlineEditFields.filter(ef =>
      names.some(n => (ef.name ?? ef) === n),
    )
    updateOptions('inlineEditFields', filtered)
  }
}

function onEditFieldPickerUpdate(names) {
  selectedEditFieldNames.value = names
  updateOptions('editFields', names)
}

function onSearchableFieldPickerUpdate(names) {
  selectedSearchableFieldNames.value = names
  updateOptions('searchableFields', names)
}

function onInlineEditFieldPickerUpdate(names) {
  selectedInlineEditFieldNames.value = names
  updateOptions('inlineEditFields', names)
}

// When module changes: reset fields + editable + refField
watch(
  () => block.value.options?.moduleID,
  (newID, oldID) => {
    if (!oldID || newID === oldID) return
    updateOptions('fields', [])
    updateOptions('editable', false)
    updateOptions('refField', undefined)
    updateOptions('inlineEditFields', [])
  },
)

// Auto-pick refField when entering editable mode (first Record-kind field pointing to parent module)
watch(editable, value => {
  if (!value || !recordListModule.value || !parentModule.value) return
  if (block.value.options?.refField) return
  const match = recordListModule.value.fields?.find(
    f => f.kind === 'Record' && f.options?.moduleID === parentModule.value.moduleID,
  )
  if (match) updateOptions('refField', match.name)
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
  get: () => !!block.value.options?.hideRecordViewButton,
  set: v => updateOptions('hideRecordViewButton', v),
})
const hideRecordEditButton = computed({
  get: () => !!block.value.options?.hideRecordEditButton,
  set: v => updateOptions('hideRecordEditButton', v),
})
const hideRecordCloneButton = computed({
  get: () => !!block.value.options?.hideRecordCloneButton,
  set: v => updateOptions('hideRecordCloneButton', v),
})
const hideRecordReminderButton = computed({
  get: () => !!block.value.options?.hideRecordReminderButton,
  set: v => updateOptions('hideRecordReminderButton', v),
})
const hideRecordPermissionsButton = computed({
  get: () => !!block.value.options?.hideRecordPermissionsButton,
  set: v => updateOptions('hideRecordPermissionsButton', v),
})
const hideRecordDeleteButton = computed({
  get: () => !!block.value.options?.hideRecordDeleteButton,
  set: v => updateOptions('hideRecordDeleteButton', v),
})

// Inline editing
const inlineRecordEditEnabled = computed({
  get: () => !!block.value.options?.inlineRecordEditEnabled,
  set: v => updateOptions('inlineRecordEditEnabled', v),
})
const inlineRecordEditAllowAddField = computed({
  get: () => !!block.value.options?.inlineRecordEditAllowAddField,
  set: v => updateOptions('inlineRecordEditAllowAddField', v),
})
const bulkRecordEditEnabled = computed({
  get: () => block.value.options?.bulkRecordEditEnabled !== false,
  set: v => updateOptions('bulkRecordEditEnabled', v),
})
const inlineValueFiltering = computed({
  get: () => !!block.value.options?.inlineValueFiltering,
  set: v => updateOptions('inlineValueFiltering', v),
})
const openRecordInEditMode = computed({
  get: () => !!block.value.options?.openRecordInEditMode,
  set: v => updateOptions('openRecordInEditMode', v),
})

// Advanced display
const showDeletedRecordsOption = computed({
  get: () => !!block.value.options?.showDeletedRecordsOption,
  set: v => updateOptions('showDeletedRecordsOption', v),
})
const showRecordPerPageOption = computed({
  get: () => !!block.value.options?.showRecordPerPageOption,
  set: v => updateOptions('showRecordPerPageOption', v),
})
const fullPageNavigation = computed({
  get: () => !!block.value.options?.fullPageNavigation,
  set: v => updateOptions('fullPageNavigation', v),
})
const enableRecordPageNavigation = computed({
  get: () => block.value.options?.enableRecordPageNavigation !== false,
  set: v => updateOptions('enableRecordPageNavigation', v),
})
const addRecordDisplayOption = computed({
  get: () => block.value.options?.addRecordDisplayOption || 'sameTab',
  set: v => updateOptions('addRecordDisplayOption', v),
})

// --- Filter presets (structured) ---
const customFilterPresets = computed({
  get: () => !!block.value.options?.customFilterPresets,
  set: v => updateOptions('customFilterPresets', v),
})
const filterPresets = computed(() => block.value.options?.filterPresets || [])

// Runtime expects filter as Array of groups; older configs may have stored a string.
// Accept both and always hand Arrays to RecordListFilter.
function normalizePresetFilter(value) {
  if (Array.isArray(value)) return value
  return []
}

function addFilterPreset() {
  updateOptions('filterPresets', [...filterPresets.value, { name: '', filter: [], roles: [] }])
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

// --- Summaries ---
const customSummaries = computed({
  get: () => !!block.value.options?.customSummaries,
  set: v => updateOptions('customSummaries', v),
})
const summaries = computed(() => block.value.options?.summaries || [])

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

// --- Presort (structured) ---
// Stored as comma-separated string with optional ` DESC` suffix per field
// (same format the runtime and API accept).

const presortItems = ref([{ field: null, descending: false }])

function parsePresort(raw) {
  if (!raw) return [{ field: null, descending: false }]
  const parts = String(raw)
    .split(',')
    .map(s => s.trim())
    .filter(Boolean)
  if (!parts.length) return [{ field: null, descending: false }]
  return parts.map(p => {
    const [name, dir] = p.split(/\s+/)
    return { field: name || null, descending: (dir || '').toUpperCase() === 'DESC' }
  })
}

function serializePresort(items) {
  return items
    .filter(it => it.field)
    .map(it => (it.descending ? `${it.field} DESC` : it.field))
    .join(',')
}

watch(
  () => block.value.options?.presort,
  raw => {
    presortItems.value = parsePresort(raw)
  },
  { immediate: true },
)

function updatePresortItem(i, key, value) {
  const updated = presortItems.value.map((it, idx) => (idx === i ? { ...it, [key]: value } : it))
  presortItems.value = updated
  updateOptions('presort', serializePresort(updated))
}

function addPresortItem() {
  presortItems.value = [...presortItems.value, { field: null, descending: false }]
}

function removePresortItem(i) {
  const updated = [...presortItems.value]
  updated.splice(i, 1)
  if (!updated.length) updated.push({ field: null, descending: false })
  presortItems.value = updated
  updateOptions('presort', serializePresort(updated))
}

const sortDirections = computed(() => [
  { value: false, label: t('block.recordList.record.ascending') },
  { value: true, label: t('block.recordList.record.descending') },
])

// --- Automation buttons (selection-level) ---
const selectionButtons = computed(() => block.value.options?.selectionButtons || [])

function onSelectionButtonsUpdate(next) {
  updateOptions('selectionButtons', next)
}
</script>
