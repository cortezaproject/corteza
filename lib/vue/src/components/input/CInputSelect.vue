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
watch(
  () => props.options,
  newOptions => {
    filteredOptions.value = [...newOptions]
  },
  { immediate: true },
)

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
  // PrimeVue AutoComplete has a `searching` flag that is set to true in search(),
  // then consumed (set to false) by the suggestions watcher on the FIRST suggestions change.
  // The watcher gates this.show() behind `searching === true`.
  //
  // For sync options (e.g. Select with static data): props.options already has data,
  // so assign filteredOptions immediately to trigger the watcher and open the panel.
  //
  // For async options (e.g. Namespace/Module selectors): props.options is empty here.
  // Don't assign — it would consume the searching flag. Instead, let the parent fetch
  // and update options; the watch on props.options will set filteredOptions, which triggers
  // the AutoComplete watcher while searching is still true.
  if (props.options.length > 0) {
    filteredOptions.value = query
      ? props.options.filter(o => {
          const label = typeof props.optionLabel === 'function' ? props.optionLabel(o) : (o[props.optionLabel] || '')
          return String(label).toLowerCase().includes(query.toLowerCase())
        })
      : [...props.options]
  }
  emit('search', query)
}
</script>
