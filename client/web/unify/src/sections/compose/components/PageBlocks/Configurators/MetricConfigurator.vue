<template>
  <div class="flex flex-col gap-3">
    <CFormGroup :label="$t('block.metric.edit.listTitle')">
      <template #actions>
        <Button
          :label="$t('block.metric.add')"
          icon="pi pi-plus"
          severity="secondary"
          size="small"
          @click="addMetric"
        />
      </template>

      <CFormList
        v-model="metrics"
        draggable
        :empty-message="$t('block.metric.edit.empty')"
        :columns="[{ width: '3rem' }, { width: '1fr' }]"
        @change="expandedMetric = -1"
        @reorder="expandedMetric = -1"
      >
        <template #row="{ item: metric, index: i }">
          <Button
            :icon="expandedMetric === i ? 'pi pi-chevron-down' : 'pi pi-chevron-right'"
            text
            rounded
            size="small"
            @click="toggleMetric(i)"
          />
          <span
            class="text-sm font-semibold truncate cursor-pointer select-none w-full"
            @click="toggleMetric(i)"
          >
            {{ metric.label || $t('block.metric.defaultMetricLabel') }}
          </span>
        </template>

        <template #extra="{ item: metric, index: i }">
          <div v-if="expandedMetric === i" class="flex flex-col gap-3 border-t border-surface pt-3">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
              <CFormGroup :label="$t('block.metric.edit.labelLabel')">
                <InputText
                  v-model="metric.label"
                  :placeholder="$t('block.metric.edit.labelPlaceholder')"
                  class="w-full"
                />
              </CFormGroup>

              <CFormGroup :label="$t('block.metric.edit.moduleLabel')">
                <CInputModule
                  :model-value="metric.moduleID"
                  :namespaceID="namespace?.namespaceID"
                  :placeholder="$t('block.metric.edit.modulePlaceholder')"
                  @update:model-value="onModuleChange(metric, $event)"
                />
              </CFormGroup>

              <CFormGroup :label="$t('block.metric.edit.metricFieldLabel')">
                <CInputModuleField
                  :model-value="metric.metricField"
                  :module-i-d="metric.moduleID"
                  :kinds="['Number']"
                  :extra-options="countOption"
                  :placeholder="$t('block.metric.edit.metricFieldSelect')"
                  :disabled="!metric.moduleID"
                  @update:model-value="onMetricFieldChange(metric, $event)"
                />
              </CFormGroup>

              <CFormGroup :label="$t('block.metric.edit.metricAggregateLabel')">
                <Select
                  v-model="metric.operation"
                  :options="aggregationOperations"
                  option-label="label"
                  option-value="operation"
                  class="w-full"
                  :placeholder="$t('block.metric.edit.metricSelectAggregate')"
                  :disabled="metric.metricField === 'count'"
                />
              </CFormGroup>

              <CFormGroup
                :label="$t('block.metric.edit.transformFunctionLabel')"
                :description="$t('block.metric.edit.transformFunctionDescription')"
                class="md:col-span-2"
              >
                <InputText v-model="metric.transformFx" class="w-full" placeholder="v" />
              </CFormGroup>

              <CFormGroup :label="$t('block.metric.edit.numberFormat')">
                <InputText v-model="metric.numberFormat" class="w-full" placeholder="0,0.00" />
              </CFormGroup>

              <CFormGroup :label="$t('block.metric.edit.prefixLabel')">
                <InputText v-model="metric.prefix" class="w-full" placeholder="$" />
              </CFormGroup>

              <CFormGroup :label="$t('block.metric.edit.suffixLabel')">
                <InputText v-model="metric.suffix" class="w-full" placeholder="USD/mo" />
              </CFormGroup>

              <CFormGroup
                :label="$t('block.metric.edit.filterLabel')"
                :description="$t('block.metric.edit.filterFootnote')"
                class="md:col-span-2"
              >
                <CInputExpression
                  :ref="el => (filterInputs[i] = el)"
                  v-model="metric.filter"
                  dialect="ql"
                  :scope="scope"
                  :query-fields="moduleFields(metric.moduleID)"
                  :placeholder="$t('block.metric.edit.filterPlaceholder')"
                />
                <CExpressionHint :scope="scope" @insert="filterInputs[i]?.insert($event)" />
              </CFormGroup>
            </div>

            <Divider />

            <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
              <CFormGroup :label="$t('block.metric.editStyle.color')">
                <div class="flex items-center gap-2">
                  <CInputColorPicker
                    :model-value="metric.valueStyle?.color || ''"
                    :default-value="defaultTextColor"
                    show-text
                    :empty-label="$t('block.metric.editStyle.default')"
                    @update:model-value="onStyleChange(metric, 'color', $event)"
                  />
                  <Button
                    icon="pi pi-undo"
                    severity="secondary"
                    text
                    rounded
                    size="small"
                    :title="$t('block.metric.editStyle.resetToDefault')"
                    @click="onStyleChange(metric, 'color', '')"
                  />
                </div>
              </CFormGroup>

              <CFormGroup :label="$t('block.metric.editStyle.backgroundColor')">
                <div class="flex items-center gap-2">
                  <CInputColorPicker
                    :model-value="metric.valueStyle?.backgroundColor || ''"
                    show-text
                    :empty-label="$t('block.metric.editStyle.default')"
                    @update:model-value="onStyleChange(metric, 'backgroundColor', $event)"
                  />
                  <Button
                    icon="pi pi-undo"
                    severity="secondary"
                    text
                    rounded
                    size="small"
                    :title="$t('block.metric.editStyle.resetToDefault')"
                    @click="onStyleChange(metric, 'backgroundColor', '')"
                  />
                </div>
              </CFormGroup>
            </div>

            <Divider />

            <div class="flex flex-col gap-3">
              <div class="flex items-center gap-2">
                <Checkbox
                  :model-value="metric.comparison?.enabled || false"
                  binary
                  :input-id="`comparison-enabled-${i}`"
                  @update:model-value="onComparisonChange(metric, 'enabled', $event)"
                />
                <label :for="`comparison-enabled-${i}`" class="text-sm">
                  {{ $t('block.metric.comparison.enabled') }}
                </label>
              </div>

              <template v-if="metric.comparison?.enabled">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                  <CFormGroup :label="$t('block.metric.comparison.period')">
                    <Select
                      :model-value="metric.comparison?.period || 'month'"
                      :options="comparisonPeriods"
                      option-label="label"
                      option-value="value"
                      class="w-full"
                      @update:model-value="onComparisonChange(metric, 'period', $event)"
                    />
                  </CFormGroup>

                  <CFormGroup :label="$t('block.metric.comparison.customFilter')">
                    <CInputExpression
                      :ref="el => (comparisonFilterInputs[i] = el)"
                      :model-value="metric.comparison?.customFilter || ''"
                      dialect="ql"
                      :scope="scope"
                      :query-fields="moduleFields(metric.moduleID)"
                      :placeholder="$t('block.metric.edit.filterPlaceholder')"
                      @update:model-value="onComparisonChange(metric, 'customFilter', $event)"
                    />
                    <CExpressionHint
                      :scope="scope"
                      @insert="comparisonFilterInputs[i]?.insert($event)"
                    />
                  </CFormGroup>
                </div>
              </template>
            </div>

            <Divider />

            <div class="flex flex-col gap-3">
              <div class="flex items-center gap-2">
                <Checkbox
                  :model-value="metric.drillDown?.enabled || false"
                  binary
                  :input-id="`drilldown-enabled-${i}`"
                  @update:model-value="onDrillDownChange(metric, 'enabled', $event)"
                />
                <label :for="`drilldown-enabled-${i}`" class="text-sm">
                  {{ $t('block.metric.drillDown.enabled') }}
                </label>
              </div>

              <template v-if="metric.drillDown?.enabled">
                <CFormGroup
                  :label="$t('block.metric.drillDown.blockID')"
                  :description="$t('block.metric.drillDown.blockIDFootnote')"
                >
                  <InputText
                    :model-value="metric.drillDown?.blockID || ''"
                    :placeholder="$t('block.metric.drillDown.blockIDPlaceholder')"
                    class="w-full"
                    @update:model-value="onDrillDownChange(metric, 'blockID', $event)"
                  />
                </CFormGroup>
              </template>
            </div>
          </div>
        </template>
      </CFormList>
    </CFormGroup>
  </div>
</template>

<script setup>
import { ref, computed, inject, onMounted } from 'vue'
import { components } from '@planetcrust/human-vue'
import { useModuleStore } from '@planetcrust/human-vue'
import { useI18n } from 'vue-i18n'
import { useExpressionScope } from '@/sections/compose/composables/useExpressionScope'

const { CInputColorPicker } = components

const { t } = useI18n()
const moduleStore = useModuleStore()

const props = defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const filterInputs = ref([])
const comparisonFilterInputs = ref([])
const { scope } = useExpressionScope({ page: computed(() => props.page) })

// Fields of the module this row queries — bare identifiers in its filter.
function moduleFields(moduleID) {
  return (moduleID && moduleStore.getByID(moduleID)?.fields) || []
}

const block = inject('blockDraft')

const expandedMetric = ref(0)

function toggleMetric(i) {
  expandedMetric.value = expandedMetric.value === i ? -1 : i
}

// Resolve the theme's text color for the swatch preview
const defaultTextColor = ref('')
onMounted(() => {
  const style = getComputedStyle(document.documentElement)
  const textColor = style.getPropertyValue('--p-text-color').trim()
  if (textColor) {
    defaultTextColor.value = textColor
  }
})

const aggregationOperations = computed(() => [
  { label: t('block.metric.edit.operationSum'), operation: 'sum' },
  { label: t('block.metric.edit.operationMax'), operation: 'max' },
  { label: t('block.metric.edit.operationMin'), operation: 'min' },
  { label: t('block.metric.edit.operationAvg'), operation: 'avg' },
])

const comparisonPeriods = computed(() => [
  { label: t('block.metric.comparison.periodDay'), value: 'day' },
  { label: t('block.metric.comparison.periodWeek'), value: 'week' },
  { label: t('block.metric.comparison.periodMonth'), value: 'month' },
  { label: t('block.metric.comparison.periodQuarter'), value: 'quarter' },
  { label: t('block.metric.comparison.periodYear'), value: 'year' },
])

/**
 * Returns the metric fields for a given module.
 * Includes 'Count' as a special option plus all Number fields.
 */
// Counting rows is an aggregate, not a field, but it is picked in the same box.
const countOption = computed(() => [{ name: 'count', label: t('block.metric.edit.countField') }])

function updateOptions(key, value) {
  if (!block.value.options) block.value.options = {}
  block.value.options[key] = value
}

const metrics = computed({
  get: () => block.value.options?.metrics || [],
  set: val => updateOptions('metrics', val),
})

function addMetric() {
  // Use the block's makeMetric() if available, otherwise create a default
  let newMetric
  if (block.value.makeMetric) {
    newMetric = block.value.makeMetric()
  } else {
    newMetric = {
      label: '',
      moduleID: '',
      dimensionField: '',
      dateFormat: '',
      filter: '',
      bucketSize: '',
      metricField: '',
      operation: '',
      numberFormat: '',
      prefix: '',
      suffix: '',
      transformFx: '',
      valueStyle: {
        backgroundColor: undefined,
        color: undefined,
      },
      drillDown: {
        enabled: false,
        blockID: '',
        recordListOptions: { fields: [] },
      },
      comparison: {
        enabled: false,
        type: 'period',
        period: 'month',
        customFilter: '',
      },
    }
  }

  const m = [...metrics.value, newMetric]
  updateOptions('metrics', m)
  expandedMetric.value = m.length - 1
}

function onModuleChange(metric, moduleID) {
  metric.moduleID = moduleID
  // Reset field and operation when module changes
  metric.metricField = ''
  metric.operation = ''
}

function onMetricFieldChange(metric, field) {
  metric.metricField = field

  if (field === 'count') {
    // Count doesn't need an aggregation operation
    metric.operation = ''
  } else if (!metric.operation) {
    // Default to sum for number fields
    metric.operation = 'sum'
  }
}

function onStyleChange(metric, key, value) {
  if (!metric.valueStyle) {
    metric.valueStyle = {
      backgroundColor: undefined,
      color: undefined,
    }
  }
  metric.valueStyle[key] = value
}

function onComparisonChange(metric, key, value) {
  if (!metric.comparison) {
    metric.comparison = {
      enabled: false,
      type: 'period',
      period: 'month',
      customFilter: '',
    }
  }
  metric.comparison[key] = value
}

function onDrillDownChange(metric, key, value) {
  if (!metric.drillDown) {
    metric.drillDown = {
      enabled: false,
      blockID: '',
    }
  }
  metric.drillDown[key] = value
}
</script>
