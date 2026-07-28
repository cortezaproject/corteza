// Project section — mounted under /project. Paths mirror the legacy app 1:1
// (legacy base was /project/ with routes /projects…, so prefixed they become
// /project/projects…). Legacy `root` redirect → section index `project`.
import ProjectSidebar from './sidebar/ProjectSidebar.vue'
import { useProjectsStore } from './stores/projects'

export default {
  id: 'project',
  // Run by the shell when this section becomes active (app load or navigation
  // into it), so the projects list feeding the sidebar tree is ready regardless
  // of whether the lazily-mounted sidebar drawer has been opened.
  preload: () => useProjectsStore().load(),
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
      // with child views for overview, activity, categories and backlog; reports
      // is a stub until the custom-reports design lands.
      path: '/project/projects/:projectId',
      component: () => import('./views/dashboard/DashboardLayout.vue'),
      meta: { section: 'project' },
      children: [
        {
          path: '',
          name: 'project.overview',
          component: () => import('./views/dashboard/Overview.vue'),
          meta: {
            section: 'project',
            titleKey: 'project.dashboard.views.overview',
            icon: 'pi-gauge',
          },
        },
        {
          path: 'activity',
          name: 'project.overview.activity',
          component: () => import('./views/dashboard/AllEventsView.vue'),
          meta: {
            section: 'project',
            titleKey: 'project.dashboard.views.activity',
            icon: 'pi-list',
          },
        },
        {
          path: 'category/:category',
          name: 'project.overview.category',
          component: () => import('./views/dashboard/CategoryView.vue'),
          meta: {
            section: 'project',
            // Generic: the real heading is per-category and comes from
            // CATEGORY_CONFIG, not from route meta.
            titleKey: 'project.dashboard.nav.categories',
            icon: 'pi-list',
          },
        },
        {
          path: 'reports',
          name: 'project.overview.reports',
          component: () => import('./views/dashboard/DashboardStub.vue'),
          meta: {
            section: 'project',
            titleKey: 'project.dashboard.views.reports',
            icon: 'pi-chart-bar',
          },
        },
        {
          path: 'backlog',
          name: 'project.overview.backlog',
          component: () => import('./views/dashboard/BacklogView.vue'),
          meta: {
            section: 'project',
            titleKey: 'project.dashboard.views.backlog',
            icon: 'pi-th-large',
          },
        },
      ],
    },
  ],
  sidebar: ProjectSidebar,
  sidebarDisabledRoutes: [],
}
