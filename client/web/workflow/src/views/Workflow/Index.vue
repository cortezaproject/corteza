<template>
  <div class="h-100 py-3 flex-grow-1 overflow-auto">
    <portal to="topbar-title">
      {{ $t('general:workflow-list') }}
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
                      :workflows="workflowIDs"
                      class="float-left mr-1"
                    />

                    <c-permissions-button
                      v-if="canGrant"
                      resource="corteza::automation:workflow/*"
                      :button-label="$t('general:permissions')"
                      button-variant="light"
                      class="float-left btn-lg"
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
                      type="search"
                      debounce="300"
                      :placeholder="$t('general:search-workflows')"
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
                :items="workflows"
                :filter="query"
                :filter-included-fields="['handle']"
                :sort-by.sync="sortBy"
                :sort-desc="sortDesc"
                head-variant="light"
                tbody-tr-class="pointer"
                responsive
                hover
                @row-clicked="handleRowClicked"
              >
                <template v-slot:cell(handle)="{ item: w }">
                  {{ w.meta.name || w.handle }}
                </template>
                <template v-slot:cell(enabled)="{ item: w }">
                  <font-awesome-icon
                    :icon="['fas', w.enabled ? 'check' : 'times']"
                  />
                </template>
                <template v-slot:cell(actions)="{ item: w }">
                  <c-permissions-button
                    v-if="w.canGrant"
                    :title="w.meta.name || w.handle"
                    :target="w.meta.name || w.handle"
                    :resource="`corteza::automation:workflow/${w.workflowID}`"
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
import Import from '../../components/Import'
import Export from '../../components/Export'
import { mapGetters } from 'vuex'

export default {
  name: 'WorkflowList',

  components: {
    Import,
    Export,
  },

  data () {
    return {
      workflows: [],

      query: '',

      sortBy: 'handle',
      sortDesc: false,

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
          label: this.$t('general:name'),
          sortable: true,
          tdClass: 'align-middle text-nowrap',
          class: 'pl-4',
        },
        {
          key: 'enabled',
          sortable: true,
          tdClass: 'align-middle',
          class: 'text-center',
        },
        {
          key: 'steps',
          sortable: true,
          sortByFormatted: true,
          tdClass: 'align-middle',
          class: 'text-center',
          formatter: steps => {
            return (steps || []).length
          },
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

  created () {
    this.fetchWorkflows()
  },

  methods: {
    fetchWorkflows () {
      this.$AutomationAPI.workflowList({ disabled: 1 })
        .then(({ set = [] }) => {
          this.workflows = set
        })
        .catch(this.defaultErrorHandler(this.$t('notification:failed-fetch-workflows')))
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
            this.raiseInfoAlert(`${skippedWorkflows.join(' ')}`, this.$t('notification:import.skipped-workflows'))
          } else {
            this.raiseSuccessAlert(this.$t('notification:import.imported-workflows'))
          }
        })
        .catch(this.defaultErrorHandler(this.$t('notification:import.failed-import')))

      await this.fetchWorkflows()

      this.importProcessing = false
    },

    handleRowClicked (workflow) {
      this.$router.push({ name: 'workflow.edit', params: { workflowID: workflow.workflowID } })
    },
  },
}
</script>
