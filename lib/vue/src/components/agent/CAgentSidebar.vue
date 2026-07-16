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
      v-if="isVisible"
      class="right-sidebar flex flex-row"
      :style="{ width: `${drawerWidth}px` }"
    >
      <CResizeHandle @mousedown="startDrawerResize" />
      <div class="flex-1 flex flex-col min-w-0">
        <CAgentChat
          :translations="translations"
          :context-provider="contextProvider"
        >
          <template #actions>
            <Button
              icon="pi pi-times"
              severity="secondary"
              variant="text"
              rounded
              size="small"
              @click="isVisible = false"
            />
          </template>
        </CAgentChat>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAgentChatStore } from '../../stores/useAgentChatStore'
import { useRightSidebarStore } from '../../stores/useRightSidebarStore'
import { useRightSidebarResize } from '../../composables/useRightSidebarResize'
import CAgentChat from './CAgentChat.vue'
import { makeAgentChatTranslations } from './translations'
import CResizeHandle from '../CResizeHandle.vue'

defineProps({
  // Optional callable returning an object that gets passed as `context` to
  // agentExec on every send. Each webapp's App.vue supplies its own — at
  // minimum the current route's params; richer providers can also resolve
  // those IDs into named entities via app-local stores.
  contextProvider: {
    type: Function as unknown as () => () => Record<string, any> | undefined,
    default: null,
  },
})

const { t } = useI18n()
const translations = computed(() => makeAgentChatTranslations(t, 'agent.sidebar.'))

const agentStore = useAgentChatStore()
const rightSidebarStore = useRightSidebarStore()
const { drawerWidth, startDrawerResize } = useRightSidebarResize()

const isVisible = computed({
  get: () => rightSidebarStore.isOpen('agent'),
  set: value => value ? rightSidebarStore.open('agent') : rightSidebarStore.close('agent'),
})
</script>
