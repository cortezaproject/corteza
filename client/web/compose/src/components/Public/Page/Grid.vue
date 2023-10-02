<template>
  <grid
    :blocks="blocks"
    :editable="false"
  >
    <template
      slot-scope="{ index, block }"
    >
      <page-block
        v-bind="{ ...$attrs }"
        :page="page"
        :blocks="blocks"
        :block="block"
        :block-index="index"
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
  },
}
</script>
