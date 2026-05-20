<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('chatbot.sessions.title') }}</span>
  </Teleport>

  <div class="container mx-auto p-4 h-full overflow-hidden min-w-0 flex flex-col">
    <CChatbotInbox
      :chatbot-i-ds="chatbotIDs"
      :status-filter="['handoff_requested', 'handoff_active', 'active', 'closed']"
      :show-filter="true"
      :refresh-rate="5"
      :translations="inboxTranslations"
    />
  </div>
</template>

<script setup>
import { CChatbotInbox, makeChatbotInboxTranslations } from '@planetcrust/human-vue'
import { computed, inject, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useChatbotStore } from '@/stores/chatbot'

const { t } = useI18n()
const $SystemAPI = inject('$SystemAPI')
const chatbotStore = useChatbotStore()

// Cross-chatbot session view: watch every chatbot the user can see. The
// inbox component already issues a single list call internally and filters
// client-side, so passing the full set is the cheapest way to feed it.
const chatbotIDs = computed(() => chatbotStore.list.map(c => c.chatbotID))

const inboxTranslations = computed(() => makeChatbotInboxTranslations(t, 'chatbot.inbox.'))

onMounted(() => {
  if (!chatbotStore.list.length) {
    void chatbotStore.fetchList($SystemAPI)
  }
})
</script>
