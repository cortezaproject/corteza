<template>
  <grid
    :blocks="blocks"
    :editable="false"
  >
    <template
      slot-scope="{ blockIndex, block }"
    >
      <page-block
        v-bind="{ ...$attrs }"
        :page="page"
        :blocks="blocks"
        :block="block"
        :block-index="blockIndex"
        :loading-record="loadingRecord"
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

    blocks: {
      type: Array,
      required: true,
    },

    loadingRecord: {
      type: Boolean,
      default: false,
    },
  },
}
</script>
