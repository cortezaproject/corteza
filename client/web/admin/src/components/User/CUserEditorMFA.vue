<template>
  <b-card
    header-bg-variant="white"
    class="shadow-sm"
  >
    <div class="d-flex align-items-center flex-wrap">
      <div>
        <span
          v-if="mfa.enforcedEmailOTP"
          v-html="$t('emailOTP.enabled.text')"
        />
        <span
          v-else
          v-html="$t('emailOTP.disabled.text')"
        />
      </div>
      <div class="ml-auto">
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
      </div>
    </div>

    <div class="d-flex align-items-center justify-content-between flex-wrap mt-2 pt-2 border-top">
      <div>
        <span
          v-if="mfa.enforcedTOTP"
          v-html="$t('TOTP.enabled.text')"
        />
        <span
          v-else
          v-html="$t('TOTP.disabled.text')"
        />
      </div>
      <div class="ml-auto">
        <b-button
          :disabled="!mfa.enforcedTOTP"
          @click="$emit('patch', '/meta/securityPolicy/mfa/enforcedTOTP', false)"
        >
          {{ $t('TOTP.remove.label') }}
        </b-button>
      </div>
    </div>
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
