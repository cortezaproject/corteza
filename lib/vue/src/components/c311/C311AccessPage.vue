<template>
  <main class="c311-access-page p-4" data-c311-main tabindex="-1" :aria-labelledby="headingID">
    <h1 :id="headingID" ref="heading" tabindex="-1" class="h2">
      {{ heading }}
    </h1>
    <p>{{ message }}</p>
    <b-button variant="primary" @click="goBack">
      {{ backLabel }}
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
  },
  mounted () {
    this.$nextTick(() => this.$refs.heading?.focus())
  },
  methods: {
    goBack () {
      const returnTo = this.$route?.query?.returnTo
      if (returnTo) return this.$router.replace(decodeURIComponent(returnTo))
      return this.$router.back()
    },
  },
}
</script>
