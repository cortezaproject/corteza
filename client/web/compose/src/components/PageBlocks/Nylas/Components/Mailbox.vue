<template>
  <nylas-mailbox
    :id="componentID"
    :access_token="accessToken"
    :query_string="queryString"
    theme="light"
  />
</template>

<script>
import base from './base'
import '@nylas/components-mailbox'

export default {
  extends: base,

  props: {
    prefillValues: {
      type: Object,
      default: () => ({}),
    },
  },

  computed: {
    queryString () {
      const { queryString = '' } = this.prefillValues
      return queryString ? `from=${queryString}` : undefined
    },
  },
}
</script>
