<template>
  <!-- Multi-select mode: one component for all values -->
  <MultiSelect
    v-if="isMultipleType"
    :model-value="multiValue"
    :options="selectOptions"
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
import Tag from 'primevue/tag'
import MultiSelect from 'primevue/multiselect'

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
const isMultipleType = computed(() => props.field.isMulti && props.field.options?.selectType === 'multiple')
const selectOptions = computed(() => props.field.options?.options || [])

// Normalize array model value for MultiSelect
const multiValue = computed(() => {
  if (Array.isArray(props.modelValue)) return props.modelValue
  return props.modelValue ? [props.modelValue] : []
})

function toColor(hex) {
  if (!hex) return null
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
    color: toColor(optStyle.textColor) || undefined,
    backgroundColor: toColor(optStyle.backgroundColor) || undefined,
  }
}
</script>
