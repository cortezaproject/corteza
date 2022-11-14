<template>
  <c-translator-button
    v-if="canManageResourceTranslations && resourceTranslationsEnabled"
    button-variant="light"
    v-bind="$props"
    :size="size"
    :title="$t('tooltip')"
    :resource="resource"
    :fetcher="fetcher"
    :updater="updater"
    class="ml-auto mr-1 py-1 px-3"
  />
</template>

<script>
import { compose } from '@cortezaproject/corteza-js'
import { mapGetters } from 'vuex'
import CTranslatorButton from 'corteza-webapp-compose/src/components/Translator/CTranslatorButton'

const keyPrefix = 'meta.bool'
const keySuffix = '.label'

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

    fetcher () {
      const { moduleID, namespaceID } = this.module

      return () => {
        return this.$ComposeAPI
          .moduleListTranslations({ namespaceID, moduleID })
          // Fields do not have their own translation endpoints,
          // we'll just filter what we need here.
          .then(set => {
            set = set
              // Extract translations for this field
              .filter(({ resource }) => this.resource === resource)
              // Ignore all option translations
              .filter(({ key }) => key.startsWith(keyPrefix) && key.endsWith(keySuffix))

            // after translations are fetched, make sure we copy all (updates) values from
            // the caller so translator editor can operate with recent values
            set
              .filter(({ lang }) => this.currentLanguage === lang)

            return set
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
            // @todo move this to ModuleFieldSelect classes
            // the logic there needs to be implemented; the idea is to encode
            // values from the set of translations back to the resource object
            const find = (key) => {
              return translations.find(t => t.key === key && t.lang === this.currentLanguage && t.resource === this.resource)
            }

            let tr
            tr = find('meta.bool.true.label')
            if (tr !== undefined) {
              this.field.options.trueLabel = tr.message
            }

            tr = find('meta.bool.false.label')
            if (tr !== undefined) {
              this.field.options.falseLabel = tr.message
            }
          })
      }
    },
  },
}
</script>
