<template>
  <div class="flex flex-col gap-1 w-full">
    <InputText
      :id="id"
      :model-value="modelValue"
      @update:model-value="onInput"
      :class="['w-full', inputClass]"
      v-bind="$attrs"
      :invalid="!isValid"
    />
    <small v-if="!isValid" class="text-red-500">
      {{ errorMessage || $t('placeholder.handle') }}
    </small>
  </div>
</template>

<script setup>
import { computed, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  id: {
    type: String,
    default: undefined,
  },
  inputClass: {
    type: [String, Object, Array],
    default: '',
  },
  errorMessage: {
    type: String,
    default: '',
  },
  // Optionally allow overriding the regex
  pattern: {
    type: RegExp,
    default: () => /^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$|^[A-Za-z]$/,
  },
})

const emit = defineEmits(['update:modelValue', 'update:valid'])

const isValid = computed(() => {
  if (!props.modelValue) return true // if empty, validity depends on 'required' standard validation
  return props.pattern.test(props.modelValue)
})

watch(isValid, v => emit('update:valid', v), { immediate: true })

const onInput = val => {
  emit('update:modelValue', val)
}
</script>
