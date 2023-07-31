<template>
  <b-tab :title="$t('report.label')">
    <b-row>
      <b-col>
        <b-form-group
          :label="$t('report.label')"
          label-class="text-primary"
        >
          <c-input-select
            v-model="options.reportID"
            :options="reports"
            :get-option-label="o => o.meta.name || o.handle"
            default-value="0"
            :reduce="o => o.reportID"
            @input="handleReportChange"
          />
        </b-form-group>
      </b-col>

      <b-col
        v-if="selectedReport && selectedReport.scenarios && selectedReport.scenarios.length > 1"
      >
        <b-form-group
          :label="$t('report.scenario.label')"
          label-class="text-primary"
        >
          <c-input-select
            v-model="options.scenarioID"
            :options="selectedReport.scenarios"
            default-value="0"
            :reduce="o => o.scenarioID"
          />
        </b-form-group>
      </b-col>
    </b-row>

    <b-form-group
      v-if="selectedReport"
      :label="$t('report.element.label')"
      :description="$t('report.element.description')"
      label-class="text-primary"
    >
      <c-input-select
        v-model="options.elementID"
        :options="selectedReport.blocks"
        :reduce="o => o.elements[0].elementID"
        :get-option-label="getElementsOptionLabel"
        default-value="0"
      >
        <template #option="option">
          <div v-if="option.elements.length > 0">
            <strong>{{ option.title || `${$t('general:label.block')} ${option.key}` }}</strong>

            <ul class="list-unstyled">
              <li
                v-for="subOption in option.elements"
                :key="subOption.name"
                class="ml-2"
              >
                {{ subOption.name || subOption.kind }}
              </li>
            </ul>
          </div>

          <div v-else>
            {{ option.name }}
          </div>
        </template>
      </c-input-select>
    </b-form-group>
  </b-tab>
</template>
<script>
import base from '../base'
import { NoID } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  extends: base,

  data () {
    return {
      reports: [],
    }
  },

  computed: {
    selectedReport () {
      const { reportID = NoID } = this.options

      if (reportID !== NoID) {
        return this.reports.find(r => r.reportID === reportID)
      }

      return undefined
    },
  },

  watch: {
    'options.blockID' () {
      this.options.elementID = NoID
    },
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  created () {
    this.fetchReports()
  },

  methods: {
    fetchReports () {
      this.$SystemAPI.reportList()
        .then(({ set = [] }) => {
          this.reports = set
        })
        .catch(this.toastErrorHandler(this.$t('notification:report.listFetchFailed')))
    },

    setDefaultValues () {
      this.reports = []
    },

    getElementsOptionLabel (o) {
      const blockTitle = o.title.length > 1 ? `${o.title} - ` : ''
      return `${blockTitle}${o.elements[0].name}` || o.elements[0].kind
    },

    handleReportChange () {
      if (this.options.elementID) {
        this.options.elementID = NoID
        this.options.scenarioID = NoID
      }
    },
  },
}
</script>
