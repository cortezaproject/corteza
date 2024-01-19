<template>
  <div>
    <b-form-group
      :label-cols="3"
      :label="$t('calendar.reminderFeed.colorLabel')"
      horizontal
      breakpoint="md"
      label-class="text-primary"
    >
      <c-input-color-picker
        v-model="feed.options.color"
        :translations="{
          modalTitle: $t('calendar.recordFeed.colorPicker'),
          light: $t('general:swatchers.labels.light'),
          dark: $t('general:swatchers.labels.dark'),
          cancelBtnLabel: $t('general:label.cancel'),
          saveBtnLabel: $t('general:label.saveAndClose')
        }"
        :color-tooltips="colorSchemeTooltips"
        :swatchers="themeColors"
        :swatcher-labels="swatcherLabels"
      />
    </b-form-group>
  </div>
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

  data () {
    return {
      swatcherLabels: [
        'black',
        'white',
        'primary',
        'secondary',
        'success',
        'warning',
        'danger',
        'light',
        'extra-light',
        'body-bg',
        'sidebar-bg',
        'topbar-bg',
      ],
    }
  },

  computed: {
    colorSchemeTooltips () {
      return this.swatcherLabels.reduce((acc, label) => {
        acc[label] = this.$t(`general:swatchers.tooltips.${label}`)
        return acc
      }, {})
    },

    themeColors () {
      const theme = this.$Settings.get('ui.studio.themes', [])
      if (!theme.length) {
        return theme
      }

      return theme.map(theme => {
        theme.values = JSON.parse(theme.values)
        return theme
      })
    },
  },
}
</script>
