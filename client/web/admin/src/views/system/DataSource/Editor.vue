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
    v-else-if="dataSource"
    v-slot="$form"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4">
      <div v-if="isEdit" class="flex justify-end gap-2 shrink-0">
        <CPermissionsButton
          v-tooltip.bottom="$t('general.label.permissions')"
          :resource="`corteza::system:dal-connection/${dataSource.connectionID}`"
          :title="dataSource.meta?.name || dataSource.handle"
          :target="dataSource.meta?.name || dataSource.handle"
        />
      </div>
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
              <Tab value="basic">{{ $t('system.data-sources.editor.tabs.basic') }}</Tab>
              <Tab v-if="showDalConfig" value="dal-config">
                {{ $t('system.data-sources.editor.tabs.dal-config') }}
              </Tab>
            </TabList>

            <TabPanels class="flex-1 overflow-y-auto min-h-0">
              <TabPanel value="basic">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <FormField name="name" class="flex flex-col gap-2">
                    <label for="name" class="font-medium text-primary">
                      {{ $t('system.data-sources.editor.basic.form.name.label') }}
                      <span class="text-red-500">*</span>
                    </label>
                    <InputText id="name" name="name" v-model="dataSource.meta.name" />
                    <Message
                      v-if="$form.name?.invalid"
                      severity="error"
                      size="small"
                      variant="simple"
                    >
                      {{ $form.name.error?.message }}
                    </Message>
                  </FormField>

                  <FormField name="handle" class="flex flex-col gap-2">
                    <label for="handle" class="font-medium text-primary">
                      {{ $t('system.data-sources.editor.basic.form.handle.label') }}
                    </label>
                    <InputText id="handle" name="handle" v-model="dataSource.handle" />
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
                    <label for="ownership" class="font-medium text-primary">
                      {{ $t('system.data-sources.editor.basic.form.ownership.label') }}
                    </label>
                    <InputText id="ownership" v-model="dataSource.meta.ownership" />
                  </div>
                </div>
              </TabPanel>

              <TabPanel v-if="showDalConfig" value="dal-config">
                <div class="flex flex-col gap-4">
                  <FormField name="dalConfig" class="flex flex-col gap-2">
                    <label for="dalConfig" class="font-medium text-primary">
                      {{ $t('system.data-sources.editor.dal.form.params.label') }}
                    </label>
                    <Textarea
                      id="dalConfig"
                      name="dalConfig"
                      v-model="rawDalConfig"
                      rows="12"
                      autoResize
                      class="font-mono text-sm"
                      @change="parseDalConfig"
                    />
                    <Message
                      v-if="$form.dalConfig?.invalid"
                      severity="error"
                      size="small"
                      variant="simple"
                    >
                      {{ $form.dalConfig.error?.message }}
                    </Message>
                  </FormField>
                </div>
              </TabPanel>
            </TabPanels>
          </Tabs>
        </template>
      </Card>
    </div>

    <!-- Bottom Actions Toolbar -->
    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="p-3 flex items-center justify-between">
        <Button
          :label="$t('general.label.back')"
          icon="pi pi-arrow-left"
          severity="secondary"
          @click="$router.push({ name: 'system.dataSources' })"
        />
        <div class="flex gap-2">
          <CInputDelete
            v-if="isEdit && dataSource.canDeleteConnection"
            :label="$t('system.data-sources.editor.delete')"
            :message="$t('general.confirm.delete')"
            :header="dataSource.meta?.name || dataSource.handle"
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
import { computed, inject, nextTick, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { system } from '@cortezaproject/corteza-js-next'
import { components, useConfirmDelete, useUnsavedGuard } from '@cortezaproject/corteza-vue-next'
import { cloneDeep, isEqual } from 'lodash-es'

const { CInputDelete } = components

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const dataSource = ref(null)
const initialDataSource = ref(null)
const activeTab = ref('basic')
const rawDalConfig = ref('{}')

const isEdit = computed(() => !!route.params.connectionID)

const showDalConfig = computed(() => {
  if (!isEdit.value) return true
  return dataSource.value?.canManageDalConfig
})

const pageTitle = computed(() =>
  isEdit.value
    ? t('system.data-sources.editor.title.edit')
    : t('system.data-sources.editor.title.create'),
)

const initialValues = computed(() => ({
  name: dataSource.value?.meta?.name || '',
  handle: dataSource.value?.handle || '',
  dalConfig: rawDalConfig.value,
}))

const resolver = ref(({ values }) => {
  const errors = {}

  if (!values.name || values.name.trim().length === 0) {
    errors.name = [{ message: t('general.label.required') }]
  }

  if (values.handle && !/^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$|^[A-Za-z]$/.test(values.handle)) {
    errors.handle = [
      { message: t('system.data-sources.editor.basic.form.handle.invalid-characters') },
    ]
  }

  try {
    if (values.dalConfig) JSON.parse(values.dalConfig)
  } catch {
    errors.dalConfig = [{ message: t('system.data-sources.editor.dal.form.params.description') }]
  }

  return { errors }
})

function parseDalConfig() {
  try {
    const parsed = JSON.parse(rawDalConfig.value || '{}')
    if (parsed && dataSource.value) {
      dataSource.value.config.dal = parsed
    }
  } catch {
    // validation will catch it
  }
}

function initRawFields() {
  if (!dataSource.value) return
  rawDalConfig.value = JSON.stringify(dataSource.value.config?.dal || {}, null, 2)
}

async function loadDataSource() {
  const connectionID = route.params.connectionID
  if (!connectionID) {
    dataSource.value = new system.DalConnection({})
    initRawFields()
    initialDataSource.value = cloneDeep(dataSource.value)
    return
  }

  loading.value = true
  try {
    const raw = await $SystemAPI.dalConnectionRead({ connectionID })
    dataSource.value = new system.DalConnection(raw)
    initRawFields()
    initialDataSource.value = cloneDeep(dataSource.value)
  } catch (e) {
    $toast.toastErrorHandler(t('notification.data-source.fetch.error'))(e)
    router.push({ name: 'system.dataSources' })
  } finally {
    loading.value = false
  }
}

async function handleSubmit({ valid }) {
  if (!valid) {
    activeTab.value = 'basic'
    $toast.toastWarning(t('general.notification.formErrors'))
    nextTick(() => {
      document.querySelector('.p-message-error')?.scrollIntoView({ behavior: 'smooth', block: 'center' })
    })
    return
  }

  saving.value = true
  try {
    parseDalConfig()

    const payload = {
      handle: dataSource.value.handle,
      meta: dataSource.value.meta,
      config: dataSource.value.config,
    }

    if (isEdit.value) {
      payload.connectionID = dataSource.value.connectionID
      const raw = await $SystemAPI.dalConnectionUpdate(payload)
      dataSource.value = new system.DalConnection(raw)
      initRawFields()
      initialDataSource.value = cloneDeep(dataSource.value)
      $toast.toastSuccess(t('notification.data-source.update.success'))
    } else {
      const created = await $SystemAPI.dalConnectionCreate(payload)
      $toast.toastSuccess(t('notification.data-source.create.success'))
      markSaved()
      router.push({
        name: 'system.dataSources.edit',
        params: { connectionID: created.connectionID },
      })
    }
  } catch (e) {
    $toast.toastErrorHandler(
      t(
        `notification.data-source.${isEdit.value ? 'update' : 'create'}.error`,
        'Failed to save data source',
      ),
    )(e)
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await $SystemAPI.dalConnectionDelete({ connectionID: dataSource.value.connectionID })
    $toast.toastSuccess(t('notification.data-source.delete.success'))
    router.push({ name: 'system.dataSources' })
  } catch (e) {
    $toast.toastErrorHandler(t('notification.data-source.delete.error'))(e)
  } finally {
    deleting.value = false
  }
}

const { markSaved } = useUnsavedGuard({
  isDirty: () => !saving.value && !deleting.value && !!dataSource.value && !!initialDataSource.value && !isEqual(dataSource.value, initialDataSource.value),
  messageKey: 'general.editor.unsavedChanges',
})

onMounted(() => loadDataSource())

watch(
  () => route.params.connectionID,
  () => loadDataSource(),
)
</script>
