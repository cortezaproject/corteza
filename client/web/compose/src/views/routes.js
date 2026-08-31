// public route builder/helper
import { components } from '@cortezaproject/corteza-vue'

const C311AccessPage = components.C311AccessPage || {
  functional: true,
  render: (h, { props }) => h('main', { attrs: { tabindex: '-1', 'data-c311-main': '' }, class: 'p-4' }, [
    h('h1', { class: 'h2' }, props.heading || `Access status ${props.status}`),
    h('p', [props.message || 'This page is not available.']),
  ]),
}

function r (name, path, component, defaultProps = {}) {
  return {
    path,
    name,
    component: () => import('./' + component + '.vue'),
    props: r => {
      return { ...defaultProps, ...r.params }
    },
  }
}

export default [
  {
    name: 'c311.unauthorized',
    path: '/c311/401',
    component: C311AccessPage,
    props: { status: 401 },
  },
  {
    name: 'c311.forbidden',
    path: '/c311/403',
    component: C311AccessPage,
    props: { status: 403 },
  },
  {
    name: 'c311.not-found',
    path: '/c311/404',
    component: C311AccessPage,
    props: { status: 404 },
  },
  {
    name: 'c311.submit',
    path: '/c311/submit',
    component: () => import('./C311/Portal.vue'),
    meta: { c311: { public: true } },
  },
  {
    name: 'c311.staff.submit',
    path: '/c311/staff/submit',
    component: () => import('./C311/Portal.vue'),
    meta: { c311: { requiresAuth: true, route: 'staff_service_request_create', capabilities: ['staff_service_request_create'], scopes: ['service_requests.write'] } },
  },
  {
    name: 'c311.status',
    path: '/c311/status',
    component: () => import('./C311/Portal.vue'),
    meta: { c311: { public: true } },
  },
  {
    name: 'c311.portal',
    path: '/c311',
    component: () => import('./C311/PublicPortal.vue'),
    meta: { c311: { public: true } },
  },
  {
    name: 'c311.services',
    path: '/c311/services',
    component: () => import('./C311/PublicPortal.vue'),
    meta: { c311: { public: true, route: 'public_content_get' } },
  },
  {
    name: 'c311.help',
    path: '/c311/help',
    component: () => import('./C311/PublicPortal.vue'),
    meta: { c311: { public: true, route: 'public_help_get' } },
  },
  {
    name: 'c311.sign-in',
    path: '/c311/sign-in',
    component: () => import('./C311/PublicPortal.vue'),
    meta: { c311: { public: true, route: 'session_sign_in' } },
  },
  {
    name: 'c311.register',
    path: '/c311/register',
    component: () => import('./C311/PublicPortal.vue'),
    meta: { c311: { public: true, route: 'account_register' } },
  },
  {
    name: 'c311.forgot-password',
    path: '/c311/forgot-password',
    component: () => import('./C311/PublicPortal.vue'),
    meta: { c311: { public: true, route: 'password_reset_request' } },
  },
  {
    name: 'c311.reset-password',
    path: '/c311/reset-password',
    component: () => import('./C311/PublicPortal.vue'),
    meta: { c311: { public: true, route: 'password_reset_confirm' } },
  },
  {
    name: 'c311.account',
    path: '/c311/account',
    component: () => import('./C311/PublicPortal.vue'),
    meta: { c311: { public: true, requiresAuth: true, route: 'profile_get', capabilities: ['profile_get'] } },
  },
  {
    name: 'c311.requests',
    path: '/c311/requests',
    component: () => import('./C311/PublicPortal.vue'),
    meta: { c311: { public: true, requiresAuth: true, route: 'portal_my_requests', capabilities: ['portal_my_requests'], scopes: ['service_requests.write'] } },
  },
  {
    name: 'c311.auth.callback',
    path: '/c311/auth/callback',
    component: () => import('./C311/PublicPortal.vue'),
    meta: { c311: { public: true, route: 'federated_sign_in_callback' } },
  },
  {
    name: 'c311.auth.link.confirm',
    path: '/c311/auth/link/confirm',
    component: () => import('./C311/PublicPortal.vue'),
    meta: { c311: { public: true } },
  },
  {
    name: 'c311.logout.callback',
    path: '/c311/logout/callback',
    component: () => import('./C311/PublicPortal.vue'),
    meta: { c311: { public: true, route: 'session_sign_out' } },
  },
  {
    name: 'c311.test.interaction',
    path: '/c311/test/modal',
    component: components.C311InteractionHarness,
    meta: { c311: { public: true } },
    beforeEnter: (_to, _from, next) => {
      if (typeof window !== 'undefined' && window.C311Mode === 'mock') return next()
      return next({ name: 'c311.not-found' })
    },
  },
  {
    name: 'c311.not-found-wildcard',
    path: '/c311/*',
    component: C311AccessPage,
    props: { status: 404 },
    meta: { c311: { public: true } },
  },
  {
    name: 'root',
    path: '',
    component: () => import('./Layout.vue'),
    redirect: 'namespaces',
    children: [
      {
        name: 'namespaces',
        path: '/',
        component: () => import('./Namespace/Index.vue'),
        redirect: 'namespaces.list',
        children: [
          r('namespace.list', '/namespaces', 'Namespace/List'),
          r('namespace.manage', '/namespaces/manage', 'Namespace/Manage'),
          r('namespace.create', '/admin/namespace/create', 'Namespace/Edit'),
          r('namespace.edit', '/admin/namespace/edit/:namespaceID', 'Namespace/Edit'),
          {
            ...r('namespace', '/ns/:slug', 'Namespace/View'),
            redirect: { name: 'pages' },

            children: [
              {
                ...r('pages', 'pages', 'Public/Index'),
                children: [
                  {
                    ...r('page', ':pageID?', 'Public/Pages/View'),

                    children: [
                      r('page.record', 'record/:recordID', 'Public/Pages/Records/View', { edit: false }),
                      r('page.record.edit', 'record/:recordID/edit', 'Public/Pages/Records/View', { edit: true }),
                      r('page.record.create', 'record', 'Public/Pages/Records/View', { edit: true }),
                    ],
                  },
                ],
              },
              {
                ...r('admin', 'admin', 'Admin/Index'),
                redirect: { name: 'admin.modules' },

                children: [
                  r('admin.modules', 'modules', 'Admin/Modules/List'),
                  r('admin.modules.create', 'modules/new', 'Admin/Modules/Edit'),
                  r('admin.modules.edit', 'modules/:moduleID/edit', 'Admin/Modules/Edit'),
                  r('admin.modules.record.list', 'modules/:moduleID/record/list', 'Admin/Modules/Records/List'),
                  r('admin.modules.record.view', 'modules/:moduleID/record/:recordID', 'Admin/Modules/Records/View', { edit: false }),
                  r('admin.modules.record.create', 'modules/:moduleID/record', 'Admin/Modules/Records/View', { edit: true }),
                  r('admin.modules.record.edit', 'modules/:moduleID/record/:recordID/edit', 'Admin/Modules/Records/View', { edit: true }),

                  r('admin.pages', 'pages', 'Admin/Pages/List'),
                  r('admin.pages.edit', 'pages/:pageID/edit', 'Admin/Pages/Edit'),
                  r('admin.pages.builder', 'pages/:pageID/builder', 'Admin/Pages/Builder'),

                  r('admin.charts', 'charts', 'Admin/Charts/List'),
                  r('admin.charts.create', 'charts/new/:category?', 'Admin/Charts/Edit'),
                  r('admin.charts.edit', 'charts/:chartID/edit', 'Admin/Charts/Edit'),

                ],
              },

              { path: '*', redirect: { name: 'pages' } },
            ],
          },
        ],
      },
    ],
  },

  // When everything else fails, go to namespaces
  { path: '*', redirect: { name: 'root' } },
]
