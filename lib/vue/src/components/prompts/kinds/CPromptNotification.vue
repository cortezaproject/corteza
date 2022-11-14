<template>
  <div>
    <p
      v-html="message"
      class="mb-0"
    />
  </div>
</template>
<script lang="js">
import base from './base.vue'

export default {
  extends: base,
  name: 'c-prompt-notification',

    mounted () {
    // Non blocking, return to workflow right away
    // keep indicates if toast will be kept open until auto delay hides it
    this.$emit('submit', { keep: true })
  }
}
</script>
