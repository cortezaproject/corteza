<template>
  <InputNumber
    :model-value="numericValue"
    :disabled="disabled"
    :min-fraction-digits="0"
    :max-fraction-digits="field.options?.precision ?? 3"
    :min="field.options?.display === 'progress' ? field.options?.min : undefined"
    :max="field.options?.display === 'progress' ? field.options?.max : undefined"
    :step="field.options?.step ?? 1"
    :prefix="field.options?.prefix"
    :suffix="field.options?.suffix"
    class="w-full"
    @input="onInput"
    @update:model-value="onUpdate"
  />
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
  modelValue: {
    type: String,
    default: '',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue'])

const numericValue = computed(() => {
  if (props.modelValue === '' || props.modelValue == null) return null
  const n = parseFloat(props.modelValue)
  return isNaN(n) ? null : n
})

// InputNumber only commits its model on blur/enter; emit per keystroke so the
// record (the validation source) stays current while typing.
//
// model-value must track those emits rather than staying frozen for the
// duration of the focus. InputNumber writes the typed text straight to the DOM
// input and leaves its own d_value behind; Vue re-applies the `value` binding on
// every patch of that input, so any re-render while the prop is stale wipes what
// was typed. In a record list, where each keystroke re-renders the row, that made
// the field impossible to type in at all. Feeding the emitted value back keeps
// the two in step, and InputNumber suppresses the emit when the parsed number is
// unchanged ("5" → "5."), so trailing decimal marks still survive.
function onInput({ value }) {
  emit('update:modelValue', value == null ? '' : String(value))
}

function onUpdate(value) {
  emit('update:modelValue', value == null ? '' : String(value))
}
</script>
