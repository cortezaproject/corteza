<template>
  <div class="py-3">
    <portal to="topbar-title">
      {{ $t('navigation.chart') }}
    </portal>

    <b-container fluid="xl">
      <b-row no-gutters>
        <b-col>
          <b-card
            no-body
            class="shadow-sm"
          >
            <b-card-header
              header-bg-variant="white"
              class="py-3"
            >
              <b-row
                class="wrap-with-vertical-gutters justify-content-between"
                no-gutters
              >
                <div class="flex-grow-1">
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
                </div>
                <div class="flex-grow-1">
                  <b-input-group
                    class="h-100"
                  >
                    <b-form-input
                      v-model.trim="query"
                      class="h-100 mw-100 text-truncate"
                      type="search"
                      :placeholder="$t('chart.searchPlaceholder')"
                    />
                    <b-input-group-append>
                      <b-input-group-text class="text-primary bg-white">
                        <font-awesome-icon
                          :icon="['fas', 'search']"
                        />
                      </b-input-group-text>
                    </b-input-group-append>
                  </b-input-group>
                </div>
              </b-row>
            </b-card-header>

            <b-card-body class="p-0">
              <b-table
                :fields="tableFields"
                :items="charts"
                :filter="query"
                :filter-function="chartFilter"
                :sort-by.sync="sortBy"
                :sort-desc="sortDesc"
                head-variant="light"
                tbody-tr-class="pointer"
                :empty-text="$t('chart.noChart')"
                show-empty
                responsive
                hover
                @row-clicked="handleRowClicked"
              >
                <template v-slot:cell(updatedAt)="{ item: c }">
                  {{ (c.updatedAt || c.createdAt) | locDate }}
                </template>
                <template v-slot:cell(actions)="{ item: c }">
                  <c-permissions-button
                    v-if="c.canGrant"
                    :title="c.name"
                    :target="c.name"
                    :resource="`corteza::compose:chart/${namespace.namespaceID}/${c.chartID}`"
                    link
                    class="btn px-2"
                  />
                </template>
              </b-table>
            </b-card-body>
          </b-card>
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>
<script>
import { mapGetters, mapActions } from 'vuex'
import { compose, fmt } from '@cortezaproject/corteza-js'
import { filter } from '@cortezaproject/corteza-vue'
import Import from 'corteza-webapp-compose/src/components/Admin/Import'
import Export from 'corteza-webapp-compose/src/components/Admin/Export'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'ChartList',

  components: {
    Import,
    Export,
  },

  props: {
    namespace: {
      type: Object,
      required: false,
      default: undefined,
    },
  },

  data () {
    return {
      query: '',

      sortBy: 'name',
      sortDesc: false,

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
          tdClass: 'align-middle pl-4 text-nowrap',
          thClass: 'pl-4',
        },
        {
          key: 'handle',
          sortable: true,
          tdClass: 'align-middle',
        },
        {
          key: 'updatedAt',
          sortable: true,
          sortByFormatted: true,
          tdClass: 'align-middle',
          class: 'text-right',
          formatter: (updatedAt, key, item) => {
            return fmt.date(updatedAt || item.createdAt)
          },
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
    ...mapActions({
      createChart: 'chart/create',
    }),

    chartFilter (chart, query) {
      return filter.Assert(chart, query, 'handle', 'name')
    },

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
