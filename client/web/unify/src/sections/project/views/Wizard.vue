<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ project?.name || 'Project' }}</span>
  </Teleport>

  <div v-if="project" class="h-full flex flex-col min-h-0">
    <div class="flex-1 flex gap-4 p-4 min-h-0">
    <!-- Left: step nav -->
    <div class="w-72 shrink-0 flex flex-col gap-2 min-h-0">
      <StepNav
        class="flex-1 min-h-0"
        :steps="navSteps"
        :active-key="activeKey"
        :statuses="statuses"
        :gate-statuses="gateStatuses"
        :gate-locked="gateLocked"
        :show-gates="showGates"
        @select="goStep"
        @gate-click="onGateClick"
      />
    </div>

    <div class="flex-1 flex flex-col min-w-0 min-h-0 rounded-xl border border-surface bg-surface overflow-hidden">
      <!-- Step header -->
      <div class="shrink-0 border-b border-surface px-4 py-3 flex items-center gap-3">
        <div class="min-w-0">
          <h2 class="text-lg font-medium truncate">{{ activeStep?.title || activeStep?.label }}</h2>
          <p class="text-sm text-muted-color mt-0.5 min-h-[1.25rem]">{{ headerHint }}</p>
        </div>

        <div class="ml-auto flex items-center gap-3">
          <Tag v-if="showStatus" :value="statusLabel" :severity="statusSeverity" />
        </div>
      </div>

      <!-- Step content + dashboard -->
      <div ref="splitRef" class="flex-1 min-h-0 flex">
        <div
          class="min-h-0"
          :class="isArchitecture || isDataModel || isConnections || isDataSensitivity || isResourceManagement ? 'overflow-hidden flex flex-col' : 'overflow-y-auto p-4'"
          :style="{ width: leftPct + '%' }"
        >
          <template v-if="isArchitecture">
            <StepStatusBanner v-if="showStatus" :status="status" :review-note="reviewNote" class="m-3 mb-0 shrink-0" />
            <ArchitectureStep
              v-model="selectedGroupId"
              :project="project"
              :disabled="locked"
              class="flex-1 min-h-0"
            />
          </template>
          <template v-else-if="isDataModel">
            <StepStatusBanner v-if="showStatus" :status="status" :review-note="reviewNote" class="m-3 mb-0 shrink-0" />
            <DataModelStep :project="project" :disabled="locked" class="flex-1 min-h-0" />
          </template>
          <template v-else-if="isConnections">
            <StepStatusBanner v-if="showStatus" :status="status" :review-note="reviewNote" class="m-3 mb-0 shrink-0" />
            <ConnectionsStep :project="project" :disabled="locked" class="flex-1 min-h-0" />
          </template>
          <template v-else-if="isDataSensitivity">
            <StepStatusBanner v-if="showStatus" :status="status" :review-note="reviewNote" class="m-3 mb-0 shrink-0" />
            <DataSensitivityStep :project="project" :disabled="locked" class="flex-1 min-h-0" />
          </template>
          <template v-else-if="isResourceManagement">
            <StepStatusBanner v-if="showStatus" :status="status" :review-note="reviewNote" class="m-3 mb-0 shrink-0" />
            <ResourceManagementStep :project="project" :disabled="locked" class="flex-1 min-h-0" />
          </template>
          <template v-else>
            <StepStatusBanner v-if="showStatus" :status="status" :review-note="reviewNote" class="mb-4" />
            <ProjectSummaryStep v-if="isSummary" v-model="working" :disabled="locked" />
            <MembersStep v-else-if="isMembers" :project="project" :disabled="locked" />
            <div
              v-else
              class="h-full flex flex-col items-center justify-center gap-2 text-muted-color text-center"
            >
              <i :class="['text-3xl', placeholderIcon]" />
              <p class="text-sm">This step isn't built yet.</p>
            </div>
          </template>
        </div>

        <!-- Drag handle -->
        <div
          class="shrink-0 w-1.5 bg-surface-200 dark:bg-surface-700 hover:bg-primary relative cursor-col-resize group select-none flex items-center justify-center transition-colors"
          :class="{ '!bg-primary': resizing }"
          @pointerdown="startResize"
        >
          <span class="absolute inset-y-0 -left-1 -right-1" />
          <span class="h-8 w-0.5 rounded-full bg-surface-400 dark:bg-surface-500 group-hover:bg-primary-contrast" :class="{ '!bg-primary-contrast': resizing }" />
        </div>

        <div class="overflow-hidden p-4 min-h-0 flex-1">
          <ResourcePanel
            :project="project"
            :locked="status === 'submitted'"
            :selected-group-id="isArchitecture ? selectedGroupId : '__all__'"
            :emphasize-kind="emphasizeKind"
          />
        </div>
      </div>
    </div>
    </div>

    <WizardToolbar
      v-if="showToolbar"
      :status="status"
      :can-write="canWrite"
      :can-request="canRequest"
      :can-grant="canGrant"
      :mode="project.mode"
      :reopened="isReopened"
      :submit-disabled="isArchitecture && !architectureSaved"
      :show-actions="!isMilestone"
      :can-prev="canPrev"
      :can-next="canNext"
      :step-index="stepIndex"
      :step-count="navSteps.length"
      :next-is-gate="nextIsGate"
      @save="onSave"
      @approve="onApprove"
      @request-changes="openReason('request-changes')"
      @resubmit="onResubmit"
      @reopen="openReason('reopen')"
      @prev="onPrev"
      @next="onNext"
      @back="onBack"
    />

    <LinkConfigDialog
      v-model="configOpen"
      :project="project"
      :resource-id="configId"
      :new-kind="configNewKind"
      :new-connector="configNewConnector"
      :new-name="configNewName"
      :readonly="locked"
    />
  </div>

  <Dialog v-model:visible="reason.visible" modal :header="reasonText.header" :style="{ width: '32rem' }">
    <CFormGroup :label="reasonText.label" required>
      <Textarea v-model="reason.note" rows="3" auto-resize fluid :placeholder="reasonText.placeholder" />
    </CFormGroup>
    <template #footer>
      <div class="flex justify-end gap-2">
        <Button label="Cancel" severity="secondary" outlined size="small" @click="reason.visible = false" />
        <Button :label="reasonText.confirm" size="small" :disabled="!reason.note.trim()" @click="confirmReason" />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import ResourcePanel from '@/sections/project/components/wizard/ResourcePanel.vue'
import StepNav from '@/sections/project/components/wizard/StepNav.vue'
import StepStatusBanner from '@/sections/project/components/wizard/StepStatusBanner.vue'
import WizardToolbar from '@/sections/project/components/wizard/WizardToolbar.vue'
import ArchitectureStep from '@/sections/project/components/wizard/steps/ArchitectureStep.vue'
import ConnectionsStep from '@/sections/project/components/wizard/steps/ConnectionsStep.vue'
import DataModelStep from '@/sections/project/components/wizard/steps/DataModelStep.vue'
import DataSensitivityStep from '@/sections/project/components/wizard/steps/DataSensitivityStep.vue'
import LinkConfigDialog from '@/sections/project/components/group/LinkConfigDialog.vue'
import MembersStep from '@/sections/project/components/wizard/steps/MembersStep.vue'
import ProjectSummaryStep from '@/sections/project/components/wizard/steps/ProjectSummaryStep.vue'
import ResourceManagementStep from '@/sections/project/components/wizard/steps/ResourceManagementStep.vue'
import { STEPS, sections, stepsForTab } from '@/sections/project/config/pipeline'
import { rolePreset } from '@/sections/project/config/roles'
import { summaryDefaults } from '@/sections/project/config/summaryForm'
import { CURRENT_USER_ID } from '@/sections/project/mock/users'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import { computed, provide, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const store = useProjectsStore()
const confirm = useConfirm()
const toast = useToast()

const project = computed(() => store.findById(route.params.projectId))

// Group selected in the Architecture step; the resource graph follows it.
const selectedGroupId = ref('__all__')

// Link/permission dialog — provided to the whole wizard so both the Architecture
// step's badges and the resource graph nodes can open it.
const configOpen = ref(false)
const configId = ref(null)
const configNewKind = ref(null)
const configNewConnector = ref(null)
const configNewName = ref('')
provide('configureResource', id => {
  configId.value = id
  configNewKind.value = null
  configNewConnector.value = null
  configNewName.value = ''
  configOpen.value = true
})
// Open the dialog to create a brand-new resource (staged until Save). `opts`
// may carry a default name and (for connections) the catalog connector id.
provide('createResource', (kind, opts = {}) => {
  configId.value = null
  configNewKind.value = kind
  configNewConnector.value = opts.connector || null
  configNewName.value = opts.name || ''
  configOpen.value = true
})
// Reset the dialog target when it closes, so the next open starts clean.
watch(configOpen, open => {
  if (!open) {
    configId.value = null
    configNewKind.value = null
    configNewConnector.value = null
    configNewName.value = ''
  }
})

// Open a group from the graph: focus it on the Architecture step.
provide('openGroup', groupId => {
  selectedGroupId.value = groupId
  goStep('architecture')
})

// --- Resizable split (form | resource graph) -------------------------------
const splitRef = ref(null)
const leftPct = ref(60) // matches the former 3fr / 2fr ratio
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

// --- Current member & role -------------------------------------------------
// Tab access follows the logged-in user's role: a member who can request OR
// grant approval works in the Governance view; everyone else sees Build only.
const currentMember = computed(() => (project.value?.members || []).find(m => m.userId === CURRENT_USER_ID))
const currentRole = computed(() => rolePreset(currentMember.value?.role))
const canWrite = computed(() => !!currentRole.value.write)
const canRequest = computed(() => !!currentRole.value.requestApproval)
const canGrant = computed(() => !!currentRole.value.grantApproval)
const canGovern = computed(() => canRequest.value || canGrant.value)

// --- Tabs (derived, not chosen) --------------------------------------------
const effectiveTab = computed(() =>
  project.value?.mode === 'gated' && canGovern.value ? 'governance' : 'build',
)
const friaCtx = computed(() => ({ friaRequired: !!project.value?.friaRequired }))
const navSteps = computed(() =>
  project.value ? stepsForTab(effectiveTab.value, friaCtx.value) : [],
)
const showGates = computed(() => project.value?.mode === 'gated' && effectiveTab.value === 'governance')

// --- Active step -----------------------------------------------------------
const activeKey = computed(() => {
  const list = navSteps.value
  if (!list.length) return ''
  const q = route.query.step
  if (q && list.some(s => s.key === q)) return q
  return list[0].key
})
const activeStep = computed(
  () => navSteps.value.find(s => s.key === activeKey.value) || STEPS.find(s => s.key === activeKey.value),
)
const isSummary = computed(() => activeKey.value === 'summary')
const isMembers = computed(() => activeKey.value === 'members')
const isArchitecture = computed(() => activeKey.value === 'architecture')
const isDataModel = computed(() => activeKey.value === 'data-model')
const isConnections = computed(() => activeKey.value === 'connections')
const isDataSensitivity = computed(() => activeKey.value === 'data-sensitivity')
const isResourceManagement = computed(() => activeKey.value === 'resource-management')
// The current step's resource kind, emphasized in the live graph.
const emphasizeKind = computed(() =>
  activeStep.value?.type === 'resource' ? activeStep.value.kind : null,
)
const isMilestone = computed(() => activeStep.value?.type === 'milestone')
const placeholderIcon = computed(
  () => ({ preview: 'pi pi-eye', publish: 'pi pi-upload' })[activeKey.value] || 'pi pi-wrench',
)

// --- Per-step governance state --------------------------------------------
const stepStatus = key => project.value?.governance?.[key]?.status || 'draft'
const status = computed(() => stepStatus(activeKey.value))
const reviewNote = computed(() => project.value?.governance?.[activeKey.value]?.reviewNote || '')
// A previously-approved step that was reopened sits back in `draft` but carries
// a review note — it needs its own resubmit path (the gate only does first-time
// submission for the whole section).
const isReopened = computed(() => status.value === 'draft' && !!reviewNote.value)
const locked = computed(
  () => !canWrite.value || status.value === 'submitted' || status.value === 'approved',
)
const showStatus = computed(
  () => project.value?.mode === 'gated' && !!activeStep.value && !isMilestone.value,
)
// The toolbar always shows when a step is active so the Prev/Next stepper is
// reachable everywhere — milestone steps just hide the per-step action buttons.
const showToolbar = computed(() => !!activeStep.value)

// The Technical Architecture step edits groups live in the store, so "saving"
// snapshots the current groups. Only a saved (and unchanged-since) architecture
// may be submitted for approval at its gate.
const groupsSnapshot = computed(() => JSON.stringify(project.value?.groups || []))
const architectureSaved = computed(
  () => project.value?.governance?.architecture?.values?.snapshot === groupsSnapshot.value,
)

const statuses = computed(() => {
  const out = {}
  for (const s of navSteps.value) out[s.key] = stepStatus(s.key)
  return out
})

const statusLabel = computed(
  () =>
    ({ draft: 'Draft', submitted: 'Submitted', approved: 'Approved', 'changes-requested': 'Changes requested' })[
      status.value
    ],
)
const statusSeverity = computed(
  () =>
    ({ draft: 'secondary', submitted: 'info', approved: 'success', 'changes-requested': 'warn' })[status.value],
)

// --- Gate sections ---------------------------------------------------------
// Built from the full governance pipeline minus any conditional steps that
// don't apply to this project (e.g. FRIA when not required).
const sectionList = computed(() => sections(stepsForTab('governance', friaCtx.value)))
const sectionByGate = computed(() =>
  Object.fromEntries(sectionList.value.filter(s => s.gateKey).map(s => [s.gateKey, s])),
)
const orderedGateKeys = computed(() => sectionList.value.filter(s => s.gateKey).map(s => s.gateKey))
const gateStatuses = computed(() => {
  const out = {}
  for (const sec of sectionList.value) {
    if (!sec.gateKey) continue
    const sts = sec.steps.map(s => stepStatus(s.key))
    if (sts.some(s => s === 'changes-requested')) out[sec.gateKey] = 'changes-requested'
    else if (sts.every(s => s === 'approved')) out[sec.gateKey] = 'approved'
    else if (sts.some(s => s === 'submitted')) out[sec.gateKey] = 'submitted'
    else out[sec.gateKey] = 'draft'
  }
  return out
})

// A step is submittable when it's editable (draft/changes-requested); the
// architecture step additionally must be saved first.
function stepSubmittable(key) {
  if (!['draft', 'changes-requested'].includes(stepStatus(key))) return false
  if (key === 'architecture') return architectureSaved.value
  return true
}

// Gates unlock sequentially: a gate can only be requested once every earlier
// gate has been approved.
const gateLocked = computed(() => {
  const out = {}
  const keys = orderedGateKeys.value
  for (let i = 0; i < keys.length; i++) {
    out[keys[i]] = keys.slice(0, i).some(k => gateStatuses.value[k] !== 'approved')
  }
  return out
})

// --- Header hint -----------------------------------------------------------
function statusHint(draftText) {
  if (canGrant.value) {
    if (status.value === 'submitted') return 'Review this step, then approve or request changes.'
    if (status.value === 'approved') return 'Approved.'
    if (status.value === 'changes-requested') return 'Changes requested — waiting for the developer.'
    return 'Nothing to review yet.'
  }
  if (status.value === 'submitted') return 'Submitted — waiting for an approver.'
  if (status.value === 'approved') return 'Approved. Reopen to make changes.'
  if (status.value === 'changes-requested') return 'Changes requested — update and resubmit.'
  return draftText
}
// Short per-step blurb shown in the header (so steps don't repeat it in-body).
const STEP_BLURB = {
  summary: 'Fill in the form, then submit the section for approval at the gate.',
  members: 'Assign people to roles. The approval flags decide who works in Build vs Governance.',
  architecture: 'Group the roles and resources that work together.',
  'data-model': 'Define the modules (data tables) and their fields. A Record field links one module to another.',
  'resource-management': 'Whitelist the permitted AI providers, infrastructure and third-party connections, with continuity options.',
}
const headerHint = computed(() => {
  if (isMilestone.value) return activeStep.value?.description || ''
  const blurb = STEP_BLURB[activeKey.value] || ''
  // In gated mode a governance status message takes precedence over the blurb.
  return project.value?.mode === 'gated' ? statusHint(blurb) : blurb
})

// --- Working copy of the active form step's values -------------------------
const working = ref({})
function loadWorking() {
  if (!isSummary.value) {
    working.value = {}
    return
  }
  const vals = { ...summaryDefaults(), ...(project.value?.governance?.summary?.values || {}) }
  if (!vals.systemName) vals.systemName = project.value?.name || ''
  working.value = vals
}
watch([activeKey, project], loadWorking, { immediate: true })

function goStep(key) {
  router.replace({ query: { ...route.query, step: key } })
}

// Leave the wizard and return to the project list.
function onBack() {
  router.push({ name: 'project.list' })
}

// --- Prev / Next stepper ---------------------------------------------------
// Walk the role-visible step list (navSteps is already filtered by Build vs
// Governance). The arrows disable at the ends.
const stepPos = computed(() => navSteps.value.findIndex(s => s.key === activeKey.value))
const stepIndex = computed(() => stepPos.value + 1) // 1-based, for display
const canPrev = computed(() => stepPos.value > 0)
const canNext = computed(() => stepPos.value < navSteps.value.length - 1)

// "Next" turns into a gate submission when the current step is the last of its
// section (a gate step), we're a requester looking at gates, and the section is
// still submittable and unlocked. Otherwise it just advances.
const nextIsGate = computed(() => {
  const step = activeStep.value
  if (!step?.gate || !showGates.value || !canRequest.value) return false
  if (gateLocked.value[step.key]) return false
  return sectionByGate.value[step.key]?.steps.some(s => stepSubmittable(s.key)) || false
})

function onPrev() {
  if (canPrev.value) goStep(navSteps.value[stepPos.value - 1].key)
}
function onNext() {
  // At a ready gate boundary, submit the section (reuses the gate confirm flow)
  // instead of advancing — you cross into the next section only once approved.
  if (nextIsGate.value) {
    onGateClick(activeStep.value.key)
    return
  }
  if (canNext.value) goStep(navSteps.value[stepPos.value + 1].key)
}

// --- Gate submission -------------------------------------------------------
function onGateClick(gateKey) {
  const sec = sectionByGate.value[gateKey]
  if (!sec) return
  // Locked gates (previous gate not yet approved) can't be acted on.
  if (gateLocked.value[gateKey]) {
    toast.add({
      severity: 'warn',
      summary: 'Gate locked',
      detail: 'Approve the previous gate before requesting this one.',
      life: 2500,
    })
    return
  }
  const submittable = sec.steps.some(s => stepSubmittable(s.key))
  if (canRequest.value && submittable) {
    confirm.require({
      header: 'Request approval',
      message: 'This submits every step in this section for approval and locks them until reviewed. Continue?',
      icon: 'pi pi-lock',
      rejectProps: { label: 'Cancel', severity: 'secondary', text: true },
      acceptProps: { label: 'Request approval' },
      accept: () => {
        store.submitSection(project.value.id, sec.steps.map(s => s.key))
        toast.add({ severity: 'success', summary: 'Requested approval', life: 2000 })
      },
    })
  } else {
    goStep(sec.steps[0].key)
  }
}

// --- Per-step actions ------------------------------------------------------
function onSave() {
  const values = isArchitecture.value ? { snapshot: groupsSnapshot.value } : working.value
  store.saveStepForm(project.value.id, activeKey.value, values)
  toast.add({ severity: 'success', summary: 'Saved', life: 2000 })
}
function onApprove() {
  store.transitionStep(project.value.id, activeKey.value, 'approve')
  toast.add({ severity: 'success', summary: 'Approved', life: 2000 })
}
function onResubmit() {
  store.transitionStep(project.value.id, activeKey.value, 'submit')
  toast.add({ severity: 'success', summary: 'Resubmitted', life: 2000 })
}

// Reason dialog, shared by Request changes and Reopen.
const reason = ref({ visible: false, action: '', note: '' })
const reasonText = computed(() =>
  reason.value.action === 'reopen'
    ? {
        header: 'Reopen for changes',
        label: 'Why are you reopening this?',
        placeholder: 'Explain why this approved step is being reopened',
        confirm: 'Reopen',
      }
    : {
        header: 'Request changes',
        label: 'What needs to change?',
        placeholder: 'Explain what the developer should revise',
        confirm: 'Send back',
      },
)
function openReason(action) {
  reason.value = { visible: true, action, note: '' }
}
function confirmReason() {
  const { action, note } = reason.value
  if (!note.trim()) return
  store.transitionStep(project.value.id, activeKey.value, action, note.trim())
  reason.value.visible = false
  toast.add({
    severity: 'info',
    summary: action === 'reopen' ? 'Reopened' : 'Sent back for changes',
    life: 2000,
  })
}
</script>
