<template>
  <!-- Small uppercase label sitting above a dialog header title, e.g. "Role",
       "Module" — used across the project section's detail/create dialogs. -->
  <div class="text-[10px] uppercase tracking-wider text-muted-color leading-none mb-0.5">
    <slot />
  </div>
</template>

<script setup>
defineOptions({ name: 'DialogEyebrow' })
</script>
