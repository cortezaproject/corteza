<template>
  <b-container
    class="py-3"
  >
    <c-content-header
      :title="$t('title')"
    />

    <c-resource-list
      :primary-key="primaryKey"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :fields="fields"
      :items="items"
      :row-class="rowClass"
      :translations="{
        notFound: $t('admin:general.notFound'),
        noItems: $t('admin:general.resource-list.no-items'),
        loading: $t('loading'),
        showingPagination: 'admin:general.pagination.showing',
        singlePluralPagination: 'admin:general.pagination.single',
        prevPagination: $t('admin:general.pagination.prev'),
        nextPagination: $t('admin:general.pagination.next'),
      }"
      hide-search
    >
      <template #header>
        <c-resource-list-status-filter
          v-model="filter.completed"
          :label="$t('filterForm.inProgress.label')"
          :excluded-label="$t('filterForm.excluded.label')"
          :inclusive-label="$t('filterForm.inclusive.label')"
          :exclusive-label="$t('filterForm.exclusive.label')"
          @change="filterList"
        />

        <b-form-radio-group
          v-model="filter.status"
          :options="statusOptions"
          buttons
          button-variant="outline-primary"
          size="sm"
          @change="filterList"
        />
        <span class="ml-2 text-nowrap">
          {{ $t('filterForm.sessions.label') }}
        </span>
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
import listHelpers from 'corteza-webapp-admin/src/mixins/listHelpers'
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
    namespaces: 'automation.sessions',
    keyPrefix: 'list',
  },

  data () {
    return {
      id: 'session',

      primaryKey: 'sessionID',
      editRoute: 'automation.session.edit',

      filter: {
        status: undefined,
        completed: 1,
        sort: 'createdAt DESC',
      },

      sorting: {
        sortBy: 'createdAt',
        sortDesc: true,
      },

      fields: [
        {
          key: 'sessionID',
        },
        {
          key: 'workflowID',
        },
        {
          key: 'status',
          sortable: true,
        },
        {
          key: 'eventType',
          sortable: true,
        },
        {
          key: 'createdAt',
          sortable: true,
          formatter: (v) => new Date(v).toLocaleString('en-EN'),
        },
        {
          key: 'actions',
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
    statusOptions () {
      return [
        { value: undefined, text: this.$t('filterForm.all.label') },
        { value: 0, text: this.$t('filterForm.started.label') },
        { value: 1, text: this.$t('filterForm.prompted.label') },
        { value: 2, text: this.$t('filterForm.suspended.label') },
        { value: 3, text: this.$t('filterForm.failed.label') },
        { value: 4, text: this.$t('filterForm.completed.label') },
      ]
    },
  },

  methods: {
    items () {
      return this.procListResults(this.$AutomationAPI.sessionList(this.encodeListParams()))
    },

    rowClass (item) {
      return { 'text-primary': item && !!item.completedAt }
    },
  },
}
</script>
