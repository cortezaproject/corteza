<template>
  <grid
    v-if="page.blocks"
    :key="page.pageID"
    :blocks="page.blocks"
    :editable="false"
  >
    <template
      slot-scope="{ boundingRect, block, index }"
    >
      <page-block
        v-bind="{ ...$attrs, block, page, boundingRect, blockIndex: index }"
        class="p-2"
        v-on="$listeners"
      />
    </template>
  </grid>
</template>
<script>
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
}
</script>
