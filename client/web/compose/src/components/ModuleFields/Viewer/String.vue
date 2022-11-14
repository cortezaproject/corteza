<template>
  <div>
    <div
      v-if="formatted"
      class="rt-content"
    >
      <p
        :style="{ 'white-space': 'pre-wrap' }"
        :class="{'multiline': field.isMulti || field.options.multiLine }"
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
}
</script>
