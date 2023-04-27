<template>
  <div
    class="d-flex w-100 overflow-auto"
  >
    <portal to="topbar-title">
      {{ $t('title') }}
    </portal>

    <portal to="topbar-tools">
      <b-btn
        v-if="canManage"
        data-test-id="button-manage-namespaces"
        variant="primary"
        size="sm"
        class="mr-1 float-left"
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

    <b-container
      class="ns-wrapper"
      fluid="xl"
    >
      <b-row
        class="wrap-with-vertical-gutters my-3"
        no-gutters
      >
        <b-col
          offset-md="2"
          offset-lg="3"
          md="8"
          lg="6"
        >
          <c-input-search
            v-model.trim="query"
            :placeholder="$t('searchPlaceholder')"
          />
        </b-col>
      </b-row>

      <transition-group
        v-if="filtered && filtered.length"
        name="namespace-list"
        tag="div"
        class="row my-3 card-deck no-gutters"
      >
        <namespace-item
          v-for="n in filtered"
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
