<template>
  <div
    v-if="step.group"
  >
    <b-form-group
      :label="$t('datasources:name')"
      label-class="text-primary"
    >
      <b-form-input
        v-model="step.group.name"
        :placeholder="$t('datasources:datasource-name')"
      />
    </b-form-group>

    <hr>

    <b-form-group
      :label="$t('datasources:source')"
      label-class="text-primary"
    >
      <b-form-select
        v-model="step.group.source"
        :options="supportedSources"
        @change="reset"
      >
        <template #first>
          <b-form-select-option
            :value="undefined"
          >
            {{ $t('general:label.none') }}
          </b-form-select-option>
        </template>
      </b-form-select>
    </b-form-group>

    <div
      v-if="step.group.source"
    >
      <b-form-group
        :label="$t('datasources:group-by')"
        label-class="text-primary"
      >
        <group-by
          :group-by.sync="step.group.keys"
        />
      </b-form-group>

      <b-form-group
        :label="$t('datasources:aggregate')"
        label-class="text-primary"
      >
        <aggregate
          :aggregate.sync="step.group.columns"
        />
      </b-form-group>

      <!-- <b-form-group
        v-if="columns.length"
        label="Prefilter"
        label-class="text-primary"
      >
        <prefilter
          :filter.sync="step.group.filter"
          :columns="columns"
        />
      </b-form-group> -->

      <!-- <b-form-group
        v-if="columns.length"
        label="Presort order"
        label-class="text-primary"
      >
        <c-input-presort
          v-model="step.group.sort"
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
  </div>
</template>

<script>
import base from '../base.vue'
import GroupBy from './GroupBy'
import Aggregate from './Aggregate'
// import Prefilter from 'corteza-webapp-reporter/src/components/Common/Prefilter'
// import { components } from '@cortezaproject/corteza-vue'
// const { CInputPresort } = components

export default {
  components: {
    GroupBy,
    Aggregate,
    // Prefilter,
    // CInputPresort,
  },

  extends: base,

  props: {
    datasources: {
      type: Array,
      required: false,
      default: () => [],
    },
  },

  data () {
    return {
      columns: [],
    }
  },

  computed: {
    supportedSources () {
      const options = []

      this.datasources.forEach(({ step }, index) => {
        Object.entries(step).forEach(([kind, { name }]) => {
          if (kind === 'load') {
            options.push({ value: name || `${index}`, text: name || `${index}` })
          }
        })
      })

      return options
    },
  },

  methods: {
    async getSourceColumns () {
      const steps = this.datasources.filter(({ step }) => step.load).map(({ step }) => step)
      steps.push(this.step)
      const describe = [this.step.group.name]

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
      this.step.group.filter = {}
      this.step.group.sort = ''
      this.step.group.keys = []
      this.step.group.columns = []
    },
  },
}
</script>
