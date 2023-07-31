<template>
  <b-container fluid>
    <b-row>
      <b-col
        cols="12"
      >
        <b-button
          v-for="(type) in types"
          :key="type.label"
          :disabled="isOptionDisabled(type)"
          variant="outline-light"
          class="mr-2 mb-2 text-dark"
          @click="$emit('select', type.block)"
          @mouseover="current = type.image"
          @mouseleave="current = undefined"
        >
          {{ type.label }}
        </b-button>
      </b-col>

      <b-col
        cols="12"
        style="height: 50vh;"
        class="d-flex align-items-center"
      >
        <b-img
          v-if="current"
          :src="current"
          center
          fluid
          class="mx-auto mh-100"
        />
      </b-col>

      <hr
        v-if="existingBlocks.length"
        class="w-100"
      >

      <b-col
        v-if="existingBlocks.length"
        cols="12"
      >
        <b-input-group class="d-flex w-100">
          <c-input-select
            v-model="selectedExistingBlock"
            :get-option-label="getBlockLabel"
            :get-option-key="b => b.blockID"
            :options="existingBlocks"
            :placeholder="$t('selector.selectableBlocks.placeholder')"
          />

          <b-input-group-append>
            <b-button
              v-b-tooltip.hover="{ title: $t('selector.tooltip.clone.noRef'), container: '#body' }"
              variant="light"
              :disabled="!selectedExistingBlock"
              class="d-flex align-items-center"
              @click="$emit('select', selectedExistingBlock.clone())"
            >
              <font-awesome-icon
                :icon="['far', 'clone']"
              />
            </b-button>

            <b-button
              v-b-tooltip.hover="{ title: $t('selector.tooltip.clone.ref'), container: '#body' }"
              variant="light"
              :disabled="!selectedExistingBlock"
              class="d-flex align-items-center"
              @click="$emit('select', selectedExistingBlock)"
            >
              <font-awesome-icon
                :icon="['far', 'copy']"
              />
            </b-button>
          </b-input-group-append>
        </b-input-group>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import * as images from '../../../../assets/PageBlocks'
import { compose } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  props: {
    recordPage: {
      type: Boolean,
      default: false,
    },

    disabledKinds: {
      type: Array,
      default: () => [],
    },

    existingBlocks: {
      type: Array,
      default: () => [],
    },
  },

  data () {
    return {
      current: undefined,

      selectedExistingBlock: undefined,

      types: [
        {
          label: this.$t('automation.label'),
          block: new compose.PageBlockAutomation(),
          image: images.Automation,
        },
        {
          label: this.$t('calendar.label'),
          block: new compose.PageBlockCalendar(),
          image: images.Calendar,
        },
        {
          label: this.$t('chart.label'),
          block: new compose.PageBlockChart(),
          image: images.Chart,
        },
        {
          label: this.$t('content.label'),
          block: new compose.PageBlockContent(),
          image: images.Content,
        },
        {
          label: this.$t('comment.label'),
          block: new compose.PageBlockComment(),
          image: images.Comment,
        },
        {
          label: this.$t('file.label'),
          block: new compose.PageBlockFile(),
          image: images.File,
        },
        {
          label: this.$t('iframe.label'),
          block: new compose.PageBlockIFrame(),
          image: images.IFrame,
        },
        {
          label: this.$t('metric.label'),
          block: new compose.PageBlockMetric(),
          image: images.Metric,
        },
        {
          label: this.$t('record.label'),
          block: new compose.PageBlockRecord(),
          image: images.Record,
          recordPageOnly: true,
        },
        {
          label: this.$t('recordList.label'),
          block: new compose.PageBlockRecordList(),
          image: images.RecordList,
        },
        {
          label: this.$t('recordOrganizer.label'),
          block: new compose.PageBlockRecordOrganizer(),
          image: images.RecordOrganizer,
        },
        {
          label: this.$t('recordRevisions.label'),
          block: new compose.PageBlockRecordRevisions(),
          image: images.RecordRevisions,
          recordPageOnly: true,
        },
        {
          label: this.$t('report.label'),
          block: new compose.PageBlockReport(),
          image: images.Report,
        },
        {
          label: this.$t('socialFeed.label'),
          block: new compose.PageBlockSocialFeed(),
          image: images.SocialFeed,
        },
        {
          label: this.$t('progress.label'),
          block: new compose.PageBlockProgress(),
          image: images.Progress,
        },
        {
          label: this.$t('nylas.label'),
          block: new compose.PageBlockNylas(),
          image: images.Nylas,
        },
        {
          label: this.$t('geometry.label'),
          block: new compose.PageBlockGeometry(),
          image: images.Geometry,
        },
        {
          label: this.$t('tabs.label'),
          block: new compose.PageBlockTab(),
          image: images.Tabs,
        },
        {
          label: this.$t('navigation.label'),
          block: new compose.PageBlockNavigation(),
          image: images.Navigation,
        },
      ],
    }
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
    isOptionDisabled (type) {
      return (!this.recordPage && type.recordPageOnly) || this.disabledKinds.includes(type.block.kind)
    },

    getBlockLabel ({ title, kind }) {
      return title || kind
    },

    setDefaultValues () {
      this.current = undefined
      this.selectedExistingBlock = undefined
      this.types = []
    },
  },
}
</script>
