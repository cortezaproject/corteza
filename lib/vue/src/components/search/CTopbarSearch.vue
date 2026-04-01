<template>
  <b-modal
    v-model="showModal"
    hide-footer
    hide-header
    size="lg"
    body-class="p-0"
    content-class="overflow-hidden border-0 shadow-lg"
    dialog-class="search-modal-dialog"
    @shown="onShown"
  >
    <div
      class="search-container d-flex flex-column"
      style="max-height: 70vh;"
    >
      <!-- Header / Input -->
      <div
        class="border-bottom"
      >
        <c-input-search
          v-model="query"
          :placeholder="labels.placeholder"
          :loading="loading"
          submittable
          class="topbar-search-input"
          @search="submitSearch"
        />
      </div>

      <!-- Recent Searches -->
      <template v-if="query.length < 2 && recentSearches.length > 0">
        <div class="flex-grow-1 overflow-auto">
          <div class="px-3 py-2 bg-extra-light border-bottom d-flex align-items-center justify-content-between text-muted">
            <span class="small font-weight-bold text-uppercase">{{ labels.recentSearches }}</span>
            <b-button
              variant="outline-light"
              class="px-1 py-0 text-muted small border-0 shadow-none"
              @click="clearRecentSearches"
            >
              {{ labels.clearHistory }}
            </b-button>
          </div>
          <div
            v-for="(s, index) in recentSearches"
            :key="index"
            class="recent-search-item d-flex align-items-center justify-content-between px-3 py-2 cursor-pointer"
            @click="useRecentSearch(s)"
          >
            <div class="d-flex align-items-center">
              <font-awesome-icon
                :icon="['fas', 'history']"
                class="text-muted mr-3 small"
              />
              <span class="text-dark">{{ s }}</span>
            </div>
            <b-button
              variant="outline-extra-light"
              class="remove-btn px-2 text-muted small border-0 shadow-none"
              @click.stop="removeRecentSearch(index)"
            >
              <font-awesome-icon
                :icon="['fas', 'times']"
              />
            </b-button>
          </div>
        </div>
      </template>

      <template v-if="query.length >= 2 && (hasResults || hasSearched || !loading)">
        <!-- Results List -->
        <div
          v-if="!hasResults && !loading && hasSearched"
          class="flex-grow-1 overflow-auto p-5 text-center text-muted"
        >
          <p>{{ labels.noResults() }}</p>
        </div>

        <div
          v-else-if="hasResults"
          ref="resultsList"
          class="search-results-list d-flex flex-column flex-grow-1 overflow-auto"
        >
          <div
            v-for="ns in sortedGroups"
            :key="ns.id"
            class="border-bottom"
          >
            <item-group
              :title="ns.name"
              :items="ns.items"
              :collapse-id="`collapse-${ns.id}`"
              :expanded="ns.expanded"
              :labels="labels"
              @update:expanded="(val) => $set(expandedGroups, ns.id, val)"
            >
              <item-group
                v-for="mod in ns.sortedModules"
                :key="mod.id"
                :title="mod.name"
                :items="mod.items"
                :collapse-id="`collapse-${ns.id}-${mod.id}`"
                :expanded="mod.expanded"
                :labels="labels"
                subgroup
                @update:expanded="(val) => $set(expandedGroups, `${ns.id}-${mod.id}`, val)"
              >
                <record-item
                  v-for="hit in mod.items"
                  :key="hit.id"
                  :hit="hit"
                  :labels="labels"
                  @click="onResultClick(hit)"
                  @open-new-tab="onOpenNewTab(hit)"
                />
              </item-group>
            </item-group>
          </div>
        </div>
      </template>
    </div>
  </b-modal>
</template>

<script>
import axios from 'axios'
import { library } from '@fortawesome/fontawesome-svg-core'
import { faTimes, faHistory } from '@fortawesome/free-solid-svg-icons'
import RecordItem from './items/RecordItem.vue'
import ItemGroup from './items/ItemGroup.vue'
import { CInputSearch } from '../input'

library.add(faTimes, faHistory)

export default {
  name: 'CTopbarSearch',

  components: {
    RecordItem,
    ItemGroup,
    CInputSearch,
  },

  props: {
    labels: {
      type: Object,
      default: () => ({
        numberOfResults: (count) => `${count} results`,
      }),
    },
  },

  data () {
    return {
      showModal: false,
      query: '',
      loading: false,
      results: [],
      hasSearched: false,
      recentSearches: [],
      cancelRequest: null,
      expandedGroups: {},
    }
  },

  computed: {
    // Get the current namespace slug from route params
    currentNamespaceSlug () {
      return this.$route.params.slug || ''
    },

    // All results as a flat list (used when not on a namespace)
    allResults () {
      return this.results
        .map(hit => {
          const highlight = this.getHitHighlight(hit)
          return {
            ...hit,
            highlight,
          }
        })
        .filter(hit => hit.highlight.label)
    },

    // Groups results by namespace and then by module, sorts by relevance/first appearance
    sortedGroups () {
      const currentNs = this.currentNamespaceSlug
      const nsOrder = [] // Track order of first appearance for namespaces
      const namespaces = {}

      this.allResults.forEach(hit => {
        const ns = hit.value.namespace || {}
        const nsID = ns.namespaceID || 'unknown'
        const nsSlug = ns.slug || nsID
        const mod = hit.value.module || {}
        const modID = mod.moduleID || 'unknown'

        if (!namespaces[nsID]) {
          namespaces[nsID] = {
            id: nsID,
            name: ns.name || nsSlug,
            slug: nsSlug,
            modules: {},
            moduleOrder: [], // Track order of first appearance for modules within this NS
            expanded: this.expandedGroups[nsID] !== false,
          }
          nsOrder.push(nsID)
        }

        const nsObj = namespaces[nsID]
        if (!nsObj.modules[modID]) {
          nsObj.modules[modID] = {
            id: modID,
            name: mod.name || modID,
            items: [],
            expanded: this.expandedGroups[`${nsID}-${modID}`] !== false,
          }
          nsObj.moduleOrder.push(modID)
        }

        nsObj.modules[modID].items.push(hit)
      })

      // Convert to array and sort namespaces
      return nsOrder
        .map(nsID => {
          const ns = namespaces[nsID]
          // Convert modules to sorted array
          ns.sortedModules = ns.moduleOrder.map(modID => ns.modules[modID])
          // Add total item count for NS badge
          ns.items = ns.sortedModules.reduce((acc, mod) => acc.concat(mod.items), [])
          return ns
        })
        .sort((a, b) => {
          // Current namespace always first
          if (a.slug === currentNs || a.id === currentNs) return -1
          if (b.slug === currentNs || b.id === currentNs) return 1
          // Otherwise keep original order (relevance)
          return 0
        })
    },

    // Check if there are any results at all
    hasResults () {
      return this.allResults.length > 0
    },
  },

  watch: {
    showModal (val) {
      if (val) {
        this.query = ''
        this.results = []
        this.hasSearched = false
      }
    },
    query (newVal) {
      // Clear results when query is cleared or too short
      if (newVal.length < 2) {
        this.loading = false
        this.results = []
        this.hasSearched = false
      }
    },
  },


  mounted () {
    window.addEventListener('keydown', this.handleKeydown)
    this.$root.$on('search:open', this.openSearch)
    this.loadRecentSearches()
  },

  beforeDestroy () {
    window.removeEventListener('keydown', this.handleKeydown)
    this.$root.$off('search:open', this.openSearch)
  },

  methods: {
    openSearch () {
      this.showModal = true
    },

    onShown () {
      if (this.$refs.searchInput) {
        this.$refs.searchInput.focus()
      }
    },

    closeSearch () {
      this.showModal = false
    },

    handleKeydown (e) {
      // Global shortcut to open search
      if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault()
        this.openSearch()
        return
      }

      // Escape to close
      if (this.showModal && e.key === 'Escape') {
        this.closeSearch()
      }
    },

    loadRecentSearches () {
      const saved = localStorage.getItem('discovery-recent-searches')
      if (saved) {
        try {
          this.recentSearches = JSON.parse(saved)
        } catch (e) {
          this.recentSearches = []
        }
      }
    },

    addToRecent (q) {
      if (!q || q.length < 2) return
      const list = [q, ...this.recentSearches.filter(s => s !== q)].slice(0, 5)
      this.recentSearches = list
      localStorage.setItem('discovery-recent-searches', JSON.stringify(list))
    },

    clearRecentSearches () {
      this.recentSearches = []
      localStorage.removeItem('discovery-recent-searches')
    },

    removeRecentSearch (index) {
      this.recentSearches.splice(index, 1)
      localStorage.setItem('discovery-recent-searches', JSON.stringify(this.recentSearches))
    },

    useRecentSearch (s) {
      this.query = s
      this.submitSearch()
    },

    submitSearch () {
      if (this.query.length >= 2) {
        this.addToRecent(this.query)
        this.loading = true
        this.performSearch(this.query)
      }
    },

    async performSearch (query) {
      if (query.length < 2) {
        return
      }

      // Cancel any pending request
      if (this.cancelRequest) {
        this.cancelRequest()
      }

      this.loading = true

      const { response, cancel } = this.$DiscoveryAPI.queryCancellable({
        query,
        resourceTypes: ['compose:record'],
        size: 20,
      })

      this.cancelRequest = cancel

      try {
        const { hits = [] } = await response()

        this.results = hits || []
      } catch (e) {
        // Ignore cancel errors, they're expected when cancelling
        if (axios.isCancel(e)) {
          return
        }
        console.error('Search failed', e)
        this.results = []
      } finally {
        this.loading = false
        this.hasSearched = true
      }
    },

    getHitHighlight (hit) {
      const matchingFields = hit.value.matching_fields || {}
      const fieldName = Object.keys(matchingFields)[0]
      const values = hit.value.values || []

      // Use matching_fields to identify which field matched, then get its clean value
      const field = fieldName
        ? values.find(v => v.name === fieldName)
        : values[0]

      return {
        label: field ? (field.label || field.name) : '',
        value: field && field.value ? field.value[0] : '',
      }
    },

    async resolveRecordRoute (hit) {
      const { recordID, module, namespace } = hit.value
      const { namespaceID } = namespace
      const { moduleID } = module

      const ns = await this.$ComposeAPI.namespaceRead({ namespaceID })
      if (!ns) {
        this.toastDanger(this.labels.notFoundNamespace)
        return null
      }

      const slug = ns.slug || ns.namespaceID

      const pages = await this.$ComposeAPI.pageList({ namespaceID, moduleID })
      const recordPage = (pages.set || []).find(p => p.moduleID === moduleID)

      if (!recordPage) {
        this.toastDanger(this.labels.notFoundPage)
        return null
      }

      return {
        name: 'page.record',
        params: { slug, pageID: recordPage.pageID, recordID },
      }
    },

    async onResultClick (hit) {
      if (hit.type !== 'compose:record') return

      this.closeSearch()
      this.addToRecent(this.query)

      try {
        const route = await this.resolveRecordRoute(hit)
        if (!route) return

        const isPagesRoute = this.$route && (['pages', 'page', 'page.record', 'page.record.edit', 'page.record.create'].includes(this.$route.name))
        const sameNamespace = route.params.slug === this.$route.params.slug

        if (isPagesRoute && sameNamespace) {
          this.$root.$emit('show-record-modal', {
            recordID: route.params.recordID,
            recordPageID: route.params.pageID,
          })
        } else {
          this.$router.push(route)
        }
      } catch (e) {
        this.toastErrorHandler(this.labels.recordRedirectError)(e)
      }
    },

    async onOpenNewTab (hit) {
      if (hit.type !== 'compose:record') return

      this.addToRecent(this.query)

      try {
        const route = await this.resolveRecordRoute(hit)
        if (!route) return

        const url = this.$router.resolve(route).href
        window.open(url, '_blank')
      } catch (e) {
        this.toastErrorHandler(this.labels.recordRedirectError)(e)
      }
    },
  },
}
</script>

<style lang="scss" scoped>
  /* Help labels truncate in flex row */
.text-truncate {
  min-width: 0;
}

.cursor-pointer {
  cursor: pointer;
}

.recent-search-item {
  transition: background-color 0.2s ease;

  .remove-btn {
    opacity: 0;
    transition: opacity 0.2s ease, color 0.2s ease;
  }

  &:hover {
    background-color: var(--light);

    .remove-btn {
      opacity: 1;
    }
  }
}
</style>

<style lang="scss">
.search-modal-dialog {
  margin-top: calc(var(--topbar-height) + 1rem) !important;
  max-width: 700px;
}

.topbar-search-input {
  input {
    border: none !important;
    background: transparent !important;
  }

  .search-button {
    top: 0 !important;
    bottom: 0 !important;
    right: 0 !important;
    border-top: none !important;
    border-right: none !important;
    border-bottom: none !important;
  }
}
</style>
