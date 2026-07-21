<template>
  <Teleport to="#topbar-title" defer>
    <span class="flex items-center gap-2">
      <span>{{ project?.name || $t('project.wizard.fallbackName') }}</span>
      <Tag
        v-if="project"
        :value="versionLabel"
        :severity="isLive ? 'secondary' : 'warn'"
        class="!text-xs"
      />
    </span>
  </Teleport>

  <!-- Right-aligned topbar tools: for a live project, jump to its dashboard;
       always offer opening the project's compose namespace. -->
  <Teleport to="#topbar-tools" defer>
    <span class="flex items-center gap-2">
      <Button
        v-if="isLive"
        :label="$t('project.viewDashboard')"
        icon="pi pi-gauge"
        size="small"
        severity="secondary"
        outlined
        @click="goDashboard"
      />
      <Button
        v-if="project?.hasNamespace"
        :label="$t('project.viewProject')"
        icon="pi pi-external-link"
        size="small"
        @click="openProject"
      />
    </span>
  </Teleport>

  <div v-if="project" class="h-full flex flex-col min-h-0">
    <div class="flex-1 flex gap-4 p-3 min-h-0">
      <!-- Left: step nav -->
      <div class="w-72 shrink-0 flex flex-col gap-2 min-h-0">
        <StepNav
          class="flex-1 min-h-0"
          :steps="navSteps"
          :active-key="activeKey"
          :statuses="statuses"
          @select="goStep"
        />
      </div>

      <!-- Step panel — an outlined panel on the background (border + rounded, no
           fill); the step header's bottom border delineates it from the content. -->
      <div
        class="flex-1 flex flex-col min-w-0 min-h-0 overflow-hidden rounded-xl border border-surface"
      >
        <!-- Step header -->
        <div class="shrink-0 border-b border-surface px-4 py-3 flex items-center gap-3">
          <!-- Leading badge mirrors the sidebar/metrics strip: resource steps take
               their kind's icon and colour, other steps a neutral badge + step icon. -->
          <span
            v-if="activeStep"
            class="inline-flex items-center justify-center w-9 h-9 rounded-md shrink-0"
            :class="activeBadge.wrap"
          >
            <i :class="activeBadge.icon" />
          </span>
          <div class="min-w-0">
            <h2 class="text-lg font-medium truncate">
              {{ activeStep ? $t(activeStep.labelKey) : '' }}
            </h2>
            <p class="text-sm text-muted-color">{{ headerHint }}</p>
          </div>

          <div class="ml-auto flex items-center gap-3">
            <Tag v-if="showStatus" :value="statusLabel" :severity="statusSeverity" />
          </div>
        </div>

        <!-- Step content + live resource panel -->
        <div ref="splitRef" class="flex-1 min-h-0 flex">
          <div
            class="min-h-0"
            :class="isResourceStep ? 'overflow-hidden flex flex-col' : 'overflow-y-auto p-4'"
            :style="{ width: leftPct + '%' }"
          >
            <template v-if="isResourceStep">
              <StepStatusBanner
                v-if="showStatus"
                :status="status"
                :review-note="reviewNote"
                class="m-3 mb-0 shrink-0"
              />
              <DataModelStep
                v-if="isDataModel"
                :project="project"
                :disabled="locked"
                class="flex-1 min-h-0"
              />
              <ConnectionsStep
                v-else-if="isConnections"
                :project="project"
                :disabled="locked"
                class="flex-1 min-h-0"
              />
              <AutomationsStep
                v-else-if="isAutomations"
                :project="project"
                :disabled="locked"
                class="flex-1 min-h-0"
              />
              <AgentsStep
                v-else-if="isAgents"
                :project="project"
                :disabled="locked"
                class="flex-1 min-h-0"
              />
              <ChatbotsStep
                v-else-if="isChatbots"
                :project="project"
                :disabled="locked"
                class="flex-1 min-h-0"
              />
              <PagesStep
                v-else-if="isPages"
                :project="project"
                :disabled="locked"
                class="flex-1 min-h-0"
              />
              <RolesStep
                v-else-if="isRoles"
                :project="project"
                :disabled="locked"
                class="flex-1 min-h-0"
              />
              <PermissionsStep
                v-else-if="isPermissions"
                :project="project"
                :disabled="locked"
                class="flex-1 min-h-0"
              />
              <UsersStep v-else :project="project" :disabled="locked" class="flex-1 min-h-0" />
            </template>
            <template v-else>
              <StepStatusBanner
                v-if="showStatus"
                :status="status"
                :review-note="reviewNote"
                class="mb-4"
              />
              <ProjectSummaryStep v-if="isSummary" v-model="working" :disabled="locked" />
              <ResourceManagementStep
                v-else-if="isResourceMgmt"
                v-model="working"
                :disabled="locked"
              />
              <MembersStep v-else-if="isMembers" :project="project" :disabled="membersLocked" />
              <DataSensitivityStep
                v-else-if="isSensitivity"
                :project="project"
                :disabled="locked"
              />
              <PublishStep
                v-else-if="isPublish"
                :project="project"
                :publishing="publishing"
                :can-request="canRequest"
                :can-grant="canGrant"
                @publish="onPublish"
                @open-dashboard="goDashboard"
              />
            </template>
          </div>

          <!-- Drag handle -->
          <div
            class="shrink-0 w-1.5 bg-emphasis hover:bg-primary relative cursor-col-resize group select-none flex items-center justify-center transition-colors"
            :class="{ '!bg-primary': resizing }"
            @pointerdown="startResize"
          >
            <span class="absolute inset-y-0 -left-1 -right-1" />
            <span
              class="h-8 w-0.5 rounded-full bg-surface-400 dark:bg-surface-500 group-hover:bg-primary-contrast"
              :class="{ '!bg-primary-contrast': resizing }"
            />
          </div>

          <div class="overflow-hidden p-4 min-h-0 flex-1">
            <ResourceGraph :project="project" :locked="status === 'submitted'" />
          </div>
        </div>
      </div>
    </div>

    <WizardToolbar
      v-if="showToolbar"
      :can-write="canWrite"
      :can-grant="canGrant"
      :mode="project.mode"
      :is-publish="isPublish"
      :show-save="showStepSave"
      :can-prev="canPrev"
      :can-next="canNext"
      :step-index="stepIndex"
      :step-count="navSteps.length"
      @save="onSave"
      @request-changes="openRequestChanges"
      @prev="onPrev"
      @next="onNext"
      @back="onBack"
    />

    <!-- Per-resource dialogs — every kind's Create + Detail dialog is mounted
         here exactly once and opened through the inspectResource/createResource
         provides. Steps and the resource graph never mount their own. -->

    <!-- Module -->
    <ModuleCreateDialog
      :model-value="createKind === 'module'"
      :project="project"
      @update:model-value="onCreateToggle('module', $event)"
      @created="onMutated"
    />
    <ModuleDetailDialog
      :model-value="detailKind === 'module'"
      :project="project"
      :resource-id="detailId"
      @update:model-value="onDetailToggle('module', $event)"
      @saved="onMutated"
    />

    <!-- Connection -->
    <ConnectionCreateDialog
      :model-value="createKind === 'connection'"
      :project="project"
      @update:model-value="onCreateToggle('connection', $event)"
      @created="onMutated"
    />
    <ConnectionDetailDialog
      :model-value="detailKind === 'connection'"
      :project="project"
      :resource-id="detailId"
      @update:model-value="onDetailToggle('connection', $event)"
      @saved="onMutated"
    />

    <!-- Automation -->
    <AutomationCreateDialog
      :model-value="createKind === 'automation'"
      :project="project"
      @update:model-value="onCreateToggle('automation', $event)"
      @created="onMutated"
    />
    <AutomationDetailDialog
      :model-value="detailKind === 'automation'"
      :project="project"
      :resource-id="detailId"
      @update:model-value="onDetailToggle('automation', $event)"
      @saved="onMutated"
    />

    <!-- Agent -->
    <AgentCreateDialog
      :model-value="createKind === 'agent'"
      :project="project"
      @update:model-value="onCreateToggle('agent', $event)"
      @created="onMutated"
    />
    <AgentDetailDialog
      :model-value="detailKind === 'agent'"
      :project="project"
      :resource-id="detailId"
      @update:model-value="onDetailToggle('agent', $event)"
      @saved="onMutated"
    />

    <!-- Chatbot -->
    <ChatbotCreateDialog
      :model-value="createKind === 'chatbot'"
      :project="project"
      @update:model-value="onCreateToggle('chatbot', $event)"
      @created="onMutated"
    />
    <ChatbotDetailDialog
      :model-value="detailKind === 'chatbot'"
      :project="project"
      :resource-id="detailId"
      @update:model-value="onDetailToggle('chatbot', $event)"
      @saved="onMutated"
    />

    <!-- Page -->
    <PageCreateDialog
      :model-value="createKind === 'page'"
      :project="project"
      @update:model-value="onCreateToggle('page', $event)"
      @created="onMutated"
    />
    <PageDetailDialog
      :model-value="detailKind === 'page'"
      :project="project"
      :resource-id="detailId"
      @update:model-value="onDetailToggle('page', $event)"
      @saved="onMutated"
    />

    <!-- Role -->
    <RoleCreateDialog
      :model-value="createKind === 'role'"
      :project="project"
      @update:model-value="onCreateToggle('role', $event)"
      @created="onMutated"
    />
    <RoleDetailDialog
      :model-value="detailKind === 'role'"
      :project="project"
      :resource-id="detailId"
      @update:model-value="onDetailToggle('role', $event)"
      @saved="onMutated"
    />

    <!-- User -->
    <UserCreateDialog
      :model-value="createKind === 'user'"
      :project="project"
      @update:model-value="onCreateToggle('user', $event)"
      @created="onMutated"
    />
    <UserDetailDialog
      :model-value="detailKind === 'user'"
      :project="project"
      :resource-id="detailId"
      @update:model-value="onDetailToggle('user', $event)"
      @saved="onMutated"
    />

    <!-- Single-field editor — a Wizard-level sub-dialog stacked over the module
         detail dialog; opened via the editField/createField provides. -->
    <FieldDialog
      v-model="fieldOpen"
      :project="project"
      :module-id="fieldModuleId"
      :field-id="fieldId"
      :readonly="locked"
    />
  </div>

  <!-- Per-step "Request changes" — a granter can flag the active step (any
       step but Publish, see WizardToolbar) at any time with a required note.
       Mirrors PublishStep.vue's own request-changes dialog. -->
  <Dialog
    v-model:visible="requestChanges.visible"
    modal
    :header="requestChangesHeader"
    :style="{ width: '32rem' }"
  >
    <CFormGroup :label="$t('project.wizard.requestChanges.dialog.label')" required>
      <Textarea
        v-model="requestChanges.note"
        rows="3"
        auto-resize
        fluid
        :placeholder="$t('project.wizard.requestChanges.dialog.placeholder')"
      />
    </CFormGroup>
    <template #footer>
      <div class="flex justify-end gap-2">
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          text
          size="small"
          @click="requestChanges.visible = false"
        />
        <Button
          :label="$t('project.wizard.requestChanges.dialog.confirm')"
          size="small"
          :loading="requestChanges.sending"
          :disabled="!requestChanges.note.trim()"
          @click="confirmRequestChanges"
        />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import FieldDialog from '@/sections/project/components/datamodel/FieldDialog.vue'
import ModuleCreateDialog from '@/sections/project/components/datamodel/ModuleCreateDialog.vue'
import ModuleDetailDialog from '@/sections/project/components/datamodel/ModuleDetailDialog.vue'
import ConnectionCreateDialog from '@/sections/project/components/connections/ConnectionCreateDialog.vue'
import ConnectionDetailDialog from '@/sections/project/components/connections/ConnectionDetailDialog.vue'
import AutomationCreateDialog from '@/sections/project/components/automations/AutomationCreateDialog.vue'
import AutomationDetailDialog from '@/sections/project/components/automations/AutomationDetailDialog.vue'
import AgentCreateDialog from '@/sections/project/components/agents/AgentCreateDialog.vue'
import AgentDetailDialog from '@/sections/project/components/agents/AgentDetailDialog.vue'
import ChatbotCreateDialog from '@/sections/project/components/chatbots/ChatbotCreateDialog.vue'
import ChatbotDetailDialog from '@/sections/project/components/chatbots/ChatbotDetailDialog.vue'
import PageCreateDialog from '@/sections/project/components/pages/PageCreateDialog.vue'
import PageDetailDialog from '@/sections/project/components/pages/PageDetailDialog.vue'
import RoleCreateDialog from '@/sections/project/components/roles/RoleCreateDialog.vue'
import RoleDetailDialog from '@/sections/project/components/roles/RoleDetailDialog.vue'
import UserCreateDialog from '@/sections/project/components/users/UserCreateDialog.vue'
import UserDetailDialog from '@/sections/project/components/users/UserDetailDialog.vue'
import ResourceGraph from '@/sections/project/components/graph/ResourceGraph.vue'
import StepNav from '@/sections/project/components/wizard/StepNav.vue'
import StepStatusBanner from '@/sections/project/components/wizard/StepStatusBanner.vue'
import WizardToolbar from '@/sections/project/components/wizard/WizardToolbar.vue'
import AgentsStep from '@/sections/project/components/wizard/steps/AgentsStep.vue'
import AutomationsStep from '@/sections/project/components/wizard/steps/AutomationsStep.vue'
import ChatbotsStep from '@/sections/project/components/wizard/steps/ChatbotsStep.vue'
import ConnectionsStep from '@/sections/project/components/wizard/steps/ConnectionsStep.vue'
import PagesStep from '@/sections/project/components/wizard/steps/PagesStep.vue'
import RolesStep from '@/sections/project/components/wizard/steps/RolesStep.vue'
import PermissionsStep from '@/sections/project/components/wizard/steps/PermissionsStep.vue'
import UsersStep from '@/sections/project/components/wizard/steps/UsersStep.vue'
import DataModelStep from '@/sections/project/components/wizard/steps/DataModelStep.vue'
import DataSensitivityStep from '@/sections/project/components/wizard/steps/DataSensitivityStep.vue'
import MembersStep from '@/sections/project/components/wizard/steps/MembersStep.vue'
import ProjectSummaryStep from '@/sections/project/components/wizard/steps/ProjectSummaryStep.vue'
import ResourceManagementStep from '@/sections/project/components/wizard/steps/ResourceManagementStep.vue'
import PublishStep from '@/sections/project/components/wizard/steps/PublishStep.vue'
import { ACCESS_KINDS, kindConfig } from '@/sections/project/config/kinds'
import { STEPS, kindsThroughStep, stepsForTab } from '@/sections/project/config/pipeline'
import { resourceManagementValues } from '@/sections/project/config/resourceManagementForm'
import { rolePreset } from '@/sections/project/config/roles'
import { summaryDefaults } from '@/sections/project/config/summaryForm'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { useProjectUsersStore } from '@/sections/project/stores/users'
import { computed, inject, provide, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const store = useProjectsStore()
const usersStore = useProjectUsersStore()
const $toast = inject('$toast')

const project = computed(() => store.findById(route.params.projectId))

// A live (published) project has a dashboard to switch to; drafts are wizard-only.
const isLive = computed(() => ['active', 'published'].includes(project.value?.status))

// User-facing versions are 1-based (the original live project is v1), so we
// display the backend revision + 1 — mirrors the dashboard topbar crumb. An
// unpublished project is flagged as a draft (e.g. "v1 draft").
const versionLabel = computed(() =>
  t(isLive.value ? 'project.dashboard.version' : 'project.dashboard.versionDraft', {
    number: (project.value?.revision ?? 0) + 1,
  }),
)

// Load the full project (members for capability resolution) + user directory.
usersStore.load()
watch(
  () => route.params.projectId,
  id => {
    if (!id) return
    store.fetchProject(id).catch(err => {
      console.error('Failed to load project', err)
      $toast.toastErrorHandler(t('project.wizard.toastLoadFailed'))(err)
    })
  },
  { immediate: true },
)

// --- Current member & role -------------------------------------------------
// Tab access follows the logged-in user's role: a member who can request OR
// grant approval works in the Governance view; everyone else sees Build only.
const currentMember = computed(() =>
  store.membersFor(project.value?.projectID).find(m => m.userId === usersStore.currentUserID),
)
const currentRole = computed(() => rolePreset(currentMember.value?.role))
const canWrite = computed(() => !!currentRole.value.write)
const canRequest = computed(() => !!currentRole.value.requestApproval)
const canGrant = computed(() => !!currentRole.value.grantApproval)
const canGovern = computed(() => canRequest.value || canGrant.value)

// --- Tabs (derived, not chosen) --------------------------------------------
const effectiveTab = computed(() =>
  project.value?.mode === 'gated' && canGovern.value ? 'governance' : 'build',
)
const navSteps = computed(() => (project.value ? stepsForTab(effectiveTab.value) : []))

// --- Active step -----------------------------------------------------------
const activeKey = computed(() => {
  const list = navSteps.value
  if (!list.length) return ''
  const q = route.query.step
  if (q && list.some(s => s.key === q)) return q
  return list[0].key
})
const activeStep = computed(
  () =>
    navSteps.value.find(s => s.key === activeKey.value) ||
    STEPS.find(s => s.key === activeKey.value),
)
const isSummary = computed(() => activeKey.value === 'summary')
const isResourceMgmt = computed(() => activeKey.value === 'resource-management')
const isMembers = computed(() => activeKey.value === 'members')
const isDataModel = computed(() => activeKey.value === 'data-model')
const isConnections = computed(() => activeKey.value === 'connections')
const isAutomations = computed(() => activeKey.value === 'automations')
const isAgents = computed(() => activeKey.value === 'agents')
const isChatbots = computed(() => activeKey.value === 'chatbots')
const isPages = computed(() => activeKey.value === 'pages')
const isRoles = computed(() => activeKey.value === 'roles')
const isPermissions = computed(() => activeKey.value === 'permissions')
const isUsers = computed(() => activeKey.value === 'users')
const isPublish = computed(() => activeKey.value === 'publish')
// Resource steps render full-height with the live resource graph beside them.
const isResourceStep = computed(
  () =>
    isDataModel.value ||
    isConnections.value ||
    isAutomations.value ||
    isAgents.value ||
    isChatbots.value ||
    isPages.value ||
    isRoles.value ||
    isPermissions.value ||
    isUsers.value,
)
const isSensitivity = computed(() => activeKey.value === 'data-sensitivity')
// The graph always shows the whole-system overview for now. To narrow it to the
// active resource step's kind again, restore this computed and pass it back as
// `:filter-kind="stepKind"` on <ResourceGraph>.
// const stepKind = computed(() =>
//   activeStep.value?.type === 'resource' ? activeStep.value.kind : null,
// )

// --- Per-resource Create / Detail dialogs ------------------------------------
// One shared contract: every resource kind owns a Create dialog (minimal new-
// resource form) and a Detail dialog (full edit of an existing resource), all
// mounted once above. Steps and the resource graph open them through these two
// provides; only one create dialog and one detail dialog are ever open at a
// time, keyed by the active kind.
const createKind = ref(null)
const detailKind = ref(null)
const detailId = ref(null)

// createResource(kind) — open the Create dialog for that kind.
provide('createResource', kind => {
  detailKind.value = null
  createKind.value = kind
})
// inspectResource(kind, id) — open the Detail dialog for an existing resource.
// Graph node clicks and step-row clicks both route here; every kind opens
// editable (no readonly mode).
provide('inspectResource', (kind, id) => {
  createKind.value = null
  detailId.value = id
  detailKind.value = kind
})

// Dialog v-model close handlers: clear the active kind when a dialog closes.
function onCreateToggle(kind, open) {
  if (!open && createKind.value === kind) createKind.value = null
}
function onDetailToggle(kind, open) {
  if (!open && detailKind.value === kind) {
    detailKind.value = null
    detailId.value = null
  }
}

// Both create and detail dialogs already refetch their list inside the store
// action (which touch()es), so the graph refreshes on its own. Bumping the
// graph version again on the emit keeps the refresh loop driven from one place
// regardless of which action ran.
function onMutated() {
  store.touch()
}

// --- Single-field edit dialog ------------------------------------------------
// Provided to the step components: open the editor for one field.
const fieldOpen = ref(false)
const fieldModuleId = ref(null)
const fieldId = ref(null)
provide('editField', (moduleId, id) => {
  fieldModuleId.value = moduleId
  fieldId.value = id
  fieldOpen.value = true
})
// Open the field editor in create mode for the given module (null fieldId).
provide('createField', moduleId => {
  fieldModuleId.value = moduleId
  fieldId.value = null
  fieldOpen.value = true
})
watch(fieldOpen, open => {
  if (!open) {
    fieldModuleId.value = null
    fieldId.value = null
  }
})

// --- Resizable split (step | resource graph) ---------------------------------
const splitRef = ref(null)
const leftPct = ref(40)
const resizing = ref(false)
const MIN_PCT = 25
const MAX_PCT = 75

const onResizeMove = e => {
  const el = splitRef.value
  if (!el) return
  const rect = el.getBoundingClientRect()
  const pct = ((e.clientX - rect.left) / rect.width) * 100
  leftPct.value = Math.min(MAX_PCT, Math.max(MIN_PCT, pct))
}

const stopResize = () => {
  resizing.value = false
  window.removeEventListener('pointermove', onResizeMove)
  window.removeEventListener('pointerup', stopResize)
  document.body.style.userSelect = ''
}

const startResize = () => {
  resizing.value = true
  document.body.style.userSelect = 'none'
  window.addEventListener('pointermove', onResizeMove)
  window.addEventListener('pointerup', stopResize)
}

// --- Per-step governance state --------------------------------------------
// Per-section approval gates are gone: a non-publish step only ever carries
// 'draft' (the default) or 'changes-requested' (raised anytime by a granter
// via the per-step "Request changes" action below) — 'submitted'/'approved'
// belong to the Publish step alone, which drives its own dedicated cycle.
const stepStatus = key => project.value?.governance?.[key]?.status || 'draft'
const status = computed(() => stepStatus(activeKey.value))
const reviewNote = computed(() => project.value?.governance?.[activeKey.value]?.reviewNote || '')
const locked = computed(
  () => !canWrite.value || status.value === 'submitted' || status.value === 'approved',
)
// Member management is RBAC-checked (project members.manage), not the role
// preset's write flag.
const membersLocked = computed(
  () =>
    !project.value?.canManageMembers || status.value === 'submitted' || status.value === 'approved',
)
// The Publish step renders its own dedicated approval panel (see
// PublishStep.vue) with wording specific to publish-time approval; the generic
// header badge + StepStatusBanner (which only ever has a
// "changes-requested" note to show now) are suppressed for it.
const showStatus = computed(
  () => project.value?.mode === 'gated' && !!activeStep.value && activeKey.value !== 'publish',
)
const showToolbar = computed(() => !!activeStep.value)
// Save only makes sense on form steps — members/sensitivity/resource steps
// persist each change immediately through their own store calls instead.
const showStepSave = computed(() => activeStep.value?.type === 'form')

const statuses = computed(() => {
  const out = {}
  for (const s of navSteps.value) out[s.key] = stepStatus(s.key)
  return out
})

const statusLabel = computed(
  () =>
    ({
      draft: t('project.governance.status.draft'),
      submitted: t('project.governance.status.submitted'),
      approved: t('project.governance.status.approved'),
      'changes-requested': t('project.governance.status.changesRequested'),
    })[status.value],
)
const statusSeverity = computed(
  () =>
    ({ draft: 'secondary', submitted: 'info', approved: 'success', 'changes-requested': 'warn' })[
      status.value
    ],
)

// Leading icon badge for the active step's header — resolved exactly like the
// sidebar (StepNav): resource steps borrow their kind's icon/colour, other steps
// fall back to a neutral badge carrying the step's own icon.
const activeBadge = computed(() => {
  const s = activeStep.value
  if (s?.kind) {
    const cfg = kindConfig(s.kind)
    return { wrap: ['ring-1', cfg.bg, cfg.ring], icon: [cfg.icon, cfg.text] }
  }
  // Permissions step mirrors the app-wide permissions button (CPermissionsButton):
  // outlined secondary — surface fill, border-surface (--p-content-border-color,
  // the token the outlined-secondary Button uses), full-tone lock.
  if (s?.type === 'permissions') {
    return {
      wrap: ['border', 'border-surface', 'bg-surface'],
      icon: ['pi', s?.icon || 'pi-lock', 'text-color'],
    }
  }
  return {
    wrap: ['ring-1', 'bg-emphasis', 'ring-surface'],
    icon: ['pi', s?.icon || 'pi-circle', 'text-muted-color'],
  }
})

// Short per-step blurb shown in the header (so steps don't repeat it in-body).
const STEP_BLURB = {
  summary: 'project.wizard.blurb.summary',
  'resource-management': 'project.wizard.blurb.resourceManagement',
  members: 'project.wizard.blurb.members',
  'data-model': 'project.wizard.blurb.dataModel',
  'data-sensitivity': 'project.wizard.blurb.dataSensitivity',
  connections: 'project.wizard.blurb.connections',
  automations: 'project.wizard.blurb.automations',
  agents: 'project.wizard.blurb.agents',
  chatbots: 'project.wizard.blurb.chatbots',
  pages: 'project.wizard.blurb.pages',
  roles: 'project.wizard.blurb.roles',
  permissions: 'project.wizard.blurb.permissions',
  users: 'project.wizard.blurb.users',
  publish: 'project.wizard.blurb.publish',
}
// Per-section approval gates are gone, so no step but Publish carries real
// governance status any more (every other step key stays in `draft` forever)
// — the header hint is always just the step's own blurb.
const headerHint = computed(() => {
  const blurbKey = STEP_BLURB[activeKey.value]
  return blurbKey ? t(blurbKey) : ''
})

// --- Working copy of the active form step's values -------------------------
const working = ref({})
function loadWorking() {
  const saved = project.value?.governance?.[activeKey.value]?.values || {}
  if (isSummary.value) {
    const vals = { ...summaryDefaults(), ...saved }
    if (!vals.systemName) vals.systemName = project.value?.name || ''
    working.value = vals
  } else if (isResourceMgmt.value) {
    working.value = resourceManagementValues(saved)
  } else {
    working.value = {}
  }
}
watch([activeKey, project], loadWorking, { immediate: true })

// Scope the live resource graph to the process: each step shows only the
// resources of steps reached so far. Re-seeded on every step change; the graph's
// own layer chips still refine (or peek past) it within a step.
watch(
  activeKey,
  key => {
    store.setGraphVisibleKinds(kindsThroughStep(key))
    // Auto-reveal the role/user access overlay while on the access steps (roles,
    // permissions, users), so the roles you're managing — and the grant edges
    // you're drawing — show up in the graph beside the form.
    const step = STEPS.find(s => s.key === key)
    store.setGraphShowAccess(ACCESS_KINDS.includes(step?.kind) || key === 'permissions')
  },
  { immediate: true },
)

function goStep(key) {
  router.replace({ query: { ...route.query, step: key } })
}

// Leave the wizard and return to the project list.
function onBack() {
  router.push({ name: 'project.list' })
}

// Open the project's compose namespace (resolved by namespaceID).
function openProject() {
  if (project.value?.hasNamespace) {
    router.push({ name: 'namespace.view', params: { slug: project.value.namespaceID } })
  }
}

// --- Publish ---------------------------------------------------------------
// First publish promotes the draft to active (no record migration); on success
// we land on the project's dashboard. Guarded against double-submit.
const publishing = ref(false)
function goDashboard() {
  router.push({ name: 'project.overview', params: { projectId: project.value.projectID } })
}
async function onPublish() {
  if (publishing.value || !project.value) return
  publishing.value = true
  try {
    await store.publishProject(project.value.projectID)
    $toast.toastSuccess(t('project.publishStep.toast.success'))
    goDashboard()
  } catch (err) {
    $toast.toastErrorHandler(t('project.publishStep.toast.failed'))(err)
  } finally {
    publishing.value = false
  }
}

// --- Prev / Next stepper ---------------------------------------------------
// Walk the role-visible step list (navSteps is already filtered by Build vs
// Governance). The arrows disable at the ends.
const stepPos = computed(() => navSteps.value.findIndex(s => s.key === activeKey.value))
const stepIndex = computed(() => stepPos.value + 1) // 1-based, for display
const canPrev = computed(() => stepPos.value > 0)
const canNext = computed(() => stepPos.value < navSteps.value.length - 1)

function onPrev() {
  if (canPrev.value) goStep(navSteps.value[stepPos.value - 1].key)
}
function onNext() {
  if (canNext.value) goStep(navSteps.value[stepPos.value + 1].key)
}

// --- Per-step actions ------------------------------------------------------
// Governance mutations persist via the API; report failures instead of
// assuming success.
async function governanceAction(fn, summary) {
  try {
    await fn()
    $toast.toastSuccess(summary)
  } catch (err) {
    $toast.toastErrorHandler(t('project.wizard.toast.actionFailed'))(err)
  }
}
function onSave() {
  governanceAction(
    () => store.saveStepForm(project.value.projectID, activeKey.value, working.value),
    t('project.wizard.toast.saved'),
  )
}

// --- Per-step "Request changes" ---------------------------------------------
// Gated mode only: a member with grant-approval capability can flag the
// active step (any step but Publish — see WizardToolbar's showRequestChanges)
// at any time, with a required note. The backend flags that step to
// changes-requested regardless of its current status, and sends the Publish
// step back for review too if it was pending (see stores/projects.js
// transitionStep + server/system/service/project_governance.go). This dialog
// reuses the single-required-note pattern the retired per-step approve/reopen
// flow used, and mirrors PublishStep.vue's own request-changes dialog.
const requestChanges = ref({ visible: false, note: '', sending: false })
const requestChangesHeader = computed(() =>
  t('project.wizard.requestChanges.dialog.header', {
    step: activeStep.value ? t(activeStep.value.labelKey) : '',
  }),
)
function openRequestChanges() {
  requestChanges.value = { visible: true, note: '', sending: false }
}
async function confirmRequestChanges() {
  const note = requestChanges.value.note.trim()
  if (!note || requestChanges.value.sending) return
  requestChanges.value.sending = true
  try {
    await store.transitionStep(project.value.projectID, activeKey.value, 'request-changes', note)
    requestChanges.value.visible = false
    $toast.toastInfo(t('project.wizard.requestChanges.toast.sent'))
  } catch (err) {
    $toast.toastErrorHandler(t('project.wizard.toast.actionFailed'))(err)
  } finally {
    requestChanges.value.sending = false
  }
}
</script>
