<template>
  <Dialog
    v-model:visible="showModal"
    :header="undefined"
    :draggable="false"
    :modal="true"
    :dismissable-mask="true"
    position="top"
    :style="{ width: '700px', marginTop: 'calc(var(--topbar-height) + 1rem)' }"
    :pt="{
      header: { class: 'hidden' },
      content: { class: '!p-0 overflow-hidden' },
    }"
    @show="onShown"
  >
    <div class="flex flex-col" style="max-height: 70vh">
      <!-- Search input -->
      <div class="flex items-center">
        <i class="pi pi-search text-muted-color ml-3" />
        <InputText
          ref="searchInputRef"
          v-model="query"
          :placeholder="labels.placeholder || 'Search...'"
          class="flex-1 border-0 shadow-none bg-transparent focus:ring-0 px-3 py-3"
          @keyup.enter="submitSearch"
        />
        <ProgressSpinner
          v-if="loading"
          style="width: 1.25rem; height: 1.25rem"
          class="mr-3"
          stroke-width="4"
        />
      </div>

      <!-- Recent searches -->
      <template v-if="query.length < 2 && recentSearches.length > 0">
        <div class="flex-1 overflow-auto">
          <div
            class="flex items-center justify-between px-4 py-2 bg-surface border-b border-surface"
          >
            <span class="text-xs font-bold uppercase text-muted-color">
              {{ labels.recentSearches }}
            </span>
            <Button
              :label="labels.clearHistory"
              severity="secondary"
              variant="text"
              size="small"
              class="text-xs"
              @click="clearRecentSearches"
            />
          </div>

          <div
            v-for="(s, index) in recentSearches"
            :key="index"
            class="flex items-center justify-between px-4 py-2 cursor-pointer hover:bg-emphasis group"
            @click="useRecentSearch(s)"
          >
            <div class="flex items-center gap-3">
              <i class="pi pi-history text-muted-color text-xs" />
              <span class="text-sm text-color">{{ s }}</span>
            </div>
            <Button
              icon="pi pi-times"
              severity="secondary"
              variant="text"
              size="small"
              rounded
              class="opacity-0 group-hover:opacity-100"
              @click.stop="removeRecentSearch(index)"
            />
          </div>
        </div>
      </template>

      <!-- No results -->
      <template v-else-if="query.length >= 2 && !hasResults && !loading && hasSearched">
        <div class="flex-1 p-10 text-center text-muted-color">
          <p>{{ labels.noResults ? labels.noResults() : 'No results found' }}</p>
        </div>
      </template>

      <!-- Results -->
      <div v-else-if="hasResults" class="flex-1 overflow-auto">
        <div
          v-for="ns in sortedGroups"
          :key="ns.id"
          class="border-b border-surface"
        >
          <ItemGroup
            :title="ns.name"
            :items="ns.items"
            :collapse-id="`collapse-${ns.id}`"
            :expanded="ns.expanded"
            :labels="labels"
            @update:expanded="val => (expandedGroups[ns.id] = val)"
          >
            <ItemGroup
              v-for="mod in ns.sortedModules"
              :key="mod.id"
              :title="mod.name"
              :items="mod.items"
              :collapse-id="`collapse-${ns.id}-${mod.id}`"
              :expanded="mod.expanded"
              :labels="labels"
              subgroup
              @update:expanded="val => (expandedGroups[`${ns.id}-${mod.id}`] = val)"
            >
              <RecordItem
                v-for="hit in mod.items"
                :key="hit.id"
                :hit="hit"
                :labels="labels"
                @click="onResultClick(hit)"
                @open-new-tab="onOpenNewTab(hit)"
              />
            </ItemGroup>
          </ItemGroup>
        </div>
      </div>
    </div>
  </Dialog>
</template>

<script setup>
import { computed, inject, nextTick, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import RecordItem from './items/RecordItem.vue'
import ItemGroup from './items/ItemGroup.vue'

const props = defineProps({
  labels: {
    type: Object,
    default: () => ({
      placeholder: 'Search...',
      noResults: () => 'No results found',
      recentSearches: 'Recent searches',
      clearHistory: 'Clear history',
      openInNewTab: 'Open in new tab',
      numberOfResults: count => `${count} results`,
      notFoundNamespace: 'Namespace not found',
      notFoundPage: 'Page not found',
      recordRedirectError: 'Could not navigate to record',
    }),
  },
})

const $DiscoveryAPI = inject('$DiscoveryAPI')
const $ComposeAPI = inject('$ComposeAPI')
const $toast = inject('$toast', null)
const route = useRoute()
const router = useRouter()

const showModal = ref(false)
const query = ref('')
const loading = ref(false)
const results = ref([])
const hasSearched = ref(false)
const recentSearches = ref([])
const cancelRequest = ref(null)
const expandedGroups = reactive({})
const searchInputRef = ref(null)

const currentNamespaceSlug = computed(() => route.params.slug || '')

const allResults = computed(() => {
  return results.value
    .map(hit => {
      const highlight = getHitHighlight(hit)
      return { ...hit, highlight }
    })
    .filter(hit => hit.highlight.label)
})

const sortedGroups = computed(() => {
  const currentNs = currentNamespaceSlug.value
  const nsOrder = []
  const namespaces = {}

  allResults.value.forEach(hit => {
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
        moduleOrder: [],
        expanded: expandedGroups[nsID] !== false,
      }
      nsOrder.push(nsID)
    }

    const nsObj = namespaces[nsID]
    if (!nsObj.modules[modID]) {
      nsObj.modules[modID] = {
        id: modID,
        name: mod.name || modID,
        items: [],
        expanded: expandedGroups[`${nsID}-${modID}`] !== false,
      }
      nsObj.moduleOrder.push(modID)
    }

    nsObj.modules[modID].items.push(hit)
  })

  return nsOrder
    .map(nsID => {
      const ns = namespaces[nsID]
      ns.sortedModules = ns.moduleOrder.map(modID => ns.modules[modID])
      ns.items = ns.sortedModules.reduce((acc, mod) => acc.concat(mod.items), [])
      return ns
    })
    .sort((a, b) => {
      if (a.slug === currentNs || a.id === currentNs) return -1
      if (b.slug === currentNs || b.id === currentNs) return 1
      return 0
    })
})

const hasResults = computed(() => allResults.value.length > 0)

watch(showModal, val => {
  if (val) {
    query.value = ''
    results.value = []
    hasSearched.value = false
  }
})

watch(query, newVal => {
  if (newVal.length < 2) {
    loading.value = false
    results.value = []
    hasSearched.value = false
  }
})

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
  loadRecentSearches()
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleKeydown)
})

function openSearch() {
  showModal.value = true
}

function onShown() {
  nextTick(() => {
    searchInputRef.value?.$el?.focus?.() || searchInputRef.value?.focus?.()
  })
}

function closeSearch() {
  showModal.value = false
}

function handleKeydown(e) {
  if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
    e.preventDefault()
    openSearch()
    return
  }

  if (showModal.value && e.key === 'Escape') {
    closeSearch()
  }
}

function loadRecentSearches() {
  const saved = localStorage.getItem('discovery-recent-searches')
  if (saved) {
    try {
      recentSearches.value = JSON.parse(saved)
    } catch {
      recentSearches.value = []
    }
  }
}

function addToRecent(q) {
  if (!q || q.length < 2) return
  const list = [q, ...recentSearches.value.filter(s => s !== q)].slice(0, 5)
  recentSearches.value = list
  localStorage.setItem('discovery-recent-searches', JSON.stringify(list))
}

function clearRecentSearches() {
  recentSearches.value = []
  localStorage.removeItem('discovery-recent-searches')
}

function removeRecentSearch(index) {
  recentSearches.value.splice(index, 1)
  localStorage.setItem('discovery-recent-searches', JSON.stringify(recentSearches.value))
}

function useRecentSearch(s) {
  query.value = s
  submitSearch()
}

function submitSearch() {
  if (query.value.length >= 2) {
    addToRecent(query.value)
    loading.value = true
    performSearch(query.value)
  }
}

async function performSearch(q) {
  if (q.length < 2) return

  if (cancelRequest.value) {
    cancelRequest.value()
  }

  loading.value = true

  const { response, cancel } = $DiscoveryAPI.queryCancellable({
    query: q,
    resourceTypes: ['compose:record'],
    size: 20,
  })

  cancelRequest.value = cancel

  try {
    const { hits = [] } = await response()
    results.value = hits || []
  } catch (e) {
    if (axios.isCancel(e)) return
    results.value = []
  } finally {
    loading.value = false
    hasSearched.value = true
  }
}

function getHitHighlight(hit) {
  const matchingFields = hit.value.matching_fields || {}
  const fieldName = Object.keys(matchingFields)[0]
  const values = hit.value.values || []

  const field = fieldName ? values.find(v => v.name === fieldName) : values[0]

  return {
    label: field ? field.label || field.name : '',
    value: field && field.value ? field.value[0] : '',
  }
}

async function resolveRecordNav(hit) {
  const { recordID, module, namespace } = hit.value
  const { namespaceID } = namespace
  const { moduleID } = module

  const ns = await $ComposeAPI.namespaceRead({ namespaceID })
  if (!ns) {
    $toast?.add({ severity: 'error', summary: props.labels.notFoundNamespace, life: 3000 })
    return null
  }

  const slug = ns.slug || ns.namespaceID
  const { set: recordPages = [] } = await $ComposeAPI.pageList({ namespaceID, moduleID })
  if (!recordPages.length) {
    $toast?.add({ severity: 'error', summary: props.labels.notFoundPage, life: 3000 })
    return null
  }

  const pageID = recordPages[0].pageID
  const u = new URL(window.location.href)
  const externalUrl = `${u.origin}/compose/namespace/${slug}/pages/${pageID}/records/${recordID}`
  const routeLocation = { name: 'page.record', params: { slug, pageID, recordID } }

  return { routeLocation, externalUrl }
}

async function onResultClick(hit) {
  if (hit.type !== 'compose:record') return

  closeSearch()
  addToRecent(query.value)

  try {
    const nav = await resolveRecordNav(hit)
    if (!nav) return

    const { routeLocation, externalUrl } = nav

    if (router.hasRoute('page.record')) {
      router.push(routeLocation)
    } else {
      window.location.href = externalUrl
    }
  } catch {
    $toast?.add({ severity: 'error', summary: props.labels.recordRedirectError, life: 3000 })
  }
}

async function onOpenNewTab(hit) {
  if (hit.type !== 'compose:record') return

  addToRecent(query.value)

  try {
    const nav = await resolveRecordNav(hit)
    if (!nav) return

    const { routeLocation, externalUrl } = nav

    if (router.hasRoute('page.record')) {
      window.open(router.resolve(routeLocation).href, '_blank', 'noopener')
    } else {
      window.open(externalUrl, '_blank', 'noopener')
    }
  } catch {
    $toast?.add({ severity: 'error', summary: props.labels.recordRedirectError, life: 3000 })
  }
}

defineExpose({ open: openSearch })
</script>
