<template>
  <div class="flex flex-col min-h-0 flex-1">
    <!-- Header -->
    <div class="flex items-center gap-2 px-3 py-2 border-b border-surface shrink-0">
      <i class="pi pi-sparkles text-primary text-sm shrink-0" />
      <div class="flex-1 min-w-0">
        <Select
          v-if="agentStore.availableAgents.length > 0"
          v-model="agentStore.activeAgentID"
          :options="agentStore.availableAgents"
          optionLabel="meta.short"
          optionValue="agentID"
          size="small"
          class="!border-0 !shadow-none !bg-transparent"
        >
          <template #value="slotProps">
            <span v-if="slotProps.value" class="font-semibold text-base text-color truncate">
              {{
                agentStore.availableAgents.find(a => a.agentID === slotProps.value)?.meta?.short ||
                agentStore.availableAgents.find(a => a.agentID === slotProps.value)?.handle
              }}
            </span>
          </template>
          <template #option="slotProps">
            <span>{{ slotProps.option.meta?.short || slotProps.option.handle }}</span>
          </template>
        </Select>
        <span v-else class="font-semibold text-base text-color">{{ titleFallback || $t('agent.sidebar.title') }}</span>
      </div>
      <slot name="actions" />
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

    <!-- No agents state -->
    <div
      v-if="agentStore.availableAgents.length === 0"
      class="flex-1 flex flex-col items-center justify-center gap-2 p-6 text-center"
    >
      <i class="pi pi-sparkles text-3xl text-muted-color" />
      <span class="text-sm text-muted-color">{{ $t('agent.sidebar.noAgents') }}</span>
    </div>

    <!-- Chat -->
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

      <!-- Input -->
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
</template>

<script setup lang="ts">
import { computed, inject, nextTick, onMounted, ref, watch } from 'vue'
import { useAgentSidebarStore } from '../../stores/useAgentSidebarStore'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  titleFallback: {
    type: String,
    default: '',
  },
})

const $SystemAPI = inject<any>('$SystemAPI')
const $Auth = inject<any>('$Auth')
const agentStore = useAgentSidebarStore()

onMounted(async () => {
  if (agentStore.availableAgents.length > 0) return

  try {
    const res = await $SystemAPI.agentList({ limit: 0 })
    const allAgents = res.set || []
    const userRoles = $Auth?.user?.roles || []
    const configuredAgents = allAgents.filter((agent: any) => {
      if (!agent.invocation?.user?.enabled) return false
      const sidebarRoles = agent.meta?.sidebarRoles || []
      if (!Array.isArray(sidebarRoles) || sidebarRoles.length === 0) return false
      return userRoles.includes('2') || sidebarRoles.some((r: string) => userRoles.includes(r))
    })
    agentStore.setAvailableAgents(configuredAgents)
  } catch (err) {
    console.warn('Failed to load agents', err)
  }
})

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
  return (activeConversation.value as any).messages.filter((m: any) => m.content)
})

const hasStartedConversations = computed(() =>
  conversations.value.some((c: any) => c.messages.length > 0),
)

const chatInput = ref('')
const executing = ref(false)
const chatContainer = ref<HTMLElement | null>(null)

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
    const activeConv = activeConversation.value as any
    const res = await $SystemAPI.agentExec({
      agentID: currentAgentID,
      input,
      ...(activeConv?.conversationID ? { conversationID: activeConv.conversationID } : {}),
    })

    if (res?.conversationID) {
      agentStore.setConversationID(currentAgentID, res.conversationID)
    }

    const outputContent =
      res?.output || (typeof res === 'string' ? res : res?.response?.text || JSON.stringify(res))

    agentStore.addMessage(currentAgentID, { role: 'agent', content: outputContent })
  } catch (err: any) {
    console.error('Agent execution error:', err)
    agentStore.addMessage(currentAgentID, { role: 'agent', content: 'Error: ' + err.message })
  } finally {
    executing.value = false
    scrollToBottom()
  }
}

function renderMarkdown(text: string) {
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
.markdown-body :deep(p) {
  margin-top: 0.25rem;
  margin-bottom: 0.25rem;
}
.markdown-body :deep(p:first-child) { margin-top: 0; }
.markdown-body :deep(p:last-child) { margin-bottom: 0; }
</style>
