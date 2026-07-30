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

    <!-- What actually happened, for the record. The backend stamps updatedAt /
         updatedBy on the publish itself, so this is the publish's own evidence
         rather than something reconstructed. How many records moved is NOT
         shown: the publish endpoint does not report a count, and inventing one
         on a screen whose whole job is to be trusted about data would be
         worse than leaving it out. -->
    <dl
      v-if="publishedAt || publishedBy"
      class="grid grid-cols-1 sm:grid-cols-2 gap-px bg-surface border border-surface rounded-xl overflow-hidden"
    >
      <div v-if="publishedAt" class="bg-surface p-3">
        <dt class="text-xs font-medium uppercase tracking-wide text-muted-color">
          {{ $t('project.publish.receipt.publishedAt') }}
        </dt>
        <dd class="text-sm mt-1 ml-0">{{ publishedAt }}</dd>
      </div>
      <div v-if="publishedBy" class="bg-surface p-3">
        <dt class="text-xs font-medium uppercase tracking-wide text-muted-color">
          {{ $t('project.publish.receipt.publishedBy') }}
        </dt>
        <dd class="text-sm mt-1 ml-0">{{ publishedBy }}</dd>
      </div>
    </dl>

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
      <!-- Branching is its own permission (project.revise), so a user who was
           allowed to publish is not automatically allowed to start the next
           revision. Hidden rather than disabled: there is nothing here for
           them to act on. -->
      <Button
        v-if="canRevise"
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
  // Already formatted / resolved by the tab — this component only lays them
  // out, and either may be blank if the backend stamped neither.
  publishedAt: { type: String, default: '' },
  publishedBy: { type: String, default: '' },
  creatingRevision: { type: Boolean, default: false },
  canRevise: { type: Boolean, default: false },
})

defineEmits(['view-dashboard', 'new-revision'])
</script>
