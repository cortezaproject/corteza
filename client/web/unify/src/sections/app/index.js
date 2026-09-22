// App section — a custom application's own HTML at `/app/:applicationID`,
// shown in a sandbox that hands it data over a message bridge. See
// app.intent.md.
export default {
  id: 'app',
  // No registry application of its own: `perApp` sends the gate to the
  // application named in the route instead. See sections.intent.md.
  app: null,
  perApp: true,
  routes: [
    {
      path: '/app/:applicationID',
      name: 'app',
      component: () => import('./views/AppView.vue'),
      meta: { section: 'app' },
    },
  ],
  sidebar: null,
}
