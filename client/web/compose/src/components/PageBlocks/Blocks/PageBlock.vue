<template>
  <Card class="h-full overflow-hidden" :pt="cardPt">
    <template v-if="$slots.title || block.title" #title>
      <slot name="title">
        {{ block.title }}
      </slot>
    </template>
    <template v-if="$slots.subtitle || block.description" #subtitle>
      <slot name="subtitle">
        {{ block.description }}
      </slot>
    </template>
    <template #content>
      <slot />
    </template>
    <template v-if="$slots.footer" #footer>
      <slot name="footer" />
    </template>
  </Card>
</template>

<script setup>
import { computed, useSlots } from 'vue'

const $slots = useSlots()

const props = defineProps({
  block: {
    type: Object,
    required: true,
  },
})

const isPlain = computed(() => {
  return props.block.style?.wrap?.kind !== 'card'
})

const cardPt = computed(() => ({
  root: { class: isPlain.value ? 'bg-transparent shadow-none' : '' },
  body: { class: 'p-0 flex-1 flex flex-col overflow-hidden gap-0' },
  caption: { class: 'pl-3 py-3 border-b border-surface gap-1' },
  content: { class: 'p-0 flex-1 flex flex-col overflow-hidden' },
}))
</script>
