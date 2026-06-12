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

        <!-- Step content + live resource panel -->
        <div ref="splitRef" class="flex-1 min-h-0 flex">
          <div
            class="min-h-0"
            :class="isDataModel ? 'overflow-hidden flex flex-col' : 'overflow-y-auto p-4'"
            :style="{ width: leftPct + '%' }"
          >
            <template v-if="isDataModel">
              <StepStatusBanner v-if="showStatus" :status="status" :review-note="reviewNote" class="m-3 mb-0 shrink-0" />
              <DataModelStep :project="project" :disabled="locked" class="flex-1 min-h-0" />
            </template>
            <template v-else>
              <StepStatusBanner v-if="showStatus" :status="status" :review-note="reviewNote" class="mb-4" />
              <ProjectSummaryStep v-if="isSummary" v-model="working" :disabled="locked" />
              <ResourceManagementStep v-else-if="isResourceMgmt" v-model="working" :disabled="locked" />
              <MembersStep v-else-if="isMembers" :project="project" :disabled="membersLocked" />
              <DataSensitivityStep v-else-if="isSensitivity" :project="project" :disabled="locked" />
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
            <ResourceGraph
              :project="project"
              :locked="status === 'submitted'"
              :filter-kind="stepKind"
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
      :show-actions="showStepActions"
      :show-save="showStepSave"
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

    <ModuleDialog
      v-model="configOpen"
      :project="project"
      :module-id="configId"
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
import ModuleDialog from '@/sections/project/components/datamodel/ModuleDialog.vue'
import ResourceGraph from '@/sections/project/components/graph/ResourceGraph.vue'
import StepNav from '@/sections/project/components/wizard/StepNav.vue'
import StepStatusBanner from '@/sections/project/components/wizard/StepStatusBanner.vue'
import WizardToolbar from '@/sections/project/components/wizard/WizardToolbar.vue'
import DataModelStep from '@/sections/project/components/wizard/steps/DataModelStep.vue'
import DataSensitivityStep from '@/sections/project/components/wizard/steps/DataSensitivityStep.vue'
import MembersStep from '@/sections/project/components/wizard/steps/MembersStep.vue'
import ProjectSummaryStep from '@/sections/project/components/wizard/steps/ProjectSummaryStep.vue'
import ResourceManagementStep from '@/sections/project/components/wizard/steps/ResourceManagementStep.vue'
import { STEPS, sections, stepsForTab } from '@/sections/project/config/pipeline'
import { resourceManagementValues } from '@/sections/project/config/resourceManagementForm'
import { rolePreset } from '@/sections/project/config/roles'
import { summaryDefaults } from '@/sections/project/config/summaryForm'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { useProjectUsersStore } from '@/sections/project/stores/users'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import { computed, provide, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const store = useProjectsStore()
const usersStore = useProjectUsersStore()
const confirm = useConfirm()
const toast = useToast()

const project = computed(() => store.findById(route.params.projectId))

// Load the full project (members for capability resolution) + user directory.
usersStore.load()
watch(
  () => route.params.projectId,
  id => {
    if (!id) return
    store.fetchProject(id).catch(err => {
      console.error('Failed to load project', err)
      toast.add({
        severity: 'error',
        summary: 'Could not load project',
        detail: err.message,
        life: 4000,
      })
    })
  },
  { immediate: true },
)

// --- Current member & role -------------------------------------------------
// Tab access follows the logged-in user's role: a member who can request OR
// grant approval works in the Governance view; everyone else sees Build only.
const currentMember = computed(() =>
  (project.value?.members || []).find(m => m.userId === usersStore.currentUserID),
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
const isResourceMgmt = computed(() => activeKey.value === 'resource-management')
const isMembers = computed(() => activeKey.value === 'members')
const isDataModel = computed(() => activeKey.value === 'data-model')
const isSensitivity = computed(() => activeKey.value === 'data-sensitivity')
// On a resource step the graph narrows to that step's kind; form steps show
// the whole-system overview.
const stepKind = computed(() =>
  activeStep.value?.type === 'resource' ? activeStep.value.kind : null,
)

// --- Module config dialog ----------------------------------------------------
// Provided to the step components: open an existing module, or stage a new one.
const configOpen = ref(false)
const configId = ref(null)
provide('configureResource', id => {
  configId.value = id
  configOpen.value = true
})
provide('createResource', () => {
  configId.value = null
  configOpen.value = true
})
watch(configOpen, open => {
  if (!open) configId.value = null
})

// --- Resizable split (step | resource graph) ---------------------------------
const splitRef = ref(null)
const leftPct = ref(60)
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
// Member management is RBAC-checked (project members.manage), not the role
// preset's write flag; the gate lock still applies once the section is sent.
const membersLocked = computed(
  () =>
    !project.value?.canManageMembers ||
    status.value === 'submitted' ||
    status.value === 'approved',
)
const showStatus = computed(() => project.value?.mode === 'gated' && !!activeStep.value)
const showToolbar = computed(() => !!activeStep.value)
// Step-level approval buttons apply to form, members and sensitivity steps;
// resource steps persist each change immediately and only lock via their gate.
// Save only makes sense on form steps (members/sensitivity persist immediately).
const showStepActions = computed(() =>
  ['form', 'members', 'sensitivity'].includes(activeStep.value?.type),
)
const showStepSave = computed(() => activeStep.value?.type === 'form')

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
const sectionList = computed(() => sections(stepsForTab('governance')))
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

// A step is submittable when it's editable (draft/changes-requested).
const stepSubmittable = key => ['draft', 'changes-requested'].includes(stepStatus(key))

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
  'resource-management':
    'Declare the permitted AI providers, infrastructure and third-party connections, then save. Later steps choose from these.',
  members: 'Assign people to project roles. Changes save as you go; the section locks at the gate.',
  'data-model':
    'Define the modules (data tables) and their fields. A Record field links one module to another. Changes save as you go.',
  'data-sensitivity':
    'Classify each field by its data sensitivity level. Changes save as you go; the section locks at the gate.',
}
const headerHint = computed(() => {
  const blurb = STEP_BLURB[activeKey.value] || ''
  // In gated mode a governance status message takes precedence over the blurb.
  return project.value?.mode === 'gated' ? statusHint(blurb) : blurb
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
      accept: async () => {
        try {
          await store.submitSection(project.value.id, sec.steps.map(s => s.key))
          toast.add({ severity: 'success', summary: 'Requested approval', life: 2000 })
        } catch (err) {
          toast.add({ severity: 'error', summary: 'Request failed', detail: err.message, life: 4000 })
        }
      },
    })
  } else {
    goStep(sec.steps[0].key)
  }
}

// --- Per-step actions ------------------------------------------------------
// Governance mutations persist via the API; report failures instead of
// assuming success.
async function governanceAction(fn, summary) {
  try {
    await fn()
    toast.add({ severity: 'success', summary, life: 2000 })
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Action failed', detail: err.message, life: 4000 })
  }
}
function onSave() {
  governanceAction(() => store.saveStepForm(project.value.id, activeKey.value, working.value), 'Saved')
}
function onApprove() {
  governanceAction(() => store.transitionStep(project.value.id, activeKey.value, 'approve'), 'Approved')
}
function onResubmit() {
  governanceAction(() => store.transitionStep(project.value.id, activeKey.value, 'submit'), 'Resubmitted')
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
async function confirmReason() {
  const { action, note } = reason.value
  if (!note.trim()) return
  try {
    await store.transitionStep(project.value.id, activeKey.value, action, note.trim())
    reason.value.visible = false
    toast.add({
      severity: 'info',
      summary: action === 'reopen' ? 'Reopened' : 'Sent back for changes',
      life: 2000,
    })
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Action failed', detail: err.message, life: 4000 })
  }
}
</script>
