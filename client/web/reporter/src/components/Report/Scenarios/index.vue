<template>
  <div
    v-if="scenario"
  >
    <b-form-group
      :label="$t('scenarios:label')"
      label-class="text-primary"
    >
      <b-form-input
        v-model="scenario.label"
        :placeholder="$t('scenarios:scenario-name')"
      />
    </b-form-group>

    <b-form-group
      :label="$t('scenarios:datasource')"
      label-class="text-primary"
    >
      <b-form-select
        v-model="currentDatasourceName"
        :options="datasourceOptions"
      />
    </b-form-group>

    <b-form-group
      v-if="currentDatasourceName && scenario.filters[currentDatasourceName]"
      :label="$t('scenarios:prefilter')"
      label-class="text-primary"
    >
      <prefilter
        :filter="scenario.filters[currentDatasourceName]"
        :columns="columns"
      />
    </b-form-group>
  </div>
</template>

<script>
import Prefilter from 'corteza-webapp-reporter/src/components/Common/Prefilter'

export default {
  components: {
    Prefilter,
  },

  props: {
    currentIndex: {
      type: Number,
      default: () => -1,
    },

    scenario: {
      type: Object,
      required: true,
      default: () => ({}),
    },

    datasources: {
      type: Array,
      required: false,
      default: () => [],
    },
  },

  data () {
    return {
      currentDatasourceName: '',
      columns: [],
    }
  },

  computed: {
    datasourceOptions () {
      const options = [{ value: '', text: this.$t('general:label.none') }]

      this.datasources.forEach(({ step }, index) => {
        Object.entries(step).forEach(([kind, { name }]) => {
          if (['load'].includes(kind)) {
            options.push({ value: name || `${index}`, text: name || `${index}` })
          }
        })
      })

      return options
    },
  },

  watch: {
    currentIndex: {
      immediate: true,
      handler () {
        // Select first defined filter on switch
        const { filters = {} } = this.scenario
        const definedFilters = Object.keys(filters)

        this.currentDatasourceName = definedFilters.length ? definedFilters[0] : ''
      },
    },

    currentDatasourceName: {
      immediate: true,
      handler (name) {
        if (name && !this.scenario.filters[name]) {
          this.scenario.filters[name] = {}
        }

        this.getSourceColumns()
      },
    },
  },

  methods: {
    async getSourceColumns () {
      this.columns = []

      if (this.currentDatasourceName) {
        const steps = this.datasources.map(({ step }) => step)

        this.$SystemAPI.reportDescribe({ steps, describe: [this.currentDatasourceName] })
          .then((frames = []) => {
            this.columns = ((frames[0] || {}).columns || [])
          }).catch((e) => {
            this.toastErrorHandler(this.$t('notification:datasource.describe-failed'))(e)
          })
      }
    },
  },
}
</script>

<style>

</style>
