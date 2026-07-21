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
          <!-- Gated projects show governance status; free build shows each
               step's own icon in the same badge used by the metrics strip and
               resource item lists (resource steps reuse the kind icon/colour). -->
          <i
            v-if="gated"
            :class="['pi', icon(statuses[step.key]).name, icon(statuses[step.key]).class]"
          />
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

defineProps({
  steps: { type: Array, required: true },
  activeKey: { type: String, default: '' },
  statuses: { type: Object, default: () => ({}) },
  // Gated projects render governance-status icons; free build renders per-step icons.
  gated: { type: Boolean, default: false },
})
defineEmits(['select'])

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

function icon(status) {
  switch (status) {
    case 'approved':
      return { name: 'pi-check-circle', class: 'text-emerald-500' }
    case 'submitted':
      return { name: 'pi-clock', class: 'text-sky-500' }
    case 'changes-requested':
      return { name: 'pi-exclamation-circle', class: 'text-amber-500' }
    default:
      return { name: 'pi-circle', class: 'text-surface-400' }
  }
}
</script>
