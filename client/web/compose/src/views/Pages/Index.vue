<template>
  <div v-if="!redirecting" class="flex flex-col items-center justify-center h-full gap-6 p-8">
    <div class="text-center max-w-lg">
      <h1 class="text-2xl font-semibold mb-2">{{ $t('onboarding.label.welcome') }}</h1>
      <p class="text-muted-color mb-6">{{ $t('onboarding.message.noPages') }}</p>

      <div class="flex flex-col gap-3 items-center">
        <Button
          :label="$t('onboarding.step.page.create')"
          icon="pi pi-objects-column"
          @click="goToPageAdmin"
        />
        <Button
          :label="$t('onboarding.step.module.create')"
          icon="pi pi-database"
          severity="secondary"
          @click="goToModuleAdmin"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { usePageStore } from '@/stores/page'
import { NoID } from '@cortezaproject/corteza-js-next'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  namespace: {
    type: Object,
    required: true,
  },
})

const router = useRouter()
const pageStore = usePageStore()

const redirecting = ref(true)

onMounted(() => {
  const pages = pageStore.set

  // Find home page: first visible, first-level, non-record page, sorted by weight
  const homePage = [...pages]
    .filter(p => p.visible && p.selfID === NoID && p.moduleID === NoID)
    .sort((a, b) => a.weight - b.weight)[0]

  if (homePage) {
    router.replace({ name: 'page', params: { pageID: homePage.pageID } })
  } else {
    redirecting.value = false
  }
})

function goToPageAdmin() {
  router.push({ name: 'admin.pages' })
}

function goToModuleAdmin() {
  router.push({ name: 'admin.modules' })
}
</script>
