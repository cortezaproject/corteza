<template>
  <WorkflowEditor
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
    @change-detected="onChangeDetected"
    @delete="deleteWorkflow"
    @undelete="undeleteWorkflow"
  />
</template>

<script setup>
import WorkflowEditor from '@/components/WorkflowEditor.vue'
import { automation } from '@cortezaproject/corteza-js-next'
import { throttle } from 'lodash-es'
import { useRBACStore, useUnsavedGuard } from '@cortezaproject/corteza-vue-next'
import { computed, inject, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { useWorkflowStore } from '@/stores/workflow'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const toast = useToast()
const workflowStore = useWorkflowStore()

const $AutomationAPI = inject('$AutomationAPI')
const $Auth = inject('$Auth')

// State
const processing = ref(true)
const processingSave = ref(false)
const processingDelete = ref(false)
const workflow = ref({})
const triggers = ref([])
const changeDetected = ref(false)

// Computed
const rbacStore = useRBACStore()
const canCreate = rbacStore.can('automation/', 'workflow.create')

const workflowID = computed(() => {
  return route.params.workflowID || (workflow.value.workflowID !== '0' ? workflow.value.workflowID : undefined)
})

const userID = computed(() => {
  return $Auth?.user?.userID
})

useUnsavedGuard({
  isDirty: () => changeDetected.value && !workflow.value.deletedAt,
  messageKey: 'general.editor.unsavedChanges',
})

// Lifecycle
onMounted(async () => {
  if (workflowID.value) {
    await fetchTriggers()
    await fetchWorkflow()
  } else {
    workflow.value = new automation.Workflow({
      ownedBy: userID.value,
      runAs: '0',
      enabled: true,
      handle: '',
    })
  }

  processing.value = false
})

// Methods
async function fetchWorkflow() {
  try {
    const wf = await $AutomationAPI.workflowRead({ workflowID: workflowID.value })
    workflow.value = new automation.Workflow(wf)
  } catch (e) {
    toast.add({ severity: 'error', summary: t('notification.failed-fetch-workflow'), life: 5000 })
  }
}

async function fetchTriggers(wfID = workflowID.value) {
  try {
    const { set = [] } = await $AutomationAPI.triggerList({ workflowID: wfID, disabled: 1 })
    triggers.value = set
  } catch (e) {
    toast.add({ severity: 'error', summary: t('notification.failed-fetch-triggers'), life: 5000 })
  }
}

// Change detection
function onChangeDetected() {
  changeDetected.value = true
}

// Expose for WorkflowEditor to use
defineExpose({ onChangeDetected })

const saveWorkflow = throttle(async function (wf) {
  try {
    processingSave.value = true

    const isNew = wf.workflowID === '0'
    const { triggers: wfTriggers = [] } = wf

    // For new workflows, create the workflow first to get a real workflowID
    // before creating triggers — otherwise triggers are created with workflowID='0'
    if (isNew) {
      wf = await $AutomationAPI.workflowCreate(wf)
      workflowStore.updateInList(wf)
    }

    // Handle trigger updates - delete removed triggers, then create/update remaining
    await Promise.all(triggers.value.filter(({ triggerID }) => {
      return !wfTriggers.find(t => triggerID === t.triggerID)
    }).map(({ triggerID }) => {
      return $AutomationAPI.triggerDelete({ triggerID })
    })).then(async () => {
      await Promise.all(wfTriggers.map(t => {
        if (t.triggerID) {
          return $AutomationAPI.triggerUpdate({
            ...t,
            workflowStepID: t.stepID,
          })
        } else {
          return $AutomationAPI.triggerCreate({
            ...t,
            workflowID: wf.workflowID,
            workflowStepID: t.stepID,
            ownedBy: userID.value,
          })
        }
      })).catch(() => {
        throw new Error(t('notification.configure-triggers'))
      })
    })

    // For existing workflows, update after triggers are saved
    if (!isNew) {
      wf = await $AutomationAPI.workflowUpdate(wf)
      workflowStore.updateInList(wf)
    }

    // Refresh triggers
    await fetchTriggers(wf.workflowID)

    changeDetected.value = false
    window.onbeforeunload = null

    workflow.value = new automation.Workflow(wf)
    toast.add({ severity: 'success', summary: t('notification.update.success'), life: 3000 })

    if (isNew) {
      router.push({ name: 'workflow.edit', params: { workflowID: workflow.value.workflowID } })
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('notification.failed-save'), detail: e.message, life: 5000 })
  }

  processingSave.value = false
}, 500)

function deleteWorkflow() {
  if (workflow.value.workflowID) {
    processingDelete.value = true

    $AutomationAPI.workflowDelete(workflow.value)
      .then(() => {
        workflowStore.removeFromList(workflow.value.workflowID)
        workflow.value = {}
        workflow.value.deletedAt = new Date()
        router.push({ name: 'workflow.list' })
        toast.add({ severity: 'success', summary: t('notification.delete.success'), life: 3000 })
      })
      .catch(() => {
        toast.add({ severity: 'error', summary: t('notification.delete.failed'), life: 5000 })
      })
      .finally(() => {
        processingDelete.value = false
      })
  }
}

function undeleteWorkflow() {
  if (workflow.value.workflowID) {
    processingDelete.value = true

    $AutomationAPI.workflowUndelete(workflow.value)
      .then(() => {
        workflow.value.deletedAt = undefined
        workflow.value.deletedBy = undefined
        toast.add({ severity: 'success', summary: t('notification.undelete.success'), life: 3000 })
      })
      .catch(() => {
        toast.add({ severity: 'error', summary: t('notification.undelete.failed'), life: 5000 })
      })
      .finally(() => {
        processingDelete.value = false
      })
  }
}
</script>

<style scoped>
#workflow-editor :deep(tr.p-datatable-row-expansion > td) {
  padding-top: 0;
}
</style>
