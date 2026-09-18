// The tab icon is the configured app icon (`ui.iconLogo`), with a dot over it
// while there is anything unread.
import { renderIcon } from '@planetcrust/human-vue'
import { toValue, watch } from 'vue'

export const BUILT_IN_ICON = `${import.meta.env.BASE_URL}icon.svg`

function dotColor() {
  const color = getComputedStyle(document.documentElement).getPropertyValue('--p-red-500').trim()
  return color || undefined
}

// PNG data URL of `src`, or of the built-in icon when `src` cannot be drawn.
export async function appIconImage(src, { dot = false } = {}) {
  const options = { dot, dotColor: dotColor() }

  try {
    return await renderIcon(src || BUILT_IN_ICON, options)
  } catch (err) {
    console.warn('app icon not drawn, using the built-in one:', err?.message || err)
    return renderIcon(BUILT_IN_ICON, options)
  }
}

function setFavicon(href) {
  let link = document.querySelector('link[rel="icon"]')
  if (!link) {
    link = document.createElement('link')
    link.rel = 'icon'
    document.head.appendChild(link)
  }

  link.type = 'image/png'
  link.href = href
}

export function useFavicon(iconUrl, hasUnread) {
  let latest = 0

  const apply = async () => {
    const run = ++latest
    let href

    try {
      href = await appIconImage(toValue(iconUrl), { dot: !!toValue(hasUnread) })
    } catch (err) {
      console.warn('tab icon not drawn:', err?.message || err)
      return
    }

    // A slower draw started earlier must not overwrite a newer one.
    if (run === latest) setFavicon(href)
  }

  watch([() => toValue(iconUrl), () => !!toValue(hasUnread)], apply, { immediate: true })
}
