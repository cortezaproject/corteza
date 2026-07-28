// Bucketing + quick time-range presets for the created-over-time trend charts
// (Overview, CategoryView) and the audit-activity pulses (ActivityPanel,
// EventsActivityPanel via useEventActivity). Report data arrives as daily
// { date: 'YYYY-MM-DD', value } (optionally { group } for stacked series);
// governance volumes are low, so day/week/month bars read better than a spiky
// daily line — and empty buckets render as honest zero bars rather than a
// line interpolated across gaps.

const DAY = 86400000

// Monday 00:00 (local) on or before d.
function weekStart(d) {
  const x = new Date(d)
  x.setHours(0, 0, 0, 0)
  const dow = (x.getDay() + 6) % 7 // 0 = Monday
  x.setDate(x.getDate() - dow)
  return x
}

// Ordered week-start Dates spanning [from, to] inclusive. Walks the calendar
// (setDate +7) rather than adding WEEK milliseconds: a DST transition inside
// the span shifts ms-stepped "Mondays" by an hour, which silently dropped the
// final (current) week bucket — and with it anything created this week.
export function weekStarts(from, to) {
  const end = weekStart(to).getTime()
  const out = []
  const cur = weekStart(from)
  while (cur.getTime() <= end) {
    out.push(new Date(cur))
    cur.setDate(cur.getDate() + 7)
  }
  return out
}

// Sum daily points into week buckets aligned to `starts`. Returns a number[]
// the same length/order as `starts`. points: [{ date, value }]. Buckets are
// matched by exact week-start lookup (not ms division — same DST hazard as
// above).
export function bucketWeekly(points, starts) {
  const values = new Array(starts.length).fill(0)
  if (!starts.length) return values
  const idxByStart = new Map(starts.map((s, i) => [s.getTime(), i]))
  for (const p of points) {
    const d = new Date(p.date)
    if (Number.isNaN(d.getTime())) continue
    const idx = idxByStart.get(weekStart(d).getTime())
    if (idx !== undefined) values[idx] += Number(p.value || 0)
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

// Calendar-walked for the same DST reason as weekStarts above.
export function dayStarts(from, to) {
  const end = dayStart(to).getTime()
  const out = []
  const cur = dayStart(from)
  while (cur.getTime() <= end) {
    out.push(new Date(cur))
    cur.setDate(cur.getDate() + 1)
  }
  return out
}

export function bucketDaily(points, starts) {
  const values = new Array(starts.length).fill(0)
  if (!starts.length) return values
  const idxByStart = new Map(starts.map((s, i) => [s.getTime(), i]))
  for (const p of points) {
    const d = new Date(p.date)
    if (Number.isNaN(d.getTime())) continue
    const idx = idxByStart.get(dayStart(d).getTime())
    if (idx !== undefined) values[idx] += Number(p.value || 0)
  }
  return values
}

// Month-granularity variant, for the widest ranges (1Y/5Y/All). Calendar-
// accurate — months vary in length, so bucketing walks year/month pairs
// rather than dividing by a fixed duration like the day/week helpers above.
function monthStart(d) {
  const x = new Date(d)
  x.setHours(0, 0, 0, 0)
  x.setDate(1)
  return x
}

// Ordered first-of-month Dates spanning [from, to] inclusive.
export function monthStarts(from, to) {
  const start = monthStart(from)
  const end = monthStart(to)
  const out = []
  const cur = new Date(start)
  while (cur.getTime() <= end.getTime()) {
    out.push(new Date(cur))
    cur.setMonth(cur.getMonth() + 1)
  }
  return out
}

// Sum daily points into month buckets aligned to `starts`. Index is the
// year/month distance from the first bucket (division by a fixed duration
// does not work for months).
export function bucketMonthly(points, starts) {
  const values = new Array(starts.length).fill(0)
  if (!starts.length) return values
  const base = starts[0]
  const baseIdx = base.getFullYear() * 12 + base.getMonth()
  for (const p of points) {
    const d = new Date(p.date)
    if (Number.isNaN(d.getTime())) continue
    const idx = d.getFullYear() * 12 + d.getMonth() - baseIdx
    if (idx >= 0 && idx < values.length) values[idx] += Number(p.value || 0)
  }
  return values
}

// "Jul 2026" — month buckets can span multiple years (1Y/5Y/All), so the
// label needs the year to disambiguate (unlike day/week's "Jul 1").
export function monthLabel(d) {
  return d.toLocaleDateString(undefined, { month: 'short', year: 'numeric' })
}

// --- Tooltip-facing range labels ---------------------------------------------
// Full, unambiguous strings for one bucket — axis labels above stay short, but
// a tooltip has room to spell out the exact range being summed. One helper
// per granularity, mirrored to *Label above.

// "Jul 13, 2026" — day-bucket tooltip range.
function dayRangeLabel(d) {
  return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' })
}

// "Jul 13 – 19, 2026" / "Jun 29 – Jul 5, 2026" / "Dec 29, 2025 – Jan 4, 2026"
// — week-bucket tooltip range. `end` is `start` + 6 days, clamped to `toDate`
// so the current, still-partial week doesn't claim days that haven't
// happened yet.
function weekRangeLabel(start, toDate) {
  const end = new Date(start)
  end.setDate(end.getDate() + 6)
  const clampEnd = dayStart(toDate)
  if (end.getTime() > clampEnd.getTime()) end.setTime(clampEnd.getTime())

  const startMonthDay = start.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })
  if (start.getFullYear() === end.getFullYear()) {
    if (start.getMonth() === end.getMonth()) return `${startMonthDay} – ${end.getDate()}, ${end.getFullYear()}`
    const endMonthDay = end.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })
    return `${startMonthDay} – ${endMonthDay}, ${end.getFullYear()}`
  }
  const startFull = start.toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' })
  const endFull = end.toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' })
  return `${startFull} – ${endFull}`
}

// "July 2026" — month-bucket tooltip range.
function monthRangeLabel(d) {
  return d.toLocaleDateString(undefined, { month: 'long', year: 'numeric' })
}

// --- Quick time-range presets ------------------------------------------------
// Shared by ActivityPanel's filter window and the Overview/CategoryView trend
// charts (via TimeRangeSelect). Each preset resolves to a `from` Date (`to`
// stays open = "until now"); 'all' resolves to null (no lower bound — see
// adaptiveWindow below for how callers handle that). `months`/`years` walk the
// calendar (setMonth/setFullYear) rather than approximating with day counts;
// 'ytd' is Jan 1 of the current year.
export const RANGES = [
  { key: 'd1', days: 1 },
  { key: 'd5', days: 5 },
  { key: 'm1', months: 1 },
  { key: 'm6', months: 6 },
  { key: 'ytd', ytd: true },
  { key: 'y1', years: 1 },
  { key: 'y5', years: 5 },
  { key: 'all', all: true },
]

export function rangeFrom(r) {
  if (!r || r.all) return null
  const d = new Date()
  if (r.ytd) return new Date(d.getFullYear(), 0, 1)
  if (r.days) d.setDate(d.getDate() - r.days)
  if (r.months) d.setMonth(d.getMonth() - r.months)
  if (r.years) d.setFullYear(d.getFullYear() - r.years)
  return d
}

const DAY_BUCKET_MAX_SPAN_DAYS = 70
const WEEK_BUCKET_MAX_SPAN_DAYS = 400
const EMPTY_LOOKBACK_MONTHS = 6

// Picks day/week/month buckets for a [from, to] span so charts stay legible
// at any preset width: day buckets ≤ ~70 days, week buckets ≤ ~400 days,
// month buckets beyond that. `to` defaults to now; `from` defaults to a
// 6-month lookback, but that only matters when the caller has no data to
// size the window from — the 'all' preset carries no lower bound (rangeFrom
// returns null), so its callers fetch unbounded first and pass in the
// earliest fetched point's date instead (see earliestPointDate below, and
// Overview/CategoryView's loadTrend). 6 empty months reads as a sane empty
// chart for a project with nothing in it, where a multi-year fallback did not.
//
// Returns { fromISO, toISO, starts, labels, rangeLabels, bucket }: `labels`
// is the prebuilt axis-label array, `rangeLabels` is the tooltip-facing
// full-range string per bucket (aligned to `labels` — see the range-label
// helpers above), `bucket` is `points => number[]` aligned to it — callers
// just call it against a fetched point series instead of hand-rolling the
// day/week/month choice themselves.
export function adaptiveWindow(from, to) {
  const toDate = to || new Date()
  let fromDate = from
  if (!fromDate) {
    fromDate = new Date(toDate)
    fromDate.setMonth(fromDate.getMonth() - EMPTY_LOOKBACK_MONTHS)
  }

  const spanDays = (toDate.getTime() - fromDate.getTime()) / DAY

  let starts, labels, rangeLabels, bucket
  if (spanDays <= DAY_BUCKET_MAX_SPAN_DAYS) {
    starts = dayStarts(fromDate, toDate)
    labels = starts.map(dayLabel)
    rangeLabels = starts.map(dayRangeLabel)
    bucket = points => bucketDaily(points, starts)
  } else if (spanDays <= WEEK_BUCKET_MAX_SPAN_DAYS) {
    starts = weekStarts(fromDate, toDate)
    labels = starts.map(weekLabel)
    rangeLabels = starts.map(s => weekRangeLabel(s, toDate))
    bucket = points => bucketWeekly(points, starts)
  } else {
    starts = monthStarts(fromDate, toDate)
    labels = starts.map(monthLabel)
    rangeLabels = starts.map(monthRangeLabel)
    bucket = points => bucketMonthly(points, starts)
  }

  return { fromISO: fromDate.toISOString(), toISO: toDate.toISOString(), starts, labels, rangeLabels, bucket }
}

// Earliest point date across one or more fetched series, as a Date — how the
// 'all' preset sizes its window (rangeFrom returns null for it, so the chart
// starts at the first real data point instead of an arbitrary lookback).
// Points arrive date-ascending from the report store, so the first element
// of each series is its earliest. Null when every series is empty.
export function earliestPointDate(...seriesPoints) {
  let min = null
  for (const points of seriesPoints) {
    const d = points?.[0]?.date
    if (d && (!min || d < min)) min = d
  }
  return min ? new Date(min) : null
}
