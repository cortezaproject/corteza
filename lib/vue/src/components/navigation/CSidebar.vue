<template>
  <Drawer
    v-model:visible="expanded"
    :modal="isMobile"
    :dismissable="isMobile"
    :pt="{
      header: { class: 'px-2 py-1', header: { class: 'max-h-15' } },
      content: 'p-3',
    }"
    :style="{ width: 'var(--sidebar-width)' }"
  >
    <template #header>
      <div class="grow">
        <img :src="logo" class="w-auto h-full max-h-16 object-contain" />
      </div>
    </template>

    <!-- Sidebar content area with teleport targets -->
    <div class="flex flex-col h-full gap-2">
      <!-- Header teleport target (e.g., namespace switcher) -->
      <div id="sidebar-header-expanded" />

      <!-- Body teleport target (e.g., navigation items) -->
      <div id="sidebar-body-expanded" class="flex-1 overflow-auto" />

      <!-- Footer teleport target -->
      <div id="sidebar-footer-expanded" />

      <!-- Default slot for direct content (fallback) -->
      <slot />
    </div>
  </Drawer>
</template>

<script setup>
import { throttle } from 'lodash-es'
import Drawer from 'primevue/drawer'
import { computed, inject, onBeforeUnmount, onMounted, ref } from 'vue'

const expanded = defineModel()

const $Settings = inject('$Settings')

const isMobile = ref(false)

const logo = computed(() => {
  return $Settings.attachment('ui.mainLogo')
})

const checkIfMobile = throttle(() => {
  isMobile.value = window.innerWidth < 1024
}, 500)

onMounted(() => {
  checkIfMobile()
  window.addEventListener('resize', checkIfMobile)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', checkIfMobile)
})
</script>
