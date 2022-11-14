<template>
  <b-container
    class="py-3"
  >
    <c-content-header
      :title="$t('title')"
    >
      <span
        class="text-nowrap"
      >
        <b-button
          variant="primary"
          data-test-id="button-new-user"
          class="mr-2 float-left"
          :to="{ name: 'system.user.new' }"
        >
          {{ $t('new') }}
        </b-button>
        <c-user-import-modal
          class="mr-1 float-left"
          @imported="onImported"
        />
        <c-user-export-modal
          class="mr-1 float-left"
          @export="onExport"
        />
        <c-permissions-button
          v-if="canGrant"
          resource="corteza::system:user/*"
          button-variant="light"
        >
          <font-awesome-icon :icon="['fas', 'lock']" />
          {{ $t('permissions') }}
        </c-permissions-button>
      </span>
      <b-dropdown
        v-if="false"
        variant="link"
        right
        menu-class="shadow-sm"
        :text="$t('export')"
      >
        <b-dropdown-item-button variant="link">
          {{ $t('yaml') }}
        </b-dropdown-item-button>
      </b-dropdown>
      <c-corredor-manual-buttons
        ui-page="user/list"
        ui-slot="toolbar"
        resource-type="system"
        class="mr-1"
        @click="dispatchCortezaSystemEvent($event)"
      />
    </c-content-header>

    <c-resource-list
      :primary-key="primaryKey"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :fields="fields"
      :items="items"
      :row-class="rowClass"
      :translations="{
        searchPlaceholder: $t('filterForm.query.placeholder'),
        notFound: $t('admin:general.notFound'),
        noItems: $t('admin:general.resource-list.no-items'),
        loading: $t('loading'),
        showingPagination: 'admin:general.pagination.showing',
        singlePluralPagination: 'admin:general.pagination.single',
        prevPagination: $t('admin:general.pagination.prev'),
        nextPagination: $t('admin:general.pagination.next'),
      }"
      @search="filterList"
    >
      <template #header>
        <c-resource-list-status-filter
          v-model="filter.deleted"
          data-test-id="filter-deleted-users"
          :label="$t('filterForm.deleted.label')"
          :excluded-label="$t('filterForm.excluded.label')"
          :inclusive-label="$t('filterForm.inclusive.label')"
          :exclusive-label="$t('filterForm.exclusive.label')"
          @change="filterList"
        />
        <c-resource-list-status-filter
          v-model="filter.suspended"
          data-test-id="filter-suspended-users"
          :label="$t('filterForm.suspended.label')"
          :excluded-label="$t('filterForm.excluded.label')"
          :inclusive-label="$t('filterForm.inclusive.label')"
          :exclusive-label="$t('filterForm.exclusive.label')"
          @change="filterList"
        />
      </template>

      <template #actions="{ item }">
        <b-button
          size="sm"
          variant="link"
          :to="{ name: editRoute, params: { [primaryKey]: item[primaryKey] } }"
        >
          <font-awesome-icon
            :icon="['fas', 'pen']"
          />
        </b-button>
      </template>
    </c-resource-list>
  </b-container>
</template>

<script>
import { system } from '@cortezaproject/corteza-js'
import moment from 'moment'
import listHelpers from 'corteza-webapp-admin/src/mixins/listHelpers'
import CUserExportModal from 'corteza-webapp-admin/src/components/User/CUserExportModal'
import CUserImportModal from 'corteza-webapp-admin/src/components/User/CUserImportModal'
import { mapGetters } from 'vuex'
import { url, components } from '@cortezaproject/corteza-vue'
const { CResourceList } = components

export default {
  name: 'UserList',

  components: {
    CUserExportModal,
    CUserImportModal,
    CResourceList,
  },

  mixins: [
    listHelpers,
  ],

  i18nOptions: {
    namespaces: 'system.users',
    keyPrefix: 'list',
  },

  data () {
    return {
      id: 'users',

      primaryKey: 'userID',
      editRoute: 'system.user.edit',

      filter: {
        query: '',
        suspended: 0,
        deleted: 0,
      },

      sorting: {
        sortBy: 'createdAt',
        sortDesc: true,
      },

      fields: [
        {
          key: 'name',
          sortable: true,
        },
        {
          key: 'email',
          sortable: true,
        },
        {
          key: 'handle',
          sortable: true,
        },
        {
          key: 'createdAt',
          label: 'Created',
          sortable: true,
          formatter: (v) => moment(v).fromNow(),
        },
        {
          key: 'actions',
          label: '',
          tdClass: 'text-right',
        },
      ].map(c => ({
        ...c,
        // Generate column label translation key
        label: this.$t(`columns.${c.key}`),
      })),
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canGrant () {
      return this.can('system/', 'grant')
    },
  },

  methods: {
    makeEvent () {
      return system.SystemEvent()
    },

    onExport (e) {
      const params = {
        filename: 'export',
        ...e,
      }

      const exportUrl = url.Make({
        url: `${this.$SystemAPI.baseURL}${this.$SystemAPI.userExportEndpoint(params)}`,
        query: {
          jwt: this.$auth.accessToken,
          inclRoleMembership: e.inclRoleMembership || false,
          inclRoles: e.inclRoles || false,
        },
      })

      window.open(exportUrl)
    },

    onImported () {
      this.toastSuccess(this.$t('notification:user.import.success'))
      this.filterList()
    },

    items () {
      return this.procListResults(this.$SystemAPI.userList(this.encodeListParams()))
    },

    rowClass (item) {
      return { 'text-secondary': item && (!!item.deletedAt || !!item.suspendedAt) }
    },
  },
}
</script>
