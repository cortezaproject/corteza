<template>
  <Select
    :model-value="modelValue"
    @update:model-value="onSelect"
    :options="options"
    :option-label="optionLabel"
    :placeholder="placeholder"
    :disabled="disabled"
    :loading="loading"
    class="w-full"
    fluid
    :filter="!hideSearch"
    :showClear="showClear"
    @filter="onFilter"
  >
    <template #option="slotProps">
      <slot name="option" :option="slotProps.option">
        {{ getLabel(slotProps.option) }}
      </slot>
    </template>
  </Select>
</template>

<script setup>
const props = defineProps({
  modelValue: {
    type: [String, Number, Object],
    default: null,
  },
  options: {
    type: Array,
    default: () => [],
  },
  optionLabel: {
    type: [String, Function],
    default: 'label',
  },
  placeholder: {
    type: String,
    default: '',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  loading: {
    type: Boolean,
    default: false,
  },
  showClear: {
    type: Boolean,
    default: true,
  },
  hideSearch: {
    type: Boolean,
    default: true,
  },
})

const emit = defineEmits(['update:modelValue', 'search'])

function getLabel(option) {
  if (!option) return ''
  if (typeof props.optionLabel === 'function') {
    return props.optionLabel(option)
  }
  return option[props.optionLabel] || ''
}

function onSelect(value) {
  emit('update:modelValue', value)
}

function onFilter(event) {
  emit('search', event.value || '')
}
</script>
