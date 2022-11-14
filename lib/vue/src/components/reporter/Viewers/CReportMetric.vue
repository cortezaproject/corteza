<template>
  <div
    v-if="options"
    :style="style"
    class="d-flex align-items-center justify-content-center overflow-hidden h-100 px-2 rounded"
  >
    <svg
      :viewBox="viewbox"
      class="h-100 w-100 align-items-end d-flex"
      width="100%"
      height="100%"
    >
      <text
        ref="metricItem"
        y="50%"
        x="50%"
        text-anchor="middle"
        dominant-baseline="central"
        text-rendering="geometricPrecision"
      >
        {{ displayedMetric }}
      </text>
    </svg>
  </div>
</template>

<script>
import base from './base.vue'

export default {
  extends: base,

  data () {
    return {
      vvb: ['0', '0', '0', '0'],
    }
  },

  computed: {
    viewbox () {
      return this.vvb.join(' ')
    },

    style () {
      return {
        fill: this.options.color || '#000000',
        backgroundColor: this.options.backgroundColor || '#ffffff',
      }
    },

    value () {
      if (this.dataframes.length) {
        const { rows = [], columns = [] } = this.dataframes[0]

        if (columns.length) {
          const columnIndex = columns.findIndex(({ name }) => name === this.options.valueColumn)

          if (rows.length) {
            return rows[0] ? rows[0][columnIndex] || '' : ''
          }
        }
      }
      return ''
    },

    displayedMetric () {
      const { prefix = '', suffix = '' } = this.options

      if (this.value) {
        return `${prefix}${this.value}${suffix}`
      }

      return ''
    },
  },

  watch: {
    value: {
      immediate: true,
      handler: function () {
        this.update()
      },
    },
  },

  methods: {
    update () {
      this.$nextTick(function () {
        const { width, height } = this.$refs.metricItem.getBBox()
        const tmp = [...this.vvb]
        tmp[2] = parseInt(Math.ceil(width))
        tmp[3] = parseInt(Math.ceil(height))
        this.vvb = tmp
      })
    },
  },
}
</script>
