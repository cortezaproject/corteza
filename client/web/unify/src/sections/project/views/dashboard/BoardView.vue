<template>
  <BoardPanel :project-id="projectId" />
</template>

<script setup>
// Thin route wrapper — routed target for `project.overview.board` (the
// dashboard view set is locked, see DashboardLayout.intent.md; PINNED
// 2026-07-28 to include Board). Its body lives in the reusable BoardPanel
// (components/dashboard/BoardPanel.vue), which the wizard's Manage & Monitor
// Board section also mounts, revision-scoped (see components/wizard/manage/
// ManageBoard.vue) — ruled 2026-07-28: the project dashboard and the
// wizard's Manage & Monitor tab are ONE surface differing only in scope (see
// views/views.intent.md / DashboardLayout.intent.md). Mirrors
// CategoryView.vue/AllEventsView.vue's own extraction exactly. No
// `revisionId` here — the dashboard is chain-wide, so BoardPanel reads
// eventsStore/backlogStore as already loaded chain-wide by
// DashboardLayout.vue rather than loading them itself.
import BoardPanel from '@/sections/project/components/dashboard/BoardPanel.vue'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const projectId = computed(() => route.params.projectId || undefined)
</script>
