<template>
  <wrap
    v-bind="$props"
    v-on="$listeners"
    @refreshBlock="refresh"
  >
    <iframe
      v-if="src"
      ref="iframe"
      class="h-100 w-100 border-0"
      :src="src | checkValidURL"
    />
  </wrap>
</template>
<script>
import base from './base'
import { NoID } from '@cortezaproject/corteza-js'
import { evaluatePrefilter, isFieldInFilter } from 'corteza-webapp-compose/src/lib/record-filter'

export default {
  extends: base,

  computed: {
    src () {
      const { srcField, src } = this.options
      const blank = 'about:blank'
      let url = src

      if (this.options.srcField) {
        if (this.record) {
          url = this.record.values[srcField]
        }
      }

      const interpolatedURL = evaluatePrefilter(url, {
        record: this.record,
        user: this.$auth.user || {},
        recordID: (this.record || {}).recordID || NoID,
        ownerID: (this.record || {}).ownedBy || NoID,
        userID: (this.$auth.user || {}).userID || NoID,
      })

      return interpolatedURL || blank
    },
  },

  mounted () {
    this.refreshBlock(this.refresh)
    this.createEvents()
  },

  beforeDestroy () {
    this.destroyEvents()
  },

  methods: {
    refresh () {
      this.$refs.iframe.src = this.src
    },

    createEvents () {
      this.$root.$on('record-field-change', this.refetchOnPrefilterValueChange)
    },

    refetchOnPrefilterValueChange ({ fieldName }) {
      const { src } = this.options

      if (isFieldInFilter(fieldName, src)) {
        this.refresh()
      }
    },

    destroyEvents () {
      this.$root.$off('record-field-change', this.refetchOnPrefilterValueChange)
    },
  },
}
</script>
