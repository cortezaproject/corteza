<template>
  <b-container
    fluid="xl"
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
        resourceSingle: $t('general:label.user.single'),
        resourcePlural: $t('general:label.user.plural'),
      }"
      clickable
      sticky-header
      class="custom-resource-list-height"
      @search="filterList"
      @row-clicked="handleRowClicked"
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

      <template #actions="{ item: u }">
        <b-dropdown
          v-if="(areActionsVisible({ resource: u, conditions: ['canDeleteUser', 'canGrant'] }))"
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
            v-if="canGrant"
            link-class="p-0"
            variant="light"
          >
            <c-permissions-button
              :title="u.name || u.handle || u.email || u.userID"
              :target="u.name || u.handle || u.email || u.userID"
              :resource="`corteza::system:user/${u.userID}`"
              button-variant="link dropdown-item text-decoration-none text-dark regular-font rounded-0"
            >
              <font-awesome-icon :icon="['fas', 'lock']" />
              {{ $t('permissions') }}
            </c-permissions-button>
          </b-dropdown-item>

          <c-input-confirm
            v-if="u.canDeleteUser"
            borderless
            variant="link"
            size="md"
            button-class="dropdown-item text-decoration-none text-dark regular-font rounded-0"
            class="w-100"
            @confirmed="handleDelete(u)"
          >
            <font-awesome-icon
              :icon="['far', 'trash-alt']"
              class="text-danger"
            />

            <span
              v-if="!u.deletedAt"
              class="p-1"
            >{{ $t('delete') }}</span>

            <span
              v-else
              class="p-1"
            >{{ $t('undelete') }}</span>
          </c-input-confirm>
        </b-dropdown>
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
          class: 'actions',
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

    handleDelete (user) {
      this.handleItemDelete({
        resource: user,
        resourceName: 'user',
      })
    },
  },
}
</script>
