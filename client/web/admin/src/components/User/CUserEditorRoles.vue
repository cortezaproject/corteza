<template>
  <b-card
    data-test-id="card-role-membership"
    header-bg-variant="white"
    footer-bg-variant="white"
    footer-class="d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <b-form
      @submit.prevent="$emit('submit')"
    >
      <c-role-picker v-model="value" />
    </b-form>

    <template #header>
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>

    <template #footer>
      <c-button-submit
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="$emit('submit')"
      />
    </template>
  </b-card>
</template>

<script>
import CRolePicker from 'corteza-webapp-admin/src/components/CRolePicker'

export default {
  i18nOptions: {
    namespaces: 'system.users',
    keyPrefix: 'editor.roles',
  },

  components: {
    CRolePicker,
  },

  props: {
    value: {
      type: Array,
      required: true,
      default: () => [],
    },

    processing: {
      type: Boolean,
      value: false,
    },

    success: {
      type: Boolean,
      value: false,
    },
  },

  computed: {
    roles: {
      get () {
        return this.currentRoles
      },

      set (roles) {
        this.$emit('update:current-roles', roles)
      },
    },
  },
}
</script>
