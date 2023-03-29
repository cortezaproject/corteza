<template>
  <grid
    v-if="layout"
    :key="layout.layoutID"
    :blocks="blocks"
    :editable="false"
  >
    <template
      slot-scope="{ boundingRect, block, index }"
    >
      <page-block
        v-bind="{ ...$attrs }"
        :page="page"
        :blocks="page.blocks"
        :block="block"
        :bounding-rect="boundingRect"
        :block-index="index"
        class="p-2"
        v-on="$listeners"
      />
    </template>
  </grid>
</template>
<script>
import { mapGetters } from 'vuex'
import Grid from '../../Common/Grid'
import PageBlock from '../../PageBlocks'
import { compose } from '@cortezaproject/corteza-js'

export default {
  name: 'PublicGrid',

  components: {
    Grid,
    PageBlock,
  },

  props: {
    page: {
      type: compose.Page,
      required: true,
    },
  },

  data () {
    return {
      layouts: [],
      layout: undefined,

      blocks: [],
    }
  },

  computed: {
    ...mapGetters({
      getPageLayouts: 'pageLayout/getByPageID',
    }),
  },

  watch: {
    'page.pageID': {
      immediate: true,
      handler (pageID) {
        this.layouts = this.getPageLayouts(pageID)
        const { layoutID } = this.$route.query

        if (layoutID) {
          this.layout = this.layouts.find(({ pageLayoutID }) => pageLayoutID === layoutID)
        } else {
          this.layout = this.layouts[0]
          this.$router.replace({ ...this.$route, query: { ...this.$route.query, layoutID: this.layout.pageLayoutID } })
        }

        this.blocks = (this.layout || {}).blocks.map(({ blockID, xywh }) => {
          const block = this.page.blocks.find(b => b.blockID === blockID)
          block.xywh = xywh
          return block
        })
      },
    },
  },
}
</script>
