<template>
  <div class="flex flex-col gap-8">
    <!-- Permitted AI Providers and Continuity -->
    <section class="flex flex-col gap-3">
      <header class="flex items-center justify-between gap-3">
        <div>
          <h3 class="text-base font-medium">{{ $t('project.resourceManagement.ai.title') }}</h3>
          <p class="text-sm text-muted-color">
            {{ $t('project.resourceManagement.ai.description') }}
          </p>
        </div>
        <Button
          v-if="!disabled"
          icon="pi pi-plus"
          :label="$t('project.resourceManagement.ai.addProvider')"
          severity="secondary"
          size="small"
          @click="addProvider"
        />
      </header>

      <CFormList
        v-model="providers"
        :columns="PROVIDER_COLUMNS"
        :empty-message="$t('project.resourceManagement.ai.empty')"
        :hide-remove="disabled"
      >
        <template #row="{ item: p }">
          <CInputLLM
            :model-value="p.provider"
            :disabled="disabled"
            size="small"
            :placeholder="$t('project.resourceManagement.ai.providerPlaceholder')"
            @update:model-value="v => setProvider(p, v)"
          >
            <template #footer>
              <div class="border-t border-surface p-1">
                <Button
                  :label="$t('project.resourceManagement.ai.createProvider')"
                  icon="pi pi-plus"
                  text
                  fluid
                  size="small"
                  class="justify-start"
                  :disabled="disabled"
                  @click="openLlmDialog(p.id)"
                />
              </div>
            </template>
          </CInputLLM>
          <CInputModel
            v-model="p.model"
            :llm-provider-i-d="p.provider"
            :disabled="disabled"
            size="small"
            :placeholder="$t('project.resourceManagement.ai.modelPlaceholder')"
          />
        </template>
      </CFormList>
    </section>

    <!-- Permitted Infrastructure and Continuity -->
    <section class="flex flex-col gap-3">
      <header class="flex items-center justify-between gap-3">
        <div>
          <h3 class="text-base font-medium">{{ $t('project.resourceManagement.infra.title') }}</h3>
          <p class="text-sm text-muted-color">
            {{ $t('project.resourceManagement.infra.description') }}
          </p>
        </div>
        <Button
          v-if="!disabled"
          icon="pi pi-plus"
          :label="$t('project.resourceManagement.infra.add')"
          severity="secondary"
          size="small"
          @click="addInfra"
        />
      </header>

      <div v-if="infraProviders.length" class="flex flex-col gap-3">
        <div v-for="p in infraProviders" :key="p.id" class="rounded-lg border border-surface p-4">
          <div class="flex items-start gap-3">
            <div class="flex flex-col gap-4 flex-1 min-w-0">
              <CInputToggleCard
                :label="$t('project.resourceManagement.infra.selfHosting.label')"
                :description="$t('project.resourceManagement.infra.selfHosting.description')"
                :model-value="p.selfHosting"
                :disabled="disabled"
                @update:model-value="v => updateInfra(p.id, { selfHosting: v })"
              />
              <CFormGroup v-if="!p.selfHosting" :label="$t('project.resourceManagement.infra.providerName')">
                <InputText
                  :model-value="p.name"
                  fluid
                  :disabled="disabled"
                  :placeholder="$t('project.resourceManagement.infra.providerNamePlaceholder')"
                  @update:model-value="v => updateInfra(p.id, { name: v })"
                />
              </CFormGroup>
              <CFormGroup :label="$t('project.resourceManagement.infra.specification')">
                <Textarea
                  :model-value="p.specification"
                  rows="2"
                  auto-resize
                  fluid
                  :disabled="disabled"
                  @update:model-value="v => updateInfra(p.id, { specification: v })"
                />
              </CFormGroup>
              <CFormGroup :label="$t('project.resourceManagement.infra.regions')">
                <InputText
                  :model-value="p.regions"
                  fluid
                  :disabled="disabled"
                  :placeholder="$t('project.resourceManagement.infra.regionsPlaceholder')"
                  @update:model-value="v => updateInfra(p.id, { regions: v })"
                />
              </CFormGroup>
              <CInputToggleCard
                :label="$t('project.resourceManagement.infra.failover.label')"
                :description="$t('project.resourceManagement.infra.failover.description')"
                :model-value="p.failover"
                :disabled="disabled"
                @update:model-value="v => updateInfra(p.id, { failover: v })"
              />
              <CFormGroup v-if="p.failover" :label="$t('project.resourceManagement.infra.failoverRegion')">
                <InputText
                  :model-value="p.failoverRegion"
                  fluid
                  :disabled="disabled"
                  :placeholder="$t('project.resourceManagement.infra.failoverRegionPlaceholder')"
                  @update:model-value="v => updateInfra(p.id, { failoverRegion: v })"
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
              @click="removeInfra(p.id)"
            />
          </div>
        </div>
      </div>
      <div
        v-else
        class="text-muted-color text-sm p-3 border border-surface rounded-border bg-highlight text-center"
      >
        {{ $t('project.resourceManagement.infra.empty') }}
      </div>

      <!-- Backup/Restore is project-wide, not per provider. -->
      <CInputToggleCard
        :label="$t('project.resourceManagement.infra.backupRestore.label')"
        :description="$t('project.resourceManagement.infra.backupRestore.description')"
        :model-value="infra.backupRestore"
        :disabled="disabled"
        @update:model-value="v => setInfra({ backupRestore: v })"
      />
    </section>

    <!-- Permitted Third-Party Connections (whitelist / catalogue) -->
    <section class="flex flex-col gap-3">
      <header class="flex items-center justify-between gap-3">
        <div>
          <h3 class="text-base font-medium">{{ $t('project.resourceManagement.connections.title') }}</h3>
          <p class="text-sm text-muted-color">
            {{ $t('project.resourceManagement.connections.description') }}
          </p>
        </div>
        <Button
          v-if="!disabled"
          icon="pi pi-plus"
          :label="$t('project.resourceManagement.connections.add')"
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
              <CFormGroup :label="$t('project.resourceManagement.connections.name')">
                <div class="flex items-center gap-2 h-10 px-3 rounded-md border border-surface bg-emphasis min-w-0">
                  <i :class="[iconForConnector(c.connector), 'text-primary shrink-0']" />
                  <span class="truncate">{{ c.name }}</span>
                </div>
              </CFormGroup>
              <CFormGroup :label="$t('project.resourceManagement.connections.type')">
                <Select
                  :model-value="c.type"
                  :options="CONNECTION_TYPES"
                  option-label="label"
                  option-value="value"
                  fluid
                  :disabled="disabled"
                  :placeholder="$t('project.resourceManagement.connections.typePlaceholder')"
                  @update:model-value="v => updateConn(c.id, { type: v })"
                />
              </CFormGroup>
              <CFormGroup :label="$t('project.resourceManagement.connections.actionIfUnavailable')">
                <SelectButton
                  :model-value="c.actionIfUnavailable"
                  :options="ACTIONS"
                  option-label="label"
                  option-value="value"
                  :allow-empty="false"
                  :disabled="disabled"
                  @update:model-value="v => updateConn(c.id, { actionIfUnavailable: v })"
                />
              </CFormGroup>
              <CFormGroup v-if="c.actionIfUnavailable === 'Replace'" :label="$t('project.resourceManagement.connections.replacement')">
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
                  <span v-else class="text-muted-color flex-1">{{ $t('project.resourceManagement.connections.selectReplacement') }}</span>
                  <i class="pi pi-chevron-down text-xs text-muted-color shrink-0" />
                </button>
              </CFormGroup>
              <CInputToggleCard
                class="xl:col-span-2"
                :label="$t('project.resourceManagement.connections.isAiSystem.label')"
                :description="$t('project.resourceManagement.connections.isAiSystem.description')"
                :model-value="c.isAiSystem"
                :disabled="disabled"
                @update:model-value="v => updateConn(c.id, { isAiSystem: v })"
              />
              <CFormGroup :label="$t('general.label.description')" class="xl:col-span-2">
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
              @click="removeConn(c.id)"
            />
          </div>
        </div>
      </div>
      <!-- Same empty-state styling as CFormList's empty-message (providers above). -->
      <div
        v-else
        class="text-muted-color text-sm p-3 border border-surface rounded-border bg-highlight text-center"
      >
        {{ $t('project.resourceManagement.connections.empty') }}
      </div>
    </section>

    <ConnectorPicker v-model="pickerOpen" @pick="onPick" />
    <LlmProviderDialog v-model="llmDialogOpen" @created="onLlmCreated" />
  </div>
</template>

<script setup>
import ConnectorPicker from '@/sections/project/components/connections/ConnectorPicker.vue'
import LlmProviderDialog from '@/sections/project/components/wizard/LlmProviderDialog.vue'
import { connector } from '@/sections/project/config/connectors'
import { kindConfig } from '@/sections/project/config/kinds'
import { components, useConfirmDelete } from '@planetcrust/human-vue'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()
const { CInputLLM, CInputModel, CInputToggleCard } = components

// Working copy of the governance step values ({ ai, infra, connections });
// the wizard owns loading and Save, this form only edits the copy.
const props = defineProps({
  modelValue: { type: Object, default: () => ({}) },
  disabled: { type: Boolean, default: false },
})
const emit = defineEmits(['update:modelValue'])

// Persisted enum values stay literal; only the displayed label is localized.
const CONNECTION_TYPES = computed(() => [
  { label: t('project.resourceManagement.connections.typeDatabase'), value: 'Database (DAL)' },
  { label: t('project.resourceManagement.connections.typeApplication'), value: 'Application (TAQ)' },
])
const ACTIONS = computed(() => [
  { label: t('project.resourceManagement.connections.actionDeactivate'), value: 'Deactivate' },
  { label: t('project.resourceManagement.connections.actionReplace'), value: 'Replace' },
])

const cfg = kindConfig('connection')
const iconForConnector = id => connector(id)?.icon || cfg.icon

const ai = computed(() => props.modelValue.ai || {})
const infra = computed(() => props.modelValue.infra || {})
const connections = computed(() => props.modelValue.connections || [])

const setInfra = patch =>
  emit('update:modelValue', { ...props.modelValue, infra: { ...infra.value, ...patch } })

const localId = prefix =>
  `${prefix}-${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 6)}`

// Permitted providers. CFormList removes through the setter; row edits mutate
// the (working-copy) row objects in place.
const providers = computed({
  get: () => ai.value.providers || [],
  set: v => emit('update:modelValue', { ...props.modelValue, ai: { ...ai.value, providers: v } }),
})

const PROVIDER_COLUMNS = computed(() => [
  { label: t('project.resourceManagement.ai.columnProvider'), width: 'minmax(16rem, 1fr)' },
  { label: t('project.resourceManagement.ai.columnModel'), width: 'minmax(16rem, 1fr)' },
])

// Permitted infrastructure providers (cards). Edits emit a fresh array so the
// working copy never mutates the cached governance values.
const infraProviders = computed(() => infra.value.providers || [])
const setInfraProviders = next =>
  emit('update:modelValue', { ...props.modelValue, infra: { ...infra.value, providers: next } })
const addInfra = () =>
  setInfraProviders([
    ...infraProviders.value,
    {
      id: localId('infra'),
      selfHosting: true,
      name: '',
      specification: '',
      regions: '',
      failover: false,
      failoverRegion: '',
    },
  ])
const updateInfra = (id, patch) =>
  setInfraProviders(infraProviders.value.map(p => (p.id === id ? { ...p, ...patch } : p)))
const removeInfra = id =>
  confirmDelete({
    header: t('project.resourceManagement.infra.removeConfirm.header'),
    message: t('project.resourceManagement.infra.removeConfirm.message'),
    onConfirm: () => setInfraProviders(infraProviders.value.filter(p => p.id !== id)),
  })

const addProvider = () => {
  providers.value = [...providers.value, { id: localId('llm'), provider: null, model: null }]
}

// Switching a row's provider invalidates its model pick.
const setProvider = (row, provider) => {
  row.provider = provider
  row.model = null
}

const setConnections = next => emit('update:modelValue', { ...props.modelValue, connections: next })
const updateConn = (id, patch) =>
  setConnections(connections.value.map(c => (c.id === id ? { ...c, ...patch } : c)))
const removeConn = id =>
  confirmDelete({
    header: t('project.resourceManagement.connections.removeConfirm.header'),
    message: t('project.resourceManagement.connections.removeConfirm.message', {
      name: connections.value.find(c => c.id === id)?.name || '',
    }),
    onConfirm: () => setConnections(connections.value.filter(c => c.id !== id)),
  })
const addConn = c =>
  setConnections([
    ...connections.value,
    {
      id: `pconn-${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 6)}`,
      name: c.label,
      connector: c.id,
      type: null,
      actionIfUnavailable: 'Deactivate',
      replacement: '',
      replacementConnector: null,
      isAiSystem: false,
      description: '',
    },
  ])

// New LLM providers are created inline via a dialog shared by all rows; on
// success the new provider is auto-selected into the row that opened it.
const llmDialogOpen = ref(false)
const llmDialogTarget = ref(null) // row id
function openLlmDialog(rowId) {
  llmDialogTarget.value = rowId
  llmDialogOpen.value = true
}
const onLlmCreated = created => {
  if (!created?.llmProviderID) return
  const row = providers.value.find(p => p.id === llmDialogTarget.value)
  if (row) setProvider(row, created.llmProviderID)
}

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
    addConn(c)
  }
  pickerTarget.value = null
}
</script>
