<template>
  <!-- Teleport sidebar content to CSidebar targets -->
  <Teleport to="#sidebar-header-expanded" :defer="true">
    <!-- Namespace Switcher -->
    <Select
      :model-value="currentNamespaceObject"
      :options="enabledNamespaces"
      option-label="name"
      data-key="namespaceID"
      :placeholder="$t('sidebar.pickNamespace')"
      class="w-full"
      @update:model-value="handleNamespaceChange"
    >
      <template #option="{ option }">
        <div class="flex items-center gap-2">
          <span>{{ option.name }}</span>
        </div>
      </template>
    </Select>
  </Teleport>

  <!-- Navigation Items at the bottom of sidebar -->
  <Teleport to="#sidebar-footer-expanded" :defer="true">
    <div v-if="namespace" class="flex flex-col gap-1 p-2">
      <Divider class="my-0" />
      <Button
        :label="$t('sidebar.modules')"
        icon="pi pi-database"
        severity="secondary"
        text
        class="justify-start"
        :class="{ 'bg-primary/10': isModulesRoute }"
        @click="goToModules"
      />
      <Button
        :label="$t('sidebar.pages')"
        icon="pi pi-file"
        severity="secondary"
        text
        class="justify-start"
        :class="{ 'bg-primary/10': isPagesRoute }"
        @click="goToPages"
      />
      <Button
        :label="$t('sidebar.charts')"
        icon="pi pi-chart-bar"
        severity="secondary"
        text
        class="justify-start"
        :class="{ 'bg-primary/10': isChartsRoute }"
        @click="goToCharts"
      />
    </div>
  </Teleport>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const props = defineProps({
  namespace: {
    type: Object,
    required: true,
  },
  namespaces: {
    type: Array,
    default: () => [],
  },
})

const route = useRoute()
const router = useRouter()

// Computed
const enabledNamespaces = computed(() => {
  return props.namespaces.filter(ns => ns.enabled)
})

const currentNamespaceSlug = computed(() => {
  return route.params.slug || null
})

// For Select model-value - needs the full namespace object
// Match by slug first, fall back to namespaceID (some namespaces have empty slug)
const currentNamespaceObject = computed(() => {
  const urlPart = currentNamespaceSlug.value
  return enabledNamespaces.value.find(ns => 
    ns.slug === urlPart || ns.namespaceID === urlPart
  ) || null
})

const isAdminRoute = computed(() => {
  return route.name?.toString().includes('admin.')
})

const isModulesRoute = computed(() => {
  return route.name?.toString().startsWith('admin.modules')
})

const isPagesRoute = computed(() => {
  return route.name?.toString().startsWith('admin.pages')
})

const isChartsRoute = computed(() => {
  return route.name?.toString().startsWith('admin.charts')
})

// Methods
function handleNamespaceChange(namespace) {
  // Use slug if available, otherwise fall back to namespaceID
  const newSlug = namespace?.slug || namespace?.namespaceID
  
  if (newSlug && newSlug !== currentNamespaceSlug.value) {
    const routeName = isAdminRoute.value ? 'admin.modules' : 'pages'
    router.push({ name: routeName, params: { slug: newSlug } })
  }
}

function goToModules() {
  if (currentNamespaceSlug.value) {
    router.push({ name: 'admin.modules' })
  }
}

function goToPages() {
  if (currentNamespaceSlug.value) {
    router.push({ name: 'admin.pages' })
  }
}

function goToCharts() {
  if (currentNamespaceSlug.value) {
    router.push({ name: 'admin.charts' })
  }
}
</script>
