<template>
  <div class="h-full w-full flex flex-row overflow-hidden">
    <!-- Chat area -->
    <div class="flex-1 flex flex-col border-r border-surface min-w-0">
      <slot name="header"></slot>
      <div class="flex-1 p-4 overflow-y-auto flex flex-col gap-4">
        <template v-for="(msg, index) in visibleMessages" :key="index">
          <div
            v-if="['user', 'agent', 'assistant'].includes(msg.role)"
            class="flex flex-col gap-1"
            :class="msg.role === 'user' ? 'items-end' : 'items-start'"
          >
            <div
              :ref="
                el => {
                  if (msg.traceIndex !== undefined) {
                    if (msg.role === 'agent' || msg.role === 'assistant') chatMsgRefs[msg.traceIndex] = el
                    else if (msg.role === 'user') chatPromptRefs[msg.traceIndex] = el
                  }
                }
              "
              :class="[
                'p-3 xl:p-4 rounded-xl max-w-[85%] text-sm md:text-base transition-all duration-200',
                msg.role === 'user'
                  ? 'bg-primary text-primary-contrast shadow-sm whitespace-pre-wrap'
                  : 'bg-emphasis text-color shadow-sm',
                msg.role === 'user' ? 'cursor-pointer' : '',
                ['agent', 'assistant'].includes(msg.role) ? 'cursor-pointer hover:shadow-md' : '',
                selectedTraceIndex === msg.traceIndex &&
                selectedTraceType === (msg.role === 'user' ? 'prompt' : 'response')
                  ? 'ring-2 ring-primary ring-offset-1'
                  : '',
              ]"
              @click="
                msg.traceIndex !== undefined
                  ? selectTrace(
                      msg.traceIndex,
                      msg.role === 'user' ? 'prompt' : 'response',
                    )
                  : null
              "
            >
              <div
                v-if="['agent', 'assistant'].includes(msg.role)"
                class="rt-content"
                v-html="renderMarkdown(msg.content)"
              />
              <template v-else>{{ msg.content }}</template>
            </div>
          </div>
        </template>
        <div v-if="executing" class="flex items-start">
          <div
            class="bg-emphasis text-color shadow-sm p-3 rounded-xl flex items-center gap-2 text-sm"
          >
            <ProgressSpinner style="width: 16px; height: 16px" strokeWidth="4" />
            <span class="text-muted-color">
              {{ $t('agent.editor.playground.thinking') }}
            </span>
          </div>
        </div>
      </div>
      <div class="p-3 border-t border-surface flex gap-2 shrink-0 bg-surface" v-if="!readonly">
        <InputText
          v-model="chatInput"
          :placeholder="$t('agent.editor.playground.placeholder')"
          class="flex-1"
          @keyup.enter="sendChatMessage"
          :disabled="isCreate || executing"
        />
        <Button
          icon="pi pi-send"
          @click="sendChatMessage"
          :disabled="isCreate || executing || !chatInput.trim()"
          :loading="executing"
        />
      </div>
    </div>

    <!-- Structured trace panel -->
    <div class="w-1/3 flex flex-col bg-surface border-l border-surface min-w-0" v-if="showTrace">
      <div class="flex-1 overflow-y-auto" ref="traceScrollContainer">
        <!-- Empty state -->
        <div
          v-if="!conversation.traceHistory.length"
          class="p-4 text-sm text-muted-color text-center"
        >
          {{ $t('agent.editor.playground.emptyTrace') }}
        </div>

        <!-- Context section -->
        <div v-if="conversation.context" class="p-3 pb-0">
          <Panel
            toggleable
            :collapsed="!contextExpanded"
            @toggle="contextExpanded = !contextExpanded"
          >
            <template #header>
              <div class="flex items-center gap-2">
                <i class="pi pi-book text-muted-color" />
                {{ $t('agent.editor.playground.traceContext') }}
              </div>
            </template>
            <pre
              class="text-xs font-mono bg-surface-ground rounded px-2.5 py-2 overflow-x-auto max-h-60 whitespace-pre-wrap break-all text-color"
              >{{ conversation.context }}</pre
            >
          </Panel>
        </div>

        <!-- Card-based trace view -->
        <div
          v-if="conversation.traceHistory.length"
          class="p-3 flex flex-col gap-4"
        >
          <div
            v-for="(traceEntry, tIdx) in conversation.traceHistory"
            :key="'t' + tIdx"
            :ref="
              el => {
                traceCardRefs[tIdx] = el
              }
            "
            class="transition-all duration-200 cursor-pointer hover:shadow-md"
            @click="selectMessage(tIdx)"
          >
            <!-- Agent Response Panel -->
            <Panel
              toggleable
              :class="
                selectedTraceIndex === tIdx && selectedTraceType === 'response'
                  ? 'border-primary shadow-md'
                  : ''
              "
              @click.stop="selectMessage(tIdx, 'response')"
            >
              <template #header>
                <div class="flex items-center gap-2">
                  <i
                    class="pi pi-sparkles"
                    :class="
                      selectedTraceIndex === tIdx && selectedTraceType === 'response'
                        ? 'text-primary'
                        : 'text-muted-color'
                    "
                  />
                  {{ $t('agent.editor.playground.traceAgentResponse') }}
                </div>
              </template>

              <div class="flex flex-col gap-3">
                <!-- Error state -->
                <div
                  v-if="traceEntry.error"
                  class="flex items-start gap-2 text-xs text-red-600 bg-red-50 rounded-md px-2.5 py-2"
                >
                  <i class="pi pi-exclamation-triangle shrink-0 mt-0.5 text-xs" />
                  <span class="break-all">{{ traceEntry.error }}</span>
                </div>

                <!-- Decision rows -->
                <div
                  v-for="(decision, dIdx) in traceEntry.decisions || []"
                  :key="'d' + tIdx + '-' + dIdx"
                  class="border border-surface rounded-md overflow-hidden"
                >
                  <!-- Decision header -->
                  <div class="flex items-center gap-2 px-2.5 py-1.5 bg-surface-ground">
                    <!-- Step icon -->
                    <div
                      class="w-5 h-5 rounded-full flex items-center justify-center shrink-0"
                      :class="
                        decision.decision === 'tool_call'
                          ? 'bg-blue-100 text-blue-600'
                          : 'bg-green-100 text-green-600'
                      "
                    >
                      <i
                        :class="
                          decision.decision === 'tool_call'
                            ? 'pi pi-wrench'
                            : 'pi pi-comment'
                        "
                        class="text-xs"
                      />
                    </div>

                    <span class="text-xs font-medium text-color">
                      {{
                        decision.decision === 'tool_call'
                          ? $t('agent.editor.playground.traceToolCall')
                          : $t('agent.editor.playground.traceResponse')
                      }}
                    </span>

                    <span
                      v-if="decision.usage?.contextWindow"
                      class="ml-auto text-xs text-muted-color"
                    >
                      {{ decision.usage.contextWindow }}
                      {{ $t('agent.editor.playground.traceTokens') }}
                    </span>
                  </div>

                  <!-- Reasoning (expandable) -->
                  <div
                    v-if="decision.reasoning"
                    class="border-t border-surface px-2.5 py-1.5"
                  >
                    <div
                      class="flex items-start gap-1.5 cursor-pointer"
                      @click.stop="toggleReasoningExpand(tIdx, dIdx)"
                    >
                      <i
                        class="pi text-xs text-muted-color mt-0.5 transition-transform duration-200"
                        :class="
                          isReasoningExpanded(tIdx, dIdx)
                            ? 'pi-chevron-down'
                            : 'pi-chevron-right'
                        "
                      />
                      <p
                        class="text-xs text-muted-color italic flex-1"
                        :class="isReasoningExpanded(tIdx, dIdx) ? '' : 'line-clamp-1'"
                      >
                        {{ decision.reasoning }}
                      </p>
                    </div>
                  </div>

                  <!-- Tool calls detail -->
                  <div v-if="decision.tools?.length" class="flex flex-col">
                    <div
                      v-for="(toolName, toolIdx) in decision.tools"
                      :key="'tool-' + tIdx + '-' + dIdx + '-' + toolIdx"
                      class="border-t border-surface"
                    >
                      <!-- Tool header (clickable) -->
                      <div
                        class="flex items-center justify-between px-2.5 py-1.5 cursor-pointer hover:bg-surface-ground/50 transition-colors"
                        @click.stop="toggleToolExpand(tIdx, dIdx, toolIdx)"
                      >
                        <div class="flex items-center gap-1.5 min-w-0">
                          <i
                            class="pi text-xs text-muted-color transition-transform duration-200"
                            :class="
                              isToolExpanded(tIdx, dIdx, toolIdx)
                                ? 'pi-chevron-down'
                                : 'pi-chevron-right'
                            "
                          />
                          <span class="font-mono text-xs text-color truncate">
                            {{ toolName }}
                          </span>
                        </div>
                        <div class="flex items-center gap-2 shrink-0 text-xs">
                          <span
                            v-if="getToolCallError(traceEntry, toolName)"
                            class="text-red-500"
                          >
                            <i class="pi pi-exclamation-triangle text-xs" />
                          </span>
                          <span
                            v-if="getToolCallDuration(traceEntry, toolName)"
                            class="text-muted-color"
                          >
                            {{ getToolCallDuration(traceEntry, toolName) }}ms
                          </span>
                        </div>
                      </div>

                      <!-- Expanded tool detail -->
                      <div
                        v-if="isToolExpanded(tIdx, dIdx, toolIdx)"
                        class="border-t border-surface bg-surface-ground/30 px-2.5 py-2 flex flex-col gap-2"
                        @click.stop
                      >
                        <!-- Tool error (inline) -->
                        <div
                          v-if="getToolCallError(traceEntry, toolName)"
                          class="text-xs text-red-600 bg-red-50 rounded px-2 py-1.5"
                        >
                          <i class="pi pi-exclamation-triangle mr-1 text-xs" />
                          {{ getToolCallError(traceEntry, toolName) }}
                        </div>

                        <!-- Args -->
                        <div v-if="getToolCallArgs(traceEntry, toolName)">
                          <div class="flex items-center justify-between mb-1">
                            <span class="text-xs font-medium text-muted-color">
                              {{ $t('agent.editor.playground.traceArgs') }}
                            </span>
                            <button
                              class="text-xs text-muted-color hover:text-color transition-colors p-0.5"
                              @click.stop="
                                copyToClipboard(
                                  JSON.stringify(
                                    getToolCallArgs(traceEntry, toolName),
                                    null,
                                    2,
                                  ),
                                )
                              "
                              v-tooltip.left="$t('agent.editor.playground.traceCopy')"
                            >
                              <i class="pi pi-copy text-xs" />
                            </button>
                          </div>
                          <pre
                            class="text-xs font-mono bg-surface-ground rounded px-2 py-1.5 overflow-x-auto max-h-40 whitespace-pre-wrap break-all text-color"
                            >{{
                              JSON.stringify(
                                getToolCallArgs(traceEntry, toolName),
                                null,
                                2,
                              )
                            }}</pre
                          >
                        </div>

                        <!-- Result -->
                        <div v-if="getToolCallResult(traceEntry, toolName) !== null">
                          <div class="flex items-center justify-between mb-1">
                            <span class="text-xs font-medium text-muted-color">
                              {{ $t('agent.editor.playground.traceResult') }}
                            </span>
                            <button
                              class="text-xs text-muted-color hover:text-color transition-colors p-0.5"
                              @click.stop="
                                copyToClipboard(
                                  JSON.stringify(
                                    getToolCallResult(traceEntry, toolName),
                                    null,
                                    2,
                                  ),
                                )
                              "
                              v-tooltip.left="$t('agent.editor.playground.traceCopy')"
                            >
                              <i class="pi pi-copy text-xs" />
                            </button>
                          </div>
                          <pre
                            class="text-xs font-mono bg-surface-ground rounded px-2 py-1.5 overflow-x-auto max-h-40 whitespace-pre-wrap break-all text-color"
                            >{{
                              JSON.stringify(
                                getToolCallResult(traceEntry, toolName),
                                null,
                                2,
                              )
                            }}</pre
                          >
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Token footer -->
                <div
                  v-if="traceEntry.usage?.contextWindow"
                  class="flex items-center justify-end text-xs text-muted-color mt-1"
                >

                  <div class="flex items-center gap-1">
                    <i class="pi pi-chart-bar text-xs" />
                    <span class="font-medium text-color">
                      {{ traceEntry.usage.contextWindow }}
                      {{ $t('agent.editor.playground.traceContextTokens') }}
                    </span>
                    <template v-if="agent.execution.limits.contextWindow > 0">
                      <span class="text-muted-color">
                        ({{
                          Math.round(
                            (traceEntry.usage.contextWindow /
                              agent.execution.limits.contextWindow) *
                              100,
                          )
                        }}%)
                      </span>
                    </template>
                  </div>
                </div>
              </div>
            </Panel>
          </div>
        </div>
      </div>

      <!-- Cumulative total footer -->
      <div
        v-if="conversation.traceHistory.length"
        class="px-3 py-2 flex items-center justify-end text-xs text-muted-color border-t border-surface shrink-0"
      >
        <div class="flex items-center gap-1">
          <i class="pi pi-chart-bar text-xs" />
          <span class="font-medium text-color">
            {{
              conversation.traceHistory[
                conversation.traceHistory.length - 1
              ]?.usage?.contextWindow || 0
            }}
            {{ $t('agent.editor.playground.traceContextTokens') }}
          </span>
          <template v-if="agent.execution.limits.contextWindow > 0">
            <span class="text-muted-color">
              ({{
                Math.round(
                  ((conversation.traceHistory[
                    conversation.traceHistory.length - 1
                  ]?.usage?.contextWindow || 0) /
                    agent.execution.limits.contextWindow) *
                    100,
                )
              }}%)
            </span>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
/* eslint-disable vue/no-mutating-props */
import { ref, inject, nextTick, computed } from 'vue'
import { useI18n } from 'vue-i18n'

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
const chatInput = ref('')
const contextExpanded = ref(false)

const selectedTraceIndex = ref(null)
const selectedTraceType = ref(null)
const chatMsgRefs = ref({})
const chatPromptRefs = ref({})
const traceCardRefs = ref({})
const traceScrollContainer = ref(null)

const expandedTools = ref(new Set())
const expandedReasoning = ref(new Set())

async function sendChatMessage() {
  if (!chatInput.value.trim() || executing.value) return

  const input = chatInput.value
  chatInput.value = ''

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

    // Append structured trace
    props.conversation.traceHistory.push({
      prompt: input,
      decisions: res?.decisions || [],
      toolCalls: res?.toolCalls || [],
      usage: res?.usage || null,
      conversationTokens: res?.conversationTokens || 0,
    })

    // Auto-select the new trace card
    selectedTraceIndex.value = props.conversation.traceHistory.length - 1
  } catch (err) {
    console.error(err)
    $toast.toastDanger(t('notification.agent.execFailed'))
    props.conversation.messages.push({ role: 'agent', content: 'Error: ' + err.message })
    props.conversation.traceHistory.push({ error: err.message })
  } finally {
    executing.value = false
  }
}

function getToolCallDuration(traceEntry, toolName) {
  if (!traceEntry?.toolCalls) return null
  const tc = traceEntry.toolCalls.find(tc => tc.tool === toolName)
  return tc?.durationMs || null
}

function getToolCallError(traceEntry, toolName) {
  if (!traceEntry?.toolCalls) return null
  const tc = traceEntry.toolCalls.find(tc => tc.tool === toolName)
  return tc?.error || null
}

function getToolCallArgs(traceEntry, toolName) {
  if (!traceEntry?.toolCalls) return null
  const tc = traceEntry.toolCalls.find(tc => tc.tool === toolName)
  return tc?.args && Object.keys(tc.args).length ? tc.args : null
}

function getToolCallResult(traceEntry, toolName) {
  if (!traceEntry?.toolCalls) return null
  const tc = traceEntry.toolCalls.find(tc => tc.tool === toolName)
  return tc?.result !== undefined ? tc.result : null
}

function toolKey(tIdx, dIdx, toolIdx) {
  return `${tIdx}-${dIdx}-${toolIdx}`
}

function reasoningKey(tIdx, dIdx) {
  return `${tIdx}-${dIdx}`
}

function toggleToolExpand(tIdx, dIdx, toolIdx) {
  const key = toolKey(tIdx, dIdx, toolIdx)
  if (expandedTools.value.has(key)) {
    expandedTools.value.delete(key)
  } else {
    expandedTools.value.add(key)
  }
  expandedTools.value = new Set(expandedTools.value)
}

function isToolExpanded(tIdx, dIdx, toolIdx) {
  return expandedTools.value.has(toolKey(tIdx, dIdx, toolIdx))
}

function toggleReasoningExpand(tIdx, dIdx) {
  const key = reasoningKey(tIdx, dIdx)
  if (expandedReasoning.value.has(key)) {
    expandedReasoning.value.delete(key)
  } else {
    expandedReasoning.value.add(key)
  }
  expandedReasoning.value = new Set(expandedReasoning.value)
}

function isReasoningExpanded(tIdx, dIdx) {
  return expandedReasoning.value.has(reasoningKey(tIdx, dIdx))
}

function copyToClipboard(text) {
  navigator.clipboard.writeText(text)
}

function selectMessage(traceIndex, type = 'response') {
  if (selectedTraceIndex.value === traceIndex && selectedTraceType.value === type) {
    selectedTraceIndex.value = null
    selectedTraceType.value = null
  } else {
    selectedTraceIndex.value = traceIndex
    selectedTraceType.value = type
    nextTick(() => {
      const refs = type === 'prompt' ? chatPromptRefs : chatMsgRefs
      const msgEl = refs.value[traceIndex]
      if (msgEl) msgEl.scrollIntoView({ behavior: 'smooth', block: 'nearest' })
    })
  }
}

function selectTrace(traceIndex, type = 'response') {
  if (selectedTraceIndex.value === traceIndex && selectedTraceType.value === type) {
    selectedTraceIndex.value = null
    selectedTraceType.value = null
  } else {
    selectedTraceIndex.value = traceIndex
    selectedTraceType.value = type
    nextTick(() => {
      const card = traceCardRefs.value[traceIndex]
      if (card) card.scrollIntoView({ behavior: 'smooth', block: 'nearest' })
    })
  }
}

// Simple markdown rendering
function renderMarkdown(text) {
  if (!text) return ''

  let html = text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')

  // Code blocks
  html = html.replace(/```([\s\S]*?)```/g, (match, p1) => {
    return `<pre class="bg-surface-ground p-3 rounded-lg overflow-x-auto text-sm font-mono my-2 text-color border border-surface shadow-sm"><code>${p1.trim()}</code></pre>`
  })

  // Inline code
  html = html.replace(/`([^`]+)`/g, '<code class="bg-surface-ground px-1.5 py-0.5 rounded text-sm font-mono text-primary border border-surface shadow-sm">$1</code>')

  // Bold
  html = html.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')

  // Italic
  html = html.replace(/\*([^*]+)\*/g, '<em>$1</em>')

  // Headers
  html = html.replace(/^### (.*$)/gim, '<h3 class="text-base font-bold mt-3 mb-1 text-primary-emphasis">$1</h3>')
  html = html.replace(/^## (.*$)/gim, '<h2 class="text-lg font-bold mt-4 mb-2 pb-1 border-b border-surface text-primary-emphasis">$1</h2>')
  html = html.replace(/^# (.*$)/gim, '<h1 class="text-xl font-bold mt-5 mb-3 pb-1 border-b border-surface text-primary-emphasis">$1</h1>')

  // Lists
  html = html.replace(/^\s*[-*+]\s+(.*)$/gim, '<li class="ml-4 list-disc marker:text-primary">$1</li>')
  html = html.replace(/(<li.*<\/li>)\n(<li.*<\/li>)/g, '$1$2')
  html = html.replace(/(<li.*<\/li>)/g, '<ul class="my-2 space-y-1 text-color">$1</ul>')
  html = html.replace(/<\/ul>\n<ul[^>]*>/g, '')

  // Line breaks
  html = html.replace(/\n\n/g, '</p><p class="my-2 text-color">')
  html = html.replace(/\n/g, '<br>')

  return `<div class="prose-sm max-w-none break-words">${html}</div>`
}
</script>
