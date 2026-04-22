import './config-check'

import App from './App.vue'

import { createApp } from 'vue'
const app = createApp(App)

import { setupAndAuthenticate } from './plugins'
setupAndAuthenticate(app)
  .then(() => {
    app.mount('body')
  })
  .catch(err => {
    console.error('App setup failed', err)
  })
