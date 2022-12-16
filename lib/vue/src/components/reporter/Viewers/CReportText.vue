<template>
  <div
    v-if="options"
    class="rt-content px-3 py-2"
    v-html="options.value"
  />
</template>

<script>
import base from './base.vue'

export default {
  extends: base,
}
</script>
