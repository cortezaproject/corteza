<template>
  <div
    :style="genStyle(metric.valueStyle)"
    class="h-100 text-center"
  >
    <!--    This div is here because .svg metrics dont render with print to PDF option-->
    <div
      class="d-none d-print-flex h-100 w-100 align-items-center justify-content-center overflow-hidden"
      :style="genStyle(metric.valueStyle)"
    >
      <template v-if="metric.prefix">
        {{ metric.prefix }}
      </template>
      {{ value.value }}
      <template v-if="metric.suffix">
        {{ metric.suffix }}
      </template>
    </div>

    <svg
      :viewBox="getVB"
      class="h-100 w-100 d-flex d-print-none"
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
        <template v-if="metric.prefix">
          {{ metric.prefix }}
        </template>
        {{ value.value }}
        <template v-if="metric.suffix">
          {{ metric.suffix }}
        </template>
      </text>
    </svg>
  </div>
</template>

<script>
export default {
  props: {
    metric: {
      type: Object,
      required: false,
      default: () => ({}),
    },
    value: {
      type: Object,
      required: false,
      default: () => ({}),
    },
  },

  data () {
    return {
      vvb: ['0', '0', '0', '0'],
    }
  },

  computed: {
    getVB () {
      return this.vvb.join(' ')
    },
  },

  watch: {
    metric: {
      handler () {
        this.update()
      },
      immediate: true,
    },
    value: {
      handler () {
        this.update()
      },
      immediate: true,
    },
  },

  methods: {
    update () {
      this.$nextTick(() => {
        const { width, height } = this.$refs.metricItem.getBBox()
        const tmp = [...this.vvb]
        tmp[2] = parseInt(Math.ceil(width))
        tmp[3] = parseInt(Math.ceil(height))
        this.vvb = tmp
      })
    },

    genStyle (s = {}) {
      const d = {
        fill: s.color,
        backgroundColor: s.backgroundColor,
        fontSize: s.fontSize ? s.fontSize + 'px' : undefined,
        color: s.color,
      }

      for (const v of Object.keys(d)) {
        if (d[v] === undefined) {
          delete d[v]
        }
      }

      return d
    },
  },
}
</script>
