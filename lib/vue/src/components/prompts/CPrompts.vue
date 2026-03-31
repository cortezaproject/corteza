<template>
  <div>
    <CPromptToast :hide-toasts="hideToasts" />
    <CPromptModal />
  </div>
</template>

<script setup>
import CPromptToast from './CPromptToast.vue'
import CPromptModal from './CPromptModal.vue'

defineProps({
  hideToasts: {
    type: Boolean,
    default: false,
  },
})
</script>
