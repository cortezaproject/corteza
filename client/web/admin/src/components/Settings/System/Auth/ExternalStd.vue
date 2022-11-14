<template>
  <div>
    <b-form-group label-cols="3">
      <b-form-checkbox
        v-model="value.enabled"
        :value="true"
        :unchecked-value="false"
      >
        {{ $t('enabled') }}
      </b-form-checkbox>
    </b-form-group>
    <b-form-group
      :label="$t('clientKey')"
      label-cols="3"
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
      label-cols="3"
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
}
</script>
