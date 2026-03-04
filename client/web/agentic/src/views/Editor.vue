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
                    <div class="flex flex-col gap-2">
                      <label for="name" class="font-medium text-primary">
                        {{ $t('agent.editor.fields.name') }}
                      </label>
                      <InputText id="name" v-model="agent.meta.short" />
                    </div>
                    <div class="flex flex-col gap-2">
                      <label for="handle" class="font-medium text-primary">
                        {{ $t('agent.editor.fields.handle') }}
                      </label>
                      <InputText id="handle" v-model="agent.handle" />
                    </div>
                    <div class="flex flex-col gap-2 md:col-span-2">
                      <label for="description" class="font-medium text-primary">
                        {{ $t('agent.editor.fields.description') }}
                      </label>
                      <Textarea
                        id="description"
                        v-model="agent.meta.description"
                        rows="3"
                        autoResize
                      />
                    </div>
                    <div class="flex flex-col gap-2">
                      <label for="status" class="font-medium text-primary">
                        {{ $t('agent.editor.fields.status') }}
                      </label>
                      <Select
                        id="status"
                        v-model="agent.status"
                        :options="statusOptions"
                        optionLabel="label"
                        optionValue="value"
                      />
                    </div>
                  </div>
                </Panel>

                <Panel :header="$t('agent.editor.panels.behavior')" toggleable>
                  <div class="grid grid-cols-1 gap-4">
                    <div class="flex flex-col gap-2">
                      <label for="systemPrompt" class="font-medium text-primary">
                        {{ $t('agent.editor.fields.systemPrompt') }}
                      </label>
                      <Textarea
                        id="systemPrompt"
                        v-model="agent.behavior.systemPrompt"
                        rows="6"
                        autoResize
                      />
                    </div>
                  </div>
                </Panel>

                <Panel :header="$t('agent.editor.panels.execution')" toggleable>
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div class="flex flex-col gap-2">
                      <label for="provider" class="font-medium text-primary">
                        {{ $t('agent.editor.fields.provider') }}
                      </label>
                      <CInputLLM id="provider" v-model="agent.execution.model.llmProviderID" />
                    </div>
                    <div class="flex flex-col gap-2">
                      <label for="model" class="font-medium text-primary">
                        {{ $t('agent.editor.fields.model') }}
                      </label>
                      <CInputModel
                        id="model"
                        v-model="agent.execution.model.model"
                        :llmProviderID="agent.execution.model.llmProviderID"
                      />
                    </div>
                    <div class="flex flex-col gap-2">
                      <label for="temperature" class="font-medium text-primary">
                        {{ $t('agent.editor.fields.temperature') }} ({{
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
                    </div>
                    <div class="flex flex-col gap-2">
                      <label for="maxTokens" class="font-medium text-primary">
                        {{ $t('agent.editor.fields.maxTokens') }}
                      </label>
                      <InputNumber
                        id="maxTokens"
                        v-model="agent.execution.model.maxTokens"
                        mode="decimal"
                        :useGrouping="false"
                      />
                    </div>
                  </div>
                </Panel>
                <Panel :header="$t('agent.editor.panels.access')" toggleable>
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div class="flex flex-col gap-2">
                      <label for="contextNamespace" class="font-medium text-primary">
                        {{ $t('agent.editor.fields.contextNamespace') }}
                      </label>
                      <InputText id="contextNamespace" v-model="agent.access.context.namespace" />
                    </div>
                    <div class="flex flex-col gap-2">
                      <label for="contextModule" class="font-medium text-primary">
                        {{ $t('agent.editor.fields.contextModule') }}
                      </label>
                      <InputText id="contextModule" v-model="agent.access.context.module" />
                    </div>
                    <div class="md:col-span-2 pt-2">
                      <div class="flex flex-col gap-2 mb-2">
                        <label class="font-medium text-primary">
                          {{ $t('agent.editor.fields.tools') }}
                        </label>
                        <MultiSelect
                          :modelValue="selectedToolNames"
                          @update:modelValue="onToolSelectionChange"
                          :options="availableTools"
                          optionLabel="name"
                          optionValue="name"
                          :placeholder="$t('agent.editor.fields.tools')"
                          :loading="loadingTools"
                          filter
                          display="chip"
                          class="w-full"
                        >
                          <template #option="{ option }">
                            <span class="font-medium">{{ option.description }}</span>
                          </template>
                        </MultiSelect>
                      </div>
                      <DataTable
                        v-if="agent.access.tools.length"
                        :value="agent.access.tools"
                        class="border border-surface rounded overflow-hidden"
                      >
                        <Column
                          field="name"
                          :header="$t('agent.editor.fields.toolName')"
                          class="w-1/3"
                        />
                        <Column field="hints" :header="$t('agent.editor.fields.toolHints')">
                          <template #body="{ data }">
                            <InputText v-model="data.hints" class="w-full" size="small" />
                          </template>
                        </Column>
                        <Column headerStyle="width: 3rem">
                          <template #body="{ data }">
                            <Button
                              icon="pi pi-trash"
                              severity="danger"
                              text
                              @click="removeTool(data.name)"
                            />
                          </template>
                        </Column>
                      </DataTable>
                    </div>
                  </div>
                </Panel>

                <Panel :header="$t('agent.editor.panels.invocation')" toggleable>
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div class="flex flex-col gap-2 justify-center">
                      <div class="flex items-center gap-2">
                        <Checkbox
                          v-model="agent.invocation.user.enabled"
                          inputId="userEnabled"
                          :binary="true"
                        />
                        <label for="userEnabled" class="font-medium text-primary">
                          {{ $t('agent.editor.fields.userEnabled') }}
                        </label>
                      </div>
                    </div>
                    <div class="flex flex-col gap-2 justify-center">
                      <div class="flex items-center gap-2">
                        <Checkbox
                          v-model="agent.invocation.system.enabled"
                          inputId="systemEnabled"
                          :binary="true"
                        />
                        <label for="systemEnabled" class="font-medium text-primary">
                          {{ $t('agent.editor.fields.systemEnabled') }}
                        </label>
                      </div>
                    </div>
                    <div class="flex flex-col gap-2">
                      <label for="serviceAccount" class="font-medium text-primary">
                        {{ $t('agent.editor.fields.serviceAccount') }}
                      </label>
                      <InputText
                        id="serviceAccount"
                        v-model="agent.invocation.system.serviceAccount"
                      />
                    </div>
                    <div class="flex flex-col gap-2">
                      <label for="outputFormat" class="font-medium text-primary">
                        {{ $t('agent.editor.fields.outputFormat') }}
                      </label>
                      <InputText id="outputFormat" v-model="agent.invocation.system.outputFormat" />
                    </div>
                  </div>
                </Panel>
              </TabPanel>

              <TabPanel value="exec" class="h-full p-0 flex flex-row">
                <div class="flex-1 flex flex-col border-r border-surface min-w-0">
                  <div class="flex-1 p-4 overflow-y-auto flex flex-col gap-4">
                    <div
                      v-for="(msg, index) in chatHistory"
                      :key="index"
                      :class="[
                        'p-3 xl:p-4 rounded-xl max-w-[85%] whitespace-pre-wrap text-sm md:text-base',
                        msg.role === 'user'
                          ? 'bg-primary text-primary-contrast self-end shadow-sm'
                          : 'bg-emphasis text-color self-start shadow-sm',
                      ]"
                    >
                      {{ msg.content }}
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
                <div class="w-1/3 flex flex-col bg-surface border-l border-surface min-w-0">
                  <div
                    class="p-3 border-b border-surface font-semibold text-sm shrink-0 flex items-center justify-between"
                  >
                    <span>{{ $t('agent.editor.playground.traceTitle') }}</span>
                    <Button
                      v-if="execTrace"
                      icon="pi pi-refresh"
                      text
                      rounded
                      severity="secondary"
                      size="small"
                      @click="execTrace = ''"
                    />
                  </div>
                  <div
                    class="flex-1 p-4 overflow-y-auto text-xs font-mono whitespace-pre-wrap break-all text-muted-color"
                  >
                    {{ execTrace || $t('agent.editor.playground.emptyTrace') }}
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
import { computed, inject, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
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
import MultiSelect from 'primevue/multiselect'
import Tab from 'primevue/tab'
import TabList from 'primevue/tablist'
import TabPanel from 'primevue/tabpanel'
import TabPanels from 'primevue/tabpanels'
import Tabs from 'primevue/tabs'
import Textarea from 'primevue/textarea'
import { components } from '@cortezaproject/corteza-vue-next'
const { CInputLLM, CInputModel } = components

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(false)
const saving = ref(false)
const agent = ref(null)
const activeTab = ref('config')

const chatHistory = ref([])
const chatInput = ref('')
const executing = ref(false)
const execTrace = ref('')
const conversationID = ref(null)
const isCreate = computed(() => !route.params.agentID)

const statusOptions = ref([
  { label: 'Active', value: 'active' },
  { label: 'Inactive', value: 'inactive' },
])

const availableTools = ref([])
const loadingTools = ref(false)

const selectedToolNames = computed(() => {
  if (!agent.value?.access?.tools) return []
  return agent.value.access.tools.map(t => t.name)
})

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

async function loadAgent() {
  const agentID = route.params.agentID
  if (!agentID) {
    agent.value = emptyAgent()
    return
  }

  loading.value = true
  try {
    const res = await $SystemAPI.agentRead({ agentID })

    // Fallbacks just in case the object is sparse
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

    if (!agent.value.execution.model.llmProviderID) {
      agent.value.execution.model.llmProviderID = '0'
    }
  } catch (err) {
    $toast.toastDanger('Failed to load agent')
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
      $toast.toastSuccess('Agent created')
      router.push({ name: 'agent.edit', params: { agentID: created.agentID } })
    } else {
      await $SystemAPI.agentUpdate({ agentID: route.params.agentID, ...agent.value })
      $toast.toastSuccess('Agent saved')
    }
  } catch (err) {
    console.error(err)
    $toast.toastDanger('Failed to save agent')
  } finally {
    saving.value = false
  }
}

async function sendChatMessage() {
  if (!chatInput.value.trim() || executing.value) return

  const input = chatInput.value
  chatInput.value = ''

  chatHistory.value.push({ role: 'user', content: input })
  executing.value = true
  execTrace.value = 'Executing...\n'

  try {
    const res = await $SystemAPI.agentExec({
      agentID: route.params.agentID,
      input: input,
      ...(conversationID.value ? { conversationID: conversationID.value } : {}),
    })

    if (res?.conversationID) {
      conversationID.value = res.conversationID
    }

    // We expect the agent execution to return a response structure.
    chatHistory.value.push({
      role: 'agent',
      content:
        res?.output || (typeof res === 'string' ? res : res?.response?.text || JSON.stringify(res)),
    })
    execTrace.value = JSON.stringify(res, null, 2)
  } catch (err) {
    console.error(err)
    $toast.toastDanger('Execution failed')
    chatHistory.value.push({ role: 'agent', content: 'Error: ' + err.message })
    execTrace.value = String(err)
  } finally {
    executing.value = false
  }
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
}

function onToolSelectionChange(selectedNames) {
  const currentTools = agent.value.access.tools
  // Build a map of existing tools to preserve hints
  const existingMap = {}
  currentTools.forEach(t => {
    existingMap[t.name] = t
  })

  // Rebuild tools array from selected names
  agent.value.access.tools = selectedNames.map(name => {
    if (existingMap[name]) {
      return existingMap[name]
    }
    // New tool — use description from available tools as default hint
    const available = availableTools.value.find(t => t.name === name)
    return { name, hints: available?.description || '' }
  })
}

function removeTool(toolName) {
  agent.value.access.tools = agent.value.access.tools.filter(t => t.name !== toolName)
}
</script>
