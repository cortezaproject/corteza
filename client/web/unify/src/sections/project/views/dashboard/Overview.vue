<template>
  <OverviewPanel :project-id="route.params.projectId" />
</template>

<script setup>
// Thin route wrapper — the routed target for `project.overview` (the dashboard
// view set is locked, see DashboardLayout.intent.md). The body lives in
// components/dashboard/OverviewPanel.vue, which the wizard's Manage & Monitor
// tab also mounts, revision-scoped (components/wizard/manage/ManageOverview.vue).
// This wrapper's only job is resolving the route's project id into
// OverviewPanel's `projectId` prop. No `revisionId` is passed, so OverviewPanel
// runs its chain-wide path (every revision's rows, via the report endpoint).
import OverviewPanel from '@/sections/project/components/dashboard/OverviewPanel.vue'
import { useRoute } from 'vue-router'

const route = useRoute()
</script>
