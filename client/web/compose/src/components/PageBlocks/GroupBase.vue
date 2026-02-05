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
        :margin="[0, 0]"
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
    }
  },

  computed: {
    groupedBlocks () {
      return this.block.options.blocks.map(({ blockID }) => {
        const unparsedBlock = blockID ? this.blocks.find(b => fetchID(b) === blockID) : undefined
        if (!unparsedBlock) return undefined

        const block = JSON.parse(JSON.stringify(unparsedBlock))
        // Set some default styling for nested blocks
        block.style.wrap.kind = 'Plain'
        block.style.border.enabled = false
        return compose.PageBlockMaker(block)
      }).filter(b => !!b)
    },
  },

  watch: {
    'block.options.blocks': {
      handler (blocks) {
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

  methods: {
    onLayoutUpdated (layout) {
      if (!this.editable) return

      this.block.options.blocks.forEach((b, idx) => {
        const { x, y, w, h } = layout[idx]
        b.xywh = [x, y, w, h]
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
