<template>
  <div v-if="buttons.length" :class="containerClass">
    <Button
      v-for="button in buttons"
      :key="button.script"
      v-tooltip.bottom="button.description"
      :label="button.label"
      :severity="mapVariant(button.variant || defaultVariant)"
      :size="size"
      @click.prevent="emit('click', button)"
    />
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'

const props = defineProps({
  // Resource the scripts are bound to, e.g. system:user
  resourceType: { type: [String, Array], required: true },
  // Trigger's ui page prop, e.g. user/editor
  uiPage: { type: String, default: undefined },
  // Trigger's ui slot prop, e.g. toolbar
  uiSlot: { type: String, default: undefined },
  app: { type: String, default: 'admin' },
  containerClass: { type: String, default: 'flex flex-wrap gap-2' },
  defaultVariant: { type: String, default: 'secondary' },
  size: { type: String, default: 'small' },
})

const emit = defineEmits(['click'])

const $UIHooks = inject('$UIHooks', null)

const variantSeverityMap = {
  primary: undefined,
  secondary: 'secondary',
  light: 'secondary',
  dark: 'contrast',
  success: 'success',
  danger: 'danger',
  warning: 'warn',
  info: 'info',
}

const mapVariant = key => variantSeverityMap[key]

const buttons = computed(() => {
  if (!$UIHooks) return []
  return $UIHooks.Find(props.resourceType, props.uiPage, props.uiSlot, props.app)
})
</script>
