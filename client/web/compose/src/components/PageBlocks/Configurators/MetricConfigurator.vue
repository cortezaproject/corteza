<template>
  <div class="flex flex-col gap-3">
    <!-- Metrics list -->
    <div class="flex flex-col gap-3">
      <h5 class="text-lg font-semibold text-primary m-0">
        {{ $t('block.metric.edit.tabTitle') }}
      </h5>

      <div
        v-for="(metric, i) in metrics"
        :key="i"
        class="flex flex-col gap-2 p-3 border rounded-lg"
      >
        <div class="flex items-center gap-2 justify-between">
          <span class="text-sm font-semibold">{{ metric.label || $t('block.metric.defaultMetricLabel') }}</span>
          <div class="flex gap-1">
            <Button
              :icon="expandedMetric === i ? 'pi pi-chevron-up' : 'pi pi-chevron-down'"
              text
              size="small"
              @click="expandedMetric = expandedMetric === i ? -1 : i"
            />
            <Button
              icon="pi pi-trash"
              severity="danger"
              text
              size="small"
              @click="removeMetric(i)"
            />
          </div>
        </div>

        <template v-if="expandedMetric === i">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
            <!-- Label -->
            <div class="flex flex-col gap-1">
              <label class="text-primary font-medium text-sm">{{ $t('block.metric.edit.labelLabel') }}</label>
              <InputText v-model="metric.label" :placeholder="$t('block.metric.edit.labelPlaceholder')" class="w-full" />
            </div>

            <!-- Module (per-metric) -->
            <div class="flex flex-col gap-1">
              <label class="text-primary font-medium text-sm">{{ $t('block.metric.edit.moduleLabel') }}</label>
              <Select
                :model-value="metric.moduleID"
                :options="modules"
                option-label="name"
                option-value="moduleID"
                :placeholder="$t('block.metric.edit.modulePlaceholder')"
                class="w-full"
                filter
                @update:model-value="onModuleChange(metric, $event)"
              />
            </div>

            <!-- Metric Field -->
            <div class="flex flex-col gap-1">
              <label class="text-primary font-medium text-sm">{{ $t('block.metric.edit.metricFieldLabel') }}</label>
              <Select
                :model-value="metric.metricField"
                :options="getMetricFields(metric.moduleID)"
                option-label="label"
                option-value="name"
                class="w-full"
                :placeholder="$t('block.metric.edit.metricFieldSelect')"
                :disabled="!metric.moduleID"
                @update:model-value="onMetricFieldChange(metric, $event)"
              />
            </div>

            <!-- Aggregation operation -->
            <div class="flex flex-col gap-1">
              <label class="text-primary font-medium text-sm">{{ $t('block.metric.edit.metricAggregateLabel') }}</label>
              <Select
                v-model="metric.operation"
                :options="aggregationOperations"
                option-label="label"
                option-value="operation"
                class="w-full"
                :placeholder="$t('block.metric.edit.metricSelectAggregate')"
                :disabled="metric.metricField === 'count'"
              />
            </div>

            <!-- Transform function -->
            <div class="flex flex-col gap-1 md:col-span-2">
              <label class="text-primary font-medium text-sm">{{ $t('block.metric.edit.transformFunctionLabel') }}</label>
              <InputText
                v-model="metric.transformFx"
                class="w-full"
                placeholder="v"
              />
              <small class="text-muted-color">{{ $t('block.metric.edit.transformFunctionDescription') }}</small>
            </div>

            <!-- Number format -->
            <div class="flex flex-col gap-1">
              <label class="text-primary font-medium text-sm">{{ $t('block.metric.edit.numberFormat') }}</label>
              <InputText
                v-model="metric.numberFormat"
                class="w-full"
                placeholder="0,0.00"
              />
            </div>

            <!-- Prefix -->
            <div class="flex flex-col gap-1">
              <label class="text-primary font-medium text-sm">{{ $t('block.metric.edit.prefixLabel') }}</label>
              <InputText v-model="metric.prefix" class="w-full" placeholder="$" />
            </div>

            <!-- Suffix -->
            <div class="flex flex-col gap-1">
              <label class="text-primary font-medium text-sm">{{ $t('block.metric.edit.suffixLabel') }}</label>
              <InputText v-model="metric.suffix" class="w-full" placeholder="USD/mo" />
            </div>

            <!-- Filter -->
            <div class="flex flex-col gap-1 md:col-span-2">
              <label class="text-primary font-medium text-sm">{{ $t('block.metric.edit.filterLabel') }}</label>
              <InputText v-model="metric.filter" class="w-full" placeholder="field1 = 1 AND field2 > 0" />
            </div>
          </div>

          <!-- Value Style -->
          <Divider />
          <h6 class="text-base font-semibold text-primary m-0">
            {{ $t('block.metric.editStyle.valueLabel') }}
          </h6>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
            <!-- Text color -->
            <div class="flex flex-col gap-1">
              <label class="text-primary font-medium text-sm">{{ $t('block.metric.editStyle.color') }}</label>
              <div class="flex items-center gap-2">
                <input
                  type="color"
                  :value="metric.valueStyle?.color || '#000000'"
                  class="w-8 h-8 border rounded cursor-pointer"
                  @input="onStyleChange(metric, 'color', $event.target.value)"
                >
                <InputText
                  :model-value="metric.valueStyle?.color || '#000000'"
                  class="w-full"
                  @update:model-value="onStyleChange(metric, 'color', $event)"
                />
              </div>
            </div>

            <!-- Background color -->
            <div class="flex flex-col gap-1">
              <label class="text-primary font-medium text-sm">{{ $t('block.metric.editStyle.backgroundColor') }}</label>
              <div class="flex items-center gap-2">
                <input
                  type="color"
                  :value="metric.valueStyle?.backgroundColor || '#FFFFFF'"
                  class="w-8 h-8 border rounded cursor-pointer"
                  @input="onStyleChange(metric, 'backgroundColor', $event.target.value)"
                >
                <InputText
                  :model-value="metric.valueStyle?.backgroundColor || '#FFFFFF00'"
                  class="w-full"
                  @update:model-value="onStyleChange(metric, 'backgroundColor', $event)"
                />
              </div>
            </div>

            <!-- Font size -->
            <div class="flex flex-col gap-1">
              <label class="text-primary font-medium text-sm">{{ $t('block.metric.editStyle.fontSize') }}</label>
              <InputNumber
                :model-value="metric.valueStyle?.fontSize ? Number(metric.valueStyle.fontSize) : undefined"
                placeholder="16"
                :min="1"
                :step="1"
                suffix=" px"
                class="w-full"
                @update:model-value="onStyleChange(metric, 'fontSize', $event)"
              />
            </div>
          </div>
        </template>
      </div>

      <Button
        :label="$t('block.metric.add')"
        icon="pi pi-plus"
        severity="secondary"
        size="small"
        class="self-start"
        @click="addMetric"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useModuleStore } from '@/stores/module'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const moduleStore = useModuleStore()

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const emit = defineEmits(['update:block'])

const expandedMetric = ref(0)

const modules = computed(() => moduleStore.set || [])

const aggregationOperations = computed(() => [
  { label: t('block.metric.edit.operationSum'), operation: 'sum' },
  { label: t('block.metric.edit.operationMax'), operation: 'max' },
  { label: t('block.metric.edit.operationMin'), operation: 'min' },
  { label: t('block.metric.edit.operationAvg'), operation: 'avg' },
])

/**
 * Returns the metric fields for a given module.
 * Includes 'Count' as a special option plus all Number fields.
 */
function getMetricFields (moduleID) {
  if (!moduleID) return []

  const mod = moduleStore.getByID(moduleID)
  if (!mod) return []

  const numberFields = mod.fields
    .filter(f => f.kind === 'Number')
    .map(f => ({ name: f.name, label: f.label || f.name }))
    .sort((a, b) => a.label.localeCompare(b.label))

  return [
    { name: 'count', label: 'Count' },
    ...numberFields,
  ]
}

function updateOptions (key, value) {
  emit('update:block', {
    ...props.block,
    options: { ...props.block.options, [key]: value },
  })
}

const metrics = computed(() => {
  return props.block.options?.metrics || []
})

function addMetric () {
  // Use the block's makeMetric() if available, otherwise create a default
  let newMetric
  if (props.block.makeMetric) {
    newMetric = props.block.makeMetric()
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
        backgroundColor: '#FFFFFF00',
        color: '#000000',
        fontSize: undefined,
      },
      drillDown: {
        enabled: false,
        blockID: '',
        recordListOptions: { fields: [] },
      },
    }
  }

  const m = [...metrics.value, newMetric]
  updateOptions('metrics', m)
  expandedMetric.value = m.length - 1
}

function removeMetric (index) {
  const m = [...metrics.value]
  m.splice(index, 1)
  updateOptions('metrics', m)

  if (expandedMetric.value >= m.length) {
    expandedMetric.value = m.length - 1
  }
}

function onModuleChange (metric, moduleID) {
  metric.moduleID = moduleID
  // Reset field and operation when module changes
  metric.metricField = ''
  metric.operation = ''
}

function onMetricFieldChange (metric, field) {
  metric.metricField = field

  if (field === 'count') {
    // Count doesn't need an aggregation operation
    metric.operation = ''
  } else if (!metric.operation) {
    // Default to sum for number fields
    metric.operation = 'sum'
  }
}

function onStyleChange (metric, key, value) {
  if (!metric.valueStyle) {
    metric.valueStyle = {
      backgroundColor: '#FFFFFF00',
      color: '#000000',
      fontSize: undefined,
    }
  }
  metric.valueStyle[key] = value
}
</script>
