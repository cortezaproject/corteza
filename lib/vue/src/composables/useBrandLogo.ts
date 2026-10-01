import {
  computed,
  inject,
  ref,
  toValue,
  type ComputedRef,
  type MaybeRefOrGetter,
  type Ref,
} from 'vue'
import { currentTheme, type Theme } from './useTheme'

export type BrandLogoKind = 'main' | 'icon'

interface LogoSettings {
  attachment(key: string, fallback?: any): string | undefined
}

/**
 * URL of the configured logo (`main`) or icon (`icon`) for a theme.
 *
 * The dark theme reads the `*Dark` setting and falls back to the light one
 * when it is empty, so a single uploaded logo serves both themes until a
 * dark one exists.
 */
export function brandLogoUrl(
  settings: LogoSettings,
  kind: BrandLogoKind,
  theme: Theme,
): string | undefined {
  const light = `ui.${kind}Logo`

  if (theme === 'dark') {
    return settings.attachment(`${light}Dark`) || settings.attachment(light)
  }

  return settings.attachment(light)
}

/**
 * Reactive logo URL for the app's own chrome: follows the theme the preset
 * was built for unless another scheme is given.
 */
export function useBrandLogo(
  kind: BrandLogoKind,
  {
    scheme = currentTheme,
    settings = inject<LogoSettings>('$Settings') as LogoSettings,
  }: { scheme?: MaybeRefOrGetter<Theme>; settings?: LogoSettings } = {},
): ComputedRef<string | undefined> {
  return computed(() => brandLogoUrl(settings, kind, toValue(scheme)))
}

let osScheme: Ref<Theme> | undefined

/**
 * The OS colour scheme (`prefers-color-scheme`), for surfaces the browser or
 * OS draws rather than the app: the tab icon, desktop notifications.
 */
export function useOsColorScheme(): Ref<Theme> {
  if (osScheme) return osScheme

  osScheme = ref<Theme>('light')

  if (typeof window !== 'undefined' && typeof window.matchMedia === 'function') {
    const query = window.matchMedia('(prefers-color-scheme: dark)')
    const apply = () => {
      osScheme!.value = query.matches ? 'dark' : 'light'
    }

    apply()
    query.addEventListener?.('change', apply)
  }

  return osScheme
}
