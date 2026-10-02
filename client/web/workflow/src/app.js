import Vue from 'vue'

import './config-check'
import './console-splash'

import './filters'
import './plugins'
import './mixins'
import './components'
import store from './store'
import router from './router'

import { i18n, websocket } from '@cortezaproject/corteza-vue'

export default (options = {}) => {
  options = {
    el: '#app',
    name: 'workflow',
    template: '<router-view v-if="loaded && i18nLoaded" />',

    data: () => ({
      loaded: false,
      i18nLoaded: false,
    }),

    async created () {
      this.$i18n.i18next.on('initialized', () => {
        this.i18nLoaded = true
      })

      this.websocket()

      return this.$auth.vue(this).handle().then(async ({ user }) => {
        // switch the favicon based on the settings
        await this.$Settings.init({ api: this.$SystemAPI }).then(() => {
          const icon = this.$Settings.attachment('ui.iconLogo') || '/icon.svg'

          const favicon = document.getElementById('favicon')

          if (favicon) {
            favicon.href = icon
          }
        })

        if (user.meta.preferredLanguage) {
          // After user is authenticated, get his preferred language
          // and instruct i18next to change it
          this.$i18n.i18next.changeLanguage(user.meta.preferredLanguage)
        }

        // switch the webapp theme based on user preference
        // disabled for now as it's not working properly
        if (user.meta.theme) {
          document.getElementsByTagName('html')[0].setAttribute('data-color-mode', user.meta.theme)
        }

        // Load effective permissions
        this.$store.dispatch('rbac/load')

        // Initialize notifications
        this.$store.dispatch('notifications/fetchNotifications')

        this.loaded = true
      }).catch((err) => {
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

            case 'notification.unread':
              this.$store.dispatch('notifications/updateUnreadNotification', msg['@value'])
              break

            case 'notification.read.all':
              this.$store.dispatch('notifications/updateAllReadNotifications', msg['@value'])
              break

            case 'notification.unread.all':
              this.$store.dispatch('notifications/updateAllUnreadNotifications', msg['@value'])
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
      { app: 'corteza-webapp-workflow' },
      'configurator',
      'editor',
      'help',
      'general',
      'navigation',
      'notification',
      'notifications',
      'permissions',
      'configurator',
      'steps',
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
