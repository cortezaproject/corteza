// The published-project dashboard sub-navigation — single source of truth for
// the in-view left rail (DashboardNav), mirroring the reference demo's sidebar.
// Structure only: sections → items. Each item routes to a child view of the
// dashboard layout.
//
// `labelKey` is an i18n key; `icon` is a PrimeIcons class (no leading `pi`).
// `route` is the target route name. `badge`/`badgeTone` are the count pill and its tone.
export const DASHBOARD_NAV = [
  {
    key: 'monitor',
    labelKey: 'project.dashboard.nav.monitor',
    items: [
      {
        key: 'overview',
        labelKey: 'project.dashboard.views.overview',
        icon: 'pi-gauge',
        route: 'project.overview',
      },
      {
        key: 'activity',
        labelKey: 'project.dashboard.views.activity',
        icon: 'pi-history',
        route: 'project.overview.activity',
      },
    ],
  },
  {
    key: 'categories',
    labelKey: 'project.dashboard.nav.categories',
    items: [
      {
        key: 'incidents',
        labelKey: 'project.dashboard.categories.incidents',
        route: 'project.overview.category',
        category: 'incident',
      },
      {
        key: 'features',
        labelKey: 'project.dashboard.categories.features',
        route: 'project.overview.category',
        category: 'feature',
      },
      {
        key: 'privacy',
        labelKey: 'project.dashboard.categories.privacy',
        route: 'project.overview.category',
        category: 'privacy',
      },
      {
        key: 'tasks',
        labelKey: 'project.dashboard.categories.tasks',
        route: 'project.overview.category',
        category: 'task',
      },
      {
        key: 'reviews',
        labelKey: 'project.dashboard.categories.reviews',
        route: 'project.overview.category',
        category: 'review',
      },
    ],
  },
  {
    key: 'insights',
    labelKey: 'project.dashboard.nav.insights',
    items: [
      {
        key: 'reports',
        labelKey: 'project.dashboard.views.reports',
        icon: 'pi-chart-bar',
        route: 'project.overview.reports',
      },
      {
        key: 'backlog',
        labelKey: 'project.dashboard.views.backlog',
        icon: 'pi-th-large',
        route: 'project.overview.backlog',
      },
    ],
  },
]
