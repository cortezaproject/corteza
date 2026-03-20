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

// Map old Bootstrap variant names to PrimeVue CSS classes
const headerTextClass = computed(() => {
  const variant = props.block.style?.variants?.headerText
  const map = {
    dark: 'text-color',
    primary: 'text-primary',
    secondary: 'text-muted-color',
    success: 'text-green-500',
    warning: 'text-orange-500',
    danger: 'text-red-500',
  }
  return map[variant] || ''
})

const hasBorder = computed(() => {
  return props.block.style?.border?.enabled
})

const cardPt = computed(() => ({
  root: {
    class: [
      isPlain.value ? 'bg-transparent shadow-none' : '',
      hasBorder.value ? 'border border-surface' : '',
    ],
  },
  body: { class: 'p-0 flex-1 flex flex-col overflow-hidden gap-0' },
  caption: { class: ['pl-3 py-3 border-b border-surface gap-1', headerTextClass.value] },
  content: { class: 'p-0 flex-1 flex flex-col overflow-hidden' },
}))
</script>
