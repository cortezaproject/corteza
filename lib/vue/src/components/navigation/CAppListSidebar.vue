<template>
  <Transition
    enter-active-class="transition-transform duration-300 ease-in-out"
    enter-from-class="translate-x-full"
    enter-to-class="translate-x-0"
    leave-active-class="transition-transform duration-300 ease-in-out"
    leave-from-class="translate-x-0"
    leave-to-class="translate-x-full"
  >
    <div
      v-if="visible"
      class="right-sidebar flex"
      :style="{ width: `${drawerWidth}px` }"
    >
      <div
        class="resize-handle w-1 h-full cursor-ew-resize hover:bg-primary/20 transition-colors shrink-0"
        @mousedown="startDrawerResize"
      />
      <div class="flex-1 flex flex-col min-w-0">
        <!-- Header -->
      <div class="flex items-center justify-between pl-3 pt-3 pb-2 pr-1">
        <span class="text-lg font-semibold text-color">{{ labels.title }}</span>
        <Button
          icon="pi pi-times"
          severity="secondary"
          variant="text"
          rounded
          size="small"
          @click="close"
        />
      </div>

      <!-- Search -->
      <div class="p-4">
        <CInputSearch v-model="query" :placeholder="labels.search" size="small" />
      </div>

      <!-- App List -->
      <div v-if="areAppsVisible" class="flex-1 overflow-auto px-4 pb-4">
        <div class="flex flex-col gap-2">
          <a
            v-for="app in apps"
            :key="app.applicationID"
            v-show="isAppVisible(app)"
            :href="app.enabled ? getAppUrl(app) : '#'"
            target="_self"
            class="flex items-center gap-3 p-3 rounded-lg border border-surface hover:bg-emphasis hover:border-primary transition-all duration-150 no-underline text-color"
            @click="!app.enabled && $event.preventDefault()"
          >
            <img
              :src="getAppLogoUrl(app)"
              :alt="app.unify?.name || app.name"
              class="w-16 h-16 object-contain rounded-md shrink-0"
              loading="lazy"
            />
            <span class="font-medium text-sm truncate">
              {{ app.unify?.name || app.name }}
            </span>
          </a>
        </div>
      </div>

      <!-- Empty state -->
      <div v-else class="flex-1 flex items-center justify-center px-4">
        <span class="text-muted-color text-sm">
          {{ query ? labels.noResults : labels.noApps }}
        </span>
      </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { computed, inject, onBeforeUnmount, ref, watch } from 'vue'
import CInputSearch from '../input/CInputSearch.vue'
import { useApplicationsStore } from '../../stores/useApplicationsStore'
import { useRightSidebarResize } from '../../composables/useRightSidebarResize'
import defaultAppIcon from '../../assets/default-app.png'
import adminAreaIcon from '../../assets/admin-area.png'
import discoveryIcon from '../../assets/discovery.png'
import lowCodeCrmIcon from '../../assets/low-code-crm-app.png'
import lowCodePlatformIcon from '../../assets/low-code-platform.png'
import lowCodeServiceIcon from '../../assets/low-code-service-solution-app.png'
import privacyIcon from '../../assets/privacy.png'
import reporterIcon from '../../assets/reporter.png'
import videoConferenceIcon from '../../assets/video-conference.png'
import workflowsIcon from '../../assets/workflows.png'

const appIconMap = {
  'applications/default-app.png': defaultAppIcon,
  'applications/admin-area.png': adminAreaIcon,
  'applications/discovery.png': discoveryIcon,
  'applications/low-code-crm-app.png': lowCodeCrmIcon,
  'applications/low-code-platform.png': lowCodePlatformIcon,
  'applications/low-code-service-solution-app.png': lowCodeServiceIcon,
  'applications/privacy.png': privacyIcon,
  'applications/reporter.png': reporterIcon,
  'applications/video-conference.png': videoConferenceIcon,
  'applications/workflows.png': workflowsIcon,
}

const visible = defineModel('visible', {
  type: Boolean,
  default: false,
})

const props = defineProps({
  labels: {
    type: Object,
    required: true,
  },
})

const $SystemAPI = inject('$SystemAPI')
const $eventBus = inject('$eventBus', null)

const query = ref('')
const applicationsStore = useApplicationsStore()
const { drawerWidth, startDrawerResize } = useRightSidebarResize()

const apps = computed(() => applicationsStore.unifyOnly)

const normalizedQuery = computed(() => (query.value || '').trim().toUpperCase())

const isAppVisible = app => {
  const q = normalizedQuery.value
  if (!q) return true
  return (
    (app.name?.toUpperCase() || '').includes(q) ||
    (app.unify?.name?.toUpperCase() || '').includes(q)
  )
}

const areAppsVisible = computed(() => apps.value.some(isAppVisible))

const getAppLogoUrl = app => {
  if (!app.unify?.logo) {
    return defaultAppIcon
  }

  // Check if it's a known bundled app icon
  if (appIconMap[app.unify.logo]) {
    return appIconMap[app.unify.logo]
  }

  const apiSystem = '/api/system'
  const apiBaseUrl = new URL($SystemAPI.baseURL, window.location.origin).toString()

  if (app.unify.logo.startsWith(apiSystem)) {
    return apiBaseUrl.substring(0, apiBaseUrl.length - apiSystem.length) + app.unify.logo
  }

  return app.unify.logo
}

const getAppUrl = app => {
  const url = app.unify?.url || ''
  if (!url || url.startsWith('/') || url.startsWith('http')) return url
  return '/' + url
}

const close = () => {
  visible.value = false
  query.value = ''
}

const offSidebar = $eventBus?.on('right-sidebar:opened', name => {
  if (name !== 'app-list') {
    close()
  }
})

watch(visible, next => {
  if (next) {
    $eventBus?.emit('right-sidebar:opened', 'app-list')
  }
})

onBeforeUnmount(() => {
  offSidebar?.()
})
</script>
