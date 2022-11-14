import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './WebConsole.vue'
import router from './views'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')
