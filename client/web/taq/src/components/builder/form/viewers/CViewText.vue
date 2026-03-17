<template>
  <span
    class="truncate"
    :class="modelValue ? 'text-color-emphasis' : 'italic text-muted-color'"
    :title="String(modelValue || '')"
  >
    {{ modelValue || notSet }}
  </span>
</template>

<script setup>
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const notSet = t('builder.preview.notSet')

defineProps({
  modelValue: { type: [String, Number], default: null },
})
</script>
