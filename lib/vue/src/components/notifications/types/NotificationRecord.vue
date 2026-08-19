<template>
  <div :class="{ 'opacity-60': notification.readAt }">
    <div class="font-semibold text-color">
      {{ notification.config?.title || '' }}
    </div>
    <div
      class="rt-content mt-1 text-sm text-muted-color"
      v-html="notification.config?.description || ''"
    />
  </div>
</template>

<script setup>
defineProps({
  notification: {
    type: Object,
    required: true,
  },
})
</script>
