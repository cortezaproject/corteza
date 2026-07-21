<template>
  <div
    class="inline-flex rounded-lg border px-2.5 py-1.5 items-center gap-2 whitespace-nowrap transition-opacity bg-transparent opacity-100"
  >
    <span
      class="inline-flex items-center justify-center w-6 h-6 rounded-md ring-1 shrink-0"
      :class="[bg, ring]"
    >
      <i :class="[icon, 'text-xs']" />
    </span>
    <span class="text-xs text-muted-color leading-none">{{ label }}</span>
  </div>
</template>

<script setup>
const props = defineProps({
  icon: { type: String, default: 'pi pi-check' },
  bg: { type: String, default: 'bg-gray-200' },
  ring: { type: String, default: 'ring-gray-200' },
  label: { type: String, required: true },
})
</script>
<style scoped></style>
