<template>
  <b-card
    class="shadow-sm"
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <b-form
      @submit.prevent="$emit('submit', settings)"
    >
      <b-form-group
        :label="$t('internal.title')"
        label-size="lg"
        label-cols="2"
      >
        <b-form-checkbox
          v-model="settings['auth.internal.enabled']"
          :value="true"
          :unchecked-value="false"
          class="mt-3"
        >
          {{ $t('internal.enabled') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-model="settings['auth.internal.password-reset.enabled']"
          :value="true"
          :unchecked-value="false"
        >
          {{ $t('internal.password-reset.enabled') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-model="settings['auth.internal.signup.email-confirmation-required']"
          :value="true"
          :unchecked-value="false"
        >
          {{ $t('internal.signup.email-confirmation-required') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-model="settings['auth.internal.signup.enabled']"
          :value="true"
          :unchecked-value="false"
        >
          {{ $t('internal.signup.enabled') }}
        </b-form-checkbox>
      </b-form-group>
      <b-form-group
        label-cols="2"
        :description="$t('internal.signup.split-credentials-check.description')"
      >
        <b-form-checkbox
          v-model="settings['auth.internal.split-credentials-check']"
          :value="true"
          :unchecked-value="false"
        >
          {{ $t('internal.signup.split-credentials-check.label') }}
        </b-form-checkbox>
      </b-form-group>

      <hr>

      <h5>
        {{ $t('internal.password-constraints.title') }}
      </h5>

      <b-alert
        v-if="!$Settings.get('auth.internal.passwordConstraints.passwordSecurity')"
        variant="warning"
        show
      >
        {{ $t('internal.password-constraints.ignored-security') }}
      </b-alert>

      <b-form-group
        :label="$t('internal.password-constraints.min-upper-case-length')"
        :description="$t('internal.password-constraints.min-upper-case-description')"
        label-cols="2"
      >
        <b-form-input
          v-model.number="settings['auth.internal.password-constraints.min-upper-case']"
          type="number"
          :placeholder="`${defaultMinUppCaseChrs}`"
          :min="defaultMinUppCaseChrs"
        />
      </b-form-group>

      <b-form-group
        :label="$t('internal.password-constraints.min-lower-case-length')"
        :description="$t('internal.password-constraints.min-lower-case-description')"
        label-cols="2"
      >
        <b-form-input
          v-model.number="settings['auth.internal.password-constraints.min-lower-case']"
          type="number"
          :placeholder="`${defaultMinLowCaseChrs}`"
          :min="defaultMinLowCaseChrs"
        />
      </b-form-group>

      <b-form-group
        :label="$t('internal.password-constraints.min-length')"
        :description="$t('internal.password-constraints.min-length-description')"
        label-cols="2"
      >
        <b-form-input
          v-model.number="settings['auth.internal.password-constraints.min-length']"
          :placeholder="`${defaultMinPwd}`"
          :min="defaultMinPwd"
          type="number"
        />
      </b-form-group>

      <b-form-group
        :label="$t('internal.password-constraints.min-num-count')"
        :description="$t('internal.password-constraints.min-num-count-description')"
        label-cols="2"
      >
        <b-form-input
          v-model.number="settings['auth.internal.password-constraints.min-num-count']"
          placeholder="0"
          min="0"
          type="number"
        />
      </b-form-group>

      <b-form-group
        :label="$t('internal.password-constraints.min-special-count')"
        :description="$t('internal.password-constraints.min-special-count-description')"
        label-cols="2"
      >
        <b-form-input
          v-model.number="settings['auth.internal.password-constraints.min-special-count']"
          placeholder="0"
          min="0"
          type="number"
        />
      </b-form-group>

      <hr>

      <h5>
        {{ $t('mfa.title') }}
      </h5>

      <b-form-group label-cols="2">
        <b-form-checkbox
          v-model="settings['auth.multi-factor.email-otp.enabled']"
          :value="true"
          :unchecked-value="false"
        >
          {{ $t('mfa.emailOTP.enabled') }}
        </b-form-checkbox>
      </b-form-group>
      <b-form-group label-cols="2">
        <b-form-checkbox
          v-model="settings['auth.multi-factor.email-otp.enforced']"
          :value="true"
          :unchecked-value="false"
        >
          {{ $t('mfa.emailOTP.enforced') }}
        </b-form-checkbox>
      </b-form-group>
      <b-form-group
        :label="$t('mfa.emailOTP.expires.label')"
        :description="$t('mfa.emailOTP.expires.description')"
        label-cols="2"
      >
        <b-input-group append="seconds">
          <b-form-input
            v-model="settings['auth.multi-factor.email-otp.expires']"
            type="number"
            placeholder="60"
          />
        </b-input-group>
      </b-form-group>
      <b-form-group label-cols="2">
        <b-form-checkbox
          v-model="settings['auth.multi-factor.totp.enabled']"
          :value="true"
          :unchecked-value="false"
        >
          {{ $t('mfa.TOTP.enabled') }}
        </b-form-checkbox>
      </b-form-group>
      <b-form-group label-cols="2">
        <b-form-checkbox
          v-model="settings['auth.multi-factor.totp.enforced']"
          :value="true"
          :unchecked-value="false"
        >
          {{ $t('mfa.TOTP.enforced') }}
        </b-form-checkbox>
      </b-form-group>

      <b-form-group
        :label="$t('mfa.TOTP.issuer.label')"
        :description="$t('mfa.TOTP.issuer.description')"
        label-cols="2"
      >
        <b-input-group>
          <b-form-input
            v-model="settings['auth.multi-factor.totp.issuer']"
            placeholder="Corteza"
          />
        </b-input-group>
      </b-form-group>

      <hr>

      <h5>
        {{ $t('mail.title') }}
      </h5>

      <b-form-group
        :label="$t('mail.from-address')"
        label-cols="2"
        :description="$t('mail.validate-email')"
      >
        <b-input-group>
          <b-form-input
            v-model="settings['auth.mail.from-address']"
            type="email"
          />
        </b-input-group>
      </b-form-group>
      <b-form-group
        :label="$t('mail.from-name')"
        label-cols="2"
      >
        <b-input-group>
          <b-form-input v-model="settings['auth.mail.from-name']" />
        </b-input-group>
      </b-form-group>
    </b-form>

    <template #header>
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>

    <template #footer>
      <c-submit-button
        class="float-right"
        :disabled="!canManage"
        :processing="processing"
        :success="success"
        @submit="$emit('submit', settings)"
      />
    </template>
  </b-card>
</template>

<script>
import CSubmitButton from 'corteza-webapp-admin/src/components/CSubmitButton'

export default {
  name: 'CSystemEditorAuth',

  i18nOptions: {
    namespaces: 'system.settings',
    keyPrefix: 'editor.auth',
  },

  components: {
    CSubmitButton,
  },

  props: {
    settings: {
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

    canManage: {
      type: Boolean,
      required: true,
    },
  },

  data () {
    return {
      defaultMinPwd: 8,
      defaultMinUppCaseChrs: 0,
      defaultMinLowCaseChrs: 0,
    }
  },
}
</script>
