<template>
  <div class="d-flex w-100 py-3">
    <portal to="topbar-title">
      {{ $t('report.list') }}
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
            :items="reportList"
            :translations="{
              searchPlaceholder: $t('searchPlaceholder'),
              notFound: $t('general:resourceList.notFound'),
              noItems: $t('general:resourceList.noItems'),
              loading: $t('general:label.loading'),
              showingPagination: 'general:resourceList.pagination.showing',
              singlePluralPagination: 'general:resourceList.pagination.single',
              prevPagination: $t('general:resourceList.pagination.prev'),
              nextPagination: $t('general:resourceList.pagination.next'),
            }"
            clickable
            class="h-100"
            @search="filterList"
            @row-clicked="viewReport"
          >
            <template #header>
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
            </template>

            <template #name="{ item: r }">
              {{ r.meta.name }}
            </template>

            <template #changedAt="{ item }">
              {{ (item.deletedAt || item.updatedAt || item.createdAt) | locFullDateTime }}
            </template>

            <template #actions="{ item: r }">
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
                :tooltip="$t('permissions:resources.system.report.tooltip')"
                :title="r.meta.name || r.handle || r.reportID"
                :target="r.meta.name || r.handle || r.reportID"
                :resource="`corteza::system:report/${r.reportID}`"
                class="btn px-2"
                link
              />
            </template>
          </c-resource-list>
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import listHelpers from 'corteza-webapp-reporter/src/mixins/listHelpers'
import { components } from '@cortezaproject/corteza-vue'
const { CResourceList } = components

export default {
  name: 'ReportList',

  i18nOptions: {
    namespaces: 'list',
  },

  components: {
    CResourceList,
  },

  mixins: [
    listHelpers,
  ],

  data () {
    return {
      primaryKey: 'reportID',

      filter: {
        query: '',
      },

      sorting: {
        sortBy: 'handle',
        sortDesc: false,
      },
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
          label: this.$t('columns.name'),
          sortable: true,
          tdClass: 'text-nowrap',
        },
        {
          key: 'handle',
          label: this.$t('columns.handle'),
          sortable: true,
        },
        {
          key: 'changedAt',
          label: this.$t('columns.changedAt'),
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

    reportList () {
      return this.procListResults(this.$SystemAPI.reportList(this.encodeListParams()))
    },
  },
}
</script>
