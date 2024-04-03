<template>
  <div>
    <b-modal
      v-model="showModal"
      size="xl"
      lazy
      scrollable
      :title="title"
      no-fade
      body-class="position-static p-0"
      @hide="onHide"
    >
      <c-translator-form
        v-if="loaded"
        :primary-resource="resource"
        :translations="translations"
        :key-prettifier="keyPrettyfier"
        :languages="languages"
        :titles="titles"
        :highlight-key="highlightKey"
        @change="changes=$event"
      />

      <template #modal-footer>
        <b-button
          data-test-id="button-submit"
          type="submit"
          variant="primary"
          :tabindex="languages.length + 1"
          :disabled="disabled || changes.length === 0"
          @click="onSubmit()"
        >
          {{ $t('save-changes') }}
        </b-button>
      </template>
    </b-modal>
  </div>
</template>
<script lang="js">
import CTranslatorForm from './CTranslatorForm.vue'
import { mapGetters, mapActions } from 'vuex'

export default {
  i18nOptions: {
    namespaces: 'resource-translator',
    keyPrefix: 'translator',
  },

  components: {
    CTranslatorForm,
  },

  data () {
    return {
      resource: undefined,
      updater: undefined,
      loaded: false,

      translations: [],
      changes: [],

      titles: {},
    }
  },

  computed: {
    ...mapGetters({
      languages: 'languages/set',
    }),

    title () {
      return this.titles[this.resource] || ''
    },

    showModal: {
      get () {
        return this.loaded
      },

      set (open) {
        if (!open) {
          this.clear()
        }
      },
    },

    disabled () {
      return false
    },
  },

  created () {
    this.loadLanguages()
  },

  mounted () {
    this.loaded = false
    this.$root.$on('c-translator', this.loadModal)
  },

  beforeDestroy () {
    this.destroyEvents()
  },

  methods: {
    ...mapActions({
      loadLanguages: 'languages/load',
    }),

    loadModal (payload) {
      if (!payload) {
        // when falsy payload is received,
        // close the translator modal
        this.clear()
        return
      }

      const { resource, titles, fetcher, updater, highlightKey, keyPrettyfier } = payload

      this.resource = resource
      this.titles = titles
      this.highlightKey = highlightKey
      this.updater = updater
      this.keyPrettyfier = keyPrettyfier
      this.changes = []

      fetcher().then(tt => {
        this.translations = tt
        this.loaded = true
      })
    },

    onSubmit () {
      if (this.changes.length === 0) {
        // no translations were changed or added
        // we can close the modal w/o saving
        this.clear()
        return
      }

      this.updater(this.changes)
        .then(() => {
          this.toastSuccess(this.$t('notification:translations.saved'))
          this.clear()
        })
        .catch(this.toastErrorHandler(this.$t('notification:translations.saveFailed')))
    },

    onHide () {
      this.clear()
    },

    clear () {
      this.titles = {}
      this.changes = []
      this.resource = undefined
      this.highlightKey = undefined
      this.updater = undefined
      this.keyPrettyfier = undefined
      this.loaded = false
    },

    destroyEvents () {
      this.$root.$off('c-translator', this.loadModal)
    },
  },
}
</script>
