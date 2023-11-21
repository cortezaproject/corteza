<template>
  <c-translator-button
    v-if="canManageResourceTranslations && resourceTranslationsEnabled"
    v-bind="$props"
    :tooltip="$t('tooltip')"
    :resource="resource"
    :titles="titles"
    :fetcher="fetcher"
    :updater="updater"
    class="ml-auto py-1 px-3"
  />
</template>

<script>
import { compose } from '@cortezaproject/corteza-js'
import { mapGetters } from 'vuex'
import CTranslatorButton from 'corteza-webapp-compose/src/components/Translator/CTranslatorButton'

export default {
  components: {
    CTranslatorButton,
  },

  i18nOptions: {
    namespaces: 'resource-translator',
    keyPrefix: 'resources.chart',
  },

  props: {
    field: {
      type: String,
      default: '',
    },

    chart: {
      type: compose.Chart,
      required: true,
    },

    highlightKey: {
      type: String,
      default: '',
    },

    disabled: {
      type: Boolean,
      default: () => false,
    },
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canManageResourceTranslations () {
      return this.can('compose/', 'resource-translations.manage')
    },

    resource () {
      const { namespaceID, chartID } = this.chart
      return `compose:chart/${namespaceID}/${chartID}`
    },

    titles () {
      const { chartID, handle } = this.chart
      const titles = {}

      titles[this.resource] = this.$t('title', { handle: handle || chartID })

      return titles
    },

    fetcher () {
      const { namespaceID, chartID } = this.chart

      return () => {
        return this.$ComposeAPI.chartListTranslations({ namespaceID, chartID })
        // @todo pass set of translations to the resource object
        // The logic there needs to be implemented; the idea is to decode
        // values from the resource object to the set of translations)
      }
    },

    updater () {
      const { namespaceID, chartID } = this.chart

      return translations => {
        return this.$ComposeAPI
          .chartUpdateTranslations({ namespaceID, chartID, translations })
          // re-fetch translations, sanitized and stripped
          .then(() => this.fetcher())
          .then((translations) => {
            // When translations are successfully saved,
            // scan changes and apply them back to the passed object
            // not the most elegant solution but is saves us from
            // handling the resource on multiple places
            //
            //
            // @todo move this to Namespace* classes
            // the logic there needs to be implemented; the idea is to encode
            // values from the set of translations back to the resource object
            const find = (key) => {
              return translations.find(t => t.key === key && t.lang === this.currentLanguage && t.resource === this.resource)
            }

            let tr

            const [report = {}] = this.chart.config.reports
            tr = find('yAxis.label')
            if (tr !== undefined) {
              this.$set(report.yAxis, 'label', tr.message)
            }

            report.metrics.forEach((metric) => {
              tr = find(`metrics.${metric.metricID}.label`)
              if (tr) {
                this.$set(metric, 'label', tr.message)
              }
            })

            if (Array.isArray(report.dimensions)) {
              report.dimensions.forEach(d => {
                if (!Array.isArray(d.meta.steps)) {
                  return
                }

                d.meta.steps.forEach((step) => {
                  tr = find(`dimensions.${d.dimensionID}.meta.steps.${step.stepID}.label`)
                  if (tr) {
                    this.$set(step, 'label', tr.message)
                  }
                })
              })
            }

            return this.chart
          })
          .then(chart => {
            this.$emit('update:chart', chart)
          })
      }
    },
  },
}
</script>
