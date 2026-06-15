<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ project?.name || $t('project.overview.fallbackName') }}</span>
  </Teleport>

  <div class="h-full w-full flex flex-col items-center justify-center gap-3 text-surface-500">
    <i class="pi pi-folder-open text-4xl" />
    <h1 class="text-xl font-medium">{{ $t('project.overview.title', { name: project?.name || $t('project.overview.fallbackName') }) }}</h1>
    <p class="text-sm">{{ $t('project.overview.blurb') }}</p>
    <div class="flex gap-2">
      <Button :label="$t('project.overview.backToProjects')" icon="pi pi-arrow-left" text size="small" @click="router.push({ name: 'project.list' })" />
      <Button :label="$t('project.overview.openWizard')" icon="pi pi-sliders-h" size="small" @click="router.push({ name: 'project.wizard', params: { projectId: route.params.projectId } })" />
    </div>
  </div>
</template>

<script setup>
import { useProjectsStore } from '@/sections/project/stores/projects'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const store = useProjectsStore()

const project = computed(() => store.findById(route.params.projectId))
</script>
