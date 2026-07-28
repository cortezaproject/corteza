<template>
  <CategoryPanel :category="category" :project-id="route.params.projectId" />
</template>

<script setup>
// Thin route wrapper — this view stays the routed target for
// `project.overview.category` (the dashboard view set is locked, see
// DashboardLayout.intent.md), but its body now lives in the reusable
// CategoryPanel (components/dashboard/CategoryPanel.vue), which the wizard's
// Manage & Monitor category sections also mount, revision-scoped (see
// components/wizard/manage/Manage{Incident,Feature,Privacy,Task,Review}.vue).
// This wrapper's only job is resolving the route's category/project params
// into CategoryPanel's props — no `revisionId` is passed, so CategoryPanel
// runs its chain-wide path (report-endpoint metrics, every revision's rows),
// exactly this view's behaviour before the extraction. CategoryPanel itself
// stays free of route assumptions (an "unknown category" prop renders the
// same muted note it always has) so it mounts the same way outside a routed
// view too.
import CategoryPanel from '@/sections/project/components/dashboard/CategoryPanel.vue'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const category = computed(() => route.params.category)
</script>
