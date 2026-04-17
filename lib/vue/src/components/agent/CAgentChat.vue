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
    <CChatMessages
      v-else
      ref="chatMessagesRef"
      :messages="filteredMessages"
      :executing="executing"
      :disabled="!agentStore.activeAgentID"
      :placeholder="$t('agent.sidebar.placeholder')"
      :thinking-label="$t('agent.sidebar.thinking')"
      @send="onSend"
    >
      <template #empty>{{ $t('agent.sidebar.empty') }}</template>
    </CChatMessages>
  </div>
</template>

<script setup lang="ts">
import { computed, inject, nextTick, onMounted, ref, watch } from 'vue'
import { useAgentSidebarStore } from '../../stores/useAgentSidebarStore'
import CChatMessages from './CChatMessages.vue'

defineProps({
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
  } catch {
    // silent
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

const executing = ref(false)
const chatMessagesRef = ref<InstanceType<typeof CChatMessages> | null>(null)

watch(activeConversationIndex, () => {
  nextTick(() => chatMessagesRef.value?.scrollToBottom())
})

async function onSend(input: string) {
  if (!agentStore.activeAgentID) return

  const currentAgentID = agentStore.activeAgentID
  agentStore.addMessage(currentAgentID, { role: 'user', content: input })

  executing.value = true

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
    agentStore.addMessage(currentAgentID, { role: 'agent', content: 'Error: ' + err.message })
  } finally {
    executing.value = false
  }
}
</script>
