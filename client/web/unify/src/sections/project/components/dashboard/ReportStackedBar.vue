<template>
  <div :class="bare ? '' : 'rounded-lg border border-surface bg-surface p-4'">
    <div v-if="!bare && titleKey" class="flex items-center gap-2 mb-2">
      <div class="text-sm font-medium text-color truncate">{{ $t(titleKey) }}</div>
    </div>

    <!-- Empty state when every row totals zero. -->
    <div
      v-if="empty"
      class="flex items-center justify-center text-muted-color text-sm"
      :style="{ minHeight: height + 'px' }"
    >
      —
    </div>

    <!-- One horizontal bar per row: label left, stacked segments middle, a
         caller-supplied right-side label (total count, or a direct
         "N open · M overdue" string). `scale` picks how a row's segments size:
         'share' sizes each row to its OWN total (a 100%-stacked composition —
         used for severity-by-category, where the point is the mix within a
         category, not its volume vs others); 'value' sizes every row against
         the shared max across all rows (a leaderboard, like CategoryRankBar —
         used for open/overdue, where bar length itself is meant to compare
         volumes). -->
    <div v-else class="flex flex-col justify-center gap-2" :style="{ minHeight: height + 'px' }">
      <div v-for="row in rows" :key="row.key" class="flex items-center gap-2">
        <span class="w-24 shrink-0 text-xs text-muted-color truncate">{{ row.label }}</span>
        <span class="flex-1 h-4 rounded-full bg-emphasis overflow-hidden flex">
          <span
            v-for="(seg, i) in row.segments"
            :key="`${seg.label}-${i}`"
            class="h-full"
            :style="segStyle(row, seg, i)"
            :title="seg.title"
          />
        </span>
        <span class="w-32 shrink-0 text-right text-xs text-color tabular-nums">{{ row.rightText }}</span>
      </div>

      <!-- Legend — segment identity must never be color+tooltip alone. -->
      <div v-if="legend.length" class="flex flex-wrap items-center gap-x-3 gap-y-1 pt-1">
        <span v-for="item in legend" :key="item.label" class="inline-flex items-center gap-1.5">
          <span class="w-2 h-2 rounded-full shrink-0" :style="{ background: item.color }" />
          <span class="text-xs text-muted-color">{{ item.label }}</span>
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t: $t } = useI18n()

const props = defineProps({
  titleKey: { type: String, default: '' },
  // [{ key, label, total, rightText, segments: [{ label, value, color, title }] }]
  rows: { type: Array, default: () => [] },
  // 'share' (each row scales to its own total) | 'value' (rows scale to the
  // shared max total, for cross-row comparison).
  scale: { type: String, default: 'share' },
  // [{ label, color }] — rendered as a dot legend under the bars. Required
  // whenever segments carry more than one series.
  legend: { type: Array, default: () => [] },
  // Chrome-less variant: no card wrapper/title, for embedding in a band.
  bare: { type: Boolean, default: false },
  height: { type: Number, default: 200 },
})

const empty = computed(() => !props.rows.some(r => (r.total || 0) > 0))

// Shared denominator for 'value' scale — same "scale to the largest row"
// leaderboard idiom as CategoryRankBar, so bar length is directly comparable
// across categories.
const maxTotal = computed(() => Math.max(1, ...props.rows.map(r => r.total || 0)))

// The 2px gap between stacked segments is a border in the surrounding surface
// colour (a live CSS var — plain DOM, unlike the echarts canvas components,
// can read it directly) — same "surface gap, not a stroke" idiom as
// CategoryTrendChart/CategoryDonutChart's segment borders.
function segStyle(row, seg, i) {
  const denom = props.scale === 'value' ? maxTotal.value : row.total || 1
  const width = denom ? ((seg.value || 0) / denom) * 100 : 0
  const isLast = i === row.segments.length - 1
  return {
    width: `${width}%`,
    background: seg.color,
    borderRight: isLast || !seg.value ? 'none' : '2px solid var(--p-content-background)',
  }
}
</script>
