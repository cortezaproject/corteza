<template>
  <!-- Revision completeness — completed work items ÷ items assigned to the
       OPEN revision, across all six work-item types (five categories +
       backlog items). Sits at the FOOT OF THE MANAGE & MONITOR RAIL and
       nowhere else (ruled 2026-07-28): it measures work items, so it only
       appears on the tab that is about work items — it used to live in the
       tab row and therefore also sat over Build and Govern, next to work it
       wasn't measuring. Entirely board-endpoint-driven (GET /project-board/,
       see the script below), never the events/backlog stores (nor, any more,
       the report endpoint — see the script's own comment for why this moved
       off six report() calls onto one board call), so it stays accurate past
       those stores' 200-row-per-category cap. -->
  <div v-if="hasData" class="flex flex-col gap-1.5" v-tooltip.top="tooltip">
    <div class="flex items-center justify-between gap-2">
      <span class="text-xs font-medium text-muted-color uppercase tracking-wide">
        {{ $t('project.wizard.completeness.title') }}
      </span>
      <span class="text-xs text-muted-color tabular-nums whitespace-nowrap">
        {{
          hasItems
            ? $t('project.wizard.completeness.label', { percent })
            : $t('project.wizard.completeness.empty')
        }}
      </span>
    </div>
    <!-- Nothing assigned still renders the track, drawn at zero, so the rail
         keeps its shape whether or not work exists. The tooltip carries the
         explanation the old text-only empty state used to spell out inline. -->
    <ProgressBar :value="percent" :show-value="false" :style="progressStyle" />
  </div>
</template>

<script setup>
// One board call (GET /project-board/, see server/system/types/project_board.go),
// asking for every column's true total and no cards (`limit: 1` — the board
// endpoint decouples its total-probe from the caller's requested item-page
// limit specifically for callers like this one, see service/project_board.go's
// loadColumn doc comment). `assigned` = the sum of all four column totals;
// `completed` = the Completed column's own total. The backend's board service
// verified this sum is mathematically identical to the six-report-call
// count/open arithmetic this component used before — the single Completed
// status bucket IS the "not open" set, so summing every column equals the old
// `count` total and the Completed column alone equals the old `count - open`.
// Nothing here re-derives that rule; it's just addition over the endpoint's
// own totals.
import { colorFor } from '@/sections/project/config/chartColors'
import { computed, inject, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const $SystemAPI = inject('$SystemAPI')

const props = defineProps({
  // Chain-root project id — the board endpoint's required ProjectID scope.
  projectId: { type: [String, Number], required: true },
  // The open revision (the wizard's current project row's own projectID).
  // This bar only ever renders revision-scoped — there is no chain-wide mode
  // (unlike OverviewPanel, which this mirrors the report-call shape of).
  revisionId: { type: [String, Number], required: true },
})

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
    const { columns = [] } = await $SystemAPI.projectBoardBoard({
      projectID: pid,
      revisionID: revId,
      limit: 1,
    })
    if (mySeq !== seq) return // stale — a newer revision switch is in flight
    let assigned = 0
    let completed = 0
    for (const col of columns) {
      const total = Number(col.total || 0)
      assigned += total
      if (col.status === 'Completed') completed += total
    }
    totals.assigned = assigned
    totals.completed = completed
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
