<template>
  <div class="flex flex-col gap-3">
    <Message :severity="noteVisible ? 'warn' : 'secondary'" :closable="false">
      {{ hint }}
    </Message>

    <!-- The reviewer's note is the most important thing on the screen when a
         publish has been sent back, so it gets the room. -->
    <blockquote v-if="noteVisible" class="border-l-2 border-amber-500 pl-3 py-1 text-sm">
      <span class="block text-xs uppercase tracking-wide text-muted-color mb-0.5">
        {{ $t('project.publish.note.label') }}
      </span>
      {{ note }}
    </blockquote>

    <!-- The reviewer writes here; the two decision buttons live in the tab's
         action bar with everything else that acts. A note is optional on an
         approval and the whole point of a send-back, so one field serves both
         rather than a dialog that only appears for one of them. -->
    <Textarea
      v-if="canDecide"
      :model-value="decisionNote"
      rows="3"
      auto-resize
      class="w-full"
      :placeholder="$t('project.publish.note.placeholder')"
      @update:model-value="$emit('update:decisionNote', $event)"
    />

    <p v-if="blockedByFlaggedStep" class="text-sm text-amber-500">
      {{ $t('project.publish.actions.approveBlockedTooltip') }}
    </p>

    <!-- Two people minimum for anything going live. Said here as well as beside
         the button, because this is the stage a submitter comes to when they
         wonder why nothing is happening. -->
    <p v-if="isOwnRequest && status === 'submitted'" class="text-sm text-muted-color">
      {{ $t('project.publish.blocked.ownRequest') }}
    </p>
  </div>
</template>

<script setup>
// Publish stage 3 — the project-level review. The actions themselves live in
// the tab's action bar (one primary control, labelled for whatever the next
// step actually is), so this stage carries the state and the reviewer's words.
import Textarea from 'primevue/textarea'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  // draft | submitted | approved | changes-requested
  status: { type: String, default: 'draft' },
  note: { type: String, default: '' },
  canGrant: { type: Boolean, default: false },
  canRequestApproval: { type: Boolean, default: false },
  // Whether the caller is the person who submitted the open request. The
  // server refuses their decision, so the screen must not imply otherwise.
  isOwnRequest: { type: Boolean, default: false },
  blockedByFlaggedStep: { type: Boolean, default: false },
  decisionNote: { type: String, default: '' },
})

defineEmits(['update:decisionNote'])

// Whether this caller is the one being asked to decide right now.
const canDecide = computed(
  () => props.status === 'submitted' && props.canGrant && !props.isOwnRequest,
)

const { t } = useI18n()

const noteVisible = computed(() => props.status === 'changes-requested' && !!props.note)

// Same hint set the old tooltip resolved, minus the tooltip.
const hint = computed(() => {
  switch (props.status) {
    case 'approved':
      return t('project.publish.hint.approved')
    case 'submitted':
      return t(
        props.canGrant && !props.isOwnRequest
          ? 'project.publish.hint.submitted.grantor'
          : 'project.publish.hint.submitted.viewer',
      )
    case 'changes-requested':
      return t(
        props.canRequestApproval
          ? 'project.publish.hint.changesRequested.requester'
          : 'project.publish.hint.changesRequested.viewer',
      )
    default:
      return t(
        props.canRequestApproval
          ? 'project.publish.hint.draft.requester'
          : 'project.publish.hint.draft.viewer',
      )
  }
})
</script>
