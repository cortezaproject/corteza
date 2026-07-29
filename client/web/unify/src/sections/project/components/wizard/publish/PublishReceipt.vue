<template>
  <div class="flex flex-col gap-4">
    <div class="rounded-xl border border-green-500/40 bg-green-500/5 p-5 flex gap-4 items-start">
      <span
        class="inline-flex items-center justify-center w-9 h-9 rounded-lg shrink-0 bg-green-500 text-white"
      >
        <i class="pi pi-check" />
      </span>
      <div class="min-w-0">
        <h3 class="text-lg font-medium">
          {{ $t('project.publish.receipt.heading', { version }) }}
        </h3>
        <p class="text-sm text-muted-color mt-1">
          {{ $t('project.publish.receipt.blurb') }}
        </p>
      </div>
    </div>

    <!-- A publish is the one moment a project's data is rewritten. Redirecting
         away from it throws out the only confirmation that it did what it
         promised, so the receipt stays put and the ways onward are explicit. -->
    <div class="flex flex-wrap gap-2">
      <Button
        :label="$t('project.publish.receipt.viewDashboard')"
        icon="pi pi-chart-line"
        size="small"
        @click="$emit('view-dashboard')"
      />
      <Button
        :label="$t('project.publish.receipt.newRevision')"
        icon="pi pi-plus"
        severity="secondary"
        outlined
        size="small"
        :loading="creatingRevision"
        @click="$emit('new-revision')"
      />
    </div>
  </div>
</template>

<script setup>
// Shown in place of the stages once the revision is live, so the tab still has
// something to say about a published revision instead of going blank.
defineProps({
  version: { type: Number, default: 1 },
  creatingRevision: { type: Boolean, default: false },
})

defineEmits(['view-dashboard', 'new-revision'])
</script>
