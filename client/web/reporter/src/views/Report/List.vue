<template>
  <b-container
    fluid="xl"
    class="d-flex flex-column py-3"
  >
    <portal to="topbar-title">
      {{ $t('report.list') }}
    </portal>

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
        resourceSingle: $t('general:label.report.single'),
        resourcePlural: $t('general:label.report.plural')
      }"
      sticky-header
      clickable
      class="h-100 flex-fill"
      @search="filterList"
      @row-clicked="viewReport"
    >
      <template #header>
        <b-button
          v-if="canCreate"
          data-test-id="button-create-report"
          variant="primary"
          size="lg"
          :to="{ name: 'report.create' }"
        >
          {{ $t('report.new') }}
        </b-button>

        <c-permissions-button
          v-if="canGrant"
          resource="corteza::system:report/*"
          :button-label="$t('permissions')"
          size="lg"
        />
      </template>

      <template #name="{ item: r }">
        {{ r.meta.name }}
      </template>

      <template #changedAt="{ item }">
        {{ (item.deletedAt || item.updatedAt || item.createdAt) | locFullDateTime }}
      </template>

      <template #actions="{ item: r }">
        <b-button-group
          v-if="r.canUpdateReport"
          size="sm"
        >
          <b-button
            data-test-id="button-report-builder"
            variant="primary"
            size="sm"
            :to="{ name: 'report.builder', params: { reportID: r.reportID } }"
          >
            {{ $t('report.builder') }}
            <font-awesome-icon
              :icon="['fas', 'tools']"
              class="ml-2"
            />
          </b-button>

          <b-button
            v-b-tooltip.hover="{ title: $t('report.edit'), container: '#body' }"
            data-test-id="button-report-edit"
            variant="primary"
            :to="{ name: 'report.edit', params: { reportID: r.reportID } }"
            class="d-flex align-items-center"
            style="margin-left:2px;"
          >
            <font-awesome-icon
              :icon="['far', 'edit']"
            />
          </b-button>
        </b-button-group>
      </template>

      <template #moreActions="{ item: r }">
        <b-dropdown
          v-if="r.canUpdateReport || r.canGrant || r.canDeleteReport"
          variant="outline-light"
          toggle-class="d-flex align-items-center justify-content-center text-primary border-0 py-2"
          no-caret
          lazy
          menu-class="m-0"
        >
          <template #button-content>
            <font-awesome-icon
              :icon="['fas', 'ellipsis-v']"
            />
          </template>

          <b-dropdown-item
            v-if="r.canGrant"
            link-class="p-0"
            variant="light"
          >
            <c-permissions-button
              :tooltip="$t('permissions:resources.system.report.tooltip')"
              :title="r.meta.name || r.handle || r.reportID"
              :target="r.meta.name || r.handle || r.reportID"
              :resource="`corteza::system:report/${r.reportID}`"
              class="text-dark d-print-none border-0"
              :button-label="$t('permissions:ui.label')"
              button-variant="link dropdown-item text-decoration-none text-dark regular-font rounded-0"
            />
          </b-dropdown-item>

          <c-input-confirm
            v-if="r.canDeleteReport"
            :processing="processingDelete"
            :text="$t('report.delete')"
            borderless
            variant="link"
            size="md"
            show-icon
            text-class="p-1"
            button-class="dropdown-item text-decoration-none regular-font rounded-0"
            icon-class="text-danger"
            class="w-100"
            @confirmed="handleDelete(r)"
          />
        </b-dropdown>
      </template>
    </c-resource-list>
  </b-container>
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

      processingDelete: false,

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
          sortable: false,
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
        {
          key: 'moreActions',
          label: '',
          tdClass: 'text-right text-nowrap actions',
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

    handleDelete (report) {
      this.processingDelete = true

      return this.$SystemAPI.reportDelete(report)
        .then(() => {
          this.toastSuccess(this.$t('notification:report.delete'))
          this.filterList()
        })
        .catch(this.toastErrorHandler(this.$t('notification:report.deleteFailed')))
        .finally(() => {
          this.processingDelete = false
        })
    },
  },
}
</script>
