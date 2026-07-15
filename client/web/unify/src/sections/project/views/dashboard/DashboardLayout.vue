<template>
  <!-- Topbar: the current project's name + which revision we're viewing.
       (Back-to-projects lives in the global sidebar's projects list.) -->
  <Teleport to="#topbar-title" defer>
    <span class="flex items-center gap-2">
      <span>{{ project?.name || $t('project.overview.fallbackName') }}</span>
      <Tag v-if="project" :value="versionLabel" severity="secondary" class="!text-xs" />
    </span>
  </Teleport>

  <!-- Right-aligned topbar tools: jump back to the build wizard (project
       overview); always offer opening the project's compose namespace. -->
  <Teleport to="#topbar-tools" defer>
    <span class="flex items-center gap-2">
      <CRouterLinkButton
        v-if="project"
        :to="{ name: 'project.wizard', params: { projectId: route.params.projectId } }"
        :label="$t('project.viewOverview')"
        icon="pi pi-sitemap"
        size="small"
        severity="secondary"
        outlined
      />
      <CRouterLinkButton
        v-if="project?.hasNamespace"
        :to="{ name: 'namespace.view', params: { slug: project.namespaceID } }"
        :label="$t('project.viewProject')"
        icon="pi pi-external-link"
        size="small"
      />
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
import { useEventsStore } from '@/sections/project/stores/events'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { components } from '@planetcrust/human-vue'
import { computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

const { CRouterLinkButton } = components

const { t } = useI18n()
const route = useRoute()
const store = useProjectsStore()
const eventsStore = useEventsStore()

const project = computed(() => store.findById(route.params.projectId))

// Load the project itself + its events on entry and whenever the active project
// changes. The project must be fetched here (not left to the sidebar's store
// load) — the sidebar lives in a lazily-mounted drawer, so landing on the
// dashboard with it collapsed would otherwise leave `project` undefined and hide
// the topbar buttons until the drawer is first opened.
watch(
  () => route.params.projectId,
  id => {
    if (!id) return
    store.fetchProject(id).catch(err => console.error('Failed to load project', err))
    eventsStore.load(id)
  },
  { immediate: true },
)

// Which version of the project we're viewing. User-facing versions are 1-based
// (the original live project is v1), so we display the backend revision + 1.
const versionLabel = computed(() =>
  t('project.dashboard.version', { number: (project.value?.revision ?? 0) + 1 }),
)
</script>
