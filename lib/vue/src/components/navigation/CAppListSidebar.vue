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
        <div class="flex items-center gap-2 px-3 py-2 border-b border-surface shrink-0">
          <i class="pi pi-th-large text-primary text-sm" />
          <span class="font-semibold text-base text-color flex-1">{{ labels.title }}</span>
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
        <div class="flex-1 overflow-auto px-4 pb-4">
          <CAppList
            :query="query"
            variant="list"
            :no-apps-text="labels.noApps"
            :no-results-text="labels.noResults"
          />
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { inject, onBeforeUnmount, ref, watch } from 'vue'
import CInputSearch from '../input/CInputSearch.vue'
import CAppList from './CAppList.vue'
import { useRightSidebarResize } from '../../composables/useRightSidebarResize'

const visible = defineModel('visible', {
  type: Boolean,
  default: false,
})

defineProps({
  labels: {
    type: Object,
    required: true,
  },
})

const $eventBus = inject('$eventBus', null)

const query = ref('')
const { drawerWidth, startDrawerResize } = useRightSidebarResize()

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
