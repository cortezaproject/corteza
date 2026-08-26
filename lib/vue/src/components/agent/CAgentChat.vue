<template>
  <div class="flex flex-col min-h-0 flex-1">
    <!-- Header -->
    <div class="flex items-center gap-2 px-3 py-2 border-b border-surface shrink-0">
      <i class="pi pi-sparkles text-primary text-sm shrink-0" />
      <div class="flex-1 min-w-0">
        <Select
          v-if="visibleAgents.length > 0"
          v-model="agentStore.activeAgentID"
          :options="visibleAgents"
          optionLabel="meta.short"
          optionValue="agentID"
          size="small"
          class="!border-0 !shadow-none !bg-transparent"
        >
          <template #value="slotProps">
            <span v-if="slotProps.value" class="font-semibold text-base text-color truncate">
              {{
                visibleAgents.find(a => a.agentID === slotProps.value)?.meta?.short ||
                visibleAgents.find(a => a.agentID === slotProps.value)?.handle
              }}
            </span>
          </template>
          <template #option="slotProps">
            <span>{{ slotProps.option.meta?.short || slotProps.option.handle }}</span>
          </template>
        </Select>
        <span v-else class="font-semibold text-base text-color truncate block">
          {{ translations.noAgents }}
        </span>
      </div>
      <slot name="actions" />
    </div>

    <!-- Conversation tabs + actions. The tab row is visible whenever an agent
         is selected so the user can switch / browse history before sending
         their first message. Empty tabs render as "New chat"; tabs with
         messages use a truncation of the first user message. -->
    <CConversationTabs
      v-if="agentStore.activeAgentID && visibleAgents.length > 0"
      :tabs="conversationTabs"
      :active-index="activeConversationIndex"
      @update:active-index="
        idx => agentStore.setActiveConversationIndex(agentStore.activeAgentID, idx)
      "
      @close="idx => agentStore.closeConversation(agentStore.activeAgentID, idx)"
    >
      <template #actions>
        <Button
          icon="pi pi-plus"
          severity="secondary"
          variant="text"
          rounded
          size="small"
          class="!w-7 !h-7"
          v-tooltip.bottom="{ value: translations.newChat, showDelay: 500 }"
          @click="agentStore.startNewConversation(agentStore.activeAgentID)"
        />
        <Divider layout="vertical" class="!mx-1 !my-0 !h-5" />
        <Button
          icon="pi pi-history"
          severity="secondary"
          variant="text"
          rounded
          size="small"
          class="!w-7 !h-7"
          v-tooltip.bottom="{ value: translations.history.button, showDelay: 500 }"
          @click="openHistory"
        />
      </template>
    </CConversationTabs>

    <!-- History popover -->
    <Popover ref="historyPopover" :pt="{ content: { class: 'max-h-96 overflow-y-auto' } }">
      <div class="w-72 flex flex-col gap-2">
        <div
          v-if="historyLoading"
          class="flex items-center justify-center gap-2 p-4 text-sm text-muted-color"
        >
          <ProgressSpinner style="width: 14px; height: 14px" strokeWidth="4" />
          <span>{{ translations.history.loading }}</span>
        </div>
        <div
          v-else-if="historyEntries.length === 0"
          class="p-4 text-sm text-muted-color text-center"
        >
          {{ translations.history.empty }}
        </div>
        <button
          v-for="conv in historyEntries"
          v-else
          :key="conv.aiConversationID"
          class="group w-full flex items-start gap-2 pl-3 pr-2 py-2 text-left border border-surface rounded-border hover:bg-emphasis transition-colors"
          @click="onPickHistory(conv)"
        >
          <div class="flex-1 min-w-0">
            <div class="text-sm text-color truncate">{{ historyLabel(conv) }}</div>
            <div class="text-xs text-muted-color">{{ formatHistoryDate(conv) }}</div>
          </div>
          <i
            class="pi pi-trash text-sm opacity-0 group-hover:opacity-100 p-1 hover:text-red-500 rounded-full transition-all shrink-0"
            @click.stop="onDeleteHistory(conv)"
          />
        </button>
      </div>
    </Popover>

    <!-- No agents state -->
    <div
      v-if="visibleAgents.length === 0"
      class="flex-1 flex flex-col items-center justify-center gap-2 p-6 text-center"
    >
      <i class="pi pi-sparkles text-3xl text-muted-color" />
      <span class="text-sm text-muted-color">{{ translations.noAgents }}</span>
    </div>

    <!-- Chat -->
    <CChatMessages
      v-else
      ref="chatMessagesRef"
      :messages="filteredMessages"
      :executing="executing"
      :disabled="!agentStore.activeAgentID"
      :placeholder="translations.placeholder"
      :thinking-label="translations.thinking"
      @send="onSend"
    >
      <template #empty>{{ translations.empty }}</template>

      <!-- What the agent stopped to ask about. It sits above the composer
           rather than in the message list: it is a decision to make, not a
           turn that happened. -->
      <template v-if="pending" #beforeComposer>
        <div
          class="mx-4 mb-2 rounded-border border border-surface bg-emphasis p-3 flex flex-col gap-2"
          data-testid="agent-approval"
        >
          <div class="flex items-start gap-2">
            <i
              class="pi pi-shield text-sm mt-0.5 shrink-0"
              :class="pending.risk === 'destructive' ? 'text-red-500' : 'text-primary'"
            />
            <div class="min-w-0">
              <div class="text-sm text-color font-semibold">{{ translations.approval.title }}</div>
              <div class="text-sm text-muted-color">
                {{ translations.approval.body(pending.tool) }}
              </div>
              <div v-if="pending.risk === 'destructive'" class="text-sm text-red-500 mt-1">
                {{ translations.approval.destructive }}
              </div>
            </div>
          </div>
          <div class="flex gap-2 justify-end">
            <Button
              :label="translations.approval.deny"
              size="small"
              severity="secondary"
              text
              @click="onDeny"
            />
            <Button
              :label="translations.approval.allow"
              size="small"
              severity="secondary"
              outlined
              @click="onApprove(false)"
            />
            <Button
              :label="translations.approval.allowChat"
              size="small"
              @click="onApprove(true)"
            />
          </div>
        </div>
      </template>
    </CChatMessages>
  </div>
</template>

<script setup lang="ts">
import { computed, inject, nextTick, onMounted, ref, watch } from 'vue'
import { useAgentChatStore } from '../../stores/useAgentChatStore'
import CChatMessages from './CChatMessages.vue'
import CConversationTabs from './CConversationTabs.vue'
import type { AgentChatTranslations } from './translations'

const props = defineProps({
  translations: {
    type: Object as () => AgentChatTranslations,
    required: true,
  },
  // Restrict the agent dropdown to this subset of agent IDs (string match).
  // Empty/undefined = no restriction (show whatever the store has).
  allowedAgentIDs: {
    type: Array as () => string[],
    default: () => [],
  },
  // Preferred default agent on first mount (overrides the store's auto-pick if
  // present in visibleAgents). Ignored if the user already has an active agent.
  defaultAgentID: {
    type: String,
    default: '',
  },
  // When true, on first mount lazily load and hydrate the most recent past
  // conversation for the active agent into the (single) tab.
  autoResume: {
    type: Boolean,
    default: false,
  },
  // Optional callable returning an object that gets passed as `context` to
  // agentExec on every send. Lets callers (e.g. the page block) attach
  // record/page identifiers without changing the chat UI.
  contextProvider: {
    type: Function as unknown as () => () => Record<string, any> | undefined,
    default: null,
  },
})

const $SystemAPI = inject<any>('$SystemAPI')
const $Auth = inject<any>('$Auth')
const agentStore = useAgentChatStore()

// When the caller (e.g. the page block) passes an explicit allowlist, we
// resolve agents directly via the API and bypass sidebarRoles entirely — the
// editor's choice in the configurator is authoritative. The shared store is
// reserved for sidebar mode (no allowlist), so it isn't polluted by
// block-only agents.
const localAgents = ref<any[]>([])

const usingAllowlist = computed(
  () => Array.isArray(props.allowedAgentIDs) && props.allowedAgentIDs.length > 0,
)

const visibleAgents = computed(() => {
  if (!usingAllowlist.value) {
    return agentStore.availableAgents
  }
  const set = new Set(props.allowedAgentIDs.map(String))
  return localAgents.value.filter((a: any) => set.has(String(a.agentID)))
})

async function fetchAvailableAgents() {
  try {
    const res = await $SystemAPI.agentList({ limit: 0 })
    const allAgents = res.set || []

    if (usingAllowlist.value) {
      // Block mode — only require user-invocable; allowlist does the rest.
      localAgents.value = allAgents.filter((agent: any) => !!agent.invocation?.user?.enabled)
      return
    }

    // Sidebar mode — apply sidebarRoles gating before publishing to the store.
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
}

onMounted(async () => {
  // Allowlist callers always fetch their own list; sidebar callers reuse the
  // store across the session.
  if (usingAllowlist.value) {
    if (localAgents.value.length === 0) await fetchAvailableAgents()
  } else if (agentStore.availableAgents.length === 0) {
    await fetchAvailableAgents()
  }

  // Reconcile activeAgentID against what the caller allows.
  if (visibleAgents.value.length > 0) {
    const currentOK =
      agentStore.activeAgentID &&
      visibleAgents.value.find(a => String(a.agentID) === String(agentStore.activeAgentID))

    if (!currentOK) {
      const preferred =
        (props.defaultAgentID &&
          visibleAgents.value.find(a => String(a.agentID) === String(props.defaultAgentID))) ||
        visibleAgents.value[0]
      if (preferred) agentStore.setActiveAgentID(preferred.agentID)
    }
  }

  if (props.autoResume && agentStore.activeAgentID) {
    await maybeResumeLatest(agentStore.activeAgentID)
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

function tabLabel(conv: any): string {
  const firstUser = (conv?.messages || []).find((m: any) => m.role === 'user' && m.content)
  if (!firstUser) return props.translations.newChat
  const text = String(firstUser.content).replace(/\s+/g, ' ').trim()
  if (text.length <= 24) return text
  return text.slice(0, 21) + '…'
}

const conversationTabs = computed(() =>
  conversations.value.map((conv: any) => ({
    label: tabLabel(conv),
  })),
)

const executing = ref(false)
const chatMessagesRef = ref<InstanceType<typeof CChatMessages> | null>(null)
const historyPopover = ref<any>(null)

const historyEntries = computed(() => {
  if (!agentStore.activeAgentID) return []
  return agentStore.getHistory(agentStore.activeAgentID)
})

const historyLoading = computed(() => {
  if (!agentStore.activeAgentID) return false
  return agentStore.isHistoryLoading(agentStore.activeAgentID)
})

watch(activeConversationIndex, () => {
  nextTick(() => chatMessagesRef.value?.scrollToBottom())
})

async function openHistory(event: MouseEvent) {
  historyPopover.value?.toggle(event)
  if (agentStore.activeAgentID) {
    // Refresh on each open — cheap and avoids stale views.
    await agentStore.loadHistory(agentStore.activeAgentID)
  }
}

function onPickHistory(conv: any) {
  if (!agentStore.activeAgentID) return
  agentStore.openConversationFromHistory(agentStore.activeAgentID, conv)
  historyPopover.value?.hide()
}

async function onDeleteHistory(conv: any) {
  if (!agentStore.activeAgentID || !conv) return
  await agentStore.deleteConversation(agentStore.activeAgentID, String(conv.aiConversationID))
}

function historyLabel(conv: any): string {
  const firstUser = (conv.messages || []).find((m: any) => m.role === 'user' && m.content)
  if (!firstUser) return props.translations.history.untitled
  const text = String(firstUser.content).replace(/\s+/g, ' ').trim()
  if (text.length <= 60) return text
  return text.slice(0, 57) + '…'
}

function formatHistoryDate(conv: any): string {
  const raw = conv.updatedAt || conv.createdAt
  if (!raw) return ''
  try {
    const d = new Date(raw)
    return d.toLocaleString()
  } catch {
    return ''
  }
}

async function maybeResumeLatest(agentID: string) {
  await agentStore.loadHistory(agentID)
  const list = agentStore.getHistory(agentID)
  if (list.length === 0) return
  agentStore.openConversationFromHistory(agentID, list[0])
}

watch(
  () => agentStore.activeAgentID,
  async (id, prev) => {
    if (!id || id === prev) return
    if (props.autoResume) await maybeResumeLatest(id)
  },
)

// Tools the user has approved, per conversation. An approval is a UX memory,
// not a control: the server asks again on a conversation it has not been told
// about, and nothing here can grant what the invoking user could not do anyway.
const approvedTools = ref<Record<string, string[]>>({})

// What the agent stopped to ask about, if anything.
const pending = ref<{ tool: string; risk?: string; agentID: string } | null>(null)

function approvalKey(agentID: string, conversationID?: string) {
  return `${agentID}:${conversationID || 'new'}`
}

async function onSend(input: string) {
  if (!agentStore.activeAgentID) return

  const currentAgentID = agentStore.activeAgentID
  agentStore.addMessage(currentAgentID, { role: 'user', content: input })

  await runAgent(currentAgentID, input)
}

// runAgent covers both a fresh question and the resume that follows an
// approval; a resume carries no input of its own, only the tool it may use.
async function runAgent(currentAgentID: string, input: string) {
  executing.value = true
  pending.value = null

  try {
    const activeConv = activeConversation.value as any
    const context = props.contextProvider ? props.contextProvider() : undefined
    const key = approvalKey(currentAgentID, activeConv?.conversationID)
    const approved = approvedTools.value[key] || []

    const res = await $SystemAPI.agentExec({
      agentID: currentAgentID,
      input,
      ...(activeConv?.conversationID ? { conversationID: activeConv.conversationID } : {}),
      ...(context && Object.keys(context).length > 0 ? { context } : {}),
      ...(approved.length ? { approvedTools: approved } : {}),
    })

    if (res?.conversationID) {
      agentStore.setConversationID(currentAgentID, res.conversationID)
    }

    const outputContent =
      res?.output || (typeof res === 'string' ? res : res?.response?.text || JSON.stringify(res))

    if (outputContent) {
      agentStore.addMessage(currentAgentID, { role: 'agent', content: outputContent })
    }

    if (res?.status === 'awaiting_approval' && res?.pendingApproval?.tool) {
      pending.value = {
        tool: res.pendingApproval.tool,
        risk: res.pendingApproval.risk,
        agentID: currentAgentID,
      }
    }
  } catch (err: any) {
    agentStore.addMessage(currentAgentID, { role: 'agent', content: 'Error: ' + err.message })
  } finally {
    executing.value = false
  }
}

async function onApprove(forChat: boolean) {
  const p = pending.value
  if (!p) return

  const conv = activeConversation.value as any
  if (forChat) {
    const key = approvalKey(p.agentID, conv?.conversationID)
    approvedTools.value = {
      ...approvedTools.value,
      [key]: [...(approvedTools.value[key] || []), p.tool],
    }
  }

  pending.value = null
  // A one-off approval is sent with the call and not remembered, so the next
  // use of the same tool asks again.
  await runAgentApproving(p.agentID, forChat ? [] : [p.tool])
}

async function runAgentApproving(currentAgentID: string, once: string[]) {
  executing.value = true
  try {
    const activeConv = activeConversation.value as any
    const context = props.contextProvider ? props.contextProvider() : undefined
    const key = approvalKey(currentAgentID, activeConv?.conversationID)
    const approved = [...(approvedTools.value[key] || []), ...once]

    const res = await $SystemAPI.agentExec({
      agentID: currentAgentID,
      input: '',
      ...(activeConv?.conversationID ? { conversationID: activeConv.conversationID } : {}),
      ...(context && Object.keys(context).length > 0 ? { context } : {}),
      approvedTools: approved,
    })

    if (res?.output) {
      agentStore.addMessage(currentAgentID, { role: 'agent', content: res.output })
    }
    if (res?.status === 'awaiting_approval' && res?.pendingApproval?.tool) {
      pending.value = {
        tool: res.pendingApproval.tool,
        risk: res.pendingApproval.risk,
        agentID: currentAgentID,
      }
    }
  } catch (err: any) {
    agentStore.addMessage(currentAgentID, { role: 'agent', content: 'Error: ' + err.message })
  } finally {
    executing.value = false
  }
}

function onDeny() {
  const p = pending.value
  if (!p) return
  pending.value = null
  agentStore.addMessage(p.agentID, { role: 'agent', content: props.translations.approval.denied })
}

defineExpose({
  reloadAvailableAgents: fetchAvailableAgents,
  refreshActiveConversation: async () => {
    if (!agentStore.activeAgentID) return
    const conv = activeConversation.value as any
    if (!conv?.conversationID) return
    try {
      const fresh = await $SystemAPI.aiConversationRead({ aiConversationID: conv.conversationID })
      conv.messages = (fresh.messages || []).map((m: any) => ({
        role: m.role === 'assistant' ? 'agent' : m.role,
        content: m.content || '',
      }))
    } catch {
      // silent
    }
  },
  reloadHistory: async () => {
    if (!agentStore.activeAgentID) return
    await agentStore.loadHistory(agentStore.activeAgentID)
  },
})
</script>
