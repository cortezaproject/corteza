<template>
  <PageBlock :block="block" :record="record" @refreshBlock="onManualRefresh">
    <CAgentChat
      ref="chatRef"
      :translations="translations"
      :allowed-agent-i-ds="options.allowedAgentIDs"
      :default-agent-i-d="options.defaultAgentID"
      :auto-resume="options.autoResume"
      :context-provider="contextProvider"
    />
  </PageBlock>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { CAgentChat, makeAgentChatTranslations } from '@planetcrust/human-vue'
import PageBlock from './PageBlock.vue'

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const { t } = useI18n()
const translations = computed(() => makeAgentChatTranslations(t, 'block.agentChat.'))

const options = computed(() => props.block.options)

// Pass structured caller-context (namespace, module, record, page) to every
// agentExec call from this block. Returned object lands on the request as
// `context: map[string]any` and is appended to the agent's system prompt.
function contextProvider() {
  const ctx = {}
  if (props.namespace?.namespaceID) ctx.namespaceID = String(props.namespace.namespaceID)
  if (props.record?.moduleID) ctx.moduleID = String(props.record.moduleID)
  if (props.record?.recordID) ctx.recordID = String(props.record.recordID)
  if (props.page?.pageID) ctx.pageID = String(props.page.pageID)
  if (props.record?.values) {
    // Record.values may be a {name: value} object (compose.Record) or an
    // array [{name, value}] (raw API). Normalize to a flat object.
    if (Array.isArray(props.record.values)) {
      ctx.recordValues = props.record.values.reduce((acc, v) => {
        if (v && v.name) acc[v.name] = v.value
        return acc
      }, {})
    } else if (typeof props.record.values === 'object') {
      ctx.recordValues = { ...props.record.values }
    }
  }
  return ctx
}

const chatRef = ref(null)

async function onManualRefresh() {
  await Promise.allSettled([
    chatRef.value?.reloadAvailableAgents?.(),
    chatRef.value?.reloadHistory?.(),
    chatRef.value?.refreshActiveConversation?.(),
  ])
}
</script>
