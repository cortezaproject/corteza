<template>
  <div class="flex flex-col h-full">
    <CSidebarNamespaceNav v-if="$route.meta.sidebar === 'namespaces'" />
    <template v-else>
      <CSidebarNamespaceSwitcher />
      <CSidebarNavigation />
    </template>
  </div>
</template>

<script setup>
// Compose's sidebar body, in two shapes the route picks between (`meta.sidebar`,
// so a route rename cannot silently change which one renders): inside a
// namespace it is the switcher + that namespace's nav tree; on the namespace
// editor it is the list of namespaces.
import CSidebarNamespaceNav from '@/sections/compose/components/CSidebarNamespaceNav.vue'
import CSidebarNamespaceSwitcher from '@/sections/compose/components/CSidebarNamespaceSwitcher.vue'
import CSidebarNavigation from '@/sections/compose/components/CSidebarNavigation.vue'
</script>
