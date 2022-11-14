<template>
  <div class="py-3">
    <portal to="topbar-title">
      {{ $t('navigation.chart') }}
    </portal>

    <b-container fluid="xl">
      <b-row no-gutters>
        <b-col>
          <c-resource-list
            :primary-key="primaryKey"
            :filter="filter"
            :sorting="sorting"
            :pagination="pagination"
            :fields="tableFields"
            :items="chartList"
            :translations="{
              searchPlaceholder: $t('chart.searchPlaceholder'),
              notFound: $t('general:resourceList.notFound'),
              noItems: $t('general:resourceList.noItems'),
              loading: $t('general:label.loading'),
              showingPagination: 'general:resourceList.pagination.showing',
              singlePluralPagination: 'general:resourceList.pagination.single',
              prevPagination: $t('general:resourceList.pagination.prev'),
              nextPagination: $t('general:resourceList.pagination.next'),
            }"
            clickable
            @search="filterList"
            @row-clicked="handleRowClicked"
          >
            <template #header>
              <div
                class="wrap-with-vertical-gutters"
              >
                <b-dropdown
                  v-if="namespace.canCreateChart"
                  variant="primary"
                  size="lg"
                  class="float-left mr-1"
                  :text="$t('chart.add')"
                >
                  <b-dropdown-item-button
                    variant="dark"
                    @click="$router.push({ name: 'admin.charts.create', params: { category: 'generic' } })"
                  >
                    {{ $t('chart.addGeneric') }}
                  </b-dropdown-item-button>
                  <b-dropdown-item-button
                    variant="dark"
                    @click="$router.push({ name: 'admin.charts.create', params: { category: 'funnel' } })"
                  >
                    {{ $t('chart.addFunnel') }}
                  </b-dropdown-item-button>
                  <b-dropdown-item-button
                    variant="dark"
                    @click="$router.push({ name: 'admin.charts.create', params: { category: 'gauge' } })"
                  >
                    {{ $t('chart.addGauge') }}
                  </b-dropdown-item-button>
                </b-dropdown>

                <import
                  v-if="namespace.canCreateChart"
                  :namespace="namespace"
                  type="chart"
                  class="float-left mr-1"
                />

                <export
                  :list="charts"
                  type="chart"
                  class="float-left mr-1"
                />
                <c-permissions-button
                  v-if="namespace.canGrant"
                  :resource="`corteza::compose:chart/${namespace.namespaceID}/*`"
                  :button-label="$t('general.label.permissions')"
                  button-variant="light"
                  class="btn-lg"
                />
              </div>
            </template>

            <template #actions="{ item: c }">
              <c-permissions-button
                v-if="c.canGrant"
                :title="c.name"
                :target="c.name"
                :resource="`corteza::compose:chart/${namespace.namespaceID}/${c.chartID}`"
                link
                class="btn px-2"
              />
            </template>

            <template #updatedAt="{ item }">
              {{ (item.updatedAt || item.createdAt) | locFullDateTime }}
            </template>
          </c-resource-list>
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>
<script>
import { mapGetters, mapActions } from 'vuex'
import { compose } from '@cortezaproject/corteza-js'
import Import from 'corteza-webapp-compose/src/components/Admin/Import'
import Export from 'corteza-webapp-compose/src/components/Admin/Export'
import listHelpers from 'corteza-webapp-compose/src/mixins/listHelpers'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'ChartList',

  components: {
    Import,
    Export,
  },

  mixins: [
    listHelpers,
  ],

  props: {
    namespace: {
      type: Object,
      required: false,
      default: undefined,
    },
  },

  data () {
    return {
      primaryKey: 'chartID',

      filter: {
        query: '',
        namespaceID: this.namespace.namespaceID,
      },

      sorting: {
        sortBy: 'name',
        sortDesc: false,
      },

      newChart: new compose.Chart({}),
    }
  },

  computed: {
    ...mapGetters({
      charts: 'chart/set',
    }),

    tableFields () {
      return [
        {
          key: 'name',
          sortable: true,
          tdClass: 'text-nowrap',
        },
        {
          key: 'handle',
          sortable: true,
        },
        {
          key: 'updatedAt',
          sortable: true,
          class: 'text-right text-nowrap',
        },
        {
          key: 'actions',
          label: '',
          tdClass: 'text-right text-nowrap',
        },
      ]
    },
  },

  methods: {
    chartList () {
      return this.procListResults(this.$ComposeAPI.chartList(this.encodeListParams()))
    },

    ...mapActions({
      createChart: 'chart/create',
    }),

    create (subType) {
      let c = new compose.Chart({ ...this.newChart, namespaceID: this.namespace.namespaceID })
      switch (subType) {
        case 'gauge':
          c = new compose.GaugeChart(c)
          break

        case 'funnel':
          c = new compose.FunnelChart(c)
          break
      }

      this.createChart(c).then((chart) => {
        this.$router.push({ name: 'admin.charts.edit', params: { chartID: chart.chartID } })
      }).catch(this.toastErrorHandler(this.$t('notification:chart.createFailed')))
    },

    encodeRouteParams () {
      const { query } = this.filter
      const { limit, pageCursor, page } = this.pagination

      return {
        query: {
          limit,
          ...this.sorting,
          query,
          page,
          pageCursor,
        },
      }
    },

    handleRowClicked ({ chartID, canUpdateChart, canDeleteChart }) {
      if (!(canUpdateChart || canDeleteChart)) {
        return
      }
      this.$router.push({
        name: 'admin.charts.edit',
        params: { chartID },
        query: null,
      })
    },
  },
}
</script>
<style lang="scss" scoped>
$input-height: 42px;

.chart-name-input {
  height: $input-height;
}
</style>
