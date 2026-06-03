// Workflow section — mounted under /workflow, paths mirror the legacy app 1:1
// (prefixed). Route names were already namespaced (`workflow.*`); the legacy
// `root` redirect becomes the section index `workflow`.
import WorkflowSidebar from './sidebar/WorkflowSidebar.vue'

export default {
  id: 'workflow',
  routes: [
    {
      path: '/workflow',
      name: 'workflow',
      redirect: { name: 'workflow.list' },
      meta: { section: 'workflow' },
    },
    {
      path: '/workflow/list',
      name: 'workflow.list',
      component: () => import('./views/Home.vue'),
      meta: { section: 'workflow' },
    },
    {
      path: '/workflow/new',
      name: 'workflow.create',
      component: () => import('./views/Editor.vue'),
      meta: { section: 'workflow' },
    },
    {
      path: '/workflow/:workflowID/edit',
      name: 'workflow.edit',
      component: () => import('./views/Editor.vue'),
      meta: { section: 'workflow' },
    },
  ],
  sidebar: WorkflowSidebar,
  // List view keeps the sidebar collapsed, like the legacy app.
  sidebarDisabledRoutes: ['workflow', 'workflow.list'],
}
