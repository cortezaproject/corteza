<template>
  <div class="h-full flex flex-col min-h-0">
    <div class="shrink-0 p-4 pb-3 border-b border-surface">
      <h2 class="text-lg font-medium mb-1">{{ $t('project.manage.views.overview') }}</h2>
      <p class="text-sm text-muted-color">{{ $t('project.manage.overview.blurb') }}</p>
    </div>

    <OverviewPanel
      :project-id="rootProjectId"
      :revision-id="revisionId"
      class="flex-1 min-h-0"
      @category-selected="emit('category-selected', $event)"
    />
  </div>
</template>

<script setup>
// Manage & Monitor > Monitor > Overview — the dashboard's overview screen
// (components/dashboard/OverviewPanel.vue), scoped to the ONE revision open in
// the wizard (route.params.projectId), same revision-scoping contract as
// ManageBoard.vue/the Manage<Category>.vue files. Absorbs what used to be the
// separate ManageMetrics.vue panel (ruled 2026-07-28 — see
// views/views.intent.md / DashboardLayout.intent.md: the dashboard and the
// wizard's Manage & Monitor tab are ONE surface differing only in scope).
//
// UNLIKE ManageBoard.vue/the Manage<Category>.vue files, this section loads
// NO stores of its own — OverviewPanel is entirely report-endpoint-driven,
// chain-wide or revision-scoped alike (the report endpoint now accepts an
// optional revisionID that scopes aggregation server-side — see
// stores/report.js), so there is nothing here to watch/load. This file's only
// job is deriving (rootProjectId, revisionId) from the `project` prop, same
// as every sibling section.
import OverviewPanel from '@/sections/project/components/dashboard/OverviewPanel.vue'
import { computed } from 'vue'

const props = defineProps({
  // The OPEN revision (a lib system.Project instance) — route.params.projectId
  // IS the revision row; its rootProjectID is the chain root work items are
  // filed against (see stores/projects.js's Project class).
  project: { type: Object, required: true },
  // Passed by Wizard.vue to every M&M section, but accepted only for the
  // shared dispatch prop set (see ManagePlaceholder.vue) — this section is
  // read-only (same as the old ManageMetrics.vue) so there is nothing to
  // disable.
  disabled: { type: Boolean, default: false },
})

// Re-emitted up to Wizard.vue's <component :is> dispatch, which sets
// activeSection — a card click switches the M&M section the same way
// ManageNav's own item clicks do (see Wizard.vue).
const emit = defineEmits(['category-selected'])

const rootProjectId = computed(() => props.project?.rootProjectID || props.project?.projectID)
const revisionId = computed(() => props.project?.projectID)
</script>
