// TAQ section — mounted under /taq. Paths mirror the legacy app 1:1 (prefixed);
// the generic legacy names (list/builder/builder-edit) are namespaced under
// `taq.*` to avoid collisions with other sections.
import TaqSidebar from './sidebar/TaqSidebar.vue'

export default {
  id: 'taq',
  routes: [
    {
      path: '/taq',
      name: 'taq',
      component: () => import('./views/List.vue'),
      meta: { section: 'taq' },
    },
    {
      path: '/taq/builder',
      name: 'taq.builder',
      component: () => import('./views/Builder.vue'),
      meta: { section: 'taq' },
    },
    {
      path: '/taq/builder/:id',
      name: 'taq.builder-edit',
      component: () => import('./views/Builder.vue'),
      meta: { section: 'taq' },
    },
  ],
  sidebar: TaqSidebar,
  // List view keeps the sidebar collapsed, like the legacy app.
  sidebarDisabledRoutes: ['taq'],
}
