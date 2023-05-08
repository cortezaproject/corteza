<template>
  <b-container
    v-if="workflow"
    class="py-3"
  >
    <c-content-header
      :title="title"
    >
      <span
        class="text-nowrap"
      >
        <b-button
          v-if="workflowID && canCreate"
          variant="primary"
          :to="{ name: 'automation.workflow.new' }"
        >
          {{ $t('new') }}
        </b-button>
        <c-permissions-button
          v-if="workflowID && canGrant"
          :title="workflow.meta.name || workflow.handle || workflowID"
          :target="workflow.meta.name || workflow.handle || workflowID"
          :resource="`corteza::automation:workflow/${workflowID}`"
          button-variant="light"
          class="ml-2"
        >
          <font-awesome-icon :icon="['fas', 'lock']" />
          {{ $t('permissions') }}
        </c-permissions-button>
      </span>
    </c-content-header>

    <c-workflow-editor-info
      :workflow="workflow"
      :processing="info.processing"
      :success="info.success"
      :can-create="canCreate"
      @submit="onInfoSubmit"
      @delete="onDelete"
    />

    <c-workflow-editor-triggers
      v-if="workflowID"
      :triggers="triggers"
      :processing="info.processing"
      :success="info.success"
    />
  </b-container>
</template>
<script>
import { isEqual, cloneDeep } from 'lodash'
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CWorkflowEditorInfo from 'corteza-webapp-admin/src/components/Workflow/CWorkflowEditorInfo'
import CWorkflowEditorTriggers from 'corteza-webapp-admin/src/components/Workflow/CWorkflowEditorTriggers'
import { mapGetters } from 'vuex'

export default {
  components: {
    CWorkflowEditorInfo,
    CWorkflowEditorTriggers,
  },

  i18nOptions: {
    namespaces: 'automation.workflows',
    keyPrefix: 'editor',
  },

  mixins: [
    editorHelpers,
  ],

  props: {
    workflowID: {
      type: String,
      required: false,
      default: undefined,
    },
  },

  data () {
    return {
      workflow: undefined,
      initialWorkflowState: undefined,
      triggers: [],

      info: {
        processing: false,
        success: false,
      },
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canCreate () {
      return this.can('automation/', 'workflow.create')
    },

    canGrant () {
      return this.can('automation/', 'grant')
    },

    userID () {
      if (this.$auth.user) {
        return this.$auth.user.userID
      }
      return undefined
    },

    title () {
      return this.workflowID ? this.$t('title.edit') : this.$t('title.create')
    },
  },

  beforeRouteUpdate (to, from, next) {
    this.checkUnsavedChanges(next, to)
  },

  beforeRouteLeave (to, from, next) {
    this.checkUnsavedChanges(next, to)
  },

  watch: {
    workflowID: {
      immediate: true,
      handler () {
        if (this.workflowID) {
          this.fetchWorkflow()
          this.fetchTriggers()
        } else {
          this.workflow = {
            ownedBy: this.userID,
            runAs: this.userID,
            enabled: true,
            meta: {
              name: '',
            },
          }

          this.initialWorkflowState = cloneDeep(this.workflow)
        }
      },
    },
  },

  methods: {
    fetchWorkflow () {
      this.incLoader()

      this.$AutomationAPI.workflowRead({ workflowID: this.workflowID })
        .then(this.prepare)
        .catch(this.toastErrorHandler(this.$t('notification:workflow.fetch.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    fetchTriggers () {
      this.incLoader()

      this.$AutomationAPI.triggerList({ workflowID: this.workflowID, disabled: 1 })
        .then(({ set = [] }) => { this.triggers = set })
        .catch(this.toastErrorHandler(this.$t('notification:workflow.trigger.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    onInfoSubmit (workflow) {
      this.info.processing = true

      if (this.workflowID) {
        this.$AutomationAPI.workflowUpdate(workflow)
          .then(() => {
            this.fetchWorkflow()

            this.animateSuccess('info')
            this.toastSuccess(this.$t('notification:workflow.update.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:workflow.update.error')))
          .finally(() => {
            this.info.processing = false
          })
      } else {
        this.$AutomationAPI.workflowCreate(workflow)
          .then(({ workflowID }) => {
            this.animateSuccess('info')
            this.toastSuccess(this.$t('notification:workflow.create.success'))

            this.$router.push({ name: 'automation.workflow.edit', params: { workflowID } })
          })
          .catch(this.toastErrorHandler(this.$t('notification:workflow.create.error')))
          .finally(() => {
            this.info.processing = false
          })
      }
    },

    onDelete () {
      this.incLoader()

      if (this.workflow.deletedAt) {
        this.$AutomationAPI.workflowUndelete({ workflowID: this.workflowID })
          .then(() => {
            this.fetchWorkflow()

            this.toastSuccess(this.$t('notification:workflow.undelete.success'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:workflow.undelete.error')))
          .finally(() => {
            this.decLoader()
          })
      } else {
        this.$AutomationAPI.workflowDelete({ workflowID: this.workflowID })
          .then(() => {
            this.fetchWorkflow()

            this.toastSuccess(this.$t('notification:workflow.delete.success'))
            this.$router.push({ name: 'automation.workflow' })
          })
          .catch(this.toastErrorHandler(this.$t('notification:workflow.delete.error')))
          .finally(() => {
            this.decLoader()
          })
      }
    },

    prepare (workflow = {}) {
      this.workflow = workflow
      this.initialWorkflowState = cloneDeep(this.workflow)
    },

    checkUnsavedChanges (next, to) {
      const isNewPage = this.$route.path.includes('/new') && to.name.includes('edit')

      if (isNewPage) {
        next(true)
      } else if (!to.name.includes('edit')) {
        next(!isEqual(this.workflow, this.initialWorkflowState) ? window.confirm(this.$t('general:editor.unsavedChanges')) : true)
      }
    },
  },
}
</script>
