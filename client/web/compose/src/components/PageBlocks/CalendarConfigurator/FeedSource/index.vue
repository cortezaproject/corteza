<template>
  <div>
    <!-- Feed list -->
    <div
      v-for="(feed, i) in options.feeds"
      :key="i"
    >
      <div
        v-if="feed.resource"
        class="d-flex justify-content-end mb-3"
      >
        <c-input-confirm
          v-if="feed.resource"
          show-icon
          size="md"
          @confirmed="onRemoveFeed(i)"
        />
      </div>

      <!-- define feed resource; eg. module, reminders, google calendar, ... -->
      <b-form-group
        :label="$t('calendar.eventSource')"
        :label-cols="3"
        horizontal
        breakpoint="md"
        label-class="text-primary"
      >
        <b-form-select
          v-model="feed.resource"
          :options="feedSources"
        >
          <template slot="first">
            <option
              value=""
              :disabled="true"
            >
              {{ $t('calendar.feedPlaceholder') }}
            </option>
          </template>
        </b-form-select>
      </b-form-group>

      <b-form-group horizontal>
        <!-- source configurator -->
        <component
          :is="configurator(feed)"
          v-if="feed.resource && configurator(feed)"
          :feed="feed"
          :modules="modules"
        />
      </b-form-group>

      <hr>
    </div>

    <b-button
      variant="primary"
      class="test-feed-add"
      @click.prevent="handleAddButton"
    >
      {{ $t('calendar.addEventsSource') }}
    </b-button>
  </div>
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

  components: {
  },

  extends: base,

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
