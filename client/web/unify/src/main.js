import './config-check'

import App from './App.vue'

import { createApp } from 'vue'
const app = createApp(App)

import { setupAndAuthenticate } from './plugins'
// Await the async setup before mounting
setupAndAuthenticate(app)
  .then(ready => {
    // false means setup handed over to the auth flow and installed none of the
    // plugins, so mounting would render a broken page for exactly as long as the
    // redirect takes to land — and log a "Need to install with app.use" crash on
    // the way. Unauthenticated users go to the auth flow, never to that page.
    if (ready) {
      app.mount('body')
    }
  })
  .catch(err => {
    console.error('App setup failed', err)
  })
