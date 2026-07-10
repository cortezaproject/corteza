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

// Pip count + fill color per risk level.
const MAP = {
  Critical: { n: 5, c: 'bg-red-500' },
  High: { n: 4, c: 'bg-orange-500' },
  Medium: { n: 3, c: 'bg-amber-500' },
  Low: { n: 2, c: 'bg-emerald-500' },
}
const count = computed(() => MAP[props.level]?.n ?? 0)
const color = computed(() => MAP[props.level]?.c ?? 'bg-surface-400')
</script>
