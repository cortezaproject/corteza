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
  },

  computed: {
    maxValue () {
      return Math.abs(this.max - this.min)
    },

    progressValue () {
      if (this.value < this.min && this.max > this.min) {
        return this.value - this.min
      } else if (this.value > this.min && this.max < this.min) {
        return this.min - this.value
      }

      return Math.abs(this.value - this.min)
    },

    progressLabel () {
      let value = this.value

      if (!this.labeled) {
        return
      }

      if (this.relative) {
        // https://stackoverflow.com/a/21907972/17926309
        value = `${Math.round((((this.progressValue) / this.maxValue) * 100) * 100) / 100}%`
      }

      if (this.progress) {
        value = `${value} / ${this.relative ? '100' : this.max}${this.relative ? '%' : ''}`
      }

      return value
    },

    progressVariant () {
      const value = Math.round((((this.progressValue < 0 ? 0 : this.progressValue) / this.maxValue) * 100) * 100) / 100

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
