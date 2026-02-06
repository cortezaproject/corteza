<template>
  <!-- Loading state -->
  <div v-if="loading" class="flex flex-col items-center justify-center h-full gap-4">
    <h1 class="text-2xl text-muted-color">
      {{ namespace?.name || namespace?.slug || '...' }}
    </h1>
    <div class="flex items-center gap-3">
      <ProgressSpinner style="width: 24px; height: 24px" />
      <span class="text-lg text-muted-color">{{ $t('general.loading') }}</span>
    </div>
  </div>

  <!-- Loaded content -->
  <div v-else-if="namespace" class="flex h-full">
    <!-- Namespace Sidebar (teleports content to CSidebar) -->
    <CNamespaceSidebar :namespace="namespace" :namespaces="namespaceStore.set" />

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
import CNamespaceSidebar from '@/components/CNamespaceSidebar.vue'
import { useChartStore } from '@/stores/chart'
import { useModuleStore } from '@/stores/module'
import { useNamespaceStore } from '@/stores/namespace'
import { usePageStore } from '@/stores/page'
import { compose } from '@cortezaproject/corteza-js-next'
import { inject, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const props = defineProps({
  slug: {
    type: String,
    required: true,
  },
})

const router = useRouter()
const { t } = useI18n()
const $toast = inject('$toast')

// Stores
const namespaceStore = useNamespaceStore()
const moduleStore = useModuleStore()
const pageStore = usePageStore()
const chartStore = useChartStore()

// State
const loading = ref(true)
const namespace = ref(null)

// Methods
async function loadNamespace() {
  loading.value = true

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
  } catch (error) {
    console.error('Failed to load namespace:', error)
    $toast.toastDanger(t('notification.namespace.loadFailed'))
    loading.value = false
  }
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
    ])
  } catch (error) {
    console.error('Failed to prepare namespace context:', error)
    // Continue anyway - individual stores will handle their errors
  } finally {
    loading.value = false
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
