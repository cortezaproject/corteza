<template>
  <PageBlock :block="block" @refreshBlock="onManualRefresh">
    <CChatbotInbox
      ref="inboxRef"
      :chatbot-i-ds="options.chatbotIDs"
      :status-filter="options.statusFilter"
      :auto-open-first="options.autoOpenFirst"
      :show-filter="!!options.showFilter"
      :translations="translations"
    />
  </PageBlock>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { CChatbotInbox, makeChatbotInboxTranslations } from '@planetcrust/human-vue'
import PageBlock from './PageBlock.vue'

// Thin Compose wrapper around the shared CChatbotInbox. All operator UX +
// data flow lives in the lib component; this file owns the PageBlock chrome,
// forwards manual-refresh events, and resolves user-facing strings from the
// compose locale tree so the lib component itself stays i18n-agnostic.
const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const { t } = useI18n()
const options = computed(() => props.block.options)

// Reactive so changes to the active language flow through without remount.
const translations = computed(() => makeChatbotInboxTranslations(t, 'block.chatbotInbox.'))

const inboxRef = ref(null)
function onManualRefresh() {
  inboxRef.value?.refresh?.()
}
</script>
