<template>
  <!-- Checkbox group mode (multi + list selectType) -->
  <div v-if="isListType && field.isMulti" class="flex flex-col gap-2">
    <div v-for="option in selectOptions" :key="option.value" class="flex items-center gap-2">
      <Checkbox
        :input-id="`select-${field.fieldID}-${option.value}`"
        :model-value="multiValue"
        :value="option.value"
        :disabled="disabled"
        @update:model-value="$emit('update:modelValue', $event ?? [])"
      />
      <label :for="`select-${field.fieldID}-${option.value}`" class="cursor-pointer">
        <span
          :class="{ 'inline-flex': isBadgeDisplay }"
          :style="isBadgeDisplay ? badgeInlineStyle(option.value) : {}"
        >
          {{ option.text }}
        </span>
      </label>
    </div>
  </div>

  <!-- Radio button group mode (single + list selectType) -->
  <div v-else-if="isListType" class="flex flex-col gap-2">
    <div v-for="option in selectOptions" :key="option.value" class="flex items-center gap-2">
      <RadioButton
        :input-id="`select-${field.fieldID}-${option.value}`"
        :model-value="modelValue"
        :value="option.value"
        :disabled="disabled"
        @update:model-value="$emit('update:modelValue', $event ?? '')"
      />
      <label :for="`select-${field.fieldID}-${option.value}`" class="cursor-pointer">
        <span
          :class="{ 'inline-flex': isBadgeDisplay }"
          :style="isBadgeDisplay ? badgeInlineStyle(option.value) : {}"
        >
          {{ option.text }}
        </span>
      </label>
    </div>
  </div>

  <!-- Multi-select mode: one component for all values -->
  <MultiSelect
    v-else-if="isMultipleType"
    :model-value="multiValue"
    :options="availableOptions"
    option-label="text"
    option-value="value"
    :disabled="disabled"
    :show-clear="!field.isRequired"
    class="w-full"
    @update:model-value="$emit('update:modelValue', $event ?? [])"
  >
    <template v-if="isBadgeDisplay" #value="{ value: vals }">
      <div class="flex flex-wrap gap-1">
        <Tag
          v-for="val in vals"
          :key="val"
          :value="optionLabel(val)"
          :pt="{ root: { style: badgeStyle(val) } }"
        />
      </div>
    </template>
    <template v-if="isBadgeDisplay" #option="{ option }">
      <Tag :value="option.text" :pt="{ root: { style: badgeStyle(option.value) } }" />
    </template>
  </MultiSelect>

  <!-- Single-select mode (default, each) -->
  <Select
    v-else
    :model-value="modelValue"
    :options="selectOptions"
    option-label="text"
    option-value="value"
    :show-clear="!field.isRequired"
    :disabled="disabled"
    class="w-full"
    @update:model-value="$emit('update:modelValue', $event ?? '')"
  >
    <template v-if="isBadgeDisplay" #value="{ value: val }">
      <Tag v-if="val" :value="optionLabel(val)" :pt="{ root: { style: badgeStyle(val) } }" />
    </template>

    <template v-if="isBadgeDisplay" #option="{ option }">
      <Tag :value="option.text" :pt="{ root: { style: badgeStyle(option.value) } }" />
    </template>
  </Select>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
  modelValue: {
    type: [String, Array],
    default: '',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

defineEmits(['update:modelValue'])

const isBadgeDisplay = computed(() => props.field.options?.displayType === 'badge')
const isMultipleType = computed(
  () => props.field.isMulti && props.field.options?.selectType === 'multiple',
)
const isListType = computed(() => props.field.options?.selectType === 'list')
const selectOptions = computed(() =>
  (props.field.options?.options || []).filter(o => o.value && o.text),
)

// Filter out already-selected values when isUniqueMultiValue is set
const availableOptions = computed(() => {
  if (!props.field.isMulti || !props.field.options?.isUniqueMultiValue) return selectOptions.value
  return selectOptions.value.filter(o => !multiValue.value.includes(o.value))
})

// Normalize array model value for MultiSelect / Checkbox group
const multiValue = computed(() => {
  if (Array.isArray(props.modelValue)) return props.modelValue
  return props.modelValue ? [props.modelValue] : []
})

// Theme-independent defaults so badges are always visible
const DEFAULT_BADGE_TEXT_COLOR = '#FFFFFFFF'
const DEFAULT_BADGE_BG_COLOR = '#09344EFF'

function toColor(hex) {
  if (!hex || hex === 'transparent') return null
  return hex.startsWith('#') ? hex : '#' + hex
}

function optionLabel(val) {
  const opt = selectOptions.value.find(o => o.value === val)
  return opt?.text || val
}

function badgeStyle(val) {
  const opt = selectOptions.value.find(o => o.value === val)
  const optStyle = opt?.style || {}
  return {
    color: toColor(optStyle.textColor) || DEFAULT_BADGE_TEXT_COLOR,
    backgroundColor: toColor(optStyle.backgroundColor) || DEFAULT_BADGE_BG_COLOR,
  }
}

function badgeInlineStyle(val) {
  const opt = selectOptions.value.find(o => o.value === val)
  const optStyle = opt?.style || {}
  return {
    color: toColor(optStyle.textColor) || DEFAULT_BADGE_TEXT_COLOR,
    backgroundColor: toColor(optStyle.backgroundColor) || DEFAULT_BADGE_BG_COLOR,
    padding: '0.15rem 0.5rem',
    borderRadius: '999px',
    fontSize: '0.9rem',
  }
}
</script>
