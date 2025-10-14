<template>
  <div class="rt-content">
    <p
      v-if="formatted && field.options.useRichTextEditor"
      :style="{ 'white-space': 'pre-line' }"
      :class="[ 'multiline' && field.isMulti || field.options.multiLine, ...classes ]"
      v-html="formatted"
    />

    <p
      v-else-if="formatted"
      :class="[ 'multiline' && field.isMulti || field.options.multiLine, ...classes ]"
    >
      {{ formatted }}
    </p>
  </div>
</template>
<script>
import base from './base'

export default {
  extends: base,
}
</script>
