<template>
  <Drawer
    v-model:visible="expanded"
    :modal="isMobile"
    :dismissable="isMobile"
    :pt="{
      root: 'border-r border-surface',
      header: { class: 'pl-3 pt-3 pb-2 pr-1 h-15 gap-2' },
      content: 'p-3',
    }"
    :style="{ width: 'var(--sidebar-width)' }"
  >
    <template #header>
      <div class="grow">
        <img :src="logo" class="flex-1 object-contain" />
      </div>
    </template>

    <!-- Sidebar content area with named slots -->
    <div class="flex flex-col h-full gap-2">
      <slot name="header" />

      <div class="flex-1 overflow-auto">
        <slot name="body" />
      </div>

      <!-- Default slot for any extra content -->
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
