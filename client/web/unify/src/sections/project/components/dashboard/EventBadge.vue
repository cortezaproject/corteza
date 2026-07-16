<template>
  <!-- Flat tinted pill used for the type / severity / status columns, matching
       the demo's badges. Severity carries a leading colored dot. -->
  <span
    class="inline-flex items-center gap-1.5 rounded-full px-2 py-0.5 text-[11px] font-medium whitespace-nowrap"
    :class="cls"
  >
    <span v-if="dot" class="w-1.5 h-1.5 rounded-full shrink-0" :class="dot" />
    {{ value || '—' }}
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  value: { type: String, default: '' },
  // 'type' (neutral bordered) | 'severity' (dot + tint) | 'status' (tint)
  variant: { type: String, default: 'type' },
})

const MUTED = 'bg-surface-200 text-muted-color dark:bg-surface-700'

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

const cls = computed(() => {
  if (props.variant === 'severity') return SEVERITY[props.value] || MUTED
  if (props.variant === 'status') return STATUS[props.value] || MUTED
  return 'border border-surface text-muted-color' // type: flat neutral pill
})
const dot = computed(() =>
  props.variant === 'severity' ? SEVERITY_DOT[props.value] || 'bg-surface-400' : null,
)
</script>
