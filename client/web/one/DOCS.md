# One (Launcher/Dashboard)

The main entry point for Corteza. Displays a grid of available applications for users to launch.

## Routes

| Path | Name | View | Description |
|------|------|------|-------------|
| `/` | — | `AppList.vue` | Application grid |
| `/:pathMatch(.*)*` | — | redirect to `/` | Catch-all |

## Stores

### useApplicationsStore (`stores/applications.js`)

- **State:** `apps[]`, `loading`, `error`
- **Getters:** `unifyOnly` — apps with `unify.listed === true`
- **Actions:** `fetchApplications()` — loads enabled apps from SystemAPI

## Views

### AppList.vue

Responsive grid of application cards (320px wide). Features:
- `CInputSearch` for real-time filtering by app name
- Cards with logo, name, hover animation (scale + shadow)
- Filters by `unify.listed` and search query
- Falls back to `applications/default-app.png` for apps without logos

## Components

No app-specific components. Uses only shared components from `lib/vue`:
- `CTopbar` (with `hide-app-selector`)
- `CLoaderLogo`
- `CInputSearch`
