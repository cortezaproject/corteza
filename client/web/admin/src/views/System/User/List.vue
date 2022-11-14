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
      primary-key="userID"
      edit-route="system.user.edit"
      :loading-text="$t('loading')"
      :paging="paging"
      :sorting="sorting"
      :items="items"
      :fields="fields"
    >
      <template #filter>
        <b-form-group
          class="p-0 m-0"
        >
          <b-input-group>
            <b-form-input
              v-model.trim="filter.query"
              :placeholder="$t('filterForm.query.placeholder')"
              @keyup="filterList"
            />
          </b-input-group>
        </b-form-group>
        <b-row
          no-gutters
          class="mt-3"
        >
          <c-resource-list-status-filter
            v-model="filter.deleted"
            class="col-12 col-lg-6 mb-1 mb-lg-0"
            :label="$t('filterForm.deleted.label')"
            :excluded-label="$t('filterForm.excluded.label')"
            :inclusive-label="$t('filterForm.inclusive.label')"
            :exclusive-label="$t('filterForm.exclusive.label')"
            @change="filterList"
          />
          <c-resource-list-status-filter
            v-model="filter.suspended"
            class="col-12 col-lg-6"
            :label="$t('filterForm.suspended.label')"
            :excluded-label="$t('filterForm.excluded.label')"
            :inclusive-label="$t('filterForm.inclusive.label')"
            :exclusive-label="$t('filterForm.exclusive.label')"
            @change="filterList"
          />
        </b-row>
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
import { url } from '@cortezaproject/corteza-vue'

export default {
  name: 'UserList',
  components: {
    CUserExportModal,
    CUserImportModal,
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

      filter: {
        query: '',
        suspended: 0,
        deleted: 0,
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
  },
}
</script>
