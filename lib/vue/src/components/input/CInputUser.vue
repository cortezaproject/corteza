<template>
  <AutoComplete
    :model-value="selectedUser"
    @update:model-value="onUpdateModelValue"
    @item-select="onItemSelect"
    :suggestions="suggestions"
    :option-label="getOptionLabel"
    :placeholder="placeholder"
    :disabled="disabled"
    :loading="loading"
    @complete="search"
    class="w-full"
    dropdown
    showClear
    :complete-on-focus="true"
  >
    <template #option="{ option }">
      {{ getOptionLabel(option) }}
    </template>
  </AutoComplete>
</template>

<script setup>
import { debounce } from 'lodash-es'
import { inject, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useUserResolver } from '../../composables/useUserResolver'

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
  clearOnSelect: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue', 'select'])

const $SystemAPI = inject('$SystemAPI')
const { formatUser, resolveUser, cacheUsers } = useUserResolver()

const suggestions = ref([])
const selectedUser = ref(null)
const loading = ref(false)

// Store cancel function for current request
let cancelCurrentRequest = null

function getOptionLabel(user) {
  return formatUser(user)
}

const search = debounce(async event => {
  // Cancel previous request if pending
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
    cancelCurrentRequest = null
  }

  loading.value = true
  try {
    const { response, cancel } = $SystemAPI.userListCancellable({
      query: event.query,
      limit: 20,
    })
    cancelCurrentRequest = cancel

    const result = await response()
    suggestions.value = result.set || []
    cacheUsers(suggestions.value)
  } catch (e) {
    // Ignore cancelled requests
    if (e?.message !== 'canceled') {
      suggestions.value = []
    }
  } finally {
    loading.value = false
    cancelCurrentRequest = null
  }
}, 300)

function onUpdateModelValue(value) {
  // AutoComplete emits string when typing, and object when selected
  selectedUser.value = value
  emit('update:modelValue', value?.userID || null)
}

function onItemSelect(event) {
  emit('select', event.value)
  if (props.clearOnSelect) {
    // defer clearing to allow event to propagate and component to finish updates
    setTimeout(() => {
      selectedUser.value = null
      emit('update:modelValue', null)
    }, 0)
  }
}

async function loadUserById(userID) {
  if (!userID) return

  loading.value = true
  try {
    const user = await resolveUser(userID)
    if (user) selectedUser.value = user
  } catch {
    // User not found or API error - leave selectedUser as null
  } finally {
    loading.value = false
  }
}

watch(
  () => props.modelValue,
  newVal => {
    if (newVal && (!selectedUser.value || selectedUser.value.userID !== newVal)) {
      loadUserById(newVal)
    } else if (!newVal) {
      selectedUser.value = null
    }
  },
  { immediate: true },
)

onMounted(() => {
  if (props.modelValue) {
    loadUserById(props.modelValue)
  }
})

onBeforeUnmount(() => {
  // Cancel any pending request
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
  search.cancel()
})
</script>
