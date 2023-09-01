<template>
  <div>
    <hr>
    <b-form-group
      :label="$t('permitted-roles.label')"
      label-cols="3"
      label-class="text-primary"
      class="mb-0"
    >
      <c-role-picker
        v-model="value.permittedRoles"
        class="mb-3"
      >
        <template #description>
          {{ $t('permitted-roles.description') }}
        </template>
      </c-role-picker>
    </b-form-group>

    <b-form-group
      :label="$t('prohibited-roles.label')"
      label-cols="3"
      label-class="text-primary"
      class="mb-0"
    >
      <c-role-picker
        v-model="value.prohibitedRoles"
        class="mb-3"
      >
        <template #description>
          {{ $t('prohibited-roles.description') }}
        </template>
      </c-role-picker>
    </b-form-group>

    <b-form-group
      :label="$t('forced-roles.label')"
      label-cols="3"
      label-class="text-primary"
      class="mb-0"
    >
      <c-role-picker
        v-model="value.forcedRoles"
        class="mb-3"
      >
        <template #description>
          {{ $t('forced-roles.description') }}
        </template>
      </c-role-picker>
    </b-form-group>
  </div>
</template>

<script>
import CRolePicker from 'corteza-webapp-admin/src/components/CRolePicker'
export default {
  name: 'StandardExternalAuthProvider',

  i18nOptions: {
    namespaces: 'system.settings',
    keyPrefix: 'editor.external.security',
  },

  components: {
    CRolePicker,
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
