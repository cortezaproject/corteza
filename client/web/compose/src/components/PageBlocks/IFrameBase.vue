<template>
  <wrap
    v-bind="$props"
    v-on="$listeners"
  >
    <iframe
      v-if="src"
      class="h-100 w-100 border-0"
      :src="src"
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
}
</script>
