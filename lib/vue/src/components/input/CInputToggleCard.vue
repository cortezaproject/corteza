<template>
  <div
    class="flex items-center gap-3 p-3 border border-surface rounded-lg transition-colors"
    :class="[
      disabled ? 'cursor-not-allowed' : 'cursor-pointer hover:bg-emphasis',
    ]"
    @click="!disabled && $emit('update:modelValue', !modelValue)"
  >
    <div
      class="flex flex-col gap-0.5 flex-1 min-w-0"
      :class="textOpacityClass"
    >
      <span class="font-medium text-primary text-sm">
        <slot name="label">{{ label }}</slot>
      </span>
      <small v-if="description || $slots.description" class="text-muted-color text-xs">
        <slot name="description">{{ description }}</slot>
      </small>
    </div>
    <ToggleSwitch
      :modelValue="modelValue"
      class="shrink-0"
      :class="disabled ? 'opacity-50 pointer-events-none' : ''"
      @click.stop
      @update:modelValue="!disabled && $emit('update:modelValue', $event)"
    />
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  label: {
    type: String,
    default: '',
  },
  description: {
    type: String,
    default: '',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  dimWhenOff: {
    type: Boolean,
    default: false,
  },
  dim: {
    type: Boolean,
    default: false,
  },
})

defineEmits(['update:modelValue'])

const textOpacityClass = computed(() => {
  if (props.dim) return 'opacity-50'
  if (props.dimWhenOff && !props.modelValue) return 'opacity-50'
  return ''
})
</script>
