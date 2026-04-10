<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.labels.editor.title') }} — {{ labelName }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-auto flex flex-col gap-4">
    <!-- Basic Info Panel -->
    <Panel :header="$t('system.labels.editor.info.title')">
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
      >
        <template #header>
          <Button
            :label="$t('system.labels.editor.namespaces.add')"
            icon="pi pi-plus"
            size="small"
            @click="showNsDialog = true"
          />
        </template>
        <template #body-name="{ data }">
          <a :href="buildNamespaceLink(data)" class="text-primary hover:underline" @click.prevent="navigateTo(buildNamespaceLink(data))">
            {{ data.name || '—' }}
          </a>
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
      >
        <template #header>
          <Button
            :label="$t('system.labels.editor.agents.add')"
            icon="pi pi-plus"
            size="small"
            @click="showAgentDialog = true"
          />
        </template>
        <template #body-name="{ data }">
          <a :href="buildAgentLink(data)" class="text-primary hover:underline" @click.prevent="navigateTo(buildAgentLink(data))">
            {{ data.meta?.short || data.handle || '—' }}
          </a>
        </template>
        <template #body-status="{ data }">
          <Tag
            :value="data.status || '—'"
            :severity="data.status === 'active' ? 'success' : 'secondary'"
          />
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
      >
        <template #header>
          <Button
            :label="$t('system.labels.editor.automations.add')"
            icon="pi pi-plus"
            size="small"
            @click="showTaqDialog = true"
          />
        </template>
        <template #body-name="{ data }">
          <a :href="buildTaqLink(data)" class="text-primary hover:underline" @click.prevent="navigateTo(buildTaqLink(data))">
            {{ data.meta?.short || data.handle || '—' }}
          </a>
        </template>
        <template #body-enabled="{ data }">
          <Tag
            :value="data.enabled ? $t('general.label.enabled') : $t('general.label.disabled')"
            :severity="data.enabled ? 'success' : 'secondary'"
          />
        </template>
      </CResourceTable>
    </Panel>

    <!-- Back button -->
    <div>
      <Button
        :label="$t('general.label.back')"
        icon="pi pi-arrow-left"
        severity="secondary"
        @click="$router.push({ name: 'system.labels' })"
      />
    </div>

    <!-- Add Namespace Dialog -->
    <Dialog
      v-model:visible="showNsDialog"
      :header="$t('system.labels.editor.namespaces.addDialog')"
      modal
      :style="{ width: '32rem' }"
    >
      <div class="flex flex-col gap-2">
        <label class="font-medium text-primary text-sm">
          {{ $t('system.labels.editor.namespaces.selectLabel') }}
        </label>
        <Select
          v-model="selectedNamespace"
          :options="availableNamespaces"
          option-label="name"
          :placeholder="$t('system.labels.editor.namespaces.selectPlaceholder')"
          :loading="nsDialogLoading"
          filter
          fluid
          @show="fetchAvailableNamespaces"
        />
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="$t('general.label.cancel')" severity="secondary" text size="small" @click="showNsDialog = false" />
          <Button :label="$t('general.label.save')" size="small" :disabled="!selectedNamespace" @click="addNamespaceLabel" />
        </div>
      </template>
    </Dialog>

    <!-- Add Agent Dialog -->
    <Dialog
      v-model:visible="showAgentDialog"
      :header="$t('system.labels.editor.agents.addDialog')"
      modal
      :style="{ width: '32rem' }"
    >
      <div class="flex flex-col gap-2">
        <label class="font-medium text-primary text-sm">
          {{ $t('system.labels.editor.agents.selectLabel') }}
        </label>
        <Select
          v-model="selectedAgent"
          :options="availableAgents"
          :option-label="agentOptionLabel"
          :placeholder="$t('system.labels.editor.agents.selectPlaceholder')"
          :loading="agentDialogLoading"
          filter
          fluid
          @show="fetchAvailableAgents"
        />
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="$t('general.label.cancel')" severity="secondary" text size="small" @click="showAgentDialog = false" />
          <Button :label="$t('general.label.save')" size="small" :disabled="!selectedAgent" @click="addAgentLabel" />
        </div>
      </template>
    </Dialog>

    <!-- Add TAQ Dialog -->
    <Dialog
      v-model:visible="showTaqDialog"
      :header="$t('system.labels.editor.automations.addDialog')"
      modal
      :style="{ width: '32rem' }"
    >
      <div class="flex flex-col gap-2">
        <label class="font-medium text-primary text-sm">
          {{ $t('system.labels.editor.automations.selectLabel') }}
        </label>
        <Select
          v-model="selectedTaq"
          :options="availableTaqs"
          :option-label="taqOptionLabel"
          :placeholder="$t('system.labels.editor.automations.selectPlaceholder')"
          :loading="taqDialogLoading"
          filter
          fluid
          @show="fetchAvailableTaqs"
        />
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="$t('general.label.cancel')" severity="secondary" text size="small" @click="showTaqDialog = false" />
          <Button :label="$t('general.label.save')" size="small" :disabled="!selectedTaq" @click="addTaqLabel" />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { components } from '@cortezaproject/corteza-vue-next'

const { CResourceTable } = components

const { t } = useI18n()
const route = useRoute()

const $ComposeAPI = inject('$ComposeAPI')
const $AutomationAPI = inject('$AutomationAPI')
const $SystemAPI = inject('$SystemAPI')
const $toast = inject('$toast')

const labelName = computed(() => decodeURIComponent(route.params.labelID || ''))

// ──────────────────────────────────────────────
// Cross-app link helpers
// ──────────────────────────────────────────────
function getAppBase(app) {
  const u = new URL(window.location)
  // In production: /compose/, /agentic/, /taq/
  // In dev: different ports, same origin with /
  return `${u.origin}/${app}/`
}

function buildNamespaceLink(ns) {
  const slug = ns.slug || ns.namespaceID
  return `${getAppBase('compose')}namespace/${slug}`
}

function buildAgentLink(agent) {
  return `${getAppBase('agentic')}${agent.agentID}`
}

function buildTaqLink(taq) {
  return `${getAppBase('taq')}builder/${taq.automationID}`
}

function navigateTo(url) {
  window.location = url
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

// Add namespace dialog
const showNsDialog = ref(false)
const selectedNamespace = ref(null)
const availableNamespaces = ref([])
const nsDialogLoading = ref(false)

async function fetchAvailableNamespaces() {
  nsDialogLoading.value = true
  try {
    const result = await $ComposeAPI.namespaceList({ limit: 200 })
    const set = result.set || []
    // Exclude already-assigned
    const assignedIds = new Set(nsItems.value.map(n => n.namespaceID))
    availableNamespaces.value = set.filter(n => !assignedIds.has(n.namespaceID))
  } catch (e) {
    availableNamespaces.value = []
  } finally {
    nsDialogLoading.value = false
  }
}

async function addNamespaceLabel() {
  if (!selectedNamespace.value) return
  try {
    const ns = selectedNamespace.value
    const labels = { ...(ns.labels || {}) }
    labels[labelName.value] = ''
    await $ComposeAPI.namespaceUpdate({
      namespaceID: ns.namespaceID,
      name: ns.name,
      slug: ns.slug,
      enabled: ns.enabled,
      meta: ns.meta,
      labels,
    })
    showNsDialog.value = false
    selectedNamespace.value = null
    $toast?.toastSuccess(t('system.labels.editor.toast.added'))
    fetchNamespaces()
  } catch (e) {
    console.error('Failed to add label to namespace:', e)
    $toast?.toastDanger(t('system.labels.editor.toast.addError'))
  }
}

// ──────────────────────────────────────────────
// AGENTS
// ──────────────────────────────────────────────
const agentItems = ref([])
const agentLoading = ref(false)

const agentFields = [
  { key: 'name', header: t('system.labels.editor.agents.columns.name') },
  { key: 'status', header: t('system.labels.editor.agents.columns.status') },
]

function agentOptionLabel(agent) {
  return agent.meta?.short || agent.handle || agent.agentID
}

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

// Add agent dialog
const showAgentDialog = ref(false)
const selectedAgent = ref(null)
const availableAgents = ref([])
const agentDialogLoading = ref(false)

async function fetchAvailableAgents() {
  agentDialogLoading.value = true
  try {
    const result = await $SystemAPI.agentList({ limit: 200 })
    const set = result.set || []
    const assignedIds = new Set(agentItems.value.map(a => a.agentID))
    availableAgents.value = set.filter(a => !assignedIds.has(a.agentID))
  } catch (e) {
    availableAgents.value = []
  } finally {
    agentDialogLoading.value = false
  }
}

async function addAgentLabel() {
  if (!selectedAgent.value) return
  try {
    const agent = selectedAgent.value
    const labels = { ...(agent.labels || {}) }
    labels[labelName.value] = ''
    await $SystemAPI.agentUpdate({
      agentID: agent.agentID,
      handle: agent.handle,
      status: agent.status,
      meta: agent.meta,
      behavior: agent.behavior,
      execution: agent.execution,
      access: agent.access,
      invocation: agent.invocation,
      labels,
    })
    showAgentDialog.value = false
    selectedAgent.value = null
    $toast?.toastSuccess(t('system.labels.editor.toast.added'))
    fetchAgents()
  } catch (e) {
    console.error('Failed to add label to agent:', e)
    $toast?.toastDanger(t('system.labels.editor.toast.addError'))
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

function taqOptionLabel(taq) {
  return taq.meta?.short || taq.handle || taq.automationID
}

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

// Add TAQ dialog
const showTaqDialog = ref(false)
const selectedTaq = ref(null)
const availableTaqs = ref([])
const taqDialogLoading = ref(false)

async function fetchAvailableTaqs() {
  taqDialogLoading.value = true
  try {
    const result = await $AutomationAPI.ngAutomationList({ limit: 200 })
    const set = result.set || []
    const assignedIds = new Set(taqItems.value.map(a => a.automationID))
    availableTaqs.value = set.filter(a => !assignedIds.has(a.automationID))
  } catch (e) {
    availableTaqs.value = []
  } finally {
    taqDialogLoading.value = false
  }
}

async function addTaqLabel() {
  if (!selectedTaq.value) return
  try {
    const taq = selectedTaq.value
    const labels = { ...(taq.labels || {}) }
    labels[labelName.value] = ''
    await $AutomationAPI.ngAutomationUpdate({
      automationID: taq.automationID,
      handle: taq.handle,
      labels,
      meta: taq.meta,
      enabled: taq.enabled,
      triggers: taq.triggers,
      steps: taq.steps,
      paths: taq.paths,
      runAs: taq.runAs,
      ownedBy: taq.ownedBy,
    })
    showTaqDialog.value = false
    selectedTaq.value = null
    $toast?.toastSuccess(t('system.labels.editor.toast.added'))
    fetchAutomations()
  } catch (e) {
    console.error('Failed to add label to automation:', e)
    $toast?.toastDanger(t('system.labels.editor.toast.addError'))
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
