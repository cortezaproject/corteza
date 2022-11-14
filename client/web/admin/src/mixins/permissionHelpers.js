import editorHelpers from './editorHelpers'

const systemRoles = ['1', '2']
const lsKey = 'permissionList.roles'

/**
 * Determines what roles should be included in the table
 * @returns {Array<String>}
 */
function getIncludedRoles () {
  const udRoles = JSON.parse(localStorage.getItem(lsKey) || '[]')
  return systemRoles.concat(udRoles)
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
      // Array of roleID's included in the permission list
      rolesIncluded: getIncludedRoles(),
      // Array of roleID's not included in the permission list
      rolesExcluded: [],

      roles: [],
      rolePermissions: [],
      permissions: {},
      effective: {},

      api: this.$SystemAPI,
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
  },

  methods: {
    fetchRoles () {
      this.incLoader()

      return this.$SystemAPI.roleList()
        .then(({ set }) => set.filter(({ isBypass }) => !isBypass))
        .then(set => {
          this.roles = set.filter(r => this.rolesIncluded.includes(r.roleID))
          this.rolePermissions = []
          this.rolesExcluded = []

          // We read permissions for included roles
          return Promise.all(set.map(role => {
            const { roleID } = role

            if (this.rolesIncluded.includes(roleID)) {
              return this.api.permissionsRead({ roleID })
                .then(rr => {
                  this.rolePermissions.push({ roleID, rules: this.roleRules(rr) })
                })
            } else {
              // Keep track of excluded roles that can be added to the list
              this.rolesExcluded.push(role)
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
          this.permissions = (permissions || [])
            .reduce((map, { type, any, op }) => {
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
      Promise.all(roleRules.map(({ roleID, rules }) => {
        const externalRules = []
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
          this.permission.processing = false
        })
    },

    addRole (role) {
      this.loaded.roles = false
      const { roleID } = role
      this.rolesIncluded.push(roleID)
      this.rolesExcluded = this.rolesExcluded.filter(r => r.roleID !== roleID)

      // Store for next time
      setIncludedRoles(this.rolesIncluded)

      this.api.permissionsRead({ roleID })
        .then(rr => {
          this.rolePermissions.push({ roleID, rules: this.roleRules(rr) })

          // Add new role
          this.roles.push(role)
          this.loaded.roles = true
        })
        .catch(this.toastErrorHandler(this.$t('notification:permissions.role.error')))
        .finally(() => {
          this.loaded.roles = true
        })
    },

    hideRole (role) {
      this.loaded.roles = false
      const { roleID } = role
      this.rolesExcluded.push(role)
      this.rolesIncluded = this.rolesIncluded.filter(r => r !== roleID)

      // Store for next time
      setIncludedRoles(this.rolesIncluded)

      this.loaded.roles = true
      this.roles = this.roles.filter(r => r.roleID !== roleID)
    },

    roleRules (rules) {
      return (rules || [])
        .reduce((map, { resource, operation, access }) => {
          const [ type ] = resource.split('/', 2)
          if ((this.permissions[type] || { ops: [] }).ops.indexOf(operation) > -1) {
            map[`${operation}@${resource}`] = access
          }

          return map
        }, {})
    },
  },
}
