<template>
  <Teleport to="#topbar-title" defer>
    <span v-if="isCreate">{{ $t('agent.editor.titleCreate') }}</span>
    <span v-else class="font-semibold">
      {{ agent?.meta?.short || agent?.handle || $t('agent.editor.titleEdit') }}
    </span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <Form
    v-else-if="agent"
    ref="formRef"
    v-slot="$form"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full overflow-hidden"
  >
    <div class="flex-1 flex flex-row min-h-0 overflow-hidden">
      <div class="flex-1 min-w-0 p-4 pr-2 flex flex-col gap-4 overflow-hidden">
        <Card
          :pt="{
            body: { class: 'p-0 h-full flex flex-col' },
            content: { class: 'p-0 h-full flex flex-col min-h-0' },
          }"
          class="flex-1 min-h-0 overflow-hidden"
        >
          <template #content>
            <Tabs v-model:value="activeTab" class="flex flex-col h-full min-h-0">
              <div class="flex items-center shrink-0 rounded-t-lg">
                <TabList class="flex-1 min-w-0 rounded-t-lg">
                  <Tab value="config">{{ $t('agent.editor.tabs.config') }}</Tab>
                  <Tab value="exec">{{ $t('agent.editor.tabs.exec') }}</Tab>
                  <Tab value="history" v-if="!isCreate">{{ $t('agent.editor.tabs.history') }}</Tab>
                </TabList>
                <div
                  v-if="!isCreate"
                  class="shrink-0 px-2 border-b border-surface self-stretch flex items-center"
                >
                  <CPermissionsButton
                    v-tooltip.bottom="$t('general.label.permissions')"
                    :resource="`corteza::system:agent/${agent.agentID}`"
                    :title="agent.meta?.short || agent.handle || agent.agentID"
                    :target="agent.meta?.short || agent.handle || agent.agentID"
                  />
                </div>
              </div>

              <TabPanels class="flex-1 min-h-0 p-0">
                <TabPanel value="config" class="h-full p-4 flex flex-col gap-4 overflow-auto">
                  <Panel :header="$t('agent.editor.panels.general')" toggleable>
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div class="flex flex-col gap-1">
                        <label for="name" class="font-medium text-primary">
                          {{ $t('agent.editor.name.label') }}
                        </label>
                        <small class="text-muted-color">{{ $t('agent.editor.name.help') }}</small>
                        <InputText id="name" v-model="agent.meta.short" />
                      </div>
                      <div class="flex flex-col gap-1">
                        <label for="handle" class="font-medium text-primary">
                          {{ $t('agent.editor.handle.label') }}
                        </label>
                        <small class="text-muted-color">{{ $t('agent.editor.handle.help') }}</small>
                        <InputText id="handle" v-model="agent.handle" />
                      </div>
                      <div class="flex flex-col gap-1 md:col-span-2">
                        <label for="description" class="font-medium text-primary">
                          {{ $t('agent.editor.description.label') }}
                        </label>
                        <small class="text-muted-color">
                          {{ $t('agent.editor.description.help') }}
                        </small>
                        <Textarea
                          id="description"
                          v-model="agent.meta.description"
                          rows="3"
                          autoResize
                        />
                      </div>
                      <div class="flex flex-col gap-1">
                        <label for="status" class="font-medium text-primary">
                          {{ $t('agent.editor.status.label') }}
                        </label>
                        <small class="text-muted-color">{{ $t('agent.editor.status.help') }}</small>
                        <Select
                          id="status"
                          v-model="agent.status"
                          :options="statusOptions"
                          optionLabel="label"
                          optionValue="value"
                        />
                      </div>
                      <div class="flex flex-col gap-1">
                        <label class="font-medium text-primary">
                          {{ $t('agent.editor.labels.label') }}
                        </label>
                        <small class="text-muted-color">{{ $t('agent.editor.labels.help') }}</small>
                        <CInputLabel
                          v-model="agent.labels"
                          :placeholder="$t('agent.editor.labels.placeholder')"
                          :create-label="$t('agent.editor.labels.createNew')"
                          :create-dialog-label="$t('agent.editor.labels.dialogCreate')"
                          :name-label="$t('agent.editor.labels.name')"
                          :save-btn-label="$t('general.label.save')"
                          :cancel-btn-label="$t('general.label.cancel')"
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
                        <small class="text-muted-color">
                          {{ $t('agent.editor.provider.help') }}
                        </small>
                        <CInputLLM id="provider" v-model="agent.execution.model.llmProviderID" />
                      </div>
                      <div class="flex flex-col gap-1">
                        <label for="model" class="font-medium text-primary">
                          {{ $t('agent.editor.model.label') }}
                        </label>
                        <small class="text-muted-color">{{ $t('agent.editor.model.help') }}</small>
                        <CInputModel
                          id="model"
                          v-model="agent.execution.model.model"
                          :llmProviderID="agent.execution.model.llmProviderID"
                        />
                      </div>
                      <div class="flex flex-col gap-1">
                        <div class="flex items-center justify-between">
                          <label for="temperature" class="font-medium text-primary">
                            {{ $t('agent.editor.temperature.label') }}
                            <span v-if="temperatureEnabled">
                              ({{ agent.execution.model.temperature }})
                            </span>
                          </label>
                          <ToggleSwitch
                            :model-value="temperatureEnabled"
                            @update:model-value="toggleTemperature"
                          />
                        </div>
                        <small class="text-muted-color">
                          {{ $t('agent.editor.temperature.help') }}
                        </small>
                        <Slider
                          id="temperature"
                          v-model="agent.execution.model.temperature"
                          :disabled="!temperatureEnabled"
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
                            <small class="text-muted-color">
                              {{ $t('agent.editor.maxIterations.help') }}
                            </small>
                            <InputNumber
                              id="maxIterations"
                              v-model="agent.execution.limits.maxIterations"
                              mode="decimal"
                              :useGrouping="false"
                              :min="1"
                              :max="50"
                            />
                          </div>
                          <div class="flex flex-col gap-1">
                            <label for="contextWindow" class="font-medium text-primary">
                              {{ $t('agent.editor.contextWindow.label') }}
                            </label>
                            <small class="text-muted-color">
                              {{ $t('agent.editor.contextWindow.help') }}
                            </small>
                            <InputNumber
                              id="contextWindow"
                              v-model="agent.execution.limits.contextWindow"
                              mode="decimal"
                              :useGrouping="false"
                            />
                          </div>
                          <div class="flex flex-col gap-1">
                            <label for="outputTokens" class="font-medium text-primary">
                              {{ $t('agent.editor.outputTokens.label') }}
                            </label>
                            <small class="text-muted-color">
                              {{ $t('agent.editor.outputTokens.help') }}
                            </small>
                            <InputNumber
                              id="outputTokens"
                              v-model="agent.execution.limits.outputTokens"
                              mode="decimal"
                              :useGrouping="false"
                            />
                          </div>
                          <div class="flex flex-col gap-1">
                            <label for="timeout" class="font-medium text-primary">
                              {{ $t('agent.editor.timeout.label') }}
                            </label>
                            <small class="text-muted-color">
                              {{ $t('agent.editor.timeout.help') }}
                            </small>
                            <InputText
                              id="timeout"
                              v-model="agent.execution.limits.timeout"
                              placeholder="30s"
                            />
                          </div>
                        </div>
                      </div>
                    </div>
                  </Panel>

                  <Panel :header="$t('agent.editor.panels.behavior')" toggleable>
                    <div class="grid grid-cols-1 gap-4">
                      <FormField name="systemPrompt" class="flex flex-col gap-1">
                        <label for="systemPrompt" class="font-medium text-primary">
                          {{ $t('agent.editor.systemPrompt.label') }}
                          <span class="text-red-500">*</span>
                        </label>
                        <small class="text-muted-color">
                          {{ $t('agent.editor.systemPrompt.help') }}
                        </small>
                        <Textarea
                          id="systemPrompt"
                          name="systemPrompt"
                          v-model="agent.behavior.systemPrompt"
                          rows="6"
                          autoResize
                        />
                        <Message
                          v-if="$form.systemPrompt?.invalid"
                          severity="error"
                          size="small"
                          variant="simple"
                        >
                          {{ $form.systemPrompt.error?.message }}
                        </Message>
                      </FormField>

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
                            dimWhenOff
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
                          :modules-placeholder="
                            $t('agent.editor.knowledgeBases.modulesPlaceholder')
                          "
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
                          v-for="(tool, toolIdx) in selectedTools"
                          :key="tool.name"
                          class="flex items-center gap-3 p-3 border border-surface rounded-lg"
                        >
                          <div class="flex flex-col gap-0.5 flex-1 min-w-0">
                            <span class="font-medium text-color text-sm truncate">
                              {{ tool.title }}
                            </span>
                            <small
                              v-if="getToolHints(tool.name)"
                              class="text-muted-color text-xs truncate"
                            >
                              {{ getToolHints(tool.name) }}
                            </small>
                            <!-- Access info -->
                            <div
                              v-if="!hasToolAllowFromName(tool.name)"
                              class="flex items-center gap-1.5"
                            >
                              <i class="pi pi-lock text-xs text-muted-color" />
                              <span class="text-xs text-muted-color">
                                {{ $t('agent.editor.tools.restricted') }}
                              </span>
                            </div>
                            <div v-else class="flex items-center gap-1.5 flex-wrap mt-0.5">
                              <i class="pi pi-lock text-xs text-muted-color" />
                              <template
                                v-for="detail in getToolAllowDetails(tool.name)"
                                :key="detail.namespaceID"
                              >
                                <Tag severity="secondary" rounded>
                                  <template #default>
                                    <span class="text-xs">
                                      {{ detail.namespaceName }}
                                      <span v-if="detail.modules.length" class="text-muted-color">
                                        · {{ detail.modules.join(', ') }}
                                      </span>
                                      <span v-else class="text-muted-color">
                                        · {{ $t('agent.editor.tools.allModules') }}
                                      </span>
                                    </span>
                                  </template>
                                </Tag>
                              </template>
                            </div>
                          </div>

                          <div class="flex items-center gap-1 shrink-0">
                            <Button
                              icon="pi pi-pencil"
                              severity="secondary"
                              text
                              size="small"
                              @click="openToolDialog(toolIdx)"
                            />
                            <Button
                              icon="pi pi-trash"
                              severity="danger"
                              text
                              size="small"
                              @click="removeTool(tool)"
                            />
                          </div>
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
                              v-model="agent.access.taqs[idx].description"
                              class="w-full mt-1"
                              size="small"
                              :placeholder="$t('agent.editor.taqs.descriptionPlaceholder')"
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
                        <small class="text-muted-color">
                          {{ $t('agent.editor.workflows.help') }}
                        </small>
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
                              v-model="agent.access.workflows[idx].description"
                              class="w-full mt-1"
                              size="small"
                              :placeholder="$t('agent.editor.workflows.descriptionPlaceholder')"
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

                        <div
                          class="flex flex-col gap-1"
                          :class="{
                            'opacity-50 pointer-events-none': !agent.invocation.user.enabled,
                          }"
                        >
                          <label for="sidebarRoles" class="font-medium text-primary">
                            {{ $t('agent.editor.sidebarRoles.label') }}
                          </label>
                          <small class="text-muted-color">
                            {{ $t('agent.editor.sidebarRoles.help') }}
                          </small>
                          <CInputRole
                            id="sidebarRoles"
                            v-model="agent.meta.sidebarRoles"
                            :multiple="true"
                            :disabled="!agent.invocation.user.enabled"
                          />
                        </div>
                      </div>

                      <!-- System Invocation Group -->
                      <div class="flex flex-col gap-4">
                        <CInputToggleCard
                          v-model="agent.invocation.system.enabled"
                          :label="$t('agent.editor.systemEnabled.label')"
                          :description="$t('agent.editor.systemEnabled.help')"
                        />

                        <div
                          class="flex flex-col gap-1"
                          :class="{
                            'opacity-50 pointer-events-none': !agent.invocation.system.enabled,
                          }"
                        >
                          <label for="serviceAccount" class="font-medium text-primary">
                            {{ $t('agent.editor.serviceAccount.label') }}
                          </label>
                          <small class="text-muted-color">
                            {{ $t('agent.editor.serviceAccount.help') }}
                          </small>
                          <CInputUser
                            id="serviceAccount"
                            v-model="agent.invocation.system.serviceAccount"
                            :disabled="!agent.invocation.system.enabled"
                          />
                        </div>
                      </div>
                    </div>
                  </Panel>
                </TabPanel>

                <TabPanel value="exec" class="h-full p-0 overflow-hidden">
                  <AiTrace
                    :key="'trace-' + activeConvIndex"
                    :agent="agent"
                    :conversation="activeConversation"
                  />
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

      <!-- Resize handle -->
      <div v-if="showChat" class="resize-handle group" @mousedown="startChatResize">
        <div class="resize-handle-bar group-hover:opacity-100" />
      </div>

      <!-- Right: permanent chatbox -->
      <div
        v-if="showChat"
        class="shrink-0 flex flex-col p-4 pl-0 overflow-hidden"
        :style="{ width: chatWidth + 'px' }"
      >
        <Card
          :pt="{
            body: { class: 'p-0 h-full flex flex-col' },
            content: { class: 'p-0 h-full flex flex-col min-h-0' },
          }"
          class="flex-1 min-h-0 overflow-hidden"
        >
          <template #content>
            <div
              v-if="isCreate"
              class="flex items-center justify-center h-full p-6 text-muted-color text-sm text-center"
            >
              {{ $t('agent.editor.split.chatPlaceholder') }}
            </div>
            <AiChat
              v-else
              :key="activeConvIndex"
              :agent="agent"
              :conversation="activeConversation"
              :isCreate="isCreate"
              :showTrace="false"
            >
              <template #header>
                <div
                  class="flex items-center gap-0 border-b border-surface shrink-0 bg-surface-ground"
                >
                  <div class="flex items-center gap-0 flex-1 overflow-x-auto">
                    <button
                      v-for="(conv, idx) in conversations"
                      :key="idx"
                      type="button"
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
                    class="shrink-0 mr-1"
                  />
                </div>
              </template>
            </AiChat>
          </template>
        </Card>
      </div>
    </div>

    <!-- Tool Configuration Dialog -->
    <Dialog
      v-model:visible="toolDialogVisible"
      :header="
        editingToolMeta?.title || editingToolForm?.name || $t('agent.editor.tools.configure')
      "
      modal
      :style="{ width: '40rem' }"
    >
      <div v-if="editingToolForm" class="flex flex-col gap-4">
        <!-- Hints -->
        <div class="flex flex-col gap-1">
          <label class="font-medium text-primary text-sm">
            {{ $t('agent.editor.tools.toolDescription') }}
          </label>
          <InputText
            v-model="editingToolForm.hints"
            :placeholder="$t('agent.editor.tools.descriptionPlaceholder')"
          />
        </div>

        <!-- Restrict access toggle -->
        <div class="flex flex-col gap-3">
          <CInputToggleCard
            :modelValue="hasToolAllow(editingToolForm)"
            :label="$t('agent.editor.tools.configureAccess')"
            :description="$t('agent.editor.tools.configureAccessHelp')"
            @update:modelValue="toggleToolAllow(editingToolForm, $event)"
          />

          <!-- Namespace / Module rows -->
          <template v-if="hasToolAllow(editingToolForm)">
            <div class="flex flex-col gap-3">
              <div
                v-for="(rule, ruleIdx) in editingToolForm.allow"
                :key="ruleIdx"
                class="flex items-start gap-2 border border-surface rounded-lg p-3"
              >
                <div class="flex flex-col gap-2 flex-1">
                  <CInputNamespace
                    :model-value="rule.namespaceID"
                    @update:model-value="onToolAllowNamespaceChange(rule, $event)"
                    :placeholder="$t('agent.editor.tools.namespacePlaceholder')"
                  />
                  <CInputModule
                    v-if="rule.namespaceID"
                    :model-value="rule.moduleIDs || []"
                    @update:model-value="rule.moduleIDs = $event"
                    :namespace-i-d="rule.namespaceID"
                    :placeholder="$t('agent.editor.tools.modulesPlaceholder')"
                    :multiple="true"
                  />
                </div>
                <Button
                  icon="pi pi-trash"
                  severity="danger"
                  text
                  size="small"
                  @click="removeToolAllowEntry(editingToolForm, ruleIdx)"
                />
              </div>

              <Button
                :label="$t('agent.editor.tools.addNamespace')"
                icon="pi pi-plus"
                severity="secondary"
                outlined
                size="small"
                @click="addToolAllowEntry(editingToolForm)"
              />
            </div>
          </template>
        </div>
      </div>

      <template #footer>
        <div class="flex justify-end gap-2">
          <Button
            :label="$t('general.label.cancel')"
            severity="secondary"
            text
            size="small"
            @click="toolDialogVisible = false"
          />
          <Button
            :label="$t('general.label.save')"
            severity="primary"
            size="small"
            @click="saveToolDialog"
          />
        </div>
      </template>
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
            v-if="!isCreate && agent.canDeleteAgent"
            :label="$t('general.label.delete')"
            :message="$t('agent.list.delete')"
            :header="agent?.meta?.short || agent?.handle || $t('general.label.delete')"
            @confirm="handleDelete"
          />
          <Button
            type="button"
            :label="$t('general.label.save')"
            icon="pi pi-save"
            :loading="saving"
            @click="submitForm"
          />
        </div>
      </div>
    </div>
  </Form>
</template>

<script setup>
import { computed, inject, nextTick, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { cloneDeep, isEqual } from 'lodash-es'

import { useAgentStore } from '@/stores/agent'
import { useRoute, useRouter } from 'vue-router'
import { useComposeResourceStore, useUnsavedGuard } from '@planetcrust/human-vue'
import { system } from '@planetcrust/human-js'

// Components (not globally registered)
import { components } from '@planetcrust/human-vue'
import AiChat from '@/components/AiChat.vue'
import AiTrace from '@/components/AiTrace.vue'
import { useEditorSplit } from '@/composables/useEditorSplit'

const {
  CInputLLM,
  CInputModel,
  CInputDelete,
  CInputUser,
  CInputKnowledgeBase,
  CInputToggleCard,
  CInputRole,
  CInputTAQ,
  CInputWorkflow,
  CResourceList,
  CInputNamespace,
  CInputModule,
  CInputLabel,
} = components

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')
const $AutomationAPI = inject('$AutomationAPI')
const agentStore = useAgentStore()
const composeStore = useComposeResourceStore()

const loading = ref(false)
const saving = ref(false)
const agent = ref(null)
const initialAgent = ref(null)
const activeTab = ref('config')
const formRef = ref(null)

function submitForm() {
  formRef.value?.submit?.()
}

const { chatWidth, showChat, startChatResize } = useEditorSplit()

const { markSaved } = useUnsavedGuard({
  isDirty: () =>
    !saving.value &&
    !!agent.value &&
    !!initialAgent.value &&
    !isEqual(agent.value, initialAgent.value),
  messageKey: 'general.editor.unsavedChanges',
})

// Conversation tabs
let convCounter = 1
const conversations = ref([
  { label: `Chat 1`, messages: [], conversationID: null, traceHistory: [], context: '' },
])
const activeConvIndex = ref(0)
const activeConversation = computed(() => conversations.value[activeConvIndex.value])

watch(
  () => activeConversation.value?.traceHistory?.length || 0,
  (len, prev) => {
    if (len > 0 && (prev || 0) === 0) activeTab.value = 'exec'
  },
)

const isCreate = computed(() => !route.params.agentID)

const temperatureEnabled = computed(
  () =>
    agent.value?.execution?.model?.temperature !== null &&
    agent.value?.execution?.model?.temperature !== undefined,
)

const historyTableFields = computed(() => [
  {
    key: 'snippet',
    header: t('agent.editor.history.columns.snippet') || 'First Message',
    style: 'width: 60%',
  },
  { key: 'createdAt', header: t('agent.editor.history.columns.createdAt') },
])

function getWarningSnippet(messages) {
  if (!messages || messages.length === 0) return '—'
  const firstUserMsg = messages.find(m => m.role === 'user')
  const content = firstUserMsg ? firstUserMsg.content : messages[0].content || '—'
  return content.replace(/\n/g, ' ').substring(0, 80) + (content.length > 80 ? '...' : '')
}

async function openHistoryChat(event) {
  const rowData = event.data || event
  if (!rowData.aiConversationID) return

  try {
    const res = await $SystemAPI.aiConversationRead({ aiConversationID: rowData.aiConversationID })
    const messages = res?.messages || rowData.messages || []
    const convID = res?.aiConversationID || rowData.aiConversationID
    const existingIdx = conversations.value.findIndex(
      c => (c.conversationID || c.aiConversationID) === convID,
    )
    if (existingIdx >= 0) {
      activeConvIndex.value = existingIdx
      return
    }
    const snippet = getWarningSnippet(messages).slice(0, 24)
    conversations.value.push({
      label: snippet || 'History',
      messages,
      conversationID: convID,
      aiConversationID: convID,
      traceHistory: res?.traceHistory || [],
      context: res?.context || '',
    })
    activeConvIndex.value = conversations.value.length - 1
  } catch (e) {
    if (e.message !== 'canceled') {
      $toast.toastDanger(t('agent.editor.history.loadError'))
    }
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

// Tool configuration dialog
const toolDialogVisible = ref(false)
const editingToolIndex = ref(-1)
const editingToolForm = ref(null)
const editingToolMeta = computed(() => {
  if (!editingToolForm.value) return null
  return availableTools.value.find(t => t.name === editingToolForm.value.name)
})

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

function toggleTemperature(enabled) {
  agent.value.execution.model.temperature = enabled ? 0.7 : null
}

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

function defaultTclArticleIds() {
  const ids = []
  for (const group of tclArticleOptions.value) {
    for (const a of group.items) {
      if (a.defaultSelected || a.hardwired) ids.push(a.id)
    }
  }
  return ids
}

function applyTclArticleDefaults() {
  const defaults = defaultTclArticleIds()
  const current = agent.value.behavior.tclArticles || []
  const merged = Array.from(new Set([...current, ...defaults]))
  agent.value.behavior.tclArticles = merged
}

watch(
  () => agent.value?.behavior?.treatyCLEnabled,
  (enabled, prev) => {
    if (enabled && !prev) applyTclArticleDefaults()
  },
)

watch(tclArticleOptions, () => {
  if (agent.value?.behavior?.treatyCLEnabled) applyTclArticleDefaults()
})

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
function applyAgentData(res) {
  agent.value = new system.Agent(res)
  initialAgent.value = cloneDeep(agent.value)
  initToolSelection()
}

const initialValues = computed(() => ({
  systemPrompt: agent.value?.behavior?.systemPrompt || '',
}))

const resolver = ref(({ values }) => {
  const errors = {}
  if (!values.systemPrompt || values.systemPrompt.trim().length === 0) {
    errors.systemPrompt = [{ message: t('agent.editor.systemPrompt.required') }]
  }
  return { errors }
})

async function loadAgent() {
  const agentID = route.params.agentID
  if (!agentID) {
    agent.value = new system.Agent()
    initialAgent.value = cloneDeep(agent.value)
    return
  }

  loading.value = true

  try {
    const res = await $SystemAPI.agentRead({ agentID })
    applyAgentData(res)
  } catch {
    $toast.toastDanger(t('notification.agent.loadFailed'))
    router.push({ name: 'root' })
  } finally {
    loading.value = false
  }
}

async function handleSubmit({ valid }) {
  if (!valid) {
    activeTab.value = 'config'
    $toast.toastWarning(t('general.notification.formErrors'))
    nextTick(() => {
      document
        .querySelector('.p-message-error')
        ?.scrollIntoView({ behavior: 'smooth', block: 'center' })
    })
    return
  }

  saving.value = true
  try {
    // Clean internal helper properties from tool allow entries before saving
    const payload = JSON.parse(
      JSON.stringify(agent.value, (key, value) => {
        if (key === '_moduleOptions' || key === '_loadingModules') return undefined
        return value
      }),
    )

    if (isCreate.value) {
      const created = await $SystemAPI.agentCreate(payload)
      agentStore.updateInList(created)
      $toast.toastSuccess(t('notification.agent.created'))
      markSaved()
      router.push({ name: 'agent.edit', params: { agentID: created.agentID } })
    } else {
      const updated = await $SystemAPI.agentUpdate({
        agentID: route.params.agentID,
        ...payload,
      })
      agentStore.updateInList(updated)
      applyAgentData(updated)
      $toast.toastSuccess(t('notification.agent.saved'))
    }
  } catch (err) {
    console.error(err)
    $toast.toastErrorHandler(t('notification.agent.saveFailed'))(err)
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  try {
    await $SystemAPI.agentDelete({ agentID: route.params.agentID })
    agentStore.removeFromList(route.params.agentID)
    $toast.toastSuccess(t('notification.agent.deleted'))
    initialAgent.value = cloneDeep(agent.value)
    router.push({ name: 'root' })
  } catch (err) {
    console.error(err)
    $toast.toastErrorHandler(t('notification.agent.deleteFailed'))(err)
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
    agentConversations.value = []
    filterConversations.value = {}
    sortingConversations.value = { sortBy: 'createdAt', sortDesc: true }
    editingToolForm.value = null
    editingToolIndex.value = -1
    toolDialogVisible.value = false

    loadAgent()
  },
  { immediate: true },
)

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

watch(activeTab, val => {
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

  const enabledNames = new Set((agent.value?.access?.tools || []).map(t => t.name))
  selectedTools.value = availableTools.value.filter(t => enabledNames.has(t.name))
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
  } catch {
    loadedTaqNames.value[id] = 'Unknown TAQ'
  }
}

watch(
  () => agent.value?.access?.taqs,
  taqs => {
    if (!taqs) return
    taqs.forEach(t => resolveTaqName(t.id))
  },
  { deep: true, immediate: true },
)

const loadedWorkflowNames = ref({})
const loadingWorkflowNames = ref({})

async function resolveWorkflowName(id) {
  if (!id || loadedWorkflowNames.value[id] || loadingWorkflowNames.value[id]) return
  loadingWorkflowNames.value[id] = true
  try {
    const res = await $AutomationAPI.workflowRead({ workflowID: id })
    loadedWorkflowNames.value[id] = res.meta?.name || res.handle || res.workflowID
  } catch {
    loadedWorkflowNames.value[id] = 'Unknown Workflow'
  }
}

watch(
  () => agent.value?.access?.workflows,
  workflows => {
    if (!workflows) return
    workflows.forEach(w => resolveWorkflowName(w.id))
  },
  { deep: true, immediate: true },
)

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
    agent.value.access.taqs.push({ id, description: '' })
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
    agent.value.access.workflows.push({ id, description: '' })
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

// --- Tool configuration dialog helpers ---

function openToolDialog(toolIdx) {
  const tool = selectedTools.value[toolIdx]
  if (!tool) return
  const accessIdx = agent.value.access.tools.findIndex(t => t.name === tool.name)
  if (accessIdx < 0) return

  // Deep copy so edits don't leak until Save
  editingToolIndex.value = accessIdx
  editingToolForm.value = JSON.parse(JSON.stringify(agent.value.access.tools[accessIdx]))
  toolDialogVisible.value = true
}

function saveToolDialog() {
  if (editingToolIndex.value < 0 || !editingToolForm.value) return
  agent.value.access.tools[editingToolIndex.value] = editingToolForm.value
  toolDialogVisible.value = false
}

// Resolved namespace and module names for tool allow summaries
const resolvedNsNames = ref({})
const resolvedModNames = ref({})

function resolveToolAllowResources() {
  const tools = agent.value?.access?.tools || []
  for (const tool of tools) {
    if (!tool.allow) continue
    for (const rule of tool.allow) {
      if (rule.namespaceID && !resolvedNsNames.value[rule.namespaceID]) {
        composeStore
          .resolveNamespace(String(rule.namespaceID))
          .then(ns => {
            if (ns) {
              resolvedNsNames.value = {
                ...resolvedNsNames.value,
                [rule.namespaceID]: ns.name || ns.slug || rule.namespaceID,
              }
            }
          })
          .catch(() => {})
      }
      if (rule.namespaceID && rule.moduleIDs?.length) {
        for (const modID of rule.moduleIDs) {
          if (!resolvedModNames.value[modID]) {
            composeStore
              .resolveModule(String(rule.namespaceID), String(modID))
              .then(mod => {
                if (mod) {
                  resolvedModNames.value = {
                    ...resolvedModNames.value,
                    [modID]: mod.name || mod.handle || modID,
                  }
                }
              })
              .catch(() => {})
          }
        }
      }
    }
  }
}

watch(
  () => agent.value?.access?.tools,
  () => {
    resolveToolAllowResources()
  },
  { deep: true, immediate: true },
)

function hasToolAllowFromName(toolName) {
  const tool = agent.value.access.tools.find(t => t.name === toolName)
  return Array.isArray(tool?.allow)
}

function getToolAllowDetails(toolName) {
  const tool = agent.value.access.tools.find(t => t.name === toolName)
  if (!tool?.allow?.length) return []

  return tool.allow.map(r => ({
    namespaceID: r.namespaceID,
    namespaceName: resolvedNsNames.value[r.namespaceID] || r.namespaceID,
    modules: (r.moduleIDs || []).map(id => resolvedModNames.value[id] || id),
  }))
}

function hasToolAllow(toolData) {
  return Array.isArray(toolData?.allow)
}

function toggleToolAllow(toolData, enabled) {
  if (!toolData) return
  if (enabled) {
    toolData.allow = []
  } else {
    toolData.allow = null
  }
}

function addToolAllowEntry(toolData) {
  if (!toolData || !Array.isArray(toolData.allow)) return
  toolData.allow.push({
    namespaceID: null,
    moduleIDs: [],
  })
}

function removeToolAllowEntry(toolData, idx) {
  if (!toolData?.allow) return
  toolData.allow.splice(idx, 1)
}

function onToolAllowNamespaceChange(rule, namespaceID) {
  rule.namespaceID = namespaceID
  rule.moduleIDs = []
}
</script>

<style scoped>
.resize-handle {
  width: 8px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: ew-resize;
}

.resize-handle-bar {
  width: 2px;
  height: 2rem;
  border-radius: 9999px;
  background-color: var(--p-content-border-color);
  opacity: 0.4;
  transition: opacity 150ms;
}
</style>
