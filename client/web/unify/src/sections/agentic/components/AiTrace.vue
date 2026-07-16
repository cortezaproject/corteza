<template>
  <div class="h-full w-full flex flex-col bg-surface min-w-0">
    <div
      v-if="!conversation.traceHistory.length"
      class="flex-1 flex flex-col items-center justify-center p-6 text-sm text-muted-color text-center gap-2"
    >
      <i class="pi pi-comments text-3xl opacity-50" />
      <span>{{ $t('agent.editor.inspect.noActiveChat') }}</span>
    </div>
    <div v-else class="flex-1 overflow-y-auto" ref="traceScrollContainer">

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
              {{ $t('agent.editor.inspect.context') }}
            </div>
          </template>
          <div
            class="text-xs bg-emphasis rounded px-2.5 py-2 overflow-x-auto max-h-60 break-words text-color"
            v-html="renderMarkdown(conversation.context)"
          />
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
          <Panel
            toggleable
            :class="
              selectedIndex === tIdx && selectedType === 'response'
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
                    selectedIndex === tIdx && selectedType === 'response'
                      ? 'text-primary'
                      : 'text-muted-color'
                  "
                />
                {{ $t('agent.editor.inspect.agentResponse') }}
              </div>
            </template>

            <div class="flex flex-col gap-3">
              <div
                v-if="traceEntry.error"
                class="flex items-start gap-2 text-xs text-red-600 bg-red-50 rounded-md px-2.5 py-2"
              >
                <i class="pi pi-exclamation-triangle shrink-0 mt-0.5 text-xs" />
                <span class="break-all">{{ traceEntry.error }}</span>
              </div>

              <div
                v-for="(decision, dIdx) in traceEntry.decisions || []"
                :key="'d' + tIdx + '-' + dIdx"
                class="border border-surface rounded-md overflow-hidden"
              >
                <div class="flex items-center gap-2 px-2.5 py-1.5 bg-emphasis">
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
                        ? $t('agent.editor.inspect.toolCall')
                        : $t('agent.editor.inspect.response')
                    }}
                  </span>

                  <span
                    v-if="decision.usage?.contextWindow"
                    class="ml-auto text-xs text-muted-color"
                  >
                    {{ decision.usage.contextWindow }}
                    {{ $t('agent.editor.inspect.tokens') }}
                  </span>
                </div>

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

                <div v-if="decision.tools?.length" class="flex flex-col">
                  <div
                    v-for="(toolName, toolIdx) in decision.tools"
                    :key="'tool-' + tIdx + '-' + dIdx + '-' + toolIdx"
                    class="border-t border-surface"
                  >
                    <div
                      class="flex items-center justify-between px-2.5 py-1.5 cursor-pointer hover:bg-emphasis transition-colors"
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

                    <div
                      v-if="isToolExpanded(tIdx, dIdx, toolIdx)"
                      class="border-t border-surface bg-emphasis px-2.5 py-2 flex flex-col gap-2"
                      @click.stop
                    >
                      <div
                        v-if="getToolCallError(traceEntry, toolName)"
                        class="text-xs text-red-600 bg-red-50 rounded px-2 py-1.5"
                      >
                        <i class="pi pi-exclamation-triangle mr-1 text-xs" />
                        {{ getToolCallError(traceEntry, toolName) }}
                      </div>

                      <div v-if="getToolCallArgs(traceEntry, toolName)">
                        <div class="flex items-center justify-between mb-1">
                          <span class="text-xs font-medium text-muted-color">
                            {{ $t('agent.editor.inspect.args') }}
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
                            v-tooltip.left="$t('agent.editor.inspect.copy')"
                          >
                            <i class="pi pi-copy text-xs" />
                          </button>
                        </div>
                        <pre
                          class="text-xs font-mono bg-emphasis rounded px-2 py-1.5 overflow-x-auto max-h-40 whitespace-pre-wrap break-all text-color"
                          >{{
                            JSON.stringify(
                              getToolCallArgs(traceEntry, toolName),
                              null,
                              2,
                            )
                          }}</pre
                        >
                      </div>

                      <div v-if="getToolCallResult(traceEntry, toolName) !== null">
                        <div class="flex items-center justify-between mb-1">
                          <span class="text-xs font-medium text-muted-color">
                            {{ $t('agent.editor.inspect.result') }}
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
                            v-tooltip.left="$t('agent.editor.inspect.copy')"
                          >
                            <i class="pi pi-copy text-xs" />
                          </button>
                        </div>
                        <pre
                          class="text-xs font-mono bg-emphasis rounded px-2 py-1.5 overflow-x-auto max-h-40 whitespace-pre-wrap break-all text-color"
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

              <div
                v-if="traceEntry.usage?.contextWindow"
                class="flex items-center justify-end text-xs text-muted-color mt-1"
              >
                <div class="flex items-center gap-1">
                  <i class="pi pi-chart-bar text-xs" />
                  <span class="font-medium text-color">
                    {{ traceEntry.usage.contextWindow }}
                    {{ $t('agent.editor.inspect.contextTokens') }}
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
          {{ $t('agent.editor.inspect.contextTokens') }}
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
</template>

<script setup>
import { ref, nextTick } from 'vue'
import { renderMarkdown } from '@planetcrust/human-js'

const props = defineProps({
  agent: { type: Object, required: true },
  conversation: { type: Object, required: true },
})

const emit = defineEmits(['select'])

const selectedIndex = ref(null)
const selectedType = ref(null)
const contextExpanded = ref(false)
const traceCardRefs = ref({})
const traceScrollContainer = ref(null)
const expandedTools = ref(new Set())
const expandedReasoning = ref(new Set())

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

function selectMessage(idx, type = 'response') {
  if (selectedIndex.value === idx && selectedType.value === type) {
    selectedIndex.value = null
    selectedType.value = null
  } else {
    selectedIndex.value = idx
    selectedType.value = type
  }
  emit('select', selectedIndex.value, selectedType.value)
}

function selectExternal(idx, type = 'response') {
  if (selectedIndex.value === idx && selectedType.value === type) {
    selectedIndex.value = null
    selectedType.value = null
  } else {
    selectedIndex.value = idx
    selectedType.value = type
    nextTick(() => {
      const card = traceCardRefs.value[idx]
      if (card) card.scrollIntoView({ behavior: 'smooth', block: 'nearest' })
    })
  }
  emit('select', selectedIndex.value, selectedType.value)
}

function highlightLatest() {
  if (!props.conversation.traceHistory.length) return
  selectedIndex.value = props.conversation.traceHistory.length - 1
  selectedType.value = 'response'
  emit('select', selectedIndex.value, selectedType.value)
}

defineExpose({ selectExternal, highlightLatest })
</script>
