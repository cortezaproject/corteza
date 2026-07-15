import { defineStore } from 'pinia'
import { inject } from 'vue'

// Project category report store — a thin wrapper over the server-side
// aggregation endpoint (GET /project-report). Unlike the events store it never
// pulls full rowsets: it asks the backend for grouped counts, so overview
// metrics stay correct regardless of a category's volume (the events store caps
// its per-category load at 200 rows; aggregates here do not).
export const useReportStore = defineStore('projectReport', () => {
  const $SystemAPI = inject('$SystemAPI')

  // Run one report and return its raw row set: [{ dimensions, metrics }].
  // An empty `dimensions` yields a single grand-total row.
  async function report(projectId, resource, { dimensions = [], metrics = [], from, to } = {}) {
    const pid = String(projectId || '')
    if (!pid) return []
    const { set = [] } = await $SystemAPI.projectReportReport({
      resource,
      projectID: pid,
      dimensions,
      metrics,
      from,
      to,
    })
    return set || []
  }

  // Created-over-time points for a category: [{ date:'YYYY-MM-DD', group, value }]
  // sorted ascending. `from`/`to` are ISO strings windowing on created-at. When
  // `groupBy` (e.g. 'severity') is set the series splits by that dimension;
  // otherwise `group` is ''.
  async function trend(projectId, resource, { from, to, groupBy } = {}) {
    const dimensions = groupBy ? ['day', groupBy] : ['day']
    const rows = await report(projectId, resource, { dimensions, from, to })
    return rows
      .map(r => ({
        date: String(r.dimensions?.day || ''),
        group: groupBy ? String(r.dimensions?.[groupBy] ?? '') : '',
        value: Number(r.metrics?.count || 0),
      }))
      .filter(p => p.date)
      .sort((a, b) => a.date.localeCompare(b.date))
  }

  return { report, trend }
})
