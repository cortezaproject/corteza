<template>
  <div class="flex flex-col min-h-full">
    <!-- `min-h-full`, not `h-full`: the column fills the shell's scroller when
         there is little to show, so the nav can hold its admin panel at the
         foot, and grows past it when there is more — the shell's scroller is
         the only one, and everything above scrolls with the rest. -->
    <CSidebarNamespaceNav v-if="$route.meta.sidebar === 'namespaces'" class="flex-1" />
    <template v-else>
      <CSidebarNamespaceSwitcher />
      <CSidebarNavigation class="flex-1" />
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
