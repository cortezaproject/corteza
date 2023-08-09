<template>
  <div id="app">
    <div>
      <router-view
        v-if="loaded && i18nLoaded"
        class="h-100"
      />
    </div>
  </div>
</template>

<script>
export default {

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
      }

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
}
</script>
