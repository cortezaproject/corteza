<template>
  <section
    class="c311-data-state"
    :data-state="state"
    :aria-busy="state === 'loading' ? 'true' : 'false'"
    data-c311-data-state
  >
    <slot v-if="state === 'populated'" name="populated" />
    <div v-else class="c311-data-state__message" role="status">
      <h2 class="h5 mb-2">
        {{ headingForState }}
      </h2>
      <p class="mb-3">
        {{ messageForState }}
      </p>
      <p v-if="serverVersionMessage" class="small text-muted" data-c311-server-version>
        {{ serverVersionMessage }}
      </p>
      <b-button
        v-if="state === 'retryable-error' && retryable"
        variant="primary"
        @click="$emit('retry')"
      >
        {{ displayRetryLabel }}
      </b-button>
      <slot name="action" />
    </div>
  </section>
</template>

<script>
export default {
  name: 'C311DataState',
  props: {
    state: {
      type: String,
      required: true,
    },
    messages: {
      type: Object,
      default: () => ({}),
    },
    error: {
      type: Object,
      default: null,
    },
    retryable: {
      type: Boolean,
      default: true,
    },
    retryLabel: {
      type: String,
      default: 'Try again',
    },
  },
  computed: {
    serverVersionMessage () {
      const version = this.error && (this.error.current_version || this.error.currentVersion)
      return version === undefined || version === null
        ? ''
        : this.translate('status.versionConflict', `Current server version: ${version}. Reload before trying again.`, { version })
    },
    headingForState () {
      const fallback = {
        loading: 'Loading',
        empty: 'No matching records',
        forbidden: 'Access denied',
        'not-found': 'Not found',
        'validation-error': 'Check your information',
        'retryable-error': 'Temporarily unavailable',
        'terminal-error': 'Unable to complete this operation',
      }[this.state] || ''
      return this.messages[`${this.state}.heading`] || this.translate(`status.${this.state}`, fallback)
    },
    messageForState () {
      const fallback = {
        loading: 'Loading data.',
        empty: 'Zero matching records were found.',
        forbidden: 'You do not have access to this information.',
        'not-found': 'The requested information could not be found.',
        'validation-error': 'Correct the highlighted fields and try again.',
        'retryable-error': 'The service is temporarily unavailable.',
        'terminal-error': 'The operation could not be completed.',
      }[this.state] || ''
      return this.messages[`${this.state}.message`] || this.translate(`status.${this.state}.message`, fallback)
    },
    displayRetryLabel () {
      return this.translate('action.retry', this.retryLabel)
    },
  },
  methods: {
    translate (key, fallback, options = {}) {
      const translated = this.$t?.(`c311:${key}`, options)
      return translated && translated !== `c311:${key}` && translated !== key ? translated : fallback
    },
  },
}
</script>

<style scoped>
.c311-data-state__message {
  padding: 1.5rem 0;
}
</style>
