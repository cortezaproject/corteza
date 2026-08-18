<template>
  <div class="flex flex-col gap-3">
    <CFormGroup
      :label="$t('block.agentChat.config.allowedAgents')"
      :description="$t('block.agentChat.config.allowedAgentsHint')"
      required
    >
      <MultiSelect
        v-model="allowedAgentIDs"
        :options="agents"
        :option-label="agentLabel"
        option-value="agentID"
        :placeholder="$t('block.agentChat.config.allowedAgentsPlaceholder')"
        :loading="loading"
        filter
        display="chip"
        class="w-full"
      />
    </CFormGroup>

    <CFormGroup
      :label="$t('block.agentChat.config.defaultAgent')"
      :description="$t('block.agentChat.config.defaultAgentHint')"
    >
      <Select
        v-model="defaultAgentID"
        :options="allowedAgents"
        :option-label="agentLabel"
        option-value="agentID"
        :placeholder="$t('block.agentChat.config.defaultAgentPlaceholder')"
        :disabled="allowedAgents.length === 0"
        show-clear
        class="w-full"
      />
    </CFormGroup>

    <CInputToggleCard
      v-model="autoResume"
      class="max-w-xl"
      :label="$t('block.agentChat.config.autoResume')"
      :description="$t('block.agentChat.config.autoResumeHint')"
    />
  </div>
</template>

<script setup>
import { computed, inject, onMounted, ref } from 'vue'

const block = inject('blockDraft')
const $SystemAPI = inject('$SystemAPI', null)

const loading = ref(false)
const agents = ref([])

const opts = computed(() => block.value.options)

const allowedAgentIDs = computed({
  get: () => opts.value.allowedAgentIDs,
  set: v => {
    opts.value.allowedAgentIDs = (v || []).map(String)
  },
})

const defaultAgentID = computed({
  get: () => opts.value.defaultAgentID,
  set: v => {
    opts.value.defaultAgentID = v ? String(v) : ''
  },
})

const autoResume = computed({
  get: () => opts.value.autoResume,
  set: v => {
    opts.value.autoResume = !!v
  },
})

// Default agent picker is constrained to whatever the allowlist contains so
// editors cannot pick a default that viewers won't see.
const allowedAgents = computed(() => {
  const ids = new Set((allowedAgentIDs.value || []).map(String))
  if (ids.size === 0) return []
  return agents.value.filter(a => ids.has(String(a.agentID)))
})

function agentLabel(a) {
  return a?.meta?.short || a?.handle || a?.agentID || ''
}

onMounted(async () => {
  if (!$SystemAPI) return
  loading.value = true
  try {
    const res = await $SystemAPI.agentList({ limit: 0 })
    agents.value = (res?.set || []).filter(a => a.invocation?.user?.enabled)
  } catch (err) {
    console.warn('[agent-chat-configurator] agentList failed', err)
  } finally {
    loading.value = false
  }
})
</script>
