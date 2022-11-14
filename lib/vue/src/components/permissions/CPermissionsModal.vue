<template>
  <div>
    <b-modal
      id="permissions-modal"
      :visible="showModal"
      size="xl"
      :title="translatedTitle"
      lazy
      scrollable
      :ok-disabled="submitDisabled"
      :ok-title="labels.save"
      :cancel-title="labels.cancel"
      cancel-variant="light"
      body-class="d-flex flex-column p-0"
      class="h-100 overflow-hidden"
      @hide="onHide"
      @ok="onSubmit"
    >
      <b-row
        no-gutters
        class="bg-extra-light border-bottom"
      >
        <b-col
          lg="4"
          class="p-3"
        >
          {{ labels.edit.description }}
        </b-col>
        <b-col
          class="d-none d-lg-block border-left p-3"
        >
          {{ labels.evaluate.description }}
        </b-col>
      </b-row>

      <b-row
        no-gutters
        align-v="stretch"
        class="bg-extra-light"
      >
        <b-col
          lg="4"
          class="d-flex align-items-center p-3 border-bottom"
        >
          <b-form-group
            class="mb-0 w-100"
            :label="labels.edit.label"
            label-class="text-primary"
          >
            <vue-select
              v-model="currentRole"
              data-test-id="select-user-list-roles"
              key="roleID"
              label="name"
              :clearable="false"
              :options="roles"
              append-to-body
              class="h-100 bg-white"
              @input="onRoleChange"
            />
          </b-form-group>
        </b-col>

        <b-col
          v-for="(e, i) in evaluate"
          data-test-id="icon-remove"
          :key="i"
          lg="2"
          class="pointer hide-eval border-bottom d-none d-lg-flex flex-column align-items-center justify-content-center overflow-hidden border-left p-3"
          @click="onHideEval(i)"
        >
          <label
            v-for="(n, index) in getEvalName(e)"
            :key="index"
            :title="n"
            class="pointer text-center text-primary text-break mb-1"
          >
            {{ n }}
          </label>
          <font-awesome-icon
            :icon="['fas', 'plus']"
            class="text-secondary rotate mt-1"
          />
        </b-col>

        <b-col
          v-if="evaluate.length < 4"
          v-b-modal.permissions-modal-eval
          data-test-id="icon-add"
          class="d-none d-lg-flex pointer border-bottom flex-column align-items-center justify-content-center overflow-hidden border-left p-3"
        >
          <label
            class="pointer text-center text-primary text-break mb-1"
          >
            {{ labels.add.label }}
          </label>
          <font-awesome-icon
            :icon="['fas', 'plus']"
            class="text-success d-block mx-auto mt-1"
          />
        </b-col>
      </b-row>

      <div
        v-if="processing"
        class="d-flex flex-column align-items-center justify-content-center h-100 py-4"
        style="min-height: 50vh;"
      >
        <b-spinner />
        <div>
          {{ labels.loading }}
        </div>
      </div>

      <b-row
        v-else
        no-gutters
      >
        <b-col
          lg="4"
          class="p-3"
        >
          <rules
            :rules.sync="rules"
          />
        </b-col>
        <b-col
          v-for="(e, i) in evaluate"
          :key="i"
          lg="2"
          class="d-none d-lg-flex border-left p-3 bg-extra-light not-allowed"
        >
          <div
            class="d-flex flex-column align-items-center justify-content-between mt-4 w-100"
          >
            <h5
              v-for="r in e.rules"
              :key="r.operation"
              :title="getRuleTooltip(r.access === 'unknown-context', !!e.userID)"
              class="text-center mb-1 mt-2 w-100"
            >
              <font-awesome-icon
                v-if="r.access === 'unknown-context'"
                :icon="['fas', 'question']"
                class="text-secondary"
              />
              <font-awesome-icon
                v-else-if="r.access === 'allow'"
                :icon="['fas', 'check']"
                class="text-success"
              />
              <font-awesome-icon
                v-else
                :icon="['fas', 'times']"
                class="text-danger"
              />
            </h5>
          </div>
        </b-col>

        <b-col
          class="d-none d-lg-block pt-4 border-left"
        />
      </b-row>
    </b-modal>

    <b-modal
      id="permissions-modal-eval"
      :title="labels.add.title"
      centered
      ok-only
      :ok-title="labels.add.save"
      :ok-disabled="!addEnabled"
      @ok="onAddEval"
    >
      <b-form-group
        :label="labels.add.role.label"
        label-class="text-primary"
        class="mb-0"
      >
        <vue-select
          data-test-id="select-role"
          key="roleID"
          v-model="add.roleID"
          :options="roles"
          label="name"
          multiple
          clearable
          :disabled="!!add.userID"
          :placeholder="labels.add.role.placeholder"
          class="bg-white"
        />
      </b-form-group>

      <b-form-group
        :label="labels.add.user.label"
        label-class="text-primary"
        class="mt-3 mb-0"
      >
        <vue-select
          data-test-id="select-user"
          key="userID"
          v-model="add.userID"
          :disabled="!!add.roleID.length"
          :options="userOptions"
          :get-option-label="getUserLabel"
          label="name"
          clearable
          :placeholder="labels.add.user.placeholder"
          class="bg-white"
          @search="searchUsers"
        />
      </b-form-group>
    </b-modal>
  </div>
</template>
<script lang="js">
import { modalOpenEventName, split } from './def.ts'
import { VueSelect } from 'vue-select'
import Rules from './form/Rules.vue'

export default {
  i18nOptions: {
    namespaces: 'permissions',
  },

  components: {
    Rules,
    VueSelect,
  },

  props: {
    labels: {
      required: false,
      type: Object,
      default: () => ({}),
    },
  },

  data () {
    return {
      processing: false,

      backendComponentName: undefined,

      resource: undefined,
      title: undefined,
      target: undefined,
      allSpecific: false,

      userOptions: [],

      // List of available permissions (for a specific resource)
      permissions: [],

      // List of rules for the current role
      rules: [],

      // List of all available roles
      roles: [],

      // Current role object
      currentRole: undefined,

      evaluate: [],

      add: {
        roleID: [],
        userID: undefined,
      },
    }
  },

  computed: {
    api () {
      const s = this.backendComponentName
      return s ? this['$' + s.charAt(0).toUpperCase() + s.slice(1) + 'API'] : undefined
    },

    showModal () {
      return !!(this.resource && this.api)
    },

    dirty () {
      return this.collectChangedRules().length > 0
    },

    submitDisabled () {
      return !this.dirty
    },

    addEnabled () {
      const { roleID = [], userID } = this.add
      return (roleID && roleID.length) || userID
    },

    translatedTitle () {
      if (this.resource) {
        const { i18nPrefix } = split(this.resource)

        let target
        if (this.allSpecific) {
          target = this.$t(`permissions:${i18nPrefix}.all-specific`, { target: this.title, interpolation: { escapeValue: false } })
        } else if (this.title) {
          target = this.$t(`permissions:${i18nPrefix}.specific`, { target: this.title, interpolation: { escapeValue: false } })
        } else {
          target = this.$t(`permissions:${i18nPrefix}.all`,)
        }

        return this.$t('permissions:ui.set-for', { target: target, interpolation: { escapeValue: false } })
      }

      return undefined
    },
  },

  mounted () {
    this.searchUsers('', () => {})

    this.$root.$on(modalOpenEventName, ({ resource, title, target, allSpecific }) => {
      this.resource = resource
      this.title = title
      this.target = target
      this.allSpecific = allSpecific

      // <schema>::<backend-component-name>:...
      this.backendComponentName = resource.split(':')[2]

      this.fetchPermissions().then(() => {
        if (!this.roles.length) {
          return this.fetchRoles()
        } else if (this.currentRole) {
          const { roleID } = this.currentRole
          this.processing = true
  
          return this.fetchRules(roleID).then(() => {
            return Promise.all(this.evaluate.map(e => {
              const { userID } = e.userID || {}
              let { roleID = [] } = e
              roleID = roleID.map(({ roleID }) => roleID)
    
              return this.evaluatePermissions({ roleID, userID }).then(rules => {
                return {
                  ...e,
                  rules,
                }
              })
            })).then(evaluate => {
              this.evaluate = evaluate
            })
          })
        }
      })

    })
  },

  destroyed () {
    this.$root.$off(modalOpenEventName)
  },

  methods: {
    onHide () {
      this.clear()
    },

    clear () {
      this.resource = undefined
      this.title = undefined
      this.target = undefined
    },

    onRoleChange ({ roleID }) {
      this.fetchRules(roleID)
    },

    onSubmit () {
      this.processing = true

      const rules = this.collectChangedRules()
      const { roleID } = this.currentRole

      this.api.permissionsUpdate({ roleID, rules }).then(() => {
        this.fetchRules(roleID)
      }).finally(() => {
        this.processing = false
      })
    },

    async fetchPermissions () {
      this.processing = true

      // clean loaded rules
      this.rules = []
      this.permissions = []

      return this.api.permissionsList().then((pp) => {
        this.permissions = this.filterPermissions(pp)
      }).finally(() => {
        this.processing = false
      })
    },

    async fetchRules (roleID) {
      this.processing = true

      return this.api.permissionsRead({ roleID, resource: this.resource }).then((rules) => {
        this.rules = this.normalizeRules(rules)
      }).finally(() => {
        this.processing = false
      })
    },

    async fetchRoles () {
      this.processing = true
      // Roles are always fetched from $SystemAPI.
      return this.$SystemAPI.roleList().then(({ set }) => {
        this.roles = set
          .filter(({ isBypass }) => !isBypass)
          .sort((a, b) => a.roleID.localeCompare(b.roleID))

        if (this.roles.length > 0) {
          this.currentRole = this.roles[0]
          this.onRoleChange(this.currentRole)
        }
      }).finally(() => {
        this.processing = false
      })
    },

    async evaluatePermissions ({ resource = this.resource, roleID, userID }) {
      this.processing = true

      return this.api.permissionsTrace({ resource, roleID, userID })
        .then(this.normalizeRules, true)
        .finally(() => {
          this.processing = false
        })
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

    onAddEval () {
      const { userID } = this.add.userID || {}
      let { roleID = [] } = this.add
      roleID = roleID.map(({ roleID }) => roleID)

      this.evaluatePermissions({ roleID, userID }).then(rules => {
        this.evaluate.push({
          ...this.add,
          rules,
        })

        this.add = {
          roleID: [],
          userID: undefined,
        }
      })
    },

    onHideEval (i) {
      this.evaluate.splice(i, 1)
    },

    getEvalName ({ roleID, userID }) {
      if (userID) {
        const { name, username, email, handle } = userID
        return [name || username || email || handle || userID || '']
      } else {
        return roleID.map(({ name }) => name)
      }
    },

    normalizeRules (rr, fallback = false) {
      const inherit = 'inherit'

      // merges roleRules (subset) with list of all permissions
      const findCurrent = ({ resource, operation }) => {
        if (!rr) {
          return inherit
        }

        let { resolution, access = inherit } = (rr.find(r => r.resource === resource && r.operation === operation) || {})

        if (resolution === 'unknown-context') {
          access = 'unknown-context'
        } else if (fallback && access === inherit) {
          access = 'deny'
        }

        return access
      }

      return this.permissions.map((p) => {
        const current = findCurrent(p)
        return { ...p, access: current, current }
      })
    },

    // Removes unneeded permissions (ones that do not match resource prop)
    // and translates the rest
    filterPermissions (pp) {
      const [ resourceType ] = this.resource.split('/', 2)
      return pp
        .filter(({ type }) => resourceType === type)
        .map(({ type, op: operation }) => {
          return {
            ...this.describePermission({ resource: type, operation }),
            operation,
            // override resource-type with the actual resource-ID
            resource: this.resource
          }
        })
    },

    collectChangedRules () {
      return this.rules.filter(r => r.access !== r.current).map(({ resource, operation, access }) => {
        return { resource, operation, access }
      })
    },

    describePermission ({ resource, operation }) {
      const i18nPrefix = split(resource).i18nPrefix + `.operations.${operation}`

      let title = ''
      if (this.allSpecific) {
        title = this.$t(`permissions:${i18nPrefix}.all-specific`, { target: this.target, interpolation: { escapeValue: false } })
      } else if (this.target) {
        title = this.$t(`permissions:${i18nPrefix}.specific`, { target: this.target, interpolation: { escapeValue: false } })
      } else {
        title = this.$t(`permissions:${i18nPrefix}.title`)
      }

      return {
        title,
        description: this.$t(`permissions:${i18nPrefix}.description`),
      }
    },

    getRuleTooltip (isUnknown = false, isUser) {
      if (!isUnknown) {
        return ''
      }

      return this.$t(`permissions:ui.tooltip.unknown-context.${isUser ? 'user' : 'role'}`)
    },
  },
}
</script>

<style scoped lang="scss">
.not-allowed {
  cursor: not-allowed;
}
.bg-extra-light {
  background-color: #F3F5F7;
}
.pointer {
  cursor: pointer;
}
.rotate {
  transform: rotate(45deg);
}

.hide-eval:hover {
  .rotate {
    color: #162425 !important;
  }
}
</style>

<style lang="scss">
#permissions-modal, #permissions-modal-eval {
  .v-select {
    min-width: 100%;

    .vs__selected-options {
      flex-wrap: wrap;
    }
  }

}
</style>
