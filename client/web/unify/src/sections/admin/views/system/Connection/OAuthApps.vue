<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.oauthApps.title') }}</span>
  </Teleport>

  <CViewContainer>
    <div class="h-full overflow-y-auto">
      <div class="flex flex-col gap-4 pb-2">
        <Card :pt="{ content: { class: 'p-0' } }">
          <template #content>
            <div class="flex flex-col gap-4">
              <p class="text-sm text-muted-color">{{ $t('system.oauthApps.intro') }}</p>

              <div class="flex flex-col gap-1.5">
                <label class="block font-medium text-sm">
                  {{ $t('system.oauthApps.redirectUrl') }}
                </label>
                <span class="block text-xs text-muted-color">
                  {{ $t('system.oauthApps.redirectUrlHint') }}
                </span>
                <div class="flex items-center gap-2 mt-1">
                  <InputText :model-value="redirectURL" readonly class="w-full font-mono text-sm" />
                  <Button
                    icon="pi pi-copy"
                    severity="secondary"
                    outlined
                    size="small"
                    v-tooltip.bottom="$t('general.label.copy')"
                    @click="copyRedirectURL"
                  />
                </div>
              </div>
            </div>
          </template>
        </Card>

        <div class="flex items-center justify-between">
          <span class="font-medium">{{ $t('system.oauthApps.providers') }}</span>
          <Button
            :label="$t('system.oauthApps.add')"
            icon="pi pi-plus"
            size="small"
            @click="openEditor(null)"
          />
        </div>

        <div v-if="loading" class="flex justify-center p-6">
          <ProgressSpinner style="width: 2rem; height: 2rem" />
        </div>

        <div v-else-if="!apps.length" class="text-sm text-muted-color italic p-4 text-center">
          {{ $t('system.oauthApps.empty') }}
        </div>

        <div v-else class="flex flex-col gap-2">
          <div
            v-for="app in apps"
            :key="app.handle"
            class="flex items-center gap-3 rounded-md border border-surface-200 dark:border-surface-700 p-3"
          >
            <span
              class="flex items-center justify-center w-9 h-9 rounded-md bg-surface-100 dark:bg-surface-800 text-muted-color shrink-0"
            >
              <i class="pi pi-key" />
            </span>
            <div class="flex flex-col min-w-0 flex-1">
              <span class="font-medium capitalize">{{ app.handle }}</span>
              <span class="text-xs text-muted-color truncate">
                {{ $t('system.oauthApps.credentialsSet') }}
              </span>
            </div>
            <Button
              icon="pi pi-pencil"
              text
              severity="secondary"
              size="small"
              @click="openEditor(app)"
            />
            <CInputDelete
              :message="$t('system.oauthApps.deleteConfirm')"
              :header="app.handle"
              size="small"
              @confirm="handleDelete(app)"
            />
          </div>
        </div>
      </div>
    </div>
  </CViewContainer>

  <Dialog
    v-model:visible="modal"
    modal
    :header="
      editing.isNew
        ? $t('system.oauthApps.editor.createTitle')
        : $t('system.oauthApps.editor.editTitle')
    "
    :style="{ width: '42rem' }"
    :breakpoints="{ '575px': '90vw' }"
  >
    <div class="flex flex-col gap-4">
      <div
        class="flex flex-col gap-3 rounded-md border border-surface-200 dark:border-surface-700 bg-surface-50 dark:bg-surface-800 p-3"
      >
        <div class="flex items-center justify-between gap-2">
          <span class="font-medium text-sm">{{ $t('system.oauthApps.editor.setupTitle') }}</span>
          <a
            v-if="selectedProvider?.docsURL"
            :href="selectedProvider.docsURL"
            target="_blank"
            rel="noopener"
            class="text-xs text-primary hover:underline whitespace-nowrap"
          >
            {{ $t('system.oauthApps.editor.docsLink') }}
            <i class="pi pi-external-link text-[10px]" />
          </a>
        </div>

        <Button
          v-if="selectedProvider?.consoleURL"
          :label="$t('system.oauthApps.editor.openConsole', { provider: editing.handle })"
          icon="pi pi-external-link"
          severity="secondary"
          outlined
          size="small"
          class="self-start capitalize"
          @click="openURL(selectedProvider.consoleURL)"
        />

        <ol
          v-if="selectedProvider?.setupSteps?.length"
          class="list-decimal ml-4 flex flex-col gap-1.5 text-sm text-muted-color"
        >
          <li v-for="(step, i) in selectedProvider.setupSteps" :key="i">{{ step }}</li>
        </ol>
        <p v-else class="text-sm text-muted-color">
          {{ $t('system.oauthApps.editor.genericStep') }}
        </p>

        <div class="flex flex-col gap-1">
          <label class="block text-xs font-medium">
            {{ $t('system.oauthApps.editor.redirectLabel') }}
          </label>
          <div class="flex items-center gap-2">
            <InputText :model-value="redirectURL" readonly class="w-full font-mono text-xs" />
            <Button
              icon="pi pi-copy"
              severity="secondary"
              outlined
              size="small"
              v-tooltip.bottom="$t('general.label.copy')"
              @click="copyRedirectURL"
            />
          </div>
        </div>
      </div>

      <CFormGroup :label="$t('system.oauthApps.editor.handle')" required>
        <Select
          v-if="editing.isNew && availableProviders.length"
          v-model="editing.handle"
          :options="availableProviders"
          option-label="handle"
          option-value="handle"
          :placeholder="$t('system.oauthApps.editor.selectProvider')"
          class="w-full"
        />
        <InputText
          v-else
          v-model="editing.handle"
          :disabled="!editing.isNew"
          placeholder="google"
          class="w-full"
        />
      </CFormGroup>

      <CFormGroup :label="$t('system.oauthApps.editor.clientId')" :required="editing.isNew">
        <InputText
          v-model="editing.clientID"
          :placeholder="editing.isNew ? '' : $t('system.oauthApps.editor.keepPlaceholder')"
          class="w-full"
        />
      </CFormGroup>

      <CFormGroup :label="$t('system.oauthApps.editor.clientSecret')" :required="editing.isNew">
        <Password
          v-model="editing.clientSecret"
          :placeholder="editing.isNew ? '' : $t('system.oauthApps.editor.keepPlaceholder')"
          toggle-mask
          :feedback="false"
          input-class="w-full"
          class="w-full"
        />
      </CFormGroup>
    </div>

    <template #footer>
      <Button
        :label="$t('general.label.cancel')"
        severity="secondary"
        text
        size="small"
        @click="modal = false"
      />
      <Button
        :label="$t('general.label.save')"
        size="small"
        :loading="saving"
        @click="handleSave"
      />
    </template>
  </Dialog>
</template>

<script setup>
import { computed, inject, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { components } from '@planetcrust/human-vue'

const { CInputDelete, CViewContainer } = components

const { t } = useI18n()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')
const $Auth = inject('$Auth')

const REDIRECT_PATH = '/oauth2/connection/callback'
const APP_KEYS = [
  'client-id',
  'client-secret',
  'auth-url',
  'token-url',
  'auth-params',
  'pkce',
  'identity-url',
  'identity-email-field',
]

const loading = ref(false)
const saving = ref(false)
const modal = ref(false)
const apps = ref([])
const redirectURL = ref('')
const catalogProviders = ref([])

const editing = reactive({
  isNew: true,
  handle: '',
  clientID: '',
  clientSecret: '',
})

// All providers the catalog defines — the create picker. Picking one that is
// already configured just re-enters its credentials.
const availableProviders = computed(() => catalogProviders.value.filter(p => p.handle))

// The blueprint for the provider being edited — drives the setup guidance.
const selectedProvider = computed(() =>
  catalogProviders.value.find(p => p.handle === editing.handle),
)

function openURL(url) {
  if (url) window.open(url, '_blank', 'noopener')
}

function appPrefix(handle) {
  return `connection.oauth.apps.${handle}.`
}

// Catalog provider blueprints, via the raw client (typed method not generated).
async function loadProviders() {
  try {
    const res = await $SystemAPI.api().get('/connections/oauth-apps')
    const payload = res?.data?.response ?? res?.data
    catalogProviders.value = Array.isArray(payload) ? payload : []
  } catch {
    // Catalog unavailable — the form falls back to a free-text provider name.
    catalogProviders.value = []
  }
}

async function load() {
  loading.value = true
  try {
    const current = await $SystemAPI.settingsCurrent()
    const oauth = current?.connection?.oauth || {}
    // Endpoints come from the provider catalog now; settings hold only the
    // credentials. The backend drops a provider whose keys were all deleted.
    apps.value = (oauth.apps || []).filter(a => a && a.handle)
    const authBase = ($Auth?.authURL || `${window.location.origin}/auth`).replace(/\/$/, '')
    redirectURL.value = oauth.redirectURL || `${authBase}${REDIRECT_PATH}`
  } catch (e) {
    $toast.toastErrorHandler(t('system.oauthApps.loadError'))(e)
  } finally {
    loading.value = false
  }
}

function openEditor(app) {
  Object.assign(editing, {
    isNew: !app,
    handle: app?.handle || '',
    clientID: '',
    clientSecret: '',
  })
  modal.value = true
}

function validate() {
  if (!editing.handle.trim()) return t('system.oauthApps.editor.handleRequired')
  if (editing.isNew && (!editing.clientID.trim() || !editing.clientSecret.trim())) {
    return t('system.oauthApps.editor.credentialsRequired')
  }
  return ''
}

async function handleSave() {
  const error = validate()
  if (error) {
    $toast.toastWarning(error)
    return
  }

  const prefix = appPrefix(editing.handle.trim())
  // Only credentials + the shared redirect live in settings; endpoints come
  // from the provider catalog. Secrets are write-only — sent only when typed.
  const values = [{ name: 'connection.oauth.redirect-url', value: redirectURL.value }]
  if (editing.clientID.trim()) {
    values.push({ name: `${prefix}client-id`, value: editing.clientID.trim() })
  }
  if (editing.clientSecret.trim()) {
    values.push({ name: `${prefix}client-secret`, value: editing.clientSecret.trim() })
  }

  saving.value = true
  try {
    await $SystemAPI.settingsUpdate({ values })
    $toast.toastSuccess(t('system.oauthApps.saveSuccess'))
    modal.value = false
    await load()
  } catch (e) {
    $toast.toastErrorHandler(t('system.oauthApps.saveError'))(e)
  } finally {
    saving.value = false
  }
}

async function handleDelete(app) {
  const prefix = appPrefix(app.handle)
  const values = APP_KEYS.map(k => ({ name: `${prefix}${k}`, value: null }))
  try {
    await $SystemAPI.settingsUpdate({ values })
    $toast.toastSuccess(t('system.oauthApps.deleteSuccess'))
    await load()
  } catch (e) {
    $toast.toastErrorHandler(t('system.oauthApps.deleteError'))(e)
  }
}

async function copyRedirectURL() {
  try {
    await navigator.clipboard.writeText(redirectURL.value)
    $toast.toastSuccess(t('general.label.copied'))
  } catch {
    // Clipboard may be blocked; the field is selectable as a fallback.
  }
}

onMounted(() => {
  load()
  loadProviders()
})
</script>
