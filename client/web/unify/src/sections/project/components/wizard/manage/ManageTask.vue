<template>
  <CategoryPanel category="task" :revision-id="revisionId" />
</template>

<script setup>
// Manage & Monitor > Categories > Task — the dashboard's task category screen
// (components/dashboard/CategoryPanel.vue), scoped to the ONE revision open in
// the wizard (route.params.projectId), same revision-scoping contract as
// ManageBoard.vue/ManageMetrics.vue: eventsStore/backlogStore are loaded with
// (root project, this revision), never the whole project (that's the live
// dashboard's scope; see project.intent.md's "Dashboards" locked contract).
// Own file (rather than a shared dispatch entry) so this category stays
// independently editable; see views/Wizard.vue for the section dispatch map.
// The load-effect duplication across the five Manage<Category>.vue files (and
// ManageBoard.vue/ManageMetrics.vue) is deliberate, not an oversight — see
// wizard.intent.md's "one file per section" rationale.
import CategoryPanel from '@/sections/project/components/dashboard/CategoryPanel.vue'
import { useBacklogItemsStore } from '@/sections/project/stores/backlogItems'
import { useEventsStore } from '@/sections/project/stores/events'
import { computed, watch } from 'vue'

const props = defineProps({
  // The OPEN revision (a lib system.Project instance) — route.params.projectId
  // IS the revision row; its rootProjectID is the chain root work items are
  // filed against (see stores/projects.js's Project class).
  project: { type: Object, required: true },
  // Passed by Wizard.vue to every M&M section, but deliberately NOT enforced
  // here: creating/editing/deleting items is ungated, matching the dashboard's
  // own CategoryView and ManageBoard.vue's quick-add.
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
