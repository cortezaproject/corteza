<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.settings.editor.title', 'System settings') }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <div v-else class="flex flex-col h-full">
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <!-- General Authentication -->
      <Panel
        v-if="groups['auth'].length"
        :header="$t('system.settings.editor.auth.title', 'Authentication')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <div class="flex flex-col divide-y">
          <div
            v-for="setting in groups['auth']"
            :key="setting.name"
            class="flex items-center justify-between py-3 gap-4"
          >
            <div class="flex flex-col flex-1 min-w-0">
              <span class="font-medium text-sm">{{ setting.label }}</span>
              <span class="text-xs text-surface-500 font-mono">{{ setting.name }}</span>
            </div>
            <div class="flex-shrink-0 w-64">
              <ToggleSwitch
                v-if="typeof setting.value === 'boolean'"
                v-model="setting.value"
                @update:modelValue="markDirty(setting)"
              />
              <InputText
                v-else
                v-model="setting.value"
                size="small"
                class="w-full"
                @update:modelValue="markDirty(setting)"
              />
            </div>
          </div>
        </div>
      </Panel>

      <!-- Internal Authentication -->
      <Panel
        v-if="groups['auth.internal'].length"
        :header="$t('system.settings.editor.auth.internal.title', 'Internal')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <div class="flex flex-col divide-y">
          <div
            v-for="setting in groups['auth.internal']"
            :key="setting.name"
            class="flex items-center justify-between py-3 gap-4"
          >
            <div class="flex flex-col flex-1 min-w-0">
              <span class="font-medium text-sm">{{ setting.label }}</span>
              <span class="text-xs text-surface-500 font-mono">{{ setting.name }}</span>
            </div>
            <div class="flex-shrink-0 w-64">
              <ToggleSwitch
                v-if="typeof setting.value === 'boolean'"
                v-model="setting.value"
                @update:modelValue="markDirty(setting)"
              />
              <InputNumber
                v-else-if="typeof setting.value === 'number'"
                v-model="setting.value"
                size="small"
                class="w-full"
                @update:modelValue="markDirty(setting)"
              />
              <InputText
                v-else
                v-model="setting.value"
                size="small"
                class="w-full"
                @update:modelValue="markDirty(setting)"
              />
            </div>
          </div>
        </div>
      </Panel>

      <!-- Authentication email sender mail -->
      <Panel
        v-if="groups['auth.mail'].length"
        :header="$t('system.settings.editor.auth.mail.title', 'Authentication email sender mail')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <div class="flex flex-col divide-y">
          <div
            v-for="setting in groups['auth.mail']"
            :key="setting.name"
            class="flex items-center justify-between py-3 gap-4"
          >
            <div class="flex flex-col flex-1 min-w-0">
              <span class="font-medium text-sm">{{ setting.label }}</span>
              <span class="text-xs text-surface-500 font-mono">{{ setting.name }}</span>
            </div>
            <div class="flex-shrink-0 w-64">
              <InputText
                v-model="setting.value"
                size="small"
                class="w-full"
                @update:modelValue="markDirty(setting)"
              />
            </div>
          </div>
        </div>
      </Panel>

      <!-- Multi-factor authentication -->
      <Panel
        v-if="groups['auth.mfa'].length"
        :header="$t('system.settings.editor.auth.mfa.title', 'Multi-factor authentication')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <div class="flex flex-col divide-y">
          <div
            v-for="setting in groups['auth.mfa']"
            :key="setting.name"
            class="flex items-center justify-between py-3 gap-4"
          >
            <div class="flex flex-col flex-1 min-w-0">
              <span class="font-medium text-sm">{{ setting.label }}</span>
              <span class="text-xs text-surface-500 font-mono">{{ setting.name }}</span>
            </div>
            <div class="flex-shrink-0 w-64">
              <ToggleSwitch
                v-if="typeof setting.value === 'boolean'"
                v-model="setting.value"
                @update:modelValue="markDirty(setting)"
              />
              <InputText
                v-else
                v-model="setting.value"
                size="small"
                class="w-full"
                @update:modelValue="markDirty(setting)"
              />
            </div>
          </div>
        </div>
      </Panel>

      <!-- Auto logout -->
      <Panel
        v-if="groups['auth.auto-logout'].length"
        :header="$t('system.settings.editor.auth.auto-logout.title', 'Auto logout')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <div class="flex flex-col divide-y">
          <div
            v-for="setting in groups['auth.auto-logout']"
            :key="setting.name"
            class="flex items-center justify-between py-3 gap-4"
          >
            <div class="flex flex-col flex-1 min-w-0">
              <span class="font-medium text-sm">{{ setting.label }}</span>
              <span class="text-xs text-surface-500 font-mono">{{ setting.name }}</span>
            </div>
            <div class="flex-shrink-0 w-64">
              <ToggleSwitch
                v-if="typeof setting.value === 'boolean'"
                v-model="setting.value"
                @update:modelValue="markDirty(setting)"
              />
              <InputNumber
                v-else-if="typeof setting.value === 'number'"
                v-model="setting.value"
                size="small"
                class="w-full"
                @update:modelValue="markDirty(setting)"
              />
              <InputText
                v-else
                v-model="setting.value"
                size="small"
                class="w-full"
                @update:modelValue="markDirty(setting)"
              />
            </div>
          </div>
        </div>
      </Panel>

      <!-- External Authentication Providers -->
      <Panel
        v-if="groups['auth.external'].length"
        :header="$t('system.settings.editor.external.title', 'External Authentication Providers')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <div class="flex flex-col divide-y">
          <div
            v-for="setting in groups['auth.external']"
            :key="setting.name"
            class="flex items-center justify-between py-3 gap-4"
          >
            <div class="flex flex-col flex-1 min-w-0">
              <span class="font-medium text-sm">{{ setting.label }}</span>
              <span class="text-xs text-surface-500 font-mono">{{ setting.name }}</span>
            </div>
            <div class="flex-shrink-0 w-64">
              <ToggleSwitch
                v-if="typeof setting.value === 'boolean'"
                v-model="setting.value"
                @update:modelValue="markDirty(setting)"
              />
              <InputText
                v-else
                v-model="setting.value"
                size="small"
                class="w-full"
                @update:modelValue="markDirty(setting)"
              />
            </div>
          </div>
        </div>
      </Panel>
    </div>

    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="p-3 flex items-center justify-end">
        <Button
          :label="$t('general.label.save')"
          icon="pi pi-save"
          :loading="saving"
          :disabled="dirtySettings.size === 0"
          @click="handleSave"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(false)
const saving = ref(false)
const rawSettings = ref([])
const dirtySettings = ref(new Set())

// Human-readable labels for known auth setting keys
const settingLabels = computed(() => ({
  'auth.url': t('system.settings.editor.auth.url', 'URL'),
  'auth.internal.enabled': t(
    'system.settings.editor.auth.internal.enabled',
    'Internal authentication enabled',
  ),
  'auth.internal.signup.enabled': t(
    'system.settings.editor.auth.internal.signup.enabled',
    'Signup enabled',
  ),
  'auth.internal.signup.email-confirmation-required': t(
    'system.settings.editor.auth.internal.signup.email-confirmation-required',
    'Email confirmation required',
  ),
  'auth.internal.signup.split-credentials-check': t(
    'system.settings.editor.auth.internal.signup.split-credentials-check.label',
    'Enable split-credentials check',
  ),
  'auth.internal.password-reset.enabled': t(
    'system.settings.editor.auth.internal.password-reset.enabled',
    'Password reset enabled',
  ),
  'auth.internal.password-constraints.min-length': t(
    'system.settings.editor.auth.internal.password-constraints.min-length',
    'Minimum length',
  ),
  'auth.internal.password-constraints.min-num-count': t(
    'system.settings.editor.auth.internal.password-constraints.min-num-count',
    'Minimum number of digits',
  ),
  'auth.internal.password-constraints.min-special-count': t(
    'system.settings.editor.auth.internal.password-constraints.min-special-count',
    'Minimum number of special characters',
  ),
  'auth.internal.password-constraints.min-upper-case-length': t(
    'system.settings.editor.auth.internal.password-constraints.min-upper-case-length',
    'Minimum number of upper case characters',
  ),
  'auth.internal.password-constraints.min-lower-case-length': t(
    'system.settings.editor.auth.internal.password-constraints.min-lower-case-length',
    'Minimum number of lower case characters',
  ),
  'auth.internal.profile-avatar.enabled': t(
    'system.settings.editor.auth.internal.profile-avatar.enabled',
    'Profile avatar enabled',
  ),
  'auth.internal.send-user-invite-email.enabled': t(
    'system.settings.editor.auth.internal.send-user-invite-email.enabled',
    'Send invite email on user creation',
  ),
  'auth.internal.send-user-invite-email.expires': t(
    'system.settings.editor.auth.internal.send-user-invite-email.expires.label',
    'Valid for',
  ),
  'auth.mail.from-address': t('system.settings.editor.auth.mail.from-address', "Sender's address"),
  'auth.mail.from-name': t('system.settings.editor.auth.mail.from-name', "Sender's name"),
  'auth.mfa.totp.enabled': t('system.settings.editor.auth.mfa.TOTP.enabled', 'Enable TOTP'),
  'auth.mfa.totp.enforced': t(
    'system.settings.editor.auth.mfa.TOTP.enforced',
    'Require all users to use TOTP',
  ),
  'auth.mfa.totp.issuer': t('system.settings.editor.auth.mfa.TOTP.issuer.label', 'Issuer'),
  'auth.mfa.email-otp.enabled': t(
    'system.settings.editor.auth.mfa.emailOTP.enabled',
    'Enable email OTP',
  ),
  'auth.mfa.email-otp.enforced': t(
    'system.settings.editor.auth.mfa.emailOTP.enforced',
    'Require all users to use email OTP',
  ),
  'auth.mfa.email-otp.expires': t(
    'system.settings.editor.auth.mfa.emailOTP.expires.label',
    'Valid for',
  ),
  'auth.auto-logout.enabled': t('system.settings.editor.auth.auto-logout.enabled.label', 'Enabled'),
  'auth.auto-logout.timeout': t('system.settings.editor.auth.auto-logout.timeout.label', 'Timeout'),
  'auth.session-lifetime': t(
    'system.settings.editor.label.auth-session-lifetime',
    'Session lifetime',
  ),
  'auth.external.enabled': t(
    'system.settings.editor.external.enabled',
    'Enable external authentication',
  ),
}))

// Group prefixes in display order (most specific first for matching)
const groupPrefixes = [
  'auth.internal',
  'auth.mail',
  'auth.mfa',
  'auth.auto-logout',
  'auth.external',
  'auth',
]

// Build grouped settings
const groups = computed(() => {
  const result = {}
  for (const prefix of groupPrefixes) {
    result[prefix] = []
  }

  for (const setting of rawSettings.value) {
    const prefix = getGroupPrefix(setting.name)
    if (prefix && result[prefix]) {
      result[prefix].push(setting)
    }
  }

  return result
})

function getGroupPrefix(name) {
  // Try longest matching prefix first
  const sorted = [...groupPrefixes].sort((a, b) => b.length - a.length)
  for (const prefix of sorted) {
    if (name.startsWith(prefix + '.') || name === prefix) {
      return prefix
    }
  }
  return 'auth'
}

function markDirty(setting) {
  dirtySettings.value = new Set([...dirtySettings.value, setting.name])
}

async function loadSettings() {
  loading.value = true
  try {
    const result = await $SystemAPI.settingsList({ prefix: 'auth.' })
    rawSettings.value = (result || []).map(s => ({
      name: s.name,
      label: settingLabels.value[s.name] || s.name.split('.').pop(),
      value: parseValue(s.value),
      originalValue: s.value,
    }))
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.fetch.error', 'Failed to load settings'))(e)
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

function serializeValue(value) {
  if (typeof value === 'boolean') return value
  if (typeof value === 'number') return value
  return value
}

async function handleSave() {
  saving.value = true
  try {
    const values = rawSettings.value
      .filter(s => dirtySettings.value.has(s.name))
      .map(s => ({
        name: s.name,
        value: serializeValue(s.value),
      }))

    await $SystemAPI.settingsUpdate({ values })
    dirtySettings.value = new Set()
    $toast.toastSuccess(t('notification.settings.update.success', 'Settings saved'))
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.update.error', 'Failed to save settings'))(e)
  } finally {
    saving.value = false
  }
}

onMounted(() => loadSettings())
</script>
