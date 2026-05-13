<template>
  <Dialog
    v-model:visible="isVisible"
    :header="$t('chatbot.sessions.detail.title')"
    modal
    :draggable="false"
    :dismissable-mask="true"
    :style="{ width: '90vw', height: '90vh' }"
    content-class="h-full p-0 !pb-0 flex flex-col"
    @hide="onHide"
  >
    <div v-if="loading" class="flex items-center justify-center flex-1 h-full">
      <ProgressSpinner />
    </div>

    <div v-else-if="messages.length === 0" class="flex items-center justify-center flex-1 h-full text-muted-color">
      {{ $t('chatbot.sessions.detail.empty') }}
    </div>

    <CChatMessages
      v-else
      :messages="messages"
      :readonly="true"
      class="flex-1 min-h-0"
    />
  </Dialog>
</template>

<script setup>
import { components } from '@planetcrust/human-vue'
import { inject, ref, watch } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false,
  },
  session: {
    type: Object,
    default: null,
  },
})

const emit = defineEmits(['update:visible'])

const { CChatMessages } = components
const $SystemAPI = inject('$SystemAPI')

const loading = ref(false)
const messages = ref([])

const isVisible = ref(props.visible)

watch(() => props.visible, v => { isVisible.value = v })
watch(isVisible, v => emit('update:visible', v))

watch(() => props.session, async session => {
  if (!session) return
  messages.value = []
  loading.value = true
  try {
    const res = await $SystemAPI.chatbotSessionRead({ sessionID: session.id })
    const steps = (res.steps || []).sort((a, b) => (a.scenarioIndex ?? 0) - (b.scenarioIndex ?? 0))

    const allMessages = []
    for (const step of steps) {
      if (!step.conversationID || step.conversationID === '0') continue
      try {
        const conv = await $SystemAPI.aiConversationRead({ aiConversationID: step.conversationID })
        const msgs = conv.messages || []
        allMessages.push(...msgs)
      } catch {
        // step may have no conversation yet
      }
    }
    messages.value = allMessages
  } catch (e) {
    console.error('Failed to load session detail:', e)
  } finally {
    loading.value = false
  }
}, { immediate: false })

function onHide() {
  messages.value = []
}
</script>
