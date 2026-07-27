<template>
  <ManagePlaceholder label-key="project.manage.views.activity" v-bind="$attrs" />
</template>

<script setup>
// Manage & Monitor > Monitor > Activity — not built yet. Own file (rather
// than a shared dispatch entry) so the activity build-out can replace this in
// isolation; see ManagePlaceholder.vue for the shared shell and
// views/Wizard.vue for the section dispatch map.
import ManagePlaceholder from './ManagePlaceholder.vue'
</script>
