// The Manage & Monitor tab's left rail — sections → items, structurally
// mirroring config/dashboard.js (the live dashboard's left rail). The
// deliberate difference: items carry NO `route`. The dashboard rail routes to
// dedicated child routes, but Manage & Monitor lives entirely inside one
// Wizard tab and switches its content via a `section` query param instead
// (see views/Wizard.vue). Keeping this nav config — not `pipeline.js` STEPS —
// as the single source of the M&M section set is what keeps the locked
// "Manage & Monitor has no steps" contract true: giving these items routes
// would pull in the step-shaped routing/registration STEPS uses, blurring
// the two mechanisms the wizard is built to keep apart.
//
// `labelKey` is an i18n key; `icon` is a PrimeIcons class (no leading `pi`).
// Category items reuse the shared `project.dashboard.categories.*` labels —
// same five categories as the live dashboard rail (config/categories.js
// CATEGORY_ORDER), so the wording stays identical across both surfaces
// without duplicating translation strings.
export const MANAGE_NAV = [
  {
    key: 'monitor',
    labelKey: 'project.manage.nav.monitor',
    items: [
      // Shares components/dashboard/OverviewPanel.vue with the live dashboard's
      // Overview page, revision-scoped via
      // components/wizard/manage/ManageOverview.vue. Leads the group, mirroring
      // the dashboard rail's own "Dashboard" entry (see config/dashboard.js).
      { key: 'overview', labelKey: 'project.manage.views.overview', icon: 'pi-gauge' },
      { key: 'board', labelKey: 'project.manage.views.board', icon: 'pi-objects-column' },
      { key: 'activity', labelKey: 'project.manage.views.activity', icon: 'pi-history' },
    ],
  },
  {
    key: 'categories',
    labelKey: 'project.manage.nav.categories',
    items: [
      { key: 'incident', labelKey: 'project.dashboard.categories.incident', category: 'incident' },
      { key: 'feature', labelKey: 'project.dashboard.categories.feature', category: 'feature' },
      { key: 'privacy', labelKey: 'project.dashboard.categories.privacy', category: 'privacy' },
      { key: 'task', labelKey: 'project.dashboard.categories.task', category: 'task' },
      { key: 'review', labelKey: 'project.dashboard.categories.review', category: 'review' },
    ],
  },
]
