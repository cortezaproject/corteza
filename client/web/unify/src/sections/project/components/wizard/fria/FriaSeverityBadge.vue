<template>
  <span
    v-if="severity"
    class="inline-flex items-center rounded-full px-2.5 py-1 text-xs font-semibold uppercase tracking-wide shrink-0"
    :class="FRIA_SEVERITY_BADGE_CLASSES[severity] || 'bg-emphasis text-muted-color'"
  >
    {{ $t(`fria.severity.${severity}.label`) }}
  </span>
  <span
    v-else
    class="inline-flex items-center rounded-full px-2.5 py-1 text-xs font-medium text-muted-color bg-emphasis shrink-0"
  >
    {{ $t('fria.severity.unset') }}
  </span>
</template>

<script setup>
import { FRIA_SEVERITY_BADGE_CLASSES } from '@/sections/project/config/friaScenario'

defineProps({
  severity: { type: String, default: null },
})
</script>
