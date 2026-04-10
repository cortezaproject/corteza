<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.labels.editor.title') }} — {{ labelName }}</span>
  </Teleport>

  <div class="flex flex-col h-full">
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <!-- Basic Info Panel -->
      <Panel :header="$t('system.labels.editor.info.title')" toggleable :collapsed="false">
        <div class="flex flex-col gap-2">
          <label class="font-medium text-primary text-sm">
            {{ $t('system.labels.editor.info.name') }}
          </label>
          <InputText :model-value="labelName" disabled />
        </div>
      </Panel>

      <!-- Namespaces with this label -->
      <Panel :header="$t('system.labels.editor.namespaces.title')" toggleable>
        <CResourceTable
          primary-key="namespaceID"
          :fields="nsFields"
          :items="nsItems"
          :loading="nsLoading"
          :empty-message="$t('system.labels.editor.namespaces.empty')"
          :action-items="getNsActions"
        >
          <template #header>
            <Button
              :label="$t('system.labels.editor.namespaces.create')"
              icon="pi pi-plus"
              size="small"
              @click="showNsDialog = true"
            />
          </template>
          <template #body-enabled="{ data }">
            <Tag
              :value="data.enabled ? $t('general.label.enabled') : $t('general.label.disabled')"
              :severity="data.enabled ? 'success' : 'secondary'"
            />
          </template>
        </CResourceTable>
      </Panel>

      <!-- Agents with this label -->
      <Panel :header="$t('system.labels.editor.agents.title')" toggleable>
        <CResourceTable
          primary-key="agentID"
          :fields="agentFields"
          :items="agentItems"
          :loading="agentLoading"
          :empty-message="$t('system.labels.editor.agents.empty')"
          :action-items="getAgentActions"
        >
          <template #header>
            <Button
              :label="$t('system.labels.editor.agents.create')"
              icon="pi pi-plus"
              size="small"
              @click="showAgentDialog = true"
            />
          </template>
          <template #body-name="{ data }">
            {{ data.meta?.short || data.handle || '—' }}
          </template>
        </CResourceTable>
      </Panel>

      <!-- TAQ Automations with this label -->
      <Panel :header="$t('system.labels.editor.automations.title')" toggleable>
        <CResourceTable
          primary-key="automationID"
          :fields="taqFields"
          :items="taqItems"
          :loading="taqLoading"
          :empty-message="$t('system.labels.editor.automations.empty')"
          :action-items="getTaqActions"
        >
          <template #header>
            <Button
              :label="$t('system.labels.editor.automations.create')"
              icon="pi pi-plus"
              size="small"
              @click="showTaqDialog = true"
            />
          </template>
          <template #body-name="{ data }">
            {{ data.meta?.short || data.handle || '—' }}
          </template>
          <template #body-enabled="{ data }">
            <Tag
              :value="data.enabled ? $t('general.label.enabled') : $t('general.label.disabled')"
              :severity="data.enabled ? 'success' : 'secondary'"
            />
          </template>
        </CResourceTable>
      </Panel>
    </div>

    <!-- Bottom toolbar -->
    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="p-3 flex items-center justify-between">
        <Button
          :label="$t('general.label.back')"
          icon="pi pi-arrow-left"
          severity="secondary"
          @click="$router.push({ name: 'system.labels' })"
        />
      </div>
    </div>

    <!-- Create Namespace Dialog -->
    <Dialog
      v-model:visible="showNsDialog"
      :header="$t('system.labels.editor.namespaces.createDialog')"
      modal
      :style="{ width: '32rem' }"
    >
      <div class="flex flex-col gap-4">
        <div class="flex flex-col gap-2">
          <label class="font-medium text-primary text-sm">
            {{ $t('system.labels.editor.namespaces.nameLabel') }}
          </label>
          <InputText
            v-model="nsForm.name"
            :placeholder="$t('system.labels.editor.namespaces.namePlaceholder')"
            fluid
          />
        </div>
        <div class="flex flex-col gap-2">
          <label class="font-medium text-primary text-sm">
            {{ $t('system.labels.editor.namespaces.slugLabel') }}
          </label>
          <InputText
            v-model="nsForm.slug"
            :placeholder="$t('system.labels.editor.namespaces.slugPlaceholder')"
            fluid
          />
        </div>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button
            :label="$t('general.label.cancel')"
            severity="secondary"
            text
            size="small"
            @click="showNsDialog = false"
          />
          <Button
            :label="$t('general.label.create')"
            size="small"
            :disabled="!nsForm.name.trim()"
            :loading="nsCreating"
            @click="createNamespace"
          />
        </div>
      </template>
    </Dialog>

    <!-- Create Agent Dialog -->
    <Dialog
      v-model:visible="showAgentDialog"
      :header="$t('system.labels.editor.agents.createDialog')"
      modal
      :style="{ width: '32rem' }"
    >
      <div class="flex flex-col gap-4">
        <div class="flex flex-col gap-2">
          <label class="font-medium text-primary text-sm">
            {{ $t('system.labels.editor.agents.nameLabel') }}
          </label>
          <InputText
            v-model="agentForm.name"
            :placeholder="$t('system.labels.editor.agents.namePlaceholder')"
            fluid
          />
        </div>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button
            :label="$t('general.label.cancel')"
            severity="secondary"
            text
            size="small"
            @click="showAgentDialog = false"
          />
          <Button
            :label="$t('general.label.create')"
            size="small"
            :disabled="!agentForm.name.trim()"
            :loading="agentCreating"
            @click="createAgent"
          />
        </div>
      </template>
    </Dialog>

    <!-- Create TAQ Dialog -->
    <Dialog
      v-model:visible="showTaqDialog"
      :header="$t('system.labels.editor.automations.createDialog')"
      modal
      :style="{ width: '32rem' }"
    >
      <div class="flex flex-col gap-4">
        <div class="flex flex-col gap-2">
          <label class="font-medium text-primary text-sm">
            {{ $t('system.labels.editor.automations.nameLabel') }}
          </label>
          <InputText
            v-model="taqForm.name"
            :placeholder="$t('system.labels.editor.automations.namePlaceholder')"
            fluid
          />
        </div>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button
            :label="$t('general.label.cancel')"
            severity="secondary"
            text
            size="small"
            @click="showTaqDialog = false"
          />
          <Button
            :label="$t('general.label.create')"
            size="small"
            :disabled="!taqForm.name.trim()"
            :loading="taqCreating"
            @click="createAutomation"
          />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { components, useConfirmDelete } from '@cortezaproject/corteza-vue-next'

const { CResourceTable } = components

const { t } = useI18n()
const route = useRoute()
const { confirmDelete } = useConfirmDelete()

const $ComposeAPI = inject('$ComposeAPI')
const $AutomationAPI = inject('$AutomationAPI')
const $SystemAPI = inject('$SystemAPI')
const $toast = inject('$toast')

const labelName = computed(() => decodeURIComponent(route.params.labelID || ''))

function buildCrossAppUrl(appBase, path) {
  const u = new URL(window.location)
  return `${u.origin}/${appBase}/${path}`
}

function openInNewTab(appBase, path) {
  window.open(buildCrossAppUrl(appBase, path), '_blank', 'noopener')
}

function getNsActions(data) {
  const items = [{
    label: t('system.labels.editor.openNewTab'),
    icon: 'pi pi-external-link',
    command: () => openInNewTab('compose', `namespace/${data.slug || data.namespaceID}`),
  }]

  if (data.canDeleteNamespace) {
    items.push({ separator: true })
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => onDeleteNamespace(data),
    })
  }

  return items
}

function getAgentActions(data) {
  const items = [{
    label: t('system.labels.editor.openNewTab'),
    icon: 'pi pi-external-link',
    command: () => openInNewTab('agentic', data.agentID),
  }]

  if (data.canDeleteAgent) {
    items.push({ separator: true })
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => onDeleteAgent(data),
    })
  }

  return items
}

function getTaqActions(data) {
  const items = [{
    label: t('system.labels.editor.openNewTab'),
    icon: 'pi pi-external-link',
    command: () => openInNewTab('taq', `builder/${data.automationID}`),
  }]

  if (data.canDeleteNgAutomation) {
    items.push({ separator: true })
    items.push({
      label: t('general.label.delete'),
      icon: 'pi pi-trash',
      class: 'text-red-500',
      command: () => onDeleteAutomation(data),
    })
  }

  return items
}

function onDeleteNamespace(data) {
  confirmDelete({
    message: t('general.confirm.delete'),
    header: data.name || data.namespaceID,
    onConfirm: () => handleDeleteNamespace(data),
  })
}

function onDeleteAgent(data) {
  confirmDelete({
    message: t('general.confirm.delete'),
    header: data.meta?.short || data.handle || data.agentID,
    onConfirm: () => handleDeleteAgent(data),
  })
}

function onDeleteAutomation(data) {
  confirmDelete({
    message: t('general.confirm.delete'),
    header: data.meta?.short || data.handle || data.automationID,
    onConfirm: () => handleDeleteAutomation(data),
  })
}

async function handleDeleteNamespace(data) {
  try {
    await $ComposeAPI.namespaceDelete({ namespaceID: data.namespaceID })
    $toast?.toastSuccess(t('system.labels.editor.toast.deleted'))
    fetchNamespaces()
  } catch (e) {
    console.error('Failed to delete namespace:', e)
    $toast?.toastDanger(t('system.labels.editor.toast.deleteError'))
  }
}

async function handleDeleteAgent(data) {
  try {
    await $SystemAPI.agentDelete({ agentID: data.agentID })
    $toast?.toastSuccess(t('system.labels.editor.toast.deleted'))
    fetchAgents()
  } catch (e) {
    console.error('Failed to delete agent:', e)
    $toast?.toastDanger(t('system.labels.editor.toast.deleteError'))
  }
}

async function handleDeleteAutomation(data) {
  try {
    await $AutomationAPI.ngAutomationDelete({ automationID: data.automationID })
    $toast?.toastSuccess(t('system.labels.editor.toast.deleted'))
    fetchAutomations()
  } catch (e) {
    console.error('Failed to delete automation:', e)
    $toast?.toastDanger(t('system.labels.editor.toast.deleteError'))
  }
}

// ──────────────────────────────────────────────
// NAMESPACES
// ──────────────────────────────────────────────
const nsItems = ref([])
const nsLoading = ref(false)

const nsFields = [
  { key: 'name', header: t('system.labels.editor.namespaces.columns.name') },
  { key: 'enabled', header: t('system.labels.editor.namespaces.columns.enabled') },
]

async function fetchNamespaces() {
  nsLoading.value = true
  try {
    const labelFilter = {}
    labelFilter[labelName.value] = ''
    const result = await $ComposeAPI.namespaceList({ labels: labelFilter, limit: 100 })
    nsItems.value = result.set || []
  } catch (e) {
    console.error('Failed to load namespaces:', e)
    nsItems.value = []
  } finally {
    nsLoading.value = false
  }
}

// Create namespace
const showNsDialog = ref(false)
const nsCreating = ref(false)
const nsForm = reactive({ name: '', slug: '' })

async function createNamespace() {
  if (!nsForm.name.trim()) return
  nsCreating.value = true
  try {
    const labels = {}
    labels[labelName.value] = ''
    await $ComposeAPI.namespaceCreate({
      name: nsForm.name.trim(),
      slug: nsForm.slug.trim() || undefined,
      enabled: false,
      meta: {},
      labels,
    })
    showNsDialog.value = false
    nsForm.name = ''
    nsForm.slug = ''
    $toast?.toastSuccess(t('system.labels.editor.toast.created'))
    fetchNamespaces()
  } catch (e) {
    console.error('Failed to create namespace:', e)
    $toast?.toastDanger(t('system.labels.editor.toast.createError'))
  } finally {
    nsCreating.value = false
  }
}

// ──────────────────────────────────────────────
// AGENTS
// ──────────────────────────────────────────────
const agentItems = ref([])
const agentLoading = ref(false)

const agentFields = [{ key: 'name', header: t('system.labels.editor.agents.columns.name') }]

async function fetchAgents() {
  agentLoading.value = true
  try {
    // Agents don't support label filtering — fetch all and filter client-side
    const result = await $SystemAPI.agentList({ limit: 200 })
    const set = result.set || []
    agentItems.value = set.filter(a => {
      if (!a.labels || typeof a.labels !== 'object') return false
      return Object.prototype.hasOwnProperty.call(a.labels, labelName.value)
    })
  } catch (e) {
    console.error('Failed to load agents:', e)
    agentItems.value = []
  } finally {
    agentLoading.value = false
  }
}

// Create agent
const showAgentDialog = ref(false)
const agentCreating = ref(false)
const agentForm = reactive({ name: '' })

async function createAgent() {
  if (!agentForm.name.trim()) return
  agentCreating.value = true
  try {
    const labels = {}
    labels[labelName.value] = ''
    await $SystemAPI.agentCreate({
      meta: { short: agentForm.name.trim() },
      labels,
    })
    showAgentDialog.value = false
    agentForm.name = ''
    $toast?.toastSuccess(t('system.labels.editor.toast.created'))
    fetchAgents()
  } catch (e) {
    console.error('Failed to create agent:', e)
    $toast?.toastDanger(t('system.labels.editor.toast.createError'))
  } finally {
    agentCreating.value = false
  }
}

// ──────────────────────────────────────────────
// TAQ AUTOMATIONS
// ──────────────────────────────────────────────
const taqItems = ref([])
const taqLoading = ref(false)

const taqFields = [
  { key: 'name', header: t('system.labels.editor.automations.columns.name') },
  { key: 'enabled', header: t('system.labels.editor.automations.columns.enabled') },
]

async function fetchAutomations() {
  taqLoading.value = true
  try {
    const labelFilter = {}
    labelFilter[labelName.value] = ''
    const result = await $AutomationAPI.ngAutomationList({ labels: labelFilter, limit: 100 })
    taqItems.value = result.set || []
  } catch (e) {
    console.error('Failed to load automations:', e)
    taqItems.value = []
  } finally {
    taqLoading.value = false
  }
}

// Create TAQ automation
const showTaqDialog = ref(false)
const taqCreating = ref(false)
const taqForm = reactive({ name: '' })

async function createAutomation() {
  if (!taqForm.name.trim()) return
  taqCreating.value = true
  try {
    const labels = {}
    labels[labelName.value] = ''
    await $AutomationAPI.ngAutomationCreate({
      meta: { short: taqForm.name.trim() },
      enabled: false,
      labels,
      triggers: [],
      steps: [],
      paths: [],
    })
    showTaqDialog.value = false
    taqForm.name = ''
    $toast?.toastSuccess(t('system.labels.editor.toast.created'))
    fetchAutomations()
  } catch (e) {
    console.error('Failed to create automation:', e)
    $toast?.toastDanger(t('system.labels.editor.toast.createError'))
  } finally {
    taqCreating.value = false
  }
}

// ──────────────────────────────────────────────
// INIT
// ──────────────────────────────────────────────
onMounted(() => {
  fetchNamespaces()
  fetchAgents()
  fetchAutomations()
})
</script>
