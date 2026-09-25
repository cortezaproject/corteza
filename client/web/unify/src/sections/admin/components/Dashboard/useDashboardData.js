import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { filters } from '@planetcrust/human-vue'
import { SERIES, seriesColor } from './chartTheme'
import { RESOURCES, TAQ_OUTCOMES, WORKFLOW_OUTCOMES, foldOutcomes, liveCount } from './resources'
import { sum } from './useSystemStats'
import { describeEvent } from './eventLabel'

const { locFullDateTime } = filters

// Everything a layout draws, derived once from the stats payload so the
// three layouts differ in arrangement only. Every field is null when the
// caller may not see that section.
export function useDashboardData(stats, isDark, bucketLabels, rangeLabels) {
  const i18n = useI18n()
  const { t } = i18n

  const n = computed(() => stats.value?.range?.buckets?.length || 0)
  const color = key => seriesColor(key, isDark.value)

  const resources = computed(() => stats.value?.resources || {})

  const tiles = computed(() =>
    RESOURCES.filter(r => resources.value[r.key]).map(r => {
      const s = resources.value[r.key]
      return {
        key: r.key,
        label: t(`dashboard.resources.${r.key}`),
        route: r.route,
        total: s.total,
        live: liveCount(s.status),
        status: s.status,
        created: s.createdInRange,
        spark: s.created,
        resource: r,
      }
    }),
  )

  const activity = computed(() => {
    const a = stats.value?.activity
    if (!a) return null
    return {
      total: a.total,
      errors: a.errors,
      series: [
        {
          key: 'total',
          name: t('dashboard.activity.entries'),
          color: color('activity'),
          data: a.series.total,
        },
        {
          key: 'errors',
          name: t('dashboard.activity.errors'),
          color: SERIES.error,
          data: a.series.errors,
        },
      ],
      errorSpark: a.series.errors,
      recentErrors: a.recentErrors.map(e => ({
        id: e.actionID,
        at: e.timestamp,
        title: describeEvent(i18n, e.resource, e.action),
        detail: e.error || e.description,
        time: locFullDateTime(e.timestamp),
        to: { name: 'system.actionLog' },
        color: SERIES.error,
      })),
    }
  })

  const signins = computed(() => {
    const s = stats.value?.signins
    if (!s) return null
    return {
      total: s.total,
      users: s.users,
      live: s.live,
      series: [
        {
          key: 'signins',
          name: t('dashboard.signins.title'),
          color: color('signins'),
          data: s.series,
        },
      ],
      spark: s.series,
    }
  })

  function runs(raw, outcomes, labelOf, failures) {
    if (!raw) return null
    const folded = foldOutcomes(raw, outcomes, n.value)
    return {
      total: raw.total,
      failed: folded.byStatus.failed,
      completed: folded.byStatus.completed,
      byStatus: folded.byStatus,
      // stacked worst-first from the bottom is hard to read; completed sits
      // at the base, failures on top where they stand out
      series: [...outcomes].reverse().map(o => ({
        key: o.key,
        name: labelOf(o.key),
        color: color(o.key),
        data: folded.series[o.key],
      })),
      failSpark: folded.series.failed,
      failures,
    }
  }

  const workflows = computed(() => {
    const w = stats.value?.automation?.workflows
    return runs(
      w,
      WORKFLOW_OUTCOMES,
      k => t(`dashboard.runs.${k}`),
      (w?.failures || []).map(f => ({
        id: f.sessionID,
        at: f.createdAt,
        title:
          f.workflowName || f.workflowHandle || t('dashboard.runs.workflow', { id: f.workflowID }),
        detail: f.error,
        time: locFullDateTime(f.createdAt),
        to: { name: 'automation.sessions.view', params: { sessionID: f.sessionID } },
        color: SERIES.failed,
      })),
    )
  })

  const taqs = computed(() => {
    const q = stats.value?.automation?.taqs
    return runs(
      q,
      TAQ_OUTCOMES,
      k => t(`dashboard.runs.${k}`),
      (q?.failures || []).map(f => ({
        id: f.actionID,
        at: f.timestamp,
        title: t('dashboard.runs.taq', { id: f.meta?.automationID || '' }),
        detail: f.error,
        time: locFullDateTime(f.timestamp),
        to: f.meta?.automationID
          ? { name: 'automation.taq.edit', params: { automationID: f.meta.automationID } }
          : { name: 'automation.sessions' },
        color: SERIES.failed,
      })),
    )
  })

  // Failed runs and errors, newest first, for the attention lists.
  const attention = computed(() =>
    [
      ...(workflows.value?.failures || []),
      ...(taqs.value?.failures || []),
      ...(activity.value?.recentErrors || []),
    ]
      .sort((a, b) => (a.at < b.at ? 1 : -1))
      .slice(0, 10),
  )

  return {
    tiles,
    activity,
    signins,
    workflows,
    taqs,
    attention,
    bucketLabels,
    rangeLabels,
    n,
    color,
    sum,
  }
}
