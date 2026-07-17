<template>
  <!-- Flat tinted pill used for the type / severity / status columns, matching
       the demo's badges. Severity carries a leading colored dot. -->
  <span
    class="inline-flex items-center gap-1.5 rounded-full font-medium whitespace-nowrap"
    :class="[cls, size === 'md' ? 'px-2.5 py-1 text-xs' : 'px-2 py-0.5 text-[11px]']"
  >
    <span
      v-if="dot"
      class="rounded-full shrink-0"
      :class="[dot, size === 'md' ? 'w-2 h-2' : 'w-1.5 h-1.5']"
    />
    {{ value || '—' }}
  </span>
</template>

<script setup>
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
// Lifecycle ramp — mirrors STATUS_COLORS in config/chartColors so pills and
// charts read as one system. Open is blue (not muted) so it never looks blank.
const STATUS = {
  Open: 'bg-blue-100 text-blue-700 dark:bg-blue-500/15 dark:text-blue-300',
  'In Progress': 'bg-amber-100 text-amber-700 dark:bg-amber-500/15 dark:text-amber-300',
  'Ready to Test': 'bg-violet-100 text-violet-700 dark:bg-violet-500/15 dark:text-violet-300',
  Completed: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-500/15 dark:text-emerald-300',
}

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

// Risk (Critical/High/Medium/Low — used in selects, where RiskPips' bare pip
// meter reads poorly without its label alongside). Same red/orange/amber/
// yellow family as RiskPips' own MAP, so the select and the list pips agree
// on what each level means even though they're two different mark types.
const RISK = {
  Critical: 'bg-red-100 text-red-700 dark:bg-red-500/15 dark:text-red-300',
  High: 'bg-orange-100 text-orange-700 dark:bg-orange-500/15 dark:text-orange-300',
  Medium: 'bg-amber-100 text-amber-700 dark:bg-amber-500/15 dark:text-amber-300',
  Low: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-500/15 dark:text-yellow-300',
}
const RISK_DOT = {
  Critical: 'bg-red-500',
  High: 'bg-orange-500',
  Medium: 'bg-amber-500',
  Low: 'bg-yellow-500',
}

const cls = computed(() => {
  if (props.variant === 'severity') return SEVERITY[props.value] || MUTED
  if (props.variant === 'status') return STATUS[props.value] || MUTED
  if (props.variant === 'priority') return PRIORITY[props.value] || MUTED
  if (props.variant === 'risk') return RISK[props.value] || MUTED
  return 'border border-surface text-muted-color' // type: flat neutral pill
})
const dot = computed(() => {
  if (props.variant === 'severity') return SEVERITY_DOT[props.value] || 'bg-surface-400'
  if (props.variant === 'priority') return PRIORITY_DOT[props.value] || 'bg-surface-400'
  if (props.variant === 'risk') return RISK_DOT[props.value] || 'bg-surface-400'
  return null
})
</script>
