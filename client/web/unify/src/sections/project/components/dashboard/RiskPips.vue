<template>
  <!-- Risk as a 5-pip meter, filled + colored by level (mirrors the demo's
       risk-indicator). Empty when the level is unknown/blank. -->
  <span class="inline-flex items-center gap-[3px]" :title="level || ''">
    <span
      v-for="i in 5"
      :key="i"
      class="w-2 h-2 rounded-sm"
      :class="i <= count ? color : 'bg-surface-300 dark:bg-surface-600'"
    />
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({ level: { type: String, default: '' } })

// Pip count + fill color per risk level. Families track config/chartColors's
// RISK_COLORS ramp (same red→orange→amber→yellow family as EventBadge's
// severity tints — the two ordinal scales read as one system). Very Low
// shares Low's yellow — one filled pip vs two carries the distinction, since
// the family ramp has no lighter legible step left. None renders zero filled
// pips, visually the same empty track as unknown/blank, but the tooltip
// (:title) still names the level.
const MAP = {
  Critical: { n: 5, c: 'bg-red-500' },
  High: { n: 4, c: 'bg-orange-500' },
  Medium: { n: 3, c: 'bg-amber-500' },
  Low: { n: 2, c: 'bg-yellow-500' },
  'Very Low': { n: 1, c: 'bg-yellow-500' },
  None: { n: 0, c: 'bg-surface-400' },
}
const count = computed(() => MAP[props.level]?.n ?? 0)
const color = computed(() => MAP[props.level]?.c ?? 'bg-surface-400')
</script>
