<template>
  <b-card
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <b-form-group
      :label="$t('recordList.export.selectFields')"
      :description="$t('recordList.export.limitations')"
      label-class="text-primary"
    >
      <field-picker
        v-if="module"
        :module="module"
        :system-fields="systemFields"
        :disabled-types="disabledTypes"
        :fields.sync="selectedFields"
        style="max-height: 45vh;"
      />
    </b-form-group>

    <b-form-group>
      <b-form-checkbox
        v-model="forTimezone"
        class="mb-1"
      >
        {{ $t('recordList.export.specifyTimezone') }}
      </b-form-checkbox>

      <vue-select
        v-if="forTimezone"
        v-model="exportTimezone"
        :options="timezones"
        class="bg-white"
        :placeholder="$t('recordList.export.timezonePlaceholder')"
      />
    </b-form-group>

    <b-form-group>
      <b-form-radio-group
        v-model="rangeType"
        :options="rangeTypeOptions"
        stacked
      />
    </b-form-group>

    <template v-if="rangeType === 'range'">
      <b-row
        v-if="rangeType === 'range'"
      >
        <b-col
          md="6"
        >
          <b-form-group
            :label="$t('recordList.export.rangeBy')"
            label-class="text-primary"
          >
            <b-form-select
              v-model="rangeBy"
              :options="rangeByOptions"
            />
          </b-form-group>
        </b-col>
        <b-col
          md="6"
        >
          <b-form-group
            :label="$t('recordList.export.dateRange')"
            label-class="text-primary"
          >
            <b-form-select
              v-model="range"
              :options="dateRangeOptions"
            />
          </b-form-group>
        </b-col>
      </b-row>
    </template>

    <b-row
      v-if="rangeType === 'range'"
      class="mb-3"
    >
      <b-col
        cols="6"
      >
        <b-form-input
          v-model="start"
          :state="dateRangeValid ? null : false"
          type="date"
          :max="end"
          @keydown.prevent
        />
      </b-col>
      <b-col
        cols="6"
      >
        <b-form-input
          v-model="end"
          :state="dateRangeValid ? null : false"
          type="date"
          :min="start"
          @keydown.prevent
        />
      </b-col>
    </b-row>

    <template
      v-if="rangeType !== 'selection'"
    >
      <b-form-group
        :label="$t('recordList.export.query.label')"
        label-class="text-primary"
      >
        <b-form-input
          v-model="exportQuery"
          :placeholder="$t('recordList.export.query.placeholder')"
          debounce="500"
        />
      </b-form-group>

      <b-form-group
        :label="$t('recordList.export.filter.label')"
        label-class="text-primary"
        class="mb-0"
      >
        <b-form-textarea
          v-model="exportFilter"
          :placeholder="$t('recordList.export.filter.placeholder')"
          debounce="500"
        />
      </b-form-group>
    </template>

    <div
      slot="footer"
      class="d-flex"
    >
      <span
        class="my-auto"
      >
        <b-spinner
          v-if="processingCount"
          small
        />
        <span
          v-else
        >
          {{ $t('recordList.export.recordCount', { count: getExportableCount || 0 }) }}
        </span>
      </span>
      <span class="ml-auto">
        <c-input-processing
          v-if="allowJSON"
          :processing="processing"
          :disabled="exportDisabled"
          variant="light"
          size="lg"
          @click="doExport('json')"
        >
          {{ $t('recordList.export.json') }}
        </c-input-processing>
        <c-input-processing
          v-if="allowCSV"
          :processing="processing"
          :disabled="exportDisabled"
          variant="light"
          size="lg"
          class="ml-2"
          @click="doExport('csv')"
        >
          {{ $t('recordList.export.csv') }}
        </c-input-processing>
      </span>
    </div>
  </b-card>
</template>

<script>
import { compose, NoID } from '@cortezaproject/corteza-js'
import FieldPicker from 'corteza-webapp-compose/src/components/Common/FieldPicker'
import moment from 'moment'
import tz from 'compact-timezone-list'
import { evaluatePrefilter, queryToFilter } from 'corteza-webapp-compose/src/lib/record-filter'
import { VueSelect } from 'vue-select'
const fmtDate = (d) => d.format('YYYY-MM-DD')

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    FieldPicker,
    VueSelect,
  },

  props: {
    allowJSON: {
      type: Boolean,
      default: true,
    },
    allowCSV: {
      type: Boolean,
      default: true,
    },
    module: {
      type: compose.Module,
      required: true,
    },
    preselectedFields: {
      type: Array,
      default: () => [],
    },
    recordCount: {
      type: Number,
      default: 0,
    },
    query: {
      type: String,
      required: false,
      default: undefined,
    },
    prefilter: {
      type: String,
      required: false,
      default: undefined,
    },
    selection: {
      type: Array,
      required: false,
      default: () => [],
    },
    filterRangeType: {
      type: String,
      default: 'all',
    },
    filterRangeBy: {
      type: String,
      default: 'createdAt',
    },
    dateRange: {
      type: String,
      default: 'lastMonth',
    },
    startDate: {
      type: String,
      default: null,
    },
    endDate: {
      type: String,
      default: null,
    },
    systemFields: {
      type: Array,
      default: () => ['ownedBy', 'createdAt', 'createdBy', 'updatedAt', 'updatedBy'],
    },
    disabledTypes: {
      type: Array,
      default: () => ['User', 'Record', 'File'],
    },
    processing: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      fields: [],
      forTimezone: false,
      exportTimezone: undefined,
      filter: {
        rangeType: null,
        query: this.query,
        filter: this.prefilter,
        rangeBy: null,
        date: {
          range: null,
          start: null,
          end: null,
        },
      },
      rangeRecordCount: 0,
      processingCount: false,
    }
  },

  computed: {
    timezones () {
      return tz.map(({ label, tzCode, offset }) => ({ label, tzCode, offset }))
    },

    // These should be computed, because of i18n
    rangeTypeOptions () {
      const options = [
        {
          value: 'all',
          text: this.$t('recordList.export.all'),
        },
        {
          value: 'range',
          text: this.$t('recordList.export.inRange'),
        },
      ]

      if (this.hasSelection) {
        options.push({
          value: 'selection',
          text: this.$t('recordList.export.selection'),
        })
      }

      return options
    },

    /**
     * checks if the given date-range is valid
     * @returns {Boolean}
     */
    dateRangeValid () {
      if (this.end < this.start) {
        return false
      }
      return true
    },

    rangeByOptions () {
      return [
        {
          value: 'createdAt',
          text: this.$t('recordList.export.filter.createdAt'),
        },
        {
          value: 'updatedAt',
          text: this.$t('recordList.export.filter.updatedAt'),
        },
      ]
    },

    dateRangeOptions () {
      return [
        {
          value: 'lastMonth',
          text: this.$t('recordList.export.filter.lastMonth'),
        },
        {
          value: 'thisMonth',
          text: this.$t('recordList.export.filter.thisMonth'),
        },
        {
          value: 'lastWeek',
          text: this.$t('recordList.export.filter.lastWeek'),
        },
        {
          value: 'thisWeek',
          text: this.$t('recordList.export.filter.thisWeek'),
        },
        {
          value: 'today',
          text: this.$t('recordList.export.filter.today'),
        },
        {
          value: 'custom',
          text: this.$t('recordList.export.filter.custom'),
        },
      ]
    },

    hasSelection () {
      return !!this.selection.length
    },

    getExportableCount () {
      // when exporting selection, only selected records are applicable
      if (this.rangeType === 'selection') {
        return this.selection.length
      }

      return this.rangeRecordCount
    },

    exportDisabled () {
      return !this.dateRangeValid || this.fields.length === 0 || !this.getExportableCount
    },

    selectedFields: {
      get () {
        return this.fields
      },

      set (selectedFields) {
        this.fields = selectedFields
      },
    },

    rangeBy: {
      get () {
        return this.filter.rangeBy
      },

      set (rangeBy) {
        this.filter.rangeBy = rangeBy
      },
    },

    exportQuery: {
      get () {
        return this.filter.query
      },
      set (v) {
        this.filter.query = v
      },
    },

    exportFilter: {
      get () {
        return this.filter.filter
      },
      set (v) {
        this.filter.filter = v
      },
    },

    rangeType: {
      get () {
        return this.filter.rangeType
      },

      set (rangeType) {
        this.filter.rangeType = rangeType
      },
    },

    range: {
      get () {
        return this.filter.date.range
      },

      set (range) {
        this.filter.date.range = range
        if (range !== 'custom') {
          this.filter.date.start = this.calcStart(moment(), range)
          this.filter.date.end = this.calcEnd(moment(), range)
        }
      },
    },

    start: {
      get () {
        return this.filter.date.start
      },

      set (start) {
        this.filter.date.start = start
        this.filter.date.range = 'custom'
      },
    },

    end: {
      get () {
        return this.filter.date.end
      },

      set (end) {
        this.filter.date.end = end
        this.filter.date.range = 'custom'
      },
    },

    currentRange () {
      if (this.rangeType === 'range') {
        return { start: this.start, end: this.end, rangeBy: this.rangeBy }
      }
      return undefined
    },
  },

  // Watchers needed for storybook
  watch: {
    filter: {
      handler () {
        this.$emit('change', this.filter)
        this.getTotalCount()
      },
      deep: true,
    },

    preselectedFields: {
      handler (value) {
        this.fields = value.filter(f => this.disabledTypes.indexOf(f.kind) < 0)
      },
      immediate: true,
    },

    filterRangeType: {
      immediate: true,
      handler (value) {
        this.filter.rangeType = value
      },
    },

    filterRangeBy: {
      handler (value) {
        this.filter.rangeBy = value
      },
      immediate: true,
    },

    startDate: {
      handler (value) {
        this.start = value
      },
    },

    endDate: {
      handler (value) {
        this.end = value
      },
    },

    dateRange: {
      handler (value) {
        this.range = value
      },
    },

    hasSelection: {
      handler (h) {
        if (h) {
          this.rangeType = 'selection'
        }
      },
      immediate: true,
    },
  },

  mounted () {
    if (this.startDate || this.endDate) {
      this.start = this.startDate
      this.end = this.endDate
    } else {
      this.range = this.dateRange
    }
  },

  methods: {
    calcStart (m, range) {
      if (range === 'lastMonth') {
        return fmtDate(m.subtract('1', 'months').startOf('month'))
      } else if (range === 'thisMonth') {
        return fmtDate(m.startOf('month'))
      } else if (range === 'lastWeek') {
        return fmtDate(m.subtract('1', 'week').startOf('week'))
      } else if (range === 'thisWeek') {
        return fmtDate(m.startOf('week'))
      } else if (range === 'today') {
        return fmtDate(m.startOf('day'))
      } else {
        throw new Error(this.$t('recordList.export.datePresetUndefined'))
      }
    },

    calcEnd (m, range) {
      if (range === 'lastMonth') {
        return fmtDate(m.subtract('1', 'months').endOf('month'))
      } else if (range === 'thisMonth') {
        return fmtDate(m.endOf('month'))
      } else if (range === 'lastWeek') {
        return fmtDate(m.subtract('1', 'week').endOf('week'))
      } else if (range === 'thisWeek') {
        return fmtDate(m.endOf('week'))
      } else if (range === 'today') {
        return fmtDate(m.endOf('day'))
      } else {
        throw new Error(this.$t('recordList.export.datePresetUndefined'))
      }
    },

    makeFilter ({ query, filter, rangeType, rangeBy, date }) {
      if (filter) {
        filter = evaluatePrefilter(filter, {
          userID: (this.$auth.user || {}).userID || NoID,
        })
      }

      query = queryToFilter(query, filter, this.module.fields)

      if (rangeType === 'all') {
        return query
      }

      if (rangeType === 'selection') {
        // @todo improve with IN operator when supported
        return this.selection.map(r => `recordID='${r}'`).join(' OR ')
      }

      let dateQuery, start, end
      if (date.start) {
        start = `DATE(${rangeBy})>='${date.start}'`
      }

      if (date.end) {
        end = `DATE(${rangeBy})<='${date.end}'`
      }

      if (start && end) {
        dateQuery = `(${start}) AND (${end})`
      } else {
        dateQuery = start || end
      }

      return query ? `(${query}) AND ${dateQuery}` : dateQuery
    },

    doExport (kind) {
      this.$emit('export', {
        ext: kind,
        fields: encodeURIComponent(this.fields.map(({ name }) => name)),
        filter: encodeURIComponent(this.makeFilter(this.filter)),
        filterRaw: encodeURIComponent(this.filter),
        timezone: encodeURIComponent(this.forTimezone ? this.exportTimezone : undefined),
      })
    },

    getTotalCount () {
      const { moduleID, namespaceID } = this.module || {}

      if (moduleID && namespaceID) {
        this.processingCount = true

        const query = this.makeFilter(this.filter)

        this.$ComposeAPI.recordList({ namespaceID, moduleID, query, limit: 1, incTotal: true })
          .then(({ filter = {} }) => {
            this.rangeRecordCount = filter.total || 0
          })
          .catch(() => {
            this.rangeRecordCount = 0
          })
          .finally(() => {
            this.processingCount = false
          })
      }
    },
  },
}
</script>
