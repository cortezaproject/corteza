// The published-project dashboard sub-navigation — single source of truth for
// the in-view left rail (DashboardNav), mirroring the reference demo's sidebar.
// Structure only: sections → items. Each item routes to a child view of the
// dashboard layout. Badge counts are mock placeholders until wired to real data.
//
// `labelKey` is an i18n key; `icon` is a PrimeIcons class (no leading `pi`).
// `route` is the target route name; `query` narrows a shared view (the Category
// items all open the Events view pre-filtered). `badge`/`badgeTone` are the
// mock count pill and its tone (default | muted | warning).
export const DASHBOARD_NAV = [
  {
    key: 'monitor',
    labelKey: 'project.dashboard.nav.monitor',
    items: [
      { key: 'dashboard', labelKey: 'project.dashboard.views.dashboard', icon: 'pi-gauge', route: 'project.overview' },
      { key: 'events', labelKey: 'project.dashboard.views.events', icon: 'pi-list', route: 'project.overview.events', iconTone: 'primary' },
    ],
  },
  {
    key: 'categories',
    labelKey: 'project.dashboard.nav.categories',
    items: [
      { key: 'incidents', labelKey: 'project.dashboard.categories.incidents', icon: 'pi-exclamation-triangle', route: 'project.overview.category', category: 'incident' },
      { key: 'features', labelKey: 'project.dashboard.categories.features', icon: 'pi-sparkles', route: 'project.overview.category', category: 'feature' },
      { key: 'privacy', labelKey: 'project.dashboard.categories.privacy', icon: 'pi-shield', route: 'project.overview.category', category: 'privacy' },
      { key: 'tasks', labelKey: 'project.dashboard.categories.tasks', icon: 'pi-check-square', route: 'project.overview.category', category: 'task' },
      { key: 'reviews', labelKey: 'project.dashboard.categories.reviews', icon: 'pi-eye', route: 'project.overview.category', category: 'review' },
    ],
  },
  {
    key: 'insights',
    labelKey: 'project.dashboard.nav.insights',
    items: [
      { key: 'reports', labelKey: 'project.dashboard.views.reports', icon: 'pi-chart-bar', route: 'project.overview.reports' },
      { key: 'backlog', labelKey: 'project.dashboard.views.backlog', icon: 'pi-th-large', route: 'project.overview.backlog', badge: 18, badgeTone: 'muted' },
    ],
  },
]
