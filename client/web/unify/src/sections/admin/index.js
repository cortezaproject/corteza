// Admin section — mounted under /admin. The 66 legacy routes keep their
// (already well-namespaced) names; we only prefix paths with /admin, rename
// the legacy `root` redirect to the section index `admin`, and tag each route
// with meta.section so the shell resolves the right sidebar.
import { adminRoutes } from './routes'
import AdminSidebar from './sidebar/AdminSidebar.vue'

const PREFIX = '/admin'

export default {
  id: 'admin',
  routes: adminRoutes.map(route => ({
    ...route,
    path: route.path === '/' ? PREFIX : PREFIX + route.path,
    name: route.name === 'root' ? 'admin' : route.name,
    meta: { ...(route.meta || {}), section: 'admin' },
  })),
  sidebar: AdminSidebar,
  // Admin shows its sidebar on every route.
  sidebarDisabledRoutes: [],
  // Open the sidebar automatically when navigating into the admin section.
  autoExpandSidebar: true,
}
