<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ project?.name || 'Project' }}</span>
  </Teleport>

  <div class="h-full w-full flex flex-col items-center justify-center gap-3 text-surface-500">
    <i class="pi pi-folder-open text-4xl" />
    <h1 class="text-xl font-medium">{{ project?.name || 'Project' }} — overview</h1>
    <p class="text-sm">Read-only overview / versions for published projects (later screen).</p>
    <div class="flex gap-2">
      <Button label="Back to projects" icon="pi pi-arrow-left" text size="small" @click="router.push({ name: 'project.list' })" />
      <Button label="Open wizard" icon="pi pi-sliders-h" size="small" @click="router.push({ name: 'project.wizard', params: { projectId: route.params.projectId } })" />
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
