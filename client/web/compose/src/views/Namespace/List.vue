<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('namespace.title') }}</span>
  </Teleport>
  <Teleport to="#topbar-tools" defer>
    <CRouterLinkButton
      :to="{ name: 'namespace.manage' }"
      :label="$t('namespace.manage-view.label')"
      size="small"
    />
  </Teleport>

  <div class="flex flex-col h-full py-4 gap-7 my-3">
    <div class="flex items-center justify-center gap-3 mx-4">
      <CInputSearch
        v-model="query"
        :placeholder="$t('namespace.searchPlaceholder')"
        class="w-full max-w-2xl"
      />
    </div>

    <div v-if="areNamespacesVisible" class="flex-1 overflow-auto">
      <div class="container mx-auto p-4">
        <div class="flex flex-wrap justify-center gap-7 p-4 xl:mx-32">
          <RouterLink
            v-for="namespace in namespaceStore.set"
            :key="namespace.namespaceID"
            :to="{
              name: 'namespace.view',
              params: { slug: namespace.slug || namespace.namespaceID },
            }"
            class="block"
            v-show="isNamespaceVisible(namespace)"
          >
            <Card
              :pt="{
                body: {
                  class: 'grow justify-center gap-0 py-1',
                },
                title: {
                  class: 'text-center line-clamp-2 group-hover:line-clamp-none',
                },
                subtitle: {
                  class: 'text-center line-clamp-2 group-hover:line-clamp-none',
                },
              }"
              class="group cursor-pointer hover:shadow-lg hover:scale-105 hover:text-primary transition-all duration-100 w-80 min-h-72 hover:h-full overflow-hidden"
            >
              <template #header>
                <div class="flex items-center justify-center w-full h-full pt-7 shrink-0">
                  <Avatar
                    :label="namespace.meta.logoEnabled ? null : namespace.initials"
                    :image="
                      namespace.meta.logoEnabled
                        ? namespace.meta.logo || $Settings.attachment('ui.mainLogo')
                        : null
                    "
                    :pt="{
                      image: { class: 'object-contain' },
                    }"
                    shape="circle"
                    size="xlarge"
                    :class="{ 'text-muted-color bg-emphasis': !namespace.meta.logoEnabled }"
                    class="w-32 h-32 font-bold"
                  />
                </div>
              </template>

              <template #title>
                {{ namespace.name }}
              </template>

              <template v-if="namespace.meta.description" #subtitle>
                {{ namespace.meta.description }}
              </template>
            </Card>
          </RouterLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useNamespaceStore } from '@/stores/namespace'
import { components } from '@cortezaproject/corteza-vue-next'
import { computed, ref } from 'vue'
const { CInputSearch, CRouterLinkButton } = components

const namespaceStore = useNamespaceStore()

const query = ref('')

const normalizedQuery = computed(() => (query.value || '').trim().toUpperCase())

const isNamespaceVisible = namespace => {
  const q = normalizedQuery.value
  if (!q) return true
  return (namespace.slug + namespace.name).toUpperCase().indexOf(q) > -1
}

const areNamespacesVisible = computed(() => namespaceStore.set.some(isNamespaceVisible))
</script>
