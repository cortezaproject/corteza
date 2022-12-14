<template>
  <fieldset class="form-group">
    <div
      v-for="(feed, i) in options.feeds"
      :key="i"
    >
      <div
        class="d-flex justify-content-between mb-3"
      >
        <h5>
          {{ $t('geometry.feedLabel') }}
        </h5>

        <template
          v-if="feed.resource"
        >
          <b-button
            variant="outline-danger"
            class="border-0"
            @click="onRemoveFeed(i)"
          >
            <font-awesome-icon :icon="['far', 'trash-alt']" />
          </b-button>
        </template>
      </div>

      <b-form-group horizontal>
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
      class="btn btn-url test-feed-add"
      @click.prevent="handleAddButton"
    >
      {{ $t('geometry.addSource') }}
    </b-button>
  </fieldset>
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
      return Object.entries(compose.PageBlockGeometry.feedResources).map(([key, value]) => ({
        value,
        text: this.$t(`geometry.${key}Feed.optionLabel`),
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
      this.block.options.feeds.push(compose.PageBlockGeometry.makeFeed())
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
