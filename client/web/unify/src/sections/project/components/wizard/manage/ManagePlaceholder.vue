<template>
  <!-- Shared shell for the Manage & Monitor sections that don't have real
       content yet (metrics/activity/the five per-category panels — board is
       the only section with a real panel so far, see ManageBoard.vue). Each
       section still gets its own thin wrapper file (Manage<Key>.vue) rather
       than sharing one dispatch entry, so the agents building those out next
       can each replace their own file without touching this one or Wizard.vue. -->
  <div class="h-full overflow-y-auto p-4">
    <h2 class="text-lg font-medium mb-1">{{ $t(labelKey) }}</h2>
    <p class="text-sm text-muted-color">{{ $t('project.manage.panel.comingSoon') }}</p>
  </div>
</template>

<script setup>
defineProps({
  // i18n key for the section's own nav label (config/manageNav.js labelKey).
  labelKey: { type: String, required: true },
  // Accepted (unused) so Wizard.vue's dispatch can pass the same prop set to
  // every Manage & Monitor section regardless of which one is mounted — real
  // panels like ManageBoard.vue use these, placeholders don't yet.
  project: { type: Object, default: null },
  disabled: { type: Boolean, default: false },
})
</script>
