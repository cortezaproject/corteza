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
    v-else-if="connection"
    v-slot="$form"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4">
      <div v-if="isEdit" class="flex justify-end gap-2 shrink-0">
        <Button
          :label="$t('system.connections.editor.manageConfigured')"
          icon="pi pi-cog"
          severity="secondary"
          size="small"
          @click="
            $router.push({
              name: 'system.connections.configure',
              params: { connectionID: connection.connectionID },
            })
          "
        />
        <CPermissionsButton
          v-tooltip.bottom="$t('general.label.permissions')"
          :resource="`corteza::system:dal-connection/${connection.connectionID}`"
          :title="connection.meta?.short || connection.handle || connection.connectionID"
          :target="connection.meta?.short || connection.handle || connection.connectionID"
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
              <Tab value="general">{{ $t('system.connections.editor.general') }}</Tab>
              <Tab value="configuration">{{ $t('system.connections.editor.configuration') }}</Tab>
            </TabList>

            <TabPanels class="flex-1 overflow-y-auto min-h-0 p-0">
              <TabPanel value="general" class="p-4">
                <div class="flex flex-col gap-6">
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <FormField name="name" class="flex flex-col gap-2">
                      <label for="name" class="font-medium text-primary">
                        {{ $t('system.connections.editor.info.name') }}
                        <span class="text-red-500">*</span>
                      </label>
                      <InputText id="name" name="name" v-model="connection.meta.short" />
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
                        {{ $t('system.connections.editor.info.handle') }}
                      </label>
                      <InputText id="handle" name="handle" v-model="connection.handle" />
                      <Message
                        v-if="$form.handle?.invalid"
                        severity="error"
                        size="small"
                        variant="simple"
                      >
                        {{ $form.handle.error?.message }}
                      </Message>
                    </FormField>

                    <div class="flex flex-col gap-2 md:col-span-2">
                      <label for="description" class="font-medium text-primary">
                        {{ $t('system.connections.editor.info.description') }}
                      </label>
                      <Textarea id="description" v-model="connection.meta.description" rows="3" />
                    </div>
                  </div>

                  <FormField name="service" class="flex flex-col gap-2">
                    <label for="service" class="font-medium text-primary">
                      {{ $t('system.connections.editor.service.title') }}
                    </label>
                    <Textarea
                      id="service"
                      name="service"
                      v-model="rawJSON.service"
                      rows="10"
                      autoResize
                      class="font-mono text-sm"
                      @change="() => parseJSONField('service')"
                    />
                    <Message
                      v-if="$form.service?.invalid"
                      severity="error"
                      size="small"
                      variant="simple"
                    >
                      {{ $form.service.error?.message }}
                    </Message>
                  </FormField>
                </div>
              </TabPanel>

              <TabPanel value="configuration" class="p-4">
                <div class="flex flex-col gap-4">
                  <Panel
                    :header="$t('system.connections.editor.configurations.resources')"
                    toggleable
                    :collapsed="false"
                    class="shadow"
                  >
                    <FormField name="resources" class="flex flex-col gap-2">
                      <Textarea
                        id="resources"
                        name="resources"
                        v-model="rawJSON.resources"
                        rows="10"
                        class="font-mono text-sm max-h-[50vh] overflow-y-auto"
                        @change="() => parseJSONField('resources')"
                      />
                      <Message
                        v-if="$form.resources?.invalid"
                        severity="error"
                        size="small"
                        variant="simple"
                      >
                        {{ $form.resources.error?.message }}
                      </Message>
                    </FormField>
                  </Panel>

                  <Panel
                    :header="$t('system.connections.editor.configurations.operations')"
                    toggleable
                    :collapsed="false"
                    class="shadow"
                  >
                    <FormField name="operations" class="flex flex-col gap-2">
                      <Textarea
                        id="operations"
                        name="operations"
                        v-model="rawJSON.operations"
                        rows="10"
                        class="font-mono text-sm max-h-[50vh] overflow-y-auto"
                        @change="() => parseJSONField('operations')"
                      />
                      <Message
                        v-if="$form.operations?.invalid"
                        severity="error"
                        size="small"
                        variant="simple"
                      >
                        {{ $form.operations.error?.message }}
                      </Message>
                    </FormField>
                  </Panel>
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
          @click="$router.push({ name: 'system.connections' })"
        />
        <div class="flex gap-2">
          <CInputDelete
            v-if="isEdit && connection.canDeleteConnection"
            :label="$t('system.connections.editor.delete')"
            :message="$t('system.connections.editor.deleteConfirm')"
            :header="connection.meta?.short || connection.handle"
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
import { computed, inject, nextTick, onMounted, ref, watch, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { system } from '@planetcrust/human-js'
import { components, useUnsavedGuard } from '@planetcrust/human-vue'
import { cloneDeep, isEqual } from 'lodash-es'

const { CInputDelete } = components

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

// State
const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const connection = ref(null)
const initialConnection = ref(null)
const initialRawJSON = ref(null)
const activeTab = ref('general')

const rawJSON = reactive({
  service: '{}',
  resources: '[]',
  standardOperations: '{}',
  operations: '[]',
})

// Computed
const isEdit = computed(() => !!route.params.connectionID)

const pageTitle = computed(() => {
  return isEdit.value
    ? t('system.connections.editor.title.edit')
    : t('system.connections.editor.title.create')
})

const initialValues = computed(() => {
  return {
    name: connection.value?.meta?.short || '',
    handle: connection.value?.handle || '',
    service: rawJSON.service,
    resources: rawJSON.resources,
    standardOperations: rawJSON.standardOperations,
    operations: rawJSON.operations,
  }
})

const resolver = ref(({ values }) => {
  const errors = {}

  if (!values.name || values.name.trim().length === 0) {
    errors.name = [{ message: t('general.label.required') }]
  }

  if (values.handle && !/^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$|^[A-Za-z]$/.test(values.handle)) {
    errors.handle = [{ message: t('system.connections.editor.info.invalid-characters') }]
  }

  const jsonFields = ['service', 'resources', 'standardOperations', 'operations']
  jsonFields.forEach(field => {
    try {
      if (values[field]) {
        JSON.parse(values[field])
      }
    } catch {
      errors[field] = [{ message: t('system.connections.editor.configurations.invalidJSON') }]
    }
  })

  return { errors }
})

// Methods
function parseJSONField(field) {
  try {
    const parsed = JSON.parse(rawJSON[field] || 'null')

    if (parsed) {
      if (field === 'service') connection.value.service = parsed
      if (field === 'resources') connection.value.resources = parsed
      if (field === 'standardOperations') connection.value.standardOperations = parsed
      if (field === 'operations') connection.value.operations = parsed
    }
  } catch {
    // Leave alone, parser failure will be caught by resolver
  }
}

function initJSONFields() {
  if (!connection.value) return

  rawJSON.service = JSON.stringify(connection.value.service || {}, null, 2)
  rawJSON.resources = JSON.stringify(connection.value.resources || [], null, 2)
  rawJSON.standardOperations = JSON.stringify(connection.value.standardOperations || {}, null, 2)
  rawJSON.operations = JSON.stringify(connection.value.operations || [], null, 2)
}

async function loadConnection() {
  const connectionID = route.params.connectionID
  if (!connectionID) {
    // Create new
    connection.value = new system.Connection({
      service: {
        baseURL: { value: '' },
        auth: { method: 'none' },
      },
    })
    initJSONFields()
    initialConnection.value = cloneDeep(connection.value)
    initialRawJSON.value = cloneDeep(rawJSON)
    return
  }

  loading.value = true
  try {
    const raw = await $SystemAPI.connectionRead({ connectionID })
    connection.value = new system.Connection(raw)
    initJSONFields()
    initialConnection.value = cloneDeep(connection.value)
    initialRawJSON.value = cloneDeep(rawJSON)
  } catch (e) {
    console.error('Failed to load connection:', e)
    $toast.toastErrorHandler(t('notification.connection.fetch.error'))(e)
    router.push({ name: 'system.connections' })
  } finally {
    loading.value = false
  }
}

async function handleSubmit({ valid }) {
  if (!valid) {
    activeTab.value = 'general'
    $toast.toastWarning(t('general.notification.formErrors'))
    nextTick(() => {
      document.querySelector('.p-message-error')?.scrollIntoView({ behavior: 'smooth', block: 'center' })
    })
    return
  }

  if (isEdit.value && !connection.value?.canUpdateConnection) return

  saving.value = true
  try {
    // Attempt parse one last time to ensure connection has latest data
    parseJSONField('service')
    parseJSONField('resources')
    parseJSONField('standardOperations')
    parseJSONField('operations')

    const payload = {
      handle: connection.value.handle,
      status: connection.value.status,
      meta: connection.value.meta,
      service: connection.value.service,
      resources: connection.value.resources,
      standardOperations: connection.value.standardOperations,
      operations: connection.value.operations,
      labels: connection.value.labels,
    }

    if (isEdit.value) {
      payload.connectionID = connection.value.connectionID
      const raw = await $SystemAPI.connectionUpdate(payload)
      connection.value = new system.Connection(raw)
      initJSONFields()
      initialConnection.value = cloneDeep(connection.value)
      initialRawJSON.value = cloneDeep(rawJSON)
      $toast.toastSuccess(t('notification.connection.update.success'))
    } else {
      const created = await $SystemAPI.connectionCreate(payload)
      $toast.toastSuccess(t('notification.connection.create.success'))
      markSaved()
      router.push({
        name: 'system.connections.edit',
        params: { connectionID: created.connectionID },
      })
    }
  } catch (e) {
    console.error('Failed to save connection:', e)
    $toast.toastErrorHandler(
      t(`notification.connection.${isEdit.value ? 'update' : 'create'}.error`),
    )(e)
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await $SystemAPI.connectionDelete({ connectionID: connection.value.connectionID })
    $toast.toastSuccess(t('notification.connection.delete.success'))
    router.push({ name: 'system.connections' })
  } catch (e) {
    console.error('Failed to delete connection:', e)
    $toast.toastErrorHandler(t('notification.connection.delete.error'))(e)
  } finally {
    deleting.value = false
  }
}

const { markSaved } = useUnsavedGuard({
  isDirty: () => !saving.value && !deleting.value && !!connection.value && !!initialConnection.value && (!isEqual(connection.value, initialConnection.value) || !isEqual({ ...rawJSON }, initialRawJSON.value)),
  messageKey: 'general.editor.unsavedChanges',
})

// Lifecycle
onMounted(() => {
  loadConnection()
})

watch(
  () => route.params.connectionID,
  () => {
    loadConnection()
  },
)
</script>
