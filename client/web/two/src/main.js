import './config-check'

import App from './App.vue'

import { createApp } from 'vue'
const app = createApp(App)

import { setupAndAuthenticate } from './plugins'
// Await the async setup before mounting
setupAndAuthenticate(app).then(() => {
  app.mount('body')
})
