<template>
  <div class="h-full w-full flex flex-row overflow-hidden">
    <!-- Chat area -->
    <div class="flex-1 flex flex-col border-r border-surface min-w-0">
      <slot name="header"></slot>
      <CChatMessages
        ref="chatMessagesRef"
        :messages="visibleMessages"
        :executing="executing"
        :readonly="readonly"
        :disabled="isCreate"
        :placeholder="$t('agent.editor.chat.placeholder')"
        :thinking-label="$t('agent.editor.chat.thinking')"
        :selected-trace-index="selectedTraceIndex"
        :selected-trace-type="selectedTraceType"
        @send="sendChatMessage"
        @trace-select="selectTrace"
      >
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
                <div class="text-sm text-color font-semibold">
                  {{ $t('agent.editor.chat.approval.title') }}
                </div>
                <div class="text-sm text-muted-color">
                  {{ $t('agent.editor.chat.approval.body', { tool: pending.label }) }}
                </div>
                <div v-if="pending.risk === 'destructive'" class="text-sm text-red-500 mt-1">
                  {{ $t('agent.editor.chat.approval.destructive') }}
                </div>
              </div>
            </div>
            <div class="flex gap-2 justify-end">
              <Button
                :label="$t('agent.editor.chat.approval.deny')"
                size="small"
                severity="secondary"
                text
                @click="onDeny"
              />
              <Button
                :label="$t('agent.editor.chat.approval.allow')"
                size="small"
                severity="secondary"
                outlined
                @click="onApprove(false)"
              />
              <Button
                :label="$t('agent.editor.chat.approval.allowChat')"
                size="small"
                @click="onApprove(true)"
              />
            </div>
          </div>
        </template>
      </CChatMessages>
    </div>

    <!-- Structured trace panel -->
    <div class="w-1/3 flex flex-col bg-surface border-l border-surface min-w-0" v-if="showTrace">
      <AiTrace
        ref="aiTraceRef"
        :agent="agent"
        :conversation="conversation"
        @select="onTraceSelect"
      />
    </div>
  </div>
</template>

<script setup>
/* eslint-disable vue/no-mutating-props */
import { ref, inject, nextTick, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { components } from '@planetcrust/human-vue'
import AiTrace from './AiTrace.vue'
const { CChatMessages } = components

const props = defineProps({
  agent: {
    type: Object,
    required: true,
  },
  conversation: {
    type: Object,
    required: true,
  },
  isCreate: {
    type: Boolean,
    default: false,
  },
  showTrace: {
    type: Boolean,
    default: true,
  },
  readonly: {
    type: Boolean,
    default: false,
  },
})

const visibleMessages = computed(() => {
  return (props.conversation?.messages || []).filter(
    m => ['user', 'agent', 'assistant'].includes(m.role) && m.content,
  )
})

const { t } = useI18n()
const $SystemAPI = inject('$SystemAPI')
const $toast = inject('$toast')

const executing = ref(false)

// What the agent stopped to ask about, if anything.
const pending = ref(null)

// Tools the user has approved, per conversation. An approval is a UX memory,
// not a control: the server asks again on a conversation it has not been told
// about, and nothing here can grant what the invoking user could not do anyway.
const approvedTools = ref({})

function awaitingApproval(res) {
  return res?.status === 'awaiting_approval'
}

function approvalKey() {
  return props.conversation?.conversationID || 'new'
}

function setPending(res) {
  pending.value =
    awaitingApproval(res) && res?.pendingApproval?.tool
      ? {
          tool: res.pendingApproval.tool,
          label: res.pendingApproval.title || res.pendingApproval.tool,
          risk: res.pendingApproval.risk,
        }
      : null
}

const selectedTraceIndex = ref(null)
const selectedTraceType = ref(null)
const chatMessagesRef = ref(null)
const aiTraceRef = ref(null)

async function sendChatMessage(input) {
  if (!input || executing.value) return

  props.conversation.messages.push({
    role: 'user',
    content: input,
    traceIndex: props.conversation.traceHistory.length,
  })
  executing.value = true

  try {
    const activeConvId =
      props.conversation.conversationID || props.conversation.aiConversationID || null
    const res = await $SystemAPI.agentExec({
      agentID: props.agent.agentID,
      input: input,
      ...(activeConvId ? { conversationID: activeConvId } : {}),
    })

    if (res?.conversationID) {
      props.conversation.conversationID = res.conversationID
      props.conversation.aiConversationID = res.conversationID || res.conversationID
    }

    if (res?.context) {
      props.conversation.context = res.context
    }

    // A run that stopped to ask usually has nothing to say yet, and the
    // fallback below would print the whole response object as the agent's
    // answer. Only fall back when the run actually finished.
    const content = awaitingApproval(res)
      ? res?.output
      : res?.output || (typeof res === 'string' ? res : res?.response?.text || JSON.stringify(res))

    if (content) {
      props.conversation.messages.push({
        role: 'agent',
        content,
        usage: res?.usage || null,
        traceIndex: props.conversation.traceHistory.length,
      })
    }

    setPending(res)

    props.conversation.traceHistory.push({
      prompt: input,
      decisions: res?.decisions || [],
      toolCalls: res?.toolCalls || [],
      usage: res?.usage || null,
      conversationTokens: res?.conversationTokens || 0,
    })

    selectedTraceIndex.value = props.conversation.traceHistory.length - 1
    selectedTraceType.value = 'response'
    nextTick(() => {
      aiTraceRef.value?.highlightLatest?.()
    })
  } catch (err) {
    console.error(err)
    $toast.toastErrorHandler(t('notification.agent.execFailed'))(err)
    props.conversation.messages.push({ role: 'agent', content: 'Error: ' + err.message })
    props.conversation.traceHistory.push({ error: err.message })
  } finally {
    executing.value = false
  }
}

async function onApprove(forChat) {
  const p = pending.value
  if (!p) return

  if (forChat) {
    const key = approvalKey()
    approvedTools.value = {
      ...approvedTools.value,
      [key]: [...(approvedTools.value[key] || []), p.tool],
    }
  }

  pending.value = null
  // A one-off approval is sent with the call and not remembered, so the next
  // use of the same tool asks again.
  await resumeAfterApproval(forChat ? [] : [p.tool])
}

function onDeny() {
  pending.value = null
  props.conversation.messages.push({
    role: 'agent',
    content: t('agent.editor.chat.approval.denied'),
  })
}

// The run continues where it paused: an empty input on the same conversation,
// carrying what has been approved.
async function resumeAfterApproval(once) {
  executing.value = true

  try {
    const activeConvId =
      props.conversation.conversationID || props.conversation.aiConversationID || null
    const approved = [...(approvedTools.value[approvalKey()] || []), ...once]

    const res = await $SystemAPI.agentExec({
      agentID: props.agent.agentID,
      input: '',
      ...(activeConvId ? { conversationID: activeConvId } : {}),
      approvedTools: approved,
    })

    if (res?.output) {
      props.conversation.messages.push({
        role: 'agent',
        content: res.output,
        usage: res?.usage || null,
        traceIndex: props.conversation.traceHistory.length,
      })
    }

    props.conversation.traceHistory.push({
      prompt: '',
      decisions: res?.decisions || [],
      toolCalls: res?.toolCalls || [],
      usage: res?.usage || null,
      conversationTokens: res?.conversationTokens || 0,
    })

    setPending(res)

    selectedTraceIndex.value = props.conversation.traceHistory.length - 1
    selectedTraceType.value = 'response'
    nextTick(() => {
      aiTraceRef.value?.highlightLatest?.()
    })
  } catch (err) {
    console.error(err)
    $toast.toastErrorHandler(t('notification.agent.execFailed'))(err)
    props.conversation.messages.push({ role: 'agent', content: 'Error: ' + err.message })
  } finally {
    executing.value = false
  }
}

function onTraceSelect(idx, type) {
  selectedTraceIndex.value = idx
  selectedTraceType.value = type
}

function selectTrace(traceIndex, type = 'response') {
  aiTraceRef.value?.selectExternal(traceIndex, type)
}
</script>
