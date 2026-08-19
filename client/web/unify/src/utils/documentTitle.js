// The topbar heading is the app's one statement of where you are, so the tab
// title is that text: every view teleports its heading into `#topbar-title`
// (an id lib/vue's navigation contract calls API), and the shell mirrors it.
// Chrome that sits beside the heading — a revision switcher, a rename button —
// carries `data-title-exclude` and stays out of the title.
import { onBeforeUnmount, onMounted, toValue, watch } from 'vue'

export const TITLE_EXCLUDE_ATTR = 'data-title-exclude'

// Matches index.html's static <title>, so a route with no heading reads the
// same before and after boot.
export const BASE_TITLE = import.meta.env.VITE_APP_TITLE || 'Human'

function collect(node) {
  let out = ''

  for (const child of node.childNodes) {
    if (child.nodeType === 3) {
      out += child.nodeValue
    } else if (child.nodeType === 1 && !child.hasAttribute(TITLE_EXCLUDE_ATTR)) {
      // Padded: sibling elements carry no whitespace between them once Vue has
      // condensed the template, and "SalesPipeline" is not a title.
      out += ` ${collect(child)} `
    }
  }

  return out
}

// Rendered text of `node`, minus every subtree marked `data-title-exclude`.
export function readTitleText(node) {
  return node ? collect(node).replace(/\s+/g, ' ').trim() : ''
}

export function composeTitle(heading, unreadCount = 0) {
  const title = heading || BASE_TITLE
  return unreadCount > 0 ? `(${unreadCount}) ${title}` : title
}

// Keeps `document.title` on the heading currently teleported into `#topbar-title`.
export function useDocumentTitle(unreadCount, targetId = 'topbar-title') {
  let target = null
  let observer = null

  const apply = () => {
    document.title = composeTitle(readTitleText(target), toValue(unreadCount))
  }

  onMounted(() => {
    // The target is in the shell's own template, so it exists by now; the
    // headings arrive later (Teleport `defer`) and come in as mutations.
    target = document.getElementById(targetId)
    apply()

    if (!target) return

    observer = new MutationObserver(apply)
    observer.observe(target, { childList: true, subtree: true, characterData: true })
  })

  onBeforeUnmount(() => {
    observer?.disconnect()
    observer = null
  })

  watch(() => toValue(unreadCount), apply)
}
