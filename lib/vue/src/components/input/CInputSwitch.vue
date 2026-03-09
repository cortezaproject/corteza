<template>
  <div class="flex flex-col gap-1">
    <label v-if="label" class="font-medium text-sm">
      {{ label }}
    </label>
    <span v-if="description" class="text-xs text-surface-500">
      {{ description }}
    </span>
    <div class="flex items-center gap-2">
      <span class="text-sm text-muted-color">{{ computedNoLabel }}</span>
      <ToggleSwitch
        :model-value="modelValue"
        :disabled="disabled"
        @update:model-value="$emit('update:modelValue', $event)"
      />
      <span class="text-sm text-muted-color">{{ computedYesLabel }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import ToggleSwitch from 'primevue/toggleswitch'

const { t } = useI18n()

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
  noLabel: {
    type: String,
    default: '',
  },
  yesLabel: {
    type: String,
    default: '',
  },
})

defineEmits(['update:modelValue'])

const computedNoLabel = computed(() => props.noLabel || t('general.label.general.no', 'No'))
const computedYesLabel = computed(() => props.yesLabel || t('general.label.general.yes', 'Yes'))
</script>
