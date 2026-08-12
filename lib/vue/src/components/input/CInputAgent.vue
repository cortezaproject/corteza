<template>
  <Select
    :model-value="modelValue"
    @update:model-value="onSelect"
    :options="options"
    option-label="label"
    option-value="agentID"
    :placeholder="placeholder"
    :disabled="disabled"
    :loading="loading"
    class="w-full"
    filter
    :filter-fields="['label', 'handle', 'agentID']"
    fluid
    showClear
    @show="onShow"
  />
</template>

<script setup>
import { inject, onBeforeUnmount, onMounted, ref, watch } from 'vue'

defineOptions({ inheritAttrs: false })

const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: null,
  },
  placeholder: {
    type: String,
    default: '',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue'])

const $SystemAPI = inject('$SystemAPI')

const options = ref([])
const loading = ref(false)

let cancelCurrentRequest = null

function getOptionLabel(agent) {
  if (!agent) return ''
  return agent.meta?.short || agent.handle || agent.agentID
}

async function fetchAgents() {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
    cancelCurrentRequest = null
  }

  loading.value = true
  try {
    const { response, cancel } = $SystemAPI.agentListCancellable({
      limit: 100,
      sort: 'handle ASC',
    })
    cancelCurrentRequest = cancel

    const result = await response()
    const agents = Array.isArray(result) ? result : result.set || []

    options.value = agents
      .map(a => ({ ...a, label: getOptionLabel(a) }))
      .sort((a, b) => (a.label || '').localeCompare(b.label || ''))

    if (props.modelValue && props.modelValue !== '0') {
      if (!options.value.some(a => a.agentID === props.modelValue)) {
        loadAgentById(props.modelValue)
      }
    }
  } finally {
    loading.value = false
    cancelCurrentRequest = null
  }
}

function onShow() {
  if (options.value.length === 0) {
    fetchAgents()
  }
}

function onSelect(agentID) {
  emit('update:modelValue', agentID || null)
}

async function loadAgentById(agentID) {
  if (!agentID || agentID === '0' || !$SystemAPI) return
  loading.value = true
  try {
    const agent = await $SystemAPI.agentRead({ agentID })
    if (agent) {
      agent.label = getOptionLabel(agent)
      if (!options.value.find(a => a.agentID === agentID)) {
        options.value = [...options.value, agent]
      }
    }
  } catch {
    // Agent not found or API error
  } finally {
    loading.value = false
  }
}

watch(
  () => props.modelValue,
  newVal => {
    if (newVal && newVal !== '0') {
      loadAgentById(newVal)
    }
  },
  { immediate: true },
)

onMounted(() => {
  fetchAgents()
  if (props.modelValue && props.modelValue !== '0') {
    loadAgentById(props.modelValue)
  }
})

onBeforeUnmount(() => {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
})
</script>
