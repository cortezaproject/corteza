<template>
  <AutoComplete
    :model-value="modelValue"
    @update:model-value="onSelect"
    :suggestions="filteredOptions"
    :option-label="optionLabel"
    :placeholder="placeholder"
    :disabled="disabled"
    :loading="loading"
    class="w-full"
    dropdown
    :showClear="showClear"
    :complete-on-focus="completeOnFocus"
    @complete="onComplete"
  >
    <template #option="slotProps">
      <slot name="option" :option="slotProps.option">
        {{ getLabel(slotProps.option) }}
      </slot>
    </template>
  </AutoComplete>
</template>

<script setup>
import AutoComplete from 'primevue/autocomplete'
import { ref, watch } from 'vue'

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
  completeOnFocus: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue', 'search'])

// Internal filtered options that AutoComplete uses
const filteredOptions = ref([])

// Sync with parent options
watch(() => props.options, (newOptions) => {
  filteredOptions.value = [...newOptions]
}, { immediate: true })

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

function onComplete(event) {
  const query = event.query || ''
  // Always reassign for empty queries so AutoComplete sees a new reference and opens the panel
  if (!query) {
    filteredOptions.value = [...props.options]
  }
  emit('search', query)
}
</script>

