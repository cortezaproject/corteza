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
              {{ $t('automation.sessions.editor.info.createdByUserName') }}
            </span>
            <span class="font-mono text-sm">{{ createdByName }}</span>
          </div>

          <div v-if="session.eventType" class="flex flex-col gap-1">
            <span class="text-xs text-muted-color">
              {{ $t('automation.sessions.editor.info.eventType') }}
            </span>
            <span class="font-mono text-sm">{{ session.eventType }}</span>
          </div>

          <div v-if="session.resourceType" class="flex flex-col gap-1">
            <span class="text-xs text-muted-color">
              {{ $t('automation.sessions.editor.info.resourceType') }}
            </span>
            <span class="font-mono text-sm">{{ session.resourceType }}</span>
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
        <div class="flex gap-2">
          <Button
            v-if="session.workflowID"
            :label="$t('automation.sessions.editor.info.openWorkflow')"
            icon="pi pi-arrow-right"
            severity="secondary"
            @click="$router.push({ name: 'automation.workflows.edit', params: { workflowID: session.workflowID } })"
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
const $SystemAPI = inject('$SystemAPI')

const loading = ref(false)
const canceling = ref(false)
const session = ref(null)
const user = ref(null)

const createdByName = computed(() => {
  if (!user.value) return session.value?.createdBy
  const { name, username, email, userID } = user.value
  return name || username || email || `<@${userID}>`
})

const isActive = computed(() => session.value && !session.value.completedAt)

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
    if (session.value.createdBy) {
      try {
        user.value = await $SystemAPI.userRead({ userID: session.value.createdBy })
      } catch {
        // non-fatal, fall back to raw ID
      }
    }
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
