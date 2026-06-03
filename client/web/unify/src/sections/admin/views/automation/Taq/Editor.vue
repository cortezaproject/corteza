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
    v-else-if="taq"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <div v-if="isEdit" class="flex justify-end gap-2 shrink-0">
        <CPermissionsButton
          v-if="taq.canGrant"
          v-tooltip.bottom="$t('general.label.permissions')"
          :resource="`corteza::automation:ng-automation/${taq.automationID}`"
          :title="taq.meta?.short || taq.handle || taq.automationID"
          :target="taq.meta?.short || taq.handle || taq.automationID"
        />
      </div>

      <Panel :header="$t('automation.taq.editor.info.title')" toggleable :collapsed="false" class="shadow">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <CFormGroup name="name" :label="$t('automation.taq.editor.info.name')" required>
            <InputText id="name" name="name" v-model="taq.meta.short" />
          </CFormGroup>

          <CFormGroup name="handle" :label="$t('automation.taq.editor.info.handle')">
            <InputText id="handle" name="handle" v-model="taq.handle" />
          </CFormGroup>

          <CFormGroup name="description" :label="$t('automation.taq.editor.info.description')" class="md:col-span-2">
            <Textarea
              id="description"
              name="description"
              v-model="taq.meta.description"
              rows="3"
            />
          </CFormGroup>

          <CFormGroup :label="$t('automation.taq.editor.info.enabled')">
            <ToggleSwitch v-model="taq.enabled" />
          </CFormGroup>
        </div>
      </Panel>

      <Panel
        v-if="isEdit"
        :header="$t('automation.taq.editor.meta.title')"
        toggleable
        :collapsed="true"
        class="shadow"
      >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="flex flex-col gap-1">
            <span class="text-xs text-muted-color">{{ $t('automation.taq.editor.meta.id') }}</span>
            <span class="font-mono text-sm">{{ taq.automationID }}</span>
          </div>

          <div class="flex flex-col gap-1">
            <span class="text-xs text-muted-color">{{ $t('automation.taq.editor.meta.createdAt') }}</span>
            <span>{{ locFullDateTime(taq.createdAt) }}</span>
          </div>

          <div v-if="taq.updatedAt" class="flex flex-col gap-1">
            <span class="text-xs text-muted-color">{{ $t('automation.taq.editor.meta.updatedAt') }}</span>
            <span>{{ locFullDateTime(taq.updatedAt) }}</span>
          </div>
        </div>
      </Panel>
    </div>

    <CEditorActions :back-to="{ name: 'automation.taq' }">
      <Button
        v-if="isEdit"
        :label="$t('automation.taq.editor.info.openBuilder')"
        icon="pi pi-external-link"
        severity="secondary"
        @click="openInBuilder"
      />
      <CInputDelete
        v-if="isEdit && taq.canDeleteNgAutomation && !taq.deletedAt"
        :label="$t('automation.taq.editor.info.delete')"
        :message="$t('general.confirm.delete')"
        :header="taq.meta?.short || taq.handle || taq.automationID"
        :disabled="deleting"
        @confirm="handleDelete"
      />
      <Button
        type="submit"
        :label="$t('general.label.save')"
        icon="pi pi-save"
        :loading="saving"
      />
    </CEditorActions>
  </Form>
</template>

<script setup>
import { computed, inject, nextTick, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { automation } from '@planetcrust/human-js'
import { components, filters, useUnsavedGuard } from '@planetcrust/human-vue'
import { cloneDeep, isEqual } from 'lodash-es'

const { CInputDelete } = components
const { locFullDateTime } = filters

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const $toast = inject('$toast')
const $AutomationAPI = inject('$AutomationAPI')

const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const taq = ref(null)
const initialTaq = ref(null)

const isEdit = computed(() => !!route.params.automationID)

const pageTitle = computed(() => {
  return isEdit.value
    ? t('automation.taq.editor.title.edit')
    : t('automation.taq.editor.title.create')
})

const initialValues = computed(() => ({
  name: taq.value?.meta?.short || '',
  handle: taq.value?.handle || '',
}))

const resolver = ref(({ values }) => {
  const errors = {}

  if (!values.name || values.name.trim().length === 0) {
    errors.name = [{ message: t('general.label.required') }]
  }

  if (values.handle && !/^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$|^[A-Za-z]$/.test(values.handle)) {
    errors.handle = [{ message: t('automation.taq.editor.info.invalidHandle') }]
  }

  return { errors }
})

function openInBuilder() {
  window.open(window.location.origin + '/taq/builder/' + taq.value.automationID, '_blank')
}

async function loadTaq() {
  const automationID = route.params.automationID
  if (!automationID) {
    taq.value = new automation.TAQ({ enabled: true, meta: { short: '' } })
    initialTaq.value = cloneDeep(taq.value)
    return
  }

  loading.value = true
  try {
    const raw = await $AutomationAPI.ngAutomationRead({ automationID })
    taq.value = new automation.TAQ(raw)
    initialTaq.value = cloneDeep(taq.value)
  } catch (e) {
    $toast.toastErrorHandler(t('notification.taq.fetch.error'))(e)
    router.push({ name: 'automation.taq' })
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
  if (isEdit.value && !taq.value?.canUpdateNgAutomation) return

  saving.value = true
  try {
    const payload = {
      handle: taq.value.handle,
      enabled: taq.value.enabled,
      meta: taq.value.meta,
    }

    if (isEdit.value) {
      payload.automationID = taq.value.automationID
      const raw = await $AutomationAPI.ngAutomationUpdate(payload)
      taq.value = new automation.TAQ(raw)
      initialTaq.value = cloneDeep(taq.value)
      $toast.toastSuccess(t('notification.taq.update.success'))
    } else {
      payload.triggers = []
      payload.steps = []
      payload.paths = []
      const created = await $AutomationAPI.ngAutomationCreate(payload)
      taq.value = new automation.TAQ(created)
      initialTaq.value = cloneDeep(taq.value)
      $toast.toastSuccess(t('notification.taq.create.success'))
      markSaved()
      router.push({ name: 'automation.taq.edit', params: { automationID: created.automationID } })
    }
  } catch (e) {
    $toast.toastErrorHandler(
      isEdit.value ? t('notification.taq.update.error') : t('notification.taq.create.error'),
    )(e)
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await $AutomationAPI.ngAutomationDelete({ automationID: taq.value.automationID })
    $toast.toastSuccess(t('notification.taq.delete.success'))
    router.push({ name: 'automation.taq' })
  } catch (e) {
    $toast.toastErrorHandler(t('notification.taq.delete.error'))(e)
  } finally {
    deleting.value = false
  }
}

const { markSaved } = useUnsavedGuard({
  isDirty: () =>
    !saving.value &&
    !deleting.value &&
    !!taq.value &&
    !!initialTaq.value &&
    !isEqual(taq.value, initialTaq.value),
  messageKey: 'general.editor.unsavedChanges',
})

watch(
  () => route.params.automationID,
  (newID, oldID) => {
    if (newID !== oldID) {
      loadTaq()
    }
  },
)

onMounted(() => {
  loadTaq()
})
</script>
