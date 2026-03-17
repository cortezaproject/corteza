<template>
  <Select
    :model-value="selectedUser"
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

const options = ref([])
const selectedUser = ref(null)
const loading = ref(false)

let cancelCurrentRequest = null

function getOptionLabel(user) {
  return formatUser(user)
}

async function fetchUsers() {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
    cancelCurrentRequest = null
  }

  loading.value = true
  try {
    const { response, cancel } = $SystemAPI.userListCancellable({
      query: '',
      limit: 100,
    })
    cancelCurrentRequest = cancel

    const result = await response()
    options.value = result.set || []
    cacheUsers(options.value)
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
    fetchUsers()
  }
}

function onSelect(value) {
  selectedUser.value = value
  emit('update:modelValue', value?.userID || null)

  if (value) {
    emit('select', value)
  }

  if (props.clearOnSelect && value) {
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
    if (user) {
      selectedUser.value = user
      if (!options.value.find(u => u.userID === userID)) {
        options.value = [...options.value, user]
      }
    }
  } catch {
    // User not found or API error
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
  fetchUsers()
  if (props.modelValue) {
    loadUserById(props.modelValue)
  }
})

onBeforeUnmount(() => {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
})
</script>
