<template>
  <div>
    <portal to="sidebar-header-expanded">
      <b-input-group
        v-if="!hideNamespaceList"
        class="d-flex w-100 mt-2"
      >
        <c-input-select
          data-test-id="select-namespace"
          :value="currentNamespaceID"
          :options="filteredNamespaces"
          :get-option-label="getNamespaceLabel"
          :selectable="option => option.namespaceID !== namespace.namespaceID"
          :placeholder="$t('pickNamespace')"
          :clearable="false"
          :autoscroll="false"
          :append-to-body="false"
          @input="namespaceSelected"
        >
          <template #list-header>
            <li
              v-if="showNamespaceListLink"
              class="border-bottom text-center mb-1"
            >
              <router-link
                :to="{ name: 'namespace.manage' }"
                data-test-id="button-manage-namespaces"
                class="d-block my-1 font-weight-bold text-decoration-none"
              >
                {{ $t('manageNamespaces') }}
              </router-link>
            </li>
          </template>
        </c-input-select>

        <b-input-group-append v-if="canManageNamespaces">
          <b-button
            v-b-tooltip.noninteractive.hover="{ title: $t('editNamespace'), container: '#body' }"
            data-test-id="button-namespace-edit"
            :disabled="!canUpdateNamespace"
            variant="primary"
            class="d-flex align-items-center"
            :to="{ name: 'namespace.edit', params: { namespaceID: namespaceID } }"
          >
            <font-awesome-icon :icon="['far', 'edit']" />
          </b-button>
        </b-input-group-append>
      </b-input-group>
    </portal>

    <portal
      to="sidebar-body-expanded"
    >
      <div
        v-if="namespace"
        class="d-flex flex-column flex-grow-1"
      >
        <div
          class="sticky-top w-100 py-2"
          style="background-color: var(--sidebar-bg);"
        >
          <b-button
            v-if="isAdminPage"
            data-test-id="button-public"
            variant="light"
            class="w-100 mb-2"
            :to="{ name: 'pages', params: { slug: namespace.slug || namespace.namespaceID } }"
          >
            {{ $t('publicPages') }}
          </b-button>

          <b-button
            v-else-if="namespace.canManageNamespace"
            data-test-id="button-admin"
            variant="light"
            class="w-100 mb-2"
            :to="{ name: 'admin.modules', params: { slug: namespace.slug || namespace.namespaceID } }"
          >
            {{ $t('adminPanel') }}
          </b-button>

          <c-input-search
            v-model.trim="query"
            :disabled="loading"
            :placeholder="$t(`searchPlaceholder.${isAdminPage ? 'admin' : 'public'}`)"
            :autocomplete="'off'"
          />
        </div>

        <div
          v-if="!loading"
        >
          <c-sidebar-nav-items
            :items="navItems"
            :start-expanded="!!query"
            default-route-name="page"
            class="overflow-auto h-100"
          />

          <div
            v-if="!navItems.length"
            class="d-flex justify-content-center mt-5"
          >
            {{ $t('noResults') }}
          </div>
        </div>

        <div
          v-else
          class="d-flex align-items-center justify-content-center mt-5"
        >
          <b-spinner />
        </div>
      </div>
    </portal>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { NoID } from '@cortezaproject/corteza-js'
import { components, filter } from '@cortezaproject/corteza-vue'
import { Portal } from 'portal-vue'
const { CSidebarNavItems, CInputSearch } = components

const moduleWrap = (module, pageName) => {
  return {
    page: {
      name: pageName,
      pageID: `module-${module.moduleID}`,
      selfID: 'modules',
      rootSelfID: 'modules',
      title: module.name || module.handle,
      visible: true,
    },
    children: [],
    params: {
      moduleID: module.moduleID,
    },
  }
}

const chartWrap = (chart) => {
  return {
    page: {
      name: 'admin.charts.edit',
      pageID: `chart-${chart.chartID}`,
      selfID: 'charts',
      rootSelfID: 'charts',
      title: chart.name || chart.handle,
      visible: true,
    },
    children: [],
    params: {
      chartID: chart.chartID,
    },
  }
}

export default {
  i18nOptions: {
    namespaces: 'sidebar',
  },

  components: {
    Portal,
    CSidebarNavItems,
    CInputSearch,
  },

  props: {
    namespaces: {
      type: Array,
      required: true,
      default: () => [],
    },
  },

  data () {
    return {
      namespace: undefined,

      query: '',
    }
  },

  computed: {
    ...mapGetters({
      moduleLoading: 'module/loading',
      chartLoading: 'chart/loading',
      pageLoading: 'page/loading',
      modules: 'module/set',
      pages: 'page/set',
      charts: 'chart/set',
      can: 'rbac/can',
    }),

    currentNamespaceID () {
      return this.namespace ? this.namespace.namespaceID : NoID
    },

    // Loading is true only when a resource is being force loaded (API call)
    loading () {
      return this.moduleLoading || this.chartLoading || this.pageLoading
    },

    hideNamespaceList () {
      const { hideNamespaceList } = this.$Settings.get('compose.ui.sidebar', {})
      return hideNamespaceList
    },

    canManageNamespaces () {
      if (this.can('compose/', 'namespace.create') || this.can('compose/', 'grant')) {
        return true
      }

      return this.namespaces.reduce((acc, ns) => {
        return acc || ns.canUpdateNamespace || ns.canDeleteNamespace
      }, false)
    },

    showNamespaceListLink () {
      const { hideNamespaceListLink } = this.$Settings.get('compose.ui.sidebar', {})
      return !hideNamespaceListLink && this.canManageNamespaces
    },

    isAdminPage () {
      return this.$route.name.includes('admin.')
    },

    publicRoutes () {
      return this.pages.filter(({ moduleID, visible }) => visible && moduleID === NoID)
    },

    filteredPages () {
      if (this.namespace) {
        // If on admin page, show admin pages. Otherwise show public pages
        const pages = [...(this.isAdminPage ? this.adminRoutes() : this.publicPageWrap(this.publicRoutes))]

        if (!this.query) {
          return pages
        }

        return pages.filter(({ page }) => !['pages', 'modules', 'charts'].includes(page.pageID) && filter.Assert(page, this.query, 'title'))
      }

      return []
    },

    filteredNamespaces () {
      return this.namespaces.filter(({ enabled }) => enabled)
    },

    navItems () {
      const current = this.filteredPages
      const ax = this.pageIndex(this.isAdminPage ? this.adminRoutes() : this.publicPageWrap(this.pages))

      // Correct potentially missing parent references
      for (const cp of current) {
        if (cp.page.selfID && cp.page.selfID !== NoID) {
          if (!ax[cp.page.selfID]) {
            cp.page.selfID = cp.page.rootSelfID
          }
        }
      }

      const cx = this.pageIndex(current)

      for (let i = current.length - 1; i >= 0; i--) {
        const cp = current[i]

        // Here, we'll need to nest our pages.
        // If the requested page isn't in the current index, check in the all index.
        // If still not there, just place it somewhere...
        // Remove hidden pages if not in admin pages section
        if (!this.isAdminPage && !cp.page.visible) {
          current.splice(i, 1)
        } else if (cp.page.selfID && cp.page.selfID !== NoID) {
          let p = cx[cp.page.selfID]
          if (!p) {
            if (ax[cp.page.selfID]) {
              current.splice(i, 1, ax[cp.page.selfID])
              p = ax[cp.page.selfID]
              cx[p.page.pageID] = p
              i++
            } else {
              current.splice(i, 0, cp)
              p = cp
              cx[p.page.pageID] = p
            }
          } else {
            current.splice(i, 1)
          }
          if (cp.page.visible) {
            p.children.unshift(cp)
          }
        }
      }

      return current
    },

    canUpdateNamespace () {
      return this.namespace ? this.namespace.canUpdateNamespace : false
    },

    namespaceID () {
      return this.namespace ? this.namespace.namespaceID : NoID
    },
  },

  watch: {
    isAdminPage: {
      handler () {
        this.query = ''
      },
    },

    '$route.params.slug': {
      immediate: true,
      handler (slug = '') {
        this.query = ''
        this.namespace = this.$store.getters['namespace/getByUrlPart'](slug)
      },
    },
  },

  methods: {
    namespaceSelected ({ namespaceID, canManageNamespace, slug = '' }) {
      let { name, params } = this.$route

      // Try to match page, otherwise redirect to pages root
      if (name.includes('admin.modules')) {
        name = 'admin.modules'
      } else if (name.includes('admin.pages')) {
        name = 'admin.pages'
      } else if (name.includes('admin.charts')) {
        name = 'admin.charts'
      }

      name = !params.pageID && canManageNamespace && !name.includes('namespace.') ? name : 'pages'

      this.$router.push({ name, params: { slug: slug || namespaceID } })
    },

    pageIndex (wraps) {
      const ix = {}

      for (const w of wraps) {
        ix[w.page.pageID] = w
      }

      return ix
    },

    adminRoutes () {
      const routeName = this.$route.name
      const pageName = routeName.startsWith('admin.modules.record') ? 'admin.modules.record.list' : 'admin.modules.edit'

      return [
        {
          page: {
            pageID: 'modules',
            selfID: NoID,
            name: 'admin.modules',
            title: this.$t('module'),
            visible: true,
          },
          children: [],
        },
        ...this.modules.map((m) => moduleWrap(m, pageName)),
        {
          page: {
            pageID: 'pages',
            selfID: NoID,
            name: 'admin.pages',
            title: this.$t('page'),
            visible: true,
          },
          children: [],
        },
        ...this.adminPageWrap(this.pages),
        {
          page: {
            pageID: 'charts',
            selfID: NoID,
            name: 'admin.charts',
            title: this.$t('chart'),
            visible: true,
          },
          children: [],
        },
        ...this.charts.map(chartWrap),
      ]
    },

    publicPageWrap (pages) {
      return pages.map(({ pageID, selfID, title, visible, config }) => {
        const { navItem = {} } = config
        const { icon = {}, expanded = '' } = navItem
        const { type = '', src = '' } = icon

        const iconType = 'attachment'
        let iconSrc = src

        if (type === iconType) {
          iconSrc = `${this.$ComposeAPI.baseURL}${src}`
        }

        return {
          page: {
            // name omitted as default is provided
            pageID,
            selfID,
            title,
            visible,
            expanded,
            icon: iconSrc,
          },
          children: [],
          params: {
            pageID,
          },
        }
      })
    },

    adminPageWrap (pages) {
      return pages.map(({ pageID, selfID, title, handle, config }) => {
        const { navItem = {} } = config
        const { icon = {} } = navItem
        const { type = '', src = '' } = icon

        const iconType = 'attachment'
        let iconSrc = src

        if (type === iconType) {
          iconSrc = `${this.$ComposeAPI.baseURL}${src}`
        }

        const pageName = this.$route.name === 'admin.pages.edit' ? 'admin.pages.edit' : 'admin.pages.builder'

        return {
          page: {
            name: pageName,
            pageID: `page-${pageID}`,
            selfID: selfID !== NoID ? `page-${selfID}` : 'pages',
            rootSelfID: 'pages',
            title: title || handle,
            visible: true,
            icon: iconSrc,
          },
          children: [],
          params: {
            pageID,
          },
        }
      })
    },

    getNamespaceLabel (value) {
      if (typeof value === 'string') {
        value = this.filteredNamespaces.find(({ namespaceID }) => namespaceID === value) || {}
      }

      return value.name
    },
  },
}
</script>
