<template>
  <div>
    <b-tab :title="$t('calendar.viewLabel')">
      <calendar-display
        v-bind="$attrs"
        v-on="$listeners"
      />
    </b-tab>

    <b-tab :title="$t('calendar.feedLabel')">
      <feed-source
        v-bind="$attrs"
        v-on="$listeners"
      />
    </b-tab>
  </div>
</template>

<script>
import FeedSource from './FeedSource'
import CalendarDisplay from './CalendarDisplay'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    CalendarDisplay,
    FeedSource,
  },

  inheritAttrs: false,
}
</script>
