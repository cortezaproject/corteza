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
import { components, useAgentTurn } from '@planetcrust/human-vue'
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

// The turn itself is shared with the sidebar and the page block; what this
// chat does with it — a conversation held as a prop, and a trace beside it — is
// its own.
const { executing, pending, send, approve, deny } = useAgentTurn({
  agentID: () => props.agent?.agentID,
  conversationID: () =>
    props.conversation.conversationID || props.conversation.aiConversationID || null,
  exec: req => $SystemAPI.agentExec(req),

  onReply: content =>
    props.conversation.messages.push({
      role: 'agent',
      content,
      traceIndex: props.conversation.traceHistory.length,
    }),

  onResponse: res => {
    if (res?.conversationID) {
      props.conversation.conversationID = res.conversationID
      props.conversation.aiConversationID = res.conversationID
    }
    if (res?.context) {
      props.conversation.context = res.context
    }

    props.conversation.traceHistory.push({
      prompt: lastPrompt.value,
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
  },

  onError: err => {
    console.error(err)
    $toast.toastErrorHandler(t('notification.agent.execFailed'))(err)
    props.conversation.messages.push({ role: 'agent', content: 'Error: ' + err.message })
    props.conversation.traceHistory.push({ error: err.message })
  },

  deniedMessage: () => t('agent.editor.chat.approval.denied'),
})

// What the trace entry for the run in flight should be labelled with.
const lastPrompt = ref('')

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

  lastPrompt.value = input
  await send(input)
}

async function onApprove(forChat) {
  lastPrompt.value = ''
  await approve(forChat)
}

function onDeny() {
  deny()
}

function onTraceSelect(idx, type) {
  selectedTraceIndex.value = idx
  selectedTraceType.value = type
}

function selectTrace(traceIndex, type = 'response') {
  aiTraceRef.value?.selectExternal(traceIndex, type)
}
</script>
