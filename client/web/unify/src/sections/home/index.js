// Home section — the `/` landing. No left sidebar; it renders its own
// columns (apps / agent / notifications) inline, so we hide the topbar's
// agent + notification toggles for this section.
export default {
  id: 'home',
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('./views/Home.vue'),
      meta: { section: 'home' },
    },
  ],
  sidebar: null,
  topbar: {
    // Home lists the apps inline (left column), so the topbar's app-menu
    // button is redundant here — hidden, like the agent + notification toggles.
    hideAppSelector: true,
    hideAgentSidebar: true,
    hideNotifications: true,
  },
}
