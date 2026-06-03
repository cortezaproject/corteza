<template>
  <div>
    <div class="flex flex-col gap-5 pb-16">
      <!-- Value Section -->
      <Fieldset :legend="$t('block.progress.value.label')" toggleable>
        <div class="flex flex-col gap-3">
          <CFormGroup
            :label="$t('block.progress.source.static.label')"
            :description="$t('block.progress.source.static.description')"
            :class="['transition-opacity', { 'opacity-40': !!valueModuleID }]"
          >
            <InputNumber
              v-model="fixedValue"
              :placeholder="$t('block.progress.value.fixed')"
              :min="0"
              :disabled="!!valueModuleID"
              class="md:w-1/2"
              fluid
            />
          </CFormGroup>

          <div class="flex items-center gap-2">
            <div class="flex-1 border-t border-surface" />
            <span class="text-muted-color text-xs font-medium">{{ $t('general.label.or') }}</span>
            <div class="flex-1 border-t border-surface" />
          </div>

          <CFormGroup
            :label="$t('block.progress.source.module.label')"
            :description="$t('block.progress.source.module.description')"
          >
            <Select
              v-model="valueModuleID"
              :options="modules"
              option-label="name"
              option-value="moduleID"
              :placeholder="$t('block.progress.module.select')"
              class="w-full md:w-1/2"
              filter
              show-clear
            />
          </CFormGroup>

          <template v-if="valueModuleID">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-2">
              <Select
                v-model="valueField"
                :options="getNumberFields(valueModuleID)"
                option-label="label"
                option-value="name"
                :placeholder="$t('block.progress.value.fieldPlaceholder')"
                class="w-full"
              />
              <Select
                v-model="valueOperation"
                :options="operationOptions"
                option-label="label"
                option-value="value"
                :placeholder="$t('block.progress.value.operationPlaceholder')"
                class="w-full"
                :disabled="!valueField || valueField === 'count'"
              />
            </div>
            <Textarea
              v-model="valueFilter"
              :placeholder="$t('block.progress.value.filterPlaceholder')"
              rows="2"
              class="w-full"
            />
          </template>
        </div>
      </Fieldset>

      <Divider />

      <!-- Min Value Section -->
      <Fieldset :legend="$t('block.progress.minValue.label')" toggleable>
        <div class="flex flex-col gap-3">
          <CFormGroup
            :label="$t('block.progress.source.static.label')"
            :description="$t('block.progress.source.static.description')"
            :class="['transition-opacity', { 'opacity-40': !!minValueModuleID }]"
          >
            <InputNumber
              v-model="fixedMinValue"
              :placeholder="$t('block.progress.minValue.fixed')"
              :min="0"
              :disabled="!!minValueModuleID"
              class="md:w-1/2"
              fluid
            />
          </CFormGroup>

          <div class="flex items-center gap-2">
            <div class="flex-1 border-t border-surface" />
            <span class="text-muted-color text-xs font-medium">{{ $t('general.label.or') }}</span>
            <div class="flex-1 border-t border-surface" />
          </div>

          <CFormGroup
            :label="$t('block.progress.source.module.label')"
            :description="$t('block.progress.source.module.description')"
          >
            <Select
              v-model="minValueModuleID"
              :options="modules"
              option-label="name"
              option-value="moduleID"
              :placeholder="$t('block.progress.module.select')"
              class="w-full md:w-1/2"
              filter
              show-clear
            />
          </CFormGroup>

          <template v-if="minValueModuleID">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-2">
              <Select
                v-model="minValueField"
                :options="getNumberFields(minValueModuleID)"
                option-label="label"
                option-value="name"
                :placeholder="$t('block.progress.value.fieldPlaceholder')"
                class="w-full"
              />
              <Select
                v-model="minValueOperation"
                :options="operationOptions"
                option-label="label"
                option-value="value"
                :placeholder="$t('block.progress.value.operationPlaceholder')"
                class="w-full"
                :disabled="!minValueField || minValueField === 'count'"
              />
            </div>
            <Textarea
              v-model="minValueFilter"
              :placeholder="$t('block.progress.value.filterPlaceholder')"
              rows="2"
              class="w-full"
            />
          </template>
        </div>
      </Fieldset>

      <Divider />

      <!-- Max Value Section -->
      <Fieldset :legend="$t('block.progress.maxValue.label')" toggleable>
        <div class="flex flex-col gap-3">
          <CFormGroup
            :label="$t('block.progress.source.static.label')"
            :description="$t('block.progress.source.static.description')"
            :class="['transition-opacity', { 'opacity-40': !!maxValueModuleID }]"
          >
            <InputNumber
              v-model="fixedMaxValue"
              :placeholder="$t('block.progress.maxValue.fixed')"
              :min="0"
              :disabled="!!maxValueModuleID"
              class="md:w-1/2"
              fluid
            />
          </CFormGroup>

          <div class="flex items-center gap-2">
            <div class="flex-1 border-t border-surface" />
            <span class="text-muted-color text-xs font-medium">{{ $t('general.label.or') }}</span>
            <div class="flex-1 border-t border-surface" />
          </div>

          <CFormGroup
            :label="$t('block.progress.source.module.label')"
            :description="$t('block.progress.source.module.description')"
          >
            <Select
              v-model="maxValueModuleID"
              :options="modules"
              option-label="name"
              option-value="moduleID"
              :placeholder="$t('block.progress.module.select')"
              class="w-full md:w-1/2"
              filter
              show-clear
            />
          </CFormGroup>

          <template v-if="maxValueModuleID">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-2">
              <Select
                v-model="maxValueField"
                :options="getNumberFields(maxValueModuleID)"
                option-label="label"
                option-value="name"
                :placeholder="$t('block.progress.value.fieldPlaceholder')"
                class="w-full"
              />
              <Select
                v-model="maxValueOperation"
                :options="operationOptions"
                option-label="label"
                option-value="value"
                :placeholder="$t('block.progress.value.operationPlaceholder')"
                class="w-full"
                :disabled="!maxValueField || maxValueField === 'count'"
              />
            </div>
            <Textarea
              v-model="maxValueFilter"
              :placeholder="$t('block.progress.value.filterPlaceholder')"
              rows="2"
              class="w-full"
            />
          </template>
        </div>
      </Fieldset>

      <Divider />

      <!-- Display Options -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <CFormGroup :label="$t('block.progress.variant.label')">
          <Select
            v-model="variant"
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
                class="pointer-events-none !py-0.5 !px-2 !text-sm"
              />
              <span v-else>{{ placeholder }}</span>
            </template>
            <template #option="{ option }">
              <Button
                :label="option.label"
                :severity="mapVariantSeverity(option.value)"
                size="small"
                class="pointer-events-none !py-0.5 !px-2 !text-sm"
              />
            </template>
          </Select>
        </CFormGroup>

        <CFormGroup :label="$t('block.progress.display-options')">
          <div class="flex items-center gap-2">
            <Checkbox v-model="showValue" binary input-id="showValue" />
            <label for="showValue" class="text-sm">{{ $t('block.progress.show.value') }}</label>
          </div>

          <template v-if="showValue">
            <div class="flex items-center gap-2">
              <Checkbox v-model="showRelative" binary input-id="showRelative" />
              <label for="showRelative" class="text-sm">
                {{ $t('block.progress.show.relative') }}
              </label>
            </div>

            <div class="flex items-center gap-2">
              <Checkbox v-model="showProgress" binary input-id="showProgress" />
              <label for="showProgress" class="text-sm">
                {{ $t('block.progress.show.progress') }}
              </label>
            </div>
          </template>
        </CFormGroup>
      </div>

      <Divider />

      <!-- Thresholds -->
      <CFormGroup :label="$t('block.progress.thresholds.label')">
        <template #actions>
          <Button
            :label="$t('general.label.add')"
            icon="pi pi-plus"
            size="small"
            severity="secondary"
            @click="addThreshold"
          />
        </template>

        <!-- Column headers (shown once above the first card) -->
        <div v-if="thresholds.length" class="flex items-end gap-3 px-2">
          <label class="flex-1 text-muted-color text-xs font-semibold uppercase tracking-wide">
            {{ $t('block.progress.thresholds.column.value') }}
          </label>
          <label class="flex-1 text-muted-color text-xs font-semibold uppercase tracking-wide">
            {{ $t('block.progress.thresholds.column.variant') }}
          </label>
          <!-- Spacer matching the delete button column width -->
          <div class="w-8 shrink-0" />
        </div>

        <Card
          v-for="(threshold, i) in thresholds"
          :key="i"
          :pt="{ body: { class: 'p-3' }, content: { class: 'p-0' } }"
          class="border border-surface"
        >
          <template #content>
            <div class="flex items-center gap-3">
              <InputNumber
                :model-value="threshold.value"
                :min="0"
                :max="100"
                suffix="%"
                fluid
                size="small"
                class="flex-1"
                @update:model-value="updateThreshold(i, 'value', $event)"
              />
              <Select
                :model-value="threshold.variant"
                :options="variantOptions"
                option-label="label"
                option-value="value"
                size="small"
                class="flex-1"
                @update:model-value="updateThreshold(i, 'variant', $event)"
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

        <small v-if="!thresholds.length" class="text-muted-color">
          {{ $t('block.progress.thresholds.empty') }}
        </small>
      </CFormGroup>
    </div>

    <!-- Sticky live preview -->
    <div
      class="sticky bottom-0 left-0 w-full bg-surface rounded-border shadow p-3 z-10"
    >
      <CFormGroup :label="$t('block.progress.preview')">
        <div class="h-14">
          <ProgressBlock :key="previewFetchKey" :block="previewBlock" :namespace="namespace" />
        </div>
      </CFormGroup>
    </div>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import { useModuleStore } from '@/sections/compose/stores/module'
import ProgressBlock from '../Blocks/ProgressBlock.vue'

const { t } = useI18n()
const moduleStore = useModuleStore()

defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const block = inject('blockDraft')

const previewBlock = computed(() => ({
  ...block.value,
  options: { ...block.value?.options, magnifyOption: '' },
  fetch: block.value?.fetch?.bind(block.value),
}))

// Key encodes all options that affect the fetch — changes here remount the preview
// and guarantee a fresh onMounted → refresh() regardless of watcher subtleties.
const previewFetchKey = computed(() => {
  const o = block.value?.options
  if (!o) return ''
  return [
    o.value?.moduleID,
    o.value?.field,
    o.value?.operation,
    o.value?.filter,
    o.minValue?.moduleID,
    o.minValue?.field,
    o.minValue?.operation,
    o.minValue?.filter,
    o.maxValue?.moduleID,
    o.maxValue?.field,
    o.maxValue?.operation,
    o.maxValue?.filter,
  ].join('\0')
})

const modules = computed(() => moduleStore.set || [])

const operationOptions = [
  { value: 'sum', label: t('block.progress.operation.sum') },
  { value: 'max', label: t('block.progress.operation.max') },
  { value: 'min', label: t('block.progress.operation.min') },
  { value: 'avg', label: t('block.progress.operation.avg') },
]

const variantKeys = ['primary', 'secondary', 'success', 'warning', 'danger', 'info', 'dark']

const variantLabel = key => t(`block.progress.variant.${key}`)

const variantOptions = variantKeys.map(value => ({ value, label: variantLabel(value) }))

// Every variant maps directly to a PrimeVue Button severity — so the dropdown tag
// and the progress bar fill both read the same --p-button-{severity}-* tokens.
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

function getNumberFields(moduleID) {
  const countOption = { name: 'count', label: t('block.progress.count') }
  if (!moduleID) return [countOption]
  const mod = moduleStore.getByID(moduleID)
  if (!mod) return [countOption]
  const numberFields = (mod.fields || [])
    .filter(f => f.kind === 'Number')
    .map(f => ({ name: f.name, label: f.label || f.name }))
    .sort((a, b) => a.label.localeCompare(b.label))
  return [countOption, ...numberFields]
}

function updateOptions(key, value) {
  if (!block.value.options) block.value.options = {}
  block.value.options[key] = value
}

function updateNestedOptions(section, key, value) {
  const current = { ...(block.value.options?.[section] || {}) }
  current[key] = value
  updateOptions(section, current)
}

function updateDisplay(key, value) {
  const display = { ...block.value.options?.display, [key]: value }
  updateOptions('display', display)
}

// --- Value ---
const fixedValue = computed({
  get: () => block.value.options?.value?.default ?? 0,
  set: v => updateNestedOptions('value', 'default', v),
})

const valueModuleID = computed({
  // Return null (not '') when empty so PrimeVue Select hides its clear icon properly
  get: () => block.value.options?.value?.moduleID || null,
  set: v => {
    const next = v || ''
    updateNestedOptions('value', 'moduleID', next)
    // Default field to 'count' when a module is freshly picked
    if (next && !block.value.options?.value?.field) {
      updateNestedOptions('value', 'field', 'count')
    }
  },
})

const valueField = computed({
  get: () => block.value.options?.value?.field || 'count',
  set: v => {
    updateNestedOptions('value', 'field', v || 'count')
    if (!v || v === 'count') updateNestedOptions('value', 'operation', '')
  },
})

const valueOperation = computed({
  get: () => block.value.options?.value?.operation || '',
  set: v => updateNestedOptions('value', 'operation', v),
})

const valueFilter = computed({
  get: () => block.value.options?.value?.filter || '',
  set: v => updateNestedOptions('value', 'filter', v),
})

// --- Min Value ---
const fixedMinValue = computed({
  get: () => block.value.options?.minValue?.default ?? 0,
  set: v => updateNestedOptions('minValue', 'default', v),
})

const minValueModuleID = computed({
  get: () => block.value.options?.minValue?.moduleID || null,
  set: v => {
    const next = v || ''
    updateNestedOptions('minValue', 'moduleID', next)
    if (next && !block.value.options?.minValue?.field) {
      updateNestedOptions('minValue', 'field', 'count')
    }
  },
})

const minValueField = computed({
  get: () => block.value.options?.minValue?.field || 'count',
  set: v => {
    updateNestedOptions('minValue', 'field', v || 'count')
    if (!v || v === 'count') updateNestedOptions('minValue', 'operation', '')
  },
})

const minValueOperation = computed({
  get: () => block.value.options?.minValue?.operation || '',
  set: v => updateNestedOptions('minValue', 'operation', v),
})

const minValueFilter = computed({
  get: () => block.value.options?.minValue?.filter || '',
  set: v => updateNestedOptions('minValue', 'filter', v),
})

// --- Max Value ---
const fixedMaxValue = computed({
  get: () => block.value.options?.maxValue?.default ?? 100,
  set: v => updateNestedOptions('maxValue', 'default', v),
})

const maxValueModuleID = computed({
  get: () => block.value.options?.maxValue?.moduleID || null,
  set: v => {
    const next = v || ''
    updateNestedOptions('maxValue', 'moduleID', next)
    if (next && !block.value.options?.maxValue?.field) {
      updateNestedOptions('maxValue', 'field', 'count')
    }
  },
})

const maxValueField = computed({
  get: () => block.value.options?.maxValue?.field || 'count',
  set: v => {
    updateNestedOptions('maxValue', 'field', v || 'count')
    if (!v || v === 'count') updateNestedOptions('maxValue', 'operation', '')
  },
})

const maxValueOperation = computed({
  get: () => block.value.options?.maxValue?.operation || '',
  set: v => updateNestedOptions('maxValue', 'operation', v),
})

const maxValueFilter = computed({
  get: () => block.value.options?.maxValue?.filter || '',
  set: v => updateNestedOptions('maxValue', 'filter', v),
})

// --- Display options ---
const showValue = computed({
  get: () => block.value.options?.display?.showValue !== false,
  set: v => updateDisplay('showValue', v),
})

const showRelative = computed({
  get: () => block.value.options?.display?.showRelative || false,
  set: v => updateDisplay('showRelative', v),
})

const showProgress = computed({
  get: () => block.value.options?.display?.showProgress || false,
  set: v => updateDisplay('showProgress', v),
})

const variant = computed({
  get: () => block.value.options?.display?.variant || 'primary',
  set: v => updateDisplay('variant', v),
})

// --- Thresholds ---
// Stored at options.display.thresholds — same location the block and type def read from
const thresholds = computed(() => block.value.options?.display?.thresholds || [])

function addThreshold() {
  updateDisplay('thresholds', [...thresholds.value, { value: 50, variant: 'warning' }])
}

function removeThreshold(i) {
  const updated = [...thresholds.value]
  updated.splice(i, 1)
  updateDisplay('thresholds', updated)
}

function updateThreshold(i, key, value) {
  const updated = [...thresholds.value]
  updated[i] = { ...updated[i], [key]: value }
  updateDisplay('thresholds', updated)
}
</script>
