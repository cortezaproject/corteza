<template>
  <div class="flex flex-col min-h-0 flex-1">
    <!-- Message list -->
    <div class="flex-1 p-4 overflow-y-auto flex flex-col gap-4" ref="chatContainer">
      <template v-if="visibleMessages.length > 0">
        <div
          v-for="(msg, index) in visibleMessages"
          :key="index"
          class="flex flex-col gap-1"
          :class="msg.role === 'user' ? 'items-end' : 'items-start'"
        >
          <div
            :ref="el => assignMsgRef(el, msg)"
            :class="[
              'p-3 xl:p-4 rounded-xl max-w-[85%] text-sm md:text-base transition-all duration-200',
              msg.role === 'user'
                ? 'bg-primary text-primary-contrast shadow-sm whitespace-pre-wrap'
                : 'bg-emphasis text-color shadow-sm',
              msg.traceIndex !== undefined ? 'cursor-pointer' : '',
              msg.traceIndex !== undefined && ['agent', 'assistant'].includes(msg.role)
                ? 'hover:shadow-md'
                : '',
              selectedTraceIndex === msg.traceIndex &&
              msg.traceIndex !== undefined &&
              selectedTraceType === (msg.role === 'user' ? 'prompt' : 'response')
                ? 'ring-2 ring-primary ring-offset-1'
                : '',
            ]"
            @click="onBubbleClick(msg)"
          >
            <div
              v-if="['agent', 'assistant'].includes(msg.role)"
              class="rt-content"
              v-html="renderMarkdown(msg.content)"
            />
            <template v-else>{{ msg.content }}</template>
          </div>
        </div>
      </template>

      <!-- Empty state -->
      <div
        v-else
        class="flex h-full items-center justify-center text-muted-color p-4 text-center text-sm flex-1"
      >
        <slot name="empty" />
      </div>

      <!-- Thinking bubble -->
      <div v-if="executing" class="flex items-start">
        <div
          class="bg-emphasis text-color shadow-sm p-3 rounded-xl flex items-center gap-2 text-sm"
        >
          <ProgressSpinner style="width: 16px; height: 16px" strokeWidth="4" />
          <span class="text-muted-color">{{ thinkingLabel }}</span>
        </div>
      </div>
    </div>

    <!-- Input row -->
    <CChatComposer
      v-if="!readonly"
      :placeholder="placeholder"
      :disabled="disabled || executing"
      :loading="executing"
      @send="submit"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import { renderMarkdown } from '@planetcrust/human-js'
import CChatComposer from './CChatComposer.vue'

interface ChatMessage {
  role: 'user' | 'agent' | 'assistant' | string
  content: string
  traceIndex?: number
}

const props = defineProps({
  messages: {
    type: Array as () => ChatMessage[],
    default: () => [],
  },
  executing: {
    type: Boolean,
    default: false,
  },
  readonly: {
    type: Boolean,
    default: false,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  placeholder: {
    type: String,
    default: '',
  },
  thinkingLabel: {
    type: String,
    default: 'Thinking...',
  },
  selectedTraceIndex: {
    type: Number,
    default: null,
  },
  selectedTraceType: {
    type: String as () => 'prompt' | 'response' | null,
    default: null,
  },
})

const emit = defineEmits<{
  (_e: 'send', _input: string): void
  (_e: 'trace-select', _traceIndex: number, _type: 'prompt' | 'response'): void
}>()

const chatContainer = ref<HTMLElement | null>(null)
const chatMsgRefs = ref<Record<number, HTMLElement>>({})
const chatPromptRefs = ref<Record<number, HTMLElement>>({})

const visibleMessages = computed(() =>
  props.messages.filter(m => ['user', 'agent', 'assistant'].includes(m.role) && m.content),
)

function assignMsgRef(el: any, msg: ChatMessage) {
  if (msg.traceIndex === undefined) return
  if (msg.role === 'user') {
    chatPromptRefs.value[msg.traceIndex] = el
  } else {
    chatMsgRefs.value[msg.traceIndex] = el
  }
}

function onBubbleClick(msg: ChatMessage) {
  if (msg.traceIndex === undefined) return
  const type = msg.role === 'user' ? 'prompt' : 'response'
  emit('trace-select', msg.traceIndex, type)
}

function submit(input: string) {
  if (!input.trim() || props.executing || props.disabled) return
  emit('send', input)
}

async function scrollToBottom() {
  await nextTick()
  if (chatContainer.value) {
    chatContainer.value.scrollTop = chatContainer.value.scrollHeight
  }
}

// Scroll when new messages arrive
watch(
  () => props.messages.length,
  () => scrollToBottom(),
)

// Scroll to the relevant bubble when trace selection changes
watch(
  () => [props.selectedTraceIndex, props.selectedTraceType],
  ([idx, type]) => {
    if (idx === null || idx === undefined) return
    nextTick(() => {
      const refs = type === 'prompt' ? chatPromptRefs : chatMsgRefs
      const el = refs.value[idx as number]
      if (el) el.scrollIntoView({ behavior: 'smooth', block: 'nearest' })
    })
  },
)

defineExpose({ scrollToBottom })
</script>

<style scoped>
.rt-content :deep(p) {
  margin-top: 0.25rem;
  margin-bottom: 0.25rem;
}
.rt-content :deep(p:first-child) {
  margin-top: 0;
}
.rt-content :deep(p:last-child) {
  margin-bottom: 0;
}
</style>
