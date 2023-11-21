<template>
  <c-translator-button
    v-if="canManageResourceTranslations && resourceTranslationsEnabled"
    v-bind="$props"
    :tooltip="$t('tooltip')"
    :resource="resource"
    :titles="titles"
    :fetcher="fetcher"
    :updater="updater"
  />
</template>

<script>
import { compose, NoID } from '@cortezaproject/corteza-js'
import { mapGetters } from 'vuex'
import CTranslatorButton from 'corteza-webapp-compose/src/components/Translator/CTranslatorButton'

export default {
  components: {
    CTranslatorButton,
  },

  i18nOptions: {
    namespaces: 'resource-translator',
    keyPrefix: 'resources.page-layout',
  },

  props: {
    pageLayout: {
      type: compose.PageLayout,
      required: true,
    },

    buttonVariant: {
      type: String,
      default: () => 'primary',
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
      const { pageID, namespaceID, pageLayoutID } = this.pageLayout
      return `compose:page-layout/${namespaceID}/${pageID}/${pageLayoutID}`
    },

    titles () {
      const titles = {}

      const { pageID, handle, meta } = this.pageLayout
      titles[this.resource] = this.$t('title', { handle: handle || meta.title || pageID })

      return titles
    },

    fetcher () {
      const { pageLayoutID, pageID, namespaceID } = this.pageLayout

      return () => {
        return this.$ComposeAPI
          .pageLayoutListTranslations({ namespaceID, pageID, pageLayoutID })
          .then(set => {
            return set
          })
      }
    },

    updater () {
      const { pageID, namespaceID, pageLayoutID } = this.pageLayout

      return translations => {
        return this.$ComposeAPI
          .pageLayoutUpdateTranslations({ namespaceID, pageID, pageLayoutID, translations })
          // re-fetch translations, sanitized and stripped
          .then(() => this.fetcher())
          .then((translations) => {
            // When translations are successfully saved,
            // scan changes and apply them back to the passed object
            // not the most elegant solution but is saves us from
            // handling the resource on multiple places
            //
            // @todo move this to Namespace* classes
            // the logic there needs to be implemented; the idea is to encode
            // values from the set of translations back to the resource object

            const find = (key) => {
              return translations.find(t => t.key === key && t.lang === this.currentLanguage && t.resource === this.resource)
            }

            let tr = find('title')
            if (tr !== undefined) {
              this.pageLayout.meta.title = tr.message
            }

            tr = find('description')
            if (tr !== undefined) {
              this.pageLayout.meta.description = tr.message
            }

            // Refresh page buttons for record pages
            if (this.pageLayout.moduleID && this.pageLayout.moduleID !== NoID) {
              tr = find('config.buttons.new.label')
              if (tr) {
                this.$set(this.pageLayout.config.buttons.new, 'label', tr.message)
              }

              tr = find('config.buttons.edit.label')
              if (tr) {
                this.$set(this.pageLayout.config.buttons.edit, 'label', tr.message)
              }

              tr = find('config.buttons.submit.label')
              if (tr) {
                this.$set(this.pageLayout.config.buttons.submit, 'label', tr.message)
              }

              tr = find('config.buttons.delete.label')
              if (tr) {
                this.$set(this.pageLayout.config.buttons.delete, 'label', tr.message)
              }

              tr = find('config.buttons.clone.label')
              if (tr) {
                this.$set(this.pageLayout.config.buttons.clone, 'label', tr.message)
              }

              tr = find('config.buttons.back.label')
              if (tr) {
                this.$set(this.pageLayout.config.buttons.back, 'label', tr.message)
              }
            }

            return this.page
          })
          .then(page => {
            this.$emit('update:pageLayout', page)
          })
      }
    },
  },
}
</script>
