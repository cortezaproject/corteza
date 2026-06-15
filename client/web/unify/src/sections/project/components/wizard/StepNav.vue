<template>
  <nav class="w-full h-full overflow-y-auto rounded-xl border border-surface bg-surface-50 dark:bg-surface-950">
    <ul class="p-3 flex flex-col gap-0.5">
      <template v-for="step in steps" :key="step.key">
        <li>
          <button
            type="button"
            class="w-full text-left rounded-md px-3 py-2 flex items-center gap-2 text-sm transition-colors"
            :class="
              step.key === activeKey
                ? 'bg-primary/10 text-primary font-medium'
                : 'hover:bg-surface-100 dark:hover:bg-surface-800'
            "
            @click="$emit('select', step.key)"
          >
            <i :class="['pi', icon(statuses[step.key]).name, icon(statuses[step.key]).class]" />
            <span class="flex-1 min-w-0 truncate">{{ $t(step.labelKey) }}</span>
            <i v-if="step.gatedOnly" class="pi pi-shield text-[10px] text-amber-500" :title="$t('project.wizard.gate.gatedOnlyStep')" />
          </button>
        </li>

        <!-- Clickable gate (Governance tab only) -->
        <li v-if="showGates && step.gate">
          <button
            type="button"
            class="w-full rounded-md px-3 py-1.5 flex items-center gap-2 text-[11px] uppercase tracking-wide transition-colors"
            :class="
              gateLocked[step.key]
                ? 'opacity-50 cursor-not-allowed'
                : 'hover:bg-surface-100 dark:hover:bg-surface-800'
            "
            :disabled="gateLocked[step.key]"
            :title="gateLocked[step.key] ? $t('project.wizard.gate.lockedTitle') : ''"
            @click="$emit('gate-click', step.key)"
          >
            <span class="h-px flex-1 bg-surface-200 dark:bg-surface-700" />
            <i
              v-if="gateLocked[step.key]"
              class="pi pi-lock text-[10px] text-surface-400"
            />
            <i
              v-else
              :class="['pi text-[10px]', icon(gateStatuses[step.key]).name, icon(gateStatuses[step.key]).class]"
            />
            <span class="text-muted-color">{{ $t('project.wizard.gate.label', { number: gateNumber[step.key] }) }}</span>
            <span class="h-px flex-1 bg-surface-200 dark:bg-surface-700" />
          </button>
        </li>
      </template>
    </ul>
  </nav>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  steps: { type: Array, required: true },
  activeKey: { type: String, default: '' },
  statuses: { type: Object, default: () => ({}) },
  gateStatuses: { type: Object, default: () => ({}) },
  gateLocked: { type: Object, default: () => ({}) },
  showGates: { type: Boolean, default: false },
})
defineEmits(['select', 'gate-click'])

// Sequential 1-based number for each gate, in step order.
const gateNumber = computed(() => {
  const out = {}
  let n = 0
  for (const step of props.steps) {
    if (step.gate) out[step.key] = ++n
  }
  return out
})

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
