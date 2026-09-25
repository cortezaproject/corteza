<template>
  <div class="flex flex-col gap-1.5 min-w-0">
    <!-- Proportional split; a 2px surface gap separates segments -->
    <div class="flex h-1.5 w-full gap-0.5 overflow-hidden rounded-full">
      <div
        v-for="s in segments"
        :key="s.key"
        :style="{ width: s.pct + '%', backgroundColor: s.color }"
        :title="`${s.label}: ${s.count}`"
        class="h-full rounded-full"
      />
    </div>
    <div v-if="!compact" class="flex flex-wrap gap-x-3 gap-y-0.5 text-xs text-muted-color">
      <span
        v-for="s in segments"
        :key="s.key"
        class="inline-flex items-center gap-1 whitespace-nowrap"
      >
        <span class="inline-block h-2 w-2 rounded-full" :style="{ backgroundColor: s.color }" />
        {{ s.count.toLocaleString() }} {{ s.label }}
      </span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { statusColor } from './chartTheme'

const props = defineProps({
  // { statusKey: count }
  status: { type: Object, required: true },
  // display order; unknown keys follow, sorted
  order: { type: Array, default: () => [] },
  compact: { type: Boolean, default: false },
})

const { t, te } = useI18n()

const segments = computed(() => {
  const keys = props.order.length ? props.order : Object.keys(props.status).sort()
  const total = keys.reduce((sum, k) => sum + (props.status[k] || 0), 0)
  return keys
    .filter(k => (props.status[k] || 0) > 0)
    .map(k => ({
      key: k,
      count: props.status[k],
      pct: total ? (props.status[k] / total) * 100 : 0,
      color: statusColor(k),
      label: te(`dashboard.status.${k}`) ? t(`dashboard.status.${k}`) : k,
    }))
})
</script>
