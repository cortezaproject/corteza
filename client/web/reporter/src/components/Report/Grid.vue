<template>
  <div
    v-if="layout.length"
    class="w-100"
  >
    <grid-layout
      :layout.sync="layout"
      :col-num="48"
      :row-height="10"
      vertical-compact
      :is-resizable="editable"
      :is-draggable="editable"
      :cols="{ lg: 48, md: 48, sm: 1, xs: 1, xxs: 1 }"
      :margin="[0, 0]"
      :responsive="!editable"
      :use-css-transforms="false"
    >
      <grid-item
        v-for="(item, index) in grid"
        :key="item.i"
        ref="items"
        :min-w="3"
        :min-h="3"
        :i="item.i"
        :h="item.h"
        :w="item.w"
        :x="item.x"
        :y="item.y"
        :class="{ 'editable-grid-item': editable }"
        drag-ignore-from=".gutter"
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
    </grid-layout>
  </div>
  <div
    v-else
    class="d-flex align-items-center justify-content-center h-50 w-100"
  >
    <h4>
      {{ $t('builder:no-blocks-added') }}
    </h4>
  </div>
</template>

<script>
import { GridLayout, GridItem } from 'vue-grid-layout'

export default {
  name: 'Grid',

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
      default: false,
    },
  },

  data () {
    return {
      // All blocks in vue-grid friendly structure
      grid: undefined,

      // Grid items bounding rect info
      boundingRects: [],
    }
  },

  computed: {
    layout: {
      get () {
        return this.grid || this.blocks
      },

      set (layout) {
        // Only update parent blocks when editable to avoid unnecessary updates
        if (this.editable) {
          this.$emit('update:blocks', layout)
        } else {
          this.grid = layout
        }
      },
    },
  },

  watch: {
    blocks: {
      immediate: true,
      deep: true,
      handler (blocks) {
        if (this.editable) {
          this.grid = blocks
        }
      },
    },
  },

  methods: {
    onBlockUpdated (index) {
      this.$emit('item-updated', index)
    },
  },
}
</script>

<style lang="scss">
.vue-grid-item.vue-grid-placeholder {
  background: $primary !important;
}
</style>

<style lang="scss" scoped>
.editable-grid-item {
  touch-action: none;
  background-image: linear-gradient(45deg, $gray-200 25%, $white 25%, $white 50%, $gray-200 50%, $gray-200 75%, $white 75%, $white 100%);
  background-size: 28.28px 28.28px;
}
</style>
