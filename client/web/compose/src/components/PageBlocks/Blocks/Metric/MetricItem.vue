<template>
  <div
    :style="genStyle(metric.valueStyle)"
    class="h-full text-center"
  >
    <svg
      :viewBox="viewBox"
      class="h-full w-full flex"
      width="100%"
      height="100%"
    >
      <text
        ref="metricItemRef"
        y="50%"
        x="50%"
        text-anchor="middle"
        dominant-baseline="central"
        text-rendering="geometricPrecision"
      >
        <template v-if="metric.prefix">
          {{ metric.prefix }}
        </template>
        {{ value.value }}
        <template v-if="metric.suffix">
          {{ metric.suffix }}
        </template>
      </text>
    </svg>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick, onBeforeUnmount } from 'vue'

const props = defineProps({
  metric: {
    type: Object,
    default: () => ({}),
  },
  value: {
    type: Object,
    default: () => ({}),
  },
})

const metricItemRef = ref(null)
const vvb = ref(['0', '0', '0', '0'])

const viewBox = computed(() => vvb.value.join(' '))

function update () {
  nextTick(() => {
    if (!metricItemRef.value) return
    try {
      const { width, height } = metricItemRef.value.getBBox()
      const tmp = [...vvb.value]
      tmp[2] = String(Math.ceil(width))
      tmp[3] = String(Math.ceil(height))
      vvb.value = tmp
    } catch {
      // SVG not rendered yet
    }
  })
}

function genStyle (s = {}) {
  const d = {}
  if (s?.color) {
    d.color = s.color
    d.fill = s.color
  }
  if (s?.backgroundColor) d.backgroundColor = s.backgroundColor
  if (s?.fontSize) d.fontSize = s.fontSize + 'px'
  return d
}

watch(() => props.metric, update, { immediate: true, deep: true })
watch(() => props.value, update, { immediate: true, deep: true })

onBeforeUnmount(() => {
  vvb.value = []
})
</script>
