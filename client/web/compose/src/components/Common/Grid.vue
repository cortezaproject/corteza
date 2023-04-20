<template>
  <div
    v-if="layout.length"
    class="w-100"
    :class="{
      'editable': editable,
      'flex-grow-1 d-flex': isStretchable,
    }"
  >
    <grid-layout
      :layout.sync="layout"
      :col-num="48"
      :row-height="10"
      vertical-compact
      :is-resizable="editable"
      :is-draggable="editable"
      :cols="columnNumber"
      :margin="[0, 0]"
      :responsive="!editable"
      :use-css-transforms="false"
      class="flex-grow-1 d-flex w-100 h-100"
    >
      <template
        v-for="(item, index) in layout"
      >
        <grid-item
          v-if="!blocks[item.i].meta.hidden"
          :key="item.i"
          ref="items"
          :i="item.i"
          :h="item.h"
          :w="item.w"
          :x="item.x"
          :y="item.y"
          :min-w="6"
          :min-h="5"
          :class="{ 'h-100': isStretchable }"
          class="grid-item"
          style="touch-action: none;"
          @moved="onBlockUpdated(index)"
          @resized="onBlockUpdated(index)"
        >
          <slot
            :block="blocks[item.i]"
            :index="index"
            :block-index="item.i"
            :bounding-rect="boundingRects[index]"
            v-on="$listeners"
          />
        </grid-item>
      </template>
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
import { GridLayout, GridItem } from 'vue-grid-layout'
import { throttle } from 'lodash'
import { compose } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'page',
  },

  components: {
    GridLayout,
    GridItem,
  },

  props: {
    blocks: {
      type: Array,
      default: () => ([]),
    },

    editable: {
      type: Boolean,
    },
  },

  data () {
    return {
      // all blocks in vue-grid friendly structure
      layout: [],

      // Grid items bounding rect info
      boundingRects: [],
    }
  },

  computed: {
    oneBlockLayout () {
      return this.blocks.filter(({ meta }) => !meta.hidden).length === 1
    },

    isStretchable () {
      return !this.editable && this.oneBlockLayout
    },

    columnNumber () {
      if (this.oneBlockLayout) {
        return { lg: 1, md: 1, sm: 1, xs: 1, xxs: 1 }
      }
      return { lg: 48, md: 48, sm: 1, xs: 1, xxs: 1 }
    },
  },

  watch: {
    blocks: {
      immediate: true,
      deep: true,
      handler (blocks) {
        blocks = blocks.map(({ meta, xywh: [x, y, w, h] }, i) => {
          // To avoid collision with hidden elements
          return meta.hidden ? { i, x: 0, y: 0, w: 0, h: 0 } : { i, x, y, w, h }
        })

        // Only use next tick if previous layout exists and not editing
        if (this.layout.length && !this.editable) {
          this.layout = []

          // Next tick is important, otherwise it can lead to overlapping blocks
          this.$nextTick(() => {
            this.layout = blocks
          })
        } else {
          this.layout = blocks
        }
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

    onBlockUpdated (index) {
      this.$emit('item-updated', index)
      this.$emit('update:blocks', this.layout.map(
        ({ x, y, w, h, i }) => new compose.PageBlockMaker({ ...this.blocks[i], xywh: [x, y, w, h] }),
      ))
    },
  },
}
</script>
<style lang="scss">
.editable {
  .grid-item {
    background-image: linear-gradient(45deg, $gray-200 25%, $white 25%, $white 50%, $gray-200 50%, $gray-200 75%, $white 75%, $white 100%);
    background-size: 28.28px 28.28px;
  }
}
</style>
