<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ headerTitle }}</span>
  </Teleport>

  <div class="flex flex-col h-full">
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4">
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
          <CPermissionsButton
            v-if="canGrant && connection?.connectionID && connection?.status === 'active'"
            v-tooltip.bottom="$t('general.label.permissions')"
            :resource="`corteza::system:dal-connection/${connection.connectionID}`"
            :title="connection.meta?.short || connection.handle || connection.connectionID"
            :target="connection.meta?.short || connection.handle || connection.connectionID"
          />
        </div>
      </template>
    </Card>

    <ConfiguredConnectionsPanel v-if="connection" :connection="connection" class="flex-1 min-h-0" />
    </div>

    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="p-3 flex items-center">
        <Button
          :label="$t('general.label.back')"
          icon="pi pi-arrow-left"
          severity="secondary"
          @click="$router.push({ name: 'system.connections' })"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useRBACStore } from '@planetcrust/human-vue'
import ConfiguredConnectionsPanel from './ConfiguredConnectionsPanel.vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const rbac = useRBACStore()
const canGrant = computed(() => rbac.can('system/', 'grant'))

const connection = ref(null)

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
  } catch (e) {
    $toast.toastErrorHandler(t('notification.connection.fetch.error'))(e)
    router.push({ name: 'system.connections' })
  }
}

onMounted(() => loadConnection())

watch(() => route.params.connectionID, loadConnection)
</script>
