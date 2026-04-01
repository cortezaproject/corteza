<template>
  <Transition
    enter-active-class="transition-transform duration-300 ease-in-out"
    enter-from-class="translate-x-full"
    enter-to-class="translate-x-0"
    leave-active-class="transition-transform duration-300 ease-in-out"
    leave-from-class="translate-x-0"
    leave-to-class="translate-x-full"
  >
    <div
      v-if="isVisible"
      class="right-sidebar flex flex-row"
      :style="{ width: `${drawerWidth}px` }"
    >
      <div
        class="resize-handle w-1 h-full cursor-ew-resize hover:bg-primary/20 transition-colors shrink-0"
        @mousedown="startDrawerResize"
      />
      <div class="flex-1 flex flex-col min-w-0">
        <div class="flex items-center justify-between pl-0 pt-2 pb-2 pr-1">
          <div class="flex-1 mr-2 min-w-0">
            <Select
              v-if="agentStore.availableAgents.length > 0"
              v-model="agentStore.activeAgentID"
              :options="agentStore.availableAgents"
              optionLabel="meta.short"
              optionValue="agentID"
              class="!border-0 !shadow-none !bg-transparent"
            >
              <template #value="slotProps">
                <div v-if="slotProps.value" class="flex items-center min-w-0">
                  <i class="pi pi-sparkles text-primary mr-2" />
                  <span class="font-semibold truncate">
                    {{
                      agentStore.availableAgents.find(a => a.agentID === slotProps.value)?.meta
                        ?.short ||
                      agentStore.availableAgents.find(a => a.agentID === slotProps.value)?.handle
                    }}
                  </span>
                </div>
              </template>
              <template #option="slotProps">
                <span>{{ slotProps.option.meta?.short || slotProps.option.handle }}</span>
              </template>
            </Select>
          </div>

          <Button
            icon="pi pi-times"
            severity="secondary"
            variant="text"
            rounded
            size="small"
            @click="isVisible = false"
          />
        </div>

        <!-- Conversations Tabs Bar -->
        <div
          v-if="conversations && conversations.length > 0"
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
              @click="agentStore.setActiveConversationIndex(agentStore.activeAgentID!, idx)"
            >
              <span class="whitespace-nowrap truncate">
                {{ $t('agent.sidebar.chatTab', { id: idx + 1 }) }}
              </span>
              <i
                class="pi pi-times text-sm opacity-0 group-hover:opacity-100 p-1 hover:bg-surface rounded-full transition-all shrink-0"
                @click.stop="agentStore.closeConversation(agentStore.activeAgentID!, idx)"
              />
            </button>
          </div>
          <div class="px-1 border-l border-surface shrink-0">
            <Button
              icon="pi pi-plus"
              severity="secondary"
              variant="text"
              rounded
              size="small"
              class="!w-7 !h-7"
              v-tooltip.bottom="$t('agent.sidebar.newChat')"
              @click="agentStore.startNewConversation(agentStore.activeAgentID!)"
            />
          </div>
        </div>

        <div class="flex min-h-0 flex-1 flex-col">
          <div class="flex-1 overflow-y-auto p-3 space-y-4 flex flex-col" ref="chatContainer">
            <!-- Chat messages -->
            <template v-if="activeConversation && activeConversation.messages.length > 0">
              <div
                v-for="(msg, index) in activeConversation.messages"
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
                  <!-- Render markdown or raw text -->
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
              {{ fallbackEmptyText }}
            </div>

            <div v-if="executing" class="self-start items-start mt-auto">
              <div
                class="px-3 py-2 rounded-2xl shadow-sm text-sm bg-surface text-color border border-surface-border flex gap-1 items-center h-[34px]"
              >
                <span class="w-1.5 h-1.5 bg-muted-color rounded-full animate-bounce"></span>
                <span
                  class="w-1.5 h-1.5 bg-muted-color rounded-full animate-bounce"
                  style="animation-delay: 0.2s"
                ></span>
                <span
                  class="w-1.5 h-1.5 bg-muted-color rounded-full animate-bounce"
                  style="animation-delay: 0.4s"
                ></span>
              </div>
            </div>
          </div>

          <!-- Chat input area -->
          <div class="p-3 border-t border-surface flex items-end gap-2 shrink-0 bg-surface">
            <Textarea
              v-model="chatInput"
              autoResize
              rows="1"
              class="flex-1 text-sm block"
              style="min-height: 2.5rem; max-height: 10rem"
              :placeholder="fallbackInputPlaceholder"
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
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { computed, inject, watch, onBeforeUnmount, ref, nextTick } from 'vue'
import { useAgentSidebarStore } from '../../stores/useAgentSidebarStore'
import { useRightSidebarResize } from '../../composables/useRightSidebarResize'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const fallbackEmptyText = computed(() => t('agent.sidebar.empty'))
const fallbackInputPlaceholder = computed(() => t('agent.sidebar.placeholder'))

const agentStore = useAgentSidebarStore()
const $eventBus = inject<any>('$eventBus', null)
const $SystemAPI = inject<any>('$SystemAPI')
const { drawerWidth, startDrawerResize } = useRightSidebarResize()

const isVisible = computed({
  get: () => agentStore.visible,
  set: value => agentStore.setVisible(value),
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

const chatInput = ref('')
const executing = ref(false)
const chatContainer = ref<HTMLElement | null>(null)

// When sidebar becomes visible, emit an event so other sidebars close
watch(isVisible, visible => {
  if (visible) {
    $eventBus?.emit('right-sidebar:opened', 'agent')
    scrollToBottom()
  }
})

// Listen for other sidebars opening
const offSidebar = $eventBus?.on('right-sidebar:opened', (name: string) => {
  if (name !== 'agent') {
    agentStore.setVisible(false)
  }
})

onBeforeUnmount(() => {
  offSidebar?.()
})

const scrollToBottom = async () => {
  await nextTick()
  if (chatContainer.value) {
    chatContainer.value.scrollTop = chatContainer.value.scrollHeight
  }
}

watch(activeConversationIndex, () => {
  scrollToBottom()
})

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

    agentStore.addMessage(currentAgentID, {
      role: 'agent',
      content: outputContent,
    })
  } catch (err: any) {
    console.error('Agent execution error:', err)
    agentStore.addMessage(currentAgentID, { role: 'agent', content: 'Error: ' + err.message })
  } finally {
    executing.value = false
    scrollToBottom()
  }
}

// Simple markdown rendering
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
.markdown-body :deep(p:first-child) {
  margin-top: 0;
}
.markdown-body :deep(p:last-child) {
  margin-bottom: 0;
}
</style>
