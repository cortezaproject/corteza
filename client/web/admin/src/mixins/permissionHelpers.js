import editorHelpers from './editorHelpers'

const systemRoles = ['1', '2']
const lsKey = 'permissionList.roles'

/**
 * Determines what roles should be included in the table
 * @returns {Array<String>}
 */
function getIncludedRoles () {
  const udRoles = JSON.parse(localStorage.getItem(lsKey) || '[]')
  return udRoles
}

/**
 * Filters out system roles and stores to localStorage
 * @param {Array<String>} roles Roles to store
 */
function setIncludedRoles (roles = []) {
  roles = roles.filter(r => !systemRoles.includes(r))
  localStorage.setItem(lsKey, JSON.stringify(roles))
}

export default {
  mixins: [
    editorHelpers,
  ],

  data () {
    return {
      roles: [],
      allRoles: [],
      rolePermissions: [],
      permissions: {},
      effective: {},

      api: this.$SystemAPI,
      resources: [],

      loaded: {
        roles: false,
        permissions: false,
      },
      permission: {
        processing: false,
        success: false,
      },
    }
  },

  computed: {
    isLoaded () {
      return this.loaded.roles && this.loaded.permissions
    },

    sortedRoles () {
      return this.roles.sort((a, b) => a.mode.localeCompare(b.mode))
    },
  },

  methods: {
    fetchRoles () {
      this.incLoader()

      return this.$SystemAPI.roleList()
        .then(({ set }) => {
          this.allRoles = set
          this.rolePermissions = []

          // We read permissions for included roles
          return Promise.all(getIncludedRoles().map(({ mode, name, roleID, userID }) => {
            if (mode === 'edit') {
              return this.readPermissions({ name, roleID })
            } else {
              return this.evaluatePermissions({ name, roleID, userID })
            }
          }))
        })
        .catch(this.toastErrorHandler(this.$t('notification:user.roles.error')))
        .finally(() => {
          this.loaded.roles = true
          this.decLoader()
        })
    },

    fetchPermissions () {
      this.incLoader()

      this.api.permissionsList()
        .then(permissions => {
          this.resources = new Set()
          this.permissions = (permissions || [])
            .reduce((map, { type, any, op }) => {
              this.resources.add(any)
              if (!map[type]) {
                map[type] = { any, ops: [] }
              }

              map[type].ops.push(op)
              return map
            }, {})
        })
        .then(() => this.fetchRoles())
        .catch(this.toastErrorHandler(this.$t('notification:permissions.fetch.system')))
        .finally(() => {
          this.loaded.permissions = true
          this.decLoader()
        })
    },

    onSubmit (roleRules) {
      this.permission.processing = true
      Promise.all(roleRules.filter(({ ID }) => ID.includes('edit')).map(({ ID, rules }) => {
        const externalRules = []
        const roleID = ID.split('-')[1]
        Object.entries(rules).forEach(([ key, value ]) => {
          let [ operation, resource ] = key.split('@', 2)
          externalRules.push({ roleID, resource, operation, access: value })
        })

        return this.api.permissionsUpdate({ roleID, rules: externalRules })
      })).then(() => {
        this.animateSuccess('permission')
        this.toastSuccess(this.$t('notification:permissions.update.success'))
      }).catch(this.toastErrorHandler(this.$t('notification:permissions.update.error')))
        .finally(() => {
          Promise.all(this.roles.filter(({ mode }) => mode === 'eval').map(({ roleID, userID }) => {
            return this.api.permissionsTrace({ roleID, userID }).then(rr => {
              let ID = userID ? `eval-${userID}` : `eval-${roleID.join('-')}`

              this.rolePermissions = [
                ...this.rolePermissions.filter(rp => rp.ID !== ID),
                { resource: '', ID, rules: this.roleRules(rr, 'eval') },
              ]
            })
          })).finally(() => {
            this.permission.processing = false
          })
        })
    },

    addRole (add) {
      this.loaded.roles = false
      const { mode } = add || {}

      if (mode === 'edit') {
        const { roleID, name } = add.roleID || {}

        this.readPermissions({ roleID, name: [name] })
          .finally(() => {
            this.loaded.roles = true
          })
      } else if (mode === 'eval') {
        let { userID, roleID } = add
        let name = ''

        if (userID) {
          const { name: uname, username, email, handle } = userID
          name = [uname || username || email || handle || userID || '']
          userID = userID.userID
        } else {
          name = roleID.map(({ name }) => name)
          roleID = roleID.map(({ roleID }) => roleID)
        }

        this.evaluatePermissions({ name, roleID, userID })
          .finally(() => {
            this.loaded.roles = true
          })
      }
    },

    hideRole (role) {
      this.loaded.roles = false
      const { ID } = role

      this.roles = this.roles.filter(r => r.ID !== ID)
      setIncludedRoles(this.roles)

      this.loaded.roles = true
    },

    async readPermissions ({ name, roleID }) {
      const resource = [...this.resources]

      return this.api.permissionsRead({ resource, roleID })
        .then(rr => {
          this.rolePermissions.push({ resource: '', ID: `edit-${roleID}`, rules: this.roleRules(rr, 'edit') })
          this.roles.push({
            mode: 'edit',
            ID: `edit-${roleID}`,
            roleID,
            name,
          })
        })
        .catch(this.toastErrorHandler(this.$t('notification:permissions.role.error')))
        .finally(() => {
          setIncludedRoles(this.roles)
        })
    },

    async evaluatePermissions ({ name, roleID, userID }) {
      const resource = [...this.resources]

      return this.api.permissionsTrace({ resource, roleID, userID })
        .then(rr => {
          const ID = userID ? `eval-${userID}` : `eval-${roleID.join('-')}`

          this.rolePermissions.push({ resource: '', ID, rules: this.roleRules(rr, 'eval') })
          this.roles.push({
            mode: 'eval',
            ID,
            roleID,
            userID,
            name,
          })
        })
        .catch(this.toastErrorHandler(this.$t('notification:permissions.eval.error')))
        .finally(() => {
          setIncludedRoles(this.roles)
        })
    },

    roleRules (rules, mode = 'edit') {
      return (rules || [])
        .reduce((map, { resource, operation, access, resolution }) => {
          const [ type ] = resource.split('/', 2)
          if ((this.permissions[type] || { ops: [] }).ops.indexOf(operation) > -1) {
            if (mode === 'eval') {
              if (resolution === 'unknown-context') {
                access = 'unknown-context'
              } else if (access === 'inherit') {
                access = 'deny'
              }
            }

            map[`${operation}@${resource}`] = access
          }

          return map
        }, {})
    },
  },
}
