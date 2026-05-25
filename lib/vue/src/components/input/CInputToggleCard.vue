<template>
  <div
    class="flex gap-3 p-3 border border-surface rounded-lg transition-colors max-w-2xl"
    :class="[disabled ? 'cursor-not-allowed' : 'cursor-pointer hover:bg-emphasis']"
    @click="!disabled && $emit('update:modelValue', !modelValue)"
  >
    <div class="flex flex-col flex-1 min-w-0" :class="textOpacityClass">
      <div class="flex items-center gap-1">
        <span class="font-medium text-primary text-sm">
          <slot name="label">{{ label }}</slot>
        </span>
        <i
          v-if="warning"
          v-tooltip.top="warning"
          class="pi pi-exclamation-triangle text-orange-500 text-xs cursor-help"
          @click.stop
        />
        <i
          v-if="hint"
          v-tooltip.top="hint"
          class="pi pi-question-circle text-muted-color text-xs cursor-help"
          @click.stop
        />
      </div>
      <small v-if="description || $slots.description" class="text-muted-color">
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
  warning: {
    type: String,
    default: '',
  },
  hint: {
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
