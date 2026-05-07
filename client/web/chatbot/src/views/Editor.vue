<template>
  <Teleport to="#topbar-title" defer>
    <span v-if="isCreate">{{ $t('chatbot.editor.titleCreate') }}</span>
    <span v-else class="font-semibold">
      {{ chatbot?.name || chatbot?.handle || $t('chatbot.editor.titleEdit') }}
    </span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <div v-else-if="chatbot" class="flex flex-col h-full overflow-hidden">
    <div class="flex-1 flex min-h-0 overflow-hidden">
      <div class="flex-1 overflow-y-auto">
        <div class="p-4 flex flex-col gap-4">
          <div v-if="!isCreate" class="flex justify-end">
            <CPermissionsButton
              v-tooltip.bottom="$t('general.label.permissions')"
              :resource="`corteza::system:chatbot/${route.params.chatbotID}`"
              :title="chatbot?.name || chatbot?.handle || route.params.chatbotID"
              :target="chatbot?.name || chatbot?.handle || route.params.chatbotID"
            />
          </div>
          <General :chatbot="chatbot" :is-create="isCreate" @regenerate-key="handleRegenerateKey" />
          <Scenarios :scenarios="chatbot.scenarios" :agents="agents" />
          <Styling :styling="chatbot.styling" :chatbot-id="chatbot.chatbotID" />
        </div>
      </div>
      <div class="hidden md:block w-[390px] shrink-0 relative border-l border-surface">
        <Preview :chatbot="chatbot" />
      </div>
    </div>

    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="flex items-center justify-between p-3">
        <Button
          :label="$t('general.label.back')"
          icon="pi pi-arrow-left"
          severity="secondary"
          @click="router.back()"
        />
        <div class="flex gap-2">
          <CInputDelete
            v-if="!isCreate && chatbot.canDeleteChatbot"
            :label="$t('general.label.delete')"
            :message="$t('chatbot.list.delete')"
            :header="chatbot?.name || chatbot?.handle || $t('general.label.delete')"
            @confirm="handleDelete"
          />
          <Button
            type="button"
            :label="$t('general.label.save')"
            icon="pi pi-save"
            :loading="saving"
            :disabled="!canSave"
            @click="handleSubmit"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { components, useUnsavedGuard } from '@planetcrust/human-vue'
import { system } from '@planetcrust/human-js'
import { cloneDeep, isEqual } from 'lodash-es'
import { computed, inject, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useChatbotStore } from '@/stores/chatbot'

import General from './Editor/General.vue'
import Styling from './Editor/Styling.vue'
import Scenarios from './Editor/Scenarios.vue'
import Preview from './Editor/Preview.vue'

const { CInputDelete } = components

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')
const chatbotStore = useChatbotStore()

const loading = ref(false)
const saving = ref(false)
const chatbot = ref(null)
const initialChatbot = ref(null)
const agents = ref([])

const isCreate = computed(() => !route.params.chatbotID)

const canSave = computed(() => {
  if (!chatbot.value) return false
  if (saving.value) return false
  if (!chatbot.value.name?.trim() && !chatbot.value.handle?.trim()) return false
  const ids = (chatbot.value.scenarios || []).map(s => s.id)
  if (ids.some(id => !id)) return false
  if (new Set(ids).size !== ids.length) return false
  return true
})

const { markSaved } = useUnsavedGuard({
  isDirty: () =>
    !saving.value &&
    !!chatbot.value &&
    !!initialChatbot.value &&
    !isEqual(chatbot.value, initialChatbot.value),
  messageKey: 'general.editor.unsavedChanges',
})

function applyChatbot(res) {
  chatbot.value = new system.Chatbot(res)
  initialChatbot.value = cloneDeep(chatbot.value)
}

async function loadChatbot() {
  const chatbotID = route.params.chatbotID
  if (!chatbotID) {
    applyChatbot({})
    return
  }

  loading.value = true
  try {
    const res = await $SystemAPI.chatbotRead({ chatbotID })
    applyChatbot(res)
  } catch {
    $toast.toastDanger(t('notification.chatbot.loadFailed'))
    router.push({ name: 'root' })
  } finally {
    loading.value = false
  }
}

async function loadAgents() {
  try {
    const res = await $SystemAPI.agentList({ limit: 500, sort: 'name ASC' })
    agents.value = res?.set || []
  } catch (e) {
    if (e?.message !== 'canceled') {
      console.error('Failed to load agents:', e)
    }
  }
}

async function handleSubmit() {
  if (!canSave.value) return
  saving.value = true
  try {
    const payload = JSON.parse(JSON.stringify(chatbot.value))
    if (payload.styling) {
      delete payload.styling.logoURL
      if (payload.styling.logoAttachmentID === '') delete payload.styling.logoAttachmentID
      if (payload.styling.launcher) {
        delete payload.styling.launcher.iconURL
        if (payload.styling.launcher.iconAttachmentID === '')
          delete payload.styling.launcher.iconAttachmentID
      }
    }
    if (isCreate.value) {
      const created = await $SystemAPI.chatbotCreate(payload)
      chatbotStore.updateInList(created)
      $toast.toastSuccess(t('notification.chatbot.created'))
      markSaved()
      router.push({ name: 'chatbot.edit', params: { chatbotID: created.chatbotID } })
    } else {
      const updated = await $SystemAPI.chatbotUpdate({
        chatbotID: route.params.chatbotID,
        ...payload,
      })
      applyChatbot(updated)
      chatbotStore.updateInList(updated)
      $toast.toastSuccess(t('notification.chatbot.saved'))
    }
  } catch (err) {
    console.error(err)
    $toast.toastErrorHandler(t('notification.chatbot.saveFailed'))(err)
  } finally {
    saving.value = false
  }
}

async function handleRegenerateKey() {
  if (!route.params.chatbotID) return
  try {
    const updated = await $SystemAPI.chatbotRegenerateWidgetKey({
      chatbotID: route.params.chatbotID,
    })
    applyChatbot(updated)
    $toast.toastSuccess(t('chatbot.editor.widgetKey.regenerated'))
  } catch (err) {
    $toast.toastErrorHandler(t('chatbot.editor.widgetKey.regenerateFailed'))(err)
  }
}

async function handleDelete() {
  try {
    await $SystemAPI.chatbotDelete({ chatbotID: route.params.chatbotID })
    chatbotStore.removeFromList(route.params.chatbotID)
    $toast.toastSuccess(t('notification.chatbot.deleted'))
    initialChatbot.value = cloneDeep(chatbot.value)
    router.push({ name: 'root' })
  } catch (err) {
    console.error(err)
    $toast.toastErrorHandler(t('notification.chatbot.deleteFailed'))(err)
  }
}

watch(
  () => route.params.chatbotID,
  () => {
    loadChatbot()
  },
  { immediate: true },
)

onMounted(() => {
  loadAgents()
})
</script>
