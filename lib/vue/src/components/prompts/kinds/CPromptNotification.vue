<template>
  <div>
    <p
      class="mb-0"
      v-html="message"
    />
  </div>
</template>
<script lang="js">
import base from './base.vue'

export default {
  name: 'CPromptNotification',
  extends: base,

  mounted () {
    // Non blocking, return to workflow right away
    // keep indicates if toast will be kept open until auto delay hides it
    this.$emit('submit', { keep: true })
  },
}
</script>
