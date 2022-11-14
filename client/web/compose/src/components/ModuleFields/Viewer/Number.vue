<template>
  <div>
    <!-- Extra empty line is added thanks to white-space: pre-line (multivalue) if we write div in multiple lines  -->
    <!-- eslint-disable-next-line -->
    <div v-if="field.options.display === 'number'" :class="classes">{{ formatted }}</div>

    <div v-else>
      <b-progress
        v-for="(v, i) in formatted"
        :key="i"
        :max="field.options.max"
        height="1.5rem"
        class="bg-light"
        :class="{ 'mt-2': i }"
      >
        <b-progress-bar
          :value="v"
          :striped="field.options.striped"
          :animated="field.options.animated"
          :variant="getProgressVariant(v)"
        >
          {{ getProgressLabel(v) }}
        </b-progress-bar>
      </b-progress>
    </div>

    <errors :errors="errors" />
  </div>
</template>

<script>
import base from './base'

export default {
  extends: base,

  computed: {
    formatted () {
      if (!this.value) {
        return
      }

      const value = this.field.isMulti ? this.value : [this.value]

      if (this.field.options.display === 'number') {
        return value.map(v => this.field.formatValue(v)).join(this.field.options.multiDelimiter)
      } else {
        return value
      }
    },

    sortedVariants () {
      return [...this.field.options.thresholds].filter(t => t.value >= 0).sort((a, b) => b.value - a.value)
    },
  },

  methods: {
    getProgressLabel (value) {
      const { max, showValue, showRelative, showProgress } = this.field.options

      if (!showValue) {
        return
      }

      if (showRelative) {
        // https://stackoverflow.com/a/21907972/17926309
        value = `${Math.round(((value / max) * 100) * 100) / 100}%`
      }

      if (showProgress) {
        value = `${value} / ${showRelative ? '100' : max}${showRelative ? '%' : ''}`
      }

      return value
    },

    getProgressVariant (value) {
      let progressVariant = this.field.options.variant

      if (this.field.options.thresholds.length) {
        const { variant } = this.sortedVariants.find(t => value >= t.value) || {}
        progressVariant = variant || progressVariant
      }

      return progressVariant
    },
  },
}
</script>
