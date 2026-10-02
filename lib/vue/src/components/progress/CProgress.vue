<template>
  <b-progress
    :max="maxValue"
    class="bg-light position-relative"
  >
    <b-progress-bar
      :value="progressValue < 0 ? 0 : progressValue"
      :striped="striped"
      :animated="animated"
      :variant="progressVariant"
    >
      <strong
        :class="textVariant"
        class="d-flex align-items-center justify-content-center position-absolute mb-0 w-100"
        :style="textStyle"
      >
        {{ progressLabel }}
      </strong>
    </b-progress-bar>
  </b-progress>
</template>

<script>
export default {
  props: {
    value: {
      type: Number,
      required: true,
    },

    min: {
      type: Number,
      default: 0,
    },

    max: {
      type: Number,
      default: 100,
    },

    labeled: {
      type: Boolean,
      default: true,
    },

    relative: {
      type: Boolean,
      default: true,
    },

    progress: {
      type: Boolean,
    },

    striped: {
      type: Boolean,
    },

    animated: {
      type: Boolean,
    },

    variant: {
      type: String,
      default: 'success',
    },

    thresholds: {
      type: Array,
      default: () => [],
    },

    textStyle: {
      type: String,
      default: '',
    },
  },

  computed: {
    // Values come from reports that may yield nothing; never render NaN
    safeMin () {
      return Number.isFinite(this.min) ? this.min : 0
    },

    safeMax () {
      return Number.isFinite(this.max) ? this.max : 100
    },

    safeValue () {
      return Number.isFinite(this.value) ? this.value : this.safeMin
    },

    maxValue () {
      return Math.abs(this.safeMax - this.safeMin)
    },

    progressValue () {
      if (this.safeValue < this.safeMin && this.safeMax > this.safeMin) {
        return this.safeValue - this.safeMin
      } else if (this.safeValue > this.safeMin && this.safeMax < this.safeMin) {
        return this.safeMin - this.safeValue
      }

      return Math.abs(this.safeValue - this.safeMin)
    },

    // Share of the range as a rounded percentage; 0 when the range is empty
    percentage () {
      if (!this.maxValue) {
        return 0
      }

      // https://stackoverflow.com/a/21907972/17926309
      return Math.round((((this.progressValue < 0 ? 0 : this.progressValue) / this.maxValue) * 100) * 100) / 100
    },

    progressLabel () {
      let value = this.safeValue

      if (!this.labeled) {
        return
      }

      if (this.relative) {
        value = `${this.percentage}%`
      }

      if (this.progress) {
        value = `${value} / ${this.relative ? '100' : this.safeMax}${this.relative ? '%' : ''}`
      }

      return value
    },

    progressVariant () {
      const value = this.percentage

      let progressVariant = this.variant

      if (this.thresholds.length) {
        const { variant } = this.sortedVariants.find(t => value >= t.value) || {}
        progressVariant = variant || progressVariant
      }

      return progressVariant
    },

    sortedVariants () {
      return [...this.thresholds].filter(t => t.value >= 0).sort((a, b) => b.value - a.value)
    },

    textVariant () {
      return ['dark', 'primary'].includes(this.progressVariant) ? 'text-white' : 'text-dark'
    },
  },
}
</script>
