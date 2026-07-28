<template>
  <OverviewPanel :project-id="route.params.projectId" />
</template>

<script setup>
// Thin route wrapper — this view stays the routed target for
// `project.overview` (the dashboard view set is locked, see
// DashboardLayout.intent.md), but its body now lives in the reusable
// OverviewPanel (components/dashboard/OverviewPanel.vue), which the wizard's
// Manage & Monitor tab also mounts, revision-scoped (see
// components/wizard/manage/ManageOverview.vue) — ruled 2026-07-28: the
// project dashboard and the wizard's Manage & Monitor tab are ONE surface
// differing only in scope (see views/views.intent.md / DashboardLayout.intent.md).
// This wrapper's only job is resolving the route's project id into
// OverviewPanel's `projectId` prop — no `revisionId` is passed, so
// OverviewPanel runs its chain-wide path (every revision's rows, via the
// report endpoint), exactly this view's behaviour before the extraction.
// OverviewPanel itself stays free of route assumptions so it mounts the same
// way outside a routed view too.
import OverviewPanel from '@/sections/project/components/dashboard/OverviewPanel.vue'
import { useRoute } from 'vue-router'

const route = useRoute()
</script>
