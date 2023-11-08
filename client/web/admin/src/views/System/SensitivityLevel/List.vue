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
        notFound: $t('admin:general.notFound'),
        noItems: $t('admin:general.resource-list.no-items'),
        loading: $t('loading'),
        showingPagination: 'admin:general.pagination.showing',
        singlePluralPagination: 'admin:general.pagination.single',
        prevPagination: $t('admin:general.pagination.prev'),
        nextPagination: $t('admin:general.pagination.next'),
        resourceSingle: $t('general:label.sensitivity_level.single'),
        resourcePlural: $t('general:label.sensitivity_level.plural'),
      }"
      clickable
      sticky-header
      hide-search
      class="custom-resource-list-height flex-fill"
      @row-clicked="handleRowClicked"
    >
      <template #header>
        <b-button
          v-if="canCreate"
          data-test-id="button-new-sens-lvl"
          variant="primary"
          size="lg"
          :to="{ name: 'system.sensitivityLevel.new' }"
        >
          {{ $t('new') }}
        </b-button>
      </template>

      <template #toolbar>
        <c-resource-list-status-filter
          v-model="filter.deleted"
          :label="$t('filterForm.deleted.label')"
          :excluded-label="$t('filterForm.excluded.label')"
          :inclusive-label="$t('filterForm.inclusive.label')"
          :exclusive-label="$t('filterForm.exclusive.label')"
          @change="filterList"
        />
      </template>

      <template #actions="{ item: s }">
        <b-dropdown
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

          <c-input-confirm
            :text="getActionText(s)"
            show-icon
            :icon="getActionIcon(s)"
            borderless
            variant="link"
            size="md"
            button-class="dropdown-item text-decoration-none text-dark regular-font rounded-0"
            icon-class="text-danger"
            class="w-100"
            @confirmed="handleDelete(s)"
          />
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
    namespaces: 'system.sensitivityLevel',
    keyPrefix: 'list',
  },

  data () {
    return {
      id: 'sensitivityLevel',

      primaryKey: 'sensitivityLevelID',
      editRoute: 'system.sensitivityLevel.edit',

      filter: {
        query: '',
        deleted: 0,
      },

      sorting: {
        sortBy: 'level',
        sortDesc: true,
      },

      fields: [
        {
          key: 'meta.name',
        },
        {
          key: 'level',
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

    canCreate () {
      return this.can('system/', 'dal-sensitivity-level.manage')
    },
  },

  methods: {
    items () {
      return this.procListResults(this.$SystemAPI.dalSensitivityLevelList(this.encodeListParams()))
    },

    handleDelete (sensitivityLevel) {
      this.handleItemDelete({
        resource: sensitivityLevel,
        resourceName: 'dalSensitivityLevel',
        locale: 'sensitivityLevel',
      })
    },
  },
}
</script>
