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
      path: '/project/projects/:projectId',
      name: 'project.overview',
      component: () => import('./views/ProjectOverview.vue'),
      meta: { section: 'project' },
    },
  ],
  sidebar: ProjectSidebar,
  sidebarDisabledRoutes: [],
  autoExpandSidebar: true,
}
