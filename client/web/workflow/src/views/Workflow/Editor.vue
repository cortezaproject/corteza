<template>
  <workflow-editor
    v-if="!processing"
    id="workflow-editor"
    :workflow-object="workflow"
    :workflow-triggers="triggers"
    :change-detected="changeDetected"
    :can-create="canCreate"
    :processing-save="processingSave"
    :processing-delete="processingDelete"
    class="overflow-hidden"
    @save="saveWorkflow"
    @delete="deleteWorkflow"
    @undelete="undeleteWorkflow"
  />
</template>

<script>
import WorkflowEditor from '../../components/WorkflowEditor'
import { automation } from '@cortezaproject/corteza-js'
import { throttle } from 'lodash'
import { mapGetters } from 'vuex'

export default {
  name: 'Editor',

  components: {
    WorkflowEditor,
  },

  beforeRouteLeave (to, from, next) {
    if (this.changeDetected && !this.workflow.deletedAt) {
      next(window.confirm(this.$t('notification:confirm-unsaved-changes')))
    } else {
      window.onbeforeunload = null
      next()
    }
  },

  data () {
    return {
      processing: true,
      processingSave: false,
      processingDelete: false,

      workflow: {},
      triggers: [],

      changeDetected: false,
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canCreate () {
      return this.can('automation/', 'workflow.create')
    },

    workflowID () {
      return this.$route.params.workflowID || (this.workflow.workflowID !== '0' ? this.workflow.workflowID : undefined)
    },

    userID () {
      if (this.$auth.user) {
        return this.$auth.user.userID
      }
      return undefined
    },
  },

  async mounted () {
    window.onbeforeunload = null

    this.$root.$on('change-detected', () => {
      if (!this.changeDetected) {
        window.onbeforeunload = () => {
          return true
        }
      }

      this.changeDetected = true
    })

    if (this.workflowID) {
      await this.fetchTriggers()
      await this.fetchWorkflow()
    } else {
      this.workflow = new automation.Workflow({
        ownedBy: this.userID,
        runAs: '0',
        enabled: true,
        handle: '',
      })
    }

    this.processing = false
  },

  beforeDestroy () {
    window.onbeforeunload = null
  },

  methods: {
    async fetchWorkflow () {
      return this.$AutomationAPI.workflowRead({ workflowID: this.workflowID })
        .then(wf => {
          this.workflow = new automation.Workflow(wf)
        })
        .catch(this.toastErrorHandler(this.$t('notification:failed-fetch-workflow')))
    },

    async fetchTriggers (workflowID = this.workflowID) {
      return this.$AutomationAPI.triggerList({ workflowID, disabled: 1 })
        .then(({ set = [] }) => {
          this.triggers = set
        })
        .catch(this.toastErrorHandler(this.$t('notification:failed-fetch-triggers')))
    },

    saveWorkflow: throttle(async function (wf) {
      // Issue #687: previously this handler would swallow trigger-update
      // failures behind a misleading "configure triggers" message and
      // could crash on `undefined.workflowID` if the workflow API call
      // rejected (e.g. backend unreachable). It now distinguishes
      // network failures from validation errors and surfaces the real
      // underlying error to the user.
      this.processingSave = true

      const isNetworkError = e => {
        if (!e) return false
        // axios populates `code` for network errors and leaves
        // `response` undefined; corteza API client rejects with the
        // raw response.data.error string for backend-returned errors,
        // so the absence of any backend payload is the signal.
        if (typeof e === 'string') return false
        if (e.response) return false
        if (e.message && /Network Error|ECONN|ETIMEDOUT|Failed to fetch/i.test(e.message)) return true
        // axios timeout / connection refused has no response field
        if (e.code && /ECONN|ETIMEDOUT|ENETUNREACH|ENOTFOUND/.test(e.code)) return true
        return false
      }

      const reportSaveError = e => {
        if (isNetworkError(e)) {
          this.toastDanger(
            this.$t('notification:failed-save-network'),
            this.$t('notification:failed-save'),
          )
          return
        }
        // Defer to the shared handler — it knows how to render workflow
        // error step rich payloads and falls back to plain toasts.
        this.toastErrorHandler(this.$t('notification:failed-save'))(e)
      }

      try {
        const isNew = wf.workflowID === '0'

        const { triggers = [] } = wf

        // Firstly handle trigger updates
        // Delete triggers of steps that were deleted
        try {
          await Promise.all(this.triggers.filter(({ triggerID }) => {
            return !triggers.find(t => triggerID === t.triggerID)
          }).map(({ triggerID }) => {
            return this.$AutomationAPI.triggerDelete({ triggerID })
          }))

          await Promise.all(triggers.map(t => {
            // Update triggers that already have an ID
            if (t.triggerID) {
              return this.$AutomationAPI.triggerUpdate({
                ...t,
                workflowStepID: t.stepID,
              })
            } else {
              // Create the other triggers
              return this.$AutomationAPI.triggerCreate({
                ...t,
                workflowID: wf.workflowID,
                workflowStepID: t.stepID,
                ownedBy: this.userID,
              })
            }
          }))
        } catch (e) {
          // Surface the actual trigger error instead of the generic
          // "configure triggers" message that used to mask it.
          reportSaveError(e)
          this.processingSave = false
          return
        }

        // Secondly handle workflow updates
        let saved
        try {
          saved = isNew
            ? await this.$AutomationAPI.workflowCreate(wf)
            : await this.$AutomationAPI.workflowUpdate(wf)
        } catch (e) {
          reportSaveError(e)
          this.processingSave = false
          return
        }

        if (!saved || !saved.workflowID) {
          // Defensive: API returned an unexpected shape. Don't crash on
          // saved.workflowID later down the chain.
          this.toastDanger(
            this.$t('notification:failed-save-unexpected-response'),
            this.$t('notification:failed-save'),
          )
          this.processingSave = false
          return
        }

        // Lastly update all of the bits
        await this.fetchTriggers(saved.workflowID)

        this.changeDetected = false
        window.onbeforeunload = null

        this.workflow = new automation.Workflow(saved)
        this.toastSuccess(this.$t('notification:update.success'))

        if (isNew) {
          // Redirect to edit route if new
          this.$router.push({ name: 'workflow.edit', params: { workflowID: this.workflow.workflowID } })
        }
      } catch (e) {
        reportSaveError(e)
      }

      this.processingSave = false
    }, 500),

    deleteWorkflow () {
      if (this.workflow.workflowID) {
        this.processingDelete = true

        this.$AutomationAPI.workflowDelete(this.workflow)
          .then(() => {
            // Disable unsaved changes prompt
            this.workflow = {}
            this.workflow.deletedAt = new Date()
            this.$router.push({ name: 'workflow.list' })
            this.toastSuccess(this.$t('notification:delete.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:delete.failed')))
          .finally(() => {
            this.processingDelete = false
          })
      }
    },

    undeleteWorkflow () {
      if (this.workflow.workflowID) {
        this.processingDelete = true

        this.$AutomationAPI.workflowUndelete(this.workflow)
          .then(() => {
            this.workflow.deletedAt = undefined
            this.workflow.deletedBy = undefined
            this.toastSuccess(this.$t('notification:undelete.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:undelete.failed')))
          .finally(() => {
            this.processingDelete = false
          })
      }
    },
  },
}
</script>

<style lang="scss">
#workflow-editor {
  tr.b-table-details > td {
    padding-top: 0;
  }

  .arrow-up {
    width: 0;
    height: 0;
    margin: 0 auto;
    border-left: 10px solid transparent;
    border-right: 10px solid transparent;
    border-bottom: 10px solid var(--light);
  }
}
</style>
