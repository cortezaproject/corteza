<template>
  <div
    class="flex items-center gap-3 p-3 border border-surface rounded-lg transition-colors"
    :class="[
      dimWhenOff && !modelValue ? 'opacity-50' : '',
      disabled ? 'cursor-not-allowed' : 'cursor-pointer hover:bg-emphasis',
    ]"
    @click="!disabled && $emit('update:modelValue', !modelValue)"
  >
    <div class="flex flex-col gap-0.5 flex-1 min-w-0">
      <span class="font-medium text-primary text-sm">
        <slot name="label">{{ label }}</slot>
      </span>
      <small v-if="description || $slots.description" class="text-muted-color text-xs">
        <slot name="description">{{ description }}</slot>
      </small>
    </div>
    <ToggleSwitch
      :modelValue="modelValue"
      :disabled="disabled"
      class="shrink-0"
      @click.stop
      @update:modelValue="$emit('update:modelValue', $event)"
    />
  </div>
</template>

<script setup>
defineProps({
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
})

defineEmits(['update:modelValue'])
</script>
