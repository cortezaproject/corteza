<template>
  <small v-if="message" class="block mt-1 text-red-500">{{ message }}</small>
</template>

<script setup>
// The one way an inline validation message looks in this section. Renders
// nothing when there's no message, so call sites don't need a v-if.
defineProps({
  message: { type: String, default: '' },
})
</script>
