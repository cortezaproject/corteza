<template>
  <div>
    <b-modal
      id="permissions-modal"
      :visible="showModal"
      size="xl"
      :title="translatedTitle"
      lazy
      scrollable
      no-fade
      body-class="d-flex flex-column p-0"
      class="h-100 overflow-hidden"
      @hide="onHide"
    >
      <b-row
        no-gutters
        class="bg-light border-bottom"
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
        class="bg-light"
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
            <c-input-select
              v-model="currentRoleID"
              data-test-id="select-user-list-roles"
              label="name"
              :disabled="!currentRoleID"
              :clearable="false"
              :options="roles"
              :get-option-key="getOptionRoleKey"
              :reduce="o => o.roleID"
              append-to-body
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
            class="pointer text-center text-primary mb-1"
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
            class="pointer text-center text-primary mb-1"
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
          class="d-none d-lg-flex border-left p-3 bg-light not-allowed"
        >
          <div
            class="d-flex flex-column align-items-center justify-content-between mt-4 w-100"
          >
            <h5
              v-for="r in e.rules"
              :key="r.operation"
              v-b-tooltip.noninteractive.hover="{ title: getRuleTooltip(r.access === 'unknown-context', !!e.userID), container: '#body' }"
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

      <template #modal-footer>
        <b-button
          data-test-id="button-cancel"
          variant="light"
          @click="onHide"
        >
          {{ labels.cancel }}
        </b-button>

        <c-button-submit
          data-test-id="button-save"
          :disabled="submitDisabled"
          :processing="submitting"
          :text="labels.save"
          @submit="onSubmit"
        />
      </template>
    </b-modal>

    <b-modal
      id="permissions-modal-eval"
      :title="labels.add.title"
      centered
      no-fade
    >
      <b-form-group
        :label="labels.add.role.label"
        label-class="text-primary"
        class="mb-0"
      >
        <c-input-select
          data-test-id="select-role"
          v-model="add.roleID"
          :options="roles"
          :get-option-key="getOptionRoleKey"
          label="name"
          multiple
          :disabled="!!add.userID"
          :placeholder="labels.add.role.placeholder"
        />
      </b-form-group>

      <b-form-group
        :label="labels.add.user.label"
        label-class="text-primary"
        class="mt-3 mb-0"
      >
        <c-input-select
          data-test-id="select-user"
          v-model="add.userID"
          :disabled="!!add.roleID.length"
          :options="userOptions"
          :get-option-label="getUserLabel"
          :get-option-key="getOptionUserKey"
          :placeholder="labels.add.user.placeholder"
          :filterable="false"
          @search="searchUsers"
        />
      </b-form-group>

      <template #modal-footer>
        <c-button-submit
          data-test-id="button-save"
          :disabled="!addEnabled"
          :processing="processing"
          :text="labels.add.save"
          @submit="onAddEval"
        />
      </template>
    </b-modal>
  </div>
</template>
<script lang="js">
import { modalOpenEventName, split } from './def.ts'
import CInputSelect from '../input/CInputSelect.vue'
import Rules from './form/Rules.vue'

export default {
  i18nOptions: {
    namespaces: 'permissions',
  },

  components: {
    Rules,
    CInputSelect,
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
      submitting: false,

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

      currentRoleID: undefined,

      evaluate: [],

      fetchedUsers: {},

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
      return !this.dirty || this.processing || this.submitting
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

    this.$root.$on(modalOpenEventName, this.loadModal)
  },

  beforeDestroy () {
    this.destroyEvents()
    this.setDefaultValues()
  },

  methods: {
    loadModal ({ resource, title, target, allSpecific }) {
      this.processing = true

      this.resource = resource
      this.title = title
      this.target = target
      this.allSpecific = allSpecific

      // <schema>::<backend-component-name>:...
      this.backendComponentName = resource.split(':')[2]

      this.fetchPermissions().then(() => {
        if (!this.roles.length) {
          return this.fetchRoles()
        } else if (this.currentRoleID) {
          return this.reEvaluatePermissions(this.currentRoleID)
        }
      }).finally(() => {
        this.processing = false
      })
    },

    onHide () {
      this.resource = undefined
      this.title = undefined
      this.target = undefined
    },

    onRoleChange (roleID) {
      this.processing = true

      this.fetchRules(roleID).finally(() => {
        this.processing = false
      })
    },

    onSubmit () {
      this.submitting = true

      const rules = this.collectChangedRules()
      const roleID = this.currentRoleID

      this.api.permissionsUpdate({ roleID, rules }).then(() => {
        this.reEvaluatePermissions(roleID)
        this.toastSuccess(this.$t('permissions:ui.notification.save.success'))
      }).catch(this.toastErrorHandler(this.$t('permissions:ui.notification.save.failed')))
        .finally(() => {
          setTimeout(() => {
            this.submitting = false
          }, 300)
        })
    },

    async fetchPermissions () {
      // clean loaded rules
      this.rules = []
      this.permissions = []

      return this.api.permissionsList().then((pp) => {
        this.permissions = this.filterPermissions(pp)
      })
    },

    async fetchRules (roleID) {
      return this.api.permissionsRead({ roleID, resource: this.resource }).then((rules) => {
        this.rules = this.normalizeRules(rules)
      })
    },

    async fetchRoles () {
      // Roles are always fetched from $SystemAPI.
      return this.$SystemAPI.roleList().then(({ set }) => {
        this.roles = set
          .filter(({ isBypass }) => !isBypass)
          .sort((a, b) => a.roleID.localeCompare(b.roleID))

        if (this.roles.length > 0) {
          this.currentRoleID = this.roles[0].roleID
          this.onRoleChange(this.currentRoleID)
        }
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

    async reEvaluatePermissions (roleID) {
      return this.fetchRules(roleID).then(() => {
        return Promise.all(this.evaluate.map(e => {
          let { roleID = [], userID } = e
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
    },

    searchUsers (query = '', loading) {
      loading(true)

      this.$SystemAPI.userList({ query, limit: 15 })
        .then(({ set }) => {
          this.userOptions = set.reduce((acc, { userID, name, username, email, handle }) => {
            if (!this.fetchedUsers[userID]) {
              this.fetchedUsers[userID] = name || username || email || `<@${userID}>`
            }

            acc.push(userID)

            return acc
          }, [])
        })
        .finally(() => {
          loading(false)
        })
    },

    getUserLabel (userID) {
      return this.fetchedUsers[userID]
    },

    onAddEval () {
      const userID = this.add.userID || {}
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

        this.$bvModal.hide('permissions-modal-eval')
      })
    },

    onHideEval (i) {
      this.evaluate.splice(i, 1)
    },

    getEvalName ({ userID, roleID }) {
      if (userID) {
        return [this.fetchedUsers[userID]]
      } else {
        return roleID.map(({ name }) => name)
      }
    },

    normalizeRules (rr, fallback = false) {
      const inherit = 'inherit'

      // merges roleRules (subset) with list of all permissions
      const findCurrent = ({ operation }) => {
        if (!rr) {
          return inherit
        }

        let { resolution, access = inherit } = (rr.find(r => r.operation === operation) || {})

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

    getOptionRoleKey ({ roleID }) {
      return roleID
    },

    setDefaultValues () {
      this.processing = false
      this.submitting = false
      this.backendComponentName = undefined
      this.resource = undefined
      this.title = undefined
      this.target = undefined
      this.allSpecific = false
      this.userOptions = []
      this.permissions = []
      this.rules = []
      this.roles = []
      this.currentRoleID = undefined
      this.evaluate = []
      this.add = {}
      this.fetchedUsers = {}
    },

    destroyEvents() {
      this.$root.$off(modalOpenEventName)
    },
  },
}
</script>

<style scoped lang="scss">
.not-allowed {
  cursor: not-allowed;
}

.pointer {
  cursor: pointer;
}
.rotate {
  transform: rotate(45deg);
}

.hide-eval:hover {
  .rotate {
    color: var(--primary) !important;
  }
}
</style>
