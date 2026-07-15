import {
  bucketDaily,
  bucketWeekly,
  dayLabel,
  dayStarts,
  weekLabel,
  weekStarts,
} from '@/sections/project/config/trend'
import { useUserStore } from '@planetcrust/human-vue'
import { inject } from 'vue'

// Project audit-event activity, from the actionlog. Scoped to the project that
// owns the affected resource (resourceProjectID — works without a request-scope
// project, unlike scope.ProjectID). Reading needs the global action-log.read
// permission, so callers should treat a thrown metrics call as "unavailable"
// and hide the widget for non-admins.

const DAY = 86400000
const EVENTS_COLOR = '#6366f1'

export function useEventActivity() {
  const $SystemAPI = inject('$SystemAPI')
  const userStore = useUserStore()

  // Pulse (events/day) + grand-total stats over a window. Defaults to 30 days;
  // adapts day→week buckets past ~10 weeks so wide ranges stay legible.
  async function loadMetrics(projectID, { from, to } = {}) {
    const end = to || new Date()
    const start = from || new Date(end.getTime() - 30 * DAY)
    const base = {
      from: start.toISOString(),
      to: end.toISOString(),
      resourceProjectID: projectID || undefined,
    }

    // Two calls: a day series for the sparkline, and a grand total for the
    // stats (distinct actors can't be summed across day buckets).
    const [daily, totals] = await Promise.all([
      $SystemAPI.actionlogReport({ ...base, dimensions: ['day'], metrics: ['count'] }),
      $SystemAPI.actionlogReport({ ...base, metrics: ['count', 'actors', 'errors'] }),
    ])

    const points = (daily.set || [])
      .map(r => ({ date: String(r.dimensions?.day || ''), value: Number(r.metrics?.count || 0) }))
      .filter(p => p.date)

    const span = (end - start) / DAY
    const useWeek = span > 70
    const starts = useWeek ? weekStarts(start, end) : dayStarts(start, end)
    const g = (totals.set && totals.set[0]) || {}

    return {
      labels: starts.map(useWeek ? weekLabel : dayLabel),
      series: [
        {
          name: 'events',
          color: EVENTS_COLOR,
          data: (useWeek ? bucketWeekly : bucketDaily)(points, starts),
        },
      ],
      total: Number(g.metrics?.count || 0),
      actors: Number(g.metrics?.actors || 0),
      errors: Number(g.metrics?.errors || 0),
    }
  }

  // The most recent events, with their actors batch-resolved to names.
  async function loadRecent(projectID, limit = 6) {
    const { set } = await $SystemAPI.actionlogList({
      resourceProjectID: projectID || undefined,
      limit,
    })
    const items = set || []
    const ids = items.map(r => r.actorID).filter(id => id && id !== '0')
    await userStore.resolveUsers(ids).catch(() => {})
    return items
  }

  function actorName(data) {
    const user = data.actorID && userStore.findByID(data.actorID)
    if (user) return user.name || user.handle || user.username || user.email || user.userID
    return data.actor || data.actorID || ''
  }

  return { loadMetrics, loadRecent, actorName }
}
