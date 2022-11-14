<template>
  <workflow-editor
    v-if="!processing"
    :workflow-object="workflow"
    :workflow-triggers="triggers"
    :change-detected="changeDetected"
    :can-create="canCreate"
    :processing-save="processingSave"
    class="overflow-hidden"
    @save="saveWorkflow"
    @delete="deleteWorkflow"
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

  data () {
    return {
      processing: true,
      processingSave: false,

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
        meta: {
          name: this.$t('general:unnamed-workflow'),
        },
      })
    }

    this.processing = false
  },

  beforeRouteLeave (to, from, next) {
    if (this.changeDetected) {
      next(window.confirm(this.$t('notification:confirm-unsaved-changes')))
    } else {
      window.onbeforeunload = null
      next()
    }
  },

  beforeDestroy () {
    window.onbeforeunload = null
  },

  methods: {
    async fetchWorkflow () {
      return this.$AutomationAPI.workflowRead({ workflowID: this.workflowID })
        .then(wf => {
          this.workflow = wf
        })
        .catch(this.defaultErrorHandler(this.$t('notification:failed-fetch-workflow')))
    },

    async fetchTriggers (workflowID = this.workflowID) {
      return this.$AutomationAPI.triggerList({ workflowID, disabled: 1 })
        .then(({ set = [] }) => {
          this.triggers = set
        })
        .catch(this.defaultErrorHandler(this.$t('notification:failed-fetch-triggers')))
    },

    saveWorkflow: throttle(async function (wf) {
      try {
        this.processingSave = true

        const isNew = wf.workflowID === '0'

        const { steps = [], paths = [], triggers = [] } = wf
        this.workflow.steps = steps
        this.workflow.paths = paths

        // Firstly handle trigger updates
        // Delete triggers of steps that were deleted
        await Promise.all(this.triggers.filter(({ triggerID }) => {
          return !triggers.find(t => triggerID === t.triggerID)
        }).map(({ triggerID }) => {
          return this.$AutomationAPI.triggerDelete({ triggerID })
        }),
        ).then(async () => {
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
          })).catch(() => {
            throw new Error(this.$t('notification:configure-triggers'))
          })
        })

        // Secondly handle workflow updates
        if (isNew) {
          wf = await this.$AutomationAPI.workflowCreate(this.workflow)
        } else {
          wf = await this.$AutomationAPI.workflowUpdate(this.workflow)
        }

        // Lastly update all of the bits
        await this.fetchTriggers(wf.workflowID)
        this.changeDetected = false
        window.onbeforeunload = null

        this.workflow = wf
        this.raiseSuccessAlert(this.$t('notification:updated-workflow'))

        if (isNew) {
          // Redirect to edit route if new
          this.$router.push({ name: 'workflow.edit', params: { workflowID: this.workflow.workflowID } })
        }
      } catch (e) {
        this.defaultErrorHandler(this.$t('notification:failed-save'))(e)
      }

      this.processingSave = false
    }, 500),

    deleteWorkflow () {
      if (this.workflow.workflowID) {
        this.$AutomationAPI.workflowDelete(this.workflow)
          .then(() => {
            // Disable unsaved changes prompt
            this.workflow = {}
            this.$router.push({ name: 'workflow.list' })

            this.raiseSuccessAlert(this.$t('notification:deleted-workflow'))
          })
          .catch(this.defaultErrorHandler(this.$t('notification:delete-failed')))
      }
    },
  },
}
</script>
