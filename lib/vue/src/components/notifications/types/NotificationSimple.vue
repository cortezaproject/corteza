<template>
  <div class="d-flex flex-column">
    <h5 class="font-weight-bold text-break">
      {{ title }}
    </h5>

    <div
      class="rt-content text-secondary mb-1 text-break"
      v-html="description"
    />
  </div>
</template>

<script>
export default {
  props: {
    notification: {
      type: Object,
      required: true,
    },
  },

  computed: {
    title () {
      if (!this.notification || !this.notification.config) {
        return ''
      }

      return this.notification.config.title
    },

    description () {
      if (!this.notification || !this.notification.config) {
        return ''
      }

      return this.notification.config.description || ''
    },
  },
}
</script>
