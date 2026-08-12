import { onBeforeUnmount, onMounted, ref } from 'vue'

const STORAGE_KEY = 'agent-editor-chat-width'

const DEFAULT_CHAT_WIDTH = 480
const MIN_CHAT_WIDTH = 320
const MAX_CHAT_WIDTH = 800
const HIDE_BELOW_VIEWPORT = 1024

function loadWidth() {
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored) {
      const parsed = parseInt(stored, 10)
      if (!isNaN(parsed)) return Math.max(MIN_CHAT_WIDTH, Math.min(MAX_CHAT_WIDTH, parsed))
    }
  } catch {
    // localStorage unavailable
  }
  return DEFAULT_CHAT_WIDTH
}

function saveWidth(value) {
  try {
    localStorage.setItem(STORAGE_KEY, String(value))
  } catch {
    // localStorage unavailable
  }
}

export function useEditorSplit() {
  const chatWidth = ref(loadWidth())
  const showChat = ref(
    typeof window !== 'undefined' ? window.innerWidth >= HIDE_BELOW_VIEWPORT : true,
  )
  let dragging = false

  function onResize() {
    showChat.value = window.innerWidth >= HIDE_BELOW_VIEWPORT
  }

  onMounted(() => {
    window.addEventListener('resize', onResize)
    onResize()
  })

  function onMouseMove(e) {
    if (!dragging) return
    const newWidth = Math.max(
      MIN_CHAT_WIDTH,
      Math.min(MAX_CHAT_WIDTH, window.innerWidth - e.clientX),
    )
    chatWidth.value = newWidth
  }

  function onMouseUp() {
    if (dragging) saveWidth(chatWidth.value)
    dragging = false
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
    document.body.style.cursor = ''
    document.body.style.userSelect = ''
  }

  function startChatResize() {
    dragging = true
    document.addEventListener('mousemove', onMouseMove)
    document.addEventListener('mouseup', onMouseUp)
    document.body.style.cursor = 'ew-resize'
    document.body.style.userSelect = 'none'
  }

  onBeforeUnmount(() => {
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
    window.removeEventListener('resize', onResize)
    document.body.style.cursor = ''
    document.body.style.userSelect = ''
  })

  return {
    chatWidth,
    showChat,
    startChatResize,
  }
}
