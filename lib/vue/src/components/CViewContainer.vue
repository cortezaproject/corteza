<template>
  <div v-if="scroll" class="flex-1 min-h-0 min-w-0 overflow-y-auto">
    <div :class="[column, 'flex flex-col', gapClass]" v-bind="$attrs">
      <slot />
    </div>
  </div>
  <div v-else :class="[column, 'h-full overflow-hidden']" v-bind="$attrs">
    <slot />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

// Top-level page wrapper shared by section views: a centred column with the
// page padding and a width cap. Two shapes: the default is fixed-height and
// its content scrolls internally (a CResourceList that scrolls itself); the
// `scroll` shape grows to fill the page and scrolls as a whole (edit forms
// stacked in a column). The scroller is the full-width element and the cap
// applies to the column inside it, so a wheel over the gutters still scrolls.
// `width` picks the cap: `list` for tables and grids, `form` for editors;
// it defaults by shape. `gap` only applies to the `scroll` shape. The caps
// are px because the root font-size is 15px, which shrinks rem scales.
defineOptions({ inheritAttrs: false })

const props = defineProps({
  scroll: { type: Boolean, default: false },
  gap: { type: [Number, String], default: 4 },
  width: { type: String as () => 'list' | 'form' | undefined, default: undefined },
})

const CAP = { list: 'max-w-[1440px]', form: 'max-w-[1024px]' }

const column = computed(() => {
  const width = props.width ?? (props.scroll ? 'form' : 'list')
  return `w-full ${CAP[width]} mx-auto p-4 md:p-6 min-w-0`
})

const gapClass = computed(() => (String(props.gap) === '5' ? 'gap-5' : 'gap-4'))
</script>
