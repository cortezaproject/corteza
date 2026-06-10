<template>
  <div class="h-full overflow-auto p-6">
    <div class="flex flex-col gap-8">
      <!-- Permitted AI Providers and Continuity -->
      <section class="flex flex-col gap-3">
        <h3 class="text-base font-medium">Permitted AI Providers and Continuity</h3>
        <div class="grid grid-cols-1 xl:grid-cols-2 gap-x-6 gap-y-4">
          <CFormGroup label="LLM Provider">
            <CInputLLM
              :model-value="ai.llmProvider"
              :disabled="disabled"
              placeholder="Select a provider"
              @update:model-value="v => setAi({ llmProvider: v, llmModel: null })"
            >
              <template #footer>
                <div class="border-t border-surface p-1">
                  <Button
                    label="Add provider"
                    icon="pi pi-plus"
                    text
                    fluid
                    size="small"
                    class="justify-start"
                    :disabled="disabled"
                    @click="addLlmProvider"
                  />
                </div>
              </template>
            </CInputLLM>
          </CFormGroup>
          <CFormGroup label="LLM Model">
            <CInputModel
              :model-value="ai.llmModel"
              :llm-provider-i-d="ai.llmProvider"
              :disabled="disabled"
              placeholder="Select a model"
              @update:model-value="v => setAi({ llmModel: v })"
            />
          </CFormGroup>
          <CFormGroup label="Alternative Provider">
            <CInputLLM
              :model-value="ai.altProvider"
              :disabled="disabled"
              placeholder="Select a provider"
              @update:model-value="v => setAi({ altProvider: v, altModel: null })"
            />
          </CFormGroup>
          <CFormGroup label="Alternative Model">
            <CInputModel
              :model-value="ai.altModel"
              :llm-provider-i-d="ai.altProvider"
              :disabled="disabled"
              placeholder="Select a model"
              @update:model-value="v => setAi({ altModel: v })"
            />
          </CFormGroup>
        </div>
      </section>

      <!-- Permitted Infrastructure and Continuity -->
      <section class="flex flex-col gap-4">
        <h3 class="text-base font-medium">Permitted Infrastructure and Continuity</h3>

        <CFormGroup label="Cloud Provider(s)">
          <MultiSelect
            :model-value="infra.cloudProviders || []"
            :options="CLOUD_PROVIDERS"
            display="chip"
            fluid
            :disabled="disabled"
            placeholder="Select cloud providers"
            @update:model-value="v => setInfra({ cloudProviders: v })"
          />
        </CFormGroup>

        <div class="flex flex-col gap-3">
          <CInputToggleCard
            label="Self-Hosting"
            description="The platform is hosted on your own infrastructure."
            :model-value="infra.selfHosting === 'Yes'"
            :disabled="disabled"
            @update:model-value="v => setInfra({ selfHosting: v ? 'Yes' : 'No' })"
          />
          <div
            v-if="infra.selfHosting === 'Yes'"
            class="grid grid-cols-1 xl:grid-cols-2 gap-x-6 gap-y-4 pl-1 max-w-2xl"
          >
            <CFormGroup label="Specification" class="xl:col-span-2">
              <Textarea
                :model-value="infra.specification"
                rows="2"
                auto-resize
                fluid
                :disabled="disabled"
                @update:model-value="v => setInfra({ specification: v })"
              />
            </CFormGroup>
            <CFormGroup label="Region(s)" class="xl:col-span-2">
              <InputText
                :model-value="infra.regions"
                fluid
                :disabled="disabled"
                placeholder="e.g. EU-West, EU-Central"
                @update:model-value="v => setInfra({ regions: v })"
              />
            </CFormGroup>
          </div>

          <CInputToggleCard
            label="Backup/Restore Available?"
            description="Data and configuration can be backed up and restored."
            :model-value="infra.backupRestore === 'Yes'"
            :disabled="disabled"
            @update:model-value="v => setInfra({ backupRestore: v ? 'Yes' : 'No' })"
          />

          <CInputToggleCard
            label="Failover Available?"
            description="Traffic can fail over to a secondary region if the primary is unavailable."
            :model-value="infra.failover === 'Yes'"
            :disabled="disabled"
            @update:model-value="v => setInfra({ failover: v ? 'Yes' : 'No' })"
          />
          <CFormGroup v-if="infra.failover === 'Yes'" label="Failover Region" class="pl-1 max-w-2xl">
            <InputText
              :model-value="infra.failoverRegion"
              fluid
              :disabled="disabled"
              placeholder="e.g. EU-North"
              @update:model-value="v => setInfra({ failoverRegion: v })"
            />
          </CFormGroup>
        </div>
      </section>

      <!-- Permitted Third-Party Connections (whitelist / catalogue) -->
      <section class="flex flex-col gap-3">
        <header class="flex items-center justify-between gap-3">
          <div>
            <h3 class="text-base font-medium">Permitted Third-Party Connections</h3>
            <p class="text-sm text-muted-color">
              The whitelist of connections the Connections step may later choose from.
            </p>
          </div>
          <Button
            v-if="!disabled"
            icon="pi pi-plus"
            label="Add connection"
            severity="secondary"
            size="small"
            @click="openPicker()"
          />
        </header>

        <div v-if="connections.length" class="flex flex-col gap-3">
          <div
            v-for="c in connections"
            :key="c.id"
            class="rounded-lg border border-surface p-4 flex flex-col gap-4"
          >
            <div class="flex items-start gap-3">
              <div class="grid grid-cols-1 xl:grid-cols-2 gap-x-6 gap-y-4 flex-1 min-w-0">
                <CFormGroup label="External Connection Name">
                  <div class="flex items-center gap-2 h-10 px-3 rounded-md border border-surface bg-emphasis min-w-0">
                    <i :class="[connectorIcon(c), 'text-primary shrink-0']" />
                    <span class="truncate">{{ c.name }}</span>
                  </div>
                </CFormGroup>
                <CFormGroup label="Connection Type">
                  <Select
                    :model-value="c.type"
                    :options="CONNECTION_TYPES"
                    fluid
                    :disabled="disabled"
                    placeholder="Select type"
                    @update:model-value="v => updateConn(c.id, { type: v })"
                  />
                </CFormGroup>
                <CFormGroup label="Action if Unavailable">
                  <SelectButton
                    :model-value="c.actionIfUnavailable"
                    :options="ACTIONS"
                    :allow-empty="false"
                    :disabled="disabled"
                    @update:model-value="v => updateConn(c.id, { actionIfUnavailable: v })"
                  />
                </CFormGroup>
                <CFormGroup v-if="c.actionIfUnavailable === 'Replace'" label="Replacement Connection">
                  <button
                    type="button"
                    class="flex items-center gap-2 w-full h-10 px-3 rounded-md border border-surface bg-surface min-w-0 text-left transition-colors"
                    :class="disabled ? 'opacity-60 cursor-not-allowed' : 'cursor-pointer hover:border-primary'"
                    :disabled="disabled"
                    @click="openPicker(c.id)"
                  >
                    <template v-if="c.replacement">
                      <i :class="[iconForConnector(c.replacementConnector), 'text-primary shrink-0']" />
                      <span class="truncate flex-1">{{ c.replacement }}</span>
                    </template>
                    <span v-else class="text-muted-color flex-1">Select replacement…</span>
                    <i class="pi pi-chevron-down text-xs text-muted-color shrink-0" />
                  </button>
                </CFormGroup>
                <CInputToggleCard
                  class="xl:col-span-2"
                  label="Is Connection an AI System?"
                  description="The connected system is itself an AI system."
                  :model-value="c.isAiSystem === 'Yes'"
                  :disabled="disabled"
                  @update:model-value="v => updateConn(c.id, { isAiSystem: v ? 'Yes' : 'No' })"
                />
                <CFormGroup label="Description" class="xl:col-span-2">
                  <Textarea
                    :model-value="c.description"
                    rows="2"
                    auto-resize
                    fluid
                    :disabled="disabled"
                    @update:model-value="v => updateConn(c.id, { description: v })"
                  />
                </CFormGroup>
              </div>
              <Button
                v-if="!disabled"
                icon="pi pi-times"
                severity="secondary"
                text
                rounded
                size="small"
                @click="store.removePermittedConnection(project.id, c.id)"
              />
            </div>
          </div>
        </div>
        <p v-else class="text-sm text-muted-color italic">No permitted connections yet.</p>
      </section>
    </div>

    <ConnectorPicker v-model="pickerOpen" @pick="onPick" />
    <LlmProviderDialog v-model="llmDialogOpen" @created="onLlmCreated" />
  </div>
</template>

<script setup>
import ConnectorPicker from '@/sections/project-poc/components/connections/ConnectorPicker.vue'
import LlmProviderDialog from '@/sections/project-poc/components/wizard/LlmProviderDialog.vue'
import { connector } from '@/sections/project-poc/config/connectors'
import { kindConfig } from '@/sections/project-poc/config/kinds'
import { useProjectsStore } from '@/sections/project-poc/stores/projects'
import { components } from '@planetcrust/human-vue'
import { computed, ref } from 'vue'

const { CInputLLM, CInputModel, CInputToggleCard } = components

const props = defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const store = useProjectsStore()

// New LLM providers are created inline via a dialog; on success we auto-select
// the new provider so the user doesn't have to find it in the list.
const llmDialogOpen = ref(false)
const addLlmProvider = () => {
  llmDialogOpen.value = true
}
const onLlmCreated = created => {
  if (created?.llmProviderID) setAi({ llmProvider: created.llmProviderID, llmModel: null })
}

const CLOUD_PROVIDERS = [
  'AWS',
  'Microsoft Azure',
  'Google Cloud',
  'Oracle Cloud',
  'IBM Cloud',
  'DigitalOcean',
  'On-premises',
]
const CONNECTION_TYPES = ['Database (DAL)', 'Application (TAQ)']
const ACTIONS = ['Deactivate', 'Replace']

const cfg = kindConfig('connection')
const iconForConnector = id => connector(id)?.icon || cfg.icon
const connectorIcon = c => iconForConnector(c.connector)

const ai = computed(() => props.project.resourceManagement?.ai || {})
const infra = computed(() => props.project.resourceManagement?.infra || {})
const connections = computed(() => props.project.resourceManagement?.connections || [])

const setAi = patch => store.updateResourceManagement(props.project.id, 'ai', patch)
const setInfra = patch => store.updateResourceManagement(props.project.id, 'infra', patch)
const updateConn = (id, patch) => store.updatePermittedConnection(props.project.id, id, patch)

// The picker is shared: `pickerTarget` is null when adding a new whitelist
// entry, or a connection id when choosing that row's replacement.
const pickerOpen = ref(false)
const pickerTarget = ref(null)
function openPicker(target = null) {
  pickerTarget.value = target
  pickerOpen.value = true
}
function onPick(c) {
  pickerOpen.value = false
  if (pickerTarget.value) {
    updateConn(pickerTarget.value, { replacement: c.label, replacementConnector: c.id })
  } else {
    store.addPermittedConnection(props.project.id, { name: c.label, connector: c.id })
  }
  pickerTarget.value = null
}
</script>
