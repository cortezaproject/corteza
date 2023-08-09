import Vue from 'vue'

import './config-check'
import './console-splash'

import './plugins'
import './mixins'
import './components'
import store from './store'
import router from './router'

import { i18n, websocket } from '@cortezaproject/corteza-vue'

export default (options = {}) => {
  options = {
    el: '#app',
    name: 'one',
    template: '<router-view v-if="loaded && i18nLoaded" />',

    data: () => ({
      loaded: false,
      i18nLoaded: false,
    }),

    async created () {
      this.$i18n.i18next.on('loaded', () => {
        this.i18nLoaded = true
      })

      return this.$auth.vue(this).handle().then(({ user }) => {
        // switch the page directionality on body based on language
        document.body.setAttribute('dir', this.textDirectionality(user.meta.preferredLanguage))

        if (user.meta.preferredLanguage) {
          // After user is authenticated, get his preferred language
          // and instruct i18next to change it
          this.$i18n.i18next.changeLanguage(user.meta.preferredLanguage)
        }

        this.$store.dispatch('wfPrompts/update')

        return this.$Settings.init({ api: this.$SystemAPI }).finally(() => {
          this.websocket()

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
            case 'workflowSessionPrompt': {
              this.$store.dispatch('wfPrompts/new', msg['@value'])
              break
            }

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
      { app: 'corteza-webapp-one' },
      'app',
      'layout',
      'navigation',
    ),

    // Any additional options we want to merge
    ...options,
  }

  return new Vue(options)
}
