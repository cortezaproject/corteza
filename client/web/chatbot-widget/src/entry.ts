import { mountChatbot } from './mount'

function findSelf(): HTMLScriptElement | null {
  const current = document.currentScript as HTMLScriptElement | null
  if (current && current.getAttribute('data-widget-key')) return current
  const withKey = document.querySelector<HTMLScriptElement>(
    'script[data-widget-key][src*="widget.js"]',
  )
  if (withKey) return withKey
  return document.querySelector<HTMLScriptElement>('script[src*="widget.js"]')
}

async function boot() {
  console.log('[human-chatbot] boot start')
  const script = findSelf()
  if (!script) {
    console.warn('[human-chatbot] no <script> tag found')
    return
  }
  let widgetKey = script.getAttribute('data-widget-key') || ''
  if (!widgetKey && script.src) {
    try {
      widgetKey = new URL(script.src).searchParams.get('k') || ''
    } catch {
      /* ignore */
    }
  }
  if (!widgetKey) {
    console.warn('[human-chatbot] missing data-widget-key (and no ?k= in src)')
    return
  }

  await mountChatbot({ scriptSrc: script.src, widgetKey })
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', () => {
    void boot()
  })
} else {
  void boot()
}
