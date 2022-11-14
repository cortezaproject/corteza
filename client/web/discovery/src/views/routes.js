export default [
  {
    name: 'root',
    path: '',
    component: () => import('./Layout.vue'),
  },

  // When everything else fails, go to root
  { path: '*', redirect: { name: 'root' } },
]
