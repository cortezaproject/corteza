export default [
  {
    path: '',
    name: 'root',
    redirect: { name: 'workflow.list' },
    component: () => import('./Layout.vue'),
    children: [
      { name: 'workflow.list', path: '/list', component: () => import('./Workflow/Index.vue') },
      { name: 'workflow.create', path: 'new', component: () => import('./Workflow/Editor.vue') },
      { name: 'workflow.edit', path: ':workflowID/edit', component: () => import('./Workflow/Editor.vue') },
    ],
  },

  { path: '*', redirect: { name: 'root' } },
]
