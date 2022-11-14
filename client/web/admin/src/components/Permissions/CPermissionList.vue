<template>
  <div>
    <b-card
      class="shadow-sm"
      body-class="py-0 px-3"
      header-bg-variant="white"
      footer-bg-variant="white"
    >
      <div
        v-if="loaded && canGrant"
      >
        <b-row
          class="text-center"
        >
          <b-col
            class="border-bottom py-2"
            cols="4"
          />
          <b-col
            v-for="role in roles"
            :key="role.roleID"
            class="border-bottom border-left py-2 overflow-hidden text-nowrap"
          >
            <div
              class="text-truncate"
            >
              {{ role.name }}
            </div>
            <b-button
              variant="link"
              class="hide-role text-light text-decoration-none py-0"
              @click="onHideRole(role)"
            >
              <font-awesome-icon
                :icon="['fas', 'plus']"
                class="rotate"
              />
            </b-button>
          </b-col>
          <b-col
            v-if="roles.length < 8"
            class="border-bottom border-left py-2 overflow-hidden text-nowrap"
          >
            <div
              class="text-truncate"
            >
              {{ $t('ui.role.add') }}
            </div>
            <b-button
              v-b-modal.addRole
              variant="link"
              class="text-primary text-truncate text-decoration-none py-0"
            >
              <font-awesome-icon
                :icon="['fas', 'plus']"
                class="d-block ml-auto mr-auto"
              />
            </b-button>
          </b-col>
        </b-row>
        <div
          v-for="type in sortedPermissions"
          :key="type"
        >
          <b-row
            class="bg-secondary"
          >
            <b-col
              cols="4"
              class="py-2 text-left font-weight-bold"
            >
              {{ getTranslation(type) }}
            </b-col>
            <b-col
              v-for="role in roles"
              :key="role.roleID"
              class="py-2 not-allowed"
            />
          </b-row>
          <b-row
            v-for="operation in permissions[type].ops"
            :key="operation"
            class="text-center"
          >
            <b-col
              class="border-bottom text-left py-2 text-truncate"
              cols="4"
            >
              <span :title="getTranslation(type, operation)">{{ getTranslation(type, operation) }}</span>
            </b-col>
            <b-col
              v-for="role in roles"
              :key="role.roleID"
              class="border-bottom border-left py-2 pointer active-cell"
              :class="{
                'not-allowed bg-light': role.roleID.includes('-'),
                'bg-warning': checkChange(role.roleID, permissions[type].any, operation)
              }"
              @click="ruleChange($event, role.roleID, permissions[type].any, operation)"
            >
              <font-awesome-icon
                v-if="checkRule(role.roleID, permissions[type].any, operation, 'allow')"
                :icon="['fas', 'check']"
                class="text-success"
              />
              <font-awesome-icon
                v-if="checkRule(role.roleID, permissions[type].any, operation, 'deny')"
                :icon="['fas', 'times']"
                class="text-danger"
              />
            </b-col>
            <b-col
              v-if="roles.length < 8"
              class="border-bottom border-left py-2 not-allowed bg-light"
            />
          </b-row>
        </div>
      </div>

      <div
        v-else
        class="text-center m-5"
      >
        <div
          v-if="!loaded"
        >
          <b-spinner
            class="align-middle m-5"
          />
          <div>
            {{ $t('ui.loading') }}
          </div>
        </div>
        <div
          v-else-if="!canGrant"
          class="text-danger"
        >
          {{ $t('ui.not-allowed') }}
        </div>
      </div>

      <template
        v-if="loaded && canGrant"
        #footer
      >
        <small class="float-left text-primary">
          {{ $t('ui.click-on-cell-to-allow') }}
          <br>
          {{ $t('ui.alt-click-to-deny') }}
        </small>
        <c-submit-button
          class="float-right"
          :processing="processing"
          :success="success"
          @submit="onSubmit"
        >
          {{ $t('ui.save-changes') }}
        </c-submit-button>
      </template>
    </b-card>

    <b-modal
      id="addRole"
      :title="$t('ui.role.add')"
      :ok-only="true"
      :ok-title="$t('ui.role.add-new')"
      @ok="onAddRole"
    >
      <vue-select
        key="roleID"
        v-model="newRole"
        :options="rolesExcluded"
        label="name"
        class="bg-white"
        :placeholder="$t('ui.role.no-role-selected')"
      />
    </b-modal>
  </div>
</template>

<script>
import CSubmitButton from 'corteza-webapp-admin/src/components/CSubmitButton'
import { VueSelect } from 'vue-select'
import _ from 'lodash'

export default {
  i18nOptions: {
    namespaces: 'permissions',
  },

  components: {
    CSubmitButton,
    VueSelect,
  },

  props: {
    roles: {
      type: Array,
      required: true,
    },

    rolesExcluded: {
      type: Array,
      required: true,
    },

    permissions: {
      type: Object,
      required: true,
    },

    rolePermissions: {
      type: Array,
      required: true,
    },

    canGrant: {
      type: Boolean,
      value: false,
    },

    loaded: {
      type: Boolean,
      value: false,
    },

    processing: {
      type: Boolean,
      value: false,
    },

    success: {
      type: Boolean,
      value: false,
    },

    component: {
      type: String,
      required: true,
    },
  },

  data () {
    return {
      newRole: null,
      permissionChanges: [],
    }
  },

  computed: {
    sortedPermissions () {
      return Object.keys(this.permissions).sort()
    },
  },

  methods: {
    checkRule (roleID, res, op, access) {
      const key = `${op}@${res}`
      return (this.rolePermissions.find(r => r.roleID === roleID) || { rules: {} }).rules[key] === access
    },

    checkChange (roleID, res, op) {
      const key = `${op}@${res}`
      const current = (this.rolePermissions.find(r => r.roleID === roleID) || { rules: {} }).rules[key]
      const initial = (this.permissionChanges.find(r => r.roleID === roleID) || { rules: {} }).rules[key]

      if (initial) {
        return current !== initial
      } else {
        return false
      }
    },

    ruleChange (event, roleID, res, op) {
      const key = `${op}@${res}`
      let access = (this.rolePermissions.find(r => r.roleID === roleID) || { rules: {} }).rules[key]

      // Keep track of permission changes, record initial value before it changes
      if (!(this.permissionChanges.find(r => r.roleID === roleID) || { rules: {} }).rules[key]) {
        this.permissionChanges.push({ roleID, rules: { } })

        if (!access) {
          access = 'inherit'
        }
        this.$set(this.permissionChanges.find(r => r.roleID === roleID).rules, key, access)
      }

      if (event.altKey) {
        if (access === 'deny') {
          access = 'inherit'
        } else {
          access = 'deny'
        }
      } else {
        if (access === 'allow') {
          access = 'inherit'
        } else {
          access = 'allow'
        }
      }

      this.$set(this.rolePermissions.find(r => r.roleID === roleID).rules, key, access)
    },

    getTranslation (resource, operation = '') {
      resource = _.kebabCase(resource.split(':')[3]) || 'component'

      if (operation) {
        return this.$t(`resources.${this.component}.${resource}.operations.${operation}.title`)
      } else {
        return this.$t(`resources.${this.component}.${resource}.label`)
      }
    },

    onSubmit () {
      this.$emit('submit', this.rolePermissions)
      this.permissionChanges = []
    },

    onAddRole () {
      const { roleID } = (this.newRole || {})
      if (roleID) {
        this.$emit('add', this.newRole)
      }
      this.newRole = null
    },

    onHideRole (role) {
      this.$emit('hide', role)
    },
  },
}
</script>
<style lang="scss" scoped>
.pointer {
  cursor: pointer;
}
.not-allowed {
  cursor: not-allowed;
}
.active-cell:hover {
  background-color: #F3F3F5;
}
.rotate {
  transform: rotate(45deg);
}
.hide-role:hover {
  color: $dark !important;
}
</style>
