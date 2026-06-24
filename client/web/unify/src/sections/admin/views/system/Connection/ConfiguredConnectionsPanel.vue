<template>
  <CResourceList
    ref="resourceListRef"
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
      <div class="flex gap-2">
        <Button
          :label="$t('system.configuredConnections.list.createLabel')"
          icon="pi pi-plus"
          size="small"
          @click="createConfiguredConnection"
        />
        <CPermissionsButton
          v-if="canGrant && connection?.connectionID && (connection.source === 'catalog' ? connection.status === 'active' : true)"
          v-tooltip.bottom="$t('general.label.permissions')"
          resource="corteza::system:configured-connection/*"
          :title="connection.meta?.short || connection.handle || connection.connectionID"
          :target="connection.meta?.short || connection.handle || connection.connectionID"
        />
      </div>
    </template>

    <template #body-status="{ data }">
      <Tag
        :value="
          $t(`system.connections.editor.statusValues.${data.status || 'draft'}`, data.status || '-')
        "
        :severity="data.status === 'active' ? 'success' : 'secondary'"
      />
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

  <TieredMenu
    ref="configuredConnectionActionsMenu"
    :model="configuredConnectionActionsMenuItems"
    popup
  >
    <template #item="{ item, props }">
      <a v-ripple v-bind="props.action" :class="item.class">
        <span :class="item.icon" />
        <span class="ml-2">{{ item.label }}</span>
      </a>
    </template>
  </TieredMenu>

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
      ref="configuredConnectionFormRef"
      :resolver="configuredConnectionResolver"
      :initialValues="configuredConnectionInitialValues"
      @submit="handleConfiguredConnectionSubmit"
      class="flex flex-col h-full min-h-0"
    >
      <div class="p-4 flex flex-col gap-4 flex-1 min-h-0 overflow-y-auto">
        <CFormGroup
          name="name"
          input-id="ccName"
          :label="$t('system.configuredConnections.editor.info.name')"
          required
        >
          <InputText id="ccName" name="name" v-model="activeConfiguredConnection.name" />
        </CFormGroup>

        <CFormGroup
          v-if="activeConfiguredConnection.configurationID"
          :label="$t('system.configuredConnections.editor.info.status')"
        >
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
        </CFormGroup>

        <CFormGroup
          v-for="param in uniqueDerivedParams"
          :key="param.name"
          :label="param.label || param.name"
          :description="param.description || ''"
          :required="param.required"
          :input-id="`param-${param.name}`"
        >
          <InputText :id="`param-${param.name}`" v-model="paramValues[param.name]" />
        </CFormGroup>
      </div>
    </Form>

    <template #footer>
      <div class="flex items-center gap-2 w-full">
        <div v-if="activeConfiguredConnection.configurationID" class="flex items-center gap-2">
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

          <Divider layout="vertical" />
          <Button
            :label="$t('system.configuredConnections.editor.check')"
            size="small"
            severity="info"
            outlined
            :loading="checkingConfiguredConnection"
            @click="handleConfiguredConnectionCheck"
          />

          <template v-if="activeConfiguredConnection.status !== 'active'">
            <Divider layout="vertical" />
            <Button
              :label="$t('system.configuredConnections.editor.enable')"
              size="small"
              severity="success"
              outlined
              :loading="enablingConfiguredConnection"
              @click="handleConfiguredConnectionEnable"
            />
          </template>
        </div>

        <div class="flex gap-2 ml-auto">
          <Button
            :label="$t('general.label.cancel')"
            severity="secondary"
            text
            size="small"
            @click="configuredConnectionModal = false"
          />
          <Button
            :label="$t('general.label.save')"
            size="small"
            :loading="savingConfiguredConnection"
            @click="configuredConnectionFormRef?.submit?.()"
          />
        </div>
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { inject, nextTick, onMounted, reactive, ref, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { components, useConfirmDelete, useResourceList, useRBACStore } from '@planetcrust/human-vue'

const { CInputDelete, CResourceList } = components

const props = defineProps({
  connection: { type: Object, required: true },
})

const router = useRouter()
const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const rbac = useRBACStore()
const canGrant = computed(() => rbac.can('system/', 'grant'))

const configuredConnectionModal = ref(false)
const configuredConnectionFormRef = ref(null)
const savingConfiguredConnection = ref(false)
const enablingConfiguredConnection = ref(false)
const checkingConfiguredConnection = ref(false)
const activeConfiguredConnection = ref(null)
const configuredConnectionRawLabels = ref('{}')
const paramValues = reactive({})
const configuredConnectionActionsMenu = ref()
const configuredConnectionActionsMenuItems = ref([])
const resourceListRef = ref()

const configuredConnectionFields = [
  {
    key: 'name',
    sortable: true,
    header: t('system.configuredConnections.list.columns.name'),
  },
  {
    key: 'status',
    sortable: true,
    header: t('system.connections.list.columns.status'),
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

const uniqueDerivedParams = computed(() => {
  const params = props.connection?.derivedParams || []
  const seen = new Set()
  return params.filter(p => {
    if (seen.has(p.name)) return false
    seen.add(p.name)
    return true
  })
})

const configuredConnectionInitialValues = computed(() => ({
  name: activeConfiguredConnection.value?.name || '',
  labels: configuredConnectionRawLabels.value,
}))

const configuredConnectionResolver = ({ values }) => {
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
}

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
      connectionID: props.connection?.connectionID,
    })
  },
  {
    filter: { query: '' },
    sorting: { sortBy: 'status', sortDesc: false },
    pagination: { limit: 10 },
    immediate: false,
  },
)

function initParamValues(existingParams = []) {
  Object.keys(paramValues).forEach(k => delete paramValues[k])
  for (const dp of uniqueDerivedParams.value) {
    const existing = existingParams.find(p => p.name === dp.name)
    paramValues[dp.name] = existing?.value || dp.default || ''
  }
}

function collectParamValues() {
  const params = []
  for (const dp of uniqueDerivedParams.value) {
    const value = paramValues[dp.name] || ''
    if (value) {
      params.push({ scope: dp.scope, name: dp.name, value })
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
    // resolver catches invalid JSON
  }
}

function handleConfiguredConnectionClick({ data }) {
  activeConfiguredConnection.value = { ...data }
  configuredConnectionRawLabels.value = JSON.stringify(data.labels || {}, null, 2)
  initParamValues(data.config?.params || [])
  configuredConnectionModal.value = true
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

function toggleConfiguredConnectionActionsMenu(event, conn) {
  configuredConnectionActionsMenuItems.value = [
    {
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => onConfirmConfiguredConnectionDelete(conn),
    },
  ]
  configuredConnectionActionsMenu.value.toggle(event)
}

function onConfirmConfiguredConnectionDelete(conn) {
  confirmDelete({
    message: t('system.configuredConnections.list.deleteConfirm'),
    header: conn.name || t('system.configuredConnections.list.resourceSingle'),
    onConfirm: () => handleConfiguredConnectionDelete(conn),
  })
}

async function handleConfiguredConnectionDelete(conn) {
  try {
    await $SystemAPI.configuredConnectionDelete({ connectionID: conn.configurationID })
    $toast.toastSuccess(t('notification.connection.delete.success'))
    filterConfiguredConnectionsList()
  } catch (e) {
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
    $toast.toastErrorHandler(t('notification.connection.delete.error'))(e)
  }
}

function summarizeCheckIssues(result) {
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
  return issues
}

function isCheckResultOk(result) {
  return result.connectivity?.ok && result.auth?.ok && (!result.probe || result.probe.ok)
}

async function checkAndEnableConfiguredConnection(configurationID) {
  if (!configurationID) return false
  try {
    const result = await $SystemAPI.configuredConnectionCheck({ connectionID: configurationID })
    if (!isCheckResultOk(result)) {
      $toast.toastWarning(summarizeCheckIssues(result).join('; '))
      return false
    }
    await $SystemAPI.configuredConnectionEnable({ connectionID: configurationID })
    return true
  } catch (e) {
    $toast.toastErrorHandler(t('system.configuredConnections.editor.checkResult.error'))(e)
    return false
  }
}

async function handleConfiguredConnectionCheck() {
  if (!activeConfiguredConnection.value?.configurationID) return
  checkingConfiguredConnection.value = true
  try {
    const result = await $SystemAPI.configuredConnectionCheck({
      connectionID: activeConfiguredConnection.value.configurationID,
    })
    if (isCheckResultOk(result)) {
      $toast.toastSuccess(t('system.configuredConnections.editor.checkResult.success'))
    } else {
      $toast.toastWarning(summarizeCheckIssues(result).join('; '))
    }
  } catch (e) {
    $toast.toastErrorHandler(t('system.configuredConnections.editor.checkResult.error'))(e)
  } finally {
    checkingConfiguredConnection.value = false
  }
}

async function handleConfiguredConnectionEnable() {
  if (!activeConfiguredConnection.value?.configurationID) return
  enablingConfiguredConnection.value = true
  try {
    const enabled = await checkAndEnableConfiguredConnection(
      activeConfiguredConnection.value.configurationID,
    )
    if (enabled) {
      $toast.toastSuccess(t('notification.connection.update.success'))
      configuredConnectionModal.value = false
      filterConfiguredConnectionsList()
    }
  } finally {
    enablingConfiguredConnection.value = false
  }
}

async function handleConfiguredConnectionSubmit({ valid }) {
  if (!valid) {
    $toast.toastWarning(t('general.notification.formErrors'))
    nextTick(() => {
      document
        .querySelector('.p-message-error')
        ?.scrollIntoView({ behavior: 'smooth', block: 'center' })
    })
    return
  }

  savingConfiguredConnection.value = true
  try {
    activeConfiguredConnection.value.config = {
      ...activeConfiguredConnection.value.config,
      params: collectParamValues(),
    }
    parseConfiguredConnectionLabels()

    if (activeConfiguredConnection.value.configurationID) {
      await $SystemAPI.connectionUpdateConfiguration({
        connectionID: props.connection.connectionID,
        configuredConnectionID: activeConfiguredConnection.value.configurationID,
        name: activeConfiguredConnection.value.name,
        config: activeConfiguredConnection.value.config,
        labels: activeConfiguredConnection.value.labels,
      })
      $toast.toastSuccess(t('notification.connection.update.success'))
      await checkAndEnableConfiguredConnection(activeConfiguredConnection.value.configurationID)
      configuredConnectionModal.value = false
      filterConfiguredConnectionsList()
    } else {
      const saved = await $SystemAPI.connectionConfigure({
        connectionID: props.connection.connectionID,
        name: activeConfiguredConnection.value.name,
        config: activeConfiguredConnection.value.config,
        labels: activeConfiguredConnection.value.labels,
      })
      $toast.toastSuccess(t('notification.connection.create.success'))
      await checkAndEnableConfiguredConnection(saved?.configurationID)
      configuredConnectionModal.value = false
      const newConnectionID = saved?.connectionID
      if (newConnectionID && newConnectionID !== props.connection.connectionID) {
        router.push({
          name: 'system.connections.configure',
          params: { connectionID: newConnectionID },
        })
      } else {
        filterConfiguredConnectionsList()
      }
    }
  } catch (e) {
    $toast.toastErrorHandler(t('notification.connection.update.error'))(e)
  } finally {
    savingConfiguredConnection.value = false
  }
}

onMounted(() => {
  if (props.connection?.connectionID) {
    filterConfiguredConnectionsList()
  }
})

watch(
  () => props.connection?.connectionID,
  (id, oldId) => {
    if (id && id !== oldId) {
      filterConfiguredConnectionsList()
    }
  },
)
</script>
