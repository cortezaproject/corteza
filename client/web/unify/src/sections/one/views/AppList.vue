<template>
  <div class="flex flex-col h-full py-4 gap-7 my-3 md:mt-7">
    <div class="flex flex-col justify-center items-center mx-4 gap-4">
      <img v-if="logoUrl" :src="logoUrl" class="px-4 max-h-lg max-w-xl w-auto mb-6" alt="Logo" />

      <CInputSearch
        v-model="query"
        :placeholder="$t('one.apps.search')"
        class="w-full max-w-2xl mx-auto"
      />
    </div>

    <div class="flex-1 overflow-auto">
      <div class="container mx-auto p-4">
        <CAppList
          variant="grid"
          :query="query"
          :no-apps-text="$t('one.apps.empty')"
          :no-results-text="$t('one.apps.noResults')"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { components } from '@planetcrust/human-vue'
import { computed, inject, ref } from 'vue'

const { CInputSearch, CAppList } = components

const query = ref('')

const $Settings = inject('$Settings')

const logoUrl = computed(() => $Settings.attachment('ui.mainLogo'))
</script>
