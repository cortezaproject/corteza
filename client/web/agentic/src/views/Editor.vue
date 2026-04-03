<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ isCreate ? $t('agent.editor.titleCreate') : $t('agent.editor.titleEdit') }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <div v-else-if="agent" class="flex flex-col h-full overflow-hidden">
    <div class="container mx-auto p-4 flex-1 overflow-hidden min-h-0 flex flex-col gap-4">
      <div v-if="!isCreate" class="flex justify-end gap-2 shrink-0">
        <CPermissionsButton
          v-tooltip.bottom="$t('general.label.permissions')"
          :resource="`corteza::system:agent/${agent.agentID}`"
          :title="agent.meta?.short || agent.handle || agent.agentID"
          :target="agent.meta?.short || agent.handle || agent.agentID"
        />
      </div>
      <Card
        :pt="{
          body: { class: 'p-0 h-full flex flex-col' },
          content: { class: 'p-0 h-full flex flex-col min-h-0' },
        }"
        class="flex-1 min-h-0 overflow-hidden"
      >
        <template #content>
          <Tabs v-model:value="activeTab" class="flex flex-col h-full min-h-0">
            <TabList class="rounded-t-lg shrink-0">
              <Tab value="config">{{ $t('agent.editor.tabs.config') }}</Tab>
              <Tab value="exec">{{ $t('agent.editor.tabs.exec') }}</Tab>
              <Tab value="history" v-if="!isCreate">{{ $t('agent.editor.tabs.history') }}</Tab>
            </TabList>

            <TabPanels class="flex-1 min-h-0 p-0">
              <TabPanel value="config" class="h-full p-4 flex flex-col gap-4 overflow-auto">
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
                      <small class="text-muted-color">
                        {{ $t('agent.editor.temperature.help') }}
                      </small>
                      <Slider
                        id="temperature"
                        v-model="agent.execution.model.temperature"
                        :min="0"
                        :max="1"
                        :step="0.1"
                        class="w-full mt-2"
                      />
                    </div>

                    <!-- Execution limits -->
                    <div class="md:col-span-2">
                      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <div class="flex flex-col gap-1">
                          <label for="maxIterations" class="font-medium text-primary">
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
                          <label for="contextWindow" class="font-medium text-primary">
                            {{ $t('agent.editor.contextWindow.label') }}
                          </label>
                          <InputNumber
                            id="contextWindow"
                            v-model="agent.execution.limits.contextWindow"
                            mode="decimal"
                            :useGrouping="false"
                          />
                          <small class="text-muted-color">
                            {{ $t('agent.editor.contextWindow.help') }}
                          </small>
                        </div>
                        <div class="flex flex-col gap-1">
                          <label for="outputTokens" class="font-medium text-primary">
                            {{ $t('agent.editor.outputTokens.label') }}
                          </label>
                          <InputNumber
                            id="outputTokens"
                            v-model="agent.execution.limits.outputTokens"
                            mode="decimal"
                            :useGrouping="false"
                          />
                          <small class="text-muted-color">
                            {{ $t('agent.editor.outputTokens.help') }}
                          </small>
                        </div>
                        <div class="flex flex-col gap-1">
                          <label for="timeout" class="font-medium text-primary">
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

                <Panel :header="$t('agent.editor.panels.behavior')" toggleable>
                  <div class="grid grid-cols-1 gap-4">
                    <div class="flex flex-col gap-1">
                      <label for="systemPrompt" class="font-medium text-primary">
                        {{ $t('agent.editor.systemPrompt.label') }}
                        <span class="text-red-500">*</span>
                      </label>
                      <Textarea
                        id="systemPrompt"
                        v-model="agent.behavior.systemPrompt"
                        rows="6"
                        autoResize
                        :invalid="submitted && !agent.behavior.systemPrompt?.trim()"
                      />
                      <small
                        v-if="submitted && !agent.behavior.systemPrompt?.trim()"
                        class="text-red-500"
                      >
                        {{ $t('agent.editor.systemPrompt.required') }}
                      </small>
                      <small v-else class="text-muted-color">
                        {{ $t('agent.editor.systemPrompt.help') }}
                      </small>
                    </div>

                    <!-- Guardrails -->
                    <div class="flex flex-col gap-1">
                      <label class="font-medium text-primary">
                        {{ $t('agent.editor.guardrails.label') }}
                      </label>
                      <small class="text-muted-color">
                        {{ $t('agent.editor.guardrails.help') }}
                      </small>
                      <div class="flex flex-col gap-2 mt-2">
                        <div
                          v-for="(rule, idx) in agent.behavior.guardrails"
                          :key="idx"
                          class="flex items-center gap-2"
                        >
                          <InputText
                            v-model="agent.behavior.guardrails[idx]"
                            class="flex-1"
                            :placeholder="$t('agent.editor.guardrails.placeholder')"
                          />
                          <Button
                            icon="pi pi-trash"
                            severity="danger"
                            text
                            size="small"
                            @click="removeGuardrail(idx)"
                          />
                        </div>
                        <div>
                          <Button
                            :label="$t('agent.editor.guardrails.add')"
                            icon="pi pi-plus"
                            severity="secondary"
                            outlined
                            size="small"
                            @click="addGuardrail"
                          />
                        </div>
                      </div>
                    </div>
                  </div>
                </Panel>

                <Panel :header="$t('agent.editor.panels.tcl')" toggleable>
                  <div class="flex flex-col gap-3">
                    <CInputToggleCard
                      v-model="agent.behavior.treatyCLEnabled"
                      :label="$t('agent.editor.tcl.enabledLabel')"
                      :description="$t('agent.editor.tcl.enabledHelp')"
                    />

                    <template v-if="agent.behavior.treatyCLEnabled">
                      <div class="flex flex-col gap-1">
                        <label class="font-medium text-primary">
                          {{ $t('agent.editor.tcl.temperature.label') }} ({{
                            agent.behavior.tclTemperature
                          }})
                        </label>
                        <small class="text-muted-color">
                          {{ $t('agent.editor.tcl.temperature.help') }}
                        </small>
                        <Slider
                          v-model="agent.behavior.tclTemperature"
                          :min="1"
                          :max="10"
                          :step="1"
                          class="w-full mt-2"
                        />
                      </div>

                      <div class="flex flex-col gap-1">
                        <label class="font-medium text-primary">
                          {{ $t('agent.editor.tcl.articles.label') }}
                        </label>
                        <small class="text-muted-color">
                          {{ $t('agent.editor.tcl.articles.help') }}
                        </small>
                      </div>

                      <div
                        v-for="group in tclSortedArticleGroups"
                        :key="group.treatyLabel"
                        class="flex flex-col gap-2 mt-2"
                      >
                        <span
                          class="text-sm font-semibold text-muted-color uppercase tracking-wide"
                        >
                          {{ group.treatyLabel }}
                        </span>
                        <CInputToggleCard
                          v-for="article in group.items"
                          :key="article.id"
                          :modelValue="agent.behavior.tclArticles.includes(article.id)"
                          :label="article.label"
                          :description="article.interpretation"
                          :disabled="article.hardwired"
                          @update:modelValue="toggleTclArticle(article.id, $event)"
                        />
                      </div>
                    </template>
                  </div>
                </Panel>

                <Panel :header="$t('agent.editor.panels.knowledgeBase')" toggleable>
                  <div class="flex flex-col gap-3">
                    <div class="flex flex-col gap-1">
                      <label class="font-medium text-primary">
                        {{ $t('agent.editor.knowledgeBases.label') }}
                      </label>
                      <small class="text-muted-color">
                        {{ $t('agent.editor.knowledgeBases.help') }}
                      </small>
                      <CInputKnowledgeBase
                        v-model="agent.behavior.knowledgeBases"
                        :placeholder="$t('agent.editor.knowledgeBases.placeholder')"
                        :create-label="$t('agent.editor.knowledgeBases.createNew')"
                        :create-dialog-label="$t('agent.editor.knowledgeBases.dialogCreate')"
                        :edit-label="$t('agent.editor.knowledgeBases.dialogEdit')"
                        :title-label="$t('agent.editor.knowledgeBases.title')"
                        :description-label="$t('agent.editor.knowledgeBases.description')"
                        :description-help="$t('agent.editor.knowledgeBases.descriptionHelp')"
                        :compose-context-label="$t('agent.editor.knowledgeBases.composeContext')"
                        :namespace-placeholder="
                          $t('agent.editor.knowledgeBases.namespacePlaceholder')
                        "
                        :modules-placeholder="$t('agent.editor.knowledgeBases.modulesPlaceholder')"
                        :add-namespace-label="$t('agent.editor.knowledgeBases.addNamespace')"
                        :save-label="$t('general.label.save')"
                        :cancel-label="$t('general.label.cancel')"
                      />
                    </div>

                    <!-- System Context as a KB-style entry in the list -->
                    <CInputToggleCard
                      v-model="agent.behavior.injectSystemContext"
                      :label="$t('agent.editor.injectSystemContext.label')"
                      :description="$t('agent.editor.injectSystemContext.help')"
                      dimWhenOff
                    />
                  </div>
                </Panel>
                <Panel :header="$t('agent.editor.panels.tools')" toggleable>
                  <div class="flex flex-col gap-3">
                    <div class="flex flex-col gap-1">
                      <label class="font-medium text-primary">
                        {{ $t('agent.editor.tools.label') }}
                      </label>
                      <small class="text-muted-color">{{ $t('agent.editor.tools.help') }}</small>
                    </div>

                    <Select
                      v-model="toolPickerSelection"
                      :options="unselectedTools"
                      option-label="description"
                      :placeholder="$t('agent.editor.tools.selectPlaceholder')"
                      :loading="loadingTools"
                      class="w-full"
                      filter
                      fluid
                      showClear
                      @update:model-value="onToolPickerSelect"
                    >
                      <template #option="{ option }">
                        <span>{{ option.title }}</span>
                      </template>
                    </Select>

                    <div v-if="selectedTools.length" class="flex flex-col gap-2">
                      <div
                        v-for="tool in selectedTools"
                        :key="tool.name"
                        class="group flex items-start gap-3 p-3 border border-surface rounded-lg hover:bg-emphasis transition-colors"
                      >
                        <div class="flex flex-col gap-1 flex-1 min-w-0">
                          <span class="font-medium text-color text-sm truncate">
                            {{ tool.title }}
                          </span>
                          <InputText
                            :modelValue="getToolHints(tool.name)"
                            @update:modelValue="setToolHints(tool.name, $event)"
                            class="w-full mt-1"
                            size="small"
                            :placeholder="$t('agent.editor.tools.hintPlaceholder')"
                          />
                        </div>
                        <Button
                          icon="pi pi-trash"
                          severity="danger"
                          text
                          size="small"
                          class="shrink-0 opacity-0 group-hover:opacity-100 focus-within:opacity-100 transition-opacity"
                          @click="removeTool(tool)"
                        />
                      </div>
                    </div>

                    <Divider class="my-4" />

                    <!-- TAQs -->
                    <div class="flex flex-col gap-1">
                      <label class="font-medium text-primary">
                        {{ $t('agent.editor.taqs.label') }}
                      </label>
                      <small class="text-muted-color">{{ $t('agent.editor.taqs.help') }}</small>
                    </div>

                    <CInputTAQ
                      v-model="taqPickerSelection"
                      :placeholder="$t('agent.editor.taqs.selectPlaceholder')"
                      @update:model-value="onTaqPickerSelect"
                    />

                    <div v-if="agent.access.taqs?.length" class="flex flex-col gap-2">
                      <div
                        v-for="(taq, idx) in agent.access.taqs"
                        :key="taq.id"
                        class="group flex items-start gap-3 p-3 border border-surface rounded-lg hover:bg-emphasis transition-colors"
                      >
                        <div class="flex flex-col gap-1 flex-1 min-w-0">
                          <span class="font-medium text-color text-sm truncate">
                            {{ loadedTaqNames[taq.id] || $t('general.label.loading') }}
                          </span>
                          <InputText
                            v-model="agent.access.taqs[idx].hints"
                            class="w-full mt-1"
                            size="small"
                            :placeholder="$t('agent.editor.taqs.hintPlaceholder')"
                          />
                        </div>
                        <Button
                          icon="pi pi-trash"
                          severity="danger"
                          text
                          size="small"
                          class="shrink-0 opacity-0 group-hover:opacity-100 focus-within:opacity-100 transition-opacity"
                          @click="removeTaq(idx)"
                        />
                      </div>
                    </div>

                    <Divider class="my-4" />

                    <!-- Workflows -->
                    <div class="flex flex-col gap-1">
                      <label class="font-medium text-primary">
                        {{ $t('agent.editor.workflows.label') }}
                      </label>
                      <small class="text-muted-color">{{ $t('agent.editor.workflows.help') }}</small>
                    </div>

                    <CInputWorkflow
                      v-model="workflowPickerSelection"
                      :placeholder="$t('agent.editor.workflows.selectPlaceholder')"
                      @update:model-value="onWorkflowPickerSelect"
                    />

                    <div v-if="agent.access.workflows?.length" class="flex flex-col gap-2">
                      <div
                        v-for="(workflow, idx) in agent.access.workflows"
                        :key="workflow.id"
                        class="group flex items-start gap-3 p-3 border border-surface rounded-lg hover:bg-emphasis transition-colors"
                      >
                        <div class="flex flex-col gap-1 flex-1 min-w-0">
                          <span class="font-medium text-color text-sm truncate">
                            {{ loadedWorkflowNames[workflow.id] || $t('general.label.loading') }}
                          </span>
                          <InputText
                            v-model="agent.access.workflows[idx].hints"
                            class="w-full mt-1"
                            size="small"
                            :placeholder="$t('agent.editor.workflows.hintPlaceholder')"
                          />
                        </div>
                        <Button
                          icon="pi pi-trash"
                          severity="danger"
                          text
                          size="small"
                          class="shrink-0 opacity-0 group-hover:opacity-100 focus-within:opacity-100 transition-opacity"
                          @click="removeWorkflow(idx)"
                        />
                      </div>
                    </div>
                  </div>
                </Panel>

                <Panel :header="$t('agent.editor.panels.invocation')" toggleable>
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <!-- User Invocation Group -->
                    <div class="flex flex-col gap-4">
                      <CInputToggleCard
                        v-model="agent.invocation.user.enabled"
                        :label="$t('agent.editor.userEnabled.label')"
                        :description="$t('agent.editor.userEnabled.help')"
                      />
                      
                      <div class="flex flex-col gap-1" :class="{ 'opacity-50 pointer-events-none': !agent.invocation.user.enabled }">
                        <label for="sidebarRoles" class="font-medium text-primary">
                          {{ $t('agent.editor.sidebarRoles.label') }}
                        </label>
                        <CInputRole
                          id="sidebarRoles"
                          v-model="agent.meta.sidebarRoles"
                          :multiple="true"
                          :disabled="!agent.invocation.user.enabled"
                        />
                        <small class="text-muted-color">
                          {{ $t('agent.editor.sidebarRoles.help') }}
                        </small>
                      </div>
                    </div>

                    <!-- System Invocation Group -->
                    <div class="flex flex-col gap-4">
                      <CInputToggleCard
                        v-model="agent.invocation.system.enabled"
                        :label="$t('agent.editor.systemEnabled.label')"
                        :description="$t('agent.editor.systemEnabled.help')"
                      />

                      <div class="flex flex-col gap-1" :class="{ 'opacity-50 pointer-events-none': !agent.invocation.system.enabled }">
                        <label for="serviceAccount" class="font-medium text-primary">
                          {{ $t('agent.editor.serviceAccount.label') }}
                        </label>
                        <CInputUser
                          id="serviceAccount"
                          v-model="agent.invocation.system.serviceAccount"
                          :disabled="!agent.invocation.system.enabled"
                        />
                        <small class="text-muted-color">
                          {{ $t('agent.editor.serviceAccount.help') }}
                        </small>
                      </div>
                    </div>
                  </div>
                </Panel>
              </TabPanel>

              <TabPanel value="exec" class="h-full p-0 flex flex-row overflow-hidden">
                <AiChat :key="activeConvIndex" :agent="agent" :conversation="activeConversation" :isCreate="isCreate">
                  <template #header>
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
                        severity="secondary"
                        size="small"
                        @click="addConversation"
                        :disabled="isCreate"
                        class="shrink-0"
                      />
                    </div>
                  </template>
                </AiChat>
              </TabPanel>

              <TabPanel value="history" class="h-full p-4 flex flex-col gap-4 overflow-auto">
                <CResourceList
                  :items="agentConversations"
                  :fields="historyTableFields"
                  :loading="loadingConversations"
                  primary-key="aiConversationID"
                  :pagination="pagingConversations"
                  :sorting="sortingConversations"
                  :filter="filterConversations"
                  hide-search
                  clickable
                  :translations="{
                    showingPagination: 'general.resourceList.pagination.showing',
                    singlePluralPagination: 'general.resourceList.pagination.single',
                    prevPagination: $t('general.resourceList.pagination.prev'),
                    nextPagination: $t('general.resourceList.pagination.next'),
                    recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
                    resourceSingle: $t('agent.editor.history.resourceSingle'),
                    resourcePlural: $t('agent.editor.history.resourcePlural'),
                  }"
                  @row-click="openHistoryChat($event)"
                  @page-change="onConversationsPageChange"
                  @sort="onConversationsSort"
                  class="w-full border border-surface rounded-xl"
                >
                  <template #empty>
                    {{ $t('agent.editor.history.empty') }}
                  </template>
                  <template #body-snippet="{ data }">
                    <span class="text-color truncate block max-w-[250px] xl:max-w-[400px]">
                      {{ getWarningSnippet(data.messages) }}
                    </span>
                  </template>
                  <template #body-createdAt="{ data }">
                    {{ data.createdAt ? new Date(data.createdAt).toLocaleString() : '' }}
                  </template>
                </CResourceList>
              </TabPanel>
            </TabPanels>
          </Tabs>
        </template>
      </Card>
    </div>

    <!-- History Dialog -->
    <Dialog
      v-model:visible="showHistoryDialog"
      :header="$t('agent.editor.tabs.history')"
      modal
      class="max-w-[70vw] w-full"
      :style="{ height: '80vh' }"
      :pt="{
        content: { class: 'p-0 flex flex-col h-full overflow-hidden' },
      }"
    >
      <div v-if="loadingHistoryChat" class="flex items-center justify-center p-4 h-full">
        <ProgressSpinner />
      </div>
      <AiChat
        v-else-if="selectedHistoryConversation"
        :agent="agent"
        :conversation="selectedHistoryConversation"
        :isCreate="false"
        :showTrace="false"
        :readonly="true"
      />
    </Dialog>

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

// Components (not globally registered)
import { components } from '@cortezaproject/corteza-vue-next'
import AiChat from '@/components/AiChat.vue'

const { CInputLLM, CInputModel, CInputDelete, CInputUser, CInputKnowledgeBase, CInputToggleCard, CInputRole, CInputTAQ, CInputWorkflow, CResourceList } =
  components

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')
const $AutomationAPI = inject('$AutomationAPI')
const agentStore = useAgentStore()

const loading = ref(false)
const saving = ref(false)
const submitted = ref(false)
const agent = ref(null)
const activeTab = ref('config')

// Conversation tabs
let convCounter = 1
const contextExpanded = ref(false)
const conversations = ref([
  { label: `Chat 1`, messages: [], conversationID: null, traceHistory: [], context: '' },
])
const activeConvIndex = ref(0)
const activeConversation = computed(() => conversations.value[activeConvIndex.value])

const isCreate = computed(() => !route.params.agentID)

const showHistoryDialog = ref(false)
const loadingHistoryChat = ref(false)
const selectedHistoryConversation = ref(null)

const historyTableFields = computed(() => [
  { key: 'snippet', header: t('agent.editor.history.columns.snippet') || 'First Message', style: 'width: 60%' },
  { key: 'createdAt', header: t('agent.editor.history.columns.createdAt') }
])

function getWarningSnippet(messages) {
  if (!messages || messages.length === 0) return '—'
  const firstUserMsg = messages.find(m => m.role === 'user')
  const content = firstUserMsg ? firstUserMsg.content : (messages[0].content || '—')
  return content.replace(/\n/g, ' ').substring(0, 80) + (content.length > 80 ? '...' : '')
}

async function openHistoryChat(event) {
  const rowData = event.data || event
  if (!rowData.aiConversationID) return
  showHistoryDialog.value = true
  loadingHistoryChat.value = true
  selectedHistoryConversation.value = null
  
  try {
    const res = await $SystemAPI.aiConversationRead({ aiConversationID: rowData.aiConversationID })
    if (res && res.messages) {
      if (!res.traceHistory) res.traceHistory = []
      selectedHistoryConversation.value = res
    } else {
      selectedHistoryConversation.value = { ...rowData, messages: rowData.messages || [], traceHistory: [] }
    }
  } catch (e) {
    if (e.message !== 'canceled') {
      $toast.toastDanger(t('agent.editor.history.loadError'))
      showHistoryDialog.value = false
    }
  } finally {
    loadingHistoryChat.value = false
  }
}

const statusOptions = computed(() => [
  { label: t('agent.list.status.active'), value: 'active' },
  { label: t('agent.list.status.inactive'), value: 'inactive' },
])

const availableTools = ref([])
const loadingTools = ref(false)
const toolPickerSelection = ref(null)
const selectedTools = ref([])

// TCL
const tclMasterList = ref(null)
const loadingTcl = ref(false)

const tclArticleOptions = computed(() => {
  if (!tclMasterList.value) return []
  const treaties = tclMasterList.value.treaties || []
  const groups = tclMasterList.value.articles || []
  return groups.map(g => {
    const treaty = treaties.find(t => t.id === g.treatyId)
    return {
      treatyLabel: treaty?.label || g.treatyId,
      items: g.items.map(a => ({
        id: a.id,
        label: a.label,
        hardwired: a.hardwired,
        defaultSelected: a.defaultSelected,
        interpretation: a.interpretation || '',
      })),
    }
  })
})

const tclSortedArticleGroups = computed(() => tclArticleOptions.value)

function toggleTclArticle(articleId, enabled) {
  if (enabled) {
    if (!agent.value.behavior.tclArticles.includes(articleId)) {
      agent.value.behavior.tclArticles.push(articleId)
    }
  } else {
    agent.value.behavior.tclArticles = agent.value.behavior.tclArticles.filter(
      id => id !== articleId,
    )
  }
}

async function fetchTclMasterList() {
  loadingTcl.value = true
  try {
    tclMasterList.value = await $SystemAPI.agentTclMasterList()
  } catch (e) {
    console.error('Failed to fetch TCL master list:', e)
    tclMasterList.value = null
  } finally {
    loadingTcl.value = false
  }
}
const emptyAgent = () => ({
  handle: '',
  status: 'active',
  meta: {
    short: '',
    description: '',
    sidebarRoles: [],
  },
  behavior: {
    systemPrompt: '',
    guardrails: [],
    treatyCLEnabled: true,
    tclTemperature: 5,
    tclArticles: [],
    injectSystemContext: true,
    knowledgeBases: [],
  },
  execution: {
    model: {
      llmProviderID: '0',
      model: '',
      temperature: 0.7,
    },
    limits: {
      maxIterations: 10,
      contextWindow: 10000,
      outputTokens: 500,
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
    taqs: [],
    workflows: [],
    allow: [],
  },
  invocation: {
    user: { enabled: true },
    system: {
      enabled: false,
      serviceAccount: '0',
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

  agent.value.meta.sidebarRoles = agent.value.meta.sidebarRoles || []
  agent.value.access.taqs = agent.value.access.taqs || []
  agent.value.access.workflows = agent.value.access.workflows || []

  initToolSelection()
}

async function loadAgent() {
  const agentID = route.params.agentID
  if (!agentID) {
    agent.value = emptyAgent()
    return
  }

  // Only show spinner on first load; subsequent switches update values in-place
  const isFirstLoad = !agent.value
  if (isFirstLoad) loading.value = true

  try {
    const res = await $SystemAPI.agentRead({ agentID })
    applyAgentData(res)
  } catch (err) {
    $toast.toastDanger(t('notification.agent.loadFailed'))
    router.push({ name: 'root' })
  } finally {
    if (isFirstLoad) loading.value = false
  }
}

async function handleSubmit() {
  submitted.value = true

  if (!agent.value.behavior.systemPrompt?.trim()) {
    activeTab.value = 'config'
    return
  }

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


// Guardrails helpers
function addConversation() {
  convCounter++
  conversations.value.push({
    label: `Chat ${convCounter}`,
    messages: [],
    conversationID: null,
    traceHistory: [],
    context: '',
  })
  activeConvIndex.value = conversations.value.length - 1
}

function closeConversation(idx) {
  conversations.value.splice(idx, 1)
  if (activeConvIndex.value >= conversations.value.length) {
    activeConvIndex.value = conversations.value.length - 1
  }
}

function addGuardrail() {
  agent.value.behavior.guardrails.push('')
}

function removeGuardrail(index) {
  agent.value.behavior.guardrails.splice(index, 1)
}

// Watch for agentID changes (handles both initial mount and sidebar navigation)
watch(
  () => route.params.agentID,
  () => {
    // Reset editor state for the new agent
    convCounter = 1
    conversations.value = [
      { label: `Chat 1`, messages: [], conversationID: null, traceHistory: [], context: '' },
    ]
    activeConvIndex.value = 0
    activeTab.value = 'config'

    loadAgent()
  },
  { immediate: true },
)

const agentConversations = ref([])
const loadingConversations = ref(false)

const filterConversations = ref({})
const sortingConversations = ref({
  sortBy: 'createdAt',
  sortDesc: true,
})
const pagingConversations = ref({
  limit: 50,
  pageCursor: '',
  nextPage: '',
  prevPage: '',
  page: 1,
})

async function loadConversations() {
  if (!agent.value || !agent.value.agentID) return
  loadingConversations.value = true
  try {
    const res = await $SystemAPI.aiConversationList({ 
      agentID: agent.value.agentID, 
      sort: `${sortingConversations.value.sortBy} ${sortingConversations.value.sortDesc ? 'DESC' : 'ASC'}`,
      limit: pagingConversations.value.limit,
      pageCursor: pagingConversations.value.pageCursor,
    })
    agentConversations.value = res.set || []
    
    pagingConversations.value = {
      ...pagingConversations.value,
      nextPage: res.filter?.nextPage || '',
      prevPage: res.filter?.prevPage || '',
    }
  } catch (e) {
    if (e.message !== 'canceled') {
      agentConversations.value = []
      $toast.toastDanger(t('agent.editor.history.loadError'))
    }
  } finally {
    loadingConversations.value = false
  }
}

function onConversationsPageChange(event) {
  pagingConversations.value = {
    ...pagingConversations.value,
    ...event,
  }
  loadConversations()
}

function onConversationsSort(event) {
  sortingConversations.value = {
    sortBy: event.sortField,
    sortDesc: event.sortOrder < 0,
  }
  loadConversations()
}


watch(activeTab, (val) => {
  if (val === 'history') {
    loadConversations()
  }
})

onMounted(() => {
  fetchAvailableTools()
  fetchTclMasterList()
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

  if (isCreate.value) {
    // New agent: enable all tools by default (hints empty)
    agent.value.access.tools = availableTools.value.map(t => ({
      name: t.name,
      hints: '',
    }))
    selectedTools.value = []
  } else {
    // Existing agent: select only the saved tools
    const enabledNames = new Set((agent.value?.access?.tools || []).map(t => t.name))
    selectedTools.value = availableTools.value.filter(t => enabledNames.has(t.name))
  }
}

const unselectedTools = computed(() => {
  const selectedNames = new Set(selectedTools.value.map(t => t.name))
  return availableTools.value.filter(t => !selectedNames.has(t.name))
})

function onToolPickerSelect(tool) {
  if (!tool) return
  selectedTools.value = [...selectedTools.value, tool]
  // Add to agent.access.tools with empty hints
  agent.value.access.tools.push({ name: tool.name, hints: '' })
  // Clear in nextTick so the Select component sees the v-model change
  // after the current update cycle completes
  nextTick(() => {
    toolPickerSelection.value = null
  })
}

const loadedTaqNames = ref({})
const loadingTaqNames = ref({})

async function resolveTaqName(id) {
  if (!id || loadedTaqNames.value[id] || loadingTaqNames.value[id]) return
  loadingTaqNames.value[id] = true
  try {
    const res = await $AutomationAPI.ngAutomationRead({ automationID: id })
    loadedTaqNames.value[id] = res.meta?.short || res.handle || res.automationID
  } catch (e) {
    loadedTaqNames.value[id] = 'Unknown TAQ'
  }
}

watch(() => agent.value?.access?.taqs, (taqs) => {
  if (!taqs) return
  taqs.forEach(t => resolveTaqName(t.id))
}, { deep: true, immediate: true })

const loadedWorkflowNames = ref({})
const loadingWorkflowNames = ref({})

async function resolveWorkflowName(id) {
  if (!id || loadedWorkflowNames.value[id] || loadingWorkflowNames.value[id]) return
  loadingWorkflowNames.value[id] = true
  try {
    const res = await $AutomationAPI.workflowRead({ workflowID: id })
    loadedWorkflowNames.value[id] = res.meta?.name || res.handle || res.workflowID
  } catch (e) {
    loadedWorkflowNames.value[id] = 'Unknown Workflow'
  }
}

watch(() => agent.value?.access?.workflows, (workflows) => {
  if (!workflows) return
  workflows.forEach(w => resolveWorkflowName(w.id))
}, { deep: true, immediate: true })

function removeTool(tool) {
  selectedTools.value = selectedTools.value.filter(t => t.name !== tool.name)
  agent.value.access.tools = agent.value.access.tools.filter(t => t.name !== tool.name)
}

const taqPickerSelection = ref(null)
const workflowPickerSelection = ref(null)

function onTaqPickerSelect(id) {
  if (!id) return
  if (!agent.value.access.taqs) {
    agent.value.access.taqs = []
  }
  if (!agent.value.access.taqs.some(t => t.id === id)) {
    agent.value.access.taqs.push({ id, hints: '' })
  }
  nextTick(() => {
    taqPickerSelection.value = null
  })
}

function removeTaq(idx) {
  agent.value.access.taqs.splice(idx, 1)
}

function onWorkflowPickerSelect(id) {
  if (!id) return
  if (!agent.value.access.workflows) {
    agent.value.access.workflows = []
  }
  if (!agent.value.access.workflows.some(w => w.id === id)) {
    agent.value.access.workflows.push({ id, hints: '' })
  }
  nextTick(() => {
    workflowPickerSelection.value = null
  })
}

function removeWorkflow(idx) {
  agent.value.access.workflows.splice(idx, 1)
}

function getToolHints(name) {
  const tool = agent.value.access.tools.find(t => t.name === name)
  return tool?.hints || ''
}

function setToolHints(name, value) {
  const tool = agent.value.access.tools.find(t => t.name === name)
  if (tool) tool.hints = value
}
</script>
