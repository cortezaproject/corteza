<template>
  <div
    v-if="step.load"
  >
    <b-form-group
      :label="$t('datasources:name')"
      label-class="text-primary"
    >
      <b-form-input
        v-model="step.load.name"
        :placeholder="$t('datasources:datasource-name')"
      />
    </b-form-group>

    <hr>

    <b-form-group
      :label="$t('datasources:source')"
      label-class="text-primary"
    >
      <b-form-select
        v-model="step.load.source"
        :options="supportedSources"
        @change="reset"
      />
    </b-form-group>

    <component
      :is="sourceTypeComponent(step.load.source)"
      v-if="step.load.source"
      :definition.sync="stepDefinition"
    />

    <b-form-group
      v-if="columns.length"
      :label="$t('datasources:prefilter')"
      label-class="text-primary"
    >
      <prefilter
        :filter.sync="step.load.filter"
        :columns="columns"
      />
    </b-form-group>

    <!-- <b-form-group
      v-if="columns.length"
      label="Presort order"
      label-class="text-primary"
    >
      <c-input-presort
        v-model="step.load.sort"
        :fields="columns"
        :labels="{
          add: $t('general:label.add'),
          ascending: $t('general:label.ascending'),
          descending: $t('general:label.descending'),
          none: $t('general:label.none'),
        }"
      />
    </b-form-group> -->
  </div>
</template>

<script>
import base from '../base.vue'
import loader from './loader'
import Prefilter from 'corteza-webapp-reporter/src/components/Common/Prefilter'
// import { components } from '@cortezaproject/corteza-vue'
// const { CInputPresort } = components
export default {
  components: {
    Prefilter,
    // CInputPresort,
  },

  extends: base,

  data () {
    return {
      // @todo get this from the API
      supportedSources: [
        {
          text: this.$t('datasources:compose-records'),
          value: 'composeRecords',
          definition: [{
            label: 'namespace',
            key: 'namespace',
          }, {
            label: 'module',
            key: 'module',
          }],
        },
      ],

      columns: [],
    }
  },

  computed: {
    stepDefinition: {
      get () {
        return this.step.load ? this.step.load.definition : undefined
      },

      set (definition) {
        this.step.load.definition = definition
        this.$emit('update:step', this.step)
      },
    },
  },

  watch: {
    stepDefinition: {
      immediate: true,
      deep: true,
      handler ({ moduleID, namespaceID }) {
        if (moduleID && namespaceID) {
          this.getSourceColumns()
        }
      },
    },
  },

  methods: {
    sourceTypeComponent: loader,

    async getSourceColumns () {
      const steps = [this.step]
      const describe = [this.step.load.name]

      if (steps.length && describe.length) {
        this.$SystemAPI.reportDescribe({ steps, describe })
          .then((frames = []) => {
            const { columns = [] } = frames.find(({ source }) => describe.includes(source)) || {}
            this.columns = columns
          }).catch((e) => {
            this.toastErrorHandler(this.$t('notification:datasource.describe-failed'))(e)
          })
      }
    },

    reset () {
      this.step.load.filter = {}
      this.step.load.sort = ''
    },
  },
}
</script>
