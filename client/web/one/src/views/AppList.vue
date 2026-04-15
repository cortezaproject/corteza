<template>
  <div class="app-selector flex flex-col md:align-items-center h-full py-4 gap-7 my-3 md:mt-7">
    <div class="flex flex-col justify-center items-center mx-4 gap-4">
      <img v-if="logoUrl" :src="logoUrl" class="px-4 max-h-lg max-w-xl w-auto mb-6" alt="Logo" />

      <CInputSearch
        v-model="query"
        placeholder="Search applications..."
        class="w-full max-w-2xl mx-auto"
      />
    </div>

    <div class="flex-1 overflow-auto">
      <div class="container mx-auto p-4">
        <CAppList
          :query="query"
          variant="grid"
          :no-apps-text="$t('layout.no-applications')"
          :no-results-text="$t('layout.no-applications-found')"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { components } from '@cortezaproject/corteza-vue-next'
import { computed, inject, ref } from 'vue'

const { CInputSearch, CAppList } = components

const query = ref('')

const $Settings = inject('$Settings')

const logoUrl = computed(() => $Settings.attachment('ui.mainLogo'))
</script>
