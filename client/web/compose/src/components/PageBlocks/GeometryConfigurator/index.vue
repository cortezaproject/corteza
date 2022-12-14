<template>
  <div>
    <b-tab :title="$t('geometry.viewLabel')">
      <configurator
        v-bind="$attrs"
        v-on="$listeners"
      />
    </b-tab>

    <b-tab :title="$t('geometry.feedLabel')">
      <feed-source
        v-bind="$attrs"
        v-on="$listeners"
      />
    </b-tab>
  </div>
</template>

<script>
import FeedSource from './FeedSource'
import Configurator from './Configurator'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    FeedSource,
    Configurator,
  },

  inheritAttrs: false,
}
</script>
