<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ pageTitle }}</span>
  </Teleport>

  <!-- Loading -->
  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <!-- Form -->
  <Form
    v-else-if="workflow"
    v-slot="$form"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <div v-if="isEdit" class="flex justify-end gap-2 shrink-0">
        <CPermissionsButton
          v-if="workflow.canGrant"
          v-tooltip.bottom="$t('general.label.permissions')"
          :resource="`corteza::automation:workflow/${workflow.workflowID}`"
          :title="workflow.meta?.name || workflow.handle || workflow.workflowID"
          :target="workflow.meta?.name || workflow.handle || workflow.workflowID"
        />
      </div>
      <Panel :header="$t('automation.workflows.editor.info.title')" toggleable :collapsed="false" class="shadow">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <CFormGroup name="name" :label="$t('automation.workflows.editor.info.name')" required>
            <InputText id="name" name="name" v-model="workflow.meta.name" />
          </CFormGroup>

          <CFormGroup name="handle" :label="$t('automation.workflows.editor.info.handle')">
            <InputText id="handle" name="handle" v-model="workflow.handle" />
          </CFormGroup>

          <CFormGroup name="description" :label="$t('automation.workflows.editor.info.description')" class="md:col-span-2">
            <Textarea
              id="description"
              name="description"
              v-model="workflow.meta.description"
              rows="3"
            />
          </CFormGroup>

          <CFormGroup :label="$t('automation.workflows.editor.info.enabled')">
            <ToggleSwitch v-model="workflow.enabled" />
          </CFormGroup>

          <CFormGroup
            :label="$t('automation.workflows.editor.info.trace')"
            :description="$t('automation.workflows.editor.info.traceHint')"
          >
            <ToggleSwitch v-model="workflow.trace" />
          </CFormGroup>
        </div>
      </Panel>

      <Panel
        v-if="isEdit"
        :header="$t('automation.workflows.editor.triggers.title')"
        toggleable
        class="shadow"
      >
        <WorkflowTriggers :triggers="triggers" />
      </Panel>
    </div>

    <CEditorActions :back-to="{ name: 'automation.workflows' }">
      <CInputDelete
        v-if="isEdit && workflow.canDeleteWorkflow && !workflow.deletedAt"
        :label="$t('automation.workflows.editor.info.delete')"
        :message="$t('general.confirm.delete')"
        :header="workflow.meta?.name || workflow.handle || workflow.workflowID"
        :disabled="deleting"
        @confirm="handleDelete"
      />
      <Button
        v-if="!isEdit || workflow.canUpdateWorkflow"
        type="submit"
        :label="$t('general.label.save')"
        icon="pi pi-save"
        :loading="saving"
      />
    </CEditorActions>
  </Form>
</template>

<script setup>
import { computed, inject, nextTick, onMounted, ref } from 'vue'
import WorkflowTriggers from '@/components/Workflow/WorkflowTriggers.vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { automation } from '@planetcrust/human-js'
import { components, useUnsavedGuard } from '@planetcrust/human-vue'
import { cloneDeep, isEqual } from 'lodash-es'

const { CInputDelete } = components

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const $toast = inject('$toast')
const $AutomationAPI = inject('$AutomationAPI')

const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const workflow = ref(null)
const initialWorkflow = ref(null)
const triggers = ref([])

const isEdit = computed(() => !!route.params.workflowID)

const pageTitle = computed(() => {
  return isEdit.value
    ? t('automation.workflows.editor.title.edit')
    : t('automation.workflows.editor.title.create')
})

const initialValues = computed(() => ({
  name: workflow.value?.meta?.name || '',
  handle: workflow.value?.handle || '',
}))

const resolver = ref(({ values }) => {
  const errors = {}

  if (!values.name || values.name.trim().length === 0) {
    errors.name = [{ message: t('general.label.required') }]
  }

  if (values.handle && !/^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$|^[A-Za-z]$/.test(values.handle)) {
    errors.handle = [{ message: t('automation.workflows.editor.info.invalid-handle-characters') }]
  }

  return { errors }
})

async function loadWorkflow() {
  const workflowID = route.params.workflowID
  if (!workflowID) {
    workflow.value = new automation.Workflow({ enabled: true, trace: false, meta: {} })
    initialWorkflow.value = cloneDeep(workflow.value)
    return
  }

  loading.value = true
  try {
    const raw = await $AutomationAPI.workflowRead({ workflowID })
    workflow.value = new automation.Workflow(raw)
    initialWorkflow.value = cloneDeep(workflow.value)

    // Load triggers for the workflow
    const triggersResult = await $AutomationAPI.triggerList({ workflowID })
    triggers.value = triggersResult?.set || []
  } catch (e) {
    $toast.toastErrorHandler(t('notification.workflow.fetch.error'))(e)
    router.push({ name: 'automation.workflows' })
  } finally {
    loading.value = false
  }
}

async function handleSubmit({ valid }) {
  if (!valid) {
    $toast.toastWarning(t('general.notification.formErrors'))
    nextTick(() => {
      document.querySelector('.p-message-error')?.scrollIntoView({ behavior: 'smooth', block: 'center' })
    })
    return
  }
  if (isEdit.value && !workflow.value?.canUpdateWorkflow) return

  saving.value = true
  try {
    const payload = {
      handle: workflow.value.handle,
      enabled: workflow.value.enabled,
      trace: workflow.value.trace,
      meta: workflow.value.meta,
    }

    if (isEdit.value) {
      payload.workflowID = workflow.value.workflowID
      const raw = await $AutomationAPI.workflowUpdate(payload)
      workflow.value = new automation.Workflow(raw)
      initialWorkflow.value = cloneDeep(workflow.value)
      $toast.toastSuccess(t('notification.workflow.update.success'))
    } else {
      const created = await $AutomationAPI.workflowCreate(payload)
      $toast.toastSuccess(t('notification.workflow.create.success'))
      markSaved()
      router.push({ name: 'automation.workflows.edit', params: { workflowID: created.workflowID } })
    }
  } catch (e) {
    $toast.toastErrorHandler(
      isEdit.value
        ? t('notification.workflow.update.error')
        : t('notification.workflow.create.error'),
    )(e)
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await $AutomationAPI.workflowDelete({ workflowID: workflow.value.workflowID })
    $toast.toastSuccess(t('notification.workflow.delete.success'))
    router.push({ name: 'automation.workflows' })
  } catch (e) {
    $toast.toastErrorHandler(t('notification.workflow.delete.error'))(e)
  } finally {
    deleting.value = false
  }
}

const { markSaved } = useUnsavedGuard({
  isDirty: () => !saving.value && !deleting.value && !!workflow.value && !!initialWorkflow.value && !isEqual(workflow.value, initialWorkflow.value),
  messageKey: 'general.editor.unsavedChanges',
})

onMounted(() => {
  loadWorkflow()
})
</script>
