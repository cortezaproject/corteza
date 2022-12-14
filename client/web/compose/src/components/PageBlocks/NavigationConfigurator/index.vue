<template>
  <div>
    <b-tab :title="$t('navigation.label')">
      <navigation-tab
        v-bind="$attrs"
        v-on="$listeners"
      />
    </b-tab>
  </div>
</template>

<script>
import NavigationTab from './Navigation'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: { NavigationTab },

  inheritAttrs: false,
}
</script>
