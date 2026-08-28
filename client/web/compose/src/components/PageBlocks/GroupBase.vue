<template>
  <wrap
    v-bind="$props"
    :scrollable-body="false"
    card-class="group-base-block-container"
    v-on="$listeners"
  >
    <div
      v-if="!groupedBlocks.length"
      class="d-flex h-100 align-items-center justify-content-center"
    >
      <p class="mb-0">
        {{ $t('noBlocks') }}
      </p>
    </div>

    <div
      v-else
      class="h-100 overflow-auto"
    >
      <grid-layout
        v-if="layout.length"
        :layout.sync="layout"
        :col-num="48"
        :row-height="10"
        vertical-compact
        :is-resizable="editable"
        :is-draggable="editable"
        :margin="gridMargin"
        :use-css-transforms="false"
        class="w-100 h-100"
        @layout-updated="onLayoutUpdated"
      >
        <grid-item
          v-for="(item, index) in layout"
          :key="item.i"
          v-bind="item"
          :min-w="1"
          :min-h="1"
          class="grid-item"
        >
          <page-block
            v-if="groupedBlocks[index]"
            v-bind="{ ...$attrs, ...$props, block: groupedBlocks[index], blockIndex: index }"
            header-class="border-0 border-white"
          />
        </grid-item>
      </grid-layout>
    </div>
  </wrap>
</template>

<script>
import base from './base'
import { GridLayout, GridItem } from 'vue-grid-layout'
import { compose } from '@cortezaproject/corteza-js'
import { fetchID } from 'corteza-webapp-compose/src/lib/block'

export default {
  i18nOptions: {
    namespaces: 'block',
    keyPrefix: 'group',
  },

  name: 'GroupBase',

  components: {
    GridLayout,
    GridItem,
    PageBlock: () => import('corteza-webapp-compose/src/components/PageBlocks'),
  },

  extends: base,

  data () {
    return {
      layout: [],
      resizeObserver: null,
      updatingLayout: false,
      lastWidth: 0,
    }
  },

  computed: {
    groupedBlocks () {
      return this.block.options.blocks.map(({ blockID }) => {
        const unparsedBlock = blockID ? this.blocks.find(b => fetchID(b) === blockID) : undefined
        if (!unparsedBlock) return undefined

        const block = JSON.parse(JSON.stringify(unparsedBlock))
        return compose.PageBlockMaker(block)
      }).filter(b => !!b)
    },

    gridMargin () {
      const p = this.block.options.padding || 0
      return [p, p]
    },
  },

  watch: {
    'block.options.blocks': {
      handler (blocks) {
        // Guard against infinite loop: onLayoutUpdated modifies xywh which
        // triggers this deep watcher, which rebuilds layout, which triggers
        // vue-grid-layout to re-layout, which fires layout-updated again.
        if (this.updatingLayout) return

        this.layout = blocks.map(({ blockID, xywh: [x, y, w, h] }, i) => ({
          i: blockID || `temp-${i}`,
          x,
          y,
          w,
          h,
        }))
      },
      immediate: true,
      deep: true,
    },
  },

  mounted () {
    // vue-grid-layout only listens to window resize, not parent container resize.
    // When the Group block is resized in the page builder grid, we need to
    // trigger a recalculation so inner blocks reflow accordingly.
    // Track width to avoid infinite loop (resize dispatch → grid recalc → size change → observer fires).
    this.lastWidth = this.$el.clientWidth
    this.resizeObserver = new ResizeObserver((entries) => {
      const width = entries[0].contentRect.width
      if (Math.abs(width - this.lastWidth) > 1) {
        this.lastWidth = width
        window.dispatchEvent(new Event('resize'))
      }
    })
    this.resizeObserver.observe(this.$el)
  },

  beforeDestroy () {
    if (this.resizeObserver) {
      this.resizeObserver.disconnect()
    }
  },

  methods: {
    onLayoutUpdated (layout) {
      if (!this.editable) return

      this.updatingLayout = true

      // When rendered inside a Tab, this.block is a deep clone — changes to it
      // are discarded on save. Find the original block in this.blocks to persist xywh.
      const originalBlock = this.blocks.find(b => fetchID(b) === fetchID(this.block))
      const targetBlocks = originalBlock ? originalBlock.options.blocks : this.block.options.blocks

      targetBlocks.forEach((b, idx) => {
        if (!layout[idx]) return
        const { x, y, w, h } = layout[idx]
        b.xywh = [x, y, w, h]
      })

      this.$nextTick(() => {
        this.updatingLayout = false
      })
    },
  },
}
</script>

<style lang="scss" scoped>
.group-base-block-container {
  .grid-layout {
    height: auto !important;
    min-height: 100%;
  }

  .grid-item {
    touch-action: none;
  }
}
</style>
