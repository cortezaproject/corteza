<template>
  <div class="flex flex-col h-full">
    <CSidebarNav
      :items="navItems"
      id-key="_id"
      parent-key="_parentId"
      label-key="_label"
      icon-key="_icon"
      divider-key="_divider"
      route-key="_route"
      match-type="exact"
      expand-all
    />
  </div>
</template>

<script setup>
import { useChatbotStore } from '@planetcrust/human-vue'
import { components } from '@planetcrust/human-vue'
import { computed, inject, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'

const { CSidebarNav } = components
const { t } = useI18n()

const $SystemAPI = inject('$SystemAPI')
const chatbotStore = useChatbotStore()

// The shell does global setup only; the section loads the list its nav renders.
onMounted(() => {
  if (!chatbotStore.list.length) {
    chatbotStore.fetchList()
  }
})

const navItems = computed(() => [
  {
    _id: 'chatbots',
    _parentId: '0',
    _label: t('navigation.chatbots'),
    _icon: 'pi pi-comments',
    _route: { name: 'chatbot' },
  },
  {
    _id: 'sessions',
    _parentId: '0',
    _label: t('navigation.sessions'),
    _icon: 'pi pi-list',
    _route: { name: 'chatbot.sessions' },
  },
  ...[...chatbotStore.list]
    .sort((a, b) => {
      const labelA = a.name || a.handle || ''
      const labelB = b.name || b.handle || ''
      return labelA.localeCompare(labelB)
    })
    .map(c => ({
      _id: c.chatbotID,
      _parentId: 'chatbots',
      _label: c.name || c.handle || t('chatbot.list.unnamed') || c.chatbotID,
      _route: { name: 'chatbot.edit', params: { chatbotID: c.chatbotID } },
    })),
])
</script>
