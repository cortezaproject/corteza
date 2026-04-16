<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ pageTitle }}</span>
  </Teleport>

  <Teleport to="#topbar-tools" defer>
    <div v-if="isEdit" class="flex gap-1">
      <Button
        :label="$t('namespace.visit')"
        icon="pi pi-external-link"
        size="small"
        :disabled="!namespace?.enabled"
        @click="visitNamespace"
      />
    </div>
  </Teleport>

  <!-- Loading -->
  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <!-- Form -->
  <Form
    v-else-if="namespace"
    v-slot="$form"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1 overflow-auto">
      <div v-if="isEdit && (namespace?.canExportNamespace || namespace?.canGrant)" class="flex justify-end gap-2 mb-4">
        <Button
          v-if="namespace?.canExportNamespace"
          :label="$t('namespace.export')"
          icon="pi pi-download"
          size="small"
          severity="secondary"
          outlined
          @click="exportNamespace"
        />
        <CPermissionsButton
          v-if="namespace?.canGrant"
          :resource="`corteza::compose:namespace/${namespace.namespaceID}`"
          :title="namespace.name || namespace.slug || namespace.namespaceID"
          :target="namespace.name || namespace.slug || namespace.namespaceID"
          v-tooltip.bottom="$t('general.label.permissions')"
          severity="secondary"
          size="small"
        />
      </div>

      <Card :pt="{ body: { class: 'p-0' }, content: { class: 'p-0' } }" class="overflow-hidden">
        <template #content>
          <div class="flex flex-col gap-5 p-5">
            <!-- Name + Slug inline -->
            <div class="flex gap-4">
              <FormField name="name" class="flex flex-col gap-2 flex-1">
                <label for="name" class="font-medium text-primary">
                  {{ $t('namespace.name.label') }}
                  <span class="text-red-500">*</span>
                </label>
                <InputText
                  id="name"
                  name="name"
                  v-model="namespace.name"
                  :placeholder="$t('namespace.name.placeholder')"
                />
                <Message v-if="$form.name?.invalid" severity="error" size="small" variant="simple">
                  {{ $form.name.error?.message }}
                </Message>
              </FormField>

              <FormField name="slug" class="flex flex-col gap-2 flex-1">
                <label for="slug" class="font-medium text-primary">
                  {{ $t('namespace.slug.label') }}
                </label>
                <InputText
                  id="slug"
                  name="slug"
                  v-model="namespace.slug"
                  :placeholder="$t('namespace.slug.placeholder')"
                />
                <small class="text-muted-color">
                  {{ $t('namespace.slug.description') }}
                </small>
                <Message v-if="$form.slug?.invalid" severity="error" size="small" variant="simple">
                  {{ $form.slug.error?.message }}
                </Message>
              </FormField>
            </div>

            <!-- Labels -->
            <div class="flex flex-col gap-2">
              <label class="font-medium text-primary">{{ $t('namespace.labels.label') }}</label>
              <CInputLabel
                v-model="namespace.labels"
                :placeholder="$t('namespace.labels.placeholder')"
                :create-label="$t('namespace.labels.createNew')"
                :create-dialog-label="$t('namespace.labels.dialogCreate')"
                :name-label="$t('namespace.labels.name')"
                :save-btn-label="$t('general.label.save')"
                :cancel-btn-label="$t('general.label.cancel')"
              />
            </div>

            <!-- Enabled -->
            <div class="flex items-center gap-2">
              <Checkbox id="enabled" v-model="namespace.enabled" binary />
              <label for="enabled">{{ $t('namespace.enabled.label') }}</label>
            </div>

            <Divider />

            <!-- Logo -->
            <div class="flex flex-col gap-2">
              <div class="flex items-center gap-2">
                <Checkbox id="logoEnabled" v-model="namespace.meta.logoEnabled" binary />
                <label for="logoEnabled">{{ $t('namespace.logo.show') }}</label>
              </div>
              <CFileDropZone
                v-if="namespace.meta.logoEnabled"
                accept="image/*"
                :uploading="logoUploading"
                :error="logoError"
                :preview-url="logoPreviewUrl"
                :clearable="!!namespace.meta.logo"
                :drop-label="$t('namespace.logo.upload')"
                compact
                preview-max-width="100%"
                preview-max-height="200px"
                :label="$t('namespace.logo.show')"
                @select="onLogoSelect"
                @clear="onLogoClear"
              />
            </div>

            <!-- Subtitle -->
            <div class="flex flex-col gap-2">
              <label for="subtitle" class="font-medium text-primary">
                {{ $t('namespace.subtitle.label') }}
              </label>
              <InputText
                id="subtitle"
                v-model="namespace.meta.subtitle"
                :placeholder="$t('namespace.subtitle.placeholder')"
              />
            </div>

            <!-- Description -->
            <div class="flex flex-col gap-2">
              <label for="description" class="font-medium text-primary">
                {{ $t('namespace.description.label') }}
              </label>
              <Textarea
                id="description"
                v-model="namespace.meta.description"
                :placeholder="$t('namespace.description.placeholder')"
                rows="3"
                auto-resize
              />
            </div>

            <Divider />

            <!-- Sidebar -->
            <div class="flex items-center gap-2">
              <Checkbox id="hideSidebar" v-model="namespace.meta.hideSidebar" binary />
              <label for="hideSidebar">{{ $t('namespace.sidebar.hide') }}</label>
            </div>

          </div>
        </template>
      </Card>
    </div>

    <!-- Toolbar -->
    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="flex items-center justify-between p-3">
        <Button
          :label="$t('general.label.back')"
          icon="pi pi-arrow-left"
          severity="secondary"
          @click="$router.back()"
        />
        <div class="flex gap-2">
          <CInputDelete
            v-if="isEdit && namespace.canDeleteNamespace"
            :label="$t('general.label.delete')"
            :message="$t('namespace.deleteConfirm')"
            :header="namespace.name"
            :disabled="deleting"
            @confirm="handleDelete"
          />
          <Button
            v-if="isEdit"
            :label="$t('namespace.clone')"
            icon="pi pi-copy"
            severity="secondary"
            :loading="cloning"
            @click="handleClone"
          />
          <Button
            type="submit"
            :label="$t('general.label.save')"
            icon="pi pi-save"
            :loading="saving"
            :disabled="!canSave"
          />
        </div>
      </div>
    </div>
  </Form>
</template>

<script setup>
import { useNamespaceStore } from '@/stores/namespace'
import { compose } from '@cortezaproject/corteza-js-next'
import { components, useFileUpload, useUnsavedGuard } from '@cortezaproject/corteza-vue-next'
import { cloneDeep, isEqual } from 'lodash-es'
import { computed, inject, nextTick, onMounted, ref, watch } from 'vue'

const { CInputDelete, CFileDropZone, CInputLabel } = components
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const $toast = inject('$toast')
const $ComposeAPI = inject('$ComposeAPI')
const $Settings = inject('$Settings')
const namespaceStore = useNamespaceStore()

// State
const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const cloning = ref(false)
const namespace = ref(null)
const initialNamespace = ref(null)

useUnsavedGuard({
  isDirty: () => !saving.value && !deleting.value && !!namespace.value && !!initialNamespace.value && !isEqual(namespace.value, initialNamespace.value),
  messageKey: 'general.editor.unsavedChanges',
})

// Logo upload
const { uploading: logoUploading, uploadError: logoError, uploadFileRaw: uploadLogoRaw, reset: resetLogoUpload } = useFileUpload()

const logoPreviewUrl = computed(() => {
  const logo = namespace.value?.meta?.logo
  if (logo) {
    if (logo.startsWith('http')) return logo
    return $ComposeAPI.baseURL + logo
  }
  // Fall back to global main logo (same as namespace list view)
  return $Settings.attachment('ui.mainLogo') || ''
})

// Computed
const isEdit = computed(() => !!route.params.slug)

const pageTitle = computed(() => {
  return isEdit.value ? t('namespace.edit') : t('namespace.create')
})

const initialValues = computed(() => {
  return {
    name: namespace.value?.name || '',
    slug: namespace.value?.slug || '',
  }
})

const resolver = ref(({ values }) => {
  const errors = {}

  if (!values.name || values.name.trim().length === 0) {
    errors.name = [{ message: t('general.label.required') }]
  }

  if (values.slug && !/^[a-zA-Z][a-zA-Z0-9_]*$/.test(values.slug)) {
    errors.slug = [{ message: t('namespace.slug.invalid-handle-characters') }]
  }

  return { errors }
})

const canSave = computed(() => {
  if (isEdit.value && !namespace.value?.canUpdateNamespace) return false
  return true
})

// Methods
async function loadNamespace() {
  const slug = route.params.slug
  if (!slug) {
    // Create new
    namespace.value = new compose.Namespace({
      enabled: true,
      meta: {
        subtitle: '',
        description: '',
        hideSidebar: false,
        logoEnabled: false,
      },
    })
    initialNamespace.value = cloneDeep(namespace.value)
    return
  }

  loading.value = true
  try {
    // Find by slug or ID
    const found = namespaceStore.getByUrlPart(slug)
    if (found) {
      namespace.value = new compose.Namespace({ ...found })
    } else {
      // Load from API
      await namespaceStore.load({ force: true })
      const ns = namespaceStore.getByUrlPart(slug)
      if (ns) {
        namespace.value = new compose.Namespace({ ...ns })
      } else {
        $toast.toastDanger(t('notification.namespace.loadFailed'))
        router.push({ name: 'namespace.manage' })
      }
    }
    initialNamespace.value = cloneDeep(namespace.value)
  } catch (e) {
    console.error('Failed to load namespace:', e)
    $toast.toastDanger(t('notification.namespace.loadFailed'))
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
  if (!canSave.value) return

  saving.value = true
  try {
    const payload = {
      name: namespace.value.name,
      slug: namespace.value.slug,
      enabled: namespace.value.enabled,
      meta: namespace.value.meta,
      labels: namespace.value.labels || {},
    }

    if (isEdit.value) {
      payload.namespaceID = namespace.value.namespaceID
      const updated = await namespaceStore.update(payload)
      namespace.value = new compose.Namespace({ ...updated })
      initialNamespace.value = cloneDeep(namespace.value)
      $toast.toastSuccess(t('notification.namespace.saved'))
    } else {
      const created = await namespaceStore.create(payload)
      $toast.toastSuccess(t('notification.namespace.saved'))
      // Navigate to edit view
      router.push({
        name: 'namespace.edit',
        params: { slug: created.slug || created.namespaceID },
      })
    }
  } catch (e) {
    console.error('Failed to save namespace:', e)
    $toast.toastDanger(t('notification.namespace.saveFailed'))
  } finally {
    saving.value = false
  }
}

async function handleClone() {
  cloning.value = true
  try {
    const cloned = await namespaceStore.clone({
      ...namespace.value,
      name: `${namespace.value.name} (${t('namespace.cloneSuffix')})`,
      slug: '',
    })
    $toast.toastSuccess(t('notification.namespace.cloned'))
    router.push({
      name: 'namespace.edit',
      params: { slug: cloned.slug || cloned.namespaceID },
    })
  } catch (e) {
    console.error('Failed to clone namespace:', e)
    $toast.toastDanger(t('notification.namespace.cloneFailed'))
  } finally {
    cloning.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await namespaceStore.delete({ namespaceID: namespace.value.namespaceID })
    $toast.toastSuccess(t('notification.namespace.deleted'))
    router.push({ name: 'namespace.manage' })
  } catch (e) {
    console.error('Failed to delete namespace:', e)
    $toast.toastDanger(t('notification.namespace.deleteFailed'))
  } finally {
    deleting.value = false
  }
}

function visitNamespace() {
  router.push({
    name: 'namespace.view',
    params: { slug: namespace.value.slug || namespace.value.namespaceID },
  })
}

function exportNamespace() {
  const params = {
    namespaceID: namespace.value.namespaceID,
    filename: encodeURIComponent((namespace.value.name || 'namespace').replace(/\./g, '-')),
  }

  const token = $ComposeAPI.accessTokenFn ? $ComposeAPI.accessTokenFn() : ''
  const exportUrl = `${$ComposeAPI.baseURL}${$ComposeAPI.namespaceExportEndpoint(params)}?jwt=${encodeURIComponent(token)}`
  window.open(exportUrl)
}

async function onLogoSelect(files) {
  const file = files[0]
  if (!file) return

  try {
    const endpoint = $ComposeAPI.baseURL + $ComposeAPI.namespaceUploadEndpoint()
    const token = $ComposeAPI.accessTokenFn ? $ComposeAPI.accessTokenFn() : ''
    const att = await uploadLogoRaw(file, { url: endpoint, token })

    if (att?.attachmentID) {
      // Build the relative URL to store in meta.logo
      const url = $ComposeAPI.attachmentOriginalEndpoint({
        kind: 'namespace',
        namespaceID: namespace.value.namespaceID || '0',
        attachmentID: att.attachmentID,
        name: att.name || file.name,
      })
      namespace.value.meta.logo = url
    }
  } catch (err) {
    // logoError is set by composable
  }
}

function onLogoClear() {
  namespace.value.meta.logo = ''
  resetLogoUpload()
}

// Lifecycle
onMounted(() => {
  loadNamespace()
})

watch(
  () => route.params.slug,
  () => {
    loadNamespace()
  },
)
</script>
