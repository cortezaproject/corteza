<template>
  <div
    class="d-flex w-100 overflow-auto"
  >
    <portal to="topbar-title">
      {{ $t('title') }}
    </portal>

    <b-container
      class="ns-wrapper"
      fluid="xl"
    >
      <b-row
        class="justify-content-between wrap-with-vertical-gutters mt-1 mb-4"
        no-gutters
      >
        <div
          class="flex-grow-1"
        >
          <div
            class="wrap-with-vertical-gutters"
          >
            <b-btn
              data-test-id="button-create"
              :to="{ name: 'namespace.create' }"
              variant="primary"
              size="lg"
              class="mr-1 float-left"
            >
              {{ $t('create') }}
            </b-btn>

            <importer-modal
              class="mr-1 float-left"
              @imported="onImported"
              @failed="onFailed"
            />

            <c-permissions-button
              v-if="canGrant"
              resource="corteza::compose:namespace/*"
              button-variant="light"
              :button-label="$t('label.permissions')"
              class="btn-lg float-left"
            />
          </div>
        </div>
        <div class="flex-grow-1">
          <b-input-group
            class="h-100"
          >
            <b-form-input
              v-model.trim="query"
              data-test-id="input-search"
              class="h-100 mw-100 text-truncate"
              type="search"
              :placeholder="$t('searchPlaceholder')"
            />
            <b-input-group-append>
              <b-input-group-text class="text-primary bg-white">
                <font-awesome-icon
                  :icon="['fas', 'search']"
                />
              </b-input-group-text>
            </b-input-group-append>
          </b-input-group>
        </div>
      </b-row>

      <transition-group
        v-if="namespacesFiltered && namespacesFiltered.length"
        name="namespace-list"
        tag="div"
        class="card-columns my-2"
      >
        <namespace-item
          v-for="n in namespacesFiltered"
          :key="n.namespaceID"
          :namespace="n"
        />
      </transition-group>

      <div
        v-else
        class="d-flex justify-content-center align-items-center h-50 w-100"
      >
        <h3
          class="text-left"
        >
          {{ $t('noResults') }}
        </h3>
      </div>
    </b-container>
  </div>
</template>
<script>
import { mapGetters } from 'vuex'
import NamespaceItem from 'corteza-webapp-compose/src/components/Namespaces/NamespaceItem'
import ImporterModal from 'corteza-webapp-compose/src/components/Namespaces/Importer'

export default {
  i18nOptions: {
    namespaces: 'namespace',
  },

  components: {
    NamespaceItem,
    ImporterModal,
  },

  data () {
    return {
      query: '',
    }
  },

  computed: {
    ...mapGetters({
      namespaces: 'namespace/set',
      can: 'rbac/can',
    }),

    canGrant () {
      return this.can('compose/', 'grant')
    },

    canCreateNamespace () {
      return this.can('compose/', 'namespace.create')
    },

    canImportNamespace () {
      // If a user is allowed to create a namespace, they are considered to be allowed
      // to create any underlying resource when it comes to importing.
      //
      // This was agreed upon internally and may change in the future.

      return this.can('compose/', 'namespace.create')
    },

    importNamespaceEndpoint () {
      return this.$ComposeAPI.namespaceImportEndpoint({})
    },

    namespacesFiltered () {
      return this.namespaces.filter(ns => {
        const query = this.query.toLowerCase()
        return (ns.slug).toLowerCase().indexOf(query) > -1 || (ns.name).toLowerCase().indexOf(query) > -1
      })
    },
  },

  mounted () {
    document.title = this.$t('general:label.app-name.public')
  },

  methods: {
    onImported () {
      this.$store.dispatch('namespace/load', { force: true })
        .then(() => this.toastSuccess(this.$t('notification:namespace.imported')))
        .catch(this.toastErrorHandler(this.$t('notification:namespace.importFailed')))
    },

    onFailed (err) {
      this.toastErrorHandler(this.$t('notification:namespace.importFailed'))(err)
    },
  },
}
</script>
