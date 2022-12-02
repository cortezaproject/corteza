<template>
  <b-card
    header-bg-variant="white"
    footer-bg-variant="white"
    class="shadow-sm"
  >
    <b-form
      @submit.prevent="submit()"
    >
      <b-form-group
        :label="$t('host.label')"
        :description="$t('host.description')"
        label-cols="2"
      >
        <b-input-group>
          <b-input
            v-model="server.host"
            data-test-id="input-server"
            :disabled="disabled"
            placeholder="host.domain.tld"
            autocomplete="off"
            required
          />
          <b-input-group-append is-text>
            :
          </b-input-group-append>
          <b-input-group-append>
            <b-input
              v-model="server.port"
              data-test-id="input-server-port"
              type="number"
              :disabled="disabled"
              step="1"
              required
            />
          </b-input-group-append>
        </b-input-group>
      </b-form-group>

      <b-form-group
        :label="$t('user.label')"
        :description="$t('user.description')"
        label-cols="2"
      >
        <b-input
          v-model="server.user"
          data-test-id="input-user"
          :disabled="disabled"
          autocomplete="off"
        />
      </b-form-group>

      <b-form-group
        :label="$t('password.label')"
        :description="$t('password.description')"
        label-cols="2"
      >
        <b-input
          v-model="server.pass"
          data-test-id="input-password"
          type="password"
          :disabled="disabled"
          autocomplete="off"
        />
      </b-form-group>

      <hr>

      <b-form-group
        :label="$t('from.label')"
        :description="$t('from.description')"
        label-cols="2"
      >
        <b-input
          v-model="server.from"
          data-test-id="input-sender-address"
          type="email"
          :disabled="disabled"
          autocomplete="off"
        />
      </b-form-group>

      <hr>

      <b-form-group
        :description="$t('tlsInsecure.description')"
        label-cols="2"
      >
        <b-form-checkbox
          v-model="server.tlsInsecure"
          data-test-id="checkbox-allow-invalid-certificates"
          :disabled="disabled"
        >
          {{ $t('tlsInsecure.label') }}
        </b-form-checkbox>
      </b-form-group>

      <b-form-group
        :label="$t('tlsServerName.label')"
        :description="$t('tlsServerName.description')"
        label-cols="2"
      >
        <b-input
          v-model="server.tlsServerName"
          data-test-id="input-tls-server-name"
          :disabled="disabled"
        />
      </b-form-group>
    </b-form>

    <template #header>
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>

    <template #footer>
      <c-submit-button
        :disabled="disabled"
        :processing="processing"
        :success="success"
        :cypress-i-d="'button-smtp'"
        variant="light"
        class="float-left"
        @submit="smtpConnectionCheck()"
      >
        {{ $t('testSmtpConfigs.button') }}
      </c-submit-button>

      <c-submit-button
        :disabled="disabled"
        :processing="processing"
        :success="success"
        class="float-right"
        @submit="submit()"
      />
    </template>
  </b-card>
</template>

<script>
import Vue from 'vue'
import CSubmitButton from 'corteza-webapp-admin/src/components/CSubmitButton'

export default {
  name: 'CComposeEditorBasic',

  i18nOptions: {
    namespaces: 'system.email',
    keyPrefix: 'editor.server',
  },

  components: {
    CSubmitButton,
  },

  props: {
    value: {
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

    disabled: {
      type: Boolean,
      value: true,
    },
  },

  data () {
    return {
      server: Vue.util.extend({
        host: '',
        port: 25,
        user: '',
        pass: '',
        from: '',
        tlsInsecure: false,
        tlsServerName: '',
      }, this.value),
    }
  },

  methods: {
    submit () {
      this.$emit('submit', this.server)
    },

    smtpConnectionCheck () {
      this.$emit('smtpConnectionCheck', this.server)
    },
  },
}
</script>
