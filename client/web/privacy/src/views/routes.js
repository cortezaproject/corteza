export default [
  {
    path: '',
    name: 'root',
    redirect: 'dashboard',
    component: () => import('./Layout'),
    children: [
      { name: 'dashboard', path: '/dashboard', component: () => import('./Privacy/Dashboard') },
      { name: 'sensitive-data', path: '/sensitive-data', component: () => import('./Privacy/SensitiveData') },
      { name: 'data-overview', path: '/data-overview', component: () => import('./Privacy/DataOverview/') },
      { name: 'data-overview.application', path: '/data-overview/application', component: () => import('./Privacy/DataOverview/Application') },
      { name: 'request.list', path: '/request/list', component: () => import('./Privacy/Request/List') },
      { name: 'request.view', path: '/request/:requestID?', component: () => import('./Privacy/Request/View'), props: true },
      { name: 'request.create', path: '/request/:kind/new', component: () => import('./Privacy/Request/Create'), props: true },
    ],
  },

  { path: '*', redirect: { name: 'root' } },
]
