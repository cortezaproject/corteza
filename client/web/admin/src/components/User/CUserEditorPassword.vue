<template>
  <b-card
    data-test-id="card-user-password"
    class="shadow-sm"
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <b-form
      @submit.prevent="onPasswordSubmit"
    >
      <b-form-group
        :label="$t('new')"
        label-cols="2"
        :description="getPasswordWarning"
      >
        <b-form-input
          v-model="password"
          data-test-id="input-new-password"
          :state="passwordState"
          autocomplete="new-password"
          required
          type="password"
        />
      </b-form-group>

      <b-form-group
        :label="$t('confirm')"
        label-cols="2"
        :description="getConfirmPasswordWarning"
        class="mb-0"
      >
        <b-form-input
          v-model="confirmPassword"
          data-test-id="input-confirm-password"
          type="password"
          autocomplete="new-password"
          required
          :disabled="!passwordState"
          :state="confirmPasswordState"
        />
      </b-form-group>
    </b-form>

    <template #header>
      <h3
        data-test-id="card-title"
        class="m-0"
      >
        {{ $t('title') }}
      </h3>
    </template>

    <template #footer>
      <confirmation-toggle
        data-test-id="button-remove-password"
        class="ml-1"
        cta-class="light"
        @confirmed="$emit('submit')"
      >
        {{ $t('removePassword') }}
      </confirmation-toggle>

      <c-submit-button
        class="float-right"
        :processing="processing"
        :success="success"
        :disabled="!passwordState || !confirmPasswordState"
        @submit="onPasswordSubmit"
      />

      <c-corredor-manual-buttons
        ui-page="user/editor"
        ui-slot="passwordFooter"
        default-variant="secondary"
        @click="dispatchCortezaSystemEvent($event)"
      />
    </template>
  </b-card>
</template>

<script>
import ConfirmationToggle from 'corteza-webapp-admin/src/components/ConfirmationToggle'
import CSubmitButton from 'corteza-webapp-admin/src/components/CSubmitButton'

export default {
  name: 'CUserEditorPassword',

  i18nOptions: {
    namespaces: 'system.users',
    keyPrefix: 'editor.password',
  },

  components: {
    ConfirmationToggle,
    CSubmitButton,
  },

  props: {
    processing: {
      type: Boolean,
      value: false,
    },
    success: {
      type: Boolean,
      value: false,
    },
    userID: {
      type: String,
      required: false,
      default: undefined,
    },
  },

  data () {
    return {
      password: '',
      confirmPassword: '',
      minPasswordLength: this.$Settings.get('auth.internal.passwordConstraints.minLength', 8),
    }
  },

  computed: {
    passwordState () {
      if (this.password.length > 0) {
        return this.password.length >= this.minPasswordLength
      }

      return null
    },

    confirmPasswordState () {
      if (this.passwordState && this.confirmPassword.length > 0) {
        return this.password === this.confirmPassword
      }

      return null
    },

    getPasswordWarning () {
      if (this.passwordState === false) {
        return this.$t('length', { length: this.minPasswordLength })
      }

      return null
    },

    getConfirmPasswordWarning () {
      if (this.confirmPasswordState === false) {
        return this.$t('missmatch')
      }

      return null
    },
  },

  methods: {
    onPasswordSubmit () {
      this.$emit('submit', this.password)

      this.password = ''
      this.confirmPassword = ''
    },
  },
}
</script>

<style scoped>

</style>
