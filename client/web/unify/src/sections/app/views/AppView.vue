<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ application?.name || $t('app.title') }}</span>
  </Teleport>

  <div class="h-full w-full overflow-auto">
    <div v-if="loading" class="h-full flex items-center justify-center">
      <ProgressSpinner style="width: 2.5rem; height: 2.5rem" />
    </div>

    <Message v-else-if="problem" severity="error" class="m-4">
      {{ problem }}
    </Message>

    <template v-else>
      <Message
        v-if="application?.enabled === false"
        severity="info"
        :closable="false"
        class="m-4"
        data-test-id="custom-app-preview"
      >
        {{ $t('app.state.preview') }}
      </Message>

      <CustomAppFrame
        :source="source"
        :source-meta="sourceMeta"
        :name="application?.name || ''"
        @navigated="problem = $t('app.state.navigated')"
      />
    </template>
  </div>
</template>

<script setup>
import { useApplicationsStore } from '@planetcrust/human-vue'
import { inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import CustomAppFrame from '../components/CustomAppFrame.vue'

const { t } = useI18n()
const route = useRoute()

const $SystemAPI = inject('$SystemAPI')

const applicationsStore = useApplicationsStore()

const loading = ref(true)
const problem = ref('')
const application = ref(null)
const source = ref('')
const sourceMeta = ref({})

// One application per route, and the route changes without the view being
// remounted when the user goes from one custom app straight to another.
async function load(applicationID) {
  loading.value = true
  problem.value = ''
  application.value = null
  source.value = ''
  sourceMeta.value = {}

  try {
    application.value = await applicationsStore.findByID(applicationID)
  } catch {
    application.value = null
  }

  // The route guard has already refused what this catches; a plain message is
  // for the case where the view is reached some other way.
  if (
    !application.value ||
    application.value.unify?.kind !== 'custom' ||
    !application.value.canAccessApplication
  ) {
    problem.value = t('app.state.refused')
    loading.value = false
    return
  }

  try {
    const read = await $SystemAPI.applicationSourceRead({ applicationID })

    if (!read.source) {
      problem.value = t('app.state.empty')
      return
    }

    sourceMeta.value = read.sourceMeta || {}
    source.value = read.source
  } catch (error) {
    problem.value = t('app.state.error', { reason: error?.message || String(error) })
  } finally {
    loading.value = false
  }
}

watch(
  () => String(route.params.applicationID || ''),
  applicationID => {
    if (applicationID) load(applicationID)
  },
  { immediate: true },
)
</script>
