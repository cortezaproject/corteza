import Vue from 'vue'

import './config-check'
import './console-splash'

import './components'
import './filters'
import './mixins'
import './plugins'

import router from './router'
import store from './store'

import { system } from '@cortezaproject/corteza-js'
import { corredor, i18n, mixins, websocket } from '@cortezaproject/corteza-vue'
import { mapGetters } from 'vuex'

const notProduction = (process.env.NODE_ENV !== 'production')

export default (options = {}) => {
  options = {
    el: '#app',
    name: 'admin',
    template: '<router-view v-if="loaded && i18nLoaded && isRbacLoaded" />',

    mixins: [
      mixins.corredor,
    ],

    data: () => ({
      loaded: false,
      i18nLoaded: false,
    }),

    computed: {
      ...mapGetters({
        isRbacLoaded: 'rbac/isLoaded',
      }),
    },

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

        // switch the page directionality on body based on language
        document.body.setAttribute('dir', this.textDirectionality(user.meta.preferredLanguage))

        if (user.meta.preferredLanguage) {
          // After user is authenticated, get his preferred language
          // and instruct i18next to change it
          this.$i18n.i18next.changeLanguage(user.meta.preferredLanguage)
        }

        // switch the webapp theme based on user preference
        if (user.meta.theme) {
          document.getElementsByTagName('html')[0].setAttribute('data-color-mode', user.meta.theme)
        }

        // ref to vue is needed inside compose helper
        // load and register bundle and list of client/server scripts
        const bundleLoaderOpt = {
          // Name of the bundle to load
          bundle: 'admin',

          // Debug logging
          verbose: notProduction,

          // Context for exec function (client scripts only!)
          //
          // Extended with additional helpers
          ctx: new corredor.WebappCtx({
            $invoker: user,
            authToken: this.$auth.accessToken,
          }),
        }

        // Load all pending prompts:
        this.$store.dispatch('wfPrompts/update')

        // Only use enabled apis
        const enabledApis = [this.$SystemAPI, this.$ComposeAPI, this.$AutomationAPI]
        if (this.$Settings.get('federation.enabled', false)) {
          enabledApis.push(this.$FederationAPI)
        }

        // Load effective permissions
        this.$store.dispatch('rbac/load', enabledApis)

        // Initialize notifications
        this.$store.dispatch('notifications/fetchNotifications')

        return this.loadBundle(bundleLoaderOpt)
          .then(() => this.$SystemAPI.automationList({ excludeInvalid: true }))
          .then(this.makeAutomationScriptsRegistrator(
            // compose specific handler that routes  onManual events for server-scripts
            // to the proper endpoint on the API
            system.TriggerSystemServerScriptOnManual(this.$SystemAPI),
          ))
          .then(() => {
            this.loaded = true
          })
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

        // register event listener for workflow messages
        this.$on('websocket-message', ({ data }) => {
          const msg = JSON.parse(data)
          switch (msg['@type']) {
            case 'workflowSessionPrompt':
              this.$store.dispatch('wfPrompts/new', msg['@value'])
              break

            case 'workflowSessionResumed':
              this.$store.dispatch('wfPrompts/clear', msg['@value'])
              break

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
      { app: 'corteza-webapp-admin' },
      'admin',
      'automation.scripts',
      'automation.sessions',
      'automation.workflows',
      'compose.automation',
      'compose.settings',
      'dashboard',
      'federation.nodes',
      'general',
      'navigation',
      'notification',
      'notifications',
      'permissions',
      'system.actionlog',
      'system.apigw',
      'system.applications',
      'system.authclients',
      'system.code-snippets',
      'system.connections',
      'system.email',
      'system.queues',
      'system.roles',
      'system.sensitivityLevel',
      'system.settings',
      'system.templates',
      'system.user-groups',
      'system.users',
      'ui.settings',
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
