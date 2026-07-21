<template>
  <nav class="w-full h-full overflow-y-auto rounded-xl border border-surface bg-surface">
    <ul class="p-3 flex flex-col gap-0.5">
      <li v-for="step in steps" :key="step.key">
        <button
          type="button"
          class="w-full text-left rounded-md px-3 py-2 flex items-center gap-2 text-sm transition-colors"
          :class="
            step.key === activeKey ? 'bg-primary/10 text-primary font-medium' : 'hover:bg-emphasis'
          "
          @click="$emit('select', step.key)"
        >
          <!-- Every step shows its own icon in the same badge used by the
               metrics strip and resource item lists (resource steps reuse the
               kind icon/colour). Approval only exists at publish now, so the
               single governance signal surfaced here is "changes requested"
               on the reviewed step. -->
          <span
            v-if="changesRequested(step)"
            class="inline-flex items-center justify-center w-6 h-6 rounded-md shrink-0 ring-1 bg-amber-500/10 ring-amber-500/30"
            :title="$t('project.publish.approval.status.changesRequested')"
          >
            <i class="pi pi-exclamation-circle text-amber-500 text-xs" />
          </span>
          <span
            v-else
            class="inline-flex items-center justify-center w-6 h-6 rounded-md shrink-0"
            :class="stepBadge(step)"
          >
            <i :class="[...stepIconClass(step), 'text-xs']" />
          </span>
          <span class="flex-1 min-w-0 truncate">{{ $t(step.labelKey) }}</span>
        </button>
      </li>
    </ul>
  </nav>
</template>

<script setup>
import { kindConfig } from '@/sections/project/config/kinds'

const props = defineProps({
  steps: { type: Array, required: true },
  activeKey: { type: String, default: '' },
  statuses: { type: Object, default: () => ({}) },
})
defineEmits(['select'])

function changesRequested(step) {
  return props.statuses[step.key] === 'changes-requested'
}

// Free-build leading icon — rendered in the same rounded badge the metrics strip
// and resource item lists use. Resource steps borrow the kind's icon, colour, bg
// and ring; other steps fall back to a neutral badge with their step icon.
function stepBadge(step) {
  if (step.kind) {
    const cfg = kindConfig(step.kind)
    return ['ring-1', cfg.bg, cfg.ring]
  }
  // Permissions mirrors the app-wide permissions button (outlined secondary):
  // border-surface is bound to --p-content-border-color, the same token the
  // outlined-secondary Button resolves its border to.
  if (step.type === 'permissions') return ['border', 'border-surface', 'bg-surface']
  return ['ring-1', 'bg-emphasis', 'ring-surface']
}
function stepIconClass(step) {
  if (step.kind) {
    const cfg = kindConfig(step.kind)
    return [cfg.icon, cfg.text]
  }
  if (step.type === 'permissions') return ['pi', step.icon || 'pi-lock', 'text-color']
  return ['pi', step.icon || 'pi-circle', 'text-muted-color']
}
</script>
