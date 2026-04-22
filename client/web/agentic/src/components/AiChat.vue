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
      />
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
  return (props.conversation?.messages || []).filter(m => ['user', 'agent', 'assistant'].includes(m.role) && m.content)
})

const { t } = useI18n()
const $SystemAPI = inject('$SystemAPI')
const $toast = inject('$toast')

const executing = ref(false)

const selectedTraceIndex = ref(null)
const selectedTraceType = ref(null)
const chatMessagesRef = ref(null)
const aiTraceRef = ref(null)

async function sendChatMessage(input) {
  if (!input || executing.value) return

  props.conversation.messages.push({ role: 'user', content: input, traceIndex: props.conversation.traceHistory.length })
  executing.value = true

  try {
    const activeConvId = props.conversation.conversationID || props.conversation.aiConversationID || null
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

    props.conversation.messages.push({
      role: 'agent',
      content:
        res?.output || (typeof res === 'string' ? res : res?.response?.text || JSON.stringify(res)),
      usage: res?.usage || null,
      traceIndex: props.conversation.traceHistory.length,
    })

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

function onTraceSelect(idx, type) {
  selectedTraceIndex.value = idx
  selectedTraceType.value = type
}

function selectTrace(traceIndex, type = 'response') {
  aiTraceRef.value?.selectExternal(traceIndex, type)
}
</script>
