<template>
  <BoardPanel
    :project-id="rootProjectId"
    :revision-id="revisionId"
    :disabled="disabled"
  />
</template>

<script setup>
// Manage & Monitor > Monitor > Board — the dashboard's board screen
// (components/dashboard/BoardPanel.vue), scoped to the ONE revision open in
// the wizard (route.params.projectId), same revision-scoping contract as
// ManageOverview.vue/ManageActivity.vue/the Manage<Category>.vue files.
// Column/card components (BoardColumn.vue/BoardCard.vue) moved alongside
// BoardPanel under components/dashboard/ in the same extraction — see
// dashboard.intent.md; there is exactly one board implementation.
//
// Loads its own events/backlog stores, scoped to (root project, this
// revision) — same as the Manage<Category>.vue files (the "one file per
// section" load duplication is deliberate, see wizard.intent.md). BoardPanel
// itself never calls store.load(); it only reads their already-loaded state.
import BoardPanel from '@/sections/project/components/dashboard/BoardPanel.vue'
import { useBacklogItemsStore } from '@/sections/project/stores/backlogItems'
import { useEventsStore } from '@/sections/project/stores/events'
import { computed, watch } from 'vue'

const props = defineProps({
  // The OPEN revision (a lib system.Project instance) — route.params.projectId
  // IS the revision row; its rootProjectID is the chain root work items are
  // filed against (see stores/projects.js's Project class).
  project: { type: Object, required: true },
  // Passed by Wizard.vue to every M&M section, but deliberately NOT enforced
  // here: creating and moving work items is ungated, matching the
  // dashboard's own New-event button (CategoryPanel.vue), which has never
  // been capability-gated. Ruled 2026-07-28 — revisit when member roles and
  // what each may do are refined.
  disabled: { type: Boolean, default: false },
})

const eventsStore = useEventsStore()
const backlogStore = useBacklogItemsStore()

const rootProjectId = computed(() => props.project?.rootProjectID || props.project?.projectID)
const revisionId = computed(() => props.project?.projectID)

// Scope both dashboard stores to (root project, this revision) — filed
// against the root, assigned to the revision, like an issue against a
// milestone. Re-fetches whenever the open revision changes (e.g. switching
// revisions via the wizard-header revision switcher).
watch(
  [rootProjectId, revisionId],
  ([pid, rid]) => {
    if (!pid || !rid) return
    eventsStore.load(pid, rid)
    backlogStore.load(pid, rid)
  },
  { immediate: true },
)
</script>
