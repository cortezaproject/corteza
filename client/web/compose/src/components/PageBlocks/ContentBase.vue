<template>
  <wrap
    v-bind="$props"
    v-on="$listeners"
  >
    <div
      class="rt-content px-3 py-2"
    >
      <p
        :style="{ 'white-space': 'pre-wrap' }"
        v-html="options.body"
      />
    </div>
  </wrap>
</template>
<script>
import base from './base'

export default {
  extends: base,
}
</script>
