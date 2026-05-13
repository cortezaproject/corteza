<template>
  <div class="h-full overflow-hidden">
    <CResourceList
      ref="resourceListRef"
      primary-key="id"
      :fields="fields"
      :items="sessions"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :loading="loading"
      :translations="{
        searchPlaceholder: $t('chatbot.sessions.searchPlaceholder'),
        showingPagination: 'general.resourceList.pagination.showing',
        singlePluralPagination: 'general.resourceList.pagination.single',
        prevPagination: $t('general.resourceList.pagination.prev'),
        nextPagination: $t('general.resourceList.pagination.next'),
        recordsPerPage: $t('general.resourceList.pagination.recordsPerPage'),
        resourceSingle: $t('chatbot.sessions.title'),
        resourcePlural: $t('chatbot.sessions.title'),
      }"
      clickable
      class="h-full"
      @update:filter="Object.assign(filter, $event)"
      @sort="handleSort"
      @row-click="handleRowClick"
      @page-change="handlePageChange"
    >
      <template v-if="!chatbotId" #filter>
        <Button
          icon="pi pi-filter"
          :class="{ 'text-primary': filter.chatbotID }"
          severity="secondary"
          size="small"
          text
          @click="filterMenu.toggle($event)"
        />
      </template>

      <template v-if="!chatbotId" #body-chatbot="{ data }">
        {{ chatbotName(data.chatbotID) }}
      </template>

      <template #body-status="{ data }">
        <Tag
          :value="$t(`chatbot.sessions.status.${data.status}`, data.status)"
          :severity="statusSeverity(data.status)"
          rounded
        />
      </template>

      <template #body-createdAt="{ data }">
        {{ locFullDateTime(data.createdAt) }}
      </template>
    </CResourceList>

    <Popover ref="filterMenu">
      <div class="flex flex-col gap-3 p-2 w-64">
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium text-color">{{ $t('chatbot.sessions.filter.chatbot') }}</label>
          <Select
            v-model="selectedChatbotID"
            :options="[{ chatbotID: '', name: $t('chatbot.sessions.filter.allChatbots') }, ...chatbotStore.list]"
            option-label="name"
            option-value="chatbotID"
            class="w-full"
            @change="applyChatbotFilter"
          />
        </div>
      </div>
    </Popover>

    <SessionDetail
      v-model:visible="showDetail"
      :session="selectedSession"
    />
  </div>
</template>

<script setup>
import { components, filters, useResourceList } from '@planetcrust/human-vue'
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useChatbotStore } from '@/stores/chatbot'
import SessionDetail from './SessionDetail.vue'

const props = defineProps({
  chatbotId: {
    type: String,
    default: '',
  },
})

const { CResourceList } = components
const { locFullDateTime } = filters
const { t } = useI18n()
const $SystemAPI = inject('$SystemAPI')
const chatbotStore = useChatbotStore()

const filterMenu = ref()
const selectedSession = ref(null)
const showDetail = ref(false)
const selectedChatbotID = ref(props.chatbotId || '')

const fields = computed(() => [
  ...(props.chatbotId ? [] : [{ key: 'chatbot', header: t('chatbot.sessions.columns.chatbot') }]),
  { key: 'status', sortable: true, header: t('chatbot.sessions.columns.status') },
  {
    key: 'createdAt',
    sortable: true,
    header: t('chatbot.sessions.columns.createdAt'),
    class: 'text-right',
    pt: { columnHeaderContent: 'justify-end' },
  },
])

function chatbotName(chatbotID) {
  const cb = chatbotStore.list.find(c => c.chatbotID === chatbotID)
  return cb?.name || cb?.handle || chatbotID
}

const apiCall = computed(() => {
  if (props.chatbotId) {
    return params => $SystemAPI.chatbotSessionListByChatbotCancellable({ ...params, chatbotID: props.chatbotId })
  }
  return params => $SystemAPI.chatbotSessionListCancellable(params)
})

const {
  items: sessions,
  loading,
  filter,
  sorting,
  pagination,
  handleSort,
  handlePageChange,
  filterList,
} = useResourceList(
  params => apiCall.value(params),
  {
    filter: { query: '', chatbotID: props.chatbotId || undefined },
    sorting: { sortBy: 'createdAt', sortDesc: true },
    pagination: { limit: 50 },
  },
)

function handleRowClick({ data }) {
  selectedSession.value = data
  showDetail.value = true
}

function applyChatbotFilter() {
  filter.chatbotID = selectedChatbotID.value || undefined
  filterList()
}

function statusSeverity(status) {
  switch (status) {
    case 'active': return 'success'
    case 'completed': return 'secondary'
    case 'handoff': return 'warn'
    case 'error': return 'danger'
    default: return 'secondary'
  }
}
</script>
