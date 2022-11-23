<template>
  <div class="h-100 py-3 flex-grow-1 overflow-auto">
    <portal to="topbar-title">
      {{ $t('general:workflow-list') }}
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
            :items="workflowList"
            :translations="{
              searchPlaceholder: $t('general:searchPlaceholder'),
              notFound: $t('general:resourceList.notFound'),
              noItems: $t('general:resourceList.noItems'),
              loading: $t('general:loading'),
              showingPagination: 'general:resourceList.pagination.showing',
              singlePluralPagination: 'general:resourceList.pagination.single',
              prevPagination: $t('general:resourceList.pagination.prev'),
              nextPagination: $t('general:resourceList.pagination.next'),
            }"
            clickable
            class="h-100"
            @search="filterList"
            @row-clicked="handleRowClicked"
          >
            <template #header>
              <b-button
                v-if="canCreate"
                data-test-id="button-create-workflow"
                variant="primary"
                size="lg"
                class="float-left mr-1"
                :to="{ name: 'workflow.create' }"
              >
                {{ $t('general:new-workflow') }}
              </b-button>

              <import
                v-if="canCreate"
                :disabled="importProcessing"
                class="float-left mr-1"
                @import="importJSON"
              />

              <export
                class="float-left mr-1"
              />

              <c-permissions-button
                v-if="canGrant"
                resource="corteza::automation:workflow/*"
                :button-label="$t('general:permissions')"
                button-variant="light"
                class="float-left btn-lg"
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

            <template #handle="{ item: w }">
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
              <c-permissions-button
                v-if="w.canGrant"
                :tooltip="$t('permissions:resources.automation.workflow.tooltip')"
                :title="w.meta.name || w.handle"
                :target="w.meta.name || w.handle"
                :resource="`corteza::automation:workflow/${w.workflowID}`"
                link
                class="btn px-2"
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
import Import from '../../components/Import'
import Export from '../../components/Export'
import listHelpers from '../../mixins/listHelpers'
import { components } from '@cortezaproject/corteza-vue'
const { CResourceList } = components

export default {
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
          key: 'handle',
          label: this.$t('general:columns.name'),
          sortable: true,
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
          tdClass: 'text-right text-nowrap',
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
  },
}
</script>
