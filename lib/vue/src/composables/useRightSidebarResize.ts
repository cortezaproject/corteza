import { ref } from 'vue'

export function useRightSidebarResize(defaultWidth = 360) {
  const drawerWidth = ref(defaultWidth)
  const isResizingDrawer = ref(false)

  function startDrawerResize() {
    isResizingDrawer.value = true
    document.addEventListener('mousemove', resizeDrawer)
    document.addEventListener('mouseup', stopDrawerResize)
  }

  function resizeDrawer(e: MouseEvent) {
    if (!isResizingDrawer.value) return
    const newWidth = window.innerWidth - e.clientX
    drawerWidth.value = Math.max(280, Math.min(800, newWidth))
  }

  function stopDrawerResize() {
    isResizingDrawer.value = false
    document.removeEventListener('mousemove', resizeDrawer)
    document.removeEventListener('mouseup', stopDrawerResize)
  }

  return {
    drawerWidth,
    isResizingDrawer,
    startDrawerResize,
  }
}
