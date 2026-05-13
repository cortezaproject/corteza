<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('chatbot.sessions.title') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0">
    <SessionsList />
  </div>
</template>

<script setup>
import SessionsList from '@/components/SessionsList.vue'
</script>
