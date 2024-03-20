<template>
  <wrap
    v-bind="$props"
    v-on="$listeners"
    @refreshBlock="refresh"
  >
    <div class="d-flex flex-column calendar-container p-2 h-100">
      <div v-if="!header.hide">
        <div
          v-if="!header.hidePrevNext || !header.hideTitle"
          class="d-flex align-items-baseline justify-content-center mb-2"
        >
          <b-btn
            v-if="!header.hidePrevNext"
            variant="link"
            class="text-dark"
            @click="api().prev()"
          >
            <font-awesome-icon :icon="['fas', 'angle-left']" />
          </b-btn>
          <span
            v-if="!header.hideTitle"
            class="h5"
          >
            {{ title }}
          </span>
          <b-btn
            v-if="!header.hidePrevNext"
            variant="link"
            class="text-dark"
            @click="api().next()"
          >
            <font-awesome-icon :icon="['fas', 'angle-right']" />
          </b-btn>
        </div>
        <b-row
          no-gutters
        >
          <b-col
            cols="12"
            sm="10"
            md="9"
            lg="8"
            xl="9"
            class="d-flex justify-content-sm-start justify-content-center flex-wrap"
          >
            <b-btn
              v-for="view in views"
              :key="view"
              variant="light"
              class="mr-1 mb-1"
              @click="api().changeView(view)"
            >
              {{ $t(`calendar.view.${view}`) }}
            </b-btn>
          </b-col>
          <b-col
            v-if="!header.hideToday && !header.hide"
            cols="12"
            sm="2"
            md="3"
            lg="4"
            xl="3"
            class="d-flex justify-content-end"
          >
            <b-btn
              variant="light"
              class="mb-1 w-100"
              @click="api().today()"
            >
              {{ $t(`calendar.today`) }}
            </b-btn>
          </b-col>
        </b-row>
      </div>

      <div
        :ref="`cc-${blockIndex}`"
        class="d-flex flex-column flex-fill"
      >
        <div
          v-if="processing"
          class="d-flex align-items-center justify-content-center h-100"
        >
          <b-spinner />
        </div>

        <full-calendar
          v-show="!processing"
          :ref="`fc-${blockIndex}`"
          :key="key"
          :height="getHeight()"
          :events="events"
          v-bind="config"
          class="flex-fill"
          @eventClick="handleEventClick"
        />
      </div>
    </div>
  </wrap>
</template>

<script>
import moment from 'moment'
import { mapGetters, mapActions } from 'vuex'
import axios from 'axios'
import base from './base'
import FullCalendar from '@fullcalendar/vue'
import dayGridPlugin from '@fullcalendar/daygrid'
import timeGridPlugin from '@fullcalendar/timegrid'
import listPlugin from '@fullcalendar/list'
import { compose, NoID } from '@cortezaproject/corteza-js'
import { BootstrapTheme } from '@fullcalendar/bootstrap'
import { createPlugin } from '@fullcalendar/core'
import { evaluatePrefilter, isFieldInFilter } from 'corteza-webapp-compose/src/lib/record-filter'

/**
 * FullCalendar Corteza theme definition.
 */
export class CortezaTheme extends BootstrapTheme {}
CortezaTheme.prototype.classes.widget = 'corteza-unthemed'
CortezaTheme.prototype.classes.button = 'btn btn-outline-primary'

CortezaTheme.prototype.baseIconClass = 'fc-icon'
CortezaTheme.prototype.iconClasses = {
  close: 'fc-icon-x',
  prev: 'fc-icon-chevron-left',
  next: 'fc-icon-chevron-right',
  prevYear: 'fc-icon-chevrons-left',
  nextYear: 'fc-icon-chevrons-right',
}

CortezaTheme.prototype.iconOverrideOption = 'buttonIcons'
CortezaTheme.prototype.iconOverrideCustomButtonOption = 'icon'
CortezaTheme.prototype.iconOverridePrefix = 'fc-icon-'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    FullCalendar,
  },

  extends: base,

  data () {
    return {
      processing: false,
      show: false,

      events: [],
      locale: undefined,
      title: '',

      loaded: {
        start: null,
        end: null,
      },

      refreshing: false,

      cancelTokenSource: axios.CancelToken.source(),
    }
  },

  computed: {
    ...mapGetters({
      pages: 'page/set',
    }),

    config () {
      return {
        header: false,
        themeSystem: 'corteza',
        defaultView: this.options.defaultView,
        editable: false,
        eventLimit: true,
        locale: this.locale,
        // @todo could be loaded on demand
        plugins: [
          dayGridPlugin,
          timeGridPlugin,
          listPlugin,
          createPlugin({
            themeClasses: {
              corteza: CortezaTheme,
            },
          }),
        ],

        // Handle event fetching when view/date-range changes
        datesRender: ({ view: { activeStart, activeEnd, title } = {} } = {}) => {
          this.loadEvents(moment(activeStart), moment(activeEnd))
          // eslint-disable-next-line vue/no-side-effects-in-computed-properties
          this.title = title
        },
      }
    },

    header () {
      return this.block.options.header
    },

    views () {
      if (this.header.hide) {
        return []
      }

      return this.block.reorderViews(this.header.views)
    },
  },

  watch: {
    options: {
      deep: true,
      handler () {
        this.updateSize()
        this.refresh()
      },
    },

    'block.xywh': {
      deep: true,
      handler () {
        this.updateSize()
      },
    },
  },

  created () {
    this.changeLocale(this.currentLanguage)
    this.refreshBlock(this.refresh)
  },

  mounted () {
    this.createEvents()
  },

  beforeDestroy () {
    this.setDefaultValues()
    this.abortRequests()
    this.destroyEvents()
  },

  methods: {
    ...mapActions({
      findModuleByID: 'module/findByID',
    }),

    createEvents () {
      this.$root.$on('module-records-updated', this.refreshOnRelatedRecordsUpdate)
      this.$root.$on('record-field-change', this.refetchOnPrefilterValueChange)
    },

    refetchOnPrefilterValueChange ({ fieldName }) {
      const { feeds } = this.options

      if (feeds.some(({ options }) => isFieldInFilter(fieldName, options.prefilter))) {
        this.refresh()
      }
    },

    updateSize () {
      this.$nextTick(() => {
        this.api() && this.api().updateSize()
      })
    },

    refreshOnRelatedRecordsUpdate ({ moduleID, notPageID }) {
      this.options.feeds.forEach((feed) => {
        const { moduleID: feedModuleID } = feed.options

        if (feedModuleID) {
          if (feedModuleID === moduleID && this.page.pageID !== notPageID) {
            this.refresh()
          }
        }
      })
    },

    /**
     * Helper method to load requested locale.
     * See https://github.com/fullcalendar/fullcalendar/tree/master/packages/core/src/locales
     * for a full list
     * @param {String} lng Locale tag.
     */
    changeLocale (lng = 'en-gb') {
      // fc doesn't provide a en locale
      if (lng === 'en') {
        lng = 'en-gb'
      }

      this.locale = require(`@fullcalendar/core/locales/${lng}`)
    },

    // Proxy to the FC API
    api () {
      if (this.$refs[`fc-${this.blockIndex}`]) {
        return this.$refs[`fc-${this.blockIndex}`].getApi()
      }
    },

    /**
     * Loads & preps fc events from `start` to `end` for all defined feeds.
     * @param {Moment} start Start date
     * @param {Moment} end End date
     */
    loadEvents (start, end) {
      if (!start || !end) {
        return
      }

      if (start.isSame(this.loaded.start) && end.isSame(this.loaded.end) && !this.refreshing) {
        return
      }

      this.loaded.start = start
      this.loaded.end = end

      this.events = []

      this.processing = true

      Promise.all(this.options.feeds.map(feed => {
        switch (feed.resource) {
          case compose.PageBlockCalendar.feedResources.record:
            return this.findModuleByID({ namespace: this.namespace, moduleID: feed.options.moduleID })
              .then(module => {
                const ff = {
                  ...feed,
                  options: { ...feed.options },
                }

                // Interpolate prefilter variables
                if (ff.options.prefilter) {
                  ff.options.prefilter = evaluatePrefilter(ff.options.prefilter, {
                    record: this.record,
                    user: this.$auth.user || {},
                    recordID: (this.record || {}).recordID || NoID,
                    ownerID: (this.record || {}).ownedBy || NoID,
                    userID: (this.$auth.user || {}).userID || NoID,
                  })
                }

                return compose.PageBlockCalendar.RecordFeed(this.$ComposeAPI, module, this.namespace, ff, this.loaded, { cancelToken: this.cancelTokenSource.token })
                  .then(events => {
                    events = this.setEventColors(events, ff)
                    this.events.push(...events)
                  })
              })
          case compose.PageBlockCalendar.feedResources.reminder:
            return compose.PageBlockCalendar.ReminderFeed(this.$SystemAPI, this.$auth.user, feed, this.loaded)
              .then(events => {
                events = this.setEventColors(events, feed)
                this.events.push(...events)
              })
        }
      }))
        .finally(() => {
          setTimeout(() => {
            this.processing = false
            this.refreshing = false

            this.updateSize()
          }, 300)
        })
    },

    /**
     * Based on event type, perform some action.
     * @param {Event} event Fullcalendar event object
     */
    handleEventClick ({ event: { extendedProps: { recordID, moduleID } } }) {
      if (!moduleID || !recordID) {
        return
      }

      const page = this.pages.find(p => p.moduleID === moduleID)
      if (!page) {
        return
      }

      const route = { name: 'page.record', params: { recordID, pageID: page.pageID } }

      if (this.options.eventDisplayOption === 'modal' || this.inModal) {
        this.$root.$emit('show-record-modal', {
          recordID,
          recordPageID: page.pageID,
        })
      } else if (this.options.eventDisplayOption === 'newTab') {
        window.open(this.$router.resolve(route).href)
      } else {
        this.$router.push(route)
      }
    },

    getHeight () {
      if (this.$refs[`cc-${this.blockIndex}`]) {
        return this.$refs[`cc-${this.blockIndex}`].clientHeight
      }
      return 'auto'
    },

    refresh () {
      this.refreshing = true
      new Promise(resolve => resolve(this.api().refetchEvents()))
        .then(() => this.key++)
        .catch(() => this.toastErrorHandler(this.$t('notification:page.block.calendar.eventFetchFailed')))
    },

    setEventColors (events, feed) {
      return events.map(event => {
        event.backgroundColor = feed.options.color
        return event
      })
    },

    setDefaultValues () {
      this.processing = false
      this.show = false
      this.events = []
      this.locale = undefined
      this.title = ''
      this.loaded = {}
      this.refreshing = false
    },

    abortRequests () {
      this.cancelTokenSource.cancel(`cancel-record-list-request-${this.block.blockID}`)
    },

    destroyEvents () {
      this.$root.$off('module-records-updated', this.refreshOnRelatedRecordsUpdate)
      this.$root.$off('record-field-change', this.refetchOnPrefilterValueChange)
    },
  },
}
</script>
<style lang="scss" scoped>
@import '~@fullcalendar/core/main.css';
@import '~@fullcalendar/daygrid/main.css';
@import '~@fullcalendar/timegrid/main.css';
@import '~@fullcalendar/list/main.css';

</style>
<style lang="scss">
.calendar-container {
  .fc-content, .event-record {
    cursor: pointer;
  }

  .fc-day-header {
    white-space: pre-wrap;
  }
}

.fc-popover {
  .fc-header {
    padding: 0.5rem;
  }

  .fc-body {
    padding: 0;

    .fc-event-container {
      display: flex;
      flex-direction: column;
      gap: 0.25rem;
    }
  }
}
</style>
