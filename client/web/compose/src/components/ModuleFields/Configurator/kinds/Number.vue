<template>
  <div class="flex flex-col gap-6">
    <!-- Display type -->
    <div>
      <label class="font-medium text-muted-color text-sm block mb-3">
        {{ $t('field.kind.number.displayType.label') }}
      </label>
      <SelectButton
        :model-value="field.options.display"
        :options="displayOptions"
        option-label="label"
        option-value="value"
        @update:model-value="field.options.display = $event"
      />
    </div>

    <!-- Precision + Step (side by side) -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="flex flex-col gap-2">
        <label class="font-medium text-muted-color text-sm">
          {{ $t('field.kind.number.precisionLabel') }} ({{ field.options.precision ?? 0 }})
        </label>
        <Slider
          v-model="field.options.precision"
          :min="0"
          :max="6"
          class="w-full mt-2"
        />
      </div>
      <div class="flex flex-col gap-2">
        <label class="font-medium text-muted-color text-sm">
          {{ $t('field.kind.number.step.label') }}
        </label>
        <InputNumber
          v-model="field.options.step"
          :max-fraction-digits="6"
          class="w-full"
        />
        <small class="text-muted-color">{{ $t('field.kind.number.step.description') }}</small>
      </div>
    </div>

    <Divider />

    <!-- Number display options -->
    <template v-if="!isProgress">
      <!-- Prefix / Suffix -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="flex flex-col gap-2">
          <label class="font-medium text-muted-color text-sm">
            {{ $t('field.kind.number.prefixLabel') }}
          </label>
          <InputText v-model="field.options.prefix" placeholder="$" class="w-full" />
        </div>
        <div class="flex flex-col gap-2">
          <label class="font-medium text-muted-color text-sm">
            {{ $t('field.kind.number.suffixLabel') }}
          </label>
          <InputText v-model="field.options.suffix" placeholder="USD" class="w-full" />
        </div>
      </div>

      <!-- Preset format + Custom format -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="flex flex-col gap-2">
          <label class="font-medium text-muted-color text-sm">
            {{ $t('field.kind.number.presetFormats.label') }}
          </label>
          <Select
            v-model="field.options.presetFormat"
            :options="formatOptions"
            option-label="label"
            option-value="value"
            class="w-full"
          />
          <small class="text-muted-color">{{ presetDescription }}</small>
        </div>
        <div class="flex flex-col gap-2">
          <label class="font-medium text-muted-color text-sm">
            {{ $t('field.kind.number.formatLabel') }}
          </label>
          <InputText
            v-model="field.options.format"
            :disabled="field.options.presetFormat !== 'custom'"
            placeholder="0.00"
            class="w-full"
          />
        </div>
      </div>

      <!-- Format examples table -->
      <div class="flex flex-col gap-2">
        <label class="font-medium text-muted-color text-sm">
          {{ $t('field.kind.number.examplesLabel') }}
        </label>
        <DataTable :value="formatExamples" size="small" class="text-sm">
          <Column field="input" :header="$t('field.kind.number.exampleInput')" />
          <Column field="format" :header="$t('field.kind.number.exampleFormat')" />
          <Column field="result" :header="$t('field.kind.number.exampleResult')" />
        </DataTable>
      </div>
    </template>

    <!-- Progress bar options -->
    <template v-else>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="flex flex-col gap-2">
          <label class="font-medium text-muted-color text-sm">
            {{ $t('field.kind.number.progress.minimumValue') }}
          </label>
          <InputNumber v-model="field.options.min" show-buttons class="w-full" />
        </div>
        <div class="flex flex-col gap-2">
          <label class="font-medium text-muted-color text-sm">
            {{ $t('field.kind.number.progress.maximumValue') }}
          </label>
          <InputNumber v-model="field.options.max" show-buttons class="w-full" />
        </div>
      </div>

      <!-- Variant selector -->
      <div class="flex flex-col gap-2">
        <label class="font-medium text-muted-color text-sm">
          {{ $t('field.kind.number.progress.variants.default') }}
        </label>
        <Select
          v-model="field.options.variant"
          :options="variantOptions"
          option-label="label"
          option-value="value"
          class="w-full md:w-1/2"
        />
      </div>

      <!-- Checkboxes -->
      <div class="flex flex-wrap gap-6">
        <div class="flex flex-col gap-3">
          <div class="flex items-center gap-2">
            <Checkbox v-model="field.options.showValue" inputId="showValue" :binary="true" />
            <label for="showValue" class="cursor-pointer">{{ $t('field.kind.number.progress.show.value') }}</label>
          </div>
          <div class="flex items-center gap-2">
            <Checkbox v-model="field.options.animated" inputId="animated" :binary="true" />
            <label for="animated" class="cursor-pointer">{{ $t('field.kind.number.progress.animated') }}</label>
          </div>
        </div>
        <div v-if="field.options.showValue" class="flex flex-col gap-3">
          <div class="flex items-center gap-2">
            <Checkbox v-model="field.options.showRelative" inputId="showRelative" :binary="true" />
            <label for="showRelative" class="cursor-pointer">{{ $t('field.kind.number.progress.show.relative') }}</label>
          </div>
          <div class="flex items-center gap-2">
            <Checkbox v-model="field.options.showProgress" inputId="showProgress" :binary="true" />
            <label for="showProgress" class="cursor-pointer">{{ $t('field.kind.number.progress.show.progress') }}</label>
          </div>
        </div>
      </div>

      <!-- Thresholds -->
      <div class="flex flex-col gap-3">
        <div class="flex items-center gap-2">
          <label class="font-medium text-muted-color text-sm">
            {{ $t('field.kind.number.progress.thresholds.label') }}
          </label>
          <Button
            :label="$t('general.label.add')"
            icon="pi pi-plus"
            text
            size="small"
            @click="addThreshold"
          />
        </div>
        <small class="text-muted-color">{{ $t('field.kind.number.progress.thresholds.description') }}</small>

        <div
          v-for="(t, i) in (field.options.thresholds || [])"
          :key="i"
          class="flex items-center gap-3"
        >
          <div class="flex items-center gap-1 flex-1">
            <InputNumber v-model="t.value" class="flex-1" />
            <span class="text-muted-color">%</span>
          </div>
          <Select
            v-model="t.variant"
            :options="variantOptions"
            option-label="label"
            option-value="value"
            class="flex-1"
          />
          <Button
            icon="pi pi-times"
            severity="danger"
            text
            rounded
            size="small"
            @click="removeThreshold(i)"
          />
        </div>
      </div>
    </template>

    <Divider />

    <!-- Live example -->
    <div class="flex flex-col gap-2">
      <label class="font-medium text-muted-color text-sm">
        {{ $t('field.kind.number.liveExample') }}
      </label>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 items-center">
        <InputNumber
          v-model="liveExampleInput"
          :step="field.options.step || 1"
          :max-fraction-digits="6"
          class="w-full"
        />
        <div class="text-color font-medium text-lg">
          {{ liveExampleOutput }}
        </div>
      </div>
    </div>

    <CConfiguratorMultiDelimiter :field="field" />
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import CConfiguratorMultiDelimiter from '../CConfiguratorMultiDelimiter.vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import ProgressBar from 'primevue/progressbar'

const { t } = useI18n()

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
})

const isProgress = computed(() => props.field.options?.display === 'progress')

const displayOptions = computed(() => [
  { value: 'number', label: t('field.kind.number.displayType.number') },
  { value: 'progress', label: t('field.kind.number.displayType.progress') },
])

const variantOptions = computed(() => [
  { value: 'primary', label: t('field.kind.number.progress.variants.primary') },
  { value: 'secondary', label: t('field.kind.number.progress.variants.secondary') },
  { value: 'success', label: t('field.kind.number.progress.variants.success') },
  { value: 'warning', label: t('field.kind.number.progress.variants.warning') },
  { value: 'danger', label: t('field.kind.number.progress.variants.danger') },
  { value: 'info', label: t('field.kind.number.progress.variants.info') },
  { value: 'light', label: t('field.kind.number.progress.variants.light') },
  { value: 'dark', label: t('field.kind.number.progress.variants.dark') },
])

const formatOptions = computed(() => [
  { value: 'custom', label: t('field.kind.number.presetFormats.options.custom') },
  { value: 'accounting', label: t('field.kind.number.presetFormats.options.accounting') },
])

const presetDescription = computed(() =>
  t(`field.kind.number.presetFormats.description.${props.field.options.presetFormat || 'custom'}`),
)

const formatExamples = computed(() => [
  { input: '1000.234', format: '0,0.00', result: '1,000.23' },
  { input: '1000.234', format: '0,0', result: '1,000' },
  { input: '0.974878234', format: '0.000%', result: '97.488%' },
  { input: '100', format: '0o', result: '100th' },
  { input: '238', format: '00:00:00', result: '0:03:58' },
])

// Live example
const liveExampleInput = ref(1234.56789)

watch(() => props.field.options?.display, (display) => {
  liveExampleInput.value = display === 'progress' ? 33.45679 : 1234.56789
})

const liveExampleOutput = computed(() => {
  const val = liveExampleInput.value
  if (val === null || val === undefined) return ''

  if (isProgress.value) {
    const min = parseFloat(props.field.options?.min || 0)
    const max = parseFloat(props.field.options?.max || 100)
    const pct = max === min ? 0 : Math.round(((val - min) / (max - min)) * 100)
    return `${pct}%`
  }

  // Touch reactive dependencies so Vue recomputes when these change
   
  void (props.field.options?.format, props.field.options?.presetFormat,
    props.field.options?.precision, props.field.options?.prefix, props.field.options?.suffix)

  // Use the field's formatValue method (from lib/js ModuleFieldNumber)
  // which handles numeral format strings, accounting mode, prefix, suffix
  if (typeof props.field.formatValue === 'function') {
    return props.field.formatValue(String(val))
  }

  // Fallback
  const { precision, prefix = '', suffix = '' } = props.field.options || {}
  const num = Number(val)
  if (isNaN(num)) return val
  const formatted = precision !== undefined ? num.toFixed(Number(precision)) : num.toLocaleString()
  return `${prefix}${formatted}${suffix}`
})

function addThreshold() {
  if (!props.field.options.thresholds) {
    props.field.options.thresholds = []
  }
  props.field.options.thresholds.push({ value: 0, variant: 'success' })
}

function removeThreshold(index) {
  if (index > -1) {
    props.field.options.thresholds.splice(index, 1)
  }
}

onMounted(() => {
  if (!props.field.options.display) {
    props.field.options.display = 'number'
  }
  if (!props.field.options.presetFormat) {
    props.field.options.presetFormat = 'custom'
  }
  if (!props.field.options.thresholds) {
    props.field.options.thresholds = []
  }
})
</script>
