<template>
  <div
    class="inline-flex items-center gap-2 rounded-lg ring-1 bg-surface font-medium text-sm select-none"
    :class="[
      cfg.ring,
      removable ? 'pl-3 pr-1.5 py-1.5' : 'px-3 py-1.5',
      clickable ? 'cursor-pointer hover:bg-emphasis transition-colors' : '',
    ]"
    @click="clickable && $emit('click')"
  >
    <i v-if="!hideIcon" :class="[cfg.icon, cfg.text, 'shrink-0']" />
    <span class="whitespace-nowrap">{{ resource.name }}</span>
    <i
      v-if="removable"
      class="pi pi-times text-xs text-muted-color hover:text-color cursor-pointer ml-0.5"
      role="button"
      aria-label="Remove"
      @click.stop="$emit('remove')"
    />
  </div>
</template>

<script setup>
import { kindConfig } from '@/sections/project/config/kinds'
import { computed } from 'vue'

defineOptions({ name: 'ResourceBadge' })

const props = defineProps({
  resource: { type: Object, required: true },
  removable: { type: Boolean, default: false },
  clickable: { type: Boolean, default: false },
  hideIcon: { type: Boolean, default: false },
})

defineEmits(['click', 'remove'])

const cfg = computed(() => kindConfig(props.resource.kind))
</script>
