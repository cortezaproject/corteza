<template>
  <b-modal
    v-model="showModal"
    :title="discoveryModalTitle"
    :ok-title="$t('general.label.saveAndClose')"
    ok-only
    ok-variant="primary"
    scrollable
    size="lg"
    body-class="p-0 border-top-0"
    header-class="p-3 pb-0 border-bottom-0"
    @ok="onSave()"
  >
    <b-tabs
      v-if="modal"
      v-model="currentTabIndex"
      active-nav-item-class="bg-grey"
      nav-wrapper-class="bg-white border-bottom"
      card
    >
      <b-tab
        v-for="scope in scopeOptions"
        :key="scope.value"
        :title="scope.title"
        class="mh-tab"
      >
        <field-picker
          :module="module"
          :fields.sync="currentFields"
          disable-system-fields
          style="max-height: 70vh;"
        />
      </b-tab>
    </b-tabs>
  </b-modal>
</template>

<script>

import { mapGetters } from 'vuex'
import { compose } from '@cortezaproject/corteza-js'
import FieldPicker from 'corteza-webapp-compose/src/components/Common/FieldPicker'

export default {
  i18nOptions: {
    namespaces: 'module',
  },

  components: {
    FieldPicker,
  },

  props: {
    modal: {
      type: Boolean,
      required: false,
    },

    module: {
      type: compose.Module,
      required: true,
    },
  },

  data () {
    return {
      public: {},
      private: {},
      protected: {},

      currentTabIndex: 0,
      currentLang: undefined,
    }
  },

  computed: {
    ...mapGetters({
      languages: 'languages/set',
    }),

    showModal: {
      get () {
        return this.modal
      },

      set (showModal) {
        this.$emit('update:modal', showModal)
      },
    },

    currentFields: {
      get () {
        if (this[this.currentScope] && this[this.currentScope].result[this.currentLanguageIndex]) {
          return this[this.currentScope].result[this.currentLanguageIndex].fields
        }
        return []
      },

      set (currentFields) {
        if (this.currentScope && this[this.currentScope].result[this.currentLanguageIndex]) {
          this[this.currentScope].result[this.currentLanguageIndex].fields = [...currentFields]
        }
      },
    },

    discoveryModalTitle () {
      const { handle } = this.module
      return handle ? `${this.$t('edit.discoverySettings.title')} (${handle})` : this.$t('edit.discoverySettings.title')
    },

    currentScope () {
      return this.currentTabIndex >= 0 ? (this.scopeOptions[this.currentTabIndex] || {}).value : undefined
    },

    currentLanguageIndex () {
      return this.currentScope ? this[this.currentScope].result.findIndex(({ lang }) => lang === this.currentLang) : -1
    },

    moduleFields () {
      return new Set(this.module.fields.map(({ name }) => name))
    },

    scopeOptions () {
      return [
        // { value: 'public', title: this.$t('edit.discoverySettings.public') },
        { value: 'private', title: this.$t('edit.discoverySettings.private') },
        // { value: 'protected', title: this.$t('edit.discoverySettings.protected') },
      ]
    },
  },

  watch: {
    modal: {
      immediate: true,
      handler (modal) {
        if (modal && this.module) {
          this.currentLang = this.defaultTranslationLanguage

          this.scopeOptions.forEach(({ value }) => {
            this[value] = {
              result: [],
            }

            this.languages.forEach(({ tag: lang }) => {
              let existingFields = new Set()

              if (this.module.config.discovery && this.module.config.discovery[value]) {
                const indexOfLanguage = this.module.config.discovery[value].result.findIndex(r => r.lang === lang)
                if (indexOfLanguage >= 0) {
                  existingFields = new Set(this.module.config.discovery[value].result[indexOfLanguage].fields.filter(name => this.moduleFields.has(name)))
                }
              }

              const fields = [...existingFields].map(name => this.module.fields.find(field => field.name === name))

              this[value].result.push({ lang, fields })
            })
          })
        }
      },
    },
  },

  methods: {
    onSave () {
      const discovery = {
        public: {},
        private: {},
        protected: {},
      }

      this.scopeOptions.forEach(({ value }) => {
        discovery[value].result = this[value].result.map(({ lang, fields }) => {
          return {
            lang,
            fields: fields.map(({ name }) => name),
          }
        })
      })

      this.$emit('save', {
        ...this.module.meta,
        discovery,
      })
    },
  },
}
</script>

<style scoped>
.mh-tab {
  max-height: calc(100vh - 16rem);
}
</style>
