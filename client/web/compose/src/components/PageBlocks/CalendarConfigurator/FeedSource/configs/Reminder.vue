<template>
  <b-row>
    <b-col
      cols="12"
      lg="6"
    >
      <b-form-group
        :label="$t('calendar.reminderFeed.colorLabel')"
        label-class="text-primary"
      >
        <c-input-color-picker
          v-model="feed.options.color"
          :translations="{
            modalTitle: $t('calendar.recordFeed.colorPicker'),
            light: $t('general:themes.labels.light'),
            dark: $t('general:themes.labels.dark'),
            cancelBtnLabel: $t('general:label.cancel'),
            saveBtnLabel: $t('general:label.saveAndClose')
          }"
          :theme-settings="themeSettings"
        />
      </b-form-group>
    </b-col>
  </b-row>
</template>

<script>
import base from './base'
import { components } from '@cortezaproject/corteza-vue'
const { CInputColorPicker } = components

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    CInputColorPicker,
  },

  extends: base,

  computed: {
    themeSettings () {
      return this.$Settings.get('ui.studio.themes', [])
    },
  },
}
</script>
