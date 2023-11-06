<template>
  <b-card
    class="shadow-sm"
    body-class="p-0"
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <template
      #header
    >
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>

    <c-resource-list
      :primary-key="primaryKey"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :fields="fields"
      :items="items"
      :row-class="genericRowClass"
      :translations="{
        searchPlaceholder: $t('filter-form.query.placeholder'),
        notFound: $t('admin:general.notFound'),
        noItems: $t('admin:general.resource-list.no-items'),
        loading: $t('loading'),
        showingPagination: 'admin:general.pagination.showing',
        singlePluralPagination: 'admin:general.pagination.single',
        prevPagination: $t('admin:general.pagination.prev'),
        nextPagination: $t('admin:general.pagination.next'),
        resourceSingle: $t('general:label.connection.single'),
        resourcePlural: $t('general:label.connection.plural')
      }"
      clickable
      hide-search
      card-header-class="rounded-0"
      class="h-100 bg-transparent"
      @row-clicked="handleRowClicked"
    >
      <template #header>
        <b-button
          variant="primary"
          size="lg"
          :to="{ name: 'system.connection.new' }"
        >
          {{ $t('add-button') }}
        </b-button>
      </template>

      <template #toolbar>
        <c-resource-list-status-filter
          v-model="filter.deleted"
          data-test-id="filter-deleted-connections"
          :label="$t('filterForm.deleted.label')"
          :excluded-label="$t('filterForm.excluded.label')"
          :inclusive-label="$t('filterForm.inclusive.label')"
          :exclusive-label="$t('filterForm.exclusive.label')"
          class="mt-2"
          @change="filterList"
        />
      </template>

      <template #actions="{ item: c }">
        <b-dropdown
          v-if="c.canDeleteConnection"
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

          <c-input-confirm
            :text="getActionText(c)"
            show-icon
            :icon="getActionIcon(c)"
            borderless
            variant="link"
            size="md"
            button-class="dropdown-item text-decoration-none text-dark regular-font rounded-0"
            icon-class="text-danger"
            class="w-100"
            @confirmed="handleDelete(c)"
          />
        </b-dropdown>
      </template>
    </c-resource-list>
  </b-card>
</template>

<script>
import listHelpers from 'corteza-webapp-admin/src/mixins/listHelpers'
import moment from 'moment'
import { components } from '@cortezaproject/corteza-vue'
const { CResourceList } = components

export default {
  components: {
    CResourceList,
  },

  mixins: [
    listHelpers,
  ],

  i18nOptions: {
    namespaces: 'system.connections',
    keyPrefix: 'list',
  },

  data () {
    return {
      id: 'connections',

      primaryKey: 'connectionID',
      editRoute: 'system.connection.edit',

      filter: {
        type: 'corteza::system:dal-connection',
        query: '',
        deleted: 0,
      },

      sorting: {
        sortBy: 'createdAt',
        sortDesc: true,
      },

      fields: [
        {
          key: 'name',
          sortable: false,
          formatter: (value, col, conn) => conn.meta.name || conn.handle,
        },
        {
          key: 'location',
          sortable: false,
          formatter: (value, col, conn) => conn.meta.location.properties.name,
        },
        {
          key: 'ownership',
          sortable: false,
          formatter: (value, col, conn) => conn.meta.ownership,
        },
        {
          key: 'createdAt',
          sortable: false,
          formatter: (v) => moment(v).fromNow(),
        },
        {
          key: 'actions',
          class: 'actions',
        },
      ].map(c => ({
        // Generate column label translation key
        label: c.label || this.$t(`columns.${c.key}`),
        ...c,
      })),
    }
  },

  methods: {

    items () {
      return this.procListResults(this.$SystemAPI.dalConnectionList(this.encodeListParams()))
    },

    handleDelete (connection) {
      this.handleItemDelete({
        resource: connection,
        resourceName: 'dalConnection',
        locale: 'connection',
      })
    },
  },
}
</script>
