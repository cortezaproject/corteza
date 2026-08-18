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
        <Slider v-model="field.options.precision" :min="0" :max="6" class="w-full mt-2" />
      </div>
      <div class="flex flex-col gap-2">
        <label class="font-medium text-muted-color text-sm">
          {{ $t('field.kind.number.step.label') }}
        </label>
        <InputNumber v-model="field.options.step" :max-fraction-digits="6" class="w-full" />
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
      <!-- Min / Max -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="flex flex-col gap-2">
          <label class="text-primary font-medium text-sm">
            {{ $t('field.kind.number.progress.minimumValue') }}
          </label>
          <InputNumber v-model="field.options.min" show-buttons fluid />
        </div>
        <div class="flex flex-col gap-2">
          <label class="text-primary font-medium text-sm">
            {{ $t('field.kind.number.progress.maximumValue') }}
          </label>
          <InputNumber v-model="field.options.max" show-buttons fluid />
        </div>
      </div>

      <!-- Variant + Display options (two columns) -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <!-- Variant (left) -->
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('field.kind.number.progress.variants.default') }}
          </label>
          <Select
            v-model="field.options.variant"
            :options="variantOptions"
            option-label="label"
            option-value="value"
            class="w-full"
          >
            <template #value="{ value, placeholder }">
              <Button
                v-if="value"
                :label="variantLabel(value)"
                :severity="mapVariantSeverity(value)"
                size="small"
                class="pointer-events-none !py-0.5 !px-2 !text-xs"
              />
              <span v-else>{{ placeholder }}</span>
            </template>
            <template #option="{ option }">
              <Button
                :label="option.label"
                :severity="mapVariantSeverity(option.value)"
                size="small"
                class="pointer-events-none !py-0.5 !px-2 !text-xs"
              />
            </template>
          </Select>
        </div>

        <!-- Display options (right) -->
        <div class="flex flex-col gap-2">
          <label class="text-primary font-medium text-sm">
            {{ $t('field.kind.number.progress.show.value') }}
          </label>
          <div class="flex items-center gap-2">
            <Checkbox v-model="field.options.showValue" inputId="showValue" :binary="true" />
            <label for="showValue" class="text-sm cursor-pointer">
              {{ $t('field.kind.number.progress.show.value') }}
            </label>
          </div>
          <template v-if="field.options.showValue">
            <div class="flex items-center gap-2">
              <Checkbox
                v-model="field.options.showRelative"
                inputId="showRelative"
                :binary="true"
              />
              <label for="showRelative" class="text-sm cursor-pointer">
                {{ $t('field.kind.number.progress.show.relative') }}
              </label>
            </div>
            <div class="flex items-center gap-2">
              <Checkbox
                v-model="field.options.showProgress"
                inputId="showProgress"
                :binary="true"
              />
              <label for="showProgress" class="text-sm cursor-pointer">
                {{ $t('field.kind.number.progress.show.progress') }}
              </label>
            </div>
          </template>
        </div>
      </div>

      <!-- Thresholds -->
      <div class="flex flex-col gap-3">
        <div class="flex flex-col gap-1">
          <div class="flex items-center gap-2">
            <label class="text-primary font-medium text-sm">
              {{ $t('field.kind.number.progress.thresholds.label') }}
            </label>
            <Button
              :label="$t('general.label.add')"
              icon="pi pi-plus"
              size="small"
              severity="secondary"
              @click="addThreshold"
            />
          </div>
          <small class="text-muted-color">
            {{ $t('field.kind.number.progress.thresholds.description') }}
          </small>
        </div>

        <!-- Column headers (shown once above the first card) -->
        <div v-if="(field.options.thresholds || []).length" class="flex items-end gap-3 px-2 mt-2">
          <label class="flex-1 text-muted-color text-xs font-semibold uppercase tracking-wide">
            {{ $t('block.progress.thresholds.column.value') }}
          </label>
          <label class="flex-1 text-muted-color text-xs font-semibold uppercase tracking-wide">
            {{ $t('block.progress.thresholds.column.variant') }}
          </label>
          <div class="w-8 shrink-0" />
        </div>

        <Card
          v-for="(t, i) in field.options.thresholds || []"
          :key="i"
          :pt="{ body: { class: 'p-3' }, content: { class: 'p-0' } }"
          class="border border-surface"
        >
          <template #content>
            <div class="flex items-center gap-3">
              <InputNumber
                v-model="t.value"
                :min="0"
                :max="100"
                suffix="%"
                fluid
                class="flex-1"
                size="small"
              />
              <Select
                v-model="t.variant"
                :options="variantOptions"
                option-label="label"
                option-value="value"
                size="small"
                class="flex-1"
              >
                <template #value="{ value, placeholder }">
                  <Button
                    v-if="value"
                    :label="variantLabel(value)"
                    :severity="mapVariantSeverity(value)"
                    size="small"
                    class="pointer-events-none !py-0.5 !px-2 !text-xs"
                  />
                  <span v-else>{{ placeholder }}</span>
                </template>
                <template #option="{ option }">
                  <Button
                    :label="option.label"
                    :severity="mapVariantSeverity(option.value)"
                    size="small"
                    class="pointer-events-none !py-0.5 !px-2 !text-xs"
                  />
                </template>
              </Select>
              <Button
                v-tooltip.top="$t('general.label.delete')"
                icon="pi pi-trash"
                severity="danger"
                text
                size="small"
                class="shrink-0"
                @click="removeThreshold(i)"
              />
            </div>
          </template>
        </Card>
      </div>
    </template>

    <Divider />

    <!-- Live example -->
    <div class="flex flex-col gap-2">
      <label class="text-primary font-medium text-sm">
        {{ $t('field.kind.number.liveExample') }}
      </label>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 items-center">
        <InputNumber
          v-model="liveExampleInput"
          :step="field.options.step || 1"
          :max-fraction-digits="6"
          fluid
        />
        <!-- Progress mode: render the actual viewer so styling/thresholds match exactly -->
        <CFieldViewer
          v-if="isProgress"
          :key="previewKey"
          :field="field"
          :record="previewRecord"
          value-only
          disable-click
        />
        <!-- Number mode: show formatted output -->
        <div v-else class="text-color font-medium text-lg">
          {{ liveExampleOutput }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { components } from '@planetcrust/human-vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

const { CFieldViewer } = components

const { t } = useI18n()

const field = inject('fieldDraft')

const isProgress = computed(() => field.value.options?.display === 'progress')

const displayOptions = computed(() => [
  { value: 'number', label: t('field.kind.number.displayType.number') },
  { value: 'progress', label: t('field.kind.number.displayType.progress') },
])

// Variants aligned with block progress — 'light' is migrated to 'secondary' on load
const variantKeys = ['primary', 'secondary', 'success', 'warning', 'danger', 'info', 'dark']
const variantLabel = key => t(`field.kind.number.progress.variants.${key}`)
const variantOptions = computed(() =>
  variantKeys.map(value => ({ value, label: variantLabel(value) })),
)

// Every variant maps to a PrimeVue Button severity — same color tokens as the viewer's progress fill
const variantSeverityMap = {
  primary: undefined,
  secondary: 'secondary',
  success: 'success',
  warning: 'warn',
  danger: 'danger',
  info: 'info',
  dark: 'contrast',
}
const mapVariantSeverity = key => variantSeverityMap[key]

const formatOptions = computed(() => [
  { value: 'custom', label: t('field.kind.number.presetFormats.options.custom') },
  { value: 'accounting', label: t('field.kind.number.presetFormats.options.accounting') },
])

const presetDescription = computed(() =>
  t(`field.kind.number.presetFormats.description.${field.value.options.presetFormat || 'custom'}`),
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

watch(
  () => field.value.options?.display,
  display => {
    liveExampleInput.value = display === 'progress' ? 33.45679 : 1234.56789
  },
)

// Fake record passed to CFieldViewer for the progress-mode live preview.
// The viewer reads record.values[field.name], so the key must match field.name
// (even if empty, '' is a valid object key).
const previewRecord = computed(() => ({
  values: { [field.value.name || '']: liveExampleInput.value },
}))

// Key forces a re-mount when options that affect rendering change, so the preview
// reliably re-renders even if reactivity doesn't propagate deep into the viewer.
const previewKey = computed(() => {
  const o = field.value.options || {}
  return [
    liveExampleInput.value,
    o.min,
    o.max,
    o.variant,
    o.showValue,
    o.showRelative,
    o.showProgress,
    JSON.stringify(o.thresholds || []),
  ].join('|')
})

const liveExampleOutput = computed(() => {
  const val = liveExampleInput.value
  if (val === null || val === undefined) return ''

  if (isProgress.value) {
    const min = parseFloat(field.value.options?.min || 0)
    const max = parseFloat(field.value.options?.max || 100)
    const pct = max === min ? 0 : Math.round(((val - min) / (max - min)) * 100)
    return `${pct}%`
  }

  // Touch reactive dependencies so Vue recomputes when these change

  void (field.value.options?.format,
  field.value.options?.presetFormat,
  field.value.options?.precision,
  field.value.options?.prefix,
  field.value.options?.suffix)

  // Use the field's formatValue method (from lib/js ModuleFieldNumber)
  // which handles numeral format strings, accounting mode, prefix, suffix
  if (typeof field.value.formatValue === 'function') {
    return field.value.formatValue(String(val))
  }

  // Fallback
  const { precision, prefix = '', suffix = '' } = field.value.options || {}
  const num = Number(val)
  if (isNaN(num)) return val
  const formatted = precision !== undefined ? num.toFixed(Number(precision)) : num.toLocaleString()
  return `${prefix}${formatted}${suffix}`
})

function addThreshold() {
  if (!field.value.options.thresholds) {
    field.value.options.thresholds = []
  }
  field.value.options.thresholds.push({ value: 0, variant: 'success' })
}

function removeThreshold(index) {
  if (index > -1) {
    field.value.options.thresholds.splice(index, 1)
  }
}

onMounted(() => {
  if (!field.value.options.display) {
    field.value.options.display = 'number'
  }
  if (!field.value.options.presetFormat) {
    field.value.options.presetFormat = 'custom'
  }
  if (!field.value.options.thresholds) {
    field.value.options.thresholds = []
  }
})
</script>
