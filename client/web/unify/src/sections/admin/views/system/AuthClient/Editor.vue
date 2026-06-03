<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ pageTitle }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <Form
    v-else-if="authClient"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <div v-if="isEdit" class="flex justify-end gap-2 shrink-0">
        <CPermissionsButton
          v-if="authClient.canGrant"
          v-tooltip.bottom="$t('general.label.permissions')"
          :resource="`corteza::system:auth-client/${authClient.authClientID}`"
          :title="authClient.meta?.name || authClient.handle || authClient.authClientID"
          :target="authClient.meta?.name || authClient.handle || authClient.authClientID"
        />
      </div>

      <!-- Info -->
      <Panel :header="$t('system.authclients.editor.info.title')" toggleable :collapsed="false" class="shadow">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <CFormGroup name="name" :label="$t('system.authclients.editor.info.name')" required>
            <InputText id="name" name="name" v-model="authClient.meta.name" />
          </CFormGroup>

          <CFormGroup name="handle" :label="$t('system.authclients.editor.info.handle.label')">
            <InputText id="handle" name="handle" v-model="authClient.handle" />
          </CFormGroup>

          <CFormGroup
            :label="$t('system.authclients.editor.info.redirectURI')"
            class="md:col-span-2"
          >
            <template #actions>
              <Button
                :label="$t('general.label.add')"
                icon="pi pi-plus"
                severity="secondary"
                size="small"
                @click="addURI"
              />
            </template>
            <CFormList
              v-model="redirectURIs"
              :columns="[{ label: $t('system.authclients.editor.info.uri'), width: '1fr' }]"
              @change="syncRedirectURIs"
            >
              <template #row="{ index }">
                <InputText
                  v-model="redirectURIs[index]"
                  size="small"
                  class="w-full"
                  :placeholder="$t('system.authclients.editor.info.redirectURIPlaceholder')"
                  @update:modelValue="syncRedirectURIs"
                />
              </template>
            </CFormList>
          </CFormGroup>

          <CFormGroup
            v-if="isEdit"
            :label="$t('system.authclients.editor.info.secret')"
            class="md:col-span-2"
          >
            <div class="flex items-center gap-2">
              <InputText
                :value="secretVisible ? secret : '••••••••••••••••'"
                readonly
                class="flex-1"
                :type="secretVisible ? 'text' : 'password'"
              />
              <Button
                :icon="secretVisible ? 'pi pi-eye-slash' : 'pi pi-eye'"
                text
                rounded
                size="small"
                severity="secondary"
                @click="secretVisible ? hideSecret() : showSecret()"
              />
              <Button
                icon="pi pi-refresh"
                text
                rounded
                size="small"
                severity="warning"
                :title="$t('system.authclients.editor.info.regenerateSecret')"
                @click="regenerateSecret"
              />
            </div>
          </CFormGroup>

          <CFormGroup :label="$t('system.authclients.editor.info.validGrant')">
            <div class="flex flex-col gap-2">
              <div
                v-for="opt in grantOptions"
                :key="opt.value"
                class="flex items-center gap-2"
              >
                <RadioButton
                  v-model="authClient.validGrant"
                  :inputId="`grant-${opt.value}`"
                  :value="opt.value"
                />
                <label :for="`grant-${opt.value}`" class="cursor-pointer">{{ opt.label }}</label>
              </div>
            </div>
          </CFormGroup>

          <CFormGroup :label="$t('system.authclients.editor.info.scope')">
            <div class="flex flex-col gap-2">
              <div class="flex items-center gap-2">
                <Checkbox inputId="scope-profile" v-model="scopeProfile" :binary="true" />
                <label for="scope-profile" class="cursor-pointer">
                  {{ $t('system.authclients.editor.info.profile') }}
                </label>
              </div>
              <div class="flex items-center gap-2">
                <Checkbox inputId="scope-api" v-model="scopeApi" :binary="true" />
                <label for="scope-api" class="cursor-pointer">
                  {{ $t('system.authclients.editor.info.api') }}
                </label>
              </div>
              <div class="flex items-center gap-2">
                <Checkbox inputId="scope-openid" v-model="scopeOpenid" :binary="true" />
                <label for="scope-openid" class="cursor-pointer">
                  {{ $t('system.authclients.editor.info.openid') }}
                </label>
              </div>
              <div class="flex items-center gap-2">
                <Checkbox inputId="scope-discovery" v-model="scopeDiscovery" :binary="true" />
                <label for="scope-discovery" class="cursor-pointer">
                  {{ $t('system.authclients.editor.info.discovery') }}
                </label>
              </div>
            </div>
          </CFormGroup>

          <CFormGroup
            :label="$t('system.authclients.editor.info.validFrom')"
            :description="$t('system.authclients.editor.info.validFromDescription')"
          >
            <DatePicker v-model="authClient.validFrom" showTime hourFormat="24" showIcon showButtonBar fluid />
          </CFormGroup>

          <CFormGroup
            :label="$t('system.authclients.editor.info.expiresAt')"
            :description="$t('system.authclients.editor.info.expiresAtDescription')"
          >
            <DatePicker v-model="authClient.expiresAt" showTime hourFormat="24" showIcon showButtonBar fluid />
          </CFormGroup>

          <CInputToggleCard
            v-model="authClient.enabled"
            :label="$t('system.authclients.editor.info.enabled.label')"
            :description="$t('system.authclients.editor.info.enabled.description')"
          />

          <CInputToggleCard
            v-model="authClient.trusted"
            :label="$t('system.authclients.editor.info.trusted.label')"
            :description="$t('system.authclients.editor.info.trusted.description')"
          />
        </div>
      </Panel>

      <!-- Security -->
      <Panel :header="$t('system.authclients.editor.tabs.security')" toggleable :collapsed="false" class="shadow">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <CFormGroup
            :label="$t('system.authclients.editor.info.security.permittedRoles.label')"
            :description="$t('system.authclients.editor.info.security.permittedRoles.description')"
          >
            <CInputRole
              :placeholder="$t('system.authclients.editor.info.security.selectRole')"
              clear-on-select
              filter-context-roles
              @select="role => addRoleToList('permittedRoles', role)"
            />
            <CFormItemList
              :items="permittedRoles"
              item-key="roleID"
              @remove="role => removeRoleFromList('permittedRoles', role)"
            >
              <template #default="{ item }">
                <span class="text-muted-color font-medium">
                  {{ item.name || item.handle || item.roleID }}
                </span>
              </template>
            </CFormItemList>
          </CFormGroup>

          <CFormGroup
            :label="$t('system.authclients.editor.info.security.prohibitedRoles.label')"
            :description="$t('system.authclients.editor.info.security.prohibitedRoles.description')"
          >
            <CInputRole
              :placeholder="$t('system.authclients.editor.info.security.selectRole')"
              clear-on-select
              filter-context-roles
              @select="role => addRoleToList('prohibitedRoles', role)"
            />
            <CFormItemList
              :items="prohibitedRoles"
              item-key="roleID"
              @remove="role => removeRoleFromList('prohibitedRoles', role)"
            >
              <template #default="{ item }">
                <span class="text-muted-color font-medium">
                  {{ item.name || item.handle || item.roleID }}
                </span>
              </template>
            </CFormItemList>
          </CFormGroup>

          <CFormGroup
            :label="$t('system.authclients.editor.info.security.forcedRoles.label')"
            :description="$t('system.authclients.editor.info.security.forcedRoles.description')"
          >
            <CInputRole
              :placeholder="$t('system.authclients.editor.info.security.selectRole')"
              clear-on-select
              filter-context-roles
              @select="role => addRoleToList('forcedRoles', role)"
            />
            <CFormItemList
              :items="forcedRoles"
              item-key="roleID"
              @remove="role => removeRoleFromList('forcedRoles', role)"
            >
              <template #default="{ item }">
                <span class="text-muted-color font-medium">
                  {{ item.name || item.handle || item.roleID }}
                </span>
              </template>
            </CFormItemList>
          </CFormGroup>

          <CFormGroup
            v-if="authClient.validGrant === 'client_credentials'"
            :label="$t('system.authclients.editor.info.security.impersonateUser.label')"
            :description="$t('system.authclients.editor.info.security.impersonateUser.description')"
          >
            <CInputUser
              v-model="authClient.security.impersonateUser"
              :placeholder="$t('system.authclients.editor.info.security.impersonateUser.placeholder')"
              class="w-full"
            />
          </CFormGroup>

          <CFormGroup :label="$t('system.authclients.editor.info.security.defaultUserGroup.label')">
            <CInputUserGroup
              v-model="authClient.security.userGroup"
              :placeholder="$t('system.authclients.editor.info.security.defaultUserGroup.placeholder')"
              class="w-full"
            />
          </CFormGroup>
        </div>
      </Panel>

      <!-- Developer (only meaningful for client_credentials) -->
      <Panel
        v-if="isEdit && authClient.validGrant === 'client_credentials'"
        :header="$t('system.authclients.editor.tabs.developer')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <div class="flex flex-col gap-6">
          <CFormGroup :label="$t('system.authclients.editor.info.curl')">
            <template #actions>
              <Button icon="pi pi-copy" text size="small" severity="secondary" @click="copyToClipboard(curlExample)" />
            </template>
            <Textarea :value="curlExample" readonly rows="3" class="font-mono text-sm" />
          </CFormGroup>

          <div class="flex flex-col gap-2">
            <CFormGroup
              v-if="tokenRequest.token"
              :label="$t('system.authclients.editor.info.accessToken')"
            >
              <template #actions>
                <Button icon="pi pi-copy" text size="small" severity="secondary" @click="copyToClipboard(tokenRequest.token)" />
              </template>
              <Textarea :value="tokenRequest.token" readonly rows="4" class="font-mono text-sm" />
            </CFormGroup>
            <div>
              <Button
                :label="$t('system.authclients.editor.info.generateAccessToken')"
                severity="secondary"
                :loading="tokenRequest.loading"
                @click="generateToken"
              />
              <p v-if="tokenRequest.error" class="text-red-500 text-sm mt-2">{{ tokenRequest.error }}</p>
            </div>
          </div>
        </div>
      </Panel>
    </div>

    <CEditorActions :back-to="{ name: 'system.authClients' }">
      <CInputDelete
        v-if="isEdit && authClient.canDeleteAuthClient"
        :label="$t('system.authclients.editor.info.delete')"
        :message="$t('general.confirm.delete')"
        :header="authClient.meta?.name || authClient.handle"
        :disabled="deleting"
        @confirm="handleDelete"
      />
      <Button
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
import { system, NoID } from '@planetcrust/human-js'
import { components, useUnsavedGuard } from '@planetcrust/human-vue'
import { cloneDeep, isEqual } from 'lodash-es'
import axios from 'axios'

const { CInputDelete, CInputRole, CInputUser, CInputUserGroup, CInputToggleCard } = components

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')
const $Auth = inject('$Auth')

const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const authClient = ref(null)
const initialAuthClient = ref(null)

// Local role lists with full role objects for display
const permittedRoles = ref([])
const prohibitedRoles = ref([])
const forcedRoles = ref([])

// Secret management
const secret = ref('')
const secretVisible = ref(false)

// Redirect URIs list
const redirectURIs = ref([])

// Token request state
const tokenRequest = ref({ token: '', error: '', loading: false })

const isEdit = computed(() => !!route.params.authClientID)

const pageTitle = computed(() =>
  isEdit.value
    ? t('system.authclients.editor.title.edit')
    : t('system.authclients.editor.title.create'),
)

const grantOptions = computed(() => [
  {
    label: t('system.authclients.editor.info.grant.authorization_code'),
    value: 'authorization_code',
  },
  {
    label: t('system.authclients.editor.info.grant.client_credentials'),
    value: 'client_credentials',
  },
])

const initialValues = computed(() => ({
  name: authClient.value?.meta?.name || '',
  handle: authClient.value?.handle || '',
}))

const resolver = ref(({ values }) => {
  const errors = {}

  if (!values.name || values.name.trim().length === 0) {
    errors.name = [{ message: t('general.label.required') }]
  }

  if (values.handle && !/^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$|^[A-Za-z]$/.test(values.handle)) {
    errors.handle = [
      { message: t('system.authclients.editor.info.handle.invalid-handle-characters') },
    ]
  }

  return { errors }
})

// Scope computed properties
function hasScope(s) {
  return (authClient.value?.scope || '').split(' ').includes(s)
}

function toggleScope(s, val) {
  let parts = (authClient.value?.scope || '').split(' ').filter(Boolean)
  if (val) {
    if (!parts.includes(s)) parts.push(s)
  } else {
    parts = parts.filter(p => p !== s)
  }
  authClient.value.scope = parts.join(' ')
}

const scopeProfile = computed({ get: () => hasScope('profile'), set: v => toggleScope('profile', v) })
const scopeApi = computed({ get: () => hasScope('api'), set: v => toggleScope('api', v) })
const scopeOpenid = computed({ get: () => hasScope('openid'), set: v => toggleScope('openid', v) })
const scopeDiscovery = computed({ get: () => hasScope('discovery'), set: v => toggleScope('discovery', v) })

// Redirect URIs management
function syncRedirectURIs() {
  authClient.value.redirectURI = redirectURIs.value.filter(Boolean).join(' ')
}

function addURI() {
  redirectURIs.value.push('')
}

// Secret management
async function showSecret() {
  if (!secret.value) {
    try {
      const result = await $SystemAPI.authClientExposeSecret({ clientID: authClient.value.authClientID })
      secret.value = result
    } catch (e) {
      $toast.toastErrorHandler('Failed to get secret')(e)
    }
  }
  secretVisible.value = true
}

function hideSecret() {
  secretVisible.value = false
}

async function regenerateSecret() {
  try {
    const result = await $SystemAPI.authClientRegenerateSecret({ clientID: authClient.value.authClientID })
    secret.value = result
    secretVisible.value = true
    $toast.toastSuccess('Secret regenerated')
  } catch (e) {
    $toast.toastErrorHandler('Failed to regenerate secret')(e)
  }
}

// Developer tab — cURL and token generation
const tokenURL = computed(() => `${$Auth?.authURL || ''}/oauth2/token`)

const curlExample = computed(() => {
  if (!authClient.value) return ''
  return `curl -X POST ${tokenURL.value} -d grant_type=${authClient.value.validGrant} -d scope='${authClient.value.scope || ''}' -u ${authClient.value.authClientID}:${secret.value || 'YOUR-CLIENT-SECRET'}`
})

async function generateToken() {
  tokenRequest.value.loading = true
  tokenRequest.value.error = ''
  try {
    if (!secret.value) await showSecret()
    const params = new URLSearchParams()
    params.append('grant_type', authClient.value.validGrant)
    params.append('scope', authClient.value.scope || '')
    const resp = await axios.post(tokenURL.value, params, {
      auth: { username: authClient.value.authClientID, password: secret.value },
    })
    tokenRequest.value.token = resp.data?.access_token || ''
  } catch (e) {
    tokenRequest.value.error = e?.response?.data?.error || e.message
  } finally {
    tokenRequest.value.loading = false
  }
}

function copyToClipboard(text) {
  navigator.clipboard.writeText(text).catch(() => {})
}


function addRoleToList(listName, role) {
  if (!role) return
  const list = { permittedRoles, prohibitedRoles, forcedRoles }[listName]
  if (!list.value.find(r => r.roleID === role.roleID)) {
    list.value.push(role)
    authClient.value.security[listName] = list.value.map(r => r.roleID)
  }
}

function removeRoleFromList(listName, role) {
  const list = { permittedRoles, prohibitedRoles, forcedRoles }[listName]
  const idx = list.value.findIndex(r => r.roleID === role.roleID)
  if (idx !== -1) {
    list.value.splice(idx, 1)
    authClient.value.security[listName] = list.value.map(r => r.roleID)
  }
}

async function loadRolesForList(ids, targetRef) {
  if (!ids || ids.length === 0) return
  const roles = await Promise.all(
    ids.map(id => $SystemAPI.roleRead({ roleID: id }).catch(() => null)),
  )
  targetRef.value = roles.filter(Boolean).map(r => new system.Role(r))
}

async function fetchDefaultUserGroup() {
  try {
    const result = await $SystemAPI.userGroupList({ limit: 100 })
    if (result?.set?.length > 0) {
      const defaultGroup = result.set.find(
        g => g.handle === 'default-root' || g.handle === 'users' || g.meta?.short?.includes('Default'),
      ) || result.set[0]
      if (defaultGroup && authClient.value) {
        authClient.value.security.userGroup = defaultGroup.userGroupID
      }
    }
  } catch (e) {
    console.warn('Silent fail fetching default user group', e)
  }
}

async function loadAuthClient() {
  const authClientID = route.params.authClientID
  if (!authClientID) {
    authClient.value = new system.AuthClient({ enabled: true })
    redirectURIs.value = []
    await fetchDefaultUserGroup()
    initialAuthClient.value = cloneDeep(authClient.value)
    return
  }

  loading.value = true
  try {
    const raw = await $SystemAPI.authClientRead({ clientID: authClientID })
    authClient.value = new system.AuthClient(raw)

    // Initialize redirect URIs list
    redirectURIs.value = (authClient.value.redirectURI || '').split(' ').filter(Boolean)

    // Load role objects for security lists
    await Promise.all([
      loadRolesForList(authClient.value.security.permittedRoles, permittedRoles),
      loadRolesForList(authClient.value.security.prohibitedRoles, prohibitedRoles),
      loadRolesForList(authClient.value.security.forcedRoles, forcedRoles),
    ])
    initialAuthClient.value = cloneDeep(authClient.value)
  } catch (e) {
    $toast.toastErrorHandler(t('notification.authclient.fetch.error'))(e)
    router.push({ name: 'system.authClients' })
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

  saving.value = true
  try {
    // Force impersonateUser to NoID unless client_credentials with a real user picked
    const isCC = authClient.value.validGrant === 'client_credentials'
    if (!isCC || !authClient.value.security.impersonateUser) {
      authClient.value.security.impersonateUser = NoID
    }

    const payload = {
      handle: authClient.value.handle,
      meta: authClient.value.meta,
      scope: authClient.value.scope,
      redirectURI: redirectURIs.value.filter(Boolean).join(' '),
      validGrant: authClient.value.validGrant,
      validFrom: authClient.value.validFrom,
      expiresAt: authClient.value.expiresAt,
      enabled: authClient.value.enabled,
      trusted: authClient.value.trusted,
      security: authClient.value.security,
    }

    if (isEdit.value) {
      payload.clientID = authClient.value.authClientID
      const raw = await $SystemAPI.authClientUpdate(payload)
      authClient.value = new system.AuthClient(raw)
      initialAuthClient.value = cloneDeep(authClient.value)
      $toast.toastSuccess(t('notification.authclient.update.success'))
    } else {
      const created = await $SystemAPI.authClientCreate(payload)
      $toast.toastSuccess(t('notification.authclient.create.success'))
      markSaved()
      router.push({
        name: 'system.authClients.edit',
        params: { authClientID: created.authClientID },
      })
    }
  } catch (e) {
    $toast.toastErrorHandler(
      t(
        `notification.authclient.${isEdit.value ? 'update' : 'create'}.error`,
        'Failed to save auth client',
      ),
    )(e)
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await $SystemAPI.authClientDelete({ clientID: authClient.value.authClientID })
    $toast.toastSuccess(t('notification.authclient.delete.success'))
    router.push({ name: 'system.authClients' })
  } catch (e) {
    $toast.toastErrorHandler(t('notification.authclient.delete.error'))(e)
  } finally {
    deleting.value = false
  }
}

const { markSaved } = useUnsavedGuard({
  isDirty: () => !saving.value && !deleting.value && !!authClient.value && !!initialAuthClient.value && !isEqual(authClient.value, initialAuthClient.value),
  messageKey: 'general.editor.unsavedChanges',
})

onMounted(() => loadAuthClient())
watch(
  () => route.params.authClientID,
  () => loadAuthClient(),
)

// Auto-default impersonateUser when switching to client_credentials
watch(
  () => authClient.value?.validGrant,
  (grant) => {
    if (!authClient.value) return
    if (
      grant === 'client_credentials' &&
      (!authClient.value.security.impersonateUser ||
        authClient.value.security.impersonateUser === NoID)
    ) {
      authClient.value.security.impersonateUser = $Auth?.user?.userID || NoID
    }
  },
)
</script>
