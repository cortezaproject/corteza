<template>
  <div
    v-if="grid.length"
    class="w-100"
    :class="{
      editable: !!editable,
      'flex-grow-1 d-flex': isStretchable,
    }"
  >
    <grid-layout
      class="flex-grow-1 d-flex w-100 h-100"
      :layout.sync="layout"
      :row-height="50"
      :is-resizable="!!editable"
      :is-draggable="!!editable"
      :cols="columnNumber"
      :margin="[0, 0]"
      :responsive="!editable"
      :use-css-transforms="false"
      @layout-updated="handleLayoutUpdate"
    >
      <grid-item
        v-for="(item, index) in gridCollection"
        :key="item.i"
        ref="items"
        class="grid-item"
        :class="{
          'h-100': isStretchable,
        }"
        style="touch-action: none;"
        v-bind="{ ...item }"
      >
        <slot
          :block="blocks[item.i]"
          :index="index"
          :block-index="item.i"
          :bounding-rect="boundingRects[index]"
          v-on="$listeners"
        />
      </grid-item>
    </grid-layout>
  </div>
  <div
    v-else
    class="no-builder-grid h-100 pt-5 container text-center"
  >
    <h4>
      {{ $t('noBlock') }}
    </h4>
  </div>
</template>

<script>
import VueGridLayout from 'vue-grid-layout'
import { compose } from '@cortezaproject/corteza-js'
import { throttle } from 'lodash'

export default {
  i18nOptions: {
    namespaces: 'page',
  },

  components: {
    GridLayout: VueGridLayout.GridLayout,
    GridItem: VueGridLayout.GridItem,
  },

  props: {
    editable: {
      type: Boolean,
    },

    blocks: {
      type: Array,
      required: true,
    },
  },

  data () {
    return {
      // all blocks in vue-grid friendly structure
      grid: [],

      // Grid items bounding rect info
      boundingRects: [],
    }
  },

  computed: {
    layout: {
      get () {
        return this.grid
      },

      set (layout) {
        this.grid = layout
        this.handleLayoutUpdate(layout)
      },
    },

    sortedGrid () {
      return Array.from(this.grid).sort((a, b) => Math.sqrt(a.x * a.x + a.y * a.y) - Math.sqrt(b.x * b.x + b.y * b.y))
    },

    isStretchable () {
      if (this.editable) {
        // When in-edit mode do not stretch the blocks
        return false
      }

      const minHeight = 10
      let heightCheck = -1

      for (let b = 0; b < this.blocks.length; b++) {
        const { xywh: [, y, , h] } = this.blocks[b]

        if (y > 0) {
          // If block is not positioned at the top,
          // do not try to make it stretchable
          return false
        }

        if (heightCheck === -1) {
          // Set block height for the next check
          heightCheck = h
        }

        if (heightCheck !== h && minHeight > h) {
          // Not full height
          return false
        }
      }

      return true
    },

    columnNumber () {
      if (this.grid.length === 1) {
        return { lg: 1, md: 1, sm: 1, xs: 1, xxs: 1 }
      }
      return { lg: 12, md: 12, sm: 1, xs: 1, xxs: 1 }
    },

    gridCollection () {
      if (this.grid.length === 1) {
        return this.sortedGrid
      }
      return this.grid
    },
  },

  watch: {
    blocks: {
      immediate: true,
      deep: true,
      handler (blocks) {
        if (blocks.length === 0) this.$emit('change', [])
        this.grid = blocks.map(({ xywh: [x, y, w, h] }, i) => ({ i, x, y, w, h }))
      },
    },
  },

  mounted () {
    window.addEventListener('resize', this.windowResizeThrottledHandler)
    this.recalculateBoundingRect()
  },

  destroyed () {
    window.removeEventListener('resize', this.windowResizeThrottledHandler)
  },

  methods: {
    windowResizeThrottledHandler: throttle(function () { this.recalculateBoundingRect() }, 500),

    // Fetch bounding boxes of all grid items
    recalculateBoundingRect () {
      this.boundingRects = (this.$refs.items || [])
        .map(({ $el }) => {
          const { x, y, width: w, height: h } = $el.getBoundingClientRect()
          return { x, y, w, h }
        })
    },

    handleLayoutUpdate (layout) {
      this.$emit('change', layout.map(
        ({ x, y, w, h, i }) => new compose.PageBlockMaker({ ...this.blocks[i], xywh: [x, y, w, h] }),
      ))
      this.recalculateBoundingRect()
    },
  },
}
</script>
<style lang="scss">
.editable {
  .grid-item {
    background-image: linear-gradient(45deg, #f3f3f5 25%, #ffffff 25%, #ffffff 50%, #f3f3f5 50%, #f3f3f5 75%, #ffffff 75%, #ffffff 100%);
    background-size: 28.28px 28.28px;
  }
}
</style>
