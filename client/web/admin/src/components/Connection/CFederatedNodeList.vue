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
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>

    <c-resource-list
      primary-key="federationID"
      :loading-text="$t('loading')"
      :pagination="pagination"
      :sorting="sorting"
      :items="items"
      :fields="fields"
    >
      <template #actions>
        <b-button
          variant="primary"
          size="lg"
          :to="{ name: 'system.connections.new' }"
        >
          {{ $t('add-button') }}
        </b-button>
      </template>

      <template #filter>
        <b-input-group
          class="h-100"
        >
          <b-form-input
            v-model.trim="filter.query"
            :placeholder="$t('filterForm.query.placeholder')"
            class="text-truncate border-right-0 h-100"
            @search="filterList"
          />
          <b-input-group-append>
            <b-input-group-text class="text-primary bg-white">
              <font-awesome-icon
                :icon="['fas', 'search']"
              />
            </b-input-group-text>
          </b-input-group-append>
        </b-input-group>
      </template>
    </c-resource-list>
  </b-card>
</template>

<script>
import { fmt } from '@cortezaproject/corteza-js'
import listHelpers from 'corteza-webapp-admin/src/mixins/listHelpers'

export default {
  mixins: [
    listHelpers,
  ],

  i18nOptions: {
    namespaces: 'system.connections',
    keyPrefix: 'federation',
  },

  data () {
    return {
      id: 'federation',

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
          key: 'url',
          sortable: true,
          tdClass: 'text-info',
        },
        {
          key: 'location',
          sortable: true,
        },
        {
          key: 'ownership',
          sortable: true,
        },
        {
          key: 'createdBy',
          sortable: true,
        },
        {
          key: 'createdAt',
          sortable: true,
          formatter: v => v ? fmt.fullDateTime(v) : v,
        },
        {
          key: 'actions',
          class: 'text-right',
        },
      ].map(c => ({
        ...c,
        // Generate column label translation key
        label: c.label || this.$t(`columns.${c.key}`),
      })),
    }
  },

  methods: {
    items () {
      const set = [
        { federationID: '1', name: 'ACME France', url: 'https://corteza.acme.fr', location: 'France', ownership: 'ACME SARL', createdBy: 'John Doe', createdAt: new Date() },
      ]

      const filter = {
        count: set.length,
        limit: 10,
      }

      return this.procListResults(new Promise(resolve => setTimeout(resolve({ filter, set })), 200), false)
    },
  },
}
</script>
