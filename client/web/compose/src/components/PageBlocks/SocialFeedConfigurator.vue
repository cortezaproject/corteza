<template>
  <b-tab :title="$t('socialFeed.label')">
    <b-form-group
      v-if="page.moduleID && page.moduleID !== '0'"
      :label="$t('socialFeed.twitterProfileField')"
      label-class="text-primary"
    >
      <b-form-select
        v-model="options.profileSourceField"
        :options="selectOptions"
      />
    </b-form-group>

    <b-form-group
      :label="$t('socialFeed.twitterProfileLabel')"
      label-class="text-primary"
    >
      <b-input
        v-model="options.profileUrl"
      />
    </b-form-group>
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
