<template>
  <b-card
    class="shadow-sm"
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <b-row>
      <b-col cols="8">
        <span
          v-if="mfa.enforcedEmailOTP"
          v-html="$t('emailOTP.enabled.text')"
        />
        <span
          v-else
          v-html="$t('emailOTP.disabled.text')"
        />
      </b-col>
      <b-col
        cols="4"
        class="text-right"
      >
        <b-button
          v-if="mfa.enforcedEmailOTP"
          @click="$emit('patch', '/meta/securityPolicy/mfa/enforcedEmailOTP', false)"
        >
          {{ $t('emailOTP.disable.label') }}
        </b-button>
        <b-button
          v-else
          @click="$emit('patch', '/meta/securityPolicy/mfa/enforcedEmailOTP', true)"
        >
          {{ $t('emailOTP.enable.label') }}
        </b-button>
      </b-col>
    </b-row>
    <b-row class="mt-2 pt-2 border-top">
      <b-col cols="8">
        <span
          v-if="mfa.enforcedTOTP"
          v-html="$t('TOTP.enabled.text')"
        />
        <span
          v-else
          v-html="$t('TOTP.disabled.text')"
        />
      </b-col>
      <b-col
        cols="4"
        class="text-right"
      >
        <b-button
          :disabled="!mfa.enforcedTOTP"
          @click="$emit('patch', '/meta/securityPolicy/mfa/enforcedTOTP', false)"
        >
          {{ $t('TOTP.remove.label') }}
        </b-button>
      </b-col>
    </b-row>
    <template #header>
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>
  </b-card>
</template>

<script>
/**
 * @todo find a way to get this number from the backend
 * @type {number}
 */
// const minPasswordLength = 5

export default {
  name: 'CUserEditorMFA',

  i18nOptions: {
    namespaces: 'system.users',
    keyPrefix: 'editor.mfa',
  },

  props: {
    mfa: {
      type: Object,
      required: true,
    },

    processing: {
      type: Boolean,
      value: false,
    },
    success: {
      type: Boolean,
      value: false,
    },
  },
}
</script>

<style scoped>

</style>
