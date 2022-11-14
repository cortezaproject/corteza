<template>
  <b-button
    variant="link"
    class="text-dark font-weight-bold text-decoration-none"
    @click="onClick(format.type, format.attrs)">
    <span :class="activeClasses(format.attrs)">
      <font-awesome-icon
        v-if="format.icon"
        :icon="format.icon" />

      <span v-else>
        {{ format.label }}
      </span>
    </span>
  </b-button>
</template>

<script>
import base from './base.vue'

/**
 * Component is used to display simple formatter options such as bold, italic, ...
 */
export default {
  name: 't-mark-item',
  extends: base,
}
</script>
