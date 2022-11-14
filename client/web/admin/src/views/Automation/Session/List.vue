<template>
  <b-container
    class="py-3"
  >
    <c-content-header
      :title="$t('title')"
    />

    <c-resource-list
      primary-key="sessionID"
      edit-route="automation.session.edit"
      :loading-text="$t('loading')"
      :paging="paging"
      :sorting="sorting"
      :items="items"
      :fields="fields"
    >
      <template #filter>
        <b-row
          no-gutters
        >
          <c-resource-list-status-filter
            v-model="filter.completed"
            class="mb-2"
            :label="$t('filterForm.inProgress.label')"
            :excluded-label="$t('filterForm.excluded.label')"
            :inclusive-label="$t('filterForm.inclusive.label')"
            :exclusive-label="$t('filterForm.exclusive.label')"
            @change="filterList"
          />
        </b-row>

        <b-form-radio-group
          v-model="filter.status"
          :options="statusOptions"
          buttons
          button-variant="outline-primary"
          size="sm"
          name="radio-btn-outline"
          @change="filterList"
        />
        <span class="mt-1 ml-2 text-nowrap">
          {{ $t('filterForm.sessions.label') }}
        </span>
      </template>
    </c-resource-list>
  </b-container>
</template>

<script>
import listHelpers from 'corteza-webapp-admin/src/mixins/listHelpers'

export default {
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
  },
}
</script>
