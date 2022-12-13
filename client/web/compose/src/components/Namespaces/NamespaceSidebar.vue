<template>
  <div>
    <portal to="sidebar-header-expanded">
      <vue-select
        v-if="!hideNamespaceList"
        key="namespaceID"
        data-test-id="select-namespace"
        label="name"
        class="namespace-selector sticky-top bg-white mt-2"
        :clearable="false"
        :options="namespaces"
        :value="namespace"
        :selectable="option => option.namespaceID !== namespace.namespaceID"
        :placeholder="$t('pickNamespace')"
        @option:selected="namespaceSelected"
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
      </vue-select>
    </portal>

    <portal
      to="sidebar-body-expanded"
    >
      <div
        v-if="namespace"
        class="d-flex flex-column flex-grow-1"
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
            {{ $t('noPages') }}
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
import { VueSelect } from 'vue-select'
const { CSidebarNavItems, CInputSearch } = components

const publicPageWrap = ({ pageID, selfID, title, visible }) => ({
  page: {
    // name omitted as default is provided
    pageID,
    selfID,
    title,
    visible,
  },
  children: [],
  params: {
    pageID,
  },
})

const adminPageWrap = (page) => {
  return {
    page: {
      name: 'admin.pages.builder',
      pageID: `page-${page.pageID}`,
      selfID: page.selfID !== NoID ? `page-${page.selfID}` : 'pages',
      rootSelfID: 'pages',
      title: page.title || page.handle,
      visible: true,
    },
    children: [],
    params: {
      pageID: page.pageID,
    },
  }
}
const moduleWrap = (module) => {
  return {
    page: {
      name: 'admin.modules.edit',
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
    VueSelect,
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
        const pages = [...(this.isAdminPage ? this.adminRoutes() : this.publicRoutes.map(publicPageWrap))]

        if (!this.query) {
          return pages
        }

        return pages.filter(({ page }) => !['pages', 'modules', 'charts'].includes(page.pageID) && filter.Assert(page, this.query, 'title'))
      }

      return []
    },

    navItems () {
      const current = this.filteredPages
      const ax = this.pageIndex(this.isAdminPage ? this.adminRoutes() : this.pages.map(publicPageWrap))

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
        ...this.modules.map(moduleWrap),
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
        ...this.pages.map(adminPageWrap),
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
  },
}
</script>

<style lang="scss">
.namespace-selector {
  font-size: 1rem;
  min-width: 100%;

  .vs__dropdown-menu {
    min-width: 100%;
  }

  .vs__dropdown-option {
    text-overflow: ellipsis;
    overflow-x: hidden;
  }

  .vs__selected-options {
    flex-wrap: nowrap;
  }

  .vs__selected {
    max-width: 230px;
    display: inline-block;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .vs__open-indicator {
    fill: $primary;
  }
}
</style>
