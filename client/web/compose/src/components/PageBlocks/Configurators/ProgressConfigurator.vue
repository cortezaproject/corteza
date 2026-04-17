<template>
  <div class="flex flex-col gap-5">
    <!-- Value Section -->
    <div class="flex flex-col gap-3">
      <h5 class="text-lg font-semibold text-primary m-0">{{ $t('block.progress.value.label') }}</h5>

      <div class="flex flex-col gap-2">
        <div class="flex items-center gap-2">
          <Checkbox v-model="useModuleValue" binary input-id="useModuleValue" />
          <label for="useModuleValue" class="text-sm">{{ $t('block.progress.value.useModule') }}</label>
        </div>

        <template v-if="useModuleValue">
          <Select
            v-model="valueModuleID"
            :options="modules"
            option-label="name"
            option-value="moduleID"
            :placeholder="$t('block.progress.module.select')"
            class="w-full"
            filter
            show-clear
          />
          <div class="grid grid-cols-1 md:grid-cols-2 gap-2">
            <Select
              v-model="valueField"
              :options="getNumberFields(valueModuleID)"
              option-label="label"
              option-value="name"
              :placeholder="$t('block.progress.value.fieldPlaceholder')"
              class="w-full"
              :disabled="!valueModuleID"
              show-clear
            />
            <Select
              v-model="valueOperation"
              :options="operationOptions"
              option-label="label"
              option-value="value"
              :placeholder="$t('block.progress.value.operationPlaceholder')"
              class="w-full"
            />
          </div>
          <Textarea
            v-model="valueFilter"
            :placeholder="$t('block.progress.value.filterPlaceholder')"
            rows="2"
            class="w-full"
          />
        </template>

        <template v-else>
          <InputNumber
            v-model="fixedValue"
            :placeholder="$t('block.progress.value.fixed')"
            :min="0"
            class="w-full"
          />
        </template>
      </div>
    </div>

    <Divider />

    <!-- Min Value Section -->
    <div class="flex flex-col gap-3">
      <h5 class="text-lg font-semibold text-primary m-0">{{ $t('block.progress.minValue.label') }}</h5>

      <div class="flex flex-col gap-2">
        <div class="flex items-center gap-2">
          <Checkbox v-model="useModuleMinValue" binary input-id="useModuleMinValue" />
          <label for="useModuleMinValue" class="text-sm">{{ $t('block.progress.value.useModule') }}</label>
        </div>

        <template v-if="useModuleMinValue">
          <Select
            v-model="minValueModuleID"
            :options="modules"
            option-label="name"
            option-value="moduleID"
            :placeholder="$t('block.progress.module.select')"
            class="w-full"
            filter
            show-clear
          />
          <div class="grid grid-cols-1 md:grid-cols-2 gap-2">
            <Select
              v-model="minValueField"
              :options="getNumberFields(minValueModuleID)"
              option-label="label"
              option-value="name"
              :placeholder="$t('block.progress.value.fieldPlaceholder')"
              class="w-full"
              :disabled="!minValueModuleID"
              show-clear
            />
            <Select
              v-model="minValueOperation"
              :options="operationOptions"
              option-label="label"
              option-value="value"
              :placeholder="$t('block.progress.value.operationPlaceholder')"
              class="w-full"
            />
          </div>
          <Textarea
            v-model="minValueFilter"
            :placeholder="$t('block.progress.value.filterPlaceholder')"
            rows="2"
            class="w-full"
          />
        </template>

        <template v-else>
          <InputNumber
            v-model="fixedMinValue"
            :placeholder="$t('block.progress.minValue.fixed')"
            :min="0"
            class="w-full"
          />
        </template>
      </div>
    </div>

    <Divider />

    <!-- Max Value Section -->
    <div class="flex flex-col gap-3">
      <h5 class="text-lg font-semibold text-primary m-0">{{ $t('block.progress.maxValue.label') }}</h5>

      <div class="flex flex-col gap-2">
        <div class="flex items-center gap-2">
          <Checkbox v-model="useModuleMaxValue" binary input-id="useModuleMaxValue" />
          <label for="useModuleMaxValue" class="text-sm">{{ $t('block.progress.value.useModule') }}</label>
        </div>

        <template v-if="useModuleMaxValue">
          <Select
            v-model="maxValueModuleID"
            :options="modules"
            option-label="name"
            option-value="moduleID"
            :placeholder="$t('block.progress.module.select')"
            class="w-full"
            filter
            show-clear
          />
          <div class="grid grid-cols-1 md:grid-cols-2 gap-2">
            <Select
              v-model="maxValueField"
              :options="getNumberFields(maxValueModuleID)"
              option-label="label"
              option-value="name"
              :placeholder="$t('block.progress.value.fieldPlaceholder')"
              class="w-full"
              :disabled="!maxValueModuleID"
              show-clear
            />
            <Select
              v-model="maxValueOperation"
              :options="operationOptions"
              option-label="label"
              option-value="value"
              :placeholder="$t('block.progress.value.operationPlaceholder')"
              class="w-full"
            />
          </div>
          <Textarea
            v-model="maxValueFilter"
            :placeholder="$t('block.progress.value.filterPlaceholder')"
            rows="2"
            class="w-full"
          />
        </template>

        <template v-else>
          <InputNumber
            v-model="fixedMaxValue"
            :placeholder="$t('block.progress.maxValue.fixed')"
            :min="0"
            class="w-full"
          />
        </template>
      </div>
    </div>

    <Divider />

    <!-- Display Options -->
    <div class="flex flex-col gap-3">
      <h5 class="text-lg font-semibold text-primary m-0">{{ $t('block.progress.display-options') }}</h5>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div class="flex items-center gap-2">
          <Checkbox v-model="showValue" binary input-id="showValue" />
          <label for="showValue" class="text-sm">{{ $t('block.progress.show.value') }}</label>
        </div>

        <div class="flex items-center gap-2">
          <Checkbox v-model="showRelative" binary input-id="showRelative" />
          <label for="showRelative" class="text-sm">{{ $t('block.progress.show.relative') }}</label>
        </div>

        <div class="flex items-center gap-2">
          <Checkbox v-model="showProgress" binary input-id="showProgress" />
          <label for="showProgress" class="text-sm">{{ $t('block.progress.show.progress') }}</label>
        </div>

        <div class="flex items-center gap-2">
          <Checkbox v-model="animated" binary input-id="animated" />
          <label for="animated" class="text-sm">{{ $t('block.progress.animated') }}</label>
        </div>

        <div class="flex items-center gap-2">
          <Checkbox v-model="striped" binary input-id="striped" />
          <label for="striped" class="text-sm">{{ $t('block.progress.striped') }}</label>
        </div>
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-primary font-medium text-sm">{{ $t('block.progress.default-variant') }}</label>
        <Select
          v-model="variant"
          :options="variantOptions"
          option-label="label"
          option-value="value"
          class="w-full"
        />
      </div>
    </div>

    <Divider />

    <!-- Thresholds -->
    <div class="flex flex-col gap-3">
      <div class="flex items-center justify-between">
        <h5 class="text-lg font-semibold text-primary m-0">{{ $t('block.progress.thresholds.label') }}</h5>
        <Button
          :label="$t('general.label.add')"
          icon="pi pi-plus"
          size="small"
          severity="secondary"
          @click="addThreshold"
        />
      </div>

      <div
        v-for="(threshold, i) in thresholds"
        :key="i"
        class="flex items-center gap-2"
      >
        <InputNumber
          :model-value="threshold.value"
          :min="0"
          :max="100"
          class="flex-1"
          @update:model-value="updateThreshold(i, 'value', $event)"
        />
        <Select
          :model-value="threshold.variant"
          :options="variantOptions"
          option-label="label"
          option-value="value"
          class="flex-1"
          @update:model-value="updateThreshold(i, 'variant', $event)"
        />
        <Button
          icon="pi pi-trash"
          severity="danger"
          text
          size="small"
          @click="removeThreshold(i)"
        />
      </div>

      <small v-if="!thresholds.length" class="text-muted-color">
        {{ $t('block.progress.thresholds.empty') }}
      </small>
    </div>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import { useModuleStore } from '@/stores/module'

const { t } = useI18n()
const moduleStore = useModuleStore()

defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const block = inject('blockDraft')

const modules = computed(() => moduleStore.set || [])

const operationOptions = [
  { value: 'sum', label: t('block.progress.operation.sum') },
  { value: 'max', label: t('block.progress.operation.max') },
  { value: 'min', label: t('block.progress.operation.min') },
  { value: 'avg', label: t('block.progress.operation.avg') },
  { value: 'count', label: t('block.progress.operation.count') },
]

const variantOptions = [
  { value: 'primary', label: t('block.progress.variant.primary') },
  { value: 'secondary', label: t('block.progress.variant.secondary') },
  { value: 'success', label: t('block.progress.variant.success') },
  { value: 'warning', label: t('block.progress.variant.warning') },
  { value: 'danger', label: t('block.progress.variant.danger') },
  { value: 'info', label: t('block.progress.variant.info') },
  { value: 'light', label: t('block.progress.variant.light') },
  { value: 'dark', label: t('block.progress.variant.dark') },
]

function getNumberFields(moduleID) {
  if (!moduleID) return []
  const mod = moduleStore.getByID(moduleID)
  if (!mod) return []
  return (mod.fields || [])
    .filter(f => f.kind === 'Number')
    .map(f => ({ name: f.name, label: f.label || f.name }))
    .sort((a, b) => a.label.localeCompare(b.label))
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
const useModuleValue = computed({
  get: () => !!(block.value.options?.value?.moduleID),
  set: v => {
    if (!v) updateOptions('value', { ...block.value.options?.value, moduleID: '', field: '', operation: '', filter: '' })
  },
})

const fixedValue = computed({
  get: () => block.value.options?.value?.default ?? 0,
  set: v => updateNestedOptions('value', 'default', v),
})

const valueModuleID = computed({
  get: () => block.value.options?.value?.moduleID || '',
  set: v => updateNestedOptions('value', 'moduleID', v),
})

const valueField = computed({
  get: () => block.value.options?.value?.field || '',
  set: v => updateNestedOptions('value', 'field', v),
})

const valueOperation = computed({
  get: () => block.value.options?.value?.operation || 'count',
  set: v => updateNestedOptions('value', 'operation', v),
})

const valueFilter = computed({
  get: () => block.value.options?.value?.filter || '',
  set: v => updateNestedOptions('value', 'filter', v),
})

// --- Min Value ---
const useModuleMinValue = computed({
  get: () => !!(block.value.options?.minValue?.moduleID),
  set: v => {
    if (!v) updateOptions('minValue', { ...block.value.options?.minValue, moduleID: '', field: '', operation: '', filter: '' })
  },
})

const fixedMinValue = computed({
  get: () => block.value.options?.minValue?.default ?? 0,
  set: v => updateNestedOptions('minValue', 'default', v),
})

const minValueModuleID = computed({
  get: () => block.value.options?.minValue?.moduleID || '',
  set: v => updateNestedOptions('minValue', 'moduleID', v),
})

const minValueField = computed({
  get: () => block.value.options?.minValue?.field || '',
  set: v => updateNestedOptions('minValue', 'field', v),
})

const minValueOperation = computed({
  get: () => block.value.options?.minValue?.operation || 'count',
  set: v => updateNestedOptions('minValue', 'operation', v),
})

const minValueFilter = computed({
  get: () => block.value.options?.minValue?.filter || '',
  set: v => updateNestedOptions('minValue', 'filter', v),
})

// --- Max Value ---
const useModuleMaxValue = computed({
  get: () => !!(block.value.options?.maxValue?.moduleID),
  set: v => {
    if (!v) updateOptions('maxValue', { ...block.value.options?.maxValue, moduleID: '', field: '', operation: '', filter: '' })
  },
})

const fixedMaxValue = computed({
  get: () => block.value.options?.maxValue?.default ?? 100,
  set: v => updateNestedOptions('maxValue', 'default', v),
})

const maxValueModuleID = computed({
  get: () => block.value.options?.maxValue?.moduleID || '',
  set: v => updateNestedOptions('maxValue', 'moduleID', v),
})

const maxValueField = computed({
  get: () => block.value.options?.maxValue?.field || '',
  set: v => updateNestedOptions('maxValue', 'field', v),
})

const maxValueOperation = computed({
  get: () => block.value.options?.maxValue?.operation || 'count',
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

const animated = computed({
  get: () => block.value.options?.animated || false,
  set: v => updateOptions('animated', v),
})

const striped = computed({
  get: () => block.value.options?.striped || false,
  set: v => updateOptions('striped', v),
})

const variant = computed({
  get: () => block.value.options?.display?.variant || 'primary',
  set: v => updateDisplay('variant', v),
})

// --- Thresholds ---
const thresholds = computed(() => block.value.options?.thresholds || [])

function addThreshold() {
  updateOptions('thresholds', [...thresholds.value, { value: 50, variant: 'warning' }])
}

function removeThreshold(i) {
  const updated = [...thresholds.value]
  updated.splice(i, 1)
  updateOptions('thresholds', updated)
}

function updateThreshold(i, key, value) {
  const updated = [...thresholds.value]
  updated[i] = { ...updated[i], [key]: value }
  updateOptions('thresholds', updated)
}
</script>
