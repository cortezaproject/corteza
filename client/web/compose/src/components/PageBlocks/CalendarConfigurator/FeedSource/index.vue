<template>
  <b-row>
    <!-- Feed list -->
    <b-col
      v-for="(feed, i) in options.feeds"
      :key="i"
      cols="12"
      class="p-0"
    >
      <b-card
        class="list-background mx-3 mb-3"
      >
        <h5 class="d-flex align-items-center mb-3">
          {{ $t('calendar.source.label') }} {{ i + 1 }}

          <c-input-confirm
            show-icon
            class="ml-auto mt-1"
            @confirmed="onRemoveFeed(i)"
          />
        </h5>

        <!-- define feed resource; eg. module, reminders, google calendar, ... -->
        <b-form-group
          :label="$t('calendar.eventSource')"
          label-class="text-primary"
        >
          <c-input-select
            v-model="feed.resource"
            :options="feedSources"
            :clearable="false"
            label="text"
            :reduce="o => o.value"
          />
        </b-form-group>

        <component
          :is="configurator(feed)"
          v-if="feed.resource && configurator(feed)"
          :feed="feed"
          :modules="modules"
          :page="page"
          :record="record"
          :module="module"
        />
      </b-card>
    </b-col>

    <b-col cols="12">
      <b-button
        variant="primary"
        class="test-feed-add"
        @click.prevent="handleAddButton"
      >
        {{ $t('calendar.addEventsSource') }}
      </b-button>
    </b-col>
  </b-row>
</template>
<script>
import { mapGetters } from 'vuex'
import base from '../../base'
import * as configs from './configs'
import { compose } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  extends: base,

  props: {
    page: {
      type: compose.Page,
      required: true,
    },
  },

  computed: {
    ...mapGetters({
      modules: 'module/set',
    }),

    /**
     * Provides a set of available feed sources.
     * @returns {Array}
     */
    feedSources () {
      return Object.entries(compose.PageBlockCalendar.feedResources).map(([key, value]) => ({
        value,
        text: this.$t(`calendar.${key}Feed.optionLabel`),
      }))
    },
  },

  created () {
    if (this.options.feeds.length === 0) {
      this.block.options.feeds = []
    }
  },

  methods: {
    /**
     * Handles feed removal
     * @param {Number} i Feed's index
     */
    onRemoveFeed (i) {
      this.block.options.feeds.splice(i, 1)
    },

    /**
     * Handles feed's addition
     */
    handleAddButton () {
      this.block.options.feeds.push(compose.PageBlockCalendar.makeFeed())
    },

    /**
     * configurator uses feed's resource to determine what configurator to use.
     * If it can't find an apropriate component, undefined is returned
     * @param {Feed} feed Feed in qestion
     * @returns {Component|undefined}
     */
    configurator (feed) {
      if (!feed.resource) {
        return
      }
      const r = feed.resource.split(':').pop()
      return configs[r]
    },
  },
}
</script>

<style lang="scss" scoped>
.list-background {
  background-color: var(--body-bg);
}
</style>
