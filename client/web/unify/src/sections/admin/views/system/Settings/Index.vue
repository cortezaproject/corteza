<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.settings.editor.title') }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <div v-else class="flex flex-col h-full">
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <!-- Internal Authentication -->
      <Panel
        :header="$t('system.settings.editor.auth.internal.title')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <CInputSwitch
            v-model="settings['auth.internal.enabled']"
            :label="$t('system.settings.editor.auth.internal.enabled')"
          />

          <CInputSwitch
            v-model="settings['auth.internal.password-reset.enabled']"
            :label="$t('system.settings.editor.auth.internal.password-reset.enabled')"
          />

          <CInputSwitch
            v-model="settings['auth.internal.signup.email-confirmation-required']"
            :label="$t('system.settings.editor.auth.internal.signup.email-confirmation-required')"
          />

          <CInputSwitch
            v-model="settings['auth.internal.signup.enabled']"
            :label="$t('system.settings.editor.auth.internal.signup.enabled')"
          />

          <CInputSwitch
            v-model="settings['auth.internal.profile-avatar.enabled']"
            :label="$t('system.settings.editor.auth.internal.profile-avatar.enabled')"
          />

          <CInputSwitch
            v-model="settings['auth.internal.split-credentials-check']"
            :label="$t('system.settings.editor.auth.internal.signup.split-credentials-check.label')"
            :description="
              $t('system.settings.editor.auth.internal.signup.split-credentials-check.description')
            "
          />
        </div>
      </Panel>

      <!-- Password Constraints -->
      <Panel
        :header="$t('system.settings.editor.auth.internal.password-constraints.title')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <Message v-if="!passwordSecurityEnabled" severity="warn" :closable="false" class="mb-4">
          {{ $t('system.settings.editor.auth.internal.password-constraints.ignored-security') }}
        </Message>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{ $t('system.settings.editor.auth.internal.password-constraints.min-length') }}
            </label>
            <span class="text-xs text-muted-color">
              {{
                $t(
                  'system.settings.editor.auth.internal.password-constraints.min-length-description',
                )
              }}
            </span>
            <InputNumber
              v-model="settings['auth.internal.password-constraints.min-length']"
              :min="8"
              placeholder="8"
              class="w-full"
            />
          </div>

          <div class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{ $t('system.settings.editor.auth.internal.password-constraints.min-num-count') }}
            </label>
            <span class="text-xs text-muted-color">
              {{
                $t(
                  'system.settings.editor.auth.internal.password-constraints.min-num-count-description',
                )
              }}
            </span>
            <InputNumber
              v-model="settings['auth.internal.password-constraints.min-num-count']"
              :min="0"
              placeholder="0"
              class="w-full"
            />
          </div>

          <div class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{
                $t(
                  'system.settings.editor.auth.internal.password-constraints.min-upper-case-length',
                )
              }}
            </label>
            <span class="text-xs text-muted-color">
              {{
                $t(
                  'system.settings.editor.auth.internal.password-constraints.min-upper-case-description',
                )
              }}
            </span>
            <InputNumber
              v-model="settings['auth.internal.password-constraints.min-upper-case']"
              :min="0"
              placeholder="0"
              class="w-full"
            />
          </div>

          <div class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{
                $t(
                  'system.settings.editor.auth.internal.password-constraints.min-lower-case-length',
                )
              }}
            </label>
            <span class="text-xs text-muted-color">
              {{
                $t(
                  'system.settings.editor.auth.internal.password-constraints.min-lower-case-description',
                )
              }}
            </span>
            <InputNumber
              v-model="settings['auth.internal.password-constraints.min-lower-case']"
              :min="0"
              placeholder="0"
              class="w-full"
            />
          </div>

          <div class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{
                $t('system.settings.editor.auth.internal.password-constraints.min-special-count')
              }}
            </label>
            <span class="text-xs text-muted-color">
              {{
                $t(
                  'system.settings.editor.auth.internal.password-constraints.min-special-count-description',
                )
              }}
            </span>
            <InputNumber
              v-model="settings['auth.internal.password-constraints.min-special-count']"
              :min="0"
              placeholder="0"
              class="w-full"
            />
          </div>
        </div>
      </Panel>

      <!-- Multi-factor authentication -->
      <Panel
        :header="$t('system.settings.editor.auth.mfa.title')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
          <CInputSwitch
            v-model="settings['auth.multi-factor.email-otp.enabled']"
            :label="$t('system.settings.editor.auth.mfa.emailOTP.enabled')"
            @update:modelValue="onEmailOtpToggle"
          />

          <div class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{ $t('system.settings.editor.auth.mfa.emailOTP.expires.label') }}
            </label>
            <span class="text-xs text-muted-color">
              {{ $t('system.settings.editor.auth.mfa.emailOTP.expires.description') }}
            </span>
            <InputGroup>
              <InputNumber
                v-model="settings['auth.multi-factor.email-otp.expires']"
                placeholder="60"
              />
              <InputGroupAddon>{{ $t('general.label.seconds') }}</InputGroupAddon>
            </InputGroup>
          </div>

          <CInputSwitch
            v-if="settings['auth.multi-factor.email-otp.enabled']"
            v-model="settings['auth.multi-factor.email-otp.enforced']"
            :label="$t('system.settings.editor.auth.mfa.emailOTP.enforced')"
          />
        </div>

        <Divider />

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <CInputSwitch
            v-model="settings['auth.multi-factor.totp.enabled']"
            :label="$t('system.settings.editor.auth.mfa.TOTP.enabled')"
            @update:modelValue="onTotpToggle"
          />

          <div class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{ $t('system.settings.editor.auth.mfa.TOTP.issuer.label') }}
            </label>
            <span class="text-xs text-muted-color">
              {{ $t('system.settings.editor.auth.mfa.TOTP.issuer.description') }}
            </span>
            <InputText
              v-model="settings['auth.multi-factor.totp.issuer']"
              placeholder="Human"
              class="w-full"
            />
          </div>

          <CInputSwitch
            v-if="settings['auth.multi-factor.totp.enabled']"
            v-model="settings['auth.multi-factor.totp.enforced']"
            :label="$t('system.settings.editor.auth.mfa.TOTP.enforced')"
          />
        </div>
      </Panel>

      <!-- Authentication email sender mail -->
      <Panel
        :header="$t('system.settings.editor.auth.mail.title')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{ $t('system.settings.editor.auth.mail.from-address') }}
            </label>
            <InputText v-model="settings['auth.mail.from-address']" type="email" class="w-full" />
          </div>

          <div class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{ $t('system.settings.editor.auth.mail.from-name') }}
            </label>
            <InputText v-model="settings['auth.mail.from-name']" class="w-full" />
          </div>
        </div>
      </Panel>

      <!-- Invite email -->
      <Panel
        :header="$t('system.settings.editor.auth.internal.send-user-invite-email.title')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <CInputSwitch
            v-model="settings['auth.internal.send-user-invite-email.enabled']"
            :label="$t('system.settings.editor.auth.internal.send-user-invite-email.enabled')"
            :description="
              $t('system.settings.editor.auth.internal.send-user-invite-email.description')
            "
          />

          <div class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{ $t('system.settings.editor.auth.internal.send-user-invite-email.expires.label') }}
            </label>
            <span class="text-xs text-muted-color">
              {{
                $t(
                  'system.settings.editor.auth.internal.send-user-invite-email.expires.description',
                )
              }}
            </span>
            <InputGroup>
              <InputNumber v-model="settings['auth.internal.send-user-invite-email.expires']" />
              <InputGroupAddon>{{ $t('general.label.hours') }}</InputGroupAddon>
            </InputGroup>
          </div>
        </div>
      </Panel>

      <!-- Auto logout -->
      <Panel
        :header="$t('system.settings.editor.auth.auto-logout.title')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <CInputSwitch
            v-model="settings['auth.auto-logout.enabled']"
            :label="$t('system.settings.editor.auth.auto-logout.enabled.label')"
            :description="$t('system.settings.editor.auth.auto-logout.enabled.description')"
          />

          <div class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{ $t('system.settings.editor.auth.auto-logout.timeout.label') }}
            </label>
            <span class="text-xs text-muted-color">
              {{ $t('system.settings.editor.auth.auto-logout.timeout.description') }}
            </span>
            <InputGroup>
              <InputNumber v-model="settings['auth.auto-logout.timeout']" />
              <InputGroupAddon>{{ $t('general.label.seconds') }}</InputGroupAddon>
            </InputGroup>
          </div>
        </div>
      </Panel>

      <!-- External Authentication Providers -->
      <Panel
        :header="$t('system.settings.editor.external.title')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <div class="flex flex-col gap-4">
          <div class="flex items-center justify-between flex-wrap gap-2">
            <Button
              :label="$t('system.settings.editor.external.oidc.add')"
              icon="pi pi-plus"
              severity="primary"
              @click="newOIDC"
            />
            <CInputSwitch
              v-model="external.enabled"
              :label="$t('system.settings.editor.external.enabled')"
            />
          </div>

          <CResourceTable
            :items="providerItems"
            :fields="providerFields"
            :action-items="getProviderActions"
            primary-key="tag"
          >
            <template #body-enabled="{ data }">
              <Checkbox
                :model-value="data.enabled"
                :binary="true"
                @update:model-value="data.enable($event)"
              />
            </template>

            <template #body-provider="{ data }">
              <span :class="{ 'line-through opacity-40': data.deleted }">
                {{ data.provider || data.tag }}
              </span>
            </template>

            <template #body-info="{ data }">
              <span :class="{ 'line-through opacity-40': data.deleted }">
                {{ data.info }}
              </span>
            </template>
          </CResourceTable>
        </div>
      </Panel>

      <!-- Auth background -->
      <Panel
        :header="$t('system.settings.editor.bgScreen.title')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <div class="flex flex-col gap-6">
          <div class="flex flex-col gap-2">
            <label class="font-medium text-sm text-primary">
              {{ $t('system.settings.editor.bgScreen.image.uploader.label') }}
            </label>
            <CFileDropZone
              accept="image/*"
              :uploading="bgUploading"
              :error="bgUploadError"
              :preview-url="bgImageUrl"
              :clearable="bgImageIsCustom"
              compact
              preview-max-width="100%"
              preview-max-height="200px"
              @select="onBgImageSelect"
              @clear="onBgImageClear"
            />
          </div>

          <div class="flex flex-col gap-2">
            <label class="font-medium text-sm text-primary">
              {{ $t('system.settings.editor.bgScreen.image.editor.label') }}
            </label>
            <Textarea
              v-model="settings['auth.ui.styles']"
              rows="16"
              class="w-full font-mono text-sm"
            />
          </div>
        </div>
      </Panel>

      <!-- Provider editor dialog -->
      <Dialog
        v-model:visible="modal.open"
        :header="modal.title"
        modal
        :style="{ width: '50rem' }"
        :breakpoints="{ '768px': '90vw' }"
      >
        <component
          :is="modal.component"
          v-if="modal.component && modal.data"
          v-model="modal.data"
        />
        <template #footer>
          <div class="flex justify-end gap-2">
            <Button
              :label="$t('general.label.cancel')"
              severity="secondary"
              text
              size="small"
              @click="modal.open = false"
            />
            <Button
              :label="$t('general.label.save')"
              size="small"
              @click="applyModal"
            />
          </div>
        </template>
      </Dialog>
    </div>

    <CEditorActions>
      <Button
        :label="$t('general.label.save')"
        icon="pi pi-save"
        :loading="saving"
        @click="handleSave"
      />
    </CEditorActions>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, reactive, ref, shallowRef } from 'vue'
import { useI18n } from 'vue-i18n'
import { isEqual } from 'lodash-es'
import { components, useFileUpload } from '@planetcrust/human-vue'
import ExternalStd from './auth/ExternalStd.vue'
import ExternalOIDC from './auth/ExternalOIDC.vue'
import ExternalSAML from './auth/ExternalSAML.vue'

const { CResourceTable, CFileDropZone } = components

const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')
const $Settings = inject('$Settings')

const loading = ref(false)
const saving = ref(false)
const settings = reactive({})
const passwordSecurityEnabled = ref(true)

// Auth background image upload state
const {
  uploading: bgUploading,
  uploadError: bgUploadError,
  uploadFileRaw: uploadBgRaw,
} = useFileUpload()

const bgImageUrl = ref('')
const bgImageIsCustom = ref(false)

function refreshBgImage() {
  const raw = settings['auth.ui.background-image-src']
  bgImageIsCustom.value = typeof raw === 'string' && raw.startsWith('attachment:')
  if (bgImageIsCustom.value) {
    const m = /^attachment:(\d+)/.exec(raw)
    if (m) {
      bgImageUrl.value =
        $SystemAPI.baseURL +
        $SystemAPI.attachmentOriginalEndpoint({
          attachmentID: m[1],
          kind: 'settings',
          name: 'auth.ui.background-image-src',
        })
      return
    }
  }
  bgImageUrl.value = ''
}

async function onBgImageSelect(files) {
  const file = files[0]
  if (!file) return
  try {
    const endpoint =
      $SystemAPI.baseURL + $SystemAPI.settingsSetEndpoint({ key: 'auth.ui.background-image-src' })
    const token = $SystemAPI.accessTokenFn ? $SystemAPI.accessTokenFn() : ''
    await uploadBgRaw(file, { url: endpoint, token })
    if ($Settings?.fetch) await $Settings.fetch()
    await loadSettings()
    $toast.toastSuccess(t('notification.settings.update.success'))
  } catch {
    // upload error tracked by composable
  }
}

async function onBgImageClear() {
  try {
    await $SystemAPI.settingsUpdate({
      values: [{ name: 'auth.ui.background-image-src', value: null }],
    })
    if ($Settings?.fetch) await $Settings.fetch()
    await loadSettings()
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.update.error'))(e)
  }
}

function onEmailOtpToggle(v) {
  if (!v) settings['auth.multi-factor.email-otp.enforced'] = false
}

function onTotpToggle(v) {
  if (!v) settings['auth.multi-factor.totp.enforced'] = false
}

// External auth state
const external = reactive({
  enabled: false,
  saml: {},
  oidc: [],
  standard: [],
})

// Original state for dirty tracking
let originalExternal = null

// Modal state
const modal = reactive({
  open: false,
  component: null,
  title: '',
  data: null,
  updater: null,
})

const idpStandard = ['google', 'github', 'facebook', 'linkedin']

const idpSecurity = {
  permittedRoles: [],
  prohibitedRoles: [],
  forcedRoles: [],
}

// Table items
const providerItems = computed(() => {
  const items = []

  // SAML row
  items.push({
    provider: external.saml.name,
    info: external.saml.idp?.url || '',
    tag: 'SAML',
    enabled: external.saml.enabled,
    enable: val => {
      external.saml.enabled = val
    },
    canDelete: false,
    deleted: false,
    editor: {
      component: ExternalSAML,
      data: external.saml,
      title: t('system.settings.editor.external.saml.title'),
      updater: changed => updater('saml', changed),
    },
  })

  // OIDC rows
  external.oidc.forEach((p, i) => {
    items.push({
      provider: p.handle,
      tag: 'OIDC',
      info: p.issuer,
      enabled: p.enabled,
      deleted: p.deleted,
      enable: val => {
        external.oidc[i].enabled = val
      },
      canDelete: true,
      toggleDelete: () => {
        external.oidc[i].deleted = !p.deleted
      },
      editor: {
        component: ExternalOIDC,
        data: p,
        title: p.handle || t('system.settings.editor.external.oidc.title'),
        updater: changed => updater('oidc', changed, i),
      },
    })
  })

  // Standard rows
  external.standard.forEach((p, i) => {
    items.push({
      provider: p.handle,
      info: p.key,
      enabled: p.enabled,
      enable: val => {
        external.standard[i].enabled = val
      },
      canDelete: false,
      deleted: false,
      editor: {
        component: ExternalStd,
        data: p,
        title: p.handle,
        updater: changed => updater('standard', changed, i),
      },
    })
  })

  return items
})

const providerFields = [
  {
    key: 'enabled',
    header: t('system.settings.editor.external.table.header.enabled'),
    headerStyle: 'width: 80px',
  },
  {
    key: 'provider',
    header: t('system.settings.editor.external.table.header.provider'),
    headerStyle: 'width: 200px',
  },
  { key: 'info', header: t('system.settings.editor.external.table.header.info') },
]

function getProviderActions(data) {
  const items = []

  items.push({
    label: t('general.label.edit'),
    icon: 'pi pi-pencil',
    command: () => openEditor(data.editor),
  })

  if (data.canDelete) {
    items.push({ separator: true })
    items.push({
      label: data.deleted ? t('general.undelete') : t('general.label.delete'),
      icon: data.deleted ? 'pi pi-undo' : 'pi pi-trash',
      class: data.deleted ? '' : 'text-red-500',
      command: () => data.toggleDelete(),
    })
  }

  return items
}

function openEditor({ component, title, data, updater }) {
  modal.open = true
  modal.component = shallowRef(component)
  modal.title = title
  modal.updater = updater
  // Deep clone to avoid mutating original until confirmed
  modal.data = JSON.parse(JSON.stringify(data))
}

function applyModal() {
  if (modal.updater) {
    modal.updater(modal.data)
  }
  modal.open = false
}

function newOIDC() {
  const data = {
    handle: '',
    enabled: true,
    issuer: '',
    key: '',
    secret: '',
    scope: '',
    fresh: true,
    security: { ...idpSecurity },
  }

  openEditor({
    component: ExternalOIDC,
    title: t('system.settings.editor.external.oidc.add'),
    data,
    updater: changed => updater('oidc', changed, -1),
  })
}

function updater(key, val, i = undefined) {
  if (i === undefined) {
    Object.assign(external[key], val)
  } else if (i < 0) {
    external[key].push(val)
  } else {
    external[key][i] = val
  }
}

// Parse settings into structured external auth data
function prepareExternal(allSettings) {
  const extractKey = (name, type = 'string') => {
    const v = allSettings.find(s => s.name === `auth.external.${name}`)
    switch (type) {
      case 'boolean':
        return !!(v || { value: null }).value
      case 'array':
        return (v || { value: [] }).value || []
      case undefined:
        return v ? v.value : undefined
      default:
        return (v || { value: null }).value || ''
    }
  }

  const extractKeys = (provider, base = {}) => {
    const out = { ...base }
    for (const k in base) {
      out[k] = extractKey(
        `providers.${provider}.${k}`,
        Array.isArray(out[k]) ? 'array' : typeof out[k],
      )
    }
    return out
  }

  const extractSec = prefix => {
    return { ...idpSecurity, ...(extractKey(`${prefix}.security`, undefined) || {}) }
  }

  const data = {
    enabled: !!(allSettings.find(v => v.name === 'auth.external.enabled') || {}).value,

    saml: {
      enabled: extractKey('saml.enabled'),
      cert: extractKey('saml.cert'),
      name: extractKey('saml.name'),
      key: extractKey('saml.key'),
      'sign-method': extractKey('saml.sign-method'),
      'sign-requests': extractKey('saml.sign-requests', 'boolean'),
      binding: extractKey('saml.binding'),
      idp: {
        url: extractKey('saml.idp.url'),
        'ident-name': extractKey('saml.idp.ident-name'),
        'ident-handle': extractKey('saml.idp.ident-handle'),
        'ident-identifier': extractKey('saml.idp.ident-identifier'),
      },
      security: extractSec('saml'),
    },

    oidc: [],
    standard: [],
  }

  // Standard providers
  data.standard = idpStandard.map(handle => ({
    handle,
    ...extractKeys(handle, {
      enabled: false,
      secret: '',
      key: '',
      security: {},
      usage: [],
    }),
    security: extractSec(`providers.${handle}`),
  }))

  // OIDC providers (dynamic)
  const prefix = 'auth.external.providers.openid-connect.'
  const oidcHandles = [
    ...new Set(
      allSettings
        .filter(v => v.name.indexOf(prefix) === 0)
        .map(({ name }) => name.substring(prefix.length).split('.', 2)[0]),
    ),
  ]

  data.oidc = oidcHandles.map(handle => ({
    ...extractKeys('openid-connect.' + handle, {
      enabled: false,
      issuer: '',
      key: '',
      secret: '',
      scope: '',
      security: {},
    }),
    handle,
    security: extractSec('providers.openid-connect.' + handle),
    deleted: false,
  }))

  return data
}

// Convert structured external data back to flat settings for saving
function getExternalChanges() {
  const c = []
  const prefix = 'auth.external.providers'
  const o = originalExternal
  const e = external

  if (!isEqual(o.enabled, e.enabled)) {
    c.push({ name: 'auth.external.enabled', value: e.enabled })
  }

  const mapKeys = (pfx, wc, org, keys) => {
    for (const k of keys) {
      if (!isEqual(wc[k], org[k])) {
        c.push({ name: `${pfx}.${k}`, value: wc[k] })
      }
    }
  }

  // Standard providers
  e.standard.forEach((p, i) => {
    mapKeys(`${prefix}.${p.handle}`, p, o.standard[i], [
      'key',
      'secret',
      'enabled',
      'security',
      'usage',
    ])
  })

  // OIDC providers
  const oidcKeys = ['key', 'secret', 'enabled', 'issuer', 'scope', 'security']
  e.oidc.forEach((p, i) => {
    if (p.deleted) {
      ;[...oidcKeys, 'weight', 'redirect', 'label'].forEach(name =>
        c.push({ name: `${prefix}.openid-connect.${p.handle}.${name}`, value: null }),
      )
    } else {
      mapKeys(`${prefix}.openid-connect.${p.handle}`, p, o.oidc[i] || {}, oidcKeys)
    }
  })

  // SAML
  mapKeys('auth.external.saml', e.saml, o.saml, [
    'enabled',
    'name',
    'key',
    'cert',
    'sign-method',
    'sign-requests',
    'binding',
    'security',
  ])

  mapKeys('auth.external.saml.idp', e.saml.idp, o.saml.idp, [
    'url',
    'ident-name',
    'ident-handle',
    'ident-identifier',
  ])

  return c
}

async function loadSettings() {
  loading.value = true
  try {
    const authResult = await $SystemAPI.settingsList({ prefix: 'auth.' })
    const allSettings = authResult || []

    for (const s of allSettings) {
      settings[s.name] = parseValue(s.value)
    }

    // Check password security from settings
    passwordSecurityEnabled.value = !!settings['auth.internal.passwordConstraints.passwordSecurity']

    refreshBgImage()

    // Parse external auth into structured data
    const parsed = prepareExternal(allSettings)
    Object.assign(external, parsed)
    originalExternal = JSON.parse(JSON.stringify(parsed))
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.fetch.error'))(e)
  } finally {
    loading.value = false
  }
}

function parseValue(raw) {
  if (raw === null || raw === undefined) return ''
  if (typeof raw === 'boolean') return raw
  if (typeof raw === 'number') return raw
  if (typeof raw === 'object') {
    if ('@value' in raw) return raw['@value']
    return JSON.stringify(raw)
  }
  return raw
}

async function handleSave() {
  saving.value = true
  try {
    const values = Object.entries(settings).map(([name, value]) => ({
      name,
      value,
    }))

    // Add external auth changes
    const externalChanges = getExternalChanges()
    values.push(...externalChanges)

    await $SystemAPI.settingsUpdate({ values })
    $toast.toastSuccess(t('notification.settings.update.success'))

    // Reload to refresh original state
    await loadSettings()
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.update.error'))(e)
  } finally {
    saving.value = false
  }
}

onMounted(() => loadSettings())
</script>
