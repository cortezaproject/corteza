import '@/assets/styles.css'
import { ToastPlugin } from '@cortezaproject/corteza-vue-next'

import 'primeicons/primeicons.css'
import PrimeVue from 'primevue/config'
import Ripple from 'primevue/ripple'

import ConfirmationService from 'primevue/confirmationservice'
import DialogService from 'primevue/dialogservice'
import ToastService from 'primevue/toastservice'

import { getTheme } from '@cortezaproject/corteza-vue-next'

export const UIPlugin = {
  install(app, options = {}) {
    app.use(PrimeVue, {
      theme: {
        preset: getTheme(options.theme),
        options: {
          darkModeSelector: '.dark',
          cssLayer: {
            name: 'primevue',
            order: 'tailwind-base, primevue, tailwind-utilities',
          },
        },
      },
      ripple: true,
    })

    app.directive('ripple', Ripple)

    app.use(ConfirmationService)
    app.use(ToastService)
    app.use(DialogService)

    app.use(ToastPlugin)
  },
}
