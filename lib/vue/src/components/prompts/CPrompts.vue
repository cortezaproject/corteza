<template>
  <div>
    <c-prompt-toast
      :hide-toasts="hideToasts"
    />
    <c-prompt-modal />
  </div>
</template>
<script>
import CPromptToast from './CPromptToast.vue'
import CPromptModal from './CPromptModal.vue'

/**
 * Container for prompt toast and modal
 *
 * Toast serves as a notifier and placeholder for light prompts
 * Modal serves as prompt lister and placeholder all prompt types
 *
 * All communication between the two components go through wf-prompt state
 */
export default {
  name: 'c-prompts',

  components: {
    CPromptToast,
    CPromptModal
  },

  props: {
    hideToasts: {
      type: Boolean,
    },
  },
}
</script>
