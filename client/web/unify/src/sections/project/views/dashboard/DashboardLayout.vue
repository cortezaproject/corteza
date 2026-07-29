<template>
  <!-- Topbar: the current project's name + the shared revision switcher
       (components/project/RevisionSwitcher.vue — same component the Wizard
       header mounts). Replaces the old static version Tag and the old
       jump-back-to-wizard button: the switcher's own "Dashboard" entry
       (marked current here) and its per-revision entries cover both. This
       layout owns the topbar — children must not Teleport into it (see
       DashboardLayout.intent.md). (Back-to-projects lives in the global
       sidebar's projects list.) -->
  <Teleport to="#topbar-title" defer>
    <span class="flex items-center gap-2">
      <span>{{ project?.name || $t('project.overview.fallbackName') }}</span>
      <RevisionSwitcher :project="project" />
    </span>
  </Teleport>

  <!-- Right-aligned topbar tools: the shared Members/View project cluster
       (components/project/ProjectTopbarTools.vue — same component the
       Wizard header mounts), replacing this view's own bespoke namespace
       link. The dashboard gains a Members button as a result. -->
  <Teleport to="#topbar-tools" defer>
    <span class="flex items-center gap-2">
      <ProjectTopbarTools :project="project" />
    </span>
  </Teleport>

  <!-- Left rail (in-view sub-nav) + the active view. The app's global sidebar
       is unchanged; this mirrors the Wizard's StepNav-beside-content layout. -->
  <div class="h-full w-full flex gap-4 p-3 min-h-0">
    <aside class="w-72 shrink-0 h-full">
      <DashboardNav />
    </aside>
    <section class="flex-1 min-w-0 min-h-0 overflow-hidden rounded-xl border border-surface">
      <RouterView />
    </section>
  </div>
</template>

<script setup>
import DashboardNav from '@/sections/project/components/dashboard/DashboardNav.vue'
import ProjectTopbarTools from '@/sections/project/components/project/ProjectTopbarTools.vue'
import RevisionSwitcher from '@/sections/project/components/project/RevisionSwitcher.vue'
import { useBacklogItemsStore } from '@/sections/project/stores/backlogItems'
import { useEventsStore } from '@/sections/project/stores/events'
import { chainHasPublished } from '@/sections/project/config/publishState'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const store = useProjectsStore()
const eventsStore = useEventsStore()
const backlogStore = useBacklogItemsStore()

const project = computed(() => store.findById(route.params.projectId))

// Load the project itself on entry and whenever the active project changes.
// The project must be fetched here (not left to the sidebar's store load) —
// the sidebar lives in a lazily-mounted drawer, so landing on the dashboard
// with it collapsed would otherwise leave `project` undefined and hide the
// topbar buttons until the drawer is first opened.
watch(
  () => route.params.projectId,
  id => {
    if (!id) return
    store.fetchProject(id).catch(err => console.error('Failed to load project', err))
  },
  { immediate: true },
)

// A project that has never published has nothing to report on, so it has no
// dashboard: send it to its wizard, which IS the whole product until the first
// publish. The list, the sidebar and the revision switcher all route around
// this already (same rule, config/publishState.js) — this covers the door they
// can't: a typed or bookmarked URL. `replace` so Back doesn't bounce off the
// redirect, and it waits for the fetch above rather than firing on the
// undefined project of a cold landing.
watch(project, p => {
  if (!p || chainHasPublished(p)) return
  router.replace({ name: 'project.wizard', params: { projectId: p.projectID } })
})

// CHAIN-WIDE SCOPE (project.intent.md "Dashboards"): the dashboard covers the
// whole revision chain, not just the one revision named in the route — the
// route's projectId is a chain HEAD, which is only the chain ROOT for an
// original, never-revised project (see system.Project#rootProjectID, which
// falls back to its own id). Load events/backlog for that root with no
// revisionId, which both stores read as "every revision in the project" (see
// stores/events.js#load / stores/backlogItems.js#load) — the wizard's Manage
// & Monitor tab is the one-revision counterpart, scoped by revisionId
// instead. Watches the resolved root rather than the raw route param, so a
// route landing on a non-root head still reloads once the project fetch
// above resolves its rootProjectID.
const chainRootId = computed(() => project.value?.rootProjectID || route.params.projectId)
watch(
  chainRootId,
  id => {
    if (!id) return
    eventsStore.load(id)
    backlogStore.load(id)
  },
  { immediate: true },
)
</script>
