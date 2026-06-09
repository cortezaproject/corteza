<template>
  <div class="shrink-0 border-t border-surface bg-surface p-3 flex items-center gap-2">
    <!-- Left: back to the project list (balances the right-side actions and
         keeps the stepper centered). -->
    <div class="flex-1 flex items-center">
      <Button
        icon="pi pi-arrow-left"
        label="Projects"
        severity="secondary"
        text
        size="small"
        @click="$emit('back')"
      />
    </div>

    <!-- Center: Previous / Next stepper -->
    <div class="flex items-center gap-2">
      <Button
        icon="pi pi-chevron-left"
        label="Previous"
        severity="secondary"
        text
        size="small"
        :disabled="!canPrev"
        @click="$emit('prev')"
      />
      <span class="text-xs text-muted-color tabular-nums whitespace-nowrap">{{ stepIndex }} / {{ stepCount }}</span>
      <Button
        :title="nextIsGate ? 'Submit this section for approval' : undefined"
        :label="nextIsGate ? 'Request approval' : 'Next'"
        icon="pi pi-chevron-right"
        icon-pos="right"
        severity="secondary"
        text
        size="small"
        :disabled="!nextIsGate && !canNext"
        @click="$emit('next')"
      />
    </div>

    <!-- Right: per-step actions (hidden on milestone steps) -->
    <div class="flex-1 flex items-center justify-end gap-2">
      <template v-if="showActions">
        <!-- Free mode: just save, no approval -->
        <template v-if="mode !== 'gated'">
          <Button label="Save" icon="pi pi-save" @click="$emit('save')" />
        </template>

        <!-- Approved: anyone who works in governance can reopen -->
        <Button
          v-else-if="status === 'approved'"
          v-show="canRequest || canGrant"
          label="Reopen"
          icon="pi pi-lock-open"
          severity="secondary"
          outlined
          @click="$emit('reopen')"
        />

        <!-- Submitted: approvers (grant) act on it -->
        <template v-else-if="status === 'submitted' && canGrant">
          <Button
            label="Request changes"
            icon="pi pi-replay"
            severity="secondary"
            outlined
            @click="$emit('request-changes')"
          />
          <Button label="Approve" icon="pi pi-check" severity="success" @click="$emit('approve')" />
        </template>

        <!-- Editable (draft / changes-requested): writers save, requesters submit -->
        <template v-else-if="status === 'draft' || status === 'changes-requested'">
          <Button
            v-if="canWrite"
            label="Save"
            icon="pi pi-save"
            :severity="showResubmit ? 'secondary' : undefined"
            :outlined="showResubmit"
            @click="$emit('save')"
          />
          <Button
            v-if="showResubmit"
            label="Resubmit"
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
  // A reopened step sits in `draft` but was previously approved; it needs a
  // resubmit path of its own (the gate only handles first-time submission).
  reopened: { type: Boolean, default: false },
  // Blocks resubmission until prerequisites are met (e.g. architecture unsaved).
  submitDisabled: { type: Boolean, default: false },
  // --- Stepper ---
  // Whether the per-step action buttons (Save/Approve/etc.) are shown; false on
  // milestone steps, which only get the stepper.
  showActions: { type: Boolean, default: true },
  canPrev: { type: Boolean, default: false },
  canNext: { type: Boolean, default: false },
  stepIndex: { type: Number, default: 1 },
  stepCount: { type: Number, default: 1 },
  // When true, "Next" submits the current gate section instead of advancing.
  nextIsGate: { type: Boolean, default: false },
})
defineEmits(['save', 'approve', 'request-changes', 'resubmit', 'reopen', 'prev', 'next', 'back'])

// Resubmit appears for a member who can request approval when a step is being
// re-worked (changes requested, or a reopened-then-draft step).
const showResubmit = computed(
  () => props.canRequest && (props.status === 'changes-requested' || (props.status === 'draft' && props.reopened)),
)
</script>
