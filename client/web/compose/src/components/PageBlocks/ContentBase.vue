<template>
  <wrap
    v-bind="$props"
    v-on="$listeners"
  >
    <div
      class="rt-content px-3 py-2"
    >
      <p
        :style="{ 'white-space': 'pre-wrap' }"
        v-html="contentBody"
      />
    </div>
  </wrap>
</template>
<script>
import base from './base'
import { evaluatePrefilter } from 'corteza-webapp-compose/src/lib/record-filter'
import { NoID } from '@cortezaproject/corteza-js'

export default {
  extends: base,

  computed: {
    contentBody () {
      const { body = '' } = this.options

      return evaluatePrefilter(body, {
        record: this.record,
        recordID: (this.record || {}).recordID || NoID,
        ownerID: (this.record || {}).ownedBy || NoID,
        userID: (this.$auth.user || {}).userID || NoID,
      })
    },
  },
}
</script>
