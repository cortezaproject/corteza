<template>
  <wrap
    v-bind="$props"
    v-on="$listeners"
  >
    <div
      v-if="processing"
      class="d-flex align-items-center justify-content-center h-100"
    >
      <b-spinner />
    </div>

    <div
      v-else
      class="d-flex h-100"
      :class="{ 'p-2': block.style.wrap.kind === 'card' }"
    >
      <c-progress
        :value="value"
        :min="min"
        :max="max"
        :labeled="options.display.showValue"
        :relative="options.display.showRelative"
        :progress="options.display.showProgress"
        :striped="options.display.striped"
        :animated="options.display.animated"
        :variant="options.display.variant"
        :thresholds="options.display.thresholds"
        class="flex-fill h-100"
      />
    </div>
  </wrap>
</template>

<script>
import base from './base'
import { NoID } from '@cortezaproject/corteza-js'
import { components } from '@cortezaproject/corteza-vue'
import { evaluatePrefilter } from 'corteza-webapp-compose/src/lib/record-filter'
const { CProgress } = components

export default {
  components: {
    CProgress,
  },

  extends: base,

  data () {
    return {
      processing: false,

      value: undefined,
      min: undefined,
      max: undefined,
    }
  },

  watch: {
    'record.recordID': {
      immediate: true,
      handler () {
        this.update()
      },
    },

    options: {
      deep: true,
      handler () {
        this.update()
      },
    },
  },

  mounted () {
    this.$root.$on(`refetch-non-record-blocks:${this.page.pageID}`, this.update)
  },

  beforeDestroy () {
    this.$root.$off(`refetch-non-record-blocks:${this.page.pageID}`)
  },

  methods: {
    /**
     * Pulls fresh data from the API
     */
    async update () {
      this.processing = true

      const { namespaceID } = this.namespace || {}

      const additionalOptions = {
        value: {
          filter: evaluatePrefilter(this.options.value.filter, {
            record: this.record,
            recordID: (this.record || {}).recordID || NoID,
            ownerID: (this.record || {}).ownedBy || NoID,
            userID: (this.$auth.user || {}).userID || NoID,
          }),
        },
        minValue: {
          filter: evaluatePrefilter(this.options.minValue.filter, {
            record: this.record,
            recordID: (this.record || {}).recordID || NoID,
            ownerID: (this.record || {}).ownedBy || NoID,
            userID: (this.$auth.user || {}).userID || NoID,
          }),
        },
        maxValue: {
          filter: evaluatePrefilter(this.options.maxValue.filter, {
            record: this.record,
            recordID: (this.record || {}).recordID || NoID,
            ownerID: (this.record || {}).ownedBy || NoID,
            userID: (this.$auth.user || {}).userID || NoID,
          }),
        },
      }

      return this.block.fetch(additionalOptions, this.$ComposeAPI, namespaceID)
        .then(({ value, min = 0, max = 100 }) => {
          this.min = min
          this.max = max
          this.value = value
        }).catch(this.toastErrorHandler(this.$t('progress.fetch-failed')))
        .finally(() => {
          this.processing = false
        })
    },
  },
}
</script>
