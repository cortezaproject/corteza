<template>
  <!-- Loading state -->
  <div v-if="loading" class="flex flex-col items-center justify-center h-full gap-4">
    <h1 class="text-2xl text-muted-color">
      {{ namespace?.name || namespace?.slug || '...' }}
    </h1>
    <div class="flex items-center gap-3">
      <ProgressSpinner style="width: 24px; height: 24px" />
      <span class="text-lg text-muted-color">{{ $t('general.label.loading') }}</span>
    </div>
  </div>

  <div v-else-if="namespace" class="flex h-full">
    <!-- Main content area -->
    <div class="flex-1 overflow-auto">
      <RouterView :namespace="namespace" />
    </div>

    <!-- Global Modals -->
    <RecordModal :namespace="namespace" />
  </div>

  <!-- Error state -->
  <div v-else class="flex items-center justify-center h-full">
    <Message severity="error" :closable="false">
      {{ $t('notification.namespace.loadFailed') }}
    </Message>
  </div>
</template>

<script setup>
import { useChartStore } from '@planetcrust/human-vue'
import { useModuleStore } from '@planetcrust/human-vue'
import { useNamespaceStore } from '@planetcrust/human-vue'
import { usePageStore } from '@planetcrust/human-vue'
import { usePageLayoutStore } from '@planetcrust/human-vue'
import { compose, NoID } from '@planetcrust/human-js'
import { useMinDuration } from '@planetcrust/human-vue'
import { inject, onMounted, provide, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import RecordModal from '@/sections/compose/components/Record/RecordModal.vue'

const props = defineProps({
  slug: {
    type: String,
    required: true,
  },
})

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const $toast = inject('$toast')

// Stores
const namespaceStore = useNamespaceStore()
const moduleStore = useModuleStore()
const pageStore = usePageStore()
const chartStore = useChartStore()
const pageLayoutStore = usePageLayoutStore()

// State with minimum loading duration to prevent spinner flash
const { loading, run } = useMinDuration(1000)
const namespace = ref(null)

// Provide namespace so deeply nested components (e.g. CFieldRecordEditor) can access it
provide('$namespace', namespace)

// Provide page store for record field viewer navigation
provide('$pageStore', pageStore)

// Methods
async function loadNamespace() {
  await run(async () => {
    try {
      // First try to find in store
      let ns = namespaceStore.getByUrlPart(props.slug)

      if (!ns) {
        // Load all namespaces if not found
        await namespaceStore.load()
        ns = namespaceStore.getByUrlPart(props.slug)
      }

      if (!ns) {
        $toast.toastDanger(t('notification.namespace.loadFailed'))
        router.push({ name: 'root' })
        return
      }

      if (!ns.enabled) {
        const isAdminRoute = route.name?.toString().startsWith('admin.')
        if (ns.canUpdateNamespace && isAdminRoute) {
          // allow through — admin configuring a disabled namespace
        } else if (ns.canUpdateNamespace) {
          $toast.toastWarning(t('notification.namespace.disabled'))
          router.push({ name: 'namespace.edit', params: { slug: props.slug } })
          return
        } else {
          $toast.toastWarning(t('notification.namespace.disabled'))
          router.push({ name: 'root' })
          return
        }
      }

      namespace.value = new compose.Namespace({ ...ns })

      // Prepare namespace context - clear and preload all stores
      await prepareNamespace()

      // Redirect to first page on initial namespace load
      // Now guaranteed to have pages loaded because we awaited prepareNamespace!
      redirectToFirstPage()
    } catch (error) {
      console.error('Failed to load namespace:', error)
      $toast.toastDanger(t('notification.namespace.loadFailed'))
    }
  })
}

/**
 * Prepare namespace context by clearing and preloading all namespace-scoped stores.
 */
async function prepareNamespace() {
  if (!namespace.value) return

  const nsID = namespace.value.namespaceID

  // Clear all namespace-scoped stores for fresh load
  moduleStore.clearSet()
  pageStore.clearSet()
  chartStore.clearSet()
  pageLayoutStore.clearSet()

  try {
    // Preload all namespace data in parallel
    await Promise.all([
      moduleStore.load({ namespaceID: nsID }),
      pageStore.load({ namespaceID: nsID }),
      chartStore.load({ namespaceID: nsID }),
      pageLayoutStore.load({ namespaceID: nsID }),
    ])
  } catch (error) {
    console.error('Failed to prepare namespace context:', error)
    // Continue anyway - individual stores will handle their errors
  }
}

/**
 * Redirect to the first visible root-level page when landing on the default namespace route.
 */
function redirectToFirstPage() {
  // Only redirect if on the default 'pages' route or the parent 'namespace.view'
  if (route.name !== 'pages' && route.name !== 'namespace.view') return

  const pages = pageStore.set

  // Find first visible, root-level, non-record page (sorted by weight)
  const firstPage = [...pages]
    .filter(p => p.visible && p.selfID === NoID && p.moduleID === NoID)
    .sort((a, b) => a.weight - b.weight)[0]

  if (firstPage) {
    router.replace({
      name: 'page',
      params: { slug: props.slug, pageID: firstPage.pageID },
      query: route.query,
    })
  }
}

// Lifecycle
onMounted(() => {
  loadNamespace()
})

// Watch for slug changes (namespace switch)
watch(
  () => props.slug,
  () => {
    loadNamespace()
  },
)
</script>
