<template>
  <div class="d-flex">
    <b-button
      size="lg"
      variant="light"
      class="flex-fill"
      @click="toggleModal"
    >
      {{ $t('general:label.export') }}
    </b-button>

    <b-modal
      :visible="showExportModal"
      size="lg"
      :title="$t('general:label.export')"
      no-fade
      scrollable
      footer-class="d-flex align-items-center justify-content-between"
      @hide="toggleModal"
    >
      <template v-if="showExportModal">
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
            style="height: 45vh;"
          />
        </b-form-group>

        <b-form-group>
          <b-form-radio-group
            v-model="rangeType"
            :options="rangeTypeOptions"
            stacked
            @change="getTotalCount()"
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
                  @change="getTotalCount()"
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
                  @change="getTotalCount()"
                />
              </b-form-group>
            </b-col>
          </b-row>
        </template>

        <b-row
          v-if="rangeType === 'range'"
        >
          <b-col
            md="6"
          >
            <b-form-group
              :label="$t('recordList.export.filter.from')"
              label-class="text-primary"
            >
              <c-input-date-time
                v-model="start"
                no-time
                only-past
                :labels="{
                  clear: $t('general:label.clear'),
                  none: $t('general:label.none'),
                  now: $t('general:label.now'),
                  today: $t('general:label.today'),
                }"
              />
            </b-form-group>
          </b-col>

          <b-col
            md="6"
          >
            <b-form-group
              :label="$t('recordList.export.filter.to')"
              label-class="text-primary"
            >
              <c-input-date-time
                v-model="end"
                no-time
                only-past
                :labels="{
                  clear: $t('general:label.clear'),
                  none: $t('general:label.none'),
                  now: $t('general:label.now'),
                  today: $t('general:label.today'),
                }"
              />
            </b-form-group>
          </b-col>
        </b-row>

        <template
          v-if="rangeType !== 'selection'"
        >
          <b-form-group
            :label="$t('recordList.export.filter.label')"
            label-class="text-primary"
          >
            <b-form-textarea
              v-model="exportFilter"
              :placeholder="$t('recordList.export.filter.placeholder')"
              debounce="500"
            />
          </b-form-group>
        </template>

        <b-form-group>
          <b-form-checkbox
            v-model="forTimezone"
            class="mb-2"
          >
            {{ $t('recordList.export.specifyTimezone') }}
          </b-form-checkbox>

          <c-input-select
            v-if="forTimezone"
            v-model="exportTimezone"
            :options="timezones"
            :get-option-key="getOptionKey"
            :placeholder="$t('recordList.export.timezonePlaceholder')"
          />
        </b-form-group>
      </template>

      <template #modal-footer>
        <div>
          <b-spinner
            v-if="processingCount"
            small
          />
          <span
            v-else
          >
            {{ $t('recordList.export.recordCount', { count: getExportableCount || 0 }) }}
          </span>
        </div>

        <div>
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
        </div>
      </template>
    </b-modal>
  </div>
</template>

<script>
import { compose } from '@cortezaproject/corteza-js'
import { components } from '@cortezaproject/corteza-vue'
import moment from 'moment'
import { throttle } from 'lodash'
import tz from 'compact-timezone-list'
import FieldPicker from 'corteza-webapp-compose/src/components/Common/FieldPicker'
import { getFieldFilter } from 'corteza-webapp-compose/src/lib/record-filter'

const { CInputDateTime } = components
const fmtDate = (d) => d.format('YYYY-MM-DD')

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    FieldPicker,
    CInputDateTime,
  },

  inheritAttrs: true,

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
    filter: {
      type: String,
      required: false,
      default: '',
    },
    selection: {
      type: Array,
      required: false,
      default: () => [],
    },
    selectedAllRecords: {
      type: Boolean,
      default: false,
    },
    filterRangeType: {
      type: String,
      default: 'all',
    },
    filterRangeBy: {
      type: String,
      default: 'createdAt',
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
      showExportModal: false,

      fields: [],
      forTimezone: false,
      exportTimezone: undefined,
      exportConfig: {
        rangeType: 'all',
        query: this.query,
        filter: this.filter,
        rangeBy: null,
        date: {
          range: 'lastMonth',
          start: this.calcStart(moment(), 'lastMonth'),
          end: this.calcEnd(moment(), 'lastMonth'),
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
        return this.selectedAllRecords ? this.rangeRecordCount : this.selection.length
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
        return this.exportConfig.rangeBy
      },

      set (rangeBy) {
        this.exportConfig.rangeBy = rangeBy
      },
    },

    exportFilter: {
      get () {
        return this.exportConfig.filter
      },
      set (v) {
        this.exportConfig.filter = v
        this.getTotalCount()
      },
    },

    rangeType: {
      get () {
        return this.exportConfig.rangeType
      },

      set (rangeType) {
        this.exportConfig.rangeType = rangeType
      },
    },

    range: {
      get () {
        return this.exportConfig.date.range
      },

      set (range) {
        this.exportConfig.date.range = range
        if (range !== 'custom') {
          this.exportConfig.date.start = this.calcStart(moment(), range)
          this.exportConfig.date.end = this.calcEnd(moment(), range)
        }
      },
    },

    start: {
      get () {
        return this.exportConfig.date.start
      },

      set (start) {
        this.exportConfig.date.start = start
        this.exportConfig.date.range = 'custom'
        this.getTotalCount()
      },
    },

    end: {
      get () {
        return this.exportConfig.date.end
      },

      set (end) {
        this.exportConfig.date.end = end
        this.exportConfig.date.range = 'custom'
        this.getTotalCount()
      },
    },

    currentRange () {
      if (this.rangeType === 'range') {
        return { start: this.start, end: this.end, rangeBy: this.rangeBy }
      }
      return undefined
    },
  },

  watch: {
    filter: {
      immediate: true,
      deep: true,
      handler (filter) {
        this.exportConfig.filter = filter
      },
    },

    preselectedFields: {
      handler (value) {
        if (!this.fields.length) {
          this.fields = value.filter(f => this.disabledTypes.indexOf(f.kind) < 0)
        }
      },
      immediate: true,
    },

    filterRangeType: {
      immediate: true,
      handler (value) {
        this.exportConfig.rangeType = value
      },
    },

    filterRangeBy: {
      immediate: true,
      handler (value) {
        this.exportConfig.rangeBy = value
      },
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

    hasSelection: {
      immediate: true,
      handler (h) {
        if (h) {
          this.rangeType = 'selection'
        }
      },
    },
  },

  methods: {
    toggleModal () {
      this.showExportModal = !this.showExportModal

      if (this.showExportModal) {
        this.getTotalCount()
      }
    },

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

    makeFilter ({ filter, rangeType, rangeBy, date }) {
      if (rangeType === 'all') {
        return filter
      }

      if (rangeType === 'selection') {
        return this.selectedAllRecords ? filter : this.selection.map(r => `recordID='${r}'`).join(' OR ')
      }

      let dateRangeQuery = ''

      if (date.start && date.end) {
        // If dates are the same, set range to that date
        date = { ...date }
        if (date.start === date.end) {
          date.start = moment(date.start, 'YYYY-MM-DD HH:mm').utc().format()
          date.end = moment(date.end, 'YYYY-MM-DD HH:mm').add(1, 'days').utc().format()
        }

        dateRangeQuery = getFieldFilter(rangeBy, 'DateTime', date, 'BETWEEN')
      } else if (date.start) {
        dateRangeQuery = getFieldFilter(rangeBy, 'DateTime', date.start, '>=')
      } else if (date.end) {
        dateRangeQuery = getFieldFilter(rangeBy, 'DateTime', date.end, '<=')
      }

      return filter && dateRangeQuery ? `(${filter}) AND ${dateRangeQuery}` : dateRangeQuery
    },

    doExport (kind) {
      this.$emit('export', {
        ext: kind,
        fields: encodeURIComponent(this.fields.map(({ name }) => name)),
        filter: encodeURIComponent(this.makeFilter(this.exportConfig)),
        filterRaw: encodeURIComponent(this.exportConfig),
        timezone: encodeURIComponent(this.forTimezone ? this.exportTimezone : undefined),
      })
    },

    getTotalCount: throttle(function () {
      const { moduleID, namespaceID } = this.module || {}

      if (moduleID && namespaceID) {
        const query = this.makeFilter(this.exportConfig)

        this.processingCount = true

        this.$ComposeAPI.recordList({ namespaceID, moduleID, query, limit: 1, incTotal: true })
          .then(({ filter = {} }) => {
            this.rangeRecordCount = filter.total || 0
          })
          .catch(e => {
            this.rangeRecordCount = 0
            this.toastErrorHandler(this.$t('notification:record.countFailed'))(e)
          })
          .finally(() => {
            this.processingCount = false
          })
      }
    }, 500),

    getOptionKey ({ tzCode }) {
      return tzCode
    },
  },
}
</script>
