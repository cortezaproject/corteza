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

  <!-- Loaded content -->
  <div v-else-if="namespace" class="flex h-full">
    <!-- Main content area -->
    <div class="flex-1 overflow-auto">
      <RouterView :namespace="namespace" />
    </div>
  </div>

  <!-- Error state -->
  <div v-else class="flex items-center justify-center h-full">
    <Message severity="error" :closable="false">
      {{ $t('notification.namespace.loadFailed') }}
    </Message>
  </div>
</template>

<script setup>
import { useChartStore } from '@/stores/chart'
import { useModuleStore } from '@/stores/module'
import { useNamespaceStore } from '@/stores/namespace'
import { usePageStore } from '@/stores/page'
import { usePageLayoutStore } from '@/stores/page-layout'
import { compose, NoID } from '@cortezaproject/corteza-js-next'
import { useMinDuration } from '@cortezaproject/corteza-vue-next'
import { inject, onMounted, provide, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

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

// Methods
async function loadNamespace() {
  await run(async () => {
    try {
      // First try to find in store
      let ns = namespaceStore.getByUrlPart(props.slug)

      if (!ns) {
        // Load all namespaces if not found
        await namespaceStore.load({ force: true })
        ns = namespaceStore.getByUrlPart(props.slug)
      }

      if (!ns) {
        $toast.toastDanger(t('notification.namespace.loadFailed'))
        router.push({ name: 'root' })
        return
      }

      if (!ns.enabled) {
        $toast.toastWarning(t('notification.namespace.disabled'))
        router.push({ name: 'root' })
        return
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
 * Similar to old Corteza's prepareNamespace pattern.
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
      moduleStore.load({
        namespaceID: nsID,
        force: true,
      }),
      pageStore.load({
        namespaceID: nsID,
        force: true,
      }),
      chartStore.load({
        namespaceID: nsID,
        force: true,
      }),
      pageLayoutStore.load({
        namespaceID: nsID,
        force: true,
      }),
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
    router.replace({ name: 'page', params: { slug: props.slug, pageID: firstPage.pageID } })
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
