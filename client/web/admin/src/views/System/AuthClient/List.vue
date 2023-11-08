<template>
  <b-container
    fluid="xl"
    class="d-flex flex-column flex-fill pt-2 pb-3"
  >
    <c-content-header :title="$t('title')" />

    <c-resource-list
      :primary-key="primaryKey"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :fields="fields"
      :items="items"
      :row-class="genericRowClass"
      :translations="{
        searchPlaceholder: $t('filterForm.query.placeholder'),
        notFound: $t('admin:general.notFound'),
        noItems: $t('admin:general.resource-list.no-items'),
        loading: $t('loading'),
        showingPagination: 'admin:general.pagination.showing',
        singlePluralPagination: 'admin:general.pagination.single',
        prevPagination: $t('admin:general.pagination.prev'),
        nextPagination: $t('admin:general.pagination.next'),
        resourceSingle: $t('general:label.auth_client.single'),
        resourcePlural: $t('general:label.auth_client.plural'),
      }"
      clickable
      sticky-header
      hide-search
      class="custom-resource-list-height flex-fill"
      @search="filterList"
      @row-clicked="handleRowClicked"
    >
      <template #header>
        <b-button
          v-if="canCreate"
          data-test-id="button-new-auth-client"
          variant="primary"
          size="lg"
          :to="{ name: 'system.authClient.new' }"
        >
          {{ $t('new') }}
        </b-button>

        <c-permissions-button
          v-if="canGrant"
          resource="corteza::system:auth-client/*"
          :button-label="$t('permissions')"
          size="lg"
        />
      </template>

      <template #toolbar>
        <c-resource-list-status-filter
          v-model="filter.deleted"
          data-test-id="filter-deleted-auth-clients"
          :label="$t('filterForm.deleted.label')"
          :excluded-label="$t('filterForm.excluded.label')"
          :inclusive-label="$t('filterForm.inclusive.label')"
          :exclusive-label="$t('filterForm.exclusive.label')"
          @change="filterList"
        />
      </template>

      <template #actions="{ item: a }">
        <b-dropdown
          v-if="(areActionsVisible({ resource: a, conditions: ['canDeleteAuthClient', 'canGrant'] }) && a.authClientID)"
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
            v-if="a.authClientID && canGrant"
            link-class="p-0"
          >
            <c-permissions-button
              :title="a.meta.name || a.handle || a.authClientID"
              :target="a.meta.name || a.handle || a.authClientID"
              :resource="`corteza::system:auth-client/${a.authClientID}`"
              button-variant="link dropdown-item text-decoration-none text-dark regular-font rounded-0"
              class="text-dark d-print-none border-0"
            >
              <font-awesome-icon :icon="['fas', 'lock']" />
              {{ $t('permissions') }}
            </c-permissions-button>
          </b-dropdown-item>

          <b-dropdown-item
            v-if="!a.isDefault && a.canDeleteAuthClient"
            link-class="p-0"
          >
            <c-input-confirm
              :text="getActionText(a)"
              show-icon
              :icon="getActionIcon(a)"
              borderless
              variant="link"
              size="md"
              button-class="dropdown-item text-decoration-none text-dark regular-font rounded-0"
              icon-class="text-danger"
              class="w-100"
              @confirmed="handleDelete(a)"
            />
          </b-dropdown-item>
        </b-dropdown>
      </template>
    </c-resource-list>
  </b-container>
</template>

<script>
import * as moment from 'moment'
import listHelpers from 'corteza-webapp-admin/src/mixins/listHelpers'
import { mapGetters } from 'vuex'
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
    namespaces: 'system.authclients',
    keyPrefix: 'list',
  },

  data () {
    return {
      id: 'authclient',

      primaryKey: 'authClientID',
      editRoute: 'system.authClient.edit',

      filter: {
        query: '',
        deleted: 0,
      },

      sorting: {
        sortBy: 'createdAt',
        sortDesc: true,
      },

      fields: [
        {
          key: 'meta.name',
          sortable: false,
        },
        {
          key: 'handle',
          sortable: true,
        },
        {
          key: 'enabled',
          formatter: (v) => v ? 'Yes' : 'No',
        },
        {
          key: 'createdAt',
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

    canCreate () {
      return this.can('system/', 'auth-client.create')
    },

    canGrant () {
      return this.can('system/', 'grant')
    },
  },

  methods: {
    items () {
      return this.procListResults(this.$SystemAPI.authClientList(this.encodeListParams()))
    },

    handleDelete (authclient) {
      this.handleItemDelete({
        resource: { clientID: authclient.authClientID },
        resourceName: 'authClient',
        locale: 'authclient',
      })
    },
  },
}
</script>
