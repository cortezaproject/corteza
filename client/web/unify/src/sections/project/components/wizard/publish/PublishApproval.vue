<template>
  <div class="flex flex-col gap-3">
    <Message :severity="noteVisible ? 'warn' : 'secondary'" :closable="false">
      {{ hint }}
    </Message>

    <!-- The reviewer's note used to live only in a button tooltip, where it was
         effectively unreadable. It is the single most important thing on the
         screen when a publish has been sent back, so it gets the room. -->
    <blockquote v-if="noteVisible" class="border-l-2 border-amber-500 pl-3 py-1 text-sm">
      <span class="block text-xs uppercase tracking-wide text-muted-color mb-0.5">
        {{ $t('project.publish.note.label') }}
      </span>
      {{ note }}
    </blockquote>

    <p v-if="blockedByFlaggedStep" class="text-sm text-amber-500">
      {{ $t('project.publish.actions.approveBlockedTooltip') }}
    </p>
  </div>
</template>

<script setup>
// Publish stage 3 — the project-level review. The actions themselves live in
// the tab's action bar (one primary control, labelled for whatever the next
// step actually is), so this stage carries the state and the reviewer's words.
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  // draft | submitted | approved | changes-requested
  status: { type: String, default: 'draft' },
  note: { type: String, default: '' },
  canGrant: { type: Boolean, default: false },
  canRequestApproval: { type: Boolean, default: false },
  blockedByFlaggedStep: { type: Boolean, default: false },
})

const { t } = useI18n()

const noteVisible = computed(() => props.status === 'changes-requested' && !!props.note)

// Same hint set the old tooltip resolved, minus the tooltip.
const hint = computed(() => {
  switch (props.status) {
    case 'approved':
      return t('project.publish.hint.approved')
    case 'submitted':
      return t(
        props.canGrant
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
