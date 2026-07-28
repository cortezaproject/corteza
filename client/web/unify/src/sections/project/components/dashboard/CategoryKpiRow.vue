<template>
  <!-- Responsive row of stat cards; wraps to a single column on narrow screens. -->
  <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
    <div
      v-for="kpi in kpis"
      :key="kpi.labelKey"
      class="rounded-xl border border-surface bg-surface px-4 py-3 flex flex-col gap-1"
    >
      <div class="text-[11px] font-semibold uppercase tracking-wide text-muted-color">
        {{ $t(kpi.labelKey) }}
      </div>
      <div class="text-2xl font-semibold text-color leading-none">
        {{ kpi.value }}
      </div>

      <!-- Optional created-per-week pulse for the category — decorative
           reinforcement only (the number above stays the primary readout, so
           the strip is aria-hidden). Every tile shares the same single series
           (the category's overall weekly volume, not a per-KPI breakdown), so
           there's one accent colour and no legend. Absent/empty `spark` ⇒ no
           strip, tile renders exactly as before. -->
      <div v-if="spark && spark.length" class="h-8 flex items-end gap-0.5 mt-1" aria-hidden="true">
        <div
          v-for="(v, i) in spark"
          :key="i"
          v-tooltip.top="`${sparkLabels[i] || ''} · ${v}`"
          class="flex-1 rounded-sm opacity-60"
          :class="accent ? [accent, 'bg-current'] : 'bg-emphasis'"
          :style="{ height: barHeight(v) + '%' }"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

// Presentational KPI row. Parent maps its category config to { labelKey, value }.
const props = defineProps({
  kpis: {
    type: Array,
    required: true,
  },
  // Optional weekly created-per-week totals (number[]), shared by every tile
  // (the category's overall weekly pulse — tiles differ by number, not by
  // series). Omitted/empty ⇒ no sparkline strip is rendered.
  spark: {
    type: Array,
    default: () => [],
  },
  // Full range strings aligned to `spark` (e.g. "Jul 13 – 19, 2026"), used
  // for the bars' tooltips only.
  sparkLabels: {
    type: Array,
    default: () => [],
  },
  // Category accent for the bars — the badge icon's text classes (e.g.
  // 'text-purple-600 dark:text-purple-400'), applied via bg-current so the
  // pulse matches the page's identity colour exactly in both themes (same
  // move as Overview's card rail). One colour, no per-bar hue variation
  // (it's one series, not a comparison). NOT a chart-palette hex: those
  // deliberately diverge from the badge colours (privacy is teal there for
  // CVD separation in stacked charts).
  accent: {
    type: String,
    default: '',
  },
})

// Bars scale to the series max; a floor keeps zero (or near-zero) weeks
// visible as a thin baseline sliver instead of disappearing entirely. When
// every week is zero, every bar lands on that same floor — reading as an
// even, empty track rather than a spike.
const maxVal = computed(() => Math.max(0, ...props.spark))
const barHeight = v => (maxVal.value > 0 ? Math.max(8, Math.round((v / maxVal.value) * 100)) : 8)
</script>
