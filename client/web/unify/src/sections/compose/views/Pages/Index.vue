<template>
  <div v-if="!redirecting" class="flex flex-col items-center justify-center h-full gap-6 p-8">
    <div class="text-center max-w-lg">
      <h1 class="text-2xl font-semibold mb-2">{{ $t('onboarding.label.welcome') }}</h1>
      <p class="text-muted-color mb-6">{{ $t('onboarding.message.noPages') }}</p>

      <div class="flex flex-col gap-3 items-center">
        <Button
          v-if="namespace?.canCreateModule"
          :label="$t('onboarding.step.module.create')"
          icon="pi pi-database"
          @click="goToModuleAdmin"
        />
        <Button
          v-if="namespace?.canCreatePage"
          :label="$t('onboarding.step.page.create')"
          icon="pi pi-objects-column"
          severity="secondary"
          @click="goToPageAdmin"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { usePageStore } from '@planetcrust/human-vue'
import { NoID } from '@planetcrust/human-js'
import { computed, watchEffect } from 'vue'
import { useRoute, useRouter } from 'vue-router'

defineProps({
  namespace: {
    type: Object,
    required: true,
  },
})

const router = useRouter()
const route = useRoute()
const pageStore = usePageStore()

const homePage = computed(
  () =>
    [...pageStore.set]
      .filter(p => p.visible && p.selfID === NoID && p.moduleID === NoID)
      .sort((a, b) => a.weight - b.weight)[0],
)

const redirecting = computed(() => !!homePage.value)

watchEffect(() => {
  if (homePage.value) {
    router.replace({
      name: 'page',
      params: { slug: route.params.slug, pageID: homePage.value.pageID },
    })
  }
})

function goToPageAdmin() {
  router.push({ name: 'admin.pages' })
}

function goToModuleAdmin() {
  router.push({ name: 'admin.modules' })
}
</script>
