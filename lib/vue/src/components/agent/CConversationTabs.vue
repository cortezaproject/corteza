<template>
  <div class="flex items-center gap-0 border-b border-surface shrink-0 bg-emphasis">
    <div class="flex items-center gap-0 flex-1 overflow-x-auto no-scrollbar">
      <button
        v-for="(tab, idx) in tabs"
        :key="idx"
        type="button"
        class="group flex items-center gap-1.5 py-2 pl-3 pr-1 border-b border-r border-r-surface whitespace-nowrap transition-colors duration-200 outline-none select-none max-w-[150px]"
        :class="[
          activeIndex === idx
            ? 'border-b-primary text-primary font-medium'
            : 'border-b-transparent text-muted-color hover:text-color hover:border-b-surface',
        ]"
        @click="emit('update:activeIndex', idx)"
      >
        <span class="whitespace-nowrap truncate">{{ tab.label }}</span>
        <i
          class="pi pi-times text-xs p-1 hover:bg-emphasis rounded-full transition-all shrink-0"
          :class="closeIconClass"
          @click.stop="onClose(idx)"
        />
      </button>
    </div>
    <div class="flex items-center px-1 border-l border-surface shrink-0">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

// Shared conversation tab strip used by both the agent editor's split-view
// chat header and the sidebar agent chat. Markup/classes here are shared
// verbatim between both call sites — keep it that way.
const props = defineProps({
  tabs: {
    type: Array as () => Array<{ label: string }>,
    default: () => [],
  },
  activeIndex: {
    type: Number,
    default: 0,
  },
  // When true, the close icon is hidden (and closing disabled) whenever only
  // a single tab remains, so the user is never left without a conversation.
  // When false (default), the close icon is always available on hover.
  guardCloseWhenLast: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits<{
  (_e: 'update:activeIndex', _idx: number): void
  (_e: 'close', _idx: number): void
}>()

const canClose = computed(() => !props.guardCloseWhenLast || props.tabs.length > 1)

const closeIconClass = computed(() =>
  canClose.value ? 'opacity-0 group-hover:opacity-100' : 'invisible',
)

function onClose(idx: number) {
  if (!canClose.value) return
  emit('close', idx)
}
</script>
