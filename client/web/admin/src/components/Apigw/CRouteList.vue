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
      <h3 class="mb-0">
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
        searchPlaceholder: $t('filterForm.query.placeholder'),
        notFound: $t('admin:general.notFound'),
        noItems: $t('admin:general.resource-list.no-items'),
        loading: $t('loading'),
        showingPagination: 'admin:general.pagination.showing',
        singlePluralPagination: 'admin:general.pagination.single',
        prevPagination: $t('admin:general.pagination.prev'),
        nextPagination: $t('admin:general.pagination.next'),
        resourceSingle: $t('general:label.route.single'),
        resourcePlural: $t('general:label.route.plural')

      }"
      clickable
      card-header-class="rounded-0"
      class="h-100 bg-transparent"
      @search="filterList"
      @row-clicked="handleRowClicked"
    >
      <template #header>
        <b-button
          v-if="canCreate"
          data-test-id="button-add"
          variant="primary"
          size="lg"
          :to="{ name: 'system.apigw.new' }"
        >
          {{ $t('new') }}
        </b-button>

        <b-button
          v-if="$Settings.get('apigw.profiler.enabled', false)"
          data-test-id="button-profiler"
          variant="info"
          size="lg"
          :to="{ name: 'system.apigw.profiler' }"
        >
          {{ $t('profiler') }}
        </b-button>

        <c-permissions-button
          v-if="canGrant"
          data-test-id="button-permissions"
          resource="corteza::system:apigw-route/*"
          :button-label="$t('permissions')"
          size="lg"
        />
      </template>

      <template #toolbar>
        <c-resource-list-status-filter
          v-model="filter.deleted"
          data-test-id="filter-deleted-routes"
          :label="$t('filterForm.deleted.label')"
          :excluded-label="$t('filterForm.excluded.label')"
          :inclusive-label="$t('filterForm.inclusive.label')"
          :exclusive-label="$t('filterForm.exclusive.label')"
          @change="filterList"
        />
      </template>

      <template #actions="{ item: r }">
        <b-dropdown
          v-if="(areActionsVisible({ resource: r, conditions: ['canDeleteApigwRoute', 'canGrant'] }))"
          boundary="viewport"
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
            v-if="r.routeID && canGrant"
            link-class="p-0"
          >
            <c-permissions-button
              :title="r.endpoint || r.routeID"
              :target="r.endpoint || r.routeID"
              :resource="`corteza::system:apigw-route/${r.routeID}`"
              button-variant="link dropdown-item text-decoration-none text-dark regular-font rounded-0"
            >
              <font-awesome-icon :icon="['fas', 'lock']" />

              {{ $t('permissions') }}
            </c-permissions-button>
          </b-dropdown-item>

          <c-input-confirm
            v-if="r.canDeleteApigwRoute"
            :text="getActionText(r)"
            show-icon
            :icon="getActionIcon(r)"
            borderless
            variant="link"
            size="md"
            button-class="dropdown-item text-decoration-none text-dark regular-font rounded-0"
            icon-class="text-danger"
            class="w-100"
            @confirmed="handleDelete(r)"
          />
        </b-dropdown>
      </template>
    </c-resource-list>
  </b-card>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
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
    namespaces: 'system.apigw',
    keyPrefix: 'list',
  },

  data () {
    return {
      primaryKey: 'routeID',
      editRoute: 'system.apigw.edit',

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
          key: 'endpoint',
          sortable: true,
        },
        {
          key: 'method',
          sortable: false,
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
      return this.can('system/', 'apigw-route.create')
    },

    canGrant () {
      return this.can('system/', 'grant')
    },
  },

  methods: {
    ...mapActions({
      incLoader: 'ui/incLoader',
      decLoader: 'ui/decLoader',
    }),

    items () {
      return this.procListResults(this.$SystemAPI.apigwRouteList(this.encodeListParams()))
    },

    handleDelete (route) {
      this.handleItemDelete({
        resource: route,
        resourceName: 'apigwRoute',
        locale: 'gateway',
      })
    },
  },
}
</script>

<style lang="scss">
.route-list {
  .card-header {
    border-radius: 0;
  }
}
</style>
