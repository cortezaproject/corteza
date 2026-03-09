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
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4">
      <Card
        :pt="{
          body: { class: 'p-0 flex flex-col h-full min-h-0' },
          content: { class: 'p-0 flex flex-col h-full min-h-0' },
        }"
        class="overflow-hidden flex-1 min-h-0 flex flex-col"
      >
        <template #content>
          <Tabs v-model:value="activeTab" class="flex flex-col h-full min-h-0">
            <TabList class="rounded-t-lg shrink-0">
              <Tab value="basic">{{ $t('system.llmProviders.editor.tabs.basic') }}</Tab>
              <Tab value="config">{{ $t('system.llmProviders.editor.tabs.config') }}</Tab>
            </TabList>

            <TabPanels class="flex-1 overflow-y-auto min-h-0">
              <TabPanel value="basic">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <FormField name="short" class="flex flex-col gap-2">
                    <label for="short" class="font-medium text-primary">
                      {{ $t('system.llmProviders.editor.info.short') }}
                    </label>
                    <InputText id="short" name="short" v-model="llmProvider.meta.short" />
                  </FormField>

                  <FormField name="handle" class="flex flex-col gap-2">
                    <label for="handle" class="font-medium text-primary">
                      {{ $t('system.llmProviders.editor.info.handle') }} *
                    </label>
                    <InputText id="handle" name="handle" v-model="llmProvider.handle" />
                    <Message
                      v-if="$form.handle?.invalid"
                      severity="error"
                      size="small"
                      variant="simple"
                    >
                      {{ $form.handle.error?.message }}
                    </Message>
                  </FormField>

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

                  <div class="flex flex-col gap-2">
                    <label for="status" class="font-medium text-primary">
                      {{ $t('system.llmProviders.editor.info.status') }}
                    </label>
                    <Select
                      id="status"
                      v-model="llmProvider.status"
                      :options="statusOptions"
                      option-label="label"
                      option-value="value"
                    />
                  </div>

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

                  <div v-if="!isEdit" class="flex flex-col gap-2 md:col-span-2">
                    <label for="apiKey" class="font-medium text-primary">
                      {{ $t('system.llmProviders.editor.info.apiKey') }} *
                    </label>
                    <InputText id="apiKey" v-model="apiKey" />
                  </div>
                </div>
              </TabPanel>

              <TabPanel value="config">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div class="flex flex-col gap-2 md:col-span-2">
                    <label for="promptURL" class="font-medium text-primary">
                      {{ $t('system.llmProviders.editor.config.promptURL') }}
                    </label>
                    <InputText id="promptURL" v-model="llmProvider.config.promptURL" />
                  </div>

                  <div class="flex flex-col gap-2">
                    <label for="model" class="font-medium text-primary">
                      {{ $t('system.llmProviders.editor.config.model') }}
                    </label>
                    <InputText id="model" v-model="llmProvider.config.model" />
                  </div>

                  <div class="flex flex-col gap-2">
                    <label for="temperature" class="font-medium text-primary">
                      {{ $t('system.llmProviders.editor.config.temperature') }}
                    </label>
                    <InputNumber
                      id="temperature"
                      v-model="llmProvider.config.temperature"
                      :min="0"
                      :max="2"
                      :step="0.1"
                      :minFractionDigits="1"
                      :maxFractionDigits="2"
                    />
                  </div>

                  <div class="flex flex-col gap-2">
                    <label for="maxTokens" class="font-medium text-primary">
                      {{ $t('system.llmProviders.editor.config.maxTokens') }}
                    </label>
                    <InputNumber id="maxTokens" v-model="llmProvider.config.maxTokens" :min="0" />
                  </div>

                  <div class="flex flex-col gap-2">
                    <label for="timeout" class="font-medium text-primary">
                      {{ $t('system.llmProviders.editor.config.timeout') }}
                    </label>
                    <InputText id="timeout" v-model="llmProvider.config.timeout" />
                  </div>
                </div>
              </TabPanel>
            </TabPanels>
          </Tabs>
        </template>
      </Card>
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
</template>

<script setup>
import { computed, inject, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { system } from '@cortezaproject/corteza-js-next'
import { components } from '@cortezaproject/corteza-vue-next'

const { CInputDelete } = components

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const llmProvider = ref(null)
const apiKey = ref('')
const activeTab = ref('basic')

const isEdit = computed(() => !!route.params.llmProviderID)

const pageTitle = computed(() =>
  isEdit.value
    ? t('system.llmProviders.editor.title.edit')
    : t('system.llmProviders.editor.title.create'),
)

const providerOptions = computed(() => [
  { label: 'OpenAI', value: 'openai' },
  { label: 'Anthropic', value: 'anthropic' },
  { label: 'Mistral', value: 'mistral' },
])

const statusOptions = computed(() => [
  { label: t('system.llmProviders.editor.info.statusOptions.active'), value: 'active' },
  { label: t('system.llmProviders.editor.info.statusOptions.inactive'), value: 'inactive' },
])

const initialValues = computed(() => ({
  handle: llmProvider.value?.handle || '',
}))

const resolver = ref(({ values }) => {
  const errors = {}

  if (!values.handle || values.handle.trim().length === 0) {
    errors.handle = [{ message: t('general.label.required') }]
  }

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
    llmProvider.value = new system.LlmProvider({ status: 'active' })
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

onMounted(() => loadLlmProvider())
watch(
  () => route.params.llmProviderID,
  () => loadLlmProvider(),
)
</script>
