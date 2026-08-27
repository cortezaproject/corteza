<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.connections.list.title') }}</span>
  </Teleport>

  <CViewContainer>
    <div class="h-full overflow-y-auto">
      <div class="flex flex-col gap-6 pb-2">
        <div class="flex items-center justify-between gap-2 flex-wrap">
          <IconField class="w-full max-w-sm">
            <InputIcon class="pi pi-search" />
            <InputText
              v-model="search"
              :placeholder="$t('system.connections.list.searchConnectors')"
              class="w-full"
            />
          </IconField>

          <div class="flex gap-2">
            <CRouterLinkButton
              :to="{ name: 'system.connections.oauthApps' }"
              :label="$t('system.oauthApps.title')"
              severity="secondary"
              icon="pi pi-key"
              size="small"
            />
            <CRouterLinkButton
              v-if="canCreate"
              :to="{ name: 'system.connections.create' }"
              :label="$t('system.connections.list.createCustom')"
              severity="secondary"
              icon="pi pi-plus"
              size="small"
            />
            <CPermissionsButton
              v-if="canGrant"
              resource="corteza::system:connection/*"
              v-tooltip.bottom="$t('general.label.permissions')"
            />
          </div>
        </div>

        <div v-if="loading" class="flex justify-center p-10">
          <ProgressSpinner style="width: 2rem; height: 2rem" />
        </div>

        <div
          v-else-if="!filteredCatalog.length && !customConnections.length"
          class="text-sm text-muted-color italic p-6 text-center"
        >
          {{ $t('system.connections.list.noConnectors') }}
        </div>

        <template v-else>
          <div
            v-if="filteredCatalog.length"
            class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3"
          >
            <div
              v-for="c in filteredCatalog"
              :key="c.connectionID"
              class="flex flex-col gap-3 rounded-lg border border-surface-200 dark:border-surface-700 p-4 hover:border-primary transition-colors"
            >
              <div class="flex items-center gap-3">
                <ConnectorLogo :icon="c.meta?.icon" :name="c.meta?.short || c.handle" size="lg" />
                <div class="min-w-0 flex-1">
                  <div class="font-medium truncate">{{ c.meta?.short || c.handle }}</div>
                  <div v-if="c.meta?.tags?.length" class="flex gap-1 flex-wrap mt-1">
                    <span
                      v-for="tag in c.meta.tags.slice(0, 3)"
                      :key="tag"
                      class="text-[10px] uppercase tracking-wide text-muted-color bg-surface-100 dark:bg-surface-800 rounded px-1.5 py-0.5"
                    >
                      {{ tag }}
                    </span>
                  </div>
                </div>
              </div>

              <p class="text-sm text-muted-color line-clamp-2 flex-1 min-h-[2.5rem]">
                {{ c.meta?.description || '' }}
              </p>

              <div class="flex items-center justify-between">
                <span
                  v-if="isInstalled(c)"
                  class="inline-flex items-center gap-1 text-xs text-green-600"
                >
                  <i class="pi pi-check-circle" />
                  {{ $t('system.connections.list.installed') }}
                </span>
                <span v-else />
                <Button
                  :label="
                    isInstalled(c)
                      ? $t('system.connections.list.setUp')
                      : $t('system.connections.list.install')
                  "
                  :icon="isInstalled(c) ? 'pi pi-cog' : 'pi pi-download'"
                  :outlined="isInstalled(c)"
                  size="small"
                  :loading="installing === c.catalogID"
                  @click="handleInstall(c)"
                />
              </div>
            </div>
          </div>

          <div v-if="customConnections.length" class="flex flex-col gap-2">
            <span class="font-medium">{{ $t('system.connections.list.customTitle') }}</span>
            <div class="flex flex-col gap-2">
              <div
                v-for="c in customConnections"
                :key="c.connectionID"
                class="flex items-center gap-3 rounded-md border border-surface-200 dark:border-surface-700 p-3 cursor-pointer hover:border-primary transition-colors"
                @click="handleCustomClick(c)"
              >
                <div class="flex flex-col min-w-0 flex-1">
                  <span class="font-medium truncate">{{ c.meta?.short || c.handle }}</span>
                  <span v-if="c.meta?.description" class="text-xs text-muted-color truncate">
                    {{ c.meta.description }}
                  </span>
                </div>
                <Tag
                  :value="$t(`system.connections.editor.statusValues.${c.status || 'draft'}`)"
                  :severity="c.status === 'active' ? 'success' : 'secondary'"
                />
                <CInputDelete
                  v-if="c.canDeleteConnection"
                  :message="$t('system.connections.list.deleteConfirm')"
                  :header="c.meta?.short || c.handle"
                  size="small"
                  @click.stop
                  @confirm="handleDelete(c)"
                />
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>
  </CViewContainer>
</template>

<script setup>
import { components, useResourceList, useRBACStore } from '@planetcrust/human-vue'
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import ConnectorLogo from './ConnectorLogo.vue'

const { CInputDelete, CRouterLinkButton, CViewContainer } = components

const router = useRouter()
const { t } = useI18n()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const rbac = useRBACStore()
const canGrant = computed(() => rbac.can('system/', 'grant'))
const canCreate = computed(() => rbac.can('system/', 'connection.create'))

const search = ref('')
const installing = ref(null)

const {
  items: connectionList,
  loading,
  filterList,
} = useResourceList(params => $SystemAPI.connectionListCancellable(params), {
  filter: { query: '', deleted: '0' },
  sorting: { sortBy: 'status', sortDesc: false },
  pagination: { limit: 200 },
})

const catalogConnectors = computed(() =>
  (connectionList.value || []).filter(c => c.source === 'catalog'),
)
const customConnections = computed(() =>
  (connectionList.value || []).filter(c => c.source !== 'catalog'),
)

const filteredCatalog = computed(() => {
  const q = search.value.trim().toLowerCase()
  const list = catalogConnectors.value
  if (!q) return list
  return list.filter(c => {
    const m = c.meta || {}
    return (
      (m.short || '').toLowerCase().includes(q) ||
      (m.description || '').toLowerCase().includes(q) ||
      (m.tags || []).some(tag => String(tag).toLowerCase().includes(q))
    )
  })
})

// A catalog connector is installed once it is imported and active in the store.
function isInstalled(c) {
  return c.status === 'active'
}

function goConfigure(connectionID) {
  router.push({ name: 'system.connections.configure', params: { connectionID } })
}

async function handleInstall(c) {
  if (isInstalled(c)) {
    goConfigure(c.connectionID)
    return
  }
  installing.value = c.catalogID
  try {
    const imported = await $SystemAPI.connectionImport({ catalogID: c.catalogID })
    goConfigure(imported?.connectionID || c.connectionID)
  } catch (e) {
    $toast.toastErrorHandler(t('notification.connection.fetch.error'))(e)
  } finally {
    installing.value = null
  }
}

function handleCustomClick(c) {
  if (!c.canUpdateConnection && !c.canDeleteConnection) return
  router.push({ name: 'system.connections.edit', params: { connectionID: c.connectionID } })
}

async function handleDelete(c) {
  try {
    await $SystemAPI.connectionDelete({ connectionID: c.connectionID })
    $toast.toastSuccess(t('notification.connection.delete.success'))
    filterList()
  } catch (e) {
    $toast.toastErrorHandler(t('notification.connection.delete.error'))(e)
  }
}
</script>
