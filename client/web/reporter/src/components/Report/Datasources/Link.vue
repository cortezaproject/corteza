<template>
  <div
    v-if="step.link"
  >
    <b-row>
      <b-col>
        <b-form-group
          :label="$t('datasources:name')"
          label-class="text-primary"
        >
          <b-form-input
            v-model="step.link.name"
            :placeholder="$t('datasources:datasource-name')"
          />
        </b-form-group>
      </b-col>
    </b-row>

    <hr>

    <b-row>
      <b-col cols="6">
        <b-form-group
          :label="$t('datasources:primary.source')"
          label-class="text-primary"
        >
          <b-form-select
            v-model="step.link.localSource"
            :options="supportedSources"
            @change="onSourceChange('local')"
          >
            <template #first>
              <b-form-select-option
                value=""
              >
                {{ $t('general:label.none') }}
              </b-form-select-option>
            </template>
          </b-form-select>
        </b-form-group>
      </b-col>
      <b-col cols="6">
        <b-form-group
          :label="$t('datasources:secondary.source')"
          label-class="text-primary"
        >
          <b-form-select
            v-model="step.link.foreignSource"
            :options="supportedSources"
            @change="onSourceChange('foreign')"
          >
            <template #first>
              <b-form-select-option
                value=""
              >
                {{ $t('general:label.none') }}
              </b-form-select-option>
            </template>
          </b-form-select>
        </b-form-group>
      </b-col>
    </b-row>

    <b-row>
      <b-col cols="6">
        <b-form-group
          v-if="step.link.localSource"
          :label="$t('datasources:primary.column')"
          label-class="text-primary"
        >
          <b-form-select
            v-model="step.link.localColumn"
            :options="localColumns"
            value-field="name"
            text-field="label"
          >
            <template #first>
              <b-form-select-option
                value=""
              >
                {{ $t('general:label.none') }}
              </b-form-select-option>
            </template>
          </b-form-select>
        </b-form-group>
      </b-col>
      <b-col cols="6">
        <b-form-group
          v-if="step.link.foreignSource"
          :label="$t('datasources:secondary.column')"
          label-class="text-primary"
        >
          <b-form-select
            v-model="step.link.foreignColumn"
            :options="foreignColumns"
            value-field="name"
            text-field="label"
          >
            <template #first>
              <b-form-select-option
                value=""
              >
                {{ $t('general:label.none') }}
              </b-form-select-option>
            </template>
          </b-form-select>
        </b-form-group>
      </b-col>
    </b-row>
  </div>
</template>

<script>
import base from './base.vue'

export default {
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
      localColumns: [],
      foreignColumns: [],
    }
  },

  computed: {
    supportedSources () {
      const options = []

      this.datasources.forEach(({ step }, index) => {
        Object.entries(step).forEach(([kind, { name }]) => {
          if (['load', 'group'].includes(kind)) {
            options.push({ value: name || `${index}`, text: name || `${index}` })
          }
        })
      })

      return options
    },
  },

  watch: {
    'step.link.name': {
      immediate: true,
      handler (newStep, oldStep) {
        if (!oldStep && newStep) {
          this.getSourceColumns(['local', 'foreign'])
        }
      },
    },

    'step.link.localSource': {
      handler () {
        this.getSourceColumns(['local'])
      },
    },

    'step.link.foreignSource': {
      handler () {
        this.getSourceColumns(['foreign'])
      },
    },
  },

  methods: {
    onSourceChange (source) {
      this.step.link[`${source}Column`] = ''
    },

    async getSourceColumns (sources = []) {
      sources.forEach(source => {
        this[`${source}Columns`] = []

        const sourceType = this.step.link[`${source}Source`]

        if (sourceType) {
          const steps = this.datasources.filter(({ step }) => step.load || step.aggregate).map(({ step }) => step)
          const describe = [sourceType]

          if (steps.length && describe.length) {
            this.$SystemAPI.reportDescribe({ steps, describe })
              .then((frames = []) => {
                const { columns = [] } = frames.find(({ source }) => describe.includes(source)) || {}
                this[`${source}Columns`] = columns
              }).catch((e) => {
                this.toastErrorHandler(this.$t('notification:datasource.describe-failed'))(e)
              })
          }
        }
      })
    },
  },
}
</script>
