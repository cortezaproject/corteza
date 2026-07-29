<template>
  <!-- Flat tinted pill used for the type / severity / status columns, matching
       the demo's badges. Severity carries a leading colored dot. -->
  <span
    class="inline-flex items-center gap-1.5 rounded-full font-medium whitespace-nowrap"
    :class="[cls, size === 'md' ? 'px-2.5 py-1 text-xs' : 'px-2 py-0.5 text-[11px]']"
    :style="tint"
  >
    <span
      v-if="dot || dotTint"
      class="rounded-full shrink-0"
      :class="[dot, size === 'md' ? 'w-2 h-2' : 'w-1.5 h-1.5']"
      :style="dotTint"
    />
    {{ value || '—' }}
  </span>
</template>

<script setup>
import { STATUS_COLORS } from '@/sections/project/config/chartColors'
import { computed } from 'vue'

const props = defineProps({
  value: { type: String, default: '' },
  // 'type' (neutral bordered) | 'severity' (dot + tint) | 'status' (tint) |
  // 'priority' (dot + tint) | 'risk' (dot + tint)
  variant: { type: String, default: 'type' },
  // 'sm' (default — table/list density) | 'md' (form selects, where the list
  // pill reads undersized against a ~2.5rem control).
  size: { type: String, default: 'sm' },
})

const MUTED = 'bg-emphasis text-muted-color'

// Tint families track config/chartColors's SEVERITY_COLORS ramp (Tailwind's
// closest family per hue — exact hex match isn't the goal, family coherence
// is, since these are labeled pills, not colour-only). The ramp reads
// red → orange → amber → yellow (worst → least severe warm), Informational
// stays sky (unchanged — it's the chart ramp's own cool outlier too).
const SEVERITY = {
  Critical: 'bg-red-100 text-red-700 dark:bg-red-500/15 dark:text-red-300',
  Serious: 'bg-orange-100 text-orange-700 dark:bg-orange-500/15 dark:text-orange-300',
  Major: 'bg-amber-100 text-amber-700 dark:bg-amber-500/15 dark:text-amber-300',
  Minor: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-500/15 dark:text-yellow-300',
  Informational: 'bg-sky-100 text-sky-700 dark:bg-sky-500/15 dark:text-sky-300',
}
const SEVERITY_DOT = {
  Critical: 'bg-red-500',
  Serious: 'bg-orange-500',
  Major: 'bg-amber-500',
  Minor: 'bg-yellow-500',
  Informational: 'bg-sky-500',
}
// Lifecycle ramp — the status variant is the one variant styled INLINE off
// STATUS_COLORS in config/chartColors, using the EXACT chart hexes (ruled
// 2026-07-29), not a nearest Tailwind family: status pills sit right next to
// the status donuts and the wizard's stacked progress bar, so "same hue
// family" still read as two different colours side by side. Text + dot carry
// the exact hex in both themes (each hex was palette-validated against both
// surfaces); the pill's tint is the same hex at 15% alpha, which lands close
// to the old bg-*-100 / dark:bg-*-500/15 pair on both surfaces.

// Backlog item priority (High/Medium/Low — no Critical tier, unlike
// severity/risk). Reuses RiskPips's family (red/orange/amber/yellow) but
// shifted one step warmer per tier since priority tops out at High rather
// than Critical: High reads as the urgent/red tier, Medium amber, Low
// yellow — same three-family ramp readers already know from severity/risk.
const PRIORITY = {
  High: 'bg-red-100 text-red-700 dark:bg-red-500/15 dark:text-red-300',
  Medium: 'bg-amber-100 text-amber-700 dark:bg-amber-500/15 dark:text-amber-300',
  Low: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-500/15 dark:text-yellow-300',
}
const PRIORITY_DOT = {
  High: 'bg-red-500',
  Medium: 'bg-amber-500',
  Low: 'bg-yellow-500',
}

// (No `risk` variant: risk keeps its 5-pip meter identity everywhere —
// RiskPips — with a text label alongside in selects/summaries.)

// The exact chart hex for a known status; null for every other variant (or
// an unknown status value, which keeps the MUTED fallback pill).
const statusHex = computed(() =>
  props.variant === 'status' ? STATUS_COLORS[props.value] || null : null,
)

const cls = computed(() => {
  if (props.variant === 'severity') return SEVERITY[props.value] || MUTED
  if (props.variant === 'status') return statusHex.value ? '' : MUTED
  if (props.variant === 'priority') return PRIORITY[props.value] || MUTED
  return 'border border-surface text-muted-color' // type: flat neutral pill
})
// '26' = 15% alpha as an 8-digit-hex suffix.
const tint = computed(() =>
  statusHex.value ? { color: statusHex.value, background: `${statusHex.value}26` } : null,
)
const dotTint = computed(() => (statusHex.value ? { background: statusHex.value } : null))

const dot = computed(() => {
  if (props.variant === 'severity') return SEVERITY_DOT[props.value] || 'bg-surface-400'
  if (props.variant === 'status') return statusHex.value ? null : 'bg-surface-400'
  if (props.variant === 'priority') return PRIORITY_DOT[props.value] || 'bg-surface-400'
  return null
})
</script>
