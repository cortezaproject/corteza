// Weekly bucketing for the created-over-time trend charts. Report data arrives
// as daily { date: 'YYYY-MM-DD', value }; governance volumes are low, so weekly
// bars read better than a spiky daily line — and empty weeks render as honest
// zero bars rather than a line interpolated across gaps.

const DAY = 86400000
const WEEK = 7 * DAY

// Monday 00:00 (local) on or before d.
function weekStart(d) {
  const x = new Date(d)
  x.setHours(0, 0, 0, 0)
  const dow = (x.getDay() + 6) % 7 // 0 = Monday
  x.setDate(x.getDate() - dow)
  return x
}

// Ordered week-start Dates spanning [from, to] inclusive.
export function weekStarts(from, to) {
  const start = weekStart(from).getTime()
  const end = weekStart(to).getTime()
  const out = []
  for (let t = start; t <= end; t += WEEK) out.push(new Date(t))
  return out
}

// Sum daily points into week buckets aligned to `starts`. Returns a number[]
// the same length/order as `starts`. points: [{ date, value }].
export function bucketWeekly(points, starts) {
  const values = new Array(starts.length).fill(0)
  if (!starts.length) return values
  const base = starts[0].getTime()
  for (const p of points) {
    const d = new Date(p.date)
    if (Number.isNaN(d.getTime())) continue
    const idx = Math.round((weekStart(d).getTime() - base) / WEEK)
    if (idx >= 0 && idx < values.length) values[idx] += Number(p.value || 0)
  }
  return values
}

// Short axis label, e.g. "Jul 1". Shared by day + week buckets.
export function weekLabel(d) {
  return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })
}
export const dayLabel = weekLabel

// Day-granularity variant of the helpers above, for shorter ranges.
function dayStart(d) {
  const x = new Date(d)
  x.setHours(0, 0, 0, 0)
  return x
}

export function dayStarts(from, to) {
  const start = dayStart(from).getTime()
  const end = dayStart(to).getTime()
  const out = []
  for (let t = start; t <= end; t += DAY) out.push(new Date(t))
  return out
}

export function bucketDaily(points, starts) {
  const values = new Array(starts.length).fill(0)
  if (!starts.length) return values
  const base = starts[0].getTime()
  for (const p of points) {
    const d = new Date(p.date)
    if (Number.isNaN(d.getTime())) continue
    const idx = Math.round((dayStart(d).getTime() - base) / DAY)
    if (idx >= 0 && idx < values.length) values[idx] += Number(p.value || 0)
  }
  return values
}

// Default trend window: the last `weeks` weeks ending today. Returns ISO
// strings for the report `from`/`to` params plus the aligned week starts.
export function trendWindow(weeks = 12) {
  const to = new Date()
  const from = new Date(to.getTime() - weeks * WEEK)
  return { fromISO: from.toISOString(), toISO: to.toISOString(), starts: weekStarts(from, to) }
}
