<template>
  <!-- Empty scaffold placeholder for a dashboard view. Content is built later;
       for now each view shows its title so the nav + routing are demonstrable. -->
  <div class="h-full flex flex-col items-center justify-center gap-3 text-center px-6 text-muted-color">
    <i :class="['pi text-4xl', icon]" />
    <h1 class="text-xl font-medium text-color">{{ $t(titleKey) }}</h1>
    <p class="text-sm">{{ $t('project.dashboard.stub.comingSoon') }}</p>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const titleKey = computed(() => route.meta.titleKey || 'project.dashboard.views.dashboard')
const icon = computed(() => route.meta.icon || 'pi-gauge')
</script>
