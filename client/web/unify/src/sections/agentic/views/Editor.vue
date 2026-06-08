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
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full overflow-hidden"
  >
    <div class="flex-1 flex flex-row min-h-0 overflow-hidden">
      <div
        class="flex-1 min-w-0 p-4 pr-2 flex flex-col gap-4 overflow-hidden w-full max-w-screen-2xl mx-auto"
      >
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
                      <CFormGroup
                        :label="$t('agent.editor.name.label')"
                        :description="$t('agent.editor.name.help')"
                        input-id="name"
                      >
                        <InputText id="name" v-model="agent.meta.short" />
                      </CFormGroup>
                      <CFormGroup
                        :label="$t('agent.editor.handle.label')"
                        :description="$t('agent.editor.handle.help')"
                        input-id="handle"
                      >
                        <InputText id="handle" v-model="agent.handle" />
                      </CFormGroup>
                      <CFormGroup
                        class="md:col-span-2"
                        :label="$t('agent.editor.description.label')"
                        :description="$t('agent.editor.description.help')"
                        input-id="description"
                      >
                        <Textarea
                          id="description"
                          v-model="agent.meta.description"
                          rows="3"
                          autoResize
                        />
                      </CFormGroup>
                      <CFormGroup
                        :label="$t('agent.editor.status.label')"
                        :description="$t('agent.editor.status.help')"
                        input-id="status"
                      >
                        <Select
                          id="status"
                          v-model="agent.status"
                          :options="statusOptions"
                          optionLabel="label"
                          optionValue="value"
                        />
                      </CFormGroup>
                      <CFormGroup
                        :label="$t('agent.editor.labels.label')"
                        :description="$t('agent.editor.labels.help')"
                      >
                        <CInputLabel
                          v-model="agent.labels"
                          :placeholder="$t('agent.editor.labels.placeholder')"
                          :create-label="$t('agent.editor.labels.createNew')"
                          :create-dialog-label="$t('agent.editor.labels.dialogCreate')"
                          :name-label="$t('agent.editor.labels.name')"
                          :save-btn-label="$t('general.label.save')"
                          :cancel-btn-label="$t('general.label.cancel')"
                        />
                      </CFormGroup>
                    </div>
                  </Panel>

                  <Panel :header="$t('agent.editor.panels.execution')" toggleable>
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <CFormGroup
                        :label="$t('agent.editor.provider.label')"
                        :description="$t('agent.editor.provider.help')"
                        input-id="provider"
                      >
                        <CInputLLM id="provider" v-model="agent.execution.model.llmProviderID" />
                      </CFormGroup>
                      <CFormGroup
                        :label="$t('agent.editor.model.label')"
                        :description="$t('agent.editor.model.help')"
                        input-id="model"
                      >
                        <CInputModel
                          id="model"
                          v-model="agent.execution.model.model"
                          :llmProviderID="agent.execution.model.llmProviderID"
                        />
                      </CFormGroup>
                      <CInputToggleCard
                        :model-value="temperatureEnabled"
                        @update:model-value="toggleTemperature"
                      >
                        <template #label>
                          {{ $t('agent.editor.temperature.label') }}
                          <span v-if="temperatureEnabled">
                            ({{ agent.execution.model.temperature }})
                          </span>
                        </template>
                        <template #description>
                          {{ $t('agent.editor.temperature.help') }}
                          <Slider
                            v-if="temperatureEnabled"
                            id="temperature"
                            v-model="agent.execution.model.temperature"
                            :min="0"
                            :max="1"
                            :step="0.1"
                            class="w-full mt-3"
                            @click.stop.prevent
                          />
                        </template>
                      </CInputToggleCard>

                      <!-- Execution limits -->
                      <div class="md:col-span-2">
                        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                          <CFormGroup
                            :label="$t('agent.editor.maxIterations.label')"
                            :description="$t('agent.editor.maxIterations.help')"
                            input-id="maxIterations"
                          >
                            <InputNumber
                              id="maxIterations"
                              v-model="agent.execution.limits.maxIterations"
                              mode="decimal"
                              :useGrouping="false"
                              :min="1"
                              :max="50"
                            />
                          </CFormGroup>
                          <CFormGroup
                            :label="$t('agent.editor.contextWindow.label')"
                            :description="$t('agent.editor.contextWindow.help')"
                            input-id="contextWindow"
                          >
                            <InputNumber
                              id="contextWindow"
                              v-model="agent.execution.limits.contextWindow"
                              mode="decimal"
                              :useGrouping="false"
                            />
                          </CFormGroup>
                          <CFormGroup
                            :label="$t('agent.editor.outputTokens.label')"
                            :description="$t('agent.editor.outputTokens.help')"
                            input-id="outputTokens"
                          >
                            <InputNumber
                              id="outputTokens"
                              v-model="agent.execution.limits.outputTokens"
                              mode="decimal"
                              :useGrouping="false"
                            />
                          </CFormGroup>
                          <CFormGroup
                            :label="$t('agent.editor.timeout.label')"
                            :description="$t('agent.editor.timeout.help')"
                            input-id="timeout"
                          >
                            <InputText
                              id="timeout"
                              v-model="agent.execution.limits.timeout"
                              placeholder="30s"
                            />
                          </CFormGroup>
                        </div>
                      </div>
                    </div>
                  </Panel>

                  <Panel :header="$t('agent.editor.panels.behavior')" toggleable>
                    <div class="grid grid-cols-1 gap-4">
                      <CFormGroup
                        name="systemPrompt"
                        :label="$t('agent.editor.systemPrompt.label')"
                        :description="$t('agent.editor.systemPrompt.help')"
                        required
                      >
                        <Textarea
                          id="systemPrompt"
                          name="systemPrompt"
                          v-model="agent.behavior.systemPrompt"
                          rows="6"
                          autoResize
                        />
                      </CFormGroup>

                      <!-- Guardrails -->
                      <CFormGroup
                        :label="$t('agent.editor.guardrails.label')"
                        :description="$t('agent.editor.guardrails.help')"
                      >
                        <template #actions>
                          <Button
                            :label="$t('agent.editor.guardrails.add')"
                            icon="pi pi-plus"
                            severity="secondary"
                            size="small"
                            @click="addGuardrail"
                          />
                        </template>
                        <CFormList v-model="agent.behavior.guardrails">
                          <template #row="{ index }">
                            <InputText
                              v-model="agent.behavior.guardrails[index]"
                              size="small"
                              class="w-full"
                              :placeholder="$t('agent.editor.guardrails.placeholder')"
                            />
                          </template>
                        </CFormList>
                      </CFormGroup>
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
                        <CFormGroup :description="$t('agent.editor.tcl.temperature.help')">
                          <template #label>
                            {{ $t('agent.editor.tcl.temperature.label') }} ({{
                              agent.behavior.tclTemperature
                            }})
                          </template>
                          <Slider
                            v-model="agent.behavior.tclTemperature"
                            :min="1"
                            :max="10"
                            :step="1"
                            class="w-full mt-2"
                          />
                        </CFormGroup>

                        <CFormGroup
                          :label="$t('agent.editor.tcl.articles.label')"
                          :description="$t('agent.editor.tcl.articles.help')"
                        />

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
                      <CFormGroup
                        :label="$t('agent.editor.knowledgeBases.label')"
                        :description="$t('agent.editor.knowledgeBases.help')"
                      >
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
                      </CFormGroup>

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
                      <CFormGroup
                        :label="$t('agent.editor.tools.label')"
                        :description="$t('agent.editor.tools.help')"
                      />

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

                      <CFormItemList
                        :items="selectedTools"
                        item-key="name"
                        @remove="tool => removeTool(tool)"
                      >
                        <template #default="{ item: tool }">
                          <span class="font-medium text-color text-sm truncate">
                            {{ tool.title }}
                          </span>
                          <small
                            v-if="getToolHints(tool.name)"
                            class="text-muted-color text-xs truncate block"
                          >
                            {{ getToolHints(tool.name) }}
                          </small>
                          <div
                            v-if="!hasToolAllowFromName(tool.name)"
                            class="flex items-center gap-1.5 mt-0.5"
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
                        </template>
                        <template #actions="{ index }">
                          <Button
                            icon="pi pi-pencil"
                            severity="secondary"
                            text
                            size="small"
                            @click="openToolDialog(index)"
                          />
                        </template>
                      </CFormItemList>

                      <Divider class="my-4" />

                      <!-- TAQs -->
                      <CFormGroup
                        :label="$t('agent.editor.taqs.label')"
                        :description="$t('agent.editor.taqs.help')"
                      />

                      <CInputTAQ
                        v-model="taqPickerSelection"
                        :placeholder="$t('agent.editor.taqs.selectPlaceholder')"
                        @update:model-value="onTaqPickerSelect"
                      />

                      <CFormItemList
                        v-if="agent.access.taqs?.length"
                        :items="agent.access.taqs"
                        item-key="id"
                        @remove="(_, idx) => removeTaq(idx)"
                      >
                        <template #default="{ item, index }">
                          <span class="font-medium text-color text-sm truncate block">
                            {{ loadedTaqNames[item.id] || $t('general.label.loading') }}
                          </span>
                          <InputText
                            v-model="agent.access.taqs[index].description"
                            class="w-full mt-1"
                            size="small"
                            :placeholder="$t('agent.editor.taqs.descriptionPlaceholder')"
                          />
                        </template>
                      </CFormItemList>

                      <Divider class="my-4" />

                      <!-- Workflows -->
                      <CFormGroup
                        :label="$t('agent.editor.workflows.label')"
                        :description="$t('agent.editor.workflows.help')"
                      />

                      <CInputWorkflow
                        v-model="workflowPickerSelection"
                        :placeholder="$t('agent.editor.workflows.selectPlaceholder')"
                        @update:model-value="onWorkflowPickerSelect"
                      />

                      <CFormItemList
                        v-if="agent.access.workflows?.length"
                        :items="agent.access.workflows"
                        item-key="id"
                        @remove="(_, idx) => removeWorkflow(idx)"
                      >
                        <template #default="{ item, index }">
                          <span class="font-medium text-color text-sm truncate block">
                            {{ loadedWorkflowNames[item.id] || $t('general.label.loading') }}
                          </span>
                          <InputText
                            v-model="agent.access.workflows[index].description"
                            class="w-full mt-1"
                            size="small"
                            :placeholder="$t('agent.editor.workflows.descriptionPlaceholder')"
                          />
                        </template>
                      </CFormItemList>
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

                        <CFormGroup
                          :label="$t('agent.editor.sidebarRoles.label')"
                          :description="$t('agent.editor.sidebarRoles.help')"
                          input-id="sidebarRoles"
                          :class="{
                            'opacity-50 pointer-events-none': !agent.invocation.user.enabled,
                          }"
                        >
                          <CInputRole
                            id="sidebarRoles"
                            v-model="agent.meta.sidebarRoles"
                            :multiple="true"
                            :disabled="!agent.invocation.user.enabled"
                          />
                        </CFormGroup>
                      </div>

                      <!-- System Invocation Group -->
                      <div class="flex flex-col gap-4">
                        <CInputToggleCard
                          v-model="agent.invocation.system.enabled"
                          :label="$t('agent.editor.systemEnabled.label')"
                          :description="$t('agent.editor.systemEnabled.help')"
                        />

                        <CFormGroup
                          :label="$t('agent.editor.serviceAccount.label')"
                          :description="$t('agent.editor.serviceAccount.help')"
                          input-id="serviceAccount"
                          :class="{
                            'opacity-50 pointer-events-none': !agent.invocation.system.enabled,
                          }"
                        >
                          <CInputUser
                            id="serviceAccount"
                            v-model="agent.invocation.system.serviceAccount"
                            :disabled="!agent.invocation.system.enabled"
                          />
                        </CFormGroup>
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
                  <div class="flex items-center gap-0 flex-1 overflow-x-auto no-scrollbar">
                    <button
                      v-for="(conv, idx) in conversations"
                      :key="idx"
                      type="button"
                      class="group flex items-center gap-1.5 py-2 pl-3 pr-1 border-b border-r border-r-surface whitespace-nowrap transition-colors duration-200 outline-none select-none max-w-[150px]"
                      :class="
                        activeConvIndex === idx
                          ? 'border-b-primary text-primary font-medium'
                          : 'border-b-transparent text-muted-color hover:text-color hover:border-b-surface-border'
                      "
                      @click="activeConvIndex = idx"
                    >
                      <span class="whitespace-nowrap truncate">{{ conv.label }}</span>
                      <i
                        class="pi pi-times text-xs p-1 hover:bg-surface rounded-full transition-all shrink-0"
                        :class="
                          conversations.length > 1
                            ? 'opacity-0 group-hover:opacity-100'
                            : 'invisible'
                        "
                        @click.stop="conversations.length > 1 && closeConversation(idx)"
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
                      :disabled="isCreate"
                      @click="addConversation"
                    />
                  </div>
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
        <CFormGroup :label="$t('agent.editor.tools.toolDescription')">
          <InputText
            v-model="editingToolForm.hints"
            :placeholder="$t('agent.editor.tools.descriptionPlaceholder')"
          />
        </CFormGroup>

        <!-- Restrict access toggle -->
        <div class="flex flex-col gap-3">
          <CInputToggleCard
            :modelValue="hasToolAllow(editingToolForm)"
            :label="$t('agent.editor.tools.configureAccess')"
            :description="$t('agent.editor.tools.configureAccessHelp')"
            @update:modelValue="toggleToolAllow(editingToolForm, $event)"
          />

          <!-- Namespace / Module rows -->
          <CFormGroup v-if="hasToolAllow(editingToolForm)">
            <template #actions>
              <Button
                :label="$t('agent.editor.tools.addNamespace')"
                icon="pi pi-plus"
                severity="secondary"
                size="small"
                @click="addToolAllowEntry(editingToolForm)"
              />
            </template>
            <CFormList v-model="editingToolForm.allow">
              <template #row="{ item }">
                <div class="flex flex-col gap-2 w-full">
                  <CInputNamespace
                    :model-value="item.namespaceID"
                    @update:model-value="onToolAllowNamespaceChange(item, $event)"
                    :placeholder="$t('agent.editor.tools.namespacePlaceholder')"
                  />
                  <CInputModule
                    v-if="item.namespaceID"
                    :model-value="item.moduleIDs || []"
                    @update:model-value="item.moduleIDs = $event"
                    :namespace-i-d="item.namespaceID"
                    :placeholder="$t('agent.editor.tools.modulesPlaceholder')"
                    :multiple="true"
                  />
                </div>
              </template>
            </CFormList>
          </CFormGroup>
        </div>
      </div>

      <template #footer>
        <div class="flex justify-end gap-2">
          <Button
            :label="$t('general.label.cancel')"
            severity="secondary"
            outlined
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

    <CEditorActions :back-to="true" @back="$router.back()">
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
    </CEditorActions>
  </Form>
</template>

<script setup>
import { computed, inject, nextTick, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { cloneDeep, isEqual } from 'lodash-es'

import { useAgentStore } from '@planetcrust/human-vue'
import { useRoute, useRouter } from 'vue-router'
import { useNamespaceStore, useModuleStore, useUnsavedGuard } from '@planetcrust/human-vue'
import { system } from '@planetcrust/human-js'

// Components (not globally registered)
import { components } from '@planetcrust/human-vue'
import AiChat from '../components/AiChat.vue'
import AiTrace from '../components/AiTrace.vue'
import { useEditorSplit } from '../composables/useEditorSplit'

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
const $ComposeAPI = inject('$ComposeAPI')
const agentStore = useAgentStore()
const namespaceStore = useNamespaceStore()
const moduleStore = useModuleStore()

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
    router.push({ name: 'agentic' })
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
      router.push({ name: 'agentic.edit', params: { agentID: created.agentID } })
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
    router.push({ name: 'agentic' })
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
        namespaceStore
          .findByID({ namespaceID: String(rule.namespaceID) })
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
            moduleStore
              .findByID($ComposeAPI, { namespaceID: String(rule.namespaceID), moduleID: String(modID) })
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
