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
      <Panel
        :header="$t('automation.workflows.editor.info.title', 'Basic information')"
        toggleable
        :collapsed="false"
      >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField name="name" class="flex flex-col gap-2">
            <label for="name" class="font-medium text-primary">
              {{ $t('automation.workflows.editor.info.name', 'Workflow name') }} *
            </label>
            <InputText id="name" name="name" v-model="workflow.meta.name" />
            <Message v-if="$form.name?.invalid" severity="error" size="small" variant="simple">
              {{ $form.name.error?.message }}
            </Message>
          </FormField>

          <FormField name="handle" class="flex flex-col gap-2">
            <label for="handle" class="font-medium text-primary">
              {{ $t('automation.workflows.editor.info.handle', 'Handle') }}
            </label>
            <InputText id="handle" name="handle" v-model="workflow.handle" />
            <Message v-if="$form.handle?.invalid" severity="error" size="small" variant="simple">
              {{ $form.handle.error?.message }}
            </Message>
          </FormField>

          <FormField name="description" class="flex flex-col gap-2 md:col-span-2">
            <label for="description" class="font-medium text-primary">
              {{ $t('automation.workflows.editor.info.description', 'Description') }}
            </label>
            <Textarea
              id="description"
              name="description"
              v-model="workflow.meta.description"
              rows="3"
            />
          </FormField>

          <div class="flex flex-col gap-2">
            <label class="font-medium text-primary">
              {{ $t('automation.workflows.editor.info.enabled', 'Enabled') }}
            </label>
            <ToggleSwitch v-model="workflow.enabled" />
          </div>

          <div class="flex flex-col gap-2">
            <label class="font-medium text-primary">
              {{ $t('automation.workflows.editor.info.trace', 'Trace') }}
            </label>
            <ToggleSwitch v-model="workflow.trace" />
            <span class="text-xs text-surface-500">
              {{ $t('automation.workflows.editor.info.traceHint', 'Enable execution tracing for debugging') }}
            </span>
          </div>
        </div>
      </Panel>
    </div>

    <!-- Bottom Actions Toolbar -->
    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="p-3 flex items-center justify-between">
        <Button
          :label="$t('general.label.back')"
          icon="pi pi-arrow-left"
          severity="secondary"
          @click="$router.push({ name: 'automation.workflows' })"
        />
        <div class="flex gap-2">
          <CInputDelete
            v-if="isEdit && workflow.canDeleteWorkflow && !workflow.deletedAt"
            :label="$t('automation.workflows.editor.info.delete', 'Delete')"
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
        </div>
      </div>
    </div>
  </Form>
</template>

<script setup>
import { computed, inject, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { automation } from '@cortezaproject/corteza-js-next'
import { components } from '@cortezaproject/corteza-vue-next'

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

const isEdit = computed(() => !!route.params.workflowID)

const pageTitle = computed(() => {
  return isEdit.value
    ? t('automation.workflows.editor.title.edit', 'Edit workflow')
    : t('automation.workflows.editor.title.create', 'Create workflow')
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
    errors.handle = [{ message: t('automation.workflows.editor.info.invalid-handle-characters', 'Invalid handle characters') }]
  }

  return { errors }
})

async function loadWorkflow() {
  const workflowID = route.params.workflowID
  if (!workflowID) {
    workflow.value = new automation.Workflow({ enabled: true, trace: false, meta: {} })
    return
  }

  loading.value = true
  try {
    const raw = await $AutomationAPI.workflowRead({ workflowID })
    workflow.value = new automation.Workflow(raw)
  } catch (e) {
    $toast.toastErrorHandler(t('notification.workflow.fetch.error', 'Failed to load workflow'))(e)
    router.push({ name: 'automation.workflows' })
  } finally {
    loading.value = false
  }
}

async function handleSubmit({ valid }) {
  if (!valid) return
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
      $toast.toastSuccess(t('notification.workflow.update.success', 'Workflow updated'))
    } else {
      const created = await $AutomationAPI.workflowCreate(payload)
      $toast.toastSuccess(t('notification.workflow.create.success', 'Workflow created'))
      router.push({ name: 'automation.workflows.edit', params: { workflowID: created.workflowID } })
    }
  } catch (e) {
    $toast.toastErrorHandler(
      isEdit.value
        ? t('notification.workflow.update.error', 'Failed to update workflow')
        : t('notification.workflow.create.error', 'Failed to create workflow'),
    )(e)
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await $AutomationAPI.workflowDelete({ workflowID: workflow.value.workflowID })
    $toast.toastSuccess(t('notification.workflow.delete.success', 'Workflow deleted'))
    router.push({ name: 'automation.workflows' })
  } catch (e) {
    $toast.toastErrorHandler(t('notification.workflow.delete.error', 'Failed to delete workflow'))(e)
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  loadWorkflow()
})
</script>
