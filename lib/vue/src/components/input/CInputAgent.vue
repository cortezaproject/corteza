<template>
  <Select
    :model-value="selectedAgent"
    @update:model-value="onSelect"
    :options="options"
    :option-label="getOptionLabel"
    :placeholder="placeholder"
    :disabled="disabled"
    :loading="loading"
    class="w-full"
    filter
    fluid
    showClear
    @show="onShow"
  >
    <template #option="{ option }">
      {{ getOptionLabel(option) }}
    </template>
  </Select>
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
const selectedAgent = ref(null)
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
    })
    cancelCurrentRequest = cancel

    const result = await response()
    options.value = Array.isArray(result) ? result : result.set || []
  } catch (e) {
    if (e?.message !== 'canceled') {
      options.value = []
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

function onSelect(value) {
  selectedAgent.value = value
  emit('update:modelValue', value?.agentID || null)
}

async function loadAgentById(agentID) {
  if (!agentID || agentID === '0' || !$SystemAPI) return
  loading.value = true
  try {
    const agent = await $SystemAPI.agentRead({ agentID })
    selectedAgent.value = agent
    if (!options.value.find(a => a.agentID === agentID)) {
      options.value = [...options.value, agent]
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
    if (
      newVal &&
      newVal !== '0' &&
      (!selectedAgent.value || selectedAgent.value.agentID !== newVal)
    ) {
      loadAgentById(newVal)
    } else if (!newVal || newVal === '0') {
      selectedAgent.value = null
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
