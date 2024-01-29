<template>
  <div
    v-if="displayElement"
  >
    <b-form-group
      :label="$t('general:label.name')"
      label-class="text-primary"
    >
      <b-form-input
        v-model="displayElement.name"
      />
    </b-form-group>

    <b-button
      v-b-toggle.datasources
      block
      :disabled="!usesDatasources"
      variant="primary"
      class="mb-2"
    >
      {{ $t('builder:datasources.label') }}
    </b-button>

    <b-collapse
      v-if="usesDatasources"
      id="datasources"
      :visible="usesDatasources"
      accordion
    >
      <b-form-group
        :label="$t('builder:datasources.label')"
        label-class="text-primary"
      >
        <b-form-select
          v-model="options.source"
          :options="sources"
          @change="setConfigurableSources"
        />
      </b-form-group>

      <div
        v-if="currentConfigurableDatasourceName"
      >
        <b-form-group
          v-if="hasMultipleConfigurableDatasources"
          :label="$t('builder:joined-datasource-handling')"
          label-class="text-primary"
        >
          <b-form-select
            v-model="currentConfigurableDatasourceName"
            :options="options.datasources"
            text-field="name"
            value-field="name"
            @change="configurableDatasourceChanged"
          />
        </b-form-group>

        <div
          v-if="currentConfigurableDatasourceIndex >= 0"
        >
          <b-form-group
            v-if="columns.length"
            :label="$t('builder:prefilter')"
            label-class="text-primary"
          >
            <prefilter
              :filter.sync="options.datasources[currentConfigurableDatasourceIndex].filter"
              :columns="columns[currentConfigurableDatasourceIndex]"
            />
          </b-form-group>

          <b-form-group
            v-if="columns.length"
            :label="$t('builder:presort-order')"
            label-class="text-primary"
          >
            <c-input-presort
              v-model="options.datasources[currentConfigurableDatasourceIndex].sort"
              :fields="columns[currentConfigurableDatasourceIndex]"
              :labels="{
                add: $t('general:label.add'),
                ascending: $t('general:label.ascending'),
                descending: $t('general:label.descending'),
                none: $t('general:label.none'),
                addButton: $t('general:label.add'),
              }"
            />
          </b-form-group>

          <b-form-group
            v-if="displayElement.kind === 'Table'"
            :label="$t('builder:limit')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="pagingLimit"
              type="number"
            />
          </b-form-group>
        </div>
      </div>
    </b-collapse>

    <b-button
      v-b-toggle.display
      block
      :disabled="!showDisplayElementConfigurator"
      variant="primary"
      class="mb-2"
    >
      {{ $t('builder:element') }}
    </b-button>

    <b-collapse
      id="display"
      :visible="showDisplayElementConfigurator"
      accordion
    >
      <component
        :is="displayElementConfigurator"
        :display-element-options.sync="options"
        :columns="columns"
        :datasource="currentDatasource"
      />
    </b-collapse>
  </div>
</template>

<script>
import getDisplayElementConfigurator from './loader'
import Prefilter from 'corteza-webapp-reporter/src/components/Common/Prefilter'
import { components } from '@cortezaproject/corteza-vue'
const { CInputPresort } = components

export default {
  components: {
    Prefilter,
    CInputPresort,
  },

  props: {
    displayElement: {
      type: Object,
      required: true,
    },

    block: {
      type: Object,
      required: true,
    },

    datasources: {
      type: Array,
      required: true,
    },
  },

  data () {
    return {
      columns: [],

      currentConfigurableDatasourceName: undefined,
      currentConfigurableDatasourceIndex: undefined,
    }
  },

  computed: {
    usesDatasources () {
      return !['Text'].includes(this.displayElement.kind)
    },

    displayElementConfigurator () {
      return getDisplayElementConfigurator(this.displayElement.kind)
    },

    showDisplayElementConfigurator () {
      return this.usesDatasources ? !!this.currentDatasource : true
    },

    sources () {
      const options = [{ value: '', text: this.$t('general:label.none') }]

      this.datasources.forEach(({ step }, index) => {
        Object.values(step).forEach(({ name }) => {
          options.push({ value: name || `${index}`, text: name || `${index}` })
        })
      })

      return options
    },

    currentDatasource () {
      if (this.options.source) {
        return this.datasources.find(({ step: { load = {}, link = {}, join = {}, aggregate = {} } }) => [load.name, link.name, join.name, aggregate.name].includes(this.options.source))
      }

      return undefined
    },

    hasMultipleConfigurableDatasources () {
      return this.currentDatasource && this.currentDatasource.step.link && this.options.datasources.length > 1
    },

    options: {
      get () {
        return this.displayElement ? this.displayElement.options : undefined
      },

      set (options = {}) {
        if (this.displayElement) {
          this.$emit('update:displayElement', { ...this.displayElement, options })
        }
      },
    },

    pagingLimit: {
      get () {
        let limit = 0

        if (this.currentConfigurableDatasourceIndex >= 0) {
          const { paging = {} } = this.options.datasources[this.currentConfigurableDatasourceIndex] || {}
          limit = paging.limit || 0
        }

        return limit
      },

      set (limit = 0) {
        if (this.currentConfigurableDatasourceIndex >= 0) {
          if (!this.options.datasources[this.currentConfigurableDatasourceIndex].paging) {
            this.options.datasources[this.currentConfigurableDatasourceIndex].paging = {}
          }

          this.options.datasources[this.currentConfigurableDatasourceIndex].paging.limit = limit || 0
        }
      },
    },
  },

  watch: {
    'options.source': {
      immediate: true,
      handler (source) {
        this.describeReport(source)
      },
    },

    'displayElement.elementID': {
      immediate: true,
      handler () {
        this.currentConfigurableDatasourceIndex = this.datasources.length ? 0 : -1
        if (this.usesDatasources) {
          this.currentConfigurableDatasourceName = (this.options.datasources[0] || {}).name
        }
      },
    },
  },

  methods: {
    describeReport (source) {
      this.columns = []

      if (source) {
        const steps = this.datasources.map(({ step }) => step)
        this.$SystemAPI.reportDescribe({ steps, describe: [this.options.source] })
          .then((frames = []) => {
            this.columns = frames.filter(({ source }) => source === this.options.source).map(({ columns = [] }) => columns) || []
          }).catch((e) => {
            this.toastErrorHandler(this.$t('notification:datasource.describe-failed'))(e)
          })
      }
    },

    setConfigurableSources (source) {
      this.options.datasources = []
      this.currentConfigurableDatasourceName = undefined

      let configurableDatasources = []
      if (source) {
        configurableDatasources = [source]

        const { link } = this.currentDatasource.step
        if (link) {
          configurableDatasources = [link.localSource, link.foreignSource]
        }
      }

      this.currentConfigurableDatasourceName = configurableDatasources[0]
      this.options.datasources = configurableDatasources.map(s => {
        return {
          name: s,
          sort: '',
          filter: {},
        }
      })
    },

    configurableDatasourceChanged (source) {
      if (source) {
        this.currentConfigurableDatasourceIndex = this.options.datasources.findIndex(({ name }) => source === name)
      }
    },
  },
}
</script>

<style>

</style>
