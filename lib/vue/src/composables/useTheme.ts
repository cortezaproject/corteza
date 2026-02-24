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
    primary: '#FF9661',
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
    'topbar-bg': '#fafafa',
    'sidebar-bg': '#f4f4f5',
  },
  dark: {
    primary: '#FF9661',
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
    'topbar-bg': '#18181b',
    'sidebar-bg': '#18181b',
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
        --body-bg: ${variables['body-bg']};
        --p-drawer-border-color: #ffffff00;
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

      .p-datatable-column-resizer {
        width: 1rem !important;
        right: -0.5rem !important;
      }

      .p-datatable-column-resizer:hover {
        background-color: var(--p-highlight-focus-background) !important;
      }
    `,
  })
}

export function useTheme(theme: Theme) {
  usePreset(getTheme(theme))
}

export function getThemeVariables(theme: Theme) {
  return {
    ...defaultVariables[theme],
  }
}
