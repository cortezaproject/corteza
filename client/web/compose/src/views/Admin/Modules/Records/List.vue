<template>
  <div class="flex flex-col h-full p-5">
    <h1 class="text-2xl font-semibold mb-4">Module Records</h1>
    <p class="text-muted-color">Module ID: {{ moduleID }}</p>
    <!-- TODO: Implement module records list -->
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'

const route = useRoute()
const moduleID = route.params.moduleID
</script>
