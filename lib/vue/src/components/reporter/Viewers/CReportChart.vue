<template>
  <div class="d-flex h-100 position-relative">
    <c-chart
      v-if="chart"
      :chart="chart"
      class="flex-fill p-1"
    />
  </div>
</template>
<script>
import base from './base.vue'
import { CChart } from '../../chart/index.ts'

export default {
  extends: base,

  title: 'CReportChart',

  components: {
    CChart,
  },

  data () {
    return {
      chart: undefined,
    }
  },

  watch: {
    dataframes: {
      deep: true,
      immediate: true,
      handler () {
        this.$nextTick(() => {
          this.renderChart()
        })
      },
    },

    options: {
      deep: true,
      handler () {
        this.chart = undefined

        this.$nextTick(() => {
          this.renderChart()
        })
      },
    }
  },

  methods: {
    renderChart () {
      const meta = {
        themeVariables: this.getThemeVariables(),
      }

      this.chart = this.options.getChartConfiguration(this.dataframes, meta)
    },

    getThemeVariables () {
      const getCssVariable = (variableName) => {
        return getComputedStyle(document.documentElement).getPropertyValue(variableName).trim()
      }

      // Turn below into an object with key value pairs
      return ['white', 'black', 'primary', 'secondary', 'success', 'warning', 'danger', 'light', 'extra-light', 'dark', 'font-regular'].reduce((acc, variable) => {
        acc[variable] = getCssVariable(`--${variable}`)
        return acc
      }, {})
    },
  },
}
</script>
