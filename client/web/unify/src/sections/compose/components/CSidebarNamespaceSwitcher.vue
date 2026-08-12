<template>
  <div class="flex items-center gap-1 my-1">
    <FloatLabel variant="on" class="flex-1 min-w-0">
      <Select
        id="namespace-selector"
        :model-value="currentNamespaceObject"
        :options="selectOptions"
        option-label="name"
        data-key="namespaceID"
        :placeholder="$t('sidebar.namespaceSelector.placeholder')"
        size="small"
        class="w-full"
        @update:model-value="handleNamespaceChange"
      >
        <template #option="{ option }">
          <div class="flex items-center gap-2">
            <span>{{ option.name }}</span>
          </div>
        </template>
        <template v-if="canManageNamespaces" #footer>
          <RouterLink
            :to="{ name: 'namespace.list' }"
            class="block p-2 text-sm text-muted-color hover:text-primary transition-colors text-center border-t"
          >
            {{ $t('sidebar.namespaceSelector.manage') }}
          </RouterLink>
        </template>
      </Select>
      <label for="namespace-selector">{{ $t('sidebar.namespaceSelector.label') }}</label>
    </FloatLabel>

    <RouterLink
      v-if="currentNamespaceObject?.canUpdateNamespace"
      :to="{
        name: 'namespace.edit',
        params: { slug: currentNamespaceObject.slug || currentNamespaceObject.namespaceID },
      }"
      v-tooltip.top="$t('sidebar.editNamespace')"
    >
      <Button icon="pi pi-pencil" size="small" text rounded severity="secondary" />
    </RouterLink>
  </div>
</template>

<script setup>
import { useNamespaceStore } from '@planetcrust/human-vue'
import { computed, inject } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const $Settings = inject('$Settings')
const namespaceStore = useNamespaceStore()

const enabledNamespaces = computed(() => {
  return namespaceStore.set.filter(ns => ns.enabled).sort((a, b) => a.name.localeCompare(b.name))
})

const selectOptions = computed(() => {
  const current = currentNamespaceObject.value
  if (current && !current.enabled) {
    return [current, ...enabledNamespaces.value]
  }
  return enabledNamespaces.value
})

const canManageNamespaces = computed(() => {
  const { hideNamespaceListLink } = $Settings.get('compose.ui.sidebar', {})
  if (hideNamespaceListLink) return false
  return namespaceStore.set.some(ns => ns.canManageNamespace)
})

const currentNamespaceSlug = computed(() => {
  return route.params.slug || null
})

const currentNamespaceObject = computed(() => {
  const urlPart = currentNamespaceSlug.value
  return namespaceStore.set.find(ns => ns.slug === urlPart || ns.namespaceID === urlPart) || null
})

const isAdminRoute = computed(() => {
  return route.name?.toString().includes('admin.')
})

function handleNamespaceChange(namespace) {
  const newSlug = namespace?.slug || namespace?.namespaceID

  if (newSlug && newSlug !== currentNamespaceSlug.value) {
    const routeName = isAdminRoute.value ? 'admin.modules' : 'pages'
    router.push({ name: routeName, params: { slug: newSlug } })
  }
}
</script>
