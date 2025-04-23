<template>
  <div class="d-flex flex-column gap-1">
    <h5 class="font-weight-bold text-break">
      {{ title }}
    </h5>

    <div
      class="text-secondary mb-1 text-break"
    >
      {{ description }}
    </div>
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

<style lang="scss" scoped>
</style>
