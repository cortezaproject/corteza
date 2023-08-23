import Vue from 'vue'

import './config-check'
import './console-splash'

import './filters'
import './plugins'
import './mixins'
import './components'

import store from './store'
import router from './router'

import { compose } from '@cortezaproject/corteza-js'
import { mixins, corredor, websocket, i18n } from '@cortezaproject/corteza-vue'

const notProduction = (process.env.NODE_ENV !== 'production')
const verboseEventbus = window.location.search.includes('verboseEventbus')

export default (options = {}) => {
  options = {
    el: '#app',
    name: 'compose',
    template: '<router-view v-if="loaded && i18nLoaded" />',

    mixins: [
      mixins.corredor,
    ],

    data: () => ({
      loaded: false,
      i18nLoaded: false,
    }),

    async created () {
      this.$i18n.i18next.on('loaded', () => {
        this.i18nLoaded = true
      })

      this.websocket()

      return this.$auth.vue(this).handle().then(({ user }) => {
        // switch the page directionality on body based on language
        document.body.setAttribute('dir', this.textDirectionality(user.meta.preferredLanguage))

        if (user.meta.preferredLanguage) {
          // After user is authenticated, get his preferred language
          // and instruct i18next to change it
          this.$i18n.i18next.changeLanguage(user.meta.preferredLanguage)

          /**
           * Let the API know what kind of language do we accept and send
           */
          this.$ComposeAPI
            .setHeader('Accept-Language', user.meta.preferredLanguage)
            .setHeader('Content-Language', user.meta.preferredLanguage)
        }

        // ref to vue is needed inside compose helper
        // load and register bundle and list of client/server scripts

        const bundleLoaderOpt = {
          // Name of the bundle to load
          bundle: 'compose',

          // Debug logging
          verbose: notProduction || verboseEventbus,

          // Context for exec function (client scripts only!)
          //
          // Extended with additional helpers
          ctx: new corredor.ComposeCtx(
            {
              $invoker: this.$auth.user,
              authToken: this.$auth.accessToken,
            },
            this,
          ),
        }

        // Load all pending prompts:
        this.$store.dispatch('wfPrompts/update')

        // Load effective permissions
        this.$store.dispatch('rbac/load')

        // Initializes reminders subsystems, do prefetch of all pending reminders
        this.$Reminder.init(this, { filter: { assignedTo: user.userID } })

        this.loadBundle(bundleLoaderOpt)
          .then(() => this.$ComposeAPI.automationList({ excludeInvalid: true }))
          .then(this.makeAutomationScriptsRegistrator(
            // compose specific handler that routes onManual events for server-scripts
            // to the proper endpoint on the API
            compose.TriggerComposeServerScriptOnManual(this.$ComposeAPI),
          ))

        this.$Settings.init({ api: this.$SystemAPI }).finally(() => {
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

            case 'reminder':
              this.$Reminder.enqueueRaw(msg['@value'])
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
      { app: 'corteza-webapp-compose' },
      'block',
      'chart',
      'field',
      'general',
      'module',
      'namespace',
      'navigation',
      'notification',
      'onboarding',
      'page',
      'permissions',
      'preview',
      'sidebar',
      'resource-translator',
    ),

    // Any additional options we want to merge
    ...options,
  }

  options.router.beforeEach((to, from, next) => {
    store.dispatch('ui/setPreviousPage', from)

    next()
  })

  return new Vue(options)
}
