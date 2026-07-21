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

    <!-- Right: Save (form steps only — see showSave) plus the per-step
         "Request changes" flag (gated mode, every step but Publish, any
         status, any time — see showRequestChanges). Per-section approval
         (submit/approve/reopen) is gone: the only step with its own
         submit/approve cycle is Publish, which renders its panel inline
         instead of using this toolbar (see Wizard.vue's showStepSave, which
         gates this toolbar's Save button to form-type steps only). -->
    <div class="flex-1 flex items-center justify-end gap-2">
      <!-- Free mode: just save, no approval concepts at all. -->
      <template v-if="mode !== 'gated'">
        <Button
          v-if="showSave"
          :label="$t('general.label.save')"
          icon="pi pi-save"
          @click="$emit('save')"
        />
      </template>
      <!-- Gated mode: writers save; a step never locks on its own status any
           more (only Publish does), so Save just needs write access. -->
      <template v-else>
        <Button
          v-if="canWrite && showSave"
          :label="$t('general.label.save')"
          icon="pi pi-save"
          @click="$emit('save')"
        />
      </template>

      <Button
        v-if="showRequestChanges"
        :label="$t('project.wizard.requestChanges.action')"
        icon="pi pi-replay"
        severity="secondary"
        outlined
        @click="$emit('request-changes')"
      />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  // Capability flags of the current member (from their role preset).
  canWrite: { type: Boolean, default: false },
  canGrant: { type: Boolean, default: false },
  mode: { type: String, default: 'gated' },
  // The Publish step carries its own approval panel (see PublishStep.vue)
  // rather than the generic per-step "Request changes" flag — it's excluded
  // here so a granter can't flag the very step that panel already governs.
  isPublish: { type: Boolean, default: false },
  // Whether the Save button is shown; false on steps that persist each
  // change immediately (members, sensitivity, resource steps) rather than
  // through a save button.
  showSave: { type: Boolean, default: true },
  // --- Stepper ---
  canPrev: { type: Boolean, default: false },
  canNext: { type: Boolean, default: false },
  stepIndex: { type: Number, default: 1 },
  stepCount: { type: Number, default: 1 },
})
defineEmits(['save', 'request-changes', 'prev', 'next', 'back'])

// A granter may flag any non-publish step at any time, gated mode only —
// unlike the old per-step approval flow, this doesn't depend on the step's
// current status (see server/system/service/project_governance.go's
// "flag anytime" behaviour).
const showRequestChanges = computed(
  () => props.mode === 'gated' && props.canGrant && !props.isPublish,
)
</script>
