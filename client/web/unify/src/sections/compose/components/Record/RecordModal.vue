<template>
  <Dialog
    v-model:visible="showModal"
    modal
    :draggable="false"
    :dismissableMask="true"
    :showHeader="true"
    :style="{ width: '90vw', height: '90vh' }"
    class="record-modal"
    content-class="h-full p-0 !pb-0"
    @update:visible="onVisibleChange"
  >
    <template #header>
      <span class="font-bold text-lg truncate">{{ pageTitle }}</span>
    </template>
    <template #closebutton="{ closeCallback }">
      <Button
        v-if="recordPageID && recordID"
        v-tooltip.bottom="$t('block.record.openFullPage')"
        icon="pi pi-external-link"
        text
        severity="secondary"
        size="small"
        @click="openFullPage"
      />
      <Button
        icon="pi pi-times"
        text
        severity="secondary"
        size="small"
        @click="closeCallback"
      />
    </template>
    
    <div v-if="showModal && recordPageID && recordID && namespace" class="h-full">
      <RecordView
        :namespace="namespace"
        :in-modal="true"
        :modal-page-i-d="recordPageID"
        :modal-record-i-d="recordID"
        @close="closeModal"
      />
    </div>
  </Dialog>
</template>

<script setup>
import { computed, inject, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import RecordView from '@/sections/compose/views/Pages/RecordView.vue'

const props = defineProps({
  namespace: {
    type: Object,
    required: true
  }
})

const route = useRoute()
const router = useRouter()
const $pageStore = inject('$pageStore', null)

const showModal = ref(false)

const recordID = computed(() => route.query.recordID)
const recordPageID = computed(() => route.query.recordPageID)
const pageTitle = computed(() => {
  if (!recordPageID.value || !$pageStore) return ''
  const page = $pageStore.getByID(recordPageID.value)
  return page?.title || ''
})

watch(
  () => [route.query.recordID, route.query.recordPageID],
  ([newRecordID, newRecordPageID]) => {
    if (newRecordID && newRecordPageID) {
      showModal.value = true
    } else {
      showModal.value = false
    }
  },
  { immediate: true }
)

function onVisibleChange(isVisible) {
  if (!isVisible) {
    closeModal()
  }
}

function closeModal() {
  const q = { ...route.query }
  delete q.recordID
  delete q.recordPageID
  delete q.edit
  delete q.cloneFromID
  router.push({ query: q })
}

function openFullPage() {
  router.push({
    name: 'page.record',
    params: {
      slug: props.namespace.slug,
      pageID: recordPageID.value,
      recordID: recordID.value,
    },
  })
}
</script>

<style>
.record-modal .p-dialog-content {
  overflow: hidden;
  display: flex;
  flex-direction: column;
  background: var(--body-bg);
}
</style>
