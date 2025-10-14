<template>
  <div>
    <div
      v-if="formatted"
      class="rt-content"
    >
      <p
        v-if="field.options.useRichTextEditor"
        :style="{ 'white-space': field.options.useRichTextEditor && 'pre-line' }"
        :class="[ 'multiline' && field.isMulti || field.options.multiLine, ...classes ]"
        v-html="formatted"
      />

      <p
        v-else
        :style="{ 'white-space': field.options.useRichTextEditor && 'pre-line' }"
        :class="[ 'multiline' && field.isMulti || field.options.multiLine, ...classes ]"
      >
        {{ formatted }}
      </p>
    </div>
  </div>
</template>
<script>
import base from './base'

export default {
  extends: base,
}
</script>
