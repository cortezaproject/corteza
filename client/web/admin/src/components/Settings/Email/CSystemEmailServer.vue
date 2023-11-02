<template>
  <b-card
    header-bg-variant="white"
    footer-bg-variant="white"
    footer-class="d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <b-form
      @submit.prevent="submit()"
    >
      <b-row>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('host.label')"
            :description="$t('host.description')"
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
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('user.label')"
            :description="$t('user.description')"
          >
            <b-input
              v-model="server.user"
              data-test-id="input-user"
              :disabled="disabled"
              autocomplete="off"
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
            :label="$t('from.label')"
            :description="$t('from.description')"
          >
            <b-input
              v-model="server.from"
              data-test-id="input-sender-address"
              type="email"
              :disabled="disabled"
              autocomplete="off"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :description="$t('tlsInsecure.description')"
          >
            <b-form-checkbox
              v-model="server.tlsInsecure"
              data-test-id="checkbox-allow-invalid-certificates"
              :disabled="disabled"
            >
              {{ $t('tlsInsecure.label') }}
            </b-form-checkbox>
          </b-form-group>
        </b-col>
      </b-row>

      <b-row>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('tlsServerName.label')"
            :description="$t('tlsServerName.description')"
          >
            <b-input
              v-model="server.tlsServerName"
              data-test-id="input-tls-server-name"
              :disabled="disabled"
            />
          </b-form-group>
        </b-col>
      </b-row>
    </b-form>

    <template #header>
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>

    <template #footer>
      <c-button-submit
        data-test-id="button-smtp"
        :disabled="disabled"
        :processing="processingSMTPTest"
        :success="successSMTPTest"
        :text="$t('testSmtpConfigs.button')"
        variant="light"
        @submit="smtpConnectionCheck()"
      />

      <c-button-submit
        :disabled="disabled"
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="submit()"
      />
    </template>
  </b-card>
</template>

<script>
import Vue from 'vue'

export default {
  name: 'CComposeEditorBasic',

  i18nOptions: {
    namespaces: 'system.email',
    keyPrefix: 'editor.server',
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

    processingSMTPTest: {
      type: Boolean,
      value: false,
    },

    success: {
      type: Boolean,
      value: false,
    },

    successSMTPTest: {
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
