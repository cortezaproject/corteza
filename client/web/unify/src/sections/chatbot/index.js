// Chatbot section — mounted under /chatbot. Paths mirror the legacy app 1:1
// (prefixed); legacy `root`→`chatbot` and `sessions`→`chatbot.sessions` to
// namespace the generic names.
import ChatbotSidebar from './sidebar/ChatbotSidebar.vue'

export default {
  id: 'chatbot',
  // Registry application gating entry to this section: the user must hold
  // `access` on it. See sections.intent.md.
  app: 'chatbot/',
  routes: [
    {
      path: '/chatbot',
      name: 'chatbot',
      component: () => import('./views/List.vue'),
      meta: { section: 'chatbot' },
    },
    {
      path: '/chatbot/create',
      name: 'chatbot.create',
      component: () => import('./views/Editor.vue'),
      meta: { section: 'chatbot' },
    },
    {
      path: '/chatbot/:chatbotID/edit',
      name: 'chatbot.edit',
      component: () => import('./views/Editor.vue'),
      meta: { section: 'chatbot' },
    },
    {
      path: '/chatbot/sessions',
      name: 'chatbot.sessions',
      component: () => import('./views/Sessions.vue'),
      meta: { section: 'chatbot' },
    },
  ],
  sidebar: ChatbotSidebar,
}
