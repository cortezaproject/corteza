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
            <InputText id="name" name="name" v-model="application.name" :disabled="!canEdit" />
          </CFormGroup>

          <CFormGroup
            :label="$t('system.applications.editor.info.description')"
            input-id="description"
            class="md:col-span-2"
          >
            <Textarea
              id="description"
              v-model="application.meta.description"
              rows="3"
              autoResize
              :disabled="!canEdit"
            />
          </CFormGroup>

          <CInputToggleCard
            v-model="application.enabled"
            :label="$t('system.applications.editor.info.enabled')"
            :description="$t('system.applications.editor.info.enabledDescription')"
            :disabled="!canEdit"
          />
        </div>
      </Panel>

      <Panel :header="$t('system.applications.editor.unify.title')" toggleable :collapsed="false">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <CFormGroup
            :label="$t('system.applications.editor.unify.kind.label')"
            :description="$t('system.applications.editor.unify.kind.description')"
            input-id="unifyKind"
            class="md:col-span-2"
          >
            <SelectButton
              id="unifyKind"
              v-model="application.unify.kind"
              :options="kindOptions"
              option-label="label"
              option-value="value"
              :allow-empty="false"
              :disabled="!canEdit"
            />
          </CFormGroup>

          <CFormGroup
            :label="$t('system.applications.editor.unify.name.label')"
            input-id="unifyName"
          >
            <InputText id="unifyName" v-model="application.unify.name" :disabled="!canEdit" />
          </CFormGroup>

          <CFormGroup
            :label="$t('system.applications.editor.unify.url.label')"
            :description="
              isCustom ? $t('system.applications.editor.unify.url.customDescription') : ''
            "
            input-id="unifyUrl"
          >
            <InputText
              v-if="isCustom"
              id="unifyUrl"
              :model-value="customUrl"
              :placeholder="$t('system.applications.editor.unify.url.customPlaceholder')"
              disabled
            />
            <InputText v-else id="unifyUrl" v-model="application.unify.url" :disabled="!canEdit" />
          </CFormGroup>

          <CInputToggleCard
            v-model="application.unify.listed"
            :label="$t('system.applications.editor.unify.listed')"
            :description="$t('system.applications.editor.unify.listedDescription')"
            class="self-start"
            :disabled="!canEdit"
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
      <Panel
        v-if="isEdit && isCustom"
        :header="$t('system.applications.editor.custom.title')"
        toggleable
        :collapsed="false"
      >
        <div class="flex flex-col gap-4">
          <div class="flex flex-wrap items-center gap-2">
            <Button
              :label="$t('system.applications.editor.custom.open')"
              icon="pi pi-external-link"
              severity="secondary"
              size="small"
              data-test-id="button-open-custom-app"
              @click="router.push(`/${customUrl}`)"
            />
            <span v-if="!application.enabled" class="text-sm text-muted-color">
              {{ $t('system.applications.editor.custom.previewNote') }}
            </span>
          </div>

          <dl class="grid grid-cols-1 md:grid-cols-2 gap-x-4 gap-y-2 text-sm">
            <div>
              <dt class="text-muted-color">
                {{ $t('system.applications.editor.custom.namespace') }}
              </dt>
              <dd>{{ sourceMeta.namespace || $t('system.applications.editor.custom.none') }}</dd>
            </div>
            <div>
              <dt class="text-muted-color">
                {{ $t('system.applications.editor.custom.modules') }}
              </dt>
              <dd class="flex flex-wrap gap-1">
                <Tag
                  v-for="module in sourceMeta.modules || []"
                  :key="module"
                  :value="module"
                  severity="secondary"
                />
                <span v-if="!(sourceMeta.modules || []).length">
                  {{ $t('system.applications.editor.custom.none') }}
                </span>
              </dd>
            </div>
            <div>
              <dt class="text-muted-color">{{ $t('system.applications.editor.custom.writes') }}</dt>
              <dd class="flex flex-wrap gap-1">
                <Tag
                  v-for="module in sourceMeta.writes || []"
                  :key="module"
                  :value="module"
                  severity="warn"
                />
                <span v-if="!(sourceMeta.writes || []).length">
                  {{ $t('system.applications.editor.custom.readOnly') }}
                </span>
              </dd>
            </div>
            <div>
              <dt class="text-muted-color">{{ $t('system.applications.editor.custom.size') }}</dt>
              <dd>{{ sourceSizeLabel }}</dd>
            </div>
            <div>
              <dt class="text-muted-color">
                {{ $t('system.applications.editor.custom.updated') }}
              </dt>
              <dd>{{ sourceUpdatedLabel }}</dd>
            </div>
          </dl>

          <Message v-if="sourceProblem" severity="warn" :closable="false">
            {{ sourceProblem }}
          </Message>
          <Message v-else-if="!source && !canEditSource" severity="info" :closable="false">
            {{ $t('system.applications.editor.custom.empty') }}
          </Message>
          <div v-else class="flex flex-col gap-2">
            <div class="flex flex-wrap items-center justify-between gap-2">
              <span class="text-sm text-muted-color">
                {{
                  canEditSource
                    ? $t('system.applications.editor.custom.edit')
                    : $t('system.applications.editor.custom.source')
                }}
              </span>
              <div class="flex items-center gap-1">
                <Button
                  :label="$t('system.applications.editor.custom.copy')"
                  icon="pi pi-copy"
                  severity="secondary"
                  variant="text"
                  size="small"
                  @click="copySource"
                />
                <Button
                  v-if="canEditSource"
                  :label="$t('system.applications.editor.custom.savePage')"
                  icon="pi pi-save"
                  size="small"
                  :loading="savingSource"
                  :disabled="draftSource === source"
                  data-test-id="button-save-page"
                  @click="handleSaveSource"
                />
              </div>
            </div>
            <CCodeEditor
              id="customSource"
              v-model="draftSource"
              language="html"
              min-height="360px"
              :read-only="!canEditSource"
            />
            <small class="text-muted-color">
              {{ $t('system.applications.editor.custom.editNote') }}
            </small>
          </div>
        </div>
      </Panel>

      <Message v-if="!canEdit" severity="warn" :closable="false">
        {{ $t('general.editor.readOnly') }}
      </Message>
    </CViewContainer>

    <CEditorActions :back-to="{ name: 'system.applications' }">
      <CInputDelete
        v-if="isEdit && application.canDeleteApplication && !application.deletedAt"
        :label="$t('system.applications.editor.info.delete')"
        :message="$t('general.confirm.delete')"
        :header="application.name"
        :disabled="deleting"
        @confirm="handleDelete"
      />
      <Button
        v-if="isEdit && application.canDeleteApplication && application.deletedAt"
        :label="$t('general.label.restore')"
        icon="pi pi-replay"
        severity="warn"
        :loading="restoring"
        :disabled="saving"
        @click="handleRestore"
      />
      <Button
        v-if="canEdit"
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
import { fmt, system } from '@planetcrust/human-js'
import {
  components,
  useFileUpload,
  resolveAppLogoUrl,
  useDraftGuard,
  useApplicationsStore,
} from '@planetcrust/human-vue'
import { appIconMap } from '@/utils/appIcons'

const { CCodeEditor, CFileDropZone, CInputDelete, CInputToggleCard, CViewContainer } = components

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')
const applicationsStore = useApplicationsStore()

const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const restoring = ref(false)
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

const kindOptions = computed(() => [
  { label: t('system.applications.editor.unify.kind.section'), value: '' },
  { label: t('system.applications.editor.unify.kind.custom'), value: 'custom' },
])

const isCustom = computed(() => application.value?.unify?.kind === 'custom')

// The server stores a custom application's url as its app-view route; before
// the application exists there is no ID to put in it.
const customUrl = computed(() =>
  application.value?.applicationID ? `app/${application.value.applicationID}` : '',
)

const source = ref('')
// What the editor holds, against what the server last stored: saving is off
// until they differ, and an edit made here is one Claude's next patch has to
// find, so the note under the editor says so.
const draftSource = ref('')
const savingSource = ref(false)
const sourceProblem = ref('')

const canEditSource = computed(() => !!application.value?.canManageSourceOnApplication)
const sourceMeta = computed(() => application.value?.sourceMeta || { size: 0 })

const sourceSizeLabel = computed(() => {
  const size = sourceMeta.value.size || 0
  return size ? `${(size / 1024).toFixed(1)} KB` : t('system.applications.editor.custom.none')
})

const sourceUpdatedLabel = computed(() => {
  const at = sourceMeta.value.updatedAt
  return at
    ? fmt.fullDateTime(at, { dateStyle: 'medium', timeStyle: 'short' })
    : t('system.applications.editor.custom.none')
})

async function loadSource() {
  source.value = ''
  draftSource.value = ''
  sourceProblem.value = ''
  if (!isEdit.value || !isCustom.value) return

  try {
    const rsp = await $SystemAPI.applicationSourceRead({
      applicationID: application.value.applicationID,
    })
    source.value = rsp.source || ''
    draftSource.value = source.value
  } catch (e) {
    sourceProblem.value = t('system.applications.editor.custom.unreadable', {
      reason: e?.message || String(e),
    })
  }
}

// The page goes through the endpoint the MCP tool uses, so the same rules
// refuse the same documents however they arrive. The declaration is sent back
// as it stands: leaving it out would empty what the app may read.
async function handleSaveSource() {
  savingSource.value = true
  try {
    const meta = application.value.sourceMeta || {}
    await $SystemAPI.applicationSourceSet({
      applicationID: application.value.applicationID,
      source: draftSource.value,
      namespace: meta.namespace,
      modules: meta.modules || [],
      writes: meta.writes || [],
    })
    application.value = new system.Application(
      await applicationsStore.findByID(application.value.applicationID),
    )
    await loadSource()
    $toast.toastSuccess(t('system.applications.editor.custom.saved'))
  } catch (e) {
    $toast.toastErrorHandler(t('system.applications.editor.custom.saveError'))(e)
  } finally {
    savingSource.value = false
  }
}

function copySource() {
  navigator.clipboard
    ?.writeText(source.value)
    .then(() => $toast.toastSuccess(t('system.applications.editor.custom.copied')))
    .catch(() => {})
}

// Read-only is one condition, used by the fields, the banner and Save alike —
// a form the user cannot save must not invite them to fill it in.
const canEdit = computed(() => !isEdit.value || !!application.value?.canUpdateApplication)

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
    await loadSource()
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
      meta: application.value.meta,
      unify: application.value.unify,
    }

    if (isEdit.value) {
      payload.applicationID = application.value.applicationID
      const raw = await applicationsStore.update(payload)
      application.value = new system.Application(raw)
      capture()
      await loadSource()
      $toast.toastSuccess(t('notification.application.update.success'))
    } else {
      let created = await applicationsStore.create(payload)
      // Only an update can set a custom application's url: it needs the ID.
      if (created.unify?.kind === 'custom') {
        created = await applicationsStore.update({
          ...payload,
          applicationID: created.applicationID,
        })
      }
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

async function handleRestore() {
  restoring.value = true
  try {
    await applicationsStore.restore(application.value.applicationID)
    await loadApplication()
    $toast.toastSuccess(t('notification.application.restore.success'))
  } catch (e) {
    console.error('Failed to restore application:', e)
    $toast.toastErrorHandler(t('notification.application.restore.error'))(e)
  } finally {
    restoring.value = false
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
