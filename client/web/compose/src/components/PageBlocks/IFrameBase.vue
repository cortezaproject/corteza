<template>
  <wrap
    v-bind="$props"
    v-on="$listeners"
    @refreshBlock="refresh"
  >
    <iframe
      v-if="src"
      ref="iframe"
      class="h-100 w-100 border-0"
      :src="src | checkValidURL"
    />
  </wrap>
</template>
<script>
import base from './base'

export default {
  extends: base,

  computed: {
    src () {
      const { srcField, src } = this.options
      const blank = 'about:blank'

      if (this.options.srcField) {
        if (this.record) {
          return this.record.values[srcField] || blank
        }
      }

      return src || blank
    },
  },

  mounted () {
    this.refreshBlock(this.refresh)
  },

  methods: {
    refresh () {
      this.$refs.iframe.src = this.src
    },
  },
}
</script>
