// Raw compose route tree. The section index wraps these children under a
// /compose layout route (ComposeHost) and prefixes absolute paths.
export const composeRoutes = [
    {
      path: '/',
      name: 'root',
      redirect: '/namespaces',
      children: [
        {
          path: '/namespaces',
          name: 'namespace.list',
          component: () => import('./views/Namespace/List.vue'),
          // Namespace chooser screens keep the sidebar collapsed.
          meta: { hideSidebar: true },
        },
        {
          path: '/namespaces/manage',
          redirect: { name: 'namespace.list' },
        },
        {
          path: '/namespaces/create',
          name: 'namespace.create',
          component: () => import('./views/Namespace/Edit.vue'),
          meta: { hideSidebar: true },
        },
        {
          path: '/namespaces/edit/:slug',
          name: 'namespace.edit',
          component: () => import('./views/Namespace/Edit.vue'),
          meta: { hideSidebar: true },
        },
        {
          path: '/namespace/:slug',
          name: 'namespace.view',
          component: () => import('./views/Namespace/View.vue'),
          props: true,
          children: [
            // Public pages
            {
              path: '',
              name: 'pages',
              component: () => import('./views/Pages/Index.vue'),
            },
            {
              path: 'pages/:pageID',
              name: 'page',
              component: () => import('./views/Pages/View.vue'),
            },

            // Admin - Modules
            {
              path: 'admin/modules',
              name: 'admin.modules',
              component: () => import('./views/Admin/Modules/List.vue'),
            },
            {
              path: 'admin/modules/create',
              name: 'admin.modules.create',
              component: () => import('./views/Admin/Modules/Edit.vue'),
            },
            {
              path: 'admin/modules/:moduleID/edit',
              name: 'admin.modules.edit',
              component: () => import('./views/Admin/Modules/Edit.vue'),
            },

            // Admin - Pages
            {
              path: 'admin/pages',
              name: 'admin.pages',
              component: () => import('./views/Admin/Pages/List.vue'),
            },
            {
              path: 'admin/pages/create',
              name: 'admin.pages.create',
              component: () => import('./views/Admin/Pages/Edit.vue'),
            },
            {
              path: 'admin/pages/:pageID/edit',
              name: 'admin.pages.edit',
              component: () => import('./views/Admin/Pages/Edit.vue'),
            },
            {
              path: 'admin/pages/:pageID/builder',
              name: 'admin.pages.builder',
              component: () => import('./views/Admin/Pages/Builder.vue'),
            },

            // Admin - Charts
            {
              path: 'admin/charts',
              name: 'admin.charts',
              component: () => import('./views/Admin/Charts/List.vue'),
            },
            {
              path: 'admin/charts/create',
              name: 'admin.charts.create',
              component: () => import('./views/Admin/Charts/Edit.vue'),
            },
            {
              path: 'admin/charts/:chartID/edit',
              name: 'admin.charts.edit',
              component: () => import('./views/Admin/Charts/Edit.vue'),
            },

            // Admin - Module Records
            {
              path: 'admin/modules/:moduleID/records',
              name: 'admin.modules.record.list',
              component: () => import('./views/Admin/Modules/Records/List.vue'),
            },
            {
              path: 'admin/modules/:moduleID/records/create',
              name: 'admin.modules.record.create',
              component: () => import('./views/Admin/Modules/Records/Create.vue'),
            },
            {
              path: 'admin/modules/:moduleID/records/:recordID',
              name: 'admin.modules.record.view',
              component: () => import('./views/Admin/Modules/Records/Edit.vue'),
            },
            {
              path: 'admin/modules/:moduleID/records/:recordID/edit',
              name: 'admin.modules.record.edit',
              component: () => import('./views/Admin/Modules/Records/Edit.vue'),
            },

            // Public - Record View
            {
              path: 'pages/:pageID/records/:recordID',
              name: 'page.record',
              component: () => import('./views/Pages/RecordView.vue'),
            },
          ],
        },
      ],
    },
]
