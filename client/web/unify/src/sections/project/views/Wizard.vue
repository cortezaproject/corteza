<template>
  <Teleport to="#topbar-title" defer>
    <span class="flex items-center gap-2">
      <span>{{ project?.name || $t('project.wizard.fallbackName') }}</span>
      <!-- Revision switcher — shared with the dashboard topbar (see
           components/project/RevisionSwitcher.vue); replaces what used to be
           an inline trigger + Menu here. -->
      <RevisionSwitcher :project="project" />
    </span>
  </Teleport>

  <!-- Members + View project — shared with the dashboard topbar, right-
       aligned tools slot (see components/project/ProjectTopbarTools.vue);
       owns the Members button, the MembersDialog mount, and the "View
       project" compose namespace link that used to live in the tab-row tool
       cluster below. auto-open-members carries the `?new=1`
       just-created-project deep link (see membersAutoOpen below) into the
       shared dialog. -->
  <Teleport to="#topbar-tools" defer>
    <span class="flex items-center gap-2">
      <ProjectTopbarTools :project="project" :auto-open-members="membersAutoOpen" />
    </span>
  </Teleport>

  <div v-if="project" class="h-full flex flex-col min-h-0">
    <!-- Header row: the fixed three-tab switcher — Build / Govern / Manage &
         Monitor, visible to every member regardless of capability (only the
         per-step and per-project review ACTIONS are capability-gated, not tab
         access) — with the wizard tool cluster right-aligned beside it. -->
    <!-- flex-wrap + ml-auto on the tool cluster: the tools sit beside the
         tabs while they fit and drop to their own right-aligned row when
         they don't (no fixed breakpoint). -->
    <div class="shrink-0 flex flex-wrap items-center justify-between gap-x-3 gap-y-2 px-3 pt-2">
      <Tabs v-model:value="activeTab" class="wizard-tabs shrink-0">
        <TabList>
          <Tab value="build" class="flex items-center gap-2">
            <i class="pi pi-wrench" />
            <span>{{ $t('project.wizard.tabs.build') }}</span>
          </Tab>
          <Tab value="govern" class="flex items-center gap-2">
            <i class="pi pi-shield" />
            <span>{{ $t('project.wizard.tabs.govern') }}</span>
          </Tab>
          <Tab value="manage" class="flex items-center gap-2">
            <i class="pi pi-chart-line" />
            <span>{{ $t('project.wizard.tabs.manageMonitor') }}</span>
          </Tab>
        </TabList>
      </Tabs>

      <!-- Publish approval cluster — right-aligned in the tab row (LOCKED,
           ruled 2026-07-24 — see Wizard.intent.md). Members and the "View
           project"/"View dashboard" navigation used to live in this cluster
           too; they've moved to the shared ProjectTopbarTools in the topbar
           (see the Teleport above) and the RevisionSwitcher's Dashboard
           entry, respectively. -->
      <span class="ml-auto flex items-center justify-end gap-2 flex-wrap">
        <!-- One state-driven primary control reading the well-known
             'publish' governance step (see publishAction). No status Tag
             (ruled 2026-07-24): the button's state carries the status; its
             tooltip carries the hint + the reviewer's note. Only relevant
             pre-publish: once a project is live, the RevisionSwitcher's
             Dashboard entry takes over and a fresh cycle only resumes with a
             future revision. -->
        <template v-if="!isLive">
          <Button
            v-if="publishAction === 'request'"
            :label="
              publishStatus === 'changes-requested'
                ? $t('project.publish.actions.resubmit')
                : $t('project.publish.actions.requestApproval')
            "
            icon="pi pi-send"
            size="small"
            :severity="publishStatus === 'changes-requested' ? 'warn' : undefined"
            :loading="requestApprovalLoading"
            v-tooltip.bottom="publishStatusTooltip"
            @click="requestApproval"
          />
          <Button
            v-else-if="publishAction === 'approve'"
            :label="$t('project.publish.actions.approveProject')"
            icon="pi pi-check"
            severity="success"
            size="small"
            :disabled="anyStepFlagged"
            :loading="approveProjectLoading"
            v-tooltip.bottom="
              anyStepFlagged
                ? $t('project.publish.actions.approveBlockedTooltip')
                : publishStatusTooltip
            "
            @click="approveProject"
          />
          <Button
            v-else-if="publishAction === 'publish'"
            :label="$t('project.publish.actions.publish')"
            icon="pi pi-cloud-upload"
            size="small"
            :loading="publishing"
            v-tooltip.bottom="publishStatusTooltip"
            @click="confirmPublish"
          />
        </template>
      </span>
    </div>

    <div v-if="activeTab !== 'manage'" class="flex-1 flex gap-4 p-3 min-h-0">
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
            :style="{ width: activeTab === 'build' ? leftPct + '%' : '100%' }"
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
              <DataSensitivityStep
                v-else-if="isSensitivity"
                :project="project"
                :disabled="locked"
              />
              <FriaDeterminationStep
                v-else-if="isFriaDetermination"
                v-model="working"
                :disabled="locked"
              />
              <FriaScenariosStep
                v-else-if="isFriaScenarios"
                :project="project"
                :disabled="locked"
              />
            </template>
          </div>

          <!-- Drag handle — Build only, like the graph it resizes. -->
          <div
            v-if="activeTab === 'build'"
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

          <!-- Resource graph — Build only. The Build tab is canvas-centric by
               locked contract; Govern is form work and the graph added nothing
               there, so the step panel takes the full width instead. -->
          <div v-if="activeTab === 'build'" class="overflow-hidden p-4 min-h-0 flex-1">
            <ResourceGraph :project="project" />
          </div>
        </div>
      </div>
    </div>

    <!-- Manage & Monitor: its own left rail (grouped nav, never STEPS — see
         config/manageNav.js) beside a content pane. No step sidebar — this
         tab carries no wizard steps. The content pane is a per-section
         COMPONENT DISPATCH (key → component, see MANAGE_SECTION_COMPONENTS
         below), never inline markup here — each section (board, metrics,
         activity, the five categories) owns its own file under
         components/wizard/manage/, so they can be built out independently
         without every one of them editing this region of Wizard.vue. -->
    <div v-else class="flex-1 flex gap-4 p-3 min-h-0">
      <!-- Rail + revision progress. The bar lives at the FOOT OF THIS RAIL and
           on no other tab (ruled 2026-07-28): it measures work items assigned
           to the open revision, so it belongs beside the sections that show
           them, not in the tab row where it also sat over Build and Govern. -->
      <aside class="w-72 shrink-0 h-full flex flex-col gap-3 min-h-0">
        <ManageNav
          class="flex-1 min-h-0"
          :active-key="activeSection"
          @select="activeSection = $event"
        />
        <RevisionCompletenessBar
          v-if="project"
          :project-id="rootProjectId"
          :revision-id="project.projectID"
          class="shrink-0 rounded-xl border border-surface bg-surface p-3"
        />
      </aside>

      <div
        class="flex-1 min-w-0 min-h-0 overflow-hidden rounded-xl border border-surface flex flex-col"
      >
        <component
          :is="activeSectionComponent"
          v-if="activeSectionComponent"
          :project="project"
          :disabled="locked"
          class="flex-1 min-h-0"
          @category-selected="activeSection = $event"
        />
      </div>
    </div>

    <WizardToolbar
      v-if="showToolbar"
      :can-write="canWrite"
      :can-grant="canGrant"
      :show-save="showStepSave"
      :can-prev="canPrev"
      :can-next="canNext"
      :step-index="stepIndex"
      :step-count="navSteps.length"
      @save="onSave"
      @request-changes="openRequestChanges"
      @approve="onApprove"
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
       Build or Govern step, any status, any time) with a required note.
       Approving directly (WizardToolbar's Approve button) needs no dialog. -->
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
import ProjectTopbarTools from '@/sections/project/components/project/ProjectTopbarTools.vue'
import RevisionSwitcher from '@/sections/project/components/project/RevisionSwitcher.vue'
import ManageNav from '@/sections/project/components/wizard/ManageNav.vue'
import ManageActivity from '@/sections/project/components/wizard/manage/ManageActivity.vue'
import ManageBoard from '@/sections/project/components/wizard/manage/ManageBoard.vue'
import ManageFeature from '@/sections/project/components/wizard/manage/ManageFeature.vue'
import ManageIncident from '@/sections/project/components/wizard/manage/ManageIncident.vue'
import ManageOverview from '@/sections/project/components/wizard/manage/ManageOverview.vue'
import ManagePrivacy from '@/sections/project/components/wizard/manage/ManagePrivacy.vue'
import ManageReview from '@/sections/project/components/wizard/manage/ManageReview.vue'
import ManageTask from '@/sections/project/components/wizard/manage/ManageTask.vue'
import RevisionCompletenessBar from '@/sections/project/components/wizard/RevisionCompletenessBar.vue'
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
import FriaDeterminationStep from '@/sections/project/components/wizard/steps/FriaDeterminationStep.vue'
import FriaScenariosStep from '@/sections/project/components/wizard/steps/FriaScenariosStep.vue'
import ProjectSummaryStep from '@/sections/project/components/wizard/steps/ProjectSummaryStep.vue'
import ResourceManagementStep from '@/sections/project/components/wizard/steps/ResourceManagementStep.vue'
import { ACCESS_KINDS, kindConfig } from '@/sections/project/config/kinds'
import { friaDeterminationValues } from '@/sections/project/config/friaDeterminationForm'
import { MANAGE_NAV } from '@/sections/project/config/manageNav'
import {
  PUBLISH_GOVERNANCE_STEP_KEY,
  STEPS,
  kindsThroughStep,
  stepsForTab,
} from '@/sections/project/config/pipeline'
import { resourceManagementValues } from '@/sections/project/config/resourceManagementForm'
import { rolePreset } from '@/sections/project/config/roles'
import { summaryDefaults } from '@/sections/project/config/summaryForm'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { useProjectUsersStore } from '@/sections/project/stores/users'
import { useConfirm } from 'primevue/useconfirm'
import { computed, inject, provide, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const store = useProjectsStore()
const usersStore = useProjectUsersStore()
const $toast = inject('$toast')
const confirm = useConfirm()

const project = computed(() => store.findById(route.params.projectId))

// Chain-root id — the report endpoint's required ProjectID scope, regardless
// of the revisionId that narrows it (see RevisionCompletenessBar.vue / the
// same derivation ManageOverview.vue uses for OverviewPanel).
const rootProjectId = computed(() => project.value?.rootProjectID || project.value?.projectID)

// A live (published) project has a dashboard to switch to; drafts are wizard-only.
const isLive = computed(() => ['active', 'published'].includes(project.value?.status))

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
// Every project behaves identically now (no build modes): capability flags
// only gate individual review ACTIONS (Approve / Request changes on a step,
// and the topbar toolbar's project approval cluster below), never tab or
// step visibility.
const currentMember = computed(() =>
  store.membersFor(project.value?.projectID).find(m => m.userId === usersStore.currentUserID),
)
const currentRole = computed(() => rolePreset(currentMember.value?.role))
const canWrite = computed(() => !!currentRole.value.write)
const canGrant = computed(() => !!currentRole.value.grantApproval)
const canRequestApproval = computed(() => !!currentRole.value.requestApproval)

// --- Members dialog auto-open signal ----------------------------------------
// Members are no longer a wizard step — the shared ProjectTopbarTools (see
// the Teleport above) owns both the "Members" button and the MembersDialog
// mount now; it wraps exactly the table the old MembersStep held. Just-
// created-project flag: ProjectList.vue's onCreated() lands here with
// `?new=1` right after creating a project — members are the first thing to
// define on a fresh one, so this signals the shared component to auto-open
// the dialog once (its auto-open-members prop). Derived directly from the
// query param (not a manually-toggled ref) so the transition re-fires
// correctly if this same wizard instance is reused for a different
// just-created project. Stripped immediately via router.replace so a
// refresh or back-navigation never re-opens it.
const membersAutoOpen = computed(() => route.query.new === '1')
watch(
  membersAutoOpen,
  val => {
    if (!val) return
    router.replace({ query: { ...route.query, new: undefined } })
  },
  { immediate: true },
)

// --- Project-level publish approval cluster (wizard header) -----------------
// Drives the single state-machine button in the header tool row above, built
// around the well-known 'publish' governance step
// key (PUBLISH_GOVERNANCE_STEP_KEY): draft/changes-requested → submitted →
// approved, then a successful publish resets it to draft for the next cycle.
// Distinct from the direct per-step Approve/Request changes review on
// Build/Govern steps below (no submit stage there) — see the governance
// section of stores/projects.js for the shared rules (session-local
// scaffolding pending a governance redesign).
const publishStatus = computed(() =>
  project.value
    ? store.governanceStatus(project.value.projectID, PUBLISH_GOVERNANCE_STEP_KEY)
    : 'draft',
)
const publishReviewNote = computed(() =>
  project.value ? store.governanceNote(project.value.projectID, PUBLISH_GOVERNANCE_STEP_KEY) : '',
)

// Whether any OTHER (Build/Govern) step currently has changes requested —
// the "publish" step can't be approved while true (mirrored client-side by
// stores/projects.js transitionStep; this only keeps the button from firing a
// request that's certain to be rejected, and explains why via the disabled
// button's tooltip).
const anyStepFlagged = computed(
  () =>
    !!project.value &&
    STEPS.some(s => store.governanceStatus(project.value.projectID, s.key) === 'changes-requested'),
)

// Which action the header button represents right now, or null when the
// current member holds neither capability for it (no control renders then —
// such members follow status via the per-step review chips instead).
const publishAction = computed(() => {
  switch (publishStatus.value) {
    case 'approved':
      return 'publish'
    case 'submitted':
      return canGrant.value ? 'approve' : null
    default: // draft | changes-requested
      return canRequestApproval.value ? 'request' : null
  }
})

// Tooltip on the publish-action button — the only place the review note
// surfaces for a changes-requested publish (there is no per-step banner for
// it, unlike Build/Govern steps), plus a status+capability-aware hint
// mirroring what the old full-page PublishStep panel used to spell out
// inline.
const publishStatusTooltip = computed(() => {
  let hint
  switch (publishStatus.value) {
    case 'approved':
      hint = t('project.publish.hint.approved')
      break
    case 'submitted':
      hint = t(
        canGrant.value
          ? 'project.publish.hint.submitted.grantor'
          : 'project.publish.hint.submitted.viewer',
      )
      break
    case 'changes-requested':
      hint = t(
        canRequestApproval.value
          ? 'project.publish.hint.changesRequested.requester'
          : 'project.publish.hint.changesRequested.viewer',
      )
      break
    default:
      hint = t(
        canRequestApproval.value
          ? 'project.publish.hint.draft.requester'
          : 'project.publish.hint.draft.viewer',
      )
  }
  if (publishStatus.value === 'changes-requested' && publishReviewNote.value) {
    return `${hint} ${t('project.publish.note.label')} "${publishReviewNote.value}"`
  }
  return hint
})

function publishFail(err) {
  $toast.toastErrorHandler(t('project.publish.toast.actionFailed'))(err)
}

const requestApprovalLoading = ref(false)
async function requestApproval() {
  if (requestApprovalLoading.value) return
  requestApprovalLoading.value = true
  try {
    await store.transitionStep(project.value.projectID, PUBLISH_GOVERNANCE_STEP_KEY, 'submit')
    $toast.toastSuccess(t('project.publish.toast.submitted'))
  } catch (err) {
    publishFail(err)
  } finally {
    requestApprovalLoading.value = false
  }
}

const approveProjectLoading = ref(false)
async function approveProject() {
  if (approveProjectLoading.value || anyStepFlagged.value) return
  approveProjectLoading.value = true
  try {
    await store.transitionStep(project.value.projectID, PUBLISH_GOVERNANCE_STEP_KEY, 'approve')
    $toast.toastSuccess(t('project.publish.toast.approved'))
  } catch (err) {
    publishFail(err)
  } finally {
    approveProjectLoading.value = false
  }
}

// Publish itself needs a confirmation — it's a one-way door (the project goes
// live and further changes are made through revisions). Reuses the exact
// store action + success/dashboard handoff PublishStep.vue used to run.
function confirmPublish() {
  confirm.require({
    header: t('project.publish.confirm.header'),
    message: t('project.publish.confirm.message'),
    icon: 'pi pi-cloud-upload',
    rejectProps: {
      label: t('general.label.cancel'),
      severity: 'secondary',
      text: true,
      size: 'small',
    },
    acceptProps: { label: t('project.publish.confirm.accept'), size: 'small' },
    accept: doPublish,
  })
}
const publishing = ref(false)
async function doPublish() {
  if (publishing.value) return
  publishing.value = true
  try {
    await store.publishProject(project.value.projectID)
    $toast.toastSuccess(t('project.publish.toast.published'))
    // Dashboard handoff — a project only gets a dashboard the moment it goes
    // live, so a successful publish lands here automatically (see
    // Wizard.intent.md: "publish (confirmed, then dashboard handoff)").
    router.push({ name: 'project.overview', params: { projectId: project.value.projectID } })
  } catch (err) {
    $toast.toastErrorHandler(t('project.publish.toast.publishFailed'))(err)
  } finally {
    publishing.value = false
  }
}

// --- Tabs --------------------------------------------------------------
// Three fixed tabs, all visible to every member: Build, Govern, Manage &
// Monitor. The active tab is persisted in its own `tab` query param (Manage &
// Monitor has no steps, so it can't be recovered from route.query.step alone
// — a reload or deep link without it falls back to whichever tab owns the
// current step, or Build.
const VALID_TABS = ['build', 'govern', 'manage']
const tabForStepKey = key => STEPS.find(s => s.key === key)?.tab
// An explicit `tab`/`step` in the query always wins. Otherwise the landing tab
// depends on what the revision is FOR: a live revision is being run, so open
// Manage & Monitor; a draft is being built, so open Build. `project` resolves
// asynchronously, so the ref seeds with the query (or Build) and the watcher
// below corrects it once the project arrives — one shot, so it never fights a
// tab the user picked afterwards.
const tabFromQuery = (VALID_TABS.includes(route.query.tab) && route.query.tab) || null
const activeTab = ref(tabFromQuery || tabForStepKey(route.query.step) || 'build')
const landingTabResolved = ref(!!tabFromQuery || !!tabForStepKey(route.query.step))
// Keyed to the project's STATUS arriving, not to isLive: isLive is false both
// while the project is still loading and for a draft, so watching it would
// never resolve for drafts — and would then fire on a mid-session publish and
// yank the user to another tab.
watch(
  () => project.value?.status,
  status => {
    if (landingTabResolved.value || !status) return
    landingTabResolved.value = true
    if (isLive.value) activeTab.value = 'manage'
  },
  { immediate: true },
)
const navSteps = computed(() => {
  if (!project.value || activeTab.value === 'manage') return []
  return stepsForTab(activeTab.value)
})
// Switching tabs away from the step the URL currently points at snaps to
// that tab's first step, so the header/content/URL stay in sync. Manage &
// Monitor carries no steps, so its step is left alone (nothing to snap to) —
// the previous Build/Govern step stays in the query for when the member
// returns. `tab` itself is always written so the choice survives a reload.
watch(activeTab, tab => {
  if (tab === 'manage') {
    router.replace({ query: { ...route.query, tab } })
    return
  }
  const list = stepsForTab(tab)
  const query = { ...route.query, tab }
  if (!list.some(s => s.key === route.query.step)) query.step = list[0]?.key || undefined
  router.replace({ query })
})

// --- Manage & Monitor section --------------------------------------------
// The M&M tab's own left rail (ManageNav) switches its content pane via a
// `section` query param — the tab carries no wizard steps (config/manageNav.js),
// so this is the only navigation state it needs. Same ref + watch shape as
// activeTab above: ManageNav's @select sets activeSection directly (like
// Tabs' v-model:value does for activeTab), and the watcher keeps the query in
// sync so a reload or deep link resumes on the same section.
const VALID_SECTIONS = MANAGE_NAV.flatMap(section => section.items.map(item => item.key))
const activeSection = ref(
  (VALID_SECTIONS.includes(route.query.section) && route.query.section) || 'overview',
)
watch(activeSection, section => {
  router.replace({ query: { ...route.query, section } })
})
// Section key → panel component (see the template comment above). Every
// MANAGE_NAV item key needs an entry here; each component owns its own file
// under components/wizard/manage/ so the sections due to be built out next
// (metrics, activity, the five categories) don't collide editing this file.
const MANAGE_SECTION_COMPONENTS = {
  overview: ManageOverview,
  board: ManageBoard,
  activity: ManageActivity,
  incident: ManageIncident,
  feature: ManageFeature,
  privacy: ManagePrivacy,
  task: ManageTask,
  review: ManageReview,
}
const activeSectionComponent = computed(
  () => MANAGE_SECTION_COMPONENTS[activeSection.value] || null,
)

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
const isFriaDetermination = computed(() => activeKey.value === 'fria-determination')
const isFriaScenarios = computed(() => activeKey.value === 'fria-scenarios')
const isDataModel = computed(() => activeKey.value === 'data-model')
const isConnections = computed(() => activeKey.value === 'connections')
const isAutomations = computed(() => activeKey.value === 'automations')
const isAgents = computed(() => activeKey.value === 'agents')
const isChatbots = computed(() => activeKey.value === 'chatbots')
const isPages = computed(() => activeKey.value === 'pages')
const isRoles = computed(() => activeKey.value === 'roles')
const isPermissions = computed(() => activeKey.value === 'permissions')
const isUsers = computed(() => activeKey.value === 'users')
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
// Direct review, no submit stage: a Build/Govern step only ever carries
// 'draft' (the default), 'changes-requested' (raised anytime by a granter via
// the per-step "Request changes" action) or 'approved' (raised anytime by a
// granter via the per-step "Approve" action — see WizardToolbar). 'submitted'
// belongs to the Publish governance step alone, which lives in the topbar
// toolbar cluster now, not in this step dispatch.
const stepStatus = key =>
  project.value ? store.governanceStatus(project.value.projectID, key) : 'draft'
const status = computed(() => stepStatus(activeKey.value))
const reviewNote = computed(() =>
  project.value ? store.governanceNote(project.value.projectID, activeKey.value) : '',
)
// Editing is gated purely on the write capability — a step's governance
// status never locks it, so a granter can always re-review after a change
// (there is no reopen/unlock action to undo an approve otherwise).
const locked = computed(() => !canWrite.value)
const showStatus = computed(() => !!activeStep.value)
const showToolbar = computed(() => !!activeStep.value)
// Save only makes sense on form steps — sensitivity/resource steps persist
// each change immediately through their own store calls instead.
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
  'data-model': 'project.wizard.blurb.dataModel',
  'data-sensitivity': 'project.wizard.blurb.dataSensitivity',
  'fria-determination': 'fria.wizard.blurb.determination',
  'fria-scenarios': 'fria.wizard.blurb.scenarios',
  connections: 'project.wizard.blurb.connections',
  automations: 'project.wizard.blurb.automations',
  agents: 'project.wizard.blurb.agents',
  chatbots: 'project.wizard.blurb.chatbots',
  pages: 'project.wizard.blurb.pages',
  roles: 'project.wizard.blurb.roles',
  permissions: 'project.wizard.blurb.permissions',
  users: 'project.wizard.blurb.users',
}
const headerHint = computed(() => {
  const blurbKey = STEP_BLURB[activeKey.value]
  return blurbKey ? t(blurbKey) : ''
})

// --- Working copy of the active form step's values -------------------------
const working = ref({})
function loadWorking() {
  const saved = project.value
    ? store.governanceValues(project.value.projectID, activeKey.value)
    : {}
  if (isSummary.value) {
    const vals = { ...summaryDefaults(), ...saved }
    if (!vals.systemName) vals.systemName = project.value?.name || ''
    working.value = vals
  } else if (isResourceMgmt.value) {
    working.value = resourceManagementValues(saved)
  } else if (isFriaDetermination.value) {
    working.value = friaDeterminationValues(saved)
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

// --- Prev / Next stepper ---------------------------------------------------
// Walk the active tab's step list (navSteps). The arrows disable at the ends.
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
// Governance mutations are session-local (see stores/projects.js), but still
// routed through a try/catch — a store call can throw (e.g. the "publish"
// approval gate while a step is flagged), so report that instead of assuming
// success.
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

// Direct "Approve" — a member with grant-approval capability can approve the
// active step at any time, from any status (including re-approving after a
// changes-requested flag was addressed). See WizardToolbar's Approve button.
function onApprove() {
  governanceAction(
    () => store.transitionStep(project.value.projectID, activeKey.value, 'approve'),
    t('project.wizard.approve.toast.approved'),
  )
}

// --- Per-step "Request changes" ---------------------------------------------
// A member with grant-approval capability can flag the active Build/Govern
// step at any time, with a required note. This flags that step to
// changes-requested regardless of its current status, and sends the Publish
// step back for review too if it was pending (see stores/projects.js
// transitionStep).
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

<style scoped>
/* Segmented pill switcher — restyles PrimeVue's default underline tabs into
   a rounded track with a raised pill for the active tab. Visual only; the
   tab set and its always-visible behavior are the locked contract. */
.wizard-tabs :deep(.p-tablist) {
  background: transparent;
}

.wizard-tabs :deep(.p-tablist-tab-list) {
  display: inline-flex;
  gap: 0.25rem;
  padding: 0.3rem;
  border: 1px solid var(--p-content-border-color);
  border-radius: 9999px;
  background: var(--p-content-hover-background);
}

.wizard-tabs :deep(.p-tab) {
  border: 1px solid transparent;
  border-radius: 9999px;
  padding: 0.5rem 1.125rem;
  color: var(--p-text-muted-color);
  transition:
    background-color 150ms,
    color 150ms,
    box-shadow 150ms;
}

.wizard-tabs :deep(.p-tab:not(.p-tab-active):hover) {
  color: var(--p-text-color);
  background: var(--p-content-background);
}

.wizard-tabs :deep(.p-tab-active) {
  background: var(--p-content-background);
  border-color: var(--p-content-border-color);
  color: var(--p-primary-color);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
  font-weight: 600;
}

.wizard-tabs :deep(.p-tablist-active-bar) {
  display: none;
}
</style>
