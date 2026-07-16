<template>
  <div>
    <div
      class="flex items-center justify-between px-4 py-2 cursor-pointer select-none sticky top-0 z-10"
      :class="subgroup ? 'bg-surface pl-6' : 'bg-emphasis'"
      @click="toggle"
    >
      <div class="flex items-center gap-2 min-w-0">
        <i
          class="pi text-xs text-muted-color shrink-0"
          :class="isExpanded ? 'pi-chevron-down' : 'pi-chevron-right'"
        />
        <span
          class="font-semibold truncate"
          :class="subgroup ? 'text-xs text-muted-color' : 'text-sm text-color'"
        >{{ title }}</span>
      </div>
      <Badge
        :value="items.length"
        severity="secondary"
        class="shrink-0"
      />
    </div>

    <div v-if="isExpanded">
      <slot />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  title: {
    type: String,
    required: true,
  },
  items: {
    type: Array,
    default: () => [],
  },
  collapseId: {
    type: String,
    required: true,
  },
  expanded: {
    type: Boolean,
    default: true,
  },
  subgroup: {
    type: Boolean,
    default: false,
  },
  labels: {
    type: Object,
    default: () => ({}),
  },
})

const emit = defineEmits(['update:expanded'])

const isExpanded = computed(() => props.expanded)

const toggle = () => {
  emit('update:expanded', !props.expanded)
}
</script>
