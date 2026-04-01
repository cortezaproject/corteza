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
              <Tab v-if="isEdit" value="configured">
                {{ $t('system.configuredConnections.list.title') }}
              </Tab>
            </TabList>

            <TabPanels class="flex-1 overflow-y-auto min-h-0 p-0">
              <TabPanel value="general" class="p-4">
                <div class="flex flex-col gap-6">
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <FormField name="name" class="flex flex-col gap-2">
                      <label for="name" class="font-medium text-primary">
                        {{ $t('system.connections.editor.info.name') }}
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

              <TabPanel v-if="isEdit" value="configured" class="flex flex-col h-full min-h-0">
                <CResourceList
                  primary-key="configurationID"
                  :fields="configuredConnectionFields"
                  :items="configuredConnectionList"
                  :filter="configuredConnectionsFilter"
                  @update:filter="Object.assign(configuredConnectionsFilter, $event)"
                  :sorting="configuredConnectionsSorting"
                  :pagination="configuredConnectionsPagination"
                  :loading="configuredConnectionsLoading"
                  :translations="{
                    searchPlaceholder: $t('system.configuredConnections.list.searchPlaceholder'),
                    showingPagination: 'general.resourceList.pagination.showing',
                    singlePluralPagination: 'general.resourceList.pagination.single',
                    prevPagination: $t('general.resourceList.pagination.prev'),
                    nextPagination: $t('general.resourceList.pagination.next'),
                    recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
                    resourceSingle: $t('system.configuredConnections.list.resourceSingle'),
                    resourcePlural: $t('system.configuredConnections.list.resourcePlural'),
                  }"
                  clickable
                  class="h-full"
                  @sort="handleConfiguredConnectionsSort"
                  @row-click="handleConfiguredConnectionClick"
                  @page-change="handleConfiguredConnectionsPageChange"
                >
                  <template #header>
                    <Button
                      :label="$t('system.configuredConnections.list.createLabel')"
                      icon="pi pi-plus"
                      size="small"
                      @click="createConfiguredConnection"
                    />
                  </template>

                  <template #body-labels="{ data }">
                    <div class="flex gap-1 flex-wrap">
                      <Tag
                        v-for="(val, key) in data.labels"
                        :key="key"
                        :value="`${key}: ${val}`"
                        severity="info"
                      />
                    </div>
                  </template>

                  <template #body-status="{ data }">
                    <Tag
                      v-if="data.status"
                      :value="
                        $t(`system.connections.editor.statusValues.${data.status}`, data.status)
                      "
                      :severity="data.status === 'active' ? 'success' : 'secondary'"
                    />
                    <span v-else>-</span>
                  </template>

                  <template #body-updatedAt="{ data }">
                    {{ locFullDateTime(data.deletedAt || data.updatedAt || data.createdAt) }}
                  </template>

                  <template #body-actions="{ data }">
                    <Button
                      icon="pi pi-ellipsis-v"
                      text
                      severity="secondary"
                      size="small"
                      class="row-action-btn w-full"
                      @click.stop="toggleConfiguredConnectionActionsMenu($event, data)"
                    />
                  </template>
                </CResourceList>
              </TabPanel>
            </TabPanels>
          </Tabs>
        </template>
      </Card>
    </div>

    <TieredMenu
      ref="configuredConnectionActionsMenu"
      :model="configuredConnectionActionsMenuItems"
      popup
    >
      <template #item="{ item, props }">
        <router-link v-if="item.route" v-slot="{ href, navigate }" :to="item.route" custom>
          <a v-ripple :href="href" v-bind="props.action" @click="navigate">
            <span :class="item.icon" />
            <span class="ml-2">{{ item.label }}</span>
          </a>
        </router-link>
        <a v-else v-ripple v-bind="props.action" :class="item.class">
          <span :class="item.icon" />
          <span class="ml-2">{{ item.label }}</span>
        </a>
      </template>
    </TieredMenu>

    <!-- Configured Connection Modal -->
    <Dialog
      v-model:visible="configuredConnectionModal"
      modal
      :header="
        activeConfiguredConnection?.configurationID
          ? $t('system.configuredConnections.editor.title.edit')
          : $t('system.configuredConnections.editor.title.create')
      "
      :style="{ width: '50vw', maxHeight: '80vh' }"
      :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"
      :pt="{
        content: { class: 'p-0 flex flex-col !overflow-hidden' },
      }"
    >
      <Form
        v-if="configuredConnectionModal"
        v-slot="$modalForm"
        :resolver="configuredConnectionResolver"
        :initialValues="configuredConnectionInitialValues"
        @submit="handleConfiguredConnectionSubmit"
        class="flex flex-col h-full min-h-0"
      >
        <div class="p-4 flex flex-col gap-4 flex-1 min-h-0 overflow-y-auto">
          <FormField name="name" class="flex flex-col gap-2">
            <label for="ccName" class="font-medium text-primary">
              {{ $t('system.configuredConnections.editor.info.name') }}
            </label>
            <InputText id="ccName" name="name" v-model="activeConfiguredConnection.name" />
            <Message v-if="$modalForm.name?.invalid" severity="error" size="small" variant="simple">
              {{ $modalForm.name.error?.message }}
            </Message>
          </FormField>

          <div v-if="activeConfiguredConnection.configurationID" class="flex flex-col gap-2">
            <label class="font-medium text-primary">
              {{ $t('system.configuredConnections.editor.info.status') }}
            </label>
            <div>
              <Tag
                :value="
                  $t(
                    `system.connections.editor.statusValues.${activeConfiguredConnection.status}`,
                    activeConfiguredConnection.status || '-',
                  )
                "
                :severity="activeConfiguredConnection.status === 'active' ? 'success' : 'secondary'"
              />
            </div>
          </div>

          <div v-for="param in uniqueDerivedParams" :key="param.name" class="flex flex-col gap-2">
            <label :for="`param-${param.name}`" class="font-medium text-primary">
              {{ param.label || param.name }}
              <span v-if="param.required" class="text-red-500">*</span>
            </label>
            <InputText :id="`param-${param.name}`" v-model="paramValues[param.name]" />
            <small v-if="param.description" class="text-muted-color">
              {{ param.description }}
            </small>
          </div>
        </div>

        <div class="border-t border-surface p-3 flex gap-2 shrink-0">
          <div v-if="activeConfiguredConnection.configurationID" class="flex">
            <CInputDelete
              :label="$t('general.label.delete')"
              :message="$t('system.configuredConnections.list.deleteConfirm')"
              :header="
                activeConfiguredConnection.name ||
                $t('system.configuredConnections.list.resourceSingle')
              "
              size="small"
              @confirm="handleConfiguredConnectionDeleteFromModal"
            />

            <div class="flex">
              <Divider layout="vertical" />

              <Button
                :label="$t('system.configuredConnections.editor.check')"
                size="small"
                severity="info"
                outlined
                :loading="checkingConfiguredConnection"
                @click="handleConfiguredConnectionCheck"
              />
            </div>

            <div v-if="activeConfiguredConnection.status !== 'active'" class="flex">
              <Divider layout="vertical" />

              <Button
                :label="$t('system.configuredConnections.editor.enable')"
                icon="pi pi-check-circle"
                size="small"
                severity="success"
                outlined
                :loading="enablingConfiguredConnection"
                @click="handleConfiguredConnectionEnable"
              />
            </div>
          </div>
          <div class="flex gap-2 ml-auto">
            <Button
              :label="$t('general.label.cancel')"
              severity="secondary"
              size="small"
              outlined
              @click="configuredConnectionModal = false"
            />
            <Button
              type="submit"
              :label="$t('general.label.save')"
              size="small"
              :loading="savingConfiguredConnection"
            />
          </div>
        </div>
      </Form>
    </Dialog>

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
import { computed, inject, onMounted, ref, watch, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { system } from '@cortezaproject/corteza-js-next'
import {
  components,
  filters,
  useConfirmDelete,
  useResourceList,
} from '@cortezaproject/corteza-vue-next'

const { CInputDelete, CResourceList } = components
const { locFullDateTime } = filters

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

// State
const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const connection = ref(null)
const activeTab = ref('general')

const configuredConnectionModal = ref(false)
const savingConfiguredConnection = ref(false)
const enablingConfiguredConnection = ref(false)
const checkingConfiguredConnection = ref(false)
const activeConfiguredConnection = ref(null)
const configuredConnectionRawLabels = ref('{}')
const paramValues = reactive({})
const configuredConnectionActionsMenu = ref()
const configuredConnectionActionsMenuItems = ref([])

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

const configuredConnectionFields = [
  {
    key: 'name',
    sortable: true,
    header: t('system.configuredConnections.list.columns.name'),
  },
  {
    key: 'labels',
    sortable: false,
    header: t('system.configuredConnections.list.columns.labels'),
  },
  {
    key: 'status',
    sortable: true,
    header: t('system.configuredConnections.list.columns.status'),
  },
  {
    key: 'updatedAt',
    sortable: true,
    header: t('system.configuredConnections.list.columns.updatedAt'),
    class: 'text-right',
    pt: {
      columnHeaderContent: 'justify-end',
    },
  },
  {
    key: 'actions',
    class: 'text-right w-12',
    header: '',
    frozen: true,
    alignFrozen: 'right',
    pt: {
      headerCell: { class: 'border-l-0' },
      bodyCell: { class: 'px-2 py-1 border-l-0' },
    },
  },
]

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

const uniqueDerivedParams = computed(() => {
  const params = connection.value?.derivedParams || []
  const seen = new Set()
  return params.filter(p => {
    if (seen.has(p.name)) return false
    seen.add(p.name)
    return true
  })
})

const configuredConnectionInitialValues = computed(() => {
  return {
    name: activeConfiguredConnection.value?.name || '',
    labels: configuredConnectionRawLabels.value,
  }
})

const configuredConnectionResolver = ref(({ values }) => {
  const errors = {}

  if (!values.name || values.name.trim().length === 0) {
    errors.name = [{ message: t('general.label.required') }]
  }

  try {
    if (values.labels && values.labels.trim().length > 0) {
      JSON.parse(values.labels)
    }
  } catch {
    errors.labels = [{ message: t('system.connections.editor.configurations.invalidJSON') }]
  }

  return { errors }
})

// Resource List
const {
  items: configuredConnectionList,
  loading: configuredConnectionsLoading,
  filter: configuredConnectionsFilter,
  sorting: configuredConnectionsSorting,
  pagination: configuredConnectionsPagination,
  handleSort: handleConfiguredConnectionsSort,
  handlePageChange: handleConfiguredConnectionsPageChange,
  filterList: filterConfiguredConnectionsList,
} = useResourceList(
  params => {
    return $SystemAPI.configuredConnectionListCancellable({
      ...params,
      connectionID: route.params.connectionID,
    })
  },
  {
    filter: { query: '' },
    sorting: { sortBy: 'createdAt', sortDesc: true },
    pagination: { limit: 10 },
    immediate: false,
  },
)

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

function initParamValues(existingParams = []) {
  // Clear old values
  Object.keys(paramValues).forEach(k => delete paramValues[k])

  for (const dp of uniqueDerivedParams.value) {
    // Find existing value if editing
    const existing = existingParams.find(p => p.name === dp.name)
    paramValues[dp.name] = existing?.value || dp.default || ''
  }
}

function collectParamValues() {
  const params = []
  for (const dp of uniqueDerivedParams.value) {
    const value = paramValues[dp.name] || ''
    if (value) {
      params.push({
        scope: dp.scope,
        name: dp.name,
        value,
      })
    }
  }
  return params
}

function parseConfiguredConnectionLabels() {
  try {
    const parsed = JSON.parse(configuredConnectionRawLabels.value || '{}')
    if (parsed) {
      activeConfiguredConnection.value.labels = parsed
    }
  } catch {
    // Leave alone, parser failure will be caught by resolver
  }
}

function handleConfiguredConnectionClick({ data }) {
  activeConfiguredConnection.value = { ...data }
  configuredConnectionRawLabels.value = JSON.stringify(data.labels || {}, null, 2)
  initParamValues(data.config?.params || [])
  configuredConnectionModal.value = true
}

// Configured Connection Actions menu methods
function toggleConfiguredConnectionActionsMenu(event, connection) {
  configuredConnectionActionsMenuItems.value = getConfiguredConnectionActionsMenuItems(connection)
  configuredConnectionActionsMenu.value.toggle(event)
}

function getConfiguredConnectionActionsMenuItems(connection) {
  const items = []

  items.push({
    label: t('general.label.delete'),
    icon: 'pi pi-trash',
    command: () => onConfirmConfiguredConnectionDelete(connection),
  })

  return items
}

function onConfirmConfiguredConnectionDelete(connection) {
  confirmDelete({
    message: t('system.configuredConnections.list.deleteConfirm'),
    header: connection.name || t('system.configuredConnections.list.resourceSingle'),
    onConfirm: () => handleConfiguredConnectionDelete(connection),
  })
}

async function handleConfiguredConnectionDelete(connection) {
  try {
    await $SystemAPI.configuredConnectionDelete({
      connectionID: connection.configurationID,
    })
    $toast.toastSuccess(t('notification.connection.delete.success'))
    filterConfiguredConnectionsList()
  } catch (e) {
    console.error('Failed to delete configured connection:', e)
    $toast.toastErrorHandler(t('notification.connection.delete.error'))(e)
  }
}

async function handleConfiguredConnectionDeleteFromModal() {
  try {
    await $SystemAPI.configuredConnectionDelete({
      connectionID: activeConfiguredConnection.value.configurationID,
    })
    $toast.toastSuccess(t('notification.connection.delete.success'))
    configuredConnectionModal.value = false
    filterConfiguredConnectionsList()
  } catch (e) {
    console.error('Failed to delete configured connection:', e)
    $toast.toastErrorHandler(t('notification.connection.delete.error'))(e)
  }
}

async function handleConfiguredConnectionCheck() {
  if (!activeConfiguredConnection.value?.configurationID) return

  checkingConfiguredConnection.value = true
  try {
    const result = await $SystemAPI.configuredConnectionCheck({
      connectionID: activeConfiguredConnection.value.configurationID,
    })

    const allOk = result.connectivity?.ok && result.auth?.ok && (!result.probe || result.probe.ok)
    if (allOk) {
      $toast.toastSuccess(t('system.configuredConnections.editor.checkResult.success'))
    } else {
      const issues = []
      if (!result.connectivity?.ok)
        issues.push(
          `${t('system.configuredConnections.editor.checkResult.connectivity')}: ${result.connectivity?.message || 'failed'}`,
        )
      if (!result.auth?.ok)
        issues.push(
          `${t('system.configuredConnections.editor.checkResult.auth')}: ${result.auth?.message || 'failed'}`,
        )
      if (result.probe && !result.probe.ok)
        issues.push(
          `${t('system.configuredConnections.editor.checkResult.probe')}: ${result.probe?.message || 'failed'}`,
        )
      $toast.toastWarning(issues.join('; '))
    }
  } catch (e) {
    console.error('Failed to check configured connection:', e)
    $toast.toastErrorHandler(t('system.configuredConnections.editor.checkResult.error'))(e)
  } finally {
    checkingConfiguredConnection.value = false
  }
}

async function handleConfiguredConnectionEnable() {
  if (!activeConfiguredConnection.value?.configurationID) return

  enablingConfiguredConnection.value = true
  try {
    // Check first
    const result = await $SystemAPI.configuredConnectionCheck({
      connectionID: activeConfiguredConnection.value.configurationID,
    })

    const allOk = result.connectivity?.ok && result.auth?.ok && (!result.probe || result.probe.ok)
    if (!allOk) {
      const issues = []
      if (!result.connectivity?.ok)
        issues.push(
          `${t('system.configuredConnections.editor.checkResult.connectivity')}: ${result.connectivity?.message || 'failed'}`,
        )
      if (!result.auth?.ok)
        issues.push(
          `${t('system.configuredConnections.editor.checkResult.auth')}: ${result.auth?.message || 'failed'}`,
        )
      if (result.probe && !result.probe.ok)
        issues.push(
          `${t('system.configuredConnections.editor.checkResult.probe')}: ${result.probe?.message || 'failed'}`,
        )
      $toast.toastWarning(issues.join('; '))
      return
    }

    // All checks passed, enable
    await $SystemAPI.configuredConnectionEnable({
      connectionID: activeConfiguredConnection.value.configurationID,
    })
    $toast.toastSuccess(t('notification.connection.update.success'))
    configuredConnectionModal.value = false
    filterConfiguredConnectionsList()
  } catch (e) {
    console.error('Failed to enable configured connection:', e)
    $toast.toastErrorHandler(t('notification.connection.update.error'))(e)
  } finally {
    enablingConfiguredConnection.value = false
  }
}

function createConfiguredConnection() {
  activeConfiguredConnection.value = {
    name: '',
    config: { params: [] },
    labels: {},
  }
  configuredConnectionRawLabels.value = '{\n  \n}'
  initParamValues([])
  configuredConnectionModal.value = true
}

async function handleConfiguredConnectionSubmit({ valid }) {
  if (!valid) return

  savingConfiguredConnection.value = true

  try {
    // Collect param values into config
    activeConfiguredConnection.value.config = {
      ...activeConfiguredConnection.value.config,
      params: collectParamValues(),
    }
    parseConfiguredConnectionLabels()

    if (activeConfiguredConnection.value.configurationID) {
      await $SystemAPI.connectionUpdateConfiguration({
        connectionID: route.params.connectionID,
        configuredConnectionID: activeConfiguredConnection.value.configurationID,
        name: activeConfiguredConnection.value.name,
        config: activeConfiguredConnection.value.config,
        labels: activeConfiguredConnection.value.labels,
      })
      $toast.toastSuccess(t('notification.connection.update.success'))
    } else {
      await $SystemAPI.connectionConfigure({
        connectionID: route.params.connectionID,
        name: activeConfiguredConnection.value.name,
        config: activeConfiguredConnection.value.config,
        labels: activeConfiguredConnection.value.labels,
      })
      $toast.toastSuccess(t('notification.connection.create.success'))
    }

    configuredConnectionModal.value = false
    filterConfiguredConnectionsList()
  } catch (e) {
    console.error('Failed to save configured connection', e)
    $toast.toastErrorHandler(t('notification.connection.update.error'))(e)
  } finally {
    savingConfiguredConnection.value = false
  }
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
    return
  }

  loading.value = true
  try {
    const raw = await $SystemAPI.connectionRead({ connectionID })
    connection.value = new system.Connection(raw)
    initJSONFields()
  } catch (e) {
    console.error('Failed to load connection:', e)
    $toast.toastErrorHandler(t('notification.connection.fetch.error'))(e)
    router.push({ name: 'system.connections' })
  } finally {
    loading.value = false
  }

  if (isEdit.value) {
    filterConfiguredConnectionsList()
  }
}

async function handleSubmit({ valid }) {
  if (!valid) return

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
      $toast.toastSuccess(t('notification.connection.update.success'))
    } else {
      const created = await $SystemAPI.connectionCreate(payload)
      $toast.toastSuccess(t('notification.connection.create.success'))
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
