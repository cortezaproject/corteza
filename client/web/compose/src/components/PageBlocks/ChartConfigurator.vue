<template>
  <b-tab :title="$t('chart.label')">
    <b-form-group
      :label="$t('chart.display')"
    >
      <b-form-select
        v-model="block.options.chartID"
        :options="options"
        text-field="name"
        value-field="chartID"
      />
    </b-form-group>
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

    options () {
      return [
        { chartID: NoID, name: this.$t('chart.pick') },
        ...this.charts,
      ]
    },
  },
}
</script>
