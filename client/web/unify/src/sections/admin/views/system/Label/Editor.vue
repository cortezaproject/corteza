<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.labels.editor.title') }} — {{ labelName }}</span>
  </Teleport>

  <div class="flex flex-col h-full">
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <!-- Namespaces with this label -->
      <Panel :header="$t('system.labels.editor.namespaces.title')" toggleable>
        <CResourceTable
          primary-key="namespaceID"
          :fields="nsFields"
          :items="nsItems"
          :loading="nsLoading"
          :empty-message="$t('system.labels.editor.namespaces.empty')"
          :action-items="getNsActions"
          :row-class="() => 'cursor-pointer'"
          @row-click="({ data }) => openInNewTab('compose', `namespace/${data.slug || data.namespaceID}`)"
        >
          <template #header>
            <div class="flex gap-2">
              <Button
                :label="$t('system.labels.editor.namespaces.create')"
                icon="pi pi-plus"
                size="small"
                @click="showNsDialog = true"
              />

              <CPermissionsButton
                v-if="canGrantCompose"
                v-tooltip.bottom="$t('general.label.permissions')"
                resource="corteza::compose:namespace/*"
              />
            </div>
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
          :row-class="() => 'cursor-pointer'"
          @row-click="({ data }) => openInNewTab('agentic', `${data.agentID}/edit`)"
        >
          <template #header>
            <div class="flex gap-2">
              <Button
                :label="$t('system.labels.editor.agents.create')"
                icon="pi pi-plus"
                size="small"
                @click="showAgentDialog = true"
              />

              <CPermissionsButton
                v-if="canGrantSystem"
                v-tooltip.bottom="$t('general.label.permissions')"
                resource="corteza::system:agent/*"
              />
            </div>
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
          :row-class="() => 'cursor-pointer'"
          @row-click="({ data }) => openInNewTab('taq', `builder/${data.automationID}`)"
        >
          <template #header>
            <div class="flex gap-2">
              <Button
                :label="$t('system.labels.editor.automations.create')"
                icon="pi pi-plus"
                size="small"
                @click="showTaqDialog = true"
              />

              <CPermissionsButton
                v-if="canGrantAutomation"
                v-tooltip.bottom="$t('general.label.permissions')"
                resource="corteza::automation:ng-automation/*"
              />
            </div>
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

    <CEditorActions :back-to="{ name: 'system.labels' }" />

    <!-- Create Namespace Dialog -->
    <Dialog
      v-model:visible="showNsDialog"
      :header="$t('system.labels.editor.namespaces.createDialog')"
      modal
      :style="{ width: '32rem' }"
    >
      <div class="flex flex-col gap-4">
        <CFormGroup :label="$t('system.labels.editor.namespaces.nameLabel')">
          <InputText
            v-model="nsForm.name"
            :placeholder="$t('system.labels.editor.namespaces.namePlaceholder')"
            fluid
          />
        </CFormGroup>
        <CFormGroup :label="$t('system.labels.editor.namespaces.slugLabel')">
          <InputText
            v-model="nsForm.slug"
            :placeholder="$t('system.labels.editor.namespaces.slugPlaceholder')"
            fluid
          />
        </CFormGroup>
      </div>
      <template #footer>
        <div class="flex justify-end gap-2">
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
        <CFormGroup :label="$t('system.labels.editor.agents.nameLabel')">
          <InputText
            v-model="agentForm.name"
            :placeholder="$t('system.labels.editor.agents.namePlaceholder')"
            fluid
          />
        </CFormGroup>
      </div>
      <template #footer>
        <div class="flex justify-end gap-2">
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
        <CFormGroup :label="$t('system.labels.editor.automations.nameLabel')">
          <InputText
            v-model="taqForm.name"
            :placeholder="$t('system.labels.editor.automations.namePlaceholder')"
            fluid
          />
        </CFormGroup>
      </div>
      <template #footer>
        <div class="flex justify-end gap-2">
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
import { components, useConfirmDelete, usePermissions, useRBACStore } from '@planetcrust/human-vue'

const { CResourceTable } = components

const { t } = useI18n()
const route = useRoute()
const { confirmDelete } = useConfirmDelete()
const { open: openPermissions } = usePermissions()
const rbac = useRBACStore()
const canGrantCompose = computed(() => rbac.can('compose/', 'grant'))
const canGrantSystem = computed(() => rbac.can('system/', 'grant'))
const canGrantAutomation = computed(() => rbac.can('automation/', 'grant'))

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
  const items = []

  if (data.canGrant) {
    items.push({
      label: t('general.label.permissions'),
      icon: 'pi pi-lock',
      command: () => openPermissions({
        resource: `corteza::compose:namespace/${data.namespaceID}`,
        title: data.name || data.slug || data.namespaceID,
      }),
    })
  }

  if (data.canDeleteNamespace) {
    if (items.length > 0) items.push({ separator: true })
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
  const items = []

  if (data.canGrant || canGrantSystem.value) {
    items.push({
      label: t('general.label.permissions'),
      icon: 'pi pi-lock',
      command: () => openPermissions({
        resource: `corteza::system:agent/${data.agentID}`,
        title: data.meta?.short || data.handle || data.agentID,
      }),
    })
  }

  if (data.canDeleteAgent) {
    if (items.length > 0) items.push({ separator: true })
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
  const items = []

  if (data.canGrant) {
    items.push({
      label: t('general.label.permissions'),
      icon: 'pi pi-lock',
      command: () => openPermissions({
        resource: `corteza::automation:ng-automation/${data.automationID}`,
        title: data.meta?.short || data.handle || data.automationID,
        target: data.meta?.short || data.handle || data.automationID,
      }),
    })
  }

  if (data.canDeleteNgAutomation) {
    if (items.length > 0) items.push({ separator: true })
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
    $toast?.toastErrorHandler(t('system.labels.editor.toast.deleteError'))(e)
  }
}

async function handleDeleteAgent(data) {
  try {
    await $SystemAPI.agentDelete({ agentID: data.agentID })
    $toast?.toastSuccess(t('system.labels.editor.toast.deleted'))
    fetchAgents()
  } catch (e) {
    console.error('Failed to delete agent:', e)
    $toast?.toastErrorHandler(t('system.labels.editor.toast.deleteError'))(e)
  }
}

async function handleDeleteAutomation(data) {
  try {
    await $AutomationAPI.ngAutomationDelete({ automationID: data.automationID })
    $toast?.toastSuccess(t('system.labels.editor.toast.deleted'))
    fetchAutomations()
  } catch (e) {
    console.error('Failed to delete automation:', e)
    $toast?.toastErrorHandler(t('system.labels.editor.toast.deleteError'))(e)
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
    const result = await $ComposeAPI.namespaceList({ labels: `${labelName.value}=`, limit: 100 })
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
      enabled: true,
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
    $toast?.toastErrorHandler(t('system.labels.editor.toast.createError'))(e)
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
    const result = await $SystemAPI.agentList({ labels: `${labelName.value}=`, limit: 100, sort: 'name ASC' })
    agentItems.value = result.set || []
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
    $toast?.toastErrorHandler(t('system.labels.editor.toast.createError'))(e)
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
    const result = await $AutomationAPI.ngAutomationList({ labels: `${labelName.value}=`, disabled: 1, limit: 100 })
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
    $toast?.toastErrorHandler(t('system.labels.editor.toast.createError'))(e)
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
