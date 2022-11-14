<template>
  <b-container
    fluid
    class="d-flex flex-column h-100 py-3"
  >
    <c-content-header
      :title="$t('ui.title.system')"
    />

    <c-permission-list
      :roles="sortedRoles"
      :all-roles="allRoles"
      :permissions="permissions"
      :role-permissions="rolePermissions"
      :can-grant="canGrant"
      :loaded="isLoaded"
      :processing="permission.processing"
      :success="permission.success"
      component="system"
      class="flex-fill"
      @submit="onSubmit"
      @add="addRole"
      @hide="hideRole"
    />
  </b-container>
</template>

<script>
import permissionHelpers from 'corteza-webapp-admin/src/mixins/permissionHelpers'
import CPermissionList from 'corteza-webapp-admin/src/components/Permissions/CPermissionList'
import { mapGetters } from 'vuex'

export default {
  i18nOptions: {
    namespaces: 'permissions',
  },

  components: {
    CPermissionList,
  },

  mixins: [
    permissionHelpers,
  ],

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canGrant () {
      return this.can('system/', 'grant')
    },
  },

  created () {
    /**
     * With this, we tell permissionHelpers mixin
     * what API it should use
     */
    this.api = this.$SystemAPI

    this.fetchPermissions()
  },
}
</script>
