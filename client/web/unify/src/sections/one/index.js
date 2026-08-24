// One section — the full-page application catalog at `/one`: a searchable grid
// of every application the user may open, under the instance logo. Its registry
// application ships disabled and unlisted, so the section is only reachable on
// an installation that turns it on.
export default {
  id: 'one',
  // Registry application gating entry to this section: the user must hold
  // `access` on it. See sections.intent.md.
  app: 'one/',
  routes: [
    {
      path: '/one',
      name: 'one',
      component: () => import('./views/AppList.vue'),
      meta: { section: 'one' },
    },
  ],
  sidebar: null,
  topbar: {
    // The page is the app selector, so the topbar's button would open a menu
    // of what is already on screen.
    hideAppSelector: true,
  },
}
