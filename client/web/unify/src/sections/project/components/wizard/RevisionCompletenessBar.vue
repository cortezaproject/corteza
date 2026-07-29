<template>
  <!-- Revision completeness — completed work items ÷ items assigned to the
       OPEN revision, across all six work-item types (five categories +
       backlog items). Sits at the FOOT OF THE MANAGE & MONITOR RAIL and
       nowhere else (ruled 2026-07-28): it measures work items, so it only
       appears on the tab that is about work items — it used to live in the
       tab row and therefore also sat over Build and Govern, next to work it
       wasn't measuring. Entirely board-endpoint-driven (GET /project-board/,
       via composables/revisionCompleteness.js — see the script below), never
       the events/backlog stores for its NUMBERS (nor, any more, the report
       endpoint — see the script's own comment for why this moved off six
       report() calls onto one board call), so it stays accurate past those
       stores' 200-row-per-category cap. The stores do serve as the bar's
       refresh SIGNAL (their `mutations` counters), so completing a card on
       the board beside this bar moves it — data and invalidation are
       deliberately separate channels. -->
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
// One board call per load, totals only — the fetch + column arithmetic live
// in composables/revisionCompleteness.js (shared with Wizard.vue's publish
// confirm, which warns on the same numbers); see that module's comment for
// why the board endpoint and why `limit: 1` is safe for totals.
import { fetchRevisionCompleteness } from '@/sections/project/composables/revisionCompleteness'
import { colorFor } from '@/sections/project/config/chartColors'
import { useBacklogItemsStore } from '@/sections/project/stores/backlogItems'
import { useEventsStore } from '@/sections/project/stores/events'
import { computed, inject, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const $SystemAPI = inject('$SystemAPI')

// Refresh signal only, never data (see the template comment): any successful
// work-item write anywhere in the app bumps these stores' `mutations`
// counters, and the watcher below re-fetches the board totals — so the bar
// tracks the board/section mutations happening right beside it instead of
// staying frozen until a revision switch or remount.
const eventsStore = useEventsStore()
const backlogStore = useBacklogItemsStore()

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
    const { assigned, completed } = await fetchRevisionCompleteness($SystemAPI, pid, revId)
    if (mySeq !== seq) return // stale — a newer revision switch is in flight
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

// The mutation counters ride in the same watch as the ids: a bump re-runs
// load with the current scope, and the seq guard already collapses rapid
// bumps (e.g. a burst of board drags) into "last response wins".
watch(
  [
    () => props.projectId,
    () => props.revisionId,
    () => eventsStore.mutations,
    () => backlogStore.mutations,
  ],
  () => load(props.projectId, props.revisionId),
  { immediate: true },
)

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
