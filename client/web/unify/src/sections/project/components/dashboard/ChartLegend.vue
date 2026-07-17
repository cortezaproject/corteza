<template>
  <!-- HTML legend strip below a chart canvas, built from the same badge
       components the tables/forms use (EventBadge/RiskPips/KindIcon) rather
       than echarts' own canvas legend — a slice/series reads as the exact
       pill a user already knows from a row. The trade: these pills don't
       toggle series visibility the way the canvas legend did; accepted,
       since every chart here has a small, fixed series count. -->
  <div class="flex flex-wrap items-center justify-center gap-x-3 gap-y-1.5 mt-2">
    <template v-for="item in items" :key="item.key ?? item.label">
      <EventBadge
        v-if="variant === 'status' || variant === 'severity' || variant === 'priority'"
        :value="item.label"
        :variant="variant"
      />
      <span v-else-if="variant === 'risk'" class="inline-flex items-center gap-1.5">
        <RiskPips :level="item.label" />
        <span class="text-xs text-muted-color">{{ item.label }}</span>
      </span>
      <span v-else-if="variant === 'category'" class="inline-flex items-center gap-1.5 min-w-0">
        <KindIcon :config="CATEGORY_CONFIG[item.key]?.badge" size="sm" />
        <span class="text-xs text-muted-color truncate">{{ item.label }}</span>
      </span>
      <!-- Anything else, including 'type': EventBadge's type pill is a
           neutral bordered one (no tint), which would make every legend
           entry look identical — colour is the only thing that carries
           slice identity here, so fall back to a plain coloured dot. -->
      <span v-else class="inline-flex items-center gap-1.5">
        <span class="w-2 h-2 rounded-full shrink-0" :style="{ background: item.color }" />
        <span class="text-xs text-muted-color">{{ item.label }}</span>
      </span>
    </template>
  </div>
</template>

<script setup>
import EventBadge from '@/sections/project/components/dashboard/EventBadge.vue'
import RiskPips from '@/sections/project/components/dashboard/RiskPips.vue'
import KindIcon from '@/sections/project/components/KindIcon.vue'
import { CATEGORY_CONFIG } from '@/sections/project/config/categories'

defineProps({
  // [{ label: string, key?: string, color?: string }]. `key` is the category
  // key for variant 'category'; `color` feeds the fallback dot.
  items: { type: Array, default: () => [] },
  // Badge family to render: status | severity | priority | risk | category |
  // '' (colored-dot fallback, also used for 'type' — see above).
  variant: { type: String, default: '' },
})
</script>
