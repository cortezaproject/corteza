import { afterEach, describe, expect, it, vi } from 'vitest'
import { ref } from 'vue'
import { brandLogoUrl, useBrandLogo, useOsColorScheme } from './useBrandLogo'

function settingsOf(ui: Record<string, string | undefined>) {
  return {
    attachment: (key: string) => {
      const v = ui[key.replace(/^ui\./, '')]
      return v ? `http://x${v}` : undefined
    },
  }
}

describe('brandLogoUrl', () => {
  it('light reads the light setting only', () => {
    const s = settingsOf({ mainLogo: '/light.svg', mainLogoDark: '/dark.svg' })
    expect(brandLogoUrl(s, 'main', 'light')).toBe('http://x/light.svg')
  })

  it('dark prefers the dark setting', () => {
    const s = settingsOf({ mainLogo: '/light.svg', mainLogoDark: '/dark.svg' })
    expect(brandLogoUrl(s, 'main', 'dark')).toBe('http://x/dark.svg')
  })

  it('dark falls back to the light setting when no dark one is set', () => {
    const s = settingsOf({ iconLogo: '/icon.svg' })
    expect(brandLogoUrl(s, 'icon', 'dark')).toBe('http://x/icon.svg')
  })
})

describe('useBrandLogo', () => {
  it('re-resolves when the scheme flips', () => {
    const scheme = ref<'light' | 'dark'>('light')
    const settings = settingsOf({ mainLogo: '/light.svg', mainLogoDark: '/dark.svg' })
    const url = useBrandLogo('main', { scheme, settings })

    expect(url.value).toBe('http://x/light.svg')
    scheme.value = 'dark'
    expect(url.value).toBe('http://x/dark.svg')
  })
})

describe('useOsColorScheme', () => {
  const original = window.matchMedia

  afterEach(() => {
    window.matchMedia = original
  })

  it('tracks prefers-color-scheme', () => {
    let onChange: (() => void) | undefined
    const query = {
      matches: true,
      addEventListener: vi.fn((_: string, fn: () => void) => {
        onChange = fn
      }),
    }
    window.matchMedia = vi.fn(() => query) as any

    // module-level singleton: only the first call reads matchMedia
    const scheme = useOsColorScheme()
    expect(scheme.value).toBe('dark')

    query.matches = false
    onChange?.()
    expect(scheme.value).toBe('light')
  })
})
