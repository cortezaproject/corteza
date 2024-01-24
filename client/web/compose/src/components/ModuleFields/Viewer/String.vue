<template>
  <div>
    <div
      v-if="formatted"
      class="rt-content"
    >
      <p
        :style="{ 'white-space': 'pre-wrap' }"
        :class="classes"
        v-html="formatted"
      />
    </div>

    <errors :errors="errors" />
  </div>
</template>
<script>
import base from './base'

export default {
  extends: base,

  computed: {
    classes () {
      const classes = []

      if (this.field.isMulti || this.field.options.multiLine) {
        classes.push('multiline')
      }

      if (this.extraOptions.configureTextWrap) {
        classes.push(this.extraOptions.configureTextWrap)
      }

      return classes
    },
  },
}
</script>
