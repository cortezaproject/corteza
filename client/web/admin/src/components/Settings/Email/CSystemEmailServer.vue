<template>
  <b-card
    header-class="border-bottom"
    footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <template #header>
      <h4 class="m-0">
        {{ $t('title') }}
      </h4>
    </template>

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
            label-class="text-primary"
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
      </b-row>

      <b-row>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('user.label')"
            :description="$t('user.description')"
            label-class="text-primary"
          >
            <b-input
              v-model="server.user"
              data-test-id="input-user"
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
            :label="$t('password.label')"
            :description="$t('password.description')"
            label-class="text-primary"
          >
            <b-input
              v-model="server.pass"
              data-test-id="input-password"
              type="password"
              :disabled="disabled"
              autocomplete="off"
            />
          </b-form-group>
        </b-col>
      </b-row>

      <hr>

      <b-row>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('from.label')"
            :description="$t('from.description')"
            label-class="text-primary"
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
      </b-row>

      <hr>

      <b-row>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('tlsServerName.label')"
            :description="$t('tlsServerName.description')"
            label-class="text-primary"
          >
            <b-input
              v-model="server.tlsServerName"
              data-test-id="input-tls-server-name"
              :disabled="disabled"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :description="$t('tlsInsecure.description')"
            class="mt-lg-3"
          >
            <b-form-checkbox
              v-model="server.tlsInsecure"
              data-test-id="checkbox-allow-invalid-certificates"
              :disabled="disabled"
              class="mt-lg-4 mb-2"
            >
              {{ $t('tlsInsecure.label') }}
            </b-form-checkbox>
          </b-form-group>
        </b-col>
      </b-row>
    </b-form>

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
