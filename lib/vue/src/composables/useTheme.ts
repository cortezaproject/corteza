import { definePreset, palette, usePreset } from '@primeuix/themes'
import Aura from '@primeuix/themes/aura'

type Theme = 'light' | 'dark'

interface SurfacePalette {
  0: string
  50: string
  100: string
  200: string
  300: string
  400: string
  500: string
  600: string
  700: string
  800: string
  900: string
  950: string
}
interface CortezaThemeVariables {
  primary: string
  success: string
  warning: string
  danger: string
  surface: SurfacePalette
  'body-bg': string
  'topbar-bg': string
  'sidebar-bg': string
}

const defaultVariables: {
  light: CortezaThemeVariables
  dark: CortezaThemeVariables
} = {
  light: {
    primary: '#09344E',
    success: '#43AA8B',
    warning: '#E27646',
    danger: '#E54122',
    surface: {
      0: '#ffffff',
      50: '#fafafa',
      100: '#f4f4f5',
      200: '#e4e4e7',
      300: '#d4d4d8',
      400: '#a1a1aa',
      500: '#71717a',
      600: '#52525b',
      700: '#3f3f46',
      800: '#27272a',
      900: '#18181b',
      950: '#09090b',
    },
    'body-bg': '#f4f4f5',
    'sidebar-bg': '#ffffff',
    'topbar-bg': '#f4f4f5',
  },
  dark: {
    primary: '#E56B5B',
    success: '#43AA8B',
    warning: '#E27646',
    danger: '#E54122',
    surface: {
      0: '#ffffff',
      50: '#fafafa',
      100: '#f4f4f5',
      200: '#e4e4e7',
      300: '#d4d4d8',
      400: '#a1a1aa',
      500: '#71717a',
      600: '#52525b',
      700: '#3f3f46',
      800: '#27272a',
      900: '#18181b',
      950: '#09090b',
    },
    'body-bg': '#27272a',
    'sidebar-bg': '#18181b',
    'topbar-bg': '#27272a',
  },
}

let themes: Record<string, any> = {}

export function setThemes(tt: Record<string, any>) {
  themes = tt
}

export function getTheme(theme: Theme) {
  const variables = getThemeVariables(theme)

  document.documentElement.classList.toggle('dark', theme === 'dark')

  return definePreset(Aura, {
    primitive: {
      green: palette(variables['success']),
      red: palette(variables['danger']),
      orange: palette(variables['warning']),
    },
    semantic: {
      primary: palette(variables['primary']),
      colorScheme: {
        light: {
          surface: variables.surface,
          primary: {
            contrastColor: '{surface.50}',
          },
        },
        dark: {
          surface: variables.surface,
          primary: {
            contrastColor: '{surface.50}',
          },
        },
      },
    },
    components: {},
    css: () => `
      :root {
        --topbar-height: 3.5rem;
        --topbar-bg: ${variables['topbar-bg']};
        --sidebar-width: 20rem;
        --sidebar-bg: ${variables['sidebar-bg']};
        --right-sidebar-width: 360px;
        --right-sidebar-bottom: calc(58px + 0.75rem);
        --body-bg: ${variables['body-bg']};
        --p-drawer-border-color: #ffffff00;
        --p-overlay-modal-padding: 1rem;
        --p-tabs-tabpanel-padding: 1rem;
        --p-tabs-tab-padding: 0.75rem 1rem !important;
      }

      body {
        background-color: var(--body-bg);
      }

      .body-bg {
        background-color: var(--body-bg);
      }

      .topbar-bg {
        background-color: var(--topbar-bg);
      }

      .sidebar-bg {
        background-color: var(--sidebar-bg);
      }

      .bg-surface {
        background-color: var(--p-content-background);
      }

      .p-panel {
        border-radius: var(--p-card-border-radius);
        .p-panel-header {
          padding: 0.5rem 0.5rem 0.5rem 1rem !important;

          .p-panel-title {
            font-size: 1.25rem;
            font-weight: 500;
          }
        }
        .p-panel-content {
          padding: 0.5rem 1rem 1rem 1rem !important;
        }
      }

      .p-dialog-header {
        padding-bottom: 0.5rem !important;
      }

      .p-datatable-column-resizer {
        width: 1rem !important;
        right: -0.5rem !important;
      }

      .p-datatable-column-resizer:hover {
        background-color: var(--p-highlight-focus-background) !important;
      }

      .right-sidebar {
        position: fixed;
        top: calc(var(--topbar-height) + 0.75rem);
        right: 0.75rem;
        bottom: var(--right-sidebar-bottom);
        width: var(--right-sidebar-width);
        max-width: 100%;
        z-index: 50;
        background-color: var(--p-content-background);
        border: 1px solid var(--p-content-border-color);
        border-radius: var(--p-border-radius-xl);
        box-shadow: var(--p-overlay-popover-shadow);
        overflow: hidden;
      }
    `,
  })
}

export function useTheme(theme: Theme) {
  usePreset(getTheme(theme))
}

export function getThemeVariables(theme: Theme) {
  const saved = getSavedThemeVariables(theme)
  return {
    ...defaultVariables[theme],
    ...saved,
  }
}

function getSavedThemeVariables(theme: Theme): Partial<CortezaThemeVariables> {
  if (!Array.isArray(themes)) return {}
  const entry = themes.find((t: any) => t.id === theme)
  if (!entry?.values) return {}
  try {
    const parsed = JSON.parse(entry.values)
    // Ensure hex values have # prefix
    const result: Record<string, string> = {}
    for (const [k, v] of Object.entries(parsed)) {
      if (typeof v === 'string') {
        result[k] = v.startsWith('#') ? v : `#${v}`
      }
    }
    return result
  } catch {
    return {}
  }
}
