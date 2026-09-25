import { onBeforeUnmount, onMounted, ref } from 'vue'

// Chart colours. Series colours are fixed hexes that clear the dataviz
// validator's bands on both the light (#ffffff) and dark (#18181b) card
// surfaces; chrome colours come from the live PrimeVue theme so a re-themed
// instance keeps its charts readable.
//
//   node <dataviz-skill>/scripts/validate_palette.js "#2a78d6,#E54122" --mode light --surface "#ffffff"
//   node <dataviz-skill>/scripts/validate_palette.js "#3987e5,#E54122" --mode dark --surface "#18181b"
//   node <dataviz-skill>/scripts/validate_palette.js "#059669,#d97706,#c32222" --mode light --surface "#ffffff"
//   node <dataviz-skill>/scripts/validate_palette.js "#059669,#d97706,#c32222" --mode dark --surface "#18181b"
export const SERIES = {
  activity: { light: '#2a78d6', dark: '#3987e5' },
  signins: { light: '#2a78d6', dark: '#3987e5' },
  created: { light: '#2a78d6', dark: '#3987e5' },
  error: '#E54122',
  completed: '#059669',
  running: '#d97706',
  failed: '#c32222',
  // neutral by design: a cancelled run carries no signal
  canceled: '#94a3b8',
}

// Inventory status splits share the outcome hues so "good / paused / gone"
// reads the same on a tile as on a runs chart.
export const STATUS_COLORS = {
  active: SERIES.completed,
  enabled: SERIES.completed,
  published: SERIES.completed,
  draft: '#2a78d6',
  suspended: SERIES.running,
  disabled: SERIES.running,
  archived: SERIES.running,
  deprecated: SERIES.running,
  deleted: SERIES.canceled,
}

export function seriesColor(key, dark) {
  const c = SERIES[key]
  if (!c) return STATUS_COLORS[key] || SERIES.canceled
  return typeof c === 'string' ? c : dark ? c.dark : c.light
}

export function statusColor(key) {
  return STATUS_COLORS[key] || SERIES.canceled
}

function cssVar(name, fallback) {
  const v = getComputedStyle(document.documentElement).getPropertyValue(name).trim()
  return v || fallback
}

// Text and hairline colours from the active PrimeVue theme.
export function chromeColors(dark) {
  return {
    text: cssVar('--p-text-color', dark ? '#f4f4f5' : '#18181b'),
    muted: cssVar('--p-text-muted-color', '#71717a'),
    grid: cssVar('--p-content-border-color', dark ? '#3f3f46' : '#e4e4e7'),
    surface: cssVar('--p-content-background', dark ? '#18181b' : '#ffffff'),
  }
}

// Tracks the shell's dark-mode class so charts re-theme without a reload.
export function useIsDark() {
  const isDark = ref(document.documentElement.classList.contains('dark'))
  let observer

  onMounted(() => {
    observer = new MutationObserver(() => {
      isDark.value = document.documentElement.classList.contains('dark')
    })
    observer.observe(document.documentElement, { attributes: true, attributeFilter: ['class'] })
  })

  onBeforeUnmount(() => observer?.disconnect())

  return isDark
}
