<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ pageTitle }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <Form
    v-else-if="application"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <CViewContainer scroll>
      <div v-if="isEdit" class="flex justify-end gap-2 shrink-0">
        <CPermissionsButton
          v-if="application.canGrant"
          v-tooltip.bottom="$t('general.label.permissions')"
          :resource="`corteza::system:application/${application.applicationID}`"
          :title="application.name || application.applicationID"
          :target="application.name || application.applicationID"
        />
      </div>
      <Panel :header="$t('system.applications.editor.info.title')" toggleable :collapsed="false">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <CFormGroup
            name="name"
            :label="$t('system.applications.editor.info.name')"
            required
            class="md:col-span-2"
          >
            <InputText id="name" name="name" v-model="application.name" />
          </CFormGroup>

          <CInputToggleCard
            v-model="application.enabled"
            :label="$t('system.applications.editor.info.enabled')"
            :description="$t('system.applications.editor.info.enabledDescription')"
          />
        </div>
      </Panel>

      <Panel :header="$t('system.applications.editor.unify.title')" toggleable :collapsed="false">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <CFormGroup
            :label="$t('system.applications.editor.unify.name.label')"
            input-id="unifyName"
          >
            <InputText id="unifyName" v-model="application.unify.name" />
          </CFormGroup>

          <CFormGroup :label="$t('system.applications.editor.unify.url.label')" input-id="unifyUrl">
            <InputText id="unifyUrl" v-model="application.unify.url" />
          </CFormGroup>

          <CInputToggleCard
            v-model="application.unify.listed"
            :label="$t('system.applications.editor.unify.listed')"
            :description="$t('system.applications.editor.unify.listedDescription')"
            class="self-start"
          />

          <div class="flex flex-col gap-2">
            <CFileDropZone
              accept="image/*"
              :uploading="logoUploading"
              :error="logoError"
              :preview-url="logoPreviewUrl"
              :clearable="isCustomLogo"
              :drop-label="$t('system.applications.editor.unify.logo.placeholder')"
              :label="$t('system.applications.editor.unify.logo.label')"
              compact
              preview-max-width="100%"
              preview-max-height="200px"
              @select="onLogoSelect"
              @clear="onLogoClear"
            />
          </div>
        </div>
      </Panel>
    </CViewContainer>

    <CEditorActions :back-to="{ name: 'system.applications' }">
      <CInputDelete
        v-if="isEdit && application.canDeleteApplication"
        :label="$t('system.applications.editor.info.delete')"
        :message="$t('general.confirm.delete')"
        :header="application.name"
        :disabled="deleting"
        @confirm="handleDelete"
      />
      <Button type="submit" :label="$t('general.label.save')" icon="pi pi-save" :loading="saving" />
    </CEditorActions>
  </Form>
</template>

<script setup>
import { computed, inject, nextTick, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { system } from '@planetcrust/human-js'
import {
  components,
  useFileUpload,
  resolveAppLogoUrl,
  useDraftGuard,
  useApplicationsStore,
} from '@planetcrust/human-vue'
import { appIconMap } from '@/utils/appIcons'

const { CFileDropZone, CInputDelete, CInputToggleCard, CViewContainer } = components

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')
const applicationsStore = useApplicationsStore()

const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const application = ref(null)
const { capture, markSaved } = useDraftGuard({
  draft: application,
  busy: () => saving.value || deleting.value,
})

// Logo upload
const {
  uploading: logoUploading,
  uploadError: logoUploadError,
  uploadFiles: uploadLogoFiles,
  reset: resetLogoUpload,
} = useFileUpload()

const logoError = computed(() => logoUploadError.value)

const logoPreviewUrl = computed(() => {
  if (!application.value) return ''
  return resolveAppLogoUrl(application.value, $SystemAPI.baseURL, appIconMap)
})

const isCustomLogo = computed(() => {
  const logo = application.value?.unify?.logo || application.value?.unify?.icon || ''
  if (!logo) return false
  // Built-in default icons should not be clearable
  return !appIconMap[logo]
})

const isEdit = computed(() => !!route.params.applicationID)

const pageTitle = computed(() =>
  isEdit.value
    ? t('system.applications.editor.title.edit')
    : t('system.applications.editor.title.create'),
)

const initialValues = computed(() => ({
  name: application.value?.name || '',
}))

const resolver = ref(({ values }) => {
  const errors = {}

  if (!values.name || values.name.trim().length === 0) {
    errors.name = [{ message: t('general.label.required') }]
  }

  return { errors }
})

async function loadApplication() {
  const applicationID = route.params.applicationID
  if (!applicationID) {
    application.value = new system.Application({ enabled: true })
    capture()
    return
  }

  loading.value = true
  try {
    const raw = await applicationsStore.findByID(applicationID)
    application.value = new system.Application(raw)
    capture()
  } catch (e) {
    $toast.toastErrorHandler(t('notification.application.fetch.error'))(e)
    router.push({ name: 'system.applications' })
  } finally {
    loading.value = false
  }
}

async function handleSubmit({ valid }) {
  if (!valid) {
    $toast.toastWarning(t('general.notification.formErrors'))
    nextTick(() => {
      document
        .querySelector('.p-message-error')
        ?.scrollIntoView({ behavior: 'smooth', block: 'center' })
    })
    return
  }

  saving.value = true
  try {
    const payload = {
      name: application.value.name,
      enabled: application.value.enabled,
      weight: application.value.weight,
      unify: application.value.unify,
    }

    if (isEdit.value) {
      payload.applicationID = application.value.applicationID
      const raw = await applicationsStore.update(payload)
      application.value = new system.Application(raw)
      capture()
      $toast.toastSuccess(t('notification.application.update.success'))
    } else {
      const created = await applicationsStore.create(payload)
      $toast.toastSuccess(t('notification.application.create.success'))
      markSaved()
      router.push({
        name: 'system.applications.edit',
        params: { applicationID: created.applicationID },
      })
    }
  } catch (e) {
    $toast.toastErrorHandler(
      t(
        `notification.application.${isEdit.value ? 'update' : 'create'}.error`,
        'Failed to save application',
      ),
    )(e)
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await applicationsStore.delete(application.value.applicationID)
    $toast.toastSuccess(t('notification.application.delete.success'))
    router.push({ name: 'system.applications' })
  } catch (e) {
    $toast.toastErrorHandler(t('notification.application.delete.error'))(e)
  } finally {
    deleting.value = false
  }
}

async function onLogoSelect(files) {
  const file = files[0]
  if (!file) return

  try {
    const results = await uploadLogoFiles([file], {
      api: $SystemAPI,
      endpoint: $SystemAPI.applicationUploadEndpoint(),
    })

    const rsp = results[0]
    if (rsp) {
      application.value.unify.logo = $SystemAPI.baseURL + rsp.url
      application.value.unify.logoID = rsp.attachmentID
    }
  } catch {
    // error is set by composable
  }
}

function onLogoClear() {
  application.value.unify.logo = ''
  application.value.unify.logoID = '0'
  resetLogoUpload()
}

onMounted(() => loadApplication())
watch(
  () => route.params.applicationID,
  () => loadApplication(),
)
</script>
