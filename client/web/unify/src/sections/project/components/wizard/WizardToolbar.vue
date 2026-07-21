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

    <!-- Right: Save (form steps only — see showSave) plus the per-step direct
         review actions — Approve and Request changes — for members with
         grant-approval capability, on every step, at any status, any time.
         There is no per-step submit stage any more; only the Publish
         governance step (surfaced in the topbar, not this toolbar) keeps a
         submit → approve/request-changes cycle. -->
    <div class="flex-1 flex items-center justify-end gap-2">
      <Button
        v-if="canWrite && showSave"
        :label="$t('general.label.save')"
        icon="pi pi-save"
        @click="$emit('save')"
      />

      <Button
        v-if="canGrant"
        :label="$t('project.wizard.requestChanges.action')"
        icon="pi pi-replay"
        severity="secondary"
        outlined
        @click="$emit('request-changes')"
      />
      <Button
        v-if="canGrant"
        :label="$t('project.wizard.approve.action')"
        icon="pi pi-check"
        severity="success"
        @click="$emit('approve')"
      />
    </div>
  </div>
</template>

<script setup>
defineProps({
  // Capability flags of the current member (from their role preset).
  canWrite: { type: Boolean, default: false },
  canGrant: { type: Boolean, default: false },
  // Whether the Save button is shown; false on steps that persist each
  // change immediately (sensitivity, resource steps) rather than through a
  // save button.
  showSave: { type: Boolean, default: true },
  // --- Stepper ---
  canPrev: { type: Boolean, default: false },
  canNext: { type: Boolean, default: false },
  stepIndex: { type: Number, default: 1 },
  stepCount: { type: Number, default: 1 },
})
defineEmits(['save', 'request-changes', 'approve', 'prev', 'next', 'back'])
</script>
