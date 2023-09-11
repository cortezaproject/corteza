<template>
  <div
    v-if="text"
    v-b-popover.hover.right="{ variant: 'dark', content: text }"
    class="d-flex align-items-center ml-2"
  >
    <font-awesome-icon
      :icon="['far', 'question-circle']"
    />
  </div>
</template>

<script>
export default {
  props: {
    text: {
      type: String,
      default: '',
    },
  },
}
</script>
