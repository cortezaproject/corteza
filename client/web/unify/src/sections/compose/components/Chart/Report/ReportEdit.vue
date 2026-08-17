<template>
  <div class="flex flex-col gap-4">
    <!-- Configure source module -->
    <Panel :header="$t('chart.edit.module.title')" toggleable>
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <CFormGroup name="moduleID" :label="$t('chart.edit.module.label')" required>
          <Select
            id="moduleID"
            name="moduleID"
            v-model="moduleID"
            :options="modules"
            option-label="name"
            option-value="moduleID"
            :placeholder="$t('chart.edit.module.placeholder')"
            class="w-full"
            filter
          />
        </CFormGroup>

        <CFormGroup v-if="module" :label="$t('chart.edit.filter.preset')" input-id="presetFilter">
          <Select
            id="presetFilter"
            v-model="reportFilter"
            :options="predefinedFilters"
            option-label="text"
            option-value="value"
            :placeholder="$t('chart.edit.filter.noFilter')"
            class="w-full"
            show-clear
          />
        </CFormGroup>

        <!-- Configure report filters -->
        <CFormGroup
          v-if="module"
          :label="$t('chart.edit.filter.label')"
          class="col-span-1 lg:col-span-2"
        >
          <CInputExpression
            ref="filterInput"
            v-model="reportFilter"
            dialect="ql"
            :scope="scope"
            :query-fields="module?.fields || []"
            :placeholder="$t('chart.edit.filter.placeholder')"
          />
          <!-- Syntax first, then the variables: the footnote describes what to
               type, the hint lists what can be dropped into it. -->
          <small class="text-muted-color">{{ $t('chart.edit.filter.footnote') }}</small>
          <CExpressionHint
            :scope="scope"
            depends-on-placement
            class="block"
            @insert="filterInput?.insert($event)"
          />
        </CFormGroup>
      </div>
    </Panel>

    <!-- Configure report dimensions -->
    <Panel v-if="module" :header="$t('chart.edit.dimension.label')" toggleable>
      <div v-for="(d, i) in dimensions" :key="i">
        <template v-if="usesDimensionsField">
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
            <CFormGroup :label="$t('chart.edit.dimension.fieldLabel')" input-id="dimensionField">
              <Select
                id="dimensionField"
                v-model="d.field"
                :options="dimensionFields"
                option-label="text"
                option-value="value"
                :placeholder="$t('chart.edit.dimension.fieldPlaceholder')"
                class="w-full"
                filter
                @change="e => onDimFieldChange(e.value, d)"
              />
            </CFormGroup>

            <CFormGroup
              :label="$t('chart.edit.dimension.function.label')"
              input-id="dimensionModifier"
            >
              <Select
                id="dimensionModifier"
                v-model="d.modifier"
                :options="dimensionModifiers"
                option-label="text"
                option-value="value"
                :disabled="!d.field || !isTemporalField(d.field)"
                :placeholder="$t('chart.edit.dimension.function.placeholder')"
                class="w-full"
              />
            </CFormGroup>
          </div>

          <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 mt-4">
            <CFormGroup
              :label="$t('chart.edit.dimension.defaultValueLabel')"
              :description="$t('chart.edit.dimension.defaultValueFootnote')"
              input-id="dimensionDefault"
            >
              <InputText id="dimensionDefault" v-model="d.default" class="w-full" />
            </CFormGroup>

            <CFormGroup :label="$t('chart.edit.dimension.options.label')">
              <div class="flex items-center gap-2">
                <Checkbox v-model="d.skipMissing" :binary="true" input-id="skipMissing" />
                <label for="skipMissing">{{ $t('chart.edit.dimension.skipMissingValues') }}</label>
              </div>
              <slot
                name="dimension-options-options"
                :dimension="d"
                :is-temporal="isTemporalField(d.field)"
              />
            </CFormGroup>
          </div>
        </template>

        <slot name="dimension-options" :index="i" :dimension="d" :field="getField(d)" />
      </div>
    </Panel>

    <!-- Configure report metrics -->
    <Panel v-if="module" :header="$t('chart.edit.metric.title')" toggleable>
      <template #icons>
        <Button
          v-if="canAddMetric"
          :label="'+ ' + $t('chart.edit.metric.add')"
          text
          size="small"
          @click="addMetric"
        />
      </template>

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
          <CFormGroup :label="$t('chart.edit.metric.fieldLabel')" :input-id="`metricField${i}`">
            <Select
              :id="`metricField${i}`"
              v-model="m.field"
              :options="metricFields"
              option-label="text"
              option-value="value"
              :placeholder="$t('chart.edit.metric.fieldPlaceholder')"
              class="w-full"
              filter
              @change="e => onMetricFieldChange(e.value, m)"
            />
          </CFormGroup>

          <CFormGroup
            :label="$t('chart.edit.metric.function.label')"
            :input-id="`metricAggregate${i}`"
          >
            <Select
              :id="`metricAggregate${i}`"
              v-model="m.aggregate"
              :options="metricAggregates"
              option-label="text"
              option-value="value"
              :disabled="!m.field || m.field === 'count'"
              :placeholder="$t('chart.edit.metric.function.placeholder')"
              class="w-full"
            />
          </CFormGroup>
        </div>

        <slot name="metric-options" :metric="m" :report="report" :index="i" />
      </div>
    </Panel>

    <template v-if="hasAxis">
      <slot name="y-axis" :report="report" />
    </template>

    <!-- Legend configuration -->
    <Panel
      v-if="hasLegend"
      :header="$t('chart.edit.additionalConfig.legend.label')"
      toggleable
      collapsed
    >
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <CFormGroup
          :label="$t('chart.edit.additionalConfig.legend.orientation.label')"
          input-id="legendOrientation"
        >
          <Select
            id="legendOrientation"
            v-model="report.legend.orientation"
            :options="orientations"
            option-label="text"
            option-value="value"
            class="w-full"
          />
        </CFormGroup>

        <CFormGroup :label="$t('chart.edit.additionalConfig.legend.show')" input-id="legendVisible">
          <ToggleSwitch input-id="legendVisible" v-model="legendVisible" />
        </CFormGroup>

        <CFormGroup
          :label="$t('chart.edit.additionalConfig.legend.align.label')"
          input-id="legendAlign"
        >
          <Select
            id="legendAlign"
            v-model="report.legend.align"
            :options="alignments"
            option-label="text"
            option-value="value"
            :disabled="!report.legend.position?.isDefault"
            class="w-full"
          />
        </CFormGroup>

        <CFormGroup :label="$t('chart.edit.additionalConfig.legend.options.label')">
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
        </CFormGroup>
      </div>
    </Panel>

    <slot name="additional-config" :report="report" :metrics="metrics" :has-axis="hasAxis" />
  </div>
</template>

<script setup>
import { computed, inject, ref, toRaw } from 'vue'
import { useI18n } from 'vue-i18n'
import { compose } from '@planetcrust/human-js'
import { useExpressionScope } from '@/sections/compose/composables/useExpressionScope'

const { t } = useI18n()

const props = defineProps({
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

const report = inject('reportDraft')

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
  get: () => report.value.moduleID,
  set: v => {
    report.value.moduleID = v
  },
})

const reportFilter = computed({
  get: () => report.value.filter,
  set: v => {
    report.value.filter = v
  },
})

const metrics = computed({
  get: () => report.value.metrics || [],
  set: v => {
    report.value.metrics = v
  },
})

const dimensions = computed({
  get: () => report.value.dimensions || [],
  set: v => {
    report.value.dimensions = v
  },
})

const module = computed(() => {
  return props.modules.find(m => m.moduleID === moduleID.value)
})

const filterInput = ref(null)

// A chart is namespace-level and may be placed on a record page or not, so the
// record variables are offered and the hint says they depend on placement.
const { scope } = useExpressionScope({ hasRecord: true })

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
    ...[...module.value.fields].sort((a, b) =>
      (a.label || a.name).localeCompare(b.label || b.name),
    ),
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
  get: () => !report.value.legend?.isHidden,
  set: v => {
    if (!report.value.legend) {
      report.value.legend = {}
    }
    report.value.legend.isHidden = !v
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
