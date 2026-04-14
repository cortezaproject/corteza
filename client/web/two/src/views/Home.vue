<template>
  <div class="flex h-full p-3 gap-0 overflow-hidden bg-surface-ground">

    <!-- ── Apps column (left) ──────────────────────────────── -->
    <div
      class="column-panel flex flex-col shrink-0"
      :style="{ width: menuWidth + 'px' }"
    >
      <div class="flex items-center gap-2 px-3 py-2 border-b border-surface shrink-0">
        <i class="pi pi-th-large text-color-secondary text-sm" />
        <span class="font-semibold text-sm text-color">{{ $t('two.column.apps') }}</span>
      </div>

      <div class="px-3 pt-3 pb-2 shrink-0">
        <CInputSearch v-model="appsQuery" size="small" class="w-full" />
      </div>

      <div v-if="areAppsVisible" class="flex-1 overflow-y-auto px-3 pb-3">
        <div class="flex flex-col gap-2">
          <a
            v-for="app in apps"
            :key="app.applicationID"
            v-show="isAppVisible(app)"
            :href="app.enabled ? getAppUrl(app) : '#'"
            target="_self"
            class="flex items-center gap-3 p-3 rounded-lg border border-surface hover:bg-emphasis hover:border-primary transition-all duration-150 no-underline text-color"
            @click="!app.enabled && $event.preventDefault()"
          >
            <img
              :src="getAppLogoUrl(app)"
              :alt="app.unify?.name || app.name"
              class="w-10 h-10 object-contain rounded-md shrink-0"
              loading="lazy"
            />
            <span class="font-medium text-sm truncate">
              {{ app.unify?.name || app.name }}
            </span>
          </a>
        </div>
      </div>

      <div v-else class="flex-1 flex items-center justify-center px-3">
        <span class="text-muted-color text-sm text-center">
          {{ appsQuery ? $t('two.apps.noResults') : $t('two.apps.empty') }}
        </span>
      </div>
    </div>

    <!-- ── Resize handle: menu / agent ──────────────────────── -->
    <div
      class="resize-handle group"
      @mousedown="startMenuResize"
    >
      <div class="resize-handle-bar group-hover:opacity-100" />
    </div>

    <!-- ── Agent column (middle) ────────────────────────────── -->
    <div class="column-panel flex flex-col flex-1 min-w-0">

      <!-- Agent header -->
      <div class="flex items-center gap-2 px-3 py-2 border-b border-surface shrink-0">
        <i class="pi pi-sparkles text-color-secondary text-sm shrink-0" />
        <Select
          v-if="agentStore.availableAgents.length > 0"
          v-model="agentStore.activeAgentID"
          :options="agentStore.availableAgents"
          optionLabel="meta.short"
          optionValue="agentID"
          :pt="{
            root: { style: 'min-height:0; padding:0; border:0; box-shadow:none; background:transparent;' },
            label: { style: 'padding:0 0.375rem 0 0; line-height:1.25rem;' },
            dropdown: { style: 'padding:0; width:1rem;' },
          }"
          class="!min-w-0 w-auto shrink-0"
        >
          <template #value="slotProps">
            <span class="font-semibold text-sm text-color">
              {{
                agentStore.availableAgents.find(a => a.agentID === slotProps.value)?.meta?.short ||
                agentStore.availableAgents.find(a => a.agentID === slotProps.value)?.handle ||
                $t('two.column.assistant')
              }}
            </span>
          </template>
          <template #option="slotProps">
            <span>{{ slotProps.option.meta?.short || slotProps.option.handle }}</span>
          </template>
        </Select>
        <span v-else class="font-semibold text-sm text-color">{{ $t('two.column.assistant') }}</span>
      </div>

      <!-- Conversation tabs -->
      <div
        v-if="hasStartedConversations"
        class="flex items-center gap-0 border-b border-surface shrink-0 bg-surface-ground"
      >
        <div class="flex items-center gap-0 flex-1 overflow-x-auto no-scrollbar">
          <button
            v-for="(conv, idx) in conversations"
            :key="idx"
            class="group flex items-center gap-1.5 p-3 text-sm border-b-2 whitespace-nowrap transition-colors duration-200 outline-none select-none max-w-[150px]"
            :class="[
              activeConversationIndex === idx
                ? 'border-b-primary text-primary font-medium'
                : 'border-b-transparent text-muted-color hover:text-color hover:border-b-surface-border',
            ]"
            @click="agentStore.setActiveConversationIndex(agentStore.activeAgentID, idx)"
          >
            <span class="whitespace-nowrap truncate">
              {{ $t('agent.sidebar.chatTab', { id: idx + 1 }) }}
            </span>
            <i
              class="pi pi-times text-sm opacity-0 group-hover:opacity-100 p-1 hover:bg-surface rounded-full transition-all shrink-0"
              @click.stop="agentStore.closeConversation(agentStore.activeAgentID, idx)"
            />
          </button>
        </div>
        <div class="flex items-center px-1 border-l border-surface shrink-0">
          <Button
            icon="pi pi-plus"
            severity="secondary"
            variant="text"
            rounded
            size="small"
            class="!w-7 !h-7"
            v-tooltip.bottom="$t('agent.sidebar.newChat')"
            @click="agentStore.startNewConversation(agentStore.activeAgentID)"
          />
          <Button
            v-if="conversations.length > 1"
            icon="pi pi-trash"
            severity="danger"
            variant="text"
            rounded
            size="small"
            class="!w-7 !h-7"
            v-tooltip.bottom="$t('agent.sidebar.clearAllChats')"
            @click="agentStore.clearConversation(agentStore.activeAgentID)"
          />
        </div>
      </div>

      <!-- No agents available -->
      <div
        v-if="agentStore.availableAgents.length === 0"
        class="flex-1 flex flex-col items-center justify-center gap-2 p-6 text-center"
      >
        <i class="pi pi-sparkles text-3xl text-muted-color" />
        <span class="text-sm text-muted-color">{{ $t('two.assistant.noAgents') }}</span>
      </div>

      <!-- Chat messages -->
      <div v-else class="flex min-h-0 flex-1 flex-col">
        <div class="flex-1 overflow-y-auto p-3 space-y-4 flex flex-col" ref="chatContainer">
          <template v-if="activeConversation && filteredMessages.length > 0">
            <div
              v-for="(msg, index) in filteredMessages"
              :key="index"
              class="flex flex-col max-w-[90%]"
              :class="msg.role === 'user' ? 'self-end items-end' : 'self-start items-start'"
            >
              <div
                class="px-3 py-2 rounded-2xl shadow-sm text-sm"
                :class="
                  msg.role === 'user'
                    ? 'bg-primary text-primary-contrast'
                    : 'bg-surface text-color border border-surface-border'
                "
              >
                <div
                  v-if="msg.role === 'agent'"
                  class="markdown-body"
                  v-html="renderMarkdown(msg.content)"
                />
                <div v-else class="whitespace-pre-wrap">{{ msg.content }}</div>
              </div>
            </div>
          </template>
          <div
            v-else
            class="flex h-full items-center justify-center text-muted-color p-4 text-center text-sm flex-1"
          >
            {{ $t('agent.sidebar.empty') }}
          </div>

          <div v-if="executing" class="self-start items-start mt-auto">
            <div class="px-3 py-2 rounded-2xl shadow-sm text-sm bg-surface text-color border border-surface-border flex gap-1 items-center h-[34px]">
              <span class="w-1.5 h-1.5 bg-muted-color rounded-full animate-bounce" />
              <span class="w-1.5 h-1.5 bg-muted-color rounded-full animate-bounce" style="animation-delay: 0.2s" />
              <span class="w-1.5 h-1.5 bg-muted-color rounded-full animate-bounce" style="animation-delay: 0.4s" />
            </div>
          </div>
        </div>

        <!-- Chat input -->
        <div class="p-3 border-t border-surface flex items-end gap-2 shrink-0 bg-surface">
          <Textarea
            v-model="chatInput"
            autoResize
            rows="1"
            class="flex-1 text-sm block"
            style="min-height: 2.5rem; max-height: 10rem"
            :placeholder="$t('agent.sidebar.placeholder')"
            @keydown.enter.prevent="sendMessage"
            :disabled="executing || !agentStore.activeAgentID"
          />
          <Button
            icon="pi pi-send"
            @click="sendMessage"
            :disabled="!chatInput.trim() || executing || !agentStore.activeAgentID"
            :loading="executing"
          />
        </div>
      </div>
    </div>

    <!-- ── Resize handle: agent / notifications ──────────────── -->
    <div
      class="resize-handle group"
      @mousedown="startNotificationsResize"
    >
      <div class="resize-handle-bar group-hover:opacity-100" />
    </div>

    <!-- ── Notifications column (right) ─────────────────────── -->
    <div
      class="column-panel flex flex-col shrink-0"
      :style="{ width: notificationsWidth + 'px' }"
    >
      <div class="flex items-center gap-2 px-3 py-2 border-b border-surface shrink-0">
        <i class="pi pi-bell text-color-secondary text-sm" />
        <span class="font-semibold text-sm text-color">{{ $t('two.column.notifications') }}</span>
      </div>

      <div class="flex-1 min-h-0 overflow-hidden">
        <Notifications />
      </div>
    </div>

  </div>
</template>

<script setup>
import {
  components,
  resolveAppLogoUrl,
  useApplicationsStore,
  useAgentSidebarStore,
  useNotificationsStore,
} from '@cortezaproject/corteza-vue-next'
import { computed, inject, nextTick, onMounted, ref, watch } from 'vue'
import { useColumnResize } from '../composables/useColumnResize'

const { CInputSearch, Notifications } = components

const $SystemAPI = inject('$SystemAPI')
const $Auth = inject('$Auth')

// ── Column resize ────────────────────────────────────────────
const { menuWidth, notificationsWidth, startMenuResize, startNotificationsResize } = useColumnResize()

// ── Apps column ──────────────────────────────────────────────
const applicationsStore = useApplicationsStore()
const appsQuery = ref('')

const apps = computed(() => applicationsStore.unifyOnly)
const normalizedAppsQuery = computed(() => (appsQuery.value || '').trim().toUpperCase())
const isAppVisible = app => {
  const q = normalizedAppsQuery.value
  if (!q) return true
  return (
    (app.name?.toUpperCase() || '').includes(q) ||
    (app.unify?.name?.toUpperCase() || '').includes(q)
  )
}
const areAppsVisible = computed(() => apps.value.some(isAppVisible))

const getAppLogoUrl = app => resolveAppLogoUrl(app, $SystemAPI.baseURL)

const getAppUrl = app => {
  const url = app.unify?.url || ''
  if (!url || url.startsWith('/') || url.startsWith('http')) return url
  return '/' + url
}

// ── Notifications column ─────────────────────────────────────
const notificationsStore = useNotificationsStore()

// ── Agent column ─────────────────────────────────────────────
const agentStore = useAgentSidebarStore()
const chatInput = ref('')
const executing = ref(false)
const chatContainer = ref(null)

const activeConversation = computed(() => {
  if (!agentStore.activeAgentID) return null
  return agentStore.getConversation(agentStore.activeAgentID)
})

const conversations = computed(() => {
  if (!agentStore.activeAgentID) return []
  return agentStore.getAllConversations(agentStore.activeAgentID) || []
})

const activeConversationIndex = computed(() => {
  if (!agentStore.activeAgentID) return 0
  return agentStore.activeConversationIndex[agentStore.activeAgentID] || 0
})

const filteredMessages = computed(() => {
  if (!activeConversation.value) return []
  return activeConversation.value.messages.filter(m => m.content)
})

const hasStartedConversations = computed(() => {
  return conversations.value.some(c => c.messages.length > 0)
})

const scrollToBottom = async () => {
  await nextTick()
  if (chatContainer.value) {
    chatContainer.value.scrollTop = chatContainer.value.scrollHeight
  }
}

watch(activeConversationIndex, () => scrollToBottom())

const sendMessage = async () => {
  if (!chatInput.value.trim() || executing.value || !agentStore.activeAgentID) return

  const input = chatInput.value
  chatInput.value = ''

  const currentAgentID = agentStore.activeAgentID
  agentStore.addMessage(currentAgentID, { role: 'user', content: input })

  executing.value = true
  scrollToBottom()

  try {
    const activeConv = activeConversation.value
    const res = await $SystemAPI.agentExec({
      agentID: currentAgentID,
      input: input,
      ...(activeConv?.conversationID ? { conversationID: activeConv.conversationID } : {}),
    })

    if (res?.conversationID) {
      agentStore.setConversationID(currentAgentID, res.conversationID)
    }

    const outputContent =
      res?.output || (typeof res === 'string' ? res : res?.response?.text || JSON.stringify(res))

    agentStore.addMessage(currentAgentID, { role: 'agent', content: outputContent })
  } catch (err) {
    console.error('Agent execution error:', err)
    agentStore.addMessage(currentAgentID, { role: 'agent', content: 'Error: ' + err.message })
  } finally {
    executing.value = false
    scrollToBottom()
  }
}

// Load available agents on mount (same logic as CAgentSidebarButton)
onMounted(async () => {
  if (agentStore.availableAgents.length > 0) return

  try {
    const res = await $SystemAPI.agentList({ limit: 0 })
    const allAgents = res.set || []
    const userRoles = $Auth.user?.roles || []
    const configuredAgents = allAgents.filter(agent => {
      if (!agent.invocation?.user?.enabled) return false
      const sidebarRoles = agent.meta?.sidebarRoles || []
      if (!Array.isArray(sidebarRoles) || sidebarRoles.length === 0) return false
      return userRoles.includes('2') || sidebarRoles.some(r => userRoles.includes(r))
    })
    agentStore.setAvailableAgents(configuredAgents)
  } catch (err) {
    console.warn('Failed to load agents', err)
  }
})

// Simple markdown rendering (same as CAgentSidebar)
function renderMarkdown(text) {
  if (!text) return ''

  let html = text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')

  html = html.replace(
    /```(\w*)\n([\s\S]*?)```/g,
    '<pre class="bg-surface p-2 rounded my-2 overflow-x-auto text-xs font-mono border border-surface-border"><code>$2</code></pre>',
  )
  html = html.replace(
    /`([^`]+)`/g,
    '<code class="bg-surface px-1 py-0.5 rounded text-xs font-mono border border-surface-border">$1</code>',
  )
  html = html.replace(/^### (.+)$/gm, '<h4 class="font-bold mt-2 mb-1">$1</h4>')
  html = html.replace(/^## (.+)$/gm, '<h3 class="font-bold text-base mt-2 mb-1">$1</h3>')
  html = html.replace(/^# (.+)$/gm, '<h2 class="font-bold text-lg mt-2 mb-1">$1</h2>')
  html = html.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
  html = html.replace(/\*(.+?)\*/g, '<em>$1</em>')
  html = html.replace(/^[*-] (.+)$/gm, '<li class="ml-4 list-disc">$1</li>')
  html = html.replace(/^\d+\. (.+)$/gm, '<li class="ml-4 list-decimal">$1</li>')
  html = html.replace(/\n\n/g, '</p><p class="my-1">')
  html = html.replace(/\n/g, '<br>')

  return `<div class="prose-sm">${html}</div>`
}
</script>

<style scoped>
/* Matches .right-sidebar visual style from useTheme.ts */
.column-panel {
  background-color: var(--p-content-background);
  border: 1px solid var(--p-content-border-color);
  border-radius: var(--p-border-radius-xl);
  box-shadow: var(--p-overlay-popover-shadow);
  overflow: hidden;
}

/* Slim drag handle between columns */
.resize-handle {
  width: 12px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: ew-resize;
}

.resize-handle-bar {
  width: 2px;
  height: 2rem;
  border-radius: 9999px;
  background-color: var(--p-content-border-color);
  opacity: 0.4;
  transition: opacity 150ms;
}

.markdown-body :deep(p) {
  margin-top: 0.25rem;
  margin-bottom: 0.25rem;
}
.markdown-body :deep(p:first-child) { margin-top: 0; }
.markdown-body :deep(p:last-child) { margin-bottom: 0; }
</style>
