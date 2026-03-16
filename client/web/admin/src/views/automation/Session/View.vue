<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('automation.sessions.editor.title') }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <div v-else-if="session" class="flex flex-col h-full">
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <Panel :header="$t('automation.sessions.editor.info.title')" toggleable :collapsed="false">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="flex flex-col gap-1">
            <span class="text-xs text-muted-color">
              {{ $t('automation.sessions.editor.info.id') }}
            </span>
            <span class="font-mono text-sm">{{ session.sessionID }}</span>
          </div>

          <div class="flex flex-col gap-1">
            <span class="text-xs text-muted-color">
              {{ $t('automation.sessions.editor.info.workflowID') }}
            </span>
            <span class="font-mono text-sm">{{ session.workflowID }}</span>
          </div>

          <div class="flex flex-col gap-1">
            <span class="text-xs text-muted-color">
              {{ $t('automation.sessions.editor.info.status') }}
            </span>
            <Tag
              :value="session.status"
              :severity="statusSeverity(session.status)"
              class="self-start"
            />
          </div>

          <div class="flex flex-col gap-1">
            <span class="text-xs text-muted-color">
              {{ $t('automation.sessions.editor.info.createdAt') }}
            </span>
            <span>{{ locFullDateTime(session.createdAt) }}</span>
          </div>

          <div v-if="session.completedAt" class="flex flex-col gap-1">
            <span class="text-xs text-muted-color">
              {{ $t('automation.sessions.editor.info.completedAt') }}
            </span>
            <span>{{ locFullDateTime(session.completedAt) }}</span>
          </div>

          <div v-if="session.createdBy" class="flex flex-col gap-1">
            <span class="text-xs text-muted-color">
              {{ $t('automation.sessions.editor.info.createdBy') }}
            </span>
            <span class="font-mono text-sm">{{ session.createdBy }}</span>
          </div>
        </div>
      </Panel>

      <Panel
        v-if="session.error"
        :header="$t('automation.sessions.editor.info.error')"
        toggleable
        :collapsed="false"
      >
        <pre
          class="text-sm text-red-500 whitespace-pre-wrap break-all font-mono bg-highlight rounded p-3"
          >{{ session.error }}</pre
        >
      </Panel>
    </div>

    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="p-3 flex items-center justify-between">
        <Button
          :label="$t('general.label.back')"
          icon="pi pi-arrow-left"
          severity="secondary"
          @click="$router.push({ name: 'automation.sessions' })"
        />
        <Button
          v-if="isActive"
          :label="$t('automation.sessions.editor.info.cancel')"
          icon="pi pi-times"
          severity="danger"
          :loading="canceling"
          @click="handleCancel"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { filters } from '@cortezaproject/corteza-vue-next'

const { locFullDateTime } = filters

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const $toast = inject('$toast')
const $AutomationAPI = inject('$AutomationAPI')

const loading = ref(false)
const canceling = ref(false)
const session = ref(null)

const isActive = computed(() => {
  return session.value && ['pending', 'started'].includes(session.value.status)
})

function statusSeverity(status) {
  switch (status) {
    case 'completed':
      return 'success'
    case 'failed':
      return 'danger'
    case 'canceled':
      return 'secondary'
    case 'started':
    case 'pending':
      return 'info'
    default:
      return 'secondary'
  }
}

async function loadSession() {
  loading.value = true
  try {
    session.value = await $AutomationAPI.sessionRead({ sessionID: route.params.sessionID })
  } catch (e) {
    $toast.toastErrorHandler(t('notification.session.fetch.error'))(e)
    router.push({ name: 'automation.sessions' })
  } finally {
    loading.value = false
  }
}

async function handleCancel() {
  canceling.value = true
  try {
    await $AutomationAPI.sessionCancel({ sessionID: session.value.sessionID })
    $toast.toastSuccess(t('notification.session.cancel.success'))
    await loadSession()
  } catch (e) {
    $toast.toastErrorHandler(t('notification.session.cancel.error'))(e)
  } finally {
    canceling.value = false
  }
}

onMounted(() => {
  loadSession()
})
</script>
