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
      resourceSingle: $t('general.label.configured-connection.single'),
      resourcePlural: $t('general.label.configured-connection.plural'),
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
          v-if="
            canGrant &&
            connection?.connectionID &&
            (connection.source === 'catalog' ? connection.status === 'active' : true)
          "
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

    <template #body-changedAt="{ data }">
      {{ changedAtText(data) }}
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

  <!-- One popup shared by every row's actions button -->
  <Menu ref="configuredConnectionActionsMenu" :model="configuredConnectionActionsMenuItems" popup>
    <template #item="{ item, props }">
      <a v-ripple v-bind="props.action" :class="item.class">
        <span :class="item.icon" />
        <span class="ml-2">{{ item.label }}</span>
      </a>
    </template>
  </Menu>

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
          v-if="hasMethodChoice"
          :label="$t('system.configuredConnections.editor.authMethod.label')"
        >
          <div class="flex flex-col gap-2">
            <div
              v-for="opt in orderedAuthOptions"
              :key="opt.method"
              class="flex items-start gap-3 rounded-md border p-3 cursor-pointer"
              :class="
                selectedAuthMethod === opt.method
                  ? 'border-primary bg-primary-50'
                  : 'border-surface-200 dark:border-surface-700'
              "
              @click="selectedAuthMethod = opt.method"
            >
              <RadioButton
                v-model="selectedAuthMethod"
                :value="opt.method"
                :input-id="`am-${opt.method}`"
              />
              <div class="flex flex-col">
                <label :for="`am-${opt.method}`" class="font-medium cursor-pointer">
                  {{ authMethodLabel(opt) }}
                </label>
                <span class="text-xs text-muted-color">{{ authMethodDescription(opt) }}</span>
              </div>
            </div>
          </div>
        </CFormGroup>

        <template v-if="!isOAuthMethod(selectedAuthMethod)">
          <CFormGroup
            v-for="param in uniqueDerivedParams"
            :key="param.name"
            :label="param.label || param.name"
            :description="param.description || ''"
            :required="param.required"
            :input-id="`param-${param.name}`"
          >
            <Select
              v-if="fieldKind(param) === 'select'"
              :input-id="`param-${param.name}`"
              v-model="paramValues[param.name]"
              :options="param.options"
              :placeholder="param.default || ''"
              show-clear
              class="w-full"
            />
            <Password
              v-else-if="fieldKind(param) === 'password'"
              :input-id="`param-${param.name}`"
              v-model="paramValues[param.name]"
              toggle-mask
              :feedback="false"
              input-class="w-full"
              class="w-full"
            />
            <Textarea
              v-else-if="fieldKind(param) === 'textarea'"
              :id="`param-${param.name}`"
              v-model="paramValues[param.name]"
              rows="4"
              auto-resize
              class="w-full"
            />
            <InputNumber
              v-else-if="fieldKind(param) === 'number'"
              :input-id="`param-${param.name}`"
              v-model="paramValues[param.name]"
              class="w-full"
            />
            <InputText
              v-else
              :id="`param-${param.name}`"
              v-model="paramValues[param.name]"
              class="w-full"
            />
          </CFormGroup>
        </template>

        <div v-else class="flex flex-col gap-3">
          <div v-if="isOAuthConnected" class="flex items-center gap-2 text-green-600">
            <span class="pi pi-check-circle" />
            <span>{{ $t('system.configuredConnections.editor.oauth.connected') }}</span>
          </div>
          <p class="text-sm text-muted-color">
            {{ $t('system.configuredConnections.editor.oauth.help') }}
          </p>
          <div v-if="selectedOption?.scopes?.length" class="text-xs text-muted-color">
            {{ $t('system.configuredConnections.editor.oauth.scopes') }}:
            {{ selectedOption.scopes.join(', ') }}
          </div>
        </div>
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

          <template v-if="activeConfiguredConnection.status === 'active'">
            <Divider layout="vertical" />
            <Button
              :label="$t('system.configuredConnections.editor.refreshResources', 'Refresh')"
              size="small"
              severity="secondary"
              outlined
              :loading="refreshingDiscovery"
              @click="handleRefreshDiscovery"
            />
          </template>

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
            v-if="isOAuthMethod(selectedAuthMethod)"
            :label="
              isOAuthConnected
                ? $t('system.configuredConnections.editor.oauth.reconnect')
                : $t('system.configuredConnections.editor.oauth.connect')
            "
            icon="pi pi-link"
            size="small"
            :loading="oauthConnecting"
            @click="handleOAuthConnect"
          />
          <Button
            v-else
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
import {
  changedAtField,
  changedAtText,
  components,
  useConfirmDelete,
  useResourceList,
  useRBACStore,
} from '@planetcrust/human-vue'

const { CInputDelete, CResourceList } = components

const props = defineProps({
  connection: { type: Object, required: true },
})

const router = useRouter()
const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')
const $Auth = inject('$Auth')

const rbac = useRBACStore()
const canGrant = computed(() => rbac.can('system/', 'grant'))

const configuredConnectionModal = ref(false)
const configuredConnectionFormRef = ref(null)
const savingConfiguredConnection = ref(false)
const enablingConfiguredConnection = ref(false)
const checkingConfiguredConnection = ref(false)
const refreshingDiscovery = ref(false)
const oauthConnecting = ref(false)
const selectedAuthMethod = ref('')
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
  changedAtField(t('general.columns.changedAt')),
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

// Pick the input control for a param from its declared type / options / name.
function fieldKind(param) {
  if (param.options && param.options.length) return 'select'
  const t = String(param.type || '').toLowerCase()
  if (['password', 'secret'].includes(t)) return 'password'
  if (['textarea', 'multiline', 'json'].includes(t)) return 'textarea'
  if (['number', 'integer', 'int', 'float'].includes(t)) return 'number'
  // Fall back to a masked field for secret-looking names.
  if (/(password|secret|token|apikey|api_key|private_?key|client_?secret)/i.test(param.name)) {
    return 'password'
  }
  return 'text'
}

// Pull a readable string out of an API error so toasts never show [object Object].
// Returns '' when nothing usable is found — the translated prefix stands alone.
function extractError(e) {
  if (!e) return ''
  if (typeof e === 'string') return e
  const cand = e.message ?? e.error ?? e.response?.data?.error?.message ?? e.response?.data?.error
  if (typeof cand === 'string') return cand
  if (cand && typeof cand === 'object') return cand.message || JSON.stringify(cand)
  const s = e.toString?.()
  return s && s !== '[object Object]' ? s : ''
}

// Show an error toast: translated prefix, plus the server detail when present.
function toastError(prefixKey, e) {
  const prefix = t(prefixKey)
  const detail = extractError(e)
  $toast.toastDanger(detail ? `${prefix}: ${detail}` : prefix)
}

const uniqueDerivedParams = computed(() => {
  const params = props.connection?.derivedParams || []
  const seen = new Set()
  return params.filter(p => {
    if (seen.has(p.name)) return false
    seen.add(p.name)
    return true
  })
})

const OAUTH_METHOD = 'oauth2_authorization_code'

// Auth options the connector offers. A picker shows only when there is a choice.
const authOptions = computed(() => props.connection?.service?.authOptions || [])
const hasMethodChoice = computed(() => authOptions.value.length > 1)

// Order the picker so the OAuth "Connect" option leads.
const orderedAuthOptions = computed(() =>
  [...authOptions.value].sort((a, b) => {
    if (a.method === OAUTH_METHOD) return -1
    if (b.method === OAUTH_METHOD) return 1
    return 0
  }),
)

const selectedOption = computed(() =>
  authOptions.value.find(o => o.method === selectedAuthMethod.value),
)

function isOAuthMethod(method) {
  return method === OAUTH_METHOD
}

// Provider name for OAuth labels: the app slug or the connector's short name.
const oauthProvider = computed(() => {
  const opt = authOptions.value.find(o => isOAuthMethod(o.method))
  return opt?.oauthApp || props.connection?.meta?.short || 'provider'
})

function authMethodLabel(opt) {
  return isOAuthMethod(opt.method)
    ? t('system.configuredConnections.editor.authMethod.oauthLabel', {
        provider: oauthProvider.value,
      })
    : t('system.configuredConnections.editor.authMethod.serviceAccountLabel')
}

function authMethodDescription(opt) {
  return isOAuthMethod(opt.method)
    ? t('system.configuredConnections.editor.authMethod.oauthDescription')
    : t('system.configuredConnections.editor.authMethod.serviceAccountDescription')
}

// A configured connection is connected when it holds an OAuth credential.
const isOAuthConnected = computed(() => {
  const cfg = activeConfiguredConnection.value?.config
  return cfg?.authMethod === OAUTH_METHOD && cfg?.credentialID && cfg.credentialID !== '0'
})

// Choose the method to preselect: the saved one, else the OAuth "Connect"
// method when offered, else the first option, else the connector default.
function defaultAuthMethod(cc) {
  if (cc?.config?.authMethod) return cc.config.authMethod
  const oauth = authOptions.value.find(o => isOAuthMethod(o.method))
  if (oauth) return oauth.method
  if (authOptions.value.length) return authOptions.value[0].method
  return props.connection?.service?.auth?.method || ''
}

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
    pagination: { limit: 10 },
    immediate: false,
  },
)

function initParamValues(existingParams = []) {
  Object.keys(paramValues).forEach(k => delete paramValues[k])
  for (const dp of uniqueDerivedParams.value) {
    const existing = existingParams.find(p => p.name === dp.name)
    let value = existing?.value ?? dp.default ?? ''
    if (fieldKind(dp) === 'number' && value !== '') value = Number(value)
    paramValues[dp.name] = value
  }
}

function collectParamValues() {
  const params = []
  for (const dp of uniqueDerivedParams.value) {
    const raw = paramValues[dp.name]
    const value = raw === null || raw === undefined ? '' : String(raw)
    if (value !== '') {
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
  selectedAuthMethod.value = defaultAuthMethod(data)
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
  selectedAuthMethod.value = defaultAuthMethod(null)
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
  // show() rather than toggle(): one popup serves every row, so clicking a
  // second row's button must re-anchor and open there rather than close the first.
  configuredConnectionActionsMenu.value.show(event, event.currentTarget)
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
    toastError('notification.connection.delete.error', e)
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
    toastError('notification.connection.delete.error', e)
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
    toastError('system.configuredConnections.editor.checkResult.error', e)
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
    toastError('system.configuredConnections.editor.checkResult.error', e)
  } finally {
    checkingConfiguredConnection.value = false
  }
}

// Re-run resource discovery. Called through the axios instance since the
// typed client method may not be generated yet.
async function handleRefreshDiscovery() {
  if (!activeConfiguredConnection.value?.configurationID) return
  refreshingDiscovery.value = true
  try {
    await $SystemAPI
      .api()
      .post(
        `/configured-connections/${activeConfiguredConnection.value.configurationID}/refresh-discovery`,
      )
    $toast.toastSuccess(
      t('system.configuredConnections.editor.discoveryRefreshed', 'Refreshed available resources'),
    )
  } catch (e) {
    const prefix = t(
      'system.configuredConnections.editor.discoveryError',
      'Could not refresh resources',
    )
    const detail = extractError(e)
    $toast.toastDanger(detail ? `${prefix}: ${detail}` : prefix)
  } finally {
    refreshingDiscovery.value = false
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

// Return a configuredConnectionID for the authorize endpoint to bind the
// credential to. An existing connection is used as-is; a new one is saved as a
// draft first.
async function ensureSavedForConnect() {
  const cc = activeConfiguredConnection.value
  if (cc.configurationID) {
    return cc.configurationID
  }

  const config = { ...(cc.config || {}), authMethod: selectedAuthMethod.value, params: [] }
  parseConfiguredConnectionLabels()

  const saved = await $SystemAPI.connectionConfigure({
    connectionID: props.connection.connectionID,
    name: cc.name,
    config,
    labels: cc.labels,
  })
  cc.configurationID = saved?.configurationID
  cc.connectionID = saved?.connectionID
  cc.config = { ...config, credentialID: saved?.config?.credentialID }
  return saved?.configurationID
}

// True once the callback has linked a credential AND bumped updatedAt past the
// pre-consent baseline. The updatedAt guard matters for Reconnect, where a
// credential already exists before consent runs.
async function isReconnectComplete(configurationID, baselineUpdatedAt) {
  try {
    const cc = await $SystemAPI.configuredConnectionRead({ connectionID: configurationID })
    const id = cc?.config?.credentialID
    return !!id && id !== '0' && cc?.updatedAt !== baselineUpdatedAt
  } catch {
    return false
  }
}

// Open the provider consent popup and resolve true on success. The auth host is
// a different origin than the SPA, so the popup's URL cannot be read. Instead we
// poll the configured connection for the callback's update and close the popup
// once it lands. Origin-independent by design.
function runOAuthPopup(configurationID) {
  const authBase = ($Auth?.authURL || `${window.location.origin}/auth`).replace(/\/$/, '')
  const url = `${authBase}/oauth2/connection/authorize?configuredConnectionID=${configurationID}`
  const baselineUpdatedAt = activeConfiguredConnection.value?.updatedAt || ''
  const popup = window.open(url, 'oauth2-connect', 'width=520,height=680')
  if (!popup) {
    $toast.toastWarning(t('system.configuredConnections.editor.oauth.popupBlocked'))
    return Promise.resolve(false)
  }
  return new Promise(resolve => {
    let settled = false
    const finish = ok => {
      if (settled) return
      settled = true
      clearInterval(timer)
      clearTimeout(timeout)
      try {
        popup.close()
      } catch {
        // Popup may already be closed.
      }
      resolve(ok)
    }
    const timer = setInterval(async () => {
      if (await isReconnectComplete(configurationID, baselineUpdatedAt)) {
        finish(true)
        return
      }
      // User closed the popup — check once more, then give up.
      if (popup.closed) {
        finish(await isReconnectComplete(configurationID, baselineUpdatedAt))
      }
    }, 1500)
    // Safety net so a stalled consent never spins forever.
    const timeout = setTimeout(() => finish(false), 5 * 60 * 1000)
  })
}

async function handleOAuthConnect() {
  if (!activeConfiguredConnection.value?.name?.trim()) {
    $toast.toastWarning(t('general.notification.formErrors'))
    return
  }
  oauthConnecting.value = true
  try {
    const configurationID = await ensureSavedForConnect()
    if (!configurationID) return

    const ok = await runOAuthPopup(configurationID)
    if (!ok) {
      $toast.toastWarning(t('system.configuredConnections.editor.oauth.denied'))
      return
    }

    // An active connection is a reconnect: the credential refreshed, and
    // enabling it again is rejected — so skip to success. A draft still needs
    // check + enable to go active.
    if (activeConfiguredConnection.value?.status === 'active') {
      $toast.toastSuccess(t('system.configuredConnections.editor.oauth.success'))
    } else if (await checkAndEnableConfiguredConnection(configurationID)) {
      $toast.toastSuccess(t('system.configuredConnections.editor.oauth.success'))
    }
    configuredConnectionModal.value = false

    // A fresh catalog configure imports the connector under a new connectionID;
    // route there so the list shows the new configured connection.
    const newConnectionID = activeConfiguredConnection.value?.connectionID
    if (newConnectionID && newConnectionID !== props.connection.connectionID) {
      router.push({
        name: 'system.connections.configure',
        params: { connectionID: newConnectionID },
      })
    } else {
      filterConfiguredConnectionsList()
    }
  } catch (e) {
    toastError('notification.connection.update.error', e)
  } finally {
    oauthConnecting.value = false
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
      authMethod: selectedAuthMethod.value,
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
    toastError('notification.connection.update.error', e)
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
