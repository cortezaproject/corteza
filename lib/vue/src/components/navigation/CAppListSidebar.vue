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
      v-if="rightSidebarStore.isOpen('app-list')"
      class="right-sidebar flex"
      :style="{ width: `${drawerWidth}px` }"
    >
      <CResizeHandle @mousedown="startDrawerResize" />
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
import { ref } from 'vue'
import CInputSearch from '../input/CInputSearch.vue'
import CResizeHandle from '../CResizeHandle.vue'
import CAppList from './CAppList.vue'
import { useRightSidebarResize } from '../../composables/useRightSidebarResize'
import { useRightSidebarStore } from '../../stores/useRightSidebarStore'

defineProps({
  labels: {
    type: Object,
    required: true,
  },
})

const rightSidebarStore = useRightSidebarStore()
const query = ref('')
const { drawerWidth, startDrawerResize } = useRightSidebarResize()

const close = () => {
  rightSidebarStore.close('app-list')
  query.value = ''
}
</script>
