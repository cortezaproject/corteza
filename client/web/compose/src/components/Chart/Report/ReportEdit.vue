<template>
  <div>
    <!-- Configure source module -->
    <div class="px-3">
      <h5 class="mb-3">
        {{ $t('chart.edit.module.title') }}
      </h5>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.module.label') }}
          </label>
          <Select
            v-model="moduleID"
            :options="modules"
            option-label="name"
            option-value="moduleID"
            :placeholder="$t('chart.edit.module.placeholder')"
            class="w-full"
            filter
          />
        </div>

        <div v-if="module" class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.filter.preset') }}
          </label>
          <Select
            v-model="reportFilter"
            :options="predefinedFilters"
            option-label="text"
            option-value="value"
            :placeholder="$t('chart.edit.filter.noFilter')"
            class="w-full"
            show-clear
          />
        </div>

        <!-- Configure report filters -->
        <div v-if="module" class="col-span-1 lg:col-span-2">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.filter.label') }}
          </label>
          <Textarea
            v-model="reportFilter"
            :placeholder="$t('chart.edit.filter.placeholder')"
            class="w-full"
            rows="2"
          />
          <small class="text-muted-color">{{ $t('chart.edit.filter.footnote') }}</small>
        </div>
      </div>
    </div>

    <Divider v-if="module" />

    <!-- Configure report dimensions -->
    <div v-if="module" class="px-3">
      <div v-for="(d, i) in dimensions" :key="i">
        <h5 class="mb-3">
          {{ $t('chart.edit.dimension.label') }}
        </h5>

        <template v-if="usesDimensionsField">
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
            <div class="flex flex-col gap-1">
              <label class="text-primary font-medium text-sm">
                {{ $t('chart.edit.dimension.fieldLabel') }}
              </label>
              <Select
                v-model="d.field"
                :options="dimensionFields"
                option-label="text"
                option-value="value"
                :placeholder="$t('chart.edit.dimension.fieldPlaceholder')"
                class="w-full"
                filter
                @change="e => onDimFieldChange(e.value, d)"
              />
            </div>

            <div class="flex flex-col gap-1">
              <label class="text-primary font-medium text-sm">
                {{ $t('chart.edit.dimension.function.label') }}
              </label>
              <Select
                v-model="d.modifier"
                :options="dimensionModifiers"
                option-label="text"
                option-value="value"
                :disabled="!d.field || !isTemporalField(d.field)"
                :placeholder="$t('chart.edit.dimension.function.placeholder')"
                class="w-full"
              />
            </div>
          </div>

          <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 mt-4">
            <div class="flex flex-col gap-1">
              <label class="text-primary font-medium text-sm">
                {{ $t('chart.edit.dimension.defaultValueLabel') }}
              </label>
              <InputText v-model="d.default" class="w-full" />
              <small class="text-muted-color">
                {{ $t('chart.edit.dimension.defaultValueFootnote') }}
              </small>
            </div>

            <div class="flex flex-col gap-1">
              <label class="text-primary font-medium text-sm">
                {{ $t('chart.edit.dimension.options.label') }}
              </label>
              <div class="flex items-center gap-2">
                <Checkbox v-model="d.skipMissing" :binary="true" input-id="skipMissing" />
                <label for="skipMissing">{{ $t('chart.edit.dimension.skipMissingValues') }}</label>
              </div>
              <slot
                name="dimension-options-options"
                :dimension="d"
                :is-temporal="isTemporalField(d.field)"
              />
            </div>
          </div>
        </template>

        <slot name="dimension-options" :index="i" :dimension="d" :field="getField(d)" />
      </div>
    </div>

    <Divider v-if="module" />

    <!-- Configure report metrics -->
    <div v-if="module" class="px-3">
      <div class="flex items-center mb-3">
        <h5 class="m-0">
          {{ $t('chart.edit.metric.title') }}
        </h5>
        <Button
          v-if="canAddMetric"
          :label="'+ ' + $t('chart.edit.metric.add')"
          text
          size="small"
          class="ml-2"
          @click="addMetric"
        />
      </div>

      <div v-for="(m, i) in metrics" :key="i" class="border border-surface rounded p-3 mb-3">
        <div v-if="metrics.length > 1" class="flex items-center mb-3">
          <h6 class="m-0">{{ $t('chart.edit.metric.label') }} {{ i + 1 }}</h6>
          <Button
            icon="pi pi-trash"
            text
            severity="danger"
            size="small"
            class="ml-auto"
            @click="removeMetric(i)"
          />
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">
              {{ $t('chart.edit.metric.fieldLabel') }}
            </label>
            <Select
              v-model="m.field"
              :options="metricFields"
              option-label="text"
              option-value="value"
              :placeholder="$t('chart.edit.metric.fieldPlaceholder')"
              class="w-full"
              filter
              @change="e => onMetricFieldChange(e.value, m)"
            />
          </div>

          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">
              {{ $t('chart.edit.metric.function.label') }}
            </label>
            <Select
              v-model="m.aggregate"
              :options="metricAggregates"
              option-label="text"
              option-value="value"
              :disabled="!m.field || m.field === 'count'"
              :placeholder="$t('chart.edit.metric.function.placeholder')"
              class="w-full"
            />
          </div>
        </div>

        <slot name="metric-options" :metric="m" :report="report" />
      </div>
    </div>

    <Divider v-if="module && hasAxis" />

    <template v-if="hasAxis">
      <slot name="y-axis" :report="report" />
    </template>

    <Divider v-if="hasLegend" />

    <!-- Legend configuration -->
    <div v-if="hasLegend" class="px-3">
      <h5 class="mb-3">
        {{ $t('chart.edit.additionalConfig.legend.label') }}
      </h5>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.additionalConfig.legend.orientation.label') }}
          </label>
          <Select
            v-model="report.legend.orientation"
            :options="orientations"
            option-label="text"
            option-value="value"
            class="w-full"
          />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.additionalConfig.legend.show') }}
          </label>
          <div class="flex items-center gap-2">
            <ToggleSwitch v-model="legendVisible" />
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 mt-4">
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.additionalConfig.legend.align.label') }}
          </label>
          <Select
            v-model="report.legend.align"
            :options="alignments"
            option-label="text"
            option-value="value"
            :disabled="!report.legend.position?.isDefault"
            class="w-full"
          />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('chart.edit.additionalConfig.legend.options.label') }}
          </label>
          <div class="flex flex-col gap-2">
            <div class="flex items-center gap-2">
              <Checkbox
                v-model="report.legend.isScrollable"
                :binary="true"
                :disabled="report.legend.orientation !== 'horizontal'"
                input-id="legendScrollable"
              />
              <label for="legendScrollable">
                {{ $t('chart.edit.additionalConfig.legend.scrollable') }}
              </label>
            </div>
          </div>
        </div>
      </div>
    </div>

    <slot name="additional-config" :report="report" :metrics="metrics" :has-axis="hasAxis" />
  </div>
</template>

<script setup>
import { computed, toRaw } from 'vue'
import { useI18n } from 'vue-i18n'
import { compose } from '@cortezaproject/corteza-js-next'

const { t } = useI18n()

const props = defineProps({
  report: {
    type: Object,
    required: true,
  },
  chart: {
    type: Object,
    default: () => ({}),
  },
  modules: {
    type: Array,
    required: true,
  },
  supportedMetrics: {
    type: Number,
    default: -1,
  },
  dimensionFieldKind: {
    type: Array,
    default: () => ['DateTime', 'Select', 'Number', 'Bool', 'String', 'Record', 'User'],
  },
  usesDimensionsField: {
    type: Boolean,
    default: true,
  },
})

const emit = defineEmits(['update:report'])

// Translated options
const metricAggregates = compose.chartUtil.aggregateFunctions.map(af => ({
  ...af,
  text: t(`chart.edit.metric.function.${af.text}`),
}))

const dimensionModifiers = compose.chartUtil.dimensionFunctions.map(df => ({
  ...df,
  text: t(`chart.edit.dimension.function.${df.text}`),
}))

const predefinedFilters = compose.chartUtil.predefinedFilters.map(pf => ({
  ...pf,
  text: t(`chart.edit.filter.${pf.text}`),
}))

const alignments = [
  { value: 'left', text: t('chart.edit.additionalConfig.legend.align.left') },
  { value: 'center', text: t('chart.edit.additionalConfig.legend.align.center') },
  { value: 'right', text: t('chart.edit.additionalConfig.legend.align.right') },
]

const orientations = [
  { value: 'horizontal', text: t('chart.edit.additionalConfig.legend.orientation.horizontal') },
  { value: 'vertical', text: t('chart.edit.additionalConfig.legend.orientation.vertical') },
]

// Computed properties
const moduleID = computed({
  get: () => props.report.moduleID,
  set: v => {
    props.report.moduleID = v
    emit('update:report', { ...toRaw(props.report), moduleID: v })
  },
})

const reportFilter = computed({
  get: () => props.report.filter,
  set: v => {
    props.report.filter = v
  },
})

const metrics = computed({
  get: () => props.report.metrics || [],
  set: v => {
    props.report.metrics = v
  },
})

const dimensions = computed({
  get: () => props.report.dimensions || [],
  set: v => {
    props.report.dimensions = v
  },
})

const module = computed(() => {
  return props.modules.find(m => m.moduleID === moduleID.value)
})

const metricFields = computed(() => {
  if (!module.value) return []
  return [
    { value: 'count', text: t('chart.general.label.count') },
    ...module.value.fields
      .filter(f => f.kind === 'Number')
      .sort((a, b) => (a.label || a.name).localeCompare(b.label || b.name))
      .map(({ label, name }) => ({ value: name, text: label || name })),
  ]
})

const dimensionFields = computed(() => {
  if (!module.value) return []
  return [
    ...module.value.fields.sort((a, b) => (a.label || a.name).localeCompare(b.label || b.name)),
    ...(module.value.systemFields
      ? module.value.systemFields().map(sf => {
          sf.label = t(`field.system.${sf.name}`)
          return sf
        })
      : []),
  ]
    .filter(({ kind, options = {} }) => {
      return (
        props.dimensionFieldKind.includes(kind) && !(options.useRichTextEditor || options.multiLine)
      )
    })
    .map(({ name, label, kind }) => {
      return { value: name, text: `${label || name} (${kind})`, kind }
    })
})

const hasAxis = computed(() => {
  return metrics.value.some(({ type }) => ['bar', 'line', 'scatter'].includes(type))
})

const hasLegend = computed(() => {
  return !metrics.value.some(({ type }) => ['gauge'].includes(type))
})

const canAddMetric = computed(() => {
  return (
    (props.supportedMetrics < 0 || metrics.value.length < props.supportedMetrics) && moduleID.value
  )
})

const legendVisible = computed({
  get: () => !props.report.legend?.isHidden,
  set: v => {
    if (!props.report.legend) {
      props.report.legend = {}
    }
    props.report.legend.isHidden = !v
  },
})

// Methods
function getField(d) {
  if (!d.field || !module.value) return undefined
  return module.value.fields.find(f => f.name === d.field)
}

function isTemporalField(name) {
  return dimensionFields.value.some(f => f.value === name && f.kind === 'DateTime')
}

function onDimFieldChange(f, d) {
  if (!isTemporalField(f)) {
    d.modifier = dimensionModifiers[0]?.value
    d.timeLabels = false
  }
  if (d.meta) d.meta.fields = []
}

function onMetricFieldChange(field, m) {
  if (field === 'count') {
    m.aggregate = undefined
  } else if (field) {
    const moduleField = module.value?.fields.find(f => f.name === field)
    if (moduleField && moduleField.options) {
      const { presetFormat, format, prefix, suffix } = moduleField.options
      if (!m.formatting) m.formatting = {}
      if (presetFormat !== undefined) m.formatting.presetFormat = presetFormat
      if (format !== undefined) m.formatting.format = format
      if (prefix !== undefined) m.formatting.prefix = prefix
      if (suffix !== undefined) m.formatting.suffix = suffix
    }

    if (!m.aggregate) {
      m.aggregate = metricAggregates[0]?.value
    }
  }
}

function addMetric() {
  const chart = toRaw(props.chart)
  if (chart.defMetric) {
    metrics.value.push(chart.defMetric())
  } else {
    metrics.value.push({
      field: 'count',
      formatting: { presetFormat: 'custom', prefix: '', suffix: '', format: '' },
    })
  }
}

function removeMetric(i) {
  metrics.value.splice(i, 1)
}
</script>
