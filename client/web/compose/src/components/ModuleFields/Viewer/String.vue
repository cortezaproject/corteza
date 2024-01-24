<template>
  <div>
    <div
      v-if="formatted"
      class="rt-content"
    >
      <p
        :style="{ 'white-space': field.options.useRichTextEditor && 'pre-line' }"
        :class="[ 'multiline' && field.isMulti || field.options.multiLine, ...classes ]"
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
      const { fieldID } = this.field
      const { textStyles = {} } = this.extraOptions

      if (this.field.isMulti || this.field.options.multiLine) {
        classes.push('multiline')
      } else if (textStyles.noWrapFields && textStyles.noWrapFields.includes(fieldID)) {
        classes.push('text-nowrap')
      }

      return classes
    },
  },
}
</script>
