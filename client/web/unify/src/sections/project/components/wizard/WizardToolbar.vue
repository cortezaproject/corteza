<template>
  <div class="shrink-0 border-t border-surface bg-surface p-3 flex items-center gap-2">
    <!-- Left: back to the project list (balances the right-side actions and
         keeps the stepper centered). -->
    <div class="flex-1 flex items-center">
      <Button
        icon="pi pi-arrow-left"
        :label="$t('project.wizard.toolbar.projects')"
        severity="secondary"
        text
        @click="$emit('back')"
      />
    </div>

    <!-- Center: Previous / Next stepper. Equal side columns keep the counter
         at the true center even when the button labels differ in width
         ("Previous" vs "Next" / "Request approval"). -->
    <div class="grid grid-cols-[1fr_auto_1fr] items-center gap-2">
      <Button
        icon="pi pi-chevron-left"
        :label="$t('project.wizard.toolbar.previous')"
        severity="secondary"
        text
        size="small"
        class="justify-self-end"
        :disabled="!canPrev"
        @click="$emit('prev')"
      />
      <span class="text-xs text-muted-color tabular-nums whitespace-nowrap text-center">
        {{ $t('project.wizard.stepCounter', { index: stepIndex, count: stepCount }) }}
      </span>
      <Button
        :label="$t('general.label.next')"
        icon="pi pi-chevron-right"
        icon-pos="right"
        severity="secondary"
        text
        size="small"
        class="justify-self-start"
        :disabled="!canNext"
        @click="$emit('next')"
      />
    </div>

    <!-- Right: per-step actions (hidden on milestone steps) -->
    <div class="flex-1 flex items-center justify-end gap-2">
      <template v-if="showActions">
        <!-- Free mode: just save, no approval -->
        <template v-if="mode !== 'gated'">
          <Button
            v-if="showSave"
            :label="$t('general.label.save')"
            icon="pi pi-save"
            @click="$emit('save')"
          />
        </template>

        <!-- Approved: anyone who works in governance can reopen -->
        <Button
          v-else-if="status === 'approved'"
          v-show="canRequest || canGrant"
          :label="$t('project.wizard.toolbar.reopen')"
          icon="pi pi-lock-open"
          severity="secondary"
          outlined
          @click="$emit('reopen')"
        />

        <!-- Submitted: approvers (grant) act on it -->
        <template v-else-if="status === 'submitted' && canGrant">
          <Button
            :label="$t('project.wizard.toolbar.requestChanges')"
            icon="pi pi-replay"
            severity="secondary"
            outlined
            @click="$emit('request-changes')"
          />
          <Button
            :label="$t('project.wizard.toolbar.approve')"
            icon="pi pi-check"
            severity="success"
            @click="$emit('approve')"
          />
        </template>

        <!-- Editable (draft / changes-requested): writers save, requesters submit -->
        <template v-else-if="status === 'draft' || status === 'changes-requested'">
          <Button
            v-if="canWrite && showSave"
            :label="$t('general.label.save')"
            icon="pi pi-save"
            :severity="showResubmit ? 'secondary' : undefined"
            :outlined="showResubmit"
            @click="$emit('save')"
          />
          <Button
            v-if="showResubmit"
            :label="$t('project.wizard.toolbar.resubmit')"
            icon="pi pi-send"
            :disabled="submitDisabled"
            @click="$emit('resubmit')"
          />
        </template>
      </template>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  status: { type: String, default: 'draft' },
  // Capability flags of the current member (from their role preset).
  canWrite: { type: Boolean, default: false },
  canRequest: { type: Boolean, default: false },
  canGrant: { type: Boolean, default: false },
  mode: { type: String, default: 'gated' },
  // A reopened step sits in `draft` but was previously approved; it needs its
  // own resubmit path.
  reopened: { type: Boolean, default: false },
  // Blocks resubmission until prerequisites are met (e.g. architecture unsaved).
  submitDisabled: { type: Boolean, default: false },
  // --- Stepper ---
  // Whether the per-step action buttons (Save/Approve/etc.) are shown; false on
  // milestone steps, which only get the stepper.
  showActions: { type: Boolean, default: true },
  // Whether the Save button is shown among the actions; false on steps that
  // persist each change immediately but still carry the approval flow.
  showSave: { type: Boolean, default: true },
  canPrev: { type: Boolean, default: false },
  canNext: { type: Boolean, default: false },
  stepIndex: { type: Number, default: 1 },
  stepCount: { type: Number, default: 1 },
})
defineEmits(['save', 'approve', 'request-changes', 'resubmit', 'reopen', 'prev', 'next', 'back'])

// Resubmit appears for a member who can request approval when a step is being
// re-worked (changes requested, or a reopened-then-draft step).
const showResubmit = computed(
  () =>
    props.canRequest &&
    (props.status === 'changes-requested' || (props.status === 'draft' && props.reopened)),
)
</script>
