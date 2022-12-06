<template>
  <b-container fluid>
    <b-row>
      <b-col
        cols="4"
      >
        <b-list-group>
          <b-list-group-item
            v-for="(type) in types"
            :key="type.label"
            :disabled="!recordPage && type.recordPageOnly"
            button
            @click="$emit('select', type.block)"
            @mouseover="current = type.image"
          >
            {{ type.label }}
          </b-list-group-item>
        </b-list-group>
      </b-col>
      <b-col
        cols="8"
        class="my-auto"
      >
        <b-img
          v-if="current"
          fluid
          thumbnail
          :src="current"
        />
      </b-col>
    </b-row>
    <b-row
      class="border-top mt-2"
    >
      <b-col>
        <div
          class="mt-2"
        >
          {{ $t('selectBlockFootnote') }}
        </div>
      </b-col>
    </b-row>
  </b-container>
</template>
<script>
import { compose } from '@cortezaproject/corteza-js'
import * as images from '../../../../assets/PageBlocks'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  props: {
    recordPage: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      current: undefined,
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
      ],
    }
  },
}
</script>
