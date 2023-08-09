import Vue from 'vue'

import './config-check'
import './console-splash'

import './filters'
import './plugins'
import './mixins'
import './components'
import store from './store'

import { i18n } from '@cortezaproject/corteza-vue'

import router from './router'

export default (options = {}) => {
  options = {
    el: '#app',
    name: 'reporter',

    template: '<router-view v-if="loaded && i18nLoaded" />',

    data: () => ({
      loaded: false,
      i18nLoaded: false,
    }),

    async created () {
      this.$i18n.i18next.on('loaded', () => {
        this.i18nLoaded = true
      })

      this.$auth.handle().then(({ accessTokenFn, user }) => {
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

        // Load effective permissions
        this.$store.dispatch('rbac/load')

        this.$Settings.init({ api: this.$SystemAPI }).then(() => {
          this.loaded = true
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

    router,
    store,
    i18n: i18n(Vue,
      { app: 'corteza-webapp-reporter' },
      'builder',
      'create',
      'datasources',
      'display-element',
      'edit',
      'general',
      'list',
      'navigation',
      'notification',
      'permissions',
      'scenarios',
      'sidebar',
      'view',
    ),
    // Any additional options we want to merge
    ...options,
  }

  return new Vue(options)
}
