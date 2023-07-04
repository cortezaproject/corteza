<template>
  <nylas-composer
    :id="componentID"
    ref="composer"
    :access_token="accessToken"
    theme="light"
    show_header="false"
  />
</template>

<script>
import base from './base'
import '@nylas/components-composer'

export default {
  extends: base,

  props: {
    prefillValues: {
      type: Object,
      default: () => ({}),
    },
  },

  watch: {
    prefillValues: {
      deep: true,
      immediate: true,
      handler (value) {
        this.$nextTick(() => {
          const el = this.$refs.composer
          // Set initial Values
          el.value = value
        })
      },
    },
  },
}
</script>
