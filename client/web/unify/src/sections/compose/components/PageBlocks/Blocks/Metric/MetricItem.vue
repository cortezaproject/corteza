<template>
  <div
    class="metric-item flex flex-col h-full p-3"
    :class="{ 'cursor-pointer hover:bg-emphasis': metric.drillDown?.enabled }"
    :style="containerStyle"
    @click="onClick"
  >
    <!-- Top row: label + change badge -->
    <div class="flex items-center gap-2 mb-1">
      <span
        v-if="metric.label"
        class="text-muted-color text-sm font-medium whitespace-nowrap overflow-hidden text-ellipsis"
      >
        {{ metric.label }}
      </span>

      <div
        v-if="changeValue != null"
        class="change-badge flex items-center gap-1 text-xs font-semibold px-2 py-0.5 rounded-full whitespace-nowrap ml-auto"
        :class="changeClass"
      >
        <i :class="changeIcon" class="text-xs" />
        {{ formattedChange }}
      </div>
    </div>

    <!-- Value row -->
    <span class="metric-value font-bold text-color" :style="valueStyle">
      <template v-if="metric.prefix">{{ metric.prefix }}</template>
      {{ value.value }}
      <template v-if="metric.suffix">{{ metric.suffix }}</template>
    </span>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  metric: {
    type: Object,
    default: () => ({}),
  },
  value: {
    type: Object,
    default: () => ({}),
  },
  changeValue: {
    type: Number,
    default: null,
  },
})

const emit = defineEmits(['drill-down'])

function onClick() {
  if (props.metric?.drillDown?.enabled) {
    emit('drill-down')
  }
}

const containerStyle = computed(() => {
  const s = props.metric.valueStyle || {}
  const d = {}
  if (s.backgroundColor && s.backgroundColor !== 'transparent')
    d.backgroundColor = s.backgroundColor
  return d
})

const valueStyle = computed(() => {
  const s = props.metric.valueStyle || {}
  const d = {}
  if (s.color && s.color !== 'transparent') d.color = s.color
  return d
})

// Change indicator computed properties
const changeClass = computed(() => {
  if (props.changeValue == null) return ''
  if (props.changeValue > 0) return 'change-positive'
  if (props.changeValue < 0) return 'change-negative'
  return 'change-neutral'
})

const changeIcon = computed(() => {
  if (props.changeValue == null) return ''
  if (props.changeValue > 0) return 'pi pi-arrow-up-right'
  if (props.changeValue < 0) return 'pi pi-arrow-down-right'
  return 'pi pi-arrow-right'
})

const formattedChange = computed(() => {
  if (props.changeValue == null) return ''
  if (props.changeValue === 0) return '—'
  const val = Math.abs(props.changeValue).toFixed(1)
  if (props.changeValue > 0) return `+${val}%`
  return `-${val}%`
})
</script>

<style scoped>
.metric-item {
  min-height: 0;
  container-type: size;
}

.metric-value {
  display: flex;
  align-items: center;
  justify-content: center;
  flex: 1;
  min-height: 0;
  font-size: min(90cqw, 60cqh);
  line-height: 1;
  letter-spacing: -0.02em;
}

.change-badge.change-positive {
  background-color: rgba(34, 197, 94, 0.15);
  color: rgb(34, 197, 94);
}

.change-badge.change-negative {
  background-color: rgba(239, 68, 68, 0.15);
  color: rgb(239, 68, 68);
}

.change-badge.change-neutral {
  background-color: color-mix(in srgb, var(--p-text-muted-color) 15%, transparent);
  color: var(--p-text-muted-color);
}
</style>
