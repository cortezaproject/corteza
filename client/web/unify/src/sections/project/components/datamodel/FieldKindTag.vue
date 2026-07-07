<script setup>
import { fieldType, fieldTypeLabelKey } from '@/sections/project/config/fieldTypes'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

// A badge for a data-model field's kind, styled like the metrics-strip items: an
// outer bordered container holding a small rounded-square icon chip (the field
// kind's icon) followed by the localized type name. Kept neutral — no per-type
// accent colors. Unknown kinds fall back to the raw type text with no chip.
const props = defineProps({
  type: { type: String, required: true },
})

const { t } = useI18n()

const icon = computed(() => fieldType(props.type)?.icon)
const label = computed(() => {
  const key = fieldTypeLabelKey(props.type)
  return key ? t(key) : props.type
})
</script>

<template>
  <span
    class="shrink-0 inline-flex items-center gap-1.5 rounded-lg border border-surface px-1.5 py-1"
  >
    <span
      v-if="icon"
      class="inline-flex items-center justify-center w-5 h-5 rounded-md ring-1 ring-surface bg-surface-100 dark:bg-surface-800 shrink-0"
    >
      <i :class="[icon, 'text-muted-color text-[10px]']" />
    </span>
    <span class="text-xs text-muted-color leading-none whitespace-nowrap">{{ label }}</span>
  </span>
</template>
