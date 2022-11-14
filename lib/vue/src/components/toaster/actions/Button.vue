<template>
  <b-button
    v-bind="options"
    @click="$emit('action')"
  >
    <span v-html="label" />
  </b-button>
</template>

<script lang="js">
import base from './base.vue'

export default {
  extends: base,
}
</script>
