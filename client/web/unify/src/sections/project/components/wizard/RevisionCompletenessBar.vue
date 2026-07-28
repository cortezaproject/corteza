<template>
  <!-- Revision completeness — completed work items ÷ items assigned to the
       OPEN revision, across all six work-item types (five categories +
       backlog items). Mounted once in Wizard.vue's tab row, right-aligned as
       a sibling of the locked publish/approval cluster (see Wizard.vue) —
       always visible on all three tabs, not gated on the publish flow.
       Entirely report-endpoint-driven (stores/report.js), never the events/
       backlog stores, so it stays accurate past their 200-row-per-category
       cap. -->
  <span v-if="hasData" class="flex items-center gap-2 shrink-0" v-tooltip.bottom="tooltip">
    <template v-if="hasItems">
      <ProgressBar :value="percent" :show-value="false" :style="progressStyle" class="w-24" />
      <span class="text-xs text-muted-color tabular-nums whitespace-nowrap">
        {{ $t('project.wizard.completeness.label', { percent }) }}
      </span>
    </template>
    <!-- Honest empty case: nothing assigned to this revision yet has no
         meaningful percentage — say so rather than rendering "0%" or
         dividing by zero. -->
    <span v-else class="text-xs text-muted-color whitespace-nowrap">
      {{ $t('project.wizard.completeness.empty') }}
    </span>
  </span>
</template>

<script setup>
// Six report calls (five categories + backlog items — the same
// KPI_RESOURCES set OverviewPanel.vue's revision-scoped KPI trio uses), each
// asking for the 'count'/'open' metrics. 'open' already applies the single
// Completed-status rule server-side (server/system/service/project_report.go
// projectReportCompletedStatus, the same rule stores/events.js#isOpenStatus
// names client-side) — completed = count - open is plain arithmetic on top
// of that, never a re-derivation of the open/closed rule itself.
import { CATEGORY_ORDER } from '@/sections/project/config/categories'
import { colorFor } from '@/sections/project/config/chartColors'
import { useReportStore } from '@/sections/project/stores/report'
import { computed, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const report = useReportStore()

const props = defineProps({
  // Chain-root project id — the report endpoint's required ProjectID scope.
  projectId: { type: [String, Number], required: true },
  // The open revision (the wizard's current project row's own projectID).
  // This bar only ever renders revision-scoped — there is no chain-wide mode
  // (unlike OverviewPanel, which this mirrors the report-call shape of).
  revisionId: { type: [String, Number], required: true },
})

const RESOURCES = [...CATEGORY_ORDER, 'backlog-item']

const totals = reactive({ assigned: 0, completed: 0 })
const loaded = ref(false)
const failed = ref(false)

// Sequence-guarded like OverviewPanel's loadAll: a revision switch mid-flight
// must not let a stale response land after a newer one already resolved.
let seq = 0
async function load(pid, revId) {
  if (!pid || !revId) {
    loaded.value = false
    return
  }
  const mySeq = ++seq
  try {
    const rowsets = await Promise.all(
      RESOURCES.map(key =>
        report.report(pid, key, { metrics: ['count', 'open'], revisionId: revId }),
      ),
    )
    if (mySeq !== seq) return // stale — a newer revision switch is in flight
    let assigned = 0
    let open = 0
    for (const rows of rowsets) {
      const g = rows?.[0]?.metrics || {}
      assigned += Number(g.count || 0)
      open += Number(g.open || 0)
    }
    totals.assigned = assigned
    totals.completed = assigned - open
    failed.value = false
    loaded.value = true
  } catch (err) {
    if (mySeq !== seq) return
    // A secondary header indicator, not a primary data surface: fail quiet
    // (hasData stays false, so the bar just doesn't render) rather than a
    // retry affordance competing with the tab row for space.
    console.error('Failed to load revision completeness', err)
    failed.value = true
    loaded.value = false
  }
}

watch([() => props.projectId, () => props.revisionId], ([pid, revId]) => load(pid, revId), {
  immediate: true,
})

const hasData = computed(() => loaded.value && !failed.value)
const hasItems = computed(() => totals.assigned > 0)
const percent = computed(() =>
  hasItems.value ? Math.round((totals.completed / totals.assigned) * 100) : 0,
)
const tooltip = computed(() =>
  hasItems.value
    ? t('project.wizard.completeness.tooltip', {
        completed: totals.completed,
        total: totals.assigned,
      })
    : t('project.wizard.completeness.emptyTooltip'),
)

// Slim header-row bar: shorter than PrimeVue's default 1.25rem, fill reuses
// the shared status palette's Completed hue (config/chartColors) rather than
// a bespoke colour, so it stays one system with the status donuts/badges.
const progressStyle = {
  '--p-progressbar-height': '0.375rem',
  '--p-progressbar-value-background': colorFor('status', 'Completed'),
}
</script>
