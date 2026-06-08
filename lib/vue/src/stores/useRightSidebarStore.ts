import { defineStore } from 'pinia'
import { ref } from 'vue'

const GLOBAL_PANELS = new Set(['notifications', 'agent'])

export const useRightSidebarStore = defineStore('rightSidebar', () => {
  const activePanel = ref<string | null>(null)

  function isOpen(name: string): boolean {
    return activePanel.value === name
  }

  function open(name: string) {
    activePanel.value = name
  }

  function close(name: string) {
    if (activePanel.value === name) activePanel.value = null
  }

  function toggle(name: string) {
    activePanel.value = activePanel.value === name ? null : name
  }

  function closeAll() {
    activePanel.value = null
  }

  function closeSectionPanels() {
    if (activePanel.value && !GLOBAL_PANELS.has(activePanel.value)) {
      activePanel.value = null
    }
  }

  return { activePanel, isOpen, open, close, toggle, closeAll, closeSectionPanels }
})
