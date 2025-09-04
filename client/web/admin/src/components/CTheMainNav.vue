<template>
  <div
    class="overflow-auto py-3"
  >
    <b-list-group
      v-for="(grp, g) in navigation"
      :key="g"
      tag="li"
    >
      <h6
        v-if="grp.header"
        class="mt-3 text-uppercase"
      >
        {{ $t(grp.header.label) }}
      </h6>

      <c-sidebar-nav-items
        :items="grp.items"
        default-route-name="dashboard"
        class="overflow-auto h-100"
      />
    </b-list-group>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { components } from '@cortezaproject/corteza-vue'
const { CSidebarNavItems } = components

export default {
  components: {
    CSidebarNavItems,
  },

  i18nOptions: {
    namespaces: 'navigation',
  },

  data () {
    return {
      nav: [
        {
          items: [
            {
              label: 'dashboard',
              route: 'dashboard',
              icon: 'tachometer-alt',
            },
          ],
        },
        {
          header: { label: 'system.group' },
          items: [
            {
              label: 'system.items.users',
              route: 'system.user',
              icon: 'users',
              can: [
                ['system/', 'users.search'],
                ['system/', 'user.create'],
              ],
            },
            {
              label: 'system.items.roles',
              route: 'system.role',
              icon: 'hat-cowboy',
              can: [
                ['system/', 'roles.search'],
                ['system/', 'role.create'],
              ],
            },
            {
              label: 'system.items.usergroups',
              route: 'system.userGroup',
              icon: 'chart-diagram',
              can: ['system/', 'user-groups.search'],
            },
            {
              label: 'system.items.applications',
              route: 'system.application',
              icon: 'th-large',
              can: [
                ['system/', 'applications.search'],
                ['system/', 'application.create'],
              ],
            },
            {
              label: 'system.items.templates',
              route: 'system.template',
              icon: 'file-code',
              can: ['system/', 'templates.search'],
            },
            {
              label: 'system.items.settings',
              route: 'system.settings',
              icon: 'sliders-h',
              can: ['system/', 'settings.read'],
            },
            {
              label: 'system.items.email',
              route: 'system.email',
              icon: 'envelope-open',
              // all email management is handled
              // via settings for now
              can: [
                ['system/', 'settings.read'],
                ['system/', 'settings.manage'],
              ],
            },
            {
              label: 'system.items.authclients',
              route: 'system.authClient',
              icon: 'key',
              can: [
                ['system/', 'auth-clients.search'],
                ['system/', 'auth-client.create'],
              ],
            },
            {
              label: 'system.items.actionlog',
              route: 'system.actionlog',
              icon: 'glasses',
              can: ['system/', 'action-log.read'],
            },
            {
              label: 'system.items.queues',
              route: 'system.queue',
              icon: 'stream',
              can: [
                ['system/', 'queues.search'],
                ['system/', 'queue.create'],
              ],
            },
            {
              label: 'system.items.apigw',
              route: 'system.apigw',
              icon: 'archway',
              can: [
                ['system/', 'apigw-routes.search'],
                ['system/', 'apigw-route.create'],
              ],
            },
            {
              label: 'system.items.connections',
              route: 'system.connection',
              icon: 'cloud',
              can: [
                ['system/', 'dal-connections.search'],
                ['system/', 'dal-connection.create'],
              ],
            },
            {
              label: 'system.items.code-snippets',
              route: 'system.codesnippets',
              icon: 'file-code',
              can: [
                ['system/', 'settings.read'],
                ['system/', 'settings.manage'],
              ],
            },
            {
              label: 'system.items.sensitivityLevel',
              route: 'system.sensitivityLevel',
              icon: 'stamp',
              can: ['system/', 'dal-sensitivity-level.manage'],
            },
            {
              label: 'system.items.permissions',
              route: 'system.permissions',
              icon: 'lock',
              can: ['system/', 'grant'],
            },
          ],
        },

        {
          header: { label: 'compose.group' },
          items: [
            {
              label: 'compose.items.settings',
              route: 'compose.settings',
              icon: 'sliders-h',
              can: [
                ['compose/', 'settings.read'],
                ['compose/', 'settings.manage'],
              ],
            },

            {
              label: 'compose.items.permissions',
              route: 'compose.permissions',
              icon: 'lock',
              can: ['compose/', 'grant'],
            },
          ],
        },

        {
          header: { label: 'automation.group' },
          items: [
            {
              label: 'automation.items.workflows',
              route: 'automation.workflow',
              icon: 'project-diagram',
              can: [
                ['automation/', 'workflows.search'],
                ['automation/', 'workflow.create'],
              ],
            },

            {
              label: 'automation.items.sessions',
              route: 'automation.session',
              icon: 'business-time',
              can: ['automation/', 'sessions.search'],
            },

            {
              label: 'automation.items.scripts',
              route: 'automation.scripts',
              icon: 'scroll',
              can: ['automation/', 'workflows.search'],
            },

            {
              label: 'automation.items.permissions',
              route: 'automation.permissions',
              icon: 'lock',
              can: ['automation/', 'grant'],
            },
          ],
        },
        {
          header: { label: 'federation.group' },
          items: [
            {
              label: 'federation.items.nodes',
              route: 'federation.nodes',
              icon: 'share-alt',
              can: ['federation/', 'pair'],
            },
            {
              label: 'federation.items.permissions',
              route: 'federation.permissions',
              icon: 'lock',
              can: ['federation/', 'grant'],
            },
          ],
        },
        {
          header: { label: 'ui.group' },
          items: [
            {
              label: 'ui.items.settings',
              route: 'ui.settings',
              icon: 'eye',
              can: [
                ['system/', 'settings.read'],
                ['system/', 'settings.manage'],
              ],
            },
          ],
        },
      ],
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    navigation () {
      return this.nav.map(grp => {
        grp = JSON.parse(JSON.stringify(grp))
        grp.items = grp.items
          .map(itm => {
            const page = {
              name: itm.route,
              title: this.$t(itm.label),
              icon: ['fas', itm.icon],
            }

            if (!itm.can) {
              // if not explicitly set, allow
              itm.can = true
            }

            if (Array.isArray(itm.can)) {
              // if array then call the perm checking function
              if (Array.isArray(itm.can[0])) {
                itm.can = itm.can.every(c => this.can(...c))
              } else {
                itm.can = this.can(...itm.can)
              }
            }

            return { page, can: itm.can }
          })
          .filter(({ can }) => can)

        if (grp.items.length === 0) {
          return null
        }

        return grp
      }).filter(i => i)
    },
  },
}
</script>
