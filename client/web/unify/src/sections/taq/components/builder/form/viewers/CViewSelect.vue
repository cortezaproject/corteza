<template>
  <span
    class="truncate"
    :class="
      modelValue != null && modelValue !== '' ? 'text-color-emphasis' : 'italic text-muted-color'
    "
    :title="displayLabel"
  >
    {{ displayLabel }}
  </span>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  modelValue: { type: [String, Number], default: null },
  options: { type: Array, default: () => [] },
})

const displayLabel = computed(() => {
  if (props.modelValue == null || props.modelValue === '') {
    return t('builder.preview.notSet')
  }

  // Match value to option label
  const match = props.options.find(o => o.value === props.modelValue)
  if (match) return match.label || String(props.modelValue)

  return String(props.modelValue)
})
</script>
