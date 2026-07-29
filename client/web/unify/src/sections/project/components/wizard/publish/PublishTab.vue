<template>
  <div class="h-full overflow-y-auto">
    <!-- Deliberately a single centred column with no left rail, unlike every
         other tab (ruled 2026-07-29): publishing is a rare one-way action, so
         it reads as a sequence to complete, not a place to browse. -->
    <div class="max-w-3xl mx-auto p-4 flex flex-col gap-4">
      <header class="flex items-start gap-3 flex-wrap">
        <div class="flex-1 min-w-60">
          <h2 class="text-xl font-medium">{{ heading }}</h2>
          <p class="text-sm text-muted-color mt-1">
            {{ isLive ? $t('project.publish.subheadingLive') : subheading }}
          </p>
        </div>
        <Tag :value="statusTag.label" :severity="statusTag.severity" />
      </header>

      <div v-if="loading" class="flex justify-center py-10">
        <ProgressSpinner style="width: 2rem; height: 2rem" />
      </div>

      <!-- Ahead of the error state on purpose: a live revision has no plan to
           fetch, so a failure there must never hide the receipt. -->
      <PublishReceipt
        v-else-if="isLive"
        :version="version"
        :creating-revision="creatingRevision"
        @view-dashboard="goDashboard"
        @new-revision="startNewRevision"
      />

      <Message v-else-if="loadError" severity="error" :closable="false">
        {{ $t('project.publish.loadFailed') }}
      </Message>

      <!-- A first publish has no parent to diff against and no records to
           migrate, so it gets its own screen rather than the revision flow
           with two of its four stages emptied out (ruled 2026-07-29). -->
      <PublishFirstRun
        v-else-if="!hasParent"
        :project-name="project.name"
        :inventory="inventory"
        :project-members="projectMembers"
      />

      <template v-else>
        <PublishStage
          :index="1"
          :title="$t('project.publish.stages.changes.title')"
          :hint="changesHint"
          :status-label="$t('project.publish.stages.reviewed')"
          status-severity="success"
          done
          :open="openStage === 1"
          @toggle="toggleStage(1)"
        >
          <PublishChanges :changes="plan.changes" :inventory="inventory" />
        </PublishStage>

        <PublishStage
          :index="2"
          :title="$t('project.publish.stages.migration.title')"
          :hint="migrationHint"
          :status-label="
            blockers.length
              ? $t('project.publish.stages.needsDecision')
              : $t('project.publish.stages.decided')
          "
          :status-severity="blockers.length ? 'warn' : 'success'"
          :current="!!blockers.length"
          :done="!blockers.length"
          :open="openStage === 2"
          @toggle="toggleStage(2)"
        >
          <PublishMigration
            :mappings="plan.suggestedMappings"
            :changes="plan.changes"
            :decisions="decisions"
            :disabled="!canWrite"
            @update:decisions="decisions = $event"
          />
        </PublishStage>

        <PublishStage
          :index="3"
          :title="$t('project.publish.stages.approval.title')"
          :hint="approvalHint"
          :status-label="approvalTag.label"
          :status-severity="approvalTag.severity"
          :current="!blockers.length && publishStatus !== 'approved'"
          :done="publishStatus === 'approved'"
          :locked="!!blockers.length"
          :open="openStage === 3"
          @toggle="toggleStage(3)"
        >
          <PublishApproval
            :status="publishStatus"
            :note="publishNote"
            :can-grant="canGrant"
            :can-request-approval="canRequestApproval"
            :blocked-by-flagged-step="anyStepFlagged && publishStatus === 'submitted'"
          />
        </PublishStage>

        <PublishStage
          :index="4"
          :title="$t('project.publish.stages.goLive.title')"
          :hint="$t('project.publish.stages.goLive.hint')"
          :current="publishStatus === 'approved'"
          :locked="publishStatus !== 'approved'"
          :open="openStage === 4"
          @toggle="toggleStage(4)"
        >
          <PublishGoLive
            :project-name="project.name"
            :handle="confirmationPhrase"
            :typed="typedConfirmation"
            :has-parent="hasParent"
            :losses="losses"
            :open-work-items="openWorkItems"
            :assigned-work-items="completeness.assigned"
            :disabled="!canWrite"
            @update:typed="typedConfirmation = $event"
          />
        </PublishStage>

        <!-- One primary control, labelled for the single next thing to do, with
             the reason it cannot fire stated beside it rather than hidden in a
             disabled tooltip. -->
        <div
          class="sticky bottom-0 -mx-4 px-4 py-3 mt-1 border-t border-surface bg-surface/90 backdrop-blur flex items-center gap-3 flex-wrap"
        >
          <p class="flex-1 min-w-48 text-sm text-muted-color">{{ actionReason }}</p>
          <Button
            v-if="primaryAction"
            :label="primaryLabel"
            :icon="primaryIcon"
            :severity="primaryAction === 'approve' ? 'success' : undefined"
            size="small"
            :disabled="!primaryEnabled"
            :loading="acting"
            @click="runPrimaryAction"
          />
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
// The Publish tab — the whole request → approve → publish cycle, which lived in
// the wizard's tab row until 2026-07-29 (see views/Wizard.intent.md). Self
// contained on purpose, the same way each manage/ section is: Wizard.vue mounts
// it with the open revision and the caller's capabilities, and everything else
// — the deployment plan, the migration decisions, the governance transitions
// and the publish call — is owned here.
import { fetchRevisionCompleteness } from '@/sections/project/composables/revisionCompleteness'
import { OVERVIEW_KINDS } from '@/sections/project/config/kinds'
import { PUBLISH_GOVERNANCE_STEP_KEY, STEPS } from '@/sections/project/config/pipeline'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { NoID } from '@planetcrust/human-js'
import PublishApproval from './PublishApproval.vue'
import PublishChanges from './PublishChanges.vue'
import PublishFirstRun from './PublishFirstRun.vue'
import PublishGoLive from './PublishGoLive.vue'
import PublishMigration from './PublishMigration.vue'
import PublishReceipt from './PublishReceipt.vue'
import PublishStage from './PublishStage.vue'
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const props = defineProps({
  project: { type: Object, required: true },
  canWrite: { type: Boolean, default: false },
  canGrant: { type: Boolean, default: false },
  canRequestApproval: { type: Boolean, default: false },
})

const { t } = useI18n()
const router = useRouter()
const store = useProjectsStore()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const plan = ref({ risk: 'safe', changes: [], suggestedMappings: [] })
const inventory = ref([])
const completeness = ref({ assigned: 0, completed: 0 })
const decisions = ref({})
const typedConfirmation = ref('')
const openStage = ref(0)
const loading = ref(true)
const loadError = ref(false)
const acting = ref(false)
const creatingRevision = ref(false)

const projectId = computed(() => props.project?.projectID)
const version = computed(() => (props.project?.revision || 0) + 1)
const isLive = computed(() => ['active', 'published'].includes(props.project?.status))
// Compared against NoID, never coerced to boolean: system.Project defaults an
// absent id to the STRING '0' (lib/js cast.ts), which is truthy — so `!!id`
// reports a parent on every first revision, the exact case this decides.
const hasParent = computed(
  () => !!props.project?.parentRevisionID && props.project.parentRevisionID !== NoID,
)
const rootProjectId = computed(() => props.project?.rootProjectID || projectId.value)
// Loaded by Wizard.vue's fetchProject; read here for the first-publish screen's
// "who gets access" panel.
const projectMembers = computed(() => store.membersFor(projectId.value))

// --- data ------------------------------------------------------------------
// Everything here reflects edits made seconds ago, so none of it is cached: the
// tab re-reads on open and after any publish-relevant change.
async function load() {
  if (!projectId.value) return
  // A live revision has nothing left to plan — the backend only computes a plan
  // for a draft and errors otherwise, so asking would manufacture a failure on
  // the one screen that should be showing a receipt.
  if (isLive.value) {
    loading.value = false
    return
  }
  loading.value = true
  loadError.value = false
  try {
    // A first publish has no parent to diff against, so the plan would come
    // back empty by definition — skip the round trip rather than ask a
    // question with a known answer.
    plan.value = hasParent.value
      ? await store.deploymentPlan(projectId.value)
      : { risk: 'safe', changes: [], suggestedMappings: [] }
    inventory.value = await loadInventory()
    // Work-item completeness only warns, so it must never fail the screen.
    completeness.value = await fetchRevisionCompleteness(
      $SystemAPI,
      rootProjectId.value,
      projectId.value,
    ).catch(() => ({ assigned: 0, completed: 0 }))
    openStage.value = blockers.value.length ? 2 : 3
  } catch (err) {
    console.error('Failed to load the deployment plan', err)
    loadError.value = true
  } finally {
    loading.value = false
  }
}

watch(
  () => [projectId.value, props.project?.status],
  () => {
    decisions.value = {}
    typedConfirmation.value = ''
    load()
  },
  { immediate: true },
)

// Per-kind counts of what this revision contains, from the same graph the Build
// canvas draws; `added` comes off the plan rather than a second graph read of
// the parent.
async function loadInventory() {
  const { nodes = [] } = await store.graph(projectId.value)
  const counts = {}
  for (const node of nodes) counts[node.kind] = (counts[node.kind] || 0) + 1

  const added = {}
  for (const change of plan.value.changes) {
    if (change.op === 'added' && change.kind !== 'field') {
      added[change.kind] = (added[change.kind] || 0) + 1
    }
  }

  // Ordered by OVERVIEW_KINDS — the pipeline order the Build tab's metric
  // cards use — but only what the revision actually has. Unlike that strip,
  // which is a fixed set of layer toggles, this is a list of what goes live,
  // and a kind with nothing in it is not going live.
  return OVERVIEW_KINDS.map(kind => ({
    kind,
    count: counts[kind] || 0,
    added: added[kind] || 0,
  })).filter(entry => entry.count > 0)
}

// --- migration decisions ---------------------------------------------------
// Destructive changes whose destination is gone. A retyped field is not one:
// casting is a defensible default, so it is reported, not asked about.
const blockers = computed(() =>
  plan.value.changes.filter(
    c => c.risk === 'dangerous' && c.op === 'removed' && !decisions.value[c.path],
  ),
)

const losses = computed(() =>
  plan.value.changes
    .filter(c => decisions.value[c.path] === 'drop')
    .map(c => ({
      path: c.path,
      kind: c.kind,
      name: c.name,
      module: c.module,
      records: c.records || 0,
    })),
)

// The mapping actually sent to the backend: the suggestion, plus a copy for
// every "move it into <field>" decision. Sending an empty set is what silently
// drops every record in the project, so this is never allowed to be a stub.
const resolvedMappings = computed(() =>
  plan.value.suggestedMappings.map(mapping => {
    const moved = plan.value.changes
      .filter(
        c =>
          c.kind === 'field' &&
          c.module === mapping.module &&
          decisions.value[c.path]?.startsWith('move:'),
      )
      .map(c => ({
        sourceField: c.name,
        targetField: decisions.value[c.path].slice('move:'.length),
        op: 'copy',
      }))

    if (!moved.length) return mapping

    // A moved field takes over its target, which the suggestion had starting
    // empty — drop that placeholder so the two do not both write the column.
    const targets = new Set(moved.map(m => m.targetField))
    return {
      ...mapping,
      fields: [...(mapping.fields || []).filter(f => !targets.has(f.targetField)), ...moved],
    }
  }),
)

// --- governance ------------------------------------------------------------
const publishStatus = computed(() =>
  projectId.value ? store.governanceStatus(projectId.value, PUBLISH_GOVERNANCE_STEP_KEY) : 'draft',
)
const publishNote = computed(() =>
  projectId.value ? store.governanceNote(projectId.value, PUBLISH_GOVERNANCE_STEP_KEY) : '',
)
const anyStepFlagged = computed(() =>
  STEPS.some(s => store.governanceStatus(projectId.value, s.key) === 'changes-requested'),
)

const openWorkItems = computed(() =>
  Math.max(0, completeness.value.assigned - completeness.value.completed),
)

// --- the single primary action ---------------------------------------------
const primaryAction = computed(() => {
  switch (publishStatus.value) {
    case 'approved':
      return 'publish'
    case 'submitted':
      return props.canGrant ? 'approve' : null
    default:
      return props.canRequestApproval ? 'request' : null
  }
})

const primaryLabel = computed(
  () =>
    ({
      request:
        publishStatus.value === 'changes-requested'
          ? t('project.publish.actions.resubmit')
          : t('project.publish.actions.requestApproval'),
      approve: t('project.publish.actions.approveProject'),
      publish: t('project.publish.actions.publish'),
    })[primaryAction.value] || '',
)

const primaryIcon = computed(
  () =>
    ({ request: 'pi pi-send', approve: 'pi pi-check', publish: 'pi pi-cloud-upload' })[
      primaryAction.value
    ],
)

// Publishing with data loss needs the project's own name typed back.
const confirmationPhrase = computed(() => props.project?.handle || props.project?.name || '')
const typedConfirmationOk = computed(
  () => !losses.value.length || typedConfirmation.value.trim() === confirmationPhrase.value,
)

const primaryEnabled = computed(() => {
  if (!props.canWrite) return false
  if (blockers.value.length) return false
  if (primaryAction.value === 'approve' && anyStepFlagged.value) return false
  if (primaryAction.value === 'publish' && !typedConfirmationOk.value) return false
  return true
})

// Why the button cannot fire, said out loud instead of hidden in a tooltip.
const actionReason = computed(() => {
  if (blockers.value.length) {
    return t('project.publish.blocked.undecided', {
      name: blockers.value[0].module
        ? `${blockers.value[0].module}.${blockers.value[0].name}`
        : blockers.value[0].name,
    })
  }
  if (primaryAction.value === 'approve' && anyStepFlagged.value) {
    return t('project.publish.actions.approveBlockedTooltip')
  }
  if (primaryAction.value === 'publish') {
    return typedConfirmationOk.value
      ? t('project.publish.blocked.irreversible')
      : t('project.publish.blocked.typeHandle')
  }
  if (!primaryAction.value) return t('project.publish.blocked.noCapability')
  if (!hasParent.value) return t('project.publish.blocked.readyFirst')
  return t('project.publish.blocked.ready')
})

async function runPrimaryAction() {
  if (acting.value || !primaryEnabled.value) return
  acting.value = true
  try {
    if (primaryAction.value === 'request') {
      await store.transitionStep(projectId.value, PUBLISH_GOVERNANCE_STEP_KEY, 'submit')
      $toast.toastSuccess(t('project.publish.toast.submitted'))
    } else if (primaryAction.value === 'approve') {
      await store.transitionStep(projectId.value, PUBLISH_GOVERNANCE_STEP_KEY, 'approve')
      $toast.toastSuccess(t('project.publish.toast.approved'))
      openStage.value = 4
    } else {
      await store.publishProject(projectId.value, resolvedMappings.value)
      $toast.toastSuccess(t('project.publish.toast.published'))
    }
  } catch (err) {
    $toast.toastErrorHandler(
      t(
        primaryAction.value === 'publish'
          ? 'project.publish.toast.publishFailed'
          : 'project.publish.toast.actionFailed',
      ),
    )(err)
  } finally {
    acting.value = false
  }
}

function goDashboard() {
  router.push({ name: 'project.overview', params: { projectId: projectId.value } })
}

async function startNewRevision() {
  if (creatingRevision.value) return
  creatingRevision.value = true
  try {
    const rev = await store.createRevision(projectId.value)
    router.push({ name: 'project.wizard', params: { projectId: rev.projectID } })
  } catch (err) {
    $toast.toastErrorHandler(t('project.publish.toast.actionFailed'))(err)
  } finally {
    creatingRevision.value = false
  }
}

// --- presentation ----------------------------------------------------------
function toggleStage(index) {
  openStage.value = openStage.value === index ? 0 : index
}

// A first publish is a launch, not a change-review — it says so from the title
// down rather than calling itself "revision 1".
const heading = computed(() =>
  hasParent.value
    ? t('project.publish.heading', { version: version.value })
    : t('project.publish.firstRun.title', { name: props.project?.name }),
)

const subheading = computed(() =>
  hasParent.value
    ? t('project.publish.subheading', { parent: props.project.revision })
    : t('project.publish.subheadingFirst'),
)

const changesHint = computed(() =>
  plan.value.changes.length
    ? t('project.publish.stages.changes.hint', {
        count: plan.value.changes.length,
        risky: plan.value.changes.filter(c => c.risk === 'dangerous').length,
      })
    : t('project.publish.stages.changes.hintNone'),
)

const migrationHint = computed(() => {
  if (!hasParent.value) return t('project.publish.stages.migration.hintFirst')
  return blockers.value.length
    ? t('project.publish.stages.migration.hintBlocked', { count: blockers.value.length })
    : t('project.publish.stages.migration.hintClear')
})

const approvalHint = computed(() =>
  publishStatus.value === 'approved'
    ? t('project.publish.stages.approval.hintApproved')
    : t('project.publish.stages.approval.hint'),
)

const approvalTag = computed(
  () =>
    ({
      approved: { label: t('project.governance.status.approved'), severity: 'success' },
      submitted: { label: t('project.publish.stages.awaiting'), severity: 'info' },
      'changes-requested': {
        label: t('project.governance.status.changesRequested'),
        severity: 'warn',
      },
    })[publishStatus.value] || {
      label: t('project.publish.stages.notStarted'),
      severity: 'secondary',
    },
)

const statusTag = computed(() => {
  if (isLive.value) return { label: t('project.status.active'), severity: 'success' }
  return approvalTag.value
})
</script>
