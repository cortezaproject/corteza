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
        resourceSingle: $t('general:label.user-group.single'),
        resourcePlural: $t('general:label.user-group.plural'),
      }"
      clickable
      sticky-header
      class="custom-resource-list-height flex-fill"
      @search="filterList"
      @row-clicked="handleRowClicked"
    >
      <template #header>
        <b-button
          variant="primary"
          size="lg"
          data-test-id="button-new-user-group"
          :to="{ name: 'system.userGroup.new' }"
        >
          {{ $t('new') }}
        </b-button>

        <c-permissions-button
          v-if="canGrant"
          resource="corteza::system:user-group/*"
          :button-label="$t('permissions')"
          size="lg"
        />
      </template>

      <template #toolbar>
        <c-resource-list-status-filter
          v-model="filter.deleted"
          data-test-id="filter-deleted-user-groups"
          :label="$t('filterForm.deleted.label')"
          :excluded-label="$t('filterForm.excluded.label')"
          :inclusive-label="$t('filterForm.inclusive.label')"
          :exclusive-label="$t('filterForm.exclusive.label')"
          @change="filterList"
        />

        <b-col />
      </template>

      <template #actions="{ item: u }">
        <b-dropdown
          v-if="(areActionsVisible({ resource: u, conditions: ['canDeleteUserGroup', 'canGrant'] }))"
          variant="outline-extra-light"
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

          <c-permissions-button
            v-if="canGrant"
            :title="u.name || u.handle || u.userGroupID"
            :target="u.name || u.handle || u.userGroupID"
            :resource="`corteza::system:user-group/${u.userGroupID}`"
            :button-label="$t('permissions')"
            class="dropdown-item"
          />

          <c-input-confirm
            v-if="u.canDeleteUserGroup"
            :text="getActionText(u)"
            show-icon
            :icon="getActionIcon(u)"
            borderless
            variant="link"
            size="md"
            button-class="dropdown-item"
            icon-class="text-danger"
            class="w-100"
            @confirmed="handleDelete(u)"
          />
        </b-dropdown>
      </template>
    </c-resource-list>
  </b-container>
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'
import listHelpers from 'corteza-webapp-admin/src/mixins/listHelpers'
import moment from 'moment'
import { mapGetters } from 'vuex'
const { CResourceList } = components

export default {
  name: 'UserGroupList',

  components: {
    CResourceList,
  },

  mixins: [
    listHelpers,
  ],

  i18nOptions: {
    namespaces: 'system.user-groups',
    keyPrefix: 'list',
  },

  data () {
    return {
      id: 'user-groups',

      primaryKey: 'userGroupID',
      editRoute: 'system.userGroup.edit',

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
          key: 'meta.short',
          sortable: false,
        },
        {
          key: 'handle',
          sortable: true,
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

    canGrant () {
      return this.can('system/', 'grant')
    },
  },

  methods: {
    items () {
      return this.procListResults(this.$SystemAPI.userGroupListCancellable(this.encodeListParams()))
    },

    rowClass (item) {
      return { 'text-secondary': item && !!item.deletedAt }
    },

    handleDelete (userGroup) {
      this.handleItemDelete({
        resource: userGroup,
        resourceName: 'userGroup',
      })
    },
  },
}
</script>
