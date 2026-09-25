<template>
  <!-- Decorative trend strip: the number beside it is the readout, so this
       carries no accessible name. A flat series draws as a baseline. -->
  <svg
    :viewBox="`0 0 ${width} ${height}`"
    :width="width"
    :height="height"
    preserveAspectRatio="none"
    aria-hidden="true"
    class="block overflow-visible"
  >
    <path v-if="area" :d="area" :fill="color" fill-opacity="0.12" />
    <path
      :d="line"
      fill="none"
      :stroke="color"
      stroke-width="1.5"
      stroke-linejoin="round"
      stroke-linecap="round"
    />
    <circle v-if="last" :cx="last.x" :cy="last.y" r="2" :fill="color" />
  </svg>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  values: { type: Array, default: () => [] },
  color: { type: String, default: 'currentColor' },
  width: { type: Number, default: 96 },
  height: { type: Number, default: 28 },
})

const points = computed(() => {
  const v = props.values
  if (!v.length) return []
  const max = Math.max(1, ...v)
  const pad = 2
  const w = props.width - pad * 2
  const h = props.height - pad * 2
  const step = v.length > 1 ? w / (v.length - 1) : 0
  return v.map((n, i) => ({
    x: pad + (v.length > 1 ? i * step : w / 2),
    y: pad + h - (n / max) * h,
  }))
})

const line = computed(() =>
  points.value.map((p, i) => `${i ? 'L' : 'M'}${p.x.toFixed(1)} ${p.y.toFixed(1)}`).join(' '),
)

const area = computed(() => {
  const p = points.value
  if (p.length < 2) return ''
  const base = props.height - 2
  return `${line.value} L${p[p.length - 1].x.toFixed(1)} ${base} L${p[0].x.toFixed(1)} ${base} Z`
})

const last = computed(() => points.value[points.value.length - 1])
</script>
