import Vue from 'vue'

import './config-check'
import './console-splash'

import './plugins'
import './mixins'
import './components'
import './filters'

import store from './store'
import router from './router'

import { system } from '@cortezaproject/corteza-js'
import { mixins, corredor, websocket, i18n } from '@cortezaproject/corteza-vue'
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
      this.$i18n.i18next.on('loaded', () => {
        this.i18nLoaded = true
      })

      this.websocket()

      return this.$auth.vue(this).handle().then(async ({ user }) => {
        // switch the page directionality on body based on language
        document.body.setAttribute('dir', this.textDirectionality(user.meta.preferredLanguage))

        if (user.meta.preferredLanguage) {
          // After user is authenticated, get his preferred language
          // and instruct i18next to change it
          this.$i18n.i18next.changeLanguage(user.meta.preferredLanguage)
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

        await this.$Settings.init({ api: this.$SystemAPI })

        // Load all pending prompts:
        this.$store.dispatch('wfPrompts/update')

        // Only use enabled apis
        const enabledApis = [this.$SystemAPI, this.$ComposeAPI, this.$AutomationAPI]
        if (this.$Settings.get('federation.enabled', false)) {
          enabledApis.push(this.$FederationAPI)
        }

        // Load effective permissions
        this.$store.dispatch('rbac/load', enabledApis)

        return this.loadBundle(bundleLoaderOpt)
          .then(() => this.$SystemAPI.automationList({ excludeInvalid: true }))
          .then(this.makeAutomationScriptsRegistrator(
            // compose specific handler that routes  onManual events for server-scripts
            // to the proper endpoint on the API
            system.TriggerSystemServerScriptOnManual(this.$SystemAPI),
          ))
          .then(() => {
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

            case 'error':
              console.error('websocket message with error', msg['@value'])
          }
        })
      },
    },

    router,
    store,
    i18n: i18n(Vue,
      { app: 'corteza-webapp-admin' },
      'admin',
      'dashboard',
      'general',
      'navigation',
      'notification',
      'general',
      'permissions',
      'system.stats',
      'system.applications',
      'system.users',
      'system.roles',
      'system.templates',
      'system.scripts',
      'system.settings',
      'system.authclients',
      'system.permissions',
      'system.actionlog',
      'system.queues',
      'system.apigw',
      'system.connections',
      'system.sensitivityLevel',
      'compose.settings',
      'compose.permissions',
      'compose.automation',
      'federation.nodes',
      'federation.permissions',
      'automation.workflows',
      'automation.scripts',
      'automation.sessions',
      'automation.permissions',
      'ui.settings',
    ),

    // Any additional options we want to merge
    ...options,
  }

  return new Vue(options)
}
