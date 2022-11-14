<template>
  <div class="d-flex w-100 py-3">
    <portal to="topbar-title">
      {{ $t('report.list') }}
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
                class="justify-content-between wrap-with-vertical-gutters"
                no-gutters
              >
                <div class="flex-grow-1">
                  <div
                    class="wrap-with-vertical-gutters"
                  >
                    <b-button
                      v-if="canCreate"
                      data-test-id="button-create-report"
                      variant="primary"
                      size="lg"
                      class="mr-1"
                      :to="{ name: 'report.create' }"
                    >
                      {{ $t('report.new') }}
                    </b-button>

                    <c-permissions-button
                      v-if="canGrant"
                      resource="corteza::system:report/*"
                      :button-label="$t('permissions')"
                      button-variant="light"
                      class="btn-lg"
                    />
                  </div>
                </div>

                <div class="flex-grow-1">
                  <b-input-group
                    class="h-100 mw-100"
                  >
                    <b-input
                      v-model.trim="query"
                      class="h-100 mw-100 text-truncate"
                      :placeholder="$t('searchPlaceholder')"
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
                :items="reports"
                :filter="query"
                :filter-included-fields="['handle', 'name']"
                head-variant="light"
                tbody-tr-class="pointer"
                responsive
                hover
                :class="{ 'mb-0': !!reports.length }"
                @row-clicked="viewReport"
              >
                <template v-slot:cell(name)="{ item: r }">
                  {{ r.meta.name }}
                </template>
                <template v-slot:cell(actions)="{ item: r }">
                  <b-button
                    v-if="r.canUpdateReport"
                    variant="light"
                    class="mr-2"
                    :to="{ name: 'report.builder', params: { reportID: r.reportID } }"
                  >
                    {{ $t('report.builder') }}
                  </b-button>
                  <b-button
                    v-if="r.canUpdateReport"
                    variant="link"
                    class="mr-2"
                    :to="{ name: 'report.edit', params: { reportID: r.reportID } }"
                  >
                    {{ $t('report.edit') }}
                  </b-button>
                  <c-permissions-button
                    v-if="r.canGrant"
                    :title="r.handle"
                    :target="r.handle"
                    resource="corteza::system:report/*"
                    class="btn px-2"
                    link
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
import { mapGetters } from 'vuex'

export default {
  name: 'ReportList',

  i18nOptions: {
    namespaces: 'list',
  },

  data () {
    return {
      query: '',
      reports: [],
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canGrant () {
      return this.can('system/', 'grant')
    },

    canCreate () {
      return this.can('system/', 'report.create')
    },

    tableFields () {
      return [
        {
          key: 'name',
          sortable: true,
          filterByFormatted: true,
          tdClass: 'align-middle pl-4 text-nowrap',
          thClass: 'pl-4',
          formatter: (name, key, item) => {
            return item.meta.name
          },
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
            return new Date(updatedAt || item.createdAt).toLocaleDateString('en-US')
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

  created () {
    this.fetchReports()
  },

  methods: {
    fetchReports () {
      this.$SystemAPI.reportList()
        .then(({ set = [] }) => {
          this.reports = set
        })
        .catch(this.toastErrorHandler(this.$t('notification:report.listFetchFailed')))
    },

    viewReport ({ reportID, canReadReport = false }) {
      if (reportID) {
        if (canReadReport) {
          this.$router.push({
            name: 'report.view',
            params: { reportID },
          })
        } else {
          this.toastDanger(this.$t('notification:report.notAllowed.read'))
        }
      }
    },
  },
}
</script>
