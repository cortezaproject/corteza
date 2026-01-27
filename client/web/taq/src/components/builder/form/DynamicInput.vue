<template>
  <div class="flex flex-col gap-1">
    <label v-if="label" class="text-sm font-medium text-color">
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </label>
    <component
      :is="inputComponent"
      :model-value="modelValue"
      @update:model-value="$emit('update:modelValue', $event)"
      :placeholder="placeholder"
      :disabled="disabled"
      v-bind="visual"
    />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { resolveInputComponent } from '@cortezaproject/corteza-vue-next/src/components/input/registry'

const props = defineProps({
  modelValue: {
    type: [String, Number, Object, Array],
    default: null,
  },
  type: {
    type: String,
    default: 'Text',
  },
  label: {
    type: String,
    default: '',
  },
  placeholder: {
    type: String,
    default: '',
  },
  required: {
    type: Boolean,
    default: false,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  visual: {
    type: Object,
    default: () => ({}),
  },
})

defineEmits(['update:modelValue'])

const inputComponent = computed(() => resolveInputComponent(props.type))
</script>
