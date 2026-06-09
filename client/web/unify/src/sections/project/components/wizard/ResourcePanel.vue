<template>
  <div class="w-full h-full flex flex-col gap-3">
    <!-- Resource metrics strip -->
    <div class="shrink-0">
      <h3 class="text-xs font-semibold uppercase tracking-wide text-muted-color mb-2">Resources</h3>

      <div
        class="grid grid-cols-[repeat(auto-fill,minmax(8rem,1fr))] gap-2"
        :class="locked ? 'opacity-50 pointer-events-none select-none' : ''"
      >
        <div
          v-for="m in metrics"
          :key="m.key"
          class="rounded-lg border border-surface px-2.5 py-1.5 flex items-center gap-2 min-w-0"
        >
          <span
            :class="[
              'inline-flex items-center justify-center w-6 h-6 rounded-md ring-1 shrink-0',
              m.cfg.bg,
              m.cfg.ring,
            ]"
          >
            <i :class="[m.cfg.icon, m.cfg.text, 'text-xs']" />
          </span>
          <span class="text-base font-semibold leading-none shrink-0">{{ m.count }}</span>
          <span class="text-xs text-muted-color leading-tight min-w-0 break-words">{{ m.cfg.label }}</span>
        </div>
      </div>
    </div>

    <!-- Relationship graph fills the rest -->
    <div
      class="flex-1 min-h-0 rounded-lg border border-surface overflow-hidden bg-emphasis"
      :class="locked ? 'opacity-60 pointer-events-none' : ''"
    >
      <GroupGraph :project="project" :selected-group-id="selectedGroupId" :emphasize-kind="emphasizeKind" />
    </div>
  </div>
</template>

<script setup>
import GroupGraph from '@/sections/project/components/graph/GroupGraph.vue'
import { RESOURCE_KINDS, kindConfig } from '@/sections/project/config/kinds'
import { computed } from 'vue'

const props = defineProps({
  project: { type: Object, default: null },
  // Dimmed only while the current section is awaiting approval.
  locked: { type: Boolean, default: false },
  // Whole-project overview ('__all__') or focus a single group by id.
  selectedGroupId: { type: String, default: '__all__' },
  // Resource kind to emphasize in the graph (current step's kind).
  emphasizeKind: { type: String, default: null },
})

const counts = computed(() => {
  const out = {}
  for (const r of props.project?.resources || []) out[r.kind] = (out[r.kind] || 0) + 1
  return out
})

// Metric cards: one per resource kind.
const metrics = computed(() =>
  RESOURCE_KINDS.map(kind => ({ key: kind, cfg: kindConfig(kind), count: counts.value[kind] || 0 })),
)
</script>
