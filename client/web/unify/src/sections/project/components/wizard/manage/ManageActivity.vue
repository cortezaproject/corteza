<template>
  <div class="h-full flex flex-col min-h-0">
    <div class="shrink-0 p-4 pb-3 border-b border-surface">
      <h2 class="text-lg font-medium mb-1">{{ $t('project.manage.views.activity') }}</h2>
      <p class="text-sm text-muted-color">{{ $t('project.manage.activity.blurb') }}</p>
    </div>

    <ActivityPanel :project-id="rootProjectId" :revision-id="revisionId" class="flex-1 min-h-0" />
  </div>
</template>

<script setup>
// Manage & Monitor > Monitor > Activity — the dashboard's activity screen
// (components/dashboard/ActivityPanel.vue), scoped to the ONE revision open in
// the wizard (route.params.projectId), same revision-scoping contract as
// ManageOverview.vue/ManageBoard.vue/the Manage<Category>.vue files. Absorbs
// what used to be this file's placeholder body (ruled 2026-07-28 — see
// views/views.intent.md / DashboardLayout.intent.md: the dashboard and the
// wizard's Manage & Monitor tab are ONE surface differing only in scope, "All
// Events" renamed to "Activity" on both).
//
// UNLIKE ManageBoard.vue/the Manage<Category>.vue files, this section loads
// NO stores of its own — ActivityPanel talks to the actionlog endpoints
// directly (via useEventActivity + $SystemAPI.actionlogList), so there is
// nothing here to watch/load. This file's only job is deriving
// (rootProjectId, revisionId) from the `project` prop and passing them
// through, same as ManageOverview.vue.
//
// REVISION FILTERING CAVEAT: the actionlog endpoints do not accept a
// revisionID filter yet (it is landing separately, server-side) — until it
// does, `revisionId` is wired through but has no effect, so this panel shows
// the same chain-wide events the dashboard does. Even once the filter lands,
// it will not show everything for this revision: work items are filed
// against the chain ROOT (project.rootProjectID), not a revision, so their
// audit events are never revision-attributable — only this revision's own
// build artifacts (namespaces, workflows, agents, etc.) will actually narrow.
// See ActivityPanel.vue's own prop comment for the full detail.
import ActivityPanel from '@/sections/project/components/dashboard/ActivityPanel.vue'
import { computed } from 'vue'

const props = defineProps({
  // The OPEN revision (a lib system.Project instance) — route.params.projectId
  // IS the revision row; its rootProjectID is the chain root work items are
  // filed against (see stores/projects.js's Project class).
  project: { type: Object, required: true },
  // Passed by Wizard.vue to every M&M section, but accepted only for the
  // shared dispatch prop set (see ManagePlaceholder.vue) — this section is
  // read-only, same as ManageOverview.vue, so there is nothing to disable.
  disabled: { type: Boolean, default: false },
})

const rootProjectId = computed(() => props.project?.rootProjectID || props.project?.projectID)
const revisionId = computed(() => props.project?.projectID)
</script>
