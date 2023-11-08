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
      class="h-100 flex-fill"
      @search="filterList"
      @row-clicked="handleRowClicked"
    >
      <template #header>
        <b-btn
          v-if="canCreate"
          data-test-id="button-create"
          :to="{ name: 'namespace.create' }"
          variant="primary"
          size="lg"
        >
          {{ $t('toolbar.buttons.create') }}
        </b-btn>

        <importer-modal
          v-if="canImport"
          @imported="onImported"
          @failed="onFailed"
        />

        <c-permissions-button
          v-if="canGrant"
          resource="corteza::compose:namespace/*"
          :button-label="$t('toolbar.buttons.permissions')"
          size="lg"
        />
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
        <div>
          <b-dropdown
            v-if="n.canDeleteNamespace || n.canGrant"
            variant="outline-light"
            toggle-class="d-flex align-items-center justify-content-center text-primary border-0 py-2"
            no-caret
            dropleft
            lazy
            menu-class="m-0"
          >
            <template #button-content>
              <font-awesome-icon
                :icon="['fas', 'ellipsis-v']"
              />
            </template>

            <b-dropdown-item
              v-if="n.canGrant"
              link-class="p-0"
              variant="light"
            >
              <c-permissions-button
                :title="n.name || n.slug || n.namespaceID"
                :target="n.name || n.slug || n.namespaceID"
                :resource="`corteza::compose:namespace/${n.namespaceID}`"
                :tooltip="$t('permissions:resources.compose.namespace.tooltip')"
                :button-label="$t('permissions:ui.label')"
                button-variant="link dropdown-item text-decoration-none text-dark regular-font rounded-0"
              />
            </b-dropdown-item>

            <c-input-confirm
              v-if="n.canDeleteNamespace"
              :text="$t('delete')"
              show-icon
              borderless
              variant="link"
              size="md"
              button-class="dropdown-item text-decoration-none text-dark regular-font rounded-0"
              icon-class="text-danger"
              class="w-100"
              @confirmed="handleDelete(n)"
            />
          </b-dropdown>
        </div>
      </template>
    </c-resource-list>
  </b-container>
</template>
<script>
import { mapGetters, mapActions } from 'vuex'
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
      application: undefined,
      isApplication: false,

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
          tdClass: 'text-right text-nowrap actions',
        },
      ]
    },
  },

  methods: {
    ...mapActions({
      load: 'namespace/load',
      deleteNamespace: 'namespace/delete',
    }),

    onImported () {
      this.load({ force: true })
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

    fetchApplication (namespace) {
      const { namespaceID, slug } = namespace
      return this.$SystemAPI.applicationList({ name: slug || namespaceID })
        .then(({ set = [] }) => {
          if (set.length) {
            this.application = set[0]
            this.isApplication = this.application.enabled
          }
        })
        .catch(this.toastErrorHandler(this.$t('notification:namespace.deleteFailed')))
    },

    async handleDelete (namespace) {
      this.fetchApplication(namespace).then(() => {
        const { namespaceID } = namespace
        const { applicationID } = this.application || {}
        this.deleteNamespace({ namespaceID })
          .catch(this.toastErrorHandler(this.$t('notification:namespace.deleteFailed')))
          .then(() => {
            if (applicationID) {
              return this.$SystemAPI.applicationDelete({ applicationID })
            }
          })
          .then(() => {
            this.toastSuccess(this.$t('notification:namespace.deleted'))
            this.filterList()
          })
      })
    },
  },
}
</script>
