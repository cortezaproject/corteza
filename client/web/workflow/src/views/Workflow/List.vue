<template>
  <b-container
    fluid="xl"
    class="d-flex flex-column py-3"
  >
    <portal to="topbar-title">
      {{ $t('general:workflow-list') }}
    </portal>

    <c-resource-list
      :primary-key="primaryKey"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :fields="tableFields"
      :items="workflowList"
      :row-class="genericRowClass"
      :translations="{
        searchPlaceholder: $t('general:searchPlaceholder'),
        notFound: $t('general:resourceList.notFound'),
        noItems: $t('general:resourceList.noItems'),
        loading: $t('general:loading'),
        showingPagination: 'general:resourceList.pagination.showing',
        singlePluralPagination: 'general:resourceList.pagination.single',
        prevPagination: $t('general:resourceList.pagination.prev'),
        nextPagination: $t('general:resourceList.pagination.next'),
        resourceSingle: $t('general:workflow.single'),
        resourcePlural: $t('general:workflow.plural')
      }"
      clickable
      sticky-header
      class="h-100 flex-fill"
      @search="filterList"
      @row-clicked="handleRowClicked"
    >
      <template #header>
        <b-button
          v-if="canCreate"
          data-test-id="button-create-workflow"
          variant="primary"
          size="lg"
          :to="{ name: 'workflow.create' }"
        >
          {{ $t('general:new-workflow') }}
        </b-button>

        <import
          v-if="canCreate"
          :disabled="importProcessing"
          class="d-flex"
          @import="importJSON"
        />

        <export size="lg" />

        <c-permissions-button
          v-if="canGrant"
          resource="corteza::automation:workflow/*"
          :button-label="$t('general:permissions')"
          size="lg"
        />
      </template>

      <template #toolbar>
        <b-col>
          <b-form-radio-group
            v-model="filter.subWorkflow"
            :options="[
              { value: 0, text: $t('general:without') },
              { value: 1, text: $t('general:including') },
              { value: 2, text: $t('general:only') }
            ]"
            buttons
            button-variant="outline-primary"
            size="sm"
            name="radio-btn-outline"
            @change="filterList"
          />
          {{ $t('general:subworkflows') }}
        </b-col>
        <b-col>
          <b-form-radio-group
            v-model="filter.disabled"
            :options="[
              { value: 0, text: $t('general:without') },
              { value: 1, text: $t('general:including') },
              { value: 2, text: $t('general:only') }
            ]"
            buttons
            button-variant="outline-primary"
            size="sm"
            name="radio-btn-outline"
            @change="filterList"
          />
          {{ $t('general:disabled') }}
        </b-col>
        <b-col>
          <b-form-radio-group
            v-model="filter.deleted"
            :options="[
              { value: 0, text: $t('general:without') },
              { value: 1, text: $t('general:including') },
              { value: 2, text: $t('general:only') }
            ]"
            buttons
            button-variant="outline-primary"
            size="sm"
            name="radio-btn-outline"
            @change="filterList"
          />
          {{ $t('general:deleted') }}
        </b-col>
      </template>

      <template #name="{ item: w }">
        {{ w.meta.name || w.handle }}
        <h5 class="d-inline-block ml-2">
          <b-badge
            v-if="w.meta.subWorkflow"
            variant="info"
          >
            {{ $t('general:subworkflow') }}
          </b-badge>
        </h5>
      </template>

      <template #enabled="{ item: w }">
        <font-awesome-icon
          :icon="['fas', w.enabled ? 'check' : 'times']"
        />
      </template>

      <template #changedAt="{ item }">
        {{ (item.deletedAt || item.updatedAt || item.createdAt) | locFullDateTime }}
      </template>

      <template #actions="{ item: w }">
        <b-dropdown
          v-if="w.canGrant || w.canDeleteWorkflow"
          variant="outline-light"
          toggle-class="d-flex align-items-center justify-content-center text-primary border-0 py-2 ml-1"
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
            v-if="w.canGrant"
            link-class="p-0"
            variant="light"
          >
            <c-permissions-button
              :tooltip="$t('permissions:resources.automation.workflow.tooltip')"
              :title="w.meta.name || w.handle || w.workflowID"
              :target="w.meta.name || w.handle || w.workflowID"
              :resource="`corteza::automation:workflow/${w.workflowID}`"
              :button-label="$t('permissions:ui.label')"
              button-variant="link dropdown-item text-decoration-none text-dark regular-font rounded-0"
            />
          </b-dropdown-item>

          <b-dropdown-item>
            <font-awesome-icon
              :icon="['fas', 'file-export']"
            />

            <export
              data-test-id="button-export-workflow"
              :workflows="([w.workflowID])"
              :file-name="w.meta.name || w.handle"
              variant="link"
              size="md"
              class="text-decoration-none text-dark regular-font p-0 ml-1"
            />
          </b-dropdown-item>

          <b-dropdown-item-button @click="handleStatusChange(w)">
            <font-awesome-icon
              :icon="['fas', w.enabled ? 'toggle-off' : 'toggle-on']"
            />
            {{ statusText(w) }}
          </b-dropdown-item-button>

          <c-input-confirm
            v-if="w.canDeleteWorkflow && !w.deletedAt"
            borderless
            variant="link"
            size="md"
            show-icon
            :text="$t('delete')"
            text-class="p-1"
            button-class="dropdown-item text-decoration-none text-dark regular-font rounded-0"
            icon-class="text-danger"
            class="w-100"
            @confirmed="handleDelete(w)"
          />

          <c-input-confirm
            v-if="w.canUndeleteWorkflow && w.deletedAt"
            borderless
            variant="link"
            size="md"
            show-icon
            :text="$t('undelete')"
            text-class="p-1"
            button-class="dropdown-item text-decoration-none text-dark regular-font rounded-0"
            icon-class="text-danger"
            class="w-100"
            @confirmed="handleDelete(w)"
          />
        </b-dropdown>
      </template>
    </c-resource-list>
  </b-container>
</template>

<script>
import { mapGetters } from 'vuex'
import Import from '../../components/Import'
import Export from '../../components/Export'
import listHelpers from '../../mixins/listHelpers'
import { components } from '@cortezaproject/corteza-vue'
const { CResourceList } = components

export default {
  i18nOptions: {
    namespaces: 'list',
  },

  name: 'WorkflowList',

  components: {
    Import,
    Export,
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
        deleted: 0,
        subWorkflow: 1,
        disabled: 0,
      },

      sorting: {
        sortBy: 'handle',
        sortDesc: false,
      },

      newWorkflow: {},

      importProcessing: false,
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canGrant () {
      return this.can('automation/', 'grant')
    },

    canCreate () {
      return this.can('automation/', 'workflow.create')
    },

    tableFields () {
      return [
        {
          key: 'name',
          label: this.$t('general:columns.name'),
          sortable: false,
          tdClass: 'text-nowrap',
        },
        {
          key: 'enabled',
          label: this.$t('general:columns.enabled'),
          sortable: true,
          class: 'text-center',
        },
        {
          key: 'steps',
          label: this.$t('general:columns.steps'),
          tdClass: 'align-middle',
          class: 'text-center',
          formatter: steps => {
            return (steps || []).length
          },
        },
        {
          key: 'changedAt',
          label: this.$t('general:columns.changedAt'),
          sortable: true,
          class: 'text-right text-nowrap',
        },
        {
          key: 'actions',
          label: '',
          tdClass: 'text-right text-nowrap actions',
        },
      ]
    },

    workflowIDs () {
      return this.workflows.map(({ workflowID }) => workflowID)
    },

    userID () {
      if (this.$auth.user) {
        return this.$auth.user.userID
      }
      return undefined
    },
  },

  methods: {
    async importJSON (workflows = []) {
      this.importProcessing = true

      const skippedWorkflows = []

      await Promise.all(workflows.map(({ triggers = [], ...wf }) => {
        // Create workflow
        return this.$AutomationAPI.workflowCreate({ ownedBy: this.userID, runAs: '0', ...wf })
          .then(({ workflowID }) => {
            // Create triggers
            return Promise.all(triggers.map(trigger => {
              return this.$AutomationAPI.triggerCreate({
                ...trigger,
                workflowID,
                workflowStepID: trigger.stepID,
                ownedBy: this.userID,
              })
            }))
          })
          .catch(({ message }) => {
            // Skip workflow and add to skipped list
            if (wf.handle) {
              skippedWorkflows.push(`${wf.handle}${message ? ' - ' + message : ''};`)
            }
          })
      }))
        .then(() => {
          if (skippedWorkflows.length) {
            this.toastInfo(`${skippedWorkflows.join(' ')}`, this.$t('notification:import.skipped-workflows'))
          } else {
            this.toastSuccess(this.$t('notification:import.imported-workflows'))
          }
          this.$bvModal.hide('import')
        })
        .catch(this.toastErrorHandler(this.$t('notification:import.failed-import')))

      this.$root.$emit('bv::refresh::table', 'resource-list')

      this.importProcessing = false
    },

    workflowList () {
      return this.procListResults(this.$AutomationAPI.workflowList(this.encodeListParams()))
    },

    handleRowClicked (workflow) {
      this.$router.push({ name: 'workflow.edit', params: { workflowID: workflow.workflowID } })
    },

    handleDelete (workflow) {
      const { deletedAt = '' } = workflow
      const method = deletedAt ? 'workflowUndelete' : 'workflowDelete'
      const event = deletedAt ? 'undelete' : 'delete'
      const { workflowID } = workflow
      this.$AutomationAPI[method]({ workflowID })
        .then(() => {
          this.toastSuccess(this.$t(`notification:${event}.success`))
          this.filterList()
        })
        .catch(this.toastErrorHandler(this.$t(`notification:${event}.failed`)))
    },

    statusText (w) {
      return w.enabled ? this.$t('general:disable') : this.$t('general:enable')
    },

    handleStatusChange ({ workflowID, enabled }) {
      enabled = !enabled
      const notificationKey = enabled ? 'enable' : 'disable'

      this.$AutomationAPI.workflowRead({ workflowID }).then((w) => {
        return this.$AutomationAPI.workflowUpdate({ ...w, enabled }).then((w) => {
          this.toastSuccess(this.$t(`notification:list.${notificationKey}.success`))
          this.filterList()
        })
      }).catch(this.toastErrorHandler(this.$t(`notification:list.${notificationKey}.failed`)))
    },
  },
}
</script>
