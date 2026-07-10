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

const SEVERITY = {
  Critical: 'bg-red-100 text-red-700 dark:bg-red-500/15 dark:text-red-300',
  Serious: 'bg-fuchsia-100 text-fuchsia-700 dark:bg-fuchsia-500/15 dark:text-fuchsia-300',
  Major: 'bg-orange-100 text-orange-700 dark:bg-orange-500/15 dark:text-orange-300',
  Minor: 'bg-amber-100 text-amber-700 dark:bg-amber-500/15 dark:text-amber-300',
  Informational: 'bg-sky-100 text-sky-700 dark:bg-sky-500/15 dark:text-sky-300',
}
const SEVERITY_DOT = {
  Critical: 'bg-red-500',
  Serious: 'bg-fuchsia-500',
  Major: 'bg-orange-500',
  Minor: 'bg-amber-500',
  Informational: 'bg-sky-500',
}
const STATUS = {
  Open: MUTED,
  'In Progress': 'bg-primary/15 text-primary',
  'Ready to Test': 'bg-amber-100 text-amber-700 dark:bg-amber-500/15 dark:text-amber-300',
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
