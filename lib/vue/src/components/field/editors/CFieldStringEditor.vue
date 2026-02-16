<template>
  <Textarea
    v-if="field.options?.multiLine"
    :model-value="modelValue"
    :disabled="disabled"
    :maxlength="field.options?.maxLength ?? undefined"
    rows="3"
    class="w-full"
    @update:model-value="$emit('update:modelValue', $event)"
  />
  <InputText
    v-else
    :model-value="modelValue"
    :disabled="disabled"
    :maxlength="field.options?.maxLength ?? undefined"
    class="w-full"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script setup>
defineProps({
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

defineEmits(['update:modelValue'])
</script>
