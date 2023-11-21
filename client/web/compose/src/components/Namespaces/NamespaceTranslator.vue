<template>
  <c-translator-button
    v-if="canManageResourceTranslations && resourceTranslationsEnabled"
    :tooltip="$t('tooltip')"
    v-bind="$props"
    :resource="resource"
    :titles="titles"
    :fetcher="fetcher"
    :updater="updater"
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
    keyPrefix: 'resources.namespace',
  },

  props: {
    namespace: {
      type: compose.Namespace,
      required: true,
    },

    disabled: {
      type: Boolean,
      default: () => false,
    },

    buttonVariant: {
      type: String,
      default: () => 'primary',
    },

    highlightKey: {
      type: String,
      default: '',
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
      const { namespaceID } = this.namespace
      return `compose:namespace/${namespaceID}`
    },

    titles () {
      const { namespaceID, slug: handle } = this.namespace
      const titles = {}

      titles[this.resource] = this.$t('title', { handle: handle || namespaceID })
      return titles
    },

    fetcher () {
      const { namespaceID } = this.namespace

      return () => {
        return this.$ComposeAPI.namespaceListTranslations({ namespaceID })
        // @todo pass set of translations to the resource object
        // The logic there needs to be implemented; the idea is to decode
        // values from the resource object to the set of translations)
      }
    },

    updater () {
      const { namespaceID } = this.namespace
      return translations => {
        return this.$ComposeAPI
          .namespaceUpdateTranslations({ namespaceID, translations })
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

            let tr

            tr = find('name')
            if (tr !== undefined) {
              this.namespace.name = tr.message
            }

            tr = find('meta.subtitle')
            if (tr !== undefined) {
              this.namespace.meta.subtitle = tr.message
            }

            tr = find('meta.description')
            if (tr !== undefined) {
              this.namespace.meta.description = tr.message
            }
          })
      }
    },
  },
}
</script>
