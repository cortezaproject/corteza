<template>
  <div class="d-flex flex-column w-100 py-2 overflow-hidden h-100 py-2">
    <portal to="topbar-title">
      {{ $t('title') }}
    </portal>

    <portal to="topbar-tools">
      <b-btn
        v-if="canManage"
        data-test-id="button-manage-namespaces"
        variant="primary"
        size="sm"
        :to="{ name: 'namespace.manage' }"
      >
        {{ $t('manage-view.label') }}
        <font-awesome-icon
          :icon="['far', 'edit']"
          size="sm"
          class="ml-2"
        />
      </b-btn>
    </portal>

    <div class="d-flex flex-column justify-content-center align-items-center mx-4 mb-2">
      <div class="search w-100 mx-auto my-4">
        <c-input-search
          v-model.trim="query"
          :placeholder="$t('searchPlaceholder')"
          :debounce="200"
        />
      </div>
    </div>

    <div class="flex-fill overflow-auto">
      <b-container class="ns-wrapper h-100">
        <transition-group
          v-if="filtered && filtered.length"
          name="namespace-list"
          tag="b-row"
          class="d-flex flex-wrap align-items-stretch justify-content-center mx-2"
        >
          <b-col
            v-for="n in filtered"
            :key="n.namespaceID"
            cols="12"
            md="6"
            lg="4"
            xl="3"
            class="p-2"
          >
            <namespace-item :namespace="n" />
          </b-col>
        </transition-group>

        <div
          v-else
          class="d-flex justify-content-center align-items-center mt-5 w-100"
        >
          <h3 data-test-id="no-namespaces-found">
            {{ $t('noResults') }}
          </h3>
        </div>
      </b-container>
    </div>
  </div>
</template>
<script>
import { mapGetters } from 'vuex'
import NamespaceItem from 'corteza-webapp-compose/src/components/Namespaces/NamespaceItem'
import { components } from '@cortezaproject/corteza-vue'
const { CInputSearch } = components

export default {
  i18nOptions: {
    namespaces: 'namespace',
  },

  components: {
    NamespaceItem,
    CInputSearch,
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

    canManage () {
      if (this.can('compose/', 'namespace.create') || this.can('compose/', 'grant')) {
        return true
      }

      return this.namespaces.reduce((acc, ns) => {
        return acc || ns.canUpdateNamespace || ns.canDeleteNamespace
      }, false)
    },

    importNamespaceEndpoint () {
      return this.$ComposeAPI.namespaceImportEndpoint({})
    },

    filtered () {
      const query = this.query.toLowerCase()
      return this.namespaces
        .filter(({ enabled }) => enabled)
        .filter(({ slug, name }) => (slug + name).toLowerCase().indexOf(query) > -1)
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

    handleRowClicked ({ namespace }) {
      this.$router.push({
        name: 'namespace.edit',
        params: { namespaceID: namespace.namespaceID },
      })
    },
  },
}
</script>

<style lang="scss" scoped>
.search {
  max-width: 600px;
}
</style>
