<template>
  <div>
    <b-card
      class="shadow h-100"
      header-class="p-0"
      body-class="overflow-auto p-0"
      header-bg-variant="white"
      footer-bg-variant="white"
    >
      <template
        v-if="loaded && canGrant"
        #header
      >
        <b-row
          no-gutters
          align-v="stretch"
          class="border-bottom"
        >
          <b-col
            cols="4"
            class="text-left p-3"
          >
            <small>
              {{ $t('ui.click-on-cell-to-allow') }}
              <br>
              {{ $t('ui.alt-click-to-deny') }}
            </small>
          </b-col>
          <b-col
            v-for="role in roles"
            :key="role.ID"
            class="d-flex flex-column align-items-center justify-content-center pointer hide-role border-left p-3 overflow-hidden"
            @click="onHideRole(role)"
          >
            <label
              v-for="(n, index) in role.name"
              :key="index"
              :title="n"
              class="pointer text-center text-primary text-break mb-1"
            >
              {{ n }}
            </label>
            <font-awesome-icon
              :icon="['fas', 'plus']"
              class="text-light rotate"
            />
          </b-col>
          <b-col
            v-if="roles.length < 8"
            v-b-modal.add
            class="d-flex flex-column align-items-center justify-content-center border-left p-3 overflow-hidden"
          >
            <label
              class="pointer text-center text-primary text-break mb-1"
            >
              {{ $t('ui.add.label') }}
            </label>
            <font-awesome-icon
              :icon="['fas', 'plus']"
              class="text-success"
            />
          </b-col>
        </b-row>
      </template>

      <div
        v-if="!loaded || !canGrant"
        class="d-flex align-items-center justify-content-center h-100 pb-4"
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

      <template v-else>
        <div
          v-for="(type, i) in sortedPermissions"
          :key="type"
        >
          <b-row
            class="bg-light border-bottom text-primary"
            align-v="stretch"
            no-gutters
          >
            <b-col
              cols="4"
              align-self="center"
              class="p-3 text-left"
            >
              <span class="h-100 h6 mb-0">
                {{ getTranslation(type) }}
              </span>
            </b-col>
            <b-col
              v-for="role in roles"
              :key="role.ID"
              class="d-flex align-items-center justify-content-center overflow-hidden p-3 text-center border-left not-allowed"
            >
              <p
                v-if="i === 0"
                :title="$t(`ui.${role.mode === 'edit' ? 'edit' : 'evaluate'}.title`)"
                class="mb-0"
              >
                {{ $t(`ui.${role.mode === 'edit' ? 'edit' : 'evaluate'}.title`) }}
              </p>
            </b-col>
            <b-col
              class="p-3 border-left not-allowed"
            />
          </b-row>
          <b-row
            v-for="operation in permissions[type].ops"
            :key="operation"
            no-gutters
          >
            <b-col
              class="border-bottom text-left p-3"
              cols="4"
            >
              <span :title="getTranslation(type, operation)">{{ getTranslation(type, operation) }}</span>
            </b-col>
            <b-col
              v-for="role in roles"
              :key="role.ID"
              :title="getRuleTooltip(checkRule(role.ID, permissions[type].any, operation, 'unknown-context'), !!role.userID)"
              class="d-flex align-items-center justify-content-center border-bottom border-left p-3 pointer active-cell h5 mb-0"
              :class="{
                'not-allowed bg-extra-light': role.mode === 'eval',
                'bg-warning': checkChange(role.ID, permissions[type].any, operation)
              }"
              @click="role.mode === 'edit' ? ruleChange($event, role.ID, permissions[type].any, operation) : undefined"
            >
              <font-awesome-icon
                v-if="checkRule(role.ID, permissions[type].any, operation, 'unknown-context')"
                :icon="['fas', 'question']"
                class="text-secondary"
              />
              <font-awesome-icon
                v-else-if="checkRule(role.ID, permissions[type].any, operation, 'allow')"
                :icon="['fas', 'check']"
                class="text-success"
              />
              <font-awesome-icon
                v-else
                :icon="['fas', 'times']"
                class="text-danger"
              />
            </b-col>
            <b-col
              v-if="roles.length < 8"
              class="border-bottom border-left p-3 not-allowed bg-extra-light"
            />
          </b-row>
        </div>
      </template>

      <template
        v-if="loaded && canGrant"
        #footer
      >
        <c-submit-button
          class="float-right"
          :processing="processing"
          :success="success"
          @submit="onSubmit"
        >
          {{ $t('ui.save') }}
        </c-submit-button>
      </template>
    </b-card>

    <b-modal
      id="add"
      :title="$t('ui.edit-or-eval')"
      centered
      ok-only
      :ok-title="$t('ui.add.save')"
      :ok-disabled="!addEnabled"
      @ok="onAdd"
    >
      <b-form-group>
        <b-form-radio-group
          v-model="add.mode"
          :options="modeOptions"
          buttons
          button-variant="outline-primary"
          class="mode rounded w-100"
        />
      </b-form-group>

      <p>
        {{ addModeDescription }}
      </p>

      <b-form-group
        :label="$t('ui.add.role.label')"
        label-class="text-primary"
        class="mb-0"
      >
        <vue-select
          key="roleID"
          v-model="add.roleID"
          :options="availableRoles"
          :multiple="add.mode === 'eval'"
          label="name"
          clearable
          :disabled="add.mode === 'eval' && !!add.userID"
          :placeholder="$t('ui.add.role.placeholder')"
          class="bg-white"
        />
      </b-form-group>

      <b-form-group
        v-if="add.mode === 'eval'"
        :label="$t('ui.add.user.label')"
        label-class="text-primary"
        class="mt-3 mb-0"
      >
        <vue-select
          key="userID"
          v-model="add.userID"
          :disabled="!!add.roleID.length"
          :options="userOptions"
          :get-option-label="getUserLabel"
          label="name"
          clearable
          :placeholder="$t('ui.add.user.placeholder')"
          class="bg-white"
          @search="searchUsers"
        />
      </b-form-group>
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

    allRoles: {
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
      add: {
        // Either edit or eval
        mode: 'edit',
        roleID: [],
        userID: undefined,
      },

      modeOptions: [
        { text: 'Edit', value: 'edit' },
        { text: 'Evaluate', value: 'eval' },
      ],

      userOptions: [],

      evaluatedPermissions: undefined,

      newRole: null,
      permissionChanges: [],
    }
  },

  computed: {
    editableRoles () {
      return this.roles.filter(({ mode }) => mode !== 'eval').map(({ roleID }) => roleID.roleID)
    },

    availableRoles () {
      if (this.add.mode === 'edit') {
        return this.allRoles.filter(({ roleID, isBypass }) => !isBypass && !this.editableRoles.includes(roleID))
      } else if (this.add.mode === 'eval') {
        return this.allRoles
      }

      return []
    },

    sortedPermissions () {
      return Object.keys(this.permissions).sort()
    },

    addModeDescription () {
      return this.add.mode === 'edit' ? this.$t('ui.add.edit.description') : this.$t('ui.add.evaluate.description')
    },

    addEnabled () {
      const { mode, roleID = [], userID } = this.add

      if (mode === 'edit') {
        return Array.isArray(roleID) ? roleID.length : roleID
      } else if (mode === 'eval') {
        return (roleID && roleID.length) || userID
      }

      return false
    },
  },

  watch: {
    'add.mode': {
      handler (mode) {
        this.add.roleID = mode === 'eval' ? [] : undefined
        this.add.userID = undefined
      },
    },
  },

  mounted () {
    this.searchUsers('', () => {})
  },

  methods: {
    checkRule (ID, res, op, access) {
      const key = `${op}@${res}`
      return (this.rolePermissions.find(r => r.ID === ID) || { rules: {} }).rules[key] === access
    },

    checkChange (ID, res, op) {
      const key = `${op}@${res}`
      const current = (this.rolePermissions.find(r => r.ID === ID) || { rules: {} }).rules[key]
      const initial = (this.permissionChanges.find(r => r.ID === ID) || { rules: {} }).rules[key]

      if (initial) {
        return current !== initial
      } else {
        return false
      }
    },

    ruleChange (event, ID, res, op) {
      const key = `${op}@${res}`
      let access = (this.rolePermissions.find(r => r.ID === ID) || { rules: {} }).rules[key]

      // Keep track of permission changes, record initial value before it changes
      if (!(this.permissionChanges.find(r => r.ID === ID) || { rules: {} }).rules[key]) {
        this.permissionChanges.push({ ID, rules: { } })

        if (!access) {
          access = 'inherit'
        }
        this.$set(this.permissionChanges.find(r => r.ID === ID).rules, key, access)
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

      this.$set(this.rolePermissions.find(r => r.ID === ID).rules, key, access)
    },

    searchUsers (query = '', loading) {
      loading(true)

      this.$SystemAPI.userList({ query, limit: 15 })
        .then(({ set }) => {
          this.userOptions = set.map(m => Object.freeze(m))
        })
        .finally(() => {
          loading(false)
        })
    },

    getUserLabel ({ userID, email, name, username }) {
      return name || username || email || `<@${userID}>`
    },

    getTranslation (resource, operation = '') {
      resource = _.kebabCase(resource.split(':')[3]) || 'component'

      if (operation) {
        return this.$t(`resources.${this.component}.${resource}.operations.${operation}.title`)
      } else {
        return this.$t(`resources.${this.component}.${resource}.label`)
      }
    },

    getRuleTooltip (isUnknown = false, isUser) {
      if (!isUnknown) {
        return ''
      }

      return this.$t(`ui.tooltip.unknown-context.${isUser ? 'user' : 'role'}`)
    },

    onSubmit () {
      this.$emit('submit', this.rolePermissions)
      this.permissionChanges = []
    },

    onAdd () {
      this.$emit('add', this.add)
      this.add = {
        mode: 'edit',
        roleID: [],
        userID: undefined,
      }
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
  .rotate {
    color: $dark !important;
  }
}
</style>

<style lang="scss">
.mode {
  .btn {
    background-color: #E4E9EF;
    border: none;
  }

  .btn:nth-child(2), .btn:nth-child(3) {
    margin-left: 0.2rem !important;
  }
}
</style>
