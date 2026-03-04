<template>
  <div class="flex flex-col gap-6">
    <!-- Input type (date+time / date only / time only) -->
    <div>
      <label class="font-medium text-surface-500 text-sm block mb-3">
        {{ $t('field.kind.dateTime.type.label') }}
      </label>
      <SelectButton
        :model-value="inputType"
        :options="inputTypeOptions"
        option-label="label"
        option-value="value"
        @update:model-value="onInputTypeChange"
      />
    </div>

    <!-- Constraints (not shown when time-only) -->
    <div v-if="!field.options.onlyTime">
      <label class="font-medium text-surface-500 text-sm block mb-3">
        {{ $t('field.kind.dateTime.constraints.label') }}
      </label>
      <div class="flex flex-col gap-2">
        <div
          v-for="opt in constraintOptions"
          :key="opt.value"
          class="flex items-center gap-2"
        >
          <RadioButton
            :input-id="`constraint-${opt.value}`"
            :model-value="constraintType"
            :value="opt.value"
            @update:model-value="onConstraintChange"
          />
          <label :for="`constraint-${opt.value}`" class="cursor-pointer">{{ opt.label }}</label>
        </div>
      </div>
    </div>

    <!-- Output format -->
    <div>
      <label class="font-medium text-surface-500 text-sm block mb-3">
        {{ $t('field.kind.dateTime.outputFormat') }}
      </label>
      <div class="flex items-center gap-2 mb-3">
        <Checkbox v-model="field.options.outputRelative" inputId="outputRelative" :binary="true" />
        <label for="outputRelative" class="cursor-pointer">{{ $t('field.kind.dateTime.relativeOutput') }}</label>
      </div>
      <template v-if="!field.options.outputRelative">
        <InputText
          v-model="field.options.format"
          placeholder="YYYY-MM-DD HH:mm"
          class="w-full"
        />
        <small class="text-muted-color mt-1 block">{{ $t('field.kind.dateTime.outputFormatFootnote', ['Moment.js']) }}</small>
      </template>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
})

const inputTypeOptions = computed(() => [
  { value: 'dateTime', label: t('field.kind.dateTime.type.options.dateTime') },
  { value: 'date', label: t('field.kind.dateTime.type.options.date') },
  { value: 'time', label: t('field.kind.dateTime.type.options.time') },
])

const constraintOptions = computed(() => [
  { value: 'all', label: t('field.kind.dateTime.constraints.options.all') },
  { value: 'pastValuesOnly', label: t('field.kind.dateTime.constraints.options.pastValuesOnly') },
  { value: 'futureValuesOnly', label: t('field.kind.dateTime.constraints.options.futureValuesOnly') },
])

const inputType = computed(() => {
  if (props.field.options.onlyDate) return 'date'
  if (props.field.options.onlyTime) return 'time'
  return 'dateTime'
})

const constraintType = computed(() => {
  if (props.field.options.onlyPastValues) return 'pastValuesOnly'
  if (props.field.options.onlyFutureValues) return 'futureValuesOnly'
  return 'all'
})

function onInputTypeChange(v) {
  props.field.options.onlyDate = v === 'date'
  props.field.options.onlyTime = v === 'time'
}

function onConstraintChange(v) {
  props.field.options.onlyPastValues = v === 'pastValuesOnly'
  props.field.options.onlyFutureValues = v === 'futureValuesOnly'
}
</script>
