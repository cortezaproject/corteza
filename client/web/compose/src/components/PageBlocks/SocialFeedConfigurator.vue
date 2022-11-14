<template>
  <b-tab :title="$t('socialFeed.label')">
    <fieldset
      v-if="page.moduleID && page.moduleID !== '0'"
      class="form-group"
    >
      <label>{{ $t('socialFeed.twitterProfileField') }}</label>
      <b-form-select
        v-model="options.profileSourceField"
        :options="selectOptions"
      />
    </fieldset>
    <fieldset class="form-group">
      <label>{{ $t('socialFeed.twitterProfileLabel') }}</label>
      <input
        v-model="options.profileUrl"
        class="form-control"
        type="text"
      >
    </fieldset>
  </b-tab>
</template>
<script>
import base from './base'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'SocialFeed',

  extends: base,
  computed: {
    selectOptions () {
      const options = []
      if (this.module) {
        this.module.fields
          .slice()
          .sort((a, b) => a.name.localeCompare(b.name))
          .forEach(({ name }) => {
            options.push(name)
          })
      }
      return options
    },
  },
}
</script>
