<script setup>
import { RouterLink } from 'vue-router'

defineOptions({
  inheritAttrs: false,
})

defineProps({
  to: {
    type: [String, Object],
    required: true,
  },
})
</script>

<template>
  <Button :as="RouterLink" :pt="{ root: { to } }" v-bind="$attrs">
    <slot />
  </Button>
</template>
