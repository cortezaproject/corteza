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
        notFound: $t('admin:general.notFound'),
        noItems: $t('admin:general.resource-list.no-items'),
        loading: $t('loading'),
        showingPagination: 'admin:general.pagination.showing',
        singlePluralPagination: 'admin:general.pagination.single',
        prevPagination: $t('admin:general.pagination.prev'),
        nextPagination: $t('admin:general.pagination.next'),
        resourceSingle: $t('general:label.session.single'),
        resourcePlural: $t('general:label.session.plural')
      }"
      clickable
      sticky-header
      hide-search
      :hide-total="!pagination.incTotal"
      class="custom-resource-list-height flex-fill"
      @row-clicked="handleRowClicked"
    >
      <template #header>
        <b-form-group
          :label="$t('columns.sessionID')"
          label-class="text-primary"
          class="mb-0"
        >
          <c-input-search
            :value="filter.sessionID"
            size="sm"
            @input="filterBySessionID"
          />
        </b-form-group>

        <b-form-group
          :label="$t('columns.workflowID')"
          label-class="text-primary"
          class="mb-0"
        >
          <c-input-search
            :value="filter.workflowID"
            size="sm"
            @input="filterByWorkflowID"
          />
        </b-form-group>
      </template>

      <template #toolbar>
        <b-col>
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
        </b-col>
      </template>

      <template #sessionID="{ item }">
        <a
          href="javascript:;"
          @click="filterBySessionID(item.sessionID)"
        >
          {{ item.sessionID }}
        </a>
      </template>

      <template #workflowID="{ item }">
        <a
          href="javascript:;"
          @click="filterByWorkflowID(item.workflowID)"
        >
          {{ item.workflowID }}
        </a>
      </template>

      <template #actions="{ item }">
        <b-button
          size="sm"
          variant="link"
          :to="{ name: editRoute, params: { [primaryKey]: item[primaryKey] } }"
        >
          <font-awesome-icon
            :icon="['far', 'edit']"
          />
        </b-button>
      </template>
    </c-resource-list>
  </b-container>
</template>

<script>
import listHelpers from 'corteza-webapp-admin/src/mixins/listHelpers'
import { components } from '@cortezaproject/corteza-vue'
const { CResourceList, CInputSearch } = components

export default {
  components: {
    CResourceList,
    CInputSearch,
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
        // Use null not undefined!
        sessionID: null,
        workflowID: null,
        status: null,
        completed: 1,
      },

      sorting: {
        sortBy: 'createdAt',
        sortDesc: true,
      },

      pagination: {
        ...this.pagination,
        incTotal: false,
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
        { value: null, text: this.$t('filterForm.all.label') },
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

    filterBySessionID (sessionID) {
      this.filter.sessionID = sessionID || null
      this.filterList()
    },

    filterByWorkflowID (workflowID) {
      this.filter.workflowID = workflowID || null
      this.filterList()
    },
  },
}
</script>

<style scoped>
.content-header{
  margin-bottom: 0 !important;
}
</style>
