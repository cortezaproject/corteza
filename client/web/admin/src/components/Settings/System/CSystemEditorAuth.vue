<template>
  <b-card
    data-test-id="card-edit-authentication"
    header-bg-variant="white"
    footer-bg-variant="white"
    footer-class="d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <template #header>
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>

    <b-form
      @submit.prevent="$emit('submit', settings)"
    >
      <b-row>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('internal.title')"
            label-class="text-primary"
          >
            <c-input-checkbox
              v-model="settings['auth.internal.enabled']"
              :value="true"
              switch
              :unchecked-value="false"
              :labels="checkboxLabel"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('internal.enabled')"
            label-class="text-primary"
          >
            <c-input-checkbox
              v-model="settings['auth.internal.enabled']"
              switch
              :value="true"
              :unchecked-value="false"
              :labels="checkboxLabel"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('internal.password-reset.enabled')"
            label-class="text-primary"
          >
            <c-input-checkbox
              v-model="settings['auth.internal.password-reset.enabled']"
              switch
              data-test-id="checkbox-password-reset"
              :value="true"
              :labels="checkboxLabel"
              :unchecked-value="false"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('internal.signup.email-confirmation-required')"
            label-class="text-primary"
          >
            <c-input-checkbox
              v-model="settings['auth.internal.signup.email-confirmation-required']"
              switch
              :value="true"
              :labels="checkboxLabel"
              :unchecked-value="false"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('internal.signup.enabled')"
            label-class="text-primary"
          >
            <c-input-checkbox
              v-model="settings['auth.internal.signup.enabled']"
              :value="true"
              :unchecked-value="false"
              :labels="checkboxLabel"
              switch
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('internal.profile-avatar.enabled')"
            label-class="text-primary"
          >
            <c-input-checkbox
              v-model="settings['auth.internal.profile-avatar.enabled']"
              :value="true"
              :unchecked-value="false"
              :labels="checkboxLabel"
              switch
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('internal.signup.split-credentials-check.label')"
            label-class="text-primary"
          >
            <c-input-checkbox
              v-model="settings['auth.internal.split-credentials-check']"
              :value="true"
              :unchecked-value="false"
              :labels="checkboxLabel"
              switch
            />
          </b-form-group>
        </b-col>
      </b-row>

      <hr>

      <div>
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

        <b-row>
          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('internal.password-constraints.min-upper-case-length')"
              :description="$t('internal.password-constraints.min-upper-case-description')"
              label-class="text-primary"
            >
              <b-form-input
                v-model.number="settings['auth.internal.password-constraints.min-upper-case']"
                type="number"
                :placeholder="`${defaultMinUppCaseChrs}`"
                :min="defaultMinUppCaseChrs"
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('internal.password-constraints.min-lower-case-length')"
              :description="$t('internal.password-constraints.min-lower-case-description')"
              label-class="text-primary"
            >
              <b-form-input
                v-model.number="settings['auth.internal.password-constraints.min-lower-case']"
                type="number"
                :placeholder="`${defaultMinLowCaseChrs}`"
                :min="defaultMinLowCaseChrs"
              />
            </b-form-group>
          </b-col>
        </b-row>

        <b-row>
          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('internal.password-constraints.min-length')"
              :description="$t('internal.password-constraints.min-length-description')"
              label-class="text-primary"
            >
              <b-form-input
                v-model.number="settings['auth.internal.password-constraints.min-length']"
                :placeholder="`${defaultMinPwd}`"
                :min="defaultMinPwd"
                type="number"
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('internal.password-constraints.min-num-count')"
              :description="$t('internal.password-constraints.min-num-count-description')"
              label-class="text-primary"
            >
              <b-form-input
                v-model.number="settings['auth.internal.password-constraints.min-num-count']"
                placeholder="0"
                min="0"
                type="number"
              />
            </b-form-group>
          </b-col>
        </b-row>

        <b-row>
          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('internal.password-constraints.min-special-count')"
              :description="$t('internal.password-constraints.min-special-count-description')"
              label-class="text-primary"
            >
              <b-form-input
                v-model.number="settings['auth.internal.password-constraints.min-special-count']"
                placeholder="0"
                min="0"
                type="number"
              />
            </b-form-group>
          </b-col>
        </b-row>
      </div>

      <hr>

      <div>
        <b-row>
          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('mfa.emailOTP.enabled')"
              label-class="text-primary"
            >
              <c-input-checkbox
                v-model="settings['auth.multi-factor.email-otp.enabled']"
                data-test-id="checkbox-enable-emailOTP"
                :value="true"
                :unchecked-value="false"
                :labels="checkboxLabel"
                switch
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('mfa.emailOTP.enforced')"
              label-class="text-primary"
            >
              <c-input-checkbox
                v-model="settings['auth.multi-factor.email-otp.enforced']"
                :value="true"
                :unchecked-value="false"
                :labels="checkboxLabel"
                switch
              />
            </b-form-group>
          </b-col>
        </b-row>

        <b-row>
          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('mfa.TOTP.enabled')"
              label-class="text-primary"
            >
              <c-input-checkbox
                v-model="settings['auth.multi-factor.totp.enabled']"
                data-test-id="checkbox-enable-TOTP"
                :value="true"
                :unchecked-value="false"
                :labels="checkboxLabel"
                switch
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('mfa.TOTP.enforced')"
              label-class="text-primary"
            >
              <c-input-checkbox
                v-model="settings['auth.multi-factor.totp.enforced']"
                :value="true"
                :unchecked-value="false"
                :labels="checkboxLabel"
                switch
              />
            </b-form-group>
          </b-col>
        </b-row>

        <b-row>
          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('mfa.emailOTP.expires.label')"
              :description="$t('mfa.emailOTP.expires.description')"
              label-class="text-primary"
            >
              <b-input-group append="seconds">
                <b-form-input
                  v-model="settings['auth.multi-factor.email-otp.expires']"
                  type="number"
                  placeholder="60"
                />
              </b-input-group>
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('mfa.TOTP.issuer.label')"
              :description="$t('mfa.TOTP.issuer.description')"
              label-class="text-primary"
            >
              <b-input-group>
                <b-form-input
                  v-model="settings['auth.multi-factor.totp.issuer']"
                  placeholder="Corteza"
                />
              </b-input-group>
            </b-form-group>
          </b-col>
        </b-row>
      </div>

      <hr>

      <div>
        <h5>
          {{ $t('mail.title') }}
        </h5>

        <b-row>
          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('mail.from-address')"
              :description="$t('mail.validate-email')"
              label-class="text-primary"
            >
              <b-input-group>
                <b-form-input
                  v-model="settings['auth.mail.from-address']"
                  type="email"
                />
              </b-input-group>
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('mail.from-name')"
              label-class="text-primary"
            >
              <b-input-group>
                <b-form-input v-model="settings['auth.mail.from-name']" />
              </b-input-group>
            </b-form-group>
          </b-col>
        </b-row>
      </div>

      <hr>

      <div>
        <h5>
          {{ $t('internal.send-user-invite-email.title') }}
        </h5>

        <b-row>
          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('internal.send-user-invite-email.enabled')"
              :description="$t('internal.send-user-invite-email.description')"
              label-class="text-primary"
            >
              <c-input-checkbox
                v-model="settings['auth.internal.send-user-invite-email.enabled']"
                :value="true"
                :unchecked-value="false"
                :labels="checkboxLabel"
                switch
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('internal.send-user-invite-email.expires.label')"
              :description="$t('internal.send-user-invite-email.expires.description')"
              label-class="text-primary"
            >
              <b-input-group append="hours">
                <b-form-input
                  v-model="settings['auth.internal.send-user-invite-email.expires']"
                  type="number"
                  placeholder="72"
                />
              </b-input-group>
            </b-form-group>
          </b-col>
        </b-row>
      </div>
    </b-form>

    <template #footer>
      <c-button-submit
        v-if="canManage"
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="$emit('submit', settings)"
      />
    </template>
  </b-card>
</template>

<script>

export default {
  name: 'CSystemEditorAuth',

  i18nOptions: {
    namespaces: 'system.settings',
    keyPrefix: 'editor.auth',
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
      checkboxLabel: {
        on: this.$t('general:label.general.yes'),
        off: this.$t('general:label.general.no'),
      },
    }
  },
}
</script>
