<template>
  <ActivityPanel :project-id="projectId" />
</template>

<script setup>
// Route target for the dashboard's Activity screen (project.overview.activity,
// path `activity`) — the route-target view family is locked (see
// views.intent.md), so this file stays as the thin wrapper: it just resolves
// route.params into props for the shared, route-free panel
// (components/dashboard/ActivityPanel.vue), which also backs the wizard's
// Manage & Monitor "Activity" section (components/wizard/manage/
// ManageActivity.vue), revision-scoped there. Mirrors CategoryView.vue's own
// extraction (see dashboard.intent.md). No `revisionId` here — the dashboard
// is chain-wide.
import ActivityPanel from '@/sections/project/components/dashboard/ActivityPanel.vue'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const projectId = computed(() => route.params.projectId || undefined)
</script>
