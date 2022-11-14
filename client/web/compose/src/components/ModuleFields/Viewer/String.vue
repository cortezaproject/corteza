<template>
  <div
    v-if="formatted"
    class="rt-content"
  >
    <p
      :class="{'multiline': field.isMulti || field.options.multiLine }"
      v-html="formatted"
    />

    <errors :errors="errors" />
  </div>
</template>
<script>
import base from './base'

export default {
  extends: base,
}
</script>
