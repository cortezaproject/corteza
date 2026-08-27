<template>
  <div
    class="sr-only"
    :aria-live="assertive ? 'assertive' : 'polite'"
    aria-atomic="true"
    data-c311-status-announcer
  >
    {{ message }}
  </div>
</template>

<script>
export default {
  name: 'C311StatusAnnouncer',
  props: {
    message: {
      type: String,
      default: '',
    },
    assertive: {
      type: Boolean,
      default: false,
    },
  },
}
</script>
