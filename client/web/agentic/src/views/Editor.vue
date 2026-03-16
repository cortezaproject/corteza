<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ isCreate ? $t('agent.editor.titleCreate') : $t('agent.editor.titleEdit') }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <div v-else-if="agent" class="flex flex-col h-full overflow-hidden">
    <div class="container mx-auto p-4 flex-1 overflow-hidden min-h-0">
      <Card
        :pt="{
          body: { class: 'p-0 h-full flex flex-col' },
          content: { class: 'p-0 h-full flex flex-col min-h-0' },
        }"
        class="h-full overflow-hidden"
      >
        <template #content>
          <Tabs v-model:value="activeTab" class="flex flex-col h-full min-h-0">
            <TabList class="rounded-t-lg shrink-0">
              <Tab value="config">{{ $t('agent.editor.tabs.config') }}</Tab>
              <Tab value="exec">{{ $t('agent.editor.tabs.exec') }}</Tab>
            </TabList>

            <TabPanels class="flex-1 overflow-auto min-h-0 p-0">
              <TabPanel value="config" class="h-full p-4 flex flex-col gap-4">
                <Panel :header="$t('agent.editor.panels.general')" toggleable>
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div class="flex flex-col gap-1">
                      <label for="name" class="font-medium text-primary">
                        {{ $t('agent.editor.name.label') }}
                      </label>
                      <InputText id="name" v-model="agent.meta.short" />
                      <small class="text-muted-color">{{ $t('agent.editor.name.help') }}</small>
                    </div>
                    <div class="flex flex-col gap-1">
                      <label for="handle" class="font-medium text-primary">
                        {{ $t('agent.editor.handle.label') }}
                      </label>
                      <InputText id="handle" v-model="agent.handle" />
                      <small class="text-muted-color">{{ $t('agent.editor.handle.help') }}</small>
                    </div>
                    <div class="flex flex-col gap-1 md:col-span-2">
                      <label for="description" class="font-medium text-primary">
                        {{ $t('agent.editor.description.label') }}
                      </label>
                      <Textarea
                        id="description"
                        v-model="agent.meta.description"
                        rows="3"
                        autoResize
                      />
                      <small class="text-muted-color">
                        {{ $t('agent.editor.description.help') }}
                      </small>
                    </div>
                    <div class="flex flex-col gap-1">
                      <label for="status" class="font-medium text-primary">
                        {{ $t('agent.editor.status.label') }}
                      </label>
                      <Select
                        id="status"
                        v-model="agent.status"
                        :options="statusOptions"
                        optionLabel="label"
                        optionValue="value"
                      />
                      <small class="text-muted-color">{{ $t('agent.editor.status.help') }}</small>
                    </div>
                  </div>
                </Panel>

                <Panel :header="$t('agent.editor.panels.behavior')" toggleable>
                  <div class="grid grid-cols-1 gap-4">
                    <div class="flex flex-col gap-1">
                      <label for="systemPrompt" class="font-medium text-primary">
                        {{ $t('agent.editor.systemPrompt.label') }}
                      </label>
                      <Textarea
                        id="systemPrompt"
                        v-model="agent.behavior.systemPrompt"
                        rows="6"
                        autoResize
                      />
                      <small class="text-muted-color">
                        {{ $t('agent.editor.systemPrompt.help') }}
                      </small>
                    </div>

                    <!-- Guardrails chip input -->
                    <div class="flex flex-col gap-1">
                      <label for="guardrails" class="font-medium text-primary">
                        {{ $t('agent.editor.guardrails.label') }}
                      </label>
                      <div
                        class="flex flex-wrap gap-2 p-3 border border-surface rounded-lg bg-surface-ground min-h-12"
                      >
                        <Tag
                          v-for="(rule, idx) in agent.behavior.guardrails"
                          :key="idx"
                          :value="rule"
                          severity="info"
                          rounded
                          class="flex items-center gap-1"
                        >
                          <span>{{ rule }}</span>
                          <button
                            class="ml-1 cursor-pointer text-xs opacity-70 hover:opacity-100"
                            @click="removeGuardrail(idx)"
                          >
                            <i class="pi pi-times text-xs" />
                          </button>
                        </Tag>
                        <InputText
                          v-model="guardrailInput"
                          :placeholder="
                            agent.behavior.guardrails.length === 0 ? 'e.g. No PII disclosure' : ''
                          "
                          class="border-0 shadow-none flex-1 min-w-32 bg-transparent p-0"
                          size="small"
                          @keydown.enter.prevent="addGuardrail"
                        />
                      </div>
                      <small class="text-muted-color">
                        {{ $t('agent.editor.guardrails.help') }}
                      </small>
                    </div>

                    <!-- Inject system context toggle -->
                    <div
                      class="w-full md:w-1/2 flex items-start gap-3 p-3 border border-surface rounded-lg"
                    >
                      <div class="flex flex-col gap-0.5 flex-1">
                        <label
                          for="injectSystemContext"
                          class="font-medium text-primary cursor-pointer"
                        >
                          {{ $t('agent.editor.injectSystemContext.label') }}
                        </label>
                        <small class="text-muted-color">
                          {{ $t('agent.editor.injectSystemContext.help') }}
                        </small>
                      </div>
                      <ToggleSwitch
                        v-model="agent.behavior.injectSystemContext"
                        inputId="injectSystemContext"
                        class="shrink-0"
                      />
                    </div>
                  </div>
                </Panel>

                <Panel :header="$t('agent.editor.panels.execution')" toggleable>
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div class="flex flex-col gap-1">
                      <label for="provider" class="font-medium text-primary">
                        {{ $t('agent.editor.provider.label') }}
                      </label>
                      <CInputLLM id="provider" v-model="agent.execution.model.llmProviderID" />
                      <small class="text-muted-color">{{ $t('agent.editor.provider.help') }}</small>
                    </div>
                    <div class="flex flex-col gap-1">
                      <label for="model" class="font-medium text-primary">
                        {{ $t('agent.editor.model.label') }}
                      </label>
                      <CInputModel
                        id="model"
                        v-model="agent.execution.model.model"
                        :llmProviderID="agent.execution.model.llmProviderID"
                      />
                      <small class="text-muted-color">{{ $t('agent.editor.model.help') }}</small>
                    </div>
                    <div class="flex flex-col gap-1">
                      <label for="temperature" class="font-medium text-primary">
                        {{ $t('agent.editor.temperature.label') }} ({{
                          agent.execution.model.temperature
                        }})
                      </label>
                      <Slider
                        id="temperature"
                        v-model="agent.execution.model.temperature"
                        :min="0"
                        :max="1"
                        :step="0.1"
                        class="w-full mt-2"
                      />
                      <small class="text-muted-color">
                        {{ $t('agent.editor.temperature.help') }}
                      </small>
                    </div>
                    <div class="flex flex-col gap-1">
                      <label for="maxTokens" class="font-medium text-primary">
                        {{ $t('agent.editor.maxTokens.label') }}
                      </label>
                      <InputNumber
                        id="maxTokens"
                        v-model="agent.execution.model.maxTokens"
                        mode="decimal"
                        :useGrouping="false"
                      />
                      <small class="text-muted-color">
                        {{ $t('agent.editor.maxTokens.help') }}
                      </small>
                    </div>

                    <!-- Execution limits -->
                    <div class="md:col-span-2 border-t border-surface pt-4 mt-2">
                      <p class="font-medium text-primary mb-3">
                        {{ $t('agent.editor.limits.label') }}
                      </p>
                      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <div class="flex flex-col gap-1">
                          <label for="maxIterations" class="text-sm text-muted-color">
                            {{ $t('agent.editor.maxIterations.label') }}
                          </label>
                          <InputNumber
                            id="maxIterations"
                            v-model="agent.execution.limits.maxIterations"
                            mode="decimal"
                            :useGrouping="false"
                            :min="1"
                            :max="50"
                          />
                          <small class="text-muted-color">
                            {{ $t('agent.editor.maxIterations.help') }}
                          </small>
                        </div>
                        <div class="flex flex-col gap-1">
                          <label for="totalMaxTokens" class="text-sm text-muted-color">
                            {{ $t('agent.editor.totalMaxTokens.label') }}
                          </label>
                          <InputNumber
                            id="totalMaxTokens"
                            v-model="agent.execution.limits.maxTokens"
                            mode="decimal"
                            :useGrouping="false"
                          />
                          <small class="text-muted-color">
                            {{ $t('agent.editor.totalMaxTokens.help') }}
                          </small>
                        </div>
                        <div class="flex flex-col gap-1">
                          <label for="timeout" class="text-sm text-muted-color">
                            {{ $t('agent.editor.timeout.label') }}
                          </label>
                          <InputText
                            id="timeout"
                            v-model="agent.execution.limits.timeout"
                            placeholder="30s"
                          />
                          <small class="text-muted-color">
                            {{ $t('agent.editor.timeout.help') }}
                          </small>
                        </div>
                      </div>
                    </div>
                  </div>
                </Panel>
                <Panel :header="$t('agent.editor.panels.access')" toggleable>
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div class="flex flex-col gap-1">
                      <label for="contextNamespace" class="font-medium text-primary">
                        {{ $t('agent.editor.contextNamespace.label') }}
                      </label>
                      <CInputNamespace
                        id="contextNamespace"
                        v-model="agent.access.context.namespace"
                      />
                      <small class="text-muted-color">
                        {{ $t('agent.editor.contextNamespace.help') }}
                      </small>
                    </div>
                    <div class="flex flex-col gap-1">
                      <label for="contextModule" class="font-medium text-primary">
                        {{ $t('agent.editor.contextModule.label') }}
                      </label>
                      <CInputModule
                        id="contextModule"
                        v-model="agent.access.context.module"
                        :namespaceSlug="agent.access.context.namespace"
                      />
                      <small class="text-muted-color">
                        {{ $t('agent.editor.contextModule.help') }}
                      </small>
                    </div>
                    <div class="md:col-span-2 pt-2">
                      <div class="flex flex-col gap-1 mb-2">
                        <label class="font-medium text-primary">
                          {{ $t('agent.editor.tools.label') }}
                        </label>
                        <small class="text-muted-color">{{ $t('agent.editor.tools.help') }}</small>
                      </div>
                      <DataTable
                        v-model:selection="selectedTools"
                        :value="availableTools"
                        :loading="loadingTools"
                        dataKey="name"
                        class="border border-surface rounded overflow-hidden"
                      >
                        <template #empty>
                          <div class="text-center text-muted-color py-4">
                            {{ $t('agent.editor.tools.noTools') }}
                          </div>
                        </template>
                        <Column selectionMode="multiple" headerStyle="width: 3rem" />
                        <Column
                          field="description"
                          :header="$t('agent.editor.tools.toolDescription')"
                        />
                        <Column :header="$t('agent.editor.tools.toolHints')">
                          <template #body="{ data }">
                            <InputText
                              :modelValue="getToolHints(data.name)"
                              @update:modelValue="setToolHints(data.name, $event)"
                              class="w-full"
                              size="small"
                              :disabled="!selectedTools.some(t => t.name === data.name)"
                            />
                          </template>
                        </Column>
                      </DataTable>
                    </div>
                  </div>
                </Panel>

                <Panel :header="$t('agent.editor.panels.invocation')" toggleable :collapsed="true">
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div class="flex flex-col gap-1 justify-center">
                      <div class="flex items-center gap-2">
                        <Checkbox
                          v-model="agent.invocation.user.enabled"
                          inputId="userEnabled"
                          :binary="true"
                        />
                        <label for="userEnabled">
                          {{ $t('agent.editor.userEnabled.label') }}
                        </label>
                      </div>
                      <small class="text-muted-color ml-7">
                        {{ $t('agent.editor.userEnabled.help') }}
                      </small>
                    </div>
                    <div class="flex flex-col gap-1 justify-center">
                      <div class="flex items-center gap-2">
                        <Checkbox
                          v-model="agent.invocation.system.enabled"
                          inputId="systemEnabled"
                          :binary="true"
                        />
                        <label for="systemEnabled">
                          {{ $t('agent.editor.systemEnabled.label') }}
                        </label>
                      </div>
                      <small class="text-muted-color ml-7">
                        {{ $t('agent.editor.systemEnabled.help') }}
                      </small>
                    </div>
                    <div class="flex flex-col gap-1">
                      <label for="serviceAccount" class="font-medium text-primary">
                        {{ $t('agent.editor.serviceAccount.label') }}
                      </label>
                      <InputText
                        id="serviceAccount"
                        v-model="agent.invocation.system.serviceAccount"
                      />
                      <small class="text-muted-color">
                        {{ $t('agent.editor.serviceAccount.help') }}
                      </small>
                    </div>
                    <div class="flex flex-col gap-1">
                      <label for="outputFormat" class="font-medium text-primary">
                        {{ $t('agent.editor.outputFormat.label') }}
                      </label>
                      <InputText id="outputFormat" v-model="agent.invocation.system.outputFormat" />
                      <small class="text-muted-color">
                        {{ $t('agent.editor.outputFormat.help') }}
                      </small>
                    </div>
                  </div>
                </Panel>
              </TabPanel>

              <TabPanel value="exec" class="h-full p-0 flex flex-row overflow-hidden">
                <!-- Chat area -->
                <div class="flex-1 flex flex-col border-r border-surface min-w-0">
                  <!-- Conversation tabs bar -->
                  <div
                    class="flex items-center gap-0 border-b border-surface shrink-0 bg-surface-ground px-2"
                  >
                    <div class="flex items-center gap-0 flex-1 overflow-x-auto">
                      <button
                        v-for="(conv, idx) in conversations"
                        :key="idx"
                        class="flex items-center gap-1.5 p-3 text-sm border-b-2 whitespace-nowrap transition-colors"
                        :class="
                          activeConvIndex === idx
                            ? 'border-primary text-primary font-medium'
                            : 'border-transparent text-muted-color hover:text-color hover:border-surface'
                        "
                        @click="activeConvIndex = idx"
                      >
                        <i class="pi pi-comments text-xs" />
                        <span>{{ conv.label }}</span>
                        <i
                          v-if="conversations.length > 1"
                          class="pi pi-times text-xs opacity-50 hover:opacity-100 ml-1"
                          @click.stop="closeConversation(idx)"
                        />
                      </button>
                    </div>
                    <Button
                      icon="pi pi-plus"
                      text
                      rounded
                      severity="secondary"
                      size="small"
                      @click="addConversation"
                      :disabled="isCreate"
                      class="shrink-0"
                    />
                  </div>
                  <div class="flex-1 p-4 overflow-y-auto flex flex-col gap-4">
                    <div
                      v-for="(msg, index) in activeConversation.messages"
                      :key="index"
                      class="flex flex-col gap-1"
                      :class="msg.role === 'user' ? 'items-end' : 'items-start'"
                    >
                      <div
                        :ref="
                          el => {
                            if (msg.traceIndex !== undefined) {
                              if (msg.role === 'agent') chatMsgRefs[msg.traceIndex] = el
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
                          msg.role === 'agent' ? 'cursor-pointer hover:shadow-md' : '',
                          selectedTraceIndex === msg.traceIndex && selectedTraceType === (msg.role === 'user' ? 'prompt' : 'response')
                            ? 'ring-2 ring-primary ring-offset-1'
                            : '',
                        ]"
                        @click="
                          msg.traceIndex !== undefined
                            ? selectMessage(msg.traceIndex, msg.role === 'user' ? 'prompt' : 'response')
                            : null
                        "
                      >
                        <div
                          v-if="msg.role === 'agent'"
                          class="rt-content"
                          v-html="renderMarkdown(msg.content)"
                        />
                        <template v-else>{{ msg.content }}</template>
                      </div>
                    </div>
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
                  <div class="p-3 border-t border-surface flex gap-2 shrink-0 bg-surface">
                    <InputText
                      v-model="chatInput"
                      :placeholder="$t('agent.editor.playground.placeholder')"
                      class="flex-1"
                      @keyup.enter="sendChatMessage"
                      :disabled="executing || isCreate"
                    />
                    <Button
                      icon="pi pi-send"
                      @click="sendChatMessage"
                      :disabled="executing || !chatInput.trim() || isCreate"
                      :loading="executing"
                    />
                  </div>
                </div>

                <!-- Structured trace panel -->
                <div class="w-1/3 flex flex-col bg-surface border-l border-surface min-w-0">
                  <div class="flex-1 overflow-y-auto" ref="traceScrollContainer">
                    <!-- Empty state -->
                    <div
                      v-if="!activeConversation.traceHistory.length"
                      class="p-4 text-sm text-muted-color text-center"
                    >
                      {{ $t('agent.editor.playground.emptyTrace') }}
                    </div>

                    <!-- Card-based trace view -->
                    <div v-else class="p-3 flex flex-col gap-4">
                      <div
                        v-for="(traceEntry, tIdx) in activeConversation.traceHistory"
                        :key="'t' + tIdx"
                        :ref="
                          el => {
                            traceCardRefs[tIdx] = el
                          }
                        "
                        class="transition-all duration-200 cursor-pointer hover:shadow-md"
                        @click="selectMessage(tIdx)"

                        <!-- Agent Response Panel -->
                        <Panel
                          toggleable
                          :class="selectedTraceIndex === tIdx && selectedTraceType === 'response' ? 'border-primary shadow-md' : ''"
                          @click.stop="selectMessage(tIdx, 'response')"
                        >
                          <template #header>
                            <div class="flex items-center gap-2">
                              <i
                                class="pi pi-sparkles"
                                :class="
                                  selectedTraceIndex === tIdx && selectedTraceType === 'response' ? 'text-primary' : 'text-muted-color'
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
                                  v-if="decision.usage?.totalTokens"
                                  class="ml-auto text-xs text-muted-color"
                                >
                                  {{ decision.usage.totalTokens }}
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

                            <!-- Token breakdown footer -->
                            <div
                              v-if="traceEntry.usage"
                              class="flex items-center justify-between text-xs text-muted-color mt-1"
                            >
                              <div v-if="(traceEntry.decisions || []).length" class="flex items-center gap-1">
                                <i class="pi pi-history text-xs" />
                                <span>
                                  {{ ((arr) => arr[arr.length - 1]?.usage?.inputTokens || 0)(traceEntry.decisions) }}
                                  {{ $t('agent.editor.playground.traceContextTokens') }}
                                  <template v-if="agent.execution.limits.maxTokens > 0">
                                    ({{ Math.round((((arr) => arr[arr.length - 1]?.usage?.inputTokens || 0)(traceEntry.decisions) / agent.execution.limits.maxTokens) * 100) }}%)
                                  </template>
                                </span>
                              </div>
                              <div class="flex items-center gap-1">
                                <i class="pi pi-chart-bar text-xs" />
                                <span class="font-medium text-color">
                                  {{ traceEntry.usage.totalTokens }}
                                  {{ $t('agent.editor.playground.traceTokens') }}
                                </span>
                              </div>
                            </div>
                          </div>
                        </Panel>
                      </div>
                    </div>
                  </div>

                  <!-- Cumulative total footer -->
                  <div
                    v-if="activeConversation.traceHistory.length"
                    class="px-3 py-2 flex items-center justify-end text-xs text-muted-color border-t border-surface shrink-0"
                  >
                    <div class="flex items-center gap-2">
                      <div v-if="(activeConversation.traceHistory[activeConversation.traceHistory.length - 1]?.decisions || []).length" class="flex items-center gap-1">
                        <i class="pi pi-history text-xs" />
                        <span>
                          {{
                            ((arr) => arr[arr.length - 1]?.usage?.inputTokens || 0)(
                              activeConversation.traceHistory[activeConversation.traceHistory.length - 1]?.decisions || [],
                            )
                          }}
                          {{ $t('agent.editor.playground.traceContextTokens') }}
                          <template v-if="agent.execution.limits.maxTokens > 0">
                            ({{
                              Math.round(
                                (((arr) => arr[arr.length - 1]?.usage?.inputTokens || 0)(
                                  activeConversation.traceHistory[activeConversation.traceHistory.length - 1]?.decisions || [],
                                ) /
                                  agent.execution.limits.maxTokens) *
                                  100,
                              )
                            }}%)
                          </template>
                        </span>
                      </div>
                      <span class="text-muted-color">·</span>
                      <div class="flex items-center gap-1">
                        <i class="pi pi-calculator text-xs" />
                        <span class="font-medium text-color">
                          {{
                            activeConversation.traceHistory.reduce(
                              (sum, t) => sum + (t.usage?.totalTokens || 0),
                              0,
                            )
                          }}
                          {{ $t('agent.editor.playground.traceTotalTokens') }}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
              </TabPanel>
            </TabPanels>
          </Tabs>
        </template>
      </Card>
    </div>

    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="flex items-center justify-between p-3">
        <Button
          :label="$t('general.label.back')"
          icon="pi pi-arrow-left"
          severity="secondary"
          @click="router.back()"
        />
        <div class="flex gap-2">
          <CInputDelete
            v-if="!isCreate && agent.canDeleteAgent !== false"
            :label="$t('general.label.delete')"
            :message="$t('agent.list.delete')"
            :header="agent?.meta?.short || agent?.handle || $t('general.label.delete')"
            @confirm="handleDelete"
          />
          <Button
            :label="$t('general.label.save')"
            icon="pi pi-save"
            :loading="saving"
            @click="handleSubmit"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, inject, nextTick, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { useAgentStore } from '@/stores/agent'
import { useRoute, useRouter } from 'vue-router'

// Components
import Button from 'primevue/button'
import Card from 'primevue/card'
import Checkbox from 'primevue/checkbox'
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import InputNumber from 'primevue/inputnumber'
import InputText from 'primevue/inputtext'
import Panel from 'primevue/panel'
import ProgressSpinner from 'primevue/progressspinner'
import Select from 'primevue/select'
import Slider from 'primevue/slider'

import Tab from 'primevue/tab'
import TabList from 'primevue/tablist'
import TabPanel from 'primevue/tabpanel'
import TabPanels from 'primevue/tabpanels'
import Tabs from 'primevue/tabs'
import Tag from 'primevue/tag'
import Textarea from 'primevue/textarea'
import ToggleSwitch from 'primevue/toggleswitch'
import { components } from '@cortezaproject/corteza-vue-next'
const { CInputLLM, CInputModel, CInputNamespace, CInputModule, CInputDelete } = components

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')
const agentStore = useAgentStore()

const loading = ref(false)
const saving = ref(false)
const agent = ref(null)
const activeTab = ref('config')

// Conversation tabs
let convCounter = 1
const conversations = ref([
  { label: `Chat 1`, messages: [], conversationID: null, traceHistory: [] },
])
const activeConvIndex = ref(0)
const activeConversation = computed(() => conversations.value[activeConvIndex.value])

const chatInput = ref('')
const executing = ref(false)
const showRawTrace = ref(false)
const selectedTraceIndex = ref(null)
const selectedTraceType = ref(null)
const chatMsgRefs = ref({})
const chatPromptRefs = ref({})
const traceCardRefs = ref({})
const traceScrollContainer = ref(null)
const isCreate = computed(() => !route.params.agentID)

// Guardrails
const guardrailInput = ref('')

const statusOptions = computed(() => [
  { label: t('agent.list.status.active'), value: 'active' },
  { label: t('agent.list.status.inactive'), value: 'inactive' },
])

const availableTools = ref([])
const loadingTools = ref(false)

const selectedTools = ref([])
const emptyAgent = () => ({
  handle: '',
  status: 'active',
  meta: {
    short: '',
    description: '',
  },
  behavior: {
    systemPrompt: '',
    guardrails: [],
    injectSystemContext: false,
  },
  execution: {
    model: {
      llmProviderID: '0',
      model: '',
      temperature: 0.7,
      maxTokens: 1000,
    },
    limits: {
      maxIterations: 10,
      maxTokens: 4000,
      timeout: '30s',
      softLimitRatio: 0.8,
    },
  },
  access: {
    context: {
      namespace: '',
      module: '',
    },
    tools: [],
    allow: [],
  },
  invocation: {
    user: { enabled: true },
    system: {
      enabled: false,
      serviceAccount: '',
      inputSchema: null,
      outputFormat: '',
    },
  },
})

function applyAgentData(res) {
  agent.value = {
    ...emptyAgent(),
    ...res,
    meta: { ...emptyAgent().meta, ...(res.meta || {}) },
    behavior: { ...emptyAgent().behavior, ...(res.behavior || {}) },
    execution: {
      model: { ...emptyAgent().execution.model, ...(res.execution?.model || {}) },
      limits: { ...emptyAgent().execution.limits, ...(res.execution?.limits || {}) },
    },
    access: {
      ...emptyAgent().access,
      ...(res.access || {}),
      context: { ...emptyAgent().access.context, ...(res.access?.context || {}) },
    },
    invocation: {
      ...emptyAgent().invocation,
      ...(res.invocation || {}),
      user: { ...emptyAgent().invocation.user, ...(res.invocation?.user || {}) },
      system: { ...emptyAgent().invocation.system, ...(res.invocation?.system || {}) },
    },
  }

  if (!Array.isArray(agent.value.behavior.guardrails)) {
    agent.value.behavior.guardrails = []
  }

  if (!agent.value.execution.model.llmProviderID) {
    agent.value.execution.model.llmProviderID = '0'
  }

  initToolSelection()
}

async function loadAgent() {
  const agentID = route.params.agentID
  if (!agentID) {
    agent.value = emptyAgent()
    return
  }

  loading.value = true
  try {
    const res = await $SystemAPI.agentRead({ agentID })
    applyAgentData(res)
  } catch (err) {
    $toast.toastDanger(t('notification.agent.loadFailed'))
    router.push({ name: 'root' })
  } finally {
    loading.value = false
  }
}

async function handleSubmit() {
  saving.value = true
  try {
    if (isCreate.value) {
      const created = await $SystemAPI.agentCreate(agent.value)
      $toast.toastSuccess(t('notification.agent.created'))
      router.push({ name: 'agent.edit', params: { agentID: created.agentID } })
    } else {
      const updated = await $SystemAPI.agentUpdate({
        agentID: route.params.agentID,
        ...agent.value,
      })
      applyAgentData(updated)
      $toast.toastSuccess(t('notification.agent.saved'))
    }
  } catch (err) {
    console.error(err)
    $toast.toastDanger(t('notification.agent.saveFailed'))
  } finally {
    saving.value = false
  }
}



async function handleDelete() {
  try {
    await $SystemAPI.agentDelete({ agentID: route.params.agentID })
    $toast.toastSuccess(t('notification.agent.deleted'))
    router.push({ name: 'root' })
  } catch (err) {
    console.error(err)
    $toast.toastDanger(t('notification.agent.deleteFailed'))
  }
}

async function sendChatMessage() {
  if (!chatInput.value.trim() || executing.value) return

  const conv = activeConversation.value
  const input = chatInput.value
  chatInput.value = ''

  conv.messages.push({ role: 'user', content: input, traceIndex: conv.traceHistory.length })
  executing.value = true

  try {
    const res = await $SystemAPI.agentExec({
      agentID: route.params.agentID,
      input: input,
      ...(conv.conversationID ? { conversationID: conv.conversationID } : {}),
    })

    if (res?.conversationID) {
      conv.conversationID = res.conversationID
    }

    conv.messages.push({
      role: 'agent',
      content:
        res?.output || (typeof res === 'string' ? res : res?.response?.text || JSON.stringify(res)),
      usage: res?.usage || null,
      traceIndex: conv.traceHistory.length,
    })

    // Append structured trace
    conv.traceHistory.push({
      prompt: input,
      decisions: res?.decisions || [],
      toolCalls: res?.toolCalls || [],
      usage: res?.usage || null,
      conversationTokens: res?.conversationTokens || 0,
    })

    // Auto-select the new trace card
    selectedTraceIndex.value = conv.traceHistory.length - 1
  } catch (err) {
    console.error(err)
    $toast.toastDanger(t('notification.agent.execFailed'))
    conv.messages.push({ role: 'agent', content: 'Error: ' + err.message })
    conv.traceHistory.push({ error: err.message })
  } finally {
    executing.value = false
  }
}

function addConversation() {
  convCounter++
  conversations.value.push({
    label: `Chat ${convCounter}`,
    messages: [],
    conversationID: null,
    traceHistory: [],
  })
  activeConvIndex.value = conversations.value.length - 1
}

function closeConversation(idx) {
  conversations.value.splice(idx, 1)
  if (activeConvIndex.value >= conversations.value.length) {
    activeConvIndex.value = conversations.value.length - 1
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

// Expand/collapse state for tool calls and reasoning
const expandedTools = ref(new Set())
const expandedReasoning = ref(new Set())

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

// Simple markdown rendering (basic: bold, italic, code blocks, inline code, lists, headers)
function renderMarkdown(text) {
  if (!text) return ''

  let html = text
    // Escape HTML
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')

  // Code blocks (``` ... ```)
  html = html.replace(
    /```(\w*)\n([\s\S]*?)```/g,
    '<pre class="bg-surface-ground p-3 rounded-lg my-2 overflow-x-auto text-sm"><code>$2</code></pre>',
  )

  // Inline code
  html = html.replace(
    /`([^`]+)`/g,
    '<code class="bg-surface-ground px-1.5 py-0.5 rounded text-sm">$1</code>',
  )

  // Headers
  html = html.replace(/^### (.+)$/gm, '<h4 class="font-bold mt-3 mb-1">$1</h4>')
  html = html.replace(/^## (.+)$/gm, '<h3 class="font-bold text-lg mt-3 mb-1">$1</h3>')
  html = html.replace(/^# (.+)$/gm, '<h2 class="font-bold text-xl mt-3 mb-1">$1</h2>')

  // Bold
  html = html.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')

  // Italic
  html = html.replace(/\*(.+?)\*/g, '<em>$1</em>')

  // Unordered lists
  html = html.replace(/^[*-] (.+)$/gm, '<li class="ml-4 list-disc">$1</li>')

  // Ordered lists
  html = html.replace(/^\d+\. (.+)$/gm, '<li class="ml-4 list-decimal">$1</li>')

  // Line breaks (double newline = paragraph)
  html = html.replace(/\n\n/g, '</p><p class="my-2">')
  html = html.replace(/\n/g, '<br>')

  return `<div class="prose-sm">${html}</div>`
}

// Guardrails helpers
function addGuardrail() {
  const val = guardrailInput.value.trim()
  if (val && !agent.value.behavior.guardrails.includes(val)) {
    agent.value.behavior.guardrails.push(val)
    guardrailInput.value = ''
  }
}

function removeGuardrail(index) {
  agent.value.behavior.guardrails.splice(index, 1)
}

onMounted(() => {
  loadAgent()
  fetchAvailableTools()
})

async function fetchAvailableTools() {
  loadingTools.value = true
  try {
    const response = await $SystemAPI.mcpListTools()
    availableTools.value = Array.isArray(response) ? response : response.set || []
  } catch {
    availableTools.value = []
  } finally {
    loadingTools.value = false
  }

  initToolSelection()
}

function initToolSelection() {
  // Wait for both agent and tools to be loaded
  if (!availableTools.value.length || loading.value) return

  if (!agent.value?.access?.tools?.length) {
    // New agent: enable all tools by default
    agent.value.access.tools = availableTools.value.map(t => ({
      name: t.name,
      hints: t.description || '',
    }))
    selectedTools.value = [...availableTools.value]
  } else {
    // Existing agent: select only the saved tools
    const enabledNames = new Set((agent.value?.access?.tools || []).map(t => t.name))
    selectedTools.value = availableTools.value.filter(t => enabledNames.has(t.name))
  }
}

watch(
  selectedTools,
  selected => {
    const selectedNames = new Set(selected.map(t => t.name))
    // Keep existing tool entries for selected tools, add new ones
    const existing = new Map(agent.value.access.tools.map(t => [t.name, t]))
    agent.value.access.tools = selected.map(t => {
      return existing.get(t.name) || { name: t.name, hints: t.description || '' }
    })
    // Clean up deselected
    agent.value.access.tools = agent.value.access.tools.filter(t => selectedNames.has(t.name))
  },
  { deep: true },
)

function getToolHints(name) {
  const tool = agent.value.access.tools.find(t => t.name === name)
  return tool?.hints || ''
}

function setToolHints(name, value) {
  const tool = agent.value.access.tools.find(t => t.name === name)
  if (tool) tool.hints = value
}
</script>
