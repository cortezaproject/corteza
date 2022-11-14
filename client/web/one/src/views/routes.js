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

export default [
  r('layout', '/', 'Layout'),

  {
    ...r('bridge', '/bridge', 'Bridge/index'),
    children: [
      r('bridge-jitsi', 'jitsi', 'Bridge/Jitsi'),
    ],
  },

  // When everything else fails, go to root
  { path: '*', redirect: { name: 'layout' } },
]
