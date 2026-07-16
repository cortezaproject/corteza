<template>
  <div v-if="show" class="fixed body-bg inset-0 flex items-center justify-center z-[9999] p-6">
    <img v-if="logoUrl" :src="logoUrl" class="max-h-48 max-w-xl w-auto animate-pulse" />
  </div>
</template>

<script setup>
defineProps({
  show: {
    type: Boolean,
    default: true,
  },

  logoUrl: {
    type: String,
    default: '',
  },
})
</script>
