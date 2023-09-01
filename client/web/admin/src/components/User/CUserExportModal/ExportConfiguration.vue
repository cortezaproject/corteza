<template>
  <b-card
    header-bg-variant="white"
    footer-bg-variant="white"
    body-class="pb-0"
  >
    <b-form-group>
      <b-form-checkbox
        v-model="inclRoleMembership"
        data-test-id="checkbox-include-role-membership"
        :value="true"
        :unchecked-value="false"
      >
        {{ $t('export.inclRoleMembership') }}
      </b-form-checkbox>
    </b-form-group>

    <b-form-group
      :description="!inclRoleMembership
        ? $t('export.membershipRequiredLabel')
        : ''
      "
      class="mb-0"
    >
      <b-form-checkbox
        v-model="inclRoles"
        data-test-id="checkbox-include-roles"
        :value="true"
        :unchecked-value="false"
        :disabled="!inclRoleMembership"
      >
        {{ $t('export.inclRoles') }}
      </b-form-checkbox>
    </b-form-group>

    <div slot="footer">
      <b-button
        data-test-id="button-export"
        variant="dark"
        class="float-right"
        @click="nextStep"
      >
        {{ $t('export.export') }}
      </b-button>
    </div>
  </b-card>
</template>

<script>
export default {
  i18nOptions: {
    namespaces: 'system.users',
  },

  name: 'CUserExportConfiguration',

  data () {
    return {
      inclRoleMembership: false,
      inclRoles: false,
    }
  },

  methods: {
    nextStep () {
      const rtr = {
        inclRoleMembership: this.inclRoleMembership,
        inclRoles: this.inclRoles,
      }

      this.$emit('configured', rtr)
    },
  },
}
</script>
