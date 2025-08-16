<template>
  <wrap
    v-bind="$props"
    v-on="$listeners"
  >
    <div
      :style="{ 'white-space': 'pre-wrap' }"
      class="rt-content p-3"
      v-html="contentBody"
    />
  </wrap>
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'
import { evaluatePrefilter } from 'corteza-webapp-compose/src/lib/record-filter'
import base from './base'

export default {
  extends: base,

  computed: {
    contentBody () {
      try {
        const { body = '' } = this.options

        return evaluatePrefilter(body, {
          record: this.record,
          user: this.$auth.user || {},
          recordID: (this.record || {}).recordID || NoID,
          ownerID: (this.record || {}).ownedBy || NoID,
          userID: (this.$auth.user || {}).userID || NoID,
        })
      } catch (e) {
        return e
      }
    },
  },
}
</script>
