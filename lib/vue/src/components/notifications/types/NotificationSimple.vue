<template>
  <div :class="{ 'opacity-60': notification.readAt }">
    <div class="font-semibold text-color">
      {{ notification.config?.title || '' }}
    </div>
    <div class="rt-content mt-1 text-sm text-muted-color" v-html="description" />
  </div>
</template>

<script setup>
import { sanitizeHtml } from '@planetcrust/human-js'
import { computed } from 'vue'

const props = defineProps({
  notification: {
    type: Object,
    required: true,
  },
})

// The description is HTML somebody else wrote — a workflow, or whoever made
// the notification — and the server stores it as it came.
const description = computed(() => sanitizeHtml(props.notification.config?.description || ''))
</script>
