<template>
  <div class="position-relative h-100 w-100 p-2">
    <canvas
      ref="chart"
    />
  </div>
</template>
<script>
import base from './base.vue'
import colorschemes from 'chartjs-plugin-colorschemes'
import Funnel from 'chartjs-plugin-funnel'

export default {
  extends: base,

  data () {
    return {
      chart: undefined,

      plugins: [
        Funnel,
        colorschemes,
      ],
    }
  },

  computed: {
    localDataframe () {
      return this.dataframes[0]
    },

    size () {
      return this.displayElement ? this.displayElement.meta.size || 100 : 100
    },
  },

  watch: {
    dataframes: {
      deep: true,
      handler (dataframes = []) {
        if (dataframes.length) {
          this.$nextTick(() => {
            this.renderChart()
          })
        }
      },
    },

    size: {
      immediate: true,
      deep: true,
      handler () {
        this.$nextTick(() => {
          this.renderChart()
        })
      },
    },
  },

  methods: {
    renderChart () {
      if (this.chart) {
        this.chart.destroy()
      }

      const ctx = this.$refs.chart.getContext('2d')

      const chartConfig = this.options.getChartConfiguration(this.dataframes)

      this.chart = new Chart(ctx, { ...chartConfig, plugins: this.plugins })
    },
  },
}
</script>
