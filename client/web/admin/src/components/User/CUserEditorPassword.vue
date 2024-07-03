<template>
  <b-card
    data-test-id="card-user-password"
    header-class="border-bottom"
    footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <template #header>
      <h4
        data-test-id="card-title"
        class="m-0"
      >
        {{ $t('title') }}
      </h4>
    </template>

    <b-form
      @submit.prevent="onPasswordSubmit"
    >
      <b-row>
        <b-col cols="12">
          <b-form-group
            :label="$t('new')"
            :description="getPasswordWarning"
            label-class="text-primary"
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
        </b-col>

        <b-col cols="12">
          <b-form-group
            :label="$t('confirm')"
            :description="getConfirmPasswordWarning"
            label-class="text-primary"
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
        </b-col>
      </b-row>
    </b-form>

    <template #footer>
      <c-input-confirm
        data-test-id="button-remove-password"
        variant="light"
        size="md"
        @confirmed="$emit('submit')"
      >
        {{ $t('removePassword') }}
      </c-input-confirm>

      <c-corredor-manual-buttons
        ui-page="user/editor"
        ui-slot="passwordFooter"
        default-variant="light"
        @click="dispatchCortezaSystemEvent($event)"
      />

      <c-button-submit
        :disabled="!passwordState || !confirmPasswordState"
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="onPasswordSubmit"
      />
    </template>
  </b-card>
</template>

<script>

export default {
  name: 'CUserEditorPassword',

  i18nOptions: {
    namespaces: 'system.users',
    keyPrefix: 'editor.password',
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
