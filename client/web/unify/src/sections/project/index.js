// Project section — mounted under /project. Paths mirror the legacy app 1:1
// (legacy base was /project/ with routes /projects…, so prefixed they become
// /project/projects…). Legacy `root` redirect → section index `project`.
import ProjectSidebar from './sidebar/ProjectSidebar.vue'

export default {
  id: 'project',
  routes: [
    {
      path: '/project',
      name: 'project',
      redirect: { name: 'project.list' },
      meta: { section: 'project' },
    },
    {
      path: '/project/projects',
      name: 'project.list',
      component: () => import('./views/ProjectList.vue'),
      meta: { section: 'project' },
    },
    {
      path: '/project/projects/:projectId/wizard',
      name: 'project.wizard',
      component: () => import('./views/Wizard.vue'),
      meta: { section: 'project' },
    },
    {
      // Published-project dashboard. A layout (in-view left rail + topbar crumb)
      // with one child per view; each is an empty scaffold stub for now.
      path: '/project/projects/:projectId',
      component: () => import('./views/dashboard/DashboardLayout.vue'),
      meta: { section: 'project' },
      children: [
        {
          path: '',
          name: 'project.overview',
          component: () => import('./views/dashboard/DashboardStub.vue'),
          meta: { section: 'project', titleKey: 'project.dashboard.views.dashboard', icon: 'pi-gauge' },
        },
        {
          path: 'events',
          name: 'project.overview.events',
          component: () => import('./views/dashboard/AllEventsView.vue'),
          meta: { section: 'project', titleKey: 'project.dashboard.views.events', icon: 'pi-list' },
        },
        {
          path: 'category/:category',
          name: 'project.overview.category',
          component: () => import('./views/dashboard/CategoryView.vue'),
          meta: { section: 'project', titleKey: 'project.dashboard.views.dashboard', icon: 'pi-list' },
        },
        {
          path: 'reports',
          name: 'project.overview.reports',
          component: () => import('./views/dashboard/DashboardStub.vue'),
          meta: { section: 'project', titleKey: 'project.dashboard.views.reports', icon: 'pi-chart-bar' },
        },
        {
          path: 'backlog',
          name: 'project.overview.backlog',
          component: () => import('./views/dashboard/DashboardStub.vue'),
          meta: { section: 'project', titleKey: 'project.dashboard.views.backlog', icon: 'pi-th-large' },
        },
      ],
    },
  ],
  sidebar: ProjectSidebar,
  sidebarDisabledRoutes: [],
  autoExpandSidebar: true,
}
