/**
 * Simple route generator
 *
 * @param name {String}
 * @param path {String}
 * @param component {String}
 * @returns {Object}
 */
import { components } from '@cortezaproject/corteza-vue'

const C311AccessPage = components.C311AccessPage || {
  functional: true,
  render: (h, { props }) => h('main', { attrs: { tabindex: '-1', 'data-c311-main': '' }, class: 'p-4' }, [
    h('h1', { class: 'h2' }, props.heading || `Access status ${props.status}`),
    h('p', [props.message || 'This page is not available.']),
  ]),
}

function r (name, path, component) {
  return {
    path,
    name,
    component: () => import('./' + component + '.vue'),
    props: true,
    // canReuse: false,
  }
}

/**
 * Wrap route generator
 *
 * Creates a route with simple template that contains only router-view component
 *
 * @param name
 * @param path
 * @param name {String}
 * @param path {String}
 * @param component {String}
 * @returns {Object}
 */
function wrap (name, path) {
  return {
    path,
    name,
    component: { name: name + 'Wrap', template: '<router-view />' },
    props: true,
    // canReuse: false,
  }
}

// Generates 3 routes - list, new-form, edit-form

/**
 * Combo routes generator
 *
 * Creates 4 routes - list, editor for new and existing / wrapper
 *
 * @param ns {String} namespace
 * @param name {String}
 * @param opt {Object}
 * @returns {Object}
 */
function combo (ns, name, opt = {}) {
  const cptlz = (s) => s.slice(0, 1).toUpperCase() + s.slice(1)

  opt = {
    pkey: `${name}ID`,
    plural: `${name}s`,
    cmpDir: cptlz(ns) + '/' + cptlz(name),
    ...opt,
  }

  return {
    ...wrap(`${ns}.${name}`, `/${ns}/${name}`),
    redirect: `/${ns}/${name}/list`,
    children: [
      r(`${ns}.${name}.list`, 'list', `${opt.cmpDir}/List`),
      r(`${ns}.${name}.new`, 'new', `${opt.cmpDir}/Editor`),
      r(`${ns}.${name}.edit`, `edit/:${opt.pkey}`, `${opt.cmpDir}/Editor`),
    ],
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
    name: 'c311.staff.reports',
    path: '/c311/staff/reports',
    component: () => import('./C311/Staff.vue'),
    meta: { c311: { requiresAuth: true, route: 'report_catalogue', capabilities: ['report_catalogue'] } },
  },
  {
    name: 'c311.staff.workflows',
    path: '/c311/staff/workflows',
    component: () => import('./C311/Staff.vue'),
    meta: { c311: { requiresAuth: true, route: 'workflow_list', capabilities: ['workflow_list'], scopes: ['workflow.execute'] } },
  },
  {
    name: 'c311.staff',
    path: '/c311/staff',
    component: () => import('./C311/Staff.vue'),
    meta: { c311: { requiresAuth: true, route: 'staff_request_queue', capabilities: ['staff_request_queue'], scopes: ['service_requests.write'] } },
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
    path: '/',
    component: () => import('./Layout.vue'),
    redirect: 'dashboard',
    children: [
      r('dashboard', 'dashboard', 'Dashboard'),
      {
        ...wrap('system', '/system'),

        children: [
          combo('system', 'user'),
          combo('system', 'role'),
          combo('system', 'application'),
          combo('system', 'template'),

          r('system.settings', 'settings', 'System/Settings/Index'),
          r('system.email', 'email', 'System/Email/Index'),

          combo('system', 'authClient', { pkey: 'authClientID' }),
          combo('system', 'userGroup', { pkey: 'userGroupID' }),

          r('system.apigw', 'apigw', 'System/Apigw/Index'),
          r('system.apigw.new', 'apigw/new', 'System/Apigw/Editor'),
          r('system.apigw.edit', 'apigw/edit/:routeID', 'System/Apigw/Editor'),
          r('system.apigw.profiler', 'apigw/profiler', 'System/Apigw/Profiler/Index'),
          r('system.apigw.profiler.route.list', 'apigw/profiler/route/:routeID', 'System/Apigw/Profiler/Route'),
          r('system.apigw.profiler.hit', 'apigw/profiler/hit/:hitID', 'System/Apigw/Profiler/Hit'),

          r('system.permissions', 'permissions', 'System/Permissions/Index'),
          r('system.actionlog', 'actionlog', 'System/Actionlog/Index'),

          r('system.connection', 'connection', 'System/Connection/Index'),
          r('system.connection.new', 'connection/new', 'System/Connection/Editor'),
          r('system.connection.edit', 'connection/edit/:connectionID', 'System/Connection/Editor'),

          r('system.codesnippets', 'codesnippets', 'System/CodeSnippets/Index'),

          combo('system', 'sensitivityLevel'),

          combo('system', 'queue', { pkey: 'queueID' }),
        ],
      },

      {
        ...wrap('compose', '/compose'),
        children: [
          r('compose.settings', 'settings', 'Compose/Settings/Index'),
          r('compose.permissions', 'permissions', 'Compose/Permissions/Index'),
        ],
      },

      {
        ...wrap('automation', '/automation'),
        children: [
          combo('automation', 'workflow'),
          r('automation.scripts', 'scripts', 'Automation/Scripts/Index'),
          combo('automation', 'session'),
          r('automation.permissions', 'permissions', 'Automation/Permissions/Index'),
        ],
      },

      {
        ...wrap('federation', '/federation'),
        children: [
          combo('federation', 'nodes', { pkey: 'nodeID' }),
          r('federation.permissions', 'permissions', 'Federation/Permissions/Index'),
        ],
      },

      {
        ...wrap('ui', '/ui'),
        children: [
          r('theming.settings', 'theming', 'UI/Theming/Index'),
          r('navigation.settings', 'navigation', 'UI/Navigation/Index'),
          r('location.settings', 'location', 'UI/Location/Index'),
        ],
      },
    ],
  },

  // When everything else fails, go to dashboard
  { path: '*', redirect: { name: 'dashboard' } },
]
