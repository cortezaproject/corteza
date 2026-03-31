<template>
  <p class="m-0 whitespace-pre-wrap" v-html="message" />
</template>

<script>
import base from './base.vue'

export default {
  name: 'CPromptNotification',
  extends: base,
  emits: ['submit'],
  mounted() {
    this.$emit('submit', { keep: true })
  },
}
</script>
