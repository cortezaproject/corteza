<template>
  <div class="flex flex-col h-full p-5">
    <h1 class="text-2xl font-semibold mb-4">Create Record</h1>
    <p class="text-muted-color">Module ID: {{ moduleID }}</p>
    <!-- TODO: Implement create record form -->
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'

const route = useRoute()
const moduleID = route.params.moduleID
</script>
