<template>
  <Dialog
    :visible="visible"
    modal
    :header="$t('system.connections.authSetup.title', { provider })"
    :style="{ width: '42rem' }"
    :breakpoints="{ '575px': '90vw' }"
    @update:visible="emit('update:visible', $event)"
  >
    <div class="flex flex-col gap-4">
      <p class="text-sm text-muted-color">
        {{ $t('system.connections.authSetup.intro', { provider }) }}
      </p>

      <div
        v-if="alreadyConfigured"
        class="flex items-center gap-2 rounded-md border border-green-200 dark:border-green-800 bg-green-50 dark:bg-green-950 p-3 text-sm text-green-700 dark:text-green-400"
      >
        <i class="pi pi-check-circle" />
        <span>{{ $t('system.connections.authSetup.alreadySet') }}</span>
      </div>

      <div
        class="flex flex-col gap-3 rounded-md border border-surface-200 dark:border-surface-700 bg-surface-50 dark:bg-surface-800 p-3"
      >
        <div class="flex items-center justify-between gap-2">
          <span class="font-medium text-sm">{{ $t('system.oauthApps.editor.setupTitle') }}</span>
          <a
            v-if="blueprint?.docsURL"
            :href="blueprint.docsURL"
            target="_blank"
            rel="noopener"
            class="text-xs text-primary hover:underline whitespace-nowrap"
          >
            {{ $t('system.oauthApps.editor.docsLink') }}
            <i class="pi pi-external-link text-[10px]" />
          </a>
        </div>

        <Button
          v-if="blueprint?.consoleURL"
          :label="$t('system.oauthApps.editor.openConsole', { provider })"
          icon="pi pi-external-link"
          severity="secondary"
          outlined
          size="small"
          class="self-start capitalize"
          @click="openURL(blueprint.consoleURL)"
        />

        <ol
          v-if="blueprint?.setupSteps?.length"
          class="list-decimal ml-4 flex flex-col gap-1.5 text-sm text-muted-color"
        >
          <li v-for="(step, i) in blueprint.setupSteps" :key="i">{{ step }}</li>
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

      <CFormGroup :label="$t('system.oauthApps.editor.clientId')" :required="!alreadyConfigured">
        <InputText
          v-model="clientID"
          :placeholder="alreadyConfigured ? $t('system.oauthApps.editor.keepPlaceholder') : ''"
          class="w-full"
        />
      </CFormGroup>

      <CFormGroup
        :label="$t('system.oauthApps.editor.clientSecret')"
        :required="!alreadyConfigured"
      >
        <Password
          v-model="clientSecret"
          :placeholder="alreadyConfigured ? $t('system.oauthApps.editor.keepPlaceholder') : ''"
          toggle-mask
          :feedback="false"
          input-class="w-full"
          class="w-full"
        />
      </CFormGroup>
    </div>

    <template #footer>
      <Button
        v-if="alreadyConfigured"
        :label="$t('system.connections.authSetup.continue')"
        severity="secondary"
        text
        size="small"
        @click="emit('skip')"
      />
      <Button
        :label="$t('system.connections.authSetup.save')"
        icon="pi pi-arrow-right"
        icon-pos="right"
        size="small"
        :loading="saving"
        @click="handleSave"
      />
    </template>
  </Dialog>
</template>

<script setup>
import { inject, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  provider: { type: String, required: true },
  alreadyConfigured: { type: Boolean, default: false },
  visible: { type: Boolean, default: false },
})

const emit = defineEmits(['update:visible', 'saved', 'skip'])

const { t } = useI18n()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')
const $Auth = inject('$Auth')

const REDIRECT_PATH = '/oauth2/connection/callback'

const clientID = ref('')
const clientSecret = ref('')
const saving = ref(false)
const redirectURL = ref('')
const blueprint = ref(null)

function openURL(url) {
  if (url) window.open(url, '_blank', 'noopener')
}

// Catalog blueprint for this provider — drives the setup guidance.
async function loadBlueprint() {
  try {
    const res = await $SystemAPI.api().get('/connections/oauth-apps')
    const payload = res?.data?.response ?? res?.data
    const list = Array.isArray(payload) ? payload : []
    blueprint.value = list.find(p => p.handle === props.provider) || null
  } catch {
    blueprint.value = null
  }
}

async function loadRedirect() {
  const authBase = ($Auth?.authURL || `${window.location.origin}/auth`).replace(/\/$/, '')
  try {
    const current = await $SystemAPI.settingsCurrent()
    redirectURL.value = current?.connection?.oauth?.redirectURL || `${authBase}${REDIRECT_PATH}`
  } catch {
    redirectURL.value = `${authBase}${REDIRECT_PATH}`
  }
}

async function handleSave() {
  if (!props.alreadyConfigured && (!clientID.value.trim() || !clientSecret.value.trim())) {
    $toast.toastWarning(t('system.oauthApps.editor.credentialsRequired'))
    return
  }

  // Nothing typed on an already-configured provider — keep the stored credentials.
  if (!clientID.value.trim() && !clientSecret.value.trim()) {
    emit('skip')
    return
  }

  const prefix = `connection.oauth.apps.${props.provider}.`
  const values = [{ name: 'connection.oauth.redirect-url', value: redirectURL.value }]
  if (clientID.value.trim()) {
    values.push({ name: `${prefix}client-id`, value: clientID.value.trim() })
  }
  if (clientSecret.value.trim()) {
    values.push({ name: `${prefix}client-secret`, value: clientSecret.value.trim() })
  }

  saving.value = true
  try {
    await $SystemAPI.settingsUpdate({ values })
    $toast.toastSuccess(t('system.oauthApps.saveSuccess'))
    emit('saved')
  } catch (e) {
    $toast.toastErrorHandler(t('system.oauthApps.saveError'))(e)
  } finally {
    saving.value = false
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

// Clear typed credentials each time the dialog opens so a reopen never reuses
// a half-entered secret.
watch(
  () => props.visible,
  open => {
    if (open) {
      clientID.value = ''
      clientSecret.value = ''
    }
  },
)

onMounted(() => {
  loadBlueprint()
  loadRedirect()
})
</script>
