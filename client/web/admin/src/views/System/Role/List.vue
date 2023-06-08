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
          v-if="canCreate"
          data-test-id="button-new-role"
          variant="primary"
          class="mr-2"
          :to="{ name: 'system.role.new' }"
        >
          {{ $t('new') }}
        </b-button>

        <c-permissions-button
          v-if="canGrant"
          resource="corteza::system:role/*"
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
        ui-page="role/list"
        ui-slot="toolbar"
        resource-type="system"
        default-variant="link"
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
        resourceSingle: $t('general:label.role.single'),
        resourcePlural: $t('general:label.role.plural'),
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
          data-test-id="filter-deleted-roles"
          :label="$t('filterForm.deleted.label')"
          :excluded-label="$t('filterForm.excluded.label')"
          :inclusive-label="$t('filterForm.inclusive.label')"
          :exclusive-label="$t('filterForm.exclusive.label')"
          @change="filterList"
        />

        <c-resource-list-status-filter
          v-model="filter.archived"
          data-test-id="filter-archived-roles"
          :label="$t('filterForm.archived.label')"
          :excluded-label="$t('filterForm.excluded.label')"
          :inclusive-label="$t('filterForm.inclusive.label')"
          :exclusive-label="$t('filterForm.exclusive.label')"
          @change="filterList"
        />
      </template>

      <template #actions="{ item: r }">
        <b-dropdown
          v-if="(areActionsVisible({ resource: r, conditions: ['canDeleteRole', 'canGrant'] }) && r.roleID)"
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
            v-if="r.roleID && canGrant"
            link-class="p-0"
            variant="light"
          >
            <c-permissions-button
              :title="r.name || r.handle || r.roleID"
              :target="r.name || r.handle || r.roleID"
              :resource="`corteza::system:role/${r.roleID}`"
              button-variant="link dropdown-item text-decoration-none text-dark regular-font rounded-0"
            >
              <font-awesome-icon :icon="['fas', 'lock']" />
              {{ $t('permissions') }}
            </c-permissions-button>
          </b-dropdown-item>

          <c-input-confirm
            v-if="r.canDeleteRole"
            borderless
            variant="link"
            size="md"
            button-class="dropdown-item text-decoration-none text-dark regular-font rounded-0"
            class="w-100"
            @confirmed="handleDelete(r)"
          >
            <font-awesome-icon
              :icon="['far', 'trash-alt']"
              class="text-danger"
            />

            <span
              v-if="!r.deletedAt"
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
    namespaces: 'system.roles',
    keyPrefix: 'list',
  },

  data () {
    return {
      id: 'roles',

      primaryKey: 'roleID',
      editRoute: 'system.role.edit',

      filter: {
        query: '',
        archived: 0,
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

    canCreate () {
      return this.can('system/', 'role.create')
    },

    canGrant () {
      return this.can('system/', 'grant')
    },

    userID () {
      if (this.$auth.user) {
        return this.$auth.user.userID
      }
      return undefined
    },
  },

  methods: {
    items () {
      return this.procListResults(this.$SystemAPI.roleList(this.encodeListParams()))
    },

    rowClass (item) {
      return { 'text-secondary': item && (!!item.deletedAt || !!item.archivedAt) }
    },

    handleDelete (role) {
      this.handleItemDelete({
        resource: role,
        resourceName: 'role',
      })
    },
  },
}
</script>
