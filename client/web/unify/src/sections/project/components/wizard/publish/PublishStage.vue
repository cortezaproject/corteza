<template>
  <section
    class="rounded-xl border bg-surface overflow-hidden transition-colors"
    :class="[current ? 'border-primary' : 'border-surface', locked ? 'opacity-60' : '']"
  >
    <button
      type="button"
      class="w-full flex items-center gap-3 px-4 py-3 text-left"
      :class="locked ? 'cursor-not-allowed' : 'cursor-pointer'"
      :aria-expanded="open"
      :disabled="locked"
      @click="$emit('toggle')"
    >
      <!-- The index is the stage's real position in a real sequence, not
           decoration: a revision cannot be approved while its data decisions
           are open. Done stages keep the number rather than swapping in a
           checkmark, so the sequence stays readable at a glance. -->
      <span
        class="inline-flex items-center justify-center w-7 h-7 rounded-md shrink-0 text-sm font-medium ring-1"
        :class="badgeClass"
      >
        {{ index }}
      </span>

      <span class="flex-1 min-w-0">
        <span class="block font-medium truncate">{{ title }}</span>
        <span class="block text-sm text-muted-color truncate">{{ hint }}</span>
      </span>

      <!-- Readiness is the caller's word: every stage fills this with the
           section's shared StatusChip (see PublishTab), so a stage tag can
           never drift in colour from the same status elsewhere. A stage with
           nothing to state (Go live) simply leaves it empty. -->
      <slot name="status" />
      <i
        class="pi pi-chevron-right text-muted-color transition-transform"
        :class="{ 'rotate-90': open }"
      />
    </button>

    <div v-if="open" class="border-t border-surface px-4 py-4">
      <slot />
    </div>
  </section>
</template>

<script setup>
// The collapsible shell every Publish stage shares — numbering, readiness
// chip, locked/current styling and open state. Purely presentational: it owns
// no publish state, so the stages stay independently buildable (same rule the
// manage/ sections follow).
import { computed } from 'vue'

const props = defineProps({
  index: { type: Number, required: true },
  title: { type: String, required: true },
  hint: { type: String, default: '' },
  // Reached, and the thing to act on right now.
  current: { type: Boolean, default: false },
  // Settled — nothing left to do here.
  done: { type: Boolean, default: false },
  // Not reachable yet; shown dimmed rather than hidden so the shape of what is
  // coming is visible from the start.
  locked: { type: Boolean, default: false },
  open: { type: Boolean, default: false },
})

defineEmits(['toggle'])

const badgeClass = computed(() => {
  if (props.done) return 'bg-green-500/10 text-green-500 ring-green-500/30'
  if (props.current) return 'bg-primary text-primary-contrast ring-primary'
  return 'bg-emphasis text-muted-color ring-surface'
})
</script>
