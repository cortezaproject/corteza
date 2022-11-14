<template>
  <div class="overflow-hidden h-100">
    <slot name="toolbar" />
    <slot name="default" />
    <slot name="footer" />
  </div>
</template>
<script>
import { compose } from '@cortezaproject/corteza-js'

export default {
  props: {
    block: {
      type: compose.PageBlock,
      required: true,
    },

    boundingRect: {
      type: Object,
      required: false,
      default: undefined,
    },
  },

  computed: {
    blockClass () {
      return [
        'block',
        this.block.kind,
      ]
    },
  },
}
</script>
