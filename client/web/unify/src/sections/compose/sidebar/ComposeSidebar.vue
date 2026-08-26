<template>
  <div class="flex flex-col h-full">
    <!-- The column fits the shell's scroller exactly, so that one never engages
         and the switcher stays put: the nav takes the height its siblings leave
         and scrolls what is inside it. -->
    <CSidebarNamespaceNav v-if="$route.meta.sidebar === 'namespaces'" class="flex-1 min-h-0" />
    <template v-else>
      <CSidebarNamespaceSwitcher />
      <CSidebarNavigation class="flex-1 min-h-0" />
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
