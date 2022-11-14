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
      class="h-100"
      :class="{ 'p-2': block.style.wrap.kind === 'card' }"
    >
      <b-progress
        :max="max"
        class="h-100 bg-light"
      >
        <b-progress-bar
          :value="value"
          :striped="options.display.striped"
          :animated="options.display.animated"
          :variant="progressVariant"
        >
          {{ progressLabel }}
        </b-progress-bar>
      </b-progress>
    </div>
  </wrap>
</template>

<script>
import base from './base'
import { NoID } from '@cortezaproject/corteza-js'
import { evaluatePrefilter } from 'corteza-webapp-compose/src/lib/record-filter'

export default {
  extends: base,

  props: {
  },

  data () {
    return {
      processing: false,

      value: undefined,
      max: undefined,
    }
  },

  computed: {
    progress () {
      const { value = 0, max = 100 } = this
      return 100 * (value / max)
    },

    progressLabel () {
      let { value } = this
      const { showValue, showRelative, showProgress } = this.options.display || {}

      if (!showValue) {
        return
      }

      if (showRelative) {
        // https://stackoverflow.com/a/21907972/17926309
        value = `${Math.round(((value / this.max) * 100) * 100) / 100}%`
      }

      if (showProgress) {
        value = `${value} / ${showRelative ? '100' : this.max}${showRelative ? '%' : ''}`
      }

      return value
    },

    sortedVariants () {
      return [...this.options.display.thresholds].filter(t => t.value >= 0).sort((a, b) => b.value - a.value)
    },

    progressVariant () {
      const { variant } = this.options.display || {}
      let progressVariant = variant

      if (this.options.display.thresholds.length) {
        const { variant } = this.sortedVariants.find(t => this.progress >= t.value) || {}
        progressVariant = variant || progressVariant
      }

      return progressVariant
    },
  },

  watch: {
    'record.recordID': {
      immediate: true,
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
        maxValue: {
          filter: evaluatePrefilter(this.options.maxValue.filter, {
            record: this.record,
            recordID: (this.record || {}).recordID || NoID,
            ownerID: (this.record || {}).ownedBy || NoID,
            userID: (this.$auth.user || {}).userID || NoID,
          }),
        },
      }

      return this.block.fetch(additionalOptions, this.$ComposeAPI, namespaceID).then(({ value, max }) => {
        this.value = value
        this.max = max
      }).catch(this.toastErrorHandler(this.$t('progress.fetch-failed')))
        .finally(() => {
          this.processing = false
        })
    },
  },
}
</script>
