import { ref } from 'vue'

/**
 * Position-on-hover popup for vue-flow nodes. Computes fixed coords from the
 * node's bounding rect so the popup can render via Teleport outside the
 * transformed wrapper that traps its stacking context.
 */
export function useNodePreview () {
  const rootEl = ref(null)
  const popupStyle = ref(null)

  function showPopup () {
    const r = rootEl.value?.getBoundingClientRect()
    if (!r) return
    popupStyle.value = {
      position: 'fixed',
      top: `${r.bottom + 14}px`,
      left: `${r.left + r.width / 2}px`,
      transform: 'translateX(-50%)',
    }
  }

  function hidePopup () {
    popupStyle.value = null
  }

  return { rootEl, popupStyle, showPopup, hidePopup }
}
