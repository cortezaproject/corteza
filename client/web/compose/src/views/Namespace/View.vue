<template>
  <div
    :id="namespace ? namespace.slug || namespace.namespaceID : ''"
    class="d-flex w-100"
  >
    <router-view
      v-if="loaded && namespace"
      :namespace="namespace"
    />

    <div
      v-else
      class="loader flex-column align-items-center justify-content-center w-100 h-50"
    >
      <h1>
        {{ namespace ? (namespace.name || namespace.slug || namespace.namespaceID) : '...' }}
      </h1>

      <div class="d-flex align-items-center justify-content-center mt-4">
        <b-spinner />
        <h4 class="mb-0 ml-2">
          {{ $t('general:label.loading') }}
        </h4>
      </div>
    </div>

    <attachment-modal />
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import AttachmentModal from 'corteza-webapp-compose/src/components/Public/Page/Attachment/Modal'

export default {
  i18nOptions: {
    namespaces: 'general',
  },

  name: 'Namespace',

  components: {
    AttachmentModal,
  },

  props: {
    slug: {
      required: true,
      type: String,
    },
  },

  data () {
    return {
      loaded: false,

      error: '',
      namespace: null,
    }
  },

  computed: {
    ...mapGetters({
      namespacePending: 'namespace/pending',
      namespaces: 'namespace/set',
      modulePending: 'module/pending',
      chartPending: 'chart/pending',
      pagePending: 'page/pending',
      pageLayoutPending: 'pageLayout/pending',
      pages: 'page/set',
    }),

    parts () {
      return {
        namespace: this.namespacePending,
        module: this.modulePending,
        page: this.pagePending,
        chart: this.chartPending,
        pageLayout: this.pageLayoutPending,
      }
    },
  },

  watch: {
    slug: {
      immediate: true,
      handler (slug) {
        this.loaded = false

        let namespace = this.$store.getters['namespace/getByUrlPart'](slug)

        if (!namespace) {
          this.$store.dispatch('namespace/load', { force: true }).then(() => {
            namespace = this.$store.getters['namespace/getByUrlPart'](slug)
          }).catch(this.errHandler)
        }

        this.namespace = namespace
        this.prepareNamespace()
      },
    },
  },

  created () {
    this.error = ''
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
    prepareNamespace () {
      if (!this.namespace) {
        this.$router.push({ name: 'root' })
        return
      }

      if (!this.namespace.enabled) {
        this.toastDanger(this.$t('notification:namespace.disabled'))
        this.$router.push({ name: 'root' })
        return
      }

      const p = { namespace: this.namespace, namespaceID: this.namespace.namespaceID, clear: true }

      this.$store.dispatch('module/clearSet')
      this.$store.dispatch('chart/clearSet')
      this.$store.dispatch('page/clearSet')

      this.$root.$emit('check-namespace-sidebar', !this.namespace.meta.hideSidebar)

      // Preload all data we need.
      Promise.all([
        this.$store.dispatch('module/load', p)
          .catch(this.errHandler),

        this.$store.dispatch('chart/load', p)
          .catch(this.errHandler),

        this.$store.dispatch('page/load', p)
          .catch(this.errHandler),

        this.$store.dispatch('pageLayout/load', p)
          .catch(this.errHandler),

      ]).catch(this.errHandler).then(() => {
        setTimeout(() => {
          this.loaded = true
        }, 500)
      })
    },

    // Error handler for Promise
    errHandler (error) {
      switch ((error.response || {}).status) {
        case 403:
          this.error = this.$t('notification:general.composeAccessNotAllowed')
      }

      return Promise.reject(error)
    },

    setDefaultValues () {
      this.loaded = false
      this.error = ''
      this.namespace = null
    },
  },
}
</script>
<style lang="scss" scoped>
.error {
  font-size: 24px;
  background-color: var(--white);
  width: 100vw;
  height: 20vh;
  padding: 60px;
  top: 40vh;
}
</style>
