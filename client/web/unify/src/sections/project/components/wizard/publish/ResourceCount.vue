<template>
  <span
    class="rounded-lg border border-surface px-2.5 py-1.5 flex items-center gap-2 whitespace-nowrap"
  >
    <span
      class="inline-flex items-center justify-center w-6 h-6 rounded-md ring-1 shrink-0"
      :class="[cfg.bg, cfg.ring]"
    >
      <i :class="[cfg.icon, cfg.text, 'text-xs']" />
    </span>
    <span class="text-base font-medium leading-none">{{ count }}</span>
    <span class="text-xs text-muted-color leading-none">{{ $t(cfg.labelKey) }}</span>
    <span v-if="added" class="text-xs font-medium leading-none text-green-500">+{{ added }}</span>
  </span>
</template>

<script setup>
// One resource-kind count, visually identical to the metric cards above the
// Build tab's resource graph (components/graph/ResourceGraph.vue) — same card,
// same ringed kind badge, same count-then-label order. Deliberately NOT that
// component: those cards are layer toggles carrying visibility state, these
// just report. Shared by both publish surfaces so the two can't drift.
import { kindConfig } from '@/sections/project/config/kinds'
import { computed } from 'vue'

const props = defineProps({
  kind: { type: String, required: true },
  count: { type: Number, default: 0 },
  // How many of them are new in this revision; omitted on a first publish,
  // where everything is new and the delta would say nothing.
  added: { type: Number, default: 0 },
})

const cfg = computed(() => kindConfig(props.kind))
</script>
