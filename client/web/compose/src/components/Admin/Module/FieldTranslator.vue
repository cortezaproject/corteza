<template>
  <c-translator-button
    v-if="canManageResourceTranslations && resourceTranslationsEnabled"
    v-bind="$props"
    :tooltip="$t('tooltip')"
    :size="size"
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
import moduleFieldResTr from 'corteza-webapp-compose/src/lib/resource-translations/module-field'

export default {
  components: {
    CTranslatorButton,
  },

  i18nOptions: {
    namespaces: 'resource-translator',
    keyPrefix: 'resources.module.field',
  },

  props: {
    field: {
      type: compose.ModuleField,
      required: true,
    },

    module: {
      type: compose.Module,
      required: true,
    },

    size: {
      type: String,
      default: 'lg',
    },

    disabled: {
      type: Boolean,
      default: () => false,
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
      const { fieldID } = this.field
      const { moduleID, namespaceID } = this.module
      return `compose:module-field/${namespaceID}/${moduleID}/${fieldID}`
    },

    titles () {
      const { fieldID, name } = this.field
      const titles = {}

      titles[this.resource] = this.$t('title', { name: name || fieldID })

      return titles
    },

    fetcher () {
      const { moduleID, namespaceID } = this.module

      return () => {
        return this.$ComposeAPI
          .moduleListTranslations({ namespaceID, moduleID })
          // Fields do not have their own translation endpoints,
          // we'll just filter what we need here.
          .then(set => {
            return set
              // Extract translations for this field
              .filter(({ resource }) => this.resource === resource)
              // Ignore all option translations
              .filter(({ key }) => !key.startsWith('meta.options'))
              .filter(({ key }) => !key.startsWith('meta.bool'))

            // @todo pass set of translations to the object (ModuleField* class)
            // The logic there needs to be implemented; the idea is to decode
            // values from the resource object to the set of translations)
          })
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
            //
            // @todo move this to ModuleField* classes
            // the logic there needs to be implemented; the idea is to encode
            // values from the set of translations back to the resource object

            moduleFieldResTr(this.field, translations, this.currentLanguage, this.resource)
          })
      }
    },
  },
}
</script>
