<template>
  <div
    class="d-flex w-100"
  >
    <router-view
      v-if="loaded && namespace"
      :namespace="namespace"
    />

    <div
      v-else
      class="loader flex-column w-100 h-100"
    >
      <div>
        <div class="logo w-100" />

        <h1 class="text-center">
          {{ namespace ? (namespace.name || namespace.slug || namespace.namespaceID) : '...' }}
        </h1>

        <div>
          <div
            v-for="(pending, part) in parts"
            :key="part"
            class="p-1"
          >
            <div
              class="pending pr-3 d-inline-block text-right"
            >
              <b-spinner
                v-if="pending"
                small
              />
              <font-awesome-icon
                v-else
                :icon="['fas', 'check']"
              />
            </div>
            <div
              class="d-inline-block"
            >
              {{ $t('navigation.' + part) }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  i18nOptions: {
    namespaces: 'general',
  },

  name: 'Namespace',

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
      pages: 'page/set',
    }),

    parts () {
      return {
        namespace: this.namespacePending,
        module: this.modulePending,
        page: this.pagePending,
        chart: this.chartPending,
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

      // Preload all data we need.
      Promise.all([
        this.$store.dispatch('module/load', p)
          .catch(this.errHandler),

        this.$store.dispatch('chart/load', p)
          .catch(this.errHandler),

        this.$store.dispatch('page/load', p)
          .catch(this.errHandler),

      ]).catch(this.errHandler).then(() => {
        this.loaded = true
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
  },
}
</script>
<style lang="scss" scoped>
.error {
  font-size: 24px;
  background-color: $white;
  width: 100vw;
  height: 20vh;
  padding: 60px;
  top: 40vh;
}
</style>
