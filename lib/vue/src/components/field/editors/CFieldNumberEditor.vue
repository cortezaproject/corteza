<template>
  <InputNumber
    :model-value="displayValue"
    :disabled="disabled"
    :min-fraction-digits="0"
    :max-fraction-digits="field.options?.precision ?? 3"
    :min="field.options?.display === 'progress' ? field.options?.min : undefined"
    :max="field.options?.display === 'progress' ? field.options?.max : undefined"
    :step="field.options?.step ?? 1"
    :prefix="field.options?.prefix"
    :suffix="field.options?.suffix"
    class="w-full"
    @focus="onFocus"
    @blur="onBlur"
    @input="onInput"
    @update:model-value="onUpdate"
  />
</template>

<script setup>
import { computed, ref } from 'vue'

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

// While focused, the model-value prop is frozen so per-keystroke emits don't
// write back into InputNumber and reformat the text mid-typing (e.g. "5." → "5").
const focused = ref(false)
const frozenValue = ref(null)

const displayValue = computed(() => (focused.value ? frozenValue.value : numericValue.value))

function onFocus() {
  frozenValue.value = numericValue.value
  focused.value = true
}

function onBlur() {
  focused.value = false
}

// InputNumber only commits its model on blur/enter; emit per keystroke so the
// record (the validation source) stays current while typing.
function onInput({ value }) {
  emit('update:modelValue', value == null ? '' : String(value))
}

function onUpdate(value) {
  emit('update:modelValue', value == null ? '' : String(value))
}
</script>
