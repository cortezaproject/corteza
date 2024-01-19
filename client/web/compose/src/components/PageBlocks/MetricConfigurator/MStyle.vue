<template>
  <b-card no-body>
    <slot name="title" />

    <fieldset>
      <b-form-group
        :label="$t('metric.editStyle.color')"
        label-class="text-primary"
      >
        <c-input-color-picker
          v-model="options.color"
          :translations="{
            modalTitle: $t('metric.editStyle.colorPicker'),
            light: $t('general:swatchers.labels.light'),
            dark: $t('general:swatchers.labels.dark'),
            cancelBtnLabel: $t('general:label.cancel'),
            saveBtnLabel: $t('general:label.saveAndClose')
          }"
          :color-tooltips="colorSchemeTooltips"
          :swatchers="themeColors"
          :swatcher-labels="swatcherLabels"
          class="mb-1"
        />
      </b-form-group>

      <b-form-group
        :label="$t('metric.editStyle.backgroundColor')"
        label-class="text-primary"
      >
        <c-input-color-picker
          v-model="options.backgroundColor"
          :translations="{
            modalTitle: $t('geometry.recordFeed.colorPicker'),
            light: $t('general:swatchers.labels.light'),
            dark: $t('general:swatchers.labels.dark'),
            cancelBtnLabel: $t('general:label.cancel'),
            saveBtnLabel: $t('general:label.saveAndClose')
          }"
          :color-tooltips="colorSchemeTooltips"
          :swatchers="themeColors"
          :swatcher-labels="swatcherLabels"
          class="mb-1"
        />
      </b-form-group>

      <b-form-group
        :label="$t('metric.editStyle.fontSize')"
        label-class="text-primary"
      >
        <b-form-input
          v-model="options.fontSize"
          type="number"
          placeholder="16"
          min="0.1"
          step="0.1"
          class="mb-1"
        />
      </b-form-group>
    </fieldset>
  </b-card>
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'
const { CInputColorPicker } = components

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    CInputColorPicker,
  },

  props: {
    options: {
      type: Object,
      required: true,
      default: () => ({}),
    },
  },

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
