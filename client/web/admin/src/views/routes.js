/**
 * Simple route generator
 *
 * @param name {String}
 * @param path {String}
 * @param component {String}
 * @returns {Object}
 */
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
    component: { name: name + 'Wrap', template: `<router-view />` },
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
      r(`${ns}.${name}.list`, `list`, `${opt.cmpDir}/List`),
      r(`${ns}.${name}.new`, `new`, `${opt.cmpDir}/Editor`),
      r(`${ns}.${name}.edit`, `edit/:${opt.pkey}`, `${opt.cmpDir}/Editor`),
    ],
  }
}

export default [
  {
    name: 'root',
    path: '/',
    component: () => import('./Layout.vue'),
    redirect: 'dashboard',
    children: [
      r('dashboard', 'dashboard', 'Dashboard'),
      {
        ...wrap(`system`, `/system`),

        children: [
          r('system.stats', 'stats', 'System/Stats'),

          combo('system', 'user'),
          combo('system', 'role'),
          combo('system', 'application'),
          combo('system', 'script'),
          combo('system', 'template'),

          r('system.settings', 'settings', 'System/Settings/Index'),
          r('system.email', 'email', 'System/Email/Index'),

          combo('system', 'authClient', { pkey: 'authClientID' }),

          combo('system', 'apigw', { pkey: 'routeID', plural: 'routes' }),

          r('system.apigw.profiler', 'apigw/profiler', 'System/Apigw/Profiler/Index'),
          r('system.apigw.profiler.route.list', 'apigw/profiler/route/:routeID', 'System/Apigw/Profiler/Route'),
          r('system.apigw.profiler.hit.list', 'apigw/profiler/hit/:hitID', 'System/Apigw/Profiler/Hit'),

          r('system.permissions', 'permissions', 'System/Permissions/Index'),
          r('system.actionlog', 'actionlog', 'System/Actionlog/Index'),

          combo('system', 'queue', { pkey: 'queueID' }),
        ],
      },

      {
        ...wrap(`compose`, `/compose`),
        children: [
          r('compose.settings', 'settings', 'Compose/Settings/Index'),
          r('compose.permissions', 'permissions', 'Compose/Permissions/Index'),
        ],
      },

      {
        ...wrap(`automation`, `/automation`),
        children: [
          combo('automation', 'workflow'),
          r('automation.scripts', 'scripts', 'Automation/Scripts/Index'),
          combo('automation', 'session'),
          r('automation.permissions', 'permissions', 'Automation/Permissions/Index'),
        ],
      },

      {
        ...wrap(`federation`, `/federation`),
        children: [
          combo('federation', 'nodes', { pkey: 'nodeID' }),
          r('federation.permissions', 'permissions', 'Federation/Permissions/Index'),
        ],
      },

      {
        ...r('ui.settings', 'ui', 'UI/Index'),
      },
    ],
  },

  // When everything else fails, go to dashboard
  { path: '*', redirect: { name: 'dashboard' } },
]
