<template>
  <div class="flex flex-col h-full p-5">
    <h1 class="text-2xl font-semibold mb-4">Page Builder</h1>
    <p class="text-muted-color">Page ID: {{ pageID }}</p>
    <!-- TODO: Implement page builder -->
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'

const route = useRoute()
const pageID = route.params.pageID
</script>
