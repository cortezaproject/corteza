<template>
  <div class="flex flex-col h-full">
    <CSidebarNamespaceSwitcher />
    <CSidebarNavigation />
  </div>
</template>

<script setup>
// Compose's sidebar = namespace switcher + the namespace nav tree. (In the
// standalone app these were the CSidebar #header and #body; here they share
// the section sidebar body.)
import CSidebarNamespaceSwitcher from '@/sections/compose/components/CSidebarNamespaceSwitcher.vue'
import CSidebarNavigation from '@/sections/compose/components/CSidebarNavigation.vue'
</script>
