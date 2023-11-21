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
import { compose } from '@cortezaproject/corteza-js'
import { mapGetters } from 'vuex'
import CTranslatorButton from 'corteza-webapp-compose/src/components/Translator/CTranslatorButton'
import moduleResTr from 'corteza-webapp-compose/src/lib/resource-translations/module'

export default {
  components: {
    CTranslatorButton,
  },

  i18nOptions: {
    namespaces: 'resource-translator',
    keyPrefix: 'resources.module',
  },

  props: {
    module: {
      type: compose.Module,
      required: true,
    },

    highlightKey: {
      type: String,
      default: '',
    },

    buttonVariant: {
      type: String,
      default: () => 'primary',
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
      const { moduleID, namespaceID } = this.module
      return `compose:module/${namespaceID}/${moduleID}`
    },

    titles () {
      const { moduleID, handle, namespaceID, fields } = this.module
      const titles = {}

      titles[this.resource] = this.$t('title', { handle: handle || moduleID })

      fields.forEach(({ fieldID, name }) => {
        titles[`compose:module-field/${namespaceID}/${moduleID}/${fieldID}`] = this.$t('field.title', { name })
      })

      return titles
    },

    fetcher () {
      const { moduleID, namespaceID } = this.module

      return () => {
        return this.$ComposeAPI.moduleListTranslations({ namespaceID, moduleID })
        // @todo pass set of translations to the resource object
        // The logic there needs to be implemented; the idea is to decode
        // values from the resource object to the set of translations)
      }
    },

    updater () {
      const { moduleID, namespaceID } = this.module

      return translations => {
        return this.$ComposeAPI
          .moduleUpdateTranslations({ namespaceID, moduleID, translations })
          // re-fetch translations, sanitized and stripped
          .then(() => this.fetcher())
          .then((translations) => {
            // When translations are successfully saved,
            // scan changes and apply them back to the passed object
            // not the most elegant solution but is saves us from
            // handling the resource on multiple places
            moduleResTr(this.module, translations, this.currentLanguage, this.resource)

            return this.module
          })
          .then(module => {
            this.$emit('update:module', module)
          })
      }
    },
  },
}
</script>
