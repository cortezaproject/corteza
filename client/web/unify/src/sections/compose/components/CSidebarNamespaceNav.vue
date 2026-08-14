<template>
  <div class="flex flex-col h-full">
    <CInputSearch
      id="namespace-nav-search"
      v-model="searchQuery"
      :placeholder="$t('namespace.searchPlaceholder')"
      size="small"
      class="my-2"
    />

    <div class="flex-1 overflow-auto">
      <CSidebarNav
        :items="navItems"
        id-key="_id"
        parent-key="_parentId"
        label-key="_label"
        icon-key="_icon"
        divider-key="_divider"
        route-key="_route"
        badge-key="_badge"
        match-type="exact"
        expand-all
      />
    </div>

    <div
      v-if="hasSearch && !filteredNamespaces.length"
      class="flex-1 flex items-center justify-center text-muted-color text-sm"
    >
      {{ $t('sidebar.noResults') }}
    </div>
  </div>
</template>

<script setup>
import { useNamespaceStore } from '@planetcrust/human-vue'
import { components } from '@planetcrust/human-vue'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { CSidebarNav, CInputSearch } = components
const { t } = useI18n()

const namespaceStore = useNamespaceStore()

const searchQuery = ref('')

const hasSearch = computed(() => searchQuery.value.trim().length > 0)

// Same match as the namespace list's card search: name and short name together.
const normalizedQuery = computed(() => searchQuery.value.trim().toUpperCase())

// Every namespace the user can list, disabled ones included: this nav belongs to
// the namespace editor, and a disabled namespace is exactly what gets configured
// there. (The switcher on the namespace's own screens filters them out — it
// navigates into a namespace, which a disabled one refuses.)
const filteredNamespaces = computed(() => {
  const query = normalizedQuery.value

  return [...namespaceStore.set]
    .filter(ns => !query || `${ns.slug || ''}${ns.name || ''}`.toUpperCase().includes(query))
    .sort((a, b) =>
      (a.name || a.slug || '').localeCompare(b.name || b.slug || '', undefined, {
        sensitivity: 'base',
      }),
    )
})

const navItems = computed(() => [
  {
    _id: 'all',
    _parentId: '0',
    _label: t('sidebar.allNamespaces'),
    _icon: 'pi pi-th-large',
    _route: { name: 'namespace.list' },
  },
  ...(filteredNamespaces.value.length
    ? [
        {
          _id: 'namespaces',
          _parentId: '0',
          _label: t('general.label.namespace.plural'),
          _icon: 'pi pi-folder',
          _divider: true,
        },
      ]
    : []),
  ...filteredNamespaces.value.map(ns => ({
    _id: ns.namespaceID,
    _parentId: 'namespaces',
    _label: ns.name || ns.slug || ns.namespaceID,
    _route: {
      name: 'namespace.edit',
      params: { slug: ns.slug || ns.namespaceID },
    },
    _badge: ns.enabled
      ? undefined
      : { value: t('namespace.manage.disabled'), severity: 'secondary' },
  })),
])
</script>
