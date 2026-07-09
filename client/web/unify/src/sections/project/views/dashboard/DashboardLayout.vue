<template>
  <!-- Topbar: the current project's name + which revision we're viewing.
       (Back-to-projects lives in the global sidebar's projects list.) -->
  <Teleport to="#topbar-title" defer>
    <span class="flex items-center gap-2">
      <span>{{ project?.name || $t('project.overview.fallbackName') }}</span>
      <span
        v-if="project"
        class="text-[11px] font-medium rounded-md px-2 py-0.5 bg-surface-200 text-muted-color dark:bg-surface-700"
      >
        {{ versionLabel }}
      </span>
    </span>
  </Teleport>

  <!-- Right-aligned topbar tool: open the project's compose namespace. -->
  <Teleport to="#topbar-tools" defer>
    <Button
      v-if="project?.namespaceID"
      :label="$t('project.viewProject')"
      icon="pi pi-external-link"
      severity="secondary"
      size="small"
      @click="openProject"
    />
  </Teleport>

  <!-- Left rail (in-view sub-nav) + the active view. The app's global sidebar
       is unchanged; this mirrors the Wizard's StepNav-beside-content layout. -->
  <div class="h-full w-full flex gap-4 p-4 min-h-0">
    <aside class="w-60 shrink-0 h-full">
      <DashboardNav />
    </aside>
    <section class="flex-1 min-h-0 overflow-y-auto rounded-xl border border-surface bg-surface-0 dark:bg-surface-900">
      <RouterView />
    </section>
  </div>
</template>

<script setup>
import DashboardNav from '@/sections/project/components/dashboard/DashboardNav.vue'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const store = useProjectsStore()

const project = computed(() => store.findById(route.params.projectId))

// Open the project's compose namespace (namespace.view resolves the URL part by
// slug OR namespaceID, so the ID is a reliable target).
function openProject() {
  if (project.value?.namespaceID) {
    router.push({ name: 'namespace.view', params: { slug: project.value.namespaceID } })
  }
}

// Which version of the project we're viewing. User-facing versions are 1-based
// (the original live project is v1), so we display the backend revision + 1.
const versionLabel = computed(() =>
  t('project.dashboard.version', { number: (project.value?.revision ?? 0) + 1 }),
)
</script>
