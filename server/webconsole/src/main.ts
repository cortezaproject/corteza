import { createApp } from 'vue'

import App from './WebConsole.vue'
import router from './views'

const app = createApp(App)

app.use(router)

app.mount('#app')
