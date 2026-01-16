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

        <b-button
          variant="light"
          size="lg"
          @click="$bvModal.show('workflow-filter')"
        >
          <font-awesome-icon
            :icon="['fas', 'filter']"
            :class="['mr-1', { 'text-primary': labelFilterCount > 0 }]"
          />
          {{ $t('general:filter.label') }}
        </b-button>

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
        <div>
          {{ w.meta.name || w.handle }}
          <b-badge
            v-if="w.meta.subWorkflow"
            variant="info"
            class="ml-2"
          >
            {{ $t('general:subworkflow') }}
          </b-badge>
        </div>
        <div
          v-for="group in getWorkflowLabels(w)"
          :key="'group-' + group.namespaceID"
          class="d-flex align-items-center flex-wrap gap-1 mt-2"
        >
          <b-badge
            v-b-tooltip.hover.d200
            :title="$t('general:filter.namespace.label')"
            variant="primary"
            style="font-size: 90%;"
          >
            {{ group.namespaceName }}
          </b-badge>

          <b-badge
            v-for="mod in group.modules"
            :key="'mod-' + group.namespaceID + '-' + mod.id"
            v-b-tooltip.hover.d200
            :title="$t('general:filter.module.label')"
            variant="extra-light"
            style="font-size: 90%;"
          >
            {{ mod.name }}
          </b-badge>
        </div>
      </template>

      <template #enabled="{ item: w }">
        <font-awesome-icon
          :icon="['fas', w.enabled ? 'check' : 'times']"
          :class="w.enabled ? 'text-primary' : 'text-extra-light'"
        />
      </template>

      <template #changedAt="{ item }">
        {{ (item.deletedAt || item.updatedAt || item.createdAt) | locFullDateTime }}
      </template>

      <template #actions="{ item: w }">
        <b-dropdown
          variant="outline-extra-light"
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

          <b-dropdown-item-button @click="handleStatusChange(w)">
            <font-awesome-icon
              :icon="['fas', w.enabled ? 'toggle-off' : 'toggle-on']"
            />
            {{ statusText(w) }}
          </b-dropdown-item-button>

          <export
            data-test-id="button-export-workflow"
            :workflows="([w.workflowID])"
            :file-name="w.meta.name || w.handle"
            size="md"
            class="dropdown-item"
          >
            <font-awesome-icon
              :icon="['fas', 'file-export']"
            />
          </export>

          <c-permissions-button
            v-if="w.canGrant"
            :tooltip="$t('permissions:resources.automation.workflow.tooltip')"
            :title="w.meta.name || w.handle || w.workflowID"
            :target="w.meta.name || w.handle || w.workflowID"
            :resource="`corteza::automation:workflow/${w.workflowID}`"
            :button-label="$t('permissions:ui.label')"
            class="dropdown-item"
          />

          <c-input-confirm
            v-if="w.canDeleteWorkflow && !w.deletedAt"
            borderless
            variant="link"
            size="md"
            show-icon
            :text="$t('delete')"
            text-class="p-1"
            button-class="dropdown-item"
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
            button-class="dropdown-item"
            icon-class="text-danger"
            class="w-100"
            @confirmed="handleDelete(w)"
          />
        </b-dropdown>
      </template>
    </c-resource-list>

    <workflow-filter-modal
      :namespace-labels="selectedNamespaceLabels"
      :module-labels="selectedModuleLabels"
      @apply="handleFilterApply"
    />
  </b-container>
</template>

<script>
import { mapGetters } from 'vuex'
import Import from '../../components/Import'
import Export from '../../components/Export'
import WorkflowFilterModal from '../../components/WorkflowFilterModal'
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
    WorkflowFilterModal,
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
        labels: [],
      },

      sorting: {
        sortBy: 'changedAt',
        sortDesc: true,
      },

      newWorkflow: {},

      importProcessing: false,

      // Label filter selections (resource paths)
      selectedNamespaceLabels: [],
      selectedModuleLabels: [],
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
      labelCache: 'labels/getNamespace',
      getModule: 'labels/getModule',
    }),

    canGrant () {
      return this.can('automation/', 'grant')
    },

    canCreate () {
      return this.can('automation/', 'workflow.create')
    },

    labelFilterCount () {
      return this.selectedNamespaceLabels.length + this.selectedModuleLabels.length
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
    handleFilterApply ({ namespaceLabels, moduleLabels }) {
      this.selectedNamespaceLabels = namespaceLabels || []
      this.selectedModuleLabels = moduleLabels || []
      this.updateLabelsFilter()
    },

    updateLabelsFilter () {
      const labels = []

      // Send labels directly - no transformation needed
      if (this.selectedNamespaceLabels.length > 0) {
        labels.push(`ref_namespace=${JSON.stringify(this.selectedNamespaceLabels)}`)
      }

      if (this.selectedModuleLabels.length > 0) {
        labels.push(`ref_module=${JSON.stringify(this.selectedModuleLabels)}`)
      }

      this.filter.labels = labels
      this.filterList()
    },

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
      return this.procListResults(
        this.$AutomationAPI.workflowListCancellable(this.encodeListParams()),
      ).then(workflows => {
        // Pre-resolve all namespace and module names
        if (workflows && workflows.length > 0) {
          return this.resolveLabelsForWorkflows(workflows).then(() => workflows)
        }
        return workflows
      })
    },

    async resolveLabelsForWorkflows (workflows) {
      const namespaceIDs = new Set()
      const modules = []

      // Collect all unique IDs
      workflows.forEach(workflow => {
        if (workflow.labels?.ref_namespace) {
          const nsValues = Array.isArray(workflow.labels.ref_namespace)
            ? workflow.labels.ref_namespace
            : [workflow.labels.ref_namespace]
          nsValues.forEach(label => {
            const nsID = label.split('/')[1]
            if (nsID) namespaceIDs.add(nsID)
          })
        }

        if (workflow.labels?.ref_module) {
          const modValues = Array.isArray(workflow.labels.ref_module)
            ? workflow.labels.ref_module
            : [workflow.labels.ref_module]
          modValues.forEach(label => {
            const parts = label.split('/')
            const nsID = parts[1]
            const modID = parts[2]
            if (nsID) namespaceIDs.add(nsID)
            if (modID) modules.push({ moduleID: modID, namespaceID: nsID })
          })
        }
      })

      // Resolve using store
      await Promise.all([
        this.$store.dispatch('labels/resolveMultipleNamespaces', {
          namespaceIDs: Array.from(namespaceIDs),
          api: this.$ComposeAPI,
        }),
        this.$store.dispatch('labels/resolveMultipleModules', {
          modules,
          api: this.$ComposeAPI,
        }),
      ])
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

    /**
     * Get resolved label names for a workflow, grouped by namespace
     * Returns array of { namespaceID, namespaceName, modules: [{ id, name }] }
     *
     * Labels are stored as:
     * - ref_namespace: ["corteza::compose:namespace/id1", ...]
     * - ref_module: ["corteza::compose:module/nsID/modID", ...]
     */
    getWorkflowLabels (workflow) {
      if (!workflow.labels) {
        return []
      }

      const namespaceIDs = []
      const modulesByNamespace = {}

      // Parse namespace labels - extract IDs for display
      if (workflow.labels.ref_namespace) {
        const nsValues = Array.isArray(workflow.labels.ref_namespace)
          ? workflow.labels.ref_namespace
          : [workflow.labels.ref_namespace]

        nsValues.forEach(label => {
          const nsID = label.split('/')[1] // Extract ID from corteza::compose:namespace/ID
          if (nsID && !namespaceIDs.includes(nsID)) {
            namespaceIDs.push(nsID)
          }
        })
      }

      // Parse module labels - extract IDs for display
      if (workflow.labels.ref_module) {
        const modValues = Array.isArray(workflow.labels.ref_module)
          ? workflow.labels.ref_module
          : [workflow.labels.ref_module]

        modValues.forEach(label => {
          const parts = label.split('/') // Split corteza::compose:module/nsID/modID
          const nsID = parts[1]
          const modID = parts[2]

          if (!nsID || !modID) return

          // Ensure namespace is in the list
          if (!namespaceIDs.includes(nsID)) {
            namespaceIDs.push(nsID)
          }

          if (!modulesByNamespace[nsID]) {
            modulesByNamespace[nsID] = []
          }

          const name = this.$store.getters['labels/getModule'](modID)
          modulesByNamespace[nsID].push({
            id: modID,
            name: name || modID,
          })
        })
      }

      if (namespaceIDs.length === 0) {
        return []
      }

      // Build the grouped result - names are already resolved from store
      return namespaceIDs.map(nsID => {
        const name = this.$store.getters['labels/getNamespace'](nsID)
        return {
          namespaceID: nsID,
          namespaceName: name || nsID,
          modules: modulesByNamespace[nsID] || [],
        }
      })
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
