<template>
  <BoardPanel :project-id="projectId" />
</template>

<script setup>
// Thin route wrapper — routed target for `project.overview.board` (the dashboard
// view set is locked, see DashboardLayout.intent.md). The body lives in
// components/dashboard/BoardPanel.vue, which the wizard's Manage & Monitor Board
// section also mounts, revision-scoped (components/wizard/manage/ManageBoard.vue).
// Mirrors CategoryView.vue/AllEventsView.vue. No `revisionId` here — the
// dashboard is chain-wide, so BoardPanel reads eventsStore/backlogStore as
// already loaded chain-wide by DashboardLayout.vue.
import BoardPanel from '@/sections/project/components/dashboard/BoardPanel.vue'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const projectId = computed(() => route.params.projectId || undefined)
</script>
