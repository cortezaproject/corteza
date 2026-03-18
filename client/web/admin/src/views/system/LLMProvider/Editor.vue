<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ pageTitle }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <Form
    v-else-if="llmProvider"
    v-slot="$form"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <div v-if="isEdit" class="flex justify-end items-center">
        <Tag
          :value="$t(`system.llmProviders.editor.info.statusOptions.${llmProvider.status}`)"
          :severity="
            llmProvider.status === 'active'
              ? 'success'
              : llmProvider.status === 'unauthorized'
                ? 'danger'
                : 'warn'
          "
          :icon="
            llmProvider.status === 'active'
              ? 'pi pi-check-circle'
              : llmProvider.status === 'unauthorized'
                ? 'pi pi-times-circle'
                : 'pi pi-pause-circle'
          "
        />
        <Divider layout="vertical" />
        <Button
          :label="$t('system.llmProviders.editor.info.updateKey')"
          icon="pi pi-key"
          severity="secondary"
          size="small"
          outlined
          @click="openApiKeyDialog"
        />
      </div>

      <Panel :header="$t('system.llmProviders.editor.tabs.basic')" toggleable :collapsed="false">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField name="short" class="flex flex-col gap-2">
            <label for="short" class="font-medium text-primary">
              {{ $t('system.llmProviders.editor.info.short') }}
            </label>
            <InputText id="short" name="short" v-model="llmProvider.meta.short" />
          </FormField>

          <FormField name="handle" class="flex flex-col gap-2">
            <label for="handle" class="font-medium text-primary">
              {{ $t('system.llmProviders.editor.info.handle') }}
            </label>
            <InputText id="handle" name="handle" v-model="llmProvider.handle" />
            <Message v-if="$form.handle?.invalid" severity="error" size="small" variant="simple">
              {{ $form.handle.error?.message }}
            </Message>
          </FormField>

          <div class="flex flex-col gap-2 md:col-span-2">
            <label for="description" class="font-medium text-primary">
              {{ $t('system.llmProviders.editor.info.description') }}
            </label>
            <Textarea
              id="description"
              v-model="llmProvider.meta.description"
              rows="3"
              auto-resize
            />
          </div>

          <Divider class="md:col-span-2" />

          <div class="flex flex-col gap-2">
            <label for="provider" class="font-medium text-primary">
              {{ $t('system.llmProviders.editor.info.provider') }} *
            </label>
            <Select
              id="provider"
              v-model="llmProvider.provider"
              :options="providerOptions"
              option-label="label"
              option-value="value"
            />
          </div>

          <div class="flex flex-col gap-1">
            <label for="promptURL" class="font-medium text-primary">
              {{ $t('system.llmProviders.editor.config.promptURL') }}
            </label>
            <InputText id="promptURL" v-model="llmProvider.config.promptURL" />
            <small class="text-muted-color">
              {{ $t('system.llmProviders.editor.config.promptURLHelp') }}
            </small>
          </div>

          <div v-if="!isEdit" class="flex flex-col gap-2">
            <label for="apiKey" class="font-medium text-primary">
              {{ $t('system.llmProviders.editor.info.apiKey') }} *
            </label>
            <InputText id="apiKey" v-model="apiKey" />
          </div>
        </div>
      </Panel>
    </div>

    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="p-3 flex items-center justify-between">
        <Button
          :label="$t('general.label.back')"
          icon="pi pi-arrow-left"
          severity="secondary"
          @click="$router.push({ name: 'system.llmProviders' })"
        />
        <div class="flex gap-2">
          <CInputDelete
            v-if="isEdit && llmProvider.canDeleteLlmProvider"
            :label="$t('system.llmProviders.editor.info.delete')"
            :message="$t('general.confirm.delete')"
            :header="llmProvider.meta?.short || llmProvider.handle"
            :disabled="deleting"
            @confirm="handleDelete"
          />
          <Button
            type="submit"
            :label="$t('general.label.save')"
            icon="pi pi-save"
            :loading="saving"
          />
        </div>
      </div>
    </div>
  </Form>

  <Dialog
    v-model:visible="apiKeyDialog"
    :header="$t('system.llmProviders.editor.info.updateKey')"
    modal
    class="w-full max-w-lg"
  >
    <div class="flex flex-col gap-2">
      <label for="dialogApiKey" class="font-medium text-primary">
        {{ $t('system.llmProviders.editor.info.apiKey') }}
      </label>
      <InputText id="dialogApiKey" v-model="apiKey" autocomplete="off" />
    </div>

    <template #footer>
      <Button :label="$t('general.label.cancel')" severity="secondary" @click="closeApiKeyDialog" />
      <Button
        :label="$t('general.label.save')"
        icon="pi pi-save"
        :loading="savingApiKey"
        :disabled="!apiKey"
        @click="handleUpdateApiKey"
      />
    </template>
  </Dialog>
</template>

<script setup>
import { computed, inject, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { system } from '@cortezaproject/corteza-js-next'
import { components } from '@cortezaproject/corteza-vue-next'

const { CInputDelete, CInputModel } = components

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const savingApiKey = ref(false)
const apiKeyDialog = ref(false)
const llmProvider = ref(null)
const apiKey = ref('')

const isEdit = computed(() => !!route.params.llmProviderID)

const pageTitle = computed(() =>
  isEdit.value
    ? t('system.llmProviders.editor.title.edit')
    : t('system.llmProviders.editor.title.create'),
)

const providerDefaultURLs = {
  mistral: 'https://api.mistral.ai/v1',
  anthropic: 'https://api.anthropic.com',
}

const providerOptions = computed(() => [
  { label: 'Mistral', value: 'mistral' },
  { label: 'Anthropic', value: 'anthropic' },
  { label: t('system.llmProviders.editor.info.providerOther'), value: 'other' },
])

const statusOptions = computed(() => [
  { label: t('system.llmProviders.editor.info.statusOptions.active'), value: 'active' },
  { label: t('system.llmProviders.editor.info.statusOptions.inactive'), value: 'inactive' },
  { label: t('system.llmProviders.editor.info.statusOptions.unauthorized'), value: 'unauthorized' },
])

const initialValues = computed(() => ({
  handle: llmProvider.value?.handle || '',
}))

const resolver = ref(({ values }) => {
  const errors = {}

  if (values.handle && !/^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$|^[A-Za-z]$/.test(values.handle)) {
    errors.handle = [
      { message: t('system.llmProviders.editor.info.handle.invalid-handle-characters') },
    ]
  }

  return { errors }
})

async function loadLlmProvider() {
  const llmProviderID = route.params.llmProviderID
  if (!llmProviderID) {
    const defaultProvider = providerOptions.value[0]?.value || ''
    llmProvider.value = new system.LlmProvider({
      status: 'active',
      provider: defaultProvider,
      config: { temperature: 0.7, promptURL: providerDefaultURLs[defaultProvider] || '' },
    })
    return
  }

  loading.value = true
  try {
    const raw = await $SystemAPI.llmProviderRead({ llmProviderID })
    llmProvider.value = new system.LlmProvider(raw)
  } catch (e) {
    $toast.toastErrorHandler(t('notification.llmProvider.fetch.error'))(e)
    router.push({ name: 'system.llmProviders' })
  } finally {
    loading.value = false
  }
}

async function handleSubmit({ valid }) {
  if (!valid) return

  saving.value = true
  try {
    const payload = {
      handle: llmProvider.value.handle,
      provider: llmProvider.value.provider,
      status: llmProvider.value.status,
      meta: llmProvider.value.meta,
      config: llmProvider.value.config,
    }

    if (isEdit.value) {
      payload.llmProviderID = llmProvider.value.llmProviderID
      const raw = await $SystemAPI.llmProviderUpdate(payload)
      llmProvider.value = new system.LlmProvider(raw)
      $toast.toastSuccess(t('notification.llmProvider.update.success'))
    } else {
      payload.apiKey = apiKey.value
      const created = await $SystemAPI.llmProviderCreate(payload)
      $toast.toastSuccess(t('notification.llmProvider.create.success'))
      router.push({
        name: 'system.llmProviders.edit',
        params: { llmProviderID: created.llmProviderID },
      })
    }
  } catch (e) {
    $toast.toastErrorHandler(
      t(`notification.llmProvider.${isEdit.value ? 'update' : 'create'}.error`),
    )(e)
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await $SystemAPI.llmProviderDelete({ llmProviderID: llmProvider.value.llmProviderID })
    $toast.toastSuccess(t('notification.llmProvider.delete.success'))
    router.push({ name: 'system.llmProviders' })
  } catch (e) {
    $toast.toastErrorHandler(t('notification.llmProvider.delete.error'))(e)
  } finally {
    deleting.value = false
  }
}

function openApiKeyDialog() {
  apiKey.value = ''
  apiKeyDialog.value = true
}

function closeApiKeyDialog() {
  apiKey.value = ''
  apiKeyDialog.value = false
}

async function handleUpdateApiKey() {
  savingApiKey.value = true
  try {
    const payload = {
      llmProviderID: llmProvider.value.llmProviderID,
      handle: llmProvider.value.handle,
      provider: llmProvider.value.provider,
      status: llmProvider.value.status,
      meta: llmProvider.value.meta,
      config: llmProvider.value.config,
      apiKey: apiKey.value,
    }
    const raw = await $SystemAPI.llmProviderUpdate(payload)
    llmProvider.value = new system.LlmProvider(raw)
    apiKey.value = ''
    apiKeyDialog.value = false
    $toast.toastSuccess(t('notification.llmProvider.apiKeyUpdate.success'))
  } catch (e) {
    $toast.toastErrorHandler(t('notification.llmProvider.apiKeyUpdate.error'))(e)
  } finally {
    savingApiKey.value = false
  }
}

// Prefill prompt URL when provider changes
watch(
  () => llmProvider.value?.provider,
  newProvider => {
    if (!llmProvider.value || !newProvider) return
    const currentURL = llmProvider.value.config?.promptURL || ''
    const defaultURLs = Object.values(providerDefaultURLs)
    // Only prefill if empty or is another provider's default
    if (!currentURL || defaultURLs.includes(currentURL)) {
      if (!llmProvider.value.config) llmProvider.value.config = {}
      llmProvider.value.config.promptURL = providerDefaultURLs[newProvider] || ''
    }
  },
)

// Switch provider when URL is manually changed
watch(
  () => llmProvider.value?.config?.promptURL,
  newURL => {
    if (!llmProvider.value || !newURL) return
    // Find if URL matches a known provider
    const matchedProvider = Object.entries(providerDefaultURLs).find(([, url]) => url === newURL)
    if (matchedProvider) {
      llmProvider.value.provider = matchedProvider[0]
    } else if (llmProvider.value.provider !== 'other') {
      llmProvider.value.provider = 'other'
    }
  },
)

onMounted(() => loadLlmProvider())
watch(
  () => route.params.llmProviderID,
  () => loadLlmProvider(),
)
</script>
