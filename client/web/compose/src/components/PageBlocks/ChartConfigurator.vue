<template>
  <b-tab :title="$t('chart.label')">
    <b-form-group
      :label="$t('chart.display')"
    >
      <b-form-select
        v-model="block.options.chartID"
        :options="chartOptions"
        text-field="name"
        value-field="chartID"
      />
    </b-form-group>

    <b-button
      v-if="selectedChart"
      :disabled="!selectedChart.canUpdateChart && !selectedChart.canDeleteChart"
      variant="light"
      :to="{ name: 'admin.charts.edit', params: { chartID: selectedChart.chartID }, query: null }"
    >
      {{ $t('chart.openInBuilder') }}
    </b-button>
  </b-tab>
</template>
<script>
import { mapGetters } from 'vuex'
import { NoID } from '@cortezaproject/corteza-js'
import base from './base'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'Chart',

  extends: base,

  computed: {
    ...mapGetters({
      charts: 'chart/set',
    }),

    chartOptions () {
      return [
        { chartID: NoID, name: this.$t('chart.pick') },
        ...this.charts,
      ]
    },

    selectedChart () {
      if (!this.options.chartID || this.options.chartID === NoID) {
        return
      }

      return this.chartOptions.find(({ chartID }) => chartID === this.options.chartID)
    },
  },
}
</script>
