<template>
  <b-card
    class="shadow-sm"
    header-bg-variant="white"
    footer-bg-variant="white"
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
        class="float-right"
        :disabled="disabled"
        :processing="processing"
        :success="success"
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
  },
}
</script>
