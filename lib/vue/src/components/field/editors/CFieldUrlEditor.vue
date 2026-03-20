<template>
  <InputText
    :model-value="modelValue"
    :disabled="disabled"
    type="url"
    class="w-full"
    @update:model-value="$emit('update:modelValue', $event)"
    @blur="onBlur"
  />
</template>

<script setup>
import { trimUrlFragment, trimUrlQuery, trimUrlPath, onlySecureUrl } from '../url'

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

function onBlur() {
  let value = props.modelValue
  if (!value) return

  const opts = props.field.options || {}

  try {
    if (opts.trimFragment) value = trimUrlFragment(value)
    if (opts.trimQuery) value = trimUrlQuery(value)
    if (opts.trimPath) value = trimUrlPath(value)
    if (opts.onlySecure) value = onlySecureUrl(value)
  } catch {
    // Invalid URL, leave as-is
    return
  }

  if (value !== props.modelValue) {
    emit('update:modelValue', value)
  }
}
</script>
