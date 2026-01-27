<template>
  <div>
    <Drawer
      v-model:visible="expanded"
      :modal="isMobile"
      :dismissable="isMobile"
      :pt="{ header: 'p-3 pr-2 gap-2' }"
    >
      <template #header>
        <div class="grow">
          <img :src="logo" class="w-auto h-full object-contain p-2" />
        </div>
      </template>
    </Drawer>

    <div class="tab flex items-center content-center absolute top-0 p-2">
      <img v-if="disabledRoutes.includes($route.name)" :src="icon" class="p-2" />

      <Button
        v-else-if="expandOnClick"
        data-test-id="button-sidebar-open"
        icon="pi pi-bars"
        variant="text"
        size="large"
        class="w-full"
        @click="expanded = true"
      />

      <Button v-else data-test-id="button-home" variant="text" size="large">
        <RouterLink :to="{ name: 'root' }">
          <i class="pi pi-home" />
        </RouterLink>
      </Button>
    </div>
  </div>
</template>

<script setup>
import { throttle } from 'lodash-es'
import Button from 'primevue/button'
import Drawer from 'primevue/drawer'
import { computed, inject, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

const props = defineProps({
  disabledRoutes: {
    type: Array,
    default: () => [],
  },
  expandOnClick: {
    type: Boolean,
    default: false,
  },
})

const expanded = defineModel()

const $Settings = inject('$Settings')
const route = useRoute()

const isMobile = ref(false)

const icon = computed(() => {
  return $Settings.attachment('ui.iconLogo')
})

const logo = computed(() => {
  return $Settings.attachment('ui.mainLogo')
})

const checkIfMobile = throttle(() => {
  isMobile.value = window.innerWidth < 1024
}, 500)

const checkSidebar = (initial = false) => {
  // If sidebar should be disabled on route, close and unpin when navigating to route
  if (props.disabledRoutes.includes(route.name)) {
    expanded.value = false
  } else if (!isMobile.value && initial) {
    expanded.value = true
  }
}

// Watch for route changes
watch(
  () => route.name,
  () => {
    checkSidebar()
  },
)

// Watch for disabled routes changes
watch(
  () => props.disabledRoutes,
  () => {
    checkSidebar()
  },
  { deep: true },
)

onMounted(() => {
  checkSidebar(true)
  checkIfMobile()
  window.addEventListener('resize', checkIfMobile)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', checkIfMobile)
})
</script>

<style scoped>
.tab {
  height: var(--topbar-height);
  width: 66px;
}
</style>
