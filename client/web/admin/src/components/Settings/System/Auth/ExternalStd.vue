<template>
  <div>
    <b-form-group>
      <b-form-checkbox
        v-model="value.enabled"
      >
        {{ $t('enabled') }}
      </b-form-checkbox>
    </b-form-group>

    <b-form-group
      :label="$t('clientKey')"
      label-class="text-primary"
    >
      <b-input-group>
        <b-form-input
          v-model.trim="value.key"
          :required="value.enabled"
        />
      </b-input-group>
    </b-form-group>

    <b-form-group
      :label="$t('clientSecret')"
      label-class="text-primary"
    >
      <b-input-group>
        <b-form-input
          v-model.trim="value.secret"
          :required="value.enabled"
        />
      </b-input-group>
    </b-form-group>

    <security
      v-model="value.security"
    />
  </div>
</template>

<script>
import Security from './ExternalSecurity'
export default {
  name: 'StandardExternalAuthProvider',

  i18nOptions: {
    namespaces: 'system.settings',
    keyPrefix: 'editor.external.standard',
  },

  components: {
    Security,
  },

  props: {
    value: {
      type: Object,
      required: true,
      default: () => ({}),
    },
  },

  computed: {
    // providerUsage provides a temporary checkbox implementation until something
    // more proper is done
    //
    // When checked (true), we consider it as wanting to use the API
    providerUsage: {
      get () {
        return ((this.value || {}).usage || []).includes('api')
      },
      set (value) {
        this.value.usage = value ? ['api'] : []
      },
    },
  },
}
</script>
