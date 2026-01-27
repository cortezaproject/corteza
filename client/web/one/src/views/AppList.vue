<template>
  <div class="app-selector flex flex-col h-full py-4 gap-7 my-3">
    <div class="flex flex-col justify-center items-center mx-4 gap-4">
      <img v-if="logoUrl" :src="logoUrl" class="px-4 max-h-lg max-w-xl w-auto mb-6" alt="Logo" />

      <CInputSearch
        v-model="query"
        placeholder="Search applications..."
        class="w-full max-w-2xl mx-auto"
      />
    </div>

    <div v-if="areAppsVisible" class="flex-1 overflow-auto">
      <div class="container mx-auto">
        <div class="flex flex-wrap justify-center gap-7 p-4 mx-2">
          <a
            v-for="app in apps"
            :key="app.applicationID"
            :href="app.enabled ? app.unify.url : '#'"
            target="_self"
            class="block"
            @click="!app.enabled && $event.preventDefault()"
            v-show="isAppVisible(app)"
          >
            <Card
              :pt="{
                body: {
                  class: 'grow justify-center gap-0 py-1',
                },
                title: {
                  class: 'text-center line-clamp-2 group-hover:line-clamp-none',
                },
              }"
              class="group cursor-pointer hover:shadow-lg hover:scale-105 hover:text-primary transition-all duration-100 w-80 min-h-72 hover:h-full overflow-hidden"
            >
              <template #header>
                <img
                  :src="getAppLogoUrl(app)"
                  :alt="app.unify?.name || app.name"
                  class="w-full h-full object-contain"
                  loading="lazy"
                  decoding="async"
                />
              </template>

              <template #title>
                {{ app.unify?.name || app.name }}
              </template>
            </Card>
          </a>
        </div>
      </div>
    </div>

    <div v-else class="flex justify-center items-center mt-20 w-full">
      <label class="text-muted-color text-lg">
        {{ query ? $t('layout.no-applications-found') : $t('layout.no-applications') }}
      </label>
    </div>
  </div>
</template>

<script setup>
import { useApplicationsStore } from '@/stores/applications'
import { components } from '@cortezaproject/corteza-vue-next'
import { computed, inject, ref } from 'vue'
const { CInputSearch } = components

const query = ref('')
const applicationsStore = useApplicationsStore()

const $Settings = inject('$Settings')
const $SystemAPI = inject('$SystemAPI')

const apps = computed(() => {
  return applicationsStore.unifyOnly
})

const normalizedQuery = computed(() => (query.value || '').trim().toUpperCase())

const isAppVisible = app => {
  const q = normalizedQuery.value
  if (!q) return true
  return (
    (app.name?.toUpperCase() || '').includes(q) ||
    (app.unify?.name?.toUpperCase() || '').includes(q)
  )
}

const areAppsVisible = computed(() => apps.value.some(isAppVisible))

const logoUrl = computed(() => {
  return $Settings.attachment('ui.mainLogo')
})

const getAppLogoUrl = app => {
  if (!app.unify?.logo) {
    return 'applications/default-app.png'
  }

  const apiSystem = '/api/system'
  const apiBaseUrl = new URL($SystemAPI.baseURL).toString()

  if (app.unify.logo.startsWith(apiSystem)) {
    return apiBaseUrl.substring(0, apiBaseUrl.length - apiSystem.length) + app.unify.logo
  }

  return app.unify.logo
}
</script>
