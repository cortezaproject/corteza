import Vue from 'vue'

import './config-check'
import './console-splash'

import './plugins'
import './mixins'
import './components'
import './filters'

import store from './store'

import { i18n, websocket } from '@cortezaproject/corteza-vue'

import router from './router'

export default (options = {}) => {
  options = {
    el: '#app',
    name: 'privacy',

    template: '<router-view v-if="loaded && i18nLoaded" />',

    data: () => ({
      loaded: false,
      i18nLoaded: false,
    }),

    async created () {
      this.$i18n.i18next.on('loaded', () => {
        this.i18nLoaded = true
      })

      this.websocket()

      return this.$auth.vue(this).handle().then(async ({ user }) => {
        if (user.meta.preferredLanguage) {
          // After user is authenticated, get his preferred language
          // and instruct i18next to change it
          this.$i18n.i18next.changeLanguage(user.meta.preferredLanguage)
        }

        // switch the webapp theme based on user preference
        if (user.meta.theme) {
          document.getElementsByTagName('html')[0].setAttribute('data-color-mode', user.meta.theme)
        }

        // Load effective permissions
        this.$store.dispatch('rbac/load')

        // Initialize notifications
        this.$store.dispatch('notifications/fetchNotifications')

        this.$Settings.init({ api: this.$SystemAPI }).then(() => {
          this.loaded = true

          // This bit removes code from the query params
          //
          // Vue router can't be used here because when on any child route there is no
          // guarantee that the route has loaded and so it may redirect us to the root page.
          //
          // @todo dig a bit deeper if there is a better vue-like solution; atm none were ok.
          const url = new URL(window.location.href)
          if (url.searchParams.get('code')) {
            url.searchParams.delete('code')
            window.location.replace(url.toString())
          }
        })
      })
        .catch((err) => {
          if (err instanceof Error && err.message === 'Unauthenticated') {
          // user not logged-in,
          // start with authentication flow
            this.$auth.startAuthenticationFlow()
            return
          }

          throw err
        })
    },

    methods: {
      /**
       * Registers event listener for websocket messages and
       * routes them depending on their type
       */
      websocket () {
        // cross-link auth & websocket so that ws can use the right access token
        websocket.init(this)

        // register event listener for messages
        this.$on('websocket-message', ({ data }) => {
          const msg = JSON.parse(data)
          switch (msg['@type']) {
            case 'notification':
              this.$store.dispatch('notifications/addNotification', msg['@value'])
              break

            case 'notification.read':
              this.$store.dispatch('notifications/updateReadNotification', msg['@value'])
              break

            case 'notification.read.all':
              this.$store.dispatch('notifications/updateAllReadNotifications', msg['@value'])
              break

            case 'notification.delete':
              this.$store.dispatch('notifications/removeNotification', msg['@value'])
              break

            case 'error':
              this.toastDanger('Websocket message with error', msg['@value'])
          }
        })
      },
    },

    router,
    store,
    i18n: i18n(Vue,
      { app: 'corteza-webapp-privacy' },
      'dashboard',
      'data-overview',
      'general',
      'map',
      'module-records',
      'navigation',
      'notification',
      'notifications',
      'request',
      'sensitive-data',
    ),
    // Any additional options we want to merge
    ...options,
  }

  const app = new Vue(options)

  // Simple HMR acceptance
  if (module.hot) {
    module.hot.accept()
  }

  return app
}
