// Agentic section — mounted under /agentic. Paths mirror the legacy agentic
// app 1:1 (just prefixed), and route names are namespaced under `agentic.*`
// to avoid collisions with other sections.
import AgenticSidebar from './sidebar/AgenticSidebar.vue'

export default {
  id: 'agentic',
  routes: [
    {
      path: '/agentic',
      name: 'agentic',
      component: () => import('./views/Home.vue'),
      // List (root) view keeps the sidebar collapsed, like the legacy app.
      meta: { section: 'agentic', hideSidebar: true },
    },
    {
      path: '/agentic/create',
      name: 'agentic.create',
      component: () => import('./views/Editor.vue'),
      meta: { section: 'agentic' },
    },
    {
      path: '/agentic/:agentID/edit',
      name: 'agentic.edit',
      component: () => import('./views/Editor.vue'),
      meta: { section: 'agentic' },
    },
  ],
  sidebar: AgenticSidebar,
}
