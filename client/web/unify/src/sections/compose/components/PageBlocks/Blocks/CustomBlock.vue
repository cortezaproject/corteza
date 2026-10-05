<template>
  <PageBlock :block="block" :record="record" @refreshBlock="load">
    <div class="h-full flex flex-col" data-test-id="custom-block">
      <div v-if="loading" class="h-full flex items-center justify-center">
        <ProgressSpinner style="width: 2rem; height: 2rem" />
      </div>

      <div
        v-else-if="problem"
        class="flex items-center justify-center h-full p-3 text-muted-color italic"
        data-test-id="custom-block-problem"
      >
        {{ problem }}
      </div>

      <CustomAppFrame
        v-else
        :key="frameKey"
        class="flex-1"
        :source="source"
        :source-meta="sourceMeta"
        :name="name"
        :context="context"
        :resizable="false"
        @navigated="problem = $t('app.state.navigated')"
        @changed="$eventBus?.emit('refetch-records')"
      />
    </div>
  </PageBlock>
</template>

<script setup>
// A custom HTML block: a custom application's page, or one of the block's
// own, in the custom app sandbox. See app.intent.md for what the frame keeps.
import { computed, inject, ref, watch } from 'vue'
import { useApplicationsStore } from '@planetcrust/human-vue'
import { useI18n } from 'vue-i18n'
import PageBlock from './PageBlock.vue'
import CustomAppFrame from '@/sections/app/components/CustomAppFrame.vue'

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const { t } = useI18n()
const $SystemAPI = inject('$SystemAPI')
// The rest of the page listens here for data changed under it.
const $eventBus = inject('$eventBus', null)
const applicationsStore = useApplicationsStore()

const loading = ref(true)
const problem = ref('')
const source = ref('')
const sourceMeta = ref({})
const name = ref('')
const frameKey = ref(0)

const applicationID = computed(() => {
  const id = props.block.options?.applicationID
  return id && id !== '0' ? String(id) : ''
})

// Where the block is shown, for `human.context()`.
const context = computed(() => ({
  namespaceID: props.namespace?.namespaceID || '',
  namespace: props.namespace?.slug || '',
  pageID: props.page?.pageID || '',
  moduleID: props.page?.moduleID && props.page.moduleID !== '0' ? props.page.moduleID : '',
  recordID: props.record?.recordID && props.record.recordID !== '0' ? props.record.recordID : '',
  params: props.block.options?.params || {},
}))

async function fromApplication(id) {
  const app = await applicationsStore.findByID(id).catch(() => null)
  if (!app || app.unify?.kind !== 'custom' || !app.canAccessApplication) {
    problem.value = t('block.custom.refused')
    return
  }

  const read = await $SystemAPI.applicationSourceRead({ applicationID: id })
  if (!read.source) {
    problem.value = t('app.state.empty')
    return
  }

  name.value = app.unify?.name || app.name || ''
  sourceMeta.value = read.sourceMeta || {}
  source.value = read.source
}

function fromBlock() {
  const o = props.block.options || {}
  if (!o.source) {
    problem.value = t('block.custom.noSource')
    return
  }

  name.value = props.block.title || ''
  // A page written into the block reads the namespace the block is shown in.
  sourceMeta.value = {
    namespace: props.namespace?.slug || '',
    namespaceID: props.namespace?.namespaceID || '',
    modules: o.modules || [],
    moduleIDs: o.moduleIDs || {},
    writes: o.writes || [],
    origins: o.origins || [],
    automations: o.automations || [],
    chatbots: o.chatbots || [],
  }
  source.value = o.source
}

async function load() {
  loading.value = true
  problem.value = ''
  source.value = ''

  try {
    if (applicationID.value) await fromApplication(applicationID.value)
    else fromBlock()
  } catch (error) {
    problem.value = t('app.state.error', { reason: error?.message || String(error) })
  } finally {
    loading.value = false
    frameKey.value++
  }
}

watch(
  () => [
    applicationID.value,
    props.block.options?.source,
    props.namespace?.namespaceID,
    props.block.options?.modules,
    props.block.options?.writes,
    props.block.options?.origins,
    props.block.options?.automations,
    props.block.options?.chatbots,
  ],
  load,
  { immediate: true, deep: true },
)
</script>
