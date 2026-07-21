<template>
  <div class="flex flex-col items-center justify-center gap-4 py-10 text-center">
    <!-- Already live: the project has been published; offer the dashboard. -->
    <template v-if="isPublished">
      <i class="pi pi-check-circle text-4xl text-green-500" />
      <h1 class="text-xl font-medium">{{ $t('project.publishStep.publishedTitle') }}</h1>
      <p class="max-w-prose text-sm text-muted-color">
        {{ $t('project.publishStep.publishedDescription') }}
      </p>
      <Button
        :label="$t('project.publishStep.openDashboard')"
        icon="pi pi-arrow-right"
        icon-pos="right"
        size="small"
        @click="emit('open-dashboard')"
      />
    </template>

    <!-- Draft: review and publish. -->
    <template v-else>
      <i class="pi pi-cloud-upload text-4xl text-surface-400" />
      <h1 class="text-xl font-medium">{{ $t('project.publishStep.title') }}</h1>
      <p class="max-w-prose text-sm text-muted-color">
        {{ $t('project.publishStep.description') }}
      </p>

      <!-- Gated mode only: the publish-time approval workflow around the single
           'publish' governance step (free mode publishes directly, no gate). -->
      <div
        v-if="isGated"
        class="w-full max-w-md flex flex-col items-center gap-3 rounded-lg border border-surface bg-surface p-4"
      >
        <Tag :value="statusLabel" :severity="statusSeverity" />
        <p class="text-sm text-muted-color">{{ hint }}</p>
        <p
          v-if="govStep.status === 'changes-requested' && govStep.reviewNote"
          class="max-w-prose text-sm italic text-muted-color"
        >
          <span class="font-medium not-italic">
            {{ $t('project.publish.approval.note.label') }}
          </span>
          “{{ govStep.reviewNote }}”
        </p>

        <!-- Draft / changes-requested: the requester submits (or resubmits). -->
        <Button
          v-if="canSubmit"
          :label="$t('project.publish.approval.actions.submit')"
          icon="pi pi-send"
          size="small"
          :loading="submitting"
          @click="submit"
        />

        <!-- Submitted: an approver reviews it. -->
        <div v-else-if="canAct" class="flex items-center gap-2">
          <Button
            :label="$t('project.publish.approval.actions.requestChanges')"
            icon="pi pi-replay"
            severity="secondary"
            outlined
            size="small"
            :disabled="approving"
            @click="openRequestChanges"
          />
          <Button
            :label="$t('project.publish.approval.actions.approve')"
            icon="pi pi-check"
            severity="success"
            size="small"
            :loading="approving"
            @click="approve"
          />
        </div>
      </div>

      <!-- The publish action itself stays hidden — not just disabled — until
           the approval panel above reports 'approved' (gated mode); free mode
           has no gate, so canPublish is always true there and this always
           shows. -->
      <Button
        v-if="canPublish"
        :label="$t('project.publishStep.action')"
        icon="pi pi-cloud-upload"
        size="small"
        :loading="publishing"
        @click="emit('publish')"
      />
    </template>

    <!-- Request changes: the same note-required pattern as the wizard's
         per-step "Request changes" dialog (see Wizard.vue's requestChanges
         dialog). -->
    <Dialog
      v-model:visible="requestChanges.visible"
      modal
      :header="$t('project.publish.approval.requestChangesDialog.header')"
      :style="{ width: '32rem' }"
    >
      <CFormGroup :label="$t('project.publish.approval.requestChangesDialog.label')" required>
        <Textarea
          v-model="requestChanges.note"
          rows="3"
          auto-resize
          fluid
          :placeholder="$t('project.publish.approval.requestChangesDialog.placeholder')"
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
            :label="$t('project.publish.approval.requestChangesDialog.confirm')"
            size="small"
            :loading="requestChanges.sending"
            :disabled="!requestChanges.note.trim()"
            @click="confirmRequestChanges"
          />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { PUBLISH_GOVERNANCE_STEP_KEY } from '@/sections/project/config/pipeline'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { computed, inject, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  project: { type: Object, required: true },
  publishing: { type: Boolean, default: false },
  // Current member's capabilities (from their role preset — see Wizard.vue's
  // currentRole), passed down rather than re-derived here.
  canRequest: { type: Boolean, default: false },
  canGrant: { type: Boolean, default: false },
})

const emit = defineEmits(['publish', 'open-dashboard'])

const { t } = useI18n()
const store = useProjectsStore()
const $toast = inject('$toast')

// 'active' is the backend's live state; 'published' is treated the same for
// forward-compatibility. Either means the project is already live.
const isPublished = computed(() => ['active', 'published'].includes(props.project?.status))
const isGated = computed(() => props.project?.mode === 'gated')

// The single governance step this whole panel revolves around. Defaults to a
// fresh draft entry — matches the backend's Governance.Step() lazy-init, so a
// project that has never been submitted reads exactly like one that reset
// after its last publish.
const govStep = computed(
  () =>
    props.project?.governance?.[PUBLISH_GOVERNANCE_STEP_KEY] || { status: 'draft', reviewNote: '' },
)

const STATUS_LABEL_KEY = {
  draft: 'project.publish.approval.status.draft',
  submitted: 'project.publish.approval.status.submitted',
  approved: 'project.publish.approval.status.approved',
  'changes-requested': 'project.publish.approval.status.changesRequested',
}
const STATUS_SEVERITY = {
  draft: 'secondary',
  submitted: 'info',
  approved: 'success',
  'changes-requested': 'warn',
}
const statusLabel = computed(() =>
  t(STATUS_LABEL_KEY[govStep.value.status] || STATUS_LABEL_KEY.draft),
)
const statusSeverity = computed(() => STATUS_SEVERITY[govStep.value.status] || 'secondary')

// Context + next-action copy, tailored to both the step's status and whether
// this member can actually act on it (a member with neither capability still
// sees the status, just no buttons).
const hint = computed(() => {
  switch (govStep.value.status) {
    case 'approved':
      return t('project.publish.approval.hint.approved')
    case 'submitted':
      return t(
        props.canGrant
          ? 'project.publish.approval.hint.submitted.grantor'
          : 'project.publish.approval.hint.submitted.viewer',
      )
    case 'changes-requested':
      return t(
        props.canRequest
          ? 'project.publish.approval.hint.changesRequested.requester'
          : 'project.publish.approval.hint.changesRequested.viewer',
      )
    default:
      return t(
        props.canRequest
          ? 'project.publish.approval.hint.draft.requester'
          : 'project.publish.approval.hint.draft.viewer',
      )
  }
})

const canSubmit = computed(
  () => props.canRequest && ['draft', 'changes-requested'].includes(govStep.value.status),
)
const canAct = computed(() => props.canGrant && govStep.value.status === 'submitted')

// Publish itself is only unlocked once the step is approved — free mode has no
// such gate. The backend enforces the same rule; this just keeps the button
// from firing a request that's certain to be rejected.
const canPublish = computed(() => !isGated.value || govStep.value.status === 'approved')

function fail(err) {
  $toast.toastErrorHandler(t('project.publish.approval.toast.actionFailed'))(err)
}

const submitting = ref(false)
async function submit() {
  if (submitting.value) return
  submitting.value = true
  try {
    await store.transitionStep(props.project.projectID, PUBLISH_GOVERNANCE_STEP_KEY, 'submit')
    $toast.toastSuccess(t('project.publish.approval.toast.submitted'))
  } catch (err) {
    fail(err)
  } finally {
    submitting.value = false
  }
}

const approving = ref(false)
async function approve() {
  if (approving.value) return
  approving.value = true
  try {
    await store.transitionStep(props.project.projectID, PUBLISH_GOVERNANCE_STEP_KEY, 'approve')
    $toast.toastSuccess(t('project.publish.approval.toast.approved'))
  } catch (err) {
    fail(err)
  } finally {
    approving.value = false
  }
}

const requestChanges = reactive({ visible: false, note: '', sending: false })
function openRequestChanges() {
  requestChanges.note = ''
  requestChanges.visible = true
}
async function confirmRequestChanges() {
  const note = requestChanges.note.trim()
  if (!note || requestChanges.sending) return
  requestChanges.sending = true
  try {
    await store.transitionStep(
      props.project.projectID,
      PUBLISH_GOVERNANCE_STEP_KEY,
      'request-changes',
      note,
    )
    requestChanges.visible = false
    $toast.toastInfo(t('project.publish.approval.toast.changesRequested'))
  } catch (err) {
    fail(err)
  } finally {
    requestChanges.sending = false
  }
}
</script>
