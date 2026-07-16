<template>
  <div :class="bare ? '' : 'rounded-lg border border-surface bg-surface-0 dark:bg-surface-900 p-4'">
    <div v-if="!bare && titleKey" class="flex items-center gap-2 mb-2">
      <span v-if="accent" class="w-2 h-2 rounded-full shrink-0" :style="{ background: accent }" />
      <div class="text-sm font-medium text-color truncate">{{ $t(titleKey) }}</div>
    </div>

    <!-- Empty state when there is nothing to plot -->
    <div
      v-if="!total"
      class="flex items-center justify-center text-muted-color text-sm"
      :style="{ height: height + 'px' }"
    >
      —
    </div>

    <!-- Ordinal breakdown as ranked bars (worst on top — `data` already
         arrives in canonical rank order, see CategoryView#breakdownFor), not
         a donut: order carries meaning here, so position does the identity
         work a legend would otherwise have to. -->
    <div v-else class="flex flex-col justify-center gap-0.5" :style="{ minHeight: height + 'px' }">
      <div
        v-for="row in data"
        :key="row.label"
        class="flex items-center gap-2"
        :title="`${row.label}: ${row.value} (${share(row.value)}%)`"
      >
        <span class="w-20 shrink-0 text-xs text-muted-color truncate">{{ row.label }}</span>
        <span class="flex-1 h-2 rounded-full bg-surface-100 dark:bg-surface-800 overflow-hidden">
          <span
            class="block h-full rounded-r"
            :style="{ width: widthPct(row.value) + '%', background: colorFor(variant, row.label) }"
          />
        </span>
        <span class="w-8 shrink-0 text-right text-xs text-color tabular-nums">{{ row.value }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { colorFor } from '@/sections/project/config/chartColors'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t: $t } = useI18n()

const props = defineProps({
  titleKey: { type: String, default: '' },
  // [{ label: string, value: number }], already sorted worst→best.
  data: { type: Array, default: () => [] },
  // Colour family for the bars: severity | risk.
  variant: { type: String, default: 'severity' },
  // Optional accent dot next to the title (hex).
  accent: { type: String, default: '' },
  // Chrome-less variant: no card wrapper/title, for embedding in a card.
  bare: { type: Boolean, default: false },
  height: { type: Number, default: 200 },
})

const total = computed(() => props.data.reduce((s, d) => s + (d.value || 0), 0))
// Bars scale to the largest row (a leaderboard), not to the total (that would
// read as a composition/share chart, which this isn't) — share is tooltip-only.
const max = computed(() => Math.max(1, ...props.data.map(d => d.value || 0)))

const share = value => (total.value ? Math.round((value / total.value) * 100) : 0)
const widthPct = value => Math.round((value / max.value) * 100)
</script>
