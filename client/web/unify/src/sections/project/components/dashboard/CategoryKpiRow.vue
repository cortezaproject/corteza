<template>
  <!-- Responsive row of stat cards; wraps to a single column on narrow screens. -->
  <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
    <div
      v-for="kpi in kpis"
      :key="kpi.labelKey"
      class="rounded-xl border border-surface bg-surface-50 dark:bg-surface-950 px-4 py-3 flex flex-col gap-1"
    >
      <div class="text-[11px] font-semibold uppercase tracking-wide text-muted-color">
        {{ $t(kpi.labelKey) }}
      </div>
      <div class="text-2xl font-semibold text-color leading-none">
        {{ kpi.value }}
      </div>
    </div>
  </div>
</template>

<script setup>
// Presentational KPI row. Parent maps its category config to { labelKey, value }.
defineProps({
  kpis: {
    type: Array,
    required: true,
  },
})
</script>
