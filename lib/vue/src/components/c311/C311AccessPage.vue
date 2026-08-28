<template>
  <main class="c311-access-page p-4" data-c311-main tabindex="-1" :aria-labelledby="headingID">
    <h1 :id="headingID" ref="heading" tabindex="-1" class="h2">
      {{ displayHeading }}
    </h1>
    <p>{{ displayMessage }}</p>
    <b-button variant="primary" @click="goBack">
      {{ displayBackLabel }}
    </b-button>
  </main>
</template>

<script>
export default {
  name: 'C311AccessPage',
  props: {
    status: { type: [Number, String], default: 403 },
    heading: { type: String, default: '' },
    message: { type: String, default: '' },
    backLabel: { type: String, default: 'Return' },
  },
  computed: {
    headingID () { return `c311-access-heading-${this._uid}` },
    displayHeading () {
      const fallback = { 401: 'Sign-in required', 403: 'Access denied', 404: 'Not found' }[this.status] || 'Access status'
      return this.heading || this.translate(`access.${this.status}.heading`, fallback)
    },
    displayMessage () {
      const fallback = { 401: 'Your session has expired or you are not signed in.', 403: 'You do not have permission to view this page.', 404: 'The requested page could not be found.' }[this.status] || 'This page is not available.'
      return this.message || this.translate(`access.${this.status}.message`, fallback)
    },
    displayBackLabel () {
      return this.translate('action.return', this.backLabel)
    },
  },
  mounted () {
    this.$nextTick(() => this.$refs.heading?.focus())
  },
  methods: {
    translate (key, fallback) {
      const translated = this.$t?.(`c311:${key}`)
      return translated && translated !== `c311:${key}` && translated !== key ? translated : fallback
    },
    goBack () {
      const returnTo = this.$route?.query?.returnTo
      if (returnTo) return this.$router.replace(decodeURIComponent(returnTo))
      return this.$router.back()
    },
  },
}
</script>
