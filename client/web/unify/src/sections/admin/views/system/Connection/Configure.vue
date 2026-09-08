<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ headerTitle }}</span>
  </Teleport>

  <div class="flex flex-col h-full">
    <CViewContainer scroll>
      <Card
        :pt="{
          body: { class: 'p-4 flex flex-col gap-2' },
          content: { class: 'p-0' },
        }"
        class="shrink-0"
      >
        <template #content>
          <div class="flex items-start justify-between gap-2">
            <div class="flex flex-col gap-2 min-w-0">
              <span class="text-lg font-medium truncate">
                {{ connection?.meta?.short || '-' }}
              </span>
              <span v-if="connection?.meta?.description" class="text-sm text-muted-color">
                {{ connection.meta.description }}
              </span>
            </div>
            <div class="flex items-center gap-2 shrink-0">
              <Button
                v-if="connection?.source === 'catalog'"
                :label="$t('system.connections.resync.button')"
                icon="pi pi-sync"
                severity="secondary"
                outlined
                size="small"
                :loading="resyncing"
                v-tooltip.bottom="$t('system.connections.resync.hint')"
                @click="handleResync"
              />
              <Button
                v-if="oauthAppHandle"
                :label="$t('system.connections.authSetup.credentialsButton')"
                icon="pi pi-key"
                severity="secondary"
                outlined
                size="small"
                @click="authModal = true"
              />
              <CInputDelete
                v-if="connection?.canDeleteConnection"
                :label="$t('system.connections.list.uninstall')"
                :message="$t('system.connections.list.uninstallConfirm')"
                :header="connection.meta?.short || connection.handle"
                icon="pi pi-trash"
                severity="secondary"
                outlined
                size="small"
                @confirm="handleUninstall"
              />
              <CPermissionsButton
                v-if="canGrant && connection?.connectionID && connection?.status === 'active'"
                v-tooltip.bottom="$t('general.label.permissions')"
                :resource="`corteza::system:connection/${connection.connectionID}`"
                :title="connection.meta?.short || connection.handle || connection.connectionID"
                :target="connection.meta?.short || connection.handle || connection.connectionID"
              />
            </div>
          </div>
        </template>
      </Card>

      <ConfiguredConnectionsPanel
        v-if="connection"
        :connection="connection"
        class="flex-1 min-h-0"
      />
    </CViewContainer>

    <ConnectionAuthSetup
      v-if="oauthAppHandle"
      v-model:visible="authModal"
      :provider="oauthAppHandle"
      :already-configured="providerConfigured"
      @saved="handleAuthSaved"
      @skip="authModal = false"
    />

    <CEditorActions :back-to="{ name: 'system.connections' }" />
  </div>
</template>

<script setup>
import { computed, inject, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { components, useRBACStore } from '@planetcrust/human-vue'
import ConfiguredConnectionsPanel from './ConfiguredConnectionsPanel.vue'
import ConnectionAuthSetup from './ConnectionAuthSetup.vue'

const { CViewContainer, CInputDelete } = components

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const rbac = useRBACStore()
const canGrant = computed(() => rbac.can('system/', 'grant'))

const connection = ref(null)

// OAuth connectors register their client id/secret in a setup dialog. It opens
// on its own when the provider has no credentials yet, or via the header button.
const OAUTH_METHOD = 'oauth2_authorization_code'
const authModal = ref(false)
const oauthAppHandle = ref('')
const providerConfigured = ref(false)

// The OAuth app handle the connector authenticates with, if any.
function oauthHandleOf(conn) {
  const svc = conn?.service
  if (!svc) return ''
  if (svc.auth?.method === OAUTH_METHOD) return svc.auth.oauthApp || ''
  const opt = (svc.authOptions || []).find(o => o.method === OAUTH_METHOD)
  return opt?.oauthApp || ''
}

// A provider is configured once its credentials live in settings; the backend
// only surfaces an app entry when at least one of its keys is set.
async function isProviderConfigured(handle) {
  if (!handle) return true
  try {
    const current = await $SystemAPI.settingsCurrent()
    const apps = current?.connection?.oauth?.apps || []
    return apps.some(a => a?.handle === handle)
  } catch {
    return false
  }
}

function handleAuthSaved() {
  providerConfigured.value = true
  authModal.value = false
}

const resyncing = ref(false)

// Re-import the connector definition from the catalog and refresh its automation
// functions live. Typed client method isn't generated, so call the axios instance.
async function handleResync() {
  if (!connection.value?.connectionID) return
  resyncing.value = true
  try {
    await $SystemAPI.api().post(`/connections/${connection.value.connectionID}/resync`)
    $toast.toastSuccess(t('system.connections.resync.success'))
    await loadConnection()
  } catch (e) {
    const detail = e?.response?.data?.error?.message || e?.message
    const prefix = t('system.connections.resync.error')
    $toast.toastDanger(detail ? `${prefix}: ${detail}` : prefix)
  } finally {
    resyncing.value = false
  }
}

// Uninstall removes the imported connection; the backend rejects it while any
// configured connection still references it, so surface that detail.
async function handleUninstall() {
  try {
    await $SystemAPI.connectionDelete({ connectionID: connection.value.connectionID })
    $toast.toastSuccess(t('notification.connection.delete.success'))
    router.push({ name: 'system.connections' })
  } catch (e) {
    const detail = e?.response?.data?.error?.message || e?.message
    const prefix = t('notification.connection.delete.error')
    $toast.toastDanger(detail ? `${prefix}: ${detail}` : prefix)
  }
}

const headerTitle = computed(() => {
  return (
    connection.value?.meta?.short ||
    connection.value?.handle ||
    t('system.connections.configurePage.title')
  )
})

async function loadConnection() {
  try {
    const loaded = await $SystemAPI.connectionRead({
      connectionID: route.params.connectionID,
    })
    if (loaded?.source === 'local' && loaded?.status === 'draft') {
      router.replace({
        name: 'system.connections.edit',
        params: { connectionID: route.params.connectionID },
      })
      return
    }
    connection.value = loaded

    const handle = oauthHandleOf(loaded)
    oauthAppHandle.value = handle
    providerConfigured.value = await isProviderConfigured(handle)
    // Pop the setup dialog first when the provider still needs credentials.
    authModal.value = !!handle && !providerConfigured.value
  } catch (e) {
    $toast.toastErrorHandler(t('notification.connection.fetch.error'))(e)
    router.push({ name: 'system.connections' })
  }
}

onMounted(() => loadConnection())

watch(() => route.params.connectionID, loadConnection)
</script>
