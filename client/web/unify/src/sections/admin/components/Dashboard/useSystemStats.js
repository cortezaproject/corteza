import { computed, inject, ref, watch } from 'vue'

// Range presets. The bucket is chosen here rather than left to the server's
// default so the x-axis is predictable per preset.
export const RANGES = [
  { key: '7d', days: 7, bucket: 'day' },
  { key: '30d', days: 30, bucket: 'day' },
  { key: '90d', days: 90, bucket: 'week' },
  { key: '1y', days: 365, bucket: 'month' },
]

const RANGE_KEY = 'admin.dashboard.range'

function readPref(key, fallback, allowed) {
  try {
    const v = localStorage.getItem(key)
    return allowed.includes(v) ? v : fallback
  } catch {
    return fallback
  }
}

function writePref(key, value) {
  try {
    localStorage.setItem(key, value)
  } catch {
    // a blocked store only loses the remembered preset
  }
}

export function usePreference(key, fallback, allowed) {
  const pref = ref(readPref(key, fallback, allowed))
  watch(pref, v => writePref(key, v))
  return pref
}

// The request a preset stands for: local-midnight start, now as the end,
// and the preset's bucket.
export function rangeParams(rangeKey) {
  const preset = RANGES.find(r => r.key === rangeKey) || RANGES[1]
  const to = new Date()
  const from = new Date(to)
  from.setDate(from.getDate() - preset.days)
  from.setHours(0, 0, 0, 0)
  return { from: from.toISOString(), to: to.toISOString(), bucket: preset.bucket }
}

// Short axis labels and fuller tooltip labels for a payload's range.
export function bucketLabelsFor(range) {
  const buckets = range?.buckets || []
  const bucket = range?.bucket || 'day'
  return {
    short: buckets.map(b => shortLabel(b, bucket)),
    long: buckets.map((b, i) => longLabel(b, bucket, buckets[i + 1], range?.to)),
  }
}

// One request per range change; the previous one is cancelled so a slow
// answer never lands over a newer one. Loaded data stays on screen while the
// next range loads.
export function useSystemStats() {
  const $SystemAPI = inject('$SystemAPI')

  const range = usePreference(
    RANGE_KEY,
    '30d',
    RANGES.map(r => r.key),
  )
  const stats = ref(null)
  const loading = ref(false)
  const error = ref(null)

  let pending = null

  async function load() {
    pending?.cancel()
    const req = $SystemAPI.statsListCancellable(rangeParams(range.value))
    pending = req
    loading.value = true
    error.value = null

    try {
      const result = await req.response()
      if (pending === req) stats.value = result
    } catch (e) {
      if (pending === req && !e?.__CANCEL__) error.value = e
    } finally {
      if (pending === req) loading.value = false
    }
  }

  watch(range, load, { immediate: true })

  const bucket = computed(() => stats.value?.range?.bucket || 'day')
  const buckets = computed(() => stats.value?.range?.buckets || [])
  const bucketLabels = computed(() => bucketLabelsFor(stats.value?.range).short)
  const rangeLabels = computed(() => bucketLabelsFor(stats.value?.range).long)

  return { range, stats, loading, error, reload: load, bucket, buckets, bucketLabels, rangeLabels }
}

function parseDay(s) {
  const [y, m, d] = s.split('-').map(Number)
  return new Date(y, m - 1, d)
}

function shortLabel(day, bucket) {
  const d = parseDay(day)
  if (bucket === 'month')
    return d.toLocaleDateString(undefined, { month: 'short', year: '2-digit' })
  return d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })
}

function longLabel(day, bucket, nextDay, to) {
  const d = parseDay(day)
  if (bucket === 'day') {
    return d.toLocaleDateString(undefined, {
      weekday: 'short',
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    })
  }
  if (bucket === 'month') return d.toLocaleDateString(undefined, { month: 'long', year: 'numeric' })

  const end = nextDay ? parseDay(nextDay) : to ? new Date(to) : null
  if (end) end.setDate(end.getDate() - 1)
  const a = d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })
  const b = end
    ? end.toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' })
    : ''
  return b ? `${a} – ${b}` : a
}

export function sum(values) {
  return (values || []).reduce((a, b) => a + b, 0)
}
