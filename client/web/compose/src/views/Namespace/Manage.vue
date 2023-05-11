<template>
  <b-container
    fluid="xl"
    class="d-flex flex-column py-3"
  >
    <portal to="topbar-title">
      {{ $t('title') }}
    </portal>

    <portal to="topbar-tools">
      <b-btn
        data-test-id="button-namespace-list"
        variant="primary"
        size="sm"
        class="mr-1 float-left"
        :to="{ name: 'namespace.list' }"
      >
        {{ $t('list-view') }}
        <font-awesome-icon
          :icon="['fas', 'columns']"
          class="ml-2"
        />
      </b-btn>
    </portal>

    <c-resource-list
      data-test-id="table-namespaces-list"
      :primary-key="primaryKey"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :fields="namespacesFields"
      :items="namespaceList"
      :translations="{
        searchPlaceholder: $t('namespace:searchPlaceholder'),
        notFound: $t('general:resourceList.notFound'),
        noItems: $t('general:resourceList.noItems'),
        loading: $t('general:label.loading'),
        showingPagination: 'general:resourceList.pagination.showing',
        singlePluralPagination: 'general:resourceList.pagination.single',
        prevPagination: $t('general:resourceList.pagination.prev'),
        nextPagination: $t('general:resourceList.pagination.next'),
        resourceSingle: $t('general:label.namespace.single'),
        resourcePlural: $t('general:label.namespace.plural'),
      }"
      clickable
      sticky-header
      class="h-100"
      @search="filterList"
      @row-clicked="handleRowClicked"
    >
      <template #header>
        <div
          class="wrap-with-vertical-gutters"
        >
          <b-btn
            v-if="canCreate"
            data-test-id="button-create"
            :to="{ name: 'namespace.create' }"
            variant="primary"
            size="lg"
            class="mr-1 float-left"
          >
            {{ $t('toolbar.buttons.create') }}
          </b-btn>

          <importer-modal
            v-if="canImport"
            class="mr-1 float-left"
            @imported="onImported"
            @failed="onFailed"
          />

          <c-permissions-button
            v-if="canGrant"
            resource="corteza::compose:namespace/*"
            button-variant="light"
            :button-label="$t('toolbar.buttons.permissions')"
            class="btn-lg float-left"
          />
        </div>
      </template>

      <template #enabled="{ item }">
        <font-awesome-icon
          :icon="['fas', item.enabled ? 'check' : 'times']"
        />
      </template>

      <template #changedAt="{ item }">
        {{ (item.deletedAt || item.updatedAt || item.createdAt) | locFullDateTime }}
      </template>

      <template #actions="{ item: n }">
        <b-button-group>
          <c-permissions-button
            v-if="n.canGrant"
            :title="n.name || n.slug || n.namespaceID"
            :target="n.name || n.slug || n.namespaceID"
            :resource="`corteza::compose:namespace/${n.namespaceID}`"
            :tooltip="$t('permissions:resources.compose.namespace.tooltip')"
            button-variant="outline-light"
            class="text-dark d-print-none border-0"
          />
        </b-button-group>
      </template>
    </c-resource-list>
  </b-container>
</template>
<script>
import { mapGetters } from 'vuex'
import ImporterModal from 'corteza-webapp-compose/src/components/Namespaces/Importer'
import listHelpers from 'corteza-webapp-compose/src/mixins/listHelpers'

export default {
  i18nOptions: {
    namespaces: 'namespace',
    keyPrefix: 'manage',
  },

  components: {
    ImporterModal,
  },

  mixins: [
    listHelpers,
  ],

  data () {
    return {
      primaryKey: 'namespaceID',

      filter: {
        query: '',
      },

      sorting: {
        sortBy: 'name',
        sortDesc: false,
      },
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

    canCreate () {
      return this.can('compose/', 'namespace.create')
    },

    canImport () {
      // If a user is allowed to create a namespace, they are considered to be allowed
      // to create any underlying resource when it comes to importing.
      //
      // This was agreed upon internally and may change in the future.

      return this.can('compose/', 'namespace.create')
    },

    importNamespaceEndpoint () {
      return this.$ComposeAPI.namespaceImportEndpoint({})
    },

    namespacesFields () {
      return [
        {
          key: 'name',
          sortable: true,
          label: this.$t('table.columns.name'),
        },
        {
          key: 'slug',
          sortable: true,
          label: this.$t('table.columns.slug'),
          class: 'text-nowrap',
        },
        {
          key: 'enabled',
          label: this.$t('table.columns.enabled'),
          class: 'text-center',
        },
        {
          key: 'changedAt',
          sortable: true,
          label: this.$t('table.columns.changedAt'),
          class: 'text-right text-nowrap',
        },
        {
          key: 'actions',
          label: '',
          tdClass: 'text-right text-nowrap',
        },
      ]
    },
  },

  methods: {
    onImported () {
      this.$store.dispatch('namespace/load', { force: true })
        .then(() => {
          this.filterList()
          this.toastSuccess(this.$t('notification:namespace.imported'))
        })
        .catch(this.toastErrorHandler(this.$t('notification:namespace.importFailed')))
    },

    onFailed (err) {
      this.toastErrorHandler(this.$t('notification:namespace.importFailed'))(err)
    },

    handleRowClicked ({ namespaceID }) {
      this.$router.push({
        name: 'namespace.edit',
        params: { namespaceID },
      })
    },

    namespaceList () {
      return this.procListResults(this.$ComposeAPI.namespaceList(this.encodeListParams()))
    },
  },
}
</script>
